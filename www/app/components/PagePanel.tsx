/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	children?: React.ReactNode
	className?: string;
	hidden?: boolean;
}

const css = {
	panel: {
		flex: 1,
		minWidth: '250px',
		padding: '0 10px',
	} as React.CSSProperties,
};

export default class PagePanel extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <div
			className={this.props.className}
			style={css.panel}
			hidden={this.props.hidden}
		>
			{this.props.children}
		</div>;
	}
}
