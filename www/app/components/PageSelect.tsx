/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Help from './Help';

interface Props {
	children?: React.ReactNode
	hidden?: boolean;
	disabled?: boolean;
	label: string;
	help: string;
	value: string;
	onChange: (val: string) => void;
}

const css = {
	label: {
		display: 'inline-block',
	} as React.CSSProperties,
};

export default class PageSelect extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <div hidden={this.props.hidden}>
			<label className="bp5-label" style={css.label}>
				{this.props.label}
				<Help
					title={this.props.label}
					content={this.props.help}
				/>
				<div className="bp5-select">
					<select
						disabled={this.props.disabled}
						value={this.props.value || ''}
						onChange={(evt): void => {
							this.props.onChange(evt.target.value);
						}}
					>
						{this.props.children}
					</select>
				</div>
			</label>
		</div>;
	}
}
