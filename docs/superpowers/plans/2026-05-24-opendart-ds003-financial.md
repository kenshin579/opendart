# OpenDART DS003 재무 핵심(JSON) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS003 정기보고서 재무정보의 JSON 재무 핵심 5개 API를 `report` 패키지에 추가한다 (XBRL 바이너리·택사노미 2개는 후속).

**Architecture:** 기존 `getList[T](ctx,hc,path,ReportParams)` 를 raw-map 코어 `getListParams[T](ctx,hc,path,map)` 로 리팩토링하고 `getList` 는 thin wrapper 로 유지(DS002 영향 없음). 추가 파라미터(fs_div, idx_cl_code)는 엔드포인트별 파라미터 struct + 타입 상수로 처리. 다중회사는 `corp_code` 콤마 구분. 5개 메서드가 3개 item struct(`AccountItem`/`FinancialIndexItem`/`FullStatementItem`)를 공유한다. 새 파일 `report/financial.go`. `client.Report` 는 이미 root 에 와이어링됨 → root 변경 없음.

**Tech Stack:** Go 1.25+ (제네릭), 표준 net/http (internal/httpclient 재사용), testify.

**Spec:** `docs/superpowers/specs/2026-05-24-opendart-ds003-financial-design.md`

**검증된 사실 (실 API, 삼성전자 00126380 / SK하이닉스 00164779 / 2023 / 11011):** 5개 모두 status 000 + list (단일 주요계정 30행, 다중 60, 전체재무제표 OFS 115, 단일 재무지표 수익성 15, 다중 재무지표 30). 단일/다중 주요계정 동일 키(corp_code 포함), 단일/다중 재무지표 동일 키. **주요계정 응답엔 docs 표에 없던 `corp_code` 가 실제로 존재**(AccountItem 에 포함). 주요계정 금액은 콤마("195,936,557,000,000"), 전체재무제표 금액은 콤마 없음("296857289000000") — 둘 다 string. add_amount 필드는 BS 항목엔 없고 IS/분기 항목에 존재(스키마엔 있음, struct 에 유지).

**기존 재사용 심볼:** `Client`, `ReportParams`, `ReportCode`/`AnnualReport`, `listResponse[T]`, `report/client_test.go` 의 `newTestClient`.

---

## File Structure

```
report/
  client.go         # (수정) getList → getListParams[T] 코어 + thin wrapper
  financial.go      # 5개 메서드 + 3 item struct + 2 파라미터 struct + 타입 상수 (신규)
  financial_test.go # (신규)
  testdata/         # 5개 실 응답 JSON fixture 추가
README.md           # (수정) DS003 커버리지
integration_test.go # (수정) SingleAccount 통합 케이스
```

---

### Task 1: getList → getListParams 리팩토링

**Files:**
- Modify: `report/client.go`

- [ ] **Step 1: `getList` 함수를 다음으로 교체** — `report/client.go` 의 기존 `getList` 정의:
```go
// getList 는 공통 list 조회 헬퍼. GetJSON 의 status 검사를 거친 뒤 list 만 반환한다.
// 조회 데이터 없음(013)은 httpclient 가 ErrNoData 로 변환한다.
func getList[T any](ctx context.Context, hc *httpclient.Client, path string, p ReportParams) ([]T, error) {
	var resp listResponse[T]
	if err := hc.GetJSON(ctx, path, p.toMap(), &resp); err != nil {
		return nil, err
	}
	return resp.List, nil
}
```
를 다음으로 교체한다:
```go
// getListParams 는 raw 파라미터 맵으로 list 를 조회하는 코어 헬퍼.
// GetJSON 의 status 검사를 거친 뒤 list 만 반환한다(013은 httpclient 가 ErrNoData 로 변환).
func getListParams[T any](ctx context.Context, hc *httpclient.Client, path string, params map[string]string) ([]T, error) {
	var resp listResponse[T]
	if err := hc.GetJSON(ctx, path, params, &resp); err != nil {
		return nil, err
	}
	return resp.List, nil
}

// getList 는 ReportParams 기반 thin wrapper.
func getList[T any](ctx context.Context, hc *httpclient.Client, path string, p ReportParams) ([]T, error) {
	return getListParams[T](ctx, hc, path, p.toMap())
}
```

