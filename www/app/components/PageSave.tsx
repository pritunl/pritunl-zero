/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	style?: React.CSSProperties;
	message: string;
	changed: boolean;
	disabled: boolean;
	hidden?: boolean;
	light?: boolean;
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
		let style: React.CSSProperties = this.props.light ? null : css.box;

		if (this.props.style) {
			style = {
				...style,
				...this.props.style,
			};
		}

		return <div
			className="layout horizontal"
			style={style}
			hidden={this.props.hidden && !this.props.children}
		>
			{this.props.children}
			<div className="flex"/>
			<div className="layout horizontal">
				<span style={css.message} hidden={!this.props.message}>
					{this.props.message}
				</span>
				<div style={css.buttons}>
					<button
						className="bp3-button bp3-icon-cross"
						style={css.button}
						hidden={this.props.hidden}
						type="button"
						disabled={!this.props.changed || this.props.disabled}
						onClick={this.props.onCancel}
					>
						Cancel
					</button>
					<button
						className="bp3-button bp3-intent-success bp3-icon-tick"
						style={css.button}
						hidden={this.props.hidden}
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
