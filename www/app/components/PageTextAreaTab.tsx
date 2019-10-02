/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Help from './Help';

interface Props {
	hidden?: boolean;
	disabled?: boolean;
	readOnly?: boolean;
	label: string;
	help: string;
	placeholder: string;
	rows: number;
	tabs: string[];
	values: string[];
	onChange: (tab: string, val: string) => void;
}

interface State {
	activeIndex: number;
}

const css = {
	label: {
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	textarea: {
		width: '100%',
		resize: 'none',
		fontSize: '12px',
		fontFamily: '"Lucida Console", Monaco, monospace',
	} as React.CSSProperties,
	tab: {
		fontSize: '12px',
		lineHeight: '24px',
		userSelect: 'none',
	} as React.CSSProperties,
};

export default class PageTextAreaTab extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			activeIndex: 0,
		};
	}

	render(): JSX.Element {
		let activeIndex = this.state.activeIndex || 0;

		let tabs: JSX.Element[] = [];
		for (let i = 0; i < (this.props.tabs || []).length; i++) {
			let tab = this.props.tabs[i];
			let index = i;

			tabs.push(
				<li
					key={i}
					className="bp3-tab"
					style={css.tab}
					role="tab"
					aria-selected={i == activeIndex}
					onClick={(): void => {
						this.setState({
							...this.state,
							activeIndex: index,
						});
					}}
				>{tab}</li>
			);
		}

		return <label
			className="bp3-label"
			style={css.label}
			hidden={this.props.hidden}
		>
			{this.props.label}
			<Help
				title={this.props.label}
				content={this.props.help}
			/>
			<div className="bp3-tabs">
				<ul className="bp3-tab-list .modifier" role="tablist">
					{tabs}
				</ul>
			</div>
			<textarea
				className="bp3-input"
				style={css.textarea}
				disabled={this.props.disabled}
				readOnly={this.props.readOnly}
				autoCapitalize="off"
				spellCheck={false}
				placeholder={this.props.placeholder}
				rows={this.props.rows}
				value={this.props.values[activeIndex] || ''}
				onChange={(evt): void => {
					this.props.onChange(
						this.props.tabs[this.state.activeIndex],
						evt.target.value,
					);
				}}
			/>
		</label>;
	}
}
