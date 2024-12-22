/// <reference path="../References.d.ts"/>
export interface Params {
	[key: string]: string
}

export interface State {
	path: string
	spec: string
	params: Params
	matched: boolean
}

let curState: State

export function match(spec: string, path: string): State {
	const specSpl = spec.split('/');
	const testSpl = path.split('/');

	if (spec === "/" && path === "") {
		return {
			path: "/",
			spec: "/",
			params: {},
			matched: true,
		}
	}

	if (specSpl.length !== testSpl.length) {
		return {
			path: "",
			spec: "",
			params: {},
			matched: false,
		}
	}

	const params: Params = {};

	for (let i = 0; i < specSpl.length; i++) {
		const specPart = specSpl[i];
		const testPart = testSpl[i];

		if (specPart.startsWith(':')) {
			params[specPart.substring(1)] = testPart;
		} else if (specPart !== testPart) {
			return {
				path: "",
				spec: "",
				params: {},
				matched: false,
			}
		}
	}

	return {
		path: path,
		spec: spec,
		params: params,
		matched: true,
	}
}

export function getState(): State {
	return curState
}

export function setState(data: State) {
	curState = data
}
