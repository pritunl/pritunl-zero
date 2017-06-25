/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Styles from '../Styles';
import * as MiscUtils from '../utils/MiscUtils';
import * as UserActions from '../actions/UserActions';
import * as UserTypes from '../types/UserTypes';

interface Props {
	userId: string;
}

interface State {
	changed: boolean;
	disabled: boolean;
	message: string;
	addRole: string;
	user: UserTypes.User;
}

const css = {
	input: {
		width: '100%',
		maxWidth: '310px',
	} as React.CSSProperties,
	button: {
		marginLeft: '10px',
	} as React.CSSProperties,
	role: {
		margin: '9px 5px 0 5px',
		height: '20px',
	} as React.CSSProperties,
};

export default class UserDetailed extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			changed: false,
			disabled: false,
			message: '',
			addRole: '',
			user: null,
		};
	}

	componentDidMount(): void {
		UserActions.get(this.props.userId).then((user: UserTypes.User) => {
			this.setState({
				...this.state,
				user: user,
			});
		});
	}

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		UserActions.commit(this.state.user).then((): void => {
			this.setState({
				...this.state,
				message: 'Your changes have been saved',
				changed: false,
				disabled: false,
			})
		});
	}

	set = (name: string, val: any): void => {
		let user = {
			...this.state.user,
		} as any;

		user[name] = val;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			user: user,
		});
	}

	onAddRole = (): void => {
		let roles = this.state.user.roles.slice(0);

		if (roles.indexOf(this.state.addRole) === -1) {
			roles.push(this.state.addRole);
		}

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			user: {
				...this.state.user,
				roles: roles,
			},
		});
	}

	onRemoveRole = (role: string): void => {
		let roles = this.state.user.roles.slice(0);

		let i = roles.indexOf(role);
		if (i === -1) {
			return;
		}

		roles.splice(i, 1);

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addRole: '',
			user: {
				...this.state.user,
				roles: roles,
			},
		});
	}

	render(): JSX.Element {
		let user = this.state.user;
		if (!user) {
			return <div/>;
		}

		let roles: JSX.Element[] = [];
		for (let role of user.roles) {
			roles.push(
				<div
					className="pt-tag pt-tag-removable pt-intent-primary"
					style={css.role}
				>
					{role}
					<button
						className="pt-tag-remove"
						onClick={(): void => {
							this.onRemoveRole(role);
						}}
					/>
				</div>
			);
		}

		return <div style={Styles.page}>
			<div className="pt-border" style={Styles.pageHeader}>
				<h2>User Info</h2>
			</div>
			<div className="layout horizontal">
				<div className="flex">
					<label className="pt-label">
						Username
						<input
							className="pt-input"
							style={css.input}
							type="text"
							autoCapitalize="off"
							spellCheck={false}
							placeholder="Enter Elasticsearch address"
							value={user.username}
							onChange={(evt): void => {
								this.set('username', evt.target.value);
							}}
						/>
					</label>
				</div>
				<div className="flex">
					<label className="pt-label">
						Roles
						<div>
							{roles}
						</div>
					</label>
					<div className="pt-control-group">
						<input
							className="pt-input"
							type="text"
							autoCapitalize="off"
							spellCheck={false}
							placeholder="Add role"
							value={this.state.addRole}
							onChange={(evt): void => {
								this.setState({
									...this.state,
									addRole: evt.target.value,
								});
							}}
							onKeyPress={(evt): void => {
								if (evt.key === 'Enter') {
									this.onAddRole();
								}
							}}
						/>
						<button
							className="pt-button"
							onClick={this.onAddRole}
						>Add</button>
					</div>
				</div>
			</div>
			<div className="layout horizontal">
				<div className="flex"/>
				<div>
					<span hidden={!this.state.message}>
						{this.state.message}
					</span>
					<button
						className="pt-button pt-intent-success pt-icon-tick"
						style={css.button}
						type="button"
						disabled={!this.state.changed || this.state.disabled}
						onClick={this.onSave}
					>
						Save
					</button>
				</div>
			</div>
		</div>;
	}
}
