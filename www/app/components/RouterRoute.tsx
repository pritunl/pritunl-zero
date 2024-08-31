/// <reference path="../References.d.ts"/>
import * as React from 'react';

import * as RouterTypes from '../types/RouterTypes';

interface Props {
	path: string;
	render: (data: RouterTypes.State) => JSX.Element;
}

export default class RouterRoute extends React.Component<Props, {}> {
	render(): JSX.Element {
		let data = RouterTypes.getState()

		return this.props.render(data)
	}
}
