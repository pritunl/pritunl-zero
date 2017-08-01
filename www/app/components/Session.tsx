/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SessionTypes from '../types/SessionTypes';
import * as MiscUtils from '../utils/MiscUtils';
import * as Constants from '../Constants';
import PageInfo from './PageInfo';

interface Props {
	session: SessionTypes.SessionRo;
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
		minWidth: '250px',
	} as React.CSSProperties,
};

export default class Session extends React.Component<Props, {}> {
	render(): JSX.Element {
		let session = this.props.session;
		let agent = session.agent || {};

		let continent = agent.continent && agent.continent_code ?
			agent.continent + ' (' + agent.continent_code + ')' :
			agent.continent || agent.continent_code || 'Unknown';

		let location = (agent.city ? agent.city + ', ' : '') +
			(agent.region || 'Unknown') +
			(agent.region_code ? ' (' + agent.region_code + ')' : '');
		let country = (agent.country || 'Unknown') +
			(agent.country_code ? ' (' + agent.country_code + ')' : '');

		let coordinates = agent.latitude && agent.longitude ?
			agent.latitude + ', ' + agent.longitude : 'Unknown';

		return <div
			className="pt-card"
			style={css.card}
		>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<PageInfo
						style={css.info}
						fields={[
							{
								label: 'ID',
								value: session.id || 'None',
							},
							{
								label: 'Timestamp',
								value: MiscUtils.formatDate(session.timestamp) || 'Unknown',
							},
							{
								label: 'Operating System',
								value: Constants.operatingSystems[agent.operating_system] ||
									'Unknown',
							},
							{
								label: 'Browser',
								value: Constants.browsers[agent.browser] || 'Unknown',
							},
						]}
					/>
				</div>
				<div style={css.group}>
					<PageInfo
						style={css.info}
						fields={[
							{
								label: 'ISP',
								value: agent.isp || 'Unknown',
							},
							{
								label: 'Location',
								value: [location, country, continent],
							},
							{
								label: 'Coordinates',
								value: coordinates,
							},
						]}
					/>
				</div>
			</div>
		</div>;
	}
}
