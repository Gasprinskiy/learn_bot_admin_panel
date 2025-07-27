package http_req

import (
	"encoding/json"
	"io"
	"learn_bot_admin_panel/internal/entity/global"
	"net/http"
)

func Get[T any](url string) (T, error) {
	var data T

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return data, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return data, err
	}

	if res.StatusCode > 399 {
		if res.StatusCode == http.StatusNotFound {
			return data, global.ErrNoData
		}

		return data, global.ErrInternalError
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(resBody, &data); err != nil {
		return data, err
	}

	return data, nil
}
