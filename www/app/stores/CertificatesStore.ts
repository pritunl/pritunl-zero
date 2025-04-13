/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as CertificateTypes from '../types/CertificateTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class CertificatesStore extends EventEmitter {
	_certificates: CertificateTypes.CertificatesRo = Object.freeze([]);
	_page: number;
	_pageCount: number;
	_filter: CertificateTypes.Filter = null;
	_count: number;
	_map: {[key: string]: number} = {};
	_token = Dispatcher.register((this._callback).bind(this));

	get certificates(): CertificateTypes.CertificatesRo {
		return this._certificates;
	}

	get certificatesM(): CertificateTypes.Certificates {
		let certificates: CertificateTypes.Certificates = [];
		this._certificates.forEach((
			certificate: CertificateTypes.CertificateRo): void => {

			certificates.push({
				...certificate,
			});
		});
		return certificates;
	}

	get page(): number {
		return this._page || 0;
	}

	get pageCount(): number {
		return this._pageCount || 20;
	}

	get pages(): number {
		return Math.ceil(this.count / this.pageCount);
	}

	get filter(): CertificateTypes.Filter {
		return this._filter;
	}

	get count(): number {
		return this._count || 0;
	}

	certificate(id: string): CertificateTypes.CertificateRo {
		let i = this._map[id];
		if (i === undefined) {
			return null;
		}
		return this._certificates[i];
	}

	emitChange(): void {
		this.emitDefer(GlobalTypes.CHANGE);
	}

	addChangeListener(callback: () => void): void {
		this.on(GlobalTypes.CHANGE, callback);
	}

	removeChangeListener(callback: () => void): void {
		this.removeListener(GlobalTypes.CHANGE, callback);
	}

	_traverse(page: number): void {
		this._page = Math.min(this.pages, page);
	}

	_filterCallback(filter: CertificateTypes.Filter): void {
		if ((this._filter !== null && filter === null) ||
			(!Object.keys(this._filter || {}).length && filter !== null) || (
				filter && this._filter && (
					filter.name !== this._filter.name
				))) {
			this._traverse(0);
		}
		this._filter = filter;
		this.emitChange();
	}

	_sync(certificates: CertificateTypes.Certificate[], count: number): void {
		this._map = {};
		for (let i = 0; i < certificates.length; i++) {
			certificates[i] = Object.freeze(certificates[i]);
			this._map[certificates[i].id] = i;
		}

		this._count = count;
		this._certificates = Object.freeze(certificates);
		this._page = Math.min(this.pages, this.page);

		this.emitChange();
	}

	_callback(action: CertificateTypes.CertificateDispatch): void {
		switch (action.type) {
			case CertificateTypes.TRAVERSE:
				this._traverse(action.data.page);
				break;

			case CertificateTypes.FILTER:
				this._filterCallback(action.data.filter);
				break;

			case CertificateTypes.SYNC:
				this._sync(action.data.certificates, action.data.count);
				break;
		}
	}
}

export default new CertificatesStore();
