/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SecretTypes from '../types/SecretTypes';
import * as SecretActions from '../actions/SecretActions';
import * as MiscUtils from '../utils/MiscUtils';
import PageInput from './PageInput';
import PageSelect from './PageSelect';
import PageInfo from './PageInfo';
import PageTextArea from './PageTextArea';
import PageCreate from './PageCreate';
import Help from './Help';
import * as Constants from "../Constants";

interface Props {
	onClose: () => void;
}

interface State {
	closed: boolean;
	disabled: boolean;
	changed: boolean;
	message: string;
	secret: SecretTypes.Secret;
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

export default class SecretNew extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			closed: false,
			disabled: false,
			changed: false,
			message: '',
			secret: {
				name: 'New Secret'
			},
		};
	}

	set(name: string, val: any): void {
		let secret: any = {
			...this.state.secret,
		};

		secret[name] = val;

		this.setState({
			...this.state,
			changed: true,
			secret: secret,
		});
	}

	onCreate = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});

		let secret: any = {
			...this.state.secret,
		};

		SecretActions.create(secret).then((): void => {
			this.setState({
				...this.state,
				message: 'Secret created successfully',
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

	render(): JSX.Element {
		let secr: SecretTypes.Secret = this.state.secret;

		let keyLabel = "";
		let keyHelp = "";
		let keyPlaceholder = "";
		let valLabel = "";
		let valHelp = "";
		let valPlaceholder = "";
		let regionLabel = "";
		let regionHelp = "";
		let regionPlaceholder = "";
		let publicKeyLabel = "";
		let publicKeyHelp = "";
		let publicKeyPlaceholder = "";

		switch (secr.type || "aws") {
			case "aws":
				keyLabel = "AWS Key ID";
				keyHelp = "Key for AWS API authentication.";
				keyPlaceholder = "Key ID";
				valLabel = "AWS Secret ID";
				valHelp = "Key secret for AWS API authentication.";
				valPlaceholder = "Key ID";
				regionLabel = "AWS Region";
				regionHelp = "Region for AWS API.";
				regionPlaceholder = "Region";
				publicKeyLabel = "";
				publicKeyHelp = "";
				publicKeyPlaceholder = "";
				break;
			case "cloudflare":
				keyLabel = "Cloudflare Token";
				keyHelp = "Cloudflare API token.";
				keyPlaceholder = "Token";
				valLabel = "";
				valHelp = "";
				valPlaceholder = "";
				regionLabel = "";
				regionHelp = "";
				regionPlaceholder = "";
				publicKeyLabel = "";
				publicKeyHelp = "";
				publicKeyPlaceholder = "";
				break;
			case "oracle_cloud":
				keyLabel = "Oracle Cloud Tenancy OCID";
				keyHelp = "Tenancy OCID for Oracle Cloud API authentication.";
				keyPlaceholder = "Tenancy OCID";
				valLabel = "Oracle Cloud User OCID";
				valHelp = "User OCID for Oracle Cloud API authentication.";
				valPlaceholder = "User OCID";
				regionLabel = "Oracle Cloud Region";
				regionHelp = "Region for Oracle Cloud API.";
				regionPlaceholder = "Region";
				publicKeyLabel = "Oracle Cloud Public Key";
				publicKeyHelp = "Public key for Oracle Cloud API authentication.";
				publicKeyPlaceholder = "Oracle Cloud Public Key";
				break;
			case "gcp":
				keyLabel = "GCP Service Account JSON";
				keyHelp = "Service account JSON credentials for GCP API authentication.";
				keyPlaceholder = "Service Account JSON";
				valLabel = "";
				valHelp = "";
				valPlaceholder = "";
				regionLabel = "";
				regionHelp = "";
				regionPlaceholder = "";
				publicKeyLabel = "";
				publicKeyHelp = "";
				publicKeyPlaceholder = "";
				break;
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
							help="Name of secret"
							type="text"
							placeholder="Name"
							value={secr.name}
							onChange={(val): void => {
								this.set('name', val);
							}}
						/>
						<PageTextArea
							label="Comment"
							help="Secret comment."
							placeholder="Secret comment"
							rows={3}
							value={secr.comment}
							onChange={(val: string): void => {
								this.set('comment', val);
							}}
						/>
						{secr.type === "gcp" ? (
							<PageTextArea
								label={keyLabel}
								help={keyHelp}
								placeholder={keyPlaceholder}
								rows={10}
								value={secr.key}
								onChange={(val: string): void => {
									this.set('key', val);
								}}
							/>
						) : (
							<PageInput
								label={keyLabel}
								help={keyHelp}
								hidden={keyLabel === ""}
								type="text"
								placeholder={keyPlaceholder}
								value={secr.key}
								onChange={(val: string): void => {
									this.set('key', val);
								}}
							/>
						)}
						<PageInput
							label={valLabel}
							help={valHelp}
							hidden={valLabel === ""}
							type="text"
							placeholder={valPlaceholder}
							value={secr.value}
							onChange={(val: string): void => {
								this.set('value', val);
							}}
						/>
						<PageInput
							label={regionLabel}
							help={regionHelp}
							hidden={regionLabel === ""}
							type="text"
							placeholder={regionPlaceholder}
							value={secr.region}
							onChange={(val: string): void => {
								this.set('region', val);
							}}
						/>
					</div>
					<div style={css.group}>
						<PageSelect
							label="Type"
							disabled={this.state.disabled}
							help="Secret provider."
							value={secr.type}
							onChange={(val): void => {
								this.set('type', val);
							}}
						>
							<option value="aws">AWS</option>
							<option value="cloudflare">Cloudflare</option>
							<option value="oracle_cloud">Oracle Cloud</option>
							<option value="gcp">GCP</option>
						</PageSelect>
						<PageTextArea
							disabled={this.state.disabled}
							hidden={publicKeyLabel === ""}
							label={publicKeyLabel}
							help={publicKeyHelp}
							placeholder={publicKeyPlaceholder}
							readOnly={true}
							rows={6}
							value={secr.public_key}
							onChange={(val: string): void => {
								this.set('public_key', val);
							}}
						/>
					</div>
				</div>
				<PageCreate
					style={css.save}
					hidden={!this.state.secret}
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
