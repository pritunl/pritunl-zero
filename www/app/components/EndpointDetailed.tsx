/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as EndpointTypes from '../types/EndpointTypes';
import * as AuthorityTypes from "../types/AuthorityTypes";
import * as EndpointActions from '../actions/EndpointActions';
import * as PageInfos from './PageInfo';
import PageInput from './PageInput';
import PageSave from './PageSave';
import PageInfo from './PageInfo';
import ConfirmButton from './ConfirmButton';
import PageInputButton from './PageInputButton';
import EndpointCharts from './EndpointCharts';
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
			showCharts: false,
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

		let endpointData = endpoint.data;
		if (endpointData) {
			if (endpointData.hostname) {
				fields.push({
					label: 'Hostname',
					value: endpointData.hostname,
				});
			}
			if (endpointData.uptime) {
				fields.push({
					label: 'Uptime',
					value: endpointData.uptime,
				});
			}
			if (endpointData.platform) {
				fields.push({
					label: 'Platform',
					value: endpointData.platform,
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
					value: endpointData.mem_total + 'mb',
				});
			}
			if (endpointData.swap_total) {
				fields.push({
					label: 'Swap',
					value: endpointData.swap_total + 'mb',
				});
			}
		}

		let roles: JSX.Element[] = [];
		for (let role of endpoint.roles) {
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

		return <td
			className="bp3-cell"
			colSpan={2}
			style={css.card}
		>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<div
						className="layout horizontal"
						style={css.buttons}
						onClick={(evt): void => {
							let target = evt.target as HTMLElement;

							if (target.className.indexOf('open-ignore') !== -1) {
								return;
							}

							this.props.onClose();
						}}
					>
            <div>
              <label
                className="bp3-control bp3-checkbox open-ignore"
                style={css.select}
              >
                <input
                  type="checkbox"
                  className="open-ignore"
                  checked={this.props.selected}
									onChange={(evt): void => {
									}}
                  onClick={(evt): void => {
										this.props.onSelect(evt.shiftKey);
									}}
                />
                <span className="bp3-control-indicator open-ignore"/>
              </label>
            </div>
						<div className="flex"/>
						<ConfirmButton
							className="bp3-minimal bp3-intent-danger bp3-icon-trash open-ignore"
							style={css.button}
							progressClassName="bp3-intent-danger"
							confirmMsg="Confirm endpoint remove"
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
					<PageSwitch
						label="Show charts"
						help="Show endpoint charts."
						checked={this.state.showCharts}
						onToggle={(): void => {
							this.setState({
								...this.state,
								showCharts: !this.state.showCharts,
							});
						}}
					/>
				</div>
				<div style={css.group}>
					<PageInfo
						fields={fields}
					/>
					<label className="bp3-label">
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
				</div>
			</div>
			<EndpointCharts
				endpoint={endpoint.id}
				disabled={!this.state.showCharts}
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
			/>
		</td>;
	}
}
