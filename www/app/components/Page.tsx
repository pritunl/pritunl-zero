/// <reference path="../References.d.ts"/>
import * as React from 'react';

const css = {
	page: {
		margin: '0 auto',
		padding: '30px 20px',
		minWidth: '200px',
		maxWidth: '1100px',
	} as React.CSSProperties,
};

export default class Page extends React.Component<{}, {}> {
	render(): JSX.Element {
		return <div style={css.page}>
			{this.props.children}
		</div>;
	}
}
