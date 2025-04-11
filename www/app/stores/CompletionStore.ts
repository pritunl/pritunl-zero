/// <reference path="../References.d.ts"/>
import Dispatcher from "../dispatcher/Dispatcher"
import EventEmitter from "../EventEmitter"
import * as CompletionTypes from "../types/CompletionTypes"
import * as GlobalTypes from "../types/GlobalTypes"

class CompletionStore extends EventEmitter {
	_completion: CompletionTypes.Completion = Object.freeze({})
	_filter: CompletionTypes.Filter = null;
	_token = Dispatcher.register((this._callback).bind(this))

	_reset(): void {
		this._completion = Object.freeze({})
		this._filter = null
		this.emitChange()
	}

	get completion(): CompletionTypes.Completion {
		return this._completion
	}

	get filter(): CompletionTypes.Filter {
		return this._filter;
	}

	emitChange(): void {
		this.emitDefer(GlobalTypes.CHANGE)
	}

	addChangeListener(callback: () => void): void {
		this.on(GlobalTypes.CHANGE, callback)
	}

	removeChangeListener(callback: () => void): void {
		this.removeListener(GlobalTypes.CHANGE, callback)
	}

	_filterCallback(filter: CompletionTypes.Filter): void {
		this._filter = filter
		this.emitChange()
	}

	_sync(completion: CompletionTypes.Completion): void {
		this._completion = Object.freeze(completion)
		this.emitChange()
	}

	_callback(action: CompletionTypes.CompletionDispatch): void {
		switch (action.type) {
			case GlobalTypes.RESET:
				this._reset()
				break

			case CompletionTypes.FILTER:
				this._filterCallback(action.data.filter)
				break

			case CompletionTypes.SYNC:
				this._sync(action.data.completion)
				break
		}
	}
}

export default new CompletionStore()
