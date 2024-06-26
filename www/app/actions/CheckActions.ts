/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import * as Csrf from '../Csrf';
import Loader from '../Loader';
import ChecksStore from '../stores/ChecksStore';
import * as CheckTypes from '../types/CheckTypes';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;
let dataSyncReqs: {[key: string]: SuperAgent.Request} = {};

export function sync(): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/checks')
			.query({
				...ChecksStore.filter,
				page: ChecksStore.page,
				page_count: ChecksStore.pageCount,
			})
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
					Alert.errorRes(res, 'Failed to load checks');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: CheckTypes.SYNC,
					data: {
						checks: res.body.checks,
						count: res.body.count,
					},
				});

				resolve();
			});
	});
}

export function traverse(page: number): Promise<void> {
	Dispatcher.dispatch({
		type: CheckTypes.TRAVERSE,
		data: {
			page: page,
		},
	});

	return sync();
}

export function filter(filt: CheckTypes.Filter): Promise<void> {
	Dispatcher.dispatch({
		type: CheckTypes.FILTER,
		data: {
			filter: filt,
		},
	});

	return sync();
}

export function commit(check: CheckTypes.Check): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.put('/checks/' + check.id)
			.send(check)
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
					Alert.errorRes(res, 'Failed to save check');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

export function create(check: CheckTypes.Check): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.post('/checks')
			.send(check)
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
					Alert.errorRes(res, 'Failed to create check');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

export function remove(checkId: string): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.delete('/checks/' + checkId)
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					Alert.errorRes(res, 'Failed to delete checks');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

export function removeMulti(checkIds: string[]): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.delete('/checks')
			.send(checkIds)
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
					Alert.errorRes(res, 'Failed to delete checks');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

export function chart(checkId: string, resource: string,
	period: number, interval: number): Promise<any> {
	let curDataSyncId = MiscUtils.uuid();

	let loader = new Loader().loading();

	resource = resource.replace(/[0-9]/g, '');

	return new Promise<any>((resolve, reject): void => {
		let req = SuperAgent.get('/checks/' + checkId + '/chart')
			.query({
				resource: resource,
				period: period.toString(),
				interval: interval.toString(),
			})
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.on('abort', () => {
				loader.done();
				resolve(null);
			});
		dataSyncReqs[curDataSyncId] = req;

		req.end((err: any, res: SuperAgent.Response): void => {
			delete dataSyncReqs[curDataSyncId];
			loader.done();

			if (res && res.status === 401) {
				window.location.href = '/login';
				resolve(null);
				return;
			}

			if (err) {
				Alert.errorRes(res, 'Failed to load check chart');
				reject(err);
				return;
			}

			resolve(res.body);
		});
	});
}

export function log(checkId: string, resource: string): Promise<any> {
	let curDataSyncId = MiscUtils.uuid();

	let loader = new Loader().loading();

	return new Promise<any>((resolve, reject): void => {
		let req = SuperAgent.get('/checks/' + checkId + '/log')
			.query({
				resource: resource,
			})
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.on('abort', () => {
				loader.done();
				resolve(null);
			});
		dataSyncReqs[curDataSyncId] = req;

		req.end((err: any, res: SuperAgent.Response): void => {
			delete dataSyncReqs[curDataSyncId];
			loader.done();

			if (res && res.status === 401) {
				window.location.href = '/login';
				resolve(null);
				return;
			}

			if (err) {
				Alert.errorRes(res, 'Failed to load check log');
				reject(err);
				return;
			}

			resolve(res.body);
		});
	});
}

export function dataCancel(): void {
	for (let [key, val] of Object.entries(dataSyncReqs)) {
		val.abort();
	}
}

EventDispatcher.register((action: CheckTypes.CheckDispatch) => {
	switch (action.type) {
		case CheckTypes.CHANGE:
			sync();
			break;
	}
});
