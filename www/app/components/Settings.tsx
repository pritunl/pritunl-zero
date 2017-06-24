/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Styles from '../Styles';
import * as SettingsTypes from '../types/SettingsTypes';
import SettingsStore from '../stores/SettingsStore';
import * as SettingsActions from '../actions/SettingsActions';

interface State {
	changed: boolean;
	disabled: boolean;
	message: string,
	settings: SettingsTypes.Settings;
}

const css = {
	input: {
		width: '100%',
		maxWidth: '310px',
	} as React.CSSProperties,
	button: {
		marginLeft: '10px',
	} as React.CSSProperties,
};

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
		return <div style={Styles.page}>
			<div className="pt-border" style={Styles.pageHeader}>
				<h2>Settings</h2>
			</div>
			<div className="layout horizontal">
				<div className="flex">
					<label className="pt-label">
						Elasticsearch Address
						<input
							className="pt-input"
							style={css.input}
							type="text"
							autoCapitalize="off"
							spellCheck={false}
							placeholder="Enter Elasticsearch address"
							value={this.state.settings.elastic_address}
							onChange={(evt): void => {
								this.set('elastic_address', evt.target.value);
							}}
						/>
					</label>
				</div>
				<div className="flex">
				</div>
			</div>
			<div className="layout horizontal">
				<div className="flex"/>
				<div>
					<span hidden={!this.state.message}>
						{this.state.message}
					</span>
					<button
						className="pt-button pt-icon-cross"
						style={css.button}
						type="button"
						disabled={!this.state.changed || this.state.disabled}
						onClick={(): void => {
							this.setState({
								...this.state,
								changed: false,
								message: 'Your changes have been discarded',
								settings: SettingsStore.settings,
							})
						}}
					>
						Cancel
					</button>
					<button
						className="pt-button pt-intent-success pt-icon-tick"
						style={css.button}
						type="button"
						disabled={!this.state.changed || this.state.disabled}
						onClick={this.onSave}
					>
						Save
					</button>
				</div>
			</div>
		</div>;
	}
}
