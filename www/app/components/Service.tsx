/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ServiceTypes from '../types/ServiceTypes';

interface Props {
	service: ServiceTypes.ServiceRo;
}

const css = {
	card: {
		padding: '10px',
		marginBottom: '5px',
	} as React.CSSProperties,
};

export default class Service extends React.Component<Props, {}> {
	render(): JSX.Element {
		let service = this.props.service;

		return <div
			className="pt-card"
			style={css.card}
		>
			{service.id}
		</div>;
	}
}
