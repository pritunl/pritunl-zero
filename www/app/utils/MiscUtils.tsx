/// <reference path="../References.d.ts"/>
import React from "react";

export class SyncInterval {
  private timer: number | null = null;
  private cancel: boolean = false;
  private readonly interval: number;
  private readonly action: () => Promise<any>;

  constructor(action: () => Promise<any>, interval: number) {
    this.action = action;
    this.interval = interval;
		this.start();
  }

  public start = async (): Promise<void> => {
    if (this.timer !== null) {
      clearTimeout(this.timer);
      this.timer = null;
    }

    this.cancel = false;

    const runSync = async (): Promise<void> => {
      if (this.cancel) return;

      try {
        await this.action();

        if (!this.cancel) {
          this.timer = window.setTimeout(() => {
            runSync();
          }, this.interval);
        }
      } catch (error) {
        console.error("Action error:", error);
        if (!this.cancel) {
          this.timer = window.setTimeout(() => {
            runSync();
          }, this.interval);
        }
      }
    };

    runSync();
  };

  public stop = (): void => {
    this.cancel = true;
    if (this.timer) {
      clearTimeout(this.timer);
      this.timer = null;
    }
  };
}

export function uuid(): string {
	return (+new Date() + Math.floor(Math.random() * 999999)).toString(36);
}

export function random(min: number, max: number): number {
	return Math.round(Math.random() * (max - min) + min);
}

export function zeroPad(num: number, width: number): string {
	if (num < Math.pow(10, width)) {
		return ('0'.repeat(width - 1) + num).slice(-width);
	}
	return num.toString();
}

export function capitalize(str: string): string {
	if (!str) {
		return str;
	}
	return str.charAt(0).toUpperCase() + str.slice(1);
}

export function titleCase(str: string): string {
	if (!str) {
		return str;
	}
	return str
		.toLowerCase()
		.split(' ')
		.map(word => word.charAt(0).toUpperCase() + word.slice(1))
		.join(' ');
}

export function formatAmount(amount: number): string {
	if (!amount) {
		return '-';
	}
	return '$' + (amount / 100).toFixed(2);
}

export function formatBytes(bytes: number, decimals: number): string {
	if (!bytes) {
		return '0B';
	}	else if (bytes < 1024) {
		return bytes + 'B';
	} else if (bytes < 1048576) {
		return Math.round(bytes / 1024).toFixed(decimals) + 'KB';
	} else if (bytes < 1073741824) {
		return (bytes / 1048576).toFixed(decimals) + 'MB';
	} else if (bytes < 1099511627776) {
		return (bytes / 1073741824).toFixed(decimals) + 'GB';
	} else {
		return (bytes / 1099511627776).toFixed(decimals) + 'TB';
	}
}

export function formatMs(ms: number): string {
	if (ms < 1000) {
		return ms + 'ms';
	} else {
		return (ms / 1000) + 's';
	}
}

export function formatUptime(time: number): string {
	let days = Math.floor(time / 86400);
	time -= days * 86400;
	let hours = Math.floor(time / 3600);
	time -= hours * 3600;
	let minutes = Math.floor(time / 60);
	time -= minutes * 60;
	return days + 'd ' + hours + 'h ' + minutes + 'm ' + time + 's';
}

