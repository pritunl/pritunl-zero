/// <reference path="../References.d.ts"/>
import * as React from 'react';
import StateStore from '../stores/StateStore';
import LoadingBar from './LoadingBar';
import Session from './Session';
import Validate from './Validate';
import Devices from './Devices';

interface State {
	devicesOpen: boolean;
	sshToken: string;
	sshDevice: string;
}

const css = {
	card: {
		padding: '20px 15px',
		minWidth: '260px',
		maxWidth: '320px',
		margin: '30px auto',
	} as React.CSSProperties,
	loading: {
		position: 'absolute',
		top: 0,
		left: 0,
		right: 0,
	} as React.CSSProperties,
};

export default class Main extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			devicesOpen: false,
			sshToken: StateStore.sshToken,
			sshDevice: StateStore.sshDevice,
		};
	}

	componentDidMount(): void {
		StateStore.addChangeListener(this.onChange);
	}

	componentWillUnmount(): void {
		StateStore.addChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			sshToken: StateStore.sshToken,
			sshDevice: StateStore.sshDevice,
		});
	}

	render(): JSX.Element {
		let bodyElm: JSX.Element;

		if (this.state.sshToken) {
			bodyElm = <Validate token={this.state.sshToken}/>;
		} else if (this.state.devicesOpen || this.state.sshDevice) {
			bodyElm = <Devices
				onClose={(): void => {
					this.setState({
						...this.state,
						devicesOpen: false,
					});
				}}
			/>;
		} else {
			bodyElm = <Session
				onDevices={(): void => {
					this.setState({
						...this.state,
						devicesOpen: true,
					});
				}}
			/>;
		}

		return <div className="bp3-card bp3-elevation-2" style={css.card}>
			<LoadingBar style={css.loading} intent="primary"/>
			{bodyElm}
		</div>;
	}
}
