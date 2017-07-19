/// <reference path="../References.d.ts"/>
export const SYNC = 'node.sync';
export const CHANGE = 'node.change';

export interface Node {
	id: string;
	type?: string;
	name?: string;
	port?: number;
	protocol?: string;
	timestamp?: string;
	memory?: number;
	load1?: number;
	load5?: number;
	load15?: number;
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
