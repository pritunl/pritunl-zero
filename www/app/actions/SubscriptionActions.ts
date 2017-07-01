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

export function activate(license: string): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.post('/subscription')
			.send({
				license: license,
			})
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					Alert.errorRes(res, 'Failed to activate subscription');
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

export function checkout(plan: string, card: string,
		email: string): Promise<string> {
	let loader = new Loader().loading();

	return new Promise<string>((resolve, reject): void => {
		SuperAgent
			.post('https://app-test.pritunl.net/subscription')
			.send({
				plan: plan,
				card: card,
				email: email,
			})
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					Alert.errorRes(res, 'Failed to checkout');
					reject(err);
					return;
				}

				resolve(res.body.msg);
			});
	});
}

export function payment(key: string, plan: string, card: string,
		email: string): Promise<string> {
	let loader = new Loader().loading();

	return new Promise<string>((resolve, reject): void => {
		SuperAgent
			.put('https://app-test.pritunl.net/subscription')
			.send({
				key: key,
				plan: plan,
				card: card,
				email: email,
			})
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					Alert.errorRes(res, 'Failed to update payment');
					reject(err);
					return;
				}

				resolve(res.body.msg);
			});
	});
}

export function cancel(key: string): Promise<string> {
	let loader = new Loader().loading();

	return new Promise<string>((resolve, reject): void => {
		SuperAgent
			.delete('https://app-test.pritunl.net/subscription')
			.send({
				key: key,
			})
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					Alert.errorRes(res, 'Failed to cancel subscription');
					reject(err);
					return;
				}

				resolve(res.body.msg);
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
