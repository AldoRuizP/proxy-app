package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	handlers "github.com/AldoRuizP/proxy-app/api/handlers"
	server "github.com/AldoRuizP/proxy-app/api/server"
	utils "github.com/AldoRuizP/proxy-app/api/utils"
	"github.com/stretchr/testify/assert"
)

func init() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		utils.LoadEnv()
		app := server.SetUp()
		handlers.HandleRedirection(app)
		wg.Done()
		server.RunServer(app)
	}(wg)

	wg.Wait()
	fmt.Println("Server running...")
}

type Response struct {
	Status   int    `json:"status,omitempty"`
	Response string `json:"result,omitempty"`
}

func TestAlgorithm(t *testing.T) {

	cases := []struct {
		Domain   string
		Weight   string
		Priority string
		Output   string
	}{
		{Domain: "alpha", Output: `["alpha"]`},
		{Domain: "omega", Output: `["alpha","omega"]`},
		{Domain: "beta", Output: `["alpha","beta","omega"]`},
		{Domain: "omega", Output: `["alpha","beta","omega","omega"]`},
		{Domain: "alpha", Output: `["alpha","alpha","beta","omega","omega"]`},
		{Domain: "", Output: `no domain received`},
		{Domain: "someFakeDomain", Output: `domain not found`},
	}

	valuesToCompare := &Response{}
	client := http.Client{}

	for _, singleCase := range cases {
		req, errNewRequest := http.NewRequest("GET", "http://localhost:8080/ping", nil)
		req.Header.Add("domain", singleCase.Domain)
		res, errDoRequest := client.Do(req)
		bytes, errReadAll := ioutil.ReadAll(res.Body)
		errUnmarshal := json.Unmarshal(bytes, valuesToCompare)

		assert.Nil(t, errNewRequest)
		assert.Nil(t, errDoRequest)
		assert.Nil(t, errReadAll)
		assert.Nil(t, errUnmarshal)
		assert.Equal(t, singleCase.Output, valuesToCompare.Response)
		assert.True(t, true)
	}

}
