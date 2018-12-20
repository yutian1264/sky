package imgUtils

import (
	"image"
	"os"
	"log"
	"image/png"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode"
	"strings"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/yutian1264/sky/com/utils/timeUtils"
	"github.com/yutian1264/sky/com/utils"
	"net/http"
	"strconv"
	"io/ioutil"
	"encoding/json"
)

type Images struct {
	Name       string
	City       string
	CreateTime string
	Path       string
	Lat        float64
	Lng        float64
	Angel      int
}

func writePng(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	return err
}

/**
	生成二维码
*/
func CreateQRCode(baseUrl, path, filename string, w, h int) error {

	b, err := utils.PathNotExistsCreate(path)
	if !b {
		log.Println(err)
		return err
	}

	code, err := qr.Encode(baseUrl, qr.L, qr.Unicode)
	// code, err := code39.Encode(base64)
	if err != nil {
		log.Println(err)
		return err
	}

	if baseUrl != code.Content() {
		log.Println("data differs")
	}

	code, err = barcode.Scale(code, w, h)
	if err != nil {
		log.Println(err)
		return err
	}

	err = writePng(path+"/"+filename, code)

	return err
}

/**
	获取图片 exif 信息 :时间,经纬度,图片旋转角度
*/
func GetExifMess(imgs []string) []Images {
	var imgList = make([]Images, 0, 10)
	if len(imgs) > 0 {
		for _, s := range imgs {
			var images = Images{}
			//fmt.Println(s)
			images.Path = s
			ps := strings.Split(s, "/")
			images.Name = ps[len(ps)-1]
			f1, err := os.Open(s)
			if err != nil {
				recover()
				f1.Close()
				//log.Println("error:",s)
				imgList = append(imgList, images)
				continue
			}
			defer f1.Close()
			exif, err := exif.Decode(f1)
			if err != nil {
				log.Println("Decode:", s)
				recover()
				f1.Close()
				continue
			}

			lat, lng, err := exif.LatLong()
			if err != nil {
				log.Println("LatLong:", err)
				images.Lat = 0
				images.Lng = 0
			} else {
				images.Lat = lat
				images.Lng = lng
			}
			angel := 0
			tag, e := exif.Get("Orientation")
			images.Angel = 360
			if e == nil {
				orientation := tag.Type
				if orientation == 0 || orientation == 1 {
					angel = 360
				} else if orientation == 3 {
					angel = 180
				} else if orientation == 6 {
					angel = 90
				} else if orientation == 8 {
					angel = 270
				}
				images.Angel = angel

			}

			a, e := exif.DateTime()
			a.String()
			if e != nil {
				recover()
				f1.Close()
				//log.Println("DateTime:",s)
				imgList = append(imgList, images)
				continue
			}
			t := timeUtils.Time2String(a)
			images.CreateTime = t
			imgList = append(imgList, images)
		}

	}
	return imgList
}

/**
	根据经纬度获取当前城市信息
*/

func GetCityByPoint(lati, longtu float64) map[string]interface{} {
	lat := strconv.FormatFloat(lati, 'f', 30, 32)
	lnt := strconv.FormatFloat(longtu, 'f', 30, 32)
	res, e := http.Get("http://api.map.baidu.com/geocoder/v2/?location=" + lat + "," + lnt + "&output=json&pois=1&ak=t9pXAGuuFGrB2tw4mcxBrYBj5kajm7un")

	if res.StatusCode != 200 || e != nil {
		log.Println(e)
		return make(map[string]interface{})
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return make(map[string]interface{})
	}
	mapResult := make(map[string]interface{})
	err = json.Unmarshal(result, &mapResult)
	if err != nil {
		return make(map[string]interface{})
	}
	return mapResult
}
