/// <reference path="../References.d.ts"/>
import * as React from 'react';
import RouterLink from './RouterLink';
import * as UserTypes from '../types/UserTypes';
import UsersStore from '../stores/UsersStore';
import * as UserActions from '../actions/UserActions';
import * as AuditActions from '../actions/AuditActions';
import User from './User';
import UsersFilter from './UsersFilter';
import Page from './Page';
import PageHeader from './PageHeader';
import UsersPage from './UsersPage';
import ConfirmButton from './ConfirmButton';

interface Selected {
	[key: string]: boolean;
}

interface State {
	users: UserTypes.UsersRo;
	filter: UserTypes.Filter;
	selected: Selected;
	lastSelected: string;
	disabled: boolean;
}

const css = {
	users: {
		width: '100%',
		marginTop: '-5px',
		display: 'table',
		tableLayout: 'fixed',
		borderSpacing: '0 5px',
	} as React.CSSProperties,
	usersBox: {
		width: '100%',
		overflowY: 'auto',
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

export default class Users extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			users: UsersStore.users,
			filter: UsersStore.filter,
			selected: {},
			lastSelected: null,
			disabled: false,
		};
	}

	get selected(): boolean {
		for (let key in this.state.selected) {
			if (this.state.selected[key]) {
				return true;
			}
		}
		return false;
	}

	componentDidMount(): void {
		UsersStore.addChangeListener(this.onChange);
		AuditActions.traverse(0);
		UserActions.sync();
	}

	componentWillUnmount(): void {
		UsersStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		let users = UsersStore.users;
		let selected: Selected = {};
		let curSelected = this.state.selected;

		users.forEach((user: UserTypes.User): void => {
			if (curSelected[user.id]) {
				selected[user.id] = true;
			}
		});

		this.setState({
			...this.state,
			users: users,
			filter: UsersStore.filter,
			selected: selected,
		});
	}

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		UserActions.remove(Object.keys(this.state.selected)).then((): void => {
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
		let usersDom: JSX.Element[] = [];

		this.state.users.forEach((user: UserTypes.UserRo): void => {
			usersDom.push(<User
				key={user.id}
				user={user}
				selected={!!this.state.selected[user.id]}
				onSelect={(shift: boolean): void => {
					let selected = {
						...this.state.selected,
					};

					if (shift) {
						let users = this.state.users;
						let start: number;
						let end: number;

						for (let i = 0; i < users.length; i++) {
							let usr = users[i];

							if (usr.id === user.id) {
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
								selected[users[i].id] = true;
							}

							this.setState({
								...this.state,
								lastSelected: user.id,
								selected: selected,
							});

							return;
						}
					}

					if (selected[user.id]) {
						delete selected[user.id];
					} else {
						selected[user.id] = true;
					}

					this.setState({
						...this.state,
						lastSelected: user.id,
						selected: selected,
					});
				}}
			/>);
		});

		let filterClass = 'bp5-button bp5-intent-primary bp5-icon-filter ';
		if (this.state.filter) {
			filterClass += 'bp5-active';
		}

		let selectedNames: string[] = [];
		for (let userId of Object.keys(this.state.selected)) {
			let user = UsersStore.user(userId);
			if (user) {
				selectedNames.push(user.username || userId);
			} else {
				selectedNames.push(userId);
			}
		}

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Users</h2>
					<div className="flex"/>
					<div style={css.buttons}>
						<button
							className={filterClass}
							style={css.button}
							type="button"
							onClick={(): void => {
								if (this.state.filter === null) {
									UserActions.filter({});
								} else {
									UserActions.filter(null);
								}
							}}
						>
							Filters
						</button>
						<ConfirmButton
							label="Delete Selected"
							className="bp5-intent-danger bp5-icon-delete"
							progressClassName="bp5-intent-danger"
							safe={true}
							style={css.button}
							confirmMsg="Permanently delete the selected users"
							confirmInput={true}
							items={selectedNames}
							disabled={!this.selected || this.state.disabled}
							onConfirm={this.onDelete}
						/>
						<RouterLink
							className="bp5-button bp5-intent-success bp5-icon-add"
							style={css.button}
							to="/user"
						>
							New
						</RouterLink>
					</div>
				</div>
			</PageHeader>
			<UsersFilter
				filter={this.state.filter}
				onFilter={(filter): void => {
					UserActions.filter(filter);
				}}
			/>
			<div style={css.usersBox}>
				<div style={css.users}>
					{usersDom}
				</div>
			</div>
			<UsersPage
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
