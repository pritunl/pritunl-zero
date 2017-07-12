/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ServiceTypes from '../types/ServiceTypes';
import ServicesStore from '../stores/ServicesStore';
import * as ServiceActions from '../actions/ServiceActions';
import Service from './Service';
import Page from './Page';
import PageHeader from './PageHeader';

interface State {
	services: ServiceTypes.ServicesRo;
	disabled: boolean;
}

const css = {
	header: {
		marginTop: '-19px',
	} as React.CSSProperties,
	heading: {
		margin: '19px 0 0 0',
	} as React.CSSProperties,
	button: {
		margin: '10px 0 0 10px',
	} as React.CSSProperties,
	buttonFirst: {
		margin: '10px 0 0 0',
	} as React.CSSProperties,
};

export default class Services extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			services: ServicesStore.services,
			disabled: false,
		};
	}

	componentDidMount(): void {
		ServicesStore.addChangeListener(this.onChange);
		ServiceActions.sync();
	}

	componentWillUnmount(): void {
		ServicesStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			services: ServicesStore.services,
		});
	}

	render(): JSX.Element {
		let servicesDom: JSX.Element[] = [];

		this.state.services.forEach((service: ServiceTypes.ServiceRo): void => {
			servicesDom.push(<Service
				key={service.id}
				service={service}
			/>);
		});

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Services</h2>
					<div className="flex"/>
					<div>
						<button
							className="pt-button pt-intent-success pt-icon-add"
							style={css.button}
							type="button"
							onClick={(): void => {
								ServiceActions.create(null);
							}}
						>New</button>
					</div>
				</div>
			</PageHeader>
			<div>
				{servicesDom}
			</div>
			<div
				className="pt-non-ideal-state"
				hidden={!!servicesDom.length}
			>
				<div className="pt-non-ideal-state-visual pt-non-ideal-state-icon">
					<span className="pt-icon pt-icon-folder-open"/>
				</div>
				<h4 className="pt-non-ideal-state-title">No services</h4>
				<div className="pt-non-ideal-state-description">
					Add a new service to get started.
				</div>
			</div>
		</Page>;
	}
}
