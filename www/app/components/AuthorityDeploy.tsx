/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';
import * as AuthorityTypes from '../types/AuthorityTypes';
import Help from './Help';
import PageInput from './PageInput';
import PageInputButton from './PageInputButton';
import PageSwitch from './PageSwitch';
import PageSelect from './PageSelect';
import * as NodeTypes from "../types/NodeTypes";

interface Props {
	disabled?: boolean;
	nodes: NodeTypes.NodesRo;
	authority: AuthorityTypes.AuthorityRo;
	proxy: boolean;
}

interface State {
	popover: boolean;
	route53: boolean;
	awsAccessKey: string;
	awsSecretKey: string;
	hostCertificate: boolean;
	hostname: string;
	server: string;
	addRole: string;
	roles: string[];
}

const css = {
	box: {
		marginBottom: '15px',
	} as React.CSSProperties,
	button: {
	} as React.CSSProperties,
	item: {
		margin: '9px 5px 0 5px',
		height: '20px',
	} as React.CSSProperties,
	callout: {
		marginBottom: '15px',
	} as React.CSSProperties,
	popover: {
		width: '230px',
	} as React.CSSProperties,
	popoverTarget: {
		top: '9px',
		left: '18px',
	} as React.CSSProperties,
	dialog: {
		maxWidth: '480px',
		margin: '30px 20px',
	} as React.CSSProperties,
	textarea: {
		width: '100%',
		resize: 'none',
		fontSize: '12px',
		fontFamily: '"Lucida Console", Monaco, monospace',
	} as React.CSSProperties,
};

