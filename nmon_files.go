// nmon2influxdb
// import nmon data in InfluxDB
// author: adejoux@djouxtech.net

package main

import (
	"bufio"
	"compress/gzip"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"sort"

	"github.com/pkg/sftp"

	"golang.org/x/crypto/ssh"
)

var remoteFileRegexp = regexp.MustCompile(`(\S+):(\S+)`)

const gzipfile = ".gz"
const size = 64000

// NmonFile structure used to select nmon files to import
type NmonFile struct {
	Name     string
	FileType string
	Host     string
	SSHUser  string
	SSHKey   string
	checksum string
	lines    []string
}

// NmonFiles array of NmonFile
type NmonFiles []NmonFile

//Add a file in the NmonFIles structure
func (nmonFiles *NmonFiles) Add(file string, fileType string) {
	*nmonFiles = append(*nmonFiles, NmonFile{Name: file, FileType: fileType})
}

//AddRemote a remote file in the NmonFIles structure
func (nmonFiles *NmonFiles) AddRemote(file string, fileType string, host string, user string, key string) {
	*nmonFiles = append(*nmonFiles, NmonFile{Name: file, FileType: fileType, Host: host, SSHUser: user, SSHKey: key})
}

//Valid returns only valid fiels for nmon import
func (nmonFiles *NmonFiles) Valid() (validFiles NmonFiles) {
	for _, v := range *nmonFiles {
		if v.FileType == ".nmon" || v.FileType == gzipfile {
			validFiles = append(validFiles, v)
		}
	}
	return validFiles
}

// FileScanner struct to manage
type FileScanner struct {
	*os.File
	*bufio.Scanner
}

// RemoteFileScanner struct for remote files
type RemoteFileScanner struct {
	*sftp.File
	*bufio.Scanner
}

// GetRemoteScanner open an nmon file based on file extension and provides a bufio Scanner
func (nmonFile *NmonFile) GetRemoteScanner() (*RemoteFileScanner, error) {

	sftpConn := InitSFTP(nmonFile.SSHUser, nmonFile.Host, nmonFile.SSHKey)
	file, err := sftpConn.Open(nmonFile.Name)
	if err != nil {
		return nil, err
	}

	if nmonFile.FileType == gzipfile {
		gr, err := gzip.NewReader(file)
		if err != nil {
			return nil, err
		}
		reader := bufio.NewReader(gr)
		return &RemoteFileScanner{file, bufio.NewScanner(reader)}, nil
	}

	reader := bufio.NewReader(file)
	return &RemoteFileScanner{file, bufio.NewScanner(reader)}, nil
}

//Checksum generates SHA1 file checksum
func (nmonFile *NmonFile) Checksum() (fileHash string) {
	if len(nmonFile.checksum) > 0 {
		return nmonFile.checksum
	}
	var result []byte
	if len(nmonFile.Host) > 0 {
		scanner, err := nmonFile.GetRemoteScanner()
		check(err)
		scanner.Seek(1024, 2)
		hash := sha1.New()
		if _, err = io.Copy(hash, scanner); err != nil {
			return
		}
		fileHash = hex.EncodeToString(hash.Sum(result))
	} else {
		scanner, err := nmonFile.GetScanner()
		check(err)
		scanner.Seek(1024, 2)
		hash := sha1.New()
		if _, err = io.Copy(hash, scanner); err != nil {
			return
		}
		fileHash = hex.EncodeToString(hash.Sum(result))
	}
	nmonFile.checksum = fileHash
	return
}

// GetScanner open an nmon file based on file extension and provides a bufio Scanner
func (nmonFile *NmonFile) GetScanner() (*FileScanner, error) {

	file, err := os.Open(nmonFile.Name)
	if err != nil {
		return nil, err
	}

	if nmonFile.FileType == gzipfile {
		gr, err := gzip.NewReader(file)
		if err != nil {
			return nil, err
		}
		reader := bufio.NewReader(gr)
		return &FileScanner{file, bufio.NewScanner(reader)}, nil
	}

	reader := bufio.NewReader(file)
	return &FileScanner{file, bufio.NewScanner(reader)}, nil
}

// Parse parameters
func (nmonFiles *NmonFiles) Parse(args []string, user string, key string) {
	for _, param := range args {
		if remoteFileRegexp.MatchString(param) {
			matched := remoteFileRegexp.FindStringSubmatch(param)
			host := matched[1]
			matchedParam := matched[2]

			sftpConn := InitSFTP(user, host, key)
			paraminfo, err := sftpConn.Stat(matchedParam)
			check(err)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("%s doesn't exist ! skipped.\n", param)
				}
				continue
			}
			if paraminfo.IsDir() {
				entries, err := sftpConn.ReadDir(matchedParam)
				check(err)
				for _, entry := range entries {
					if !entry.IsDir() {
						file := path.Join(matchedParam, entry.Name())
						nmonFiles.AddRemote(file, path.Ext(file), host, user, key)
					}
				}
				sftpConn.Close()
				continue
			}
			nmonFiles.AddRemote(matchedParam, path.Ext(matchedParam), host, user, key)
			sftpConn.Close()
			continue
		}

		paraminfo, err := os.Stat(param)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("%s doesn't exist ! skipped.\n", param)
			}
			continue
		}

		if paraminfo.IsDir() {
			entries, err := ioutil.ReadDir(param)
			check(err)
			for _, entry := range entries {
				if !entry.IsDir() {
					file := path.Join(param, entry.Name())
					nmonFiles.Add(file, path.Ext(file))
				}
			}
			continue
		}
		nmonFiles.Add(param, path.Ext(param))
	}
}

//SSHConfig contains SSH parameters
type SSHConfig struct {
	User string
	Key  string
}

//InitSFTP init sftp session
func InitSFTP(user string, host string, key string) *sftp.Client {
	var auths []ssh.AuthMethod
	if !IsNotFile(key) {
		pemBytes, err := ioutil.ReadFile(key)
		if err != nil {
			log.Fatal(err)
		}
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			log.Fatalf("parse key failed:%v", err)
		}

		auths = append(auths, ssh.PublicKeys(signer))
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: auths,
	}
	sshhost := fmt.Sprintf("%s:22", host)
	conn, err := ssh.Dial("tcp", sshhost, config)
	if err != nil {
		log.Fatalf("dial failed:%v", err)
	}

	c, err := sftp.NewClient(conn, sftp.MaxPacket(size))
	if err != nil {
		log.Fatalf("unable to start sftp subsytem: %v", err)
	}
	return c
}

//Content returns the nmon files content sorted in an slice of string format
func (nmonFile *NmonFile) Content() []string {
	if len(nmonFile.lines) > 0 {
		return nmonFile.lines
	}
	if len(nmonFile.Host) > 0 {
		scanner, err := nmonFile.GetRemoteScanner()
		check(err)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			nmonFile.lines = append(nmonFile.lines, scanner.Text())
		}
		scanner.Close()
	} else {
		scanner, err := nmonFile.GetScanner()
		check(err)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			nmonFile.lines = append(nmonFile.lines, scanner.Text())
		}
		scanner.Close()
	}

	sort.Strings(nmonFile.lines)

	return nmonFile.lines
}
