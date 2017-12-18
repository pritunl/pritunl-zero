/// <reference path="../References.d.ts"/>
export const LOAD = 'keybase.load';
export const UNLOAD = 'keybase.unload';
export const CHANGE = 'keybase.change';

export interface Info {
	username?: string;
	picture?: string;
	twitter?: string;
	github?: string;
}

export type InfoRo = Readonly<Info>;

export interface InfoDispatch {
	type: string;
	data?: {
		info?: Info;
	};
}
