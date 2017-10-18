/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import * as Csrf from '../Csrf';
import Loader from '../Loader';
import * as UserTypes from '../types/UserTypes';
import UserStore from '../stores/UserStore';
import UsersStore from '../stores/UsersStore';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;

export function load(userId: string): Promise<void> {
	if (!userId) {
		let user: UserTypes.User = {
			id: null,
			type: 'local',
			roles: [],
			permissions: [],
		};

		Dispatcher.dispatch({
			type: UserTypes.LOAD,
			data: {
				user: user,
			},
		});

		return Promise.resolve();
	}

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/user/' + userId)
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (res.status === 401) {
					window.location.href = '/login';
					resolve();
					return;
				}

				if (err) {
					Alert.errorRes(res, 'Failed to load user');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: UserTypes.LOAD,
					data: {
						user: res.body,
					},
				});

				resolve();
			});
	});
}

export function reload(): Promise<void> {
	return load(UserStore.user ? UserStore.user.id : null);
}

export function unload(): void {
	Dispatcher.dispatch({
		type: UserTypes.UNLOAD,
	});
}

export function sync(): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/user')
			.query({
				...UsersStore.filter,
				page: UsersStore.page,
				page_count: UsersStore.pageCount,
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
					Alert.errorRes(res, 'Failed to load users');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: UserTypes.SYNC,
					data: {
						users: res.body.users,
						count: res.body.count,
					},
				});

				resolve();
			});
	});
}

export function traverse(page: number): Promise<void> {
	Dispatcher.dispatch({
		type: UserTypes.TRAVERSE,
		data: {
			page: page,
		},
	});

	return sync();
}

export function filter(filt: UserTypes.Filter): Promise<void> {
	Dispatcher.dispatch({
		type: UserTypes.FILTER,
		data: {
			filter: filt,
		},
	});

	return sync();
}

export function commit(user: UserTypes.User): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.put('/user/' + user.id)
			.send(user)
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (res.status === 401) {
					window.location.href = '/login';
					resolve();
					return;
				}

				if (err) {
					Alert.errorRes(res, 'Failed to save user');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: UserTypes.LOAD,
					data: {
						user: res.body,
					},
				});

				resolve();
			});
	});
}

export function create(user: UserTypes.User): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.post('/user')
			.send(user)
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (res.status === 401) {
					window.location.href = '/login';
					resolve();
					return;
				}

				if (err) {
					Alert.errorRes(res, 'Failed to create user');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

export function remove(userIds: string[]): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.delete('/user')
			.send(userIds)
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (res.status === 401) {
					window.location.href = '/login';
					resolve();
					return
				}

				if (err) {
					Alert.errorRes(res, 'Failed to delete users');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

EventDispatcher.register((action: UserTypes.UserDispatch) => {
	switch (action.type) {
		case UserTypes.CHANGE:
			sync();
			break;
	}
});
