/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	children?: React.ReactNode
	hidden?: boolean;
	label: string;
	value: string;
	disabled?: boolean;
	buttonClass?: string;
	onChange: (val: string) => void;
	onSubmit: () => void;
}

const css = {
	group: {
		marginBottom: '15px',
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	select: {
		width: '100%',
		borderTopLeftRadius: '3px',
		borderBottomLeftRadius: '3px',
	} as React.CSSProperties,
	selectInner: {
		width: '100%',
	} as React.CSSProperties,
	selectBox: {
		flex: '1',
	} as React.CSSProperties,
};

export default class PageSelectButton extends React.Component<Props, {}> {
	render(): JSX.Element {
		let buttonClass = 'bp5-button';
		if (this.props.buttonClass) {
			buttonClass += ' ' + this.props.buttonClass;
		}

		return <div
			className="bp5-control-group"
			style={css.group}
			hidden={this.props.hidden}
		>
			<div style={css.selectBox}>
				<div className="bp5-select" style={css.select}>
					<select
						style={css.selectInner}
						disabled={this.props.disabled}
						value={this.props.value || ''}
						onChange={(evt): void => {
							this.props.onChange(evt.target.value);
						}}
					>
						{this.props.children}
					</select>
				</div>
			</div>
			<button
				className={buttonClass}
				disabled={this.props.disabled}
				onClick={this.props.onSubmit}
			>{this.props.label}</button>
		</div>;
	}
}