- [ ] **Step 2: 회귀 없음 확인 (기존 DS002 메서드가 getList 사용)**

Run: `go test ./report/ -v`
Expected: 기존 30개 테스트 모두 PASS (동작 동일). `go vet ./report/` clean.

- [ ] **Step 3: Commit**

```bash
git add report/client.go
git commit -m "refactor(report): extract getListParams raw-map core from getList"
```

---

### Task 2: 단일/다중 주요계정 (SingleAccount / MultiAccount)

**Files:**
- Create: `report/financial.go`, `report/financial_test.go`
- Create: `report/testdata/fnlttSinglAcnt.json`, `report/testdata/fnlttMultiAcnt.json`

- [ ] **Step 1: fixture 작성**

`report/testdata/fnlttSinglAcnt.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "reprt_code": "11011",
            "bsns_year": "2023",
            "corp_code": "00126380",
            "stock_code": "005930",
            "fs_div": "CFS",
            "fs_nm": "연결재무제표",
            "sj_div": "BS",
            "sj_nm": "재무상태표",
            "account_nm": "유동자산",
            "thstrm_nm": "제 55 기",
            "thstrm_dt": "2023.12.31 현재",
            "thstrm_amount": "195,936,557,000,000",
            "frmtrm_nm": "제 54 기",
            "frmtrm_dt": "2022.12.31 현재",
            "frmtrm_amount": "218,470,581,000,000",
            "bfefrmtrm_nm": "제 53 기",
            "bfefrmtrm_dt": "2021.12.31 현재",
            "bfefrmtrm_amount": "218,163,185,000,000",
            "ord": "1",
            "currency": "KRW"
        }
    ]
}
```

`report/testdata/fnlttMultiAcnt.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "reprt_code": "11011",
            "bsns_year": "2023",
            "corp_code": "00126380",
            "stock_code": "005930",
            "fs_div": "CFS",
            "fs_nm": "연결재무제표",
            "sj_div": "BS",
            "sj_nm": "재무상태표",
            "account_nm": "유동자산",
            "thstrm_nm": "제 55 기",
            "thstrm_dt": "2023.12.31 현재",
            "thstrm_amount": "195,936,557,000,000",
            "frmtrm_nm": "제 54 기",
            "frmtrm_dt": "2022.12.31 현재",
            "frmtrm_amount": "218,470,581,000,000",
            "bfefrmtrm_nm": "제 53 기",
            "bfefrmtrm_dt": "2021.12.31 현재",
            "bfefrmtrm_amount": "218,163,185,000,000",
            "ord": "1",
            "currency": "KRW"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 작성** — `report/financial_test.go` (기존 `newTestClient` 재사용):
```go
package report

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSingleAccount(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fnlttSinglAcnt.json": "fnlttSinglAcnt.json"})
	items, err := c.SingleAccount(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "유동자산", items[0].AccountNm)
	assert.Equal(t, "연결재무제표", items[0].FsNm)
	assert.Equal(t, "195,936,557,000,000", items[0].ThstrmAmount)
	assert.Equal(t, "00126380", items[0].CorpCode)
}

