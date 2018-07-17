package horloge

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

const (
	MalformedMessage      string = "Malformed or empty request body"
	InvalidJobRequestBody string = "\"pattern\" and \"name\" must be present"
)

type JSONMessage struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

type JobScheduledMessage struct {
	Nexts []time.Time `json:"nexts"`
	Name  string      `json:"name"`
}

type jobRegistrationMessage struct {
	Name    string  `json:"name"`
	Pattern Pattern `json:"pattern"`
}

// HTTPHandlerPing Handles GET requests to /ping.
//
// Replies with pong.
func HTTPHandlerPing() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	}
}

// HTTPHandlerVersion Handles GET requests to /version.
//
// Replies with project version
func HTTPHandlerVersion() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"version": Version})
	}
}

// HTTPHandlerHealthCheck Handles GET requests to /health_check.
//
// Replies an empty string and status code 200. This is useful if you want
// to monitor the state of the application.
func HTTPHandlerHealthCheck() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
	}
}

// HTTPHandlerRegisterJob Handles POST requests to /job.
//
// To add a job you must set a request to /job with a json body
//
func HTTPHandlerRegisterJob(runner *Runner) func(c echo.Context) error {
	return func(c echo.Context) error {
		var data = jobRegistrationMessage{}
		logger := c.Logger()

		if err := c.Bind(&data); err != nil {
			logger.Error(err)
			return c.JSON(http.StatusBadRequest, JSONMessage{Message: "bad request", Details: MalformedMessage})
		}

		if data.Name == "" || data.Pattern.IsZero() {
			return c.JSON(http.StatusBadRequest, JSONMessage{Message: "bad request", Details: InvalidJobRequestBody})
		}

		job := NewJob(data.Name, data.Pattern)
		nexts, err := runner.AddJob(job)
		if err != nil {
			logger.Error(err)
			return c.JSON(http.StatusConflict, JSONMessage{Message: "conflict", Details: err.Error()})
		}

		return c.JSON(http.StatusAccepted, JobScheduledMessage{
			Name:  data.Name,
			Nexts: nexts,
		})
	}
}
