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

func Test_GetMixVideoCommand(t *testing.T)  {
	var cmd = mixutil.GetMixVideoCommand("1.mp4", "2.mp4", "out.mp4")
	assert.Equal(t, cmd, "ffmpeg -i 1.mp4 -i 2.mp4 -filter_complex \"[0]scale=iw/5:ih/5 [pip];[1][0]amix=inputs=2:duration=longest; [1][pip] overlay=main_w-overlay_w-10:main_h-overlay_h-10\" -strict -2 out.mp4")
}

func Test_MixVideoImp(t *testing.T)  {
	err := mixutil.MixVideosImp("./testmedia/1.mp4", "./testmedia/2.mp4", "out.mp4")
	assert.Nil(t, err)
}

func Test_MixVideos(t *testing.T)  {
	thumbNail, duraion, err := mixutil.MixVideos("./testmedia/1.mp4", "./testmedia/2.mp4", "out.mp4")
	fmt.Println(thumbNail)
	fmt.Println(duraion)
	assert.Nil(t, err)
}
