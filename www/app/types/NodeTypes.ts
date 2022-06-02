/// <reference path="../References.d.ts"/>
export const SYNC = 'node.sync';
export const CHANGE = 'node.change';

export interface Node {
	id: string;
	type?: string;
	name?: string;
	port?: number;
	no_redirect_server?: boolean;
	protocol?: string;
	timestamp?: string;
	management_domain?: string;
	user_domain?: string;
	webauthn_domain?: string;
	certificates?: string[];
	requests_min?: number;
	memory?: number;
	load1?: number;
	load5?: number;
	load15?: number;
	services?: string[];
	authorities?: string[];
	forwarded_for_header?: string;
	forwarded_proto_header?: string;
	software_version?: string;
	hostname?: string;
}

export type Nodes = Node[];

export type NodeRo = Readonly<Node>;
export type NodesRo = ReadonlyArray<NodeRo>;

export interface NodeDispatch {
	type: string;
	data?: {
		id?: string;
		node?: Node;
		nodes?: Nodes;
	};
}
