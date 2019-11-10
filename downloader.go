package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// THIS DOWNLOADER IS BASED ON https://github.com/knadh/go-get-youtube

// Download downloads a video based on downloadMode
func Download(channelName, channelType, downloadMode string) error {
	videoId, videoTitle := GetLatestVideo(channelName, channelType)
	path := channelName + "/" + videoTitle + ".mp4"

	video, err := getVideoData(videoId)
	if err != nil {
		log.Error("Couldn't get video data: ", err)
	}
	if downloadMode == "Video And Audio" { // Download .mp4 with audio and video in one file
		option := &Option{
			Rename: false,
			Resume: true,
			Mp3:    false,
		}
		video.Download(0, path, option, videoId)
		log.Info("Successfully downloaded video")
		UpdateLatestDownloaded(channelName, videoId)
	} else if downloadMode == "Audio Only" { // Extract audio from the .mp4 file and remove the .mp4
		option := &Option{
			Rename: false,
			Resume: true,
			Mp3:    true,
		}
		video.Download(0, path, option, videoId)
		log.Info("Removing .mp4")
		err := os.Remove(path)
		if err != nil {
			log.Error("Error removing .mp4:", err)
		} else {
			log.Info("Successfully removed .mp4")
			log.Info("Successfully downloaded video")
			UpdateLatestDownloaded(channelName, videoId)
			didDownloadFail := CheckIfDownloadFailed(path)
			if didDownloadFail == true {
				InsertFailedDownload(videoId)
				return fmt.Errorf("Download failed, writing to failed.json")
			}
		}
	}

	return fmt.Errorf("Something went seriously wrong")
}

func parseMeta(video_id, query_string string) (*Video, error) {
	// parse the query string
	u, _ := url.Parse("?" + query_string)

	// parse url params
	query := u.Query()

	// no such video
	if query.Get("errorcode") != "" || query.Get("status") == "fail" {
		return nil, errors.New(query.Get("reason"))
	}

	var player_response PlayerResponse
	json.Unmarshal([]byte(query.Get("player_response")), &player_response)

	// collate the necessary params
	video := &Video{
		Id:            video_id,
		Title:         player_response.VideoDetails.Title,
		Author:        player_response.VideoDetails.Author,
		Keywords:      fmt.Sprint(player_response.VideoDetails.Keywords),
		Thumbnail_url: player_response.VideoDetails.Thumbnail.Thumbnails[0].URL,
	}

	v, _ := strconv.Atoi(player_response.VideoDetails.ViewCount)
	video.View_count = v

	video.Avg_rating = float32(player_response.VideoDetails.AverageRating)

	l, _ := strconv.Atoi(player_response.VideoDetails.LengthSeconds)
	video.Length_seconds = l

	// further decode the format data
	format_params := strings.Split(query.Get("url_encoded_fmt_stream_map"), ",")

	// every video has multiple format choices. collate the list.
	for _, f := range format_params {
		furl, _ := url.Parse("?" + f)
		fquery := furl.Query()

		itag, _ := strconv.Atoi(fquery.Get("itag"))

		video.Formats = append(video.Formats, Format{
			Itag:       itag,
			Video_type: fquery.Get("type"),
			Quality:    fquery.Get("quality"),
			Url:        fquery.Get("url"),
		})
	}

	return video, nil
}

func fetchMeta(video_id string) (string, error) {
	resp, err := http.Get(URL_META + video_id)

	// fetch the meta information from http
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	query_string, _ := ioutil.ReadAll(resp.Body)

	return string(query_string), nil
}

func getVideoData(videoId string) (Video, error) {
	// fetch video meta from youtube
	query_string, err := fetchMeta(videoId)

	if err != nil {
		return Video{}, err
	}

	meta, err := parseMeta(videoId, query_string)

	if err != nil {
		return Video{}, err
	}

	return *meta, nil
}

