/// <reference path="../References.d.ts"/>
export function uuid(): string {
	return (+new Date() + Math.floor(Math.random() * 999999)).toString(36);
}

export function zeroPad(num: number, width: number): string {
	if (num < Math.pow(10, width)) {
		return ('0'.repeat(width - 1) + num).slice(-width);
	}
	return num.toString();
}

export function formatDate(date: Date): string {
	let str = '';

	str += zeroPad(date.getHours(), 2) + ':';
	str += zeroPad(date.getMinutes(), 2) + ':';
	str += zeroPad(date.getSeconds(), 2) + ' ';
	str += zeroPad(date.getMonth() + 1, 2) + '-';
	str += zeroPad(date.getDate(), 2) + '-';
	str += date.getFullYear().toString();

	return str;
}
