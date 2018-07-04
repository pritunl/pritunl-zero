/// <reference path="../References.d.ts"/>
export const SSH_TOKEN = 'state.ssh_token';
export const SSH_DEVICE = 'state.ssh_device';

export interface StateDispatch {
	type: string;
	data?: string;
}
