/// <reference path="../References.d.ts"/>
export const SYNC = 'device.sync';
export const CHANGE = 'device.change';

export interface Device {
	id: string;
	user?: string;
	name?: string;
	type?: string;
	timestamp?: string;
	disabled?: boolean;
	active_until?: string;
	last_active?: string;
}

export type Devices = Device[];

export type DeviceRo = Readonly<Device>;
export type DevicesRo = ReadonlyArray<DeviceRo>;

export interface DeviceDispatch {
	type: string;
	data?: {
		id?: string;
		device?: Device;
		devices?: Devices;
	};
}
