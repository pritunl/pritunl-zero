/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ServiceTypes from '../types/ServiceTypes';
import * as AuthorityTypes from '../types/AuthorityTypes';
import ServicesStore from '../stores/ServicesStore';
import AuthoritiesStore from '../stores/AuthoritiesStore';
import * as ServiceActions from '../actions/ServiceActions';
import * as AuthorityActions from '../actions/AuthorityActions';
import Service from './Service';
import ServiceNew from './ServiceNew';
import ServicesFilter from './ServicesFilter';
import ServicesPage from './ServicesPage';
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
	services: ServiceTypes.ServicesRo;
	filter: ServiceTypes.Filter;
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

export default class Services extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			services: ServicesStore.services,
			filter: ServicesStore.filter,
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
		ServicesStore.addChangeListener(this.onChange);
		AuthoritiesStore.addChangeListener(this.onChange);
		ServiceActions.sync();
		AuthorityActions.syncNames();
	}

	componentWillUnmount(): void {
		ServicesStore.removeChangeListener(this.onChange);
		AuthoritiesStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		let services = ServicesStore.services;
		let selected: Selected = {};
		let curSelected = this.state.selected;
		let opened: Opened = {};
		let curOpened = this.state.opened;

		services.forEach((service: ServiceTypes.Service): void => {
			if (curSelected[service.id]) {
				selected[service.id] = true;
			}
			if (curOpened[service.id]) {
				opened[service.id] = true;
			}
		});

		this.setState({
			...this.state,
			services: services,
			filter: ServicesStore.filter,
			authorities: AuthoritiesStore.authoritiesName,
			selected: selected,
			opened: opened,
		});
	}

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		ServiceActions.removeMulti(
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
		let servicesDom: JSX.Element[] = [];

		this.state.services.forEach((
				service: ServiceTypes.ServiceRo): void => {
			servicesDom.push(<Service
				key={service.id}
				service={service}
				authorities={this.state.authorities}
				selected={!!this.state.selected[service.id]}
				open={!!this.state.opened[service.id]}
				onSelect={(shift: boolean): void => {
					let selected = {
						...this.state.selected,
					};

					if (shift) {
						let services = this.state.services;
						let start: number;
						let end: number;

						for (let i = 0; i < services.length; i++) {
							let usr = services[i];

							if (usr.id === service.id) {
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
								selected[services[i].id] = true;
							}

							this.setState({
								...this.state,
								lastSelected: service.id,
								selected: selected,
							});

							return;
						}
					}

					if (selected[service.id]) {
						delete selected[service.id];
					} else {
						selected[service.id] = true;
					}

					this.setState({
						...this.state,
						lastSelected: service.id,
						selected: selected,
					});
				}}
				onOpen={(): void => {
					let opened = {
						...this.state.opened,
					};

					if (opened[service.id]) {
						delete opened[service.id];
					} else {
						opened[service.id] = true;
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

		let selectedNames: string[] = [];
		for (let instId of Object.keys(this.state.selected)) {
			let inst = ServicesStore.service(instId);
			if (inst) {
				selectedNames.push(inst.name || instId);
			} else {
				selectedNames.push(instId);
			}
		}

		let newServiceDom: JSX.Element;
		if (this.state.newOpened) {
			newServiceDom = <ServiceNew
				authorities={this.state.authorities}
				onClose={(): void => {
					this.setState({
						...this.state,
						newOpened: false,
					});
				}}
			/>;
		}

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Services</h2>
					<div className="flex"/>
					<div style={css.buttons}>
						<button
							className={filterClass}
							style={css.button}
							type="button"
							onClick={(): void => {
								if (this.state.filter === null) {
									ServiceActions.filter({});
								} else {
									ServiceActions.filter(null);
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
						<ConfirmButton
							label="Delete Selected"
							className="bp5-intent-danger bp5-icon-delete"
							progressClassName="bp5-intent-danger"
							safe={true}
							style={css.button}
							confirmMsg="Permanently delete the selected services"
							confirmInput={true}
							items={selectedNames}
							disabled={!this.selected || this.state.disabled}
							onConfirm={this.onDelete}
						/>
						<button
							className="bp5-button bp5-intent-success bp5-icon-add"
							style={css.button}
							disabled={this.state.disabled || this.state.newOpened}
							type="button"
							onClick={(): void => {
								this.setState({
									...this.state,
									newOpened: true,
								});
							}}
						>New</button>
					</div>
				</div>
			</PageHeader>
			<ServicesFilter
				filter={this.state.filter}
				onFilter={(filter): void => {
					ServiceActions.filter(filter);
				}}
				authorities={this.state.authorities}
			/>
			<div style={css.itemsBox}>
				<div style={css.items}>
					{newServiceDom}
					{servicesDom}
					<tr className="bp5-card bp5-row" style={css.placeholder}>
						<td colSpan={2} style={css.placeholder}/>
					</tr>
				</div>
			</div>
			<NonState
				hidden={!!servicesDom.length}
				iconClass="bp5-icon-cloud"
				title="No services"
				description="Add a new service to get started."
			/>
			<ServicesPage
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
