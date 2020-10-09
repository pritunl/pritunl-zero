/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import * as Csrf from '../Csrf';
import Loader from '../Loader';
import ServicesStore from '../stores/ServicesStore';
import * as ServiceTypes from '../types/ServiceTypes';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;
let nameSyncId: string;

export function sync(): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/service')
			.query({
				...ServicesStore.filter,
				page: ServicesStore.page,
				page_count: ServicesStore.pageCount,
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
					Alert.errorRes(res, 'Failed to load services');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: ServiceTypes.SYNC,
					data: {
						services: res.body.services,
						count: res.body.count,
					},
				});

				resolve();
			});
	});
}

export function traverse(page: number): Promise<void> {
	Dispatcher.dispatch({
		type: ServiceTypes.TRAVERSE,
		data: {
			page: page,
		},
	});

	return sync();
}

export function filter(filt: ServiceTypes.Filter): Promise<void> {
	Dispatcher.dispatch({
		type: ServiceTypes.FILTER,
		data: {
			filter: filt,
		},
	});

	return sync();
}

export function commit(service: ServiceTypes.Service): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.put('/service/' + service.id)
			.send(service)
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
					Alert.errorRes(res, 'Failed to save service');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

export function create(service: ServiceTypes.Service): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.post('/service')
			.send(service)
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
					Alert.errorRes(res, 'Failed to create service');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

export function remove(serviceId: string): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.delete('/service/' + serviceId)
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					Alert.errorRes(res, 'Failed to delete services');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

export function removeMulti(serviceIds: string[]): Promise<void> {
	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.delete('/service')
			.send(serviceIds)
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
					Alert.errorRes(res, 'Failed to delete services');
					reject(err);
					return;
				}

				resolve();
			});
	});
}

export function syncNames(): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	nameSyncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/service')
			.query({
				service_names: "true",
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

				if (curSyncId !== nameSyncId) {
					resolve();
					return;
				}

				if (err) {
					Alert.errorRes(res, 'Failed to load service names');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: ServiceTypes.SYNC_NAMES,
					data: {
						services: res.body,
					},
				});

				resolve();
			});
	});
}

EventDispatcher.register((action: ServiceTypes.ServiceDispatch) => {
	switch (action.type) {
		case ServiceTypes.CHANGE:
			sync();
			break;
	}
});
