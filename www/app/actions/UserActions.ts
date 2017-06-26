/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import Loader from '../Loader';
import * as UserTypes from '../types/UserTypes';
import UsersStore from '../stores/UsersStore';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;

export function load(userId: string): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/user/' + userId)
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					Alert.error('Failed to load user');
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

export function unload(): void {
	Dispatcher.dispatch({
		type: UserTypes.UNLOAD,
	});
}

function _sync(): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/user')
			.query({
				page: UsersStore.page,
				page_count: UsersStore.pageCount,
			})
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (curSyncId !== syncId) {
					resolve();
					return;
				}

				if (err) {
					Alert.error('Failed to sync users');
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

	return _sync();
}

export function sync(): Promise<void> {
	return _sync();
}

export function commit(user: UserTypes.User): Promise<UserTypes.User> {
	let loader = new Loader().loading();

	return new Promise<UserTypes.User>((resolve, reject): void => {
		SuperAgent
			.put('/user/' + user.id)
			.send(user)
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					Alert.error('Failed to commit user');
					reject(err);
					return;
				}

				resolve(res.body);
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
