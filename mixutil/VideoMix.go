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
package mixutil

import (
	"os/exec"
	"fmt"
	"github.com/mind-stack-cn/golang-fileserver/model"
	"path"
	"strings"
)

const thumbExt = ".jpg"

// Mix Two Video
func MixVideos(inputAudioFilePath1 string, inputAudioFilePath2 string, outAudioFilePath string) (string, float64, error){
	if inputAudioFilePath1== "" || inputAudioFilePath2 == "" ||  outAudioFilePath == "" {
		return "", 0, fmt.Errorf("invalidate params")
	}

	// Mix audio
	err := MixVideosImp(inputAudioFilePath1, inputAudioFilePath2, outAudioFilePath)
	if err != nil {
		return "", 0, err
	}

	// Get File Duration
	duration, err := model.GetMediaDuration(outAudioFilePath)
	if err != nil {
		return "", 0, err
	}

	// Generate Thumbnail
	// Thumbnail path
	thumbNailPath := strings.TrimSuffix(outAudioFilePath, path.Ext(outAudioFilePath)) + thumbExt
	// Generate default thumbnail
	if err := model.GetVideoThumbnail(outAudioFilePath, thumbNailPath); err == nil {
		return thumbNailPath, duration, err
	}else{
		return "", duration, err
	}
}

// Use Command Line "ffmpeg" to Mix Audio
// ffmpeg -i [input0] -i [input1] -filter_complex "[1][0]amix=inputs=2:duration=longest;" -strict -2 [output]
func MixVideosImp(inputAudioFilePath1 string, inputAudioFilePath2 string, outAudioFilePath string) error {
	ffmpegcmd:= GetMixVideoCommand(inputAudioFilePath1, inputAudioFilePath2, outAudioFilePath)
	fmt.Println(ffmpegcmd)
	_, err:= exec.Command("sh", "-c", ffmpegcmd).Output()
	return err
}

func GetMixVideoCommand(inputAudioFilePath1 string, inputAudioFilePath2 string, outAudioFilePath string) string{
	filePathInput := fmt.Sprintf("-i %s -i %s", inputAudioFilePath1, inputAudioFilePath2)
	return fmt.Sprintf("ffmpeg %s -filter_complex \"[0]scale=iw/5:ih/5 [pip];[1][0]amix=inputs=2:duration=longest; [1][pip] overlay=main_w-overlay_w-10:main_h-overlay_h-10\" -strict -2 -y %s", filePathInput, outAudioFilePath)
}
