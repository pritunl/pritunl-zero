/// <reference path="../References.d.ts"/>
import * as React from 'react';

type OnChange = (val: string) => void;

interface Props {
	style: React.CSSProperties;
	placeholder: string;
	value: string;
	onChange: OnChange;
}

export default class SearchInput extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <div className="bp3-input-group" style={this.props.style}>
			<span className="bp3-icon bp3-icon-search"/>
			<input
				className="bp3-input bp3-round"
				type="text"
				autoCapitalize="off"
				spellCheck={false}
				placeholder={this.props.placeholder}
				value={this.props.value || ''}
				onChange={(evt): void => {
					this.props.onChange(evt.target.value);
				}}
			/>
		</div>;
	}
}
