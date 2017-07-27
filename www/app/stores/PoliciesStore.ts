/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as PolicyTypes from '../types/PolicyTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class PoliciesStore extends EventEmitter {
	_policies: PolicyTypes.PoliciesRo = Object.freeze([]);
	_map: {[key: string]: number} = {};
	_token = Dispatcher.register((this._callback).bind(this));

	get policies(): PolicyTypes.PoliciesRo {
		return this._policies;
	}

	get policiesM(): PolicyTypes.Policies {
		let policies: PolicyTypes.Policies = [];
		this._policies.forEach((
				policy: PolicyTypes.PolicyRo): void => {
			policies.push({
				...policy,
			});
		});
		return policies;
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

	_sync(policies: PolicyTypes.Policy[]): void {
		this._map = {};
		for (let i = 0; i < policies.length; i++) {
			policies[i] = Object.freeze(policies[i]);
			this._map[policies[i].id] = i;
		}

		this._policies = Object.freeze(policies);
		this.emitChange();
	}

	_callback(action: PolicyTypes.PolicyDispatch): void {
		switch (action.type) {
			case PolicyTypes.SYNC:
				this._sync(action.data.policies);
				break;
		}
	}
}

export default new PoliciesStore();
