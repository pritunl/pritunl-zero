/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as KeybaseTypes from '../types/KeybaseTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class KeybaseStore extends EventEmitter {
	_info: KeybaseTypes.InfoRo;
	_token = Dispatcher.register((this._callback).bind(this));

	get info(): KeybaseTypes.InfoRo {
		return this._info;
	}

	get infoM(): KeybaseTypes.Info {
		if (this._info) {
			return {
				...this._info,
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

	_load(info: KeybaseTypes.Info): void {
		this._info = Object.freeze(info);
		this.emitChange();
	}

	_unload(): void {
		this._info = null;
		this.emitChange();
	}

	_callback(action: KeybaseTypes.InfoDispatch): void {
		switch (action.type) {
			case KeybaseTypes.LOAD:
				this._load(action.data.info);
				break;
			case KeybaseTypes.UNLOAD:
				this._unload();
				break;
		}
	}
}

export default new KeybaseStore();
