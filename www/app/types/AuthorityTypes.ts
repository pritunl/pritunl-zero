/// <reference path="../References.d.ts"/>
export const SYNC = 'authority.sync';
export const CHANGE = 'authority.change';

export interface Authority {
	id: string;
	name?: string;
	type?: string[];
	roles?: string[];
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
