/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as BlueprintDateTime from '@blueprintjs/datetime2';
import Help from './Help';

interface Props {
	hidden?: boolean;
	disabled?: boolean;
	label: string;
	help: string;
	value: string;
	onChange: (val: string) => void;
}

const minDate = new Date(2018, 0, 1);
const maxDate = new Date(new Date().getFullYear() + 5, 0, 0);

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
		if (!this.props.value ||
				this.props.value.includes('0000-12-31') ||
				this.props.value.includes('0001-01-01')) {
			date = null;
		}

		if (!date || this.props.disabled) {
			dateStyle.opacity = 0.5;
		}

		return <div hidden={this.props.hidden}>
			<div style={css.group}>
				<label className="bp5-label" style={css.label}>
					{this.props.label}
					<Help
						title={this.props.label}
						content={this.props.help}
					/>
				</label>
				<div style={dateStyle}>
					<BlueprintDateTime.DatePicker3
						value={this.props.disabled ? null : date}
						showActionsBar={true}
						minDate={minDate}
						maxDate={maxDate}
						timePickerProps={{
							showArrowButtons: true,
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
