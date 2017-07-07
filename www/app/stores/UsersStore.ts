/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as UserTypes from '../types/UserTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class UsersStore extends EventEmitter {
	_users: UserTypes.UsersRo = Object.freeze([]);
	_page: number;
	_pageCount: number;
	_filter: UserTypes.Filter = null;
	_count: number;
	_token = Dispatcher.register((this._callback).bind(this));

	get users(): UserTypes.UsersRo {
		return this._users;
	}

	get usersM(): UserTypes.Users {
		let users: UserTypes.Users = [];
		this._users.forEach((user: UserTypes.UserRo): void => {
			users.push({
				...user,
			});
		});
		return users;
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

	get filter(): UserTypes.Filter {
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
		this._page = page;
	}

	_filterCallback(filter: UserTypes.Filter): void {
		this._filter = filter;
		this.emitChange();
	}

	_sync(users: UserTypes.User[], count: number): void {
		for (let i = 0; i < users.length; i++) {
			users[i] = Object.freeze(users[i]);
		}

		this._count = count;
		this._users = Object.freeze(users);
		this._page = Math.min(this.pages, this.page);

		this.emitChange();
	}

	_callback(action: UserTypes.UserDispatch): void {
		switch (action.type) {
			case UserTypes.TRAVERSE:
				this._traverse(action.data.page);
				break;

			case UserTypes.FILTER:
				this._filterCallback(action.data.filter);
				break;

			case UserTypes.SYNC:
				this._sync(action.data.users, action.data.count);
				break;
		}
	}
}

export default new UsersStore();
