package repository

// Kodek `JSONDATA` - pemetaan SATU-SATUNYA antara struct bernama (models) dan
// kunci JSON Pega. Bukti setiap kunci: `docs/PARITAS-LAYAR-DAN-AKSI.md` §3–§5.
//
// ⛔ Kunci peka huruf besar-kecil, ejaan Pega dipertahankan (`POLICYHODER`,
// `Non_Employee`, `UnderwritingLimitList`) - ketiga view membacanya persis.
// ⛔ Nilai skalar Pega adalah TEKS. Pembaca menerima teks, angka JSON, boolean,
// dan null (data lama); angka tidak pernah melewati float (`json.Number`).
// ⛔ Kunci yang tidak dikelola layar DIPERTAHANKAN: di tingkat halaman lewat
// JSON lama yang dibaca ulang saat menulis, di tingkat baris lewat `Asli`.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"nusantarare/modul/masterproductnamelife/backend/models"
)

// ErrJSONRusak - `JSONDATA` tidak dapat diurai (constraint `IS JSON` DEV
// mencegahnya; bila tetap terjadi, gagal terang, bukan produk kosong).
var ErrJSONRusak = errors.New("repository: product JSONDATA cannot be read")

// medanTeks - satu kunci JSON skalar ↔ satu medan teks struct.
type medanTeks[T any] struct {
	kunci string
	ambil func(*T) *string
	// tanggal - teks `dd/MM/yyyy` di JSON, `YYYY-MM-DD` di model.
	tanggal bool
}

// Kunci skalar halaman `ProductName` (`M_PRODUCT_LIFE.JSONDATA`).
var medanUmum = []medanTeks[models.ProdukUmum]{
	{kunci: "PRODUCTNAME", ambil: func(u *models.ProdukUmum) *string { return &u.ProductName }},
	{kunci: "CEDING", ambil: func(u *models.ProdukUmum) *string { return &u.Ceding }},
	{kunci: "CEDINGID", ambil: func(u *models.ProdukUmum) *string { return &u.CedingID }},
	{kunci: "SOBNAME", ambil: func(u *models.ProdukUmum) *string { return &u.SOBName }},
	{kunci: "SOBID", ambil: func(u *models.ProdukUmum) *string { return &u.SOBID }},
	{kunci: "RICOMM", ambil: func(u *models.ProdukUmum) *string { return &u.RIComm }},
	{kunci: "RIRISK", ambil: func(u *models.ProdukUmum) *string { return &u.RIRisk }},
	{kunci: "RIRISKID", ambil: func(u *models.ProdukUmum) *string { return &u.RIRiskID }},
	{kunci: "INWARDNAME", ambil: func(u *models.ProdukUmum) *string { return &u.InwardName }},
	{kunci: "TREATYNUMBER", ambil: func(u *models.ProdukUmum) *string { return &u.TreatyNumber }},
	{kunci: "CAUSE", ambil: func(u *models.ProdukUmum) *string { return &u.Cause }},
	{kunci: "CAUSEID", ambil: func(u *models.ProdukUmum) *string { return &u.CauseID }},
	{kunci: "POLICYHODER", ambil: func(u *models.ProdukUmum) *string { return &u.PolicyHolder }},
	{kunci: "POLICYHODERNAME", ambil: func(u *models.ProdukUmum) *string { return &u.PolicyHolderName }},
	{kunci: "CREATEOP", ambil: func(u *models.ProdukUmum) *string { return &u.CreateOp }},
	{kunci: "UPDATEOP", ambil: func(u *models.ProdukUmum) *string { return &u.UpdateOp }},
	{kunci: "Comment", ambil: func(u *models.ProdukUmum) *string { return &u.Comment }},
	// Medan layar mati - kunci view `PRODUCT_LIFE`.
	{kunci: "TYPE", ambil: func(u *models.ProdukUmum) *string { return &u.TypeBasicRider }},
	{kunci: "TYPE_CEDING", ambil: func(u *models.ProdukUmum) *string { return &u.TypeCeding }},
	{kunci: "GRUP", ambil: func(u *models.ProdukUmum) *string { return &u.Grup }},
	{kunci: "PRODUCTCODE", ambil: func(u *models.ProdukUmum) *string { return &u.ProductCode }},
	{kunci: "PRODUCTTYPE", ambil: func(u *models.ProdukUmum) *string { return &u.ProductType }},
	{kunci: "PRODUCTTYPEID", ambil: func(u *models.ProdukUmum) *string { return &u.ProductTypeID }},
	{kunci: "RIRATE", ambil: func(u *models.ProdukUmum) *string { return &u.RIRate }},
	{kunci: "RIRATEID", ambil: func(u *models.ProdukUmum) *string { return &u.RIRateID }},
	{kunci: "RICOMMID", ambil: func(u *models.ProdukUmum) *string { return &u.RICommID }},
	{kunci: "OUTWARDNAME", ambil: func(u *models.ProdukUmum) *string { return &u.OutwardName }},
	{kunci: "OUTWARDNAMEID", ambil: func(u *models.ProdukUmum) *string { return &u.OutwardNameID }},
	{kunci: "OUTWARDRATE", ambil: func(u *models.ProdukUmum) *string { return &u.OutwardRate }},
	{kunci: "OUTWARDRATEID", ambil: func(u *models.ProdukUmum) *string { return &u.OutwardRateID }},
	{kunci: "OUTWARDCOMM", ambil: func(u *models.ProdukUmum) *string { return &u.OutwardComm }},
	{kunci: "OUTWARDCOMMID", ambil: func(u *models.ProdukUmum) *string { return &u.OutwardCommID }},
	{kunci: "BENEFIT", ambil: func(u *models.ProdukUmum) *string { return &u.Benefit }},
	{kunci: "BENEFITID", ambil: func(u *models.ProdukUmum) *string { return &u.BenefitID }},
}

