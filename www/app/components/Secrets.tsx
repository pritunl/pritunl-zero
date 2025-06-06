/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SecretTypes from '../types/SecretTypes';
import SecretsStore from '../stores/SecretsStore';
import * as SecretActions from '../actions/SecretActions';
import Secret from './Secret';
import SecretNew from './SecretNew';
import SecretsFilter from './SecretsFilter';
import SecretsPage from './SecretsPage';
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
	secrets: SecretTypes.SecretsRo;
	filter: SecretTypes.Filter;
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

export default class Secrets extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			secrets: SecretsStore.secrets,
			filter: SecretsStore.filter,
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
		SecretsStore.addChangeListener(this.onChange);
		SecretActions.sync();
	}

	componentWillUnmount(): void {
		SecretsStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		let secrets = SecretsStore.secrets;
		let selected: Selected = {};
		let curSelected = this.state.selected;
		let opened: Opened = {};
		let curOpened = this.state.opened;

		secrets.forEach((secret: SecretTypes.Secret): void => {
			if (curSelected[secret.id]) {
				selected[secret.id] = true;
			}
			if (curOpened[secret.id]) {
				opened[secret.id] = true;
			}
		});

		this.setState({
			...this.state,
			secrets: secrets,
			filter: SecretsStore.filter,
			selected: selected,
			opened: opened,
		});
	}

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		SecretActions.removeMulti(
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
		let secretsDom: JSX.Element[] = [];

		this.state.secrets.forEach((
				secret: SecretTypes.SecretRo): void => {
			secretsDom.push(<Secret
				key={secret.id}
				secret={secret}
				selected={!!this.state.selected[secret.id]}
				open={!!this.state.opened[secret.id]}
				onSelect={(shift: boolean): void => {
					let selected = {
						...this.state.selected,
					};

					if (shift) {
						let secrets = this.state.secrets;
						let start: number;
						let end: number;

						for (let i = 0; i < secrets.length; i++) {
							let usr = secrets[i];

							if (usr.id === secret.id) {
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
								selected[secrets[i].id] = true;
							}

							this.setState({
								...this.state,
								lastSelected: secret.id,
								selected: selected,
							});

							return;
						}
					}

					if (selected[secret.id]) {
						delete selected[secret.id];
					} else {
						selected[secret.id] = true;
					}

					this.setState({
						...this.state,
						lastSelected: secret.id,
						selected: selected,
					});
				}}
				onOpen={(): void => {
					let opened = {
						...this.state.opened,
					};

					if (opened[secret.id]) {
						delete opened[secret.id];
					} else {
						opened[secret.id] = true;
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
			let inst = SecretsStore.secret(instId);
			if (inst) {
				selectedNames.push(inst.name || instId);
			} else {
				selectedNames.push(instId);
			}
		}

		let newSecretDom: JSX.Element;
		if (this.state.newOpened) {
			newSecretDom = <SecretNew
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
					<h2 style={css.heading}>Secrets</h2>
					<div className="flex"/>
					<div style={css.buttons}>
						<button
							className={filterClass}
							style={css.button}
							type="button"
							onClick={(): void => {
								if (this.state.filter === null) {
									SecretActions.filter({});
								} else {
									SecretActions.filter(null);
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
							confirmMsg="Permanently delete the selected secrets"
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
			<SecretsFilter
				filter={this.state.filter}
				onFilter={(filter): void => {
					SecretActions.filter(filter);
				}}
			/>
			<div style={css.itemsBox}>
				<div style={css.items}>
					{newSecretDom}
					{secretsDom}
					<tr className="bp5-card bp5-row" style={css.placeholder}>
						<td colSpan={2} style={css.placeholder}/>
					</tr>
				</div>
			</div>
			<NonState
				hidden={!!secretsDom.length}
				iconClass="bp5-icon-ip-address"
				title="No secrets"
				description="Add a new secret to get started."
			/>
			<SecretsPage
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
