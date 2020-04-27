package main

import (
	"bytes"
	"fmt"
	"github.com/nikitawootten/cmsc483-project/common"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
)



func helloWorld(w http.ResponseWriter, _ *http.Request) {
	log.Println("New request!")
	_, err := fmt.Fprintf(w, "hi there!\n")


	f, err := os.Open("./simple-http-server/img/pitt.jpg")
	if err != nil {
		log.Fatal(err)
	}


	imgIn, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	resizedImg := resizerFunc(imgIn, 640, 480)
	bytesOfImg := imgToBytes(resizedImg)
	fmt.Fprintf(w, "Resizing image to 640x480")
	err = ioutil.WriteFile("./simple-http-server/img/pitt640x480.jpg", bytesOfImg, 777)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
}


func resizerFunc(imgIn image.Image, len int, wid int) image.Image {
	minXVal := imgIn.Bounds().Min.X
	minYVal := imgIn.Bounds().Min.Y
	maxXVal := imgIn.Bounds().Max.X
	maxYVal := imgIn.Bounds().Max.Y
	for (maxXVal-minXVal) % len != 0 {
		maxXVal--
	}
	for (maxYVal-minYVal) % wid != 0 {
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
	scale := 1.0 / float64((maxX-minX) * (maxY-minY))

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


func imgToBytes(imgIn image.Image) []byte {
	var optimize jpeg.Options
	optimize.Quality = 100
	newBuffer := bytes.NewBuffer(nil)
	err := jpeg.Encode(newBuffer, imgIn, &optimize)
	if err != nil {
		log.Fatal(err)
	}

	return newBuffer.Bytes()
}


func main() {
	req, lbs, address, err := common.ParseFlagsClient()
	if err != nil {
		log.Fatal("Failed to parse args:", err)
	}
	common.ConnectToParentLBs(req, lbs)

	http.HandleFunc("/hello", helloWorld)

	fs := http.FileServer(http.Dir("./simple-http-server/img"))

	http.Handle("/imgs", http.StripPrefix("/imgs", fs))


	log.Println("Mapped routes, listening on ", address)

	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
