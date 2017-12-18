/// <reference path="../References.d.ts"/>
import * as React from 'react';

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

export default class Session extends React.Component<{}, {}> {
	render(): JSX.Element {
		return <div>
			<div className="pt-non-ideal-state" style={css.body}>
				<h4 className="pt-non-ideal-state-title">
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
					className="pt-button pt-large pt-intent-primary pt-icon-download"
					style={css.button}
					href="https://docs.pritunl.com/v1/docs/ssh-client"
				>
					Install SSH Client
				</a>
				<a
					className="pt-button pt-large pt-intent-warning pt-icon-delete"
					style={css.button}
					href="/logout"
				>
					Logout
				</a>
				<a
					className="pt-button pt-large pt-intent-danger pt-icon-trash"
					style={css.button}
					href="/logout_all"
				>
					End All Sessions
				</a>
			</div>
		</div>;
	}
}
