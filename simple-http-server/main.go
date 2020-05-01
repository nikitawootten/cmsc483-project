package main

import (
	"fmt"
	"github.com/nikitawootten/cmsc483-project/common"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

func fib(n int) int {
	varOne := 0
	varTwo := 1
	for i := 0; i < n; i++ {
		temp := varOne
		varOne = varTwo
		varTwo = temp + varOne
	}
	return varOne
}

func fibonacciEndpoint(w http.ResponseWriter, _ *http.Request) {
	log.Println("New request - Fibonacci number test")
	for i := 0; i < (rand.Intn(80-50) + 50); i++ {

		fmt.Fprintf(w, strconv.Itoa(fib(i))+" ")

	}

	_, err := fmt.Fprint(w, "\n")
	if err != nil {
		log.Fatal(err)
	}
}

func resizeImageEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("New request - Resize image test")

	imgIn, err := jpeg.Decode(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	resizedImg := resizerFunc(imgIn, 640, 480)

	err = jpeg.Encode(w, resizedImg, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Fatal(err)
	}
}

func helloWorldEndpoint(w http.ResponseWriter, _ *http.Request) {
	log.Println("New request!")

	_, err := fmt.Fprint(w, "Hello there!")
	if err != nil {
		log.Fatal(err)
	}
}

func resizerFunc(imgIn image.Image, len int, wid int) image.Image {
	minXVal := imgIn.Bounds().Min.X
	minYVal := imgIn.Bounds().Min.Y
	maxXVal := imgIn.Bounds().Max.X
	maxYVal := imgIn.Bounds().Max.Y
	for (maxXVal-minXVal)%len != 0 {
		maxXVal--
	}
	for (maxYVal-minYVal)%wid != 0 {
		maxYVal--
	}
	scaleXVal := (maxXVal - minXVal) / len
	scaleYVal := (maxYVal - minYVal) / wid

	imageRectangle := image.Rect(0, 0, len, wid)
	imageRes := image.NewRGBA(imageRectangle)
	draw.Draw(imageRes, imageRes.Bounds(), &image.Uniform{C: color.White}, image.ZP, draw.Src)
	for y := 0; y < wid; y++ {
		for x := 0; x < len; x++ {
			avgCol := avgColor(imgIn, minXVal+x*scaleXVal, minXVal+(x+1)*scaleXVal, minYVal+y*scaleYVal, minYVal+(y+1)*scaleYVal)
			imageRes.Set(x, y, avgCol)
		}
	}
	return imageRes
}

func avgColor(imgIn image.Image, minX int, maxX int, minY int, maxY int) color.Color {
	var avgR, avgG, avgB, avgAlpha float64
	scale := 1.0 / float64((maxX-minX)*(maxY-minY))

	for i := minX; i < maxX; i++ {
		for k := minY; k < maxY; k++ {
			r, g, b, a := imgIn.At(i, k).RGBA()
			avgR += float64(r) * scale
			avgG += float64(g) * scale
			avgB += float64(b) * scale
			avgAlpha += float64(a) * scale
		}
	}

	avgR = math.Sqrt(avgR)
	avgG = math.Sqrt(avgG)
	avgB = math.Sqrt(avgB)
	avgAlpha = math.Sqrt(avgAlpha)

	avgColor := color.RGBA{
		R: uint8(avgR),
		G: uint8(avgG),
		B: uint8(avgB),
		A: uint8(avgAlpha)}

	return avgColor
}

func main() {
	req, lbs, address, _, err := common.ParseFlags(false)
	if err != nil {
		log.Fatal("Failed to parse args:", err)
	}

	heartbeat := common.ClientHeartbeat{}
	connCount := common.NewConnectionCounterFromHeartbeat(&heartbeat)

	http.HandleFunc("/resize", connCount.WrapHttp(resizeImageEndpoint))
	http.HandleFunc("/fib", connCount.WrapHttp(fibonacciEndpoint))
	http.HandleFunc("/hello_world", connCount.WrapHttp(helloWorldEndpoint))

	log.Println("Mapped routes, listening on ", address)

	common.ConnectToParentLBs(req, lbs, &heartbeat)

	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
