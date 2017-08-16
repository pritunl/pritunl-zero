/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as CertificateTypes from '../types/CertificateTypes';
import CertificatesStore from '../stores/CertificatesStore';
import * as CertificateActions from '../actions/CertificateActions';
import * as Constants from '../Constants';
import Certificate from './Certificate';
import Page from './Page';
import PageHeader from './PageHeader';

interface State {
	certificates: CertificateTypes.CertificatesRo;
	disabled: boolean;
	initialized: boolean;
}

const css = {
	header: {
		marginTop: '-19px',
	} as React.CSSProperties,
	heading: {
		margin: '19px 0 0 0',
	} as React.CSSProperties,
	button: {
		margin: '10px 0 0 0',
	} as React.CSSProperties,
	noCerts: {
		height: 'auto',
	} as React.CSSProperties,
};

export default class Certificates extends React.Component<{}, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			certificates: CertificatesStore.certificates,
			disabled: false,
			initialized: false,
		};
	}

	componentDidMount(): void {
		CertificatesStore.addChangeListener(this.onChange);
		CertificateActions.sync();
		setTimeout((): void => {
			this.setState({
				...this.state,
				initialized: true,
			});
		}, Constants.loadDelay);
	}

	componentWillUnmount(): void {
		CertificatesStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			certificates: CertificatesStore.certificates,
		});
	}

	render(): JSX.Element {
		let certsDom: JSX.Element[] = [];

		this.state.certificates.forEach((
				cert: CertificateTypes.CertificateRo): void => {
			certsDom.push(<Certificate
				key={cert.id}
				certificate={cert}
			/>);
		});

		return <Page>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>Certificates</h2>
					<div className="flex"/>
					<div>
						<button
							className="pt-button pt-intent-success pt-icon-add"
							style={css.button}
							type="button"
							onClick={(): void => {
								CertificateActions.create(null);
							}}
						>New</button>
					</div>
				</div>
			</PageHeader>
			<div>
				{certsDom}
			</div>
			<div
				className="pt-non-ideal-state"
				style={css.noCerts}
				hidden={!!certsDom.length || !this.state.initialized}
			>
				<div className="pt-non-ideal-state-visual pt-non-ideal-state-icon">
					<span className="pt-icon pt-icon-endorsed"/>
				</div>
				<h4 className="pt-non-ideal-state-title">No certificates</h4>
				<div className="pt-non-ideal-state-description">
					Add a new certificate to get started.
				</div>
			</div>
		</Page>;
	}
}