export function formatDate(dateStr: string): string {
	if (!dateStr || dateStr === '0001-01-01T00:00:00Z') {
		return '';
	}

	let date = new Date(dateStr);
	let str = '';

	let hours = date.getHours();
	let period = 'AM';

	if (hours > 12) {
		period = 'PM';
		hours -= 12;
	} else if (hours === 0) {
		hours = 12;
	}

	let day;
	switch (date.getDay()) {
		case 0:
			day = 'Sun';
			break;
		case 1:
			day = 'Mon';
			break;
		case 2:
			day = 'Tue';
			break;
		case 3:
			day = 'Wed';
			break;
		case 4:
			day = 'Thu';
			break;
		case 5:
			day = 'Fri';
			break;
		case 6:
			day = 'Sat';
			break;
	}

	let month;
	switch (date.getMonth()) {
		case 0:
			month = 'Jan';
			break;
		case 1:
			month = 'Feb';
			break;
		case 2:
			month = 'Mar';
			break;
		case 3:
			month = 'Apr';
			break;
		case 4:
			month = 'May';
			break;
		case 5:
			month = 'Jun';
			break;
		case 6:
			month = 'Jul';
			break;
		case 7:
			month = 'Aug';
			break;
		case 8:
			month = 'Sep';
			break;
		case 9:
			month = 'Oct';
			break;
		case 10:
			month = 'Nov';
			break;
		case 11:
			month = 'Dec';
			break;
	}

	str += day + ' ';
	str += date.getDate() + ' ';
	str += month + ' ';
	str += date.getFullYear() + ', ';
	str += hours + ':';
	str += zeroPad(date.getMinutes(), 2) + ':';
	str += zeroPad(date.getSeconds(), 2) + ' ';
	str += period;

	return str;
}

export function formatDateShort(dateStr: string): string {
	if (!dateStr || dateStr === '0001-01-01T00:00:00Z') {
		return '';
	}

	let date = new Date(dateStr);
	let curDate = new Date();

	let month;
	switch (date.getMonth()) {
		case 0:
			month = 'Jan';
			break;
		case 1:
			month = 'Feb';
			break;
		case 2:
			month = 'Mar';
			break;
		case 3:
			month = 'Apr';
			break;
		case 4:
			month = 'May';
			break;
		case 5:
			month = 'Jun';
			break;
		case 6:
			month = 'Jul';
			break;
		case 7:
			month = 'Aug';
			break;
		case 8:
			month = 'Sep';
			break;
		case 9:
			month = 'Oct';
			break;
		case 10:
			month = 'Nov';
			break;
		case 11:
			month = 'Dec';
			break;
	}

	let str = month + ' ' + date.getDate();

	if (date.getFullYear() !== curDate.getFullYear()) {
		str += ' ' + date.getFullYear();
	}

	return str;
}

export function formatDateShortTime(dateStr: string): string {
	if (!dateStr || dateStr === '0001-01-01T00:00:00Z') {
		return '';
	}

	let date = new Date(dateStr);
	let curDate = new Date();

	let month;
	switch (date.getMonth()) {
		case 0:
			month = 'Jan';
			break;
		case 1:
			month = 'Feb';
			break;
		case 2:
			month = 'Mar';
			break;
		case 3:
			month = 'Apr';
			break;
		case 4:
			month = 'May';
			break;
		case 5:
			month = 'Jun';
			break;
		case 6:
			month = 'Jul';
			break;
		case 7:
			month = 'Aug';
			break;
		case 8:
			month = 'Sep';
			break;
		case 9:
			month = 'Oct';
			break;
		case 10:
			month = 'Nov';
			break;
		case 11:
			month = 'Dec';
			break;
	}

	let str = month + ' ' + date.getDate();

	if (date.getFullYear() !== curDate.getFullYear()) {
		str += ' ' + date.getFullYear();
	} else if (date.getMonth() === curDate.getMonth() &&
			date.getDate() === curDate.getDate()) {
		let hours = date.getHours();
		let period = 'AM';

		if (hours > 12) {
			period = 'PM';
			hours -= 12;
		} else if (hours === 0) {
			hours = 12;
		}

		str = hours + ':';
		str += zeroPad(date.getMinutes(), 2) + ':';
		str += zeroPad(date.getSeconds(), 2) + ' ';
		str += period;
	}

	return str;
}

export function highlightMatch(input: string, query: string): React.ReactNode {
	if (!query) {
		return input;
	}

	let index = input.toLowerCase().indexOf(query.toLowerCase())
	if (index === -1) {
		return input;
	}

	return <span>
		{input.substring(0, index)}
		<b>{input.substring(index, index + query.length)}</b>
		{input.substring(index + query.length)}
	</span>;
}
