package ownership

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMajorStockReports(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/majorstock.json": "majorstock.json"})
	items, err := c.MajorStockReports(context.Background(), "00126380")
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "삼성물산", items[0].Repror)
	assert.Equal(t, "1,199,285,813", items[0].Stkqy)
	assert.Equal(t, "20.09", items[0].Stkrt)
}

func TestExecutiveStockReports(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/elestock.json": "elestock.json"})
	items, err := c.ExecutiveStockReports(context.Background(), "00126380")
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "손준호", items[0].Repror)
	assert.Equal(t, "상무", items[0].IsuExctvOfcps)
	assert.Equal(t, "비등기임원", items[0].IsuExctvRgistAt)
}
