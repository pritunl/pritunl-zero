/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	buttonClass?: string;
	hidden?: boolean;
	disabled?: boolean;
	readOnly?: boolean;
	label: string;
	type: string;
	placeholder: string;
	value: string;
	onChange?: (val: string) => void;
	onSubmit: () => void;
}

const css = {
	group: {
		marginBottom: '15px',
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	input: {
		width: '100%',
	} as React.CSSProperties,
	inputBox: {
		flex: '1',
	} as React.CSSProperties,
};

export default class PageInputButton extends React.Component<Props, {}> {
	render(): JSX.Element {
		let buttonClass = 'pt-button';
		if (this.props.buttonClass) {
			buttonClass += ' ' + this.props.buttonClass;
		}

		return <div
			className="pt-control-group"
			style={css.group}
			hidden={this.props.hidden}
		>
			<div style={css.inputBox}>
				<input
					className="pt-input"
					style={css.input}
					type={this.props.type}
					disabled={this.props.disabled}
					readOnly={this.props.readOnly}
					autoCapitalize="off"
					spellCheck={false}
					placeholder={this.props.placeholder}
					value={this.props.value || ''}
					onChange={(evt): void => {
						if (this.props.onChange) {
							this.props.onChange(evt.target.value);
						}
					}}
					onKeyPress={(evt): void => {
						if (evt.key === 'Enter') {
							this.props.onSubmit();
						}
					}}
				/>
			</div>
			<div>
				<button
					className={buttonClass}
					disabled={this.props.disabled}
					onClick={this.props.onSubmit}
				>{this.props.label}</button>
			</div>
		</div>;
	}
}
