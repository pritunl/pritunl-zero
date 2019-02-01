/// <reference path="../References.d.ts"/>
import * as React from 'react';

type OnToggle = () => void;

interface Props {
	style: React.CSSProperties;
	label: string;
	checked: boolean;
	onToggle: OnToggle;
}

export default class Switch extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <label className="bp3-control bp3-switch" style={this.props.style}>
			<input
				type="checkbox"
				checked={this.props.checked}
				onChange={this.props.onToggle}
			/>
			<span className="bp3-control-indicator"/>
			{this.props.label}
		</label>;
	}
}
