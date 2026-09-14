import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { BooleanPicker } from '../core/components/boolean_picker';
import { ContentBlock } from '../core/components/content_block';
import { EnumPicker } from '../core/components/enum_picker';
import Toast from '../core/components/toast';
import { getLaunchedSims, raidSimStatus } from '../core/launched_sims';
import { ProgressMetrics, StatWeightsResult, StatWeightValues } from '../core/proto/api';
import { Spec, Stat } from '../core/proto/common';
import { statNames } from '../core/proto_utils/names';
import { UnitStat } from '../core/proto_utils/stats';
import { getSpecSiteUrl, naturalSpecOrder, specNames, textCssClassForSpec } from '../core/proto_utils/utils';
import { Sim } from '../core/sim';
import { RequestTypes } from '../core/sim_signal_manager';
import { SimUI } from '../core/sim_ui';
import { EventID, TypedEvent } from '../core/typed_event';
import { stDevToConf90 } from '../core/utils';
import { getSpecEpConfig, runSpecStatWeights } from './runner';

// Which of the metrics a stat weights run produces is on screen. One run fills in all of
// them, so switching between these never costs another sim.
enum WeightsMetric {
	Dps,
	Tps,
	Dtps,
}

const METRIC_LABELS: Record<WeightsMetric, string> = {
	[WeightsMetric.Dps]: 'Damage (DPS)',
	[WeightsMetric.Tps]: 'Threat (TPS)',
	[WeightsMetric.Dtps]: 'Damage taken (DTPS)',
};

// The stats worth putting side by side for throughput. Fixed rather than derived from the
// specs on screen, so a column means the same thing in every row.
const THROUGHPUT_STATS: Array<Stat> = [
	Stat.StatStrength,
	Stat.StatAgility,
	Stat.StatIntellect,
	Stat.StatSpirit,
	Stat.StatAttackPower,
	Stat.StatRangedAttackPower,
	Stat.StatMeleeHit,
	Stat.StatMeleeCrit,
	Stat.StatExpertise,
	Stat.StatSpellPower,
	Stat.StatSpellHit,
	Stat.StatSpellCrit,
	Stat.StatSpellHaste,
	Stat.StatMP5,
];

// Damage taken is about staying alive, so the mitigation stats replace the throughput ones
// when that metric is selected.
const MITIGATION_STATS: Array<Stat> = [
	Stat.StatStamina,
	Stat.StatHealth,
	Stat.StatArmor,
	Stat.StatBonusArmor,
	Stat.StatDefense,
	Stat.StatDodge,
	Stat.StatParry,
	Stat.StatBlock,
	Stat.StatBlockValue,
];

// Every stat any launched spec asks to have weighed, in the order the columns appear when
// 'Show All Stats' is on.
const ALL_STATS: Array<Stat> = [
	Stat.StatStrength,
	Stat.StatAgility,
	Stat.StatStamina,
	Stat.StatIntellect,
	Stat.StatSpirit,
	Stat.StatHealth,
	Stat.StatMana,
	Stat.StatAttackPower,
	Stat.StatFeralAttackPower,
	Stat.StatRangedAttackPower,
	Stat.StatMeleeHit,
	Stat.StatMeleeCrit,
	Stat.StatExpertise,
	Stat.StatSpellPower,
	Stat.StatSpellDamage,
	Stat.StatArcanePower,
	Stat.StatFirePower,
	Stat.StatFrostPower,
	Stat.StatHolyPower,
	Stat.StatNaturePower,
	Stat.StatShadowPower,
	Stat.StatHealingPower,
	Stat.StatSpellHit,
	Stat.StatSpellCrit,
	Stat.StatSpellHaste,
	Stat.StatMP5,
	Stat.StatArmor,
	Stat.StatBonusArmor,
	Stat.StatDefense,
	Stat.StatDodge,
	Stat.StatParry,
	Stat.StatBlock,
	Stat.StatBlockValue,
	Stat.StatArcaneResistance,
	Stat.StatFireResistance,
	Stat.StatFrostResistance,
	Stat.StatNatureResistance,
	Stat.StatShadowResistance,
];

