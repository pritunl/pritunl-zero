/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Bar {
	progressClass?: string;
	label: string;
	value: number;
}

interface Props {
	hidden?: boolean;
	bars: Bar[];
}

const css = {
	label: {
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	bar: {
		marginBottom: '5px',
	} as React.CSSProperties,
};

export default class PageProgress extends React.Component<Props, {}> {
	render(): JSX.Element {
		let bars: JSX.Element[] = [];

		for (let bar of this.props.bars) {
			let style: React.CSSProperties = {
				width: (bar.value || 0) + '%',
			};

			bars.push(
				<div key={bar.label} style={css.bar}>
					{bar.label}
					<div className={'pt-progress-bar ' + (bar.progressClass || '')}>
						<div className="pt-progress-meter" style={style}/>
					</div>
				</div>,
			);
		}

		return <label
			className="pt-label"
			style={css.label}
			hidden={this.props.hidden}
		>
			{bars}
		</label>;
	}
}
