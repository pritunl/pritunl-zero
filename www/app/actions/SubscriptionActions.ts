/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import Loader from '../Loader';
import * as SubscriptionTypes from '../types/SubscriptionTypes';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;

export function sync(): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/subscription')
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (curSyncId !== syncId) {
					resolve();
					return;
				}

				if (err) {
					Alert.errorRes(res, 'Failed to sync subscription');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: SubscriptionTypes.SYNC,
					data: res.body,
				});

				resolve();
			});
	});
}

EventDispatcher.register((action: SubscriptionTypes.SubscriptionDispatch) => {
	switch (action.type) {
		case SubscriptionTypes.CHANGE:
			sync();
			break;
	}
});
