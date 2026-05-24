# OpenDART DS003 XBRL Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS003 정기보고서 재무정보의 XBRL 2개 API(택사노미 양식 + 재무제표 원본파일 ZIP)를 `report` 패키지에 추가한다.

**Architecture:** 기존 헬퍼만 재사용한다 — 택사노미는 JSON list 이므로 `getListParams[T]`, 원본파일은 바이너리 ZIP 이므로 `httpclient.GetBytes`(disclosure.DownloadDocument 와 동일 패턴). 새 파일 `report/xbrl.go`. 새 추상화·root 변경 없음(`client.Report` 기존).

**Tech Stack:** Go 1.25+, 표준 net/http (internal/httpclient 재사용), testify.

**Spec:** `docs/superpowers/specs/2026-05-24-opendart-ds003-xbrl-design.md`

**검증된 사실 (실 API):** `xbrlTaxonomy.json?sj_div=BS1` → status 000, 519행(필드 sj_div/bsns_de/account_id/account_nm/label_kor/label_eng/data_tp/ifrs_ref; abstract 행엔 data_tp 없음). `fnlttXbrl.xml?rcept_no=20240312000736&reprt_code=11011` → ~990KB ZIP(PK 매직, 바이너리).

**기존 재사용 심볼:** `Client`, `ReportCode`/`AnnualReport`, `getListParams[T]`(report/client.go), `httpclient.Client.GetBytes`/`httpclient.APIError`, `report/client_test.go` 의 `newTestClient`.

---

## File Structure

```
report/
  xbrl.go        # TaxonomyItem + XbrlTaxonomy + DownloadXbrl (신규)
  xbrl_test.go   # (신규)
  testdata/      # xbrlTaxonomy_bs1.json fixture 추가
README.md        # (수정) DS003 커버리지에 XBRL
integration_test.go  # (수정) XbrlTaxonomy 통합 케이스
```

---

### Task 1: XbrlTaxonomy + DownloadXbrl (report/xbrl.go)

**Files:**
- Create: `report/xbrl.go`, `report/xbrl_test.go`
- Create: `report/testdata/xbrlTaxonomy_bs1.json`

- [ ] **Step 1: fixture 작성** — `report/testdata/xbrlTaxonomy_bs1.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "sj_div": "BS1",
            "bsns_de": "20180701",
            "account_id": "ifrs_CurrentAssets",
            "account_nm": "CurrentAssets",
            "label_kor": "유동자산",
            "label_eng": "Current assets",
            "data_tp": "X",
            "ifrs_ref": "K-IFRS 1001 문단 60"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 작성** — `report/xbrl_test.go` (택사노미는 `newTestClient` 재사용, 다운로드는 자체 httptest):
```go
package report

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/opendart/internal/httpclient"
)

func TestXbrlTaxonomy(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/xbrlTaxonomy.json": "xbrlTaxonomy_bs1.json"})
	items, err := c.XbrlTaxonomy(context.Background(), "BS1")
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "BS1", items[0].SjDiv)
	assert.Equal(t, "ifrs_CurrentAssets", items[0].AccountId)
	assert.Equal(t, "유동자산", items[0].LabelKor)
	assert.Equal(t, "X", items[0].DataTp)
}

func TestDownloadXbrl_Binary(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "20240312000736", r.URL.Query().Get("rcept_no"))
		assert.Equal(t, "11011", r.URL.Query().Get("reprt_code"))
		w.Write([]byte("PK\x03\x04xbrlzip"))
	}))
	t.Cleanup(srv.Close)
	c := New(httpclient.New(httpclient.Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()}))

	b, err := c.DownloadXbrl(context.Background(), "20240312000736", AnnualReport)
	require.NoError(t, err)
	assert.Equal(t, "PK\x03\x04xbrlzip", string(b))
}

func TestDownloadXbrl_ErrorJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"010","message":"등록되지 않은 인증키입니다."}`))
	}))
	t.Cleanup(srv.Close)
	c := New(httpclient.New(httpclient.Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()}))

	_, err := c.DownloadXbrl(context.Background(), "x", AnnualReport)
	var apiErr *httpclient.APIError
	assert.ErrorAs(t, err, &apiErr)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./report/ -run 'TestXbrlTaxonomy|TestDownloadXbrl' -v`
Expected: FAIL — `undefined: ... XbrlTaxonomy`, `DownloadXbrl`.

