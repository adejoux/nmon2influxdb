// nmon2influxdb
// import HMC data in InfluxDB
// author: adejoux@djouxtech.net

package hmc

// PCMData contains the json data structure
type PCMData struct {
	SystemUtil struct {
		UtilInfo struct {
			Version          string   `json:"version"`
			MetricType       string   `json:"metricType"`
			Frequency        int      `json:"frequency"`
			StartTimeStamp   string   `json:"startTimeStamp"`
			EndTimeStamp     string   `json:"endTimeStamp"`
			Mtms             string   `json:"mtms"`
			Name             string   `json:"name"`
			MetricArrayOrder []string `json:"metricArrayOrder"`
			UUID             string   `json:"uuid"`
		} `json:"utilInfo"`
		UtilSamples []struct {
			SampleType         string `json:"sampleType"`
			SystemFirmwareUtil struct {
				UtilizedProcUnits []interface{} `json:"utilizedProcUnits"`
				AssignedMem       []float64         `json:"assignedMem"`
			} `json:"systemFirmwareUtil"`
			ServerUtil struct {
				Processor struct {
					TotalProcUnits        []float64 `json:"totalProcUnits"`
					UtilizedProcUnits     []float64 `json:"utilizedProcUnits"`
					AvailableProcUnits    []float64 `json:"availableProcUnits"`
					ConfigurableProcUnits []float64 `json:"configurableProcUnits"`
				} `json:"processor"`
				Memory struct {
					TotalMem           []float64 `json:"totalMem"`
					AssignedMemToLpars []float64 `json:"assignedMemToLpars"`
					AvailableMem       []float64 `json:"availableMem"`
					ConfigurableMem    []float64 `json:"configurableMem"`
				} `json:"memory"`
				SharedProcessorPool []struct {
					AssignedProcUnits   []float64 `json:"assignedProcUnits"`
					UtilizedProcUnits   []float64 `json:"utilizedProcUnits"`
					AvailableProcUnits  []float64 `json:"availableProcUnits"`
					ConfiguredProcUnits []float64 `json:"configuredProcUnits"`
					BorrowedProcUnits   []float64 `json:"borrowedProcUnits"`
					ID                  int       `json:"id"`
					Name                string    `json:"name"`
				} `json:"sharedProcessorPool"`
				Network struct {
					Headapters []struct {
						DrcIndex      string `json:"drcIndex"`
						PhysicalPorts []struct {
							TransferredBytes []float64 `json:"transferredBytes"`
							ID               int       `json:"id"`
							PhysicalLocation string    `json:"physicalLocation"`
							ReceivedPackets  []float64 `json:"receivedPackets"`
							SentPackets      []float64 `json:"sentPackets"`
							DroppedPackets   []float64 `json:"droppedPackets"`
							SentBytes        []float64 `json:"sentBytes"`
							ReceivedBytes    []float64 `json:"receivedBytes"`
						} `json:"physicalPorts"`
					} `json:"headapters"`
				} `json:"network"`
			} `json:"serverUtil"`
			ViosUtil []struct {
				UUID   string `json:"uuid"`
				State  string `json:"state"`
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Memory struct {
					AssignedMem []float64 `json:"assignedMem"`
					UtilizedMem []float64 `json:"utilizedMem"`
				} `json:"memory"`
				Processor struct {
					PoolID                    int       `json:"poolId"`
					Weight                    int       `json:"weight"`
					Mode                      string    `json:"mode"`
					MaxVirtualProcessors      []float64     `json:"maxVirtualProcessors"`
					MaxProcUnits              []float64 `json:"maxProcUnits"`
					EntitledProcUnits         []float64 `json:"entitledProcUnits"`
					UtilizedProcUnits         []float64 `json:"utilizedProcUnits"`
					UtilizedCappedProcUnits   []float64 `json:"utilizedCappedProcUnits"`
					UtilizedUncappedProcUnits []float64 `json:"utilizedUncappedProcUnits"`
					IdleProcUnits             []float64 `json:"idleProcUnits"`
					DonatedProcUnits          []float64 `json:"donatedProcUnits"`
				} `json:"processor"`
				Network struct {
					GenericAdapters []struct {
						TransferredBytes []float64 `json:"transferredBytes"`
						Type             string    `json:"type"`
						ID               string    `json:"id"`
						PhysicalLocation string    `json:"physicalLocation"`
						ReceivedPackets  []float64 `json:"receivedPackets"`
						SentPackets      []float64 `json:"sentPackets"`
						DroppedPackets   []float64 `json:"droppedPackets"`
						SentBytes        []float64 `json:"sentBytes"`
						ReceivedBytes    []float64 `json:"receivedBytes"`
					} `json:"genericAdapters"`
					SharedAdapters []struct {
						TransferredBytes []float64 `json:"transferredBytes"`
						ID               string    `json:"id"`
						Type             string    `json:"type"`
						PhysicalLocation string    `json:"physicalLocation"`
						ReceivedPackets  []float64 `json:"receivedPackets"`
						SentPackets      []float64 `json:"sentPackets"`
						DroppedPackets   []float64 `json:"droppedPackets"`
						SentBytes        []float64 `json:"sentBytes"`
						ReceivedBytes    []float64 `json:"receivedBytes"`
						BridgedAdapters  []string  `json:"bridgedAdapters"`
					} `json:"sharedAdapters"`
				} `json:"network"`
				Storage struct {
					GenericPhysicalAdapters []struct {
						TransmittedBytes []float64 `json:"transmittedBytes"`
						Type             string    `json:"type"`
						ID               string    `json:"id"`
						PhysicalLocation string    `json:"physicalLocation"`
						NumOfReads       []float64 `json:"numOfReads"`
						NumOfWrites      []float64 `json:"numOfWrites"`
						ReadBytes        []float64 `json:"readBytes"`
						WriteBytes       []float64 `json:"writeBytes"`
					} `json:"genericPhysicalAdapters"`
					SharedStoragePools []struct {
						TransmittedBytes []float64 `json:"transmittedBytes"`
						ID               string    `json:"id"`
						TotalSpace       []float64     `json:"totalSpace"`
						UsedSpace        []float64     `json:"usedSpace"`
						NumOfReads       []float64 `json:"numOfReads"`
						NumOfWrites      []float64 `json:"numOfWrites"`
						ReadBytes        []float64 `json:"readBytes"`
						WriteBytes       []float64 `json:"writeBytes"`
					} `json:"sharedStoragePools"`
					FiberChannelAdapters []struct {
						TransmittedBytes []float64 `json:"transmittedBytes"`
						Wwpn             string    `json:"wwpn"`
						PhysicalLocation string    `json:"physicalLocation"`
						NumOfPorts       int       `json:"numOfPorts"`
						RunningSpeed     []float64     `json:"runningSpeed"`
						ID               string    `json:"id"`
						NumOfReads       []float64 `json:"numOfReads"`
						NumOfWrites      []float64 `json:"numOfWrites"`
						ReadBytes        []float64 `json:"readBytes"`
						WriteBytes       []float64 `json:"writeBytes"`
					} `json:"fiberChannelAdapters"`
					GenericVirtualAdapters []struct {
						TransmittedBytes []float64 `json:"transmittedBytes"`
						Type             string    `json:"type"`
						ID               string    `json:"id"`
						PhysicalLocation string    `json:"physicalLocation"`
						NumOfReads       []float64 `json:"numOfReads"`
						NumOfWrites      []float64 `json:"numOfWrites"`
						ReadBytes        []float64 `json:"readBytes"`
						WriteBytes       []float64 `json:"writeBytes"`
					} `json:"genericVirtualAdapters"`
				} `json:"storage"`
			} `json:"viosUtil"`
			LparsUtil []struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Type   string `json:"type"`
				Memory struct {
					LogicalMem        []float64 `json:"logicalMem"`
					BackedPhysicalMem []float64 `json:"backedPhysicalMem"`
				} `json:"memory"`
				Processor struct {
					PoolID                      int       `json:"poolId"`
					Weight                      int       `json:"weight"`
					Mode                        string    `json:"mode"`
					MaxVirtualProcessors        []float64 `json:"maxVirtualProcessors"`
					MaxProcUnits                []float64 `json:"maxProcUnits"`
					EntitledProcUnits           []float64 `json:"entitledProcUnits"`
					UtilizedProcUnits           []float64 `json:"utilizedProcUnits"`
					UtilizedCappedProcUnits     []float64 `json:"utilizedCappedProcUnits"`
					UtilizedUncappedProcUnits   []float64 `json:"utilizedUncappedProcUnits"`
					IdleProcUnits               []float64 `json:"idleProcUnits"`
					DonatedProcUnits            []float64 `json:"donatedProcUnits"`
					TimeSpentWaitingForDispatch []float64 `json:"timeSpentWaitingForDispatch"`
					TimePerInstructionExecution []float64 `json:"timePerInstructionExecution"`
				} `json:"processor"`
				Storage struct {
					VirtualFiberChannelAdapters []struct {
						TransmittedBytes []float64 `json:"transmittedBytes"`
						Wwpn             string    `json:"wwpn"`
						Wwpn2            string    `json:"wwpn2"`
						ViosID           int       `json:"viosId"`
						PhysicalLocation string    `json:"physicalLocation"`
						PhysicalPortWWPN string    `json:"physicalPortWWPN"`
						RunningSpeed     []float64     `json:"runningSpeed"`
						ID               string    `json:"id"`
						NumOfReads       []float64 `json:"numOfReads"`
						NumOfWrites      []float64 `json:"numOfWrites"`
						ReadBytes        []float64 `json:"readBytes"`
						WriteBytes       []float64 `json:"writeBytes"`
					} `json:"virtualFiberChannelAdapters"`
					GenericVirtualAdapters []struct {
						TransmittedBytes []float64 `json:"transmittedBytes"`
						Type             string    `json:"type"`
						ID               string    `json:"id"`
						ViosID           int       `json:"viosId"`
						PhysicalLocation string    `json:"physicalLocation"`
						NumOfReads       []float64 `json:"numOfReads"`
						NumOfWrites      []float64 `json:"numOfWrites"`
						ReadBytes        []float64 `json:"readBytes"`
						WriteBytes       []float64 `json:"writeBytes"`
					} `json:"genericVirtualAdapters"`
				} `json:"storage"`
				Network struct {
					VirtualEthernetAdapters []struct {
						TransferredPhysicalBytes []float64 `json:"transferredPhysicalBytes"`
						TransferredBytes         []float64 `json:"transferredBytes"`
						PhysicalLocation         string    `json:"physicalLocation"`
						VlanID                   int       `json:"vlanId"`
						VswitchID                int       `json:"vswitchId"`
						IsPortVlanID             bool      `json:"isPortVlanId"`
						ReceivedPackets          []float64 `json:"receivedPackets"`
						SentPackets              []float64 `json:"sentPackets"`
						DroppedPackets           []float64 `json:"droppedPackets"`
						SentBytes                []float64 `json:"sentBytes"`
						ReceivedBytes            []float64 `json:"receivedBytes"`
						ReceivedPhysicalPackets  []float64 `json:"receivedPhysicalPackets"`
						SentPhysicalPackets      []float64 `json:"sentPhysicalPackets"`
						DroppedPhysicalPackets   []float64 `json:"droppedPhysicalPackets"`
						SentPhysicalBytes        []float64 `json:"sentPhysicalBytes"`
						ReceivedPhysicalBytes    []float64 `json:"receivedPhysicalBytes"`
						ViosID                   int       `json:"viosId,omitempty"`
						SharedEthernetAdapterID  string    `json:"sharedEthernetAdapterId,omitempty"`
					} `json:"virtualEthernetAdapters"`
					SriovLogicalPorts []struct {
						DrcIndex         string    `json:"drcIndex"`
						PhysicalLocation string    `json:"physicalLocation"`
						PhysicalDrcIndex string    `json:"physicalDrcIndex"`
						PhysicalPortID   int       `json:"physicalPortId"`
						TransferredBytes []float64 `json:"transferredBytes"`
						VlanID           int       `json:"vlanId"`
						VswitchID        int       `json:"vswitchId"`
						IsPortVlanID     bool      `json:"isPortVlanId"`
						ReceivedPackets  []float64 `json:"receivedPackets"`
						SentPackets      []float64 `json:"sentPackets"`
						DroppedPackets   []float64 `json:"droppedPackets"`
						SentBytes        []float64 `json:"sentBytes"`
						ReceivedBytes    []float64 `json:"receivedBytes"`
						ViosID           int       `json:"viosId,omitempty"`
					} `json:"sriovLogicalPorts"`
				} `json:"network"`
				State string `json:"state"`
				UUID  string `json:"uuid"`
			} `json:"lparsUtil"`
			SampleInfo struct {
				TimeStamp string `json:"timeStamp"`
				Status    int    `json:"status"`
				ErrorInfo []struct {
					ErrID           string `json:"errId"`
					ErrMsg          string `json:"errMsg"`
					UUID            string `json:"uuid"`
					ReportedBy      string `json:"reportedBy"`
					OccurrenceCount int    `json:"occurrenceCount"`
				} `json:"errorInfo"`
			} `json:"sampleInfo"`
		} `json:"utilSamples"`
	} `json:"systemUtil"`
}
