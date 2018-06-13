/// <reference path="../References.d.ts"/>
import * as React from 'react';
import LoadingStore from '../stores/LoadingStore';

interface Props {
	style?: React.CSSProperties;
	size?: string;
	intent?: string;
}

interface State {
	loading: boolean;
}

const css = {
	progress: {
		width: '100%',
		height: '4px',
		borderTopLeftRadius: '3px',
		borderTopRightRadius: '3px',
		borderBottomLeftRadius: 0,
		borderBottomRightRadius: 0,
	} as React.CSSProperties,
	progressBar: {
		width: '50%',
		borderRadius: 0,
	} as React.CSSProperties,
};

export default class LoadingBar extends React.Component<Props, State> {
	constructor(props: Props, context: any) {
		super(props, context);
		this.state = {
			loading: LoadingStore.loading,
		};
	}

	componentDidMount(): void {
		LoadingStore.addChangeListener(this.onChange);
	}

	componentWillUnmount(): void {
		LoadingStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			loading: LoadingStore.loading,
		});
	}

	render(): JSX.Element {
		let progress: JSX.Element;

		if (!this.state.loading) {
			progress = <div style={css.progress}/>;
		} else {
			let className = 'pt-progress-bar pt-no-stripes pt-no-animation ';
			if (this.props.intent) {
				className += ' pt-intent-' + this.props.intent;
			}

			progress = <div className={className} style={css.progress}>
				<div
					className="pt-progress-meter pt-loading-bar"
					style={css.progressBar}
				/>
			</div>;
		}

		return <div style={this.props.style}>
			{progress}
		</div>;
	}
}
