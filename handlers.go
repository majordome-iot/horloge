package horloge

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	// MalformedMessage Returned when an unacceptable request is made
	MalformedMessage string = "Malformed or empty request body"
	// InvalidJobRequestBody Returned when tring to register a job without a pattern or a name
	InvalidJobRequestBody string = "\"pattern\" and \"name\" must be present"
	// UnableToSerializeJobs Returned when the server could not serialize jobs
	UnableToSerializeJobs string = "Unable to serialize jobs"
)

// JSONErrorMessage Used to detail an error
type JSONErrorMessage struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

// JobScheduledMessage Returned when a job was successfully registered
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
func HTTPHandlerPing() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	}
}

// HTTPHandlerVersion Handles GET requests to /version.
//
// Replies with project version
func HTTPHandlerVersion() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"version": Version})
	}
}

// HTTPHandlerHealthCheck Handles GET requests to /health_check.
//
// Replies an empty string and status code 200. This is useful if you want
// to monitor the state of the application.
func HTTPHandlerHealthCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "ok"})
	}
}

// HTTPHandlerRegisterJob Handles POST requests to /jobs.
//
// To add a job you must send a request to /jobs with a json body
func HTTPHandlerRegisterJob(r *Runner) func(c *gin.Context) {
	return func(c *gin.Context) {
		var data = jobRegistrationMessage{}

		if err := c.Bind(&data); err != nil {
			r.log.WithFields(logrus.Fields{
				"code": 400,
			}).Error(string(err.Error()))
			c.JSON(http.StatusBadRequest, JSONErrorMessage{
				Message: "bad request",
				Details: MalformedMessage,
			})
			return
		}

		if data.Name == "" || data.Pattern.IsZero() {
			r.log.WithFields(logrus.Fields{
				"code": 400,
			}).Error(InvalidJobRequestBody)

			c.JSON(http.StatusBadRequest, JSONErrorMessage{
				Message: "bad request",
				Details: InvalidJobRequestBody,
			})
			return
		}

		job := NewJob(data.Name, data.Pattern, data.Args)
		job.Description = data.Description
		nexts, err := r.AddJob(*job)

		if err != nil {
			r.log.WithFields(logrus.Fields{
				"code": 409,
			}).Error(err)
			c.JSON(http.StatusConflict, JSONErrorMessage{
				Message: "conflict",
				Details: err.Error(),
			})
			return
		}

		c.JSON(http.StatusAccepted, JobScheduledMessage{
			Name:  data.Name,
			Nexts: nexts,
		})
	}
}

// HTTPHandlerJobDetail Show job detail
func HTTPHandlerJobDetail(r *Runner) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("id")
		job := r.GetJob(name)

		if job == nil {
			c.JSON(http.StatusNotFound, JSONErrorMessage{
				Message: "not found",
				Details: "Job with name `" + name + "` does not exist",
			})
		}

		c.JSON(http.StatusOK, job)
	}
}

// HTTPHandlerDeleteJob Delete a job
func HTTPHandlerDeleteJob(r *Runner) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("id")
		job := r.GetJob(name)

		if job == nil {
			c.JSON(http.StatusNotFound, JSONErrorMessage{
				Message: "not found",
				Details: "Job with name " + name + "does not exist",
			})
		}

		r.RemoveJob(job)
		c.JSON(http.StatusNoContent, nil)
	}
}

// HTTPHandlerListJobs List Jobs
func HTTPHandlerListJobs(r *Runner) func(c *gin.Context) {
	return func(c *gin.Context) {
		jobs, err := r.ToJSON()
		if err != nil {
			c.JSON(http.StatusInternalServerError, JSONErrorMessage{
				Message: "internal server error",
				Details: UnableToSerializeJobs,
			})
		}
		c.JSON(http.StatusOK, jobs)
	}
}
