/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import * as Csrf from '../Csrf';
import Loader from '../Loader';
import * as SessionTypes from '../types/SessionTypes';
import * as MiscUtils from '../utils/MiscUtils';
import SessionsStore from '../stores/SessionsStore';

let syncId: string;

export function _load(userId: string): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/session/' + userId)
			.query({
				show_removed: SessionsStore.showRemoved,
			})
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (curSyncId !== syncId) {
					resolve();
					return;
				}

				if (err) {
					Alert.errorRes(res, 'Failed to load sessions');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: SessionTypes.SYNC,
					data: {
						userId: userId,
						sessions: res.body,
					},
				});

				resolve();
			});
	});
}

export function load(userId: string): Promise<void> {
	Dispatcher.dispatch({
		type: SessionTypes.SHOW_REMOVED,
		data: {
			showRemoved: false,
		},
	});

	return _load(userId);
}

export function reload(): Promise<void> {
	return _load(SessionsStore.userId);
}

export function showRemoved(state: boolean): Promise<void> {
	Dispatcher.dispatch({
		type: SessionTypes.SHOW_REMOVED,
		data: {
			showRemoved: state,
		},
	});

	return reload();
}

export function remove(sessionId: string): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.delete('/session/' + sessionId)
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					Alert.errorRes(res, 'Failed to delete session');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

EventDispatcher.register((action: SessionTypes.SessionDispatch) => {
	switch (action.type) {
		case SessionTypes.CHANGE:
			reload();
			break;
	}
});
