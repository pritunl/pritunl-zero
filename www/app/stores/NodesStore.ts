/// <reference path="../References.d.ts"/>
import Dispatcher from '../dispatcher/Dispatcher';
import EventEmitter from '../EventEmitter';
import * as NodeTypes from '../types/NodeTypes';
import * as GlobalTypes from '../types/GlobalTypes';

class NodesStore extends EventEmitter {
	_nodes: NodeTypes.NodesRo = Object.freeze([]);
	_page: number;
	_pageCount: number;
	_filter: NodeTypes.Filter = null;
	_count: number;
	_map: {[key: string]: number} = {};
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

	get page(): number {
		return this._page || 0;
	}

	get pageCount(): number {
		return this._pageCount || 20;
	}

	get pages(): number {
		return Math.ceil(this.count / this.pageCount);
	}

	get filter(): NodeTypes.Filter {
		return this._filter;
	}

	get count(): number {
		return this._count || 0;
	}

	node(id: string): NodeTypes.NodeRo {
		let i = this._map[id];
		if (i === undefined) {
			return null;
		}
		return this._nodes[i];
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

	_traverse(page: number): void {
		this._page = Math.min(this.pages, page);
	}

	_filterCallback(filter: NodeTypes.Filter): void {
		if ((this._filter !== null && filter === null) ||
			(!Object.keys(this._filter || {}).length && filter !== null) || (
				filter && this._filter && (
					filter.name !== this._filter.name
				))) {
			this._traverse(0);
		}
		this._filter = filter;
		this.emitChange();
	}

	_sync(nodes: NodeTypes.Node[], count: number): void {
		this._map = {};
		for (let i = 0; i < nodes.length; i++) {
			nodes[i] = Object.freeze(nodes[i]);
			this._map[nodes[i].id] = i;
		}

		this._count = count;
		this._nodes = Object.freeze(nodes);
		this._page = Math.min(this.pages, this.page);

		this.emitChange();
	}

	_callback(action: NodeTypes.NodeDispatch): void {
		switch (action.type) {
			case NodeTypes.TRAVERSE:
				this._traverse(action.data.page);
				break;

			case NodeTypes.FILTER:
				this._filterCallback(action.data.filter);
				break;

			case NodeTypes.SYNC:
				this._sync(action.data.nodes, action.data.count);
				break;
		}
	}
}

export default new NodesStore();
