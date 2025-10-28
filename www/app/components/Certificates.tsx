/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as CertificateTypes from '../types/CertificateTypes';
import CertificatesStore from '../stores/CertificatesStore';
import * as CertificateActions from '../actions/CertificateActions';
import Certificate from './Certificate';
import CertificateNew from './CertificateNew';
import CertificatesFilter from './CertificatesFilter';
import CertificatesPage from './CertificatesPage';
import * as SecretTypes from '../types/SecretTypes';
import SecretsStore from '../stores/SecretsStore';
import * as SecretActions from '../actions/SecretActions';
import Page from './Page';
import PageHeader from './PageHeader';
import NonState from './NonState';
import ConfirmButton from './ConfirmButton';

interface Selected {
	[key: string]: boolean;
}

interface Opened {
	[key: string]: boolean;
}

interface State {
	certificates: CertificateTypes.CertificatesRo;
	secrets: SecretTypes.SecretsRo;
	filter: CertificateTypes.Filter;
	selected: Selected;
	opened: Opened;
	newOpened: boolean;
	lastSelected: string;
	disabled: boolean;
}

const css = {
	items: {
		width: '100%',
		marginTop: '-5px',
		display: 'table',
		tableLayout: 'fixed',
		borderSpacing: '0 5px',
	} as React.CSSProperties,
	itemsBox: {
		width: '100%',
		overflowY: 'auto',
	} as React.CSSProperties,
	placeholder: {
		opacity: 0,
		width: '100%',
	} as React.CSSProperties,
	header: {
		marginTop: '-19px',
	} as React.CSSProperties,
	heading: {
		margin: '19px 0 0 0',
	} as React.CSSProperties,
	button: {
		margin: '8px 0 0 8px',
	} as React.CSSProperties,
	buttons: {
		marginTop: '8px',
	} as React.CSSProperties,
};

