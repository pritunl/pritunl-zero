/// <reference path="References.d.ts"/>
import * as React from 'react';
import * as ReactDOM from 'react-dom';
import * as Blueprint from '@blueprintjs/core';
import * as Csrf from './Csrf';

document.body.className = 'pt-dark';

Csrf.load().then((): void => {
	Blueprint.FocusStyleManager.onlyShowFocusOnTabs();

	ReactDOM.render(
		<div>Test</div>,
		document.getElementById('app'),
	);
});