// Kunci boolean halaman umum.
const (
	kunciIsORS  = "IsORS"
	kunciIsView = "IsView"
	kunciID     = "ID"
)

// Kunci daftar halaman umum.
const (
	kunciLien      = "LienClause"
	kunciDokumen   = "DocumentClaim"
	kunciPlan      = "PlanList"
	kunciFinUW     = "FinancialUnderwritingList"
	kunciUWLimit   = "UnderwritingLimitList"
	kunciOutward   = "OutwardList"
	kunciKomentar  = "CommentList"
	kunciProductID = "PRODUCTID"
)

// Kunci skalar halaman `ProductNameInward` (`M_PRODUCTINWARD_LIFE.JSONDATA`).
var medanInward = []medanTeks[models.ProdukInward]{
	{kunci: "POLICYHODER", ambil: func(i *models.ProdukInward) *string { return &i.PolicyHolder }},
	{kunci: "POLICYHODERNAME", ambil: func(i *models.ProdukInward) *string { return &i.PolicyHolderName }},
	{kunci: "INSURED", ambil: func(i *models.ProdukInward) *string { return &i.Insured }},
	{kunci: "ADDENDUMNO", ambil: func(i *models.ProdukInward) *string { return &i.AddendumNo }},
	{kunci: "ADDENDUMWORD", ambil: func(i *models.ProdukInward) *string { return &i.AddendumWord }},
	{kunci: "AMANDEMENTNO", ambil: func(i *models.ProdukInward) *string { return &i.AmandementNo }},
	{kunci: "AMANDEMENTSCHD", ambil: func(i *models.ProdukInward) *string { return &i.AmandementSchd }},
	{kunci: "MAXEXPIREDCLAIM", ambil: func(i *models.ProdukInward) *string { return &i.MaxExpiredClaim }},
	{kunci: "BEGIN", ambil: func(i *models.ProdukInward) *string { return &i.Begin }, tanggal: true},
	{kunci: "STNC", ambil: func(i *models.ProdukInward) *string { return &i.STNC }, tanggal: true},
	{kunci: "CEDINGRETENTIONNUM", ambil: func(i *models.ProdukInward) *string { return &i.CedingRetentionNum }},
	{kunci: "CEDINGLIMIT", ambil: func(i *models.ProdukInward) *string { return &i.CedingLimit }},
	{kunci: "BROKERAGE", ambil: func(i *models.ProdukInward) *string { return &i.Brokerage }},
	{kunci: "MINAGE", ambil: func(i *models.ProdukInward) *string { return &i.MinAge }},
	{kunci: "MAXAGE", ambil: func(i *models.ProdukInward) *string { return &i.MaxAge }},
	{kunci: "EXPIRYAGE", ambil: func(i *models.ProdukInward) *string { return &i.ExpiryAge }},
	{kunci: "EXTRAPREMI", ambil: func(i *models.ProdukInward) *string { return &i.ExtraPremi }},
	{kunci: "MINSUMINSURED", ambil: func(i *models.ProdukInward) *string { return &i.MinSumInsured }},
	{kunci: "MAXSUMINSURED", ambil: func(i *models.ProdukInward) *string { return &i.MaxSumInsured }},
	{kunci: "MAXSUMREASURED", ambil: func(i *models.ProdukInward) *string { return &i.MaxSumReasured }},
	{kunci: "RNMSHARE", ambil: func(i *models.ProdukInward) *string { return &i.RNMShare }},
	{kunci: "RNMLIMITNUM", ambil: func(i *models.ProdukInward) *string { return &i.RNMLimitNum }},
	{kunci: "PREMIUMFACTOR", ambil: func(i *models.ProdukInward) *string { return &i.PremiumFactor }},
	{kunci: "PAYMENT", ambil: func(i *models.ProdukInward) *string { return &i.Payment }},
	{kunci: "SUBJECTTO", ambil: func(i *models.ProdukInward) *string { return &i.SubjectTo }},
	{kunci: "AnnuityInterest", ambil: func(i *models.ProdukInward) *string { return &i.AnnuityInterest }},
	{kunci: "PremiumRefundFactor", ambil: func(i *models.ProdukInward) *string { return &i.PremiumRefundFactor }},
	{kunci: "MAXDATARECEIVE", ambil: func(i *models.ProdukInward) *string { return &i.MaxDataReceive }},
	{kunci: "MATURE", ambil: func(i *models.ProdukInward) *string { return &i.Mature }, tanggal: true},
	{kunci: "BIRTHDAY", ambil: func(i *models.ProdukInward) *string { return &i.Birthday }},
	{kunci: "CURRENCY", ambil: func(i *models.ProdukInward) *string { return &i.Currency }},
	{kunci: "CURRENCYID", ambil: func(i *models.ProdukInward) *string { return &i.CurrencyID }},
	{kunci: "EXTRAMORTALITY", ambil: func(i *models.ProdukInward) *string { return &i.ExtraMortality }},
	{kunci: "MAXCONTRACT", ambil: func(i *models.ProdukInward) *string { return &i.MaxContract }},
	{kunci: "PROPORTIONALTABLE", ambil: func(i *models.ProdukInward) *string { return &i.ProportionalTable }},
	// Kunci view tanpa medan form - dari JSON lama.
	{kunci: "CEDING", ambil: func(i *models.ProdukInward) *string { return &i.Ceding }},
	{kunci: "TREATYNUMBER", ambil: func(i *models.ProdukInward) *string { return &i.TreatyNumber }},
	{kunci: "INWARDTREATYNM", ambil: func(i *models.ProdukInward) *string { return &i.InwardTreatyNm }},
	{kunci: "CEDINGRETENTIONPCT", ambil: func(i *models.ProdukInward) *string { return &i.CedingRetentionPct }},
	{kunci: "CEDINGLIMITXPN", ambil: func(i *models.ProdukInward) *string { return &i.CedingLimitXPN }},
	{kunci: "RNMLIMITPCT", ambil: func(i *models.ProdukInward) *string { return &i.RNMLimitPct }},
	{kunci: "LIENCLAUSE", ambil: func(i *models.ProdukInward) *string { return &i.LienClause }},
	{kunci: "MONTHS", ambil: func(i *models.ProdukInward) *string { return &i.Months }},
}

