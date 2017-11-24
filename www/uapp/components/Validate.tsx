/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SuperAgent from 'superagent';
import * as Csrf from '../Csrf';
import * as Alert from '../Alert';

interface Props {
	token: string;
}

interface State {
	disabled: boolean;
}

const css = {
	body: {
		padding: '0 10px',
	} as React.CSSProperties,
	description: {
		opacity: 0.7,
	} as React.CSSProperties,
	buttons: {
		marginTop: '15px',
	} as React.CSSProperties,
	button: {
		margin: '5px',
		width: '116px',
	} as React.CSSProperties,
};

export default class Validate extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
		};
	}

	render(): JSX.Element {
		return <div>
			<div className="pt-non-ideal-state" style={css.body}>
				<div className="pt-non-ideal-state-visual pt-non-ideal-state-icon">
					<span className="pt-icon pt-icon-endorsed"/>
				</div>
				<h4 className="pt-non-ideal-state-title">Validate SSH Key</h4>
				<span style={css.description}>If you did not initiate this validation deny the request and report the incident to an administrator</span>
			</div>
			<div className="layout horizontal center-justified" style={css.buttons}>
				<button
					className="pt-button pt-large pt-intent-success pt-icon-add"
					style={css.button}
					type="button"
					disabled={this.state.disabled}
					onClick={(): void => {
						this.setState({
							...this.state,
							disabled: true,
						});

						SuperAgent
							.put('/ssh/validate/' + this.props.token)
							.set('Accept', 'application/json')
							.set('Csrf-Token', Csrf.token)
							.end((err: any, res: SuperAgent.Response): void => {
								this.setState({
									...this.state,
									disabled: false,
								});

								if (err) {
									Alert.errorRes(res, 'Failed to validate SSH key');
									return;
								}
							});
					}}
				>
					Approve
				</button>
				<button
					className="pt-button pt-large pt-intent-danger pt-icon-delete"
					style={css.button}
					type="button"
					disabled={this.state.disabled}
					onClick={(): void => {
					}}
				>
					Deny
				</button>
			</div>
		</div>;
	}
}
