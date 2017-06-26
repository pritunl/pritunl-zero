/// <reference path="../References.d.ts"/>
export const SYNC = 'user.sync';
export const TRAVERSE = 'user.traverse';
export const LOAD = 'user.load';
export const UNLOAD = 'user.unload';
export const CHANGE = 'user.change';

export interface User {
	id: string;
	type?: string;
	username?: string;
	password?: string;
	last_active?: string;
	roles?: string[];
	administrator: string;
	permissions?: string[];
}

export type Users = User[];

export interface UserDispatch {
	type: string;
	data?: {
		id?: string;
		user?: User;
		users?: User[];
		page?: number;
		pageCount?: number;
		count?: number;
	};
}
