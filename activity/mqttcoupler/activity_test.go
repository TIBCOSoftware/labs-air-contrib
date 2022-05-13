package mqttcoupler

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
)

func Pour(port string) {
	for {
		conn, _ := net.Dial("tcp", net.JoinHostPort("", port))
		if conn != nil {
			conn.Close()
			break
		}
	}
}

func TestParseTopic(t *testing.T) {
	test := func(input, output string, params map[string]string) {
		parsed := ParseTopic(input)
		assert.Equal(t, parsed.String(params), output)
	}
	test("/a/:x/b/:y", "/a/test/b/j/k", map[string]string{"x": "test", "y": "j/k"})
	test("/a/:/b/:", "/a/test/b/j/k", map[string]string{"0": "test", "1": "j/k"})
	test("a/:/b/:", "a/test/b/j/k", map[string]string{"0": "test", "1": "j/k"})
	test("a/:/b", "a/test/b", map[string]string{"0": "test"})
	test("a/:/b/", "a/test/b/", map[string]string{"0": "test"})
	test("", "", map[string]string{})
	test(":", "test", map[string]string{"0": "test"})
	test(":", "test", map[string]string{"0": "test"})
	test("/", "/", map[string]string{})
	test("/:", "/test", map[string]string{"0": "test"})
	test("/:", "/test", map[string]string{"0": "test"})
	test("/a/b", "/a/b", map[string]string{})
}

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func TestEval(t *testing.T) {
	settings := Settings{
		Broker:          "tcp://localhost:1883",
		Id:              "TestX",
		Topic:           "stmrequest",
		ResponseTimeout: 50,
	}
	init := test.NewActivityInitContext(settings, nil)
	act, err := New(init)
	assert.Nil(t, err)
	context := test.NewActivityContext(activityMd)
	context.SetInput("message", `{"message": "hello world"}`)

	done := false
	for {
		done, err = act.Eval(context)
		fmt.Println(context.GetOutput("data"))

		time.Sleep(time.Duration(10) * time.Second)
	}

	assert.True(t, done)
	assert.Nil(t, err)

}
