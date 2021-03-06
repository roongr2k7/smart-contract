package did

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/tendermint/abci/example/code"
	"github.com/tendermint/abci/types"
)

// TODO: unit testing
func addNodePublicKey(param string, app *DIDApplication) types.ResponseDeliverTx {
	fmt.Println("AddNodePublicKey")
	var nodePublicKey NodePublicKey
	err := json.Unmarshal([]byte(param), &nodePublicKey)
	if err != nil {
		return ReturnDeliverTxLog(err.Error())
	}

	key := "NodePublicKey" + "|" + nodePublicKey.NodeID
	value := nodePublicKey.PublicKey
	app.state.Size++
	app.state.db.Set(prefixKey([]byte(key)), []byte(value))
	return ReturnDeliverTxLog("success")
}

func registerMsqDestination(param string, app *DIDApplication) types.ResponseDeliverTx {
	fmt.Println("RegisterMsqDestination")
	var funcParam RegisterMsqDestinationParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(err.Error())
	}

	for _, user := range funcParam.Users {
		key := "MsqDestination" + "|" + user.HashID
		chkExists := app.state.db.Get(prefixKey([]byte(key)))

		if chkExists != nil {
			var nodes []Node
			err = json.Unmarshal([]byte(chkExists), &nodes)
			if err != nil {
				return ReturnDeliverTxLog(err.Error())
			}

			newNode := Node{user.Ial, funcParam.NodeID}
			// Check duplicate before add
			chkDup := false
			for _, node := range nodes {
				if newNode == node {
					chkDup = true
					break
				}
			}

			if chkDup == false {
				nodes = append(nodes, newNode)
				value, err := json.Marshal(nodes)
				if err != nil {
					return ReturnDeliverTxLog(err.Error())
				}
				app.state.Size++
				app.state.db.Set(prefixKey([]byte(key)), []byte(value))
			}

		} else {
			var nodes []Node
			newNode := Node{user.Ial, funcParam.NodeID}
			nodes = append(nodes, newNode)
			value, err := json.Marshal(nodes)
			if err != nil {
				return ReturnDeliverTxLog(err.Error())
			}
			app.state.Size++
			app.state.db.Set(prefixKey([]byte(key)), []byte(value))
		}
	}

	return ReturnDeliverTxLog("success")
}

func addAccessorMethod(param string, app *DIDApplication) types.ResponseDeliverTx {
	fmt.Println("AddAccessorMethod")
	var accessorMethod AccessorMethod
	err := json.Unmarshal([]byte(param), &accessorMethod)
	if err != nil {
		return ReturnDeliverTxLog(err.Error())
	}

	key := "AccessorMethod" + "|" + accessorMethod.AccessorID
	value, err := json.Marshal(accessorMethod)
	if err != nil {
		return ReturnDeliverTxLog(err.Error())
	}
	app.state.Size++
	app.state.db.Set(prefixKey([]byte(key)), []byte(value))
	return ReturnDeliverTxLog("success")
}

func createRequest(param string, app *DIDApplication) types.ResponseDeliverTx {
	fmt.Println("CreateRequest")
	var request Request
	err := json.Unmarshal([]byte(param), &request)
	if err != nil {
		return ReturnDeliverTxLog(err.Error())
	}

	key := "Request" + "|" + request.RequestID
	value, err := json.Marshal(request)
	if err != nil {
		return ReturnDeliverTxLog(err.Error())
	}

	existValue := app.state.db.Get(prefixKey([]byte(key)))
	if existValue != nil {
		return ReturnDeliverTxLog("Duplicate Request ID")
	}

	app.state.Size++
	app.state.db.Set(prefixKey([]byte(key)), []byte(value))

	// callback to IDP
	uri := getEnv("CALLBACK_URI", "")
	if uri != "" {
		go callBack(uri, request.RequestID, app.state.Height)
	}
	return ReturnDeliverTxLog("success")
}

