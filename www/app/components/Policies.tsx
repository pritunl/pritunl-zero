/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as PolicyTypes from '../types/PolicyTypes';
import * as ServiceTypes from '../types/ServiceTypes';
import PoliciesStore from '../stores/PoliciesStore';
import ServicesStore from '../stores/ServicesStore';
import * as PolicyActions from '../actions/PolicyActions';
import * as ServiceActions from '../actions/ServiceActions';
import NonState from './NonState';
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
			<NonState
				hidden={!!policiesDom.length}
				iconClass="pt-icon-filter"
				title="No policies"
				description="Add a new policy to get started."
			/>
		</Page>;
	}
}
