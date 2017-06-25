/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	title: string;
}

const css = {
	header: {
		fontSize: '24px',
		paddingBottom: '8px',
		marginBottom: '20px',
		borderBottomStyle: 'solid',
	} as React.CSSProperties,
};

export default class PageHeader extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <div className="pt-border" style={css.header}>
			<h2>{this.props.title}</h2>
		</div>;
	}
}
