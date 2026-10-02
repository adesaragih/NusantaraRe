package repository

// Tabel FLAT produk - keputusan work owner 02-10-2026 (K2, K5; tiket 01 bab bertanggal 02-10-2026): satu produk = satu
// baris induk `M_PRODUCTNAME_LIFE` (sisi umum + sisi inward digabung, grilling Q1b) + tujuh tabel anak berawalan
// `M_PRODUCTNAME_LIFE_` (grilling D2). DDL: `migrations/140`–`147`; kolom dan tipe: `docs/STRUKTUR-TABEL-…md`.
//
// Untuk apa berkas ini: pemetaan SATU-SATUNYA antara medan model dan kolom flat, dan bentuk KANONIK nilai yang
// dikembalikan Oracle (NormalkanFlat). Tiga pemakai, satu aturan: penulis/pembaca Oracle (`mpnl_flat_sql.go`), alat
// pindah dan rekonsiliasinya (`mpnl_pindah.go`), dan gudang tiruan uji. Uji `db` membuktikan Oracle mengembalikan
// persis bentuk yang diramalkan di sini.
//
// ⛔ Tipe patuh penjaga inti `TestNolNumberTanpaPresisi` (K6): desimal NUMBER(38,8) - paling banyak 8 angka di
// belakang koma dan 30 di depannya; bulat NUMBER(5). Nilai yang tidak muat menjadi MasalahNilai, TIDAK PERNAH
// dipotong atau dibulatkan diam-diam. Tidak satu angka pun melewati float (ADR-0003).

