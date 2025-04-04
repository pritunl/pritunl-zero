/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as AuthorityTypes from '../types/AuthorityTypes';
import * as NodeTypes from "../types/NodeTypes";
import * as AuthorityActions from '../actions/AuthorityActions';
import PageInput from './PageInput';
import PageSwitch from './PageSwitch';
import PageSelect from './PageSelect';
import PageInputButton from './PageInputButton';
import AuthorityDeploy from './AuthorityDeploy';
import PageTextAreaTab from './PageTextAreaTab';
import * as PageInfos from './PageInfo';
import PageInfo from './PageInfo';
import PageCreate from './PageCreate';
import ConfirmButton from './ConfirmButton';
import Help from './Help';
import * as MiscUtils from "../utils/MiscUtils";
import AuthoritiesStore from "../stores/AuthoritiesStore";

interface Props {
	nodes: NodeTypes.NodesRo;
	onClose: () => void;
}

interface State {
	closed: boolean;
	disabled: boolean;
	changed: boolean;
	message: string;
	authority: AuthorityTypes.Authority;
	addRole: string;
	addMatch: string;
	addSubnet: string;
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
	label: {
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	inputGroup: {
		width: '100%',
	} as React.CSSProperties,
	hostname: {
		flex: '1',
		minWidth: '160px',
	} as React.CSSProperties,
	port: {
		width: '60px',
		flex: '0 1 auto',
	} as React.CSSProperties,
	controlButton: {
		marginRight: '10px',
	} as React.CSSProperties,
	button: {
		height: '30px',
	} as React.CSSProperties,
	buttons: {
		position: 'absolute',
		top: '5px',
		right: '5px',
	} as React.CSSProperties,
	select: {
		margin: '7px 0px 0px 6px',
		paddingTop: '3px',
	} as React.CSSProperties,
};

export default class AuthorityDetailed extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			closed: false,
			disabled: false,
			changed: false,
			message: '',
			addRole: null,
			addMatch: null,
			addSubnet: null,
			authority: {
				name: 'New Authority',
			},
		};
	}

	set(name: string, val: any): void {
		let authority: any = {
			...this.state.authority,
		};

		authority[name] = val;

		this.setState({
			...this.state,
			changed: true,
			authority: authority,
		});
	}

	toggle(name: string): void {
		let authority: any = {
			...this.state.authority,
		};

		authority[name] = !authority[name];

		this.setState({
			...this.state,
			changed: true,
			authority: authority,
		});
	}

	onCreate = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});

		let authority: any = {
			...this.state.authority,
		};

		AuthorityActions.create(authority).then((): void => {
			this.setState({
				...this.state,
				message: 'Authority created successfully',
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
		let authority:  AuthorityTypes.Authority;

		authority = {
			...this.state.authority,
		};

		let roles = [
			...authority.roles,
		];

		if (!this.state.addRole) {
			return;
		}

		if (roles.indexOf(this.state.addRole) === -1) {
			roles.push(this.state.addRole);
		}

		roles.sort();

		authority.roles = roles;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			authority: authority,
		});
	}

	onRemoveRole(role: string): void {
		let authority:  AuthorityTypes.Authority;

		authority = {
			...this.state.authority,
		};

		let roles = [
			...authority.roles,
		];

		let i = roles.indexOf(role);
		if (i === -1) {
			return;
		}

		roles.splice(i, 1);

		authority.roles = roles;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			authority: authority,
		});
	}

	onAddMatch = (): void => {
		let authority:  AuthorityTypes.Authority;

		authority = {
			...this.state.authority,
		};

		let matches = [
			...(authority.host_matches || []),
		];

		if (!this.state.addMatch) {
			return;
		}

		if (matches.indexOf(this.state.addMatch) === -1) {
			matches.push(this.state.addMatch);
		}

		matches.sort();

		authority.host_matches = matches;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addMatch: '',
			authority: authority,
		});
	}

	onRemoveMatch(match: string): void {
		let authority:  AuthorityTypes.Authority;

		authority = {
			...this.state.authority,
		};

		let matches = [
			...authority.host_matches,
		];

		let i = matches.indexOf(match);
		if (i === -1) {
			return;
		}

		matches.splice(i, 1);

		authority.host_matches = matches;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addMatch: '',
			authority: authority,
		});
	}

	onAddSubnet = (): void => {
		let authority:  AuthorityTypes.Authority;

		authority = {
			...this.state.authority,
		};

		let subnets = [
			...(authority.host_subnets || []),
		];

		if (!this.state.addSubnet) {
			return;
		}

		if (subnets.indexOf(this.state.addSubnet) === -1) {
			subnets.push(this.state.addSubnet);
		}

		subnets.sort();

		authority.host_subnets = subnets;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addSubnet: '',
			authority: authority,
		});
	}

	onRemoveSubnet(subnet: string): void {
		let authority:  AuthorityTypes.Authority;

		authority = {
			...this.state.authority,
		};

		let subnets = [
			...authority.host_subnets,
		];

		let i = subnets.indexOf(subnet);
		if (i === -1) {
			return;
		}

		subnets.splice(i, 1);

		authority.host_subnets = subnets;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addSubnet: '',
			authority: authority,
		});
	}

	render(): JSX.Element {
		let authority: AuthorityTypes.Authority = this.state.authority;
		let isHsm = authority.type === 'pritunl_hsm';

		let roles: JSX.Element[] = [];
		for (let role of (authority.roles || [])) {
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

		let matches: JSX.Element[] = [];
		for (let match of authority.host_matches || []) {
			matches.push(
				<div
					className="bp5-tag bp5-tag-removable bp5-intent-primary"
					style={css.item}
					key={match}
				>
					{match}
					<button
						className="bp5-tag-remove"
						onMouseUp={(): void => {
							this.onRemoveMatch(match);
						}}
					/>
				</div>,
			);
		}

		let subnets: JSX.Element[] = [];
		for (let subnet of authority.host_subnets || []) {
			subnets.push(
				<div
					className="bp5-tag bp5-tag-removable bp5-intent-primary"
					style={css.item}
					key={subnet}
				>
					{subnet}
					<button
						className="bp5-tag-remove"
						onMouseUp={(): void => {
							this.onRemoveSubnet(subnet);
						}}
					/>
				</div>,
			);
		}

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
							help="Name of authority"
							type="text"
							placeholder="Enter name"
							value={authority.name}
							onChange={(val): void => {
								this.set('name', val);
							}}
						/>
						<PageSelect
							label="Type"
							help="Authority type"
							value={authority.type}
							onChange={(val): void => {
								this.set('type', val);
							}}
						>
							<option value="local">Local</option>
						</PageSelect>
						<PageSwitch
							label="Host certificates"
							help="Allow servers to validate and sign SSH host keys. This should be disabled for most configurations."
							checked={authority.host_certificates}
							onToggle={(): void => {
								this.toggle('host_certificates');
							}}
						/>
						<PageSwitch
							label="Strict host checking"
							help="Enable strict host checking for SSH clients connecting to servers in this domain."
							hidden={!authority.host_certificates}
							checked={authority.strict_host_checking}
							onToggle={(): void => {
								this.toggle('strict_host_checking');
							}}
						/>
						<PageInput
							label="Host Domain"
							help="Domain that will be used for SSH host certificates. All servers must have a subdomain registered on this domain. This should be empty for most configurations."
							type="text"
							placeholder="Host domain"
							value={authority.host_domain}
							onChange={(val): void => {
								this.set("host_domain", val);
							}}
						/>
						<PageSwitch
							label="Automatic bastion server"
							help="Enable automatic bastion servers on nodes using Docker containers. This should be disabled for most configurations."
							checked={authority.proxy_hosting}
							onToggle={(): void => {
								this.toggle('proxy_hosting');
							}}
						/>
						<label className="bp5-label"
							style={css.label}
							hidden={!authority.proxy_hosting}
						>
							Bastion Hostname and Port
							<Help
								title="Bastion Hostname and Port"
								content="Hostname of bastion server and port that SSH nodes will run on. This port cannot be 22 or conflict with existing services on the Pritunl Zero node. Each authority must have a unique bastion port. The bastion hostname will need to point to a Pritunl Zero bastion node or network load balancer in front of Pritunl Zero bastion nodes."
							/>
							<div className="bp5-control-group" style={css.inputGroup}>
								<input
									className="bp5-input"
									style={css.hostname}
									type="text"
									autoCapitalize="off"
									spellCheck={false}
									placeholder="Hostname"
									value={authority.proxy_hostname}
									onChange={(evt): void => {
										this.set('proxy_hostname', evt.target.value);
									}}
								/>
								<input
									className="bp5-input"
									style={css.port}
									type="text"
									autoCapitalize="off"
									spellCheck={false}
									placeholder="Port"
									value={authority.proxy_port || ''}
									onChange={(evt): void => {
										if (evt.target.value) {
											this.set('proxy_port', parseInt(evt.target.value, 10));
										} else {
											this.set('proxy_port', 0);
										}
									}}
								/>
							</div>
						</label>
						<PageInput
							hidden={authority.proxy_hosting}
							label="Bastion Host"
							help="Optional username and hostname of bastion host to proxy client connections for this domain. If the bastion station requires a specific username it must be included such as 'ec2-user@server.domain.com'. Bastion hostname does not need to be in host domain. If strict host checking is enabled bastion host must have a valid certificate. This should be empty for most configurations."
							type="text"
							placeholder="Bastion host"
							value={authority.host_proxy}
							onChange={(val): void => {
								this.set('host_proxy', val);
							}}
						/>
						<AuthorityDeploy
							disabled={this.state.disabled}
							nodes={this.props.nodes}
							authority={authority}
							proxy={false}
						/>
						<AuthorityDeploy
							hidden={authority.proxy_hosting || !authority.host_proxy}
							disabled={this.state.disabled || !authority.host_proxy}
							nodes={this.props.nodes}
							authority={authority}
							proxy={true}
						/>
					</div>
					<div style={css.group}>
						<PageInput
							hidden={authority.type !== 'pritunl_hsm'}
							label="HSM YubiKey Serial"
							help="Serial number of YubiKey that will be used to sign certificates. This number can be found on the back of the key."
							type="text"
							placeholder="HSM serial"
							value={authority.hsm_serial}
							onChange={(val): void => {
								this.set('hsm_serial', val);
							}}
						/>
						<PageInput
							hidden={!isHsm}
							readOnly={true}
							label="HSM Token"
							help="Pritunl HSM token."
							type="text"
							placeholder="Save to generate token"
							value={""}
						/>
						<PageSwitch
							hidden={!isHsm}
							label="Generate new HSM token and secret"
							help="Enable to generate a new token and secret on save. Secret can only be shown by generating new credentials."
							checked={authority.hsm_generate_secret}
							onToggle={(): void => {
								this.set('hsm_generate_secret', !authority.hsm_generate_secret);
							}}
						/>
						<PageInput
							label="Certificate Expire Minutes"
							help="Number of minutes until certificates expire. The certificate only needs to be active when initiating the SSH connection. The SSH connection will stay connected after the certificate expires. Must be greater then 1 and no more then 1440."
							type="text"
							placeholder="Certificate expire minutes"
							value={authority.expire}
							onChange={(val): void => {
								this.set('expire', parseInt(val, 10));
							}}
						/>
						<PageInput
							label="Host Certificate Expire Minutes"
							help="Number of minutes until host certificates expire. Must be greater then 14 and no more then 1440."
							type="text"
							placeholder="Host certificate expire minutes"
							hidden={!authority.host_certificates}
							value={authority.host_expire}
							onChange={(val): void => {
								this.set('host_expire', parseInt(val, 10));
							}}
						/>
						<PageSelect
							label="SSH Key ID Format"
							help="Format of the key ID field in the users SSH certificate. The user ID will include the users database ID. The username option will include the users name. The username ID option will include the user ID and username in the format userid-username. The username strip domain option will remove all characters after @ in the username."
							value={authority.key_id_format}
							onChange={(val): void => {
								this.set('key_id_format', val);
							}}
						>
							<option value="user_id">User ID</option>
							<option value="username">Username</option>
							<option value="username_id">Username ID</option>
							<option
								value="username_strip_domain"
							>Username Strip Domain</option>
						</PageSelect>
						<PageSwitch
							label="Match roles"
							help="Require a matching role with the user before giving a certificate. If disabled all users will be given a certificate from this authority. The certificate principles will only contain the users roles."
							checked={authority.match_roles}
							onToggle={(): void => {
								this.toggle('match_roles');
							}}
						/>
						<label className="bp5-label" hidden={!authority.match_roles}>
							Roles
							<Help
								title="Roles"
								content="Roles associated with this authority. If at least one role matches the user will be given a certificate from this authority. The certificate principles will only contain the users roles."
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
							hidden={!authority.match_roles}
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
							Custom Matches
							<Help
								title="Custom Matches"
								content="Custom domains that will be proxied through the bastion host. This should be empty for most configurations."
							/>
							<div>
								{matches}
							</div>
						</label>
						<PageInputButton
							buttonClass="bp5-intent-success bp5-icon-add"
							label="Add"
							type="text"
							placeholder="Add match"
							value={this.state.addMatch}
							onChange={(val): void => {
								this.setState({
									...this.state,
									addMatch: val,
								});
							}}
							onSubmit={this.onAddMatch}
						/>
						<label className="bp5-label">
							Match Subnets
							<Help
								title="Match Subnets"
								content="Subnets that will be proxied through the bastion host. All hosts in the subnets must be accessible from the bastion host. For best security match only private subnets in the same network as the bastion host. Currently only /8, /16, /24 and /32 subnets are supported. This should be empty for most configurations."
							/>
							<div>
								{subnets}
							</div>
						</label>
						<PageInputButton
							buttonClass="bp5-intent-success bp5-icon-add"
							label="Add"
							type="text"
							placeholder="Add subnet"
							value={this.state.addSubnet}
							onChange={(val): void => {
								this.setState({
									...this.state,
									addSubnet: val,
								});
							}}
							onSubmit={this.onAddSubnet}
						/>
					</div>
				</div>
				<PageCreate
					style={css.save}
					hidden={!this.state.authority}
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
