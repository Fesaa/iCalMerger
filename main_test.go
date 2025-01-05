package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Fesaa/ical-merger/config"
	"github.com/Fesaa/ical-merger/log"
	ical "github.com/arran4/golang-ical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const icsPrefix = ".ics"

var now = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

type TestSuite struct {
	suite.Suite

	sourceCalServer *httptest.Server

	cals map[string]*ical.Calendar
}

// newSourceCalendar creates a new source calendar with the given name, rules and modifiers
// this is used to create a source mock source calendar which will then be used in the testing
// of the merge process
func (s *TestSuite) newSourceCalendar(name string, rules []config.Rule, modifiers []config.Modifier) (*ical.Calendar, config.SourceInfo, error) {
	cal, ok := s.cals[name]
	if ok {
		return cal, config.SourceInfo{}, fmt.Errorf("calendar already exists")
	}

	s.cals[name] = ical.NewCalendar()
	s.cals[name].SetName(name)

	return s.cals[name], config.SourceInfo{
		Name:      name,
		Url:       s.sourceCalServer.URL + "/" + name + ".ics",
		Rules:     rules,
		Modifiers: modifiers,
	}, nil
}

// newMockICalServer creates a new mock iCal server which will be used to serve the mock source calendars
func (s *TestSuite) newMockICalServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.EqualFold(r.URL.Path[len(r.URL.Path)-4:], icsPrefix) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		calName := r.URL.Path[1 : len(r.URL.Path)-4]
		cal, ok := s.cals[calName]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var buf bytes.Buffer
		err := cal.SerializeTo(&buf)
		require.NoError(s.T(), err)
		_, err = w.Write(buf.Bytes())
		require.NoError(s.T(), err)
	})
	return httptest.NewServer(mux)
}

// Setup the test
func (s *TestSuite) SetupTest() {
	s.cals = make(map[string]*ical.Calendar)
	s.sourceCalServer = s.newMockICalServer()
}

// Cleanup after the test
func (s *TestSuite) TearDownSuite() {
	s.sourceCalServer.Close()
}

// TestMergeCals tests the merging of calendars
func (s *TestSuite) TestMergeCals() {
	// Create a new source calendar
	cal1, cal1SI, err := s.newSourceCalendar("cal1", []config.Rule{}, []config.Modifier{})
	require.NoError(s.T(), err)
	cal2, cal2SI, err := s.newSourceCalendar("cal2", []config.Rule{}, []config.Modifier{})
	require.NoError(s.T(), err)

	// Create events
	for i := 0; i < 3; i++ {
		n := fmt.Sprintf("event%d", i)
		e := cal1.AddEvent(n)
		e.SetSummary(n)
		e.SetDescription(n)
		e.SetStartAt(now.Add(24 * time.Hour).Add(time.Duration(6+i) * time.Hour))
	}

	// Create a merged calendar server
	server := newTestCalServer("test", cal1SI, cal2SI)
	defer server.Close()
	// Fetch the calendar
	actualCal, err := server.fetchCalendar()
	require.NoError(s.T(), err)
	require.NotNil(s.T(), actualCal)
	assert.Len(s.T(), actualCal.Events(), len(cal1.Events())+len(cal2.Events()))

	for _, e := range actualCal.Events() {
		assert.NotEmpty(s.T(), e.GetProperty(ical.ComponentPropertySummary).Value)
	}
}

func TestMain(t *testing.T) {
	log.Init("ERROR", config.Notification{Service: "none"})
	suite.Run(t, new(TestSuite))
}

// /////////////////////////
// Test Helper Functions //
// /////////////////////////
type testCalServer struct {
	*httptest.Server
	Source config.Source
}

func newTestCalServer(calName string, sources ...config.SourceInfo) testCalServer {
	source := config.Source{
		Name:     calName,
		EndPoint: calName,
		Info:     sources,
	}
	mux := newServerMux(&config.Config{
		Sources: []config.Source{
			source,
		},
	})
	server := httptest.NewServer(mux)

	return testCalServer{
		Server: server,
		Source: source,
	}
}

func (s *testCalServer) fetchCalendar() (*ical.Calendar, error) {
	resp, err := s.Server.Client().Get(s.Server.URL + "/" + s.Source.EndPoint + icsPrefix)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ical.ParseCalendar(resp.Body)
}

func (s *testCalServer) Close() {
	s.Server.Close()
}
