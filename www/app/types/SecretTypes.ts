/// <reference path="../References.d.ts"/>
export const SYNC = 'secret.sync';
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

export type Secrets = Secret[];

export type SecretRo = Readonly<Secret>;
export type SecretsRo = ReadonlyArray<SecretRo>;

export interface SecretDispatch {
	type: string;
	data?: {
		id?: string;
		secret?: Secret;
		secrets?: Secrets;
	};
}
