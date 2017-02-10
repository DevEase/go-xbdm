package goxbdm

import (
	"strings"

	"fmt"

	"github.com/0xdeafcafe/go-xbdm/clients"
)

// Client ..
type Client struct {
	tcpClient *clients.TCPClient
}

const (
	defaultXboxTCPPort = 730
)

// DebugName returns the debug name of the Xbox
func (client *Client) DebugName() (string, error) {
	client.tcpClient.WriteString("dbgname", true)
	return client.tcpClient.ReadString()
}

// NewXBDMClient ..
func NewXBDMClient(xboxIP string) (*Client, error) {
	return NewXBDMClientWithPort(xboxIP, defaultXboxTCPPort)
}

// NewXBDMClientWithPort ..
func NewXBDMClientWithPort(xboxIP string, port int) (*Client, error) {
	tcpClient, err := clients.NewTCPClientWithPort(xboxIP, port)
	if err != nil {
		return nil, err
	}

	client := &Client{
		tcpClient: tcpClient,
	}

	return client, nil
}

func parseMultilineResponse(str string) (map[string]string, error) {
	body := make(map[string]string)

	// Remove header and eom
	splitStr := strings.Split(str, "\r\n")
	dataLines := splitStr[1 : len(splitStr)-2]

	// Join array with spaces, then add an extra space at the end to save last entry
	dataStr := strings.Join(dataLines, " ") + " "
	currentKey := ""
	currentValue := ""
	isInsideQuote := false
	onKey := true
	for _, char := range dataStr {
		// If we detect a space, and we aren't inside a double quote, switch to `onKey`
		if char == ' ' && !isInsideQuote {
			// Before we switch we need to save the read values to the `body` map and
			// reset the `currentKey` and `currentValue` variables
			body[currentKey] = currentValue
			currentKey = ""
			currentValue = ""
			onKey = true
			continue
		}

		// If we find an equals, and we aren't inside a double quote, switch to `!onKey`
		if char == '=' && !isInsideQuote {
			onKey = false
			continue
		}

		// If we find a double quote, switch between inside and outside
		if char == '"' {
			isInsideQuote = !isInsideQuote
		}

		// Save the value to the relevant part
		if onKey {
			currentKey = fmt.Sprintf("%s%s", currentKey, string(char))
		} else {
			currentValue = fmt.Sprintf("%s%s", currentValue, string(char))
		}
	}

	return body, nil
}
