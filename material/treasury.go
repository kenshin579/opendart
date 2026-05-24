package material

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// TreasuryStockAcquisitionItem 은 자기주식 취득 결정 (tsstkAqDecsn) 한 건.
type TreasuryStockAcquisitionItem struct {
	RceptNo        string `json:"rcept_no"`           // 접수번호
	CorpCls        string `json:"corp_cls"`           // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`          // 고유번호
	CorpName       string `json:"corp_name"`          // 회사명
	AqplnStkOstk   string `json:"aqpln_stk_ostk"`     // 취득예정주식(주)(보통주식)
	AqplnStkEstk   string `json:"aqpln_stk_estk"`     // 취득예정주식(주)(기타주식)
	AqplnPrcOstk   string `json:"aqpln_prc_ostk"`     // 취득예정금액(원)(보통주식)
	AqplnPrcEstk   string `json:"aqpln_prc_estk"`     // 취득예정금액(원)(기타주식)
	AqexpdBgd      string `json:"aqexpd_bgd"`         // 취득예상기간(시작일)
	AqexpdEdd      string `json:"aqexpd_edd"`         // 취득예상기간(종료일)
	HdexpdBgd      string `json:"hdexpd_bgd"`         // 보유예상기간(시작일)
	HdexpdEdd      string `json:"hdexpd_edd"`         // 보유예상기간(종료일)
	AqPp           string `json:"aq_pp"`              // 취득목적
	AqMth          string `json:"aq_mth"`             // 취득방법
	CsIvBk         string `json:"cs_iv_bk"`           // 위탁투자중개업자
	AqWtnDivOstk   string `json:"aq_wtn_div_ostk"`    // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(보통주식))
	AqWtnDivOstkRt string `json:"aq_wtn_div_ostk_rt"` // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(비율%))
	AqWtnDivEstk   string `json:"aq_wtn_div_estk"`    // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(기타주식))
	AqWtnDivEstkRt string `json:"aq_wtn_div_estk_rt"` // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(비율%))
	EaqOstk        string `json:"eaq_ostk"`           // 취득 전 자기주식 보유현황(기타취득(주)(보통주식))
	EaqOstkRt      string `json:"eaq_ostk_rt"`        // 취득 전 자기주식 보유현황(기타취득(주)(비율%))
	EaqEstk        string `json:"eaq_estk"`           // 취득 전 자기주식 보유현황(기타취득(주)(기타주식))
	EaqEstkRt      string `json:"eaq_estk_rt"`        // 취득 전 자기주식 보유현황(기타취득(주)(비율%))
	AqDd           string `json:"aq_dd"`              // 취득결정일
	OdAAtT         string `json:"od_a_at_t"`          // 사외이사 참석여부(참석(명))
	OdAAtB         string `json:"od_a_at_b"`          // 사외이사 참석여부(불참(명))
	AdtAAtn        string `json:"adt_a_atn"`          // 감사(사외이사가 아닌 감사위원) 참석여부
	D1ProdlmOstk   string `json:"d1_prodlm_ostk"`     // 1일 매수 주문수량 한도(보통주식)
	D1ProdlmEstk   string `json:"d1_prodlm_estk"`     // 1일 매수 주문수량 한도(기타주식)
}

// TreasuryStockAcquisition 은 자기주식 취득 결정(주요사항보고서)을 조회한다.
func (c *Client) TreasuryStockAcquisition(ctx context.Context, p MaterialParams) ([]TreasuryStockAcquisitionItem, error) {
	return httpclient.GetList[TreasuryStockAcquisitionItem](ctx, c.http, "/api/tsstkAqDecsn.json", p.toMap())
}
