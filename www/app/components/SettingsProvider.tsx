/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SettingsTypes from '../types/SettingsTypes';
import SettingsStore from '../stores/SettingsStore';
import * as SettingsActions from '../actions/SettingsActions';
import Page from './Page';
import PageHeader from './PageHeader';
import PagePanel from './PagePanel';
import PageSplit from './PageSplit';
import PageInput from './PageInput';
import PageSelect from './PageSelect';
import PageSave from './PageSave';

interface Props {
	provider: SettingsTypes.ProviderAny;
	onChange: (state: SettingsTypes.ProviderAny) => void;
}

interface State {
}

export default class SettingsProvider extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {};
	}

	clone(): SettingsTypes.ProviderAny {
		return {
			...this.props.provider,
		};
	}

	google(): JSX.Element {
		let provider = this.props.provider;

		return <div>
			<PageInput
				label="Domain"
				type="text"
				placeholder="Google domain to match"
				value={provider.domain}
				onChange={(val: string): void => {
					let state = this.clone();
					state.domain = val;
					this.props.onChange(state);
				}}
			/>
		</div>;
	}

	render(): JSX.Element {
		let provider = this.props.provider;
		let label = '';
		let options: JSX.Element;

		switch (provider.type) {
			case 'google':
				label = 'Google';
				options = this.google();
				break;
			case 'onelogin':
				label = 'OneLogin';
				break;
			case 'okta':
				label = 'Okta';
				break;
		}

		return <div className="pt-card">
			<h6>{label}</h6>
			<PageInput
				label="Label"
				type="text"
				placeholder="Provider label"
				value={provider.label}
				onChange={(val: string): void => {
					let state = this.clone();
					state.label = val;
					this.props.onChange(state);
				}}
			/>
			{options}
		</div>;
	}
}