// EP for damage taken is normalised against armor by the sim itself rather than against
// the spec's own reference stat. See DTPSReferenceStat in sim/core/statweight.go.
const DTPS_REFERENCE_STAT = Stat.StatArmor;

const DEFAULT_ITERATIONS = 1000;

const STORAGE_KEY_PREFIX = '__classic_stat_weights';

const statName = (stat: Stat): string => statNames.get(stat) ?? 'Unknown';

interface SpecRow {
	readonly spec: Spec;
	readonly epStats: Array<Stat>;
	readonly epReferenceStat: Stat;
	readonly rowElem: HTMLTableRowElement;
	readonly statusElem: HTMLSpanElement;
	readonly runButton: HTMLButtonElement;
	result: StatWeightsResult | null;
	iterations: number;
	durationMs: number;
}

export class StatWeightsSimUI extends SimUI {
	private readonly rows: Array<SpecRow> = [];
	private readonly tableHeadRow: HTMLTableRowElement;
	private readonly tableBody: HTMLTableSectionElement;
	private readonly normalisationElem: HTMLParagraphElement;
	private readonly costElem: HTMLParagraphElement;
	private readonly runAllButton: HTMLButtonElement;

	private metric = WeightsMetric.Dps;
	private showAllStats = false;
	private running = false;
	private cancelling = false;

