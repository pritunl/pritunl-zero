/// <reference path="../References.d.ts"/>
export const SYNC = 'log.sync';
export const TRAVERSE = 'log.traverse';
export const FILTER = 'log.filter';
export const CHANGE = 'log.change';

export interface Log {
	id: string;
	level?: string;
	timestamp?: string;
	message?: string;
	stack?: string;
	fields?: {[key: string]: any};
}

export interface Filter {
	level?: string;
	message?: string;
}

export type Logs = Log[];

export type LogRo = Readonly<Log>;
export type LogsRo = ReadonlyArray<LogRo>;

export interface LogDispatch {
	type: string;
	data?: {
		id?: string;
		log?: Log;
		logs?: Logs;
		page?: number;
		pageCount?: number;
		filter?: Filter;
		count?: number;
	};
}
