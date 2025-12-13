/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as PolicyTypes from '../types/PolicyTypes';
import * as ServiceTypes from '../types/ServiceTypes';
import * as AuthorityTypes from '../types/AuthorityTypes';
import * as SettingsTypes from '../types/SettingsTypes';
import PoliciesStore from '../stores/PoliciesStore';
import ServicesStore from '../stores/ServicesStore';
import AuthoritiesStore from '../stores/AuthoritiesStore';
import SettingsStore from '../stores/SettingsStore';
import * as PolicyActions from '../actions/PolicyActions';
import * as ServiceActions from '../actions/ServiceActions';
import * as AuthorityActions from '../actions/AuthorityActions';
import * as SettingsActions from '../actions/SettingsActions';
import Policy from './Policy';
import PolicyNew from './PolicyNew';
import PoliciesFilter from './PoliciesFilter';
import PoliciesPage from './PoliciesPage';
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
	policies: PolicyTypes.PoliciesRo;
	services: ServiceTypes.ServicesRo;
	authorities: AuthorityTypes.AuthoritiesRo;
	providers: SettingsTypes.SecondaryProviders;
	filter: PolicyTypes.Filter;
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

export default class Policies extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			policies: PoliciesStore.policies,
			services: ServicesStore.servicesName,
			authorities: AuthoritiesStore.authorities,
			providers: SettingsStore.settings ?
				SettingsStore.settings.auth_secondary_providers : [],
			filter: PoliciesStore.filter,
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
		PoliciesStore.addChangeListener(this.onChange);
		ServicesStore.addChangeListener(this.onChange);
		AuthoritiesStore.addChangeListener(this.onChange);
		SettingsStore.addChangeListener(this.onChange);
		PolicyActions.sync();
		ServiceActions.syncNames();
		AuthorityActions.syncNames();
		SettingsActions.sync();
	}

	componentWillUnmount(): void {
		PoliciesStore.removeChangeListener(this.onChange);
		ServicesStore.removeChangeListener(this.onChange);
		AuthoritiesStore.removeChangeListener(this.onChange);
		SettingsStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		let policies = PoliciesStore.policies;
		let selected: Selected = {};
		let curSelected = this.state.selected;
		let opened: Opened = {};
		let curOpened = this.state.opened;

		policies.forEach((policy: PolicyTypes.Policy): void => {
			if (curSelected[policy.id]) {
				selected[policy.id] = true;
			}
			if (curOpened[policy.id]) {
				opened[policy.id] = true;
			}
		});

		this.setState({
			...this.state,
			policies: policies,
			services: ServicesStore.servicesName,
			authorities: AuthoritiesStore.authoritiesName,
			providers: SettingsStore.settings ?
				SettingsStore.settings.auth_secondary_providers : [],
			filter: PoliciesStore.filter,
			selected: selected,
			opened: opened,
		});
	}

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		PolicyActions.removeMulti(
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
		let policiesDom: JSX.Element[] = [];

		this.state.policies.forEach((
				policy: PolicyTypes.PolicyRo): void => {
			policiesDom.push(<Policy
				key={policy.id}
				policy={policy}
				services={this.state.services}
				authorities={this.state.authorities}
				providers={this.state.providers}
				selected={!!this.state.selected[policy.id]}
				open={!!this.state.opened[policy.id]}
				onSelect={(shift: boolean): void => {
					let selected = {
						...this.state.selected,
					};

					if (shift) {
						let policies = this.state.policies;
						let start: number;
						let end: number;

						for (let i = 0; i < policies.length; i++) {
							let usr = policies[i];

							if (usr.id === policy.id) {
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
								selected[policies[i].id] = true;
							}

							this.setState({
								...this.state,
								lastSelected: policy.id,
								selected: selected,
							});

							return;
						}
					}

					if (selected[policy.id]) {
						delete selected[policy.id];
					} else {
						selected[policy.id] = true;
					}

					this.setState({
						...this.state,
						lastSelected: policy.id,
						selected: selected,
					});
				}}
				onOpen={(): void => {
					let opened = {
						...this.state.opened,
					};

					if (opened[policy.id]) {
						delete opened[policy.id];
					} else {
						opened[policy.id] = true;
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
			let inst = PoliciesStore.policy(instId);
			if (inst) {
				selectedNames.push(inst.name || instId);
			} else {
				selectedNames.push(instId);
			}
		}

		let newPolicyDom: JSX.Element;
		if (this.state.newOpened) {
			newPolicyDom = <PolicyNew
				services={this.state.services}
				authorities={this.state.authorities}
				providers={this.state.providers}
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
					<h2 style={css.heading}>Policies</h2>
					<div className="flex"/>
					<div style={css.buttons}>
						<button
							className={filterClass}
							style={css.button}
							type="button"
							onClick={(): void => {
								if (this.state.filter === null) {
									PolicyActions.filter({});
								} else {
									PolicyActions.filter(null);
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
							confirmMsg="Permanently delete the selected policies"
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
			<PoliciesFilter
				filter={this.state.filter}
				onFilter={(filter): void => {
					PolicyActions.filter(filter);
				}}
			/>
			<div style={css.itemsBox}>
				<div style={css.items}>
					{newPolicyDom}
					{policiesDom}
					<tr className="bp5-card bp5-row" style={css.placeholder}>
						<td colSpan={2} style={css.placeholder}/>
					</tr>
				</div>
			</div>
			<NonState
				hidden={!!policiesDom.length}
				iconClass="bp5-icon-ip-filter"
				title="No policies"
				description="Add a new policy to get started."
			/>
			<PoliciesPage
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