	constructor(parentElem: HTMLElement) {
		super(parentElem, new Sim(), {
			cssClass: 'stat-weights-sim-ui',
			cssScheme: 'raid',
			spec: null,
			simStatus: raidSimStatus,
			knownIssues: [
				'Every spec here is simmed on its own sim defaults, so these weights describe that default build and nothing else. Re-run them on your own gear in the individual sim before acting on them.',
			],
		});

		this.sim.setIterations(TypedEvent.nextEventID(), DEFAULT_ITERATIONS);

		// The sidebar title comes from SimTitleDropdown, which only knows about a spec or
		// the raid sim, so it arrives naming the raid sim. The dropdown under it is the
		// app's navigation and is worth keeping; only the label is wrong.
		(this.rootElem.querySelector('.sim-link-title') as HTMLElement).textContent = 'Stat Weights';

		this.addTab('Stat Weights', 'stat-weights', '');
		const tabContent = this.rootElem.querySelector('#stat-weights-tab') as HTMLElement;

		const aboutBlock = new ContentBlock(tabContent, 'stat-weights-about', {
			header: { title: 'Which stats matter for which spec' },
		});
		const normalisationRef = ref<HTMLParagraphElement>();
		const costRef = ref<HTMLParagraphElement>();
		aboutBlock.bodyElement.appendChild(
			<>
				<p>
					Each row is one launched spec, simmed on the gear, talents, consumes, buffs, debuffs and encounter its own individual sim starts you on,
					with the auto rotation. Nothing is hand-tuned, and no spec is geared for the stats it is being asked about.
				</p>
				<p ref={normalisationRef} className="stat-weights-normalisation"></p>
				<p>
					A weight is only as good as the sample behind it. At low iteration counts a small weight is mostly noise, so hover a value for its 90%
					confidence interval before reading anything into a difference of a few hundredths. Values smaller than their own confidence interval are
					greyed out.
				</p>
				<p>
					Melee and spell hit and crit off gear are one pool under the Forever ruleset, not two, so a spec that weighs both gets the same number in
					both columns. That is the rule working, not the table repeating itself.
				</p>
				<p ref={costRef} className="stat-weights-cost"></p>
				<p className="stat-weights-missing">
					Four healer specs (Restoration Druid, Holy Paladin, Healing Priest, Restoration Shaman) and Feral Tank Druid are not implemented in this
					sim, so they are absent from the table, and healing weights (HPS) are not offered for the same reason. The resistances and the weapon
					pseudo-stats an individual sim will also weigh are left out here, since neither compares across specs.
				</p>
			</>,
		);
		this.normalisationElem = normalisationRef.value!;
		this.costElem = costRef.value!;

		const tableBlock = new ContentBlock(tabContent, 'stat-weights-results', {
			header: { title: 'Stat weights by spec' },
		});
		const optionsRef = ref<HTMLDivElement>();
		const headRowRef = ref<HTMLTableRowElement>();
		const bodyRef = ref<HTMLTableSectionElement>();
		tableBlock.bodyElement.appendChild(
			<>
				<div ref={optionsRef} className="stat-weights-options row"></div>
				<div className="stat-weights-table-container">
					<table className="stat-weights-table">
						<thead>
							<tr ref={headRowRef}></tr>
						</thead>
						<tbody ref={bodyRef}></tbody>
					</table>
				</div>
			</>,
		);
		this.tableHeadRow = headRowRef.value!;
		this.tableBody = bodyRef.value!;

		const optionsElem = optionsRef.value!;
		const metricContainer = (<div className="col col-sm-4"></div>) as HTMLElement;
		const showAllContainer = (<div className="col col-sm-4 stat-weights-show-all"></div>) as HTMLElement;
		optionsElem.appendChild(metricContainer);
		optionsElem.appendChild(showAllContainer);

		new EnumPicker(metricContainer, this, {
			id: 'stat-weights-metric',
			label: 'Metric',
			labelTooltip: 'Which metric the weights are measured against. One run produces all of them, so switching does not re-run anything.',
			values: [WeightsMetric.Dps, WeightsMetric.Tps, WeightsMetric.Dtps].map(metric => ({
				name: METRIC_LABELS[metric],
				value: metric,
			})),
			changedEvent: () => new TypedEvent(),
			getValue: () => this.metric,
			setValue: (_eventID: EventID, _ui: StatWeightsSimUI, newValue: number) => {
				this.metric = newValue as WeightsMetric;
				this.updateTable();
			},
		});

		new BooleanPicker(showAllContainer, this, {
			id: 'stat-weights-show-all',
			label: 'Show All Stats',
			labelTooltip: 'Show every stat any spec asks to have weighed, including school-specific damage and resistances.',
			inline: true,
			changedEvent: () => new TypedEvent(),
			getValue: () => this.showAllStats,
			setValue: (_eventID: EventID, _ui: StatWeightsSimUI, newValue: boolean) => {
				this.showAllStats = newValue;
				this.updateTable();
			},
		});

		getLaunchedSims()
			.sort((a, b) => naturalSpecOrder.indexOf(a) - naturalSpecOrder.indexOf(b))
			.forEach(spec => this.rows.push(this.makeSpecRow(spec)));

		const runAllButtonRef = ref<HTMLButtonElement>();
		this.simActionsContainer.appendChild(
			<button ref={runAllButtonRef} className="btn btn-primary w-100 stat-weights-run-all">
				Run All Specs
			</button>,
		);
		this.runAllButton = runAllButtonRef.value!;
		tippy(this.runAllButton, {
			content: 'Runs every spec that has no result yet, one after another. Specs already filled in are left alone.',
		});
		this.runAllButton.addEventListener('click', () => {
			if (this.running) this.cancel();
			else this.runSpecs(this.rows.filter(row => !row.result));
		});

		this.updateTable();
		this.sim.iterationsChangeEmitter.on(() => this.updateNotes());
	}

	getStorageKey(postfix: string): string {
		return STORAGE_KEY_PREFIX + postfix;
	}

	applyDefaults(eventID: EventID) {
		this.sim.setIterations(eventID, DEFAULT_ITERATIONS);
	}

	toLink(): string {
		return window.location.href;
	}

	private makeSpecRow(spec: Spec): SpecRow {
		const epConfig = getSpecEpConfig(spec);
		const rowRef = ref<HTMLTableRowElement>();
		const statusRef = ref<HTMLSpanElement>();
		const buttonRef = ref<HTMLButtonElement>();

		this.tableBody.appendChild(
			<tr ref={rowRef} className="stat-weights-row">
				<th className="stat-weights-spec-cell">
					<button ref={buttonRef} className="btn btn-sm btn-primary stat-weights-run">
						Run
					</button>
					<a className={`stat-weights-spec-name ${textCssClassForSpec(spec)}`} href={getSpecSiteUrl(spec)} target="_blank">
						{specNames[spec]}
					</a>
					<span ref={statusRef} className="stat-weights-status"></span>
				</th>
				<td className="stat-weights-reference-cell"></td>
			</tr>,
		);

		const row: SpecRow = {
			spec: spec,
			epStats: epConfig.epStats,
			epReferenceStat: epConfig.epReferenceStat,
			rowElem: rowRef.value!,
			statusElem: statusRef.value!,
			runButton: buttonRef.value!,
			result: null,
			iterations: 0,
			durationMs: 0,
		};

		row.runButton.addEventListener('click', () => this.runSpecs([row]));

		return row;
	}

