/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import * as Events from 'events';
import * as UserTypes from '../types/UserTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class UserStore extends Events.EventEmitter {
	_user: UserTypes.User;
	_token = Dispatcher.register((this._callback).bind(this));

	get user(): UserTypes.User {
		return this._user;
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

	_load(user: UserTypes.User): void {
		this._user = user;
		this.emitChange();
	}

	_unload(): void {
		this._user = null;
		this.emitChange();
	}

	_callback(action: UserTypes.UserDispatch): void {
		switch (action.type) {
			case UserTypes.LOAD:
				this._load(action.data.user);
				break;
			case UserTypes.UNLOAD:
				this._unload();
				break;
		}
	}
}

export default new UserStore();
