package tests

import (
	"fmt"
	logic "medods_jwt_service/logic"
	model "medods_jwt_service/model"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWT_possibleUserStates(t *testing.T) {
	user := &model.User{Name: "maxim", Email: "something@ya.ru", GUID: 1, Created_at: time.Now()}
	element := reflect.ValueOf(user).Elem()
	fields := element.NumField()
	for i := 0; i < fields; i++ {
		if !element.Field(i).CanSet() {
			continue
		}
		saved := element.Field(i)
		element.Field(i).Set(reflect.Zero(saved.Type()))
		RunJWTTestCase(*user, t, element.Type().Field(i).Name)
		element.Field(i).Set(saved)
	}
}

func RunJWTTestCase(user model.User, t *testing.T, field string) {
	token, err := logic.JWT_generator(user, "52")
	if err != nil {
		fmt.Printf("t: %v\n", t)
		fmt.Printf("Tested %s as empty field\n", field)
	}
	if token != "" {
		t.Errorf("JWT getter for user %s had failed", user.Name)
	}
}

func TestJWT_emptyStruct(t *testing.T) {
	_, err := logic.JWT_generator(model.User{}, "52")
	if err != nil {
		assert.Equal(t, err.Error(), "Err type invalid user. Message: empty fields")
	} else {
		t.Error("Error. JWT generation didn`t failed with empty struct")
	}
}

func BenchmarkJWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		token, err := logic.JWT_generator(model.User{Name: "maxim", Email: "bruh@ya.ru"}, "52")
		if err == nil {
			fmt.Printf("token: %v\n", token)
			b.Logf("token: %v", token)
		}
	}
}
