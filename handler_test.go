package horloge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func request(method, path, data string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func TestHTTPHandlerPing(t *testing.T) {
	rec := httptest.NewRecorder()

	context, _ := gin.CreateTestContext(rec)
	context.Request = request("GET", "/ping", "")

	fn := HTTPHandlerPing()
	fn(context)

	expectedCode := http.StatusOK
	actualCode := rec.Code

	if expectedCode != actualCode {
		t.Errorf("expected status code to be %d, got %d", expectedCode, actualCode)
	}

	expectedBody := bytes.NewBufferString("{\"message\":\"pong\"}")
	actualBody := strings.TrimSpace(rec.Body.String())

	if expectedBody.String() != actualBody {
		t.Errorf("expected body to be `%s`, got `%s`", expectedBody, actualBody)
	}
}

func TestHTTPHandlerVersion(t *testing.T) {
	rec := httptest.NewRecorder()

	context, _ := gin.CreateTestContext(rec)
	context.Request = request("GET", "/version", "")

	fn := HTTPHandlerVersion()
	fn(context)

	expectedCode := http.StatusOK
	actualCode := rec.Code

	if expectedCode != actualCode {
		t.Errorf("expected status code to be %d, got %d", expectedCode, actualCode)
	}

	message, _ := json.Marshal(map[string]string{"version": Version})
	expectedBody := bytes.NewBuffer(message)
	actualBody := strings.TrimSpace(rec.Body.String())

	if expectedBody.String() != actualBody {
		t.Errorf("expected body to be %s, got %s", expectedBody, actualBody)
	}
}

func TestHTTPHandlerHealthCheck(t *testing.T) {
	rec := httptest.NewRecorder()

	context, _ := gin.CreateTestContext(rec)
	context.Request = request("GET", "/health_check", "")

	fn := HTTPHandlerHealthCheck()
	fn(context)

	expectedCode := http.StatusOK
	actualCode := rec.Code

	if expectedCode != actualCode {
		t.Errorf("expected status code to be %d, got %d", expectedCode, actualCode)
	}

	expectedBody := bytes.NewBufferString("{\"message\":\"ok\"}")
	actualBody := strings.TrimSpace(rec.Body.String())

	if expectedBody.String() != actualBody {
		t.Errorf("expected body to be %s, got %s", expectedBody, actualBody)
	}
}

func TestHTTPHandlerRegisterJobEmptyBody(t *testing.T) {
	runner := NewRunner()
	rec := httptest.NewRecorder()

	context, _ := gin.CreateTestContext(rec)
	context.Request = request("POST", "/jobs", "")

	fn := HTTPHandlerRegisterJob(runner)
	fn(context)

	expectedCode := http.StatusBadRequest
	actualCode := rec.Code

	if expectedCode != actualCode {
		t.Errorf("expected status code to be %d, got %d", expectedCode, actualCode)
	}

	var message = &JSONMessage{}
	err := json.Unmarshal(rec.Body.Bytes(), message)

	if err != nil {
		t.Error(err)
	}

	if message.Details != MalformedMessage {
		t.Errorf("expected body to be %s, got %s", message.Details, MalformedMessage)
	}
}

func TestHTTPHandlerRegisterJobMissingParams(t *testing.T) {
	runner := NewRunner()
	rec := httptest.NewRecorder()

	context, _ := gin.CreateTestContext(rec)
	context.Request = request("POST", "/jobs", "{\"foo\": \"bar\"}")

	fn := HTTPHandlerRegisterJob(runner)
	fn(context)

	expectedCode := http.StatusBadRequest
	actualCode := rec.Code

	if expectedCode != actualCode {
		t.Errorf("expected status code to be %d, got %d", expectedCode, actualCode)
	}

	var message = &JSONMessage{}
	err := json.Unmarshal(rec.Body.Bytes(), message)

	if err != nil {
		t.Error(err)
	}

	if message.Details != InvalidJobRequestBody {
		t.Errorf("expected body to be \"%s\", got \"%s\"", InvalidJobRequestBody, message.Details)
	}
}

func TestHTTPHandlerRegisterJob(t *testing.T) {
	runner := NewRunner()
	jobName := "Test"
	body := fmt.Sprintf("{\"name\":\"%s\",\"pattern\":{\"occurence\":\"every\",\"minute\":2}}", jobName)
	rec := httptest.NewRecorder()

	context, _ := gin.CreateTestContext(rec)
	context.Request = request("POST", "/jobs", body)

	fn := HTTPHandlerRegisterJob(runner)
	fn(context)

	expectedCode := http.StatusAccepted
	actualCode := rec.Code

	if expectedCode != actualCode {
		t.Errorf("expected status code to be %d, got %d", expectedCode, actualCode)
	}

	var message = &JobScheduledMessage{}
	err := json.Unmarshal(rec.Body.Bytes(), message)

	if err != nil {
		t.Error(err)
	}

	if message.Name != jobName {
		t.Errorf("expected job name to be \"%s\", got \"%s\"", jobName, message.Name)
	}
}

func TestHTTPHandlerRegisterJobConflict(t *testing.T) {
	runner := NewRunner()
	jobName := "Test"
	body := fmt.Sprintf("{\"name\":\"%s\",\"pattern\":{\"occurence\":\"every\",\"minute\":2}}", jobName)
	rec := httptest.NewRecorder()
	rec2 := httptest.NewRecorder()

	context, _ := gin.CreateTestContext(rec)
	context2, _ := gin.CreateTestContext(rec2)

	context.Request = request("POST", "/jobs", body)
	context2.Request = request("POST", "/jobs", body)

	fn := HTTPHandlerRegisterJob(runner)
	fn(context)
	fn(context2)

	expectedCode := http.StatusConflict
	actualCode := rec2.Code

	if expectedCode != actualCode {
		t.Errorf("expected status code to be %d, got %d", expectedCode, actualCode)
	}

	var message = &JSONMessage{}
	fmt.Println(rec2.Body.String())
	err := json.Unmarshal(rec2.Body.Bytes(), message)

	if err != nil {
		t.Error(err)
	}

	expectedDetails := fmt.Sprintf(JobExistsError, jobName)
	actualDetails := message.Details

	if actualDetails != expectedDetails {
		t.Errorf("expected details to be \"%s\", got \"%s\"", expectedDetails, actualDetails)
	}
}
