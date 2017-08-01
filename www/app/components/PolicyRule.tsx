/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as PolicyTypes from '../types/PolicyTypes';
import * as Constants from '../Constants';
import PageSwitch from './PageSwitch';
import PageSelectButton from './PageSelectButton';

interface Props {
	rule: PolicyTypes.Rule;
	onChange: (state: PolicyTypes.Rule) => void;
}

interface State {
	addValue: string;
}

const css = {
	item: {
		margin: '9px 5px 0 5px',
		height: '20px',
	} as React.CSSProperties,
};

export default class PolicyRule extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			addValue: '',
		};
	}

	clone(): PolicyTypes.Rule {
		return {
			...this.props.rule,
		};
	}

	onAddValue = (value: string): void => {
		let rule = this.clone();

		let values = [
			...rule.values,
		];

		if (values.indexOf(value) === -1) {
			values.push(value);
		}

		values.sort();

		rule.values = values;

		this.props.onChange(rule);

		this.setState({
			...this.state,
		});
	}

	onRemoveValue(value: string): void {
		let rule = this.clone();

		let values = [
			...rule.values,
		];

		let i = values.indexOf(value);
		if (i === -1) {
			return;
		}

		values.splice(i, 1);

		rule.values = values;

		this.props.onChange(rule);
	}

	render(): JSX.Element {
		let rule = this.props.rule;
		let defaultOption: string;

		let label: string;
		let selectLabel: string;
		let options: {[key: string]: string};
		switch (this.props.rule.type) {
			case 'operating_system':
				label = 'Permitted Operating Systems';
				selectLabel = 'Operating Systems';
				options = Constants.operatingSystems;
				break;
			case 'browser':
				label = 'Permitted Browsers';
				selectLabel = 'Browsers';
				options = Constants.browsers;
				break;
			case 'location':
				label = 'Permitted Locations';
				selectLabel = 'Locations';
				options = Constants.locations;
				break;
		}

		let optionsSelect: JSX.Element[] = [];
		for (let option in options) {
			if (!options.hasOwnProperty(option)) {
				continue;
			}
			if (!defaultOption) {
				defaultOption = option;
			}

			optionsSelect.push(
				<option key={option} value={option}>{options[option]}</option>,
			);
		}

		let values: JSX.Element[] = [];
		for (let value of rule.values || []) {
			values.push(
				<div
					className="pt-tag pt-tag-removable pt-intent-primary"
					style={css.item}
					key={value}
				>
					{options[value] || value}
					<button
						className="pt-tag-remove"
						onMouseUp={(): void => {
							this.onRemoveValue(value);
						}}
					/>
				</div>,
			);
		}

		return <div>
			<PageSwitch
				label={selectLabel}
				checked={rule.values != null}
				onToggle={(): void => {
					let state = this.clone();
					state.values = rule.values == null ? [] : null;
					this.props.onChange(state);
				}}
			/>
			<PageSwitch
				label="Disabled user on failure"
				checked={rule.disable}
				hidden={rule.values == null}
				onToggle={(): void => {
					let state = this.clone();
					state.disable = !state.disable;
					this.props.onChange(state);
				}}
			/>
			<label
				className="pt-label"
				hidden={rule.values == null}
			>
				{label}
				<div>
					{values}
				</div>
			</label>
			<PageSelectButton
				hidden={rule.values == null}
				buttonClass="pt-intent-success pt-icon-add"
				label="Add"
				value={this.state.addValue}
				onChange={(val): void => {
					this.setState({
						...this.state,
						addValue: val,
					});
				}}
				onSubmit={(): void => {
					this.onAddValue(this.state.addValue || defaultOption);
				}}
			>
				{optionsSelect}
			</PageSelectButton>
		</div>;
	}
}
