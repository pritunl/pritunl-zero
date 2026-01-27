/// <reference path="../References.d.ts"/>
export const SYNC = 'check.sync';
export const SYNC_NAMES = 'check.sync_names';
export const TRAVERSE = 'check.traverse';
export const FILTER = 'check.filter';
export const CHANGE = 'check.change';

export interface Check {
	id?: string;
	name?: string;
	roles?: string[];
	frequency?: number;
	type?: string;
	targets?: string[];
	timeout?: number;
	status_code?: number;
	headers?: Header[];
	states?: State[];
}

export interface Header {
	key?: string;
	value?: string;
}

export interface State {
	e?: string;
	t?: string;
	x?: string[];
	l?: string[];
	r?: string[];
}

export interface Filter {
	id?: string;
	name?: string;
	role?: string;
}

export type Checks = Check[];

export type CheckRo = Readonly<Check>;
export type ChecksRo = ReadonlyArray<CheckRo>;

export interface CheckDispatch {
	type: string;
	data?: {
		id?: string;
		check?: Check;
		checks?: Checks;
		page?: number;
		pageCount?: number;
		filter?: Filter;
		count?: number;
	};
}
