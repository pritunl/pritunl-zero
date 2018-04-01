/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ServiceTypes from '../types/ServiceTypes';
import ServicesStore from '../stores/ServicesStore';
import * as ServiceActions from '../actions/ServiceActions';
import NonState from './NonState';
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
		margin: '8px 0 0 8px',
	} as React.CSSProperties,
	buttons: {
		marginTop: '8px',
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
					<div style={css.buttons}>
						<button
							className="pt-button pt-intent-success pt-icon-add"
							style={css.button}
							disabled={this.state.disabled}
							type="button"
							onClick={(): void => {
								this.setState({
									...this.state,
									disabled: true,
								});
								ServiceActions.create({
									id: null,
									share_session: true,
									websockets: true,
								}).then((): void => {
									this.setState({
										...this.state,
										disabled: false,
									});
								}).catch((): void => {
									this.setState({
										...this.state,
										disabled: false,
									});
								});
							}}
						>New</button>
					</div>
				</div>
			</PageHeader>
			<div>
				{servicesDom}
			</div>
			<NonState
				hidden={!!servicesDom.length}
				iconClass="pt-icon-cloud"
				title="No services"
				description="Add a new service to get started."
			/>
		</Page>;
	}
}
