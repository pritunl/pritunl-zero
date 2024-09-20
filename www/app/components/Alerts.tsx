/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as AlertTypes from '../types/AlertTypes';
import * as AuthorityTypes from '../types/AuthorityTypes';
import AlertsStore from '../stores/AlertsStore';
import AuthoritiesStore from '../stores/AuthoritiesStore';
import * as AlertActions from '../actions/AlertActions';
import * as AuthorityActions from '../actions/AuthorityActions';
import Alert from './Alert';
import AlertsFilter from './AlertsFilter';
import AlertsPage from './AlertsPage';
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
	alerts: AlertTypes.AlertsRo;
	filter: AlertTypes.Filter;
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

export default class Alerts extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			alerts: AlertsStore.alerts,
			filter: AlertsStore.filter,
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
		AlertsStore.addChangeListener(this.onChange);
		AuthoritiesStore.addChangeListener(this.onChange);
		AlertActions.sync();
		AuthorityActions.sync();
	}

	componentWillUnmount(): void {
		AlertsStore.removeChangeListener(this.onChange);
		AuthoritiesStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		let alerts = AlertsStore.alerts;
		let selected: Selected = {};
		let curSelected = this.state.selected;
		let opened: Opened = {};
		let curOpened = this.state.opened;

		alerts.forEach((alert: AlertTypes.Alert): void => {
			if (curSelected[alert.id]) {
				selected[alert.id] = true;
			}
			if (curOpened[alert.id]) {
				opened[alert.id] = true;
			}
		});

		this.setState({
			...this.state,
			alerts: alerts,
			filter: AlertsStore.filter,
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
		AlertActions.removeMulti(
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
		let alertsDom: JSX.Element[] = [];

		this.state.alerts.forEach((
				alert: AlertTypes.AlertRo): void => {
			alertsDom.push(<Alert
				key={alert.id}
				alert={alert}
				authorities={this.state.authorities}
				selected={!!this.state.selected[alert.id]}
				open={!!this.state.opened[alert.id]}
				onSelect={(shift: boolean): void => {
					let selected = {
						...this.state.selected,
					};

					if (shift) {
						let alerts = this.state.alerts;
						let start: number;
						let end: number;

						for (let i = 0; i < alerts.length; i++) {
							let usr = alerts[i];

							if (usr.id === alert.id) {
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
								selected[alerts[i].id] = true;
							}

							this.setState({
								...this.state,
								lastSelected: alert.id,
								selected: selected,
							});

							return;
						}
					}

					if (selected[alert.id]) {
						delete selected[alert.id];
					} else {
						selected[alert.id] = true;
					}

					this.setState({
						...this.state,
						lastSelected: alert.id,
						selected: selected,
					});
				}}
				onOpen={(): void => {
					let opened = {
						...this.state.opened,
					};

					if (opened[alert.id]) {
						delete opened[alert.id];
					} else {
						opened[alert.id] = true;
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
			let inst = AlertsStore.alert(instId);
			if (inst) {
				selectedNames.push(inst.name || instId);
			} else {
				selectedNames.push(instId);
			}
		}

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Alerts</h2>
					<div className="flex"/>
					<div style={css.buttons}>
						<button
							className={filterClass}
							style={css.button}
							type="button"
							onClick={(): void => {
								if (this.state.filter === null) {
									AlertActions.filter({});
								} else {
									AlertActions.filter(null);
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
							confirmMsg="Permanently delete the selected alerts"
							confirmInput={true}
							items={selectedNames}
							disabled={!this.selected || this.state.disabled}
							onConfirm={this.onDelete}
						/>
						<button
							className="bp5-button bp5-intent-success bp5-icon-add"
							style={css.button}
							disabled={this.state.disabled}
							type="button"
							onClick={(): void => {
								this.setState({
									...this.state,
									disabled: true,
								});
								AlertActions.create({
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
			<AlertsFilter
				filter={this.state.filter}
				onFilter={(filter): void => {
					AlertActions.filter(filter);
				}}
				authorities={this.state.authorities}
			/>
			<div style={css.itemsBox}>
				<div style={css.items}>
					{alertsDom}
					<tr className="bp5-card bp5-row" style={css.placeholder}>
						<td colSpan={5} style={css.placeholder}/>
					</tr>
				</div>
			</div>
			<NonState
				hidden={!!alertsDom.length}
				iconClass="bp5-icon-notifications"
				title="No alerts"
				description="Add a new alert to get started."
			/>
			<AlertsPage
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
