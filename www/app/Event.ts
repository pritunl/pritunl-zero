/// <reference path="./References.d.ts"/>
import EventDispatcher from './dispatcher/EventDispatcher';
import * as Csrf from './Csrf';

let connected = false;
const pendingEvents: Record<string, any> = {};

function connect(): void {
	let url = '';
	let location = window.location;

	if (location.protocol === 'https:') {
		url += 'wss';
	} else {
		url += 'ws';
	}

	url += '://' + location.host + '/event?csrf_token=' + Csrf.token;

	let socket = new WebSocket(url);

	socket.addEventListener('close', () => {
		setTimeout(() => {
			connect();
		}, 500);
	});

	socket.addEventListener('message', (evt) => {
		const eventData = JSON.parse(evt.data).data;
		const eventId = JSON.stringify(eventData);

		if (pendingEvents[eventId]) {
			return;
		}

		pendingEvents[eventId] = eventData;

		setTimeout(() => {
			if (pendingEvents[eventId]) {
				console.log(eventData);
				EventDispatcher.dispatch(eventData);

				delete pendingEvents[eventId];
			}
		}, 300);
	});
}

export function init() {
	if (connected) {
		return;
	}
	connected = true;

	connect();
}
