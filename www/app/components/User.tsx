/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ReactRouter from 'react-router-dom';
import * as MiscUtils from '../utils/MiscUtils';
import * as UserTypes from '../types/UserTypes';

interface Props {
	user: UserTypes.UserRo;
	selected: boolean;
	onSelect: (shift: boolean) => void;
}

const css = {
	card: {
		display: 'table-row',
		width: '100%',
		padding: 0,
		boxShadow: 'none',
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
					className="bp3-tag bp3-intent-primary"
					style={css.tag}
					key={role}
				>
					{role}
				</div>,
			);
		}

		let cardStyle = {
			...css.card,
		};
		if (user.disabled) {
			cardStyle.opacity = 0.6;
		}

		let userType: string;
		switch (user.type) {
			case 'local':
				userType = 'Local';
				break;
			case 'google':
				userType = 'Google';
				break;
			case 'onelogin':
				userType = 'OneLogin';
				break;
			case 'okta':
				userType = 'Okta';
				break;
			case 'authzero':
				userType = 'Auth0';
				break;
			case 'azure':
				userType = 'Azure';
				break;
			case 'jumpcloud':
				userType = 'JumpCloud';
				break;
			case 'api':
				userType = 'API';
				break;
			default:
				userType = user.type;
		}

		return <div
			className="bp3-card bp3-row"
			style={cardStyle}
		>
			<div className="bp3-cell" style={css.name}>
				<div className="layout horizontal">
					<label className="bp3-control bp3-checkbox" style={css.select}>
						<input
							type="checkbox"
							checked={this.props.selected}
							onChange={(evt): void => {
							}}
							onClick={(evt): void => {
								this.props.onSelect(evt.shiftKey);
							}}
						/>
						<span className="bp3-control-indicator"/>
					</label>
					<ReactRouter.Link to={'/user/' + user.id} style={css.nameLink}>
						{user.username}
					</ReactRouter.Link>
				</div>
			</div>
			<div className="bp3-cell" style={css.type}>
				{userType}
			</div>
			<div className="bp3-cell" style={css.lastActivity}>
				{MiscUtils.formatDateShortTime(user.last_active) || 'Inactive'}
			</div>
			<div className="flex bp3-cell" style={css.roles}>
				<span
					className="bp3-tag bp3-intent-danger"
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
