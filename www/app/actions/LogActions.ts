/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import * as Csrf from '../Csrf';
import Loader from '../Loader';
import * as LogTypes from '../types/LogTypes';
import LogsStore from '../stores/LogsStore';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;

export function sync(): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/log')
			.query({
				...LogsStore.filter,
				page: LogsStore.page,
				page_count: LogsStore.pageCount,
			})
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (res.status === 401) {
					window.location.href = '/login';
					resolve();
					return;
				}

				if (curSyncId !== syncId) {
					resolve();
					return;
				}

				if (err) {
					Alert.errorRes(res, 'Failed to load logs');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: LogTypes.SYNC,
					data: {
						logs: res.body.logs,
						count: res.body.count,
					},
				});

				resolve();
			});
	});
}

export function traverse(page: number): Promise<void> {
	Dispatcher.dispatch({
		type: LogTypes.TRAVERSE,
		data: {
			page: page,
		},
	});

	return sync();
}

export function filter(filt: LogTypes.Filter): Promise<void> {
	Dispatcher.dispatch({
		type: LogTypes.FILTER,
		data: {
			filter: filt,
		},
	});

	return sync();
}

EventDispatcher.register((action: LogTypes.LogDispatch) => {
	switch (action.type) {
		case LogTypes.CHANGE:
			sync();
			break;
	}
});
