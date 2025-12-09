/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from "@blueprintjs/core"
import CopyButton from './CopyButton';

export interface Field {
	valueClass?: string;
	valueClasses?: string[];
	label: string;
	value?: string | number | string[];
	hover?: JSX.Element;
	copy?: boolean;
	embedded?: Props;
	maxLines?: number;
}

export interface Bar {
	progressClass?: string;
	label?: string;
	value: number;
	color?: string;
}

export interface Props {
	style?: React.CSSProperties;
	hidden?: boolean;
	fields?: Field[];
	bars?: Bar[];
	compact?: boolean;
}

const css = {
	label: {
		width: '100%',
		maxWidth: '320px',
	} as React.CSSProperties,
	value: {
		wordWrap: 'break-word',
	} as React.CSSProperties,
	valueLimit: {
		wordWrap: 'break-word',
		display: '-webkit-box',
		WebkitBoxOrient: 'vertical',
		overflow: 'hidden',
		textOverflow: 'ellipsis',
	} as React.CSSProperties,
	item: {
		marginTop: '0px',
		marginBottom: '5px',
	} as React.CSSProperties,
	itemCompact: {
		marginTop: '0px',
		marginBottom: '2px',
	} as React.CSSProperties,
	bar: {
		maxWidth: '280px',
	} as React.CSSProperties,
	embedded: {
		minWidth: '300px',
		padding: '10px',
	} as React.CSSProperties,
};

export default class PageInfo extends React.Component<Props, {}> {
	render(): JSX.Element {
		let fields: JSX.Element[] = [];
		let bars: JSX.Element[] = [];
		let itemStyle = this.props.compact ? css.itemCompact : css.item;

		for (let field of this.props.fields || []) {
			if (field == null) {
				continue;
			}

			let value: string | JSX.Element[];
			let copyBtn: JSX.Element;

			if (typeof field.value === 'string') {
				value = field.value;
				if (field.copy) {
					copyBtn = <CopyButton
						value={field.value}
					/>;
				}
			} else if (typeof field.value === 'number') {
				value = field.value.toString();
				if (field.copy) {
					copyBtn = <CopyButton
						value={field.value.toString()}
					/>;
				}
			} else if (field.value) {
				value = [];
				for (let i = 0; i < field.value.length; i++) {
					let copyItemBtn: JSX.Element;

					if (field.copy) {
						copyItemBtn = <CopyButton
							value={field.value[i]}
						/>;
					}

					value.push(
						<div
							key={i}
							className={
								field.valueClasses ?
								field.valueClasses[i] :
								(field.valueClass || 'bp5-text-muted')
							}
						>
							{field.value[i]}{copyItemBtn}
						</div>
					);
				}
			}

			if (field.hover || field.embedded) {
				fields.push(
					<Blueprint.Popover
						key={field.label}
						interactionKind="hover"
						placement="bottom"
						minimal={true}
						content={field.hover || <div
							style={css.embedded}
							className="bp5-content-popover">
								<PageInfo {...field.embedded}/>
							</div>
						}
						renderTarget={({isOpen, ...targetProps}): JSX.Element => {
								return <div {...targetProps} style={itemStyle}>
								{field.label}
								<div
									className={field.valueClass || 'bp5-text-muted'}
									style={css.value}
								>
									{value}{copyBtn}
								</div>
							</div>
						}}
					/>,
				);
			} else {
				let style = css.value
				if (field.maxLines) {
					style = {...css.valueLimit}
					style.WebkitLineClamp = field.maxLines
				}

				fields.push(
					<div key={field.label} style={itemStyle}>
						{field.label}
						<div
							className={field.valueClass || 'bp5-text-muted'}
							style={style}
						>
							{value}{copyBtn}
						</div>
					</div>,
				);
			}
		}

		if (this.props.bars) {
			for (let i = 0; i < this.props.bars.length; i++) {
				let bar = this.props.bars[i]

				let style: React.CSSProperties = {
					width: (bar.value || 0) + '%',
				};

				if (bar.color) {
					style.backgroundColor = bar.color;
				}

				bars.push(
					<div key={bar.label || i} style={itemStyle}>
						{bar.label}
						<div
							className={'bp5-progress-bar ' + (bar.progressClass || '')}
							style={css.bar}
						>
							<div className="bp5-progress-meter" style={style}/>
						</div>
					</div>,
				);
			}
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
			className="bp5-label"
			style={labelStyle}
			hidden={this.props.hidden}
		>
			{fields}
			{bars}
		</label>;
	}
}
