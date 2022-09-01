/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as CheckTypes from '../types/CheckTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class ChecksStore extends EventEmitter {
	_checks: CheckTypes.ChecksRo = Object.freeze([]);
	_page: number;
	_pageCount: number;
	_filter: CheckTypes.Filter = null;
	_count: number;
	_map: {[key: string]: number} = {};
	_token = Dispatcher.register((this._callback).bind(this));

	get checks(): CheckTypes.ChecksRo {
		return this._checks;
	}

	get checksM(): CheckTypes.Checks {
		let checks: CheckTypes.Checks = [];
		this._checks.forEach((check: CheckTypes.CheckRo): void => {
			checks.push({
				...check,
			});
		});
		return checks;
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

	get filter(): CheckTypes.Filter {
		return this._filter;
	}

	get count(): number {
		return this._count || 0;
	}

	check(id: string): CheckTypes.CheckRo {
		let i = this._map[id];
		if (i === undefined) {
			return null;
		}
		return this._checks[i];
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

	_filterCallback(filter: CheckTypes.Filter): void {
		if ((this._filter !== null && filter === null) ||
			(this._filter === {} && filter !== null) || (
				filter && this._filter && (
					filter.name !== this._filter.name
				))) {
			this._traverse(0);
		}
		this._filter = filter;
		this.emitChange();
	}

	_sync(checks: CheckTypes.Check[], count: number): void {
		this._map = {};
		for (let i = 0; i < checks.length; i++) {
			checks[i] = Object.freeze(checks[i]);
			this._map[checks[i].id] = i;
		}

		this._count = count;
		this._checks = Object.freeze(checks);
		this._page = Math.min(this.pages, this.page);

		this.emitChange();
	}

	_callback(action: CheckTypes.CheckDispatch): void {
		switch (action.type) {
			case CheckTypes.TRAVERSE:
				this._traverse(action.data.page);
				break;

			case CheckTypes.FILTER:
				this._filterCallback(action.data.filter);
				break;

			case CheckTypes.SYNC:
				this._sync(action.data.checks, action.data.count);
				break;
		}
	}
}

export default new ChecksStore();
