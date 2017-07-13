/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ServiceTypes from '../types/ServiceTypes';
import PageInput from './PageInput';

interface Props {
	service: ServiceTypes.ServiceRo;
}

interface State {
	changed: boolean;
	service: ServiceTypes.Service;
}

const css = {
	card: {
		padding: '10px',
		marginBottom: '5px',
	} as React.CSSProperties,
};

export default class Service extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			changed: false,
			service: null,
		};
	}

	set = (name: string, val: any): void => {
		let service: any;

		if (this.state.changed) {
			service = {
				...this.state.service,
			};
		} else {
			service = {
				...this.props.service,
			};
		}

		service[name] = val;

		this.setState({
			...this.state,
			changed: true,
			service: service,
		});
	}

	render(): JSX.Element {
		let service: ServiceTypes.Service = this.state.changed ?
			this.state.service : this.props.service;

		return <div
			className="pt-card"
			style={css.card}
		>
			<PageInput
				label="Name"
				type="text"
				placeholder="Enter name"
				value={service.name}
				onChange={(val): void => {
					this.set('name', val);
				}}
			/>
		</div>;
	}
}
