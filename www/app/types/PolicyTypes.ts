/// <reference path="../References.d.ts"/>
export const SYNC = 'policy.sync';
export const CHANGE = 'policy.change';

export interface Rule {
	type?: string;
	value?: string;
}

export interface Policy {
	id: string;
	name?: string;
	services?: string[];
	roles?: string[];
	rules?: Rule[];
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
	};
}
