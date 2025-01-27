/// <reference path="../References.d.ts"/>
import * as React from "react"
import * as Blueprint from "@blueprintjs/core"
import * as Icons from "@blueprintjs/icons"
import Help from "./Help"
import PageSelectButton from "./PageSelectButton"

interface Props {
	hidden?: boolean
	disabled?: boolean
	title: string
	help: string
	addLabel: string
	menuLabel: string
	listMax?: number
	selected: Item[]
	options: Item[]
	icon: JSX.Element
	onAdd: (id: string) => void
	onRemove: (id: string) => void
}

interface State {
	selected: string
}

export interface Item {
	id: string
	name: string
}

const css = {
	label: {
		display: 'inline-block',
	} as React.CSSProperties,
	item: {
		margin: '9px 5px 0 5px',
		height: '20px',
	} as React.CSSProperties,
	menu: {
		maxHeight: '400px',
		overflowY: "auto",
	} as React.CSSProperties,
	menuOpen: {
		marginLeft: '0',
	} as React.CSSProperties,
	menuLabel: {
	} as React.CSSProperties,
	menuRemove: {
		opacity: 0.5,
	} as React.CSSProperties,
}

export class PageSelector extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context)
		this.state = {
			selected: "",
		}
	}

	render(): JSX.Element {
		let itemsBox: JSX.Element

		if ((this.props.selected || []).length > (this.props.listMax || 6)) {
			const itemMap = new Map<string, Array<Item>>();

			(this.props.selected || []).forEach((item) => {
				if (!itemMap.has(item.name)) {
					itemMap.set(item.name, [])
				}
				itemMap.get(item.name)!.push(item)
			})

			const items: JSX.Element[] = Array
				.from(itemMap.keys())
				.sort((a, b) => a.localeCompare(b))
				.flatMap((name) =>
					itemMap.get(name)!.map((item) =>
						<Blueprint.MenuItem
							key={item.id}
							disabled={this.props.disabled}
							selected={false}
							roleStructure="menuitem"
							icon={<Icons.Remove
								style={css.menuRemove}
							/>}
							onClick={(evt): void => {
								evt.stopPropagation()
								this.props.onRemove(item.id)
							}}
							text={item.name}
						/>
					)
				)

			itemsBox = <Blueprint.Popover
				content={<Blueprint.Menu style={css.menu}>
					{items}
				</Blueprint.Menu>}
				placement="bottom"
			>
				<Blueprint.Button
					alignText="left"
					icon={this.props.icon}
					rightIcon={<Icons.CaretDown/>}
					text={this.props.menuLabel}
					style={css.menuOpen}
				/>
			</Blueprint.Popover>
		} else {
			const itemMap = new Map<string, Array<Item>>();

			(this.props.selected || []).forEach((item) => {
				if (!itemMap.has(item.name)) {
					itemMap.set(item.name, []);
				}
				itemMap.get(item.name)!.push(item)
			})

			const items: JSX.Element[] = Array
				.from(itemMap.keys())
				.sort((a, b) => a.localeCompare(b))
				.flatMap((name) =>
					itemMap.get(name)!.map((item) =>
						<div
							className="bp5-tag bp5-tag-removable bp5-intent-primary"
							style={css.item}
							key={item.id}
						>
							{item.name}
							<button
								className="bp5-tag-remove"
								onMouseUp={(): void => {
									this.props.onRemove(item.id)
								}}
							/>
						</div>
					)
				)

			itemsBox = <div>{items}</div>
		}

		let selects: JSX.Element[] = [];
		if ((this.props.options || []).length === 0) {
			selects.push(<option key="null" value="">None</option>)
		} else {
			const optionsMap = new Map<string, Array<Item>>();

			(this.props.options || []).forEach((item) => {
				if (!optionsMap.has(item.name)) {
					optionsMap.set(item.name, [])
				}
				optionsMap.get(item.name)!.push(item)
			})

			selects = Array
				.from(optionsMap.keys())
				.sort((a, b) => a.localeCompare(b))
				.flatMap((name) =>
					optionsMap.get(name)!.map((item) =>
						<option key={item.id} value={item.id}>{item.name}</option>
					)
				)
		}

		return <div>
			<label
				className="bp5-label"
				style={css.label}
				hidden={this.props.hidden}
			>
				{this.props.title}
				<Help
					title={this.props.title}
					content={this.props.help}
				/>
				<div>
					{itemsBox}
				</div>
			</label>
			<PageSelectButton
				hidden={this.props.hidden}
				label={this.props.addLabel}
				value={this.state.selected}
				disabled={this.props.disabled}
				buttonClass="bp5-intent-success"
				onChange={(val: string): void => {
					this.setState({
						...this.state,
						selected: val,
					})
				}}
				onSubmit={() => {
					let id = this.state.selected
					if (!id && this.props.options) {
						id = this.props.options[0].id
					}

					if (id) {
						this.props.onAdd(id)
					}
				}}
			>
				{selects}
			</PageSelectButton>
		</div>
	}
}
