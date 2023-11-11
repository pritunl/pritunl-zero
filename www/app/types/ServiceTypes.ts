/// <reference path="../References.d.ts"/>
export const SYNC = 'service.sync';
export const SYNC_NAMES = 'service.sync_names';
export const TRAVERSE = 'service.traverse';
export const FILTER = 'service.filter';
export const CHANGE = 'service.change';

export interface Domain {
	domain?: string;
	host?: string;
}

export interface Path {
	path?: string;
}

export interface Server {
	protocol?: string;
	hostname?: string;
	port?: number;
}

export interface Service {
	id: string;
	name?: string;
	type?: string;
	share_session?: boolean;
	logout_path?: string;
	websockets?: boolean;
	disable_csrf_check?: boolean;
	client_authority?: string;
	domains?: Domain[];
	roles?: string[];
	servers?: Server[];
	whitelist_networks?: string[];
	whitelist_paths?: Path[];
	whitelist_options?: boolean;
}

export interface Filter {
	id?: string;
	name?: string;
	role?: string;
}

export type Services = Service[];

export type ServiceRo = Readonly<Service>;
export type ServicesRo = ReadonlyArray<ServiceRo>;

export interface ServiceDispatch {
	type: string;
	data?: {
		id?: string;
		service?: Service;
		services?: Services;
		page?: number;
		pageCount?: number;
		filter?: Filter;
		count?: number;
	};
}
