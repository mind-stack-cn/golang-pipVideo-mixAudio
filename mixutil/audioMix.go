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
)

// Use Command Line "ffmpeg" to Mix Audio
// ffmpeg -i [input0] -i [input1] -filter_complex "[1][0]amix=inputs=2:duration=longest;" -strict -2 [output]
func MixAudios(inputAudioFilePaths []string, outAudioFilePath string) error {
	// param validate
	if inputAudioFilePaths == nil || len(inputAudioFilePaths) <= 1 || outAudioFilePath == "" {
		return fmt.Errorf("invalidate params")
	}

	ffmpegcmd := GetMixAudioCommand(inputAudioFilePaths, outAudioFilePath)
	fmt.Println(ffmpegcmd)
	_, err := exec.Command("sh", "-c", ffmpegcmd).Output()
	return err
}

func GetMixAudioCommand(inputAudioFilePaths []string, outAudioFilePath string) string {
	filePathInput := ""
	mixInputs := ""
	for i := 0; i < len(inputAudioFilePaths); i++ {
		filePathInput += " -i " + inputAudioFilePaths[i]
		mixInputs += fmt.Sprintf("[%d]", i)
	}
	return fmt.Sprintf("ffmpeg%s -filter_complex \"%samix=inputs=%d:duration=longest\" -strict -2 -y %s", filePathInput, mixInputs, len(inputAudioFilePaths), outAudioFilePath)
}
