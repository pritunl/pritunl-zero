/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as CertificateTypes from '../types/CertificateTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class CertificatesStore extends EventEmitter {
	_certificates: CertificateTypes.CertificatesRo = Object.freeze([]);
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

	_sync(certificates: CertificateTypes.Certificate[]): void {
		this._map = {};
		for (let i = 0; i < certificates.length; i++) {
			certificates[i] = Object.freeze(certificates[i]);
			this._map[certificates[i].id] = i;
		}

		this._certificates = Object.freeze(certificates);
		this.emitChange();
	}

	_callback(action: CertificateTypes.CertificateDispatch): void {
		switch (action.type) {
			case CertificateTypes.SYNC:
				this._sync(action.data.certificates);
				break;
		}
	}
}

export default new CertificatesStore();
