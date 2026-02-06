/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	children?: React.ReactNode
	wide?: boolean;
	full?: boolean;
}

const css = {
	page: {
		margin: '0 auto',
		padding: '30px 20px',
		minWidth: '200px',
		maxWidth: '1100px',
	} as React.CSSProperties,
	pageWide: {
		margin: '0 auto',
		padding: '30px 20px',
		minWidth: '200px',
		maxWidth: '1400px',
	} as React.CSSProperties,
	pageFull: {
		height: '100%',
		display: 'flex',
		flexDirection: 'column',
		overflow: 'hidden',
		margin: '0 auto',
		padding: '30px 20px 0 20px',
		minWidth: '200px',
		maxWidth: '1600px',
	} as React.CSSProperties,
};

export default class Page extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <div style={this.props.full ? css.pageFull :
				(this.props.wide ? css.pageWide : css.page)}>
			{this.props.children}
		</div>;
	}
}
