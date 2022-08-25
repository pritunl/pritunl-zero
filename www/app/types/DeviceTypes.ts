/// <reference path="../References.d.ts"/>
export const SYNC = 'device.sync';
export const CHANGE = 'device.change';

export interface Device {
	id: string;
	user?: string;
	name?: string;
	type?: string;
	mode?: string;
	timestamp?: string;
	disabled?: boolean;
	active_until?: string;
	number?: string;
	last_active?: string;
	wan_rp_id?: string;
	ssh_public_key?: string;
}

export type Devices = Device[];

export type DeviceRo = Readonly<Device>;
export type DevicesRo = ReadonlyArray<DeviceRo>;

export interface DeviceDispatch {
	type: string;
	data?: {
		id?: string;
		userId?: string;
		device?: Device;
		devices?: Devices;
		showRemoved?: boolean;
	};
}
