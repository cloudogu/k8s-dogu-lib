package client

import (
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
	"testing"
)

func TestNewForConfigV3beta1(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		config := &rest.Config{}

		// when
		clientSet, err := NewForConfigV3beta1(config)

		// then
		require.NoError(t, err)
		require.NotNil(t, clientSet)
	})

}

func TestEcoSystemV3beta1Client_Dogus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		config := &rest.Config{}
		clientSet, err := NewForConfigV3beta1(config)
		require.NoError(t, err)
		require.NotNil(t, clientSet)

		// when
		client := clientSet.Dogus("ecosystem")

		// then
		require.NotNil(t, client)
	})
}
