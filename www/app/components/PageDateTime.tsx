/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as BlueprintDateTime from '@blueprintjs/datetime';
import Help from './Help';

interface Props {
	hidden?: boolean;
	disabled?: boolean;
	label: string;
	help: string;
	value: string;
	onChange: (val: string) => void;
}

const css = {
	group: {
		display: 'inline-block',
	} as React.CSSProperties,
	label: {
		marginBottom: '5px',
	} as React.CSSProperties,
};

export default class PageDateTime extends React.Component<Props, {}> {
	render(): JSX.Element {
		let dateStyle: React.CSSProperties = {};

		let date = new Date(this.props.value);
		if (!this.props.value || this.props.value === '0001-01-01T00:00:00Z') {
			date = null;
		}

		if (!date || this.props.disabled) {
			dateStyle.opacity = 0.5;
		}

		return <div hidden={this.props.hidden}>
			<div style={css.group}>
				<label className="bp3-label" style={css.label}>
					{this.props.label}
					<Help
						title={this.props.label}
						content={this.props.help}
					/>
				</label>
				<div style={dateStyle}>
					<BlueprintDateTime.DateTimePicker
						value={this.props.disabled ? null : date}
						timePickerProps={{
							showArrowButtons: true,
						}}
						datePickerProps={{
							showActionsBar: true,
						}}
						onChange={(newDate: Date): void => {
							if (this.props.disabled) {
								return;
							}

							if (newDate) {
								this.props.onChange(newDate.toJSON());
							} else {
								this.props.onChange(null);
							}
						}}
					/>
				</div>
			</div>
		</div>;
	}
}
