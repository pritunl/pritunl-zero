/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SuperAgent from 'superagent';
import * as Csrf from '../Csrf';
import * as Alert from '../Alert';
import Session from './Session';
import KeybaseStore from "../stores/KeybaseStore";
import * as KeybaseTypes from "../types/KeybaseTypes";
import * as KeybaseActions from "../actions/KeybaseActions";

interface Props {
	token: string;
	signature: string;
}

interface State {
	disabled: boolean;
	answered: boolean;
	info: KeybaseTypes.InfoRo;
}

const css = {
	body: {
		padding: 0,
	} as React.CSSProperties,
	description: {
		opacity: 0.7,
		padding: '0 10px',
	} as React.CSSProperties,
	buttons: {
		marginTop: '15px',
	} as React.CSSProperties,
	button: {
		margin: '5px',
		width: '116px',
	} as React.CSSProperties,
	picture: {
		width: '100%',
		maxWidth: '140px',
		borderRadius: '50%',
	} as React.CSSProperties,
	info: {
		margin: 0,
		textAlign: 'center',
	} as React.CSSProperties,
	value: {
		opacity: 0.7,
	} as React.CSSProperties,
};

export default class Validate extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			answered: false,
			info: KeybaseStore.info,
		};
	}

	componentDidMount(): void {
		KeybaseStore.addChangeListener(this.onChange);
		KeybaseActions.load(this.props.token);
	}

	componentWillUnmount(): void {
		KeybaseStore.removeChangeListener(this.onChange);
		KeybaseActions.unload();
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			info: KeybaseStore.info,
		});
	}

	render(): JSX.Element {
		let info = this.state.info || {};

		if (this.state.answered) {
			return <Session/>;
		}

		return <div>
			<div className="pt-non-ideal-state" style={css.body}>
				<div className="pt-non-ideal-state-visual pt-non-ideal-state-icon">
					<span className="pt-icon pt-icon-endorsed"/>
				</div>
				<h4 className="pt-non-ideal-state-title">Associate Keybase Account</h4>
				<div>
					<img hidden={!info.picture} style={css.picture} src={info.picture}/>
					<div hidden={!info.username} style={css.info}>
						Keybase: <span style={css.value}>{info.username}</span>
					</div>
					<div hidden={!info.twitter} style={css.info}>
						Twitter: <span style={css.value}>{info.twitter}</span>
					</div>
					<div hidden={!info.github} style={css.info}>
						Github: <span style={css.value}>{info.github}</span>
					</div>
				</div>
				<span style={css.description}>
					If you did not initiate this association deny the request and
					report the incident to an administrator
				</span>
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
							.put('/keybase/validate')
							.set('Accept', 'application/json')
							.set('Csrf-Token', Csrf.token)
							.send({
								token: this.props.token,
								signature: this.props.signature,
							})
							.end((err: any, res: SuperAgent.Response): void => {
								this.setState({
									...this.state,
									disabled: false,
								});

								if (res.status === 404) {
									Alert.error('Keybase association request has expired', 0);
								} else if (err) {
									Alert.errorRes(res, 'Failed to associate keybase', 0);
								} else {
									Alert.success('Successfully associated keybase', 0);
								}

								this.setState({
									...this.state,
									answered: true,
								});

								window.history.replaceState(
									null, null, window.location.pathname);
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
						this.setState({
							...this.state,
							disabled: true,
						});

						SuperAgent
							.delete('/keybase/validate')
							.set('Accept', 'application/json')
							.set('Csrf-Token', Csrf.token)
							.send({
								token: this.props.token,
								signature: this.props.signature,
							})
							.end((err: any, res: SuperAgent.Response): void => {
								this.setState({
									...this.state,
									disabled: false,
								});

								if (res.status === 404) {
									Alert.error('Keybase association request has expired', 0);
								} else if (err) {
									Alert.errorRes(res,
										'Failed to deny keybase association', 0);
									return;
								} else {
									Alert.error('Successfully denied keybase association. ' +
										'Report this incident to an administrator.', 0);
								}

								this.setState({
									...this.state,
									answered: true,
								});

								window.history.replaceState(
									null, null, window.location.pathname);
							});
					}}
				>
					Deny
				</button>
			</div>
		</div>;
	}
}
