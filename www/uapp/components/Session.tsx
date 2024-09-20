/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	onDevices: () => void;
}

const css = {
	body: {
		padding: 0,
		textAlign: 'center',
	} as React.CSSProperties,
	title: {
		margin: '10px 0 15px 0',
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
			<div style={css.body}>
				<h4 style={css.title}>
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
					className="bp5-button bp5-large bp5-intent-primary bp5-icon-download"
					style={css.button}
					href="https://docs.pritunl.com/v1/docs/ssh-client"
				>
					Install SSH Client
				</a>
				<button
					className="bp5-button bp5-large bp5-intent-success bp5-icon-id-number"
					style={css.button}
					onClick={this.props.onDevices}
				>
					Security Devices
				</button>
				<a
					className="bp5-button bp5-large bp5-intent-warning bp5-icon-delete"
					style={css.button}
					href="/logout"
				>
					Logout
				</a>
				<a
					className="bp5-button bp5-large bp5-intent-danger bp5-icon-trash"
					style={css.button}
					href="/logout_all"
				>
					End All Sessions
				</a>
			</div>
		</div>;
	}
}
