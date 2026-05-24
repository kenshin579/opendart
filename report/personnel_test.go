package report

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutives(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/exctvSttus.json": "exctvSttus.json"})
	items, err := c.Executives(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "한종희", items[0].Nm)
	assert.Equal(t, "부회장", items[0].Ofcps)
	assert.Equal(t, "사내이사", items[0].RgistExctvAt)
}

func TestEmployees(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/empSttus.json": "empSttus.json"})
	items, err := c.Employees(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "DX", items[0].FoBbm)
	assert.Equal(t, "37,962", items[0].RgllbrCo)
	assert.Equal(t, "38,286", items[0].Sm)
}

func TestUnregisteredExecutiveCompensation(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/unrstExctvMendngSttus.json": "unrstExctvMendngSttus.json"})
	items, err := c.UnregisteredExecutiveCompensation(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "미등기임원", items[0].Se)
	assert.Equal(t, "1,015", items[0].Nmpr)
}

func TestOutsideDirectorChanges(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/outcmpnyDrctrNdChangeSttus.json": "outcmpnyDrctrNdChangeSttus.json"})
	items, err := c.OutsideDirectorChanges(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "11", items[0].DrctrCo)
	assert.Equal(t, "6", items[0].OtcmpDrctrCo)
}

func TestDirectorAuditorApprovedCompensation(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/drctrAdtAllMendngSttusGmtsckConfmAmount.json": "drctrAdtAllMendngSttusGmtsckConfmAmount.json"})
	items, err := c.DirectorAuditorApprovedCompensation(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "등기이사", items[0].Se)
	assert.Equal(t, "5", items[0].Nmpr)
}

func TestDirectorAuditorTotalCompensation(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/hmvAuditAllSttus.json": "hmvAuditAllSttus.json"})
	items, err := c.DirectorAuditorTotalCompensation(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "11", items[0].Nmpr)
	assert.Equal(t, "23,227,000,000", items[0].MendngTotamt)
}

func TestDirectorAuditorCompensationByType(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/drctrAdtAllMendngSttusMendngPymntamtTyCl.json": "drctrAdtAllMendngSttusMendngPymntamtTyCl.json"})
	items, err := c.DirectorAuditorCompensationByType(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "22,009,000,000", items[0].PymntTotamt)
	assert.Equal(t, "4,402,000,000", items[0].Psn1AvrgPymntamt)
}

func TestIndividualDirectorAuditorCompensation(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/hmvAuditIndvdlBySttus.json": "hmvAuditIndvdlBySttus.json"})
	items, err := c.IndividualDirectorAuditorCompensation(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "경계현", items[0].Nm)
	assert.Equal(t, "2,403,000,000", items[0].MendngTotamt)
}

func TestIndividualTop5Compensation(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/indvdlByPay.json": "indvdlByPay.json"})
	items, err := c.IndividualTop5Compensation(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "김기남", items[0].Nm)
	assert.Equal(t, "고문", items[0].Ofcps)
	assert.Equal(t, "17,265,000,000", items[0].MendngTotamt)
}
