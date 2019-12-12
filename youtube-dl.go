package main

import (
	"encoding/json"
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func (target DownloadTarget) GetMetadata() (TargetMetadata, error) {
	cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", target.URL)
	log.Info("executing youtube-dl command: ", cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("From GetMetadata(): ", err)
		return TargetMetadata{}, fmt.Errorf("From c.GetMetadata(): %v", err)
	}
	metaData := &TargetMetadata{}
	if err = json.Unmarshal(out, metaData); err != nil {
		log.Error("From GetMetadata(): ", err)
		return TargetMetadata{}, fmt.Errorf("From c.GetMetadata(): %v", err)
	}

	return *metaData, nil
}

func (target DownloadTarget) GetLatestVideo() (string, error) {
	log.Info("fetching latest upload")
	cmd := exec.Command("youtube-dl", "-j", "--playlist-end", "1", target.URL)
	log.Info("executing youtube-dl command: ", cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		log.Errorf("c.GetLatestVideo: %s | %s", err, string(out))
	}
	metaData := &TargetMetadata{}
	if err = json.Unmarshal(out, metaData); err != nil {
		log.Error("c.GetLatestVideo: ", err)
		return "", fmt.Errorf("c.GetLatestVideo: %s", err)
	}
	log.Info("successfully fetched latest video ")
	return metaData.ID, nil
}

func (target DownloadTarget) Download(downloadQuality, fileExtension string, downloadEntire bool) error {
	log.Info("DOWNLOAD: ", target)
	var ytdlCommand YTDLCommand
	if target.DownloadMode == "Audio Only" {
		log.Info("downloading audio only")
		if downloadQuality == "best" {
			downloadQuality = "0"
		} else if downloadQuality == "medium" {
			downloadQuality = "5"
		} else if downloadQuality == "worst" {
			downloadQuality = "9"
		}
		log.Info("download quality set to: ", downloadQuality)
	}
	if target.Type == "Channel" || target.Type == "Playlist" {
		if target.DownloadPath == "" {
			if target.DownloadMode == "Audio Only" {
				if downloadEntire == false {
					ytdlCommand = YTDLCommand{
						Binary:       "youtube-dl",
						FirstFlag:    "--playlist-end",
						FirstFlagArg: "1",
						FileType:     "-f bestaudio[ext=" + fileExtension + "]",
						Output:       "-o downloads" + target.DownloadPath + " " + target.URL,
					}
				} else {
					ytdlCommand = YTDLCommand{
						Binary:       "youtube-dl",
						FirstFlag:    "--ignore-errors",
						FirstFlagArg: "",
						FileType:     "-f bestaudio[ext=" + fileExtension + "]",
						Output:       "-o downloads" + target.DownloadPath + " " + target.URL,
					}
				}
			} else if target.DownloadMode == "Video And Audio" {
				if downloadEntire == false {
					ytdlCommand = YTDLCommand{
						Binary:       "youtube-dl",
						FirstFlag:    "--playlist-end",
						FirstFlagArg: "1",
						FileType:     "-f bestaudio[ext=" + fileExtension + "]",
						Output:       "-o downloads" + target.DownloadPath,
						Target:       target.URL,
					}
				} else {
					ytdlCommand = YTDLCommand{
						Binary:       "youtube-dl",
						FirstFlag:    "--ignore-errors",
						FirstFlagArg: "",
						FileType:     "-f bestaudio[ext=" + fileExtension + "]",
						Output:       "-o downloads" + target.DownloadPath,
						Target:       target.URL,
					}
				}
			}
		} else {
			if target.DownloadMode == "Audio Only" {
				// Downloads only latest video
				if downloadEntire == false {
					ytdlCommand = YTDLCommand{
						Binary:       "youtube-dl",
						FirstFlag:    "--playlist-end",
						FirstFlagArg: "1",
						FileType:     "bestaudio[ext=" + fileExtension + "]",
						Output:       "downloads" + target.DownloadPath,
						Target:       target.URL,
					}
				} else {
					ytdlCommand = YTDLCommand{
						Binary:       "youtube-dl",
						FirstFlag:    "--ignore-errors",
						FirstFlagArg: "",
						FileType:     "bestaudio[ext=" + fileExtension + "]",
						Output:       "downloads" + target.DownloadPath,
						Target:       target.URL,
					}
				}
			} else if target.DownloadMode == "Video And Audio" {
				if downloadEntire == false {
					ytdlCommand = YTDLCommand{
						Binary:       "youtube-dl",
						FirstFlag:    "--playlist-end",
						FirstFlagArg: "1",
						FileType:     "bestaudio[ext=" + fileExtension + "]",
						Output:       "downloads" + target.DownloadPath,
						Target:       target.URL,
					}
				} else {
					ytdlCommand = YTDLCommand{
						Binary:       "youtube-dl",
						FirstFlag:    "--ignore-errors",
						FirstFlagArg: "",
						FileType:     "bestaudio[ext=" + fileExtension + "]",
						Output:       "downloads" + target.DownloadPath + " " + target.URL,
						Target:       target.URL,
					}
				}
			}
		}
	}

	err := DownloadVideo(ytdlCommand)
	if err != nil {
		return fmt.Errorf("Download: %v", err)
	}
	return nil
}

func DownloadVideo(command YTDLCommand) error {
	log.Info("downloading video")
	var cmd *exec.Cmd
	log.Info(command)
	if command.FirstFlag == "" {
		cmd = exec.Command(command.Binary, "-f", command.FileType, "-o", command.Output, command.Target)
	} else {
		cmd = exec.Command(command.Binary, command.FirstFlag, command.FirstFlagArg, "-f", command.FileType, "-o", command.Output, command.Target)
	}
	log.Info("executing youtube-dl command: ", cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Error(err, cmd.String())
		return fmt.Errorf("DownloadVideo: %s", err.Error())
	}
	return nil
}
