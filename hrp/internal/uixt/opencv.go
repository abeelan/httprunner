//go:build opencv

package uixt

import (
	"bytes"
	"image"
	"io/ioutil"
	"os"

	cvHelper "github.com/electricbubble/opencv-helper"
)

const (
	// TmSqdiff maps to TM_SQDIFF
	TmSqdiff TemplateMatchMode = iota
	// TmSqdiffNormed maps to TM_SQDIFF_NORMED
	TmSqdiffNormed
	// TmCcorr maps to TM_CCORR
	TmCcorr
	// TmCcorrNormed maps to TM_CCORR_NORMED
	TmCcorrNormed
	// TmCcoeff maps to TM_CCOEFF
	TmCcoeff
	// TmCcoeffNormed maps to TM_CCOEFF_NORMED
	TmCcoeffNormed
)

type DebugMode int

const (
	// DmOff no output
	DmOff DebugMode = iota
	// DmEachMatch output matched and mismatched values
	DmEachMatch
	// DmNotMatch output only values that do not match
	DmNotMatch
)

// Extend 获得扩展后的 Driver，
// 并指定匹配阀值，
// 获取当前设备的 Scale，
// 默认匹配模式为 TmCcoeffNormed，
// 默认关闭 OpenCV 匹配值计算后的输出
func (dExt *DriverExt) extendOpenCV(threshold float64, matchMode ...TemplateMatchMode) (err error) {
	if dExt.scale, err = dExt.Scale(); err != nil {
		return err
	}

	if len(matchMode) == 0 {
		matchMode = []TemplateMatchMode{TmCcoeffNormed}
	}
	dExt.matchMode = matchMode[0]
	cvHelper.Debug(cvHelper.DebugMode(DmOff))
	dExt.threshold = threshold
	return nil
}

func (dExt *DriverExt) Debug(dm DebugMode) {
	cvHelper.Debug(cvHelper.DebugMode(dm))
}

func (dExt *DriverExt) OnlyOnceThreshold(threshold float64) (newExt *DriverExt) {
	newExt = new(DriverExt)
	newExt.WebDriver = dExt.WebDriver
	newExt.scale = dExt.scale
	newExt.matchMode = dExt.matchMode
	newExt.threshold = threshold
	return
}

func (dExt *DriverExt) OnlyOnceMatchMode(matchMode TemplateMatchMode) (newExt *DriverExt) {
	newExt = new(DriverExt)
	newExt.WebDriver = dExt.WebDriver
	newExt.scale = dExt.scale
	newExt.matchMode = matchMode
	newExt.threshold = dExt.threshold
	return
}

// func (sExt *DriverExt) findImgRect(search string) (rect image.Rectangle, err error) {
// 	pathSource := filepath.Join(sExt.pathname, cvHelper.GenFilename())
// 	if err = sExt.driver.ScreenshotToDisk(pathSource); err != nil {
// 		return image.Rectangle{}, err
// 	}
//
// 	if rect, err = cvHelper.FindImageRectFromDisk(pathSource, search, float32(sExt.Threshold), cvHelper.TemplateMatchMode(sExt.MatchMode)); err != nil {
// 		return image.Rectangle{}, err
// 	}
// 	return
// }

func (dExt *DriverExt) FindAllImageRect(search string) (rects []image.Rectangle, err error) {
	var bufSource, bufSearch *bytes.Buffer
	if bufSearch, err = getBufFromDisk(search); err != nil {
		return nil, err
	}
	if bufSource, err = dExt.takeScreenShot(); err != nil {
		return nil, err
	}

	if rects, err = cvHelper.FindAllImageRectsFromRaw(bufSource, bufSearch, float32(dExt.threshold), cvHelper.TemplateMatchMode(dExt.matchMode)); err != nil {
		return nil, err
	}
	return
}

func (dExt *DriverExt) FindImageRectInUIKit(imagePath string) (x, y, width, height float64, err error) {
	var bufSource, bufSearch *bytes.Buffer
	if bufSearch, err = getBufFromDisk(imagePath); err != nil {
		return 0, 0, 0, 0, err
	}
	if bufSource, err = dExt.takeScreenShot(); err != nil {
		return 0, 0, 0, 0, err
	}

	var rect image.Rectangle
	if rect, err = cvHelper.FindImageRectFromRaw(bufSource, bufSearch, float32(dExt.threshold), cvHelper.TemplateMatchMode(dExt.matchMode)); err != nil {
		return 0, 0, 0, 0, err
	}

	// if rect, err = dExt.findImgRect(search); err != nil {
	// 	return 0, 0, 0, 0, err
	// }
	x, y, width, height = dExt.MappingToRectInUIKit(rect)
	return
}

func getBufFromDisk(name string) (*bytes.Buffer, error) {
	var f *os.File
	var err error
	if f, err = os.Open(name); err != nil {
		return nil, err
	}
	var all []byte
	if all, err = ioutil.ReadAll(f); err != nil {
		return nil, err
	}
	return bytes.NewBuffer(all), nil
}

func (dExt *DriverExt) MappingToRectInUIKit(rect image.Rectangle) (x, y, width, height float64) {
	x, y = float64(rect.Min.X)/dExt.scale, float64(rect.Min.Y)/dExt.scale
	width, height = float64(rect.Dx())/dExt.scale, float64(rect.Dy())/dExt.scale
	return
}
