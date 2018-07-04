/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';
import * as DeviceTypes from '../types/DeviceTypes';
import * as DeviceActions from '../actions/DeviceActions';
import * as MiscUtils from "../utils/MiscUtils";
import * as Alert from "../Alert";
import ConfirmButton from './ConfirmButton';

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
		padding: '5px 3px 2px 7px',
		marginBottom: '10px',
	} as React.CSSProperties,
	info: {
		marginTop: '5px',
	} as React.CSSProperties,
	icon: {
		marginTop: '5px',
	} as React.CSSProperties,
	name: {
		margin: '0 3px 0 7px',
	} as React.CSSProperties,
	group: {
		margin: '0 3px 0 7px',
	} as React.CSSProperties,
	nameGroup: {
		margin: 0,
	} as React.CSSProperties,
	item: {
		marginBottom: '3px',
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
			DeviceActions.sync();

			Alert.success("Device name updated");

			this.setState({
				...this.state,
				changed: false,
				disabled: false,
				device: null,
			});
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
			DeviceActions.sync();

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
		}

		let deviceMode = 'Unknown';
		switch (device.mode) {
			case 'secondary':
				deviceMode = 'Secondary';
				break;
			case 'ssh':
				deviceMode = 'SSH';
				break;
		}

		return <div
			className="pt-card"
			style={css.card}
		>
			<div className="layout horizontal">
				<Blueprint.Icon icon="id-number" iconSize={20} style={css.icon}/>
				<div
					className="pt-input-group flex"
					style={css.group}
				>
					<input
						className="pt-input"
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
						className="pt-button pt-minimal pt-intent-primary pt-icon-tick"
						hidden={!this.state.device}
						disabled={this.state.disabled}
						onClick={this.onSave}
					/>
				</div>
				<div>
					<ConfirmButton
						className="pt-minimal pt-intent-danger pt-icon-trash"
						progressClassName="pt-intent-danger"
						confirmMsg="Confirm device remove"
						disabled={this.state.disabled}
						onConfirm={this.onDelete}
					/>
				</div>
			</div>
			<div className="layout vertical" style={css.info}>
				<div style={css.item}>
					ID: <span className="pt-text-muted">
						{device.id}
					</span>
				</div>
				<div style={css.item}>
					Type: <span className="pt-text-muted">
						{deviceType}
					</span>
				</div>
				<div style={css.item}>
					Mode: <span className="pt-text-muted">
						{deviceMode}
					</span>
				</div>
				<div style={css.item}>
					Registered: <span className="pt-text-muted">
						{MiscUtils.formatDateMid(device.timestamp)}
					</span>
				</div>
				<div style={css.item}>
					Last Active: <span className="pt-text-muted">
						{MiscUtils.formatDateMid(device.last_active)}
					</span>
				</div>
			</div>
		</div>;
	}
}
