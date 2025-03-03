package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
)

var (
	geocodeUrl        = "https://geocode.maps.co/search?q="
	reverseGeocodeUrl = "https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=%s&longitude=%s&localityLanguage=en"

	defaultBitSize = 32
	earthRadius    = 6378.137
	aa             = 6378245.0
	ee             = 0.00669342162296594323

	inChina = [][]float64{
		{49.220400, 79.446200, 42.889900, 96.330000},
		{54.141500, 109.687200, 39.374200, 135.000200},
		{42.889900, 073.124600, 29.529700, 124.143255},
		{29.529700, 082.968400, 26.718600, 097.035200},
		{29.529700, 097.025300, 20.414096, 124.367395},
		{20.414096, 107.975793, 17.871542, 111.744104},
	}
	outChina = [][]float64{
		{25.398623, 119.921265, 21.785006, 122.497559},
		{22.284000, 101.865200, 20.098800, 106.665000},
		{21.542200, 106.452500, 20.487800, 108.051000},
		{55.817500, 109.032300, 50.325700, 119.127000},
		{55.817500, 127.456800, 49.557400, 137.022700},
		{44.892200, 131.266200, 42.569200, 137.022700},
	}
)

type GeocodeRes struct {
	Lat         string `json:"lat"`
	Lng         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

type ReverseGeocodeRes struct {
	CountryName          string               `json:"countryName"`
	PrincipalSubdivision string               `json:"principalSubdivision"`
	City                 string               `json:"city"`
	Locality             string               `json:"locality"`
	LocalityInfo         *ReverseLocalityInfo `json:"localityInfo"`
}
type ReverseLocalityInfo struct {
	Administrative []*ReverseAdministrative `json:"administrative"`
}
type ReverseAdministrative struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Order       uint8  `json:"order"`
}

func Geocode(address string) ([]*GeocodeRes, error) {
	ret, err := HttpGet(geocodeUrl + url.QueryEscape(address))
	if err != nil {
		return nil, err
	}
	data := []*GeocodeRes{}
	err = json.Unmarshal(ret.Body(), &data)
	return data, err
}

// 地址转地理坐标
func Addr2Coord(address string) (*GeocodeRes, error) {
	ret, err := Geocode(address)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, errors.New("No Results")
	}
	return ret[0], nil
}

func ReverseGeocode(lat, lng string) (*ReverseGeocodeRes, error) {
	ret, err := HttpGet(fmt.Sprintf(reverseGeocodeUrl, lat, lng))
	if err != nil {
		return nil, err
	}
	var data *ReverseGeocodeRes
	err = json.Unmarshal(ret.Body(), &data)
	return data, err
}

// 坐标转地址
func Coord2Addr(lat, lng string) (addr string) {
	data, err := ReverseGeocode(lat, lng)
	if err != nil {
		return addr
	}
	spliter := ", "
	if data.LocalityInfo != nil && len(data.LocalityInfo.Administrative) != 0 {
		for _, l := range data.LocalityInfo.Administrative {
			addr = spliter + l.Name + addr
		}
	} else {
		if data.Locality != "" {
			addr += spliter + data.Locality
		}
		if data.City != "" {
			addr += spliter + data.City
		}
		if data.PrincipalSubdivision != "" {
			addr += spliter + data.PrincipalSubdivision
		}
		if data.CountryName != "" {
			addr += spliter + data.CountryName
		}
	}
	return strings.TrimLeft(addr, spliter)
}

func rad(point string) float64 {
	point2, _ := strconv.ParseFloat(point, defaultBitSize)
	return point2 * math.Pi / 180.0
}

// 根据两点坐标计算距离(单位米)
func CalcDistance(lat1, lng1, lat2, lng2 string) float64 {
	radLat1 := rad(lat1)
	radLat2 := rad(lat2)
	radLng1 := rad(lng1)
	radLng2 := rad(lng2)
	cos := math.Cos(radLat1)*math.Cos(radLat2)*math.Cos(radLng2-radLng1) + math.Sin(radLat1)*math.Sin(radLat2)
	return math.Round(math.Acos(cos) * earthRadius * 1000)
}

func Ddmm2dd(point string) string {
	val, _ := strconv.ParseFloat(point, defaultBitSize)
	val2 := val / 100
	intVal := math.Floor(val2)
	ret := intVal + (val2-intVal)*1000000/60/10000
	return fmt.Sprintf("%f", ret)
}

