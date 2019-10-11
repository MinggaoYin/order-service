package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const host = "http://localhost:8080"

func TestHome(t *testing.T) {
	t.Run("GET /", func(t *testing.T) {

		t.Run("it should not return error", func(t *testing.T) {
			resp, err := http.Get(host + "/")
			if assert.Nil(t, err, "it should not return error when visiting home page ") {
				assert.Equal(t, 200, resp.StatusCode, "it should return 200 as status code")
				assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"), "it should return application/json in header")
			}
		})
	})
}

func TestNotFoundURL(t *testing.T) {
	t.Run("GET /abc", func(t *testing.T) {

		t.Run("it should return 404 with json response", func(t *testing.T) {
			resp, err := http.Get(host + "/abc")
			if assert.Nil(t, err, "it should not return error when requesting undefined URL") {
				defer resp.Body.Close()

				assert.Equal(t, 404, resp.StatusCode, "it should return 404 as status code")
				assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"), "it should return application/json in header")

				var m map[string]interface{}
				bs, _ := ioutil.ReadAll(resp.Body)
				err = json.Unmarshal(bs, &m)
				assert.Nil(t, err, "it should return json response")
			}
		})
	})
}

func TestPlaceOrder(t *testing.T) {
	t.Run("POST /orders", func(t *testing.T) {

		t.Run("it should return 400 for invalid input", func(t *testing.T) {

			t.Run("missing origin", func(t *testing.T) {
				params := getPlaceOrderParams()
				delete(params, "origin")

				bs, _ := json.Marshal(params)
				resp, err := http.Post(host+"/orders", "application/json", bytes.NewBuffer(bs))
				if assert.Nil(t, err) {
					defer resp.Body.Close()
					assert.Equal(t, 400, resp.StatusCode)
				}
			})

			t.Run("missing destination", func(t *testing.T) {
				params := getPlaceOrderParams()
				delete(params, "destination")

				bs, _ := json.Marshal(params)
				resp, err := http.Post(host+"/orders", "application/json", bytes.NewBuffer(bs))
				if assert.Nil(t, err) {
					defer resp.Body.Close()
					assert.Equal(t, 400, resp.StatusCode)
				}
			})

			t.Run("invalid value type of latitude & latitude", func(t *testing.T) {
				params := getPlaceOrderParams()
				params["origin"] = []float64{1, 1}

				bs, _ := json.Marshal(params)
				resp, err := http.Post(host+"/orders", "application/json", bytes.NewBuffer(bs))
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})

			t.Run("non-numeric value of of latitude & latitude", func(t *testing.T) {
				params := getPlaceOrderParams()
				params["origin"] = []string{"a", "b"}

				bs, _ := json.Marshal(params)
				resp, err := http.Post(host+"/orders", "application/json", bytes.NewBuffer(bs))
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})

			t.Run("invalid range of latitude & latitude", func(t *testing.T) {
				params := getPlaceOrderParams()
				params["origin"] = []string{"9999999", "-99999999"}

				bs, _ := json.Marshal(params)
				resp, err := http.Post(host+"/orders", "application/json", bytes.NewBuffer(bs))
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})

			t.Run("origin or destination is not of length 2", func(t *testing.T) {
				params := getPlaceOrderParams()
				params["origin"] = []string{"22.286681", "114.193260", "114.193260"}

				bs, _ := json.Marshal(params)
				resp, err := http.Post(host+"/orders", "application/json", bytes.NewBuffer(bs))
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})

			t.Run("cannot calculate distance between origin and destination", func(t *testing.T) {
				params := getFarAwayPlaceOrderParams()

				bs, _ := json.Marshal(params)
				resp, err := http.Post(host+"/orders", "application/json", bytes.NewBuffer(bs))
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.Equal(t, "cannot calculate distance for given location", m["error"], "it should contain error message")
					}
				}
			})

		})

		t.Run("it should return the created order for valid input", func(t *testing.T) {

			params := getPlaceOrderParams()
			bs, _ := json.Marshal(params)
			resp, err := http.Post(host+"/orders", "application/json", bytes.NewBuffer(bs))

			if assert.Nil(t, err) {
				defer resp.Body.Close()

				assert.Equal(t, 200, resp.StatusCode, "it should return status code 200")
				assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"), "it should return application/json in header")

				var m map[string]interface{}
				bs, _ := ioutil.ReadAll(resp.Body)
				err = json.Unmarshal(bs, &m)
				if assert.Nil(t, err, "it should return a valid order json object") {
					_, ok := m["id"].(float64)
					assert.True(t, ok)

					assert.Equal(t, "UNASSIGNED", m["status"], "order status should be UNASSIGNED")

					_, ok = m["distance"].(int)
					assert.False(t, ok, "distance should be an integer")
				}
			}
		})
	})
}

