/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as NodeTypes from '../types/NodeTypes';
import NodesStore from '../stores/NodesStore';
import * as NodeActions from '../actions/NodeActions';
import Node from './Node';
import Page from './Page';
import PageHeader from './PageHeader';

interface State {
	nodes: NodeTypes.NodesRo;
	disabled: boolean;
}

const css = {
	header: {
		marginTop: '-19px',
	} as React.CSSProperties,
	heading: {
		margin: '19px 0 0 0',
	} as React.CSSProperties,
	button: {
		margin: '10px 0 0 10px',
	} as React.CSSProperties,
	buttonFirst: {
		margin: '10px 0 0 0',
	} as React.CSSProperties,
	noNodes: {
		height: 'auto',
	} as React.CSSProperties,
};

export default class Nodes extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			nodes: NodesStore.nodes,
			disabled: false,
		};
	}

	componentDidMount(): void {
		NodesStore.addChangeListener(this.onChange);
		NodeActions.sync();
	}

	componentWillUnmount(): void {
		NodesStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			nodes: NodesStore.nodes,
		});
	}

	render(): JSX.Element {
		let nodesDom: JSX.Element[] = [];

		this.state.nodes.forEach((node: NodeTypes.NodeRo): void => {
			nodesDom.push(<Node
				key={node.id}
				node={node}
			/>);
		});

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Nodes</h2>
					<div className="flex"/>
				</div>
			</PageHeader>
			<div>
				{nodesDom}
			</div>
		</Page>;
	}
}
