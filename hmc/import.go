// nmon2influxdb
// import HMC data in InfluxDB
// author: adejoux@djouxtech.net

package hmc

import (
	"fmt"
	"net/url"
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

	systems, GetSysErr := hmc.GetManagedSystems()
	nmon2influxdblib.CheckError(GetSysErr)

	for _, system := range systems {
		if len(hmc.FilterManagedSystem) > 0 {
			if hmc.FilterManagedSystem != system.Name {
				fmt.Printf("Skipping system: %s\n", system.Name)
				continue
			}
		}

		//set parameters common to all points in GlobalPoint
		hmc.GlobalPoint.Server = system.Name

		fmt.Printf("MANAGED SYSTEM: %s\n", system.Name)
		pcmlinks, getPCMErr := hmc.GetSystemPCMLinks(system.UUID)
		if getPCMErr != nil {
			fmt.Printf("Error getting PCM data\n")
			continue
		}

		var lparLinks PCMLinks
		for _, link := range pcmlinks.Partitions {
			//need to parse the link because the specified hostname can be different
			//of the one specified by the user and the auth cookie will not match
			rawurl, _ := url.Parse(link)
			var lparGetPCMErr error
			lparLinks, lparGetPCMErr = hmc.GetPartitionPCMLinks(rawurl.Path)
			if lparGetPCMErr != nil {
				fmt.Println(lparGetPCMErr)
				fmt.Printf("Error getting PCM data\n")
				continue
			}

			for _, lparLink := range lparLinks.Partitions {
				hmc.GlobalPoint = Point{Server: system.Name}
				lparData, getErr := hmc.GetPCMData(lparLink)
				nmon2influxdblib.CheckError(getErr)
				fmt.Printf("partition %s:", lparData.SystemUtil.UtilSamples[0].LparsUtil[0].Name)
				for _, sample := range lparData.SystemUtil.UtilSamples {
					timestamp, timeErr := time.Parse("2006-01-02T15:04:05+0000", sample.SampleInfo.TimeStamp)
					nmon2influxdblib.CheckError(timeErr)
					//Set timestamp common to all this points
					hmc.GlobalPoint.Timestamp = timestamp

					for _, lpar := range sample.LparsUtil {
						hmc.GlobalPoint.Partition = lpar.Name
						hmc.AddPoint(Point{Name: "processor",
							Metric: "MaxVirtualProcessors",
							Value:  lpar.Processor.MaxVirtualProcessors[0]})
						hmc.AddPoint(Point{Name: "processor",
							Metric: "MaxProcUnits",
							Value:  lpar.Processor.MaxProcUnits[0]})
						hmc.AddPoint(Point{Name: "processor",
							Metric: "EntitledProcUnits",
							Value:  lpar.Processor.EntitledProcUnits[0]})
						hmc.AddPoint(Point{Name: "processor",
							Metric: "UtilizedProcUnits",
							Value:  lpar.Processor.UtilizedProcUnits[0]})
						hmc.AddPoint(Point{Name: "processor",
							Metric: "UtilizedCappedProcUnits",
							Value:  lpar.Processor.UtilizedCappedProcUnits[0]})
						hmc.AddPoint(Point{Name: "processor",
							Metric: "UtilizedUncappedProcUnits",
							Value:  lpar.Processor.UtilizedUncappedProcUnits[0]})
						hmc.AddPoint(Point{Name: "processor",
							Metric: "IdleProcUnits",
							Value:  lpar.Processor.IdleProcUnits[0]})
						hmc.AddPoint(Point{Name: "processor",
							Metric: "DonatedProcUnits",
							Value:  lpar.Processor.DonatedProcUnits[0]})
						hmc.AddPoint(Point{Name: "processor",
							Metric: "TimeSpentWaitingForDispatch",
							Value:  lpar.Processor.TimeSpentWaitingForDispatch[0]})
						hmc.AddPoint(Point{Name: "processor",
							Metric: "TimePerInstructionExecution",
							Value:  lpar.Processor.TimePerInstructionExecution[0]})

						hmc.AddPoint(Point{Name: "memory",
							Metric: "LogicalMem",
							Value:  lpar.Memory.LogicalMem[0]})
						hmc.AddPoint(Point{Name: "memory",
							Metric: "BackedPhysicalMem",
							Value:  lpar.Memory.BackedPhysicalMem[0]})
					}

				}
				fmt.Printf(" %d points\n", hmc.InfluxDB.PointsCount())
				hmc.WritePoints()
			}
		}
		// Get Managed System PCM metrics
		data, err := hmc.GetPCMData(pcmlinks.System)
		nmon2influxdblib.CheckError(err)
		for _, sample := range data.SystemUtil.UtilSamples {

			timestamp, timeErr := time.Parse("2006-01-02T15:04:05+0000", sample.SampleInfo.TimeStamp)
			nmon2influxdblib.CheckError(timeErr)

			//Set timestamp common to all this points
			hmc.GlobalPoint.Timestamp = timestamp

			hmc.AddPoint(Point{Name: "processor",
				Metric: "TotalProcUnits",
				Value:  sample.ServerUtil.Processor.TotalProcUnits[0]})
			hmc.AddPoint(Point{Name: "processor",
				Metric: "UtilizedProcUnits",
				Value:  sample.ServerUtil.Processor.UtilizedProcUnits[0]})
			hmc.AddPoint(Point{Name: "processor",
				Metric: "availableProcUnits",
				Value:  sample.ServerUtil.Processor.AvailableProcUnits[0]})
			hmc.AddPoint(Point{Name: "processor",
				Metric: "configurableProcUnits",
				Value:  sample.ServerUtil.Processor.ConfigurableProcUnits[0]})

			hmc.AddPoint(Point{Name: "memory",
				Metric: "TotalMem",
				Value:  sample.ServerUtil.Memory.TotalMem[0]})
			hmc.AddPoint(Point{Name: "memory",
				Metric: "assignedMemToLpars",
				Value:  sample.ServerUtil.Memory.AssignedMemToLpars[0]})
			hmc.AddPoint(Point{Name: "memory",
				Metric: "availableMem",
				Value:  sample.ServerUtil.Memory.AvailableMem[0]})
			hmc.AddPoint(Point{Name: "memory",
				Metric: "ConfigurableMem",
				Value:  sample.ServerUtil.Memory.ConfigurableMem[0]})

			for _, spp := range sample.ServerUtil.SharedProcessorPool {
				hmc.AddPoint(Point{Name: "sharedProcessorPool",
					Metric: "assignedProcUnits",
					Value:  spp.AssignedProcUnits[0],
					Pool:   spp.Name})
				hmc.AddPoint(Point{Name: "sharedProcessorPool",
					Metric: "utilizedProcUnits",
					Pool:   spp.Name,
					Value:  spp.UtilizedProcUnits[0]})
				hmc.AddPoint(Point{Name: "sharedProcessorPool",
					Metric: "availableProcUnits",
					Value:  spp.AvailableProcUnits[0],
					Pool:   spp.Name})
			}
			for _, vios := range sample.ViosUtil {
				hmc.GlobalPoint.Partition = vios.Name
				for _, scsi := range vios.Storage.GenericPhysicalAdapters {
					hmc.AddPoint(Point{Name: "genericPhysicalAdapters",
						Metric: "transmittedBytes",
						Value:  scsi.TransmittedBytes[0],
						Device: scsi.ID})
					hmc.AddPoint(Point{Name: "genericPhysicalAdapters",
						Metric: "numOfReads",
						Value:  scsi.NumOfReads[0],
						Device: scsi.ID})
					hmc.AddPoint(Point{Name: "genericPhysicalAdapters",
						Metric: "numOfWrites",
						Value:  scsi.NumOfWrites[0],
						Device: scsi.ID})
					hmc.AddPoint(Point{Name: "genericPhysicalAdapters",
						Metric: "readBytes",
						Value:  scsi.ReadBytes[0],
						Device: scsi.ID})
					hmc.AddPoint(Point{Name: "genericPhysicalAdapters",
						Metric: "writeBytes",
						Value:  scsi.WriteBytes[0],
						Device: scsi.ID})
				}
				for _, fc := range vios.Storage.FiberChannelAdapters {
					hmc.AddPoint(Point{Name: "fiberChannelAdapters",
						Metric: "transmittedBytes",
						Value:  fc.TransmittedBytes[0],
						Device: fc.ID})
					hmc.AddPoint(Point{Name: "fiberChannelAdapters",
						Metric: "numOfReads",
						Value:  fc.NumOfReads[0],
						Device: fc.ID})
					hmc.AddPoint(Point{Name: "fiberChannelAdapters",
						Metric: "numOfWrites",
						Value:  fc.NumOfWrites[0],
						Device: fc.ID})
					hmc.AddPoint(Point{Name: "fiberChannelAdapters",
						Metric: "readBytes",
						Value:  fc.ReadBytes[0],
						Device: fc.ID})
					hmc.AddPoint(Point{Name: "fiberChannelAdapters",
						Metric: "writeBytes",
						Value:  fc.WriteBytes[0],
						Device: fc.ID})
				}
				for _, vscsi := range vios.Storage.GenericVirtualAdapters {
					hmc.AddPoint(Point{Name: "genericVirtualAdapters",
						Metric: "transmittedBytes",
						Value:  vscsi.TransmittedBytes[0],
						Device: vscsi.ID})
					hmc.AddPoint(Point{Name: "genericVirtualAdapters",
						Metric: "numOfReads",
						Value:  vscsi.NumOfReads[0],
						Device: vscsi.ID})
					hmc.AddPoint(Point{Name: "genericVirtualAdapters",
						Metric: "numOfWrites",
						Value:  vscsi.NumOfWrites[0],
						Device: vscsi.ID})
					hmc.AddPoint(Point{Name: "genericVirtualAdapters",
						Metric: "readBytes",
						Value:  vscsi.ReadBytes[0],
						Device: vscsi.ID})
					hmc.AddPoint(Point{Name: "genericVirtualAdapters",
						Metric: "writeBytes",
						Value:  vscsi.WriteBytes[0],
						Device: vscsi.ID})
				}
				for _, ssp := range vios.Storage.SharedStoragePools {
					hmc.AddPoint(Point{Name: "sharedStoragePool",
						Metric: "transmittedBytes",
						Value:  ssp.TransmittedBytes[0],
						Pool:   ssp.ID})
					hmc.AddPoint(Point{Name: "sharedStoragePool",
						Metric: "totalSpace",
						Value:  ssp.TotalSpace[0],
						Pool:   ssp.ID})
					hmc.AddPoint(Point{Name: "sharedStoragePool",
						Metric: "usedSpace",
						Value:  ssp.UsedSpace[0],
						Pool:   ssp.ID})
					hmc.AddPoint(Point{Name: "sharedStoragePool",
						Metric: "numOfReads",
						Value:  ssp.NumOfReads[0],
						Pool:   ssp.ID})
					hmc.AddPoint(Point{Name: "sharedStoragePool",
						Metric: "numOfWrites",
						Value:  ssp.NumOfWrites[0],
						Pool:   ssp.ID})
					hmc.AddPoint(Point{Name: "sharedStoragePool",
						Metric: "readBytes",
						Value:  ssp.ReadBytes[0],
						Pool:   ssp.ID})
					hmc.AddPoint(Point{Name: "sharedStoragePool",
						Metric: "writeBytes",
						Value:  ssp.WriteBytes[0],
						Pool:   ssp.ID})
				}
				for _, net := range vios.Network.GenericAdapters {
					hmc.AddPoint(Point{Name: "genericAdapters",
						Metric: "transferredBytes",
						Value:  net.TransferredBytes[0],
						Device: net.ID,
						Type:   net.Type})
					hmc.AddPoint(Point{Name: "genericAdapters",
						Metric: "receivedPackets",
						Value:  net.ReceivedPackets[0],
						Device: net.ID,
						Type:   net.Type})
					hmc.AddPoint(Point{Name: "genericAdapters",
						Metric: "sentPackets",
						Value:  net.SentPackets[0],
						Device: net.ID,
						Type:   net.Type})
					hmc.AddPoint(Point{Name: "genericAdapters",
						Metric: "droppedPackets",
						Value:  net.DroppedPackets[0],
						Device: net.ID,
						Type:   net.Type})
					hmc.AddPoint(Point{Name: "genericAdapters",
						Metric: "sentBytes",
						Value:  net.SentBytes[0],
						Device: net.ID,
						Type:   net.Type})
					hmc.AddPoint(Point{Name: "genericAdapters",
						Metric: "ReceivedBytes",
						Value:  net.ReceivedBytes[0],
						Device: net.ID,
						Type:   net.Type})
				}

				for _, net := range vios.Network.SharedAdapters {
					hmc.AddPoint(Point{Name: "sharedAdapters",
						Metric: "transferredBytes",
						Value:  net.TransferredBytes[0],
						Device: net.ID,
						Type:   net.Type})
					hmc.AddPoint(Point{Name: "sharedAdapters",
						Metric: "receivedPackets",
						Value:  net.ReceivedPackets[0],
						Device: net.ID,
						Type:   net.Type})
					hmc.AddPoint(Point{Name: "sharedAdapters",
						Metric: "sentPackets",
						Value:  net.SentPackets[0],
						Device: net.ID,
						Type:   net.Type})
					hmc.AddPoint(Point{Name: "sharedAdapters",
						Metric: "droppedPackets",
						Value:  net.DroppedPackets[0],
						Device: net.ID,
						Type:   net.Type})
					hmc.AddPoint(Point{Name: "sharedAdapters",
						Metric: "sentBytes",
						Value:  net.SentBytes[0],
						Device: net.ID,
						Type:   net.Type})
					hmc.AddPoint(Point{Name: "sharedAdapters",
						Metric: "ReceivedBytes",
						Value:  net.ReceivedBytes[0],
						Device: net.ID,
						Type:   net.Type})
				}

			}

		}
		hmc.WritePoints()
	}
}
