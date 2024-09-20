/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';
import * as MiscUtils from '../utils/MiscUtils';
import * as LogTypes from '../types/LogTypes';

interface State {
	stack: boolean;
}

interface Props {
	log: LogTypes.LogRo;
}

const css = {
	card: {
		display: 'table-row',
		width: '100%',
		padding: 0,
		boxShadow: 'none',
	} as React.CSSProperties,
	timestamp: {
		verticalAlign: 'top',
		display: 'table-cell',
		padding: '6px',
	} as React.CSSProperties,
	level: {
		verticalAlign: 'top',
		display: 'table-cell',
		padding: '6px',
	} as React.CSSProperties,
	message: {
		verticalAlign: 'top',
		display: 'table-cell',
		padding: '6px',
	} as React.CSSProperties,
	fields: {
		verticalAlign: 'top',
		display: 'table-cell',
		padding: '6px',
	} as React.CSSProperties,
	buttons: {
		verticalAlign: 'top',
		display: 'table-cell',
		padding: '0',
		width: '30px',
	} as React.CSSProperties,
	key: {
		fontWeight: 'bold',
	} as React.CSSProperties,
	value: {
	} as React.CSSProperties,
	dialog: {
		height: '500px',
		width: '90%',
		maxWidth: '700px',
	} as React.CSSProperties,
	dialogBody: {
		height: '100%',
	} as React.CSSProperties,
	textarea: {
		resize: 'none',
		fontSize: '12px',
		fontFamily: '"Lucida Console", Monaco, monospace',
		marginBottom: 0,
	} as React.CSSProperties,
};

export default class Log extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			stack: false,
		};
	}

	render(): JSX.Element {
		let log = this.props.log;

		let className = 'bp5-cell ';
		switch (log.level) {
			case 'debug':
				className += 'bp5-text-intent-success';
				break;
			case 'info':
				className += 'bp5-text-intent-primary';
				break;
			case 'warning':
				className += 'bp5-text-intent-warning';
				break;
			case 'error':
				className += 'bp5-text-intent-danger';
				break;
			case 'fatal':
				className += 'bp5-text-intent-danger';
				break;
			case 'panic':
				className += 'bp5-text-intent-danger';
				break;
		}

		let fields: JSX.Element[] = [];
		for (let key in log.fields) {
			if (!log.fields.hasOwnProperty(key)) {
				continue;
			}

			let val = log.fields[key];

			fields.push(
				<div key={key}>
					<span style={css.key}>{key}: </span>
					<span style={css.value}>
						{JSON.stringify(val)}
					</span>
				</div>,
			);
		}

		return <div
			className="bp5-card bp5-row"
			style={css.card}
		>
			<div className={className} style={css.timestamp}>
				{MiscUtils.formatDateShortTime(log.timestamp) || 'Unknown'}
			</div>
			<div className={className} style={css.level}>
				{log.level}
			</div>
			<div className={className} style={css.message}>
				{log.message}
			</div>
			<div className="bp5-cell" style={css.fields}>
				{fields}
			</div>
			<div className="bp5-cell" style={css.buttons}>
				<button
					className="bp5-button bp5-minimal bp5-icon-document-open"
					hidden={!log.stack}
					onClick={(): void => {
						this.setState({
							...this.state,
							stack: true,
						});
					}}
				/>
			</div>
			<Blueprint.Dialog
				title="Stack Trace"
				style={css.dialog}
				isOpen={this.state.stack}
				usePortal={true}
				portalContainer={document.body}
				onClose={(): void => {
					this.setState({
						...this.state,
						stack: false,
					});
				}}
			>
				<textarea
					className="bp5-dialog-body bp5-input"
					style={css.textarea}
					autoCapitalize="off"
					spellCheck={false}
					readOnly={true}
					value={log.stack || ''}
				/>
			</Blueprint.Dialog>
		</div>;
	}
}
