/// <reference path="./References.d.ts"/>

export function setLocation(location: string) {
	window.location.hash = location
	let evt = new Event("router_update")
	window.dispatchEvent(evt)
}

export function reload() {
	let evt = new Event("router_update")
	window.dispatchEvent(evt)
}
