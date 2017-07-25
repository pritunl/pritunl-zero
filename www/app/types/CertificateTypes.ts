/// <reference path="../References.d.ts"/>
export const SYNC = 'certificate.sync';
export const CHANGE = 'certificate.change';

export interface Certificate {
	id: string;
	name?: string;
	type?: string;
	key?: string;
	certificate?: string;
	acme_domains?: string[];
	acme_account?: string;
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
	};
}
