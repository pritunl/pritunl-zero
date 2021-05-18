/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ChartJs from 'chart.js';
import * as MiscUtils from '../utils/MiscUtils';
import * as EndpointActions from '../actions/EndpointActions';
import * as EndpointTypes from '../types/EndpointTypes';
import * as Blueprint from '@blueprintjs/core';
import Help from './Help';

ChartJs.Chart.register(ChartJs.LineController);
ChartJs.Chart.register(ChartJs.CategoryScale);
ChartJs.Chart.register(ChartJs.LinearScale);
ChartJs.Chart.register(ChartJs.TimeScale);
ChartJs.Chart.register(ChartJs.PointElement);
ChartJs.Chart.register(ChartJs.LineElement);
ChartJs.Chart.register(ChartJs.Title);
ChartJs.Chart.register(ChartJs.Tooltip);
ChartJs.Chart.register(ChartJs.Filler);

interface Props {
	disabled?: boolean;
}

interface State {
	hours: number;
	disabled: boolean;
}

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
	data: EndpointTypes.SystemChart;
	chart: ChartJs.Chart;
	chartRef: React.RefObject<HTMLCanvasElement>;

	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			hours: 24,
			disabled: false,
		};

		this.chartRef = React.createRef();
	}

	ticks = (axis: ChartJs.Scale) => {
		let ticks = axis.ticks;
		let newTicks: ChartJs.Tick[] = [];
		let dataset = this.data.cpu_usage;
		let tickMod = 3600000; // 1 hour
		let len = dataset.length;

		if (len) {
			let first = dataset[0] as ChartJs.ScatterDataPoint;
			let last = dataset[len-1] as ChartJs.ScatterDataPoint;
			let range = last.x - first.x;

			if (range >= 611280000) {
				tickMod = 86400000; // 1 day
			} else if (range >= 276480000) {
				tickMod = 43200000; // 12 hours
			} else if (range >= 89280000) {
				tickMod = 21600000; // 6 hours
			} else {
				tickMod = 3600000; // 1 hours
			}
		}

		for (let i = 0; i < ticks.length; i++) {
			let tick = ticks[i];

			if (tick.value % tickMod === 0) {
				newTicks.push(tick);
			}
		}

		axis.ticks = newTicks;
	}

	config(): ChartJs.ChartConfiguration {
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
							padding: 0,
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
						beforeTickToLabelConversion: this.ticks,
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
							padding: 0,
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
					title: {
						display: true,
						text: 'System Usage',
						color: 'rgba(255, 255, 255, 1)',
						padding: 3,
						font: {
							size: 13,
						},
					},
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
						data: this.data.cpu_usage,
						fill: 'origin',
						pointRadius: 0,
						backgroundColor: 'rgba(19, 124, 189, 0.2)',
						borderColor: 'rgba(19, 124, 189, 1)',
						borderWidth: 2,
					},
					{
						label: 'Memory Usage',
						data: this.data.mem_usage,
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

	update(hours: number): void {
		this.setState({
			...this.state,
			hours: hours,
		});
		EndpointActions.chart('5facaf6119095293ebb71257', 'system',
				hours).then((data: EndpointTypes.SystemChart): void => {
			this.data = data;
			this.updateChart();
		});
	}

	updateChart(): void {
		this.chart.data.datasets[0].data = this.data.cpu_usage;
		this.chart.data.datasets[1].data = this.data.mem_usage;
		this.chart.update();
	}

	componentDidMount(): void {
		EndpointActions.chart('5facaf6119095293ebb71257', 'system',
				this.state.hours).then((data: EndpointTypes.SystemChart): void => {
			this.data = data;
			this.chart = new ChartJs.Chart(
				this.chartRef.current,
				this.config(),
			);
		});
	}

	componentWillUnmount(): void {
	}

	render(): JSX.Element {
		return <div>
			<Blueprint.NumericInput
				allowNumericCharactersOnly={true}
				min={1}
				minorStepSize={1}
				stepSize={1}
				majorStepSize={10}
				disabled={this.props.disabled}
				selectAllOnFocus={true}
				onValueChange={(val: number): void => {
					this.update(val);
				}}
				value={this.state.hours}
			/>
			<canvas
				ref={this.chartRef}
			/>
		</div>;
	}
}