func TestTakeOrder(t *testing.T) {
	t.Run("PATCH /orders/:id", func(t *testing.T) {

		t.Run("invalid order id", func(t *testing.T) {
			t.Run("missing order id", func(t *testing.T) {
				params := getTakeOrderParams()
				bs, _ := json.Marshal(params)

				req, _ := http.NewRequest("PATCH", host+"/orders", bytes.NewBuffer(bs))
				resp, err := http.DefaultClient.Do(req)
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 404, resp.StatusCode)

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})

			t.Run("invalid order id", func(t *testing.T) {
				params := getTakeOrderParams()
				bs, _ := json.Marshal(params)

				req, _ := http.NewRequest("PATCH", host+"/orders/abc", bytes.NewBuffer(bs))
				resp, err := http.DefaultClient.Do(req)
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 400, resp.StatusCode)

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.Equal(t, "order_id must be an integer", m["error"])
					}
				}
			})

			t.Run("order id not found", func(t *testing.T) {
				params := getTakeOrderParams()
				bs, _ := json.Marshal(params)

				req, _ := http.NewRequest("PATCH", host+"/orders/9999999999999", bytes.NewBuffer(bs))
				resp, err := http.DefaultClient.Do(req)
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 404, resp.StatusCode)

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.Equal(t, "Order not found", m["error"])
					}
				}
			})
		})

		t.Run("invalid request param", func(t *testing.T) {

			var orderId int

			// pre create order
			params := getPlaceOrderParams()
			bs, _ := json.Marshal(params)
			resp, err := http.Post(host+"/orders", "application/json", bytes.NewBuffer(bs))

			if assert.Nil(t, err) {
				defer resp.Body.Close()

				var m map[string]interface{}
				bs, _ := ioutil.ReadAll(resp.Body)
				err = json.Unmarshal(bs, &m)
				if assert.Nil(t, err) {
					id, ok := m["id"].(float64)
					assert.True(t, ok)
					orderId = int(id)
				}
			}

			t.Run("missing status attribute", func(t *testing.T) {
				params := getTakeOrderParams()
				delete(params, "status")

				bs, _ := json.Marshal(params)

				req, _ := http.NewRequest("PATCH", host+"/orders/" + fmt.Sprint(orderId), bytes.NewBuffer(bs))
				resp, err := http.DefaultClient.Do(req)
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 400, resp.StatusCode)

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})

			t.Run("invalid status attribute", func(t *testing.T) {
				params := getTakeOrderParams()
				params["status"] = "a"

				bs, _ := json.Marshal(params)

				req, _ := http.NewRequest("PATCH", host+"/orders/" + fmt.Sprint(orderId), bytes.NewBuffer(bs))
				resp, err := http.DefaultClient.Do(req)
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 400, resp.StatusCode)

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})
		})

		t.Run("take order with valid params", func(t *testing.T) {

			// pre create order
			var orderId int

			// pre create order
			params := getPlaceOrderParams()
			bs, _ := json.Marshal(params)
			resp, err := http.Post(host+"/orders", "application/json", bytes.NewBuffer(bs))

			if assert.Nil(t, err) {
				defer resp.Body.Close()

				var m map[string]interface{}
				bs, _ := ioutil.ReadAll(resp.Body)
				err = json.Unmarshal(bs, &m)
				if assert.Nil(t, err) {
					id, ok := m["id"].(float64)
					assert.True(t, ok)
					orderId = int(id)
				}
			}

			t.Run("SUCCESS", func(t *testing.T) {
				params := getTakeOrderParams()

				bs, _ := json.Marshal(params)

				req, _ := http.NewRequest("PATCH", host+"/orders/" + fmt.Sprint(orderId), bytes.NewBuffer(bs))
				resp, err := http.DefaultClient.Do(req)
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 200, resp.StatusCode)

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.Equal(t, "SUCCESS", m["status"])
					}
				}
			})
		})

		t.Run("race condition", func(t *testing.T) {

			// pre create order
			var orderId int

			// pre create order
			params := getPlaceOrderParams()
			bs, _ := json.Marshal(params)
			resp, err := http.Post(host+"/orders", "application/json", bytes.NewBuffer(bs))

			if assert.Nil(t, err) {
				defer resp.Body.Close()

				var m map[string]interface{}
				bs, _ := ioutil.ReadAll(resp.Body)
				err = json.Unmarshal(bs, &m)
				if assert.Nil(t, err) {
					id, ok := m["id"].(float64)
					assert.True(t, ok)
					orderId = int(id)
				}
			}

			t.Run("", func(t *testing.T) {
				params := getTakeOrderParams()

				bs, _ := json.Marshal(params)

				req1, _ := http.NewRequest("PATCH", host+"/orders/" + fmt.Sprint(orderId), bytes.NewBuffer(bs))
				req2, _ := http.NewRequest("PATCH", host+"/orders/" + fmt.Sprint(orderId), bytes.NewBuffer(bs))

				ch := make(chan *http.Response, 2)

				go func(ch chan *http.Response) {
					resp, err := http.DefaultClient.Do(req1)
					if assert.Nil(t, err) {
						defer resp.Body.Close()
						ch <- resp
					} else {
						ch <- nil
					}
				}(ch)

				go func(ch chan *http.Response) {
					resp, err := http.DefaultClient.Do(req2)
					if assert.Nil(t, err) {
						defer resp.Body.Close()
						ch <- resp
					} else {
						ch <- nil
					}
				}(ch)

				resp1 := <-ch
				resp2 := <-ch

				assert.Contains(t, []int{200, 409}, resp1.StatusCode)
				assert.Contains(t, []int{200, 409}, resp2.StatusCode)
				assert.True(t, resp1.StatusCode != resp2.StatusCode, "only one request should take the order")
			})
		})
	})
}

