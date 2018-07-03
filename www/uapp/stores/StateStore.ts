/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as StateTypes from '../types/StateTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class StateStore extends EventEmitter {
	_sshToken: string;
	_sshDevice: string;
	_token = Dispatcher.register((this._callback).bind(this));

	get sshToken(): string {
		return this._sshToken;
	}

	get sshDevice(): string {
		return this._sshDevice;
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

	_callback(action: StateTypes.StateDispatch): void {
		switch (action.type) {
			case StateTypes.SSH_TOKEN:
				this._sshToken = action.data;
				this.emitChange();
				break;
			case StateTypes.SSH_DEVICE:
				this._sshDevice = action.data;
				this.emitChange();
				break;
		}
	}
}

export default new StateStore();
