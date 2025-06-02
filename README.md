# go-librenms

A Go client library for interacting with the LibreNMS API.

## Installation

To use this library, you can import it into your Go project:

```go
import "github.com/jokelyo/go-librenms"
```

Then, run `go get` to download and install the package:

```bash
go get github.com/jokelyo/go-librenms
```

## Example Usage

Here's a basic example of how to use the library to create a new LibreNMS client and get a device:

```go
package main

import (
	"fmt"
	"log"

	"github.com/jokelyo/go-librenms"
)

func main() {
	// Replace with your LibreNMS API URL and token
	baseURL := "https://your-librenms-instance.com/"
	token := "YOUR_API_TOKEN"

	// Create a new LibreNMS client
	client, err := librenms.New(baseURL, token)
	if err != nil {
		log.Fatalf("Error creating LibreNMS client: %v", err)
	}

	// Get a device by its hostname or ID
	// Replace "device-hostname-or-id" with the actual hostname or ID
	deviceIdentifier := "device-hostname-or-id"
	deviceResp, err := client.GetDevice(deviceIdentifier)
	if err != nil {
		log.Fatalf("Error getting device: %v", err)
	}

	// Print device information
	if len(deviceResp.Devices) > 0 {
		fmt.Printf("Device ID: %d\n", deviceResp.Devices[0].DeviceID)
		fmt.Printf("Hostname: %s\n", deviceResp.Devices[0].Hostname)
		fmt.Printf("OS: %s\n", deviceResp.Devices[0].OS)
	} else {
		fmt.Println("No device found.")
	}
}
```

### Creating a Device Example

Here's an example of how to create a new device in LibreNMS:

```go
// Initialize the device creation request
deviceCreateReq := &librenms.DeviceCreateRequest{
    Hostname:      "192.168.1.10",
    Display:       "My New Router",
    SNMPCommunity: "public",
    SNMPVersion:   "v2c",
}

// Create the device
deviceResp, err := client.CreateDevice(deviceCreateReq)
if err != nil {
    log.Fatalf("Error creating device: %v", err)
}

// Handle the response
if deviceResp.Status == "ok" {
    fmt.Printf("Device created successfully! Device ID: %d\n", deviceResp.Devices[0].DeviceID)
    fmt.Printf("Hostname: %s\n", deviceResp.Devices[0].Hostname)
} else {
    fmt.Printf("Failed to create device: %s\n", deviceResp.Message)
}
```

