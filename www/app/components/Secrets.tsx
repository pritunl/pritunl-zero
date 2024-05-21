/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SecretTypes from '../types/SecretTypes';
import SecretsStore from '../stores/SecretsStore';
import * as SecretActions from '../actions/SecretActions';
import NonState from './NonState';
import Secret from './Secret';
import Page from './Page';
import PageHeader from './PageHeader';

interface State {
	secrets: SecretTypes.SecretsRo;
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
	noCerts: {
		height: 'auto',
	} as React.CSSProperties,
};

export default class Secrets extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			secrets: SecretsStore.secrets,
			disabled: false,
		};
	}

	componentDidMount(): void {
		SecretsStore.addChangeListener(this.onChange);
		SecretActions.sync();
	}

	componentWillUnmount(): void {
		SecretsStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			secrets: SecretsStore.secrets,
		});
	}

	render(): JSX.Element {
		let certsDom: JSX.Element[] = [];

		this.state.secrets.forEach((
				cert: SecretTypes.SecretRo): void => {
			certsDom.push(<Secret
				key={cert.id}
				secret={cert}
			/>);
		});

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Secrets</h2>
					<div className="flex"/>
					<div style={css.buttons}>
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
								SecretActions.create(null).then((): void => {
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
				{certsDom}
			</div>
			<NonState
				hidden={!!certsDom.length}
				iconClass="bp3-icon-key"
				title="No secrets"
				description="Add a new secret to get started."
			/>
		</Page>;
	}
}
