/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Help from './Help';

interface Props {
	hidden?: boolean;
	disabled?: boolean;
	label: string;
	help: string;
	checked: boolean;
	onToggle: () => void;
}

const css = {
	label: {
		display: 'inline-block',
	} as React.CSSProperties,
};

export default class PageSwitch extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <div hidden={this.props.hidden}>
			<label className="pt-control pt-switch" style={css.label}>
				<input
					type="checkbox"
					disabled={this.props.disabled}
					checked={!!this.props.checked}
					onChange={this.props.onToggle}
				/>
				<span className="pt-control-indicator"/>
				{this.props.label}
				<Help
					title={this.props.label}
					content={this.props.help}
				/>
			</label>
		</div>;
	}
}
