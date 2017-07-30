/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as PolicyTypes from '../types/PolicyTypes';
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

const operatingSystems: {[key: string]: string} = {
	linux: 'Linux',
	macos_1010: 'macOS 10.10',
	macos_1011: 'macOS 10.11',
	macos_1012: 'macOS 10.12',
	macos_1013: 'macOS 10.13',
	windows_xp: 'Windows XP',
	windows_7: 'Windows 7',
	windows_vista: 'Windows Vista',
	windows_8: 'Windows 8',
	windows_10: 'Windows 10',
	chrome_os: 'Chrome OS',
	ios_8: 'iOS 8',
	ios_9: 'iOS 9',
	ios_10: 'iOS 10',
	ios_11: 'iOS 11',
	ios_12: 'iOS 12',
	android_4: 'Android KitKat 4.4',
	android_5: 'Android Lollipop 5.0',
	android_6: 'Android Marshmallow 6.0',
	android_7: 'Android Nougat 7.0',
	android_8: 'Android 8.0',
	blackberry_10: 'Blackerry 10',
	windows_phone: 'Windows Phone',
	firefox_os: 'Firefox OS',
	kindle: 'Kindle',
};

const browsers: {[key: string]: string} = {
	chrome: 'Chrome',
	chrome_mobile: 'Chrome Mobile',
	safari: 'Safari',
	safari_mobile: 'Safari Mobile',
	firefox: 'Firefox',
	firefox_mobile: 'Firefox Mobile',
	edge: 'Microsoft Edge',
	internet_explorer: 'Internet Explorer',
	internet_explorer_mobile: 'Internet Explorer Mobile',
	opera: 'Opera',
	opera_mobile: 'Opera Mobile',
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

	optionsSelect(): JSX.Element {
		let rule = this.props.rule;
		let defaultOption: string;

		let label: string;
		let selectLabel: string;
		let options: {[key: string]: string};
		switch (this.props.rule.type) {
			case 'operating_system':
				label = 'Permitted Operating Systems';
				selectLabel = 'Operating Systems';
				options = operatingSystems;
				break;
			case 'browser':
				label = 'Permitted Browsers';
				selectLabel = 'Browsers';
				options = browsers;
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

	render(): JSX.Element {
		let rule = this.props.rule;

		let options: JSX.Element;
		switch (rule.type) {
			case 'operating_system':
				options = this.optionsSelect();
				break;
			case 'browser':
				options = this.optionsSelect();
				break;
		}

		return <div>
			{options}
		</div>;
	}
}
