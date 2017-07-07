/// <reference path="../References.d.ts"/>
import * as React from 'react';
import UsersStore from '../stores/UsersStore';
import * as UserActions from '../actions/UserActions';

interface State {
	page: number;
	pageCount: number;
	count: number;
}

const css = {
	button: {
		margin: '0 5px 0 0',
	} as React.CSSProperties,
	buttonLast: {
		margin: '0 0 0 0',
	} as React.CSSProperties,
};

export default class Users extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			page: UsersStore.page,
			pageCount: UsersStore.pageCount,
			count: UsersStore.count,
		};
	}

	componentDidMount(): void {
		UsersStore.addChangeListener(this.onChange);
	}

	componentWillUnmount(): void {
		UsersStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			page: UsersStore.page,
			pageCount: UsersStore.pageCount,
			count: UsersStore.count,
		});
	}

	render(): JSX.Element {
		return <div className="layout horizontal">
			<button
				className="pt-button"
				style={css.button}
				type="button"
				onClick={(): void => {
					UserActions.traverse(0);
				}}
			>
				First
			</button>

			<button
				className="pt-button"
				style={css.buttonLast}
				type="button"
				onClick={(): void => {
					UserActions.traverse(this.state.pageCount);
				}}
			>
				Last
			</button>
		</div>;
	}
}
