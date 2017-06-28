/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ReactRouter from 'react-router-dom';
import * as MiscUtils from '../utils/MiscUtils';
import * as UserTypes from '../types/UserTypes';

type OnSelect = () => void;

interface Props {
	user: UserTypes.User;
	selected: boolean;
	onSelect: OnSelect;
}

const css = {
	card: {
		display: 'table-row',
		width: '100%',
		padding: 0,
	} as React.CSSProperties,
	select: {
		margin: '2px 0 0 0',
		paddingTop: '1px',
		minHeight: '18px',
	} as React.CSSProperties,
	name: {
		verticalAlign: 'top',
		display: 'table-cell',
		padding: '8px',
	} as React.CSSProperties,
	type: {
		verticalAlign: 'top',
		display: 'table-cell',
		padding: '9px',
	} as React.CSSProperties,
	lastActivity: {
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
	tag: {
		margin: '8px 5px 0 5px',
		height: '20px',
	} as React.CSSProperties,
	nameLink: {
		margin: '0 5px 0 0',
	} as React.CSSProperties,
};

export default class User extends React.Component<Props, {}> {
	render(): JSX.Element {
		let user = this.props.user;
		let roles: JSX.Element[] = [];

		for (let role of user.roles) {
			roles.push(
				<div
					className="pt-tag pt-intent-primary"
					style={css.tag}
					key={role}
				>
					{role}
				</div>
			);
		}

		return <div
			className="pt-card"
			style={css.card}
		>
			<div style={css.name}>
				<div className="layout horizontal">
					<label className="pt-control pt-checkbox" style={css.select}>
						<input
							type="checkbox"
							checked={this.props.selected}
							onChange={this.props.onSelect}
						/>
						<span className="pt-control-indicator"/>
					</label>
					<ReactRouter.Link to={'/user/' + user.id} style={css.nameLink}>
						{user.username}
					</ReactRouter.Link>
				</div>
			</div>
			<div className="layout horizontal" style={css.type}>
				{user.type}
			</div>
			<div className="layout horizontal" style={css.lastActivity}>
				{MiscUtils.formatDate(user.last_active) || 'Inactive'}
			</div>
			<div className="flex" style={css.roles}>
				<span
					className="pt-tag pt-intent-danger"
					style={css.tag}
					hidden={!user.administrator}
				>
					admin
				</span>
				{roles}
			</div>
		</div>;
	}
}
