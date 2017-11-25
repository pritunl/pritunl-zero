/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as SshcertificateTypes from '../types/SshcertificateTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class SshcertificatesStore extends EventEmitter {
	_userId: string;
	_sshcertificates: SshcertificateTypes.SshcertificatesRo = Object.freeze([]);
	_page: number;
	_pageCount: number;
	_count: number;
	_token = Dispatcher.register((this._callback).bind(this));

	get userId(): string {
		return this._userId;
	}

	get sshcertificates(): SshcertificateTypes.SshcertificatesRo {
		return this._sshcertificates;
	}

	get sshcertificatesM(): SshcertificateTypes.Sshcertificates {
		let sshcertificates: SshcertificateTypes.Sshcertificates = [];
		this._sshcertificates.forEach((
				sshcertificate: SshcertificateTypes.SshcertificateRo): void => {
			sshcertificates.push({
				...sshcertificate,
			});
		});
		return sshcertificates;
	}

	get page(): number {
		return this._page || 0;
	}

	get pageCount(): number {
		return this._pageCount || 10;
	}

	get pages(): number {
		return Math.ceil(this.count / this.pageCount);
	}

	get count(): number {
		return this._count || 0;
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

	_sync(userId: string, certs: SshcertificateTypes.Sshcertificate[],
			count: number): void {
		this._userId = userId;

		for (let i = 0; i < certs.length; i++) {
			certs[i] = Object.freeze(certs[i]);
		}

		this._count = count;
		this._sshcertificates = Object.freeze(certs);
		this._page = Math.min(this.pages, this.page);

		this.emitChange();
	}

	_callback(action: SshcertificateTypes.SshcertificateDispatch): void {
		switch (action.type) {
			case SshcertificateTypes.TRAVERSE:
				this._traverse(action.data.page);
				break;

			case SshcertificateTypes.SYNC:
				this._sync(action.data.userId, action.data.certificates,
					action.data.count);
				break;
		}
	}
}

export default new SshcertificatesStore();
