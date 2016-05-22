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
package download

import (
	"io"
	"log"
	"net/http"
	"os"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

//download used to access file via url to get the response status
//and bytes
func Download(base, fileUri string) (string, error) {
	fmt.Println("downloading: " + fileUri)
	resp, errResp := http.Get(fileUri)
	if errResp != nil {
		return "", errResp
	}

	if resp.StatusCode == 200 {
		parsedUri, _ := url.Parse(fileUri)
		fileName := strings.TrimPrefix(fileUri, parsedUri.Scheme + "://")
		return base + fileName, buildFile(base, fileName, resp.Body)
	} else {
		return "", fmt.Errorf("status Code %d", resp.StatusCode)
	}
}

//buildFile used to create file in local disk
func buildFile(base string, file string, content io.ReadCloser) error {
	// close the source here, because this process happened
	// in other thread
	defer content.Close()

	absoluteFilePath := base + file

	// Make dir
	err := os.MkdirAll(filepath.Dir(absoluteFilePath), 0777)
	if err != nil {
		return err;
	}

	// create file
	f, err := os.Create(absoluteFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	//copying file content into disk
	size, errIO := io.Copy(f, content)
	if err != nil {
		return errIO
	}

	log.Printf("%v downloaded with size %v", file, size)
	return nil
}
