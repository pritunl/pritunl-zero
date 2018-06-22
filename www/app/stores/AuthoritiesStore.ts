/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as AuthorityTypes from '../types/AuthorityTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class AuthoritiesStore extends EventEmitter {
	_authorities: AuthorityTypes.AuthoritiesRo = Object.freeze([]);
	_secrets: {[key: string]: string} = {};
	_map: {[key: string]: number} = {};
	_token = Dispatcher.register((this._callback).bind(this));

	get authorities(): AuthorityTypes.AuthoritiesRo {
		return this._authorities;
	}

	get authoritiesM(): AuthorityTypes.Authorities {
		let authorities: AuthorityTypes.Authorities = [];
		this._authorities.forEach((
				policy: AuthorityTypes.AuthorityRo): void => {
			authorities.push({
				...policy,
			});
		});
		return authorities;
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

	_sync(authorities: AuthorityTypes.Authority[]): void {
		this._map = {};
		for (let i = 0; i < authorities.length; i++) {
			authorities[i] = Object.freeze(authorities[i]);
			this._map[authorities[i].id] = i;
		}

		this._authorities = Object.freeze(authorities);
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
			case AuthorityTypes.SYNC:
				this._sync(action.data.authorities);
				break;
			case AuthorityTypes.SYNC_SECRET:
				this._syncSecret(action.data.id, action.data.secret);
				break;
		}
	}
}

export default new AuthoritiesStore();
