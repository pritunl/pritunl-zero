/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as DeviceTypes from '../types/DeviceTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class DevicesStore extends EventEmitter {
	_userId: string;
	_devices: DeviceTypes.DevicesRo = Object.freeze([]);
	_token = Dispatcher.register((this._callback).bind(this));

	get userId(): string {
		return this._userId;
	}

	get devices(): DeviceTypes.DevicesRo {
		return this._devices;
	}

	get devicesM(): DeviceTypes.Devices {
		let devices: DeviceTypes.Devices = [];
		this._devices.forEach((device: DeviceTypes.DeviceRo): void => {
			devices.push({
				...device,
			});
		});
		return devices;
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

	_sync(userId: string, devices: DeviceTypes.Device[]): void {
		this._userId = userId;

		for (let i = 0; i < devices.length; i++) {
			devices[i] = Object.freeze(devices[i]);
		}

		this._devices = Object.freeze(devices);
		this.emitChange();
	}

	_callback(action: DeviceTypes.DeviceDispatch): void {
		switch (action.type) {
			case DeviceTypes.SYNC:
				this._sync(action.data.userId, action.data.devices);
				break;
		}
	}
}

export default new DevicesStore();
