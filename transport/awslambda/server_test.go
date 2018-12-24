package awslambda

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type key int

const (
	KeyBeforeOne key = iota
	KeyBeforeTwo key = iota
	KeyAfterOne  key = iota
	KeyEncMode   key = iota
)

func TestInvokeHappyPath(t *testing.T) {
	svc := serviceTest01{}

	helloHandler := NewServer(
		makeTest01HelloEndpoint(svc),
		decodeHelloRequest,
		encodeResponse,
		ServerErrorLogger(log.NewNopLogger()),
		ServerBefore(func(
			ctx context.Context, payload []byte,
		) context.Context {
			ctx = context.WithValue(ctx, KeyBeforeOne, "bef1")
			return ctx
		}),
		ServerBefore(func(
			ctx context.Context, payload []byte,
		) context.Context {
			ctx = context.WithValue(ctx, KeyBeforeTwo, "bef2")
			return ctx
		}),
		ServerAfter(func(
			ctx context.Context, response interface{},
		) context.Context {
			ctx = context.WithValue(ctx, KeyAfterOne, "af1")
			return ctx
		}),
		ServerAfter(func(
			ctx context.Context, response interface{},
		) context.Context {
			if _, ok := ctx.Value(KeyAfterOne).(string); !ok {
				t.Fatalf("\nValue was not set properly during multi ServerAfter")
			}
			return ctx
		}),
		ServerFinalizer(func(
			_ context.Context, resp []byte, _ error,
		) {
			apigwResp := events.APIGatewayProxyResponse{}
			err := json.Unmarshal(resp, &apigwResp)
			if err != nil {
				t.Fatalf("\nshould have no error, but got: %+v", err)
			}

			response := helloResponse{}
			err = json.Unmarshal([]byte(apigwResp.Body), &response)
			if err != nil {
				t.Fatalf("\nshould have no error, but got: %+v", err)
			}

			expectedGreeting := "hello john doe bef1 bef2"
			if response.Greeting != expectedGreeting {
				t.Fatalf(
					"\nexpect: %s\nactual: %s", expectedGreeting, response.Greeting)
			}
		}),
	)

	ctx := context.Background()
	req, _ := json.Marshal(events.APIGatewayProxyRequest{
		Body: `{"name":"john doe"}`,
	})
	resp, err := helloHandler.Invoke(ctx, req)

	if err != nil {
		t.Fatalf("\nshould have no error, but got: %+v", err)
	}

	apigwResp := events.APIGatewayProxyResponse{}
	err = json.Unmarshal(resp, &apigwResp)
	if err != nil {
		t.Fatalf("\nshould have no error, but got: %+v", err)
	}

	response := helloResponse{}
	err = json.Unmarshal([]byte(apigwResp.Body), &response)
	if err != nil {
		t.Fatalf("\nshould have no error, but got: %+v", err)
	}

	expectedGreeting := "hello john doe bef1 bef2"
	if response.Greeting != expectedGreeting {
		t.Fatalf(
			"\nexpect: %s\nactual: %s", expectedGreeting, response.Greeting)
	}
}

func TestInvokeFailDecode(t *testing.T) {
	svc := serviceTest01{}

	helloHandler := NewServer(
		makeTest01HelloEndpoint(svc),
		decodeHelloRequest,
		encodeResponse,
		ServerErrorEncoder(func(
			ctx context.Context, err error,
		) ([]byte, error) {
			apigwResp := events.APIGatewayProxyResponse{}
			apigwResp.Body = `{"error":"yes"}`
			apigwResp.StatusCode = 500
			resp, merr := json.Marshal(apigwResp)
			if merr != nil {
				return resp, merr
			}
			return resp, err
		}),
	)

	ctx := context.Background()
	req, _ := json.Marshal(events.APIGatewayProxyRequest{
		Body: `{"name":"john doe"}`,
	})
	resp, err := helloHandler.Invoke(ctx, req)

	if err == nil {
		t.Fatalf("\nshould have error, but got: %+v", err)
	}

	apigwResp := events.APIGatewayProxyResponse{}
	json.Unmarshal(resp, &apigwResp)
	if apigwResp.StatusCode != 500 {
		t.Fatalf("\nexpect status code of 500, instead of %d", apigwResp.StatusCode)
	}
}

