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

// HMC contains the base struct used by all the hmc sub command
type HMC struct {
	Session     *Session
	InfluxDB    *influxdbclient.InfluxDB
	GlobalPoint Point
}

// Point is a struct to simplify InfluxDB point creation
type Point struct {
	Name      string
	Server    string
	Metric    string
	Pool      string
	Device    string
	Partition string
	Type      string
	Value     interface{}
	Timestamp time.Time
}

//NewHMC return a new HMC struct and use the command line and config file parameters to intialize it.
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

// WritePoints send points to InfluxDB database and reset points count
func (hmc *HMC) WritePoints() (err error) {
	err = hmc.InfluxDB.WritePoints()
	nmon2influxdblib.CheckError(err)
	hmc.InfluxDB.ClearPoints()
	return
}

// AddPoint add a InfluxDB point. It's using the GlobalPoint parameter to fill some fields
func (hmc *HMC) AddPoint(point Point) {

	value, ok := point.Value.(float64)
	if ok != true {
		// Else it's int type(no other possible types)
		value = float64(point.Value.(int))
	}

	tags := map[string]string{"server": hmc.GlobalPoint.Server, "name": point.Metric}
	if len(point.Pool) > 0 {
		tags["pool"] = point.Pool
	}

	if len(point.Type) > 0 {
		tags["type"] = point.Type
	}
	if len(point.Device) > 0 {
		tags["device"] = point.Device
	}

	if len(hmc.GlobalPoint.Partition) > 0 {
		tags["partition"] = hmc.GlobalPoint.Partition
	}
	field := map[string]interface{}{"value": value}
	hmc.InfluxDB.AddPoint(point.Name, hmc.GlobalPoint.Timestamp, field, tags)
}

//
// XML parsing structures
//

// Feed base struct of Atom feed
type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Entry  `xml:"entry"`
}

// Entry is the atom feed section containing the links to PCM data and the Category
type Entry struct {
	XMLName xml.Name `xml:"entry"`
	ID      string   `xml:"id"`
	Link    struct {
		Href string `xml:"href,attr"`
	} `xml:"link,omitempty"`
	Contents []Content `xml:"content"`
	Category struct {
		Term string `xml:"term,attr"`
	} `xml:"category,omitempty"`
}

// Content feed struct containing all managed systems
type Content struct {
	XMLName xml.Name        `xml:"content"`
	System  []ManagedSystem `xml:"http://www.ibm.com/xmlns/systems/power/firmware/uom/mc/2012_10/ ManagedSystem"`
}

// ManagedSystem struct contains a managed system and his associated partitions
type ManagedSystem struct {
	XMLName                     xml.Name `xml:"http://www.ibm.com/xmlns/systems/power/firmware/uom/mc/2012_10/ ManagedSystem"`
	SystemName                  string
	AssociatedLogicalPartitions Partitions `xml:"http://www.ibm.com/xmlns/systems/power/firmware/uom/mc/2012_10/ AssociatedLogicalPartitions"`
}

// Partitions contains links to the partition informations
type Partitions struct {
	Links []Link `xml:"link,omitempty"`
}

// Link the link itself is stored in the attribute href
type Link struct {
	Href string `xml:"href,attr"`
}

// Session is the HTTP session struct
type Session struct {
	client   *http.Client
	User     string
	Password string
	url      string
}

// NewSession initialize a Session struct
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

// doLogon performs the login to the inflxudb instance
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

// PCMLinks store a system and associated partitions links to PCM data
type PCMLinks struct {
	System     string
	Partitions []string
}

// GetPCMLinks encapsulation function
func (hmc *HMC) GetPCMLinks(uuid string) (PCMLinks, error) {
	return hmc.Session.getPCMLinks(uuid)
}

func (s *Session) getPCMLinks(uuid string) (PCMLinks, error) {
	var pcmlinks PCMLinks
	pcmurl := s.url + "/rest/api/pcm/ManagedSystem/" + uuid + "/ProcessedMetrics"
	request, _ := http.NewRequest("GET", pcmurl, nil)

	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	response, requestErr := s.client.Do(request)
	if requestErr != nil {
		return pcmlinks, requestErr
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		errorMessage := fmt.Sprintf("Error getting PCM informations. status code: %d", response.StatusCode)
		statusErr := errors.New(errorMessage)
		return pcmlinks, statusErr
	}

	var feed Feed
	contents, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return pcmlinks, readErr
	}
	unmarshalErr := xml.Unmarshal(contents, &feed)

	if unmarshalErr != nil {
		return pcmlinks, unmarshalErr
	}
	for _, entry := range feed.Entries {
		if len(entry.Category.Term) == 0 {
			continue
		}
		if entry.Category.Term == "ManagedSystem" {
			pcmlinks.System = entry.Link.Href
		}

		if entry.Category.Term == "LogicalPartition" {
			pcmlinks.Partitions = append(pcmlinks.Partitions, entry.Link.Href)
		}
	}

	return pcmlinks, nil
}

// GetPCMData encapsulation function
func (hmc *HMC) GetPCMData(link string) (PCMData, error) {
	return hmc.Session.getPCMData(link)
}

// get PCMData retreives the PCM data in JSON format and returns them stored in an PCMData struct
func (s *Session) getPCMData(rawurl string) (PCMData, error) {
	var data PCMData
	u, _ := url.Parse(rawurl)
	pcmurl := s.url + u.Path
	request, err := http.NewRequest("GET", pcmurl, nil)

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

	json.Unmarshal(contents, &data)

	return data, err

}

// GetManagedSystems encapsulation function
func (hmc *HMC) GetManagedSystems() []System {
	return hmc.Session.getManagedSystems()
}

// getManagedSystems returns a list of the managed systems retrieved from the atom feed
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

// System struct store system Name and UUID
type System struct {
	Name string
	UUID string
}
