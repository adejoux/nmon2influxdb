// nmon2influxdb
// import HMC data in InfluxDB
// author: adejoux@djouxtech.net

package hmc

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/adejoux/nmon2influxdb/nmon2influxdblib"
	"github.com/codegangsta/cli"
)

const timeFormat = "2006-01-02T15:04:05-0700"

//Import is the entry point for subcommand hmc
func Import(c *cli.Context) {
	// parsing parameters
	//config := nmon2influxdblib.ParseParameters(c)
	//new hmc session

	hmc := NewHMC(c)

	if hmc.Samples > 0 {
		fmt.Printf("Fetching %d latest samples. 30 seconds interval.\n", hmc.Samples)
	} else {
		fmt.Printf("Fetching latest 2 hours performance metrics. See hmc_samples parameter.\n")
	}

	fmt.Printf("Getting list of managed systems\n")
	systems, GetSysErr := hmc.GetManagedSystems()
	nmon2influxdblib.CheckError(GetSysErr)

	for _, system := range systems {
		if len(hmc.FilterManagedSystem) > 0 {
			if hmc.FilterManagedSystem != system.Name {
				fmt.Printf("\nSkipping system: %s\n", system.Name)
				continue
			}
		}

		//set parameters common to all points in GlobalPoint
		hmc.GlobalPoint.System = system.Name

		fmt.Printf("\n%s\n", system.Name)
		pcmlinks, getPCMErr := hmc.GetSystemPCMLinks(system.UUID)
		if getPCMErr != nil {
			fmt.Printf("Error getting PCM data\n")
			continue
		}

		// Get Managed System PCM metrics
		data, err := hmc.GetPCMData(pcmlinks.System)
		nmon2influxdblib.CheckError(err)
		for _, sample := range data.SystemUtil.UtilSamples {

			timestamp, timeErr := time.Parse(timeFormat, sample.SampleInfo.TimeStamp)
			nmon2influxdblib.CheckError(timeErr)

			//Set timestamp common to all this points
			hmc.GlobalPoint.Timestamp = timestamp

			hmc.AddPoint(Point{Name: "SystemProcessor",
				Metric: "TotalProcUnits",
				Value:  sample.ServerUtil.Processor.TotalProcUnits[0]})
			hmc.AddPoint(Point{Name: "SystemProcessor",
				Metric: "UtilizedProcUnits",
				Value:  sample.ServerUtil.Processor.UtilizedProcUnits[0]})
			hmc.AddPoint(Point{Name: "SystemProcessor",
				Metric: "availableProcUnits",
				Value:  sample.ServerUtil.Processor.AvailableProcUnits[0]})
			hmc.AddPoint(Point{Name: "SystemProcessor",
				Metric: "configurableProcUnits",
				Value:  sample.ServerUtil.Processor.ConfigurableProcUnits[0]})

			hmc.AddPoint(Point{Name: "SystemMemory",
				Metric: "TotalMem",
				Value:  sample.ServerUtil.Memory.TotalMem[0]})
			hmc.AddPoint(Point{Name: "SystemMemory",
				Metric: "assignedMemToLpars",
				Value:  sample.ServerUtil.Memory.AssignedMemToLpars[0]})
			hmc.AddPoint(Point{Name: "SystemMemory",
				Metric: "availableMem",
				Value:  sample.ServerUtil.Memory.AvailableMem[0]})
			hmc.AddPoint(Point{Name: "SystemMemory",
				Metric: "ConfigurableMem",
				Value:  sample.ServerUtil.Memory.ConfigurableMem[0]})

			for _, spp := range sample.ServerUtil.SharedProcessorPool {
				hmc.GlobalPoint.Pool = spp.Name
				hmc.AddPoint(Point{Name: "SystemSharedProcessorPool",
					Metric: "assignedProcUnits",
					Value:  spp.AssignedProcUnits[0]})
				hmc.AddPoint(Point{Name: "SystemSharedProcessorPool",
					Metric: "utilizedProcUnits",
					Value:  spp.UtilizedProcUnits[0]})
				hmc.AddPoint(Point{Name: "SystemSharedProcessorPool",
					Metric: "availableProcUnits",
					Value:  spp.AvailableProcUnits[0]})
				hmc.GlobalPoint.Pool = ""
			}
			for _, vios := range sample.ViosUtil {
				hmc.GlobalPoint.Partition = vios.Name
				for _, scsi := range vios.Storage.GenericPhysicalAdapters {
					hmc.GlobalPoint.Device = scsi.ID
					hmc.AddPoint(Point{Name: "SystemgenericPhysicalAdapters",
						Metric: "transmittedBytes",
						Value:  scsi.TransmittedBytes[0]})
					hmc.AddPoint(Point{Name: "SystemGenericPhysicalAdapters",
						Metric: "numOfReads",
						Value:  scsi.NumOfReads[0]})
					hmc.AddPoint(Point{Name: "SystemGenericPhysicalAdapters",
						Metric: "numOfWrites",
						Value:  scsi.NumOfWrites[0]})
					hmc.AddPoint(Point{Name: "SystemGenericPhysicalAdapters",
						Metric: "readBytes",
						Value:  scsi.ReadBytes[0]})
					hmc.AddPoint(Point{Name: "SystemGenericPhysicalAdapters",
						Metric: "writeBytes",
						Value:  scsi.WriteBytes[0]})
					hmc.GlobalPoint.Device = ""
				}
				for _, fc := range vios.Storage.FiberChannelAdapters {
					hmc.GlobalPoint.Device = fc.ID
					hmc.AddPoint(Point{Name: "SystemFiberChannelAdapters",
						Metric: "transmittedBytes",
						Value:  fc.TransmittedBytes[0]})
					hmc.AddPoint(Point{Name: "SystemFiberChannelAdapters",
						Metric: "numOfReads",
						Value:  fc.NumOfReads[0]})
					hmc.AddPoint(Point{Name: "SystemFiberChannelAdapters",
						Metric: "numOfWrites",
						Value:  fc.NumOfWrites[0]})
					hmc.AddPoint(Point{Name: "SystemFiberChannelAdapters",
						Metric: "readBytes",
						Value:  fc.ReadBytes[0]})
					hmc.AddPoint(Point{Name: "SystemFiberChannelAdapters",
						Metric: "writeBytes",
						Value:  fc.WriteBytes[0]})
					hmc.GlobalPoint.Device = ""
				}
				for _, vscsi := range vios.Storage.GenericVirtualAdapters {
					hmc.GlobalPoint.Device = vscsi.ID
					hmc.AddPoint(Point{Name: "SystemGenericVirtualAdapters",
						Metric: "transmittedBytes",
						Value:  vscsi.TransmittedBytes[0]})
					hmc.AddPoint(Point{Name: "SystemGenericVirtualAdapters",
						Metric: "numOfReads",
						Value:  vscsi.NumOfReads[0]})
					hmc.AddPoint(Point{Name: "SystemGenericVirtualAdapters",
						Metric: "numOfWrites",
						Value:  vscsi.NumOfWrites[0]})
					hmc.AddPoint(Point{Name: "SystemGenericVirtualAdapters",
						Metric: "readBytes",
						Value:  vscsi.ReadBytes[0]})
					hmc.AddPoint(Point{Name: "SystemGenericVirtualAdapters",
						Metric: "writeBytes",
						Value:  vscsi.WriteBytes[0]})
					hmc.GlobalPoint.Device = ""
				}
				for _, ssp := range vios.Storage.SharedStoragePools {
					hmc.GlobalPoint.Pool = ssp.ID
					hmc.AddPoint(Point{Name: "SystemSharedStoragePool",
						Metric: "transmittedBytes",
						Value:  ssp.TransmittedBytes[0]})
					hmc.AddPoint(Point{Name: "SystemSharedStoragePool",
						Metric: "totalSpace",
						Value:  ssp.TotalSpace[0]})
					hmc.AddPoint(Point{Name: "SystemSharedStoragePool",
						Metric: "usedSpace",
						Value:  ssp.UsedSpace[0]})
					hmc.AddPoint(Point{Name: "SystemSharedStoragePool",
						Metric: "numOfReads",
						Value:  ssp.NumOfReads[0]})
					hmc.AddPoint(Point{Name: "SystemSharedStoragePool",
						Metric: "numOfWrites",
						Value:  ssp.NumOfWrites[0]})
					hmc.AddPoint(Point{Name: "SystemSharedStoragePool",
						Metric: "readBytes",
						Value:  ssp.ReadBytes[0]})
					hmc.AddPoint(Point{Name: "SystemSharedStoragePool",
						Metric: "writeBytes",
						Value:  ssp.WriteBytes[0]})
					hmc.GlobalPoint.Pool = ""
				}
				for _, net := range vios.Network.GenericAdapters {
					hmc.GlobalPoint.Device = net.ID
					hmc.GlobalPoint.Type = net.Type
					hmc.AddPoint(Point{Name: "SystemGenericAdapters",
						Metric: "transferredBytes",
						Value:  net.TransferredBytes[0]})
					hmc.AddPoint(Point{Name: "SystemGenericAdapters",
						Metric: "receivedPackets",
						Value:  net.ReceivedPackets[0]})
					hmc.AddPoint(Point{Name: "SystemGenericAdapters",
						Metric: "sentPackets",
						Value:  net.SentPackets[0]})
					hmc.AddPoint(Point{Name: "SystemGenericAdapters",
						Metric: "droppedPackets",
						Value:  net.DroppedPackets[0]})
					hmc.AddPoint(Point{Name: "SystemGenericAdapters",
						Metric: "sentBytes",
						Value:  net.SentBytes[0]})
					hmc.AddPoint(Point{Name: "SystemGenericAdapters",
						Metric: "ReceivedBytes",
						Value:  net.ReceivedBytes[0]})
					hmc.GlobalPoint.Device = ""
					hmc.GlobalPoint.Type = ""
				}

				for _, net := range vios.Network.SharedAdapters {
					hmc.GlobalPoint.Device = net.ID
					hmc.GlobalPoint.Type = net.Type
					hmc.AddPoint(Point{Name: "SystemSharedAdapters",
						Metric: "transferredBytes",
						Value:  net.TransferredBytes[0]})
					hmc.AddPoint(Point{Name: "SystemSharedAdapters",
						Metric: "receivedPackets",
						Value:  net.ReceivedPackets[0]})
					hmc.AddPoint(Point{Name: "SystemSharedAdapters",
						Metric: "sentPackets",
						Value:  net.SentPackets[0]})
					hmc.AddPoint(Point{Name: "SystemSharedAdapters",
						Metric: "droppedPackets",
						Value:  net.DroppedPackets[0]})
					hmc.AddPoint(Point{Name: "SystemSharedAdapters",
						Metric: "sentBytes",
						Value:  net.SentBytes[0]})
					hmc.AddPoint(Point{Name: "SystemSharedAdapters",
						Metric: "ReceivedBytes",
						Value:  net.ReceivedBytes[0]})
					hmc.GlobalPoint.Device = ""
					hmc.GlobalPoint.Type = ""
				}
			}

		}
		fmt.Printf("managed system %25s: %8d points fetched.\n", system.Name, hmc.InfluxDB.PointsCount())
		hmc.WritePoints()
		if hmc.ManagedSystemOnly {
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
				hmc.GlobalPoint = Point{System: system.Name}
				lparData, getErr := hmc.GetPCMData(lparLink)
				nmon2influxdblib.CheckError(getErr)
				fmt.Printf("Partition %30s:", lparData.SystemUtil.UtilSamples[0].LparsUtil[0].Name)
				for _, sample := range lparData.SystemUtil.UtilSamples {
					timestamp, timeErr := time.Parse(timeFormat, sample.SampleInfo.TimeStamp)
					nmon2influxdblib.CheckError(timeErr)
					//Set timestamp common to all this points
					hmc.GlobalPoint.Timestamp = timestamp

					for _, lpar := range sample.LparsUtil {
						hmc.GlobalPoint.Partition = lpar.Name
						hmc.AddPoint(Point{Name: "PartitionProcessor",
							Metric: "MaxVirtualProcessors",
							Value:  lpar.Processor.MaxVirtualProcessors[0]})
						hmc.AddPoint(Point{Name: "PartitionProcessor",
							Metric: "MaxProcUnits",
							Value:  lpar.Processor.MaxProcUnits[0]})
						hmc.AddPoint(Point{Name: "PartitionProcessor",
							Metric: "EntitledProcUnits",
							Value:  lpar.Processor.EntitledProcUnits[0]})
						hmc.AddPoint(Point{Name: "PartitionProcessor",
							Metric: "UtilizedProcUnits",
							Value:  lpar.Processor.UtilizedProcUnits[0]})
						hmc.AddPoint(Point{Name: "PartitionProcessor",
							Metric: "UtilizedCappedProcUnits",
							Value:  lpar.Processor.UtilizedCappedProcUnits[0]})
						hmc.AddPoint(Point{Name: "PartitionProcessor",
							Metric: "UtilizedUncappedProcUnits",
							Value:  lpar.Processor.UtilizedUncappedProcUnits[0]})
						hmc.AddPoint(Point{Name: "PartitionProcessor",
							Metric: "IdleProcUnits",
							Value:  lpar.Processor.IdleProcUnits[0]})
						hmc.AddPoint(Point{Name: "PartitionProcessor",
							Metric: "DonatedProcUnits",
							Value:  lpar.Processor.DonatedProcUnits[0]})
						hmc.AddPoint(Point{Name: "PartitionProcessor",
							Metric: "TimeSpentWaitingForDispatch",
							Value:  lpar.Processor.TimeSpentWaitingForDispatch[0]})
						hmc.AddPoint(Point{Name: "PartitionProcessor",
							Metric: "TimePerInstructionExecution",
							Value:  lpar.Processor.TimePerInstructionExecution[0]})

						hmc.AddPoint(Point{Name: "PartitionMemory",
							Metric: "LogicalMem",
							Value:  lpar.Memory.LogicalMem[0]})
						hmc.AddPoint(Point{Name: "PartitionMemory",
							Metric: "BackedPhysicalMem",
							Value:  lpar.Memory.BackedPhysicalMem[0]})

						for _, vfc := range lpar.Storage.VirtualFiberChannelAdapters {
							hmc.GlobalPoint.WWPN = vfc.Wwpn
							hmc.GlobalPoint.PhysicalPortWWPN = vfc.PhysicalPortWWPN
							hmc.GlobalPoint.ViosID = strconv.Itoa(vfc.ViosID)
							hmc.AddPoint(Point{Name: "PartitionVirtualFiberChannelAdapters",
								Metric: "transmittedBytes",
								Value:  vfc.TransmittedBytes[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualFiberChannelAdapters",
								Metric: "numOfReads",
								Value:  vfc.NumOfReads[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualFiberChannelAdapters",
								Metric: "numOfWrites",
								Value:  vfc.NumOfWrites[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualFiberChannelAdapters",
								Metric: "readBytes",
								Value:  vfc.ReadBytes[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualFiberChannelAdapters",
								Metric: "writeBytes",
								Value:  vfc.WriteBytes[0]})
							hmc.GlobalPoint.WWPN = ""
							hmc.GlobalPoint.PhysicalPortWWPN = ""
							hmc.GlobalPoint.ViosID = ""
						}

						for _, vscsi := range lpar.Storage.GenericVirtualAdapters {
							hmc.GlobalPoint.Device = vscsi.ID
							hmc.GlobalPoint.ViosID = strconv.Itoa(vscsi.ViosID)
							hmc.AddPoint(Point{Name: "PartitionVSCSIAdapters",
								Metric: "transmittedBytes",
								Value:  vscsi.TransmittedBytes[0]})
							hmc.AddPoint(Point{Name: "PartitionVSCSIAdapters",
								Metric: "numOfReads",
								Value:  vscsi.NumOfReads[0]})
							hmc.AddPoint(Point{Name: "PartitionVSCSIAdapters",
								Metric: "numOfWrites",
								Value:  vscsi.NumOfWrites[0]})
							hmc.AddPoint(Point{Name: "PartitionVSCSIAdapters",
								Metric: "readBytes",
								Value:  vscsi.ReadBytes[0]})
							hmc.AddPoint(Point{Name: "PartitionVSCSIAdapters",
								Metric: "writeBytes",
								Value:  vscsi.WriteBytes[0]})
							hmc.GlobalPoint.Device = ""
							hmc.GlobalPoint.ViosID = ""
						}

						for _, net := range lpar.Network.VirtualEthernetAdapters {
							hmc.GlobalPoint.VlanID = strconv.Itoa(net.VlanID)
							hmc.GlobalPoint.VswitchID = strconv.Itoa(net.VswitchID)
							hmc.GlobalPoint.SharedEthernetAdapterID = net.SharedEthernetAdapterID
							hmc.GlobalPoint.ViosID = strconv.Itoa(net.ViosID)
							hmc.AddPoint(Point{Name: "PartitionVirtualEthernetAdapters",
								Metric: "transferredBytes",
								Value:  net.TransferredBytes[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualEthernetAdapters",
								Metric: "receivedPackets",
								Value:  net.ReceivedPackets[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualEthernetAdapters",
								Metric: "sentPackets",
								Value:  net.SentPackets[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualEthernetAdapters",
								Metric: "droppedPackets",
								Value:  net.DroppedPackets[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualEthernetAdapters",
								Metric: "sentBytes",
								Value:  net.SentBytes[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualEthernetAdapters",
								Metric: "ReceivedBytes",
								Value:  net.ReceivedBytes[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualEthernetAdapters",
								Metric: "transferredPhysicalBytes",
								Value:  net.TransferredPhysicalBytes[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualEthernetAdapters",
								Metric: "receivedPhysicalPackets",
								Value:  net.ReceivedPhysicalPackets[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualEthernetAdapters",
								Metric: "sentPhysicalPackets",
								Value:  net.SentPhysicalPackets[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualEthernetAdapters",
								Metric: "droppedPhysicalPackets",
								Value:  net.DroppedPhysicalPackets[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualEthernetAdapters",
								Metric: "sentPhysicalBytes",
								Value:  net.SentPhysicalBytes[0]})
							hmc.AddPoint(Point{Name: "PartitionVirtualEthernetAdapters",
								Metric: "ReceivedPhysicalBytes",
								Value:  net.ReceivedPhysicalBytes[0]})
							hmc.GlobalPoint.VlanID = ""
							hmc.GlobalPoint.VswitchID = ""
							hmc.GlobalPoint.SharedEthernetAdapterID = ""
							hmc.GlobalPoint.ViosID = ""
						}

						for _, net := range lpar.Network.SriovLogicalPorts {
							hmc.GlobalPoint.DrcIndex = net.DrcIndex
							hmc.GlobalPoint.PhysicalLocation = net.PhysicalLocation
							hmc.GlobalPoint.PhysicalDrcIndex = net.PhysicalDrcIndex
							hmc.GlobalPoint.PhysicalPortID = strconv.Itoa(net.PhysicalPortID)
							hmc.AddPoint(Point{Name: "PartitionSriovLogicalPorts",
								Metric: "receivedPackets",
								Value:  net.ReceivedPackets[0]})
							hmc.AddPoint(Point{Name: "PartitionSriovLogicalPorts",
								Metric: "sentPackets",
								Value:  net.SentPackets[0]})
							hmc.AddPoint(Point{Name: "PartitionSriovLogicalPorts",
								Metric: "droppedPackets",
								Value:  net.DroppedPackets[0]})
							hmc.AddPoint(Point{Name: "PartitionSriovLogicalPorts",
								Metric: "sentBytes",
								Value:  net.SentBytes[0]})
							hmc.AddPoint(Point{Name: "PartitionSriovLogicalPorts",
								Metric: "ReceivedBytes",
								Value:  net.ReceivedBytes[0]})

							hmc.GlobalPoint.DrcIndex = ""
							hmc.GlobalPoint.PhysicalLocation = ""
							hmc.GlobalPoint.PhysicalDrcIndex = ""
							hmc.GlobalPoint.PhysicalPortID = ""
						}
					}
				}
				fmt.Printf(" %8d points fetched.\n", hmc.InfluxDB.PointsCount())
				hmc.WritePoints()
			}
		}
	}
}