// kodekBaris - kunci satu jenis baris daftar + tempat kunci asingnya.
type kodekBaris[T any] struct {
	medan []medanTeks[T]
	asli  func(*T) *string
}

var (
	kodekLien = kodekBaris[models.BarisLien]{
		medan: []medanTeks[models.BarisLien]{
			{kunci: "Usia", ambil: func(b *models.BarisLien) *string { return &b.Usia }},
			{kunci: "Manfaat", ambil: func(b *models.BarisLien) *string { return &b.Manfaat }},
		},
		asli: func(b *models.BarisLien) *string { return &b.Asli },
	}
	kodekDokumen = kodekBaris[models.BarisDokumen]{
		medan: []medanTeks[models.BarisDokumen]{
			{kunci: "Document", ambil: func(b *models.BarisDokumen) *string { return &b.Document }},
		},
		asli: func(b *models.BarisDokumen) *string { return &b.Asli },
	}
	kodekPlan = kodekBaris[models.BarisPlan]{
		medan: []medanTeks[models.BarisPlan]{
			{kunci: "Plan", ambil: func(b *models.BarisPlan) *string { return &b.Plan }},
			{kunci: "PlanID", ambil: func(b *models.BarisPlan) *string { return &b.PlanID }},
			{kunci: "Name", ambil: func(b *models.BarisPlan) *string { return &b.Name }},
			{kunci: "Benefit", ambil: func(b *models.BarisPlan) *string { return &b.Benefit }},
			{kunci: "RIRATE", ambil: func(b *models.BarisPlan) *string { return &b.RIRate }},
			{kunci: "RIRATEID", ambil: func(b *models.BarisPlan) *string { return &b.RIRateID }},
		},
		asli: func(b *models.BarisPlan) *string { return &b.Asli },
	}
	kodekFinUW = kodekBaris[models.BarisFinUW]{
		medan: []medanTeks[models.BarisFinUW]{
			{kunci: "MinInsured", ambil: func(b *models.BarisFinUW) *string { return &b.MinInsured }},
			{kunci: "MaxInsured", ambil: func(b *models.BarisFinUW) *string { return &b.MaxInsured }},
			{kunci: "Employee", ambil: func(b *models.BarisFinUW) *string { return &b.Employee }},
			{kunci: "Non_Employee", ambil: func(b *models.BarisFinUW) *string { return &b.NonEmployee }},
		},
		asli: func(b *models.BarisFinUW) *string { return &b.Asli },
	}
	kodekUWLimit = kodekBaris[models.BarisUWLimit]{
		medan: []medanTeks[models.BarisUWLimit]{
			{kunci: "MinInsured", ambil: func(b *models.BarisUWLimit) *string { return &b.MinInsured }},
			{kunci: "MaxInsured", ambil: func(b *models.BarisUWLimit) *string { return &b.MaxInsured }},
			{kunci: "MinAge", ambil: func(b *models.BarisUWLimit) *string { return &b.MinAge }},
			{kunci: "MaxAge", ambil: func(b *models.BarisUWLimit) *string { return &b.MaxAge }},
			{kunci: "Medical", ambil: func(b *models.BarisUWLimit) *string { return &b.Medical }},
			{kunci: "Description", ambil: func(b *models.BarisUWLimit) *string { return &b.Description }},
		},
		asli: func(b *models.BarisUWLimit) *string { return &b.Asli },
	}
	kodekOutward = kodekBaris[models.BarisOutward]{
		medan: []medanTeks[models.BarisOutward]{
			{kunci: "REINSTYPEID", ambil: func(b *models.BarisOutward) *string { return &b.ReinsTypeID }},
			{kunci: "REINSTYPENAME", ambil: func(b *models.BarisOutward) *string { return &b.ReinsTypeName }},
			{kunci: "TRANSACTIONYEAR", ambil: func(b *models.BarisOutward) *string { return &b.TransactionYear }},
			{kunci: "TREATYCONTRACTID", ambil: func(b *models.BarisOutward) *string { return &b.TreatyContractID }},
			{kunci: "UNDERWRITINGYEAR", ambil: func(b *models.BarisOutward) *string { return &b.UnderwritingYear }},
			{kunci: "OVR_COMM", ambil: func(b *models.BarisOutward) *string { return &b.OvrComm }},
		},
		asli: func(b *models.BarisOutward) *string { return &b.Asli },
	}
	kodekKomentar = kodekBaris[models.BarisKomentar]{
		medan: []medanTeks[models.BarisKomentar]{
			{kunci: "Date", ambil: func(b *models.BarisKomentar) *string { return &b.Date }},
			{kunci: "OperatorName", ambil: func(b *models.BarisKomentar) *string { return &b.OperatorName }},
			{kunci: "IsApproved", ambil: func(b *models.BarisKomentar) *string { return &b.IsApproved }},
			{kunci: "Suggest", ambil: func(b *models.BarisKomentar) *string { return &b.Suggest }},
		},
		asli: func(b *models.BarisKomentar) *string { return &b.Asli },
	}
)

