/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';
import * as AuthorityTypes from '../types/AuthorityTypes';
import Help from './Help';
import PageInput from './PageInput';
import PageInputButton from './PageInputButton';
import PageSwitch from './PageSwitch';

interface Props {
	disabled?: boolean;
	authority: AuthorityTypes.AuthorityRo;
}

interface State {
	popover: boolean;
	hostCertificate: boolean;
	hostname: string;
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
			hostCertificate: null,
			hostname: '',
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
		let content = '';
		let hostCertificate = this.state.hostCertificate;
		if (hostCertificate === null) {
			hostCertificate = this.props.authority.host_certificates;
		}
		if (!this.props.authority.host_certificates ||
				!this.props.authority.host_tokens.length) {
			hostCertificate = false;
		}

		let hostname = this.state.hostname ||
			(this.props.authority.host_domain ?
				'.' + this.props.authority.host_domain : '');

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

		if (hostCertificate) {
			content = `sudo sed -i '/^TrustedUserCAKeys/d' /etc/ssh/sshd_config
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
rm -f key.tmp
sudo yum -y install pritunl-ssh-host

sudo pritunl-ssh-host config add-token ${
	this.props.authority.host_tokens.length ?
	this.props.authority.host_tokens[0] : 'HOST_TOKEN_UNAVAILABLE'}
sudo pritunl-ssh-host config hostname ${hostname}
sudo pritunl-ssh-host config server ${window.location.host}

sudo systemctl restart sshd || true
sudo service sshd restart || true`;
		} else {
			content = `sudo sed -i '/^TrustedUserCAKeys/d' /etc/ssh/sshd_config
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

		if (this.state.popover) {
			popoverElem = <Blueprint.Dialog
				title="Deploy Test"
				style={css.dialog}
				isOpen={this.state.popover}
				onClose={(): void => {
					this.setState({
						...this.state,
						popover: false,
					});
				}}
			>
				<div className="pt-dialog-body">
					<PageSwitch
						label="Host certificate"
						hidden={!this.props.authority.host_certificates}
						disabled={!this.props.authority.host_tokens.length}
						help="Provision a host certificate to this server, requires installing Pritunl Zero host client. Authority must have at least one host token."
						checked={hostCertificate}
						onToggle={(): void => {
							this.setState({
								...this.state,
								hostCertificate: !hostCertificate,
							});
						}}
					/>
					<PageInput
						label="Server Hostname"
						hidden={!hostCertificate}
						help="Hostname of the server. The Pritunl Zero server must be able to resolve the server using this hostname to provision the host certificate. The hostname must be a subdomain of the authority domain."
						type="text"
						placeholder="Server hostname"
						value={hostname}
						onChange={(val): void => {
							this.setState({
								...this.state,
								hostname: val,
							});
						}}
					/>
					<label className="pt-label">
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

		return <div style={css.box}>
			<button
				className="pt-button pt-icon-cloud-upload pt-intent-primary"
				style={css.button}
				type="button"
				disabled={this.props.disabled}
				onClick={(): void => {
					this.setState({
						...this.state,
						popover: !this.state.popover,
					});
				}}
			>
				Deploy Script
			</button>
			{popoverElem}
		</div>;
	}
}
