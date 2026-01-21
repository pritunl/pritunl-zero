/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SecretTypes from '../types/SecretTypes';
import * as SecretActions from '../actions/SecretActions';
import * as MiscUtils from '../utils/MiscUtils';
import PageInput from './PageInput';
import PageSelect from './PageSelect';
import PageInfo from './PageInfo';
import PageTextArea from './PageTextArea';
import PageSave from './PageSave';
import ConfirmButton from './ConfirmButton';
import Help from './Help';
import * as Constants from "../Constants";

interface Props {
	secret: SecretTypes.SecretRo;
	selected: boolean;
	onSelect: (shift: boolean) => void;
	onClose: () => void;
}

interface State {
	disabled: boolean;
	changed: boolean;
	message: string;
	secret: SecretTypes.Secret;
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
		cursor: 'pointer',
		position: 'absolute',
		top: 0,
		left: 0,
		right: 0,
		padding: '4px',
		height: '39px',
	} as React.CSSProperties,
	select: {
		margin: '7px 0px 0px 6px',
		paddingTop: '3px',
	} as React.CSSProperties,
};

export default class SecretDetailed extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
			secret: null,
		};
	}

	set(name: string, val: any): void {
		let secret: any;

		if (this.state.changed) {
			secret = {
				...this.state.secret,
			};
		} else {
			secret = {
				...this.props.secret,
			};
		}

		secret[name] = val;

		this.setState({
			...this.state,
			changed: true,
			secret: secret,
		});
	}

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		SecretActions.commit(this.state.secret).then((): void => {
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
						secret: null,
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
		SecretActions.remove(this.props.secret.id).then((): void => {
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

	render(): JSX.Element {
		let secr: SecretTypes.Secret = this.state.secret ||
			this.props.secret;

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

		switch (secr.type) {
			case "aws":
			case "":
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
		}

		return <td
			className="bp5-cell"
			colSpan={2}
			style={css.card}
		>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<div
						className="layout horizontal tab-close bp5-card-header"
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
							className="bp5-minimal bp5-intent-danger bp5-icon-trash"
							style={css.button}
							safe={true}
							progressClassName="bp5-intent-danger"
							dialogClassName="bp5-intent-danger bp5-icon-delete"
							dialogLabel="Delete Secret"
							confirmMsg="Permanently delete this secret"
							confirmInput={true}
							items={[secr.name]}
							disabled={this.state.disabled}
							onConfirm={this.onDelete}
						/>
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
					<PageInfo
						fields={[
							{
								label: 'ID',
								value: this.props.secret.id || 'None',
							},
						]}
					/>
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
			<PageSave
				style={css.save}
				hidden={!this.state.secret}
				message={this.state.message}
				changed={this.state.changed}
				disabled={this.state.disabled}
				light={true}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						secret: null,
					});
				}}
				onSave={this.onSave}
			/>
		</td>;
	}
}
