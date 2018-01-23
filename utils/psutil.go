package utils

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"runtime"
)

func MemoryUsed() (used float64, err error) {
	virt, err := mem.VirtualMemory()
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "utils: Failed to read virtual memory"),
		}
		return
	}

	used = ToFixed(virt.UsedPercent, 2)

	return
}

type LoadStat struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

func LoadAverage() (ld *LoadStat, err error) {
	count := float64(runtime.NumCPU())

	avg, err := load.Avg()
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "utils: Failed to read load average"),
		}
		return
	}

	ld = &LoadStat{
		Load1:  ToFixed(avg.Load1/count*100, 2),
		Load5:  ToFixed(avg.Load5/count*100, 2),
		Load15: ToFixed(avg.Load15/count*100, 2),
	}

	return
}
