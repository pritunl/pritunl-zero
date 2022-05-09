/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as EndpointTypes from '../types/EndpointTypes';
import * as AuthorityTypes from '../types/AuthorityTypes';
import EndpointsStore from '../stores/EndpointsStore';
import AuthoritiesStore from '../stores/AuthoritiesStore';
import * as EndpointActions from '../actions/EndpointActions';
import * as AuthorityActions from '../actions/AuthorityActions';
import Endpoint from './Endpoint';
import EndpointsFilter from './EndpointsFilter';
import EndpointsPage from './EndpointsPage';
import Page from './Page';
import PageHeader from './PageHeader';
import NonState from './NonState';
import ConfirmButton from './ConfirmButton';

interface Selected {
	[key: string]: boolean;
}

interface Opened {
	[key: string]: boolean;
}

interface State {
	endpoints: EndpointTypes.EndpointsRo;
	filter: EndpointTypes.Filter;
	authorities: AuthorityTypes.AuthoritiesRo;
	selected: Selected;
	opened: Opened;
	newOpened: boolean;
	lastSelected: string;
	disabled: boolean;
}

const css = {
	items: {
		width: '100%',
		marginTop: '-5px',
		display: 'table',
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

export default class Endpoints extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			endpoints: EndpointsStore.endpoints,
			filter: EndpointsStore.filter,
			authorities: AuthoritiesStore.authorities,
			selected: {},
			opened: {},
			newOpened: false,
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
		EndpointsStore.addChangeListener(this.onChange);
		AuthoritiesStore.addChangeListener(this.onChange);
		EndpointActions.sync();
		AuthorityActions.sync();
	}

	componentWillUnmount(): void {
		EndpointsStore.removeChangeListener(this.onChange);
		AuthoritiesStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		let endpoints = EndpointsStore.endpoints;
		let selected: Selected = {};
		let curSelected = this.state.selected;
		let opened: Opened = {};
		let curOpened = this.state.opened;

		endpoints.forEach((endpoint: EndpointTypes.Endpoint): void => {
			if (curSelected[endpoint.id]) {
				selected[endpoint.id] = true;
			}
			if (curOpened[endpoint.id]) {
				opened[endpoint.id] = true;
			}
		});

		this.setState({
			...this.state,
			endpoints: endpoints,
			filter: EndpointsStore.filter,
			authorities: AuthoritiesStore.authorities,
			selected: selected,
			opened: opened,
		});
	}

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		EndpointActions.removeMulti(
				Object.keys(this.state.selected)).then((): void => {
			this.setState({
				...this.state,
				selected: {},
				disabled: false,
			});
		}).catch((): void => {
			this.setState({
				...this.state,
				disabled: false,
			});
		});
	}

	render(): JSX.Element {
		let endpointsDom: JSX.Element[] = [];

		this.state.endpoints.forEach((
				endpoint: EndpointTypes.EndpointRo): void => {
			endpointsDom.push(<Endpoint
				key={endpoint.id}
				endpoint={endpoint}
				authorities={this.state.authorities}
				selected={!!this.state.selected[endpoint.id]}
				open={!!this.state.opened[endpoint.id]}
				onSelect={(shift: boolean): void => {
					let selected = {
						...this.state.selected,
					};

					if (shift) {
						let endpoints = this.state.endpoints;
						let start: number;
						let end: number;

						for (let i = 0; i < endpoints.length; i++) {
							let usr = endpoints[i];

							if (usr.id === endpoint.id) {
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
								selected[endpoints[i].id] = true;
							}

							this.setState({
								...this.state,
								lastSelected: endpoint.id,
								selected: selected,
							});

							return;
						}
					}

					if (selected[endpoint.id]) {
						delete selected[endpoint.id];
					} else {
						selected[endpoint.id] = true;
					}

					this.setState({
						...this.state,
						lastSelected: endpoint.id,
						selected: selected,
					});
				}}
				onOpen={(): void => {
					let opened = {
						...this.state.opened,
					};

					if (opened[endpoint.id]) {
						delete opened[endpoint.id];
					} else {
						opened[endpoint.id] = true;
					}

					this.setState({
						...this.state,
						opened: opened,
					});
				}}
			/>);
		});

		let filterClass = 'bp3-button bp3-intent-primary bp3-icon-filter ';
		if (this.state.filter) {
			filterClass += 'bp3-active';
		}

		let selectedNames: string[] = [];
		for (let instId of Object.keys(this.state.selected)) {
			let inst = EndpointsStore.endpoint(instId);
			if (inst) {
				selectedNames.push(inst.name || instId);
			} else {
				selectedNames.push(instId);
			}
		}

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Endpoints</h2>
					<div className="flex"/>
					<div style={css.buttons}>
						<button
							className={filterClass}
							style={css.button}
							type="button"
							onClick={(): void => {
								if (this.state.filter === null) {
									EndpointActions.filter({});
								} else {
									EndpointActions.filter(null);
								}
							}}
						>
							Filters
						</button>
						<button
							className="bp3-button bp3-intent-warning bp3-icon-chevron-up"
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
						<ConfirmButton
							label="Delete Selected"
							className="bp3-intent-danger bp3-icon-delete"
							progressClassName="bp3-intent-danger"
							safe={true}
							style={css.button}
							confirmMsg="Permanently delete the selected endpoints"
							confirmInput={true}
							items={selectedNames}
							disabled={!this.selected || this.state.disabled}
							onConfirm={this.onDelete}
						/>
						<button
							className="bp3-button bp3-intent-success bp3-icon-add"
							style={css.button}
							disabled={this.state.disabled}
							type="button"
							onClick={(): void => {
								this.setState({
									...this.state,
									disabled: true,
								});
								EndpointActions.create({
									id: null,
								}).then((): void => {
									this.setState({
										...this.state,
										disabled: false,
									});
								}).catch((): void => {
									this.setState({
										...this.state,
										disabled: false,
									});
								});
							}}
						>New</button>
					</div>
				</div>
			</PageHeader>
			<EndpointsFilter
				filter={this.state.filter}
				onFilter={(filter): void => {
					EndpointActions.filter(filter);
				}}
				authorities={this.state.authorities}
			/>
			<div style={css.itemsBox}>
				<div style={css.items}>
					{endpointsDom}
					<tr className="bp3-card bp3-row" style={css.placeholder}>
						<td colSpan={5} style={css.placeholder}/>
					</tr>
				</div>
			</div>
			<NonState
				hidden={!!endpointsDom.length}
				iconClass="bp3-icon-shield"
				title="No endpoints"
				description="Add a new endpoint to get started."
			/>
			<EndpointsPage
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
