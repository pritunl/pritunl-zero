/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as DeviceTypes from '../types/DeviceTypes';
import * as MiscUtils from '../utils/MiscUtils';
import * as DeviceActions from '../actions/DeviceActions';
import * as PageInfos from './PageInfo';
import PageInfo from './PageInfo';
import ConfirmButton from './ConfirmButton';
import * as Alert from '../Alert';

interface Props {
	device: DeviceTypes.DeviceRo;
}

interface State {
	disabled: boolean;
	changed: boolean;
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
};

export default class Device extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
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

		let deviceType = 'Unknown';
		switch (device.type) {
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
		if (device.type === 'smart_card') {
			deviceOther = {
				label: 'SSH Public Key',
				value: device.ssh_public_key,
			};
		} else if (device.type === 'call' || device.type === 'sms') {
			deviceOther = {
				label: 'Phone Number',
				value: device.number,
			};
		}

		let cardStyle = {
			...css.card,
		};
		if (device.disabled) {
			cardStyle.opacity = 0.6;
		}

		return <div
			className="bp3-card"
			style={cardStyle}
		>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<div style={css.remove}>
						<ConfirmButton
							className="bp3-minimal bp3-intent-danger bp3-icon-trash"
							progressClassName="bp3-intent-danger"
							confirmMsg="Confirm device remove"
							disabled={this.state.disabled}
							onConfirm={this.onDelete}
						/>
					</div>
					<div
						className="bp3-input-group flex"
						style={css.inputGroup}
					>
						<input
							className="bp3-input"
							type="text"
							placeholder="Device name"
							value={device.name}
							onChange={(evt): void => {
								this.set('name', evt.target.value);
							}}
							onKeyPress={(evt): void => {
								if (evt.key === 'Enter') {
									this.onSave();
								}
							}}
						/>
						<button
							className="bp3-button bp3-minimal bp3-intent-primary bp3-icon-tick"
							hidden={!this.state.device}
							disabled={this.state.disabled}
							onClick={this.onSave}
						/>
					</div>
					<PageInfo
						style={css.info}
						fields={[
							{
								label: 'ID',
								value: device.id || 'None',
							},
							{
								label: 'Type',
								value: deviceType,
							},
							deviceOther,
						]}
					/>
				</div>
				<div style={css.group}>
					<PageInfo
						style={css.info}
						fields={[
							{
								label: 'Registered',
								value: MiscUtils.formatDate(device.timestamp) || 'Unknown',
							},
							{
								label: 'Last Active',
								value: MiscUtils.formatDate(device.last_active) || 'Unknown',
							},
							{
								label: 'Mode',
								value: deviceMode,
							},
						]}
					/>
				</div>
			</div>
		</div>;
	}
}
