/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Loading from './Loading';

document.body.className = 'root pt-dark';

const css = {
	nav: {
		overflowX: 'auto',
		overflowY: 'hidden',
	} as React.CSSProperties,
	heading: {
		marginRight: '11px',
	} as React.CSSProperties,
};

export default class Main extends React.Component<{}, {}> {
	render(): JSX.Element {
		return <div>
			<nav className="pt-navbar layout horizontal" style={css.nav}>
				<div className="pt-navbar-group pt-align-left flex">
					<div className="pt-navbar-heading"
						style={css.heading}
					>Pritunl Zero</div>
					<Loading size="small"/>
				</div>
				<div className="pt-navbar-group pt-align-right">
					<button
						className="pt-button pt-minimal pt-icon-refresh"
						onClick={() => {}}
					>Refresh</button>
				</div>
			</nav>
		</div>;
	}
}
