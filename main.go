package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	key := ""
	// c, err := maps.NewClient(maps.WithAPIKey(key))
	// if err != nil {
	// 	log.Fatalf("fatal error: %s", err)
	// }
	// r := &maps.DirectionsRequest{
	// 	Origin:      "Sydney",
	// 	Destination: "Perth",
	// }

	// near := &maps.NearbySearchRequest{
	// 	Location: &maps.LatLng{
	// 		Lat: -7.9192886,
	// 		Lng: 112.6170444,
	// 	},
	// 	Radius:  100000,
	// 	Keyword: "Kec.",
	// }
	// r := &maps.GeocodingRequest{
	// 	LatLng: &maps.LatLng{
	// 		Lat: -7.9192886,
	// 		Lng: 112.6170444,
	// 	},
	// }
	// res, err := c.NearbySearch(context.Background(), near)

	// places := &maps.TextSearchRequest{
	// 	Query:  "kecamatan",
	// 	Radius: 10000,
	// 	Location: &maps.LatLng{
	// 		Lat: -7.9192886,
	// 		Lng: 112.6170444,
	// 	},
	// 	Language: "id",
	// }
	// res2, err := c.TextSearch(context.Background(), places)

	// if err != nil {
	// 	log.Fatalf("fatal error: %s", err)
	// }

	// pretty.Println(res2)

	var finalResult PlacesSearchResponse
	query := "kecamatan"
	lat := "-7.9192886"
	lng := "112.6170444"
	langu := "id"
	radius := "10000"
	pageToken := ""

	for {
		err, res := searchPlaces(key, query, lat, lng, langu, radius, pageToken)
		if err != nil {
			log.Fatalf("fatal error: %s", err)
			return
		}
		var responseMap PlacesSearchResponse
		json.Unmarshal(res, &responseMap)
		finalResult.Results = append(finalResult.Results, responseMap.Results...)

		log.Println("next page token : ", responseMap.NextPageToken)
		if responseMap.NextPageToken == "" {
			break
		} else {
			pageToken = responseMap.NextPageToken
		}
	}

	fmt.Println("Length Result : ", len(finalResult.Results))

	b, err := json.Marshal(finalResult)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

}

func searchPlaces(key, query, lat, lng, language, radius, pageToken string) (error, []byte) {

	url := "https://maps.googleapis.com/maps/api/place/textsearch/json?"
	if pageToken != "" {
		url = url + "pagetoken=" + pageToken
		url = url + "&key=" + key
	} else {
		url = url + "query=" + query
		url = url + fmt.Sprintf("&radius=%s&location=%s,%s&language=%s&key=%s", radius, lat, lng, language, key)
	}

	log.Println("url : ", url)

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return err, nil
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err, nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err, nil
	}

	return nil, body
}
