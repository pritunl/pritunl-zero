/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import * as Csrf from '../Csrf';
import Loader from '../Loader';
import * as SshcertificateTypes from '../types/SshcertificateTypes';
import * as MiscUtils from '../utils/MiscUtils';
import SshcertificatesStore from '../stores/SshcertificatesStore';

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
			.get('/sshcertificate/' + userId)
			.query({
				page: SshcertificatesStore.page,
				page_count: SshcertificatesStore.pageCount,
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
					Alert.errorRes(res, 'Failed to load SSH certificates');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: SshcertificateTypes.SYNC,
					data: {
						userId: userId,
						certificates: res.body.certificates,
						count: res.body.count,
					},
				});

				resolve();
			});
	});
}

export function reload(): Promise<void> {
	return load(SshcertificatesStore.userId);
}

export function traverse(page: number): Promise<void> {
	Dispatcher.dispatch({
		type: SshcertificateTypes.TRAVERSE,
		data: {
			page: page,
		},
	});

	return reload();
}

EventDispatcher.register((
		action: SshcertificateTypes.SshcertificateDispatch) => {
	switch (action.type) {
		case SshcertificateTypes.CHANGE:
			reload();
			break;
	}
});
