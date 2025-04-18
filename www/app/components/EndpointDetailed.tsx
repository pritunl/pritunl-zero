/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as EndpointTypes from '../types/EndpointTypes';
import * as AuthorityTypes from "../types/AuthorityTypes";
import * as EndpointActions from '../actions/EndpointActions';
import * as PageInfos from './PageInfo';
import * as MiscUtils from '../utils/MiscUtils';
import PageInput from './PageInput';
import PageSave from './PageSave';
import PageInfo from './PageInfo';
import ConfirmButton from './ConfirmButton';
import PageInputButton from './PageInputButton';
import EndpointCharts from './EndpointCharts';
import EndpointKmsg from './EndpointKmsg';
import Help from './Help';
import PageSwitch from "./PageSwitch";

interface Props {
	endpoint: EndpointTypes.EndpointRo;
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
	showCharts: boolean;
	endpoint: EndpointTypes.Endpoint;
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
};

export default class EndpointDetailed extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
			addRole: '',
			endpoint: null,
			showCharts: true,
		};
	}

	set(name: string, val: any): void {
		let endpoint: any;

		if (this.state.changed) {
			endpoint = {
				...this.state.endpoint,
			};
		} else {
			endpoint = {
				...this.props.endpoint,
			};
		}

		endpoint[name] = val;

		this.setState({
			...this.state,
			changed: true,
			endpoint: endpoint,
		});
	}

	onResetClientKey = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});

		let endpoint = {
			...this.props.endpoint,
			reset_client_key: true,
		};

		EndpointActions.commit(endpoint).then((): void => {
			this.setState({
				...this.state,
				message: 'Client key reset',
				changed: false,
				disabled: false,
			});

			setTimeout((): void => {
				if (!this.state.changed) {
					this.setState({
						...this.state,
						endpoint: null,
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

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		EndpointActions.commit(this.state.endpoint).then((): void => {
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
						endpoint: null,
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
		EndpointActions.remove(this.props.endpoint.id).then((): void => {
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
		let endpoint: EndpointTypes.Endpoint;

		if (this.state.changed) {
			endpoint = {
				...this.state.endpoint,
			};
		} else {
			endpoint = {
				...this.props.endpoint,
			};
		}

		let roles = [
			...endpoint.roles,
		];

		if (!this.state.addRole) {
			return;
		}

		if (roles.indexOf(this.state.addRole) === -1) {
			roles.push(this.state.addRole);
		}

		roles.sort();

		endpoint.roles = roles;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			endpoint: endpoint,
		});
	}

	onRemoveRole(role: string): void {
		let endpoint: EndpointTypes.Endpoint;

		if (this.state.changed) {
			endpoint = {
				...this.state.endpoint,
			};
		} else {
			endpoint = {
				...this.props.endpoint,
			};
		}

		let roles = [
			...endpoint.roles,
		];

		let i = roles.indexOf(role);
		if (i === -1) {
			return;
		}

		roles.splice(i, 1);

		endpoint.roles = roles;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			endpoint: endpoint,
		});
	}

	render(): JSX.Element {
		let endpoint: EndpointTypes.Endpoint = this.state.endpoint ||
			this.props.endpoint;

		let fields: PageInfos.Field[] = [
			{
				label: 'ID',
				value: this.props.endpoint.id || 'None',
			},
		];

		let endpointData = endpoint.data || {};
		if (endpointData) {
			if (endpointData.version) {
				fields.push({
					label: 'Endpoint Version',
					value: endpointData.version,
				});
			}
			if (endpointData.hostname) {
				fields.push({
					label: 'Hostname',
					value: endpointData.hostname,
				});
			}
			if (endpointData.uptime) {
				fields.push({
					label: 'Uptime',
					value: MiscUtils.formatUptime(endpointData.uptime),
				});
			}
			if (endpointData.platform) {
				fields.push({
					label: 'Platform',
					value: endpointData.platform,
				});
			}
			if (endpointData.package_updates) {
				fields.push({
					label: 'System Package Updates',
					value: endpointData.package_updates,
					valueClass: 'bp5-text-intent-danger',
				});
			}
			if (endpointData.virtualization) {
				fields.push({
					label: 'Virtualization',
					value: endpointData.virtualization,
				});
			}
			if (endpointData.cpu_cores) {
				fields.push({
					label: 'CPU Cores',
					value: endpointData.cpu_cores,
				});
			}
			if (endpointData.mem_total) {
				fields.push({
					label: 'Memory',
					value: endpointData.mem_total + 'MB',
				});
			}
			if (endpointData.swap_total) {
				fields.push({
					label: 'Swap',
					value: endpointData.swap_total + 'MB',
				});
			}
			if (endpointData.huge_total) {
				fields.push({
					label: 'HugePages',
					value: endpointData.huge_total + 'MB',
				});
			}
		}

		if (endpoint.data.md_stat && endpoint.data.md_stat.length) {
			let failed = 0;
			let total = 0;

			for (let md of endpoint.data.md_stat) {
				failed += md.failed;
				total += md.total;
			}

			fields.push({
				label: 'Raid Devices',
				value: 'Failed: ' + failed + ' Total: ' + total,
			});
		}

		let roles: JSX.Element[] = [];
		for (let role of endpoint.roles) {
			roles.push(
				<div
					className="bp5-tag bp5-tag-removable bp5-intent-primary"
					style={css.item}
					key={role}
				>
					{role}
					<button
						className="bp5-tag-remove"
						onMouseUp={(): void => {
							this.onRemoveRole(role);
						}}
					/>
				</div>,
			);
		}

		let alerts: string[] = [];
		for (let alert of Object.values(endpoint.info.alerts)) {
			alerts.push(alert);
		}

		let checks: string[] = [];
		for (let check of Object.values(endpoint.info.checks)) {
			checks.push(check);
		}

		let secretKey = '';
		let secretUri = '';
		if (!endpoint.has_client_key) {
			if (endpoint.client_key) {
				secretKey = endpoint.id + '_' + endpoint.client_key.secret;
			} else {
				secretKey = 'unknown';
			}

			let secretDomain = '';
			if (this.props.endpoint.info.domain) {
				secretDomain = this.props.endpoint.info.domain;
			} else {
				secretDomain = window.location.host;
			}

			secretUri = 'pritunl://' + secretDomain + '/' + secretKey;
		}

		return <td
			className="bp5-cell"
			colSpan={3}
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
                className="bp5-control bp5-checkbox"
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
                <span className="bp5-control-indicator"/>
              </label>
            </div>
						<div className="flex tab-close"/>
						<ConfirmButton
							safe={true}
							className="bp5-minimal bp5-intent-danger bp5-icon-trash"
							progressClassName="bp5-intent-danger"
							dialogClassName="bp5-intent-danger bp5-icon-delete"
							dialogLabel="Delete Endpoint"
							confirmMsg="Permanently delete this endpoint"
							confirmInput={true}
							items={[endpoint.name]}
							disabled={this.state.disabled}
							onConfirm={this.onDelete}
						/>
					</div>
					<PageInput
						label="Name"
						help="Name of endpoint"
						type="text"
						placeholder="Enter name"
						value={endpoint.name}
						onChange={(val): void => {
							this.set('name', val);
						}}
					/>
					<PageInput
						hidden={endpoint.has_client_key}
						disabled={this.state.disabled}
						readOnly={true}
						autoSelect={true}
						label="Registration Key"
						help="Key for endpoint registration"
						type="text"
						placeholder=""
						value={secretKey}
					/>
					<label className="bp5-label">
						Roles
						<Help
							title="Roles"
							content="The user roles that will be allowed access to this endpoint. At least one role must match for the user to access the endpoint."
						/>
						<div>
							{roles}
						</div>
					</label>
					<PageInputButton
						buttonClass="bp5-intent-success bp5-icon-add"
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
					<PageSwitch
						label="Show charts and dmesg"
						help="Show endpoint charts and dmesg."
						checked={this.state.showCharts}
						hidden={!endpointData.hostname}
						onToggle={(): void => {
							this.setState({
								...this.state,
								showCharts: !this.state.showCharts,
							});
						}}
					/>
					<PageInfo
						fields={[
							{
								label: 'Alerts',
								value: alerts.length ? alerts : '-',
							},
							{
								label: 'Health Checks',
								value: checks.length ? checks : '-',
							},
						]}
					/>
				</div>
				<div style={css.group}>
					<PageInfo
						fields={fields}
					/>
					<PageInput
						hidden={endpoint.has_client_key}
						disabled={this.state.disabled}
						readOnly={true}
						autoSelect={true}
						label="Registration URI"
						help="URI for endpoint registration"
						type="text"
						placeholder=""
						value={secretUri}
					/>
				</div>
			</div>
			<EndpointCharts
				endpoint={endpoint.id}
				disabled={!endpointData.hostname || !this.state.showCharts}
			/>
			<EndpointKmsg
				endpoint={endpoint.id}
				disabled={!endpointData.hostname || !this.state.showCharts}
			/>
			<PageSave
				style={css.save}
				hidden={!this.state.endpoint && !this.state.message}
				message={this.state.message}
				changed={this.state.changed}
				disabled={this.state.disabled}
				light={true}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						endpoint: null,
					});
				}}
				onSave={this.onSave}
			>
				<ConfirmButton
					label="Reset Key"
					className="bp5-intent-danger bp5-icon-key"
					progressClassName="bp5-intent-danger"
					style={css.controlButton}
					hidden={!endpoint.has_client_key}
					disabled={this.state.disabled}
					safe={true}
					onConfirm={(): void => {
						this.onResetClientKey();
					}}
				/>
			</PageSave>
		</td>;
	}
}
