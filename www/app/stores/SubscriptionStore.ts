/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import * as Events from 'events';
import * as SubscriptionTypes from '../types/SubscriptionTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class SubscriptionStore extends Events.EventEmitter {
	_subscription: SubscriptionTypes.Subscription = {};
	_token = Dispatcher.register((this._callback).bind(this));

	get subscription(): SubscriptionTypes.Subscription {
		return {
			...this._subscription,
		};
	}

	emitChange(): void {
		this.emit(GlobalTypes.CHANGE);
	}

	addChangeListener(callback: () => void): void {
		this.on(GlobalTypes.CHANGE, callback);
	}

	removeChangeListener(callback: () => void): void {
		this.removeListener(GlobalTypes.CHANGE, callback);
	}

	_sync(settings: SubscriptionTypes.Subscription): void {
		this._subscription = settings;
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
