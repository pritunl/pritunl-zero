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
		return <div className="pt-input-group" style={this.props.style}>
			<span className="pt-icon pt-icon-search"/>
			<input
				className="pt-input pt-round"
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