	// The stats that get a column, which depends on the metric and on whether all of them
	// have been asked for.
	private displayedStats(): Array<Stat> {
		if (this.showAllStats) {
			return ALL_STATS.filter(stat => this.rows.some(row => row.epStats.includes(stat)));
		}
		return this.metric == WeightsMetric.Dtps ? MITIGATION_STATS : THROUGHPUT_STATS;
	}

	private statWeightValues(row: SpecRow): StatWeightValues | undefined {
		if (!row.result) return undefined;
		switch (this.metric) {
			case WeightsMetric.Tps:
				return row.result.tps;
			case WeightsMetric.Dtps:
				return row.result.dtps;
			default:
				return row.result.dps;
		}
	}

	private referenceStat(row: SpecRow): Stat {
		return this.metric == WeightsMetric.Dtps ? DTPS_REFERENCE_STAT : row.epReferenceStat;
	}

	private updateTable() {
		const stats = this.displayedStats();

		this.tableHeadRow.innerHTML = '';
		this.tableHeadRow.appendChild(
			<>
				<th className="stat-weights-spec-cell">Spec</th>
				<th className="stat-weights-reference-cell">Relative to</th>
				{stats.map(stat => (
					<th className="stat-weights-value-cell">{statName(stat)}</th>
				))}
			</>,
		);

		this.rows.forEach(row => this.updateRow(row, stats));
		this.updateNotes();
	}

	private updateRow(row: SpecRow, stats: Array<Stat>) {
		const referenceCell = row.rowElem.querySelector('.stat-weights-reference-cell') as HTMLElement;
		referenceCell.textContent = statName(this.referenceStat(row));

		row.rowElem.querySelectorAll('.stat-weights-value-cell').forEach(cell => cell.remove());

		const values = this.statWeightValues(row);
		const iterations = row.iterations || 1;
		const best = stats.reduce(
			(acc, stat) => (values && row.epStats.includes(stat) ? Math.max(acc, UnitStat.fromStat(stat).getProtoValue(values.epValues!)) : acc),
			0,
		);

		stats.forEach(stat => {
			const cellRef = ref<HTMLTableCellElement>();
			row.rowElem.appendChild(<td ref={cellRef} className="stat-weights-value-cell"></td>);
			const cell = cellRef.value!;

			if (!row.epStats.includes(stat)) {
				cell.classList.add('stat-weights-not-weighed');
				cell.textContent = '—';
				tippy(cell, { content: `The ${specNames[row.spec]} sim does not ask for this stat to be weighed.` });
				return;
			}
			if (!values) return;

			const unitStat = UnitStat.fromStat(stat);
			const ep = unitStat.getProtoValue(values.epValues!);
			const conf90 = stDevToConf90(unitStat.getProtoValue(values.epValuesStdev!), iterations);
			// A weight a hair under zero rounds to '-0.00', which reads as a real negative.
			cell.textContent = ep.toFixed(2).replace('-0.00', '0.00');
			if (best > 0 && ep == best) cell.classList.add('stat-weights-best');
			else if (Math.abs(ep) < conf90) cell.classList.add('stat-weights-noise');
			tippy(cell, {
				content: `${statName(stat)}: ${ep.toFixed(3)} &plusmn;${conf90.toFixed(3)} at 90% confidence, over ${row.iterations} iterations.`,
			});
		});
	}

