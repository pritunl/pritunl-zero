/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	style?: React.CSSProperties;
	message: string;
	changed: boolean;
	disabled: boolean;
	closed?: boolean;
	hidden?: boolean;
	light?: boolean;
	onCancel: () => void;
	onCreate: () => void;
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

export default class PageCreate extends React.Component<Props, {}> {
	render(): JSX.Element {
		let style: React.CSSProperties = this.props.light ? null : css.box;

		if (this.props.style) {
			style = {
				...style,
				...this.props.style,
			};
		}

		let closedDom: JSX.Element;
		if (this.props.closed) {
			closedDom = <button
				className="bp5-button bp5-intent-success bp5-icon-cross"
				style={css.button}
				type="button"
				onClick={this.props.onCancel}
			>
				Close
			</button>;
		}

		return <div
			className="layout horizontal"
			style={style}
			hidden={this.props.hidden}
		>
			<div className="flex"/>
			<div className="layout horizontal">
				<span style={css.message} hidden={!this.props.message}>
					{this.props.message}
				</span>
				<div style={css.buttons}>
					<button
						className="bp5-button bp5-icon-cross"
						style={css.button}
						type="button"
						hidden={this.props.closed}
						disabled={this.props.disabled}
						onClick={this.props.onCancel}
					>
						Cancel
					</button>
					<button
						className="bp5-button bp5-intent-success bp5-icon-tick"
						style={css.button}
						type="button"
						hidden={this.props.closed}
						disabled={!this.props.changed || this.props.disabled}
						onClick={this.props.onCreate}
					>
						Create
					</button>
					{closedDom}
				</div>
			</div>
		</div>;
	}
}
