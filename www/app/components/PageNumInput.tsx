/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';
import Help from './Help';

interface Props {
	hidden?: boolean;
	disabled?: boolean;
	min?: number;
	max?: number;
	minorStepSize?: number;
	stepSize?: number;
	majorStepSize?: number;
	selectAllOnFocus?: true;
	label: string;
	help: string;
	value: number;
	onChange: (val: number) => void;
}

const css = {
	label: {
		display: 'inline-block',
	} as React.CSSProperties,
};

export default class PageNumInput extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <div hidden={this.props.hidden}>
			<label className="bp3-label" style={css.label}>
				{this.props.label}
				<Help
					title={this.props.label}
					content={this.props.help}
				/>
				<Blueprint.NumericInput
					allowNumericCharactersOnly={true}
					min={this.props.min}
					minorStepSize={this.props.minorStepSize}
					stepSize={this.props.stepSize}
					majorStepSize={this.props.majorStepSize}
					disabled={this.props.disabled}
					selectAllOnFocus={this.props.selectAllOnFocus}
					onValueChange={this.props.onChange}
					value={this.props.value}
				/>
			</label>
		</div>;
	}
}
