/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as UserTypes from '../types/UserTypes';
import SearchInput from './SearchInput';
import SwitchNull from './SwitchNull';

type OnFilter = (filter: UserTypes.Filter) => void;

interface Props {
	filter: UserTypes.Filter;
	onFilter: OnFilter;
}

const css = {
	filters: {
		margin: '-15px 0 5px 0',
	} as React.CSSProperties,
	username: {
		width: '200px',
		margin: '5px 5px 5px 5px',
	} as React.CSSProperties,
	role: {
		width: '150px',
		margin: '5px 5px 5px 5px',
	} as React.CSSProperties,
	check: {
		margin: '12px 5px 8px 5px',
	} as React.CSSProperties,
};

export default class UsersFilter extends React.Component<Props, {}> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			menu: false,
		};
	}

	render(): JSX.Element {
		if (this.props.filter === null) {
			return <div/>;
		}

		return <div className="layout horizontal wrap" style={css.filters}>
			<SearchInput
				style={css.username}
				placeholder="Username"
				value={this.props.filter.username}
				onChange={(val: string): void => {
					let filter = {
						...this.props.filter,
					};

					if (val) {
						filter.username = val;
					} else {
						delete filter.username;
					}

					this.props.onFilter(filter);
				}}
			/>
			<SearchInput
				style={css.role}
				placeholder="Role"
				value={this.props.filter.role}
				onChange={(val: string): void => {
					let filter = {
						...this.props.filter,
					};

					if (val) {
						filter.role = val;
					} else {
						delete filter.role;
					}

					this.props.onFilter(filter);
				}}
			/>
			<SwitchNull
				style={css.check}
				label="Administrator"
				checked={this.props.filter.administrator}
				onToggle={(): void => {
					let filter = {
						...this.props.filter,
					};

					if (filter.administrator === undefined) {
						filter.administrator = true;
					} else if (filter.administrator === true) {
						filter.administrator = false;
					} else {
						delete filter.administrator;
					}

					this.props.onFilter(filter);
				}}
			/>
			<SwitchNull
				style={css.check}
				label="Disabled"
				checked={this.props.filter.disabled}
				onToggle={(): void => {
					let filter = {
						...this.props.filter,
					};

					if (filter.disabled === undefined) {
						filter.disabled = true;
					} else if (filter.disabled === true) {
						filter.disabled = false;
					} else {
						delete filter.disabled;
					}

					this.props.onFilter(filter);
				}}
			/>
		</div>;
	}
}
