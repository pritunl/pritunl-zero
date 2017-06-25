/// <reference path="../References.d.ts"/>
import * as React from 'react';

type OnChange = (val: string) => void;

interface Props {
	label: string;
	type: string;
	placeholder: string;
	value: string;
	onChange: OnChange;
}

const css = {
	label: {
		width: '100%',
		maxWidth: '310px',
	} as React.CSSProperties,
	input: {
		width: '100%',
	} as React.CSSProperties,
};

export default class PageInput extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <label className="pt-label" style={css.label}>
			{this.props.label}
			<input
				className="pt-input"
				style={css.input}
				type={this.props.type}
				autoCapitalize="off"
				spellCheck={false}
				placeholder={this.props.placeholder}
				value={this.props.value || ''}
				onChange={(evt): void => {
					this.props.onChange(evt.target.value);
				}}
			/>
		</label>;
	}
}
