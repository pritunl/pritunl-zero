/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as ServiceTypes from '../types/ServiceTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class ServicesStore extends EventEmitter {
	_services: ServiceTypes.ServicesRo = Object.freeze([]);
	_map: {[key: string]: number} = {};
	_token = Dispatcher.register((this._callback).bind(this));

	get services(): ServiceTypes.ServicesRo {
		return this._services;
	}

	get servicesM(): ServiceTypes.Services {
		let services: ServiceTypes.Services = [];
		this._services.forEach((service: ServiceTypes.ServiceRo): void => {
			services.push({
				...service,
			});
		});
		return services;
	}

	service(id: string): ServiceTypes.ServiceRo {
		let i = this._map[id];
		if (i === undefined) {
			return null;
		}
		return this._services[i];
	}

	emitChange(): void {
		this.emitDefer(GlobalTypes.CHANGE);
	}

	addChangeListener(callback: () => void): void {
		this.on(GlobalTypes.CHANGE, callback);
	}

	removeChangeListener(callback: () => void): void {
		this.removeListener(GlobalTypes.CHANGE, callback);
	}

	_sync(services: ServiceTypes.Service[]): void {
		this._map = {};
		for (let i = 0; i < services.length; i++) {
			services[i] = Object.freeze(services[i]);
			this._map[services[i].id] = i;
		}

		this._services = Object.freeze(services);
		this.emitChange();
	}

	_callback(action: ServiceTypes.ServiceDispatch): void {
		switch (action.type) {
			case ServiceTypes.SYNC:
				this._sync(action.data.services);
				break;
		}
	}
}

export default new ServicesStore();
