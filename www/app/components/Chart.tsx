/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ChartJs from 'chart.js';
import * as MiscUtils from '../utils/MiscUtils';
import * as EndpointActions from '../actions/EndpointActions';
import * as EndpointTypes from '../types/EndpointTypes';
import Help from './Help';

ChartJs.Chart.register(ChartJs.LineController);
ChartJs.Chart.register(ChartJs.CategoryScale);
ChartJs.Chart.register(ChartJs.LinearScale);
ChartJs.Chart.register(ChartJs.TimeScale);
ChartJs.Chart.register(ChartJs.PointElement);
ChartJs.Chart.register(ChartJs.LineElement);
ChartJs.Chart.register(ChartJs.Tooltip);
ChartJs.Chart.register(ChartJs.Filler);

interface Props {
	disabled?: boolean;
}

interface State {
	disabled: boolean;
}

const css = {
};

class LineTracerController extends ChartJs.LineController {
	draw(): void {
		super.draw();

		let chart = this.chart as any;
		if (chart.tooltip._active && chart.tooltip._active.length) {
			let ctx = this.chart.ctx;
			let x = chart.tooltip.caretX;
			let topY = chart.scales.y.top;
			let bottomY = chart.scales.y.bottom;

			ctx.save();
			ctx.beginPath();
			ctx.moveTo(x, topY);
			ctx.lineTo(x, bottomY);
			ctx.lineWidth = 1;
			ctx.strokeStyle = 'rgba(255, 255, 255, 0.6)';
			ctx.stroke();
			ctx.restore();
		}
	}
}
(ChartJs.Chart as any).registry.controllers.items.line = LineTracerController;

export default class Chart extends React.Component<Props, State> {
	chartRef: React.RefObject<HTMLCanvasElement>;

	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
		};

		this.chartRef = React.createRef();
	}

	config(data: EndpointTypes.SystemChart): ChartJs.ChartConfiguration {
		let dataMem = [] as ChartJs.ScatterDataPoint[];
		let cur = 40;
		for (let i = 0; i < 188; i++) {
			cur += MiscUtils.random(-2000, 2000) / 1000;
			dataMem.push({
				x: new Date(2021, 5, 12, 5, 20).getTime() + (i * 300000),
				y: cur,
			})
		}

		let dataCpu = [] as ChartJs.ScatterDataPoint[];
		cur = 20;
		for (let i = 0; i < 188; i++) {
			cur += MiscUtils.random(-2000, 2000) / 1000;
			dataCpu.push({
				x: new Date(2021, 5, 12, 5, 20).getTime() + (i * 300000),
				y: cur,
			})
		}

		return {
			type: 'line',
			options: {
				scales: {
					x: {
						type: 'time',
						title: {
							display: true,
							text: 'Time (UTC)',
							color: 'rgba(255, 255, 255, 1)',
							font: {
								weight: 'bold',
							},
						},
						time: {
							unit: 'minute',
							displayFormats: {
								minute: 'HH:mm',
							},
						},
						ticks: {
							stepSize: 1,
							count: 100,
							maxTicksLimit: 100,
							color: 'rgba(255, 255, 255, 1)',
							source: 'data',
						},
						grid: {
							color: 'rgba(255, 255, 255, 0.2)',
						},
						beforeTickToLabelConversion(axis: ChartJs.Scale) {
							let ticks = axis.ticks;
							let newTicks: ChartJs.Tick[] = [];

							for (let i = 0; i < ticks.length; i++) {
								let tick = ticks[i];
								if (tick.value % 3600000 === 0) {
									newTicks.push(tick);
								}
							}

							axis.ticks = newTicks;
						},
					},
					y: {
						min: 0,
						max: 100,
						offset: false,
						beginAtZero: true,
						title: {
							display: true,
							text: 'Percent',
							color: 'rgba(255, 255, 255, 1)',
							font: {
								weight: 'bold',
							},
						},
						ticks: {
							color: 'rgba(255, 255, 255, 1)',
						},
						grid: {
							color: 'rgba(255, 255, 255, 0.2)',
						},
					},
				},
				plugins: {
					tooltip: {
						mode: 'index',
						intersect: false,
						backgroundColor: 'rgba(0, 0, 0, 0.7)',
						callbacks: {
							label(item): string {
								let raw = item.raw as any;
								return item.dataset.label + ' ' + raw.y.toFixed(2) + '%';
							},
						},
					},
				},
			},
			data: {
				datasets: [
					{
						label: 'CPU Usage',
						data: data.cpu_usage,
						fill: 'origin',
						pointRadius: 0,
						backgroundColor: 'rgba(19, 124, 189, 0.2)',
						borderColor: 'rgba(19, 124, 189, 1)',
						borderWidth: 2,
					},
					{
						label: 'Memory Usage',
						data: data.mem_usage,
						fill: 'origin',
						pointRadius: 0,
						backgroundColor: 'rgba(255, 99, 132, 0.2)',
						borderColor: 'rgba(255, 99, 132, 1)',
						borderWidth: 2,
					},
				],
			},
		} as ChartJs.ChartConfiguration;
	}

	componentDidMount(): void {
		EndpointActions.chart('5facaf6119095293ebb71257', 'system').then((
				data: EndpointTypes.SystemChart): void => {
			let chart = new ChartJs.Chart(
				this.chartRef.current,
				this.config(data),
			);
		}).catch((): void => {
		});

	}

	componentWillUnmount(): void {
	}

	render(): JSX.Element {
		return <div>
			<div className="bp3-border">Memory Usage</div>
			<canvas
				ref={this.chartRef}
			/>
		</div>;
	}
}
