package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultOccurrences(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/dfOcr.json": "dfOcr.json"})
	items, err := c.DefaultOccurrences(context.Background(), MaterialParams{CorpCode: "00126089", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "DH오토넥스", items[0].CorpName)
	assert.Equal(t, "당사 김제지점 발행 만기어음 부도", items[0].DfCn)
	assert.Equal(t, "48,322,175", items[0].DfAmt)
}

func TestBusinessSuspensions(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/bsnSp.json": "bsnSp.json"})
	items, err := c.BusinessSuspensions(context.Background(), MaterialParams{CorpCode: "00153393", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "방적 사업", items[0].BsnspRm)
	assert.Equal(t, "97,254,982,693", items[0].BsnspAmt)
	assert.Equal(t, "2023년 08월 31일", items[0].Bsnspd)
}

func TestRehabilitationApplications(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/ctrcvsBgrq.json": "ctrcvsBgrq.json"})
	items, err := c.RehabilitationApplications(context.Background(), MaterialParams{CorpCode: "00126089", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "서울회생법원", items[0].Cpct)
	assert.Equal(t, "주식회사 대유플러스", items[0].Apcnt)
}

func TestDissolutionCauses(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/dsRsOcr.json": "dsRsOcr.json"})
	items, err := c.DissolutionCauses(context.Background(), MaterialParams{CorpCode: "00580603", BgnDe: "20200101", EndDe: "20201231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "존립기간의 만료", items[0].DsRs)
	assert.Equal(t, "2020년 03월 27일", items[0].DsRsd)
}
