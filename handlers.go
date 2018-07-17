package horloge

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type jsonMessage struct {
	Message string `json:"message"`
}

type versionMessage struct {
	Version string `json:"version"`
}

type jobScheduledMessage struct {
	Nexts []time.Time `json:"nexts"`
	Name  string      `json:"name"`
}

type jobRegistrationMessage struct {
	Name    string  `json:"name"`
	Pattern Pattern `json:"pattern"`
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
		var data = jobRegistrationMessage{}
		logger := c.Logger()

		if err := c.Bind(&data); err != nil {
			logger.Error(err)
			return c.JSON(http.StatusBadRequest, jsonMessage{Message: "bad request"})
		}

		if data.Name == "" || data.Pattern.IsZero() {
			return c.JSON(http.StatusBadRequest, jsonMessage{Message: "bad request"})
		}

		job := NewJob(data.Name, data.Pattern)
		nexts, err := runner.AddJob(job)
		if err != nil {
			logger.Error(err)
			return c.JSON(http.StatusConflict, jsonMessage{Message: "a job with this name already exists"})
		}

		return c.JSON(http.StatusAccepted, jobScheduledMessage{
			Name:  data.Name,
			Nexts: nexts,
		})
	}
}
