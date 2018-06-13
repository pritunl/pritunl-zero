/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import * as Alert from '../Alert';
import * as Csrf from '../Csrf';
import Loader from "../Loader";
import * as KeybaseTypes from '../types/KeybaseTypes';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;

export function load(token: string): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/keybase/info/' + token)
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (curSyncId !== syncId) {
					resolve();
					return;
				}

				if (err) {
					Alert.errorRes(res, 'Failed to load Keybase user information');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: KeybaseTypes.LOAD,
					data: {
						info: res.body,
					},
				});

				resolve();
			});
	});
}

export function unload(): void {
	Dispatcher.dispatch({
		type: KeybaseTypes.UNLOAD,
	});
}
