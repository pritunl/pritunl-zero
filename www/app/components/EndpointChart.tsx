/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ChartJs from 'chart.js';
import * as EndpointActions from '../actions/EndpointActions';
import * as ChartTypes from '../types/ChartTypes';
import * as MiscUtils from '../utils/MiscUtils';

interface Props {
	endpoint: string;
	resource: string;
	sync: number;
	period: number;
	interval: number;
	onLoading: () => void;
	onLoaded: () => void;
}

interface State {
	disabled: boolean;
}

// const colors = [
// 	'#d50000', // red
// 	'#c51162', // pink
// 	'#aa00ff', // purple
// 	'#6200ea', // deep purple
// 	'#304ffe', // indigo
// 	'#2962ff', // blue
// 	'#0091ea', // light blue
// 	'#00b8d4', // cyan
// 	'#00bfa5', // teal
// 	'#00c853', // green
// 	'#64dd17', // light green
// 	'#aeea00', // lime
// 	'#ffd600', // yellow
// 	'#ffab00', // amber
// 	'#ff6d00', // orange
// 	'#dd2c00', // deep orange
// 	'#5d4037', // brown
// 	'#455a64', // blue grey
// ];

const colors = [
	'#0091ea', // light blue
	'#d50000', // red
	'#00c853', // green
	'#aa00ff', // purple
	'#ffab00', // amber
	'#c51162', // pink
	'#2962ff', // blue
	'#ff6d00', // orange
	'#00bfa5', // teal
	'#304ffe', // indigo
	'#00b8d4', // cyan
	'#6200ea', // deep purple
	'#ffd600', // yellow
	'#dd2c00', // deep orange
	'#5d4037', // brown
	'#455a64', // blue grey
	'#64dd17', // light green
	'#aeea00', // lime

	'#0091ea', // light blue
	'#d50000', // red
	'#00c853', // green
	'#aa00ff', // purple
	'#ffab00', // amber
	'#c51162', // pink
	'#2962ff', // blue
	'#ff6d00', // orange
	'#00bfa5', // teal
	'#304ffe', // indigo
	'#00b8d4', // cyan
	'#6200ea', // deep purple
	'#ffd600', // yellow
	'#dd2c00', // deep orange
	'#5d4037', // brown
	'#455a64', // blue grey
	'#64dd17', // light green
	'#aeea00', // lime

	'#0091ea', // light blue
	'#d50000', // red
	'#00c853', // green
	'#aa00ff', // purple
	'#ffab00', // amber
	'#c51162', // pink
	'#2962ff', // blue
	'#ff6d00', // orange
	'#00bfa5', // teal
	'#304ffe', // indigo
	'#00b8d4', // cyan
	'#6200ea', // deep purple
	'#ffd600', // yellow
	'#dd2c00', // deep orange
	'#5d4037', // brown
	'#455a64', // blue grey
	'#64dd17', // light green
	'#aeea00', // lime

	'#0091ea', // light blue
	'#d50000', // red
	'#00c853', // green
	'#aa00ff', // purple
	'#ffab00', // amber
	'#c51162', // pink
	'#2962ff', // blue
	'#ff6d00', // orange
	'#00bfa5', // teal
	'#304ffe', // indigo
	'#00b8d4', // cyan
	'#6200ea', // deep purple
	'#ffd600', // yellow
	'#dd2c00', // deep orange
	'#5d4037', // brown
	'#455a64', // blue grey
	'#64dd17', // light green
	'#aeea00', // lime
];

