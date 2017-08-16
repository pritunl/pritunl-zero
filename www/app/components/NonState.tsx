/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Constants from '../Constants';

interface Props {
	hidden: boolean;
	iconClass: string;
	title: string;
	description?: string;
}

interface State {
	initialized: boolean;
}

const css = {
	state: {
		height: 'auto',
	} as React.CSSProperties,
};

export default class NonState extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			initialized: false,
		};
	}

	componentDidMount(): void {
		setTimeout((): void => {
			this.setState({
				...this.state,
				initialized: true,
			});
		}, Constants.loadDelay);
	}

	render(): JSX.Element {
		let description: JSX.Element;
		if (this.props.description) {
			description = <div className="pt-non-ideal-state-description">
				{this.props.description}
			</div>;
		}

		return <div
			className="pt-non-ideal-state"
			style={css.state}
			hidden={this.props.hidden || !this.state.initialized}
		>
			<div className="pt-non-ideal-state-visual pt-non-ideal-state-icon">
				<span className={'pt-icon ' + this.props.iconClass}/>
			</div>
			<h4 className="pt-non-ideal-state-title">{this.props.title}</h4>
			{description}
		</div>;
	}
}
