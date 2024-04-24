package utils

import (
	"context"
	"errors"
	"fmt"
	"googlemaps.github.io/maps"
	"strings"
)

var (
	googleDimensionNum = 25
)

func googleMapsClient(apiKey string) (*maps.Client, error) {
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	return c, err
}

// GoogleGeocoding 谷歌坐标转换
func GoogleGeocoding(apiKey, address string) (float64, float64, string, error) {
	c, err := googleMapsClient(apiKey)
	if err != nil {
		return 0, 0, "", err
	}
	res, err := c.Geocode(context.Background(), &maps.GeocodingRequest{
		Address: address,
	})
	if err != nil {
		return 0, 0, "", err
	}
	if len(res) == 0 {
		return 0, 0, "", errors.New("can't find the address")
	}
	var addr string
	// 地址校正为统一格式
	if len(res[0].AddressComponents) > 0 {
		var premise, streetNumber, route, neighborhood, locality, sublocality, areaLevel1, country, postalCode string
		for _, comp := range res[0].AddressComponents {
			if len(comp.Types) > 0 {
				for _, tp := range comp.Types {
					if tp == "premise" {
						premise = comp.ShortName
					} else if tp == "street_number" {
						streetNumber = comp.ShortName
					} else if tp == "route" {
						route = comp.ShortName
					} else if tp == "neighborhood" {
						neighborhood = comp.ShortName
					} else if tp == "locality" {
						locality = comp.ShortName
					} else if tp == "sublocality" {
						sublocality = comp.ShortName
					} else if tp == "administrative_area_level_1" {
						areaLevel1 = comp.ShortName
					} else if tp == "country" {
						country = comp.ShortName
					} else if tp == "postal_code" {
						postalCode = comp.ShortName
					}
				}
			}
		}
		var line1Arr []string
		if premise != "" {
			line1Arr = append(line1Arr, strings.ReplaceAll(premise, ",", " "))
		}
		if streetNumber != "" {
			line1Arr = append(line1Arr, strings.ReplaceAll(streetNumber, ",", " "))
		}
		if route != "" {
			line1Arr = append(line1Arr, strings.ReplaceAll(route, ",", " "))
		}
		if neighborhood != "" {
			line1Arr = append(line1Arr, strings.ReplaceAll(neighborhood, ",", " "))
		}
		if locality == "" && sublocality != "" {
			locality = sublocality
		}
		addr = fmt.Sprintf("%s, %s, %s %s, %s", strings.Join(line1Arr, " "), locality, areaLevel1, postalCode, country)
	}
	location := res[0].Geometry.Location
	return location.Lat, location.Lng, addr, nil
}

// GoogleDistanceMatrix 谷歌距离矩阵
func GoogleDistanceMatrix(apiKey string, origins, destinations [][]interface{}) ([]map[int]int, error) {
	c, err := googleMapsClient(apiKey)
	if err != nil {
		return nil, err
	}
	var oris, dests []string
	for _, o := range origins {
		oris = append(oris, o[1].(string))
	}
	for _, d := range destinations {
		dests = append(dests, d[1].(string))
	}
	res, err := c.DistanceMatrix(context.Background(), &maps.DistanceMatrixRequest{
		DepartureTime: `now`,
		Units:         maps.UnitsMetric,
		Mode:          maps.TravelModeDriving,
		Origins:       oris,
		Destinations:  dests,
	})
	fmt.Println(res.Rows[0].Elements[0].Distance.Meters)
	return nil, nil
}