export default class EndpointChart extends React.Component<Props, State> {
	data: ChartTypes.ChartData;
	sync: number;
	period: number;
	interval: number;
	chart: ChartJs.Chart;
	chartRef: React.RefObject<HTMLCanvasElement>;

	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
		};
		this.chartRef = React.createRef();
	}

	ticks = (axis: ChartJs.Scale) => {
		let ticks = axis.ticks;
		let newTicks: ChartJs.Tick[] = [];
		let dataset = Object.values(this.data)[0];
		let tickMod = 3600000; // 1 hour
		let len = dataset.length;

		if (len) {
			let first = dataset[0] as ChartJs.ScatterDataPoint;
			let last = dataset[len-1] as ChartJs.ScatterDataPoint;
			let range = last.x - first.x;

			if (range >= 2833920000) {
				tickMod = 604800000; // 7 day
			} else if (range >= 1451520000) {
				tickMod = 172800000; // 2 day
			} else if (range >= 611280000) {
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

	config = (): ChartJs.ChartConfiguration => {
		let labels = ChartTypes.getChartLabels(this.props.resource, this.data);

		let config = {
			type: 'line',
			options: {
				scales: {
					x: {
						type: 'time',
						title: {
							display: true,
							text: 'Time',
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
						min: labels.resource_min,
						max: labels.resource_max,
						offset: false,
						beginAtZero: true,
						title: {
							display: true,
							text: labels.resource_label,
							color: 'rgba(255, 255, 255, 1)',
							padding: 0,
							font: {
								weight: 'bold',
							},
						},
						ticks: {
							color: 'rgba(255, 255, 255, 1)',
							callback: function(val: number): number | string {
								switch (labels.resource_type) {
									case 'bytes':
										return MiscUtils.formatBytes(val, 0);
									case 'milliseconds':
										return MiscUtils.formatMs(val);
									default:
										return val;
								}
							}
						},
						grid: {
							color: 'rgba(255, 255, 255, 0.2)',
						},
					},
				},
				plugins: {
					title: {
						display: true,
						text: labels.title,
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

								let val = '';
								switch (labels.resource_type) {
									case 'bytes':
										val = MiscUtils.formatBytes(raw.y, labels.resource_fixed);
										break;
									case 'milliseconds':
										val = MiscUtils.formatMs(raw.y);
										break;
									case 'float':
										val = raw.y.toFixed(labels.resource_fixed);
										break;
									default:
										val = raw.y;
								}

								if (labels.resource_fixed) {
									return item.dataset.label + ' ' +
										val + labels.resource_suffix;
								}
								return item.dataset.label + ' ' + val +
									labels.resource_suffix;
							},
						},
					},
				},
			},
			data: {
				datasets: [],
			},
		} as ChartJs.ChartConfiguration;

		let data = ChartTypes.getChartData(this.props.resource, this.data);
		for (let i = 0; i < labels.datasets.length; i++) {
			let datasetLabels = labels.datasets[i];

			config.data.datasets.push({
				label: datasetLabels.label,
				data: data[i],
				fill: 'origin',
				pointRadius: 0,
				backgroundColor: colors[i] + '15',
				borderColor: colors[i],
				borderWidth: 2,
			} as ChartJs.ChartDataset);
		}

		return config;
	}

	update(sync: number, period: number, interval: number): void {
		this.sync = sync;
		this.period = period;
		this.interval = interval;

		let loading = true;
		this.props.onLoading();

		EndpointActions.chart(
			this.props.endpoint,
			this.props.resource,
			this.period,
			this.interval,
		).then((data: ChartTypes.ChartData): void => {
			if (loading) {
				loading = false;
				this.props.onLoaded();
			}

			if (data) {
				this.data = data;
				this.updateChart();
			}
		}).catch((): void => {
			if (loading) {
				loading = false;
				this.props.onLoaded();
			}
		});
	}

	updateChart(): void {
		let data = ChartTypes.getChartData(this.props.resource, this.data);

		for (let i = 0; i < data.length; i++) {
			this.chart.data.datasets[i].data = data[i];
		}

		this.chart.update();
	}

	componentDidMount(): void {
		this.sync = this.props.sync;
		this.period = this.props.period;
		this.interval = this.props.interval;

		let loading = true;
		this.props.onLoading();

		EndpointActions.chart(
			this.props.endpoint,
			this.props.resource,
			this.period,
			this.interval,
		).then((data: ChartTypes.ChartData): void => {
			if (loading) {
				loading = false;
				this.props.onLoaded();
			}

			if (data) {
				this.data = data;
				this.chart = new ChartJs.Chart(
					this.chartRef.current,
					this.config(),
				);
			}
		}).catch((): void => {
			if (loading) {
				loading = false;
				this.props.onLoaded();
			}
		});
	}

	componentWillUnmount(): void {
		if (this.chart) {
			this.chart.destroy();
		}
	}

	render(): JSX.Element {
		if ((this.sync !== undefined && this.period !== undefined &&
				this.interval !== undefined) &&
				(this.props.sync !== this.sync ||
				this.props.period !== this.period ||
				this.props.interval !== this.interval)) {
			this.update(this.props.sync, this.props.period, this.props.interval);
		}

		return <canvas
			ref={this.chartRef}
		/>;
	}
}
