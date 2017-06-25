/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Styles from '../Styles';
import * as MiscUtils from '../utils/MiscUtils';
import * as UserTypes from '../types/UserTypes';

interface Props {
	user: UserTypes.User;
}

const css = {
	card: {
		width: '100%',
		padding: '9px',
	} as React.CSSProperties,
	tag: {
		margin: '0 5px',
	},
	select: {
		margin: '2px 0 0 0',
		paddingTop: '1px',
		minHeight: '18px',
	} as React.CSSProperties,
	name: {
		flex: '0 0 30%',
	} as React.CSSProperties,
	type: {
		flex: '0 0 30%',
	} as React.CSSProperties,
	lastActivity: {
		flex: '0 0 30%',
	} as React.CSSProperties,
	roles: {
		flex: '0 0 30%',
	} as React.CSSProperties,
	nameLink: {
		fontSize: '16px',
		margin: '0 5px 0 0',
	} as React.CSSProperties,
};

export default class User extends React.Component<Props, {}> {
	render(): JSX.Element {
		let user = this.props.user;

		return <div
			className="pt-card layout horizontal"
			style={css.card}
		>
			<div className="layout horizontal" style={css.name}>
				<label className="pt-control pt-checkbox" style={css.select}>
					<input type="checkbox"/>
					<span className="pt-control-indicator"/>
				</label>
				<div>
					<a style={css.nameLink}>
						{user.username}
					</a>
				</div>
				<div
					className="pt-tag pt-intent-danger"
					style={css.tag}
					hidden={!user.administrator}
				>
					administrator
				</div>
			</div>
			<div className="layout horizontal" style={css.type}>
				{user.type}
			</div>
			<div className="layout horizontal" style={css.lastActivity}>
				{MiscUtils.formatDate(user.last_active)}
			</div>
			<div className="layout horizontal" style={css.roles}>

			</div>
		</div>;
	}
}
