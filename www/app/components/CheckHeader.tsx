/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as CheckTypes from '../types/CheckTypes';

interface Props {
	header: CheckTypes.Header;
	onChange: (state: CheckTypes.Header) => void;
	onRemove: () => void;
}

const css = {
	group: {
		width: '100%',
		maxWidth: '310px',
		marginTop: '5px',
	} as React.CSSProperties,
	header: {
		width: '100%',
		borderRadius: '0 3px 3px 0',
	} as React.CSSProperties,
	headerBox: {
		flex: '1',
	} as React.CSSProperties,
};

export default class CheckHeader extends React.Component<Props, {}> {
	clone(): CheckTypes.Header {
		return {
			...this.props.header,
		};
	}

	render(): JSX.Element {
		let header = this.props.header;

		return <div className="bp5-control-group" style={css.group}>
			<div style={css.headerBox}>
				<input
					className="bp5-input"
					style={css.header}
					type="text"
					autoCapitalize="off"
					spellCheck={false}
					placeholder="Key"
					value={header.key || ''}
					onChange={(evt): void => {
						let state = this.clone();
						state.key = evt.target.value;
						this.props.onChange(state);
					}}
				/>
			</div>
			<div style={css.headerBox}>
				<input
					className="bp5-input"
					style={css.header}
					type="text"
					autoCapitalize="off"
					spellCheck={false}
					placeholder="Value"
					value={header.value || ''}
					onChange={(evt): void => {
						let state = this.clone();
						state.value = evt.target.value;
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
