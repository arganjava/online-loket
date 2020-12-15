package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/arganjava/online-loket/src/routers"
	. "github.com/smartystreets/goconvey/convey"
	"time"

	//mockTest "github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestTransactionRouter_CreatePurchaseTransaction(t *testing.T) {
	Convey("Given Create Transaction Purchase", t, func() {
		Convey("When User Create Transaction Purchase", func() {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()
			ts := httptest.NewServer(routers.SetupServer(mockDB))
			defer ts.Close()
			requestBody, err := json.Marshal(map[string]interface{}{
				"eventTicketId": "ae5733f8-3d20-40d5-b6b2-caa56e2d36c8",
				"customerName":  "Argan Megariansyah",
				"customerPhone": "0123232",
				"customerEmail": "arganjava@gmail.com",
				"orderQuantity": 2,
			})

			requestBody404, err := json.Marshal(map[string]interface{}{
				"test": "test",
			})

			requestBody404Validate, err := json.Marshal(map[string]interface{}{
				"eventTicketId": "ae5733f8-3d20-40d5-b6b2-caa56e2d36c8",
				"customerName":  "Argan MegariansyahArgan MegariansyahArgan MegariansyahArgan MegariansyahArgan MegariansyahArgan Megariansyah",
				"customerPhone": "0123232",
				"customerEmail": "arganjava@gmail.com",
				"orderQuantity": 2,
			})

			endpoint := fmt.Sprintf("%s/api/v1/transaction/purchase", ts.URL)

			Convey("Then Fail Bad Request 400 fields not complete", func() {
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody404))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 400)
				So(string(body), ShouldEqual, `{"msg":"Key: 'TransactionRequest.EventTicketId' Error:Field validation for 'EventTicketId' failed on the 'required' tag\nKey: 'TransactionRequest.CustomerName' Error:Field validation for 'CustomerName' failed on the 'required' tag\nKey: 'TransactionRequest.CustomerPhone' Error:Field validation for 'CustomerPhone' failed on the 'required' tag\nKey: 'TransactionRequest.CustomerEmail' Error:Field validation for 'CustomerEmail' failed on the 'required' tag\nKey: 'TransactionRequest.OrderQuantity' Error:Field validation for 'OrderQuantity' failed on the 'required' tag"}`)
			})

			Convey("And Then Fail Bad Request 400 fields validation CustomerName: greater than max ", func() {
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody404Validate))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 400)
				So(string(body), ShouldEqual, `{"msg":"CustomerName: greater than max"}`)
			})

			Convey("And Then Fail Tx Begin Ticket Data Return 0 and Error", func() {
				mock.ExpectBegin().WillReturnError(fmt.Errorf("error"))
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Fail Query Tiket Return 0 and Error", func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("ae5733f8-3d20-40d5-b6b2-caa56e2d36c8").WillReturnError(fmt.Errorf("error")) //findTicket
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Query Scan Tiket Empty Data Return 0 and Error", func() {
				mock.ExpectBegin()
				rsTicket := sqlmock.NewRows([]string{"id", "ticket_type", "quantity", "price", "event_id"})
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("ae5733f8-3d20-40d5-b6b2-caa56e2d36c8").WillReturnRows(rsTicket) //findTicket
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"Event Ticket not found for id ae5733f8-3d20-40d5-b6b2-caa56e2d36c8 ","status":500}`)
			})

			Convey("And Then Query Scan Tiket Not Enough Quota Return 0 and Error", func() {
				mock.ExpectBegin()
				rsTicket := sqlmock.NewRows([]string{"id", "ticket_type", "quantity", "price", "event_id"}).AddRow("id", "ticket_type", 0, 1000, "event_id")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("ae5733f8-3d20-40d5-b6b2-caa56e2d36c8").WillReturnRows(rsTicket) //findTicket
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"Event Ticket not enough quota for id ","status":500}`)
			})

			Convey("And Then Fail ExpectExec Insert Transaction Return 0 and Error", func() {
				mock.ExpectBegin()
				rsTicket := sqlmock.NewRows([]string{"id", "ticket_type", "quantity", "price", "event_id"}).AddRow("id", "ticket_type", 100, 1000, "event_id")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("ae5733f8-3d20-40d5-b6b2-caa56e2d36c8").WillReturnRows(rsTicket) //findTicket

				mock.ExpectExec(regexp.QuoteMeta(``)).WillReturnError(fmt.Errorf("error"))
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Fail ExpectExec Update Ticket Quota Return 0 and Error", func() {
				mock.ExpectBegin()
				rsTicket := sqlmock.NewRows([]string{"id", "ticket_type", "quantity", "price", "event_id"}).AddRow("id", "ticket_type", 100, 1000, "event_id")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("ae5733f8-3d20-40d5-b6b2-caa56e2d36c8").WillReturnRows(rsTicket) //findTicket

				mock.ExpectExec(regexp.QuoteMeta(``)).WillReturnResult(sqlmock.NewResult(1, 1)) //insert transaction
				mock.ExpectExec(regexp.QuoteMeta(``)).WillReturnError(fmt.Errorf("error"))      //update quota ticket

				//mock.ExpectCommit().WillReturnError(fmt.Errorf("error"))
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Fail Commit Return 0 and Error", func() {
				mock.ExpectBegin()
				rsTicket := sqlmock.NewRows([]string{"id", "ticket_type", "quantity", "price", "event_id"}).AddRow("id", "ticket_type", 100, 1000, "event_id")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("ae5733f8-3d20-40d5-b6b2-caa56e2d36c8").WillReturnRows(rsTicket) //findTicket

				mock.ExpectExec(regexp.QuoteMeta(``)).WillReturnResult(sqlmock.NewResult(1, 1)) //insert transaction
				mock.ExpectExec(regexp.QuoteMeta(``)).WillReturnResult(sqlmock.NewResult(1, 1)) //update quota ticket

				mock.ExpectCommit().WillReturnError(fmt.Errorf("error"))
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Data Not Exist and Success Commit Return 1", func() {
				mock.ExpectBegin()
				rsTicket := sqlmock.NewRows([]string{"id", "ticket_type", "quantity", "price", "event_id"}).AddRow("id", "ticket_type", 100, 1000, "event_id")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("ae5733f8-3d20-40d5-b6b2-caa56e2d36c8").WillReturnRows(rsTicket) //findTicket

				mock.ExpectExec(regexp.QuoteMeta(``)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(``)).WillReturnResult(sqlmock.NewResult(1, 1)) //update quota ticket
				mock.ExpectCommit()
				respApi, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 201)
				So(string(body), ShouldEqual, `{"message":"Transaction created successfully!","status":201}`)
			})

		})
	})
}

