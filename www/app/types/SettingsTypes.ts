/// <reference path="../References.d.ts"/>
export const SYNC = 'settings.sync';
export const CHANGE = 'settings.change';

export interface Provider {
	id?: string;
	type?: string;
	label?: string;
	default_roles?: string[];
	auto_create?: boolean;
	role_management?: string;
}

export interface AzureProvider extends Provider {
	tenant?: string;
	client_id?: string;
	client_secret?: string;
}

export interface GoogleProvider extends Provider {
	domain?: string;
	google_key?: string;
	google_email?: string;
}

export interface SamlProvider extends Provider {
	issuer_url?: string;
	saml_url?: string;
	saml_cert?: string;
}

export type ProviderAny = Provider & AzureProvider & GoogleProvider &
	SamlProvider;
export type Providers = ProviderAny[];

export interface Settings {
	auth_providers: Providers;
	auth_admin_expire: number;
	auth_admin_max_duration: number;
	auth_proxy_expire: number;
	auth_proxy_max_duration: number;
	auth_user_expire: number;
	auth_user_max_duration: number;
	elastic_address: string;
	elastic_proxy_requests: boolean;
}

export type SettingsRo = Readonly<Settings>;

export interface SettingsDispatch {
	type: string;
	data?: Settings;
}
