package utils

import (
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
)

type MemInfo struct {
	Total                uint64
	Free                 uint64
	Available            uint64
	Buffers              uint64
	Cached               uint64
	Used                 uint64
	UsedPercent          float64
	Dirty                uint64
	SwapTotal            uint64
	SwapFree             uint64
	SwapUsed             uint64
	SwapUsedPercent      float64
	HugePagesTotal       uint64
	HugePagesFree        uint64
	HugePagesReserved    uint64
	HugePagesUsed        uint64
	HugePagesUsedPercent float64
	HugePageSize         uint64
}

func GetMemInfo() (info *MemInfo, err error) {
	info = &MemInfo{}

	lines, err := ReadLines("/proc/meminfo")
	if err != nil {
		return
	}

	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		value = strings.Replace(value, " kB", "", -1)

		switch key {
		case "MemTotal":
			valueInt, e := strconv.ParseUint(value, 10, 64)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "utils: Failed to parse mem total"),
				}
				return
			}
			info.Total = valueInt
		case "MemFree":
			valueInt, e := strconv.ParseUint(value, 10, 64)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "utils: Failed to parse mem free"),
				}
				return
			}
			info.Free = valueInt
		case "MemAvailable":
			valueInt, e := strconv.ParseUint(value, 10, 64)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "utils: Failed to parse mem available"),
				}
				return
			}
			info.Available = valueInt
		case "Buffers":
			valueInt, e := strconv.ParseUint(value, 10, 64)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "utils: Failed to parse buffers"),
				}
				return
			}
			info.Buffers = valueInt
		case "Cached":
			valueInt, e := strconv.ParseUint(value, 10, 64)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "utils: Failed to parse cached"),
				}
				return
			}
			info.Cached = valueInt
		case "Dirty":
			valueInt, e := strconv.ParseUint(value, 10, 64)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "utils: Failed to parse dirty"),
				}
				return
			}
			info.Dirty = valueInt
		case "SwapTotal":
			valueInt, e := strconv.ParseUint(value, 10, 64)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "utils: Failed to parse swap total"),
				}
				return
			}
			info.SwapTotal = valueInt
		case "SwapFree":
			valueInt, e := strconv.ParseUint(value, 10, 64)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "utils: Failed to parse swap free"),
				}
				return
			}
			info.SwapFree = valueInt
		case "HugePages_Total":
			valueInt, e := strconv.ParseUint(value, 10, 64)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "utils: Failed to parse hugepages total"),
				}
				return
			}
			info.HugePagesTotal = valueInt
		case "HugePages_Free":
			valueInt, e := strconv.ParseUint(value, 10, 64)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "utils: Failed to parse hugepages total"),
				}
				return
			}
			info.HugePagesFree = valueInt
		case "HugePages_Rsvd":
			valueInt, e := strconv.ParseUint(value, 10, 64)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e,
						"utils: Failed to parse hugepages reserved"),
				}
				return
			}
			info.HugePagesReserved = valueInt
		case "Hugepagesize":
			valueInt, e := strconv.ParseUint(value, 10, 64)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "utils: Failed to parse hugepages size"),
				}
				return
			}
			info.HugePageSize = valueInt
		}
	}

	info.Used = info.Total - info.Free - info.Buffers - info.Cached
	info.UsedPercent = float64(info.Used) / float64(info.Total) * 100.0

	info.SwapUsed = info.SwapTotal - info.SwapFree
	if info.SwapUsed != 0 {
		info.SwapUsedPercent = float64(
			info.SwapUsed) / float64(info.SwapTotal) * 100.0
	}

	info.HugePagesUsed = (info.HugePagesTotal - info.HugePagesFree) +
		info.HugePagesReserved
	if info.HugePagesUsed != 0 {
		info.HugePagesUsedPercent = float64(
			info.HugePagesUsed) / float64(info.HugePagesTotal) * 100.0
	}

	return
}

type LoadStat struct {
	CpuUnits int
	Load1    float64
	Load5    float64
	Load15   float64
}

func LoadAverage() (ld *LoadStat, err error) {
	count := runtime.NumCPU()
	countFloat := float64(count)
	_ = countFloat

	ld = &LoadStat{}

	line, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "utils: Failed to read loadavg"),
		}
		return
	}

	values := strings.Fields(string(line))
	if len(values) < 3 {
		err = &errortypes.ParseError{
			errors.Wrap(err, "utils: Invalid loadavg data"),
		}
		return
	}

	load1, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "utils: Invalid load1 data"),
		}
		return
	}
	load5, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "utils: Invalid load5 data"),
		}
		return
	}
	load15, err := strconv.ParseFloat(values[2], 64)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "utils: Invalid load15 data"),
		}
		return
	}

	ld = &LoadStat{
		CpuUnits: count,
		Load1:    ToFixed(load1/countFloat*100, 2),
		Load5:    ToFixed(load5/countFloat*100, 2),
		Load15:   ToFixed(load15/countFloat*100, 2),
	}

	return
}
