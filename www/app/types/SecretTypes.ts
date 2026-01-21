/// <reference path="../References.d.ts"/>
export const SYNC = 'secret.sync';
export const TRAVERSE = 'secret.traverse';
export const FILTER = 'secret.filter';
export const CHANGE = 'secret.change';

export interface Secret {
	id?: string;
	name?: string;
	comment?: string;
	type?: string;
	key?: string;
	value?: string;
	region?: string;
	public_key?: string;
}

export interface Filter {
	id?: string;
	name?: string;
}

export type Secrets = Secret[];

export type SecretRo = Readonly<Secret>;
export type SecretsRo = ReadonlyArray<SecretRo>;

export interface SecretDispatch {
	type: string;
	data?: {
		id?: string;
		secret?: Secret;
		secrets?: Secrets;
		page?: number;
		pageCount?: number;
		filter?: Filter;
		count?: number;
	};
}