import (
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

// Nama tabel flat.
const (
	TabelFlatInduk    = "M_PRODUCTNAME_LIFE"
	TabelFlatLien     = "M_PRODUCTNAME_LIFE_LIEN"
	TabelFlatDokumen  = "M_PRODUCTNAME_LIFE_DOCCLAIM"
	TabelFlatPlan     = "M_PRODUCTNAME_LIFE_PLAN"
	TabelFlatFinUW    = "M_PRODUCTNAME_LIFE_FINUW"
	TabelFlatUWLimit  = "M_PRODUCTNAME_LIFE_UWLIMIT"
	TabelFlatOutward  = "M_PRODUCTNAME_LIFE_OUTWARD"
	TabelFlatKomentar = "M_PRODUCTNAME_LIFE_COMMENT"
)

// DaftarTabelFlat - induk lalu ketujuh anak, urutan migrasi 140–147.
var DaftarTabelFlat = []string{TabelFlatInduk, TabelFlatLien, TabelFlatDokumen, TabelFlatPlan, TabelFlatFinUW,
	TabelFlatUWLimit, TabelFlatOutward, TabelFlatKomentar}

// JenisKolomFlat - golongan fisik kolom flat.
type JenisKolomFlat int

const (
	// FlatTeks - VARCHAR2(Lebar); lebar dalam BYTE (semantik kolom DEV tidak tercatat - byte aman untuk keduanya).
	FlatTeks JenisKolomFlat = iota
	// FlatDesimal - NUMBER(38,8): uang, persen, rate, faktor.
	FlatDesimal
	// FlatBulat - NUMBER(5): usia, nomor, hari, kontrak, tahun.
	FlatBulat
	// FlatTanggal - DATE; model `YYYY-MM-DD`.
	FlatTanggal
	// FlatStempel - TIMESTAMP; model stempel Pega `20060102T150405.000 GMT` (`@CurrentDateTime()`), disimpan GMT.
	FlatStempel
)

// Batas NUMBER(38,8) dan NUMBER(5).
const (
	SkalaDesimalFlat = 8
	DigitBulatFlat   = 30
	DigitKecilFlat   = 5
)

// BentukStempelPega - stempel tanggal baris komentar (`AddCommentList_Act` 1 b235).
const BentukStempelPega = "20060102T150405.000 GMT"

// KolomFlat - satu kolom flat ↔ satu medan teks model T.
type KolomFlat[T any] struct {
	Nama  string
	Jenis JenisKolomFlat
	Lebar int // FlatTeks saja
	// Label - label medan VERBATIM layar (PARITAS §3–§5), untuk kalimat penolakan services.
	Label string
	ambil func(*T) *string
}

func teks[T any](nama string, lebar int, label string, ambil func(*T) *string) KolomFlat[T] {
	return KolomFlat[T]{Nama: nama, Jenis: FlatTeks, Lebar: lebar, Label: label, ambil: ambil}
}

func desimal[T any](nama, label string, ambil func(*T) *string) KolomFlat[T] {
	return KolomFlat[T]{Nama: nama, Jenis: FlatDesimal, Label: label, ambil: ambil}
}

func bulat[T any](nama, label string, ambil func(*T) *string) KolomFlat[T] {
	return KolomFlat[T]{Nama: nama, Jenis: FlatBulat, Label: label, ambil: ambil}
}

func tanggal[T any](nama, label string, ambil func(*T) *string) KolomFlat[T] {
	return KolomFlat[T]{Nama: nama, Jenis: FlatTanggal, Label: label, ambil: ambil}
}

type pr = models.Produk

// KolomFlatInduk - kolom `M_PRODUCTNAME_LIFE` selain `ID` (kunci) dan `IS_ORS` (bendera, KolomIsORS), urutan DDL 140.
//
// ⛔ `POLICYHOLDER`/`POLICYHOLDERNAME`: SATU pasang kolom untuk kedua sisi (sama di 196 produk DEV), sumbernya sisi
// inward (`SaveProductName_Act` 1 b535/b556 menyalinnya ke sisi umum). Pembaca mengisi keduanya.
var KolomFlatInduk = []KolomFlat[pr]{
	teks("PRODUCTNAME", 1000, "Product Name", func(p *pr) *string { return &p.Umum.ProductName }),
	teks("CEDINGID", 100, "Ceding ID", func(p *pr) *string { return &p.Umum.CedingID }),
	teks("CEDING", 1000, "Ceding", func(p *pr) *string { return &p.Umum.Ceding }),
	teks("SOBID", 100, "SOB ID", func(p *pr) *string { return &p.Umum.SOBID }),
	teks("SOBNAME", 1000, "SOB", func(p *pr) *string { return &p.Umum.SOBName }),
	desimal("RICOMM", "Deduction (%)", func(p *pr) *string { return &p.Umum.RIComm }),
	teks("RIRISKID", 10, "R/I Risk Name ID", func(p *pr) *string { return &p.Umum.RIRiskID }),
	teks("RIRISK", 100, "R/I Risk Name", func(p *pr) *string { return &p.Umum.RIRisk }),
	teks("RIRATEID", 10, "R/I Rate ID", func(p *pr) *string { return &p.Umum.RIRateID }),
	teks("RIRATE", 500, "R/I Rate", func(p *pr) *string { return &p.Umum.RIRate }),
	teks("INWARDNAME", 1000, "Treaty Name", func(p *pr) *string { return &p.Umum.InwardName }),
	teks("TREATYNUMBER", 100, "Treaty Number", func(p *pr) *string { return &p.Umum.TreatyNumber }),
	teks("CAUSEID", 10, "Cause Of Loss ID", func(p *pr) *string { return &p.Umum.CauseID }),
	teks("CAUSE", 100, "Cause Of Loss", func(p *pr) *string { return &p.Umum.Cause }),
	teks("POLICYHOLDER", 100, "Policy Holder ID", func(p *pr) *string { return &p.Inward.PolicyHolder }),
	teks("POLICYHOLDERNAME", 1000, "Policy Holder", func(p *pr) *string { return &p.Inward.PolicyHolderName }),
	teks("INSURED", 4000, "Insured", func(p *pr) *string { return &p.Inward.Insured }),
	bulat("ADDENDUMNO", "Addendum No.", func(p *pr) *string { return &p.Inward.AddendumNo }),
	teks("ADDENDUMWORD", 4000, "Addendum", func(p *pr) *string { return &p.Inward.AddendumWord }),
	bulat("AMANDEMENTNO", "Amandement No.", func(p *pr) *string { return &p.Inward.AmandementNo }),
	teks("AMANDEMENTSCHD", 4000, "Amandement", func(p *pr) *string { return &p.Inward.AmandementSchd }),
	bulat("MAXEXPIREDCLAIM", "Max Notification Claim Expired", func(p *pr) *string { return &p.Inward.MaxExpiredClaim }),
	tanggal("BEGIN_DATE", "Begin Date", func(p *pr) *string { return &p.Inward.Begin }),
	tanggal("STNC", "STNC", func(p *pr) *string { return &p.Inward.STNC }),
	tanggal("MATURE", "Expired Date", func(p *pr) *string { return &p.Inward.Mature }),
	desimal("CEDINGRETENTIONNUM", "Ceding Retention (%)", func(p *pr) *string { return &p.Inward.CedingRetentionNum }),
	desimal("CEDINGLIMIT", "Ceding's Limit", func(p *pr) *string { return &p.Inward.CedingLimit }),
	desimal("BROKERAGE", "Brokerage Fee (%)", func(p *pr) *string { return &p.Inward.Brokerage }),
	bulat("MINAGE", "Minimum Age (Years)", func(p *pr) *string { return &p.Inward.MinAge }),
	bulat("MAXAGE", "Maximum Age (Years)", func(p *pr) *string { return &p.Inward.MaxAge }),
	bulat("EXPIRYAGE", "Expiry Age (Years)", func(p *pr) *string { return &p.Inward.ExpiryAge }),
	desimal("EXTRAPREMI", "Extra Premium", func(p *pr) *string { return &p.Inward.ExtraPremi }),
	desimal("MINSUMINSURED", "Min Sum Insured", func(p *pr) *string { return &p.Inward.MinSumInsured }),
	desimal("MAXSUMINSURED", "Max Sum Insured", func(p *pr) *string { return &p.Inward.MaxSumInsured }),
	desimal("MAXSUMREASURED", "Max Sum Reasured", func(p *pr) *string { return &p.Inward.MaxSumReasured }),
	desimal("RNMSHARE", "Nusantara Re Share (%)", func(p *pr) *string { return &p.Inward.RNMShare }),
	desimal("RNMLIMITNUM", "Nusantara Re's Limit", func(p *pr) *string { return &p.Inward.RNMLimitNum }),
	desimal("PREMIUMFACTOR", "Premium Factor (%)", func(p *pr) *string { return &p.Inward.PremiumFactor }),
	teks("PAYMENT", 1, "Premium Payment Method", func(p *pr) *string { return &p.Inward.Payment }),
	teks("SUBJECTTO", 4000, "Subject To", func(p *pr) *string { return &p.Inward.SubjectTo }),
	desimal("ANNUITYINTEREST", "Annuity Interest (%)", func(p *pr) *string { return &p.Inward.AnnuityInterest }),
	desimal("PREMIUMREFUNDFACTOR", "Premium Refund Factor (%)", func(p *pr) *string { return &p.Inward.PremiumRefundFactor }),
	bulat("MAXDATARECEIVE", "Max Production Data Receive", func(p *pr) *string { return &p.Inward.MaxDataReceive }),
	teks("BIRTHDAY", 1, "Birthday", func(p *pr) *string { return &p.Inward.Birthday }),
	teks("CURRENCYID", 100, "Currency ID", func(p *pr) *string { return &p.Inward.CurrencyID }),
	teks("CURRENCY", 20, "Currency", func(p *pr) *string { return &p.Inward.Currency }),
	desimal("EXTRAMORTALITY", "Extra Mortality (%)", func(p *pr) *string { return &p.Inward.ExtraMortality }),
	bulat("MAXCONTRACT", "Max Contract (year)", func(p *pr) *string { return &p.Inward.MaxContract }),
	desimal("PROPORTIONALTABLE", "Proportional Table", func(p *pr) *string { return &p.Inward.ProportionalTable }),
	teks("CREATEOP", 100, "Create Operator", func(p *pr) *string { return &p.Umum.CreateOp }),
	teks("UPDATEOP", 100, "Last Updated Operator", func(p *pr) *string { return &p.Umum.UpdateOp }),
}

// Kolom kunci dan bendera induk.
const (
	KolomIDFlat    = "ID"
	KolomIsORS     = "IS_ORS"
	KolomProductID = "PRODUCTID"
	KolomUrut      = "URUT"
)

// AnakFlat - satu tabel anak: nama, label daftar (kalimat penolakan), kolom isi (urutan DDL, sesudah PRODUCTID, URUT).
type AnakFlat[T any] struct {
	Tabel, Label string
	Kolom        []KolomFlat[T]
	daftar       func(*models.Produk) *[]T
}

var (
	AnakLien = AnakFlat[models.BarisLien]{Tabel: TabelFlatLien, Label: "Lien Clause",
		daftar: func(p *pr) *[]models.BarisLien { return &p.LienClause },
		Kolom: []KolomFlat[models.BarisLien]{
			teks("USIA", 200, "Usia saat Klaim", func(b *models.BarisLien) *string { return &b.Usia }),
			teks("MANFAAT", 200, "% Manfaat yang dibayarkan", func(b *models.BarisLien) *string { return &b.Manfaat }),
		}}
	AnakDokumen = AnakFlat[models.BarisDokumen]{Tabel: TabelFlatDokumen, Label: "Document Claim",
		daftar: func(p *pr) *[]models.BarisDokumen { return &p.DocumentClaim },
		Kolom: []KolomFlat[models.BarisDokumen]{
			teks("DOCUMENT", 500, "Document List", func(b *models.BarisDokumen) *string { return &b.Document }),
		}}
	AnakPlan = AnakFlat[models.BarisPlan]{Tabel: TabelFlatPlan, Label: "Plan List",
		daftar: func(p *pr) *[]models.BarisPlan { return &p.PlanList },
		Kolom: []KolomFlat[models.BarisPlan]{
			teks("PLANID", 10, "Plan Name ID", func(b *models.BarisPlan) *string { return &b.PlanID }),
			teks("PLAN", 200, "Plan Name", func(b *models.BarisPlan) *string { return &b.Plan }),
			teks("NAME", 200, "Bussines", func(b *models.BarisPlan) *string { return &b.Name }),
			teks("BENEFIT", 500, "Benefit", func(b *models.BarisPlan) *string { return &b.Benefit }),
			teks("RIRATEID", 10, "R/I Rate ID", func(b *models.BarisPlan) *string { return &b.RIRateID }),
			teks("RIRATE", 500, "R/I Rate", func(b *models.BarisPlan) *string { return &b.RIRate }),
		}}
	AnakFinUW = AnakFlat[models.BarisFinUW]{Tabel: TabelFlatFinUW, Label: "Financial Underwriting",
		daftar: func(p *pr) *[]models.BarisFinUW { return &p.FinancialUnderwriting },
		Kolom: []KolomFlat[models.BarisFinUW]{
			desimal("MININSURED", "Min Insured", func(b *models.BarisFinUW) *string { return &b.MinInsured }),
			desimal("MAXINSURED", "Max Insured", func(b *models.BarisFinUW) *string { return &b.MaxInsured }),
			teks("EMPLOYEE", 1000, "Employee", func(b *models.BarisFinUW) *string { return &b.Employee }),
			teks("NON_EMPLOYEE", 1000, "Non-Employee", func(b *models.BarisFinUW) *string { return &b.NonEmployee }),
		}}
	AnakUWLimit = AnakFlat[models.BarisUWLimit]{Tabel: TabelFlatUWLimit, Label: "Underwriting Limit",
		daftar: func(p *pr) *[]models.BarisUWLimit { return &p.UnderwritingLimit },
		Kolom: []KolomFlat[models.BarisUWLimit]{
			desimal("MININSURED", "Min Insured", func(b *models.BarisUWLimit) *string { return &b.MinInsured }),
			desimal("MAXINSURED", "Max Insured", func(b *models.BarisUWLimit) *string { return &b.MaxInsured }),
			bulat("MINAGE", "Min Age", func(b *models.BarisUWLimit) *string { return &b.MinAge }),
			bulat("MAXAGE", "Max Age", func(b *models.BarisUWLimit) *string { return &b.MaxAge }),
			teks("MEDICAL", 200, "Medical", func(b *models.BarisUWLimit) *string { return &b.Medical }),
			teks("DESCRIPTION", 1000, "Description", func(b *models.BarisUWLimit) *string { return &b.Description }),
		}}
	AnakOutward = AnakFlat[models.BarisOutward]{Tabel: TabelFlatOutward, Label: "Outward List",
		daftar: func(p *pr) *[]models.BarisOutward { return &p.OutwardList },
		Kolom: []KolomFlat[models.BarisOutward]{
			teks("REINSTYPEID", 100, "Reins Type ID", func(b *models.BarisOutward) *string { return &b.ReinsTypeID }),
			teks("REINSTYPENAME", 100, "Reins Type", func(b *models.BarisOutward) *string { return &b.ReinsTypeName }),
			bulat("TRANSACTIONYEAR", "Transaction Year", func(b *models.BarisOutward) *string { return &b.TransactionYear }),
			teks("TREATYCONTRACTID", 100, "Treaty Contract ID", func(b *models.BarisOutward) *string { return &b.TreatyContractID }),
			bulat("UNDERWRITINGYEAR", "Underwriting Year", func(b *models.BarisOutward) *string { return &b.UnderwritingYear }),
			desimal("OVR_COMM", "Override Commission", func(b *models.BarisOutward) *string { return &b.OvrComm }),
		}}
	AnakKomentar = AnakFlat[models.BarisKomentar]{Tabel: TabelFlatKomentar, Label: "Comment List",
		daftar: func(p *pr) *[]models.BarisKomentar { return &p.CommentList },
		Kolom: []KolomFlat[models.BarisKomentar]{
			{Nama: "TANGGAL", Jenis: FlatStempel, Label: "Date",
				ambil: func(b *models.BarisKomentar) *string { return &b.Date }},
			teks("OPERATORNAME", 100, "PIC", func(b *models.BarisKomentar) *string { return &b.OperatorName }),
			teks("SUGGEST", 4000, "Comment", func(b *models.BarisKomentar) *string { return &b.Suggest }),
		}}
)

// --- bentuk kanonik ---------------------------------------------------------------

// Jenis MasalahNilai.
const (
	MasalahBukanAngka     = "bukan angka"
	MasalahBukanTanggal   = "bukan tanggal"
	MasalahBukanStempel   = "bukan stempel tanggal Pega"
	MasalahTerlaluPanjang = "melebihi lebar kolom"
	MasalahSkala          = "lebih dari 8 angka desimal"
	MasalahDigitBulat     = "lebih dari 30 digit bulat"
	MasalahBukanBulat     = "bukan bilangan bulat"
	MasalahDigitKecil     = "lebih dari 5 digit"
)

// MasalahNilai - satu nilai yang tidak muat kolom flatnya. ⛔ `Nilai` tidak pernah dicetak laporan pindah (agregat
// saja, brief T2); services boleh mengutipnya kembali kepada pengirimnya sendiri.
type MasalahNilai struct {
	Tabel, Kolom, Label string
	Urut                int // 0 = induk; n = baris ke-n daftar
	Jenis               string
	Nilai               string
	Batas               int // FlatTeks: lebar kolom
}

// AngkaKanonik - teks angka dalam bentuk yang dikembalikan Oracle lewat `db.FmtDesimal` lalu dibaca modul ini: titik
// desimal, tanpa nol ekor, tanpa notasi ilmiah, berawalan `0` di depan titik. Koma diterima sebagai titik bila tidak
// ada titik (services `desimal`); koma DAN titik, eksponen, atau teks lain = galat.
func AngkaKanonik(s string) (string, *apd.Decimal, bool) {
	s = strings.TrimSpace(s)
	if strings.Contains(s, ",") {
		if strings.Contains(s, ".") {
			return "", nil, false
		}
		s = strings.ReplaceAll(s, ",", ".")
	}
	if s == "" || strings.ContainsAny(s, "eE") {
		return "", nil, false
	}
	d, err := utils.ParseDecimal(s)
	if err != nil {
		return "", nil, false
	}
	d.Reduce(d)
	if d.IsZero() {
		d.Negative = false
		return "0", d, true
	}
	return d.Text('f'), d, true
}

// muatDesimal - NUMBER(38,8): paling banyak 8 angka di belakang koma dan 30 di depannya.
func muatDesimal(d *apd.Decimal) string {
	if d.Exponent < -SkalaDesimalFlat {
		return MasalahSkala
	}
	if int(d.NumDigits())+int(d.Exponent) > DigitBulatFlat {
		return MasalahDigitBulat
	}
	return ""
}

// muatBulat - NUMBER(5): bilangan bulat, paling banyak 5 digit.
func muatBulat(d *apd.Decimal) string {
	if d.Exponent < 0 {
		return MasalahBukanBulat
	}
	if !d.IsZero() && int(d.NumDigits())+int(d.Exponent) > DigitKecilFlat {
		return MasalahDigitKecil
	}
	return ""
}

// NilaiKanonik - satu nilai dalam bentuk kanonik kolomnya; masalah = jenis MasalahNilai ("" = muat). Kosong = NULL,
// selalu muat. Nilai yang tidak muat dikembalikan "" (kolomnya NULL).
func NilaiKanonik(jenis JenisKolomFlat, lebar int, v string) (string, string) {
	if v == "" {
		return "", ""
	}
	switch jenis {
	case FlatTeks:
		if len(v) > lebar {
			return "", MasalahTerlaluPanjang
		}
		return v, ""
	case FlatDesimal, FlatBulat:
		k, d, ok := AngkaKanonik(v)
		if !ok {
			return "", MasalahBukanAngka
		}
		cek := muatDesimal
		if jenis == FlatBulat {
			cek = muatBulat
		}
		if m := cek(d); m != "" {
			return "", m
		}
		return k, ""
	case FlatTanggal:
		t, err := time.Parse(bentukTanggalAPI, strings.TrimSpace(v))
		if err != nil {
			return "", MasalahBukanTanggal
		}
		return t.Format(bentukTanggalAPI), ""
	case FlatStempel:
		t, err := time.Parse(BentukStempelPega, strings.TrimSpace(v))
		if err != nil {
			return "", MasalahBukanStempel
		}
		return t.UTC().Format(BentukStempelPega), ""
	}
	return v, ""
}

func normalkanMedan[T any](kolom []KolomFlat[T], tabel string, urut int, b *T, masalah *[]MasalahNilai) {
	for _, k := range kolom {
		v := k.ambil(b)
		baru, m := NilaiKanonik(k.Jenis, k.Lebar, *v)
		if m != "" {
			*masalah = append(*masalah, MasalahNilai{Tabel: tabel, Kolom: k.Nama, Label: k.Label, Urut: urut, Jenis: m,
				Nilai: *v, Batas: k.Lebar})
		}
		*v = baru
	}
}

func normalkanAnak[T any](a AnakFlat[T], p *models.Produk, masalah *[]MasalahNilai) {
	daftar := *a.daftar(p)
	salinan := make([]T, len(daftar))
	copy(salinan, daftar)
	for i := range salinan {
		normalkanMedan(a.Kolom, a.Tabel, i+1, &salinan[i], masalah)
	}
	*a.daftar(p) = salinan
}

// NormalkanFlat - produk sebagaimana AKAN dibaca kembali dari tabel flat: setiap medan berkolom dalam bentuk kanonik
// kolomnya (nilai yang tidak muat = "" + MasalahNilai), medan TANPA kolom dikosongkan (lihat MedanTanpaKolom), ID inward
// dan PRODUCTID = ID produk, pemegang polis kedua sisi dari satu pasang kolom, kunci asing baris (`Asli`) dibuang (D2).
func NormalkanFlat(p models.Produk) (models.Produk, []MasalahNilai) {
	var masalah []MasalahNilai
	n := models.Produk{ID: p.ID, LienClause: p.LienClause, DocumentClaim: p.DocumentClaim, PlanList: p.PlanList,
		FinancialUnderwriting: p.FinancialUnderwriting, UnderwritingLimit: p.UnderwritingLimit,
		OutwardList: p.OutwardList, CommentList: p.CommentList}
	for _, k := range KolomFlatInduk {
		*k.ambil(&n) = *k.ambil(&p)
	}
	n.Umum.IsORS = p.Umum.IsORS
	normalkanMedan(KolomFlatInduk, TabelFlatInduk, 0, &n, &masalah)
	n.Umum.PolicyHolder, n.Umum.PolicyHolderName = n.Inward.PolicyHolder, n.Inward.PolicyHolderName
	n.Inward.ID, n.Inward.ProductID = p.ID, p.ID
	normalkanAnak(AnakLien, &n, &masalah)
	normalkanAnak(AnakDokumen, &n, &masalah)
	normalkanAnak(AnakPlan, &n, &masalah)
	normalkanAnak(AnakFinUW, &n, &masalah)
	normalkanAnak(AnakUWLimit, &n, &masalah)
	normalkanAnak(AnakOutward, &n, &masalah)
	normalkanAnak(AnakKomentar, &n, &masalah)
	buangAsli(&n)
	return n, masalah
}

// buangAsli - kunci baris JSON lama (`Asli`) tidak punya kolom (D2).
func buangAsli(p *models.Produk) {
	for i := range p.LienClause {
		p.LienClause[i].Asli = ""
	}
	for i := range p.DocumentClaim {
		p.DocumentClaim[i].Asli = ""
	}
	for i := range p.PlanList {
		p.PlanList[i].Asli = ""
	}
	for i := range p.FinancialUnderwriting {
		p.FinancialUnderwriting[i].Asli = ""
	}
	for i := range p.UnderwritingLimit {
		p.UnderwritingLimit[i].Asli = ""
	}
	for i := range p.OutwardList {
		p.OutwardList[i].Asli = ""
	}
	for i := range p.CommentList {
		p.CommentList[i].Asli = ""
	}
}

// MasalahFlat - nilai produk yang tidak muat kolom flatnya (penulis menolak produk ini; services lebih dulu menolaknya
// berkalimat).
func MasalahFlat(p models.Produk) []MasalahNilai {
	_, m := NormalkanFlat(p)
	return m
}

// MedanTanpaKolom - medan model yang TIDAK menjadi kolom flat (tiket 01 bab 02-10-2026 "Tidak menjadi kolom"): medan
// layar mati sisi umum dan kunci inward warisan - 0 terisi di DEV. Nama = kunci Pega; nilai terisi = nilainya hilang
// bila dipindah (alat pindah GAGAL). `Comment`/`IsView` bukan di sini: keadaan layar dan masukan, bukan data.
func MedanTanpaKolom(p models.Produk) []string {
	var hasil []string
	for kunci, v := range map[string]string{
		"TYPE": p.Umum.TypeBasicRider, "TYPE_CEDING": p.Umum.TypeCeding, "GRUP": p.Umum.Grup,
		"PRODUCTCODE": p.Umum.ProductCode, "PRODUCTTYPE": p.Umum.ProductType, "PRODUCTTYPEID": p.Umum.ProductTypeID,
		"RICOMMID": p.Umum.RICommID, "OUTWARDNAME": p.Umum.OutwardName, "OUTWARDNAMEID": p.Umum.OutwardNameID,
		"OUTWARDRATE": p.Umum.OutwardRate, "OUTWARDRATEID": p.Umum.OutwardRateID, "OUTWARDCOMM": p.Umum.OutwardComm,
		"OUTWARDCOMMID": p.Umum.OutwardCommID, "BENEFIT": p.Umum.Benefit, "BENEFITID": p.Umum.BenefitID,
		"inward.CEDING": p.Inward.Ceding, "inward.TREATYNUMBER": p.Inward.TreatyNumber,
		"inward.INWARDTREATYNM": p.Inward.InwardTreatyNm, "inward.CEDINGRETENTIONPCT": p.Inward.CedingRetentionPct,
		"inward.CEDINGLIMITXPN": p.Inward.CedingLimitXPN, "inward.RNMLIMITPCT": p.Inward.RNMLimitPct,
		"inward.LIENCLAUSE": p.Inward.LienClause, "inward.MONTHS": p.Inward.Months,
	} {
		if strings.TrimSpace(v) != "" {
			hasil = append(hasil, kunci)
		}
	}
	return hasil
}
