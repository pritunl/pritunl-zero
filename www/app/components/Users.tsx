/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ReactRouter from 'react-router-dom';
import * as UserTypes from '../types/UserTypes';
import UsersStore from '../stores/UsersStore';
import * as UserActions from '../actions/UserActions';
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
		borderSpacing: '0 5px',
	} as React.CSSProperties,
	header: {
		marginTop: '-19px',
	} as React.CSSProperties,
	heading: {
		margin: '19px 0 0 0',
	} as React.CSSProperties,
	button: {
		margin: '10px 0 0 10px',
	} as React.CSSProperties,
	buttonFirst: {
		margin: '10px 0 0 0',
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
		for (let val in this.state.selected) {
			if (this.state.selected[val]) {
				return true;
			}
		}
		return false;
	}

	componentDidMount(): void {
		UsersStore.addChangeListener(this.onChange);
		UserActions.sync();
	}

	componentWillUnmount(): void {
		UsersStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		let users = UsersStore.users;
		let selected: Selected = {};
		let curSelected = this.state.selected;

		this.state.users.forEach((user: UserTypes.User): void => {
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

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Users</h2>
					<div className="flex"/>
					<div>
						<button
							className="pt-button pt-intent-primary pt-icon-filter"
							style={css.buttonFirst}
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
							className="pt-intent-danger pt-icon-delete"
							progressClassName="pt-intent-danger"
							style={css.button}
							disabled={!this.selected || this.state.disabled}
							onConfirm={this.onDelete}
						/>
						<ReactRouter.Link
							className="pt-button pt-intent-success pt-icon-add"
							style={css.button}
							to="/user"
						>
							New
						</ReactRouter.Link>
					</div>
				</div>
			</PageHeader>
			<UsersFilter
				filter={this.state.filter}
				onFilter={(filter): void => {
					UserActions.filter(filter);
				}}
			/>
			<div style={css.users}>
				{usersDom}
			</div>
			<UsersPage
				onPage={(): void => {
					this.setState({
						lastSelected: null,
					});
				}}
			/>
		</Page>;
	}
}
