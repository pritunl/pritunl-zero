/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Styles from '../Styles';
import * as MiscUtils from '../utils/MiscUtils';
import * as UserActions from '../actions/UserActions';
import * as UserTypes from '../types/UserTypes';

interface Props {
	userId: string;
}

interface State {
	changed: boolean;
	disabled: boolean;
	message: string,
	user: UserTypes.User;
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

export default class UserDetailed extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			changed: false,
			disabled: false,
			message: '',
			user: null,
		};
	}

	componentDidMount(): void {
		UserActions.get(this.props.userId).then((user: UserTypes.User) => {
			this.setState({
				...this.state,
				user: user,
			});
		});
	}

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
	}

	set = (name: string, val: any): void => {
		let user = {
			...this.state.user,
		} as any;

		user[name] = val;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			user: user,
		});
	}

	render(): JSX.Element {
		let user = this.state.user;

		if (!user) {
			return <div/>;
		}

		return <div style={Styles.page}>
			<div className="pt-border" style={Styles.pageHeader}>
				<h2>User Info</h2>
			</div>
			<div className="layout horizontal">
				<div className="flex">
					<label className="pt-label">
						Username
						<input
							className="pt-input"
							style={css.input}
							type="text"
							autoCapitalize="off"
							spellCheck={false}
							placeholder="Enter Elasticsearch address"
							value={user.username}
							onChange={(evt): void => {
								this.set('username', evt.target.value);
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
