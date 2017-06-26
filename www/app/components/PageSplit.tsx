/// <reference path="../References.d.ts"/>
import * as React from 'react';

export default class PageSplit extends React.Component<{}, {}> {
	render(): JSX.Element {
		return <div className="layout horizontal wrap">
			{this.props.children}
		</div>;
	}
}
