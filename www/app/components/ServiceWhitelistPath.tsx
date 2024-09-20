/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as ServiceTypes from '../types/ServiceTypes';

interface Props {
	path: ServiceTypes.Path;
	onChange: (state: ServiceTypes.Path) => void;
	onRemove: () => void;
}

const css = {
	group: {
		width: '100%',
		maxWidth: '310px',
		marginTop: '5px',
	} as React.CSSProperties,
	path: {
		width: '100%',
		borderRadius: '0 3px 3px 0',
	} as React.CSSProperties,
	pathBox: {
		flex: '1',
	} as React.CSSProperties,
};

export default class ServiceWhitelistPath extends React.Component<Props, {}> {
	clone(): ServiceTypes.Path {
		return {
			...this.props.path,
		};
	}

	render(): JSX.Element {
		let path = this.props.path;

		return <div className="bp5-control-group" style={css.group}>
			<div style={css.pathBox}>
				<input
					className="bp5-input"
					style={css.path}
					type="text"
					autoCapitalize="off"
					spellCheck={false}
					placeholder="Permitted path"
					value={path.path || ''}
					onChange={(evt): void => {
						let state = this.clone();
						state.path = evt.target.value;
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
