/// <reference path="../References.d.ts"/>
export const SYNC = 'settings.sync';
export const CHANGE = 'settings.change';

export interface Settings {
	elastic_address: string;
}

export type SettingsRo = Readonly<Settings>;

export interface SettingsDispatch {
	type: string;
	data?: Settings;
}
