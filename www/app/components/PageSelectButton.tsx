/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	label: string;
	value: string;
	buttonClass?: string;
	onChange: (val: string) => void;
	onSubmit: () => void;
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