func (video *Video) Download(index int, filename string, option *Option, videoId string) error {
	var (
		out    *os.File
		err    error
		offset int64
		length int64
	)

	if option.Resume {
		// Resume download from last known offset
		flags := os.O_WRONLY | os.O_CREATE
		out, err = os.OpenFile(filename, flags, 0644)
		if err != nil {
			return fmt.Errorf("Unable to open file %q: %s", filename, err)
		}
		offset, err = out.Seek(0, os.SEEK_END)
		if err != nil {
			return fmt.Errorf("Unable to seek file %q: %s", filename, err)
		}
		fmt.Printf("Resuming from offset %d (%s)\n", offset, abbr(offset))

	} else {
		// Start new download
		flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		out, err = os.OpenFile(filename, flags, 0644)
		if err != nil {
			return fmt.Errorf("Unable to write to file %q: %s", filename, err)
		}
	}

	defer out.Close()

	url := video.Formats[index].Url
	video.Filename = filename

	// Get video content length
	if resp, err := http.Head(url); err != nil {
		return fmt.Errorf("Head request failed: %s", err)
	} else {
		if resp.StatusCode == 403 {
			return errors.New("Head request failed: Video is 403 forbidden")
		}

		if size := resp.Header.Get("Content-Length"); len(size) == 0 {
			return errors.New("Content-Length header is missing")
		} else if length, err = strconv.ParseInt(size, 10, 64); err != nil {
			return fmt.Errorf("Invalid Content-Length: %s", err)
		}

		if length <= offset {
			fmt.Println("Video file is already downloaded.")
			return nil
		}
	}

	if length > 0 {
		go printProgress(out, offset, length)
	}

	// Not using range requests by default, because Youtube is throttling
	// download speed. Using a single GET request for max speed.
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Request failed: %s", err)
	}
	defer resp.Body.Close()

	if length, err = io.Copy(out, resp.Body); err != nil {
		return err
	}

	// Download stats
	duration := time.Now().Sub(start)
	speed := float64(length) / float64(duration/time.Second)
	if duration > time.Second {
		duration -= duration % time.Second
	} else {
		speed = float64(length)
	}

	if option.Rename {
		// Rename output file using video title
		wspace := regexp.MustCompile(`\W+`)
		fname := strings.Split(filename, ".")[0]
		ext := filepath.Ext(filename)
		title := wspace.ReplaceAllString(video.Title, "-")
		if len(title) > 64 {
			title = title[:64]
		}
		title = strings.TrimRight(strings.ToLower(title), "-")
		video.Filename = fmt.Sprintf("%s-%s%s", fname, title, ext)
		if err := os.Rename(filename, video.Filename); err != nil {
			fmt.Println("Failed to rename output file:", err)
		}
	}

	// Extract audio from downloaded video using ffmpeg
	if option.Mp3 {
		if err := out.Close(); err != nil {
			fmt.Println("Error:", err)
		}
		ffmpeg, err := exec.LookPath("ffmpeg")
		if err != nil {
			fmt.Println("ffmpeg not found")
		} else {
			fmt.Println("Extracting audio ..")
			fname := video.Filename
			mp3 := strings.TrimRight(fname, filepath.Ext(fname)) + ".mp3"
			cmd := exec.Command(ffmpeg, "-y", "-loglevel", "quiet", "-i", fname, "-vn", mp3)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Failed to extract audio:", err)
			} else {
				fmt.Println()
				fmt.Println("Extracted audio:", mp3)
			}
		}
	}

	fmt.Printf("Download duration: %s\n", duration)
	fmt.Printf("Average speed: %s/s\n", abbr(int64(speed)))

	return nil
}

func abbr(byteSize int64) string {
	size := float64(byteSize)
	switch {
	case size > GB:
		return fmt.Sprintf("%.1fGB", size/GB)
	case size > MB:
		return fmt.Sprintf("%.1fMB", size/MB)
	case size > KB:
		return fmt.Sprintf("%.1fKB", size/KB)
	}
	return fmt.Sprintf("%d", byteSize)
}

// Measure download speed using output file offset
func printProgress(out *os.File, offset, length int64) {
	var clear string
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	start := time.Now()
	tail := offset

	var err error
	for now := range ticker.C {
		duration := now.Sub(start)
		duration -= duration % time.Second
		offset, err = out.Seek(0, os.SEEK_CUR)
		if err != nil {
			return
		}
		speed := offset - tail
		percent := int(100 * offset / length)
		progress := fmt.Sprintf(
			"%s%s\t %s/%s\t %d%%\t %s/s",
			clear, duration, abbr(offset), abbr(length), percent, abbr(speed))
		fmt.Println(progress)
		tail = offset
		if tail >= length {
			break
		}
		if clear == "" {
			switch runtime.GOOS {
			case "darwin":
				clear = "\033[A\033[2K\r"
			case "linux":
				clear = "\033[A\033[2K\r"
			case "windows":
			}
		}
	}
}
