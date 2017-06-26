/// <reference path="../References.d.ts"/>
import * as React from 'react';

const css = {
	header: {
		fontSize: '24px',
		paddingBottom: '8px',
		marginBottom: '20px',
		borderBottomStyle: 'solid',
	} as React.CSSProperties,
};

export default class PageHeader extends React.Component<{}, {}> {
	render(): JSX.Element {
		return <div className="pt-border" style={css.header}>
			<h2>{this.props.children}</h2>
		</div>;
	}
}
