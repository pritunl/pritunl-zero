/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as AuthorityTypes from '../types/AuthorityTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class AuthoritiesStore extends EventEmitter {
	_authorities: AuthorityTypes.AuthoritiesRo = Object.freeze([]);
	_secrets: {[key: string]: string} = {};
	_page: number;
	_pageCount: number;
	_filter: AuthorityTypes.Filter = null;
	_count: number;
	_map: {[key: string]: number} = {};
	_token = Dispatcher.register((this._callback).bind(this));

	get authorities(): AuthorityTypes.AuthoritiesRo {
		return this._authorities;
	}

	get authoritiesM(): AuthorityTypes.Authorities {
		let authorities: AuthorityTypes.Authorities = [];
		this._authorities.forEach((authority: AuthorityTypes.AuthorityRo): void => {
			authorities.push({
				...authority,
			});
		});
		return authorities;
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

	get filter(): AuthorityTypes.Filter {
		return this._filter;
	}

	get count(): number {
		return this._count || 0;
	}

	authority(id: string): AuthorityTypes.AuthorityRo {
		let i = this._map[id];
		if (i === undefined) {
			return null;
		}
		return this._authorities[i];
	}

	authoritySecret(id: string): string {
		return this._secrets[id];
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

	_filterCallback(filter: AuthorityTypes.Filter): void {
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

	_sync(authorities: AuthorityTypes.Authority[], count: number): void {
		this._map = {};
		for (let i = 0; i < authorities.length; i++) {
			authorities[i] = Object.freeze(authorities[i]);
			this._map[authorities[i].id] = i;
		}

		this._count = count;
		this._authorities = Object.freeze(authorities);
		this._page = Math.min(this.pages, this.page);

		this.emitChange();
	}

	_syncSecret(id: string, secret: string): void {
		if (!secret) {
			delete this._secrets[id];
		} else {
			this._secrets[id] = secret;
		}
		this.emitChange();
	}

	_callback(action: AuthorityTypes.AuthorityDispatch): void {
		switch (action.type) {
			case AuthorityTypes.TRAVERSE:
				this._traverse(action.data.page);
				break;

			case AuthorityTypes.FILTER:
				this._filterCallback(action.data.filter);
				break;

			case AuthorityTypes.SYNC:
				this._sync(action.data.authorities, action.data.count);
				break;

			case AuthorityTypes.SYNC_SECRET:
				this._syncSecret(action.data.id, action.data.secret);
				break;
		}
	}
}

export default new AuthoritiesStore();
