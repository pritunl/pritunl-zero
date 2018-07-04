/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import * as StateTypes from '../types/StateTypes';

export function setSshToken(sshToken: string): void {
	if (sshToken === null) {
		window.history.replaceState(
			null, null, window.location.pathname);
	}

	Dispatcher.dispatch({
		type: StateTypes.SSH_TOKEN,
		data: sshToken,
	});
}

export function setSshDevice(sshDevice: string): void {
	if (sshDevice === null) {
		window.history.replaceState(
			null, null, window.location.pathname);
	}

	Dispatcher.dispatch({
		type: StateTypes.SSH_DEVICE,
		data: sshDevice,
	});
}
