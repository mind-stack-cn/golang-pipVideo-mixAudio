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
package test

import (
	"testing"
	"github.com/mind-stack-cn/golang-pipVideo-mixAudio/mixutil"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func Test_GetMixAudioCommand(t *testing.T)  {
	var cmd = mixutil.GetMixAudioCommand([]string{"1.aac", "2.aac"}, "out.aac")
	assert.Equal(t, cmd, "ffmpeg -i 1.aac -i 2.aac -filter_complex \"[0][1]amix=inputs=2:duration=longest\" -strict -2 -y out.aac")

	cmd = mixutil.GetMixAudioCommand([]string{"1.aac", "2.aac", "3.aac"}, "out.aac")
	assert.Equal(t, cmd, "ffmpeg -i 1.aac -i 2.aac -i 3.aac -filter_complex \"[0][1][2]amix=inputs=3:duration=longest\" -strict -2 -y out.aac")
}

func Test_MixAudiosImp(t *testing.T)  {
	err := mixutil.MixAudiosImp([]string{"./testmedia/1.aac", "./testmedia/2.aac"}, "out.aac")
	assert.Nil(t, err)
}

func Test_MixAudios(t *testing.T)  {
	duraion, err := mixutil.MixAudios([]string{"./testmedia/1.aac", "./testmedia/2.aac"}, "out.aac")
	fmt.Println(duraion)
	assert.Nil(t, err)
}
