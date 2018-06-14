/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as DeviceTypes from '../types/DeviceTypes';
import * as MiscUtils from '../utils/MiscUtils';
import * as DeviceActions from '../actions/DeviceActions';
import PageInfo from './PageInfo';
import ConfirmButton from './ConfirmButton';
import * as Alert from "../Alert";

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
		marginTop: '15px',
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
			Alert.success("Device name updated");

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
			}, 3000);
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

		let cardStyle = {
			...css.card,
		};
		if (device.disabled) {
			cardStyle.opacity = 0.6;
		}

		return <div
			className="pt-card"
			style={cardStyle}
		>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<div style={css.remove}>
						<ConfirmButton
							className="pt-minimal pt-intent-danger pt-icon-trash"
							progressClassName="pt-intent-danger"
							confirmMsg="Confirm node remove"
							disabled={this.state.disabled}
							onConfirm={this.onDelete}
						/>
					</div>
					<PageInfo
						style={css.info}
						fields={[
							{
								label: 'ID',
								value: device.id || 'None',
							},
						]}
					/>
					<div
						className="pt-input-group flex"
						style={css.inputGroup}
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
						]}
					/>
				</div>
			</div>
		</div>;
	}
}
