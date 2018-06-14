/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import * as Csrf from '../Csrf';
import Loader from '../Loader';
import * as DeviceTypes from '../types/DeviceTypes';
import * as MiscUtils from '../utils/MiscUtils';
import DevicesStore from '../stores/DevicesStore';

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
			.get('/device/' + userId)
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
					Alert.errorRes(res, 'Failed to load devices');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: DeviceTypes.SYNC,
					data: {
						userId: userId,
						devices: res.body,
					},
				});

				resolve();
			});
	});
}

export function reload(): Promise<void> {
	return load(DevicesStore.userId);
}

export function commit(device: DeviceTypes.Device): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.put('/device/' + device.id)
			.send(device)
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
					Alert.errorRes(res, 'Failed to save device');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

export function remove(deviceId: string): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.delete('/device/' + deviceId)
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
					Alert.errorRes(res, 'Failed to delete device');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

EventDispatcher.register((action: DeviceTypes.DeviceDispatch) => {
	switch (action.type) {
		case DeviceTypes.CHANGE:
			reload();
			break;
	}
});
