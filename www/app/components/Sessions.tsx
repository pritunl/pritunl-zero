/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SessionTypes from '../types/SessionTypes';
import SessionsStore from '../stores/SessionsStore';
import * as SessionActions from '../actions/SessionActions';
import Session from './Session';
import PageHeader from './PageHeader';

interface Props {
	userId: string;
}

interface State {
	sessions: SessionTypes.SessionsRo;
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
		margin: '10px 0 0 0',
	} as React.CSSProperties,
	noCerts: {
		height: 'auto',
	} as React.CSSProperties,
};

export default class Sessions extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			sessions: SessionsStore.sessions,
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
			sessions.push(<Session
				key={session.id}
				session={session}
			/>);
		});

		let filterClass = 'pt-button pt-intent-primary pt-icon-filter ';

		return <div>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>User Sessions</h2>
					<div className="flex"/>
					<div>
						<button
							className={filterClass}
							style={css.button}
							type="button"
							onClick={(): void => {
							}}
						>
						Show Ended
						</button>
					</div>
				</div>
			</PageHeader>
			<div>
				{sessions}
			</div>
			<div
				className="pt-non-ideal-state"
				style={css.noCerts}
				hidden={!!sessions.length}
			>
				<div className="pt-non-ideal-state-visual pt-non-ideal-state-icon">
					<span className="pt-icon pt-icon-user"/>
				</div>
				<h4 className="pt-non-ideal-state-title">No sessions</h4>
			</div>
		</div>;
	}
}
