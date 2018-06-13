/// <reference path="References.d.ts"/>
import * as React from 'react';
import * as ReactDOM from 'react-dom';
import * as Blueprint from '@blueprintjs/core';
import Main from './components/Main';
import * as Alert from './Alert';
import * as Csrf from './Csrf';

document.body.className = 'pt-dark';

Csrf.load().then((): void => {
	Blueprint.FocusStyleManager.onlyShowFocusOnTabs();
	Alert.init();

	ReactDOM.render(
		<div><Main/></div>,
		document.getElementById('app'),
	);
});