func TestMultiAccount(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fnlttMultiAcnt.json": "fnlttMultiAcnt.json"})
	items, err := c.MultiAccount(context.Background(), ReportParams{CorpCode: "00126380,00164779", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "유동자산", items[0].AccountNm)
	assert.Equal(t, "00126380", items[0].CorpCode)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./report/ -run 'TestSingleAccount|TestMultiAccount' -v`
Expected: FAIL — `undefined: ... SingleAccount` 등.

- [ ] **Step 4: 구현** — `report/financial.go`:
```go
package report

import "context"

// AccountItem 은 단일/다중회사 주요계정 (fnlttSinglAcnt / fnlttMultiAcnt) 한 건.
type AccountItem struct {
	RceptNo         string `json:"rcept_no"`          // 접수번호
	ReprtCode       string `json:"reprt_code"`        // 보고서 코드
	BsnsYear        string `json:"bsns_year"`         // 사업 연도
	CorpCode        string `json:"corp_code"`         // 고유번호
	StockCode       string `json:"stock_code"`        // 종목 코드
	FsDiv           string `json:"fs_div"`            // 개별/연결구분
	FsNm            string `json:"fs_nm"`             // 개별/연결명
	SjDiv           string `json:"sj_div"`            // 재무제표구분
	SjNm            string `json:"sj_nm"`             // 재무제표명
	AccountNm       string `json:"account_nm"`        // 계정명
	ThstrmNm        string `json:"thstrm_nm"`         // 당기명
	ThstrmDt        string `json:"thstrm_dt"`         // 당기일자
	ThstrmAmount    string `json:"thstrm_amount"`     // 당기금액
	ThstrmAddAmount string `json:"thstrm_add_amount"` // 당기누적금액
	FrmtrmNm        string `json:"frmtrm_nm"`         // 전기명
	FrmtrmDt        string `json:"frmtrm_dt"`         // 전기일자
	FrmtrmAmount    string `json:"frmtrm_amount"`     // 전기금액
	FrmtrmAddAmount string `json:"frmtrm_add_amount"` // 전기누적금액
	BfefrmtrmNm     string `json:"bfefrmtrm_nm"`      // 전전기명
	BfefrmtrmDt     string `json:"bfefrmtrm_dt"`      // 전전기일자
	BfefrmtrmAmount string `json:"bfefrmtrm_amount"`  // 전전기금액
	Ord             string `json:"ord"`               // 계정과목 정렬순서
	Currency        string `json:"currency"`          // 통화 단위
}

// SingleAccount 는 단일회사 주요계정을 조회한다.
func (c *Client) SingleAccount(ctx context.Context, p ReportParams) ([]AccountItem, error) {
	return getList[AccountItem](ctx, c.http, "/api/fnlttSinglAcnt.json", p)
}

// MultiAccount 는 다중회사 주요계정을 조회한다. p.CorpCode 는 콤마로 여러 고유번호를 전달한다.
func (c *Client) MultiAccount(ctx context.Context, p ReportParams) ([]AccountItem, error) {
	return getList[AccountItem](ctx, c.http, "/api/fnlttMultiAcnt.json", p)
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./report/ -v`
Expected: 전체 PASS (기존 30 + 신규 2, 회귀 없음). `go vet ./report/` clean.

- [ ] **Step 6: Commit**

```bash
git add report/financial.go report/financial_test.go report/testdata/
git commit -m "feat(report): 단일·다중회사 주요계정"
```

---

### Task 3: 전체 재무제표 + 단일/다중 재무지표 (파라미터 struct + 상수)

**Files:**
- Modify: `report/financial.go` (3개 메서드 + 2 struct + 2 파라미터 struct + 상수)
- Modify: `report/financial_test.go` (3개 메서드 테스트 + toMap 테스트)
- Create: `report/testdata/fnlttSinglAcntAll.json`, `report/testdata/fnlttSinglIndx.json`, `report/testdata/fnlttCmpnyIndx.json`

- [ ] **Step 1: fixture 작성**

`report/testdata/fnlttSinglAcntAll.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "reprt_code": "11011",
            "bsns_year": "2023",
            "corp_code": "00126380",
            "sj_div": "BS",
            "sj_nm": "재무상태표",
            "account_id": "ifrs-full_Assets",
            "account_nm": "자산총계",
            "account_detail": "-",
            "thstrm_nm": "제 55 기",
            "thstrm_amount": "296857289000000",
            "frmtrm_nm": "제 54 기",
            "frmtrm_amount": "260083750000000",
            "bfefrmtrm_nm": "제 53 기",
            "bfefrmtrm_amount": "251112184000000",
            "ord": "7",
            "currency": "KRW"
        }
    ]
}
```

`report/testdata/fnlttSinglIndx.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "reprt_code": "11011",
            "bsns_year": "2023",
            "corp_code": "00126380",
            "stock_code": "005930",
            "stlm_dt": "2023-12-31",
            "idx_cl_code": "M210000",
            "idx_cl_nm": "수익성지표",
            "idx_code": "M211100",
            "idx_nm": "세전계속사업이익률",
            "idx_val": "5.981"
        }
    ]
}
```

`report/testdata/fnlttCmpnyIndx.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "reprt_code": "11011",
            "bsns_year": "2023",
            "corp_code": "00126380",
            "stock_code": "005930",
            "stlm_dt": "2023-12-31",
            "idx_cl_code": "M210000",
            "idx_cl_nm": "수익성지표",
            "idx_code": "M211100",
            "idx_nm": "세전계속사업이익률",
            "idx_val": "5.981"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 추가** — `report/financial_test.go` 에 추가:
