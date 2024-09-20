/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as AlertTypes from '../types/AlertTypes';
import * as AuthorityTypes from '../types/AuthorityTypes';
import AlertDetailed from './AlertDetailed';

interface Props {
	alert: AlertTypes.AlertRo;
	authorities: AuthorityTypes.AuthoritiesRo;
	selected: boolean;
	onSelect: (shift: boolean) => void;
	open: boolean;
	onOpen: () => void;
}

const css = {
	card: {
		display: 'table-row',
		width: '100%',
		padding: 0,
		boxShadow: 'none',
		cursor: 'pointer',
	} as React.CSSProperties,
	cardOpen: {
		display: 'table-row',
		width: '100%',
		padding: 0,
		boxShadow: 'none',
		position: 'relative',
	} as React.CSSProperties,
	select: {
		margin: '2px 0 0 0',
		paddingTop: '3px',
		minHeight: '18px',
	} as React.CSSProperties,
	name: {
		verticalAlign: 'top',
		display: 'table-cell',
		padding: '8px',
	} as React.CSSProperties,
	nameSpan: {
		margin: '1px 5px 0 0',
	} as React.CSSProperties,
	item: {
		verticalAlign: 'top',
		display: 'table-cell',
		padding: '9px',
		whiteSpace: 'nowrap',
	} as React.CSSProperties,
	roles: {
		verticalAlign: 'top',
		display: 'table-cell',
		padding: '0 8px 8px 8px',
	} as React.CSSProperties,
	icon: {
		marginRight: '3px',
	} as React.CSSProperties,
	tag: {
		margin: '8px 5px 0 5px',
		height: '20px',
	} as React.CSSProperties,
	bars: {
		verticalAlign: 'top',
		display: 'table-cell',
		padding: '8px',
		width: '30px',
	} as React.CSSProperties,
	bar: {
		height: '6px',
		marginBottom: '1px',
	} as React.CSSProperties,
	barLast: {
		height: '6px',
	} as React.CSSProperties,
};

export default class Alert extends React.Component<Props, {}> {
	render(): JSX.Element {
		let alert = this.props.alert;

		if (this.props.open) {
			return <div
				className="bp5-card bp5-row"
				style={css.cardOpen}
			>
				<AlertDetailed
					alert={this.props.alert}
					authorities={this.props.authorities}
					selected={this.props.selected}
					onSelect={this.props.onSelect}
					onClose={(): void => {
						this.props.onOpen();
					}}
				/>
			</div>;
		}

		let cardStyle = {
			...css.card,
		};

		let roles: JSX.Element[] = [];

		for (let role of alert.roles) {
			roles.push(
				<div
					className="bp5-tag bp5-intent-primary"
					style={css.tag}
					key={role}
				>
					{role}
				</div>,
			);
		}

		return <div
			className="bp5-card bp5-row"
			style={cardStyle}
			onClick={(evt): void => {
				let target = evt.target as HTMLElement;

				if (target.className.indexOf('open-ignore') !== -1) {
					return;
				}

				this.props.onOpen();
			}}
		>
			<div className="bp5-cell" style={css.name}>
				<div className="layout horizontal">
					<label
						className="bp5-control bp5-checkbox open-ignore"
						style={css.select}
					>
						<input
							type="checkbox"
							className="open-ignore"
							checked={this.props.selected}
							onChange={(evt): void => {
							}}
							onClick={(evt): void => {
								this.props.onSelect(evt.shiftKey);
							}}
						/>
						<span className="bp5-control-indicator open-ignore"/>
					</label>
					<div style={css.nameSpan}>
						{alert.name}
					</div>
				</div>
			</div>
			<div className="flex bp5-cell" style={css.roles}>
				{roles}
			</div>
		</div>;
	}
}
