package client

func (client *Client) DeliveryUpdate(deliveryId string, status string) error {
	body := map[string]interface{}{
		"status": status,
	}

	_, err := client.Put(client.corePath("/agent/deliveries/%s", deliveryId), body)

	return err
}

func (client *Client) DeliveryLogsUpdate(deliveryId string, logs string) error {
	body := map[string]interface{}{
		"logs": logs,
	}

	_, err := client.Put(client.corePath("/agent/deliveries/%s/logs", deliveryId), body)

	return err
}
