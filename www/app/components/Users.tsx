/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Styles from '../Styles';
import * as UserTypes from '../types/UserTypes';
import UsersStore from '../stores/UsersStore';
import * as UserActions from '../actions/UserActions';
import User from './User';

interface State {
	users: UserTypes.Users;
}

const css = {
	users: {
		width: '100%',
		display: 'table',
		borderSpacing: '0 5px',
		marginTop: '-5px',
	} as React.CSSProperties,
};

export default class Users extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			users: UsersStore.users,
		};
	}

	componentDidMount(): void {
		UserActions.sync();
		UsersStore.addChangeListener(this.onChange);
	}

	componentWillUnmount(): void {
		UsersStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			users: UsersStore.users,
		});
	}

	render(): JSX.Element {
		let usersDom: JSX.Element[] = [];

		for (let user of this.state.users) {
			usersDom.push(<User key={user.id} user={user}/>)
		}

		return <div style={Styles.page}>
			<div className="pt-border" style={Styles.pageHeader}>
				<h2>Users</h2>
			</div>
			<div style={css.users}>
				{usersDom}
			</div>
		</div>;
	}
}
