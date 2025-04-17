/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as MiscUtils from '../utils/MiscUtils';
import * as NodeTypes from '../types/NodeTypes';
import * as ServiceTypes from '../types/ServiceTypes';
import * as AuthorityTypes from '../types/AuthorityTypes';
import * as CertificateTypes from '../types/CertificateTypes';
import NodeDetailed from './NodeDetailed';

interface Props {
	node: NodeTypes.NodeRo;
	services: ServiceTypes.ServicesRo;
	authorities: AuthorityTypes.AuthoritiesRo;
	certificates: CertificateTypes.CertificatesRo;
	selected: boolean;
	open: boolean;
	onSelect: (shift: boolean) => void;
	onOpen: () => void;
}

const css = {
	card: {
		display: 'table-row',
		width: '100%',
		padding: 0,
		boxShadow: 'none',
		cursor: 'pointer',
	} as React.CSSProperties,
	cardOpen: {
		display: 'table-row',
		width: '100%',
		padding: 0,
		boxShadow: 'none',
		position: 'relative',
	} as React.CSSProperties,
	select: {
		margin: '2px 0 0 0',
		paddingTop: '3px',
		minHeight: '18px',
	} as React.CSSProperties,
	name: {
		verticalAlign: 'top',
		display: 'table-cell',
		padding: '8px',
	} as React.CSSProperties,
	nameSpan: {
		margin: '1px 5px 0 0',
	} as React.CSSProperties,
	item: {
		verticalAlign: 'top',
		display: 'table-cell',
		padding: '9px',
		whiteSpace: 'nowrap',
	} as React.CSSProperties,
	bars: {
		verticalAlign: 'top',
		display: 'table-cell',
		padding: '8px',
		width: '85px',
	} as React.CSSProperties,
	bar: {
		height: '6px',
		marginBottom: '1px',
	} as React.CSSProperties,
	barLast: {
		height: '6px',
	} as React.CSSProperties,
};

export default class Node extends React.Component<Props, {}> {
	render(): JSX.Element {
		let node = this.props.node;

		if (this.props.open) {
			return <div
				className="bp5-card bp5-row"
				style={css.cardOpen}
			>
				<NodeDetailed
					node={this.props.node}
					services={this.props.services}
					authorities={this.props.authorities}
					certificates={this.props.certificates}
					selected={this.props.selected}
					onSelect={this.props.onSelect}
					onClose={(): void => {
						this.props.onOpen();
					}}
				/>
			</div>;
		}

		let active = node.requests_min !== 0 || node.memory !== 0 ||
			node.load1 !== 0 || node.load5 !== 0 || node.load15 !== 0;

		let cardStyle = {
			...css.card,
		};
		if (!active) {
			cardStyle.opacity = 0.6;
		}

		let memoryStyle: React.CSSProperties = {
			width: (node.memory || 0) + '%',
		};
		let load1Style: React.CSSProperties = {
			width: (node.load1 || 0) + '%',
		};
		let load5Style: React.CSSProperties = {
			width: (node.load5 || 0) + '%',
		};

		return <div
			className="bp5-card bp5-row"
			style={cardStyle}
			onClick={(evt): void => {
				let target = evt.target as HTMLElement;

				if (target.className.indexOf('open-ignore') !== -1) {
					return;
				}

				this.props.onOpen();
			}}
		>
			<div className="bp5-cell" style={css.name}>
				<div className="layout horizontal">
					<label
						className="bp5-control bp5-checkbox open-ignore"
						style={css.select}
					>
						<input
							type="checkbox"
							className="open-ignore"
							checked={this.props.selected}
							onChange={(evt): void => {
							}}
							onClick={(evt): void => {
								this.props.onSelect(evt.shiftKey);
							}}
						/>
						<span className="bp5-control-indicator open-ignore"/>
					</label>
					<div style={css.nameSpan}>
						{node.name}
					</div>
				</div>
			</div>
			<div className="bp5-cell" style={css.item}>
				{node.requests_min + '/min'}
			</div>
			<div className="bp5-cell" style={css.item}>
				{MiscUtils.formatDateShortTime(node.timestamp) || 'Inactive'}
			</div>
			<div className="bp5-cell" style={css.bars}>
				<div
					className="bp5-progress-bar bp5-no-stripes bp5-intent-primary"
					style={css.bar}
				>
					<div className="bp5-progress-meter" style={memoryStyle}/>
				</div>
				<div
					className="bp5-progress-bar bp5-no-stripes bp5-intent-success"
					style={css.bar}
				>
					<div className="bp5-progress-meter" style={load1Style}/>
				</div>
				<div
					className="bp5-progress-bar bp5-no-stripes bp5-intent-warning"
					style={css.bar}
				>
					<div className="bp5-progress-meter" style={load5Style}/>
				</div>
			</div>
		</div>;
	}
}
