// nmon2influxdb
// import HMC data in InfluxDB
// author: adejoux@djouxtech.net

package hmc

import (
	"fmt"
	"time"

	"github.com/adejoux/nmon2influxdb/nmon2influxdblib"
	"github.com/codegangsta/cli"
)

//Import is the entry point for subcommand hmc
func Import(c *cli.Context) {
	// parsing parameters
	//config := nmon2influxdblib.ParseParameters(c)
	//new hmc session

	hmc := NewHMC(c)

	systems := hmc.GetManagedSystems()

	for _, system := range systems {
		fmt.Printf("Processing performance data from: %s\n", system.Name)

		pcmlinks, getPCMErr := hmc.GetPCMLinks(system.UUID)
		if getPCMErr != nil {
			fmt.Printf("Error getting PCM data\n")
			continue
		}
		for _, link := range pcmlinks {
			data, err := hmc.GetPCMData(link)
			nmon2influxdblib.CheckError(err)
			for _, sample := range data.SystemUtil.UtilSamples {

				timestamp, timeErr := time.Parse("2006-01-02T15:04:05+0000", sample.SampleInfo.TimeStamp)
				nmon2influxdblib.CheckError(timeErr)

				hmc.AddPoint(HMCPoint{Name: "processor",
					Server:    system.Name,
					Metric:    "TotalProcUnits",
					Value:     sample.ServerUtil.Processor.TotalProcUnits[0],
					Timestamp: timestamp})
				hmc.AddPoint(HMCPoint{Name: "processor",
					Server:    system.Name,
					Metric:    "UtilizedProcUnits",
					Value:     sample.ServerUtil.Processor.UtilizedProcUnits[0],
					Timestamp: timestamp})
				hmc.AddPoint(HMCPoint{Name: "processor",
					Server:    system.Name,
					Metric:    "availableProcUnits",
					Value:     sample.ServerUtil.Processor.AvailableProcUnits[0],
					Timestamp: timestamp})
				hmc.AddPoint(HMCPoint{Name: "processor",
					Server:    system.Name,
					Metric:    "configurableProcUnits",
					Value:     sample.ServerUtil.Processor.ConfigurableProcUnits[0],
					Timestamp: timestamp})

				hmc.AddPoint(HMCPoint{Name: "memory",
					Server:    system.Name,
					Metric:    "TotalMem",
					Value:     sample.ServerUtil.Memory.TotalMem[0],
					Timestamp: timestamp})
				hmc.AddPoint(HMCPoint{Name: "memory",
					Server:    system.Name,
					Metric:    "assignedMemToLpars",
					Value:     sample.ServerUtil.Memory.AssignedMemToLpars[0],
					Timestamp: timestamp})
				hmc.AddPoint(HMCPoint{Name: "memory",
					Server:    system.Name,
					Metric:    "availableMem",
					Value:     sample.ServerUtil.Memory.AvailableMem[0],
					Timestamp: timestamp})
				hmc.AddPoint(HMCPoint{Name: "memory",
					Server:    system.Name,
					Metric:    "ConfigurableMem",
					Value:     sample.ServerUtil.Memory.ConfigurableMem[0],
					Timestamp: timestamp})

				for _, spp := range sample.ServerUtil.SharedProcessorPool {
					hmc.AddPoint(HMCPoint{Name: "sharedProcessorPool",
						Server:    system.Name,
						Metric:    "assignedProcUnits",
						Value:     spp.AssignedProcUnits[0],
						Pool:      spp.Name,
						Timestamp: timestamp})
					hmc.AddPoint(HMCPoint{Name: "sharedProcessorPool",
						Server:    system.Name,
						Metric:    "utilizedProcUnits",
						Pool:      spp.Name,
						Value:     spp.UtilizedProcUnits[0],
						Timestamp: timestamp})
					hmc.AddPoint(HMCPoint{Name: "sharedProcessorPool",
						Server:    system.Name,
						Metric:    "availableProcUnits",
						Value:     spp.AvailableProcUnits[0],
						Pool:      spp.Name,
						Timestamp: timestamp})
					hmc.AddPoint(HMCPoint{Name: "sharedProcessorPool",
						Server:    system.Name,
						Metric:    "configuredProcUnits",
						Value:     spp.ConfiguredProcUnits[0],
						Pool:      spp.Name,
						Timestamp: timestamp})
					hmc.AddPoint(HMCPoint{Name: "sharedProcessorPool",
						Server:    system.Name,
						Metric:    "borrowedProcUnits",
						Value:     spp.BorrowedProcUnits[0],
						Pool:      spp.Name,
						Timestamp: timestamp})
				}
			}
		}
		hmc.WritePoints()
	}
}
