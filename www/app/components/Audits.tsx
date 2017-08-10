/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as AuditTypes from '../types/AuditTypes';
import AuditsStore from '../stores/AuditsStore';
import * as AuditActions from '../actions/AuditActions';
import Audit from './Audit';
import PageHeader from './PageHeader';

interface Props {
	userId: string;
}

interface State {
	audits: AuditTypes.AuditsRo;
	disabled: boolean;
}

const css = {
	header: {
		marginTop: '-5px',
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

export default class Audits extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			audits: AuditsStore.audits,
			disabled: false,
		};
	}

	componentDidMount(): void {
		AuditsStore.addChangeListener(this.onChange);
		if (this.props.userId) {
			AuditActions.load(this.props.userId);
		}
	}

	componentWillUnmount(): void {
		AuditsStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState({
			...this.state,
			audits: AuditsStore.audits,
		});
	}

	render(): JSX.Element {
		if (!this.props.userId) {
			return <div/>;
		}

		let audits: JSX.Element[] = [];

		this.state.audits.forEach((
				audit: AuditTypes.AuditRo): void => {
			audits.push(<Audit
				key={audit.id}
				audit={audit}
			/>);
		});

		return <div>
			<PageHeader>
				<div className="layout horizontal wrap" style={css.header}>
					<h2 style={css.heading}>User Audits</h2>
					<div className="flex"/>
				</div>
			</PageHeader>
			<div>
				{audits}
			</div>
			<div
				className="pt-non-ideal-state"
				style={css.noCerts}
				hidden={!!audits.length}
			>
				<div className="pt-non-ideal-state-visual pt-non-ideal-state-icon">
					<span className="pt-icon pt-icon-search-template"/>
				</div>
				<h4 className="pt-non-ideal-state-title">No audits</h4>
			</div>
		</div>;
	}
}