func TestInvokeFailEndpoint(t *testing.T) {
	svc := serviceTest01{}

	helloHandler := NewServer(
		makeTest01FailEndpoint(svc),
		decodeHelloRequest,
		encodeResponse,
		ServerBefore(func(
			ctx context.Context, payload []byte,
		) context.Context {
			ctx = context.WithValue(ctx, KeyBeforeOne, "bef1")
			return ctx
		}),
		ServerBefore(func(
			ctx context.Context, payload []byte,
		) context.Context {
			ctx = context.WithValue(ctx, KeyBeforeTwo, "bef2")
			return ctx
		}),
		ServerErrorEncoder(func(
			ctx context.Context, err error,
		) ([]byte, error) {
			apigwResp := events.APIGatewayProxyResponse{}
			apigwResp.Body = `{"error":"yes"}`
			apigwResp.StatusCode = 500
			resp, merr := json.Marshal(apigwResp)
			if merr != nil {
				return resp, merr
			}
			return resp, err
		}),
	)

	ctx := context.Background()
	req, _ := json.Marshal(events.APIGatewayProxyRequest{
		Body: `{"name":"john doe"}`,
	})
	resp, err := helloHandler.Invoke(ctx, req)

	if err == nil {
		t.Fatalf("\nshould have error, but got: %+v", err)
	}

	apigwResp := events.APIGatewayProxyResponse{}
	json.Unmarshal(resp, &apigwResp)
	if apigwResp.StatusCode != 500 {
		t.Fatalf("\nexpect status code of 500, instead of %d", apigwResp.StatusCode)
	}
}

func TestInvokeFailEncode(t *testing.T) {
	svc := serviceTest01{}

	helloHandler := NewServer(
		makeTest01HelloEndpoint(svc),
		decodeHelloRequest,
		encodeResponse,
		ServerBefore(func(
			ctx context.Context, payload []byte,
		) context.Context {
			ctx = context.WithValue(ctx, KeyBeforeOne, "bef1")
			return ctx
		}),
		ServerBefore(func(
			ctx context.Context, payload []byte,
		) context.Context {
			ctx = context.WithValue(ctx, KeyBeforeTwo, "bef2")
			return ctx
		}),
		ServerAfter(func(
			ctx context.Context, response interface{},
		) context.Context {
			ctx = context.WithValue(ctx, KeyEncMode, "fail_encode")
			return ctx
		}),
		ServerErrorEncoder(func(
			ctx context.Context, err error,
		) ([]byte, error) {
			apigwResp := events.APIGatewayProxyResponse{}
			apigwResp.Body = `{"error":"yes"}`
			apigwResp.StatusCode = 500
			resp, merr := json.Marshal(apigwResp)
			if merr != nil {
				return resp, merr
			}
			return resp, err
		}),
	)

	ctx := context.Background()
	req, _ := json.Marshal(events.APIGatewayProxyRequest{
		Body: `{"name":"john doe"}`,
	})
	resp, err := helloHandler.Invoke(ctx, req)

	if err == nil {
		t.Fatalf("\nshould have error, but got: %+v", err)
	}

	apigwResp := events.APIGatewayProxyResponse{}
	json.Unmarshal(resp, &apigwResp)
	if apigwResp.StatusCode != 500 {
		t.Fatalf("\nexpect status code of 500, instead of %d", apigwResp.StatusCode)
	}
}

func decodeHelloRequest(
	ctx context.Context, req []byte,
) (interface{}, error) {
	apigwReq := events.APIGatewayProxyRequest{}
	err := json.Unmarshal([]byte(req), &apigwReq)
	if err != nil {
		return apigwReq, err
	}

	request := helloRequest{}
	err = json.Unmarshal([]byte(apigwReq.Body), &request)
	if err != nil {
		return request, err
	}

	valOne, ok := ctx.Value(KeyBeforeOne).(string)
	if !ok {
		return request, fmt.Errorf(
			"Value was not set properly when multiple ServerBefores are used")
	}

	valTwo, ok := ctx.Value(KeyBeforeTwo).(string)
	if !ok {
		return request, fmt.Errorf(
			"Value was not set properly when multiple ServerBefores are used")
	}

	request.Name += " " + valOne + " " + valTwo
	return request, err
}

func encodeResponse(
	ctx context.Context, response interface{},
) ([]byte, error) {
	apigwResp := events.APIGatewayProxyResponse{}

	mode, ok := ctx.Value(KeyEncMode).(string)
	fmt.Printf("\nmode: %s ok: %+v\n", mode, ok)
	if ok && mode == "fail_encode" {
		fmt.Printf("\nEnter\n")
		return nil, fmt.Errorf("fail encoding")
	}

	respByte, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	apigwResp.Body = string(respByte)
	apigwResp.StatusCode = 200

	resp, err := json.Marshal(apigwResp)
	return resp, err
}

type helloRequest struct {
	Name string `json:"name"`
}

type helloResponse struct {
	Greeting string `json:"greeting"`
}

func makeTest01HelloEndpoint(svc serviceTest01) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(helloRequest)
		greeting := svc.hello(req.Name)
		return helloResponse{greeting}, nil
	}
}

func makeTest01FailEndpoint(_ serviceTest01) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return nil, fmt.Errorf("test error endpoint")
	}
}

type serviceTest01 struct{}

func (ts *serviceTest01) hello(name string) string {
	return fmt.Sprintf("hello %s", name)
}