func TransformLat(lat, lng float64) float64 {
	ret := -100.0 + 2.0*lng + 3.0*lat + 0.2*lat*lat + 0.1*lng*lat + 0.2*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*math.Pi) + 20.0*math.Sin(2.0*lng*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*math.Pi) + 40.0*math.Sin(lat/3.0*math.Pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(lat/12.0*math.Pi) + 320*math.Sin(lat*math.Pi/30.0)) * 2.0 / 3.0
	return ret
}

func TransformLng(lat, lng float64) float64 {
	ret := 300.0 + lng + 2.0*lat + 0.1*lng*lng + 0.1*lng*lat + 0.1*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*math.Pi) + 20.0*math.Sin(2.0*lng*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lng*math.Pi) + 40.0*math.Sin(lng/3.0*math.Pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(lng/12.0*math.Pi) + 300.0*math.Sin(lng/30.0*math.Pi)) * 2.0 / 3.0
	return ret
}

// 坐标偏移
func Wgs84ToGcj02(lat0, lng0 string) (string, string) {
	lat, _ := strconv.ParseFloat(lat0, defaultBitSize)
	lng, _ := strconv.ParseFloat(lng0, defaultBitSize)
	dlat := TransformLat(lat-35.0, lng-105.0)
	dlng := TransformLng(lat-35.0, lng-105.0)
	radlat := lat / 180.0 * math.Pi
	magic := math.Sin(radlat)
	magic = 1 - ee*magic*magic
	sqrtmagic := math.Sqrt(magic)
	dlat = (dlat * 180.0) / ((aa * (1 - ee)) / (magic * sqrtmagic) * math.Pi)
	dlng = (dlng * 180.0) / (aa / sqrtmagic * math.Cos(radlat) * math.Pi)
	glat := lat + dlat
	glon := lng + dlng
	return fmt.Sprintf("%f", glat), fmt.Sprintf("%f", glon)
}

// 坐标转换
func CalcCoord(lat, lng string) (string, string) {
	newLat := Ddmm2dd(lat)
	newLng := Ddmm2dd(lng)

	//todo 国内需要处理坐标偏移
	//if global.GVA_CONFIG.AliyunAPI.IOT.Localization == "gcj02" { // 国内需要处理坐标偏移
	//	return wgs84ToGcj02(newLat, newLng)
	//}
	return newLat, newLng
}

func GsmSignal(val string) int {
	gsm, _ := strconv.Atoi(val)
	if gsm > 25 {
		return 5
	} else if gsm > 19 && gsm <= 25 {
		return 4
	} else if gsm > 13 && gsm <= 19 {
		return 3
	} else if gsm > 7 && gsm <= 13 {
		return 2
	} else if gsm > 1 && gsm <= 7 {
		return 1
	} else {
		return 0
	}
}

func GpsSignal(val string) int {
	gps, _ := strconv.Atoi(val)
	if gps > 8 {
		return 5
	} else if gps > 5 {
		return 4
	} else if gps == 5 {
		return 3
	} else if gps == 4 {
		return 2
	} else if gps == 3 {
		return 1
	} else {
		return 0
	}
}

func Celsius2FahrenheitF(val float64) float64 {
	if val == 0 {
		return val
	}
	str := fmt.Sprintf("%0.1f", val*1.8+32)
	ret, _ := strconv.ParseFloat(str, 64)
	return ret
}

func Celsius2Fahrenheit(val string) string {
	if val == "0" {
		return val
	}
	temp, _ := strconv.ParseFloat(val, 64)
	return fmt.Sprintf("%0.1f", temp*1.8+32)
}

func Battery2Distance(battery uint) uint {
	return uint(Meter2MileF(float64(battery * 1100)))
}

func Meter2Mile(val string) string {
	if val == "0" {
		return val
	}
	meter, _ := strconv.ParseFloat(val, 64)
	return fmt.Sprintf("%0.1f", meter*0.0006214)
}

func Meter2MileF(val float64) float64 {
	if val == 0 {
		return val
	}
	return val * 0.0006214
}

func Mile2MeterF(val float64) float64 {
	return val / 0.0006214
}

func KiloMeter2Mile(val string) string {
	if val == "0" {
		return val
	}
	meter, _ := strconv.ParseFloat(val, 64)
	return fmt.Sprintf("%0.1f", meter*0.6214)
}

func KiloMeter2MileF(val float64) float64 {
	if val == 0 {
		return val
	}
	return val * 0.6214
}

