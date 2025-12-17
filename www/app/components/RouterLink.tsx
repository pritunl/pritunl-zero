/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Router from '../Router';

interface Props {
	className?: string;
	style?: React.CSSProperties;
	hidden?: boolean;
	to: string;
	children?: React.ReactNode
}

export default class RouterLink extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <a
			className={this.props.className}
			style={this.props.style}
			hidden={this.props.hidden}
			href={"#" + this.props.to}
			onClick={(): void => {
				Router.setLocation(this.props.to);
			}}
		>{this.props.children}</a>
	}
}
