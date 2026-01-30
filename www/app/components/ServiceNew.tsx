/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ServiceTypes from '../types/ServiceTypes';
import * as AuthorityTypes from "../types/AuthorityTypes";
import * as ServiceActions from '../actions/ServiceActions';
import * as MiscUtils from '../utils/MiscUtils';
import ServiceDomain from './ServiceDomain';
import ServiceServer from './ServiceServer';
import ServiceWhitelistPath from './ServiceWhitelistPath';
import PageInput from './PageInput';
import PageSelect from './PageSelect';
import PageSwitch from './PageSwitch';
import PageCreate from './PageCreate';
import PageInfo from './PageInfo';
import ConfirmButton from './ConfirmButton';
import PageInputButton from './PageInputButton';
import Help from './Help';

interface Props {
	authorities: AuthorityTypes.AuthoritiesRo;
	onClose: () => void;
}

interface State {
	closed: boolean;
	disabled: boolean;
	changed: boolean;
	message: string;
	addRole: string;
	addWhitelistNet: string;
	service: ServiceTypes.Service;
}

const css = {
	row: {
		display: 'table-row',
		width: '100%',
		padding: 0,
		boxShadow: 'none',
		position: 'relative',
	} as React.CSSProperties,
	card: {
		position: 'relative',
		padding: '10px 10px 0 10px',
		width: '100%',
	} as React.CSSProperties,
	remove: {
		position: 'absolute',
		top: '5px',
		right: '5px',
	} as React.CSSProperties,
	item: {
		margin: '9px 5px 0 5px',
		minHeight: '20px',
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
		position: 'absolute',
		top: '5px',
		right: '5px',
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

export default class ServiceNew extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			closed: false,
			disabled: false,
			changed: false,
			message: '',
			addRole: '',
			addWhitelistNet: '',
			service: {},
		};
	}

	set(name: string, val: any): void {
		let service: any = {
			...this.state.service,
		};

		service[name] = val;

		this.setState({
			...this.state,
			changed: true,
			service: service,
		});
	}

	onCreate = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});

		let service: any = {
			...this.state.service,
		};

		ServiceActions.create(service).then((): void => {
			this.setState({
				...this.state,
				message: 'Service created successfully',
				changed: false,
			});

			setTimeout((): void => {
				this.setState({
					...this.state,
					disabled: false,
					changed: true,
				});
			}, 2000);
		}).catch((): void => {
			this.setState({
				...this.state,
				message: '',
				disabled: false,
			});
		});
	}

	onAddRole = (): void => {
		let service: ServiceTypes.Service;

		service = {
			...this.state.service,
		};

		let roles = [
			...(service.roles || []),
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

		service = {
			...this.state.service,
		};

		let roles = [
			...(service.roles || []),
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

		service = {
			...this.state.service,
		};

		let whitelistNets = [
			...(service.whitelist_networks || []),
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

		service = {
			...this.state.service,
		};

		let whitelistNets = [
			...(service.whitelist_networks || []),
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

		service = {
			...this.state.service,
		};

		let servers = [
			...(service.servers || []),
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

		service = {
			...this.state.service,
		};

		let servers = [
			...(service.servers || []),
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

		service = {
			...this.state.service,
		};

		let servers = [
			...(service.servers || []),
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

		service = {
			...this.state.service,
		};

		let domains = [
			...(service.domains || []),
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

		service = {
			...this.state.service,
		};

		let domains = [
			...(service.domains || []),
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

		service = {
			...this.state.service,
		};

		let domains = [
			...(service.domains || []),
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

	onAddWhitelistPath = (): void => {
		let service: ServiceTypes.Service;

		service = {
			...this.state.service,
		};

		let paths = [
			...(service.whitelist_paths || []),
			{},
		];

		service.whitelist_paths = paths;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			service: service,
		});
	}

	onChangeWhitelistPath(i: number, state: ServiceTypes.Path): void {
		let service: ServiceTypes.Service;

		service = {
			...this.state.service,
		};

		let paths = [
			...(service.whitelist_paths || []),
		];

		paths[i] = state;

		service.whitelist_paths = paths;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			service: service,
		});
	}

	onRemoveWhitelistPath(i: number): void {
		let service: ServiceTypes.Service;

		service = {
			...this.state.service,
		};

		let paths = [
			...(service.whitelist_paths || []),
		];

		paths.splice(i, 1);

		service.whitelist_paths = paths;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			service: service,
		});
	}

	render(): JSX.Element {
		let service: ServiceTypes.Service = this.state.service;

		let domains: JSX.Element[] = [];
		(service.domains || []).forEach((domn, index) => {
			domains.push(
				<ServiceDomain
					key={index}
					domain={domn}
					onChange={(state: ServiceTypes.Domain): void => {
						this.onChangeDomain(index, state);
					}}
					onRemove={(): void => {
						this.onRemoveDomain(index);
					}}
				/>,
			);
		})

		let roles: JSX.Element[] = [];
		(service.roles || []).forEach((role) => {
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
		})

		let servers: JSX.Element[] = [];
		(service.servers || []).forEach((server, index) => {
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
		})

		let authorities: JSX.Element[] = [
			<option key="null" value="">None</option>,
		];
		for (let authority of this.props.authorities) {
			authorities.push(
				<option
					key={authority.id}
					value={authority.id}
				>{authority.name}</option>,
			);
		}

		let whitelistNets: JSX.Element[] = [];
		(service.whitelist_networks || []).forEach((whitelistNet) => {
			whitelistNets.push(
				<div
					className="bp5-tag bp5-tag-removable bp5-intent-primary"
					style={css.item}
					key={whitelistNet}
				>
					{whitelistNet}
					<button
						className="bp5-tag-remove"
						onMouseUp={(): void => {
							this.onRemoveWhitelistNet(whitelistNet);
						}}
					/>
				</div>,
			);
		});

		let whitelistPaths: JSX.Element[] = [];
		(service.whitelist_paths || []).forEach((path, index) => {
			whitelistPaths.push(
				<ServiceWhitelistPath
					key={index}
					path={path}
					onChange={(state: ServiceTypes.Path): void => {
						this.onChangeWhitelistPath(index, state);
					}}
					onRemove={(): void => {
						this.onRemoveWhitelistPath(index);
					}}
				/>,
			);
		});

		return <div
			className="bp5-card bp5-row"
			style={css.row}
		>
		<td
			className="bp5-cell"
			colSpan={2}
			style={css.card}
		>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<div style={css.buttons}>
					</div>
						<PageInput
							label="Name"
							help="Name of service"
							type="text"
							placeholder="Enter name"
							value={service.name}
							onChange={(val): void => {
								this.set('name', val);
							}}
						/>
						<PageSelect
							label="Type"
							help="Service type"
							value={service.type}
							onChange={(val): void => {
								this.set('type', val);
							}}
						>
							<option value="http">HTTP</option>
						</PageSelect>
						<label style={css.itemsLabel}>
							External Domains
							<Help
								title="External Domains"
								contents={['When a request comes into a proxy node the requests host will be used to match the request with the domain of a service. The external domain must point to either a node that has the service added or a load balancer that forwards to nodes serving the service. Some internal services will be expecting a specific host such as a web server that serves mutliple websites that is also matching the requests host to one of the mutliple websites. Services that are associated with the same node should not also have the same domains. Wildcards are supported for the first component of domain. Multifactor U2F authentication is not supported for wildcard domains. When using a wildcard with U2F authentication the domain where the user login occurs must be included in external domains.', 'If the internal service is expecting a different host set the host field, otherwise leave it blank. ']}
							/>
						</label>
						{domains}
						<button
							className="bp5-button bp5-intent-success bp5-icon-add"
							style={css.itemsAdd}
							type="button"
							onClick={this.onAddDomain}
						>
							Add Domain
						</button>
						<label style={css.itemsLabel}>
							Internal Servers
							<Help
								title="Internal Servers"
								content="After a proxy node receives an authenticated request it will be forwarded to the internal servers and the response will be sent back to the user. Multiple internal servers can be added to load balance the requests. This should only be done if outages are not expected as no health checks are preformed for each server. If outages are expected a load balancer such as AWS ELB should be used. If a domain is used with HTTPS the internal server must have a valid certificate. When an IP address is used with HTTPS the internal servers certificate will not be validated. These internal servers should ideally be configured to only accept requests from the private IP addresses of the Pritunl Zero nodes. It is important to consider that if the internal servers are configured to accept requests from other IP addresses those requests will be sent directly to the internal server and will bypass the authentication provided by Pritunl Zero."
							/>
						</label>
						{servers}
						<button
							className="bp5-button bp5-intent-success bp5-icon-add"
							style={css.itemsAdd}
							type="button"
							onClick={this.onAddServer}
						>
							Add Server
						</button>
						<PageSelect
							label="Client Certificate Authority"
							help="Certificate authority to use for internal client certificate. Only valid for HTTPS connections to internal servers."
							value={service.client_authority}
							onChange={(val): void => {
								this.set('client_authority', val);
							}}
						>
							{authorities}
						</PageSelect>
						<PageInput
							label="Logout Path"
							help="Optional, path such as '/logout' that will end the Pritunl Zero users session. Supports '*' and '?' wildcards. This can cause unintentional logouts due to Chrome prefetching pages."
							type="text"
							placeholder="Enter logout path"
							value={service.logout_path}
							onChange={(val): void => {
								this.set('logout_path', val);
							}}
						/>
					</div>
					<div style={css.group}>
						<label className="bp5-label">
							Roles
							<Help
								title="Roles"
								content="The user roles that will be allowed access to this service. At least one role must match for the user to access the service."
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
						<label className="bp5-label">
							Permitted Networks
							<Help
								title="Permitted Networks"
								content="Permitted subnets with CIDR such as 10.0.0.0/8 that can access the service without authenticating. Single IP addresses can also be used. Any request coming from an IP address on these networks will be able to access the service without any authentication. Extra care should be taken when using this with the forwarded for header option in the node settings. If the nodes forwarded for header is enabled without a load balancer the user can modify the header value to spoof an IP address."
							/>
							<div>
								{whitelistNets}
							</div>
						</label>
						<PageInputButton
							buttonClass="bp5-intent-success bp5-icon-add"
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
						<label style={css.itemsLabel}>
							Permitted Paths
							<Help
								title="Permitted Paths"
								content="Permitted paths that can be accessed without authenticating. Supports '*' and '?' wildcards. Using this feature significantly increases the attack surface of the service and is not recommended."
							/>
						</label>
						{whitelistPaths}
						<button
							className="bp5-button bp5-intent-success bp5-icon-add"
							style={css.itemsAdd}
							type="button"
							onClick={this.onAddWhitelistPath}
						>
							Add Permitted Path
						</button>
						<PageSwitch
							label="Use HTTP/2"
							help="Use HTTP/2 transport."
							checked={service.http2}
							onToggle={(): void => {
								this.set('http2', !service.http2);
							}}
						/>
						<PageSwitch
							label="Share session with subdomains"
							help="This option will allow an authenticated user to access multiple services across different subdomains without needing to authenticate at each services subdomain."
							checked={service.share_session}
							onToggle={(): void => {
								this.set('share_session', !service.share_session);
							}}
						/>
						<PageSwitch
							label="Allow WebSockets"
							help="This will allow WebSockets to be proxied to the user. If the internal service relies on WebSockets this must be enabled."
							checked={service.websockets}
							onToggle={(): void => {
								this.set('websockets', !service.websockets);
							}}
						/>
						<PageSwitch
							label="CSRF check"
							help="Check headers to block cross domain requests. The service must implement CSRF protection if disabled."
							checked={!service.disable_csrf_check}
							onToggle={(): void => {
								this.set('disable_csrf_check', !service.disable_csrf_check);
							}}
						/>
						<PageSwitch
							label="Permit unauthenticated options requests"
							help="Permit HTTP OPTIONS requests to be proxied to the internal server without authentication."
							checked={service.whitelist_options}
							onToggle={(): void => {
								this.set('whitelist_options', !service.whitelist_options);
							}}
						/>
					</div>
				</div>
				<PageCreate
					style={css.save}
					hidden={!this.state.service}
					message={this.state.message}
					changed={this.state.changed}
					disabled={this.state.disabled}
					closed={this.state.closed}
					light={true}
					onCancel={this.props.onClose}
					onCreate={this.onCreate}
				/>
			</td>
		</div>;
	}
}
