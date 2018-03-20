/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as PolicyTypes from '../types/PolicyTypes';
import * as Constants from '../Constants';
import PageSwitch from './PageSwitch';
import PageSelectButton from './PageSelectButton';
import PageInputButton from './PageInputButton';
import Help from './Help';

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
				selectLabel = 'Operating system policies';
				options = Constants.operatingSystems;
				break;
			case 'browser':
				label = 'Permitted Browsers';
				selectLabel = 'Browser policies';
				options = Constants.browsers;
				break;
			case 'location':
				label = 'Permitted Locations';
				selectLabel = 'Location policies';
				options = Constants.locations;
				break;
			case 'cidr':
				label = 'Permitted CIDR Blocks';
				selectLabel = 'CIDR Blocks';
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
					{options ? options[value] || value : value}
					<button
						className="pt-tag-remove"
						onMouseUp={(): void => {
							this.onRemoveValue(value);
						}}
					/>
				</div>,
			);
		}

		let valueAppender: JSX.Element[] = []
		if (options) {
			valueAppender.push(
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
			)
		} else {
			valueAppender.push(
				<PageInputButton
					hidden={rule.values == null}
					buttonClass="pt-intent-success pt-icon-add"
					label="Add"
					type="text"
					value={this.state.addValue}
					onChange={(val): void => {
						this.setState({
							...this.state,
							addValue: val,
						});
					}}
					onSubmit={(): void => {
						this.onAddValue(this.state.addValue);
					}}
				/>
			)
		}

		return <div>
			<PageSwitch
				label={selectLabel}
				help="Turn on to enable policy."
				checked={rule.values != null}
				onToggle={(): void => {
					let state = this.clone();
					state.values = rule.values == null ? [] : null;
					this.props.onChange(state);
				}}
			/>
			<PageSwitch
				label="Disabled user on failure"
				help="This will disable the user when the policy check fails. It is generally only useful for the location check to disable a user account when an authentication occurs from a foreign country. It is important to consider that the policy check is the last check that occurs during authentication. An authentication attempt with an incorrect password from a foreign country would not trigger a policy failure or disable the user."
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
				<Help
					title={label}
					content="One of the values must match for the check to pass."
				/>
				<div>
					{values}
				</div>
			</label>
			{valueAppender}
		</div>;
	}
}
