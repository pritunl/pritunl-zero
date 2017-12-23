/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import * as Csrf from '../Csrf';
import Loader from '../Loader';
import * as SettingsTypes from '../types/SettingsTypes';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;

export function sync(): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/settings')
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (res && res.status === 401) {
					window.location.href = '/login';
					resolve();
					return;
				}

				if (curSyncId !== syncId) {
					resolve();
					return;
				}

				if (err) {
					Alert.errorRes(res, 'Failed to sync builds');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: SettingsTypes.SYNC,
					data: res.body,
				});

				resolve();
			});
	});
}

export function commit(
		settings: SettingsTypes.Settings): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.put('/settings')
			.send(settings)
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (res && res.status === 401) {
					window.location.href = '/login';
					resolve();
					return;
				}

				if (err) {
					Alert.errorRes(res, 'Failed to commit settings');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: SettingsTypes.SYNC,
					data: res.body,
				});

				resolve();
			});
	});
}

EventDispatcher.register((action: SettingsTypes.SettingsDispatch) => {
	switch (action.type) {
		case SettingsTypes.CHANGE:
			sync();
			break;
	}
});
