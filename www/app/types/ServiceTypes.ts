/// <reference path="../References.d.ts"/>
export const SYNC = 'service.sync';
export const CHANGE = 'service.change';

export interface Domain {
	domain?: string;
	host?: string;
}

export interface Server {
	protocol?: string;
	hostname?: string;
	port?: number;
}

export interface Service {
	id: string;
	name?: string;
	share_session?: boolean;
	websockets?: boolean;
	domains?: Domain[];
	roles?: string[];
	servers?: Server[];
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
	};
}
