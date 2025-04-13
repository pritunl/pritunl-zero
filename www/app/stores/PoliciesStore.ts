/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as PolicyTypes from '../types/PolicyTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class PoliciesStore extends EventEmitter {
	_policies: PolicyTypes.PoliciesRo = Object.freeze([]);
	_page: number;
	_pageCount: number;
	_filter: PolicyTypes.Filter = null;
	_count: number;
	_map: {[key: string]: number} = {};
	_token = Dispatcher.register((this._callback).bind(this));

	get policies(): PolicyTypes.PoliciesRo {
		return this._policies;
	}

	get policiesM(): PolicyTypes.Policies {
		let policies: PolicyTypes.Policies = [];
		this._policies.forEach((policy: PolicyTypes.PolicyRo): void => {
			policies.push({
				...policy,
			});
		});
		return policies;
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

	get filter(): PolicyTypes.Filter {
		return this._filter;
	}

	get count(): number {
		return this._count || 0;
	}

	policy(id: string): PolicyTypes.PolicyRo {
		let i = this._map[id];
		if (i === undefined) {
			return null;
		}
		return this._policies[i];
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

	_filterCallback(filter: PolicyTypes.Filter): void {
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

	_sync(policies: PolicyTypes.Policy[], count: number): void {
		this._map = {};
		for (let i = 0; i < policies.length; i++) {
			policies[i] = Object.freeze(policies[i]);
			this._map[policies[i].id] = i;
		}

		this._count = count;
		this._policies = Object.freeze(policies);
		this._page = Math.min(this.pages, this.page);

		this.emitChange();
	}

	_callback(action: PolicyTypes.PolicyDispatch): void {
		switch (action.type) {
			case PolicyTypes.TRAVERSE:
				this._traverse(action.data.page);
				break;

			case PolicyTypes.FILTER:
				this._filterCallback(action.data.filter);
				break;

			case PolicyTypes.SYNC:
				this._sync(action.data.policies, action.data.count);
				break;
		}
	}
}

export default new PoliciesStore();
