package forestryroads

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"skogkursbachelor/server/internal/constants"
	"strings"
)

// implementedMethods is a list of the implemented HTTP methods for the status endpoint.
var implementedMethods = []string{http.MethodGet}

// Handler
// Currently only supports GET requests.
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleForestryRoadGet(w, r)

	default:
		// If the method is not implemented, return an error with the allowed methods
		http.Error(
			w, fmt.Sprintf(
				"REST Method '%s' not supported. Currently only '%v' are supported.", r.Method,
				implementedMethods,
			), http.StatusNotImplemented,
		)
		return
	}
}

// handleForestryRoadGet handles GET requests to the forestry road endpoint.
func handleForestryRoadGet(w http.ResponseWriter, r *http.Request) {
	// Pseudo code
	// 1. Mirror the request to the remote server
	// 2. Get the response from the remote server
	// 3. Parse the response
	// 4. Calculate trafficality
	// 5. Return the response, with the calculated trafficality as a rgb value in the geojson response

	// Get time parameter from url
	time := r.URL.Query().Get("time")
	if time == "" {
		http.Error(w, "Missing time URL parameter", http.StatusBadRequest)
		return
	}

	// Split ISO string to get date. ex: 2021-03-01T00:00:00Z -> 2021-03-01
	// Gets put into struct later
	date := strings.Split(time, "T")[0]

	// Mirror request to https://wms.geonorge.no/skwms1/wms.traktorveg_skogsbilveger
	proxyReq, err := http.NewRequest(
		r.Method,
		constants.ForestryRoadsWFS+"?"+r.URL.RawQuery,
		r.Body,
	)
	if err != nil {
		http.Error(w, "Failed to create internal request", http.StatusInternalServerError)
		log.Println("Error creating request to GeoNorge for forestry roads: ", err)
		return
	}

	// Do request
	proxyResp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		http.Error(w, "Failed to fetch data from external WMS server", http.StatusBadGateway)
		log.Println("Error fetching data from GeoNorge WMS server: ", err)
		return
	}

	// Decode into struct
	var wfsResponse WFSResponse
	err = json.NewDecoder(proxyResp.Body).Decode(&wfsResponse)
	if err != nil {
		http.Error(w, "Failed to decode external response", http.StatusInternalServerError)
		log.Println("Error decoding response from GeoNorge WMS server: ", err)
		return
	}

	// Randomize color for testing, update date
	for i := range wfsResponse.Features {
		if wfsResponse.Date == "" {
			wfsResponse.Date = date
		}

		if wfsResponse.Features[i].Properties.Farge == nil {
			wfsResponse.Features[i].Properties.Farge = make([]int, 3)
		}

		// Get middle of the road (ish)
		length := len(wfsResponse.Features[i].Geometry.Coordinates)
		middleIndex := length / 2

		isFrozen, err := GetIsGroundFrozen(wfsResponse.Features[i].Geometry.Coordinates[middleIndex], date)
		if err != nil {
			http.Error(w, "Failed to get frost data", http.StatusInternalServerError)
			log.Println("Error getting frost data: ", err)
			return
		}

		// If the ground is frozen, set the color to green
		if isFrozen {
			wfsResponse.Features[i].Properties.Farge[0] = 0
			wfsResponse.Features[i].Properties.Farge[1] = 255
			wfsResponse.Features[i].Properties.Farge[2] = 0
		} else {
			// If the ground is not frozen, set the color to red
			wfsResponse.Features[i].Properties.Farge[0] = 255
			wfsResponse.Features[i].Properties.Farge[1] = 0
			wfsResponse.Features[i].Properties.Farge[2] = 0
		}
	}

	// Encode response
	err = json.NewEncoder(w).Encode(wfsResponse)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Println("Error encoding final response: ", err)
		return
	}
}
