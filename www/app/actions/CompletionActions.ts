/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import * as Constants from '../Constants';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import * as Csrf from '../Csrf';
import Loader from '../Loader';
import * as CompletionTypes from '../types/CompletionTypes';
import CompletionStore from '../stores/CompletionStore';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;
let lastSyncTime: number | null = null;
let syncInProgress: boolean = false;

export function sync(): Promise<void> {
	if (syncInProgress) {
		return Promise.resolve();
	}

	syncInProgress = true;
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	return new Promise<void>((resolve, reject): void => {
		try {
			SuperAgent
				.get('/completion')
				.query({
					...CompletionStore.filter,
				})
				.set('Accept', 'application/json')
				.set('Csrf-Token', Csrf.token)
				.end((err: any, res: SuperAgent.Response): void => {
					syncInProgress = false;

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
						Alert.errorRes(res, 'Failed to load completion data');
						reject(err);
						return;
					}

					Dispatcher.dispatch({
						type: CompletionTypes.SYNC,
						data: {
							completion: res.body,
						},
					});

					lastSyncTime = Date.now();
					resolve();
				});
		} catch (e) {
			syncInProgress = false;
			reject(e);
		}
	});
}

export function lastSync(): number | null {
	return lastSyncTime;
}

export function filter(filt: CompletionTypes.Filter): Promise<void> {
	Dispatcher.dispatch({
		type: CompletionTypes.FILTER,
		data: {
			filter: filt,
		},
	});

	return sync();
}

EventDispatcher.register((action: CompletionTypes.CompletionDispatch) => {
	switch (action.type) {
		case CompletionTypes.CHANGE:
				sync();
			break;
	}
});
