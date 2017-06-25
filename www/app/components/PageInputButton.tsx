/// <reference path="../References.d.ts"/>
import * as React from 'react';

type OnChange = (val: string) => void;
type OnSubmit = () => void;

interface Props {
	label: string;
	type: string;
	placeholder: string;
	value: string;
	onChange: OnChange;
	onSubmit: OnSubmit;
}

const css = {
};

export default class PageInputButton extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <div className="pt-control-group">
			<input
				className="pt-input"
				type={this.props.type}
				autoCapitalize="off"
				spellCheck={false}
				placeholder={this.props.placeholder}
				value={this.props.value || ''}
				onChange={(evt): void => {
					this.props.onChange(evt.target.value);
				}}
				onKeyPress={(evt): void => {
					if (evt.key === 'Enter') {
						this.props.onSubmit();
					}
				}}
			/>
			<button
				className="pt-button"
				onClick={this.props.onSubmit}
			>{this.props.label}</button>
		</div>;
	}
}
