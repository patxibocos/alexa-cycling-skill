package alexa

import (
	"fmt"
	"io"
	"strings"
	"time"
)

func getLocation(request Request) *time.Location {
	resp, _ := doRequest("GET", fmt.Sprintf("/v2/devices/%s/settings/System.timeZone", request.Context.System.Device.DeviceID), request, nil)
	responseBytes, _ := io.ReadAll(resp.Body)
	timeZone := string(responseBytes)
	// Remove heading and trailing " from the response
	location, _ := time.LoadLocation(strings.Trim(timeZone, "\""))
	fmt.Println(location)
	return location
}
