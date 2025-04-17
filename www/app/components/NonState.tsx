/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Constants from '../Constants';

interface Props {
	hidden: boolean;
	iconClass: string;
	title: string;
	noDelay?: boolean;
	description?: string;
}

interface State {
	initialized: boolean;
}

const css = {
	state: {
		height: 'auto',
	} as React.CSSProperties,
	stateIcon: {
		fontSize: '48px',
		lineHeight: '48px',
		marginBottom: '0',
	} as React.CSSProperties,
};

export default class NonState extends React.Component<Props, State> {
	timeout: number;

	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			initialized: false,
		};
	}

	componentDidMount(): void {
		if (!this.props.noDelay) {
			this.timeout = window.setTimeout((): void => {
				this.setState({
					...this.state,
					initialized: true,
				});
			}, Constants.loadDelay);
		}
	}

	componentWillUnmount(): void {
		if (this.timeout) {
			window.clearTimeout(this.timeout);
		}
	}

	render(): JSX.Element {
		let description: JSX.Element;
		if (this.props.description) {
			description = <div className="bp5-non-ideal-state-description">
				{this.props.description}
			</div>;
		}

		return <div
			className="bp5-non-ideal-state"
			style={css.state}
			hidden={this.props.hidden || (!this.state.initialized && !this.props.noDelay)}
		>
			<div
				className="bp5-non-ideal-state-visual bp5-non-ideal-state-icon"
				style={css.stateIcon}
			>
				<span className={'bp5-icon ' + this.props.iconClass}/>
			</div>
			<h4 className="bp5-non-ideal-state-title">{this.props.title}</h4>
			{description}
		</div>;
	}
}
