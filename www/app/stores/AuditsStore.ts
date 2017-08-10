/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as AuditTypes from '../types/AuditTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class AuditsStore extends EventEmitter {
	_userId: string;
	_audits: AuditTypes.AuditsRo = Object.freeze([]);
	_token = Dispatcher.register((this._callback).bind(this));

	get userId(): string {
		return this._userId;
	}

	get audits(): AuditTypes.AuditsRo {
		return this._audits;
	}

	get auditsM(): AuditTypes.Audits {
		let audits: AuditTypes.Audits = [];
		this._audits.forEach((audit: AuditTypes.AuditRo): void => {
			audits.push({
				...audit,
			});
		});
		return audits;
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

	_sync(userId: string, audits: AuditTypes.Audit[]): void {
		this._userId = userId;

		for (let i = 0; i < audits.length; i++) {
			audits[i] = Object.freeze(audits[i]);
		}

		this._audits = Object.freeze(audits);
		this.emitChange();
	}

	_callback(action: AuditTypes.AuditDispatch): void {
		switch (action.type) {
			case AuditTypes.SYNC:
				this._sync(action.data.userId, action.data.audits);
				break;
		}
	}
}

export default new AuditsStore();
