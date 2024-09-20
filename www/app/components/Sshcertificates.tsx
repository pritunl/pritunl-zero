/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SshcertificateTypes from '../types/SshcertificateTypes';
import SshcertificatesStore from '../stores/SshcertificatesStore';
import * as SshcertificateActions from '../actions/SshcertificateActions';
import NonState from './NonState';
import Sshcertificate from './Sshcertificate';
import PageHeader from './PageHeader';
import SshcertificatesPage from './SshcertificatesPage';

interface Props {
	userId: string;
}

interface State {
	sshcertificates: SshcertificateTypes.SshcertificatesRo;
	disabled: boolean;
}

const css = {
	header: {
		marginTop: '5px',
	} as React.CSSProperties,
	heading: {
		margin: '19px 0 0 0',
	} as React.CSSProperties,
	button: {
		margin: '15px 0 0 0',
	} as React.CSSProperties,
};

export default class Sshcertificates extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			sshcertificates: SshcertificatesStore.sshcertificates,
			disabled: false,
		};
	}

	componentDidMount(): void {
		SshcertificatesStore.addChangeListener(this.onChange);
		if (this.props.userId) {
			SshcertificateActions.load(this.props.userId);
		}
	}

	componentWillUnmount(): void {
		SshcertificatesStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			sshcertificates: SshcertificatesStore.sshcertificates,
		});
	}

	render(): JSX.Element {
		if (!this.props.userId) {
			return <div/>;
		}

		let sshcertificates: JSX.Element[] = [];

		this.state.sshcertificates.forEach((
				sshcertificate: SshcertificateTypes.SshcertificateRo): void => {
			sshcertificates.push(<Sshcertificate
				key={sshcertificate.id}
				sshcertificate={sshcertificate}
			/>);
		});

		return <div>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>User SSH Certificates</h2>
					<div className="flex"/>
				</div>
			</PageHeader>
			<div>
				{sshcertificates}
			</div>
			<NonState
				hidden={!!sshcertificates.length}
				iconClass="bp5-icon-endorsed"
				title="No SSH certificates"
			/>
			<SshcertificatesPage/>
		</div>;
	}
}
