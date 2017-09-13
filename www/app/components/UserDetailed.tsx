/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as UserActions from '../actions/UserActions';
import * as UserTypes from '../types/UserTypes';
import * as MiscUtils from '../utils/MiscUtils';
import UserStore from '../stores/UserStore';
import Sessions from './Sessions';
import Audits from './Audits';
import Page from './Page';
import PageHeader from './PageHeader';
import PagePanel from './PagePanel';
import PageSplit from './PageSplit';
import PageInfo from './PageInfo';
import PageInput from './PageInput';
import PageInputButton from './PageInputButton';
import PageSwitch from './PageSwitch';
import PageSelect from './PageSelect';
import PageDateTime from './PageDateTime';
import PageSave from './PageSave';
import PageNew from './PageNew';
import ConfirmButton from './ConfirmButton';
import Help from './Help';

interface Props {
	userId?: string;
}

interface State {
	changed: boolean;
	disabled: boolean;
	locked: boolean;
	message: string;
	addRole: string;
	user: UserTypes.User;
}

const css = {
	header: {
		marginTop: '-19px',
	} as React.CSSProperties,
	heading: {
		margin: '19px 0 0 0',
	} as React.CSSProperties,
	button: {
		margin: '10px 0 0 0',
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
			locked: false,
			message: '',
			addRole: '',
			user: UserStore.userM,
		};
	}

	componentDidMount(): void {
		UserStore.addChangeListener(this.onChange);
		UserActions.load(this.props.userId);
	}

	componentWillUnmount(): void {
		UserStore.removeChangeListener(this.onChange);
		UserActions.unload();
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			user: UserStore.userM,
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
			});
		}).catch((): void => {
			this.setState({
				...this.state,
				message: '',
				disabled: false,
			});
		});
	}

	onNew = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		UserActions.create(this.state.user).then((): void => {
			this.setState({
				...this.state,
				message: 'User has been created',
				changed: false,
				disabled: false,
				locked: true,
			});
		}).catch((): void => {
			this.setState({
				...this.state,
				message: '',
				disabled: false,
			});
		});
	}

	set = (name: string, val: any): void => {
		let user: any = {
			...this.state.user,
		};

		user[name] = val;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			user: user,
		});
	}

	onAddRole = (): void => {
		let roles = [
			...this.state.user.roles,
		];

		if (!this.state.addRole) {
			return;
		}

		if (roles.indexOf(this.state.addRole) === -1) {
			roles.push(this.state.addRole);
		}

		roles.sort();

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
		let roles = [
			...this.state.user.roles,
		];

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

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		UserActions.remove([this.props.userId]).then((): void => {
			this.setState({
				...this.state,
				message: 'User has been deleted',
				changed: false,
				disabled: false,
				locked: true,
			});
		}).catch((): void => {
			this.setState({
				...this.state,
				disabled: false,
			});
		});
	}

	render(): JSX.Element {
		let userId = this.props.userId;
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
					key={role}
				>
					{role}
					<button
						className="pt-tag-remove"
						disabled={this.state.locked}
						onMouseUp={(): void => {
							this.onRemoveRole(role);
						}}
					/>
				</div>,
			);
		}

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>{userId ? 'User Info' : 'New User'}</h2>
					<div className="flex"/>
					<div>
						<ConfirmButton
							label="Delete"
							className="pt-intent-danger pt-icon-delete"
							progressClassName="pt-intent-danger"
							style={css.button}
							disabled={this.state.disabled || this.state.locked}
							hidden={!userId}
							onConfirm={this.onDelete}
						/>
					</div>
				</div>
			</PageHeader>
			<PageSplit>
				<PagePanel className="layout vertical">
					<PageInput
						disabled={this.state.locked}
						label="Username"
						help="Username, if using single sign-on username must match"
						type="text"
						placeholder="Enter username"
						value={user.username}
						onChange={(val): void => {
							this.set('username', val);
						}}
					/>
					<PageInput
						hidden={user.type !== 'local'}
						disabled={this.state.locked}
						label="Password"
						help="Password, leave blank to keep current password"
						type="password"
						placeholder="Change password"
						value={user.password}
						onChange={(val): void => {
							this.set('password', val);
						}}
					/>
					<PageSelect
						disabled={this.state.locked}
						label="Type"
						help="A local user is a user that is created on the Pritunl Zero database that has a username and password. The other user types can be used to create users for single sign-on services. Generally single sign-on users will be created automatically when the user authenticates for the first time. It can sometimes be desired to manaully create a single sign-on user to provide roles in advanced of the first login."
						value={user.type}
						onChange={(val): void => {
							this.set('type', val);
						}}
					>
						<option value="local">Local</option>
						<option value="azure">Azure</option>
						<option value="google">Google</option>
						<option value="onelogin">OneLogin</option>
						<option value="okta">Okta</option>
					</PageSelect>
					<label className="pt-label">
						Roles
						<Help
							title="Roles"
							content="User roles will be used to match with service roles. A user must have a matching role to access a service."
						/>
						<div>
							{roles}
						</div>
					</label>
					<PageInputButton
						disabled={this.state.locked}
						buttonClass="pt-intent-success pt-icon-add"
						label="Add"
						type="text"
						placeholder="Add role"
						value={this.state.addRole}
						onChange={(val): void => {
							this.setState({
								...this.state,
								addRole: val,
							});
						}}
						onSubmit={this.onAddRole}
					/>
					<PageSwitch
						label="Administrator"
						help="Enable to give user administrator access to the management console"
						disabled={this.state.locked}
						checked={user.administrator === 'super'}
						onToggle={(): void => {
							if (this.state.user.administrator === 'super') {
								this.set('administrator', '');
							} else {
								this.set('administrator', 'super');
							}
						}}
					/>
					<PageSwitch
						label="Disabled"
						help="Disables the user ending all active sessions and prevents new authentications"
						disabled={this.state.locked}
						checked={user.disabled}
						onToggle={(): void => {
							this.set('disabled', !this.state.user.disabled);
						}}
					/>
				</PagePanel>
				<PagePanel>
					<PageInfo
						fields={[
							{
								label: 'ID',
								value: user.id || 'None',
							},
							{
								label: 'Last Active',
								value: MiscUtils.formatDate(user.last_active) || 'Inactive',
							},
						]}
					/>
					<PageDateTime
						label="Active Until"
						help="Set this to schedule the user to be disabled at the set date and time. This is useful to give a user temporary access to a service."
						value={user.active_until}
						disabled={user.disabled || this.state.locked}
						onChange={(val): void => {
							this.set('active_until', val);
						}}
					/>
				</PagePanel>
			</PageSplit>
			{userId ? <PageSave
				message={this.state.message}
				changed={this.state.changed}
				disabled={this.state.disabled || this.state.locked}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						message: 'Your changes have been discarded',
						addRole: '',
						user: UserStore.userM,
					});
				}}
				onSave={this.onSave}
			/> : <PageNew
				message={this.state.message}
				changed={this.state.changed}
				disabled={this.state.disabled || this.state.locked}
				onSave={this.onNew}
			/>}
			{this.state.locked ? null : <Sessions userId={userId}/>}
			{this.state.locked ? null : <Audits userId={userId}/>}
		</Page>;
	}
}
