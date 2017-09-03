/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Help from './Help';

interface Props {
	hidden?: boolean;
	disabled?: boolean;
	label: string;
	help: string;
	type: string;
	placeholder: string;
	value: string | number;
	onChange: (val: string) => void;
}

const css = {
	label: {
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	input: {
		width: '100%',
	} as React.CSSProperties,
};

export default class PageInput extends React.Component<Props, {}> {
	render(): JSX.Element {
		let value: any = this.props.value;
		value = isNaN(value) ? this.props.value || '' : this.props.value;

		return <label
			className="pt-label"
			style={css.label}
			hidden={this.props.hidden}
		>
			{this.props.label}
			<Help
				title={this.props.label}
				content={this.props.help}
			/>
			<input
				className="pt-input"
				style={css.input}
				type={this.props.type}
				disabled={this.props.disabled}
				autoCapitalize="off"
				spellCheck={false}
				placeholder={this.props.placeholder}
				value={value}
				onChange={(evt): void => {
					this.props.onChange(evt.target.value);
				}}
			/>
		</label>;
	}
}
