/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as ServiceTypes from '../types/ServiceTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class ServicesStore extends EventEmitter {
	_services: ServiceTypes.ServicesRo = Object.freeze([]);
	_services_name: ServiceTypes.ServicesRo = Object.freeze([]);
	_page: number;
	_pageCount: number;
	_filter: ServiceTypes.Filter = null;
	_count: number;
	_map: {[key: string]: number} = {};
	_map_name: {[key: string]: number} = {};
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

	get servicesName(): ServiceTypes.ServicesRo {
		return this._services_name;
	}

	get servicesNameM(): ServiceTypes.Services {
		let services: ServiceTypes.Services = [];
		this._services_name.forEach((service: ServiceTypes.ServiceRo): void => {
			services.push({
				...service,
			});
		});
		return services;
	}

	get page(): number {
		return this._page || 0;
	}

	get pageCount(): number {
		return this._pageCount || 20;
	}

	get pages(): number {
		return Math.ceil(this.count / this.pageCount);
	}

	get filter(): ServiceTypes.Filter {
		return this._filter;
	}

	get count(): number {
		return this._count || 0;
	}

	service(id: string): ServiceTypes.ServiceRo {
		let i = this._map[id];
		if (i === undefined) {
			return null;
		}
		return this._services[i];
	}

	serviceName(id: string): ServiceTypes.ServiceRo {
		let i = this._map_name[id];
		if (i === undefined) {
			return null;
		}
		return this._services_name[i];
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

	_traverse(page: number): void {
		this._page = Math.min(this.pages, page);
	}

	_filterCallback(filter: ServiceTypes.Filter): void {
		if ((this._filter !== null && filter === null) ||
			(!Object.keys(this._filter || {}).length && filter !== null) || (
				filter && this._filter && (
					filter.name !== this._filter.name
				))) {
			this._traverse(0);
		}
		this._filter = filter;
		this.emitChange();
	}

	_sync(services: ServiceTypes.Service[], count: number): void {
		this._map = {};
		for (let i = 0; i < services.length; i++) {
			services[i] = Object.freeze(services[i]);
			this._map[services[i].id] = i;
		}

		this._count = count;
		this._services = Object.freeze(services);
		this._page = Math.min(this.pages, this.page);

		this.emitChange();
	}

	_sync_names(services: ServiceTypes.Service[]): void {
		this._map_name = {};
		for (let i = 0; i < services.length; i++) {
			services[i] = Object.freeze(services[i]);
			this._map_name[services[i].id] = i;
		}

		this._services_name = Object.freeze(services);
		this.emitChange();
	}

	_callback(action: ServiceTypes.ServiceDispatch): void {
		switch (action.type) {
			case ServiceTypes.TRAVERSE:
				this._traverse(action.data.page);
				break;

			case ServiceTypes.FILTER:
				this._filterCallback(action.data.filter);
				break;

			case ServiceTypes.SYNC:
				this._sync(action.data.services, action.data.count);
				break;

			case ServiceTypes.SYNC_NAMES:
				this._sync_names(action.data.services);
				break;
		}
	}
}

export default new ServicesStore();