export default class Certificates extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			certificates: CertificatesStore.certificates,
			secrets: SecretsStore.secretsName,
			filter: CertificatesStore.filter,
			selected: {},
			opened: {},
			newOpened: false,
			lastSelected: null,
			disabled: false,
		};
	}

	get selected(): boolean {
		return !!Object.keys(this.state.selected).length;
	}

	get opened(): boolean {
		return !!Object.keys(this.state.opened).length;
	}

	componentDidMount(): void {
		CertificatesStore.addChangeListener(this.onChange);
		SecretsStore.addChangeListener(this.onChange);
		CertificateActions.sync();
		SecretActions.syncNames();
	}

	componentWillUnmount(): void {
		CertificatesStore.removeChangeListener(this.onChange);
		SecretsStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		let certificates = CertificatesStore.certificates;
		let selected: Selected = {};
		let curSelected = this.state.selected;
		let opened: Opened = {};
		let curOpened = this.state.opened;

		certificates.forEach((certificate: CertificateTypes.Certificate): void => {
			if (curSelected[certificate.id]) {
				selected[certificate.id] = true;
			}
			if (curOpened[certificate.id]) {
				opened[certificate.id] = true;
			}
		});

		this.setState({
			...this.state,
			certificates: CertificatesStore.certificates,
			secrets: SecretsStore.secretsName,
			filter: CertificatesStore.filter,
			selected: selected,
			opened: opened,
		});
	}

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		CertificateActions.removeMulti(
				Object.keys(this.state.selected)).then((): void => {
			this.setState({
				...this.state,
				selected: {},
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
		let certificatesDom: JSX.Element[] = [];

		this.state.certificates.forEach((
				certificate: CertificateTypes.CertificateRo): void => {
			certificatesDom.push(<Certificate
				key={certificate.id}
				certificate={certificate}
				secrets={this.state.secrets}
				selected={!!this.state.selected[certificate.id]}
				open={!!this.state.opened[certificate.id]}
				onSelect={(shift: boolean): void => {
					let selected = {
						...this.state.selected,
					};

					if (shift) {
						let certificates = this.state.certificates;
						let start: number;
						let end: number;

						for (let i = 0; i < certificates.length; i++) {
							let usr = certificates[i];

							if (usr.id === certificate.id) {
								start = i;
							} else if (usr.id === this.state.lastSelected) {
								end = i;
							}
						}

						if (start !== undefined && end !== undefined) {
							if (start > end) {
								end = [start, start = end][0];
							}

							for (let i = start; i <= end; i++) {
								selected[certificates[i].id] = true;
							}

							this.setState({
								...this.state,
								lastSelected: certificate.id,
								selected: selected,
							});

							return;
						}
					}

					if (selected[certificate.id]) {
						delete selected[certificate.id];
					} else {
						selected[certificate.id] = true;
					}

					this.setState({
						...this.state,
						lastSelected: certificate.id,
						selected: selected,
					});
				}}
				onOpen={(): void => {
					let opened = {
						...this.state.opened,
					};

					if (opened[certificate.id]) {
						delete opened[certificate.id];
					} else {
						opened[certificate.id] = true;
					}

					this.setState({
						...this.state,
						opened: opened,
					});
				}}
			/>);
		});

		let filterClass = 'bp5-button bp5-intent-primary bp5-icon-filter ';
		if (this.state.filter) {
			filterClass += 'bp5-active';
		}

		let selectedNames: string[] = [];
		for (let instId of Object.keys(this.state.selected)) {
			let inst = CertificatesStore.certificate(instId);
			if (inst) {
				selectedNames.push(inst.name || instId);
			} else {
				selectedNames.push(instId);
			}
		}

		let newCertificateDom: JSX.Element;
		if (this.state.newOpened) {
			newCertificateDom = <CertificateNew
				secrets={this.state.secrets}
				onClose={(): void => {
					this.setState({
						...this.state,
						newOpened: false,
					});
				}}
			/>;
		}

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Certificates</h2>
					<div className="flex"/>
					<div style={css.buttons}>
						<button
							className={filterClass}
							style={css.button}
							type="button"
							onClick={(): void => {
								if (this.state.filter === null) {
									CertificateActions.filter({});
								} else {
									CertificateActions.filter(null);
								}
							}}
						>
							Filters
						</button>
						<button
							className="bp5-button bp5-intent-warning bp5-icon-chevron-up"
							style={css.button}
							disabled={!this.opened}
							type="button"
							onClick={(): void => {
								this.setState({
									...this.state,
									opened: {},
								});
							}}
						>
							Collapse All
						</button>
						<ConfirmButton
							label="Delete Selected"
							className="bp5-intent-danger bp5-icon-delete"
							progressClassName="bp5-intent-danger"
							safe={true}
							style={css.button}
							confirmMsg="Permanently delete the selected certificates"
							confirmInput={true}
							items={selectedNames}
							disabled={!this.selected || this.state.disabled}
							onConfirm={this.onDelete}
						/>
						<button
							className="bp5-button bp5-intent-success bp5-icon-add"
							style={css.button}
							disabled={this.state.disabled || this.state.newOpened}
							type="button"
							onClick={(): void => {
								this.setState({
									...this.state,
									newOpened: true,
								});
							}}
						>New</button>
					</div>
				</div>
			</PageHeader>
			<CertificatesFilter
				filter={this.state.filter}
				onFilter={(filter): void => {
					CertificateActions.filter(filter);
				}}
			/>
			<div style={css.itemsBox}>
				<div style={css.items}>
					{newCertificateDom}
					{certificatesDom}
					<tr className="bp5-card bp5-row" style={css.placeholder}>
						<td colSpan={2} style={css.placeholder}/>
					</tr>
				</div>
			</div>
			<NonState
				hidden={!!certificatesDom.length}
				iconClass="bp5-icon-endorsed"
				title="No certificates"
				description="Add a new certificate to get started."
			/>
			<CertificatesPage
				onPage={(): void => {
					this.setState({
						...this.state,
						selected: {},
						lastSelected: null,
					});
				}}
			/>
		</Page>;
	}
}
