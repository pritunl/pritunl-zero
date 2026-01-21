/// <reference path="../References.d.ts"/>
export const SYNC = 'alert.sync';
export const SYNC_NAMES = 'alert.sync_names';
export const TRAVERSE = 'alert.traverse';
export const FILTER = 'alert.filter';
export const CHANGE = 'alert.change';

export interface Alert {
	id: string;
	name?: string;
	roles?: string[];
	resource?: string;
	level?: number;
	frequency?: number;
	ignores?: string[];
	value_int?: number;
	value_str?: string;
}

export interface Filter {
	id?: string;
	name?: string;
	role?: string;
}

export type Alerts = Alert[];

export type AlertRo = Readonly<Alert>;
export type AlertsRo = ReadonlyArray<AlertRo>;

export interface AlertDispatch {
	type: string;
	data?: {
		id?: string;
		alert?: Alert;
		alerts?: Alerts;
		page?: number;
		pageCount?: number;
		filter?: Filter;
		count?: number;
	};
}
