/* 
 * The MIT License (MIT)
 * 
 * Copyright (c) 2016 tony<wuhaiyang1213@gmail.com>
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.

 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
package main

import (
	"net/http"
	"flag"
	"log"
	"path/filepath"
	"os"
	"sync"
	"io"
	"encoding/json"
	"github.com/mind-stack-cn/golang-fileserver/model"
	"fmt"
	"github.com/mind-stack-cn/golang-pipVideo-mixAudio/download"
	"github.com/mind-stack-cn/golang-pipVideo-mixAudio/mixutil"
	"github.com/mind-stack-cn/golang-fileserver/handle"
	"github.com/streamrail/concurrent-map"
	"io/ioutil"
	"strings"
)

var (
	dir string
	wg sync.WaitGroup
)

func init() {
	flag.StringVar(&dir, "dir", ".", "Directory path to save all files.")
	flag.Parse()
}

func main() {
	// If no path is passed to app, normalize to path formath
	if dir == "." {
		dir, _ = filepath.Abs(dir)
		dir += "/data/"
	}

	if _, err := os.Stat(dir); err != nil {
		log.Printf("Directory %s not exist, Create it", dir)
		errPath := os.MkdirAll(dir, 0777)
		if errPath != nil {
			log.Fatalf("Directory %s not exist, Create it Fail", dir)
			return
		}
	}

	server := http.Server{
		Addr:    ":8089",
		Handler: &requestHandler{},
	}
	server.ListenAndServe()
}

type requestHandler struct{}

type PostParams struct {
	AudioUri0   string	// 待叠加音频0
	AudioUri1   string	// 待叠加音频1
	VideoUri0   string	// 待叠加视频0
	VideoUri1   string	// 待叠加视频1
	CallBackUrl string	// 回调请求url,POST
}

type MixedResult struct {
	MixedAudio interface{}
	MixedVideo interface{}
}

func (*requestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var params PostParams
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	noAudioMix := params.AudioUri0 == "" || params.AudioUri1 == ""
	noVideoMix := params.VideoUri0 == "" || params.VideoUri1 == ""

	if noAudioMix && noVideoMix{
		http.Error(w, "BadRequest, audio or video param error", http.StatusBadRequest)
		return
	}

	log.Printf("audio0: %s, audio1: %s, video0: %s, video1: %s, callback: %s",
		params.AudioUri0, params.AudioUri1, params.VideoUri0, params.VideoUri1, params.CallBackUrl)
	go DownloadFilesAndMix(!noAudioMix, !noVideoMix, params)

	io.WriteString(w, "Golang pipVideo mixAudio Server, working...")
}


func DownloadFilesAndMix(forAudio bool, forVide bool, postParams PostParams) {
	var mixedResult MixedResult
	if forAudio {
		audioFilePath0, audioFilePath1, err := DownloadFiles(postParams.AudioUri0, postParams.AudioUri1)
		if err != nil {
			InvokeCallBack(postParams.CallBackUrl, fmt.Sprintf("download audio file error %s", err.Error()))
			return
		}

		absoluteFilePath, relatedFilePath, err := handle.GenerateNewFilePath(dir, filepath.Ext(audioFilePath0))
		if err != nil {
			InvokeCallBack(postParams.CallBackUrl, fmt.Sprintf("generate audio file error %s", err.Error()))
			return
		}

		mixErr := mixutil.MixAudios([]string{audioFilePath0, audioFilePath1}, absoluteFilePath)
		if mixErr != nil {
			InvokeCallBack(postParams.CallBackUrl, fmt.Sprintf("Mix audio file error %s", mixErr.Error()))
			return
		}

		mixedAudioFileInfo := model.ResFileFromFileName(absoluteFilePath, relatedFilePath, model.FileTypeAudio)
		mixedResult.MixedAudio = mixedAudioFileInfo
	}

	if forVide {
		videoFilePath0, videoFilePath1, err := DownloadFiles(postParams.VideoUri0, postParams.VideoUri1)
		if err != nil {
			InvokeCallBack(postParams.CallBackUrl, fmt.Sprintf("download video file error %s", err.Error()))
			return

		}

		absoluteFilePath, relatedFilePath, err := handle.GenerateNewFilePath(dir, filepath.Ext(videoFilePath0))
		if err != nil {
			InvokeCallBack(postParams.CallBackUrl, fmt.Sprintf("generate video file error %s", err.Error()))
			return
		}

		mixErr := mixutil.MixVideos(videoFilePath0, videoFilePath1, absoluteFilePath)
		if mixErr != nil {
			InvokeCallBack(postParams.CallBackUrl, fmt.Sprintf("Mix video file error %s", mixErr.Error()))
			return
		}

		mixedVideoFileInfo := model.ResFileFromFileName(absoluteFilePath, relatedFilePath, model.FileTypeVideo)
		mixedResult.MixedVideo = mixedVideoFileInfo
	}

	json, err := json.Marshal(mixedResult)
	if err != nil {
		InvokeCallBack(postParams.CallBackUrl, fmt.Sprintf("Mix files json error %s", err.Error()))
		return
	}

	InvokeCallBack(postParams.CallBackUrl, string(json))
}

func DownloadFiles(fileUris0 string, fileUris1 string) (string, string, error) {
	var filesUriFilePathMap = cmap.New()
	var filesUriErrorMap  = cmap.New()

	var fileUris []string
	fileUris = append(fileUris, fileUris0)
	fileUris = append(fileUris, fileUris1)
	// Start workers to grab the file only if the container not empty
	if len(fileUris) >= 1 {
		// number of workers depends on number of files
		for _, fileUri := range fileUris {
			wg.Add(1)
			// Put downloader process into another thread
			// for each file.
			go func(fileUri string) {
				defer wg.Done()
				absoluteFileName, err := download.Download(dir, fileUri)
				if err == nil {
					filesUriFilePathMap.Set(fileUri, absoluteFileName)
				} else {
					filesUriErrorMap.Set(fileUri, err.Error())
				}
			}(fileUri)
		}
	}

	// wait for all channels until they finish their jobs
	wg.Wait()

	err0, hasError0 :=filesUriErrorMap.Get(fileUris0);
	err1, hasError1 :=filesUriErrorMap.Get(fileUris1);

	if hasError0 || hasError1 {
		return "", "", fmt.Errorf("%s,%s", err0.(string), err1.(string))
	} else {
		fileAbs0, _ :=filesUriFilePathMap.Get(fileUris0);
		fileAbs1, _ :=filesUriFilePathMap.Get(fileUris1);
		return fileAbs0.(string), fileAbs1.(string), nil
	}
}

func InvokeCallBack(url string, result string)  {
	if url == "" {
		fmt.Println(result)
	} else {
		req, err := http.NewRequest("POST", url, strings.NewReader(result))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
}


