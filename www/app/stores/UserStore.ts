/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as UserTypes from '../types/UserTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class UserStore extends EventEmitter {
	_user: UserTypes.UserRo;
	_token = Dispatcher.register((this._callback).bind(this));

	get user(): UserTypes.UserRo {
		return this._user;
	}

	get userM(): UserTypes.User {
		if (this._user) {
			return {
				...this._user,
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

	_load(user: UserTypes.User): void {
		this._user = Object.freeze(user);
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
