/// <reference path="../References.d.ts"/>
import * as React from 'react';

interface Props {
	domain: string;
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
		return <div className="bp3-control-group" style={css.group}>
			<div style={css.domainBox}>
				<input
					className="bp3-input"
					style={css.domain}
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
				className="bp3-button bp3-minimal bp3-intent-danger bp3-icon-remove"
				onClick={(): void => {
					this.props.onRemove();
				}}
			/>
		</div>;
	}
}
