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

// hmc can be really slow to answer
const timeout = 30

// HMC contains the base struct used by all the hmc sub command
type HMC struct {
	Session             *Session
	InfluxDB            *influxdbclient.InfluxDB
	GlobalPoint         Point
	FilterManagedSystem string
	Debug               bool
	ManagedSystemOnly   bool
	Samples             int
	TagParsers          nmon2influxdblib.TagParsers
}

// Point is a struct to simplify InfluxDB point creation
type Point struct {
	Name                    string
	System                  string
	Metric                  string
	Pool                    string
	Device                  string
	Partition               string
	Type                    string
	WWPN                    string
	PhysicalPortWWPN        string
	ViosID                  string
	VlanID                  string
	VswitchID               string
	SharedEthernetAdapterID string
	DrcIndex                string
	PhysicalLocation        string
	PhysicalDrcIndex        string
	PhysicalPortID          string
	Instance                string
	InstanceID              string
	Value                   interface{}
	Timestamp               time.Time
}

//NewHMC return a new HMC struct and use the command line and config file parameters to intialize it.
func NewHMC(c *cli.Context) *HMC {

	var hmc HMC
	// parsing parameters
	config := nmon2influxdblib.ParseParameters(c)

	if config.Debug {
		log.Printf("configuration: %+v\n", config.Sanitized())
	}

	//getting databases connections
	hmc.InfluxDB = config.GetDB("hmc")
	hmc.ManagedSystemOnly = config.HMCManagedSystemOnly
	hmc.Samples = config.HMCSamples
	hmc.Debug = config.Debug
	hmc.FilterManagedSystem = config.HMCManagedSystem

	if len(config.Inputs) > 0 {
		//Build tag parsing
		hmc.TagParsers = nmon2influxdblib.ParseInputs(config.Inputs)
	}
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

	tags := map[string]string{"system": hmc.GlobalPoint.System, "name": point.Metric}
	if len(hmc.GlobalPoint.Pool) > 0 {
		tags["pool"] = hmc.GlobalPoint.Pool
	}

	if len(hmc.GlobalPoint.Type) > 0 {
		tags["type"] = hmc.GlobalPoint.Type
	}
	if len(hmc.GlobalPoint.Device) > 0 {
		tags["device"] = hmc.GlobalPoint.Device
	}

	if len(hmc.GlobalPoint.WWPN) > 0 {
		tags["wwpn"] = hmc.GlobalPoint.WWPN
	}

	if len(hmc.GlobalPoint.PhysicalPortWWPN) > 0 {
		tags["PhysicalPortWWPN"] = hmc.GlobalPoint.PhysicalPortWWPN
	}

	if len(hmc.GlobalPoint.ViosID) > 0 {
		tags["ViosID"] = hmc.GlobalPoint.ViosID
	}
	if len(hmc.GlobalPoint.VlanID) > 0 {
		tags["VlanID"] = hmc.GlobalPoint.VlanID
	}
	if len(hmc.GlobalPoint.VswitchID) > 0 {
		tags["VswitchID"] = hmc.GlobalPoint.VswitchID
	}
	if len(hmc.GlobalPoint.SharedEthernetAdapterID) > 0 {
		tags["SEA"] = hmc.GlobalPoint.SharedEthernetAdapterID
	}
	if len(hmc.GlobalPoint.Partition) > 0 {
		tags["partition"] = hmc.GlobalPoint.Partition
	}
	field := map[string]interface{}{"value": value}

	// Checking additional tagging
	for key, value := range tags {
		if _, ok := hmc.TagParsers[point.Name][key]; ok {
			for _, tagParser := range hmc.TagParsers[point.Name][key] {
				if tagParser.Regexp.MatchString(value) {
					tags[tagParser.Name] = tagParser.Value
				}
			}
		}
	}
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

	return &Session{client: &http.Client{Transport: tr, Jar: jar, Timeout: time.Second * timeout}, User: user, Password: password, url: url}
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
		log.Fatalf("HMC error sending auth request: %v\n", err)
	} else {
		defer response.Body.Close()
		if response.StatusCode != 200 {
			log.Fatalf("HMC authentication error: %s\n", response.Status)
		}
	}
}

// PCMLinks store a system and associated partitions links to PCM data
type PCMLinks struct {
	System     string
	Partitions []string
}

