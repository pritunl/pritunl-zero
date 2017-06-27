/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';
import * as UserTypes from '../types/UserTypes';

type OnFilter = (filter: UserTypes.Filter) => void;

interface Props {
	filter: UserTypes.Filter;
	onFilter: OnFilter;
}

interface State {
}

const css = {
	label: {
		display: 'inline-block',
	} as React.CSSProperties,
};

export default class UsersFilter extends React.Component<Props, State> {
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

		return <div className="layout horizontal">
			<label className="pt-control pt-switch" style={css.label}>
				<input
					type="checkbox"
					checked={!!this.props.filter['administrator']}
					onChange={(): void => {
						let filter = {
							...this.props.filter,
						};

						if (filter['administrator']) {
							delete filter['administrator'];
						} else {
							filter['administrator'] = true;
						}

						this.props.onFilter(filter);
					}}
				/>
				<span className="pt-control-indicator"/>
				Administrator
			</label>
		</div>;
	}
}
