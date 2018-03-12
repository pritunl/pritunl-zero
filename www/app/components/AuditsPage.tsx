/// <reference path="../References.d.ts"/>
import * as React from 'react';
import AuditsStore from '../stores/AuditsStore';
import * as AuditActions from '../actions/AuditActions';

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
		userSelect: 'none',
		margin: '0 5px 0 0',
	} as React.CSSProperties,
	buttonLast: {
		userSelect: 'none',
		margin: '0 0 0 0',
	} as React.CSSProperties,
	link: {
		cursor: 'pointer',
		userSelect: 'none',
		margin: '5px 5px 0 0',
	} as React.CSSProperties,
	current: {
		opacity: 0.5,
	} as React.CSSProperties,
};

export default class AuditsPage extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			page: AuditsStore.page,
			pageCount: AuditsStore.pageCount,
			pages: AuditsStore.pages,
			count: AuditsStore.count,
		};
	}

	componentDidMount(): void {
		AuditsStore.addChangeListener(this.onChange);
	}

	componentWillUnmount(): void {
		AuditsStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			page: AuditsStore.page,
			pageCount: AuditsStore.pageCount,
			pages: AuditsStore.pages,
			count: AuditsStore.count,
		});
	}

	render(): JSX.Element {
		let page = this.state.page;
		let pages = this.state.pages;

		if (pages <= 1) {
			return <div/>;
		}

		let links: JSX.Element[] = [];
		let start = Math.max(0, page - 7);
		let end = Math.min(pages, start + 15);

		for (let i = start; i < end; i++) {
			links.push(<span
				key={i}
				style={page === i ? {
					...css.link,
					...css.current,
				} : css.link}
				onClick={(): void => {
					AuditActions.traverse(i);
					if (this.props.onPage) {
						this.props.onPage();
					}
				}}
			>
				{i + 1}
			</span>);
		}

		return <div className="layout horizontal center-justified">
			<button
				className="pt-button pt-minimal pt-icon-chevron-backward"
				hidden={pages < 5}
				disabled={page === 0}
				type="button"
				onClick={(): void => {
					AuditActions.traverse(0);
					if (this.props.onPage) {
						this.props.onPage();
					}
				}}
			/>
			<button
				className="pt-button pt-minimal pt-icon-chevron-left"
				style={css.button}
				disabled={page === 0}
				type="button"
				onClick={(): void => {
					AuditActions.traverse(Math.max(0, this.state.page - 1));
					if (this.props.onPage) {
						this.props.onPage();
					}
				}}
			/>
			{links}
			<button
				className="pt-button pt-minimal pt-icon-chevron-right"
				style={css.button}
				disabled={page === pages - 1}
				type="button"
				onClick={(): void => {
					AuditActions.traverse(Math.min(
						this.state.pages - 1, this.state.page + 1));
					if (this.props.onPage) {
						this.props.onPage();
					}
				}}
			/>
			<button
				className="pt-button pt-minimal pt-icon-chevron-forward"
				hidden={pages < 5}
				disabled={page === pages - 1}
				type="button"
				onClick={(): void => {
					AuditActions.traverse(this.state.pages - 1);
					if (this.props.onPage) {
						this.props.onPage();
					}
				}}
			/>
		</div>;
	}
}
