/// <reference path="../References.d.ts"/>
import * as React from 'react';
import UsersStore from '../stores/UsersStore';
import * as UserActions from '../actions/UserActions';

interface State {
	page: number;
	pageCount: number;
	pages: number;
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
			pages: UsersStore.pages,
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
			pages: UsersStore.pages,
			count: UsersStore.count,
		});
	}

	render(): JSX.Element {
		let links: JSX.Element[] = [];
		let page = this.state.page;
		let pages = this.state.pages;
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
					UserActions.traverse(this.state.pages);
				}}
			>
				Last
			</button>
		</div>;
	}
}
