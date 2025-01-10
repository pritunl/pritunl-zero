/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as SecretTypes from '../types/SecretTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class SecretsStore extends EventEmitter {
	_secrets: SecretTypes.SecretsRo = Object.freeze([]);
	_map: {[key: string]: number} = {};
	_token = Dispatcher.register((this._callback).bind(this));

	get secrets(): SecretTypes.SecretsRo {
		return this._secrets;
	}

	get secretsM(): SecretTypes.Secrets {
		let secrets: SecretTypes.Secrets = [];
		this._secrets.forEach((
				secret: SecretTypes.SecretRo): void => {
			secrets.push({
				...secret,
			});
		});
		return secrets;
	}

	secret(id: string): SecretTypes.SecretRo {
		let i = this._map[id];
		if (i === undefined) {
			return null;
		}
		return this._secrets[i];
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

	_sync(secrets: SecretTypes.Secret[]): void {
		this._map = {};
		for (let i = 0; i < secrets.length; i++) {
			secrets[i] = Object.freeze(secrets[i]);
			this._map[secrets[i].id] = i;
		}

		this._secrets = Object.freeze(secrets);
		this.emitChange();
	}

	_callback(action: SecretTypes.SecretDispatch): void {
		switch (action.type) {
			case SecretTypes.SYNC:
				this._sync(action.data.secrets);
				break;
		}
	}
}

export default new SecretsStore();
