/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as EndpointTypes from '../types/EndpointTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class EndpointsStore extends EventEmitter {
	_endpoints: EndpointTypes.EndpointsRo = Object.freeze([]);
	_page: number;
	_pageCount: number;
	_filter: EndpointTypes.Filter = null;
	_count: number;
	_map: {[key: string]: number} = {};
	_token = Dispatcher.register((this._callback).bind(this));

	get endpoints(): EndpointTypes.EndpointsRo {
		return this._endpoints;
	}

	get endpointsM(): EndpointTypes.Endpoints {
		let endpoints: EndpointTypes.Endpoints = [];
		this._endpoints.forEach((endpoint: EndpointTypes.EndpointRo): void => {
			endpoints.push({
				...endpoint,
			});
		});
		return endpoints;
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

	get filter(): EndpointTypes.Filter {
		return this._filter;
	}

	get count(): number {
		return this._count || 0;
	}

	endpoint(id: string): EndpointTypes.EndpointRo {
		let i = this._map[id];
		if (i === undefined) {
			return null;
		}
		return this._endpoints[i];
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

	_filterCallback(filter: EndpointTypes.Filter): void {
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

	_sync(endpoints: EndpointTypes.Endpoint[], count: number): void {
		this._map = {};
		for (let i = 0; i < endpoints.length; i++) {
			endpoints[i] = Object.freeze(endpoints[i]);
			this._map[endpoints[i].id] = i;
		}

		this._count = count;
		this._endpoints = Object.freeze(endpoints);
		this._page = Math.min(this.pages, this.page);

		this.emitChange();
	}

	_callback(action: EndpointTypes.EndpointDispatch): void {
		switch (action.type) {
			case EndpointTypes.TRAVERSE:
				this._traverse(action.data.page);
				break;

			case EndpointTypes.FILTER:
				this._filterCallback(action.data.filter);
				break;

			case EndpointTypes.SYNC:
				this._sync(action.data.endpoints, action.data.count);
				break;
		}
	}
}

export default new EndpointsStore();
