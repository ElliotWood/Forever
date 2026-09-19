import tippy from 'tippy.js';

import { SITE_REPO_URL } from '../constants/other.js';
import { Component } from './component';

export class SocialLinks extends Component {
	static buildGitHubLink(): Element {
		const anchor = (
			<a href={SITE_REPO_URL} target="_blank" className="github-link link-alt" dataset={{ tippyContent: 'Contribute on GitHub' }}>
				<i className="fab fa-github fa-lg" />
			</a>
		);
		tippy(anchor);
		return anchor;
	}
}
