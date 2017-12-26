/// <reference path="../References.d.ts"/>
export const SYNC = 'authority.sync';
export const CHANGE = 'authority.change';

export interface Info {
	key_alg?: string;
}

export interface Authority {
	id: string;
	name?: string;
	type?: string[];
	info?: Info;
	expire?: number;
	host_expire?: number;
	match_roles?: boolean;
	roles?: string[];
	public_key?: string;
	host_domain?: string;
	host_proxy?: string;
	strict_host_checking?: boolean;
	host_tokens?: string[];
}

export type Authorities = Authority[];

export type AuthorityRo = Readonly<Authority>;
export type AuthoritiesRo = ReadonlyArray<AuthorityRo>;

export interface AuthorityDispatch {
	type: string;
	data?: {
		id?: string;
		authority?: Authority;
		authorities?: Authorities;
	};
}
