package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type listItem struct {
	Name string `json:"name"`
}

func TestGetList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "KEY", r.URL.Query().Get("crtfc_key"))
		assert.Equal(t, "x", r.URL.Query().Get("p"))
		w.Write([]byte(`{"status":"000","message":"정상","list":[{"name":"가"},{"name":"나"}]}`))
	}))
	t.Cleanup(srv.Close)
	c := New(Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()})

	items, err := GetList[listItem](context.Background(), c, "/api/x.json", map[string]string{"p": "x"})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, "가", items[0].Name)
}

func TestGetList_NoData(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"013","message":"조회된 데이타가 없습니다."}`))
	}))
	t.Cleanup(srv.Close)
	c := New(Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()})

	_, err := GetList[listItem](context.Background(), c, "/api/x.json", nil)
	assert.ErrorIs(t, err, ErrNoData)
}