// GetSystemPCMLinks encapsulation function
func (hmc *HMC) GetSystemPCMLinks(uuid string) (PCMLinks, error) {
	var pcmURL string
	if hmc.Samples > 0 {
		pcmURL = fmt.Sprintf("%s/rest/api/pcm/ManagedSystem/%s/ProcessedMetrics?NoOfSamples=%d", hmc.Session.url, uuid, hmc.Samples)
	} else {
		pcmURL = hmc.Session.url + "/rest/api/pcm/ManagedSystem/" + uuid + "/ProcessedMetrics"
	}
	return hmc.Session.getPCMLinks(pcmURL, hmc.Debug)
}

// GetPartitionPCMLinks encapsulation function
func (hmc *HMC) GetPartitionPCMLinks(link string) (PCMLinks, error) {
	var pcmURL string
	if hmc.Samples > 0 {
		pcmURL = fmt.Sprintf("%s%s?NoOfSamples=%d", hmc.Session.url, link, hmc.Samples)
	} else {
		pcmURL = hmc.Session.url + link
	}
	return hmc.Session.getPCMLinks(pcmURL, hmc.Debug)
}

func (s *Session) getPCMLinks(link string, debug bool) (PCMLinks, error) {
	if debug {
		log.Printf("getPCMLinks link: %s\n", link)
	}
	var pcmlinks PCMLinks
	request, _ := http.NewRequest("GET", link, nil)

	request.Header.Set("Accept", "*/*;q=0.8")

	if debug {
		log.Printf("getPCMLinks HTTP request: ")
		log.Printf(nmon2influxdblib.SPrintHTTPRequest(request))
	}
	response, requestErr := s.client.Do(request)
	if requestErr != nil {
		return pcmlinks, requestErr
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		errorMessage := fmt.Sprintf("Error getting PCM informations. status code: %d", response.StatusCode)
		statusErr := errors.New(errorMessage)
		if debug {
			log.Printf("getPCMLinks HTTP response: ")
			log.Printf(nmon2influxdblib.SPrintHTTPResponse(response))
		}
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
	return hmc.Session.getPCMData(link, hmc.Debug)
}

// GetEnergyPCMData encapsulation function
func (hmc *HMC) GetEnergyPCMData(link string) (PCMData, error) {
	link += "?type=Energy"
	return hmc.Session.getPCMData(link, hmc.Debug)
}

// get PCMData retreives the PCM data in JSON format and returns them stored in an PCMData struct
func (s *Session) getPCMData(rawurl string, debug bool) (PCMData, error) {
	var data PCMData
	u, _ := url.Parse(rawurl)
	pcmurl := s.url + u.Path
	if debug {
		log.Printf("getPCMData link:%s\n", pcmurl)
	}
	request, _ := http.NewRequest("GET", pcmurl, nil)

	response, err := s.client.Do(request)
	if err != nil {
		return data, err
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		if debug {
			log.Printf("getPCMData response: ")
			log.Printf(nmon2influxdblib.SPrintHTTPResponse(response))
		}
		return data, err
	}

	if debug {
		log.Printf("getPCMData JSON: ")
		log.Printf(nmon2influxdblib.SPrintPrettyJSON(contents))
	}

	if response.StatusCode != 200 {
		log.Fatalf("Error getting PCM Data informations. status code: %d", response.StatusCode)
	}

	jsonErr := json.Unmarshal(contents, &data)

	if jsonErr != nil {
		log.Printf(nmon2influxdblib.SPrintPrettyJSON(contents))
	}
	return data, jsonErr

}

// GetManagedSystems encapsulation function
func (hmc *HMC) GetManagedSystems() ([]System, error) {
	return hmc.Session.getManagedSystems()
}

// getManagedSystems returns a list of the managed systems retrieved from the atom feed
func (s *Session) getManagedSystems() (systems []System, err error) {
	mgdurl := s.url + "/rest/api/uom/ManagedSystem"
	request, _ := http.NewRequest("GET", mgdurl, nil)

	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	response, err := s.client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()
	contents, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return systems, readErr
	}

	if response.StatusCode != 200 {
		log.Fatalf("Error getting LPAR informations. status code: %d", response.StatusCode)
	}

	var feed Feed
	newErr := xml.Unmarshal(contents, &feed)

	if newErr != nil {
		return systems, newErr
	}
	for _, entry := range feed.Entries {

		for _, content := range entry.Contents {
			for _, system := range content.System {
				systems = append(systems, System{Name: system.SystemName, UUID: entry.ID})
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
