/// <reference path="../References.d.ts"/>
export const SYNC = 'policy.sync';
export const TRAVERSE = 'policy.traverse';
export const FILTER = 'policy.filter';
export const CHANGE = 'policy.change';

export interface Rule {
	type?: string;
	disable?: boolean;
	values?: string[];
}

export interface Policy {
	id?: string;
	name?: string;
	disabled?: boolean;
	services?: string[];
	authorities?: string[];
	roles?: string[];
	rules?: {[key: string]: Rule};
	admin_secondary?: string;
	user_secondary?: string;
	proxy_secondary?: string;
	authority_secondary?: string;
	admin_device_secondary?: boolean;
	user_device_secondary?: boolean;
	proxy_device_secondary?: boolean;
	authority_device_secondary?: boolean;
	authority_require_smart_card?: boolean;
}

export interface Filter {
	id?: string;
	name?: string;
}

export type Policies = Policy[];

export type PolicyRo = Readonly<Policy>;
export type PoliciesRo = ReadonlyArray<PolicyRo>;

export interface PolicyDispatch {
	type: string;
	data?: {
		id?: string;
		policy?: Policy;
		policies?: Policies;
		page?: number;
		pageCount?: number;
		filter?: Filter;
		count?: number;
	};
}
