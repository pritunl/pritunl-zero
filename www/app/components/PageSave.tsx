/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	message: string;
	changed: boolean;
	disabled: boolean;
	hidden?: boolean;
	onCancel: () => void;
	onSave: () => void;
}

const css = {
	message: {
		marginTop: '6px',
	} as React.CSSProperties,
	box: {
		marginTop: '15px',
	} as React.CSSProperties,
	button: {
		marginLeft: '10px',
	} as React.CSSProperties,
	buttons: {
		flexShrink: 0,
	} as React.CSSProperties,
};

export default class PageSave extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <div
			className="layout horizontal"
			style={css.box}
			hidden={this.props.hidden}
		>
			<div className="flex"/>
			<div className="layout horizontal">
				<span style={css.message} hidden={!this.props.message}>
					{this.props.message}
				</span>
				<div style={css.buttons}>
					<button
						className="pt-button pt-icon-cross"
						style={css.button}
						type="button"
						disabled={!this.props.changed || this.props.disabled}
						onClick={this.props.onCancel}
					>
						Cancel
					</button>
					<button
						className="pt-button pt-intent-success pt-icon-tick"
						style={css.button}
						type="button"
						disabled={!this.props.changed || this.props.disabled}
						onClick={this.props.onSave}
					>
						Save
					</button>
				</div>
			</div>
		</div>;
	}
}
