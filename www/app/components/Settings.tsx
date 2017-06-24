/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SettingsTypes from '../types/SettingsTypes';
import SettingsStore from '../stores/SettingsStore';
import * as SettingsActions from '../actions/SettingsActions';

interface State {
	settings: SettingsTypes.Settings;
}

function getState(): State {
	return {
		settings: SettingsStore.settings,
	};
}

export default class Settings extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = getState();
	}

	componentDidMount(): void {
		SettingsActions.sync();
		SettingsStore.addChangeListener(this.onChange);
	}

	componentWillUnmount(): void {
		SettingsStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState(getState());
	}

	render(): JSX.Element {
		return <div>
			{this.state.settings.elastic_address}
		</div>;
	}
}
