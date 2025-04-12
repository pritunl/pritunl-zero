/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import * as Csrf from '../Csrf';
import Loader from '../Loader';
import * as PolicyTypes from '../types/PolicyTypes';
import PoliciesStore from '../stores/PoliciesStore';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;

export function sync(): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/policy')
			.query({
				...PoliciesStore.filter,
				page: PoliciesStore.page,
				page_count: PoliciesStore.pageCount,
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
					Alert.errorRes(res, 'Failed to load policies');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: PolicyTypes.SYNC,
					data: {
						policies: res.body,
						count: res.body.count,
					},
				});

				resolve();
			});
	});
}

export function traverse(page: number): Promise<void> {
	Dispatcher.dispatch({
		type: PolicyTypes.TRAVERSE,
		data: {
			page: page,
		},
	});
	return sync();
}

export function filter(filt: PolicyTypes.Filter): Promise<void> {
	Dispatcher.dispatch({
		type: PolicyTypes.FILTER,
		data: {
			filter: filt,
		},
	});
	return sync();
}

export function commit(policy: PolicyTypes.Policy): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.put('/policy/' + policy.id)
			.send(policy)
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
					Alert.errorRes(res, 'Failed to save policy');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

export function create(policy: PolicyTypes.Policy): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.post('/policy')
			.send(policy)
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
					Alert.errorRes(res, 'Failed to create policy');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

export function remove(policyId: string): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.delete('/policy/' + policyId)
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
					Alert.errorRes(res, 'Failed to delete policies');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

export function removeMulti(authorityIds: string[]): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.delete('/authority')
			.send(authorityIds)
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
					Alert.errorRes(res, 'Failed to delete authorities');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

EventDispatcher.register((action: PolicyTypes.PolicyDispatch) => {
	switch (action.type) {
		case PolicyTypes.CHANGE:
			sync();
			break;
	}
});
