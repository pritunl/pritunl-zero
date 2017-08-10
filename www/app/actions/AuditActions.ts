/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import * as Csrf from '../Csrf';
import Loader from '../Loader';
import * as AuditTypes from '../types/AuditTypes';
import * as MiscUtils from '../utils/MiscUtils';
import AuditsStore from '../stores/AuditsStore';

let syncId: string;

export function load(userId: string): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/audit/' + userId)
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (curSyncId !== syncId) {
					resolve();
					return;
				}

				if (err) {
					Alert.errorRes(res, 'Failed to load audits');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: AuditTypes.SYNC,
					data: {
						userId: userId,
						audits: res.body,
					},
				});

				resolve();
			});
	});
}

export function reload(): Promise<void> {
	return load(AuditsStore.userId);
}

EventDispatcher.register((action: AuditTypes.AuditDispatch) => {
	switch (action.type) {
		case AuditTypes.CHANGE:
			reload();
			break;
	}
});
