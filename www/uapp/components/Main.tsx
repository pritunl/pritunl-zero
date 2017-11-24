/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Session from './Session';
import Validate from './Validate';

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
		let query = window.location.search.substring(1);
		let vals = query.split('&');
		for (let val of vals) {
			let keyval = val.split('=');
			console.log(val);
			if (keyval[0] === 'ssh-token') {
				sshToken = keyval[1];
				break;
			}
		}

		return <div className="pt-card pt-elevation-2" style={css.card}>
			{sshToken ? <Validate token={sshToken}/> : <Session/>}
		</div>;
	}
}
