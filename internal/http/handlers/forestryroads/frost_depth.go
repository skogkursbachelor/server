package forestryroads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"skogkursbachelor/server/internal/constants"
	"skogkursbachelor/server/internal/http/handlers/forestryroads/structures"
	"strings"
)

func mapGridCentersToFrozenStatus(featureMap map[string]bool, date string) (map[string]bool, error) {
	// Create request body
	// Coordinates is in format "X1 Y1, X2 Y2, ..."
	stringBuilder := strings.Builder{}

	for key := range featureMap {
		stringBuilder.WriteString(strings.Replace(key, ",", " ", -1))
		stringBuilder.WriteString(",")
	}

	coordinatesString := stringBuilder.String()

	// Remove last comma
	length := len(coordinatesString)
	if length > 0 {
		coordinatesString = coordinatesString[:length-1]
	}

	body := structures.NVEFrostDepthRequest{
		Theme:            "gwb_frd",
		StartDate:        date + "T00",
		EndDate:          date + "T00",
		Format:           "json",
		MapCoordinateCsv: coordinatesString,
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Use NVE api to get frost data
	r, err := http.NewRequest(
		http.MethodPost,
		constants.NVEFrostDepthAPI,
		bytes.NewBuffer(bodyJSON),
	)
	if err != nil {
		log.Error().Msg("Failed to create request: " + err.Error())
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")

	// Do the request
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Error().Msg("Failed to do request: " + err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	// Decode response
	response := structures.NVECellTimeSeriesFrostDepthResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Error().Msg("Failed to decode response: " + err.Error())
		return nil, err
	}

	threshold := 0.0

	// Create map of isFrozen values
	isFrozenMap := make(map[string]bool)
	for i := range response.CellTimeSeries {
		key := fmt.Sprintf("%d,%d", response.CellTimeSeries[i].X, response.CellTimeSeries[i].Y)
		isFrozenMap[key] = response.CellTimeSeries[i].Data[0] > threshold
	}

	return isFrozenMap, nil
}
