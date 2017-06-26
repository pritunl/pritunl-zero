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
import PageSave from './PageSave';

interface State {
	changed: boolean;
	disabled: boolean;
	message: string,
	settings: SettingsTypes.Settings;
}

export default class Settings extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			changed: false,
			disabled: false,
			message: '',
			settings: SettingsStore.settings,
		};
	}

	componentDidMount(): void {
		SettingsActions.sync();
		SettingsStore.addChangeListener(this.onChange);
	}

	componentWillUnmount(): void {
		SettingsStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			changed: false,
			settings: SettingsStore.settings,
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
			})
		});
	}

	set = (name: string, val: any): void => {
		let settings = {
			...this.state.settings,
		} as any;

		settings[name] = val;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			settings: settings,
		});
	}

	render(): JSX.Element {
		return <Page>
			<PageHeader title="User Info"/>
			<PageSplit>
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
				<PagePanel>
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
						settings: SettingsStore.settings,
					});
				}}
				onSave={this.onSave}
			/>
		</Page>;
	}
}
