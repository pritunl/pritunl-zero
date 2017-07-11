/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	hidden?: boolean;
	label: string;
	type: string;
	placeholder: string;
	rows: number;
	value: string;
	onChange: (val: string) => void;
}

const css = {
	label: {
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	textarea: {
		width: '100%',
		resize: 'none',
		fontFamily: '"Lucida Console", Monaco, monospace',
	} as React.CSSProperties,
};

export default class PageTextArea extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <label
			className="pt-label"
			style={css.label}
			hidden={this.props.hidden}
		>
			{this.props.label}
			<textarea
				className="pt-input"
				style={css.textarea}
				type={this.props.type}
				autoCapitalize="off"
				spellCheck={false}
				placeholder={this.props.placeholder}
				rows={this.props.rows}
				value={this.props.value || ''}
				onChange={(evt): void => {
					this.props.onChange(evt.target.value);
				}}
			/>
		</label>;
	}
}
