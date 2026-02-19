/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as MiscUtils from '../utils/MiscUtils';
import * as NodeTypes from '../types/NodeTypes';
import * as ServiceTypes from '../types/ServiceTypes';
import * as AuthorityTypes from '../types/AuthorityTypes';
import * as CertificateTypes from '../types/CertificateTypes';
import NodesStore from '../stores/NodesStore';
import ServicesStore from '../stores/ServicesStore';
import AuthoritiesStore from '../stores/AuthoritiesStore';
import CertificatesStore from '../stores/CertificatesStore';
import * as NodeActions from '../actions/NodeActions';
import * as ServiceActions from '../actions/ServiceActions';
import * as AuthorityActions from '../actions/AuthorityActions';
import * as CertificateActions from '../actions/CertificateActions';
import Node from './Node';
import NodesFilter from './NodesFilter';
import NodesPage from './NodesPage';
import Page from './Page';
import PageHeader from './PageHeader';

interface Selected {
	[key: string]: boolean;
}

interface Opened {
	[key: string]: boolean;
}

interface State {
	nodes: NodeTypes.NodesRo;
	filter: NodeTypes.Filter;
	services: ServiceTypes.ServicesRo;
	authorities: AuthorityTypes.AuthoritiesRo;
	certificates: CertificateTypes.CertificatesRo;
	selected: Selected;
	opened: Opened;
	lastSelected: string;
	disabled: boolean;
}

const css = {
	items: {
		width: '100%',
		marginTop: '-5px',
		display: 'table',
		tableLayout: 'fixed',
		borderSpacing: '0 5px',
	} as React.CSSProperties,
	itemsBox: {
		width: '100%',
		overflowY: 'auto',
	} as React.CSSProperties,
	placeholder: {
		opacity: 0,
		width: '100%',
	} as React.CSSProperties,
	header: {
		marginTop: '-19px',
	} as React.CSSProperties,
	heading: {
		margin: '19px 0 0 0',
	} as React.CSSProperties,
	button: {
		margin: '8px 0 0 8px',
	} as React.CSSProperties,
	buttons: {
		marginTop: '8px',
	} as React.CSSProperties,
};

export default class Nodes extends React.Component<{}, State> {
	sync: MiscUtils.SyncInterval;

	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			nodes: NodesStore.nodes,
			filter: NodesStore.filter,
			services: ServicesStore.servicesName,
			authorities: AuthoritiesStore.authorities,
			certificates: CertificatesStore.certificates,
			selected: {},
			opened: {},
			lastSelected: null,
			disabled: false,
		};
	}

	get selected(): boolean {
		return !!Object.keys(this.state.selected).length;
	}

	get opened(): boolean {
		return !!Object.keys(this.state.opened).length;
	}

	componentDidMount(): void {
		NodesStore.addChangeListener(this.onChange);
		ServicesStore.addChangeListener(this.onChange);
		AuthoritiesStore.addChangeListener(this.onChange);
		CertificatesStore.addChangeListener(this.onChange);
		NodeActions.sync();
		ServiceActions.syncNames();
		AuthorityActions.syncNames();
		CertificateActions.syncNames();

		this.sync = new MiscUtils.SyncInterval(
			() => NodeActions.sync(true),
			2000,
		)
	}

	componentWillUnmount(): void {
		NodesStore.removeChangeListener(this.onChange);
		ServicesStore.removeChangeListener(this.onChange);
		AuthoritiesStore.removeChangeListener(this.onChange);
		CertificatesStore.removeChangeListener(this.onChange);

		this.sync?.stop()
	}

	onChange = (): void => {
		let nodes = NodesStore.nodes;
		let selected: Selected = {};
		let curSelected = this.state.selected;
		let opened: Opened = {};
		let curOpened = this.state.opened;

		nodes.forEach((node: NodeTypes.Node): void => {
			if (curSelected[node.id]) {
				selected[node.id] = true;
			}
			if (curOpened[node.id]) {
				opened[node.id] = true;
			}
		});

		this.setState({
			...this.state,
			nodes: nodes,
			filter: NodesStore.filter,
			services: ServicesStore.servicesName,
			authorities: AuthoritiesStore.authoritiesName,
			certificates: CertificatesStore.certificatesName,
			selected: selected,
			opened: opened,
		});
	}

	render(): JSX.Element {
		let nodesDom: JSX.Element[] = [];

		this.state.nodes.forEach((node: NodeTypes.NodeRo): void => {
			nodesDom.push(<Node
				key={node.id}
				node={node}
				services={this.state.services}
				authorities={this.state.authorities}
				certificates={this.state.certificates}
				selected={!!this.state.selected[node.id]}
				open={!!this.state.opened[node.id]}
				onSelect={(shift: boolean): void => {
					let selected = {
						...this.state.selected,
					};

					if (shift) {
						let nodes = this.state.nodes;
						let start: number;
						let end: number;

						for (let i = 0; i < nodes.length; i++) {
							let usr = nodes[i];

							if (usr.id === node.id) {
								start = i;
							} else if (usr.id === this.state.lastSelected) {
								end = i;
							}
						}

						if (start !== undefined && end !== undefined) {
							if (start > end) {
								end = [start, start = end][0];
							}

							for (let i = start; i <= end; i++) {
								selected[nodes[i].id] = true;
							}

							this.setState({
								...this.state,
								lastSelected: node.id,
								selected: selected,
							});

							return;
						}
					}

					if (selected[node.id]) {
						delete selected[node.id];
					} else {
						selected[node.id] = true;
					}

					this.setState({
						...this.state,
						lastSelected: node.id,
						selected: selected,
					});
				}}
				onOpen={(): void => {
					let opened = {
						...this.state.opened,
					};

					if (opened[node.id]) {
						delete opened[node.id];
					} else {
						opened[node.id] = true;
					}

					this.setState({
						...this.state,
						opened: opened,
					});
				}}
			/>);
		});

		let filterClass = 'bp5-button bp5-intent-primary bp5-icon-filter ';
		if (this.state.filter) {
			filterClass += 'bp5-active';
		}

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Nodes</h2>
					<div className="flex"/>
					<div style={css.buttons}>
						<button
							className={filterClass}
							style={css.button}
							type="button"
							onClick={(): void => {
								if (this.state.filter === null) {
									NodeActions.filter({});
								} else {
									NodeActions.filter(null);
								}
							}}
						>
							Filters
						</button>
						<button
							className="bp5-button bp5-intent-warning bp5-icon-chevron-up"
							style={css.button}
							disabled={!this.opened}
							type="button"
							onClick={(): void => {
								this.setState({
									...this.state,
									opened: {},
								});
							}}
						>
							Collapse All
						</button>
					</div>
				</div>
			</PageHeader>
			<NodesFilter
				filter={this.state.filter}
				onFilter={(filter): void => {
					NodeActions.filter(filter);
				}}
			/>
			<div style={css.itemsBox}>
				<div style={css.items}>
					{nodesDom}
					<tr className="bp5-card bp5-row" style={css.placeholder}>
						<td colSpan={4} style={css.placeholder}/>
					</tr>
				</div>
			</div>
			<NodesPage
				onPage={(): void => {
					this.setState({
						...this.state,
						selected: {},
						lastSelected: null,
					});
				}}
			/>
		</Page>;
	}
}
