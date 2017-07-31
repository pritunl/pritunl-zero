/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	hidden?: boolean;
	label: string;
	type: string;
	placeholder: string;
	value: string | number;
	defaultValue: string;
	onChange: (val: string) => void;
}

interface State {
	checked?: boolean;
}

const css = {
	switchLabel: {
		display: 'inline-block',
		marginBottom: 0,
	} as React.CSSProperties,
	inputLabel: {
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	input: {
		width: '100%',
	} as React.CSSProperties,
};

export default class PageInputSwitch extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			checked: false,
		};
	}

	render(): JSX.Element {
		return <div hidden={this.props.hidden}>
			<label className="pt-control pt-switch" style={css.switchLabel}>
				<input
					type="checkbox"
					checked={!!this.props.value || this.state.checked}
					onChange={(): void => {
						if (!!this.props.value || this.state.checked) {
							this.setState({
								...this.state,
								checked: false,
							});
							this.props.onChange(null);
						} else {
							this.setState({
								...this.state,
								checked: true,
							});
							this.props.onChange(this.props.defaultValue);
						}
					}}
				/>
				<span className="pt-control-indicator"/>
				{this.props.label}
			</label>
			<label className="pt-label" style={css.inputLabel}>
				<input
					className="pt-input"
					style={css.input}
					hidden={!this.props.value && !this.state.checked}
					type={this.props.type}
					autoCapitalize="off"
					spellCheck={false}
					placeholder={this.props.placeholder}
					value={this.props.value || ''}
					onChange={(evt): void => {
						this.setState({
							...this.state,
							checked: true,
						});
						this.props.onChange(evt.target.value);
					}}
				/>
			</label>
		</div>;
	}
}
