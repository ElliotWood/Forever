import clsx from 'clsx';

import { getCommunityBuilds } from '../community_builds.js';
import { phaseNames } from '../constants/other.js';
import { getLaunchedSimsForClass, LaunchStatus, raidSimStatus, simLaunchStatuses } from '../launched_sims.js';
import { Class, Spec } from '../proto/common.js';
import {
	classIcons,
	classNames,
	getSpecSiteUrl,
	naturalClassOrder,
	raidSimIcon,
	raidSimLabel,
	raidSimSiteUrl,
	specNames,
	specToClass,
	textCssClassForClass,
	textCssClassForSpec,
	titleIcons,
} from '../proto_utils/utils.js';
import { Component } from './component.js';

interface ClassOptions {
	type: 'Class';
	index: Class;
}

interface SpecOptions {
	type: 'Spec';
	index: Spec;
}

interface RaidOptions {
	type: 'Raid';
}

type SimTitleDropdownConfig = {
	noDropdown?: boolean;
};

// Dropdown menu for selecting a player.
export class SimTitleDropdown extends Component {
	private readonly dropdownMenu: HTMLElement | undefined;

	private readonly specLabels: Record<Spec, string> = {
		[Spec.SpecBalanceDruid]: 'Balance',
		[Spec.SpecFeralDruid]: 'Feral DPS',
		[Spec.SpecFeralTankDruid]: 'Feral Tank',
		[Spec.SpecRestorationDruid]: 'Restoration',
		[Spec.SpecElementalShaman]: 'Elemental',
		[Spec.SpecEnhancementShaman]: 'Enhancement',
		[Spec.SpecRestorationShaman]: 'Restoration',
		[Spec.SpecHunter]: 'Hunter',
		[Spec.SpecMage]: 'Mage',
		[Spec.SpecRogue]: 'DPS',
		[Spec.SpecHolyPaladin]: 'Holy',
		[Spec.SpecProtectionPaladin]: 'Protection',
		[Spec.SpecRetributionPaladin]: 'Retribution',
		[Spec.SpecHealingPriest]: 'Healing',
		[Spec.SpecShadowPriest]: 'Shadow',
		[Spec.SpecSmitePriest]: 'Smite',
		[Spec.SpecWarlock]: 'DPS',
		[Spec.SpecWarrior]: 'DPS',
		[Spec.SpecTankWarrior]: 'Tank',
	};

	constructor(parent: HTMLElement, currentSpecIndex: Spec | null, config: SimTitleDropdownConfig = {}) {
		super(parent, 'sim-title-dropdown-root');

		const rootLinkArgs: SpecOptions | RaidOptions = currentSpecIndex === null ? { type: 'Raid' } : { type: 'Spec', index: currentSpecIndex };
		const rootLink = this.buildRootSimLink(rootLinkArgs);

		if (config.noDropdown) {
			this.rootElem.appendChild(rootLink);
			return;
		}

		const dropdownMenu = (<ul className="dropdown-menu"></ul>) as HTMLElement;
		this.rootElem.appendChild(
			<div className="dropdown sim-link-dropdown">
				{rootLink}
				{dropdownMenu}
			</div>,
		);

		this.dropdownMenu = dropdownMenu;
		this.buildDropdown();
	}

	private buildDropdown() {
		// Compares the status rather than the object it sits on, which is what kept this
		// commented out: raidSimStatus is a {phase, status} pair and never ordered against
		// a LaunchStatus.
		if (raidSimStatus.status >= LaunchStatus.Alpha) {
			// Add the raid sim to the top of the dropdown
			this.dropdownMenu?.appendChild(<li>{this.buildRaidLink()}</li>);
		}

		naturalClassOrder.forEach(classIndex => {
			const sims = getLaunchedSimsForClass(classIndex);
			if (!sims.length) return;

			// A class with a single sim and no builds to list under it is linked directly.
			// Anything else - several specs, or builds - gets a dropdown of its own.
			const directLink = sims.length == 1 && !getCommunityBuilds(sims[0]).length;
			const link = directLink ? this.buildClassLink(classIndex, false) : this.buildClassDropdown(classIndex);
			this.dropdownMenu?.appendChild(<li>{link}</li>);
		});
	}

	private buildClassDropdown(classIndex: Class): Element {
		const dropdownMenu = (<ul className="dropdown-menu"></ul>) as HTMLElement;

		// Each of the class's specs, and under each spec the community builds the landing
		// page lists for it.
		getLaunchedSimsForClass(classIndex).forEach(specIndex => {
			dropdownMenu.appendChild(<li>{this.buildSpecLink(specIndex)}</li>);
			getCommunityBuilds(specIndex).forEach(build => {
				dropdownMenu.appendChild(<li>{this.buildBuildLink(specIndex, build)}</li>);
			});
		});

		return (
			<div className="dropend sim-link-dropdown">
				{this.buildClassLink(classIndex, true)}
				{dropdownMenu}
			</div>
		);
	}

