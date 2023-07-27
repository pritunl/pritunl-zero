/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as AlertTypes from '../types/AlertTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class AlertsStore extends EventEmitter {
	_alerts: AlertTypes.AlertsRo = Object.freeze([]);
	_page: number;
	_pageCount: number;
	_filter: AlertTypes.Filter = null;
	_count: number;
	_map: {[key: string]: number} = {};
	_token = Dispatcher.register((this._callback).bind(this));

	get alerts(): AlertTypes.AlertsRo {
		return this._alerts;
	}

	get alertsM(): AlertTypes.Alerts {
		let alerts: AlertTypes.Alerts = [];
		this._alerts.forEach((alert: AlertTypes.AlertRo): void => {
			alerts.push({
				...alert,
			});
		});
		return alerts;
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

	get filter(): AlertTypes.Filter {
		return this._filter;
	}

	get count(): number {
		return this._count || 0;
	}

	alert(id: string): AlertTypes.AlertRo {
		let i = this._map[id];
		if (i === undefined) {
			return null;
		}
		return this._alerts[i];
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

	_filterCallback(filter: AlertTypes.Filter): void {
		if ((this._filter !== null && filter === null) ||
			(!Object.keys(this._filter).length && filter !== null) || (
				filter && this._filter && (
					filter.name !== this._filter.name
				))) {
			this._traverse(0);
		}
		this._filter = filter;
		this.emitChange();
	}

	_sync(alerts: AlertTypes.Alert[], count: number): void {
		this._map = {};
		for (let i = 0; i < alerts.length; i++) {
			alerts[i] = Object.freeze(alerts[i]);
			this._map[alerts[i].id] = i;
		}

		this._count = count;
		this._alerts = Object.freeze(alerts);
		this._page = Math.min(this.pages, this.page);

		this.emitChange();
	}

	_callback(action: AlertTypes.AlertDispatch): void {
		switch (action.type) {
			case AlertTypes.TRAVERSE:
				this._traverse(action.data.page);
				break;

			case AlertTypes.FILTER:
				this._filterCallback(action.data.filter);
				break;

			case AlertTypes.SYNC:
				this._sync(action.data.alerts, action.data.count);
				break;
		}
	}
}

export default new AlertsStore();