export default class AuthorityDeploy extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			popover: false,
			route53: false,
			awsAccessKey: '',
			awsSecretKey: '',
			hostCertificate: null,
			hostname: '',
			server: null,
			addRole: '',
			roles: [],
		};
	}

	onAddRole = (): void => {
		let roles = [
			...this.state.roles,
		];

		if (!this.state.addRole) {
			return;
		}

		if (roles.indexOf(this.state.addRole) === -1) {
			roles.push(this.state.addRole);
		}

		roles.sort();

		this.setState({
			...this.state,
			addRole: '',
			roles: roles,
		});
	}

	onRemoveRole(role: string): void {
		let roles = [
			...this.state.roles,
		];

		let i = roles.indexOf(role);
		if (i === -1) {
			return;
		}

		roles.splice(i, 1);

		this.setState({
			...this.state,
			addRole: '',
			roles: roles,
		});
	}

	render(): JSX.Element {
		let popoverElem: JSX.Element;

		if (this.state.popover) {
			let content = '';
			let callout = 'Use the startup script below to provision a ' +
				'Pritunl Zero host.';
			let errorMsg = '';
			let errorMsgElem: JSX.Element;
			let hostCertificate = this.state.hostCertificate;
			let hostCertificateDisabled = false;
			if (hostCertificate === null) {
				hostCertificate = this.props.authority.host_certificates;
			}

			let servers = new Set();
			let serverDefault: string = null;
			let serversElm: JSX.Element[] = [];
			if (this.props.nodes) {
				for (let node of this.props.nodes) {
					if (node.user_domain) {
						servers.add(node.user_domain);
					}
				}
			}

			if (!this.props.authority.host_tokens.length || servers.size === 0) {
				hostCertificate = false;
				hostCertificateDisabled = true;
			}

			servers.forEach((server): void => {
				if (!serverDefault) {
					serverDefault = server;
				}
				serversElm.push(<option value={server}>{server}</option>);
			});
			if (servers.size === 1) {
				serversElm = [];
			}

			let bastionUsername = '';
			let bastionHostname = '';
			if (this.props.proxy) {
				let bastionSplit = this.props.authority.host_proxy.split('@');
				if (bastionSplit.length === 2) {
					bastionUsername = this.props.authority.host_proxy.split('@')[0];
					if (bastionSplit[1].indexOf(
							this.props.authority.host_domain) !== -1) {
						bastionHostname = bastionSplit[1].replace(
							'.' + this.props.authority.host_domain, '');
					}
				}

				if (!bastionUsername) {
					errorMsg = 'Bastion host is missing username.';
				} else if (!bastionHostname) {
					errorMsg = 'Bastion hostname is not a subdomain of host domain.';
				}
			}

			let epel = '';
			let boto = '';
			let route53 = '';
			if (this.state.route53 && hostCertificate) {
				epel = '\nsudo yum -y install epel-release || ' +
					'sudo rpm -Uvh https://dl.fedoraproject.org/' +
					'pub/epel/epel-release-latest-7.noarch.rpm';
				boto = ' python2-boto3 python27-boto3';
				if (this.state.awsAccessKey) {
					route53 += '\nsudo pritunl-ssh-host config aws-access-key ' +
						this.state.awsAccessKey;
				}
				if (this.state.awsSecretKey) {
					route53 += '\nsudo pritunl-ssh-host config aws-secret-key ' +
						this.state.awsSecretKey;
				}
				route53 += '\nsudo pritunl-ssh-host config route-53-zone ' +
					this.props.authority.host_domain;
			}

			let roles: JSX.Element[] = [];
			for (let role of this.state.roles) {
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

			if (this.props.proxy) {
				callout = 'Open port 9748 and use the startup script below to ' +
					'provision a Pritunl Zero host. Provisioning may take several ' +
					'minutes if the servers DNS record was created recently.';
				content = `#!/bin/bash
sudo sed -i '/^TrustedUserCAKeys/d' /etc/ssh/sshd_config
sudo sed -i '/^AuthorizedPrincipalsFile/d' /etc/ssh/sshd_config
sudo tee -a /etc/ssh/sshd_config << EOF

Match User ${bastionUsername}
	AllowAgentForwarding no
	AllowTcpForwarding yes
	PermitOpen *:22
	GatewayPorts no
	X11Forwarding no
	PermitTunnel no
	ForceCommand echo 'Pritunl Zero Bastion Host'
	TrustedUserCAKeys /etc/ssh/trusted
	AuthorizedPrincipalsFile /etc/ssh/principals
Match all
EOF
sudo tee /etc/ssh/principals << EOF
bastion
EOF
sudo tee /etc/ssh/trusted << EOF
${this.props.authority.public_key}
EOF

sudo tee -a /etc/yum.repos.d/pritunl.repo << EOF
[pritunl]
name=Pritunl Repository
baseurl=https://repo.pritunl.com/stable/yum/centos/7/
gpgcheck=1
enabled=1
EOF

gpg --keyserver hkp://keyserver.ubuntu.com --recv-keys 7568D9BB55FF9E5287D586017AE645C0CF8E292A
gpg --armor --export 7568D9BB55FF9E5287D586017AE645C0CF8E292A > key.tmp
sudo rpm --import key.tmp
rm -f key.tmp${epel}
sudo yum -y install pritunl-ssh-host${boto}
${route53}
sudo pritunl-ssh-host config add-token ${
	this.props.authority.host_tokens.length ?
	this.props.authority.host_tokens[0] : 'HOST_TOKEN_UNAVAILABLE'}
sudo pritunl-ssh-host config hostname ${bastionHostname}
sudo pritunl-ssh-host config server ${this.state.server || serverDefault}
sudo useradd ${bastionUsername} || true

sudo systemctl restart sshd || true
sudo service sshd restart || true`;
			} else if (hostCertificate) {
				callout = 'Open port 9748 and use the startup script below to ' +
					'provision a Pritunl Zero host. Provisioning may take several ' +
					'minutes if the servers DNS record was created recently.';
				content = `#!/bin/bash
sudo sed -i '/^TrustedUserCAKeys/d' /etc/ssh/sshd_config
sudo sed -i '/^AuthorizedPrincipalsFile/d' /etc/ssh/sshd_config
sudo tee -a /etc/ssh/sshd_config << EOF

TrustedUserCAKeys /etc/ssh/trusted
AuthorizedPrincipalsFile /etc/ssh/principals
EOF
sudo tee /etc/ssh/principals << EOF
emergency
${this.state.roles.length ? this.state.roles.join('\n') + '\n' : ''}EOF
sudo tee /etc/ssh/trusted << EOF
${this.props.authority.public_key}
EOF

sudo tee -a /etc/yum.repos.d/pritunl.repo << EOF
[pritunl]
name=Pritunl Repository
baseurl=https://repo.pritunl.com/stable/yum/centos/7/
gpgcheck=1
enabled=1
EOF

gpg --keyserver hkp://keyserver.ubuntu.com --recv-keys 7568D9BB55FF9E5287D586017AE645C0CF8E292A
gpg --armor --export 7568D9BB55FF9E5287D586017AE645C0CF8E292A > key.tmp
sudo rpm --import key.tmp
rm -f key.tmp${epel}
sudo yum -y install pritunl-ssh-host${boto}
${route53}
sudo pritunl-ssh-host config add-token ${
	this.props.authority.host_tokens.length ?
	this.props.authority.host_tokens[0] : 'HOST_TOKEN_UNAVAILABLE'}
sudo pritunl-ssh-host config hostname ${this.state.hostname}
sudo pritunl-ssh-host config server ${this.state.server || serverDefault}

sudo systemctl restart sshd || true
sudo service sshd restart || true`;
			} else {
				content = `#!/bin/bash
sudo sed -i '/^TrustedUserCAKeys/d' /etc/ssh/sshd_config
sudo sed -i '/^AuthorizedPrincipalsFile/d' /etc/ssh/sshd_config
sudo tee -a /etc/ssh/sshd_config << EOF

TrustedUserCAKeys /etc/ssh/trusted
AuthorizedPrincipalsFile /etc/ssh/principals
EOF
sudo tee /etc/ssh/principals << EOF
emergency
${this.state.roles.length ? this.state.roles.join('\n') + '\n' : ''}EOF
sudo tee /etc/ssh/trusted << EOF
${this.props.authority.public_key}
EOF

sudo systemctl restart sshd || true
sudo service sshd restart || true`;
			}

			if (errorMsg) {
				errorMsgElem = <div className="pt-dialog-body">
					<div
						className="pt-callout pt-intent-danger pt-icon-ban-circle"
						style={css.callout}
					>
						{errorMsg}
					</div>
				</div>;
			}

			let title = '';
			if (this.props.proxy) {
				title = 'Generate Bastion Deploy Script';
			} else {
				title = 'Generate Deploy Script';
			}

			popoverElem = <Blueprint.Dialog
				title={title}
				style={css.dialog}
				isOpen={this.state.popover}
				onClose={(): void => {
					this.setState({
						...this.state,
						popover: false,
					});
				}}
			>
				{errorMsgElem}
				<div className="pt-dialog-body" hidden={!!errorMsgElem}>
					<div
						className="pt-callout pt-intent-primary pt-icon-info-sign"
						style={css.callout}
					>
						{callout}
					</div>
					<PageSwitch
						label="Host certificate"
						hidden={!this.props.authority.host_certificates ||
							this.props.proxy}
						disabled={hostCertificateDisabled}
						help="Provision a host certificate to this server, requires installing Pritunl Zero host client. Authority must have at least one host token and at least one node must have a user domain."
						checked={hostCertificate}
						onToggle={(): void => {
							this.setState({
								...this.state,
								hostCertificate: !hostCertificate,
							});
						}}
					/>
					<PageSelect
						hidden={!hostCertificate || serversElm.length === 0 ||
							this.props.proxy}
						label="Pritunl Zero Server"
						help="The Pritunl Zero server hostname that the client will authenticate from."
						value={this.state.server || serverDefault}
						onChange={(val): void => {
							this.setState({
								...this.state,
								server: val,
							});
						}}
					>
						{serversElm}
					</PageSelect>
					<PageInput
						label="Server Hostname"
						hidden={!hostCertificate || this.props.proxy}
						help="Hostname portion of the server domain. The Pritunl Zero server must be able to resolve the server using this hostname to provision the host certificate. The hostname will be combined with the authority domain to form the servers domain."
						type="text"
						placeholder="Server hostname"
						value={this.state.hostname}
						onChange={(val): void => {
							this.setState({
								...this.state,
								hostname: val,
							});
						}}
					/>
					<PageSwitch
						label="Auto Route53 configuration"
						hidden={!hostCertificate}
						help="Automatically update a Route53 record for this servers hostname. The authority domain must be hosted in Route53."
						checked={this.state.route53}
						onToggle={(): void => {
							this.setState({
								...this.state,
								route53: !this.state.route53,
							});
						}}
					/>
					<PageInput
						label="AWS Access Key"
						hidden={!hostCertificate || !this.state.route53}
						help="AWS access key for auto Route53 configuration. Leave blank if the instance is configured with an instance role."
						type="text"
						placeholder="Leave blank to use instance role"
						value={this.state.awsAccessKey}
						onChange={(val): void => {
							this.setState({
								...this.state,
								awsAccessKey: val,
							});
						}}
					/>
					<PageInput
						label="AWS Secret Key"
						hidden={!hostCertificate || !this.state.route53}
						help="AWS secret key for auto Route53 configuration. Leave blank if the instance is configured with an instance role."
						type="text"
						placeholder="Leave blank to use instance role"
						value={this.state.awsSecretKey}
						onChange={(val): void => {
							this.setState({
								...this.state,
								awsSecretKey: val,
							});
						}}
					/>
					<label
						className="pt-label"
						hidden={this.props.proxy}
					>
						Roles
						<Help
							title="Roles"
							content="Roles associated with this server. The user must have at least one matching role to access this server."
						/>
						<div>
							{roles}
						</div>
					</label>
					<PageInputButton
						buttonClass="pt-intent-success pt-icon-add"
						hidden={this.props.proxy}
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
					<textarea
						className="pt-input"
						style={css.textarea}
						readOnly={true}
						autoCapitalize="off"
						spellCheck={false}
						rows={18}
						value={content}
						onClick={(evt): void => {
							evt.currentTarget.select();
						}}
					/>
				</div>
				<div className="pt-dialog-footer">
					<div className="pt-dialog-footer-actions">
						<button
							className="pt-button"
							type="button"
							onClick={(): void => {
								this.setState({
									...this.state,
									popover: !this.state.popover,
								});
							}}
						>Close</button>
					</div>
				</div>
			</Blueprint.Dialog>;
		}

		let buttonLabel = '';
		if (this.props.proxy) {
			buttonLabel = 'Generate Bastion Deploy Script';
		} else {
			buttonLabel = 'Generate Deploy Script';
		}

		return <div style={css.box}>
			<button
				className="pt-button pt-icon-cloud-upload pt-intent-primary"
				style={css.button}
				type="button"
				disabled={this.props.disabled ||
					(this.props.proxy && (!this.props.authority.host_proxy ||
					!this.props.authority.host_certificates))}
				onClick={(): void => {
					this.setState({
						...this.state,
						popover: !this.state.popover,
					});
				}}
			>
				{buttonLabel}
			</button>
			{popoverElem}
		</div>;
	}
}
