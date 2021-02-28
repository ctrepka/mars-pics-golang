package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/api/pictures/", func(w http.ResponseWriter, r *http.Request) {
		limit := r.URL.Query().Get("limit")
		if len(limit) < 1 {
			limit = "10"
		}
		data, err := query(limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", nil)
}

func query(limit string) (marsPictures, error) {
	resp, err := http.Get("https://api.nasa.gov/mars-photos/api/v1/rovers/curiosity/photos?sol=" + limit + "&api_key=DEMO_KEY")
	if err != nil {
		return marsPictures{}, err
	}
	defer resp.Body.Close()
	var d marsPictures
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return marsPictures{}, err
	}
	return d, nil
}

type marsPictures struct {
	Photos []struct {
		ID     int `json:"id"`
		Sol    int `json:"sol"`
		Camera struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			RoverID  int    `json:"rover_id"`
			FullName string `json:"full_name"`
		} `json:"camera"`
		ImgSrc    string `json:"img_src"`
		EarthDate string `json:"earth_date"`
		Rover     struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			LandingDate string `json:"landing_date"`
			LaunchDate  string `json:"launch_date"`
			Status      string `json:"status"`
		} `json:"rover"`
	} `json:"photos"`
}
