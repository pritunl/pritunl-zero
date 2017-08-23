/// <reference path="../References.d.ts"/>
import * as React from 'react';
import UsersStore from '../stores/UsersStore';
import * as UserActions from '../actions/UserActions';

interface Props {
	onPage?: () => void;
}

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
		userSelect: 'none',
		margin: '5px 5px 0 0',
	} as React.CSSProperties,
	current: {
		opacity: 0.5,
	} as React.CSSProperties,
};

export default class Users extends React.Component<Props, State> {
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
		let page = this.state.page;
		let pages = this.state.pages;

		if (pages <= 1) {
			return <div/>;
		}

		let offset = 1;
		if (pages < 5) {
			offset = 0;
		}

		let links: JSX.Element[] = [];
		let start = Math.max(offset, page - 7);
		let end = Math.min(pages - offset, start + 15);

		for (let i = start; i < end; i++) {
			links.push(<a
				key={i}
				style={page === i ? {
					...css.link,
					...css.current,
				} : css.link}
				onClick={(): void => {
					UserActions.traverse(i);
					if (this.props.onPage) {
						this.props.onPage();
					}
				}}
			>
				{i + 1}
			</a>);
		}

		return <div className="layout horizontal center-justified">
			<button
				className="pt-button"
				hidden={!offset}
				style={page === 0 ? {
					...css.button,
					...css.current,
				} : css.button}
				type="button"
				onClick={(): void => {
					UserActions.traverse(0);
					if (this.props.onPage) {
						this.props.onPage();
					}
				}}
			>
				First
			</button>
			{links}
			<button
				className="pt-button"
				hidden={!offset}
				style={page === pages ? {
					...css.buttonLast,
					...css.current,
				} : css.buttonLast}
				type="button"
				onClick={(): void => {
					UserActions.traverse(this.state.pages);
					if (this.props.onPage) {
						this.props.onPage();
					}
				}}
			>
				Last
			</button>
		</div>;
	}
}