func TestListOrders(t *testing.T) {

	t.Run("GET /orders", func(t *testing.T) {

		t.Run("invalid page", func(t *testing.T) {

			t.Run("missing page", func(t *testing.T) {
				resp, err := http.Get(host + "/orders?limit=1")
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 400, resp.StatusCode, "status code should be 400")
					assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"), "it should return application/json in header")

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})

			t.Run("of type alphabetic character", func(t *testing.T) {
				resp, err := http.Get(host + "/orders?page=a&limit=1")
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 400, resp.StatusCode, "status code should be 400")
					assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"), "it should return application/json in header")

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})

			t.Run("of type float", func(t *testing.T) {
				resp, err := http.Get(host + "/orders?page=1.1&limit=1")
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 400, resp.StatusCode, "status code should be 400")
					assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"), "it should return application/json in header")

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})

			t.Run("negative page", func(t *testing.T) {
				resp, err := http.Get(host + "/orders?page=-1&limit=1")
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 400, resp.StatusCode, "status code should be 400")
					assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"), "it should return application/json in header")

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})

			t.Run("zero", func(t *testing.T) {
				resp, err := http.Get(host + "/orders?page=0&limit=1")
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 400, resp.StatusCode, "status code should be 400")
					assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"), "it should return application/json in header")

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})
		})

		t.Run("invalid limit", func(t *testing.T) {
			t.Run("missing limit", func(t *testing.T) {
				resp, err := http.Get(host + "/orders?page=1")
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 400, resp.StatusCode, "status code should be 400")
					assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"), "it should return application/json in header")

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})

			t.Run("of type alphabetic character", func(t *testing.T) {
				resp, err := http.Get(host + "/orders?page=1&limit=a")
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 400, resp.StatusCode, "status code should be 400")
					assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"), "it should return application/json in header")

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})

			t.Run("of type float", func(t *testing.T) {
				resp, err := http.Get(host + "/orders?page=1&limit=1.1")
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 400, resp.StatusCode, "status code should be 400")
					assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"), "it should return application/json in header")

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})

			t.Run("negative limit", func(t *testing.T) {
				resp, err := http.Get(host + "/orders?page=1&limit=-1")
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 400, resp.StatusCode, "status code should be 400")
					assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"), "it should return application/json in header")

					var m map[string]interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &m)
					if assert.Nil(t, err, "it should return json response") {
						assert.IsType(t, "", m["error"], "it should contain error message")
					}
				}
			})
		})

		t.Run("it should return created orders with valid input", func(t *testing.T) {
			t.Run("return empty list when there is no more orders at that page", func(t *testing.T) {
				resp, err := http.Get(host + "/orders?page=9999999999&limit=10")
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 200, resp.StatusCode)
					assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"))

					var orders []interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &orders)
					if assert.Nil(t, err, "it should return json response") {
						assert.Equal(t, 0, len(orders), "length of orders should be 0")
					}
				}
			})

			t.Run("return orders given a large limit number", func(t *testing.T) {
				resp, err := http.Get(host + "/orders?page=1&limit=999999999")
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 200, resp.StatusCode)
					assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"))
				}
			})

			t.Run("number of order should less than limit", func(t *testing.T) {
				limit := 1
				resp, err := http.Get(host + "/orders?page=1&limit=" + fmt.Sprint(limit))
				if assert.Nil(t, err) {
					defer resp.Body.Close()

					assert.Equal(t, 200, resp.StatusCode)
					assert.Equal(t, "application/json; charset=UTF-8", resp.Header.Get("Content-Type"))

					var orders []interface{}
					bs, _ := ioutil.ReadAll(resp.Body)
					err = json.Unmarshal(bs, &orders)
					if assert.Nil(t, err, "it should return json response") {
						assert.Equal(t, true, len(orders) <= limit, "number of order should less than or equal to limit")
					}
				}
			})
		})
	})
}

func getPlaceOrderParams() map[string]interface{} {
	params := make(map[string]interface{})

	origin := []string{"22.286681", "114.193260"}
	destination := []string{"22.279707", "114.186301"}

	params["origin"] = origin
	params["destination"] = destination

	return params
}

func getFarAwayPlaceOrderParams() map[string]interface{} {
	params := make(map[string]interface{})

	origin := []string{"22.780247", "113.687473"}
	destination := []string{"22.217851", "114.207989"}

	params["origin"] = origin
	params["destination"] = destination

	return params
}

func getTakeOrderParams() map[string]interface{} {
	params := make(map[string]interface{})
	params["status"] = "TAKEN"
	return params
}
