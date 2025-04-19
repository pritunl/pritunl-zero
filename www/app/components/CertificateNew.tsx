/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as CertificateTypes from '../types/CertificateTypes';
import * as SecretTypes from '../types/SecretTypes';
import * as CertificateActions from '../actions/CertificateActions';
import * as MiscUtils from '../utils/MiscUtils';
import CertificateDomain from './CertificateDomain';
import PageInput from './PageInput';
import PageSelect from './PageSelect';
import PageInfo from './PageInfo';
import PageTextArea from './PageTextArea';
import PageCreate from './PageCreate';
import ConfirmButton from './ConfirmButton';
import Help from './Help';
import * as Constants from "../Constants";

interface Props {
	secrets: SecretTypes.SecretsRo;
	onClose: () => void;
}

interface State {
	closed: boolean;
	disabled: boolean;
	changed: boolean;
	message: string;
	certificate: CertificateTypes.Certificate;
	addDomain: string;
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
	domain: {
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
		minWidth: '280px',
		margin: '0 10px',
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

export default class CertificateDetailed extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			closed: false,
			disabled: false,
			changed: false,
			message: '',
			addDomain: null,
			certificate: {
				name: 'New Certificate',
			},
		};
	}

	set(name: string, val: any): void {
		let certificate: any = {
			...this.state.certificate,
		};

		certificate[name] = val;

		this.setState({
			...this.state,
			changed: true,
			certificate: certificate,
		});
	}

	onCreate = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});

		let certificate: any = {
			...this.state.certificate,
		};

		CertificateActions.create(certificate).then((): void => {
			this.setState({
				...this.state,
				message: 'Certificate created successfully',
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

	onAddDomain = (): void => {
		let cert: CertificateTypes.Certificate;

		cert = {
			...this.state.certificate,
		};

		let acmeDomains = [
			...cert.acme_domains,
			'',
		];

		cert.acme_domains = acmeDomains;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addDomain: '',
			certificate: cert,
		});
	}

	onChangeDomain(i: number, state: string): void {
		let cert: CertificateTypes.Certificate;

		cert = {
			...this.state.certificate,
		};

		let acmeDomains = [
			...cert.acme_domains,
		];

		acmeDomains[i] = state;

		cert.acme_domains = acmeDomains;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			certificate: cert,
		});
	}

	onRemoveDomain(i: number): void {
		let cert: CertificateTypes.Certificate;

		cert = {
			...this.state.certificate,
		};

		let acmeDomains = [
			...cert.acme_domains,
		];

		acmeDomains.splice(i, 1);

		cert.acme_domains = acmeDomains;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addDomain: '',
			certificate: cert,
		});
	}

	render(): JSX.Element {
		let cert: CertificateTypes.Certificate = this.state.certificate;

		let hasSecrets = false;
		let secretsSelect: JSX.Element[] = [];
		if (this.props.secrets.length) {
			secretsSelect.push(<option key="null" value="">Select Secret</option>);

			for (let secret of this.props.secrets) {
				hasSecrets = true;
				secretsSelect.push(
					<option
						key={secret.id}
						value={secret.id}
					>{secret.name}</option>,
				);
			}
		}

		if (!hasSecrets) {
			secretsSelect = [<option key="null" value="">No Secrets</option>];
		}

		let domains: JSX.Element[] = [];
		(cert.acme_domains || []).forEach((acmeDomain, index) => {
			domains.push(
				<CertificateDomain
					key={index}
					disabled={this.state.disabled}
					domain={acmeDomain}
					onChange={(state: string): void => {
						this.onChangeDomain(index, state);
					}}
					onRemove={(): void => {
						this.onRemoveDomain(index);
					}}
				/>,
			);
		})

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
							help="Name of certificate"
							type="text"
							placeholder="Name"
							value={cert.name}
							onChange={(val): void => {
								this.set('name', val);
							}}
						/>
						<PageTextArea
							label="Comment"
							help="Certificate comment."
							placeholder="Certificate comment"
							rows={3}
							value={cert.comment}
							onChange={(val: string): void => {
								this.set('comment', val);
							}}
						/>
						<PageTextArea
							readOnly={cert.type !== 'text'}
							label="Private Key"
							help="Certificate private key in PEM format"
							placeholder="Private key"
							rows={6}
							value={cert.key}
							onChange={(val: string): void => {
								this.set('key', val);
							}}
						/>
						<PageTextArea
							readOnly={cert.type !== 'text'}
							label="Certificate Chain"
							help="Certificate followed by certificate chain in PEM format"
							placeholder="Certificate chain"
							rows={6}
							value={cert.certificate}
							onChange={(val: string): void => {
								this.set('certificate', val);
							}}
						/>
						<label
							style={css.itemsLabel}
							hidden={cert.type !== 'lets_encrypt'}
						>
							LetsEncrypt Domains
							<Help
								title="LetsEncrypt Domains"
								content="Enter domain names for the certificate. All domains names must point to a Pritunl Cloud server in the cluster. The servers must also have port 80 publicy open. The port will need to stay open to renew the certificate."
							/>
						</label>
						<div hidden={cert.type !== 'lets_encrypt'}>
							{domains}
						</div>
						<button
							className="bp5-button bp5-intent-success bp5-icon-add"
							disabled={this.state.disabled}
							style={css.itemsAdd}
							hidden={cert.type !== 'lets_encrypt'}
							type="button"
							onClick={this.onAddDomain}
						>
							Add Domain
						</button>
					</div>
					<div style={css.group}>
						<PageSelect
							label="Type"
							disabled={this.state.disabled}
							help="Certificate type, use text to provide a certificate. LetsEncrypt provides free certificates that automatically renew."
							value={cert.type}
							onChange={(val): void => {
								this.set('type', val);
							}}
						>
							<option value="text">Text</option>
							<option value="lets_encrypt">LetsEncrypt</option>
						</PageSelect>
						<PageSelect
							label="LetsEncrypt Verification Type"
							disabled={this.state.disabled}
							hidden={cert.type != "lets_encrypt"}
							help="Verification type for LetsEncrypt certificate. HTTP verification will use a HTTP request on port 80 from the host. DNS will use a DNS API provider to set a DNS TXT record."
							value={cert.acme_type}
							onChange={(val): void => {
								this.set('acme_type', val);
							}}
						>
							<option value="acme_http">HTTP</option>
							<option value="acme_dns">DNS TXT</option>
						</PageSelect>
						<PageSelect
							label="LetsEncrypt Verification Provider"
							disabled={this.state.disabled}
							hidden={cert.acme_type != "acme_dns"}
							help="API provider for LetsEncrypt verification."
							value={cert.acme_auth}
							onChange={(val): void => {
								this.set('acme_auth', val);
							}}
						>
							<option value="acme_aws">AWS</option>
							<option value="acme_cloudflare">Cloudflare</option>
							<option value="acme_oracle_cloud">Oracle Cloud</option>
						</PageSelect>
						<PageSelect
							disabled={this.state.disabled}
							hidden={cert.acme_type != "acme_dns"}
							label="LetsEncrypt Verification Secret"
							help="Secret containing API keys to use for LetsEncrypt verification."
							value={cert.acme_secret}
							onChange={(val): void => {
								this.set('acme_secret', val);
							}}
						>
							{secretsSelect}
						</PageSelect>
					</div>
				</div>
				<PageCreate
					style={css.save}
					hidden={!this.state.certificate}
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
