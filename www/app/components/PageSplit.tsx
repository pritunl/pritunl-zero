/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	children?: React.ReactNode
}

export default class PageSplit extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <div className="layout horizontal wrap">
			{this.props.children}
		</div>;
	}
}
