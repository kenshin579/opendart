package material

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// OtherAssetTransferPutbackOptionItem 은 자산양수도(기타), 풋백옵션 (astInhtrfEtcPtbkOpt) 한 건.
type OtherAssetTransferPutbackOptionItem struct {
	RceptNo      string `json:"rcept_no"`       // 접수번호
	CorpCls      string `json:"corp_cls"`       // 법인구분 (Y/K/N/E)
	CorpCode     string `json:"corp_code"`      // 고유번호
	CorpName     string `json:"corp_name"`      // 회사명
	RpRsn        string `json:"rp_rsn"`         // 보고 사유
	AstInhtrfPrc string `json:"ast_inhtrf_prc"` // 자산양수ㆍ도 가액
}

// OtherAssetTransferPutbackOption 은 자산양수도(기타), 풋백옵션(주요사항보고서)을 조회한다.
func (c *Client) OtherAssetTransferPutbackOption(ctx context.Context, p MaterialParams) ([]OtherAssetTransferPutbackOptionItem, error) {
	return httpclient.GetList[OtherAssetTransferPutbackOptionItem](ctx, c.http, "/api/astInhtrfEtcPtbkOpt.json", p.toMap())
}
