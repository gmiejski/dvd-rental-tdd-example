package infrastructure

import (
	"fmt"
	"reflect"
	"testing"
)

func TestA(t *testing.T) {

	e := getEventsList()
	for _, event := range e {
		fmt.Println(reflect.TypeOf(event).Name())
	}

}
