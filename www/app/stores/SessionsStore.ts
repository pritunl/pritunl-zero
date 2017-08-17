/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as SessionTypes from '../types/SessionTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class SessionsStore extends EventEmitter {
	_userId: string;
	_showRemoved: boolean;
	_sessions: SessionTypes.SessionsRo = Object.freeze([]);
	_token = Dispatcher.register((this._callback).bind(this));

	get userId(): string {
		return this._userId;
	}

	get sessions(): SessionTypes.SessionsRo {
		return this._sessions;
	}

	get sessionsM(): SessionTypes.Sessions {
		let sessions: SessionTypes.Sessions = [];
		this._sessions.forEach((session: SessionTypes.SessionRo): void => {
			sessions.push({
				...session,
			});
		});
		return sessions;
	}

	get showRemoved(): boolean {
		return this._showRemoved;
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

	_sync(userId: string, sessions: SessionTypes.Session[]): void {
		this._userId = userId;

		for (let i = 0; i < sessions.length; i++) {
			sessions[i] = Object.freeze(sessions[i]);
		}

		this._sessions = Object.freeze(sessions);
		this.emitChange();
	}

	_setShowRemoved(state: boolean): void {
		this._showRemoved = state;
		this.emitChange();
	}

	_callback(action: SessionTypes.SessionDispatch): void {
		switch (action.type) {
			case SessionTypes.SYNC:
				this._sync(action.data.userId, action.data.sessions);
				break;
			case SessionTypes.SHOW_REMOVED:
				this._setShowRemoved(action.data.showRemoved);
				break;
		}
	}
}

export default new SessionsStore();
