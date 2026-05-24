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
