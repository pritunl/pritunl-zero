/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as LogTypes from '../types/LogTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class LogsStore extends EventEmitter {
	_logs: LogTypes.LogsRo = Object.freeze([]);
	_page: number;
	_pageCount: number;
	_filter: LogTypes.Filter = null;
	_count: number;
	_token = Dispatcher.register((this._callback).bind(this));

	get logs(): LogTypes.LogsRo {
		return this._logs;
	}

	get logsM(): LogTypes.Logs {
		let logs: LogTypes.Logs = [];
		this._logs.forEach((log: LogTypes.LogRo): void => {
			logs.push({
				...log,
			});
		});
		return logs;
	}

	get page(): number {
		return this._page || 0;
	}

	get pageCount(): number {
		return this._pageCount || 50;
	}

	get pages(): number {
		return Math.ceil(this.count / this.pageCount);
	}

	get filter(): LogTypes.Filter {
		return this._filter;
	}

	get count(): number {
		return this._count || 0;
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

	_filterCallback(filter: LogTypes.Filter): void {
		if ((this._filter !== null && filter === null) ||
			(!Object.keys(this._filter || {}).length && filter !== null) || (
				filter && this._filter && (
					filter.level !== this._filter.level
				))) {
			this._traverse(0);
		}
		this._filter = filter;
		this.emitChange();
	}

	_sync(logs: LogTypes.Log[], count: number): void {
		for (let i = 0; i < logs.length; i++) {
			logs[i] = Object.freeze(logs[i]);
		}

		this._count = count;
		this._logs = Object.freeze(logs);
		this._page = Math.min(this.pages, this.page);

		this.emitChange();
	}

	_callback(action: LogTypes.LogDispatch): void {
		switch (action.type) {
			case LogTypes.TRAVERSE:
				this._traverse(action.data.page);
				break;

			case LogTypes.FILTER:
				this._filterCallback(action.data.filter);
				break;

			case LogTypes.SYNC:
				this._sync(action.data.logs, action.data.count);
				break;
		}
	}
}

export default new LogsStore();