- [ ] **Step 4: 구현** — `report/xbrl.go`:
```go
package report

import "context"

// TaxonomyItem 은 XBRL 택사노미 재무제표양식 (xbrlTaxonomy) 한 건.
type TaxonomyItem struct {
	SjDiv     string `json:"sj_div"`     // 재무제표구분
	AccountId string `json:"account_id"` // 계정ID
	AccountNm string `json:"account_nm"` // 계정명
	BsnsDe    string `json:"bsns_de"`    // 기준일 (YYYYMMDD)
	LabelKor  string `json:"label_kor"`  // 한글 출력명
	LabelEng  string `json:"label_eng"`  // 영문 출력명
	DataTp    string `json:"data_tp"`    // 데이터 유형 (일부 행에는 없음)
	IfrsRef   string `json:"ifrs_ref"`   // IFRS Reference
}

// XbrlTaxonomy 는 표준 XBRL 재무제표 택사노미 양식을 조회한다.
// sjDiv 는 재무제표구분 코드(BS1~4 재무상태표 / IS 손익계산서 / CIS 포괄손익 / CF 현금흐름표 /
// SCE 자본변동표 등). 잘못된 코드는 *opendart.APIError 로 반환된다.
func (c *Client) XbrlTaxonomy(ctx context.Context, sjDiv string) ([]TaxonomyItem, error) {
	return getListParams[TaxonomyItem](ctx, c.http, "/api/xbrlTaxonomy.json", map[string]string{"sj_div": sjDiv})
}

// DownloadXbrl 은 접수번호(rceptNo)+보고서코드로 재무제표 원본 XBRL(ZIP) 을 그대로 반환한다.
// 압축 해제·파싱은 호출자 몫.
func (c *Client) DownloadXbrl(ctx context.Context, rceptNo string, reprtCode ReportCode) ([]byte, error) {
	return c.http.GetBytes(ctx, "/api/fnlttXbrl.xml", map[string]string{
		"rcept_no":   rceptNo,
		"reprt_code": string(reprtCode),
	})
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./report/ -v`
Expected: 전체 PASS (기존 36 + 신규 3, 회귀 없음). `go vet ./report/` clean, `gofmt -l report/xbrl.go report/xbrl_test.go` no output.

- [ ] **Step 6: Commit**

```bash
git add report/xbrl.go report/xbrl_test.go report/testdata/
git commit -m "feat(report): XBRL 택사노미 + 재무제표 원본파일 다운로드"
```

---

### Task 2: README 커버리지 · 통합 테스트 · 최종 검증

**Files:**
- Modify: `README.md`
- Modify: `integration_test.go`

- [ ] **Step 1: README 커버리지 갱신** — `README.md` 의 DS003 줄을 다음으로 교체:
```markdown
- DS003 정기보고서 재무정보: 단일/다중회사 주요계정 · 단일회사 전체 재무제표(개별/연결) · 단일/다중회사 주요 재무지표 · XBRL 택사노미 양식 · 재무제표 원본파일(XBRL)
```
그리고 `- (예정)` 줄을 다음으로 교체(DS003 XBRL 제거):
```markdown
- (예정) DS002 개인별 보수 Ver2.0 2종 · DS004~DS006
```

- [ ] **Step 2: 통합 테스트 추가** — `integration_test.go` 에 함수 추가 (기존 `//go:build integration` · `report` import 유지):
```go
func TestIntegration_XbrlTaxonomy(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	items, err := c.Report.XbrlTaxonomy(context.Background(), "BS1")
	require.NoError(t, err)
	require.NotEmpty(t, items)
}
```
> 주: 이 테스트는 corp_code 매핑이 필요 없지만, 다른 통합 테스트와 동일하게 `NewClientFromEnv` 로
> 클라이언트를 만든다. `report` import 는 기존 통합 테스트에서 이미 존재한다(추가 import 불필요).

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
git commit -m "docs(report): DS003 XBRL 커버리지 + 통합 테스트"
```

---

## Self-Review Notes

- **Spec coverage:** TaxonomyItem + XbrlTaxonomy(getListParams) + DownloadXbrl(GetBytes) = Task1 · 테스트(택사노미 fixture + 다운로드 바이너리/에러) = Task1 · README/통합 = Task2. root 와이어링 불필요(client.Report 기존). 모두 매핑됨.
- **Type consistency:** `TaxonomyItem`(8필드, json 태그는 실 응답과 1:1) · `XbrlTaxonomy(ctx, sjDiv string) ([]TaxonomyItem, error)` · `DownloadXbrl(ctx, rceptNo string, reprtCode ReportCode) ([]byte, error)`. `getListParams[T]`/`GetBytes`/`AnnualReport`/`newTestClient` 는 기존 심볼 재사용.
- **검증된 fixture:** xbrlTaxonomy BS1 첫 데이터 행(data_tp="X" 포함) 실 응답. fnlttXbrl 은 바이너리라 httptest 로 모킹(실 ZIP 은 ~990KB 라 fixture 부적합).
- **새 추상화 없음:** 기존 getListParams/GetBytes 재사용. root 변경 없음.
- **이 그룹 완료 시 DS003 7/7 완료.** DS002 V2 2종만 보류로 남음.
