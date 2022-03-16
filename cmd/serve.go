/*
Copyright © 2022 Zhj Rong <rongzhj2020@163.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
	"syscall"

	"github.com/cwchiu/go-winapi"
	"github.com/kbinani/screenshot"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: run,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serveCmd.Flags().Int("port", 8080, "Port")
}

func handler(w http.ResponseWriter, r *http.Request) {
	all := findRect()
	img, err := screenshot.CaptureRect(all)
	if err != nil {
		panic(err)
	}

	writeImage(w, img)
}
func findRect() image.Rectangle {
	p, err := syscall.UTF16PtrFromString("WeChatMainWndForPC")
	if err != nil {
		panic(err)
	}

	hwnd := winapi.FindWindow(p, nil)
	// fmt.Println(hwnd)

	r := winapi.RECT{}
	ok := winapi.GetWindowRect(hwnd, &r)
	if !ok {
		panic("no rect")
	}

	// fmt.Println(r, r.Bottom, r.Left, r.Right, r.Top)

	// all := image.Rect(int(r.Left), int(r.Top), int(r.Right), int(r.Bottom))
	all := image.Rect(int(r.Left), int(r.Top), int(r.Left+150), int(r.Top+300))

	return all
}

func writeImage(w http.ResponseWriter, img image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func run(cmd *cobra.Command, args []string) {
	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		panic(err)
	}
}
