/// <reference path="./References.d.ts"/>
import * as SuperAgent from 'superagent';

export let token = '';

export function load(): Promise<void> {
	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/csrf')
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				if (res && res.status === 401) {
					window.location.href = '/login';
					resolve();
					return;
				}

				if (err) {
					reject(err);
					return;
				}

				token = res.body.token;

				resolve();
			});
	});
}
