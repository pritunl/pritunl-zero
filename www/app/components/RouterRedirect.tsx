/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	to: string;
}

export default class RouterRedirect extends React.Component<Props, {}> {
	render(): JSX.Element {
		window.location.hash = this.props.to
		let evt = new Event("router_update")
		window.dispatchEvent(evt)
		return <div></div>
	}
}