```go
func TestFsDivAndIndexParams_toMap(t *testing.T) {
	sm := FinancialStatementParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport, FsDiv: FsDivSeparate}.toMap()
	assert.Equal(t, "OFS", sm["fs_div"])
	assert.Equal(t, "11011", sm["reprt_code"])

	im := FinancialIndexParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport, IdxClCode: IndexProfitability}.toMap()
	assert.Equal(t, "M210000", im["idx_cl_code"])
	assert.Equal(t, "00126380", im["corp_code"])
}

func TestSingleFullStatement(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fnlttSinglAcntAll.json": "fnlttSinglAcntAll.json"})
	items, err := c.SingleFullStatement(context.Background(), FinancialStatementParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport, FsDiv: FsDivSeparate})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "ifrs-full_Assets", items[0].AccountId)
	assert.Equal(t, "자산총계", items[0].AccountNm)
	assert.Equal(t, "296857289000000", items[0].ThstrmAmount)
}

func TestSingleIndex(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fnlttSinglIndx.json": "fnlttSinglIndx.json"})
	items, err := c.SingleIndex(context.Background(), FinancialIndexParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport, IdxClCode: IndexProfitability})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "수익성지표", items[0].IdxClNm)
	assert.Equal(t, "세전계속사업이익률", items[0].IdxNm)
	assert.Equal(t, "5.981", items[0].IdxVal)
}

func TestMultiIndex(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fnlttCmpnyIndx.json": "fnlttCmpnyIndx.json"})
	items, err := c.MultiIndex(context.Background(), FinancialIndexParams{CorpCode: "00126380,00164779", BsnsYear: "2023", ReprtCode: AnnualReport, IdxClCode: IndexProfitability})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "00126380", items[0].CorpCode)
	assert.Equal(t, "5.981", items[0].IdxVal)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./report/ -run 'TestFsDivAndIndexParams|TestSingleFullStatement|TestSingleIndex|TestMultiIndex' -v`
Expected: FAIL — `undefined: FinancialStatementParams` 등.

