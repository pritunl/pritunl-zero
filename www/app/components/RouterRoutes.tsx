/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as RouterTypes from '../types/RouterTypes';

interface Props {
	children?: React.ReactNode
}

interface State {
	path: string
}

export default class RouterRoutes extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			path: "",
		}
	}

	componentDidMount() {
		window.addEventListener("router_update", this.refresh)
	}

	componentWillUnmount() {
		window.removeEventListener("router_update", this.refresh)
	}

	refresh = () => {
		this.setState({
			...this.state,
			path: window.location.hash,
		})
	}

	render(): JSX.Element {
		let path = window.location.hash.replace(/^#/, '')

		let curElem: JSX.Element;

		React.Children.forEach(this.props.children, (elem) => {
			if (React.isValidElement(elem)) {
				let data = RouterTypes.match(
					elem.props.path, path)

				if (data.matched) {
					RouterTypes.setState(data)
					curElem = elem
				}
			}
		})

		if (!curElem) {
			RouterTypes.setState(null)
			console.log(`Failed to match ${path}`)
		}

		return <>{curElem}</>
	}
}
