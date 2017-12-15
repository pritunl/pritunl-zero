/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Session from './Session';
import Validate from './Validate';
import Keybase from './Keybase';

const css = {
	card: {
		padding: '20px 10px',
		minWidth: '260px',
		maxWidth: '300px',
		margin: '0 auto',
		position: 'absolute',
		top: '50%',
		left: '50%',
		width: '100%',
		transform: 'translate(-50%, -50%)',
	} as React.CSSProperties,
};

export default class Main extends React.Component<{}, {}> {
	render(): JSX.Element {
		let sshToken = '';
		let keybaseToken = '';
		let keybaseSig = '';
		let query = window.location.search.substring(1);
		let vals = query.split('&');
		for (let val of vals) {
			let keyval = val.split('=');
			if (keyval[0] === 'ssh-token') {
				sshToken = keyval[1];
			} else if (keyval[0] === 'keybase-token') {
				keybaseToken = keyval[1];
			} else if (keyval[0] === 'keybase-sig') {
				keybaseSig = decodeURIComponent(keyval[1]).replace(/\+/g, ' ');
			}
		}

		let bodyElm: JSX.Element;

		if (sshToken) {
			bodyElm = <Validate token={sshToken}/>;
		} else if (keybaseToken && keybaseSig) {
			bodyElm = <Keybase token={keybaseToken} signature={keybaseSig}/>
		} else {
			bodyElm = <Session/>;
		}

		return <div className="pt-card pt-elevation-2" style={css.card}>
			{bodyElm}
		</div>;
	}
}
