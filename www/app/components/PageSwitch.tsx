/// <reference path="../References.d.ts"/>
import * as React from 'react';

type OnToggle = () => void;

interface Props {
	label: string;
	checked: boolean;
	onToggle: OnToggle;
}

const css = {
	label: {
		display: 'inline-block',
	} as React.CSSProperties,
};

export default class PageSwitch extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <div>
			<label className="pt-control pt-switch" style={css.label}>
				<input
					type="checkbox"
					checked={this.props.checked}
					onChange={this.props.onToggle}
				/>
				<span className="pt-control-indicator"/>
				{this.props.label}
			</label>
		</div>;
	}
}
