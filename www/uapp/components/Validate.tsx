/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	token: string;
}

const css = {
	body: {
		padding: '0 10px',
	} as React.CSSProperties,
	description: {
		opacity: 0.7,
	} as React.CSSProperties,
	buttons: {
		marginTop: '15px',
	} as React.CSSProperties,
	button: {
		margin: '5px',
		width: '116px',
	} as React.CSSProperties,
};

export default class Validate extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <div>
			<div className="pt-non-ideal-state" style={css.body}>
				<div className="pt-non-ideal-state-visual pt-non-ideal-state-icon">
					<span className="pt-icon pt-icon-endorsed"/>
				</div>
				<h4 className="pt-non-ideal-state-title">Validate SSH Key</h4>
				<span style={css.description}>If you did not initiate this validation deny the request and report the incident to an administrator</span>
			</div>
			<div className="layout horizontal center-justified" style={css.buttons}>
				<button
					className="pt-button pt-large pt-intent-success pt-icon-add"
					style={css.button}
					type="button"
					onClick={(): void => {
					}}
				>
					Approve
				</button>
				<button
					className="pt-button pt-large pt-intent-danger pt-icon-delete"
					style={css.button}
					type="button"
					onClick={(): void => {
					}}
				>
					Deny
				</button>
			</div>
		</div>;
	}
}
