package util

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
)

func GetEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func ApplyResourceLimits() {
    vcpu := GetEnvInt("SIM_VCPU", runtime.NumCPU())
    maxRAM := GetEnvInt("MAX_RAM", 0) // 0 = no limit
    
    runtime.GOMAXPROCS(vcpu)
    if maxRAM > 0 {
        debug.SetMemoryLimit(int64(maxRAM) * 1024 * 1024)
    }
    
    fmt.Printf("Resource limits - VCPU: %d, RAM: %dMB\n", vcpu, maxRAM)
}