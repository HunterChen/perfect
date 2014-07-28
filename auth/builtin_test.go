package auth

import (
	"reflect"
	"testing"
)

func TestNewBuiltinStrategy(t *testing.T) {

	strategy := NewBuiltinStrategy(mock_auth_config)

	if !reflect.DeepEqual(strategy.Config, mock_auth_config) {
		t.Fatalf("actual builtin strategy config differs from expected configuration:\n actual: %#v\n expected: %#v\n", strategy.Config, mock_auth_config)
	}
}

//tests that NewBuiltinStrategyFunc returns a BuiltinStrategy, not some other strategy
func TestNewBuiltinStrategyFunc(t *testing.T) {
	var (
		auth_strategy Strategy
		ok            bool
	)

	auth_strategy = NewBuiltinStrategyFunc(mock_auth_config)

	_, ok = auth_strategy.(*BuiltinStrategy)

	if !ok {
		t.Fatalf("auth strategy is '%#v', expected *BuiltinStrategy", auth_strategy)
	}
}
