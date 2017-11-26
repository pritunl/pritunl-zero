/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Field {
	valueClass?: string;
	label: string;
	value: string | string[];
}

interface Bar {
	progressClass?: string;
	label: string;
	value: number;
}

interface Props {
	style?: React.CSSProperties;
	hidden?: boolean;
	fields?: Field[];
	bars?: Bar[];
}

const css = {
	label: {
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	value: {
		wordWrap: 'break-word',
	} as React.CSSProperties,
	item: {
		marginBottom: '5px',
	} as React.CSSProperties,
};

export default class PageInfo extends React.Component<Props, {}> {
	render(): JSX.Element {
		let fields: JSX.Element[] = [];
		let bars: JSX.Element[] = [];

		for (let field of this.props.fields || []) {
			let value: string | JSX.Element[];

			if (typeof field.value === 'string') {
				value = field.value;
			} else {
				value = [];
				for (let i = 0; i < field.value.length; i++) {
					value.push(<div key={i}>{field.value[i]}</div>);
				}
			}

			fields.push(
				<div key={field.label} style={css.item}>
					{field.label}
					<div
						className={field.valueClass || 'pt-text-muted'}
						style={css.value}
					>
						{value}
					</div>
				</div>,
			);
		}

		for (let bar of this.props.bars || []) {
			let style: React.CSSProperties = {
				width: (bar.value || 0) + '%',
			};

			bars.push(
				<div key={bar.label} style={css.item}>
					{bar.label}
					<div className={'pt-progress-bar ' + (bar.progressClass || '')}>
						<div className="pt-progress-meter" style={style}/>
					</div>
				</div>,
			);
		}

		let labelStyle: React.CSSProperties;
		if (this.props.style) {
			labelStyle = {
				...css.label,
				...this.props.style,
			};
		} else {
			labelStyle = css.label;
		}

		return <label
			className="pt-label"
			style={labelStyle}
			hidden={this.props.hidden}
		>
			{fields}
			{bars}
		</label>;
	}
}
