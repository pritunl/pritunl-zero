/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as NodeTypes from '../types/NodeTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class NodesStore extends EventEmitter {
	_nodes: NodeTypes.NodesRo = Object.freeze([]);
	_token = Dispatcher.register((this._callback).bind(this));

	get nodes(): NodeTypes.NodesRo {
		return this._nodes;
	}

	get nodesM(): NodeTypes.Nodes {
		let nodes: NodeTypes.Nodes = [];
		this._nodes.forEach((node: NodeTypes.NodeRo): void => {
			nodes.push({
				...node,
			});
		});
		return nodes;
	}

	emitChange(): void {
		this.emitDefer(GlobalTypes.CHANGE);
	}

	addChangeListener(callback: () => void): void {
		this.on(GlobalTypes.CHANGE, callback);
	}

	removeChangeListener(callback: () => void): void {
		this.removeListener(GlobalTypes.CHANGE, callback);
	}

	_sync(nodes: NodeTypes.Node[]): void {
		for (let i = 0; i < nodes.length; i++) {
			nodes[i] = Object.freeze(nodes[i]);
		}

		this._nodes = Object.freeze(nodes);
		this.emitChange();
	}

	_callback(action: NodeTypes.NodeDispatch): void {
		switch (action.type) {
			case NodeTypes.SYNC:
				this._sync(action.data.nodes);
				break;
		}
	}
}

export default new NodesStore();
