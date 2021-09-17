package client

import "fmt"

func (client *Client) CoreDeliveryUpdate(deliveryId string, status string) error {
	body := map[string]interface{}{
		"status": status,
	}

	_, err := client.Put(client.corePath("/agent/deliveries/%s", deliveryId), body)

	return err
}

func (client *Client) CoreDeliveryLogsUpdate(deliveryId string, logs string) error {
	body := map[string]interface{}{
		"logs": logs,
	}

	_, err := client.Put(client.corePath("/agent/deliveries/%s/logs", deliveryId), body)

	return err
}

func (client *Client) corePath(path string, tokens ...interface{}) string {
	url, _ := client.CoreBaseURL.Parse(fmt.Sprintf(path, tokens...))
	return url.String()
}
