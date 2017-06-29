/// <reference path="../References.d.ts"/>
export const SYNC = 'subscription.sync';
export const CHANGE = 'subscription.change';

export interface Subscription {
	active?: boolean;
	status?: string;
	plan?: string;
	quantity?: number;
}

export type SubscriptionRo = Readonly<Subscription>

export interface SubscriptionDispatch {
	type: string;
	data?: Subscription;
}
