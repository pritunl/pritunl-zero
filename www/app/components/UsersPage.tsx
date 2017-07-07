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
	link: {
		margin: '5px 5px 0 0',
	} as React.CSSProperties,
	current: {
		opacity: 0.5,
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
		let links: JSX.Element[] = [];
		let page = this.state.page;
		let pages = Math.ceil(this.state.count / this.state.pageCount);
		let start = Math.max(1, page - 7);
		let end = Math.min(pages - 1, start + 15);

		for (let i = start; i < end; i++) {
			links.push(<a
				key={i}
				style={page === i ? {
					...css.link,
					...css.current,
				} : css.link}
				onClick={(): void => {
					UserActions.traverse(i);
				}}
			>
				{i + 1}
			</a>);
		}

		return <div className="layout horizontal center-justified">
			<button
				className="pt-button"
				style={page === 0 ? {
					...css.button,
					...css.current,
				} : css.button}
				type="button"
				onClick={(): void => {
					UserActions.traverse(0);
				}}
			>
				First
			</button>
			{links}
			<button
				className="pt-button"
				style={page === pages ? {
					...css.buttonLast,
					...css.current,
				} : css.buttonLast}
				type="button"
				onClick={(): void => {
					UserActions.traverse(
						Math.ceil(this.state.count / this.state.pageCount));
				}}
			>
				Last
			</button>
		</div>;
	}
}
