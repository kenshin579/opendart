package material

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// PaidInCapitalIncreaseItem 은 유상증자 결정 (piicDecsn) 한 건.
type PaidInCapitalIncreaseItem struct {
	RceptNo        string `json:"rcept_no"`         // 접수번호
	CorpCls        string `json:"corp_cls"`         // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`        // 고유번호
	CorpName       string `json:"corp_name"`        // 회사명
	NstkOstkCnt    string `json:"nstk_ostk_cnt"`    // 신주의 종류와 수(보통주식)
	NstkEstkCnt    string `json:"nstk_estk_cnt"`    // 신주의 종류와 수(기타주식)
	FvPs           string `json:"fv_ps"`            // 1주당 액면가액 (원)
	BficTisstkOstk string `json:"bfic_tisstk_ostk"` // 증자전 발행주식총수(보통주식)
	BficTisstkEstk string `json:"bfic_tisstk_estk"` // 증자전 발행주식총수(기타주식)
	FdppFclt       string `json:"fdpp_fclt"`        // 자금조달목적(시설자금)
	FdppBsninh     string `json:"fdpp_bsninh"`      // 자금조달목적(영업양수자금)
	FdppOp         string `json:"fdpp_op"`          // 자금조달목적(운영자금)
	FdppDtrp       string `json:"fdpp_dtrp"`        // 자금조달목적(채무상환자금)
	FdppOcsa       string `json:"fdpp_ocsa"`        // 자금조달목적(타법인 증권 취득자금)
	FdppEtc        string `json:"fdpp_etc"`         // 자금조달목적(기타자금)
	IcMthn         string `json:"ic_mthn"`          // 증자방식
	SslAt          string `json:"ssl_at"`           // 공매도 해당여부
	SslBgd         string `json:"ssl_bgd"`          // 공매도 시작일
	SslEdd         string `json:"ssl_edd"`          // 공매도 종료일
}

// PaidInCapitalIncrease 는 유상증자 결정을 조회한다.
func (c *Client) PaidInCapitalIncrease(ctx context.Context, p MaterialParams) ([]PaidInCapitalIncreaseItem, error) {
	return httpclient.GetList[PaidInCapitalIncreaseItem](ctx, c.http, "/api/piicDecsn.json", p.toMap())
}

// FreeCapitalIncreaseItem 은 무상증자 결정 (fricDecsn) 한 건.
type FreeCapitalIncreaseItem struct {
	RceptNo         string `json:"rcept_no"`           // 접수번호
	CorpCls         string `json:"corp_cls"`           // 법인구분 (Y/K/N/E)
	CorpCode        string `json:"corp_code"`          // 고유번호
	CorpName        string `json:"corp_name"`          // 회사명
	NstkOstkCnt     string `json:"nstk_ostk_cnt"`      // 신주의 종류와 수(보통주식)
	NstkEstkCnt     string `json:"nstk_estk_cnt"`      // 신주의 종류와 수(기타주식)
	FvPs            string `json:"fv_ps"`              // 1주당 액면가액 (원)
	BficTisstkOstk  string `json:"bfic_tisstk_ostk"`   // 증자전 발행주식총수(보통주식)
	BficTisstkEstk  string `json:"bfic_tisstk_estk"`   // 증자전 발행주식총수(기타주식)
	NstkAsstd       string `json:"nstk_asstd"`         // 신주배정기준일
	NstkAscntPsOstk string `json:"nstk_ascnt_ps_ostk"` // 1주당 신주배정 주식수(보통주식)
	NstkAscntPsEstk string `json:"nstk_ascnt_ps_estk"` // 1주당 신주배정 주식수(기타주식)
	NstkDividrk     string `json:"nstk_dividrk"`       // 신주의 배당기산일
	NstkDlprd       string `json:"nstk_dlprd"`         // 신주권교부예정일
	NstkLstprd      string `json:"nstk_lstprd"`        // 신주의 상장 예정일
	Bddd            string `json:"bddd"`               // 이사회결의일(결정일)
	OdAAtT          string `json:"od_a_at_t"`          // 사외이사 참석여부(참석)
	OdAAtB          string `json:"od_a_at_b"`          // 사외이사 참석여부(불참)
	AdtAAtn         string `json:"adt_a_atn"`          // 감사(감사위원) 참석여부
}

// FreeCapitalIncrease 는 무상증자 결정을 조회한다.
func (c *Client) FreeCapitalIncrease(ctx context.Context, p MaterialParams) ([]FreeCapitalIncreaseItem, error) {
	return httpclient.GetList[FreeCapitalIncreaseItem](ctx, c.http, "/api/fricDecsn.json", p.toMap())
}
