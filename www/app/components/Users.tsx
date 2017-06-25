/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Styles from '../Styles';
import * as UserTypes from '../types/UserTypes';
import UserStore from '../stores/UserStore';
import * as UserActions from '../actions/UserActions';
import User from './User';

interface State {
	users: UserTypes.Users;
}

export default class Users extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			users: UserStore.users,
		};
	}

	componentDidMount(): void {
		UserActions.sync();
		UserStore.addChangeListener(this.onChange);
	}

	componentWillUnmount(): void {
		UserStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			users: UserStore.users,
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
			<div className="layout horizontal">
				{usersDom}
			</div>
		</div>;
	}
}
