/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SettingsTypes from '../types/SettingsTypes';
import PageInput from './PageInput';
import PageSelect from './PageSelect';
import PageSwitch from './PageSwitch';
import PageInfo from './PageInfo';

interface Props {
	provider: SettingsTypes.SecondaryProviderAny;
	onChange: (state: SettingsTypes.SecondaryProviderAny) => void;
	onRemove: () => void;
}

const css = {
	label: {
		fontSize: '16px',
		margin: '0 0 7px 0',
	} as React.CSSProperties,
	card: {
		marginBottom: '5px',
	} as React.CSSProperties,
	role: {
		margin: '9px 5px 0 5px',
		minHeight: '20px',
	} as React.CSSProperties,
};

export default class SettingsSecondaryProvider extends React.Component<Props, {}> {
	clone(): SettingsTypes.SecondaryProviderAny {
		return {
			...this.props.provider,
		};
	}

	duo(): JSX.Element {
		let provider = this.props.provider;

		return <div>
			<PageInput
				label="Duo API Hostname"
				help="Duo API hostname found in Duo admin console."
				type="text"
				placeholder="Duo API hostname"
				value={provider.duo_hostname}
				onChange={(val: string): void => {
					let state = this.clone();
					state.duo_hostname = val;
					this.props.onChange(state);
				}}
			/>
			<PageInput
				label="Duo Integration Key"
				help="Duo integration key found in Duo admin console."
				type="text"
				placeholder="Duo integration key"
				value={provider.duo_key}
				onChange={(val: string): void => {
					let state = this.clone();
					state.duo_key = val;
					this.props.onChange(state);
				}}
			/>
			<PageInput
				label="Duo Secret Key"
				help="Duo secret key found in Duo admin console."
				type="text"
				placeholder="Duo secret key"
				value={provider.duo_secret}
				onChange={(val: string): void => {
					let state = this.clone();
					state.duo_secret = val;
					this.props.onChange(state);
				}}
			/>
			<PageSwitch
				label="Push authentication"
				help="Allow push authentication."
				checked={provider.push_factor}
				onToggle={(): void => {
					let state = this.clone();
					state.push_factor = !state.push_factor;
					this.props.onChange(state);
				}}
			/>
			<PageSwitch
				label="Phone authentication"
				help="Allow phone authentication."
				checked={provider.phone_factor}
				onToggle={(): void => {
					let state = this.clone();
					state.phone_factor = !state.phone_factor;
					this.props.onChange(state);
				}}
			/>
			<PageSwitch
				label="Passcode authentication"
				help="Allow passcode authentication."
				checked={provider.passcode_factor}
				onToggle={(): void => {
					let state = this.clone();
					state.passcode_factor = !state.passcode_factor;
					this.props.onChange(state);
				}}
			/>
			<PageSwitch
				label="SMS authentication"
				help="Allow SMS authentication."
				checked={provider.sms_factor}
				onToggle={(): void => {
					let state = this.clone();
					state.sms_factor = !state.sms_factor;
					this.props.onChange(state);
				}}
			/>
		</div>;
	}

	onelogin(): JSX.Element {
		let provider = this.props.provider;

		return <div>
			<PageInput
				label="OneLogin API Client ID"
				help="OneLogin API client ID found in OneLogin admin console."
				type="text"
				placeholder="OneLogin API client ID"
				value={provider.one_login_id}
				onChange={(val: string): void => {
					let state = this.clone();
					state.one_login_id = val;
					this.props.onChange(state);
				}}
			/>
			<PageInput
				label="OneLogin API Client Secret"
				help="OneLogin API client secret found in OneLogin admin console."
				type="text"
				placeholder="OneLogin API client secret"
				value={provider.one_login_secret}
				onChange={(val: string): void => {
					let state = this.clone();
					state.one_login_secret = val;
					this.props.onChange(state);
				}}
			/>
			<PageSelect
				label="OneLogin API Region"
				help="OneLogin region for API requests."
				value={provider.one_login_region}
				onChange={(val): void => {
					let state = this.clone();
					state.one_login_region = val;
					this.props.onChange(state);
				}}
			>
				<option value="us">United States</option>
				<option value="eu">Europe</option>
			</PageSelect>
			<PageSwitch
				label="Push authentication"
				help="Allow push authentication."
				checked={provider.push_factor}
				onToggle={(): void => {
					let state = this.clone();
					state.push_factor = !state.push_factor;
					this.props.onChange(state);
				}}
			/>
			<PageSwitch
				label="Passcode authentication"
				help="Allow passcode authentication."
				checked={provider.passcode_factor}
				onToggle={(): void => {
					let state = this.clone();
					state.passcode_factor = !state.passcode_factor;
					this.props.onChange(state);
				}}
			/>
		</div>;
	}

	okta(): JSX.Element {
		let provider = this.props.provider;

		return <div>
			<PageInput
				label="Okta Domain"
				help="Okta domain used to login to Okta such as 'pritunl.okta.com'."
				type="text"
				placeholder="Okta domain"
				value={provider.okta_domain}
				onChange={(val: string): void => {
					let state = this.clone();
					state.okta_domain = val;
					this.props.onChange(state);
				}}
			/>
			<PageInput
				label="Okta API Token"
				help="Okta API token found in Okta admin console."
				type="text"
				placeholder="Okta API token"
				value={provider.okta_token}
				onChange={(val: string): void => {
					let state = this.clone();
					state.okta_token = val;
					this.props.onChange(state);
				}}
			/>
			<PageSwitch
				label="Push authentication"
				help="Allow push authentication."
				checked={provider.push_factor}
				onToggle={(): void => {
					let state = this.clone();
					state.push_factor = !state.push_factor;
					this.props.onChange(state);
				}}
			/>
			<PageSwitch
				label="Passcode authentication"
				help="Allow passcode authentication."
				checked={provider.passcode_factor}
				onToggle={(): void => {
					let state = this.clone();
					state.passcode_factor = !state.passcode_factor;
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
			case 'duo':
				label = 'Duo';
				options = this.duo();
				break;
			case 'one_login':
				label = 'OneLogin';
				options = this.onelogin();
				break;
			case 'okta':
				label = 'Okta';
				options = this.okta();
				break;
		}

		return <div className="bp5-card" style={css.card}>
			<h6 style={css.label}>{label}</h6>
			<PageInfo
				fields={[
					{
						label: 'ID',
						value: provider.id || 'None',
					},
				]}
			/>
			<PageInput
				label="Name"
				help="Two-factor provider name."
				type="text"
				placeholder="Two-factor provider name"
				value={provider.name}
				onChange={(val: string): void => {
					let state = this.clone();
					state.name = val;
					this.props.onChange(state);
				}}
			/>
			<PageInput
				label="Label"
				help="Two-factor provider label that will be shown to users on the login page."
				type="text"
				placeholder="Two-factor provider label"
				value={provider.label}
				onChange={(val: string): void => {
					let state = this.clone();
					state.label = val;
					this.props.onChange(state);
				}}
			/>
			{options}
			<button
				className="bp5-button bp5-intent-danger"
				onClick={(): void => {
					this.props.onRemove();
				}}
			>Remove</button>
		</div>;
	}
}
