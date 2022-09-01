package alert

const (
	Low    = 1
	Medium = 5
	High   = 10
)

const (
	SystemOffline        = "system_offline"
	SystemCpuLevel       = "system_cpu_level"
	SystemMemoryLevel    = "system_memory_level"
	SystemSwapLevel      = "system_swap_level"
	SystemHugePagesLevel = "system_hugepages_level"
	SystemMdFailed       = "system_md_failed"
	DiskUsageLevel       = "disk_usage_level"
	KmsgKeyword          = "kmsg_keyword"
	CheckHttpFailed      = "check_http_failed"
)
