/// <reference path="../References.d.ts"/>
export const SYNC = 'certificate.sync';
export const SYNC_NAMES = 'certificate.sync_names';
export const TRAVERSE = 'certificate.traverse';
export const FILTER = 'certificate.filter';
export const CHANGE = 'certificate.change';

export interface Info {
	signature_alg?: string;
	public_key_alg?: string;
	issuer?: string;
	issued_on?: string;
	expires_on?: string;
	dns_names?: string[];
}

export interface Certificate {
	id?: string;
	name?: string;
	comment?: string;
	type?: string;
	key?: string;
	certificate?: string;
	info?: Info;
	acme_type?: string;
	acme_auth?: string;
	acme_secret?: string;
	acme_domains?: string[];
}

export interface Filter {
	id?: string;
	name?: string;
}

export type Certificates = Certificate[];

export type CertificateRo = Readonly<Certificate>;
export type CertificatesRo = ReadonlyArray<CertificateRo>;

export interface CertificateDispatch {
	type: string;
	data?: {
		id?: string;
		certificate?: Certificate;
		certificates?: Certificates;
		page?: number;
		pageCount?: number;
		filter?: Filter;
		count?: number;
	};
}
