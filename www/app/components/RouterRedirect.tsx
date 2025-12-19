/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Router from '../Router';

interface Props {
	to: string;
}

export default class RouterRedirect extends React.Component<Props, {}> {
	render(): JSX.Element {
		Router.setLocation(this.props.to);
		return <div></div>
	}
}
