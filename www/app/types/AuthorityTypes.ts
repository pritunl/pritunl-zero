/// <reference path="../References.d.ts"/>
export const SYNC = 'authority.sync';
export const SYNC_NAMES = 'authority.sync_names';
export const CHANGE = 'authority.change';
export const TRAVERSE = 'authority.traverse';
export const FILTER = 'authority.filter';
export const SYNC_SECRET = 'authority.sync_secret';

export interface Info {
	key_alg?: string;
}

export interface Authority {
	id: string;
	name?: string;
	type?: string;
	algorithm?: string;
	info?: Info;
	expire?: number;
	host_expire?: number;
	key_id_format?: string;
	match_roles?: boolean;
	roles?: string[];
	public_key?: string;
	public_key_pem?: string;
	root_certificate?: string;
	proxy_jump?: string;
	proxy_hosting?: boolean;
	proxy_hostname?: string;
	proxy_port?: number;
	host_domain?: string;
	host_subnets?: string[];
	host_matches?: string[];
	host_proxy?: string;
	host_certificates?: boolean;
	strict_host_checking?: boolean;
	host_tokens?: string[];
	hsm_status?: string;
	hsm_timestamp?: string;
	hsm_token?: string;
	hsm_secret?: string;
	hsm_serial?: string;
	hsm_generate_secret?: boolean;
	reset_proxy_host_key?: boolean;
}

export interface Filter {
	id?: string;
	name?: string;
}

export type Authorities = Authority[];

export type AuthorityRo = Readonly<Authority>;
export type AuthoritiesRo = ReadonlyArray<AuthorityRo>;

export interface AuthorityDispatch {
	type: string;
	data?: {
		id?: string;
		secret?: string;
		authority?: Authority;
		authorities?: Authorities;
		page?: number;
		pageCount?: number;
		filter?: Filter;
		count?: number;
	};
}
