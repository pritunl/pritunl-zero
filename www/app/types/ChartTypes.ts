/// <reference path="../References.d.ts"/>
import * as EndpointTypes from '../types/EndpointTypes';

export interface Point {
	x: number;
	y: number;
}
export type Points = Point[];
export type Chart = Points[];

export interface Dataset {
	label: string;
}
export type Datasets = Dataset[];

export type ChartData = {[key: string]: Points};

export interface Labels {
	title: string;
	resource_label: string;
	resource_type: string;
	resource_suffix: string;
	resource_fixed: number;
	resource_min: number;
	resource_max?: number;
	datasets: Datasets;
}

export function getChartLabels(resource: string, data: any): Labels {
	switch (resource) {
		case 'system':
			return {
				title: 'System Usage',
				resource_label: 'Percent',
				resource_type: 'float',
				resource_suffix: '%',
				resource_fixed: 3,
				resource_min: 0,
				resource_max: 100,
				datasets: [
					{
						label: 'CPU Usage',
					},
					{
						label: 'Memory Usage',
					},
					{
						label: 'Swap Usage',
					},
				],
			};
		case 'load':
			return {
				title: 'Load Average',
				resource_label: 'Load',
				resource_type: 'float',
				resource_suffix: '',
				resource_fixed: 4,
				resource_min: 0,
				datasets: [
					{
						label: 'Load1',
					},
					{
						label: 'Load5',
					},
					{
						label: 'Load15',
					},
				],
			};
		case 'disk':
			let diskData = data as EndpointTypes.DiskChart;
			let diskDatasets: Datasets = [];

			for (let key of Object.keys(diskData).sort()) {
				diskDatasets.push({
					label: key,
				} as Dataset);
			}

			return {
				title: 'Disks',
				resource_label: 'Usage',
				resource_type: 'float',
				resource_suffix: '%',
				resource_fixed: 3,
				resource_min: 0,
				resource_max: 100,
				datasets: diskDatasets,
			};
		case 'diskio0':
		case 'diskio1':
			let diskioData = data as EndpointTypes.NetworkChart;
			let diskioDatasets: Datasets = [];

			for (let key of Object.keys(diskioData).sort()) {
				let keys = key.split('-');
				let diskDevice = keys.slice(0, keys.length-1).join('-');
				let dataType = keys[keys.length-1];

				let label = '';

				if (resource === 'diskio0') {
					switch (dataType) {
						case 'br':
							label = 'Read';
							break;
						case 'bw':
							label = 'Written';
							break;
						default:
							continue;
					}
				} else {
					switch (dataType) {
						case 'tr':
							label = 'Read';
							break;
						case 'tw':
							label = 'Write';
							break;
						case 'ti':
							label = 'I/O';
							break;
						default:
							continue;
					}
				}

				diskioDatasets.push({
					label: diskDevice + ' ' + label,
				} as Dataset);
			}

			if (resource === 'diskio0') {
				return {
					title: 'Disk I/O',
					resource_label: 'Activity',
					resource_type: 'bytes',
					resource_suffix: '',
					resource_fixed: 2,
					resource_min: 0,
					datasets: diskioDatasets,
				};
			} else {
				return {
					title: 'Disk I/O Wait',
					resource_label: 'Waiting',
					resource_type: 'milliseconds',
					resource_suffix: '',
					resource_fixed: 2,
					resource_min: 0,
					datasets: diskioDatasets,
				};
			}
		case 'network':
			let netData = data as EndpointTypes.NetworkChart;
			let netDatasets: Datasets = [];

			for (let key of Object.keys(netData).sort()) {
				let keys = key.split('-');
				let iface = keys.slice(0, keys.length-1).join('-');
				let dataType = keys[keys.length-1];

				let label = '';
				switch (dataType) {
					case 'bs':
						label = 'Transmitted';
						break;
					case 'br':
						label = 'Received';
						break;
					default:
						label = 'Unknown';
				}

				netDatasets.push({
					label: iface + ' ' + label,
				} as Dataset);
			}

			return {
				title: 'Network Traffic',
				resource_label: 'Traffic',
				resource_type: 'bytes',
				resource_suffix: '',
				resource_fixed: 2,
				resource_min: 0,
				datasets: netDatasets,
			};
	}
	return undefined;
}

export function getChartData(resource: string, data: any): Chart {
	switch (resource) {
		case 'system':
			return [
				data.cpu_usage,
				data.mem_usage,
				data.swap_usage,
			];
		case 'load':
			return [
				data.load1,
				data.load5,
				data.load15,
			];
		case 'disk':
			let diskData = data as EndpointTypes.DiskChart;
			let diskChart: Chart = [];

			for (let key of Object.keys(diskData).sort()) {
				diskChart.push(diskData[key]);
			}

			return diskChart;
		case 'diskio0':
		case 'diskio1':
			let diskioData = data as EndpointTypes.DiskIoChart;
			let diskioChart: Chart = [];

			for (let key of Object.keys(diskioData).sort()) {
				let keys = key.split('-');
				let dataType = keys[keys.length-1];

				if (resource === 'diskio0') {
					if (dataType === 'br' || dataType === 'bw') {
						diskioChart.push(diskioData[key]);
					}
				} else {
					if (dataType === 'tr' || dataType === 'tw' || dataType === 'ti') {
						diskioChart.push(diskioData[key]);
					}
				}
			}

			return diskioChart;
		case 'network':
			let netData = data as EndpointTypes.NetworkChart;
			let netChart: Chart = [];

			for (let key of Object.keys(netData).sort()) {
				netChart.push(netData[key]);
			}

			return netChart;
	}

	return undefined;
}
