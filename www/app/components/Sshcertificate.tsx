/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as SshcertificateTypes from '../types/SshcertificateTypes';
import * as AgentUtils from '../utils/AgentUtils';
import * as MiscUtils from '../utils/MiscUtils';
import PageInfo from './PageInfo';

interface Props {
	sshcertificate: SshcertificateTypes.SshcertificateRo;
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

export default class Sshcertificate extends React.Component<Props, {}> {
	render(): JSX.Element {
		let sshcertificate = this.props.sshcertificate;
		let agent = sshcertificate.agent || {};

		let certsInfo: string[] = [];
		for (let info of sshcertificate.certificates_info) {
			certsInfo.push(info.serial + ': ' + MiscUtils.formatDateShortTime(
				info.expires));
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
								value: sshcertificate.id || 'None',
							},
							{
								label: 'Timestamp',
								value: MiscUtils.formatDate(
									sshcertificate.timestamp) || 'Unknown',
							},
							{
								label: 'Authority IDs',
								value: sshcertificate.authority_ids,
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
								label: 'Certificate Expirations',
								value: certsInfo,
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
