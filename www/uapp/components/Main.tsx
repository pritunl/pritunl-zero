/// <reference path="../References.d.ts"/>
import * as React from 'react';
import LoadingBar from './LoadingBar';
import Session from './Session';
import Validate from './Validate';

const css = {
	card: {
		padding: '20px 15px',
		minWidth: '260px',
		maxWidth: '320px',
		margin: '0 auto',
		position: 'absolute',
		top: '50%',
		left: '50%',
		width: '100%',
		transform: 'translate(-50%, -50%)',
	} as React.CSSProperties,
	loading: {
		position: 'absolute',
		top: 0,
		left: 0,
		right: 0,
	} as React.CSSProperties,
};

export default class Main extends React.Component<{}, {}> {
	render(): JSX.Element {
		let sshToken = '';
		let query = window.location.search.substring(1);
		let vals = query.split('&');
		for (let val of vals) {
			let keyval = val.split('=');
			if (keyval[0] === 'ssh-token') {
				sshToken = keyval[1];
			}
		}

		let bodyElm: JSX.Element;

		if (sshToken) {
			bodyElm = <Validate token={sshToken}/>;
		} else {
			bodyElm = <Session/>;
		}

		return <div className="pt-card pt-elevation-2" style={css.card}>
			<LoadingBar style={css.loading} intent="primary"/>
			{bodyElm}
		</div>;
	}
}
