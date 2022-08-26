/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';
import * as SuperAgent from 'superagent';
import * as WebAuthn from '@github/webauthn-json';
import * as DeviceTypes from '../types/DeviceTypes';
import DevicesStore from '../stores/DevicesStore';
import StateStore from '../stores/StateStore';
import * as DeviceActions from '../actions/DeviceActions';
import * as StateActions from '../actions/StateActions';
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
	sshDevice: string;
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
		textAlign: 'center',
	} as React.CSSProperties,
	bodyRelative: {
		padding: 0,
		textAlign: 'center',
		position: 'relative',
	} as React.CSSProperties,
	title: {
		margin: '10px 0 15px 0',
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
	centerButton: {
		margin: '15px auto 0 auto',
		display: 'block',
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
	group: {
		marginTop: '15px',
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
		top: '-26px',
		right: '-11px',
		width: '36px',
	} as React.CSSProperties,
};

export default class Devices extends React.Component<Props, State> {
	timeout: number;
	alertKey: string;

	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			devices: DevicesStore.devices,
			sshDevice: StateStore.sshDevice,
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
		StateStore.addChangeListener(this.onChange);
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
		StateStore.addChangeListener(this.onChange);

		if (this.timeout) {
			window.clearTimeout(this.timeout);
		}
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			devices: DevicesStore.devices,
			sshDevice: StateStore.sshDevice,
		});
	}

	wanRegister = (cred: any): void => {
		let loader = new Loader().loading();

		cred.device_type = 'webauthn';

		SuperAgent
			.post('/device/manage/register')
			.send(cred)
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

	smartCardRegistered = (): void => {
		let loader = new Loader().loading();

		SuperAgent
			.post('/device/manage/register')
			.send({
				device_type: 'smart_card',
				token: this.state.register.token,
				name: this.state.deviceName,
				ssh_public_key: this.state.sshDevice,
			})
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				StateActions.setSshDevice(null);

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

	onRegister = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});

		if (this.state.sshDevice) {
			this.smartCardRegistered();
		} else {
			WebAuthn.create(this.state.register.options).then((cred: any): void => {
				cred.name = this.state.deviceName;
				cred.token = this.state.register.token;
				this.wanRegister(cred);
			}).catch((err: any): void => {
				Alert.errorRes(err, 'Failed to register device');
				this.setState({
					...this.state,
					disabled: false,
				});
			});
		}
	}

	register(): JSX.Element {
		return <div>
			<div style={css.body}>
				<div className="bp3-non-ideal-state-visual bp3-non-ideal-state-icon">
					<span className="bp3-icon bp3-icon-key"/>
				</div>
				<h4 style={css.title}>
					Register Security Device
				</h4>
				<span style={css.description}>
					Enter a name for your new security device.
				</span>
				<div
					className="bp3-control-group"
					style={css.group}
				>
					<div style={css.inputBox}>
						<input
							className="bp3-input"
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
									this.onRegister();
								}
							}}
						/>
					</div>
					<div>
						<button
							className="bp3-button bp3-intent-success bp3-icon-add"
							disabled={this.state.disabled}
							onClick={this.onRegister}
						>Add Device</button>
					</div>
				</div>
			</div>
		</div>;
	}

	wanRespond = (resp: any): void => {
		Alert.dismiss(this.alertKey);

		let loader = new Loader().loading();

		let deviceType = 'webauthn';
		if (this.state.sshDevice) {
			deviceType = 'smart_card';
		}

		resp.device_type = deviceType;
		resp.token = this.state.secondary.token;

		SuperAgent
			.post('/device/manage/respond')
			.send(resp)
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					Alert.errorRes(res, 'Failed to complete device authentication');
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
					disabled: false,
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
			.get('/device/manage/request')
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

				WebAuthn.get(res.body).then((cred: any): void => {
					this.wanRespond(cred);
				}).catch((err: any): void => {
					Alert.errorRes(err, 'Failed to authenticate device');
					this.setState({
						...this.state,
						disabled: false,
						secondary: null,
						register: false,
					});
				});
			});
	}

	device(): JSX.Element {
		return <div>
			<div style={css.body}>
				<div className="bp3-non-ideal-state-visual bp3-non-ideal-state-icon">
					<span className="bp3-icon bp3-icon-key"/>
				</div>
				<h4 style={css.title}>
					{this.state.secondary.label}
				</h4>
				<span style={css.description}>
					A current security device is required to add new devices
				</span>
				<button
					className="bp3-button bp3-intent-success bp3-icon-id-number"
					disabled={this.state.disabled}
					onClick={this.deviceSign}
					style={css.centerButton}
				>Authenticate</button>
			</div>
		</div>;
	}

	smartCard(): JSX.Element {
		let sshDevice = this.state.sshDevice;
		sshDevice = sshDevice.replace(/-/g, '+').replace(/_/g, '/');
		while (sshDevice.length % 4) {
			sshDevice += '=';
		}
		sshDevice = atob(sshDevice);

		let cardSplit = sshDevice.split('cardno:');
		let cardSerial = 'unknown';
		if (cardSplit.length > 1) {
			cardSerial = cardSplit[1];
		}

		return <div>
			<div style={css.body}>
				<div className="bp3-non-ideal-state-visual bp3-non-ideal-state-icon">
					<span className="bp3-icon bp3-icon-sim-card"/>
				</div>
				<h4 style={css.title}>
					Register Smart Card
				</h4>
				<span style={css.description}>
					Registering Smart Card <b>{cardSerial}</b>
				</span>
				<div
					className="layout horizontal center-justified"
					style={css.buttons}
				>
					<button
						className="bp3-button bp3-large bp3-intent-danger bp3-icon-cross"
						style={css.button}
						disabled={this.state.disabled}
						onClick={(): void => {
							StateActions.setSshDevice(null);
						}}
					>Cancel</button>
					<button
						className="bp3-button bp3-large bp3-intent-success bp3-icon-tick"
						style={css.button}
						disabled={this.state.disabled}
						onClick={this.initRegister}
					>Continue</button>
				</div>
			</div>
		</div>;
	}

	secondarySubmit(factor: string): void {
		let passcode = '';
		if (factor === 'passcode') {
			passcode = this.state.passcode;
		}

		let deviceType = 'webauthn';
		if (this.state.sshDevice) {
			deviceType = 'smart_card';
		}

		SuperAgent
			.put('/device/manage/secondary')
			.send({
				device_type: deviceType,
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
					StateActions.setSshDevice(null);
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
						disabled: false,
						register: res.body,
					});
				}
			});
	}

	secondary(): JSX.Element {
		return <div>
			<div style={css.body}>
				<div className="bp3-non-ideal-state-visual bp3-non-ideal-state-icon">
					<span className="bp3-icon bp3-icon-key"/>
				</div>
				<h4 style={css.title}>
					{this.state.secondary.label}
				</h4>
				<span style={css.description}>
					Secondary authentication required
				</span>
			</div>
			<div className="layout vertical center-justified" style={css.buttons}>
				<button
					className="bp3-button"
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
					className="bp3-button"
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
					className="bp3-button"
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
					className="bp3-input"
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
					className="bp3-button"
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

	initRegister = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});

		Alert.dismiss(this.alertKey);
		let loader = new Loader().loading();

		let deviceType = 'webauthn';
		if (this.state.sshDevice) {
			deviceType = 'smart_card';
		}

		SuperAgent
			.get('/device/manage/register')
			.query({
				device_type: deviceType,
			})
			.set('Accept', 'application/json')
			.set('Csrf-Token', Csrf.token)
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					StateActions.setSshDevice(null);
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
		} else if (this.state.sshDevice) {
			return this.smartCard();
		}

		let devicesDom: JSX.Element[] = [];

		this.state.devices.forEach((device: DeviceTypes.DeviceRo): void => {
			devicesDom.push(<Device
				key={device.id}
				device={device}
			/>);
		});

		return <div style={css.bodyRelative}>
			<button
				className="bp3-button bp3-minimal bp3-intent-danger"
				style={css.close}
				onClick={this.props.onClose}
			>
				<Blueprint.Icon icon="cross" iconSize={26}/>
			</button>
			<h4 style={css.title}>
				Security Devices
			</h4>
			<div
				className="layout vertical center-justified wrap"
				style={css.buttons}
			>
				{devicesDom}
				<div
					className="bp3-non-ideal-state"
					style={css.state}
					hidden={!!devicesDom.length || !this.state.initialized}
				>
					<div
						className="bp3-non-ideal-state-visual bp3-non-ideal-state-icon"
						style={css.stateIcon}
					>
						<Blueprint.Icon icon="id-number" iconSize={80}/>
					</div>
					<h4 className="bp3-non-ideal-state-title">
						No devices registered
					</h4>
				</div>
				<button
					className="bp3-button bp3-intent-success bp3-icon-add"
					disabled={this.state.disabled}
					onClick={this.initRegister}
				>Add WebAuthn Device</button>
			</div>
		</div>;
	}
}
