/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	domain: string;
	disabled: boolean;
	onChange: (state: string) => void;
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

export default class CertificateDomain extends React.Component<Props, {}> {
	render(): JSX.Element {
		return <div className="bp5-control-group" style={css.group}>
			<div style={css.domainBox}>
				<input
					className="bp5-input"
					style={css.domain}
					disabled={this.props.disabled}
					type="text"
					autoCapitalize="off"
					spellCheck={false}
					placeholder="Domain"
					value={this.props.domain || ''}
					onChange={(evt): void => {
						this.props.onChange(evt.target.value);
					}}
				/>
			</div>
			<button
				className="bp5-button bp5-minimal bp5-intent-danger bp5-icon-remove"
				disabled={this.props.disabled}
				onClick={(): void => {
					this.props.onRemove();
				}}
			/>
		</div>;
	}
}
