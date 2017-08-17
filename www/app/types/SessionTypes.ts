/// <reference path="../References.d.ts"/>
import * as AgentTypes from './AgentTypes';

export const SYNC = 'session.sync';
export const CHANGE = 'session.change';
export const SHOW_REMOVED = 'session.show_removed';

export interface Session {
	id: string;
	user?: string;
	timestamp?: string;
	last_active?: string;
	removed?: boolean;
	agent?: AgentTypes.Agent;
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
		showRemoved?: boolean;
	};
}
