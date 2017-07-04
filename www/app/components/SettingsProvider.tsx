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
import PageInputButton from './PageInputButton';
import PageSelect from './PageSelect';
import PageSave from './PageSave';

interface Props {
	provider: SettingsTypes.ProviderAny;
	onChange: (state: SettingsTypes.ProviderAny) => void;
	onRemove: () => void;
}

interface State {
	addRole: string;
}

const css = {
	card: {
		marginBottom: '5px',
	} as React.CSSProperties,
	role: {
		margin: '9px 5px 0 5px',
		height: '20px',
	} as React.CSSProperties,
};

export default class SettingsProvider extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			addRole: '',
		};
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

		let roles: JSX.Element[] = [];
		for (let role of provider.default_roles) {
			roles.push(
				<div
					className="pt-tag pt-tag-removable pt-intent-primary"
					style={css.role}
					key={role}
				>
					{role}
					<button
						className="pt-tag-remove"
						onMouseUp={(): void => {
							let roles = [
								...this.props.provider.default_roles,
							];

							let i = roles.indexOf(role);
							if (i === -1) {
								return;
							}

							roles.splice(i, 1);

							let state = this.clone();
							state.default_roles = roles;
							this.props.onChange(state);
						}}
					/>
				</div>,
			);
		}

		return <div className="pt-card" style={css.card}>
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
			<label className="pt-label">
				Roles
				<div>
					{roles}
				</div>
			</label>
			<PageInputButton
				buttonClass="pt-intent-success"
				label="Add"
				type="text"
				placeholder="Add default role"
				value={this.state.addRole}
				onChange={(val: string): void => {
					this.setState({
						...this.state,
						addRole: val,
					});
				}}
				onSubmit={(): void => {
					let roles = [
						...this.props.provider.default_roles,
					];

					if (!this.state.addRole) {
						return;
					}

					if (roles.indexOf(this.state.addRole) === -1) {
						roles.push(this.state.addRole);
					}

					roles.sort();

					let state = this.clone();
					state.default_roles = roles;
					this.props.onChange(state);

					this.setState({
						...this.state,
						addRole: '',
					});
				}}
			/>
			{options}
			<button
				className="pt-button pt-intent-danger"
				onClick={(): void => {
					this.props.onRemove();
				}}
			>Remove</button>
		</div>;
	}
}
