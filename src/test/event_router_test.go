package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	repository "github.com/arganjava/online-loket/src/mocks"
	"github.com/arganjava/online-loket/src/routers"
	. "github.com/smartystreets/goconvey/convey"
	mockTest "github.com/stretchr/testify/mock"
	"io/ioutil"
	"regexp"
	"time"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEventRouter_CreateEvent(t *testing.T) {
	Convey("Given Create Location", t, func() {
		Convey("When User Create Event", func() {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()
			ts := httptest.NewServer(routers.SetupServer(mockDB))
			defer ts.Close()
			requestBody, err := json.Marshal(map[string]string{
				"locationId":    "1ad5ae0e-9e49-4025-90aa-295e1a4bd886",
				"eventName":     "Badut Kota",
				"description":   "Desc",
				"scheduleBegin": "2020-12-30",
				"scheduleEnd":   "2021-01-30",
			})

			requestBody404, err := json.Marshal(map[string]string{
				"test": "test",
			})

			requestBody404Validate, err := json.Marshal(map[string]string{
				"locationId":    "1ad5ae0e-9e49-4025-90aa-295e1a4bd886",
				"eventName":     "Badut KotaBadut KotaBadut KotaBadut KotaBadut KotaBadut KotaBadut KotaBadut KotaBadut KotaBadut Kota",
				"description":   "Desc",
				"scheduleBegin": "2020-12-30",
				"scheduleEnd":   "2021-01-30",
			})

			eventRepository := &repository.MockEventRepository{}
			endpoint := fmt.Sprintf("%s/api/v1/event/create", ts.URL)

			Convey("Then Fail Bad Request 400 fields not complete", func() {
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody404))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 400)
				So(string(body), ShouldEqual, `{"msg":"Key: 'EventRequest.LocationId' Error:Field validation for 'LocationId' failed on the 'required' tag\nKey: 'EventRequest.EventName' Error:Field validation for 'EventName' failed on the 'required' tag\nKey: 'EventRequest.Description' Error:Field validation for 'Description' failed on the 'required' tag\nKey: 'EventRequest.ScheduleBegin' Error:Field validation for 'ScheduleBegin' failed on the 'required' tag\nKey: 'EventRequest.ScheduleEnd' Error:Field validation for 'ScheduleEnd' failed on the 'required' tag"}`)
			})

			Convey("And Then Fail Bad Request 400 fields validation Country: greater than max ", func() {
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody404Validate))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 400)
				So(string(body), ShouldEqual, `{"msg":"EventName: greater than max"}`)
			})

			Convey("And Then Fail Query Location Data Return 0 and Error", func() {
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("1ad5ae0e-9e49-4025-90aa-295e1a4bd886").WillReturnError(fmt.Errorf("error")) //findLocation
				eventRepository.On("CreateEvent", mockTest.Anything).Return(nil, mockTest.Anything)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Fail Query Scan Location Empty Data Return 0 and Error", func() {
				rsLocation := sqlmock.NewRows([]string{"id", "country", "city_name", "village", "address"})
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("1ad5ae0e-9e49-4025-90aa-295e1a4bd886").WillReturnRows(rsLocation) //findLocation
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("1ad5ae0e-9e49-4025-90aa-295e1a4bd886",
					"Badut Kota",
					"2020-12-30",
					"2021-01-30").WillReturnError(fmt.Errorf("error")) //findEvent
				eventRepository.On("CreateEvent", mockTest.Anything).Return(nil, mockTest.Anything)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"Location not found for id 1ad5ae0e-9e49-4025-90aa-295e1a4bd886","status":500}`)
			})

			Convey("And Then Fail Query Scan Location Data Return 0 and Error", func() {
				rsLocation := sqlmock.NewRows([]string{"id", "country", "city_name"}).AddRow("id", "country", "city_name")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("1ad5ae0e-9e49-4025-90aa-295e1a4bd886").WillReturnRows(rsLocation) //findLocation
				/*mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("1ad5ae0e-9e49-4025-90aa-295e1a4bd886",
				"Badut Kota",
				"2020-12-30",
				"2021-01-30").WillReturnError(fmt.Errorf("error"))*/ //findEvent
				eventRepository.On("CreateEvent", mockTest.Anything).Return(nil, mockTest.Anything)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"sql: expected 3 destination arguments in Scan, not 5","status":500}`)
			})

			Convey("And Then Fail Query Exist Event Data Return 0 and Error", func() {
				rsLocation := sqlmock.NewRows([]string{"id", "country", "city_name", "village", "address"}).AddRow("id", "country", "city_name", "village", "address")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("1ad5ae0e-9e49-4025-90aa-295e1a4bd886").WillReturnRows(rsLocation) //findLocation
				scheduleBegin, err := time.Parse("2006-01-02", "2020-12-30")
				if err != nil {
					panic(err)
				}
				scheduleEnd, err := time.Parse("2006-01-02", "2021-01-30")
				if err != nil {
					panic(err)
				}
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs(
					"1ad5ae0e-9e49-4025-90aa-295e1a4bd886",
					"Badut Kota",
					scheduleBegin,
					scheduleEnd).WillReturnError(fmt.Errorf("error")) //findEvent
				eventRepository.On("CreateEvent", mockTest.Anything).Return(nil, mockTest.Anything)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Event Already Exist Return 0 and Error", func() {
				rsLocation := sqlmock.NewRows([]string{"id", "country", "city_name", "village", "address"}).AddRow("id", "country", "city_name", "village", "address")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("1ad5ae0e-9e49-4025-90aa-295e1a4bd886").WillReturnRows(rsLocation) //findLocation
				scheduleBegin, err := time.Parse("2006-01-02", "2020-12-30")
				if err != nil {
					panic(err)
				}
				scheduleEnd, err := time.Parse("2006-01-02", "2021-01-30")
				if err != nil {
					panic(err)
				}
				rsEvent := sqlmock.NewRows([]string{"event_name", "schedule_begin", "schedule_end", "location_id"}).AddRow("event_name", "schedule_begin", "schedule_end", "location_id")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs(
					"1ad5ae0e-9e49-4025-90aa-295e1a4bd886",
					"Badut Kota",
					scheduleBegin,
					scheduleEnd).WillReturnRows(rsEvent) //findEvent
				eventRepository.On("CreateEvent", mockTest.Anything).Return(nil, mockTest.Anything)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"Event already exist for city_name Badut Kota from 2020-12-30 to 2021-01-30 ","status":500}`)
			})

			Convey("And Then Data Not Exist and Fail BeginTx Return 0 and Error", func() {
				rsLocation := sqlmock.NewRows([]string{"id", "country", "city_name", "village", "address"}).AddRow("id", "country", "city_name", "village", "address")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("1ad5ae0e-9e49-4025-90aa-295e1a4bd886").WillReturnRows(rsLocation) //findLocation
				scheduleBegin, err := time.Parse("2006-01-02", "2020-12-30")
				if err != nil {
					panic(err)
				}
				scheduleEnd, err := time.Parse("2006-01-02", "2021-01-30")
				if err != nil {
					panic(err)
				}
				rsEvent := sqlmock.NewRows([]string{"event_name", "schedule_begin", "schedule_end", "location_id"})
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs(
					"1ad5ae0e-9e49-4025-90aa-295e1a4bd886",
					"Badut Kota",
					scheduleBegin,
					scheduleEnd).WillReturnRows(rsEvent) //findEvent
				eventRepository.On("CreateEvent", mockTest.Anything).Return(nil, mockTest.Anything)

				mock.ExpectBegin().WillReturnError(fmt.Errorf("error"))
				eventRepository.On("CreateEvent", mockTest.Anything).Return(nil, mockTest.Anything)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Data Not Exist and Fail ExpectExec Return 0 and Error", func() {
				rsLocation := sqlmock.NewRows([]string{"id", "country", "city_name", "village", "address"}).AddRow("id", "country", "city_name", "village", "address")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("1ad5ae0e-9e49-4025-90aa-295e1a4bd886").WillReturnRows(rsLocation) //findLocation
				scheduleBegin, err := time.Parse("2006-01-02", "2020-12-30")
				if err != nil {
					panic(err)
				}
				scheduleEnd, err := time.Parse("2006-01-02", "2021-01-30")
				if err != nil {
					panic(err)
				}
				rsEvent := sqlmock.NewRows([]string{"event_name", "schedule_begin", "schedule_end", "location_id"})
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs(
					"1ad5ae0e-9e49-4025-90aa-295e1a4bd886",
					"Badut Kota",
					scheduleBegin,
					scheduleEnd).WillReturnRows(rsEvent) //findEvent
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(``)).WillReturnError(fmt.Errorf("error"))
				eventRepository.On("CreateEvent", mockTest.Anything).Return(nil, mockTest.Anything)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Data Not Exist and Fail Commit Return 0 and Error", func() {
				rsLocation := sqlmock.NewRows([]string{"id", "country", "city_name", "village", "address"}).AddRow("id", "country", "city_name", "village", "address")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("1ad5ae0e-9e49-4025-90aa-295e1a4bd886").WillReturnRows(rsLocation) //findLocation
				scheduleBegin, err := time.Parse("2006-01-02", "2020-12-30")
				if err != nil {
					panic(err)
				}
				scheduleEnd, err := time.Parse("2006-01-02", "2021-01-30")
				if err != nil {
					panic(err)
				}
				rsEvent := sqlmock.NewRows([]string{"event_name", "schedule_begin", "schedule_end", "location_id"})
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs(
					"1ad5ae0e-9e49-4025-90aa-295e1a4bd886",
					"Badut Kota",
					scheduleBegin,
					scheduleEnd).WillReturnRows(rsEvent) //findEvent
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(``)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(fmt.Errorf("error"))
				eventRepository.On("CreateEvent", mockTest.Anything).Return(nil, mockTest.Anything)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Data Not Exist and Success Commit Return 1", func() {
				rsLocation := sqlmock.NewRows([]string{"id", "country", "city_name", "village", "address"}).AddRow("id", "country", "city_name", "village", "address")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("1ad5ae0e-9e49-4025-90aa-295e1a4bd886").WillReturnRows(rsLocation) //findLocation
				scheduleBegin, err := time.Parse("2006-01-02", "2020-12-30")
				if err != nil {
					panic(err)
				}
				scheduleEnd, err := time.Parse("2006-01-02", "2021-01-30")
				if err != nil {
					panic(err)
				}
				rsEvent := sqlmock.NewRows([]string{"event_name", "schedule_begin", "schedule_end", "location_id"})
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs(
					"1ad5ae0e-9e49-4025-90aa-295e1a4bd886",
					"Badut Kota",
					scheduleBegin,
					scheduleEnd).WillReturnRows(rsEvent) //findEvent
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(``)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				eventRepository.On("CreateEvent", mockTest.Anything).Return(1, nil)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 201)
				So(string(body), ShouldEqual, `{"message":"Event created successfully!","status":201}`)
			})
		})
	})
}
