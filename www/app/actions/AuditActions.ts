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
	if (!userId) {
		return Promise.resolve();
	}

	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/audit/' + userId)
			.query({
				page: AuditsStore.page,
				page_count: AuditsStore.pageCount,
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
					Alert.errorRes(res, 'Failed to load audits');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: AuditTypes.SYNC,
					data: {
						userId: userId,
						audits: res.body.audits,
						count: res.body.count,
					},
				});

				resolve();
			});
	});
}

export function reload(): Promise<void> {
	return load(AuditsStore.userId);
}

export function traverse(page: number): Promise<void> {
	Dispatcher.dispatch({
		type: AuditTypes.TRAVERSE,
		data: {
			page: page,
		},
	});

	return reload();
}

EventDispatcher.register((action: AuditTypes.AuditDispatch) => {
	switch (action.type) {
		case AuditTypes.CHANGE:
			reload();
			break;
	}
});
