/// <reference path="./References.d.ts"/>
import * as SuperAgent from 'superagent';
import * as Blueprint from '@blueprintjs/core';

let toaster = Blueprint.Toaster.create({
	position: Blueprint.Position.BOTTOM,
});

export function success(message: string, timeout?: number): void {
	if (timeout === undefined) {
		timeout = 5000;
	}

	toaster.show({
		intent: Blueprint.Intent.SUCCESS,
		message: message,
		timeout: timeout,
	});
}

export function info(message: string, timeout?: number): void {
	if (timeout === undefined) {
		timeout = 5000;
	}

	toaster.show({
		intent: Blueprint.Intent.PRIMARY,
		message: message,
		timeout: timeout,
	});
}

export function warning(message: string, timeout?: number): void {
	if (timeout === undefined) {
		timeout = 5000;
	}

	toaster.show({
		intent: Blueprint.Intent.WARNING,
		message: message,
		timeout: timeout,
	});
}

export function error(message: string, timeout?: number): void {
	if (timeout === undefined) {
		timeout = 5000;
	}

	toaster.show({
		intent: Blueprint.Intent.DANGER,
		message: message,
		timeout: timeout,
	});
}

export function errorRes(res: SuperAgent.Response, message: string,
		timeout?: number): void {
	if (timeout === undefined) {
		timeout = 5000;
	}

	try {
		message = res.body.error_msg || message;
	} catch(err) {
	}

	toaster.show({
		intent: Blueprint.Intent.DANGER,
		message: message,
		timeout: timeout,
	});
}
