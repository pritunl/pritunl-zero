/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Help from './Help';

interface Props {
	buttonClass?: string;
	hidden?: boolean;
	disabled?: boolean;
	buttonDisabled?: boolean;
	readOnly?: boolean;
	autoSelect?: boolean;
	label?: string;
	labelTop?: boolean;
	listStyle?: boolean;
	help?: string;
	type: string;
	placeholder?: string;
	value: string;
	onChange?: (val: string) => void;
	onSubmit: () => void;
}

const css = {
	group: {
		marginBottom: '15px',
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	groupList: {
		marginTop: '5px',
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	groupTop: {
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	label: {
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	input: {
		width: '100%',
	} as React.CSSProperties,
	inputBox: {
		flex: '1',
	} as React.CSSProperties,
	buttonTop: {
		marginTop: '5px',
	} as React.CSSProperties,
};

export default class PageInputButton extends React.Component<Props, {}> {
	autoSelect = (evt: React.MouseEvent<HTMLInputElement>): void => {
		evt.currentTarget.select();
	}

	render(): JSX.Element {
		let buttonClass = 'pt-button';
		if (this.props.buttonClass) {
			buttonClass += ' ' + this.props.buttonClass;
		}

		if (this.props.labelTop) {
			return <label
				className="pt-label"
				style={css.label}
				hidden={this.props.hidden}
			>
				{this.props.label}
				<Help
					title={this.props.label}
					content={this.props.help}
				/>
				<div
					className="pt-control-group"
					style={css.groupTop}
					hidden={this.props.hidden}
				>
					<div style={css.inputBox}>
						<input
							className="pt-input"
							style={css.input}
							type={this.props.type}
							disabled={this.props.disabled}
							readOnly={this.props.readOnly}
							autoCapitalize="off"
							spellCheck={false}
							placeholder={this.props.placeholder}
							value={this.props.value || ''}
							onClick={this.props.autoSelect ? this.autoSelect : null}
							onChange={(evt): void => {
								if (this.props.onChange) {
									this.props.onChange(evt.target.value);
								}
							}}
							onKeyPress={(evt): void => {
								if (evt.key === 'Enter') {
									this.props.onSubmit();
								}
							}}
						/>
					</div>
					<div>
						<button
							className={buttonClass}
							style={css.buttonTop}
							disabled={this.props.disabled || this.props.buttonDisabled}
							onClick={this.props.onSubmit}
						/>
					</div>
				</div>
			</label>;
		} else {
			return <div
				className="pt-control-group"
				style={this.props.listStyle ? css.groupList : css.group}
				hidden={this.props.hidden}
			>
				<div style={css.inputBox}>
					<input
						className="pt-input"
						style={css.input}
						type={this.props.type}
						disabled={this.props.disabled}
						readOnly={this.props.readOnly}
						autoCapitalize="off"
						spellCheck={false}
						placeholder={this.props.placeholder || ''}
						value={this.props.value || ''}
						onChange={(evt): void => {
							if (this.props.onChange) {
								this.props.onChange(evt.target.value);
							}
						}}
						onKeyPress={(evt): void => {
							if (evt.key === 'Enter') {
								this.props.onSubmit();
							}
						}}
					/>
				</div>
				<div>
					<button
						className={buttonClass}
						disabled={this.props.disabled || this.props.buttonDisabled}
						onClick={this.props.onSubmit}
					>{this.props.label || ''}</button>
				</div>
			</div>;
		}
	}
}
