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
				AssignedMem       []int         `json:"assignedMem"`
			} `json:"systemFirmwareUtil"`
			ServerUtil struct {
				Processor struct {
					TotalProcUnits        []int     `json:"totalProcUnits"`
					UtilizedProcUnits     []float64 `json:"utilizedProcUnits"`
					AvailableProcUnits    []float64 `json:"availableProcUnits"`
					ConfigurableProcUnits []int     `json:"configurableProcUnits"`
				} `json:"processor"`
				Memory struct {
					TotalMem           []int `json:"totalMem"`
					AssignedMemToLpars []int `json:"assignedMemToLpars"`
					AvailableMem       []int `json:"availableMem"`
					ConfigurableMem    []int `json:"configurableMem"`
				} `json:"memory"`
				SharedProcessorPool []struct {
					AssignedProcUnits   []int     `json:"assignedProcUnits"`
					UtilizedProcUnits   []float64 `json:"utilizedProcUnits"`
					AvailableProcUnits  []int     `json:"availableProcUnits"`
					ConfiguredProcUnits []int     `json:"configuredProcUnits"`
					BorrowedProcUnits   []int     `json:"borrowedProcUnits"`
					ID                  int       `json:"id"`
					Name                string    `json:"name"`
				} `json:"sharedProcessorPool"`
				Network struct {
					Headapters []struct {
						DrcIndex      string `json:"drcIndex"`
						PhysicalPorts []struct {
							TransferredBytes []int  `json:"transferredBytes"`
							ID               int    `json:"id"`
							PhysicalLocation string `json:"physicalLocation"`
							ReceivedPackets  []int  `json:"receivedPackets"`
							SentPackets      []int  `json:"sentPackets"`
							DroppedPackets   []int  `json:"droppedPackets"`
							SentBytes        []int  `json:"sentBytes"`
							ReceivedBytes    []int  `json:"receivedBytes"`
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
					AssignedMem []int `json:"assignedMem"`
					UtilizedMem []int `json:"utilizedMem"`
				} `json:"memory"`
				Processor struct {
					PoolID                    int       `json:"poolId"`
					Weight                    int       `json:"weight"`
					Mode                      string    `json:"mode"`
					MaxVirtualProcessors      []int     `json:"maxVirtualProcessors"`
					MaxProcUnits              []int     `json:"maxProcUnits"`
					EntitledProcUnits         []int     `json:"entitledProcUnits"`
					UtilizedProcUnits         []float64 `json:"utilizedProcUnits"`
					UtilizedCappedProcUnits   []float64 `json:"utilizedCappedProcUnits"`
					UtilizedUncappedProcUnits []float64 `json:"utilizedUncappedProcUnits"`
					IdleProcUnits             []float64 `json:"idleProcUnits"`
					DonatedProcUnits          []int     `json:"donatedProcUnits"`
				} `json:"processor"`
				Network struct {
					GenericAdapters []struct {
						TransferredBytes []float64 `json:"transferredBytes"`
						Type             string    `json:"type"`
						ID               string    `json:"id"`
						PhysicalLocation string    `json:"physicalLocation"`
						ReceivedPackets  []float64 `json:"receivedPackets"`
						SentPackets      []float64 `json:"sentPackets"`
						DroppedPackets   []int     `json:"droppedPackets"`
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
						DroppedPackets   []int     `json:"droppedPackets"`
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
						NumOfReads       []int     `json:"numOfReads"`
						NumOfWrites      []float64 `json:"numOfWrites"`
						ReadBytes        []int     `json:"readBytes"`
						WriteBytes       []float64 `json:"writeBytes"`
					} `json:"genericPhysicalAdapters"`
					SharedStoragePools []struct {
						TransmittedBytes []int  `json:"transmittedBytes"`
						ID               string `json:"id"`
						TotalSpace       []int  `json:"totalSpace"`
						UsedSpace        []int  `json:"usedSpace"`
						NumOfReads       []int  `json:"numOfReads"`
						NumOfWrites      []int  `json:"numOfWrites"`
						ReadBytes        []int  `json:"readBytes"`
						WriteBytes       []int  `json:"writeBytes"`
					} `json:"sharedStoragePools"`
					FiberChannelAdapters []struct {
						TransmittedBytes []int  `json:"transmittedBytes"`
						Wwpn             string `json:"wwpn"`
						PhysicalLocation string `json:"physicalLocation"`
						NumOfPorts       int    `json:"numOfPorts"`
						RunningSpeed     []int  `json:"runningSpeed"`
						ID               string `json:"id"`
						NumOfReads       []int  `json:"numOfReads"`
						NumOfWrites      []int  `json:"numOfWrites"`
						ReadBytes        []int  `json:"readBytes"`
						WriteBytes       []int  `json:"writeBytes"`
					} `json:"fiberChannelAdapters"`
					GenericVirtualAdapters []struct {
						TransmittedBytes []int  `json:"transmittedBytes"`
						Type             string `json:"type"`
						ID               string `json:"id"`
						PhysicalLocation string `json:"physicalLocation"`
						NumOfReads       []int  `json:"numOfReads"`
						NumOfWrites      []int  `json:"numOfWrites"`
						ReadBytes        []int  `json:"readBytes"`
						WriteBytes       []int  `json:"writeBytes"`
					} `json:"genericVirtualAdapters"`
				} `json:"storage"`
			} `json:"viosUtil"`
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
