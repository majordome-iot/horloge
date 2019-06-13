package horloge

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

const (
	MalformedMessage      string = "Malformed or empty request body"
	InvalidJobRequestBody string = "\"pattern\" and \"name\" must be present"
	UnableToSerializeJobs string = "Unable to serialize jobs"
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
	Name        string   `json:"name"`
	Pattern     Pattern  `json:"pattern"`
	Args        []string `json:"args"`
	Description string   `json:"description"`
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

// HTTPHandlerRegisterJob Handles POST requests to /jobs.
//
// To add a job you must send a request to /jobs with a json body
//
func HTTPHandlerRegisterJob(r *Runner) func(c echo.Context) error {
	return func(c echo.Context) error {
		var data = jobRegistrationMessage{}

		if err := c.Bind(&data); err != nil {
			r.log.WithFields(logrus.Fields{
				"code": 400,
			}).Error(string(err.Error()))
			return c.JSON(http.StatusBadRequest, JSONMessage{
				Message: "bad request",
				Details: MalformedMessage,
			})
		}

		if data.Name == "" || data.Pattern.IsZero() {
			r.log.WithFields(logrus.Fields{
				"code": 400,
			}).Error(InvalidJobRequestBody)
			return c.JSON(http.StatusBadRequest, JSONMessage{
				Message: "bad request",
				Details: InvalidJobRequestBody,
			})
		}

		job := NewJob(data.Name, data.Pattern, data.Args)
		job.Description = data.Description
		nexts, err := r.AddJob(*job)

		if err != nil {
			r.log.WithFields(logrus.Fields{
				"code": 409,
			}).Error(err)
			return c.JSON(http.StatusConflict, JSONMessage{
				Message: "conflict",
				Details: err.Error(),
			})
		}

		return c.JSON(http.StatusAccepted, JobScheduledMessage{
			Name:  data.Name,
			Nexts: nexts,
		})
	}
}

// HTTPHandlerJobDetail Show job detail
func HTTPHandlerJobDetail(r *Runner) func(c echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("id")
		job := r.GetJob(name)

		if job == nil {
			return c.JSON(http.StatusNotFound, JSONMessage{
				Message: "not found",
				Details: "Job with name `" + name + "` does not exist",
			})
		}

		return c.JSON(http.StatusOK, job)
	}
}

// HTTPHandlerDeleteJob Delete a job
func HTTPHandlerDeleteJob(r *Runner) func(c echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("id")
		job := r.GetJob(name)

		if job == nil {
			return c.JSON(http.StatusNotFound, JSONMessage{
				Message: "not found",
				Details: "Job with name " + name + "does not exist",
			})
		}

		r.RemoveJob(job)

		return c.NoContent(http.StatusNoContent)
	}
}

// HTTPHandlerListJobs List Jobs
func HTTPHandlerListJobs(r *Runner) func(c echo.Context) error {
	return func(c echo.Context) error {
		jobs, err := r.ToJSON()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, JSONMessage{
				Message: "internal server error",
				Details: UnableToSerializeJobs,
			})
		}
		return c.JSON(http.StatusOK, jobs)
	}
}
