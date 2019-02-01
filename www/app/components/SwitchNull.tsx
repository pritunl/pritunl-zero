/// <reference path="../References.d.ts"/>
import * as React from 'react';

type OnToggle = () => void;

interface Props {
	style?: React.CSSProperties;
	label: string;
	checked: boolean;
	onToggle: OnToggle;
}

export default class SwitchNull extends React.Component<Props, {}> {
	render(): JSX.Element {
		let style: React.CSSProperties = {
			...this.props.style,
		};

		if (this.props.checked === null || this.props.checked === undefined) {
			style.opacity = 0.5;
		}

		return <label className="bp3-control bp3-switch" style={style}>
			<input
				type="checkbox"
				checked={!!this.props.checked}
				onChange={this.props.onToggle}
			/>
			<span className="bp3-control-indicator"/>
			{this.props.label}
		</label>;
	}
}
