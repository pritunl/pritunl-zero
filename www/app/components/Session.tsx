/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SessionTypes from '../types/SessionTypes';
import * as MiscUtils from '../utils/MiscUtils';
import * as AgentUtils from '../utils/AgentUtils';
import * as Constants from '../Constants';
import * as SessionActions from '../actions/SessionActions';
import PageInfo from './PageInfo';

interface Props {
	session: SessionTypes.SessionRo;
}

interface State {
	disabled: boolean;
}

const css = {
	card: {
		position: 'relative',
		padding: '10px',
		marginBottom: '5px',
	} as React.CSSProperties,
	info: {
		marginBottom: '-5px',
	} as React.CSSProperties,
	group: {
		flex: 1,
		minWidth: '290px',
	} as React.CSSProperties,
	remove: {
		position: 'absolute',
		top: '5px',
		right: '5px',
	} as React.CSSProperties,
};

export default class Session extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
		};
	}

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		SessionActions.remove(this.props.session.id).then((): void => {
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
	}

	render(): JSX.Element {
		let session = this.props.session;
		let agent = session.agent || {};

		return <div
			className="pt-card"
			style={css.card}
		>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<div style={css.remove}>
						<button
							className="pt-button pt-minimal pt-intent-danger pt-icon-cross"
							type="button"
							disabled={this.state.disabled}
							onClick={this.onDelete}
						/>
					</div>
					<PageInfo
						style={css.info}
						fields={[
							{
								label: 'ID',
								value: session.id || 'None',
							},
							{
								label: 'Created',
								value: MiscUtils.formatDate(session.timestamp) || 'Unknown',
							},
							{
								label: 'Last Active',
								value: MiscUtils.formatDate(session.last_active) || 'Unknown',
							},
						]}
					/>
				</div>
				<div style={css.group}>
					<PageInfo
						style={css.info}
						fields={[
							{
								label: 'Operating System',
								value: Constants.operatingSystems[agent.operating_system] ||
								'Unknown',
							},
							{
								label: 'Browser',
								value: Constants.browsers[agent.browser] || 'Unknown',
							},
							{
								label: 'ISP',
								value: agent.isp || 'Unknown',
							},
						]}
					/>
				</div>
				<div style={css.group}>
					<PageInfo
						style={css.info}
						fields={[
							{
								label: 'Location',
								value: [
									AgentUtils.formatLocation(agent),
									AgentUtils.formatCountry(agent),
									AgentUtils.formatContinent(agent),
								],
							},
							{
								label: 'Coordinates',
								value: AgentUtils.formatCoordinates(agent),
							},
						]}
					/>
				</div>
			</div>
		</div>;
	}
}
