/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as LogTypes from '../types/LogTypes';
import SearchInput from './SearchInput';

interface Props {
	filter: LogTypes.Filter;
	onFilter: (filter: LogTypes.Filter) => void;
}

const css = {
	filters: {
		margin: '-15px 0 5px 0',
	} as React.CSSProperties,
	input: {
		width: '200px',
		margin: '5px',
	} as React.CSSProperties,
	role: {
		width: '150px',
		margin: '5px',
	} as React.CSSProperties,
	type: {
		margin: '5px',
	} as React.CSSProperties,
	check: {
		margin: '12px 5px 8px 5px',
	} as React.CSSProperties,
};

export default class LogsFilter extends React.Component<Props, {}> {
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
				style={css.input}
				placeholder="Message"
				value={this.props.filter.message}
				onChange={(val: string): void => {
					let filter = {
						...this.props.filter,
					};

					if (val) {
						filter.message = val;
					} else {
						delete filter.message;
					}

					this.props.onFilter(filter);
				}}
			/>
			<div className="bp3-select" style={css.type}>
				<select
					value={this.props.filter.level || 'any'}
					onChange={(evt): void => {
						let filter = {
							...this.props.filter,
						};

						let val = evt.target.value;

						if (val === 'any') {
							delete filter.level;
						} else {
							filter.level = val;
						}

						this.props.onFilter(filter);
					}}
				>
					<option value="any">Any</option>
					<option value="debug">Debug</option>
					<option value="info">Info</option>
					<option value="warning">Warning</option>
					<option value="error">Error</option>
					<option value="fatal">Fatal</option>
					<option value="panic">Panic</option>
				</select>
			</div>
		</div>;
	}
}
