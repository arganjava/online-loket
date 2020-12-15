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

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLocationRouter_CreateLocation(t *testing.T) {
	Convey("Given Create Location", t, func() {
		Convey("When User Create Location", func() {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()
			ts := httptest.NewServer(routers.SetupServer(mockDB))
			defer ts.Close()
			requestBody, err := json.Marshal(map[string]string{
				"country":  "Indonesia",
				"cityName": "Bandung",
				"village":  "Ujung Berung",
				"address":  "test",
			})

			requestBody404, err := json.Marshal(map[string]string{
				"test": "test",
			})

			requestBody404Validate, err := json.Marshal(map[string]string{
				"country":  "IndonesiaIndonesiaIndonesiaIndonesiaIndonesiaIndonesiaIndonesiaIndonesiaIndonesiaIndonesia",
				"cityName": "Bandung",
				"village":  "Ujung Berung",
				"address":  "Test",
			})

			locationRepository := &repository.MockLocationRepository{}
			endpoint := fmt.Sprintf("%s/api/v1/location/create", ts.URL)

			Convey("Then Fail Bad Request 400 fields not complete", func() {
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody404))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 400)
				So(string(body), ShouldEqual, `{"msg":"Key: 'LocationRequest.Country' Error:Field validation for 'Country' failed on the 'required' tag\nKey: 'LocationRequest.CityName' Error:Field validation for 'CityName' failed on the 'required' tag\nKey: 'LocationRequest.Village' Error:Field validation for 'Village' failed on the 'required' tag\nKey: 'LocationRequest.Address' Error:Field validation for 'Address' failed on the 'required' tag"}`)
			})

			Convey("Then Fail Bad Request 400 fields validation Country: greater than max ", func() {
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody404Validate))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 400)
				So(string(body), ShouldEqual, `{"msg":"Country: greater than max"}`)
			})

			Convey("Then Fail Query Exist Data Return 0 and Error", func() {
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("Indonesia", "Bandung").WillReturnError(fmt.Errorf("error"))
				locationRepository.On("CreateLocation", mockTest.Anything).Return(nil, mockTest.Anything)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Data Already Exist Return 0 and Error", func() {
				rs := sqlmock.NewRows([]string{"ID", "Country", "CityName"}).AddRow("ID", "Country", "CityName")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("Indonesia", "Bandung").WillReturnRows(rs)
				locationRepository.On("CreateLocation", mockTest.Anything).Return(nil, mockTest.Anything)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"Location already exist for Indonesia Bandung Ujung Berung","status":500}`)
			})

			Convey("And Then Data Not Exist and Fail BeginTx Return 0 and Error", func() {
				rs := sqlmock.NewRows([]string{"ID", "Country", "CityName"})
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("Indonesia", "Bandung").WillReturnRows(rs)
				mock.ExpectBegin().WillReturnError(fmt.Errorf("error"))
				locationRepository.On("CreateLocation", mockTest.Anything).Return(nil, mockTest.Anything)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Data Not Exist and Fail ExpectExec Return 0 and Error", func() {
				rs := sqlmock.NewRows([]string{"ID", "Country", "CityName"})
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("Indonesia", "Bandung").WillReturnRows(rs)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(``)).WillReturnError(fmt.Errorf("error"))
				locationRepository.On("CreateLocation", mockTest.Anything).Return(nil, mockTest.Anything)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Data Not Exist and Fail Commit Return 0 and Error", func() {
				rs := sqlmock.NewRows([]string{"ID", "Country", "CityName"})
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("Indonesia", "Bandung").WillReturnRows(rs)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(``)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(fmt.Errorf("error"))
				locationRepository.On("CreateLocation", mockTest.Anything).Return(nil, mockTest.Anything)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Data Not Exist and Success Commit Return 1", func() {
				rs := sqlmock.NewRows([]string{"ID", "Country", "CityName"})
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("Indonesia", "Bandung").WillReturnRows(rs)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(``)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				locationRepository.On("CreateLocation", mockTest.Anything).Return(1, nil)
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 201)
				So(string(body), ShouldEqual, `{"message":"Location created successfully!","status":201}`)
			})
		})
	})
}
