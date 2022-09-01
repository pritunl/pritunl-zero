/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as AlertTypes from '../types/AlertTypes';
import * as AuthorityTypes from "../types/AuthorityTypes";
import * as AlertActions from '../actions/AlertActions';
import * as PageInfos from './PageInfo';
import PageInput from './PageInput';
import PageSave from './PageSave';
import PageInfo from './PageInfo';
import ConfirmButton from './ConfirmButton';
import PageInputButton from './PageInputButton';
import Help from './Help';
import PageSwitch from "./PageSwitch";
import PageSelect from "./PageSelect";

interface Props {
	alert: AlertTypes.AlertRo;
	authorities: AuthorityTypes.AuthoritiesRo;
	selected: boolean;
	onSelect: (shift: boolean) => void;
	onClose: () => void;
}

interface State {
	disabled: boolean;
	changed: boolean;
	message: string;
	addRole: string;
	addIgnore: string;
	alert: AlertTypes.Alert;
}

const css = {
	card: {
		position: 'relative',
		padding: '48px 10px 0 10px',
		width: '100%',
	} as React.CSSProperties,
	remove: {
		position: 'absolute',
		top: '5px',
		right: '5px',
	} as React.CSSProperties,
	item: {
		margin: '9px 5px 0 5px',
		height: '20px',
	} as React.CSSProperties,
	itemsLabel: {
		display: 'block',
	} as React.CSSProperties,
	itemsAdd: {
		margin: '8px 0 15px 0',
	} as React.CSSProperties,
	group: {
		flex: 1,
		minWidth: '250px',
		margin: '0 10px',
	} as React.CSSProperties,
	controlButton: {
		marginRight: '10px',
	} as React.CSSProperties,
	save: {
		paddingBottom: '10px',
	} as React.CSSProperties,
	button: {
		height: '30px',
	} as React.CSSProperties,
	buttons: {
		cursor: 'pointer',
		position: 'absolute',
		top: 0,
		left: 0,
		right: 0,
		padding: '4px',
		height: '39px',
		backgroundColor: 'rgba(0, 0, 0, 0.13)',
	} as React.CSSProperties,
	label: {
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	status: {
		margin: '6px 0 0 1px',
	} as React.CSSProperties,
	icon: {
		marginRight: '3px',
	} as React.CSSProperties,
	inputGroup: {
		width: '100%',
	} as React.CSSProperties,
	protocol: {
		flex: '0 1 auto',
	} as React.CSSProperties,
	port: {
		flex: '1',
	} as React.CSSProperties,
	select: {
		margin: '7px 0px 0px 6px',
		paddingTop: '3px',
	} as React.CSSProperties,
	header: {
		fontSize: '20px',
		marginTop: '-10px',
		paddingBottom: '2px',
		marginBottom: '10px',
		borderBottomStyle: 'solid',
	} as React.CSSProperties,
	heading: {
		margin: '19px 0 0 0',
	} as React.CSSProperties,
	alertsButtons: {
		marginTop: '8px',
	} as React.CSSProperties,
	alertsAdd: {
		margin: '8px 0 0 8px',
	} as React.CSSProperties,
};

export default class AlertDetailed extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
			addRole: '',
			addIgnore: '',
			alert: null,
		};
	}

	set(name: string, val: any): void {
		let alert: any;

		if (this.state.changed) {
			alert = {
				...this.state.alert,
			};
		} else {
			alert = {
				...this.props.alert,
			};
		}

		alert[name] = val;

		this.setState({
			...this.state,
			changed: true,
			alert: alert,
		});
	}

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		AlertActions.commit(this.state.alert).then((): void => {
			this.setState({
				...this.state,
				message: 'Your changes have been saved',
				changed: false,
				disabled: false,
			});

			setTimeout((): void => {
				if (!this.state.changed) {
					this.setState({
						...this.state,
						alert: null,
						changed: false,
					});
				}
			}, 1000);

			setTimeout((): void => {
				if (!this.state.changed) {
					this.setState({
						...this.state,
						message: '',
					});
				}
			}, 3000);
		}).catch((): void => {
			this.setState({
				...this.state,
				message: '',
				disabled: false,
			});
		});
	}

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		AlertActions.remove(this.props.alert.id).then((): void => {
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

	onAddRole = (): void => {
		let alert: AlertTypes.Alert;

		if (this.state.changed) {
			alert = {
				...this.state.alert,
			};
		} else {
			alert = {
				...this.props.alert,
			};
		}

		let roles = [
			...alert.roles,
		];

		if (!this.state.addRole) {
			return;
		}

		if (roles.indexOf(this.state.addRole) === -1) {
			roles.push(this.state.addRole);
		}

		roles.sort();

		alert.roles = roles;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			alert: alert,
		});
	}

	onRemoveRole(role: string): void {
		let alert: AlertTypes.Alert;

		if (this.state.changed) {
			alert = {
				...this.state.alert,
			};
		} else {
			alert = {
				...this.props.alert,
			};
		}

		let roles = [
			...alert.roles,
		];

		let i = roles.indexOf(role);
		if (i === -1) {
			return;
		}

		roles.splice(i, 1);

		alert.roles = roles;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			alert: alert,
		});
	}

	onAddIgnore = (): void => {
		let alert: AlertTypes.Alert;

		if (this.state.changed) {
			alert = {
				...this.state.alert,
			};
		} else {
			alert = {
				...this.props.alert,
			};
		}

		let ignores = [
			...(alert.ignores || []),
		];

		if (!this.state.addIgnore) {
			return;
		}

		if (ignores.indexOf(this.state.addIgnore) === -1) {
			ignores.push(this.state.addIgnore);
		}

		ignores.sort();

		alert.ignores = ignores;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addIgnore: '',
			alert: alert,
		});
	}

	onRemoveIgnore(ignore: string): void {
		let alert: AlertTypes.Alert;

		if (this.state.changed) {
			alert = {
				...this.state.alert,
			};
		} else {
			alert = {
				...this.props.alert,
			};
		}

		let ignores = [
			...(alert.ignores || []),
		];

		let i = ignores.indexOf(ignore);
		if (i === -1) {
			return;
		}

		ignores.splice(i, 1);

		alert.ignores = ignores;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addIgnore: '',
			alert: alert,
		});
	}

	render(): JSX.Element {
		let alert: AlertTypes.Alert = this.state.alert ||
			this.props.alert;

		let fields: PageInfos.Field[] = [
			{
				label: 'ID',
				value: this.props.alert.id || 'None',
			},
		];

		let roles: JSX.Element[] = [];
		for (let role of alert.roles) {
			roles.push(
				<div
					className="bp3-tag bp3-tag-removable bp3-intent-primary"
					style={css.item}
					key={role}
				>
					{role}
					<button
						className="bp3-tag-remove"
						onMouseUp={(): void => {
							this.onRemoveRole(role);
						}}
					/>
				</div>,
			);
		}

		let ignores: JSX.Element[] = [];
		for (let ignore of (alert.ignores || [])) {
			ignores.push(
				<div
					className="bp3-tag bp3-tag-removable bp3-intent-primary"
					style={css.item}
					key={ignore}
				>
					{ignore}
					<button
						className="bp3-tag-remove"
						onMouseUp={(): void => {
							this.onRemoveIgnore(ignore);
						}}
					/>
				</div>,
			);
		}

		let valueInt = false;
		let valueStr = false;
		let valueLabel = '';
		let valueHelp = '';
		let ignoreShow = false;
		let ignoreLabel = '';
		let ignoreTitle = '';
		let ignoreHelp = '';
		switch (alert.resource) {
			case "system_cpu_level":
				valueInt = true;
				valueLabel = 'Usage Threshold';
				valueHelp = 'Maximum percent CPU usage as integer ' +
					'before alert is triggered.';
				break;
			case "system_memory_level":
				valueInt = true;
				valueLabel = 'Usage Threshold';
				valueHelp = 'Maximum percent memory usage as integer ' +
					'before alert is triggered.';
				break;
			case "system_swap_level":
				valueInt = true;
				valueLabel = 'Usage Threshold';
				valueHelp = 'Maximum percent swap usage as integer ' +
					'before alert is triggered.';
				break;
			case "system_hugepages_level":
				valueInt = true;
				valueLabel = 'Usage Threshold';
				valueHelp = 'Maximum percent hugepages usage as integer ' +
					'before alert is triggered.';
				break;
			case "system_md_failed":
				valueInt = false;
				valueStr = false;
				break;
			case "disk_usage_level":
				ignoreShow = true;
				ignoreLabel = 'Ignore Disk Paths';
				ignoreTitle = 'Ignore Disk Paths';
				ignoreHelp = 'Path of disk devices to ignore.';
				valueInt = true;
				valueLabel = 'Usage Threshold';
				valueHelp = 'Maximum percent disk space usage as integer ' +
					'before alert is triggered.';
				break;
			case "kmsg_keyword":
				valueStr = true;
				valueLabel = 'Dmesg Keyword Match';
				valueHelp = 'Case insensitive dmesg match string to trigger alert.';
				break;
			case "check_http_failed":
				valueInt = false;
				valueStr = false;
				break;
		}

		return <td
			className="bp3-cell"
			colSpan={2}
			style={css.card}
		>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<div
						className="layout horizontal tab-close"
						style={css.buttons}
						onClick={(evt): void => {
							let target = evt.target as HTMLElement;

							if (target.className.indexOf('tab-close') !== -1) {
								this.props.onClose();
							}
						}}
					>
						<div>
							<label
								className="bp3-control bp3-checkbox"
								style={css.select}
							>
								<input
									type="checkbox"
									checked={this.props.selected}
									onChange={(evt): void => {
									}}
									onClick={(evt): void => {
										this.props.onSelect(evt.shiftKey);
									}}
								/>
								<span className="bp3-control-indicator"/>
							</label>
						</div>
						<div className="flex tab-close"/>
						<ConfirmButton
							safe={true}
							className="bp3-minimal bp3-intent-danger bp3-icon-trash"
							progressClassName="bp3-intent-danger"
							dialogClassName="bp3-intent-danger bp3-icon-delete"
							dialogLabel="Delete Alert"
							confirmMsg="Permanently delete this alert"
							confirmInput={true}
							items={[alert.name]}
							disabled={this.state.disabled}
							onConfirm={this.onDelete}
						/>
					</div>
					<PageInput
						disabled={this.state.disabled}
						label="Name"
						help="Name of alert"
						type="text"
						placeholder="Enter name"
						value={alert.name}
						onChange={(val): void => {
							this.set('name', val);
						}}
					/>
					<label className="bp3-label">
						Roles
						<Help
							title="Roles"
							content="The user roles that will be allowed access to this alert. At least one role must match for the user to access the alert."
						/>
						<div>
							{roles}
						</div>
					</label>
					<PageInputButton
						disabled={this.state.disabled}
						buttonClass="bp3-intent-success bp3-icon-add"
						label="Add"
						type="text"
						placeholder="Add role"
						value={this.state.addRole}
						onChange={(val): void => {
							this.setState({
								...this.state,
								addRole: val,
							});
						}}
						onSubmit={this.onAddRole}
					/>
					<PageSelect
						disabled={this.state.disabled}
						label="Alert Type"
						help="Type of alert"
						value={alert.resource}
						onChange={(val): void => {
							this.set('resource', val);
						}}
					>
						<option
							value="system_cpu_level"
						>CPU Usage Threshold</option>
						<option
							value="system_memory_level"
						>Memory Usage Threshold</option>
						<option
							value="system_swap_level"
						>Swap Usage Threshold</option>
						<option
							value="system_hugepages_level"
						>HugePages Usage Threshold</option>
						<option
							value="system_md_failed"
						>MD RAID Device Failed</option>
						<option
							value="disk_usage_level"
						>Disk Usage Threshold</option>
						<option
							value="kmsg_keyword"
						>Dmesg Keyword Match</option>
						<option
							value="check_http_failed"
						>HTTP Health Check Failed</option>
					</PageSelect>
					<label className="bp3-label" hidden={!ignoreShow}>
						{ignoreLabel}
						<Help
							title={ignoreTitle}
							content={ignoreHelp}
						/>
						<div>
							{ignores}
						</div>
					</label>
					<PageInputButton
						disabled={this.state.disabled}
						buttonClass="bp3-intent-success bp3-icon-add"
						label="Add"
						type="text"
						placeholder="Add ignore"
						value={this.state.addIgnore}
						hidden={!ignoreShow}
						onChange={(val): void => {
							this.setState({
								...this.state,
								addIgnore: val,
							});
						}}
						onSubmit={this.onAddIgnore}
					/>
				</div>
				<div style={css.group}>
					<PageInfo
						fields={fields}
					/>
					<PageInput
						disabled={this.state.disabled}
						label={valueLabel}
						help={valueHelp}
						type="text"
						placeholder="Default"
						value={alert.value_int}
						hidden={!valueInt}
						onChange={(val): void => {
							this.set('value_int', parseInt(val, 10));
						}}
					/>
					<PageInput
						disabled={this.state.disabled}
						label={valueLabel}
						help={valueHelp}
						type="text"
						placeholder="Default"
						value={alert.value_str}
						hidden={!valueStr}
						onChange={(val): void => {
							this.set('value_str', val);
						}}
					/>
					<PageSelect
						disabled={this.state.disabled}
						label="Alert Level"
						help="Level of alert, used for matching device notifications. An endpoint role must also match a user role for ntofications."
						value={(alert.level || 0).toString()}
						onChange={(val): void => {
							this.set('level', parseInt(val, 10));
						}}
					>
						<option value="1">Low</option>
						<option value="5">Medium</option>
						<option value="10">High</option>
					</PageSelect>
					<PageInput
						disabled={this.state.disabled}
						label="Alert Frequency"
						help="Minimum duration in seconds between repeat alerts."
						type="text"
						placeholder="Enter frequency"
						value={alert.frequency}
						onChange={(val): void => {
							this.set('frequency', parseInt(val, 10));
						}}
					/>
				</div>
			</div>
			<PageSave
				style={css.save}
				hidden={!this.state.alert && !this.state.message}
				message={this.state.message}
				changed={this.state.changed}
				disabled={this.state.disabled}
				light={true}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						alert: null,
					});
				}}
				onSave={this.onSave}
			/>
		</td>;
	}
}
