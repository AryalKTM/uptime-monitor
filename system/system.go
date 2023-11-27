package system

import (
	"image/color"
	"io"
	"log"
	"time"
	"net/http"
	"github.com/jaypipes/ghw"
)

type HostInfo struct {
	Platform        string
	PlatformVersion string
	KenelArch       string
	Uptime          uint64
}

type CpuInfo struct {
	ModelName string
	Frequency float64
	Cores     int
	Usage     float64
}

type GpuInfo struct {
	Name   string
	Vendor string
}

type MemoryInfo struct {
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}

type NetworkInfo struct {
	InternetConnected bool
	PublicIP          string
}

type SystemInfo struct {
	Host    HostInfo
	Cpu     CpuInfo
	Gpu     GpuInfo
	Memory  MemoryInfo
	Network NetworkInfo
}

func PrintSystemInfo() {
	sysInfo, err := GetSystemInfo()
	if err != nil {
		log.Printf("[%s] %s \n", utilities.CreateColorString("Error", color.FgHiRed), err)
		return
	}
}

func GetSystemInfo(SystemInfo, error) {

	c, err := cpu.Info()
	if err != nil {
		log.Printf("[%s] %s \n", utilities.CreateColorString("Error", color.FgHiRed), err)
		return SystemInfo{}, err
	}

	h, err := host.Info()
	if err != nil {
		log.Printf("[%s] %s \n", utilities.CreateColorString("Error", color.FgHiRed), err)
		return SystemInfo{}, err
	}

	p, err := cpu.Percent(time.Millisecond*100, false)
	if err != nil {
		log.Printf("[%s] %s \n", utilities.CreateColorString("Error", color.FgHiRed), err)
		return SystemInfo{}, err
	}

	v, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("[%s] %s \n", utilities.CreateColorString("Error", color.FgHiRed), err)
		return	SystemInfo{}, err
	}

	_gpu := true
	gpu, err := ghw.GPU()
	if err != nil {
		log.Printf("[%s], %s \n", utilities.CreateColorString("Error", color.FgHiRed), err)
		_gpu = false
		return SystemInfo{}, err
	}

	_connected := false
	publicIp := ""
	resp, errHttp := http.Get("https://api.ipify.org")
	if errHttp != nil {
		_connected = false
	} else {
		_connected = true
		body, _ := io.ReadAll(resp.Body)
		publicIp = string(body)
	}

	gpuInf := GpuInfo {
		Name: "",
		Vendor: "",
	}

	if _gpu && len(gpu.GraphicsCards) > 0 && gpu.GraphicsCards[0].DeviceInfo != nil {
		gpuInf.Name = gpu.GraphicsCards[0].DeviceInfo.Product.Name
		gpuInf.Vendor = gpu.GraphicsCards[0].DeviceInfo.Vendor.Name
	}

	systemInfo := SystemInfo{
		Host: HostInfo{
			Platform:        h.Platform,
			PlatformVersion: h.PlatformVersion,
			KenelArch:       h.KenelArch,
			Uptime:          h.Uptime,
		},
		Cpu: CpuInfo{
			ModelName: c[0].ModelName,
			Frequency: c[0].Mhz,
			Cores:     len(c),
			Usage:     p[0],
		},
		Gpu: gpuInf,
		Memory: MemoryInfo{
			Total:       v.Total,
			Free:        v.Free,
			Used:        v.Total - v.Free,
			UsedPercent: v.UsedPercent,
		},
		Network: NetworkInfo{
			InternetConnected: _connected,
			PublicIP:          publicIp,
		},
	}
	return systemInfo
}
