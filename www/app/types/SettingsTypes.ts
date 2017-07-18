/// <reference path="../References.d.ts"/>
export const SYNC = 'settings.sync';
export const CHANGE = 'settings.change';

export interface Provider {
	id: string;
	type: string;
	label: string;
	default_roles: string[];
	auto_create: boolean;
}

export interface GoogleProvider extends Provider {
	domain?: string;
}

export interface SamlProvider extends Provider {
	issuer_url?: string;
	saml_url?: string;
	saml_cert?: string;
}

export type ProviderAny = Provider & GoogleProvider & SamlProvider;
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