func createIdpResponse(param string, app *DIDApplication) types.ResponseDeliverTx {
	fmt.Println("CreateIdpResponse")
	var response Response
	err := json.Unmarshal([]byte(param), &response)
	if err != nil {
		return ReturnDeliverTxLog(err.Error())
	}

	key := "Request" + "|" + response.RequestID
	value := app.state.db.Get(prefixKey([]byte(key)))

	if value == nil {
		return ReturnDeliverTxLog("Request ID not found")
	}
	var request Request
	err = json.Unmarshal([]byte(value), &request)
	if err != nil {
		return ReturnDeliverTxLog(err.Error())
	}

	// Check duplicate before add
	chkDup := false
	for _, oldResponse := range request.Responses {
		if response == oldResponse {
			chkDup = true
			break
		}
	}

	if chkDup == false {
		request.Responses = append(request.Responses, response)
		value, err := json.Marshal(request)
		if err != nil {
			return ReturnDeliverTxLog(err.Error())
		}
		app.state.Size++
		app.state.db.Set(prefixKey([]byte(key)), []byte(value))

		// callback to RP
		uri := getEnv("CALLBACK_URI", "")
		if uri != "" {
			go callBack(uri, request.RequestID, app.state.Height)
		}

		return ReturnDeliverTxLog("success")
	}
	return ReturnDeliverTxLog("Response duplicate")
}

func signData(param string, app *DIDApplication) types.ResponseDeliverTx {
	fmt.Println("SignData")
	var signData SignDataParam
	err := json.Unmarshal([]byte(param), &signData)
	if err != nil {
		return ReturnDeliverTxLog(err.Error())
	}

	key := "SignData" + "|" + signData.Signature
	value, err := json.Marshal(signData)
	if err != nil {
		return ReturnDeliverTxLog(err.Error())
	}

	app.state.Size++
	app.state.db.Set(prefixKey([]byte(key)), []byte(value))
	return ReturnDeliverTxLog("success")
}

func registerServiceDestination(param string, app *DIDApplication) types.ResponseDeliverTx {
	fmt.Println("RegisterServiceDestination")
	var funcParam RegisterServiceDestinationParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(err.Error())
	}

	key := "ServiceDestination" + "|" + funcParam.AsID + "|" + funcParam.AsServiceID
	var node GetServiceDestinationResult
	node.NodeID = funcParam.NodeID
	value, err := json.Marshal(node)
	if err != nil {
		return ReturnDeliverTxLog(err.Error())
	}
	app.state.Size++
	app.state.db.Set(prefixKey([]byte(key)), []byte(value))
	return ReturnDeliverTxLog("success")
}

func callBack(uri string, requestID string, height int64) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	fmt.Println("CALLBACK_URI:" + uri)
	var callback Callback
	callback.RequestID = requestID
	callback.Height = height
	data, err := json.Marshal(callback)
	if err != nil {
		fmt.Println("error:", err)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("POST", uri, strings.NewReader(string(data)))
	if err != nil {
		fmt.Println("error:", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)
	fmt.Println(resp.Status)
}

// ReturnDeliverTxLog return types.ResponseDeliverTx
func ReturnDeliverTxLog(log string) types.ResponseDeliverTx {
	return types.ResponseDeliverTx{
		Code: code.CodeTypeOK,
		Log:  fmt.Sprintf(log)}
}

// DeliverTxRouter is Pointer to function
func DeliverTxRouter(method string, param string, app *DIDApplication) types.ResponseDeliverTx {
	funcs := map[string]interface{}{
		"AddNodePublicKey":           addNodePublicKey,
		"RegisterMsqDestination":     registerMsqDestination,
		"AddAccessorMethod":          addAccessorMethod,
		"CreateRequest":              createRequest,
		"CreateIdpResponse":          createIdpResponse,
		"SignData":                   signData,
		"RegisterServiceDestination": registerServiceDestination,
	}
	value, _ := callDeliverTx(funcs, method, param, app)
	return value[0].Interface().(types.ResponseDeliverTx)
}

func callDeliverTx(m map[string]interface{}, name string, param string, app *DIDApplication) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	in := make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(param)
	in[1] = reflect.ValueOf(app)
	result = f.Call(in)
	return
}
