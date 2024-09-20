/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as DeviceTypes from '../types/DeviceTypes';
import * as MiscUtils from '../utils/MiscUtils';
import * as DeviceActions from '../actions/DeviceActions';
import * as PageInfos from './PageInfo';
import PageInfo from './PageInfo';
import ConfirmButton from './ConfirmButton';
import * as Alert from '../Alert';
import PageSwitch from "./PageSwitch";
import PageSave from "./PageSave";
import PageInput from "./PageInput";

interface Props {
	device: DeviceTypes.DeviceRo;
}

interface State {
	disabled: boolean;
	changed: boolean;
	message: string;
	device: DeviceTypes.Device;
}

const css = {
	card: {
		position: 'relative',
		padding: '10px',
		marginBottom: '5px',
	} as React.CSSProperties,
	info: {
		marginBottom: '-5px',
	} as React.CSSProperties,
	group: {
		flex: 1,
		minWidth: '250px',
	} as React.CSSProperties,
	inputGroup: {
		marginBottom: '11px',
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	remove: {
		position: 'absolute',
		top: '5px',
		right: '5px',
	} as React.CSSProperties,
	controlButton: {
		marginRight: '10px',
	} as React.CSSProperties,
	save: {
		paddingTop: '10px',
	} as React.CSSProperties,
};

export default class Device extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
			device: null,
		};
	}

	set(name: string, val: any): void {
		let device: any;

		if (this.state.changed) {
			device = {
				...this.state.device,
			};
		} else {
			device = {
				...this.props.device,
			};
		}

		device[name] = val;

		this.setState({
			...this.state,
			changed: true,
			device: device,
		});
	}

	toggleLevel(level: number) {
		let device: any;

		if (this.state.changed) {
			device = {
				...this.state.device,
			};
		} else {
			device = {
				...this.props.device,
			};
		}

		let levels: number[] = Object.assign([], (device.alert_levels || []));
		let index = levels.indexOf(level);

		if (index !== -1) {
			levels.splice(index, 1);
		} else {
			levels.push(level);
		}

		device.alert_levels = levels;

		this.setState({
			...this.state,
			changed: true,
			device: device,
		});
	}

	onTestAlert = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		DeviceActions.testAlert(this.props.device.id).then((): void => {
			Alert.success('Test alert sent');

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
	}

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		DeviceActions.commit(this.state.device).then((): void => {
			Alert.success('Device name updated');

			this.setState({
				...this.state,
				disabled: false,
				changed: false,
			});

			setTimeout((): void => {
				if (!this.state.changed) {
					this.setState({
						...this.state,
						changed: false,
						device: null,
					});
				}
			}, 1000);
		}).catch((): void => {
			this.setState({
				...this.state,
				disabled: false,
			});
		});
	}

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		DeviceActions.remove(this.props.device.id).then((): void => {
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
	}

	render(): JSX.Element {
		let device: DeviceTypes.Device = this.state.device ||
			this.props.device;

		let isPhone: boolean = this.props.device.mode === 'phone';

		let deviceType = 'Unknown';
		switch (device.type) {
			case 'webauthn':
				deviceType = 'WebAuthn';
				break;
			case 'u2f':
				deviceType = 'U2F';
				break;
			case 'smart_card':
				deviceType = 'Smart Card';
				break;
			case 'call':
				deviceType = 'Call';
				break;
			case 'message':
				deviceType = 'SMS';
				break;
		}

		let deviceMode = 'Unknown';
		switch (device.mode) {
			case 'secondary':
				deviceMode = 'Secondary';
				break;
			case 'ssh':
				deviceMode = 'SSH';
				break;
			case 'phone':
				deviceMode = 'Phone';
				break;
		}

		let deviceOther: PageInfos.Field;
		if (device.wan_rp_id) {
			deviceOther = {
				label: 'WebAuthn Domain',
				value: device.wan_rp_id,
			};
		} else if (device.type === 'smart_card') {
			deviceOther = {
				label: 'SSH Public Key',
				value: device.ssh_public_key,
			};
		} else if (device.type === 'call' || device.type === 'message') {
			deviceOther = {
				label: 'Phone Number',
				value: device.number,
			};
		}

		let alertIcon = 'bp5-icon-phone';
		if (device.type === 'message') {
			alertIcon = 'bp5-icon-mobile-phone';
		}

		let cardStyle = {
			...css.card,
		};
		if (device.disabled) {
			cardStyle.opacity = 0.6;
		}

		let fields1: PageInfos.Field[];
		let fields2: PageInfos.Field[];

		if (isPhone) {
			fields1 = [
				{
					label: 'ID',
					value: device.id || 'None',
				},
			];
			fields2 = [
				{
					label: 'Type',
					value: deviceType,
				},
				{
					label: 'Mode',
					value: deviceMode,
				},
				deviceOther,
				{
					label: 'Registered',
					value: MiscUtils.formatDate(device.timestamp) || 'Unknown',
				},
				{
					label: 'Last Active',
					value: MiscUtils.formatDate(device.last_active) || 'Unknown',
				},
			];
		} else {
			fields1 = [
				{
					label: 'ID',
					value: device.id || 'None',
				},
				{
					label: 'Type',
					value: deviceType,
				},
				deviceOther,
			];
			fields2 = [
				{
					label: 'Mode',
					value: deviceMode,
				},
				{
					label: 'Registered',
					value: MiscUtils.formatDate(device.timestamp) || 'Unknown',
				},
				{
					label: 'Last Active',
					value: MiscUtils.formatDate(device.last_active) || 'Unknown',
				},
			];
		}

		let testButton: JSX.Element;
		if (isPhone) {
			testButton = <ConfirmButton
				label="Send Test Alert"
				className={'bp5-intent-success ' + alertIcon}
				progressClassName="bp5-intent-success"
				style={css.controlButton}
				disabled={this.state.disabled}
				onConfirm={(): void => {
					this.onTestAlert();
				}}
			/>;
		}

		return <div
			className="bp5-card"
			style={cardStyle}
		>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<div style={css.remove}>
						<ConfirmButton
							className="bp5-minimal bp5-intent-danger bp5-icon-trash"
							progressClassName="bp5-intent-danger"
							confirmMsg="Confirm device remove"
							disabled={this.state.disabled}
							onConfirm={this.onDelete}
						/>
					</div>
					<PageInput
						label="Device Name"
						help="Name of device."
						type="text"
						placeholder="Enter name"
						disabled={this.state.disabled}
						value={device.name}
						onChange={(val): void => {
							this.set('name', val);
						}}
					/>
					<PageSwitch
						label="Low alerts"
						help="Recieve low level alerts on this device."
						hidden={!isPhone}
						disabled={this.state.disabled}
						checked={(device.alert_levels || []).indexOf(1) !== -1}
						onToggle={(): void => {
							this.toggleLevel(1);
						}}
					/>
					<PageSwitch
						label="Medium alerts"
						help="Recieve medium level alerts on this device."
						hidden={!isPhone}
						disabled={this.state.disabled}
						checked={(device.alert_levels || []).indexOf(5) !== -1}
						onToggle={(): void => {
							this.toggleLevel(5);
						}}
					/>
					<PageSwitch
						label="High alerts"
						help="Recieve high level alerts on this device."
						hidden={!isPhone}
						disabled={this.state.disabled}
						checked={(device.alert_levels || []).indexOf(10) !== -1}
						onToggle={(): void => {
							this.toggleLevel(10);
						}}
					/>
					<PageInfo
						style={css.info}
						fields={fields1}
					/>
				</div>
				<div style={css.group}>
					<PageInfo
						style={css.info}
						fields={fields2}
					/>
				</div>
			</div>
			<PageSave
				style={css.save}
				hidden={!this.state.device && !this.state.message}
				message={this.state.message}
				changed={this.state.changed}
				disabled={this.state.disabled}
				light={true}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						device: null,
					});
				}}
				onSave={this.onSave}
			>
				{testButton}
			</PageSave>
		</div>;
	}
}
