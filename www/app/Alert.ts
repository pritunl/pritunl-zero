/// <reference path="./References.d.ts"/>
import * as Blueprint from '@blueprintjs/core';

let toaster = Blueprint.Toaster.create({
	position: Blueprint.Position.BOTTOM,
});

export function info(message: string): void {
	toaster.show({
		intent: Blueprint.Intent.PRIMARY,
		message: message,
	});
}

export function warning(message: string): void {
	toaster.show({
		intent: Blueprint.Intent.WARNING,
		message: message,
	});
}

export function error(message: string): void {
	toaster.show({
		intent: Blueprint.Intent.DANGER,
		message: message,
	});
}
