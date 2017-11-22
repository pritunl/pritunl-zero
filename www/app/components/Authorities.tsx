/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as AuthorityTypes from '../types/AuthorityTypes';
import AuthoritiesStore from '../stores/AuthoritiesStore';
import ServicesStore from '../stores/ServicesStore';
import * as AuthorityActions from '../actions/AuthorityActions';
import NonState from './NonState';
import Authority from './Authority';
import Page from './Page';
import PageHeader from './PageHeader';

interface State {
	authorities: AuthorityTypes.AuthoritiesRo;
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

export default class Authorities extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			authorities: AuthoritiesStore.authorities,
			disabled: false,
		};
	}

	componentDidMount(): void {
		AuthoritiesStore.addChangeListener(this.onChange);
		AuthorityActions.sync();
	}

	componentWillUnmount(): void {
		AuthoritiesStore.removeChangeListener(this.onChange);
		ServicesStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			authorities: AuthoritiesStore.authorities,
		});
	}

	render(): JSX.Element {
		let authoritiesDom: JSX.Element[] = [];

		this.state.authorities.forEach((
				authority: AuthorityTypes.AuthorityRo): void => {
			authoritiesDom.push(<Authority
				key={authority.id}
				authority={authority}
			/>);
		});

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Authorities</h2>
					<div className="flex"/>
					<div>
						<button
							className="pt-button pt-intent-success pt-icon-add"
							style={css.button}
							type="button"
							onClick={(): void => {
								AuthorityActions.create(null);
							}}
						>New</button>
					</div>
				</div>
			</PageHeader>
			<div>
				{authoritiesDom}
			</div>
			<NonState
				hidden={!!authoritiesDom.length}
				iconClass="pt-icon-filter"
				title="No authorities"
				description="Add a new authority to get started."
			/>
		</Page>;
	}
}
