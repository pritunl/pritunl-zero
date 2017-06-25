/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import Loader from '../Loader';
import * as UserTypes from '../types/UserTypes';
import UserStore from '../stores/UserStore';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;

export function get(userId: string): Promise<UserTypes.User> {
	let loader = new Loader().loading();

	return new Promise<UserTypes.User>((resolve, reject): void => {
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

				resolve(res.body);
			});
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
				page: UserStore.page,
				page_count: UserStore.pageCount,
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

export function remove(id: string): void {
	Dispatcher.dispatch({
		type: UserTypes.REMOVE,
		data: {
			id: id,
		},
	});

	Alert.info('User successfully removed');
}

EventDispatcher.register((action: UserTypes.UserDispatch) => {
	switch (action.type) {
		case UserTypes.CHANGE:
			sync();
			break;
	}
});
