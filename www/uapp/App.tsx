/// <reference path="References.d.ts"/>
import * as React from 'react';
import * as ReactDOM from 'react-dom';
import * as Blueprint from '@blueprintjs/core';
import * as StateActions from './actions/StateActions';
import Main from './components/Main';
import * as Alert from './Alert';
import * as Csrf from './Csrf';

document.body.className = 'pt-dark';

Csrf.load().then((): void => {
	Blueprint.FocusStyleManager.onlyShowFocusOnTabs();
	Alert.init();

	let query = window.location.search.substring(1);

	let vals = query.split('&');
	for (let val of vals) {
		let keyval = val.split('=');
		if (keyval[0] === 'ssh-token') {
			StateActions.setSshToken(keyval[1]);
		} else if (keyval[0] === 'device') {
			StateActions.setSshDevice(keyval[1]);
		}
	}

	ReactDOM.render(
		<div><Main/></div>,
		document.getElementById('app'),
	);
});
