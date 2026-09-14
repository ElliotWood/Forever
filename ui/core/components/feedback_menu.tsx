import { ref } from 'tsx-vanilla';

import { SITE_REPO_URL, SITE_VERSION, SPEC_DIRECTORY } from '../constants/other.js';
import { Ruleset } from '../proto/api.js';
import { SimUI } from '../sim_ui.js';
import { BaseModal } from './base_modal.jsx';
import Toast from './toast.jsx';

// Feedback becomes a Github issue on the Forever repo. The site is static, so there is no
// server to post to and no token we could ship without publishing it to the world: instead
// we prefill Github's own issue composer and let the reporter press submit.
const ISSUE_URL = `${SITE_REPO_URL}/issues/new`;

// Github truncates very long issue urls, so leave room for the diagnostics and the url itself.
const MAX_DESCRIPTION_LENGTH = 5000;

const PASTE_HINT = 'Screenshot copied. Paste it into the Github issue with Ctrl+V (Cmd+V on a Mac).';
const BODY_PASTE_HINT = 'Screenshot: paste it here with Ctrl+V (Cmd+V on a Mac) - it is already on your clipboard.';
const BODY_NO_SCREENSHOT_HINT = 'Screenshot: none was attached automatically. Paste one here if you have it.';

export class FeedbackMenu extends BaseModal {
	private readonly simUI: SimUI;

	private readonly textElem: HTMLTextAreaElement;
	private readonly statusElem: HTMLElement;
	private readonly captureButton: HTMLButtonElement;
	private readonly openIssueLink: HTMLAnchorElement;

	constructor(parent: HTMLElement, simUI: SimUI) {
		super(parent, 'feedback-menu', { title: 'Send Feedback', footer: true });
		this.simUI = simUI;

		const textElemRef = ref<HTMLTextAreaElement>();
		const statusElemRef = ref<HTMLDivElement>();
		const captureButtonRef = ref<HTMLButtonElement>();
		const openIssueLinkRef = ref<HTMLAnchorElement>();

		this.body.replaceChildren(
			<>
				<div className="feedback-description">
					<p>
						Describe what you were doing, what you expected and what happened instead. Your sim page, ruleset, browser and url are added to the
						report automatically.
					</p>
					<p>
						Capturing a screenshot copies it to your clipboard - your browser will ask which tab or window to share, so pick this one. Github then
						opens in a new tab with everything filled in, and you paste the screenshot into the issue with Ctrl+V (Cmd+V on a Mac) before submitting
						it.
					</p>
				</div>
				<textarea ref={textElemRef} className="feedback-textarea form-control" placeholder="What went wrong, or what would you like to see?"></textarea>
				<div ref={statusElemRef} className="feedback-status form-text" hidden></div>
			</>,
		);

		this.footer!.appendChild(
			<>
				<a ref={openIssueLinkRef} href="javascript:void(0)" className="feedback-button btn btn-link" target="_blank">
					Open Issue
				</a>
				<button ref={captureButtonRef} className="feedback-button btn btn-primary capture-button">
					<i className="fas fa-camera me-1"></i>
					Capture Screenshot & Open Issue
				</button>
			</>,
		);

		this.textElem = textElemRef.value!;
		this.statusElem = statusElemRef.value!;
		this.captureButton = captureButtonRef.value!;
		this.openIssueLink = openIssueLinkRef.value!;
		this.openIssueLink.href = this.buildIssueURL(false);

		// The plain link is a real anchor so it can never be caught by a popup blocker, which
		// means its target has to follow whatever is currently typed in the box.
		this.textElem.addEventListener('input', () => (this.openIssueLink.href = this.buildIssueURL(false)));
		this.captureButton.addEventListener('click', () => this.onCapture());
	}

