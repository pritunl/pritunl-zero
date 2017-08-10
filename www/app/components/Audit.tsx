/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as AuditTypes from '../types/AuditTypes';
import * as Constants from '../Constants';
import PageInfo from './PageInfo';

interface Props {
	audit: AuditTypes.AuditRo;
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
};

export default class Audit extends React.Component<Props, {}> {
	render(): JSX.Element {
		let audit = this.props.audit;
		let agent = audit.agent || {};

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


		let fields: string[] = [];
		for (let key in audit.fields) {
			if (!audit.fields.hasOwnProperty(key)) {
				continue;
			}
			fields.push(key + ': ' + audit.fields[key]);
		}

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
								value: audit.id || 'None',
							},
							{
								label: 'Type',
								value: audit.type,
							},
							{
								label: 'Fields',
								value: fields,
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
