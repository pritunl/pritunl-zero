/// <reference path="../References.d.ts"/>
export const SYNC = 'subscription.sync';
export const CHANGE = 'subscription.change';

export interface Subscription {
	active?: boolean;
	status?: string;
	plan?: string;
	quantity?: number;
	amount?: number;
	period_end?: string;
	trial_end?: string;
	cancel_at_period_end?: boolean;
	balance?: number;
	url_key?: string;
}

export type SubscriptionRo = Readonly<Subscription>

export interface SubscriptionDispatch {
	type: string;
	data?: Subscription;
}
