/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SuperAgent from 'superagent';
import * as Csrf from '../Csrf';
import * as Alert from '../Alert';
import Session from './Session';
import Loader from "../Loader";

interface Secondary {
	token: string;
	label: string;
	push: boolean;
	phone: boolean;
	passcode: boolean;
	sms: boolean;
	device: boolean;
	device_register: boolean;
}

interface SecondaryState {
	push: boolean;
	phone: boolean;
	passcode: boolean;
	sms: boolean;
}

interface Props {
	token: string;
}

interface State {
	disabled: boolean;
	answered: boolean;
	passcode: string;
	secondary: Secondary;
	secondaryState: SecondaryState;
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
	secondaryButton: {
		margin: '5px auto',
		padding: '8px 15px',
		width: '75%',
	} as React.CSSProperties,
	secondaryInput: {
		margin: '5px auto',
		width: '75%',
	} as React.CSSProperties,
};

const u2fErrorCodes: {[index: number]: string} = {
	0: 'ok',
	1: 'other',
	2: 'bad request',
	3: 'configuration unsupported',
	4: 'device ineligible',
	5: 'timed out',
};

export default class Validate extends React.Component<Props, State> {
	alertKey: string;

	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			answered: false,
			passcode: '',
			secondary: null,
			secondaryState: null,
		};
	}

	u2fSigned = (resp: any): void => {
		Alert.dismiss(this.alertKey);

		if (resp.errorCode) {
			let errorMsg = 'U2F error code ' + resp.errorCode;
			let u2fMsg = u2fErrorCodes[resp.errorCode as number];
			if (u2fMsg) {
				errorMsg += ': ' + u2fMsg;
			}
			Alert.error(errorMsg);

			return
		}

		let loader = new Loader().loading();

		SuperAgent
			.post('/ssh/u2f/sign')
			.send({
				token: this.state.secondary.token,
				response: resp,
			})
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					Alert.errorRes(res, 'Failed to complete device sign');
					return;
				}

				if (res.status === 201) {
					this.setState({
						...this.state,
						secondary: res.body,
						secondaryState: {
							push: true,
							phone: true,
							passcode: true,
							sms: true,
						},
						disabled: false,
					});

					return;
				}

				this.setState({
					...this.state,
					answered: true,
					secondary: null,
				});

				window.history.replaceState(
					null, null, window.location.pathname);

				Alert.success('Successfully approved SSH key', 0);
			});
	}

	deviceSign(token: string): void {
		let loader = new Loader().loading();

		SuperAgent
			.get('/ssh/u2f/sign')
			.query({
				token: token,
			})
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					Alert.errorRes(res, 'Failed to request device sign');
					return;
				}

				this.alertKey = Alert.info(
					'Insert your security key and tap the button', 30000);

				(window as any).u2f.sign(res.body.appId,
					res.body.challenge, res.body.registeredKeys,
					this.u2fSigned, 30);
			});
	}

	device(): JSX.Element {
		return <div>
			<div className="pt-non-ideal-state" style={css.body}>
				<div className="pt-non-ideal-state-visual pt-non-ideal-state-icon">
					<span className="pt-icon pt-icon-key"/>
				</div>
				<h4 className="pt-non-ideal-state-title">
					{this.state.secondary.label}
				</h4>
				<span style={css.description}>
					Insert your security key and tap the button
				</span>
			</div>
		</div>;
	}

	secondarySubmit(factor: string): void {
		let passcode = '';
		if (factor === 'passcode') {
			passcode = this.state.passcode;
		}

		SuperAgent
			.put('/ssh/secondary')
			.send({
				token: this.state.secondary.token,
				factor: factor,
				passcode: passcode
			})
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				this.setState({
					...this.state,
					passcode: '',
					secondaryState: {
						...this.state.secondaryState,
						passcode: true,
					},
				});

				if (res && res.status === 404) {
					Alert.error('SSH verification request has expired', 0);
				} else if (err) {
					Alert.errorRes(res, 'Failed to approve SSH key', 0);
					return;
				} else if (res.status === 206 && factor === 'sms') {
					Alert.info('Text message sent', 0);
					return;
				} else {
					Alert.success('Successfully approved SSH key', 0);
				}

				this.setState({
					...this.state,
					answered: true,
					secondary: null,
				});

				window.history.replaceState(
					null, null, window.location.pathname);
			});
	}

	secondary(): JSX.Element {
		return <div>
			<div className="pt-non-ideal-state" style={css.body}>
				<div className="pt-non-ideal-state-visual pt-non-ideal-state-icon">
					<span className="pt-icon pt-icon-key"/>
				</div>
				<h4 className="pt-non-ideal-state-title">
					{this.state.secondary.label}
				</h4>
				<span style={css.description}>
					Secondary authentication required
				</span>
			</div>
			<div className="layout vertical center-justified" style={css.buttons}>
				<button
					className="pt-button"
					style={css.secondaryButton}
					type="button"
					hidden={!this.state.secondary.push}
					disabled={!this.state.secondaryState.push}
					onClick={(): void => {
						this.setState({
							...this.state,
							secondaryState: {
								...this.state.secondaryState,
								push: false,
							},
						});
						this.secondarySubmit('push');
					}}
				>
					Push
				</button>
				<button
					className="pt-button"
					style={css.secondaryButton}
					type="button"
					hidden={!this.state.secondary.phone}
					disabled={!this.state.secondaryState.phone}
					onClick={(): void => {
						this.setState({
							...this.state,
							secondaryState: {
								...this.state.secondaryState,
								phone: false,
							},
						});
						this.secondarySubmit('phone');
					}}
				>
					Call Me
				</button>
				<button
					className="pt-button"
					style={css.secondaryButton}
					type="button"
					hidden={!this.state.secondary.sms}
					disabled={!this.state.secondaryState.sms}
					onClick={(): void => {
						this.setState({
							...this.state,
							secondaryState: {
								...this.state.secondaryState,
								sms: false,
							},
						});
						this.secondarySubmit('sms');
					}}
				>
					Text Me
				</button>
				<input
					className="pt-input"
					style={css.secondaryInput}
					hidden={!this.state.secondary.passcode}
					disabled={!this.state.secondaryState.passcode}
					type="text"
					autoCapitalize="off"
					spellCheck={false}
					placeholder="Passcode"
					value={this.state.passcode || ''}
					onChange={(evt): void => {
						this.setState({
							...this.state,
							passcode: evt.target.value,
						});
					}}
					onKeyPress={(evt): void => {
						if (evt.key === 'Enter') {
							this.setState({
								...this.state,
								secondaryState: {
									...this.state.secondaryState,
									passcode: false,
								},
							});
							this.secondarySubmit('passcode');
						}
					}}
				/>
				<button
					className="pt-button"
					style={css.secondaryButton}
					type="button"
					hidden={!this.state.secondary.passcode}
					disabled={!this.state.secondaryState.passcode}
					onClick={(): void => {
						this.setState({
							...this.state,
							secondaryState: {
								...this.state.secondaryState,
								passcode: false,
							},
						});
						this.secondarySubmit('passcode');
					}}
				>
					Submit
				</button>
			</div>
		</div>;
	}

	render(): JSX.Element {
		if (this.state.answered) {
			return <Session/>;
		}

		if (this.state.secondary) {
			if (this.state.secondary.device) {
				return this.device();
			}
			return this.secondary();
		}

		return <div>
			<div className="pt-non-ideal-state" style={css.body}>
				<div className="pt-non-ideal-state-visual pt-non-ideal-state-icon">
					<span className="pt-icon pt-icon-endorsed"/>
				</div>
				<h4 className="pt-non-ideal-state-title">Validate SSH Key</h4>
				<span style={css.description}>
					If you did not initiate this validation deny the request and
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
							.put('/ssh/validate/' + this.props.token)
							.set('Accept', 'application/json')
							.set('Csrf-Token', Csrf.token)
							.end((err: any, res: SuperAgent.Response): void => {
								this.setState({
									...this.state,
									disabled: false,
								});

								if (res && res.status === 404) {
									Alert.error('SSH verification request has expired', 0);
								} else if (err) {
									Alert.errorRes(res, 'Failed to approve SSH key', 0);
								} else if (res.status === 201) {

									this.setState({
										...this.state,
										secondary: res.body,
										secondaryState: {
											push: true,
											phone: true,
											passcode: true,
											sms: true,
										},
										disabled: false,
									});

									if (res.body.device) {
										this.deviceSign(res.body.token);
									}

									return;
								} else {
									Alert.success('Successfully approved SSH key', 0);
								}

								this.setState({
									...this.state,
									answered: true,
									disabled: false,
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
							.delete('/ssh/validate/' + this.props.token)
							.set('Accept', 'application/json')
							.set('Csrf-Token', Csrf.token)
							.end((err: any, res: SuperAgent.Response): void => {
								this.setState({
									...this.state,
									disabled: false,
								});

								if (res.status === 404) {
									Alert.error('SSH verification request has expired', 0);
								} else if (err) {
									Alert.errorRes(res, 'Failed to deny SSH key', 0);
									return;
								} else {
									Alert.error('Successfully denied SSH key. Report ' +
										'this incident to an administrator.', 0);
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
