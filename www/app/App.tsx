/// <reference path="References.d.ts"/>
import 'chartjs-adapter-moment';
import * as ChartJs from 'chart.js';
import * as React from 'react';
import * as ReactDOM from 'react-dom';
import * as Blueprint from '@blueprintjs/core';
import Main from './components/Main';
import * as Alert from './Alert';
import * as Event from './Event';
import * as Theme from './Theme';
import * as Csrf from './Csrf';

ChartJs.Chart.register(ChartJs.LineController);
ChartJs.Chart.register(ChartJs.CategoryScale);
ChartJs.Chart.register(ChartJs.LinearScale);
ChartJs.Chart.register(ChartJs.TimeScale);
ChartJs.Chart.register(ChartJs.PointElement);
ChartJs.Chart.register(ChartJs.LineElement);
ChartJs.Chart.register(ChartJs.Title);
ChartJs.Chart.register(ChartJs.Tooltip);
ChartJs.Chart.register(ChartJs.Filler);

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
			ctx.lineWidth = 0.7;
			ctx.strokeStyle = Theme.chartColor3();
			ctx.stroke();
			ctx.restore();
		}
	}
}
(ChartJs.Chart as any).registry.controllers.items.line = LineTracerController;

Csrf.load().then((): void => {
	Blueprint.FocusStyleManager.onlyShowFocusOnTabs();
	Alert.init();
	Event.init();

	ReactDOM.render(
		<Blueprint.OverlaysProvider>
			<div><Main/></div>
		</Blueprint.OverlaysProvider>,
		document.getElementById('app'),
	);
});
