/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as SettingsTypes from '../types/SettingsTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class SettingsStore extends EventEmitter {
	_settings: SettingsTypes.SettingsRo;
	_token = Dispatcher.register((this._callback).bind(this));

	get settings(): SettingsTypes.SettingsRo {
		return this._settings;
	}

	get settingsM(): SettingsTypes.Settings {
		if (this._settings) {
			return {
				...this._settings,
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

	_sync(settings: SettingsTypes.Settings): void {
		this._settings = Object.freeze(settings);
		this.emitChange();
	}

	_callback(action: SettingsTypes.SettingsDispatch): void {
		switch (action.type) {
			case SettingsTypes.SYNC:
				this._sync(action.data);
				break;
		}
	}
}

export default new SettingsStore();
