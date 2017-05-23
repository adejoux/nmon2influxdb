// nmon2influxdb
// import HMC data in InfluxDB
// author: adejoux@djouxtech.net

package hmc

import (
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/adejoux/nmon2influxdb/nmon2influxdblib"
	"github.com/codegangsta/cli"
)

const timeFormat = "2006-01-02T15:04:05-0700"

//Import is the entry point for subcommand hmc
func Import(c *cli.Context) {
	//new hmc session
	hmc := NewHMC(c)

	if hmc.Samples > 0 {
		log.Printf("Fetching %d latest samples. 30 seconds interval.\n", hmc.Samples)
	} else {
		log.Printf("Fetching latest 2 hours performance metrics. See hmc_samples parameter.\n")
	}

	log.Printf("Getting list of managed systems\n")
	systems, GetSysErr := hmc.GetManagedSystems()
	nmon2influxdblib.CheckError(GetSysErr)

	for _, system := range systems {
		if len(hmc.FilterManagedSystem) > 0 {
			if hmc.FilterManagedSystem != system.Name {
				log.Printf("\nSkipping system: %s\n", system.Name)
				continue
			}
		}

		//set parameters common to all points in GlobalPoint
		hmc.GlobalPoint.System = system.Name

		log.Printf("MANAGED SYSTEM %s\n", strings.ToUpper(system.Name))
		pcmlinks, getPCMErr := hmc.GetSystemPCMLinks(system.UUID)
		if getPCMErr != nil {
			log.Printf("Error getting PCM data\n")
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

			// if sample status equal 1 we have no data in this sample
			if sample.SampleInfo.Status == 1 {
				log.Printf("Skipping sample. Error in sample collection: %s\n", sample.SampleInfo.ErrorInfo[0].ErrMsg)
				continue
			}

			hmc.AddPoint("SystemProcessor","TotalProcUnits",sample.ServerUtil.Processor.TotalProcUnits)
			hmc.AddPoint("SystemProcessor","UtilizedProcUnits",sample.ServerUtil.Processor.UtilizedProcUnits)
			hmc.AddPoint("SystemProcessor","availableProcUnits",sample.ServerUtil.Processor.AvailableProcUnits)
			hmc.AddPoint("SystemProcessor","configurableProcUnits",sample.ServerUtil.Processor.ConfigurableProcUnits)

			hmc.AddPoint("SystemMemory","TotalMem",sample.ServerUtil.Memory.TotalMem)
			hmc.AddPoint("SystemMemory","assignedMemToLpars",sample.ServerUtil.Memory.AssignedMemToLpars)
			hmc.AddPoint("SystemMemory","availableMem",sample.ServerUtil.Memory.AvailableMem)
			hmc.AddPoint("SystemMemory","ConfigurableMem",sample.ServerUtil.Memory.ConfigurableMem)

			for _, spp := range sample.ServerUtil.SharedProcessorPool {
				hmc.GlobalPoint.Pool = spp.Name
				hmc.AddPoint("SystemSharedProcessorPool","assignedProcUnits",spp.AssignedProcUnits)
				hmc.AddPoint("SystemSharedProcessorPool","utilizedProcUnits",spp.UtilizedProcUnits)
				hmc.AddPoint("SystemSharedProcessorPool","availableProcUnits",spp.AvailableProcUnits)
				hmc.GlobalPoint.Pool = ""
			}
			for _, vios := range sample.ViosUtil {
				hmc.GlobalPoint.Partition = vios.Name
				for _, scsi := range vios.Storage.GenericPhysicalAdapters {
					hmc.GlobalPoint.Device = scsi.ID
					hmc.AddPoint("SystemgenericPhysicalAdapters", "transmittedBytes",scsi.TransmittedBytes)
					hmc.AddPoint("SystemGenericPhysicalAdapters","numOfReads", scsi.NumOfReads)
					hmc.AddPoint("SystemGenericPhysicalAdapters","numOfWrites", scsi.NumOfWrites)
					hmc.AddPoint("SystemGenericPhysicalAdapters","readBytes", scsi.ReadBytes)
					hmc.AddPoint("SystemGenericPhysicalAdapters","writeBytes", scsi.WriteBytes)
					hmc.GlobalPoint.Device = ""
				}
				for _, fc := range vios.Storage.FiberChannelAdapters {
					hmc.GlobalPoint.Device = fc.ID
					if len(fc.TransmittedBytes) > 0 {
						hmc.AddPoint("SystemFiberChannelAdapters", "transmittedBytes",fc.TransmittedBytes)
				  }
					hmc.AddPoint("SystemFiberChannelAdapters","numOfReads", fc.NumOfReads)
					hmc.AddPoint("SystemFiberChannelAdapters","numOfWrites", fc.NumOfWrites)
					hmc.AddPoint("SystemFiberChannelAdapters","readBytes", fc.ReadBytes)
					hmc.AddPoint("SystemFiberChannelAdapters","writeBytes", fc.WriteBytes)
					hmc.GlobalPoint.Device = ""
				}
				for _, vscsi := range vios.Storage.GenericVirtualAdapters {
					hmc.GlobalPoint.Device = vscsi.ID
					if len(vscsi.TransmittedBytes) > 0 {
						hmc.AddPoint("SystemGenericVirtualAdapters", "transmittedBytes",vscsi.TransmittedBytes)
				  }
					hmc.AddPoint("SystemGenericVirtualAdapters","numOfReads", vscsi.NumOfReads)
					hmc.AddPoint("SystemGenericVirtualAdapters","numOfWrites", vscsi.NumOfWrites)
					hmc.AddPoint("SystemGenericVirtualAdapters","readBytes", vscsi.ReadBytes)
					hmc.AddPoint("SystemGenericVirtualAdapters","writeBytes", vscsi.WriteBytes)
					hmc.GlobalPoint.Device = ""
				}
				for _, ssp := range vios.Storage.SharedStoragePools {
					hmc.GlobalPoint.Pool = ssp.ID
					if len(ssp.TransmittedBytes) > 0 {
						hmc.AddPoint("SystemSharedStoragePool", "transmittedBytes",ssp.TransmittedBytes)
				  }
					hmc.AddPoint("SystemSharedStoragePool","totalSpace", ssp.TotalSpace)
					hmc.AddPoint("SystemSharedStoragePool","usedSpace", ssp.UsedSpace)
					hmc.AddPoint("SystemSharedStoragePool","numOfReads", ssp.NumOfReads)
					hmc.AddPoint("SystemSharedStoragePool","numOfWrites", ssp.NumOfWrites)
					hmc.AddPoint("SystemSharedStoragePool","readBytes", ssp.ReadBytes)
					hmc.AddPoint("SystemSharedStoragePool","writeBytes", ssp.WriteBytes)
					hmc.GlobalPoint.Pool = ""
				}
				for _, net := range vios.Network.GenericAdapters {
					hmc.GlobalPoint.Device = net.ID
					hmc.GlobalPoint.Type = net.Type
					if len(net.TransferredBytes) > 0 {
						hmc.AddPoint("SystemGenericAdapters", "transferredBytes",net.TransferredBytes)
				  }
					hmc.AddPoint("SystemGenericAdapters","receivedPackets", net.ReceivedPackets)
					hmc.AddPoint("SystemGenericAdapters","sentPackets", net.SentPackets)
					hmc.AddPoint("SystemGenericAdapters","droppedPackets", net.DroppedPackets)
					hmc.AddPoint("SystemGenericAdapters","sentBytes", net.SentBytes)
					hmc.AddPoint("SystemGenericAdapters","ReceivedBytes", net.ReceivedBytes)
					hmc.GlobalPoint.Device = ""
					hmc.GlobalPoint.Type = ""
				}

				for _, net := range vios.Network.SharedAdapters {
					hmc.GlobalPoint.Device = net.ID
					hmc.GlobalPoint.Type = net.Type
					if len(net.TransferredBytes) > 0 {
						hmc.AddPoint("SystemSharedAdapters", "transferredBytes",net.TransferredBytes)
				  }
					hmc.AddPoint("SystemSharedAdapters","receivedPackets", net.ReceivedPackets)
					hmc.AddPoint("SystemSharedAdapters","sentPackets", net.SentPackets)
					hmc.AddPoint("SystemSharedAdapters","droppedPackets", net.DroppedPackets)
					hmc.AddPoint("SystemSharedAdapters","sentBytes", net.SentBytes)
					hmc.AddPoint("SystemSharedAdapters","ReceivedBytes", net.ReceivedBytes)
					hmc.GlobalPoint.Device = ""
					hmc.GlobalPoint.Type = ""
				}
			}

		}
		log.Printf("managed system metrics: %8d points fetched.\n", hmc.InfluxDB.PointsCount())
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
				log.Println(lparGetPCMErr)
				log.Printf("Error getting PCM data\n")
				continue
			}

			for _, lparLink := range lparLinks.Partitions {
				hmc.GlobalPoint = Point{System: system.Name}
				lparData, getErr := hmc.GetPCMData(lparLink)
				nmon2influxdblib.CheckError(getErr)

				for _, sample := range lparData.SystemUtil.UtilSamples {
					// if sample status equal 1 we have no data in this sample
					if sample.SampleInfo.Status == 1 {
						log.Printf("Skipping sample. Error in sample collection: %s\n", sample.SampleInfo.ErrorInfo[0].ErrMsg)
						continue
					}

					timestamp, timeErr := time.Parse(timeFormat, sample.SampleInfo.TimeStamp)
					nmon2influxdblib.CheckError(timeErr)
					//Set timestamp common to all this points
					hmc.GlobalPoint.Timestamp = timestamp

					for _, lpar := range sample.LparsUtil {
						hmc.GlobalPoint.Partition = lpar.Name
						hmc.AddPoint("PartitionProcessor", "MaxVirtualProcessors",lpar.Processor.MaxVirtualProcessors)
						hmc.AddPoint("PartitionProcessor", "MaxProcUnits",lpar.Processor.MaxProcUnits)
						hmc.AddPoint("PartitionProcessor", "EntitledProcUnits",lpar.Processor.EntitledProcUnits)
						hmc.AddPoint("PartitionProcessor", "UtilizedProcUnits",lpar.Processor.UtilizedProcUnits)
						hmc.AddPoint("PartitionProcessor", "UtilizedCappedProcUnits",lpar.Processor.UtilizedCappedProcUnits)
						hmc.AddPoint("PartitionProcessor", "UtilizedUncappedProcUnits",lpar.Processor.UtilizedUncappedProcUnits)
						hmc.AddPoint("PartitionProcessor", "IdleProcUnits",lpar.Processor.IdleProcUnits)
						hmc.AddPoint("PartitionProcessor", "DonatedProcUnits",lpar.Processor.DonatedProcUnits)
						hmc.AddPoint("PartitionProcessor","TimeSpentWaitingForDispatch", lpar.Processor.TimeSpentWaitingForDispatch)
						hmc.AddPoint("PartitionProcessor", "TimePerInstructionExecution",lpar.Processor.TimePerInstructionExecution)
						hmc.AddPoint("PartitionMemory", "LogicalMem",lpar.Memory.LogicalMem)
						hmc.AddPoint("PartitionMemory", "BackedPhysicalMem",lpar.Memory.BackedPhysicalMem)

						for _, vfc := range lpar.Storage.VirtualFiberChannelAdapters {
							hmc.GlobalPoint.WWPN = vfc.Wwpn
							hmc.GlobalPoint.PhysicalPortWWPN = vfc.PhysicalPortWWPN
							hmc.GlobalPoint.ViosID = strconv.Itoa(vfc.ViosID)

							hmc.AddPoint("PartitionVirtualFiberChannelAdapters","transmittedBytes",vfc.TransmittedBytes)
							hmc.AddPoint("PartitionVirtualFiberChannelAdapters","numOfReads", vfc.NumOfReads)
							hmc.AddPoint("PartitionVirtualFiberChannelAdapters","numOfWrites", vfc.NumOfWrites)
							hmc.AddPoint("PartitionVirtualFiberChannelAdapters","readBytes", vfc.ReadBytes)
							hmc.AddPoint("PartitionVirtualFiberChannelAdapters","writeBytes", vfc.WriteBytes)
							hmc.GlobalPoint.WWPN = ""
							hmc.GlobalPoint.PhysicalPortWWPN = ""
							hmc.GlobalPoint.ViosID = ""
						}

						for _, vscsi := range lpar.Storage.GenericVirtualAdapters {
							hmc.GlobalPoint.Device = vscsi.ID
							hmc.GlobalPoint.ViosID = strconv.Itoa(vscsi.ViosID)

							hmc.AddPoint("PartitionVSCSIAdapters","transmittedBytes",vscsi.TransmittedBytes)
							hmc.AddPoint("PartitionVSCSIAdapters","numOfReads", vscsi.NumOfReads)
							hmc.AddPoint("PartitionVSCSIAdapters","numOfWrites", vscsi.NumOfWrites)
							hmc.AddPoint("PartitionVSCSIAdapters","readBytes", vscsi.ReadBytes)
							hmc.AddPoint("PartitionVSCSIAdapters","writeBytes", vscsi.WriteBytes)
							hmc.GlobalPoint.Device = ""
							hmc.GlobalPoint.ViosID = ""
						}

						for _, net := range lpar.Network.VirtualEthernetAdapters {
							hmc.GlobalPoint.VlanID = strconv.Itoa(net.VlanID)
							hmc.GlobalPoint.VswitchID = strconv.Itoa(net.VswitchID)
							hmc.GlobalPoint.SharedEthernetAdapterID = net.SharedEthernetAdapterID
							hmc.GlobalPoint.ViosID = strconv.Itoa(net.ViosID)
							hmc.AddPoint("PartitionVirtualEthernetAdapters", "transferredBytes", net.TransferredBytes)
							hmc.AddPoint("PartitionVirtualEthernetAdapters","receivedPackets", net.ReceivedPackets)
							hmc.AddPoint("PartitionVirtualEthernetAdapters","sentPackets", net.SentPackets)
							hmc.AddPoint("PartitionVirtualEthernetAdapters","droppedPackets", net.DroppedPackets)
							hmc.AddPoint("PartitionVirtualEthernetAdapters","sentBytes", net.SentBytes)
							hmc.AddPoint("PartitionVirtualEthernetAdapters","ReceivedBytes", net.ReceivedBytes)
							hmc.AddPoint("PartitionVirtualEthernetAdapters","transferredPhysicalBytes", net.TransferredPhysicalBytes)
							hmc.AddPoint("PartitionVirtualEthernetAdapters","receivedPhysicalPackets", net.ReceivedPhysicalPackets)
							hmc.AddPoint("PartitionVirtualEthernetAdapters","sentPhysicalPackets", net.SentPhysicalPackets)
							hmc.AddPoint("PartitionVirtualEthernetAdapters","droppedPhysicalPackets", net.DroppedPhysicalPackets)
							hmc.AddPoint("PartitionVirtualEthernetAdapters","sentPhysicalBytes", net.SentPhysicalBytes)
							hmc.AddPoint("PartitionVirtualEthernetAdapters","ReceivedPhysicalBytes", net.ReceivedPhysicalBytes)
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
							hmc.AddPoint("PartitionSriovLogicalPorts","receivedPackets", net.ReceivedPackets)
							hmc.AddPoint("PartitionSriovLogicalPorts","sentPackets", net.SentPackets)
							hmc.AddPoint("PartitionSriovLogicalPorts","droppedPackets", net.DroppedPackets)
							hmc.AddPoint("PartitionSriovLogicalPorts","sentBytes", net.SentBytes)
							hmc.AddPoint("PartitionSriovLogicalPorts","ReceivedBytes", net.ReceivedBytes)

							hmc.GlobalPoint.DrcIndex = ""
							hmc.GlobalPoint.PhysicalLocation = ""
							hmc.GlobalPoint.PhysicalDrcIndex = ""
							hmc.GlobalPoint.PhysicalPortID = ""
						}
					}
				}
				log.Printf("Partition %25s: %8d points fetched.\n", lparData.SystemUtil.UtilSamples[0].LparsUtil[0].Name, hmc.InfluxDB.PointsCount())
				hmc.WritePoints()
			}
		}
	}
}
