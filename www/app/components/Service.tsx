/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ServiceTypes from '../types/ServiceTypes';
import * as ServiceActions from '../actions/ServiceActions';
import ServiceDomain from './ServiceDomain';
import ServiceServer from './ServiceServer';
import PageInput from './PageInput';
import PageSwitch from './PageSwitch';
import PageSave from './PageSave';
import PageInfo from './PageInfo';
import ConfirmButton from './ConfirmButton';
import PageInputButton from './PageInputButton';

interface Props {
	service: ServiceTypes.ServiceRo;
}

interface State {
	disabled: boolean;
	changed: boolean;
	message: string;
	addRole: string;
	addWhitelistNet: string;
	service: ServiceTypes.Service;
}

const css = {
	card: {
		position: 'relative',
		padding: '10px 10px 0 10px',
		marginBottom: '5px',
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
	} as React.CSSProperties,
	save: {
		paddingBottom: '10px',
	} as React.CSSProperties,
};

export default class Service extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
			addRole: '',
			addWhitelistNet: '',
			service: null,
		};
	}

	set(name: string, val: any): void {
		let service: any;

		if (this.state.changed) {
			service = {
				...this.state.service,
			};
		} else {
			service = {
				...this.props.service,
			};
		}

		service[name] = val;

		this.setState({
			...this.state,
			changed: true,
			service: service,
		});
	}

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		ServiceActions.commit(this.state.service).then((): void => {
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
						message: '',
						changed: false,
						service: null,
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
		ServiceActions.remove(this.props.service.id).then((): void => {
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
		let service: ServiceTypes.Service;

		if (this.state.changed) {
			service = {
				...this.state.service,
			};
		} else {
			service = {
				...this.props.service,
			};
		}

		let roles = [
			...service.roles,
		];

		if (!this.state.addRole) {
			return;
		}

		if (roles.indexOf(this.state.addRole) === -1) {
			roles.push(this.state.addRole);
		}

		roles.sort();

		service.roles = roles;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			service: service,
		});
	}

	onRemoveRole(role: string): void {
		let service: ServiceTypes.Service;

		if (this.state.changed) {
			service = {
				...this.state.service,
			};
		} else {
			service = {
				...this.props.service,
			};
		}

		let roles = [
			...service.roles,
		];

		let i = roles.indexOf(role);
		if (i === -1) {
			return;
		}

		roles.splice(i, 1);

		service.roles = roles;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			service: service,
		});
	}

	onAddWhitelistNet = (): void => {
		let service: ServiceTypes.Service;

		if (this.state.changed) {
			service = {
				...this.state.service,
			};
		} else {
			service = {
				...this.props.service,
			};
		}

		let whitelistNets = [
			...service.whitelist_networks,
		];

		if (!this.state.addWhitelistNet) {
			return;
		}

		if (whitelistNets.indexOf(this.state.addWhitelistNet) === -1) {
			whitelistNets.push(this.state.addWhitelistNet);
		}

		whitelistNets.sort();

		service.whitelist_networks = whitelistNets;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addWhitelistNet: '',
			service: service,
		});
	}

	onRemoveWhitelistNet(whitelistNet: string): void {
		let service: ServiceTypes.Service;

		if (this.state.changed) {
			service = {
				...this.state.service,
			};
		} else {
			service = {
				...this.props.service,
			};
		}

		let whitelistNets = [
			...service.whitelist_networks,
		];

		let i = whitelistNets.indexOf(whitelistNet);
		if (i === -1) {
			return;
		}

		whitelistNets.splice(i, 1);

		service.whitelist_networks = whitelistNets;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addWhitelistNet: '',
			service: service,
		});
	}

	onAddServer = (): void => {
		let service: ServiceTypes.Service;

		if (this.state.changed) {
			service = {
				...this.state.service,
			};
		} else {
			service = {
				...this.props.service,
			};
		}

		let servers = [
			...service.servers,
			{
				protocol: 'https',
				hostname: '',
				port: 443,
			},
		];

		service.servers = servers;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			service: service,
		});
	}

	onChangeServer(i: number, state: ServiceTypes.Server): void {
		let service: ServiceTypes.Service;

		if (this.state.changed) {
			service = {
				...this.state.service,
			};
		} else {
			service = {
				...this.props.service,
			};
		}

		let servers = [
			...service.servers,
		];

		servers[i] = state;

		service.servers = servers;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			service: service,
		});
	}

	onRemoveServer(i: number): void {
		let service: ServiceTypes.Service;

		if (this.state.changed) {
			service = {
				...this.state.service,
			};
		} else {
			service = {
				...this.props.service,
			};
		}

		let servers = [
			...service.servers,
		];

		servers.splice(i, 1);

		service.servers = servers;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			service: service,
		});
	}

	onAddDomain = (): void => {
		let service: ServiceTypes.Service;

		if (this.state.changed) {
			service = {
				...this.state.service,
			};
		} else {
			service = {
				...this.props.service,
			};
		}

		let domains = [
			...service.domains,
			{},
		];

		service.domains = domains;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			service: service,
		});
	}

	onChangeDomain(i: number, state: ServiceTypes.Domain): void {
		let service: ServiceTypes.Service;

		if (this.state.changed) {
			service = {
				...this.state.service,
			};
		} else {
			service = {
				...this.props.service,
			};
		}

		let domains = [
			...service.domains,
		];

		domains[i] = state;

		service.domains = domains;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			service: service,
		});
	}

	onRemoveDomain(i: number): void {
		let service: ServiceTypes.Service;

		if (this.state.changed) {
			service = {
				...this.state.service,
			};
		} else {
			service = {
				...this.props.service,
			};
		}

		let domains = [
			...service.domains,
		];

		domains.splice(i, 1);

		service.domains = domains;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			service: service,
		});
	}

	render(): JSX.Element {
		let service: ServiceTypes.Service = this.state.service ||
			this.props.service;

		let domains: JSX.Element[] = [];
		for (let i = 0; i < service.domains.length; i++) {
			let index = i;

			domains.push(
				<ServiceDomain
					key={index}
					domain={service.domains[index]}
					onChange={(state: ServiceTypes.Domain): void => {
						this.onChangeDomain(index, state);
					}}
					onRemove={(): void => {
						this.onRemoveDomain(index);
					}}
				/>,
			);
		}

		let roles: JSX.Element[] = [];
		for (let role of service.roles) {
			roles.push(
				<div
					className="pt-tag pt-tag-removable pt-intent-primary"
					style={css.item}
					key={role}
				>
					{role}
					<button
						className="pt-tag-remove"
						onMouseUp={(): void => {
							this.onRemoveRole(role);
						}}
					/>
				</div>,
			);
		}

		let servers: JSX.Element[] = [];
		for (let i = 0; i < service.servers.length; i++) {
			let index = i;

			servers.push(
				<ServiceServer
					key={index}
					server={service.servers[index]}
					onChange={(state: ServiceTypes.Server): void => {
						this.onChangeServer(index, state);
					}}
					onRemove={(): void => {
						this.onRemoveServer(index);
					}}
				/>,
			);
		}

		let whitelistNets: JSX.Element[] = [];
		for (let whitelistNet of service.whitelist_networks) {
			whitelistNets.push(
				<div
					className="pt-tag pt-tag-removable pt-intent-primary"
					style={css.item}
					key={whitelistNet}
				>
					{whitelistNet}
					<button
						className="pt-tag-remove"
						onMouseUp={(): void => {
							this.onRemoveWhitelistNet(whitelistNet);
						}}
					/>
				</div>,
			);
		}

		return <div
			className="pt-card"
			style={css.card}
		>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<div style={css.remove}>
						<ConfirmButton
							className="pt-minimal pt-intent-danger pt-icon-cross"
							progressClassName="pt-intent-danger"
							confirmMsg="Confirm service remove"
							disabled={this.state.disabled}
							onConfirm={this.onDelete}
						/>
					</div>
					<PageInput
						label="Name"
						type="text"
						placeholder="Enter name"
						value={service.name}
						onChange={(val): void => {
							this.set('name', val);
						}}
					/>
					<label style={css.itemsLabel}>
						External Domains
					</label>
					{domains}
					<button
						className="pt-button pt-intent-success pt-icon-add"
						style={css.itemsAdd}
						type="button"
						onClick={this.onAddDomain}
					>
						Add Domain
					</button>
					<label style={css.itemsLabel}>
						Internal Servers
					</label>
					{servers}
					<button
						className="pt-button pt-intent-success pt-icon-add"
						style={css.itemsAdd}
						type="button"
						onClick={this.onAddServer}
					>
						Add Server
					</button>
					<PageSwitch
						label="Share session with subdomains"
						checked={service.share_session}
						onToggle={(): void => {
							this.set('share_session', !service.share_session);
						}}
					/>
					<PageSwitch
						label="Allow WebSockets"
						checked={service.websockets}
						onToggle={(): void => {
							this.set('websockets', !service.websockets);
						}}
					/>
				</div>
				<div style={css.group}>
					<PageInfo
						fields={[
							{
								label: 'ID',
								value: service.id || 'None',
							},
						]}
					/>
					<label className="pt-label">
						Roles
						<div>
							{roles}
						</div>
					</label>
					<PageInputButton
						buttonClass="pt-intent-success pt-icon-add"
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
					<label className="pt-label">
						Whitelisted Networks
						<div>
							{whitelistNets}
						</div>
					</label>
					<PageInputButton
						buttonClass="pt-intent-success pt-icon-add"
						label="Add"
						type="text"
						placeholder="Add network"
						value={this.state.addWhitelistNet}
						onChange={(val): void => {
							this.setState({
								...this.state,
								addWhitelistNet: val,
							});
						}}
						onSubmit={this.onAddWhitelistNet}
					/>
				</div>
			</div>
			<PageSave
				style={css.save}
				hidden={!this.state.service}
				message={this.state.message}
				changed={this.state.changed}
				disabled={this.state.disabled}
				light={true}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						service: null,
					});
				}}
				onSave={this.onSave}
			/>
		</div>;
	}
}
