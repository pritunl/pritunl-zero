/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';
import * as SuperAgent from 'superagent';
import * as DeviceTypes from '../types/DeviceTypes';
import DevicesStore from '../stores/DevicesStore';
import * as DeviceActions from '../actions/DeviceActions';
import * as Alert from "../Alert";
import * as Csrf from "../Csrf";
import Device from './Device';
import * as Constants from "../Constants";
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
	onClose: () => void;
}

interface State {
	devices: DeviceTypes.DevicesRo;
	deviceName: string;
	disabled: boolean;
	passcode: string;
	secondary: Secondary;
	secondaryState: SecondaryState;
	register: any;
	initialized: boolean;
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
	state: {
		marginBottom: '5px',
	} as React.CSSProperties,
	stateIcon: {
		marginBottom: '10px',
	} as React.CSSProperties,
	title: {
		textAlign: 'center',
	} as React.CSSProperties,
	group: {
		width: '100%',
	} as React.CSSProperties,
	input: {
		width: '100%',
	} as React.CSSProperties,
	inputBox: {
		flex: '1',
	} as React.CSSProperties,
	close: {
		position: 'absolute',
		top: '7px',
		right: '7px',
		width: '36px',
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

export default class Devices extends React.Component<Props, State> {
	timeout: number;
	alertKey: string;

	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			devices: DevicesStore.devices,
			deviceName: '',
			disabled: false,
			initialized: false,
			passcode: '',
			secondary: null,
			secondaryState: null,
			register: null,
		};
	}

	componentDidMount(): void {
		DevicesStore.addChangeListener(this.onChange);
		DeviceActions.sync();

		this.timeout = window.setTimeout((): void => {
			this.setState({
				...this.state,
				initialized: true,
			});
		}, Constants.loadDelay);
	}

	componentWillUnmount(): void {
		DevicesStore.removeChangeListener(this.onChange);

		if (this.timeout) {
			window.clearTimeout(this.timeout);
		}
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			devices: DevicesStore.devices,
		});
	}

	u2fRegistered = (resp: any): void => {
		Alert.dismiss(this.alertKey);

		if (resp.errorCode) {
			let errorMsg = 'U2F error code ' + resp.errorCode;
			let u2fMsg = u2fErrorCodes[resp.errorCode as number];
			if (u2fMsg) {
				errorMsg += ': ' + u2fMsg;
			}
			Alert.error(errorMsg);

			this.setState({
				...this.state,
				disabled: false,
				secondary: null,
				register: null,
			});

			return
		}

		let loader = new Loader().loading();

		SuperAgent
			.post('/device/u2f/register')
			.send({
				token: this.state.register.token,
				name: this.state.deviceName,
				response: resp,
			})
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				this.setState({
					...this.state,
					disabled: false,
					deviceName: '',
					secondary: null,
					register: null,
				});

				if (err) {
					Alert.errorRes(res, 'Failed to register device');
					return;
				}

				DeviceActions.sync();

				this.alertKey = Alert.success('Successfully registered device');
			});
	}

	registerSign = (): void => {
		this.setState({
			disabled: true,
		});

		this.alertKey = Alert.info(
			'Insert your security key and tap the button', 30000);

		(window as any).u2f.register(this.state.register.request.appId,
			this.state.register.request.registerRequests,
			this.state.register.request.registeredKeys,
			this.u2fRegistered, 30);
	}

	register(): JSX.Element {
		return <div>
			<div className="pt-non-ideal-state" style={css.body}>
				<div className="pt-non-ideal-state-visual pt-non-ideal-state-icon">
					<span className="pt-icon pt-icon-key"/>
				</div>
				<h4 className="pt-non-ideal-state-title">
					Register U2F Device
				</h4>
				<span style={css.description}>
					Enter a name for your new security device.
				</span>
				<div
					className="pt-control-group"
					style={css.group}
				>
					<div style={css.inputBox}>
						<input
							className="pt-input"
							style={css.input}
							type="text"
							placeholder="Device name"
							value={this.state.deviceName}
							onChange={(evt): void => {
								this.setState({
									...this.state,
									deviceName: evt.target.value,
								});
							}}
							onKeyPress={(evt): void => {
								if (evt.key === 'Enter') {
									this.registerSign();
								}
							}}
						/>
					</div>
					<div>
						<button
							className="pt-button pt-intent-success pt-icon-add"
							disabled={this.state.disabled}
							onClick={this.registerSign}
						>Add Device</button>
					</div>
				</div>
			</div>
		</div>;
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

			this.setState({
				...this.state,
				disabled: false,
				secondary: null,
				register: false,
			});

			return
		}

		let loader = new Loader().loading();

		SuperAgent
			.post('/device/u2f/sign')
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
					register: res.body,
				});
			});
	}

	deviceSign = (): void => {
		let loader = new Loader().loading();

		this.setState({
			...this.state,
			disabled: true,
		});

		SuperAgent
			.get('/device/u2f/sign')
			.query({
				token: this.state.secondary.token,
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
					A current security device is required to add new devices
				</span>
				<button
					className="pt-button pt-intent-success pt-icon-id-number"
					disabled={this.state.disabled}
					onClick={this.deviceSign}
				>Authenticate</button>
			</div>
		</div>;
	}

	secondarySubmit(factor: string): void {
		let passcode = '';
		if (factor === 'passcode') {
			passcode = this.state.passcode;
		}

		SuperAgent
			.put('/device/u2f/secondary')
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
					Alert.error('Device registration request has expired', 0);
					this.setState({
						...this.state,
						disabled: false,
						secondary: null,
						register: null,
					});
				} else if (err) {
					Alert.errorRes(res, 'Failed to register device', 0);
				} else if (res.status === 206 && factor === 'sms') {
					Alert.info('Text message sent', 0);
				} else {
					this.setState({
						...this.state,
						register: res.body,
					});
				}
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

	onRegister = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});

		Alert.dismiss(this.alertKey);
		let loader = new Loader().loading();

		SuperAgent
			.get('/device/u2f/register')
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					Alert.errorRes(res, 'Failed to request device registration');
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
					register: res.body,
					disabled: false,
				});
			});
	}

	render(): JSX.Element {
		if (this.state.register) {
			return this.register();
		} else if (this.state.secondary) {
			if (this.state.secondary.device) {
				return this.device();
			} else {
				return this.secondary();
			}
		}

		let devicesDom: JSX.Element[] = [];

		this.state.devices.forEach((device: DeviceTypes.DeviceRo): void => {
			devicesDom.push(<Device
				key={device.id}
				device={device}
			/>);
		});

		return <div>
			<button
				className="pt-button pt-minimal pt-intent-danger"
				style={css.close}
				onClick={this.props.onClose}
			>
				<Blueprint.Icon icon="cross" iconSize={26}/>
			</button>
			<h4 style={css.title}>
				U2F Devices
			</h4>
			<div
				className="layout vertical center-justified wrap"
				style={css.buttons}
			>
				{devicesDom}
				<div
					className="pt-non-ideal-state"
					style={css.state}
					hidden={!!devicesDom.length || !this.state.initialized}
				>
					<div
						className="pt-non-ideal-state-visual pt-non-ideal-state-icon"
						style={css.stateIcon}
					>
						<Blueprint.Icon icon="id-number" iconSize={80}/>
					</div>
					<h4 className="pt-non-ideal-state-title">
						No devices registered
					</h4>
				</div>
				<button
					className="pt-button pt-intent-success pt-icon-add"
					disabled={this.state.disabled}
					onClick={this.onRegister}
				>Add Device</button>
			</div>
		</div>;
	}
}
