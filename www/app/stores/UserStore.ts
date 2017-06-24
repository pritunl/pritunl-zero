/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import * as Events from 'events';
import * as UserTypes from '../types/UserTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class UserStore extends Events.EventEmitter {
	_users: UserTypes.Users = [];
	_page: number;
	_pageCount: number;
	_count: number;
	_map: {[key: string]: number} = {};
	_token = Dispatcher.register((this._callback).bind(this));

	get users(): UserTypes.Users {
		return this._users;
	}

	get page(): number {
		return this._page || 0;
	}

	get pageCount(): number {
		return this._pageCount || 50;
	}

	get count(): number {
		return this._count || 0;
	}

	build(id: string): UserTypes.User {
		let i = this._map[id];
		if (i === undefined) {
			return null;
		}
		return this._users[i];
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

	_update(data: UserTypes.User): void {
		let n = this._map[data.id];
		if (n === undefined) {
			return;
		}
		this._users[n] = data;
		this.emitChange();
	}

	_traverse(page: number): void {
		this._page = page;
	}

	_sync(users: UserTypes.Users, count: number): void {
		this._count = count;

		this._map = {};
		for (let i = 0; i < users.length; i++) {
			this._map[users[i].id] = i;
		}
		this._users = users;

		this.emitChange();
	}

	_remove(id: string): void {
		let n = this._map[id];
		if (n === undefined) {
			return;
		}
		delete this._map[id];

		this._users.splice(n, 1);

		for (let i = n; i < this._users.length; i++) {
			this._map[this._users[i].id] = i;
		}

		this.emitChange();
	}

	_callback(action: UserTypes.UserDispatch): void {
		switch (action.type) {
			case UserTypes.UPDATE:
				this._update(action.data.user);
				break;

			case UserTypes.TRAVERSE:
				this._traverse(action.data.page);
				break;

			case UserTypes.SYNC:
				this._sync(action.data.users, action.data.count);
				break;

			case UserTypes.REMOVE:
				this._remove(action.data.id);
				break;
		}
	}
}

export default new UserStore();
