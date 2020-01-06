package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	TypeChannel  = "Channel"
	TypePlaylist = "Playlist"
)

var confs = map[string]*Database{
	TypeChannel:  &Database{RWMutex: &sync.RWMutex{}, file: "channels.json"},
	TypePlaylist: &Database{RWMutex: &sync.RWMutex{}, file: "playlists.json"},
}

func init() {
	initLogFile()
	CreateDirIfNotExist("./config")
	initDatabase()
	for _, v := range confs {
		if err := v.load(); err != nil {
			log.Fatal(nil)
		}
	}
}

func (db *Database) write() error {
	b, err := json.MarshalIndent(&db.contents, "", "	")
	if err != nil {
		err = fmt.Errorf("database.write: could not marshal json: %w", err)
		log.Error(err)
		return err
	}
	err = ioutil.WriteFile(filepath.Join(CONFIG_ROOT, db.file), b, 0644)
	if err != nil {
		err = fmt.Errorf("database.write: could not write file %w", err)
		log.Error(err)
		return err
	}
	return nil
}

func (db *Database) load() error {
	b, err := ioutil.ReadFile(filepath.Join(CONFIG_ROOT, db.file))
	if err != nil {
		err = fmt.Errorf("database.load: could not read file: %w", err)
		log.Error(err)
		return err
	}
	err = json.Unmarshal(b, &db.contents)
	if err != nil {
		err = fmt.Errorf("database.load: could not read file: %w", err)
		log.Error(err)
		return err
	}
	db.lookup = make(map[string]*DownloadTarget)
	for i := range db.contents {
		db.lookup[db.contents[i].URL] = &db.contents[i]
	}
	return nil
}

func UpdateCheckingInterval(file, time, intervalStr string) error {
	log.Info("update checking interval")
	cf, ok := confs[file]
	if !ok {
		return fmt.Errorf("UpdateCheckingInterval: bad conf type")
	}
	if len(cf.contents) == 0 {
		return fmt.Errorf("UpdateCheckingInterval: empty config list")
	}
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		return fmt.Errorf("UpdateCheckingInterval: non-numeric interval %w", err)
	}
	_, ok = checkingIntervalMultipliers[time]
	if !ok {
		return fmt.Errorf("UpdateCheckingInterval: bad multiplier")
	}
	cf.Lock()
	cf.contents[0].CheckingInterval = interval
	cf.contents[0].CheckingIntervalTime = time
	cf.write()
	cf.Unlock()
	return nil
}

func updateDB(file, url string, f func(t *DownloadTarget) error) error {
	cf, ok := confs[file]
	if !ok {
		return fmt.Errorf("updateDB: bad conf type")
	}
	cf.Lock()
	defer cf.Unlock()
	t, ok := cf.lookup[url]
	if !ok {
		return fmt.Errorf("updateDB: no such URL")
	}
	err := f(t)
	if err != nil {
		return err
	}
	return cf.write()
}

func (t DownloadTarget) UpdateLastChecked() error {
	return updateDB(t.Type, t.URL, func(t *DownloadTarget) error {
		t.LastChecked = time.Now().Format("01-02-2006 15:04:05")
		log.Info("last checked date and time updated successfully")
		return nil
	})
}

func (t DownloadTarget) UpdateLatestDownloaded(videoId string) error {
	return updateDB(t.Type, t.URL, func(t *DownloadTarget) error {
		t.LatestDownloaded = videoId
		log.Info("latest downloaded video id updated successfully")
		return nil
	})
}

func (t DownloadTarget) DoesExist() (bool, error) {
	cf, ok := confs[t.Type]
	if !ok {
		return false, fmt.Errorf("DoesExist: bad conf type")
	}
	cf.RLock()
	_, ok = cf.lookup[t.URL]
	cf.RUnlock()
	return ok, nil
}

func (t DownloadTarget) AddToDatabase() error {
	cf, ok := confs[t.Type]
	if !ok {
		return fmt.Errorf("updateDB: bad conf type")
	}
	cf.Lock()
	cf.contents = append(cf.contents, t)
	cf.lookup[t.URL] = &cf.contents[len(cf.contents)-1]
	cf.Unlock()
	return nil
}

func (t DownloadTarget) UpdateDownloadHistory(videoId string) error {
	return updateDB(t.Type, t.URL, func(t *DownloadTarget) error {
		t.DownloadHistory = append(t.DownloadHistory, videoId)
		log.Info(t.DownloadHistory)
		return nil
	})
}

func (t DownloadTarget) Delete() error {
	cf, ok := confs[t.Type]
	if !ok {
		return fmt.Errorf("Delete: bad conf type")
	}
	cf.Lock()
	defer cf.write()
	defer cf.Unlock()
	delete(cf.lookup, t.URL)
	for i := range cf.contents {
		if cf.contents[i].URL == t.URL {
			cf.contents = RemoveAtIndex(cf.contents, i)
			log.Info("successfully removed channel from channels.json")
			return nil
		}
	}
	return fmt.Errorf("DownloadTarget.Delete: no such target")
}

func (t DownloadTarget) GetFromDatabase() (out DownloadTarget, err error) {
	cf, ok := confs[t.Type]
	if !ok {
		err = fmt.Errorf("GetFromDatabase: bad conf type")
		return
	}
	cf.RLock()
	defer cf.RUnlock()
	ft, ok := cf.lookup[t.URL]
	if !ok {
		err = fmt.Errorf("GetFromDatabase: no such URL")
		return
	}
	out = *ft
	return
}

func GetCheckingInterval(file string) (int, error) {
	interval, time, err := GetCheckingIntervalConfig(file)
	if err != nil {
		return 0, err
	}

	if interval == 0 {
		return 0, nil
	}
	return interval * checkingIntervalMultipliers[time], nil
}

func GetCheckingIntervalConfig(file string) (int, string, error) {
	cf, ok := confs[file]
	if !ok {
		return 0, "", fmt.Errorf("GetFromDatabase: bad conf type")
	}
	cf.RLock()
	defer cf.RUnlock()
	if len(cf.contents) == 0 {
		return 0, "", fmt.Errorf("GetCheckingInterval: empty target list")
	}
	target := cf.contents[0]
	return target.CheckingInterval, target.CheckingIntervalTime, nil
}
