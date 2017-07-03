/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as SubscriptionTypes from '../types/SubscriptionTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class SubscriptionStore extends EventEmitter {
	_subscription: SubscriptionTypes.SubscriptionRo;
	_token = Dispatcher.register((this._callback).bind(this));

	get subscription(): SubscriptionTypes.SubscriptionRo {
		return this._subscription;
	}

	get subscriptionM(): SubscriptionTypes.Subscription {
		if (this._subscription) {
			return {
				...this._subscription,
			};
		}
		return undefined;
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

	_sync(subscription: SubscriptionTypes.Subscription): void {
		this._subscription = Object.freeze(subscription);
		this.emitChange();
	}

	_callback(action: SubscriptionTypes.SubscriptionDispatch): void {
		switch (action.type) {
			case SubscriptionTypes.SYNC:
				this._sync(action.data);
				break;
		}
	}
}

export default new SubscriptionStore();