	private buildRootSimLink(data: SpecOptions | RaidOptions): Element {
		let label;

		if (data.type == 'Raid') label = raidSimLabel;
		else {
			const classIndex = specToClass[data.index];
			if (getLaunchedSimsForClass(classIndex).length > 1)
				// If the class has multiple sims, use the spec name
				label = specNames[data.index];
			// If the class has only 1 sim, use the class name
			else label = classNames[classIndex];
		}

		// Each class in this menu opens a menu of its own, so this one has to stay open
		// while that happens. Bootstrap otherwise closes a menu on any click it sees,
		// including a click on something inside it; autoClose outside leaves that to
		// clicks that land elsewhere.
		return (
			<a
				href="javascript:void(0)"
				className={clsx('sim-link', this.getContextualKlass(data))}
				dataset={{ bsToggle: 'dropdown', bsTrigger: 'click', bsAutoClose: 'outside' }}>
				<div className="sim-link-content">
					<img src={this.getSimIconPath(data)} className="sim-link-icon" />
					<div className="d-flex flex-column">
						<span className="sim-link-label text-white">Forever Sim (unofficial)</span>
						<span className="sim-link-title">{label}</span>
						{this.launchStatusLabel(data)}
					</div>
				</div>
			</a>
		);
	}

	private buildRaidLink(): Element {
		return (
			<a href={raidSimSiteUrl} className={clsx('sim-link', this.getContextualKlass({ type: 'Raid' }))}>
				<div className="sim-link-content">
					<img src={this.getSimIconPath({ type: 'Raid' })} className="sim-link-icon" />
					<div className="d-flex flex-column">
						<span className="sim-link-title">{raidSimLabel}</span>
						{this.launchStatusLabel({ type: 'Raid' })}
					</div>
				</div>
			</a>
		);
	}

	private buildClassLink(classIndex: Class, isDropdownToggle: boolean): Element {
		const specIndexes = getLaunchedSimsForClass(classIndex);
		const href = isDropdownToggle ? 'javascript:void(0)' : getSpecSiteUrl(specIndexes[0]);

		return (
			<a
				href={href}
				className={clsx('sim-link', this.getContextualKlass({ type: 'Class', index: classIndex }))}
				dataset={isDropdownToggle ? { bsToggle: 'dropdown' } : {}}>
				<div className="sim-link-content">
					<img src={this.getSimIconPath({ type: 'Class', index: classIndex })} className="sim-link-icon" />
					<div className="d-flex flex-column">
						<span className="sim-link-title">{classNames[classIndex]}</span>
						{!isDropdownToggle && this.launchStatusLabel({ type: 'Spec', index: specIndexes[0] })}
					</div>
				</div>
			</a>
		);
	}

	private buildSpecLink(specIndex: Spec): Element {
		const href = getSpecSiteUrl(specIndex);

		return (
			<a href={href} className={clsx('sim-link', this.getContextualKlass({ type: 'Spec', index: specIndex }))}>
				<div className="sim-link-content">
					<img src={this.getSimIconPath({ type: 'Spec', index: specIndex })} className="sim-link-icon" />
					<div className="d-flex flex-column">
						<span className="sim-link-label">{classNames[specToClass[specIndex]]}</span>
						<span className="sim-link-title">{this.specLabels[specIndex]}</span>
						{this.launchStatusLabel({ type: 'Spec', index: specIndex })}
					</div>
				</div>
			</a>
		);
	}

	// A community build listed under its spec, linking to that spec's sim with the build
	// applied. The same link the landing page lists, down to the ?build= name.
	private buildBuildLink(specIndex: Spec, build: string): Element {
		const href = `${getSpecSiteUrl(specIndex)}?build=${encodeURIComponent(build)}`;

		return (
			<a href={href} className={clsx('sim-link', 'sim-link-build', this.getContextualKlass({ type: 'Spec', index: specIndex }))}>
				<div className="sim-link-content">
					<img src={this.getSimIconPath({ type: 'Spec', index: specIndex })} className="sim-link-icon" />
					<div className="d-flex flex-column">
						<span className="sim-link-label">{this.specLabels[specIndex]}</span>
						<span className="sim-link-title">{build}</span>
						<span className="launch-status-label text-brand">Community build</span>
					</div>
				</div>
			</a>
		);
	}

	private launchStatusLabel(data: SpecOptions | RaidOptions): Element {
		const status = data.type == 'Raid' ? raidSimStatus.status : simLaunchStatuses[data.index].status;
		const phase = data.type == 'Raid' ? raidSimStatus.phase : simLaunchStatuses[data.index].phase;

		return (
			<span className="launch-status-label text-brand">
				{status === LaunchStatus.Unlaunched ? (
					<>Not Yet Supported</>
				) : (
					<>
						{phaseNames[phase]}
						{status != LaunchStatus.Launched && <> - {LaunchStatus[status]}</>}
					</>
				)}
			</span>
		);
	}

	private getSimIconPath(data: ClassOptions | SpecOptions | RaidOptions): string {
		let iconPath: string;

		if (data.type == 'Raid') {
			iconPath = raidSimIcon;
		} else if (data.type == 'Class') {
			iconPath = classIcons[data.index];
		} else {
			iconPath = titleIcons[data.index];
		}

		return iconPath;
	}

	private getContextualKlass(data: ClassOptions | SpecOptions | RaidOptions): string {
		if (data.type == 'Raid')
			// Raid link
			return 'text-white';
		else if (data.type == 'Class')
			// Class links
			return textCssClassForClass(data.index);
		else return textCssClassForSpec(data.index);
	}
}
