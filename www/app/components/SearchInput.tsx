/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';

type OnChange = (val: string) => void;

interface Props {
	style: React.CSSProperties;
	placeholder: string;
	value: string;
	dynamic?: boolean;
	exactDefault?: boolean;
	onChange: OnChange;
}

interface State {
	exact: boolean;
	value: string;
}

export default class SearchInput extends React.Component<Props, State> {
	pending: string[] = [];

	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			exact: !!this.props.exactDefault,
			value: this.parseValue(this.props.value),
		};
	}

	parseValue(value: string): string {
		let val = value || "";

		if (this.props.dynamic && val.startsWith("~")) {
			val = val.substring(1);
		}

		return val;
	}

	componentDidUpdate(prevProps: Props): void {
		let value = this.props.value || "";
		if (value === (prevProps.value || "")) {
			return;
		}

		let index = this.pending.indexOf(value);
		if (index !== -1) {
			this.pending.splice(0, index + 1);
			return;
		}

		this.pending = [];

		let val = this.parseValue(value);
		if (val !== this.state.value) {
			this.setState({
				...this.state,
				value: val,
			});
		}
	}

	emit(val: string, exact: boolean): void {
		let value: string;

		if (this.props.dynamic && !exact && val) {
			value = "~" + val;
		} else {
			value = val;
		}

		this.pending.push(value);
		this.props.onChange(value);
	}

	render(): JSX.Element {
		let val = this.state.value;

		return <div style={this.props.style}>
			<Blueprint.InputGroup
				type="text"
				autoCapitalize="off"
				spellCheck={false}
				placeholder={this.props.placeholder}
				value={val}
				onChange={(evt): void => {
					let value = evt.target.value;

					this.setState({
						...this.state,
						value: value,
					});

					this.emit(value, this.state.exact);
				}}
				leftElement={<span className="bp5-icon bp5-icon-search"/>}
				rightElement={<Blueprint.Tooltip
					content={this.state.exact ? "Exact match" : "Partial match"}
				>
					<Blueprint.Button
						hidden={!this.props.dynamic}
						icon={this.state.exact ? "link" : "unlink"}
						onClick={() => {
							let exact = !this.state.exact;

							this.setState({
								...this.state,
								exact: exact,
							});

							this.emit(val, exact);
						}}
						minimal={true}
					/>
				</Blueprint.Tooltip>}
			/>
		</div>
	}
}
