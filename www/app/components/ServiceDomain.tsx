/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ServiceTypes from '../types/ServiceTypes';

interface Props {
	domain: ServiceTypes.Domain;
	onChange: (state: ServiceTypes.Domain) => void;
	onRemove: () => void;
}

const css = {
	group: {
		width: '100%',
		maxWidth: '310px',
		marginTop: '5px',
	} as React.CSSProperties,
	domain: {
		width: '100%',
		borderRadius: '0 3px 3px 0',
	} as React.CSSProperties,
	domainBox: {
		flex: '1',
	} as React.CSSProperties,
};

export default class ServiceDomain extends React.Component<Props, {}> {
	clone(): ServiceTypes.Domain {
		return {
			...this.props.domain,
		};
	}

	render(): JSX.Element {
		let domain = this.props.domain;

		return <div className="bp5-control-group" style={css.group}>
			<div style={css.domainBox}>
				<input
					className="bp5-input"
					style={css.domain}
					type="text"
					autoCapitalize="off"
					spellCheck={false}
					placeholder="Domain"
					value={domain.domain || ''}
					onChange={(evt): void => {
						let state = this.clone();
						state.domain = evt.target.value;
						this.props.onChange(state);
					}}
				/>
			</div>
			<div style={css.domainBox}>
				<input
					className="bp5-input"
					style={css.domain}
					type="text"
					autoCapitalize="off"
					spellCheck={false}
					placeholder="Host"
					value={domain.host || ''}
					onChange={(evt): void => {
						let state = this.clone();
						state.host = evt.target.value;
						this.props.onChange(state);
					}}
				/>
			</div>
			<button
				className="bp5-button bp5-minimal bp5-intent-danger bp5-icon-remove"
				onClick={(): void => {
					this.props.onRemove();
				}}
			/>
		</div>;
	}
}