	private updateNotes() {
		const iterations = this.sim.getIterations();
		const referenceNames = [...new Set(this.rows.map(row => statName(this.referenceStat(row))))];
		this.normalisationElem.textContent =
			this.metric == WeightsMetric.Dtps
				? `Values are EP normalised against ${statName(DTPS_REFERENCE_STAT)}, which the sim fixes for damage taken rather than taking it ` +
				  `from the spec: 1.00 avoids as much damage as one point of ${statName(DTPS_REFERENCE_STAT)}, and higher is better even though ` +
				  `damage taken is what is going down.`
				: `Values are EP normalised against each spec's own reference stat, named in the second column (${referenceNames.join(', ')}): ` +
				  `1.00 is one point of that stat, 2.00 is worth twice as much per point. Reading down a column compares how much each spec wants ` +
				  `a stat relative to its own reference, not how much raw output it gets, because every row is divided by a different number.`;

		const completed = this.rows.filter(row => row.result);
		const measured = completed.filter(row => row.durationMs > 0);
		const perSpec = measured.length ? measured.reduce((acc, row) => acc + row.durationMs, 0) / measured.length / 1000 : 0;
		const preamble = `Nothing runs until asked. A stat weights run is two sims per stat plus a baseline, so filling all ${this.rows.length} rows is several hundred sims. `;
		this.costElem.textContent = measured.length
			? `${preamble}At ${iterations} iterations this browser is averaging ${perSpec.toFixed(1)}s per spec, putting all ${
					this.rows.length
			  } at roughly ${Math.max(1, Math.round((perSpec * this.rows.length) / 60))} minutes. ${completed.length} of ${this.rows.length} filled in so far.`
			: `${preamble}Run a single spec first to see what this browser does at ${iterations} iterations before committing to all ${this.rows.length}.`;
	}

	private setRunning(running: boolean) {
		this.running = running;
		this.cancelling = false;
		this.runAllButton.textContent = running ? 'Cancel' : 'Run All Specs';
		this.rows.forEach(row => (row.runButton.disabled = running));
	}

	private async cancel() {
		if (this.cancelling) return;
		this.cancelling = true;
		this.runAllButton.textContent = 'Cancelling';
		try {
			await this.sim.signalManager.abortType(RequestTypes.StatWeights);
		} catch (error) {
			console.error(error);
		}
	}

	private async runSpecs(rows: Array<SpecRow>) {
		if (this.running) return;
		if (!rows.length) {
			new Toast({ variant: 'info', body: 'Every spec already has weights. Use a spec’s own Run button to redo one.' });
			return;
		}
		this.setRunning(true);

		const startTime = performance.now();
		let index = 0;
		for (const row of rows) {
			index++;
			if (this.cancelling) break;

			const iterations = this.sim.getIterations();
			row.statusElem.textContent = 'starting';
			row.rowElem.classList.add('stat-weights-running');
			this.resultsViewer.setContent(`<div class="results-sim">Running ${specNames[row.spec]}<br>${index} of ${rows.length} specs</div>`);

			const specStart = performance.now();
			let result: StatWeightsResult;
			try {
				result = await runSpecStatWeights(this.sim, row.spec, iterations, (progress: ProgressMetrics) => {
					row.statusElem.textContent = progress.totalIterations
						? `${Math.min(100, Math.floor((100 * progress.completedIterations) / progress.totalIterations))}%`
						: 'running';
				});
			} catch (error: any) {
				console.error(error);
				row.statusElem.textContent = 'failed';
				row.rowElem.classList.remove('stat-weights-running');
				row.rowElem.classList.add('stat-weights-failed');
				new Toast({ variant: 'error', body: `${specNames[row.spec]}: ${error?.message || 'stat weights failed.'}` });
				break;
			}

			row.rowElem.classList.remove('stat-weights-running');
			if (result.error) {
				// A cancel comes back as an error on the result rather than as a throw.
				row.statusElem.textContent = '';
				break;
			}

			row.result = result;
			row.iterations = iterations;
			row.durationMs = performance.now() - specStart;
			row.statusElem.textContent = `${(row.durationMs / 1000).toFixed(1)}s`;
			this.updateRow(row, this.displayedStats());
			this.updateNotes();
		}

		console.log(`Stat weights for ${rows.length} spec(s) took ${((performance.now() - startTime) / 1000).toFixed(1)}s.`);
		this.resultsViewer.hideAll();
		this.setRunning(false);
	}
}
