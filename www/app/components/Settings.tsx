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
import PageSelectButton from './PageSelectButton';
import PageSave from './PageSave';
import SettingsProvider from './SettingsProvider';

interface State {
	changed: boolean;
	disabled: boolean;
	message: string;
	provider: string;
	settings: SettingsTypes.Settings;
}

const css = {
	providers: {
		paddingBottom: '6px',
		marginBottom: '5px',
		borderBottomStyle: 'solid',
	} as React.CSSProperties,
	providersLabel: {
		margin: 0,
	} as React.CSSProperties,
};

export default class Settings extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			changed: false,
			disabled: false,
			message: '',
			provider: 'google',
			settings: SettingsStore.settingsM,
		};
	}

	componentDidMount(): void {
		SettingsStore.addChangeListener(this.onChange);
		SettingsActions.sync();
	}

	componentWillUnmount(): void {
		SettingsStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			changed: false,
			settings: SettingsStore.settingsM,
		});
	}

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		SettingsActions.commit(this.state.settings).then((): void => {
			this.setState({
				...this.state,
				message: 'Your changes have been saved',
				changed: false,
				disabled: false,
			});
		}).catch((): void => {
			this.setState({
				...this.state,
				message: '',
				disabled: false,
			});
		});
	}

	set = (name: string, val: any): void => {
		let settings: any = {
			...this.state.settings,
		};

		settings[name] = val;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			settings: settings,
		});
	}

	render(): JSX.Element {
		let settings = this.state.settings;

		if (!settings) {
			return <div/>;
		}

		let providers: JSX.Element[] = [];
		for (let i = 0; i < settings.auth_providers.length; i++) {
			providers.push(<SettingsProvider
				key={i}
				provider={settings.auth_providers[i]}
				onChange={(state): void => {
					let providers = [
						...this.state.settings.auth_providers,
					];
					providers[i] = state;
					this.set('auth_providers', providers);
				}}
				onRemove={(): void => {
					let providers = [
						...this.state.settings.auth_providers,
					];
					providers.splice(i, 1);
					this.set('auth_providers', providers);
				}}
			/>);
		}

		return <Page>
			<PageHeader label="Settings"/>
			<PageSplit>
				<PagePanel>
					<div className="pt-border" style={css.providers}>
						<h5 style={css.providersLabel}>Authentication Providers</h5>
					</div>
					{providers}
					<PageSelectButton
						label="Add Provider"
						value={this.state.provider}
						buttonClass="pt-intent-success"
						onChange={(val: string): void => {
							this.setState({
								...this.state,
								provider: val,
							});
						}}
						onSubmit={(): void => {
							let authProviders: SettingsTypes.Providers = [
								...settings.auth_providers,
								{
									type: this.state.provider,
									label: '',
									default_roles: [],
									auto_create: true,
								},
							];
							this.set('auth_providers', authProviders);
						}}
					>
						<option value="google">Google</option>
						<option value="onelogin">OneLogin</option>
						<option value="okta">Okta</option>
					</PageSelectButton>
				</PagePanel>
				<PagePanel>
					<PageInput
						label="Elasticsearch Address"
						type="text"
						placeholder="Enter Elasticsearch address"
						value={this.state.settings.elastic_address}
						onChange={(val): void => {
							this.set('elastic_address', val);
						}}
					/>
				</PagePanel>
			</PageSplit>
			<PageSave
				message={this.state.message}
				changed={this.state.changed}
				disabled={this.state.disabled}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						message: 'Your changes have been discarded',
						settings: SettingsStore.settingsM,
					});
				}}
				onSave={this.onSave}
			/>
		</Page>;
	}
}