func TestTransactionRouter_GetInfoPurchaseTransaction(t *testing.T) {
	Convey("Given Get Info Transaction Purchase", t, func() {
		Convey("When User Get Info Transaction Purchase", func() {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()
			ts := httptest.NewServer(routers.SetupServer(mockDB))
			defer ts.Close()
			id := "ae5733f8-3d20-40d5-b6b2-caa56e2d36c8"
			endpoint := fmt.Sprintf("%s/api/v1/transaction/get_info/%v", ts.URL, id)
			//eventRepo := &repository.MockEventRepository{}
			Convey("And Then Fail Query purchase Return Error", func() {
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("ae5733f8-3d20-40d5-b6b2-caa56e2d36c8").WillReturnError(fmt.Errorf("error")) //findTransaction
				respApi, _ := http.Get(endpoint)
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Data Not Found Return Success", func() {
				rsTransaction := sqlmock.NewRows([]string{"id", "customer_name", "customer_phone", "customer_email", "order_quantity", "transaction_time", "event_ticket_id"})
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("ae5733f8-3d20-40d5-b6b2-caa56e2d36c8").WillReturnRows(rsTransaction) //findTransaction
				respApi, _ := http.Get(endpoint)
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 200)
				So(string(body), ShouldEqual, `{"data":null,"message":"Data not found","status":200}`)
			})

			Convey("And Then Data Found Return Error", func() {
				rsTransaction := sqlmock.NewRows([]string{"id", "customer_name", "customer_phone", "customer_email", "order_quantity", "transaction_time", "event_ticket_id"}).AddRow("id", "customer_name", "customer_phone", "customer_email", "order_quantity", "transaction_time", "event_ticket_id")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("ae5733f8-3d20-40d5-b6b2-caa56e2d36c8").WillReturnRows(rsTransaction) //findTransaction
				respApi, _ := http.Get(endpoint)
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"sql: Scan error on column index 4, name \"order_quantity\": converting driver.Value type string (\"order_quantity\") to a int64: invalid syntax","status":500}`)
			})

			Convey("And Then Data Event Ticket Error Return Error", func() {
				rsTransaction := sqlmock.NewRows([]string{"id", "customer_name", "customer_phone", "customer_email", "order_quantity", "transaction_time", "event_ticket_id"}).AddRow("id", "customer_name", "customer_phone", "customer_email", int64(10), time.Now(), "event_ticket_id")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("ae5733f8-3d20-40d5-b6b2-caa56e2d36c8").WillReturnRows(rsTransaction) //findTransaction

				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("event_ticket_id").WillReturnError(fmt.Errorf("error")) //findTransaction
				respApi, _ := http.Get(endpoint)
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 500)
				So(string(body), ShouldEqual, `{"message":"error","status":500}`)
			})

			Convey("And Then Data Found Return Success", func() {
				trxTime, _ := time.Parse("2006-01-02", "2020-12-30")

				rsTransaction := sqlmock.NewRows([]string{"id", "customer_name", "customer_phone", "customer_email", "order_quantity", "transaction_time", "event_ticket_id"}).AddRow("id", "customer_name", "customer_phone", "customer_email", int64(10), trxTime, "event_ticket_id")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("ae5733f8-3d20-40d5-b6b2-caa56e2d36c8").WillReturnRows(rsTransaction) //findTransaction

				rsTicket := sqlmock.NewRows([]string{"id", "ticket_type", "quantity", "price", "event_id"}).AddRow("id", "ticket_type", int64(10), float64(1000), "event_id")
				mock.ExpectQuery(regexp.QuoteMeta(``)).WithArgs("event_ticket_id").WillReturnRows(rsTicket) //findTicket
				respApi, _ := http.Get(endpoint)
				body, _ := ioutil.ReadAll(respApi.Body)
				So(respApi.StatusCode, ShouldEqual, 200)
				So(string(body), ShouldEqual, `{"data":{"id":"id","eventTicket":{"id":"id","type":"ticket_type","quantity":10,"price":1000},"customerName":"customer_name","customerPhone":"customer_phone","customerEmail":"customer_email","orderQuantity":10,"transactionTime":"`+trxTime.Format("2006-01-02")+"T00:00:00Z"+`"},"message":"Transaction Get successfully!","status":200}`)
			})

		})
	})
}