- [ ] **Step 4: 구현** — `report/financial.go` 에 추가:
```go
// FsDiv 는 개별/연결 구분.
type FsDiv string

const (
	FsDivSeparate     FsDiv = "OFS" // 재무제표(개별)
	FsDivConsolidated FsDiv = "CFS" // 연결재무제표
)

// IndexClass 는 재무지표 분류 코드.
type IndexClass string

const (
	IndexProfitability IndexClass = "M210000" // 수익성지표
	IndexStability     IndexClass = "M220000" // 안정성지표
	IndexGrowth        IndexClass = "M230000" // 성장성지표
	IndexActivity      IndexClass = "M240000" // 활동성지표
)

// FinancialStatementParams 는 전체 재무제표 요청 인자.
type FinancialStatementParams struct {
	CorpCode  string     // 고유번호 (8자리)
	BsnsYear  string     // 사업연도 (4자리, 2015 이후)
	ReprtCode ReportCode // 보고서 코드
	FsDiv     FsDiv      // 개별(OFS)/연결(CFS)
}

func (p FinancialStatementParams) toMap() map[string]string {
	return map[string]string{
		"corp_code":  p.CorpCode,
		"bsns_year":  p.BsnsYear,
		"reprt_code": string(p.ReprtCode),
		"fs_div":     string(p.FsDiv),
	}
}

// FinancialIndexParams 는 주요 재무지표 요청 인자. 다중회사는 CorpCode 를 콤마로 구분한다.
type FinancialIndexParams struct {
	CorpCode  string     // 고유번호 (8자리; 다중은 콤마 구분)
	BsnsYear  string     // 사업연도 (4자리)
	ReprtCode ReportCode // 보고서 코드
	IdxClCode IndexClass // 지표분류코드
}

func (p FinancialIndexParams) toMap() map[string]string {
	return map[string]string{
		"corp_code":   p.CorpCode,
		"bsns_year":   p.BsnsYear,
		"reprt_code":  string(p.ReprtCode),
		"idx_cl_code": string(p.IdxClCode),
	}
}

// FullStatementItem 은 단일회사 전체 재무제표 (fnlttSinglAcntAll) 한 건.
type FullStatementItem struct {
	RceptNo         string `json:"rcept_no"`          // 접수번호
	ReprtCode       string `json:"reprt_code"`        // 보고서 코드
	BsnsYear        string `json:"bsns_year"`         // 사업 연도
	CorpCode        string `json:"corp_code"`         // 고유번호
	SjDiv           string `json:"sj_div"`            // 재무제표구분
	SjNm            string `json:"sj_nm"`             // 재무제표명
	AccountId       string `json:"account_id"`        // 계정ID
	AccountNm       string `json:"account_nm"`        // 계정명
	AccountDetail   string `json:"account_detail"`    // 계정상세
	ThstrmNm        string `json:"thstrm_nm"`         // 당기명
	ThstrmAmount    string `json:"thstrm_amount"`     // 당기금액
	ThstrmAddAmount string `json:"thstrm_add_amount"` // 당기누적금액
	FrmtrmNm        string `json:"frmtrm_nm"`         // 전기명
	FrmtrmAmount    string `json:"frmtrm_amount"`     // 전기금액
	FrmtrmQNm       string `json:"frmtrm_q_nm"`       // 전기명(분/반기)
	FrmtrmQAmount   string `json:"frmtrm_q_amount"`   // 전기금액(분/반기)
	FrmtrmAddAmount string `json:"frmtrm_add_amount"` // 전기누적금액
	BfefrmtrmNm     string `json:"bfefrmtrm_nm"`      // 전전기명
	BfefrmtrmAmount string `json:"bfefrmtrm_amount"`  // 전전기금액
	Ord             string `json:"ord"`               // 계정과목 정렬순서
	Currency        string `json:"currency"`          // 통화 단위
}

// SingleFullStatement 는 단일회사 전체 재무제표를 조회한다.
func (c *Client) SingleFullStatement(ctx context.Context, p FinancialStatementParams) ([]FullStatementItem, error) {
	return getListParams[FullStatementItem](ctx, c.http, "/api/fnlttSinglAcntAll.json", p.toMap())
}

// FinancialIndexItem 은 단일/다중회사 주요 재무지표 (fnlttSinglIndx / fnlttCmpnyIndx) 한 건.
type FinancialIndexItem struct {
	ReprtCode string `json:"reprt_code"`  // 보고서 코드
	BsnsYear  string `json:"bsns_year"`   // 사업 연도
	CorpCode  string `json:"corp_code"`   // 고유번호
	StockCode string `json:"stock_code"`  // 종목 코드
	StlmDt    string `json:"stlm_dt"`     // 결산기준일
	IdxClCode string `json:"idx_cl_code"` // 지표분류코드
	IdxClNm   string `json:"idx_cl_nm"`   // 지표분류명
	IdxCode   string `json:"idx_code"`    // 지표코드
	IdxNm     string `json:"idx_nm"`      // 지표명
	IdxVal    string `json:"idx_val"`     // 지표값
}

// SingleIndex 는 단일회사 주요 재무지표를 조회한다.
func (c *Client) SingleIndex(ctx context.Context, p FinancialIndexParams) ([]FinancialIndexItem, error) {
	return getListParams[FinancialIndexItem](ctx, c.http, "/api/fnlttSinglIndx.json", p.toMap())
}

// MultiIndex 는 다중회사 주요 재무지표를 조회한다. p.CorpCode 는 콤마로 여러 고유번호를 전달한다.
func (c *Client) MultiIndex(ctx context.Context, p FinancialIndexParams) ([]FinancialIndexItem, error) {
	return getListParams[FinancialIndexItem](ctx, c.http, "/api/fnlttCmpnyIndx.json", p.toMap())
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./report/ -v`
Expected: 전체 PASS (Task 2 포함 신규 5개 메서드 + toMap + 기존). `go vet ./report/` clean.

