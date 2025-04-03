// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	"image/color" 
)

// Grayscale applies a grayscale filtering effect to the image
func (img *Image) Grayscale() {

	// Bounds returns defines the dimensions of the image. Always
	// use the bounds Min and Max fields to get out the width
	// and height for the image
	bounds := img.out.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//Returns the pixel (i.e., RGBA) value at a (x,y) position
			// Note: These get returned as int32 so based on the math you'll
			// be performing you'll need to do a conversion to float64(..)
			r, g, b, a := img.in.At(x, y).RGBA()

			//Note: The values for r,g,b,a for this assignment will range between [0, 65535].
			//For certain computations (i.e., convolution) the values might fall outside this
			// range so you need to clamp them between those values.
			greyC := clamp(float64(r+g+b) / 3)

			//Note: The values need to be stored back as uint16 (I know weird..but there's valid reasons
			// for this that I won't get into right now).
			img.out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
}

func (img *Image) Convolution(kernel []float64) {
    bounds := img.out.Bounds()
    kernelSize := 3
    padding := 1
    var r, g, b, a uint32
    var rSum, gSum, bSum float64

	// Steps :
	// 1. Iterate over (y, x) Image dimensions, (ky, kx) Kernel dimensions
	// 2. Perform same padding convolution
	// 3. Write to Image Out
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            rSum, gSum, bSum = 0, 0, 0

            for ky := -padding; ky <= padding; ky++ {
                for kx := -padding; kx <= padding; kx++ {
                    imgX := x + kx
                    imgY := y + ky

                    if imgX >= bounds.Min.X && imgX < bounds.Max.X && imgY >= bounds.Min.Y && imgY < bounds.Max.Y {
                        r, g, b, a = img.in.At(imgX, imgY).RGBA()
                        kernelValue := kernel[(ky+padding)*kernelSize+(kx+padding)]
                        rSum += float64(r) * kernelValue
                        gSum += float64(g) * kernelValue
                        bSum += float64(b) * kernelValue
                    }
                }
            }

            img.out.Set(x, y, color.RGBA64{clamp(rSum), clamp(gSum), clamp(bSum), uint16(a)})
        }
    }
}


func (img *Image) Sharpen() {
	k := []float64 {0,-1,0,-1,5,-1,0,-1,0}
	img.Convolution(k)
}


func (img *Image) EdgeDetection() {
	k := []float64 {-1,-1,-1,-1,8,-1,-1,-1,-1}
	img.Convolution(k)
}

func (img *Image) Blur() {
	k := []float64 {1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}

	img.Convolution(k)
}


func (img *Image) RunEffects(effects []string){
	// Steps: 
	// 1. Iterate over Effects
	// 2. Execute effect using Convolution
	// 3. Image in  = Previous Image out
	for i := 0; i < len(effects); i++ {
		effect := effects[i]
		if (effect == "S") {
			img.Sharpen()
		} else if (effect == "E") {
			img.EdgeDetection()
		} else if (effect == "B") {
			img.Blur()
		} else if (effect == "G") {
			img.Grayscale()
		} else {
			panic("Incorrect Effect")
		}

		img.in = img.out
	}
}