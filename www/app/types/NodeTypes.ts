/// <reference path="../References.d.ts"/>
export const SYNC = 'node.sync';
export const CHANGE = 'node.change';

export interface Node {
	id: string;
	name?: string;
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
