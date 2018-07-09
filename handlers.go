package horloge

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type jsonMessage struct {
	Message string `json:"message"`
}

type versionMessage struct {
	Version string `json:"version"`
}

func HttpHandlerPing() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, jsonMessage{Message: "pong"})
	}
}

func HttpHandlerVersion() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, versionMessage{Version: Version})
	}
}

func HttpHandlerRegisterJob(runner Runner) func(c echo.Context) error {
	return func(c echo.Context) error {
		var pattern = Pattern{}
		logger := c.Logger()

		if err := c.Bind(&pattern); err != nil {
			logger.Error(err)
			return c.JSON(http.StatusBadRequest, jsonMessage{Message: "bad request"})
		}

		job := NewJob("foobar", pattern)
		nexts, err := runner.AddJob(job)
		if err != nil {
			logger.Error(err)
			return c.JSON(http.StatusConflict, jsonMessage{Message: "a job with this name already exists"})
		}

		b, err := json.Marshal(nexts)
		if err != nil {
			logger.Error(err)
			return c.JSON(http.StatusBadRequest, jsonMessage{Message: "bad request"})
		}

		fmt.Println(b, nexts)

		return c.JSON(http.StatusAccepted, b)
	}
}
