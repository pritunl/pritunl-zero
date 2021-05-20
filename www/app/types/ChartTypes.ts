/// <reference path="../References.d.ts"/>
import * as ChartJs from 'chart.js';

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

export interface Labels {
	title: string;
	resource_label: string;
	resource_suffix: string;
	resource_fixed: number;
	resource_min: number;
	resource_max: number;
	datasets: Datasets;
}

export function getChartLabels(resource: string, data: any): Labels {
	switch (resource) {
		case 'system':
			return {
				title: 'System Usage',
				resource_label: 'Percent',
				resource_suffix: '%',
				resource_fixed: 2,
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
				]
			};
		case 'load':
			return {
				title: 'Load Average',
				resource_label: 'Load',
				resource_suffix: '',
				resource_fixed: 2,
				resource_min: 0,
				resource_max: undefined,
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
				]
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
	}
	return undefined;
}
