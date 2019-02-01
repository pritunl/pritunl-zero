/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	label?: string;
}

const css = {
	header: {
		fontSize: '24px',
		paddingBottom: '8px',
		marginBottom: '20px',
		borderBottomStyle: 'solid',
	} as React.CSSProperties,
	label: {
		margin: 0,
	} as React.CSSProperties,
};

export default class PageHeader extends React.Component<Props, {}> {
	render(): JSX.Element {
		let label: JSX.Element;
		if (this.props.label) {
			label = <h2 style={css.label}>{this.props.label}</h2>;
		}

		return <div className="bp3-border" style={css.header}>
			{label}
			{this.props.children}
		</div>;
	}
}
