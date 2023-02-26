package alexa

import (
	"fmt"
	"io"
	"strings"
	"time"
)

func getLocation(request Request) *time.Location {
	resp, err := doRequest("GET", fmt.Sprintf("/v2/devices/%s/settings/System.timeZone", request.Context.System.Device.DeviceID), request, nil)
	if err != nil {
		return nil
	}
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	timeZone := string(responseBytes)
	// Remove heading and trailing " from the response
	location, _ := time.LoadLocation(strings.Trim(timeZone, "\""))
	return location
}
