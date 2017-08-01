/// <reference path="../References.d.ts"/>
export const SYNC = 'session.sync';
export const CHANGE = 'session.change';

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

export interface Session {
	id: string;
	user_id?: string;
	timestamp?: string;
	last_active?: string;
	agent?: Agent;
}

export type Sessions = Session[];

export type SessionRo = Readonly<Session>;
export type SessionsRo = ReadonlyArray<SessionRo>;

export interface SessionDispatch {
	type: string;
	data?: {
		id?: string;
		userId?: string;
		session?: Session;
		sessions?: Sessions;
	};
}
