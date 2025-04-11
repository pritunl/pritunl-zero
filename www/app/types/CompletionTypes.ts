/// <reference path="../References.d.ts"/>

import * as NodeTypes from "./NodeTypes"
import * as CertificateTypes from "./CertificateTypes"
import * as SecretTypes from "./SecretTypes"

export const SYNC = "completion.sync"
export const FILTER = "completion.filter"
export const CHANGE = "completion.change"

export interface Completion {
	nodes?: NodeTypes.Node[]
	certificates?: CertificateTypes.Certificate[]
	secrets?: SecretTypes.Secret[]
}

export interface Filter {
}

export interface CompletionDispatch {
	type: string
	data?: {
		completion?: Completion
		filter?: Filter
	}
}
