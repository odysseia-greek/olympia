package handlers

import (
	"encoding/json"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlexandros(t *testing.T) {
	scheme := "http"
	baseUrl := "somelocalhost.com"
	uuid := "thisisnotauuid"
	response := []models.Meros{{
		Greek:   "ὄνος",
		English: "an ass",
		Dutch:   "ezel",
	},
	}

	config := service.ClientConfig{
		Ca: nil,
		Alexandros: service.OdysseiaApi{
			Url:    baseUrl,
			Scheme: scheme,
			Cert:   nil,
		},
	}

	t.Run("Get", func(t *testing.T) {
		codes := []int{
			200,
		}

		r, err := json.Marshal(response)
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient, err := service.NewFakeClient(config, codes, responses)
		assert.Nil(t, err)

		handler := HomerosHandler{
			HttpClients: testClient,
		}

		sut, err := handler.Dictionary("word", "language", "mode", uuid)
		assert.Nil(t, err)
		assert.Equal(t, response, sut)
	})
}
