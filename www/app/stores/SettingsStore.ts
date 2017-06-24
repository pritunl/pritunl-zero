/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import * as Events from 'events';
import * as SettingsTypes from '../types/SettingsTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class SettingsStore extends Events.EventEmitter {
	_settings: SettingsTypes.Settings = {
		elastic_address: '',
	};
	_token = Dispatcher.register((this._callback).bind(this));

	get settings(): SettingsTypes.Settings {
		return {
			...this._settings,
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

	_sync(settings: SettingsTypes.Settings): void {
		this._settings = settings;
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
