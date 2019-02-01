/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	onDevices: () => void;
}

const css = {
	body: {
		padding: 0,
	} as React.CSSProperties,
	description: {
		opacity: 0.7,
		padding: '0 10px',
	} as React.CSSProperties,
	buttons: {
		marginTop: '15px',
	} as React.CSSProperties,
	button: {
		margin: '5px',
		width: '182px',
	} as React.CSSProperties,
};

export default class Session extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <div>
			<div className="bp3-non-ideal-state" style={css.body}>
				<h4 className="bp3-non-ideal-state-title">
					Pritunl Zero User Console
				</h4>
				<span style={css.description}>
					To validate an SSH key install the client and run "pritunl-ssh"
				</span>
			</div>
			<div
				className="layout horizontal center-justified wrap"
				style={css.buttons}
			>
				<a
					className="bp3-button bp3-large bp3-intent-primary bp3-icon-download"
					style={css.button}
					href="https://docs.pritunl.com/v1/docs/ssh-client"
				>
					Install SSH Client
				</a>
				<button
					className="bp3-button bp3-large bp3-intent-success bp3-icon-id-number"
					style={css.button}
					onClick={this.props.onDevices}
				>
					Security Devices
				</button>
				<a
					className="bp3-button bp3-large bp3-intent-warning bp3-icon-delete"
					style={css.button}
					href="/logout"
				>
					Logout
				</a>
				<a
					className="bp3-button bp3-large bp3-intent-danger bp3-icon-trash"
					style={css.button}
					href="/logout_all"
				>
					End All Sessions
				</a>
			</div>
		</div>;
	}
}
