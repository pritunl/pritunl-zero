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
import NonState from './NonState';
import Policy from './Policy';
import Page from './Page';
import PageHeader from './PageHeader';

interface State {
	policies: PolicyTypes.PoliciesRo;
	services: ServiceTypes.ServicesRo;
	authorities: AuthorityTypes.AuthoritiesRo;
	providers: SettingsTypes.SecondaryProviders;
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
			disabled: false,
		};
	}

	componentDidMount(): void {
		PoliciesStore.addChangeListener(this.onChange);
		ServicesStore.addChangeListener(this.onChange);
		AuthoritiesStore.addChangeListener(this.onChange);
		SettingsStore.addChangeListener(this.onChange);
		PolicyActions.sync();
		ServiceActions.syncNames();
		AuthorityActions.sync();
		SettingsActions.sync();
	}

	componentWillUnmount(): void {
		PoliciesStore.removeChangeListener(this.onChange);
		ServicesStore.removeChangeListener(this.onChange);
		AuthoritiesStore.removeChangeListener(this.onChange);
		SettingsStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			policies: PoliciesStore.policies,
			services: ServicesStore.servicesName,
			authorities: AuthoritiesStore.authorities,
			providers: SettingsStore.settings ?
				SettingsStore.settings.auth_secondary_providers : [],
		});
	}

	render(): JSX.Element {
		let policiesDom: JSX.Element[] = [];

		this.state.policies.forEach((policy: PolicyTypes.PolicyRo): void => {
			policiesDom.push(<Policy
				key={policy.id}
				policy={policy}
				services={this.state.services}
				authorities={this.state.authorities}
				providers={this.state.providers}
			/>);
		});

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Policies</h2>
					<div className="flex"/>
					<div style={css.buttons}>
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
								PolicyActions.create(null).then((): void => {
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
			<div>
				{policiesDom}
			</div>
			<NonState
				hidden={!!policiesDom.length}
				iconClass="bp5-icon-filter"
				title="No policies"
				description="Add a new policy to get started."
			/>
		</Page>;
	}
}