func CliDuration(seconds uint64) (ret string) {
	if seconds >= 3600 {
		h := seconds / 3600
		remain := seconds % 3600
		m := remain / 60
		s := remain % 60
		hstr := StrPad(fmt.Sprintf("%d", h), 2, "0", "LEFT")
		mstr := StrPad(fmt.Sprintf("%d", m), 2, "0", "LEFT")
		sstr := StrPad(fmt.Sprintf("%d", s), 2, "0", "LEFT")
		ret = fmt.Sprintf(`%sh%s'%s"`, hstr, mstr, sstr)
	} else if seconds >= 60 {
		m := seconds / 60
		s := seconds % 60
		mstr := StrPad(fmt.Sprintf("%d", m), 2, "0", "LEFT")
		sstr := StrPad(fmt.Sprintf("%d", s), 2, "0", "LEFT")
		ret = fmt.Sprintf(`%s'%s"`, mstr, sstr)
	} else {
		sstr := StrPad(fmt.Sprintf("%d", seconds), 2, "0", "LEFT")
		ret = fmt.Sprintf(`00'%s"`, sstr)
	}
	return
}

func IsSoco(deviceName string) bool {
	return !strings.HasPrefix(deviceName, "STSX245") && strings.HasPrefix(deviceName, "S")
}

func IsFPRO(deviceName string) bool {
	return strings.HasPrefix(deviceName, "FPRO")
}

func ShowMileage(deviceName string) bool {
	return true
}

/**
 * 根据霍尔脉冲数_500ms、电机转动一圈的霍尔脉冲数、车轮直径(cm)、时长(s) 计算行驶距离
 * 车轮500ms转动的距离(m) = 霍尔脉冲数_500ms / 电机转动一圈的霍尔脉冲数 * (车轮直径(cm) / 100) * math.Pi
 * 车子秒速 = 车轮500ms转动的距离(m) * 2
 * 距离 = 车子秒速 * 时长
 */

func CalcHallPulse2Distance(motorSpeed0, perHallPulse0, wheelDiameter0 int64, diff0 uint64) float64 {
	motorSpeed := float64(motorSpeed0)
	perHallPulse := float64(perHallPulse0)
	wheelDiameter := float64(wheelDiameter0)
	diff := float64(diff0)
	speed := motorSpeed / perHallPulse * (wheelDiameter / 100) * math.Pi * 2 // 电动车速度 m/s
	return speed * diff
}

// IsInGeoFence 判断坐标是否在电子围栏内 point = [2]float64{lng,lat}
func IsInGeoFence(point [2]float64, points [][2]float64) bool {
	var (
		counter    int
		pointCount = len(points)
		p1         = points[0]
	)

	for i := range points {
		p2 := points[i%pointCount]
		if point[0] > math.Min(p1[0], p2[0]) && point[0] <= math.Max(p1[0], p2[0]) {
			if point[1] <= math.Max(p1[1], p2[1]) {
				if p1[0] != p2[0] {
					xinters := (point[0]-p1[0])*(p2[1]-p1[1])/(p2[0]-p1[0]) + p1[1]
					if p1[1] == p2[1] || point[1] <= xinters {
						counter++
					}
				}
			}
		}
		p1 = p2
	}
	if counter%2 == 0 {
		return false
	} else {
		return true
	}
}

func IsInChina(lat, lng string) bool {
	latFloat, _ := strconv.ParseFloat(lat, 64)
	lngFloat, _ := strconv.ParseFloat(lng, 64)
	return IsInsideChina(latFloat, lngFloat)
}

func IsInsideChina(lat, lng float64) bool {
	for i := 0; i < 6; i++ {
		if lat <= inChina[i][0] && lat >= inChina[i][2] && lng >= inChina[i][1] && lng <= inChina[i][3] {
			for j := 0; j < 6; j++ {
				if lat <= outChina[j][0] && lat >= outChina[j][2] && lng >= outChina[j][1] && lng <= outChina[j][3] {
					return false
				}
			}
			return true
		}
	}
	return false
}

// dd.mmmm转dd.dddd
func Ddmmmm2dd(point string) string {
	val, _ := strconv.ParseFloat(point, defaultBitSize)
	degrees := math.Floor(val)
	minutes := val - degrees
	degreesInDecimal := minutes / 60
	return fmt.Sprintf("%f", degrees+degreesInDecimal)
}
