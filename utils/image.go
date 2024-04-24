package utils

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/qr"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"image"
	"image/draw"
	"image/png"
	"os"
)

const (
	LabelW              = 1024
	LabelH              = 551
	LabelQRL            = 236
	LabelBCL            = 800
	LabelBCH            = 120
	LabelMargin         = 20
	LabelTextLineHeight = 70
)

func barcodeSubtitle(bc barcode.Barcode) image.Image {
	fontFace := basicfont.Face7x13
	//fontColor := color.RGBA{0, 0, 0, 255}
	//margin := 5

	bounds, _ := font.BoundString(fontFace, bc.Content())

	widthTxt := int((bounds.Max.X - bounds.Min.X) / 64)
	//heightTxt := int((bounds.Max.Y - bounds.Min.Y) / 64)

	width := widthTxt
	if bc.Bounds().Dx() > width {
		width = bc.Bounds().Dx()
	}
	//height := heightTxt + bc.Bounds().Dy() + margin
	height := bc.Bounds().Dy()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, image.Rect(0, 0, bc.Bounds().Dx(), bc.Bounds().Dy()), bc, bc.Bounds().Min, draw.Over)
	return img
}

// 生成条形码
func GenBarcode(bcFile, content string, width, height int) error {
	barCode, _ := code128.Encode(content)
	scaled, _ := barcode.Scale(barCode, width, height)
	img := barcodeSubtitle(scaled)
	f, _ := os.Create(bcFile)
	defer f.Close()
	return png.Encode(f, img)
}

// 生成二维码
func GenQRcode(qrFile, content string, length int) error {
	qrCode, _ := qr.Encode(content, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, length, length)
	file, _ := os.Create(qrFile)
	defer file.Close()
	return png.Encode(file, qrCode)
}
