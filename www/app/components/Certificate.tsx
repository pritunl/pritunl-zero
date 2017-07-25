/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as CertificateTypes from '../types/CertificateTypes';
import * as CertificateActions from '../actions/CertificateActions';
import CertificateDomain from './CertificateDomain';
import PageInput from './PageInput';
import PageSelect from './PageSelect';
import PageInfo from './PageInfo';
import PageSave from './PageSave';
import ConfirmButton from './ConfirmButton';

interface Props {
	certificate: CertificateTypes.CertificateRo;
}

interface State {
	disabled: boolean;
	changed: boolean;
	message: string;
	certificate: CertificateTypes.Certificate;
	addDomain: string;
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
};

export default class Certificate extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
			certificate: null,
			addDomain: null,
		};
	}

	set(name: string, val: any): void {
		let certificate: any;

		if (this.state.changed) {
			certificate = {
				...this.state.certificate,
			};
		} else {
			certificate = {
				...this.props.certificate,
			};
		}

		certificate[name] = val;

		this.setState({
			...this.state,
			changed: true,
			certificate: certificate,
		});
	}

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		CertificateActions.commit(this.state.certificate).then((): void => {
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
						certificate: null,
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
		CertificateActions.remove(this.props.certificate.id).then((): void => {
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

	onAddDomain = (): void => {
		let cert: CertificateTypes.Certificate = this.state.certificate ||
			this.props.certificate;

		let acmeDomains = [
			...cert.acme_domains,
			'',
		];

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addDomain: '',
			certificate: {
				...cert,
				acme_domains: acmeDomains,
			},
		});
	}

	onChangeDomain(i: number, state: string): void {
		let cert: CertificateTypes.Certificate = this.state.certificate ||
			this.props.certificate;

		let acmeDomains = [
			...cert.acme_domains,
		];

		acmeDomains[i] = state;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			certificate: {
				...cert,
				acme_domains: acmeDomains,
			},
		});
	}

	onRemoveDomain(i: number): void {
		let cert: CertificateTypes.Certificate = this.state.certificate ||
			this.props.certificate;

		let acmeDomains = [
			...cert.acme_domains,
		];

		acmeDomains.splice(i, 1);

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addDomain: '',
			certificate: {
				...cert,
				acme_domains: acmeDomains,
			},
		});
	}

	render(): JSX.Element {
		let cert: CertificateTypes.Certificate = this.state.certificate ||
			this.props.certificate;

		let domains: JSX.Element[] = [];
		for (let i = 0; i < cert.acme_domains.length; i++) {
			let index = i;

			domains.push(
				<CertificateDomain
					key={index}
					domain={cert.acme_domains[index]}
					onChange={(state: string): void => {
						this.onChangeDomain(index, state);
					}}
					onRemove={(): void => {
						this.onRemoveDomain(index);
					}}
				/>,
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
							confirmMsg="Confirm certificate remove"
							disabled={this.state.disabled}
							onConfirm={this.onDelete}
						/>
					</div>
					<PageInput
						label="Name"
						type="text"
						placeholder="Enter name"
						value={cert.name}
						onChange={(val): void => {
							this.set('name', val);
						}}
					/>
					<PageSelect
						label="Type"
						value={cert.type}
						onChange={(val): void => {
							this.set('type', val);
						}}
					>
						<option value="text">Text</option>
						<option value="lets_encrypt">LetsEncrypt</option>
					</PageSelect>
					<label
						style={css.itemsLabel}
						hidden={cert.type !== 'lets_encrypt'}
					>
						LetsEncrypt Domains
					</label>
					<div hidden={cert.type !== 'lets_encrypt'}>
						{domains}
					</div>
					<button
						className="pt-button pt-intent-success pt-icon-add"
						style={css.itemsAdd}
						hidden={cert.type !== 'lets_encrypt'}
						type="button"
						onClick={this.onAddDomain}
					>
						Add Domain
					</button>
				</div>
				<div style={css.group}>
					<PageInfo
						fields={[
							{
								label: 'ID',
								value: cert.id || 'None',
							},
						]}
					/>
				</div>
			</div>
			<PageSave
				style={css.save}
				hidden={!this.state.certificate}
				message={this.state.message}
				changed={this.state.changed}
				disabled={false}
				light={true}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						certificate: null,
					});
				}}
				onSave={this.onSave}
			/>
		</div>;
	}
}
