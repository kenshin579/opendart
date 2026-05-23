package report

import "context"

// DividendItem 은 배당에 관한 사항 (alotMatter) 한 건.
type DividendItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호 (14자리)
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호 (8자리)
	CorpName string `json:"corp_name"` // 법인명
	Se       string `json:"se"`        // 구분 (주당액면가액, 주당 현금배당금 등)
	StockKnd string `json:"stock_knd"` // 주식 종류 (보통주 등)
	Thstrm   string `json:"thstrm"`    // 당기
	Frmtrm   string `json:"frmtrm"`    // 전기
	Lwfr     string `json:"lwfr"`      // 전전기
	StlmDt   string `json:"stlm_dt"`   // 결산기준일 (YYYY-MM-DD)
}

// Dividend 는 배당에 관한 사항을 조회한다.
func (c *Client) Dividend(ctx context.Context, p ReportParams) ([]DividendItem, error) {
	return getList[DividendItem](ctx, c.http, "/api/alotMatter.json", p)
}
