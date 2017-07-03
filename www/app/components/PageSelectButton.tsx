/// <reference path="../References.d.ts"/>
import * as React from 'react';

type OnChange = (val: string) => void;
type OnSubmit = () => void;

interface Props {
	label: string;
	value: string;
	buttonClass?: string;
	onChange: OnChange;
	onSubmit: OnSubmit;
}

const css = {
	group: {
		marginBottom: '15px',
	} as React.CSSProperties,
};

export default class PageSelectButton extends React.Component<Props, {}> {
	render(): JSX.Element {
		let buttonClass = 'pt-button';
		if (this.props.buttonClass) {
			buttonClass += ' ' + this.props.buttonClass;
		}

		return <div className="pt-control-group" style={css.group}>
			<div className="pt-select">
				<select
					value={this.props.value || ''}
					onChange={(evt): void => {
						this.props.onChange(evt.target.value);
					}}
				>
					{this.props.children}
				</select>
			</div>
			<button
				className={buttonClass}
				onClick={this.props.onSubmit}
			>{this.props.label}</button>
		</div>;
	}
}
