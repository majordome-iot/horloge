package horloge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
)

func request(method, path, data string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return req
}

func TestHTTPHandlerPing(t *testing.T) {
	req := request(echo.GET, "/ping", "")
	rec := httptest.NewRecorder()

	e := echo.New()
	context := e.NewContext(req, rec)

	fn := HTTPHandlerPing()
	fn(context)

	expectedCode := 200
	actualCode := rec.Code

	if expectedCode != actualCode {
		t.Errorf("expected status code to be %d, got %d", expectedCode, actualCode)
	}

	expectedBody := bytes.NewBufferString("{\"message\":\"pong\"}")
	actualBody := rec.Body

	if expectedBody.String() != actualBody.String() {
		t.Errorf("expected body to be %s, got %s", expectedBody, actualBody)
	}
}

func TestHTTPHandlerVersion(t *testing.T) {
	req := request(echo.GET, "/version", "")
	rec := httptest.NewRecorder()

	e := echo.New()
	context := e.NewContext(req, rec)

	fn := HTTPHandlerVersion()
	fn(context)

	expectedCode := 200
	actualCode := rec.Code

	if expectedCode != actualCode {
		t.Errorf("expected status code to be %d, got %d", expectedCode, actualCode)
	}

	message, _ := json.Marshal(map[string]string{"version": Version})
	expectedBody := bytes.NewBuffer(message)
	actualBody := rec.Body

	if expectedBody.String() != actualBody.String() {
		t.Errorf("expected body to be %s, got %s", expectedBody, actualBody)
	}
}

func TestHTTPHandlerHealthCheck(t *testing.T) {
	req := request(echo.GET, "/health_check", "")
	rec := httptest.NewRecorder()

	e := echo.New()
	context := e.NewContext(req, rec)

	fn := HTTPHandlerHealthCheck()
	fn(context)

	expectedCode := 200
	actualCode := rec.Code

	if expectedCode != actualCode {
		t.Errorf("expected status code to be %d, got %d", expectedCode, actualCode)
	}

	expectedBody := bytes.NewBufferString("{\"message\":\"ok\"}")
	actualBody := rec.Body

	if expectedBody.String() != actualBody.String() {
		t.Errorf("expected body to be %s, got %s", expectedBody, actualBody)
	}
}

func TestHTTPHandlerRegisterJobEmptyBody(t *testing.T) {
	runner := NewRunner()
	req := request(echo.POST, "/job", "")
	rec := httptest.NewRecorder()

	e := echo.New()
	context := e.NewContext(req, rec)

	fn := HTTPHandlerRegisterJob(runner)
	fn(context)

	expectedCode := 400
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

func TestHTTPHandlerRegisterJobMissingPArams(t *testing.T) {
	runner := NewRunner()
	req := request(echo.POST, "/job", "{\"foo\": \"bar\"}")
	rec := httptest.NewRecorder()

	e := echo.New()
	context := e.NewContext(req, rec)

	fn := HTTPHandlerRegisterJob(runner)
	fn(context)

	expectedCode := 400
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
	req := request(echo.POST, "/job", body)
	rec := httptest.NewRecorder()

	e := echo.New()
	context := e.NewContext(req, rec)

	fn := HTTPHandlerRegisterJob(runner)
	fn(context)

	expectedCode := 202
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
	req := request(echo.POST, "/job", body)
	req2 := request(echo.POST, "/job", body)
	rec := httptest.NewRecorder()
	rec2 := httptest.NewRecorder()

	e := echo.New()
	context := e.NewContext(req, rec)
	context2 := e.NewContext(req2, rec2)

	fn := HTTPHandlerRegisterJob(runner)
	fn(context)
	fn(context2)

	expectedCode := 409
	actualCode := rec2.Code

	if expectedCode != actualCode {
		t.Errorf("expected status code to be %d, got %d", expectedCode, actualCode)
	}

	var message = &JSONMessage{}
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
