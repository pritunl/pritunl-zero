/// <reference path="../References.d.ts"/>
export const SYNC = 'settings.sync';
export const CHANGE = 'settings.change';

export interface Provider {
	type: string;
	label: string;
	default_roles: string[];
}

export interface GoogleProvider extends Provider {
	domain?: string;
}

export type ProviderAny = Provider & GoogleProvider;
export type Providers = ProviderAny[];

export interface Settings {
	auth_providers: Providers;
	elastic_address: string;
}

export type SettingsRo = Readonly<Settings>;

export interface SettingsDispatch {
	type: string;
	data?: Settings;
}
