/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SessionTypes from '../types/SessionTypes';
import SessionsStore from '../stores/SessionsStore';
import * as SessionActions from '../actions/SessionActions';
import NonState from './NonState';
import Session from './Session';
import PageHeader from './PageHeader';

interface Props {
	userId: string;
}

interface State {
	sessions: SessionTypes.SessionsRo;
	showEnded: boolean;
	disabled: boolean;
}

const css = {
	header: {
		marginTop: '5px',
	} as React.CSSProperties,
	heading: {
		margin: '19px 0 0 0',
	} as React.CSSProperties,
	button: {
		margin: '15px 0 -5px 0',
	} as React.CSSProperties,
};

export default class Sessions extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			sessions: SessionsStore.sessions,
			showEnded: false,
			disabled: false,
		};
	}

	componentDidMount(): void {
		SessionsStore.addChangeListener(this.onChange);
		if (this.props.userId) {
			SessionActions.load(this.props.userId);
		}
	}

	componentWillUnmount(): void {
		SessionsStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			sessions: SessionsStore.sessions,
		});
	}

	render(): JSX.Element {
		if (!this.props.userId) {
			return <div/>;
		}

		let sessions: JSX.Element[] = [];

		this.state.sessions.forEach((
				session: SessionTypes.SessionRo): void => {
			if (session.removed && !this.state.showEnded) {
				return;
			}
			sessions.push(<Session
				key={session.id}
				session={session}
			/>);
		});

		return <div>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>User Sessions</h2>
					<div className="flex"/>
					<div>
						<button
							className="bp5-button bp5-minimal"
							style={css.button}
							type="button"
							onClick={(): void => {
								this.setState({
									...this.state,
									showEnded: !this.state.showEnded,
								});
								SessionActions.showRemoved(!this.state.showEnded);
							}}
						>
							{(this.state.showEnded ? 'Hide' : 'Show') + ' ended sessions'}
						</button>
					</div>
				</div>
			</PageHeader>
			<div>
				{sessions}
			</div>
			<NonState
				hidden={!!sessions.length}
				iconClass="bp5-icon-user"
				title="No sessions"
			/>
		</div>;
	}
}