// --- membaca ------------------------------------------------------------------

// uraiObjek mengurai satu objek JSON tanpa melewati float.
func uraiObjek(teks string) (map[string]json.RawMessage, error) {
	hasil := map[string]json.RawMessage{}
	if strings.TrimSpace(teks) == "" {
		return hasil, nil
	}
	d := json.NewDecoder(strings.NewReader(teks))
	d.UseNumber()
	if err := d.Decode(&hasil); err != nil {
		return nil, err
	}
	return hasil, nil
}

// teksDari - nilai skalar sebagai teks: teks apa adanya, angka sebagai
// literalnya, boolean `true`/`false`, null/absen = "".
func teksDari(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	var v any
	if err := d.Decode(&v); err != nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case json.Number:
		return x.String()
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

// Bentuk tanggal Pega di JSON produk dan bentuk API.
const (
	BentukTanggalPega = "02/01/2006"
	bentukTanggalAPI  = "2006-01-02"
)

// TanggalKeAPI - `dd/MM/yyyy` → `YYYY-MM-DD`; teks lain apa adanya (data lama
// tidak disembunyikan).
func TanggalKeAPI(v string) string {
	t, err := time.Parse(BentukTanggalPega, strings.TrimSpace(v))
	if err != nil {
		return v
	}
	return t.Format(bentukTanggalAPI)
}

// TanggalKePega - `YYYY-MM-DD` → `dd/MM/yyyy`; kosong tetap kosong; teks lain
// apa adanya (services sudah menolaknya).
func TanggalKePega(v string) string {
	t, err := time.Parse(bentukTanggalAPI, strings.TrimSpace(v))
	if err != nil {
		return v
	}
	return t.Format(BentukTanggalPega)
}

func isiMedan[T any](obj map[string]json.RawMessage, medan []medanTeks[T], ke *T) {
	for _, m := range medan {
		v := teksDari(obj[m.kunci])
		if m.tanggal {
			v = TanggalKeAPI(v)
		}
		*m.ambil(ke) = v
	}
}

// uraiBaris mengurai satu larik objek; kunci asing tiap baris disimpan di `Asli`.
func uraiBaris[T any](raw json.RawMessage, k kodekBaris[T]) ([]T, error) {
	hasil := []T{}
	if len(bytes.TrimSpace(raw)) == 0 || string(bytes.TrimSpace(raw)) == "null" {
		return hasil, nil
	}
	var larik []json.RawMessage
	if err := json.Unmarshal(raw, &larik); err != nil {
		// Satu objek tunggal (Pega menulis PageList satu baris tetap sebagai larik;
		// data lama ganjil diterima, bukan dibuang diam-diam).
		larik = []json.RawMessage{raw}
	}
	dikelola := map[string]bool{}
	for _, m := range k.medan {
		dikelola[m.kunci] = true
	}
	for _, r := range larik {
		obj, err := uraiObjek(string(r))
		if err != nil {
			return nil, err
		}
		var b T
		isiMedan(obj, k.medan, &b)
		sisa := map[string]json.RawMessage{}
		for kunci, v := range obj {
			if !dikelola[kunci] {
				sisa[kunci] = v
			}
		}
		if len(sisa) > 0 {
			teks, err := rakitObjek(sisa)
			if err != nil {
				return nil, err
			}
			*k.asli(&b) = teks
		}
		hasil = append(hasil, b)
	}
	return hasil, nil
}

func benar(raw json.RawMessage) bool { return strings.EqualFold(teksDari(raw), "true") }

// UraiProduk merakit satu produk dari `JSONDATA` kedua tabel. `jsonInward`
// kosong = produk tanpa baris inward (DEV 196 vs 197 baris).
func UraiProduk(id, jsonUmum, idInward, jsonInward string) (models.Produk, error) {
	p := models.Produk{ID: id}
	rusak := func(tabel, idBaris string, err error) error {
		return fmt.Errorf("%w: %s %s: %v", ErrJSONRusak, tabel, idBaris, err)
	}
	umum, err := uraiObjek(jsonUmum)
	if err != nil {
		return p, rusak(TabelProduk, id, err)
	}
	isiMedan(umum, medanUmum, &p.Umum)
	p.Umum.IsORS = benar(umum[kunciIsORS])
	if p.LienClause, err = uraiBaris(umum[kunciLien], kodekLien); err != nil {
		return p, rusak(TabelProduk, id, err)
	}
	if p.DocumentClaim, err = uraiBaris(umum[kunciDokumen], kodekDokumen); err != nil {
		return p, rusak(TabelProduk, id, err)
	}
	if p.PlanList, err = uraiBaris(umum[kunciPlan], kodekPlan); err != nil {
		return p, rusak(TabelProduk, id, err)
	}
	if p.FinancialUnderwriting, err = uraiBaris(umum[kunciFinUW], kodekFinUW); err != nil {
		return p, rusak(TabelProduk, id, err)
	}
	if p.UnderwritingLimit, err = uraiBaris(umum[kunciUWLimit], kodekUWLimit); err != nil {
		return p, rusak(TabelProduk, id, err)
	}
	if p.OutwardList, err = uraiBaris(umum[kunciOutward], kodekOutward); err != nil {
		return p, rusak(TabelProduk, id, err)
	}
	if p.CommentList, err = uraiBaris(umum[kunciKomentar], kodekKomentar); err != nil {
		return p, rusak(TabelProduk, id, err)
	}
	inward, err := uraiObjek(jsonInward)
	if err != nil {
		return p, rusak(TabelInward, idInward, err)
	}
	isiMedan(inward, medanInward, &p.Inward)
	p.Inward.ID = idInward
	p.Inward.ProductID = teksDari(inward[kunciProductID])
	return p, nil
}

// RingkasanDari - satu baris grid daftar (`BrowseProduct_Life`) dari JSON umum.
func RingkasanDari(id, jsonUmum string) (models.RingkasanProduk, error) {
	obj, err := uraiObjek(jsonUmum)
	if err != nil {
		return models.RingkasanProduk{}, fmt.Errorf("%w: %s %s: %v", ErrJSONRusak, TabelProduk, id, err)
	}
	return models.RingkasanProduk{
		ID:           id,
		Ceding:       teksDari(obj["CEDING"]),
		TreatyNumber: teksDari(obj["TREATYNUMBER"]),
		InwardName:   teksDari(obj["INWARDNAME"]),
		CreateOp:     teksDari(obj["CREATEOP"]),
		UpdateOp:     teksDari(obj["UPDATEOP"]),
	}, nil
}

// --- menulis (bersama) ------------------------------------------------------

// rakitObjek menulis objek JSON berkunci urut tanpa meloloskan `<`, `>`, `&`
// (Pega menulis teks apa adanya).
func rakitObjek(obj map[string]json.RawMessage) (string, error) {
	kunci := make([]string, 0, len(obj))
	for k := range obj {
		kunci = append(kunci, k)
	}
	sort.Strings(kunci)
	var b strings.Builder
	b.WriteByte('{')
	for i, k := range kunci {
		if i > 0 {
			b.WriteByte(',')
		}
		kk, err := jsonTanpaLolos(k)
		if err != nil {
			return "", err
		}
		b.Write(kk)
		b.WriteByte(':')
		b.Write(obj[k])
	}
	b.WriteByte('}')
	return b.String(), nil
}

// jsonTanpaLolos - json.Marshal tanpa SetEscapeHTML.
func jsonTanpaLolos(v any) ([]byte, error) {
	var buf bytes.Buffer
	e := json.NewEncoder(&buf)
	e.SetEscapeHTML(false)
	if err := e.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

// ErrBarisAsliRusak - medan `asli` baris yang dikirim klien bukan objek JSON (400).
var ErrBarisAsliRusak = errors.New("repository: a list row carries unreadable original data (asli)")

// teksJSON - nilai skalar Pega selalu TEKS.
func teksJSON(v string) json.RawMessage {
	b, _ := jsonTanpaLolos(v)
	return b
}

func tulisMedan[T any](obj map[string]json.RawMessage, medan []medanTeks[T], dari *T) {
	for _, m := range medan {
		v := *m.ambil(dari)
		if m.tanggal {
			v = TanggalKePega(v)
		}
		obj[m.kunci] = teksJSON(v)
	}
}

// rakitBaris - satu larik; tiap baris = kunci aslinya (`Asli`) + kunci dikelola.
func rakitBaris[T any](daftar []T, k kodekBaris[T]) (json.RawMessage, error) {
	var b strings.Builder
	b.WriteByte('[')
	for i := range daftar {
		baris := daftar[i]
		obj, err := uraiObjek(*k.asli(&baris))
		if err != nil {
			return nil, fmt.Errorf("%w: row %d: %v", ErrBarisAsliRusak, i+1, err)
		}
		tulisMedan(obj, k.medan, &baris)
		teks, err := rakitObjek(obj)
		if err != nil {
			return nil, err
		}
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(teks)
	}
	b.WriteByte(']')
	return json.RawMessage(b.String()), nil
}

// bendera - TrueFalse Pega; tipe nilai lama dipertahankan (boolean JSON tetap
// boolean), selain itu teks `"true"`/`"false"`.
func bendera(lama json.RawMessage, v bool) json.RawMessage {
	teks := "false"
	if v {
		teks = "true"
	}
	if t := strings.TrimSpace(string(lama)); t == "true" || t == "false" {
		return json.RawMessage(teks)
	}
	return teksJSON(teks)
}

// RakitUmum menulis `M_PRODUCT_LIFE.JSONDATA` - padanan `@GetPageJSONString()`
// halaman `ProductName` (`SaveProductName_Act` 8 b1623). `lama` = JSON
// tersimpan (kosong untuk produk baru): kunci yang tidak dikelola layar
// dipertahankan; setiap kunci yang dibaca view `PRODUCT_LIFE` dijamin ada.
//
// `IsView`: produk baru tidak membawanya (`NewProductLife` tidak mengisinya);
// simpan sesudah `Edit` menulisnya `false` (`SetViewEdit` b145).
func RakitUmum(p models.Produk, lama string, baru bool) (string, error) {
	obj, err := uraiObjek(lama)
	if err != nil {
		return "", fmt.Errorf("%w: %s %s: %v", ErrJSONRusak, TabelProduk, p.ID, err)
	}
	obj[kunciID] = teksJSON(p.ID)
	tulisMedan(obj, medanUmum, &p.Umum)
	obj[kunciIsORS] = bendera(obj[kunciIsORS], p.Umum.IsORS)
	if !baru {
		obj[kunciIsView] = bendera(obj[kunciIsView], false)
	}
	larik := []struct {
		kunci string
		rakit func() (json.RawMessage, error)
	}{
		{kunciLien, func() (json.RawMessage, error) { return rakitBaris(p.LienClause, kodekLien) }},
		{kunciDokumen, func() (json.RawMessage, error) { return rakitBaris(p.DocumentClaim, kodekDokumen) }},
		{kunciPlan, func() (json.RawMessage, error) { return rakitBaris(p.PlanList, kodekPlan) }},
		{kunciFinUW, func() (json.RawMessage, error) { return rakitBaris(p.FinancialUnderwriting, kodekFinUW) }},
		{kunciUWLimit, func() (json.RawMessage, error) { return rakitBaris(p.UnderwritingLimit, kodekUWLimit) }},
		{kunciOutward, func() (json.RawMessage, error) { return rakitBaris(p.OutwardList, kodekOutward) }},
		{kunciKomentar, func() (json.RawMessage, error) { return rakitBaris(p.CommentList, kodekKomentar) }},
	}
	for _, l := range larik {
		raw, err := l.rakit()
		if err != nil {
			return "", err
		}
		obj[l.kunci] = raw
	}
	for _, k := range KunciViewProduk {
		if _, ada := obj[k]; !ada {
			obj[k] = teksJSON("")
		}
	}
	return rakitObjek(obj)
}
