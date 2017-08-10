/// <reference path="../References.d.ts"/>
export const SYNC = 'audit.sync';
export const CHANGE = 'audit.change';

export interface Agent {
	operating_system?: string;
	browser?: string;
	ip?: string;
	isp?: string;
	continent?: string;
	continent_code?: string;
	country?: string;
	country_code?: string;
	region?: string;
	region_code?: string;
	city?: string;
	latitude?: number;
	longitude?: number;
}

export interface Audit {
	id: string;
	user?: string;
	timestamp?: string;
	type?: string;
	fields?: {[key: string]: string};
	agent?: Agent;
}

export type Audits = Audit[];

export type AuditRo = Readonly<Audit>;
export type AuditsRo = ReadonlyArray<AuditRo>;

export interface AuditDispatch {
	type: string;
	data?: {
		id?: string;
		userId?: string;
		audit?: Audit;
		audits?: Audits;
	};
}
