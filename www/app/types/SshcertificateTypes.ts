/// <reference path="../References.d.ts"/>
import * as AgentTypes from './AgentTypes';

export const SYNC = 'sshcertificate.sync';
export const TRAVERSE = 'sshcertificate.traverse';
export const CHANGE = 'sshcertificate.change';

export interface Info {
	serial?: string;
	expires?: string;
	principals?: string[];
	extensions?: string[];
}

export interface Sshcertificate {
	id: string;
	user_id?: string;
	authority_ids?: string[];
	timestamp?: string;
	agent?: AgentTypes.Agent;
	certificates_info?: Info[];
}

export type Sshcertificates = Sshcertificate[];

export type SshcertificateRo = Readonly<Sshcertificate>;
export type SshcertificatesRo = ReadonlyArray<SshcertificateRo>;

export interface SshcertificateDispatch {
	type: string;
	data?: {
		id?: string;
		userId?: string;
		certificate?: Sshcertificate;
		certificates?: Sshcertificates;
		page?: number;
		count?: number;
	};
}