	private async onCapture() {
		this.captureButton.disabled = true;
		this.setStatus('Waiting for you to choose what to share...');

		let screenshotAttached = false;
		let status = PASTE_HINT;
		try {
			await this.copyScreenshotToClipboard();
			screenshotAttached = true;
		} catch (error: any) {
			// Never drop the screenshot silently: say what went wrong, then file the report anyway.
			const reason = error?.message || String(error);
			status = `Screenshot could not be attached automatically (${reason}), so the report was prepared without it.`;
			new Toast({ variant: 'warning', body: `Screenshot could not be attached: ${reason}` });
		}
		this.captureButton.disabled = false;

		const url = this.buildIssueURL(screenshotAttached);
		this.openIssueLink.href = url;
		if (!window.open(url, '_blank', 'noopener')) {
			status += ' Your browser blocked the new tab, so use the Open Issue button to file it.';
		}
		this.setStatus(status);
	}

	// Screenshots come from the browser's own tab capture rather than a rendering library, so
	// there is no extra dependency and what gets attached is exactly what the reporter sees.
	private async copyScreenshotToClipboard() {
		if (!navigator.mediaDevices?.getDisplayMedia) throw new Error('this browser cannot capture the screen');
		if (!navigator.clipboard?.write || typeof ClipboardItem == 'undefined') throw new Error('this browser cannot copy images to the clipboard');

		const stream = await navigator.mediaDevices.getDisplayMedia({
			video: { displaySurface: 'browser' },
			audio: false,
			// Chromium-only hint that makes the current tab the default choice in the picker.
			preferCurrentTab: true,
		} as DisplayMediaStreamOptions);

		let screenshot: Blob;
		try {
			screenshot = await this.grabFrame(stream);
		} finally {
			stream.getTracks().forEach(track => track.stop());
		}

		await navigator.clipboard.write([new ClipboardItem({ [screenshot.type]: screenshot })]);
	}

	private async grabFrame(stream: MediaStream): Promise<Blob> {
		const video = document.createElement('video');
		video.srcObject = stream;
		video.muted = true;
		await video.play();
		// The first frame is not necessarily painted when play() resolves.
		await new Promise(resolve => requestAnimationFrame(resolve));

		const canvas = document.createElement('canvas');
		canvas.width = video.videoWidth;
		canvas.height = video.videoHeight;
		const context = canvas.getContext('2d');
		if (!canvas.width || !canvas.height || !context) throw new Error('the captured frame could not be read');
		context.drawImage(video, 0, 0, canvas.width, canvas.height);
		video.srcObject = null;

		return await new Promise<Blob>((resolve, reject) =>
			canvas.toBlob(blob => (blob ? resolve(blob) : reject(new Error('the screenshot could not be encoded'))), 'image/png'),
		);
	}

	private setStatus(status: string) {
		this.statusElem.textContent = status;
		this.statusElem.hidden = false;
	}

	private buildIssueURL(screenshotAttached: boolean): string {
		const params = new URLSearchParams({
			title: this.buildIssueTitle(),
			body: this.buildIssueBody(screenshotAttached),
		});
		return `${ISSUE_URL}?${params.toString()}`;
	}

	private buildIssueTitle(): string {
		const firstLine = this.description().split('\n')[0].trim();
		if (!firstLine) return `Feedback: ${this.simName()}`;
		return firstLine.length > 72 ? `${firstLine.substring(0, 69)}...` : firstLine;
	}

	private buildIssueBody(screenshotAttached: boolean): string {
		const description = this.description();
		return [
			description || '_No description given._',
			'',
			screenshotAttached ? BODY_PASTE_HINT : BODY_NO_SCREENSHOT_HINT,
			'',
			'### Diagnostics',
			'',
			`- Sim: ${this.simName()}`,
			`- Version: ${SITE_VERSION}`,
			`- Ruleset: ${this.simUI.sim.getRuleset() == Ruleset.RulesetForever ? 'Forever' : 'Classic'}`,
			`- Url: ${window.location.href}`,
			`- Viewport: ${window.innerWidth}x${window.innerHeight}`,
			`- Browser: ${navigator.userAgent}`,
		].join('\n');
	}

	private description(): string {
		const description = this.textElem.value.trim();
		return description.length > MAX_DESCRIPTION_LENGTH ? `${description.substring(0, MAX_DESCRIPTION_LENGTH)}...` : description;
	}

	private simName(): string {
		return SPEC_DIRECTORY || this.simUI.cssClass;
	}
}
