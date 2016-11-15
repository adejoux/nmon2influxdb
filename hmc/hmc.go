// nmon2influxdb
// import HMC data in InfluxDB
// author: adejoux@djouxtech.net

package hmc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"text/template"
	"time"

	"github.com/adejoux/influxdbclient"
	"github.com/adejoux/nmon2influxdb/nmon2influxdblib"
	"github.com/codegangsta/cli"
)

type HMC struct {
	Session  *Session
	InfluxDB *influxdbclient.InfluxDB
}

type HMCPoint struct {
	Name      string
	Server    string
	Metric    string
	Pool      string
	Value     interface{}
	Timestamp time.Time
}

func NewHMC(c *cli.Context) *HMC {

	var hmc HMC
	// parsing parameters
	config := nmon2influxdblib.ParseParameters(c)

	//getting databases connections
	hmc.InfluxDB = config.GetDB("hmc")
	hmcURL := fmt.Sprintf("https://"+"%s"+":12443", config.HMCServer)
	//initialize new http session
	hmc.Session = NewSession(config.HMCUser, config.HMCPassword, hmcURL)
	hmc.Session.doLogon()

	return &hmc
}

func (hmc *HMC) WritePoints() (err error) {
	err = hmc.InfluxDB.WritePoints()
	nmon2influxdblib.CheckError(err)
	hmc.InfluxDB.ClearPoints()
	return
}

func (hmc *HMC) GetManagedSystems() []System {
	return hmc.Session.getManagedSystems()
}

func (hmc *HMC) GetPCMLinks(uuid string) ([]string, error) {
	return hmc.Session.getPCM(uuid)
}

func (hmc *HMC) GetPCMData(link string) (PCMData, error) {
	return hmc.Session.getPCMData(link)
}

func (hmc *HMC) AddPoint(point HMCPoint) {

	value, ok := point.Value.(float64)
	if ok != true {
		// Else it's int type(no other possible types)
		value = float64(point.Value.(int))
	}

	tags := map[string]string{"server": point.Server, "name": point.Metric}
	if len(point.Pool) > 0 {
		tags["pool"] = point.Pool
	}
	field := map[string]interface{}{"value": value}
	hmc.InfluxDB.AddPoint(point.Name, point.Timestamp, field, tags)
}

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

//
// XML parsing structures
//

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	XMLName xml.Name `xml:"entry"`
	ID      string   `xml:"id"`
	Link    struct {
		Href string `xml:"href,attr"`
	} `xml:"link,omitempty"`
	Contents []Content `xml:"content"`
}

type Content struct {
	XMLName xml.Name        `xml:"content"`
	System  []ManagedSystem `xml:"http://www.ibm.com/xmlns/systems/power/firmware/uom/mc/2012_10/ ManagedSystem"`
}

type ManagedSystem struct {
	XMLName    xml.Name `xml:"http://www.ibm.com/xmlns/systems/power/firmware/uom/mc/2012_10/ ManagedSystem"`
	SystemName string
}

//
// HTTP session struct
//

type Session struct {
	client   *http.Client
	User     string
	Password string
	url      string
}

func NewSession(user string, password string, url string) *Session {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Session{client: &http.Client{Transport: tr, Jar: jar}, User: user, Password: password, url: url}
}

func (s *Session) doLogon() {

	authurl := s.url + "/rest/api/web/Logon"

	// template for login request
	logintemplate := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
  <LogonRequest xmlns="http://www.ibm.com/xmlns/systems/power/firmware/web/mc/2012_10/" schemaVersion="V1_1_0">
    <Metadata>
      <Atom/>
    </Metadata>
    <UserID kb="CUR" kxe="false">{{.User}}</UserID>
    <Password kb="CUR" kxe="false">{{.Password}}</Password>
  </LogonRequest>`

	tmpl := template.New("logintemplate")
	tmpl.Parse(logintemplate)
	authrequest := new(bytes.Buffer)
	err := tmpl.Execute(authrequest, s)
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("PUT", authurl, authrequest)

	// set request headers
	request.Header.Set("Content-Type", "application/vnd.ibm.powervm.web+xml; type=LogonRequest")
	request.Header.Set("Accept", "application/vnd.ibm.powervm.web+xml; type=LogonResponse")
	request.Header.Set("X-Audit-Memento", "hmctest")

	response, err := s.client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		if response.StatusCode != 200 {
			log.Fatalf("Error status code: %d", response.StatusCode)
		}
	}
}

func (s *Session) getPCM(uuid string) ([]string, error) {
	var pcmlinks []string
	pcmurl := s.url + "/rest/api/pcm/ManagedSystem/" + uuid + "/ProcessedMetrics"
	request, _ := http.NewRequest("GET", pcmurl, nil)

	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	response, requestErr := s.client.Do(request)
	if requestErr != nil {
		return pcmlinks, requestErr
	}

	defer response.Body.Close()
	contents, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return pcmlinks, readErr
	}

	if response.StatusCode != 200 {
		errorMessage := fmt.Sprintf("Error getting PCM informations. status code: %d", response.StatusCode)
		statusErr := errors.New(errorMessage)
		return pcmlinks, statusErr
	}

	var feed Feed
	unmarshalErr := xml.Unmarshal(contents, &feed)

	if unmarshalErr != nil {
		return pcmlinks, unmarshalErr
	}
	for _, entry := range feed.Entries {
		pcmlinks = append(pcmlinks, entry.Link.Href)
	}

	return pcmlinks, nil
}

func (s *Session) getPCMData(rawurl string) (PCMData, error) {
	//the link url can use ip address instead of hostname. Authentication was not performed on it.
	//it's better to only keep the path and use the url provided by the user.
	var data PCMData
	u, _ := url.Parse(rawurl)
	pcmurl := s.url + u.Path

	request, err := http.NewRequest("GET", pcmurl, nil)

	//request.Header.Set("Accept", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return data, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return data, err
	}

	if response.StatusCode != 200 {
		log.Fatalf("Error getting PCM Data informations. status code: %d", response.StatusCode)
	}

	// var prettyJSON bytes.Buffer
	// error := json.Indent(&prettyJSON, contents, "", "\t")
	// if error != nil {
	// 	log.Println("JSON parse error: ", error)
	// 	return
	// }
	//
	// log.Println("output:", string(prettyJSON.Bytes()))

	json.Unmarshal(contents, &data)

	return data, err

}

func (s *Session) getManagedSystems() (systems []System) {
	mgdurl := s.url + "/rest/api/uom/ManagedSystem"
	request, err := http.NewRequest("GET", mgdurl, nil)

	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	response, err := s.client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		if response.StatusCode != 200 {
			log.Fatalf("Error getting LPAR informations. status code: %d", response.StatusCode)
		}

		var feed Feed
		newErr := xml.Unmarshal(contents, &feed)

		if newErr != nil {
			log.Fatal(newErr)
		}
		for _, entry := range feed.Entries {

			for _, content := range entry.Contents {
				for _, system := range content.System {
					systems = append(systems, System{Name: system.SystemName, UUID: entry.ID})
				}
			}
		}
	}
	return
}

type System struct {
	Name string
	UUID string
}
