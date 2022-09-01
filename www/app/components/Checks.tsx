/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as CheckTypes from '../types/CheckTypes';
import * as AuthorityTypes from '../types/AuthorityTypes';
import ChecksStore from '../stores/ChecksStore';
import AuthoritiesStore from '../stores/AuthoritiesStore';
import * as CheckActions from '../actions/CheckActions';
import * as AuthorityActions from '../actions/AuthorityActions';
import Check from './Check';
import ChecksFilter from './ChecksFilter';
import ChecksPage from './ChecksPage';
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
	checks: CheckTypes.ChecksRo;
	filter: CheckTypes.Filter;
	authorities: AuthorityTypes.AuthoritiesRo;
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

export default class Checks extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			checks: ChecksStore.checks,
			filter: ChecksStore.filter,
			authorities: AuthoritiesStore.authorities,
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
		ChecksStore.addChangeListener(this.onChange);
		AuthoritiesStore.addChangeListener(this.onChange);
		CheckActions.sync();
		AuthorityActions.sync();
	}

	componentWillUnmount(): void {
		ChecksStore.removeChangeListener(this.onChange);
		AuthoritiesStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		let checks = ChecksStore.checks;
		let selected: Selected = {};
		let curSelected = this.state.selected;
		let opened: Opened = {};
		let curOpened = this.state.opened;

		checks.forEach((check: CheckTypes.Check): void => {
			if (curSelected[check.id]) {
				selected[check.id] = true;
			}
			if (curOpened[check.id]) {
				opened[check.id] = true;
			}
		});

		this.setState({
			...this.state,
			checks: checks,
			filter: ChecksStore.filter,
			authorities: AuthoritiesStore.authorities,
			selected: selected,
			opened: opened,
		});
	}

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		CheckActions.removeMulti(
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
		let checksDom: JSX.Element[] = [];

		this.state.checks.forEach((
				check: CheckTypes.CheckRo): void => {
			checksDom.push(<Check
				key={check.id}
				check={check}
				authorities={this.state.authorities}
				selected={!!this.state.selected[check.id]}
				open={!!this.state.opened[check.id]}
				onSelect={(shift: boolean): void => {
					let selected = {
						...this.state.selected,
					};

					if (shift) {
						let checks = this.state.checks;
						let start: number;
						let end: number;

						for (let i = 0; i < checks.length; i++) {
							let usr = checks[i];

							if (usr.id === check.id) {
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
								selected[checks[i].id] = true;
							}

							this.setState({
								...this.state,
								lastSelected: check.id,
								selected: selected,
							});

							return;
						}
					}

					if (selected[check.id]) {
						delete selected[check.id];
					} else {
						selected[check.id] = true;
					}

					this.setState({
						...this.state,
						lastSelected: check.id,
						selected: selected,
					});
				}}
				onOpen={(): void => {
					let opened = {
						...this.state.opened,
					};

					if (opened[check.id]) {
						delete opened[check.id];
					} else {
						opened[check.id] = true;
					}

					this.setState({
						...this.state,
						opened: opened,
					});
				}}
			/>);
		});

		let filterClass = 'bp3-button bp3-intent-primary bp3-icon-filter ';
		if (this.state.filter) {
			filterClass += 'bp3-active';
		}

		let selectedNames: string[] = [];
		for (let instId of Object.keys(this.state.selected)) {
			let inst = ChecksStore.check(instId);
			if (inst) {
				selectedNames.push(inst.name || instId);
			} else {
				selectedNames.push(instId);
			}
		}

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Health Checks</h2>
					<div className="flex"/>
					<div style={css.buttons}>
						<button
							className={filterClass}
							style={css.button}
							type="button"
							onClick={(): void => {
								if (this.state.filter === null) {
									CheckActions.filter({});
								} else {
									CheckActions.filter(null);
								}
							}}
						>
							Filters
						</button>
						<button
							className="bp3-button bp3-intent-warning bp3-icon-chevron-up"
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
							className="bp3-intent-danger bp3-icon-delete"
							progressClassName="bp3-intent-danger"
							safe={true}
							style={css.button}
							confirmMsg="Permanently delete the selected checks"
							confirmInput={true}
							items={selectedNames}
							disabled={!this.selected || this.state.disabled}
							onConfirm={this.onDelete}
						/>
						<button
							className="bp3-button bp3-intent-success bp3-icon-add"
							style={css.button}
							disabled={this.state.disabled}
							type="button"
							onClick={(): void => {
								this.setState({
									...this.state,
									disabled: true,
								});
								CheckActions.create({
									id: null,
								}).then((): void => {
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
							}}
						>New</button>
					</div>
				</div>
			</PageHeader>
			<ChecksFilter
				filter={this.state.filter}
				onFilter={(filter): void => {
					CheckActions.filter(filter);
				}}
				authorities={this.state.authorities}
			/>
			<div style={css.itemsBox}>
				<div style={css.items}>
					{checksDom}
					<tr className="bp3-card bp3-row" style={css.placeholder}>
						<td colSpan={5} style={css.placeholder}/>
					</tr>
				</div>
			</div>
			<NonState
				hidden={!!checksDom.length}
				iconClass="bp3-icon-lifesaver"
				title="No checks"
				description="Add a new check to get started."
			/>
			<ChecksPage
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
