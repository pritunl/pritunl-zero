/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as LogTypes from '../types/LogTypes';
import LogsStore from '../stores/LogsStore';
import * as LogActions from '../actions/LogActions';
import Log from './Log';
import LogsFilter from './LogsFilter';
import Page from './Page';
import PageHeader from './PageHeader';
import LogsPage from './LogsPage';

interface State {
	logs: LogTypes.LogsRo;
	filter: LogTypes.Filter;
}

const css = {
	logs: {
		width: '100%',
		marginTop: '-3px',
		display: 'table',
		tableLayout: 'fixed',
		borderSpacing: '0 3px',
	} as React.CSSProperties,
	logsBox: {
		width: '100%',
		overflowY: 'auto',
	} as React.CSSProperties,
	header: {
		marginTop: '-19px',
	} as React.CSSProperties,
	heading: {
		margin: '19px 0 0 0',
	} as React.CSSProperties,
	button: {
		margin: '8px 0 0 8px',
	} as React.CSSProperties,
	buttons: {
		marginTop: '8px',
	} as React.CSSProperties,
};

export default class Logs extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			logs: LogsStore.logs,
			filter: LogsStore.filter,
		};
	}

	componentDidMount(): void {
		LogsStore.addChangeListener(this.onChange);
		LogActions.sync();
	}

	componentWillUnmount(): void {
		LogsStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			logs: LogsStore.logs,
			filter: LogsStore.filter,
		});
	}

	render(): JSX.Element {
		let logsDom: JSX.Element[] = [];

		this.state.logs.forEach((log: LogTypes.LogRo): void => {
			logsDom.push(<Log
				key={log.id}
				log={log}
			/>);
		});

		let filterClass = 'bp5-button bp5-intent-primary bp5-icon-filter ';
		if (this.state.filter) {
			filterClass += 'bp5-active';
		}

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Logs</h2>
					<div className="flex"/>
					<div style={css.buttons}>
						<button
							className={filterClass}
							style={css.button}
							type="button"
							onClick={(): void => {
								if (this.state.filter === null) {
									LogActions.filter({});
								} else {
									LogActions.filter(null);
								}
							}}
						>
							Filters
						</button>
					</div>
				</div>
			</PageHeader>
			<LogsFilter
				filter={this.state.filter}
				onFilter={(filter): void => {
					LogActions.filter(filter);
				}}
			/>
			<div style={css.logsBox}>
				<div style={css.logs}>
					{logsDom}
				</div>
			</div>
			<LogsPage/>
		</Page>;
	}
}
