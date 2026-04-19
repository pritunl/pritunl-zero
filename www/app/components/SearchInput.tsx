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
	exact: boolean
}

export default class SearchInput extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			exact: !!this.props.exactDefault,
		};
	}

	render(): JSX.Element {
		let val = this.props.value || ""

		if (this.props.dynamic) {
			if (val.startsWith("~")) {
				val = val.substring(1)
			}
		}

		return <div style={this.props.style}>
			<Blueprint.InputGroup
				type="text"
				autoCapitalize="off"
				spellCheck={false}
				placeholder={this.props.placeholder}
				value={val}
				onChange={(evt): void => {
					if (this.props.dynamic && !this.state.exact && evt.target.value) {
						this.props.onChange("~" + evt.target.value)
					} else {
						this.props.onChange(evt.target.value)
					}
				}}
				leftElement={<span className="bp5-icon bp5-icon-search"/>}
				rightElement={<Blueprint.Tooltip
					content={this.state.exact ? "Exact match" : "Partial match"}
				>
					<Blueprint.Button
						hidden={!this.props.dynamic}
						icon={this.state.exact ? "link" : "unlink"}
						onClick={() => {
							let exact = !this.state.exact

							this.setState({
								...this.state,
								exact: exact,
							})

							if (!exact && val) {
								this.props.onChange("~" + val)
							} else {
								this.props.onChange(val)
							}
						}}
						minimal={true}
					/>
				</Blueprint.Tooltip>}
			/>
		</div>
	}
}
