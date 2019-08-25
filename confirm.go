package tinkoff

import (
	"errors"
	"fmt"
	"strconv"
)

type ConfirmRequest struct {
	BaseRequest

	PaymentID uint64           `json:"PaymentId"`
	ClientIP  string           `json:"IP"`
	Amount    uint64           `json:"Amount"`
	Receipt   *Receipt `json:"Receipt"`
}

func (i *ConfirmRequest) GetValuesForToken() map[string]string {
	return map[string]string{
		"Amount":    strconv.FormatUint(i.Amount, 10),
		"IP":        i.ClientIP,
		"PaymentId": strconv.FormatUint(i.PaymentID, 10),
	}
}

type ConfirmResponse struct {
	TerminalKey  string `json:"TerminalKey"`       // Идентификатор терминала, выдается Продавцу Банком
	OrderID      string `json:"OrderId"`           // Номер заказа в системе Продавца
	Success      bool   `json:"Success"`           // Успешность операции
	Status       string `json:"Status"`            // Статус транзакции
	PaymentID    string `json:"PaymentId"`         // Уникальный идентификатор транзакции в системе Банка
	ErrorCode    string `json:"ErrorCode"`         // Код ошибки, «0» - если успешно
	ErrorMessage string `json:"Message,omitempty"` // Краткое описание ошибки
	ErrorDetails string `json:"Details,omitempty"` // Подробное описание ошибки
}

func (c *Client) Confirm(request *ConfirmRequest) (status string, err error) {
	response, err := c.postRequest("/Confirm", request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	var res ConfirmResponse
	err = c.decodeResponse(response, &res)
	if err != nil {
		return
	}

	if !res.Success || res.ErrorCode != "0" {
		err = errors.New(fmt.Sprintf("while Confirm request: code %s - %s. %s", res.ErrorCode, res.ErrorMessage, res.ErrorDetails))
	}

	return status, err
}
