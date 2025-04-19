/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import * as Csrf from '../Csrf';
import Loader from '../Loader';
import * as AuthorityTypes from '../types/AuthorityTypes';
import AuthoritiesStore from '../stores/AuthoritiesStore';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;

export function sync(): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/authority')
			.query({
				...AuthoritiesStore.filter,
				page: AuthoritiesStore.page,
				page_count: AuthoritiesStore.pageCount,
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
					Alert.errorRes(res, 'Failed to load authorities');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: AuthorityTypes.SYNC,
					data: {
						authorities: res.body.authorities,
						count: res.body.count,
					},
				});

				resolve();
			});
	});
}

export function traverse(page: number): Promise<void> {
	Dispatcher.dispatch({
		type: AuthorityTypes.TRAVERSE,
		data: {
			page: page,
		},
	});
	return sync();
}

export function filter(filt: AuthorityTypes.Filter): Promise<void> {
	Dispatcher.dispatch({
		type: AuthorityTypes.FILTER,
		data: {
			filter: filt,
		},
	});
	return sync();
}

export function clearSecret(id: string): void {
	Dispatcher.dispatch({
		type: AuthorityTypes.SYNC_SECRET,
		data: {
			id: id,
			secret: null,
		},
	});
}

export function commit(authority: AuthorityTypes.Authority): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.put('/authority/' + authority.id)
			.send(authority)
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
					Alert.errorRes(res, 'Failed to save authority');
					reject(err);
					return;
				}

				if (res.body && res.body.hsm_secret) {
					Dispatcher.dispatch({
						type: AuthorityTypes.SYNC_SECRET,
						data: {
							id: res.body.id,
							secret: res.body.hsm_secret,
						},
					});
				}

				resolve();
			});
	});
}

export function create(authority: AuthorityTypes.Authority): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.post('/authority')
			.send(authority)
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
					Alert.errorRes(res, 'Failed to create authority');
					reject(err);
					return;
				}

				if (res.body && res.body.hsm_secret) {
					Dispatcher.dispatch({
						type: AuthorityTypes.SYNC_SECRET,
						data: {
							id: res.body.id,
							secret: res.body.hsm_secret,
						},
					});
				}

				resolve();
			});
	});
}

export function remove(authorityId: string): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.delete('/authority/' + authorityId)
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

export function createToken(authorityId: string): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.post('/authority/' + authorityId + '/token')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (res && res.status === 401) {
					window.location.href = '/login';
					resolve();
					return;
				}

				if (err) {
					Alert.errorRes(res, 'Failed to create authority token');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

export function deleteToken(authorityId: string,
		token: string): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.delete('/authority/' + authorityId + '/token/' + token)
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (res && res.status === 401) {
					window.location.href = '/login';
					resolve();
					return;
				}

				if (err) {
					Alert.errorRes(res, 'Failed to delete authority token');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

EventDispatcher.register((action: AuthorityTypes.AuthorityDispatch) => {
	switch (action.type) {
		case AuthorityTypes.CHANGE:
			sync();
			break;
	}
});