- [ ] **Step 6: Commit**

```bash
git add report/financial.go report/financial_test.go report/testdata/
git commit -m "feat(report): 전체 재무제표 + 단일·다중 재무지표 (fs_div/idx_cl_code 파라미터)"
```

---

### Task 4: README 커버리지 · 통합 테스트 · 최종 검증

**Files:**
- Modify: `README.md`
- Modify: `integration_test.go`

- [ ] **Step 1: README 커버리지 갱신** — `README.md` 의 `- (예정)` 줄 **앞에** DS003 줄을 추가한다 (DS002 줄은 그대로):
```markdown
- DS003 정기보고서 재무정보: 단일/다중회사 주요계정 · 단일회사 전체 재무제표(개별/연결) · 단일/다중회사 주요 재무지표
```
그리고 `- (예정)` 줄을 다음으로 교체:
```markdown
- (예정) DS002 개인별 보수 Ver2.0 2종 · DS003 XBRL 원본파일·택사노미 · DS004~DS006
```

- [ ] **Step 2: 통합 테스트 추가** — `integration_test.go` 에 함수 추가 (기존 `//go:build integration` · `report` import 유지):
```go
func TestIntegration_SingleAccount(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)

	items, err := c.Report.SingleAccount(context.Background(), report.ReportParams{
		CorpCode:  corp,
		BsnsYear:  "2023",
		ReprtCode: report.AnnualReport,
	})
	require.NoError(t, err)
	require.NotEmpty(t, items)
}
```

- [ ] **Step 3: 통합 빌드 확인 (기본 빌드 제외)**

Run: `go vet -tags integration ./...`
Expected: clean.
Run: `go test ./...`
Expected: 전체 PASS, integration 미실행.

- [ ] **Step 4: 최종 전체 검증**

Run:
```bash
go build ./...
go vet ./...
go test ./...
gofmt -l . | grep -v '^scripts/crawl' || echo "clean"
```
Expected: build/vet 성공, 전체 PASS, gofmt 신규 파일 차이 없음("clean").

- [ ] **Step 5: Commit**

```bash
git add README.md integration_test.go
git commit -m "docs(report): DS003 재무정보 커버리지 + 통합 테스트"
```

---

## Self-Review Notes

- **Spec coverage:** getList 리팩토링 = Task1 · 5개 메서드 = Task2(2)+Task3(3) · 3 struct(AccountItem/FinancialIndexItem/FullStatementItem) + 2 파라미터 struct + FsDiv/IndexClass 상수 = Task2·3 · 테스트(fixture+toMap) = Task2·3 · README/통합 = Task4. root 와이어링 불필요(client.Report 기존). XBRL 2개 비범위(spec 명시). 모두 매핑됨.
- **Type consistency:** `getListParams[T]`/`getList[T]`(wrapper) · `AccountItem`/`FinancialIndexItem`/`FullStatementItem` · `FinancialStatementParams{...,FsDiv}`/`FinancialIndexParams{...,IdxClCode}`(+toMap) · `FsDiv`(FsDivSeparate/FsDivConsolidated)/`IndexClass`(IndexProfitability/Stability/Growth/Activity) · 메서드 `SingleAccount`/`MultiAccount`(ReportParams) · `SingleFullStatement`(FinancialStatementParams) · `SingleIndex`/`MultiIndex`(FinancialIndexParams). 시그니처 일관. 필드·json 태그는 캡처한 실 응답과 1:1.
- **spec 보정:** AccountItem 에 `CorpCode` 추가(docs 표엔 없으나 실 응답에 존재). 전체재무제표 금액은 콤마 없음(string 유지). add_amount 필드는 스키마에 있어 struct 에 유지(BS 항목엔 빈 값).
- **검증된 fixture:** 5개 모두 실 API(삼성/SK하이닉스, 2023, 사업보고서) 응답 첫 항목.
- **getList 리팩토링은 동작 보존:** 기존 DS002 메서드는 getList(ReportParams) wrapper 를 그대로 호출 → 회귀 없음(Task1 Step2 로 검증).
