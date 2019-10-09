package main

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"image/jpeg"
	"io/ioutil"
	"os"
	"time"

	log "github.com/albrow/prtty"
	"github.com/chromedp/chromedp"

	//"github.com/davecgh/go-spew/spew"

	"./config/db"
	"./model"
)

func main() {
	allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(),
		[]chromedp.ExecAllocatorOption{
			chromedp.NoFirstRun,
			chromedp.NoDefaultBrowserCheck,
			chromedp.Flag("disable-background-networking", true),
			chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
			chromedp.Flag("disable-background-timer-throttling", true),
			chromedp.Flag("disable-backgrounding-occluded-windows", true),
			chromedp.Flag("disable-breakpad", true),
			chromedp.Flag("disable-client-side-phishing-detection", true),
			chromedp.Flag("disable-default-apps", true),
			chromedp.Flag("disable-dev-shm-usage", true),
			chromedp.Flag("disable-features", "site-per-process,TranslateUI,BlinkGenPropertyTrees"),
			chromedp.Flag("disable-hang-monitor", true),
			chromedp.Flag("disable-ipc-flooding-protection", true),
			chromedp.Flag("disable-popup-blocking", true),
			chromedp.Flag("disable-prompt-on-repost", true),
			chromedp.Flag("disable-renderer-backgrounding", true),
			chromedp.Flag("disable-sync", true),
			chromedp.Flag("force-color-profile", "srgb"),
			chromedp.Flag("metrics-recording-only", true),
			chromedp.Flag("safebrowsing-disable-auto-update", true),
			chromedp.Flag("enable-automation", true),
			chromedp.Flag("password-store", "basic"),
			chromedp.Flag("use-mock-keychain", true),
		}...)
	defer cancel1()

	newCtx, cancel2 := chromedp.NewContext(allocCtx)
	defer cancel2()

	if err := chromedp.Run(newCtx); err != nil {
		log.Error.Println(err)
	}

	//games := []model.Game{}
	//db.DB.Where(&model.Game{Vendor: "Netent"}).Find(&games)

	for id := 1; id <= 3283; id++ {
	//for _, game := range games {
		game := model.Game{}
		db.DB.First(&game, id)

		var buf []byte
		if err := chromedp.Run(newCtx,
			chromedp.Navigate(fmt.Sprintf(`https://redirect/?id=%d`, game.Model.ID)),
			chromedp.Sleep(2 * time.Second),
			chromedp.CaptureScreenshot(&buf),
		); err != nil {
			log.Error.Println(err)
		}

		SaveScreenshot(&buf, game.Model.ID)
	}
}

func SaveScreenshot(buf *[]byte, id uint) {

	tmpfile, err := ioutil.TempFile("", "*.png")
	if err != nil {
		log.Error.Println(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(*buf); err != nil {
		log.Error.Println(err)
	}
	tmpfile.Seek(0, 0)

	path := fmt.Sprintf("public/system/src/screenshot_%d.jpg", id)
	file, err := os.Create(path)
	if err != nil {
		log.Error.Println(err)
	}
	defer file.Close()

	img, err := png.Decode(tmpfile)
	if err != nil {
		log.Error.Println(err)
	}

	bounds := img.Bounds()

	black := true
	minX := bounds.Min.X
	minY := bounds.Min.Y
	maxX := bounds.Max.X
	maxY := bounds.Max.Y
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		black = true
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if !(r == 65535 && g == 65535 && b == 65535) {
				black = false
				if minX < x {
					minX = x
				}
				if maxX > x {
					maxX = x
				}
			}
		}

		if !black && maxY > y {
			maxY = y
		}

		if !black && minY < y {
			minY = y
		}
	}

	x1 := minX
	y1 := minY
	x2 := maxX
	y2 := maxY
	for y := y2; y < y1; y++ {
		black = true
		for x := x2; x < x1; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if !(r < 3000 && g < 3000 && b < 3000) {
				black = false
				if minX > x {
					minX = x
				}
				if maxX < x {
					maxX = x
				}
			}
		}

		if !black && maxY < y {
			maxY = y
		}

		if !black && minY > y {
			minY = y
		}
	}

	newimg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(minX, maxY, maxX, minY))

	err = jpeg.Encode(file, newimg, &jpeg.Options{Quality: 90})
	if err != nil {
		log.Error.Println(err)
	}

}
