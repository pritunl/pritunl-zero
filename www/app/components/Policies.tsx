/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as PolicyTypes from '../types/PolicyTypes';
import * as ServiceTypes from '../types/ServiceTypes';
import PoliciesStore from '../stores/PoliciesStore';
import ServicesStore from '../stores/ServicesStore';
import * as PolicyActions from '../actions/PolicyActions';
import * as ServiceActions from '../actions/ServiceActions';
import Policy from './Policy';
import Page from './Page';
import PageHeader from './PageHeader';

interface State {
	policies: PolicyTypes.PoliciesRo;
	services: ServiceTypes.ServicesRo;
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
		margin: '10px 0 0 0',
	} as React.CSSProperties,
	noPolicies: {
		height: 'auto',
	} as React.CSSProperties,
};

export default class Policies extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			policies: PoliciesStore.policies,
			services: ServicesStore.services,
			disabled: false,
		};
	}

	componentDidMount(): void {
		PoliciesStore.addChangeListener(this.onChange);
		ServicesStore.addChangeListener(this.onChange);
		PolicyActions.sync();
		ServiceActions.sync();
	}

	componentWillUnmount(): void {
		PoliciesStore.removeChangeListener(this.onChange);
		ServicesStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			policies: PoliciesStore.policies,
			services: ServicesStore.services,
		});
	}

	render(): JSX.Element {
		let policiesDom: JSX.Element[] = [];

		this.state.policies.forEach((policy: PolicyTypes.PolicyRo): void => {
			policiesDom.push(<Policy
				key={policy.id}
				policy={policy}
				services={this.state.services}
			/>);
		});

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Policies</h2>
					<div className="flex"/>
					<div>
						<button
							className="pt-button pt-intent-success pt-icon-add"
							style={css.button}
							type="button"
							onClick={(): void => {
								PolicyActions.create(null);
							}}
						>New</button>
					</div>
				</div>
			</PageHeader>
			<div>
				{policiesDom}
			</div>
			<div
				className="pt-non-ideal-state"
				style={css.noPolicies}
				hidden={!!policiesDom.length}
			>
				<div className="pt-non-ideal-state-visual pt-non-ideal-state-icon">
					<span className="pt-icon pt-icon-cloud"/>
				</div>
				<h4 className="pt-non-ideal-state-title">No policies</h4>
				<div className="pt-non-ideal-state-description">
					Add a new policy to get started.
				</div>
			</div>
		</Page>;
	}
}
