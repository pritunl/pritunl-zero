/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as UserTypes from '../types/UserTypes';
import UsersStore from '../stores/UsersStore';
import * as UserActions from '../actions/UserActions';
import User from './User';
import Page from './Page';
import PageHeader from './PageHeader';

interface State {
	users: UserTypes.Users;
}

const css = {
	users: {
		width: '100%',
		marginTop: '-5px',
		display: 'table',
		borderSpacing: '0 5px',
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

		return <Page>
			<PageHeader>
				Users
			</PageHeader>
			<div style={css.users}>
				{usersDom}
			</div>
		</Page>;
	}
}
