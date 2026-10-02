package repository

// Kodek JSONDATA - membaca bentuk Pega (paket 1).

import (
	"errors"
	"strings"
	"testing"
)

// jsonUmumUji - halaman `ProductName` berbentuk Pega: angka sebagai angka JSON
// maupun teks, `POLICYHODER` (sic), `Non_Employee`, kunci asing di baris.
const jsonUmumUji = `{
  "ID": "UJI-01", "PRODUCTNAME": "UJI PRODUK", "CEDING": "UJI CEDING", "CEDINGID": "L0UJI",
  "SOBNAME": "UJI SOB", "SOBID": "L0SOB", "RICOMM": 0, "RIRISK": "UJI RISK", "RIRISKID": "1000117",
  "INWARDNAME": "UJI PRODUK UJI PEMEGANG", "TREATYNUMBER": "UJI/001", "CAUSE": "ANY CAUSE", "CAUSEID": "100004",
  "POLICYHODER": "UJI-ORG-1", "POLICYHODERNAME": "UJI PEMEGANG", "CREATEOP": "UJI-OP", "UPDATEOP": "UJI-OP",
  "IsORS": "true", "IsView": false, "Comment": "UJI komentar", "TYPE": "1", "KUNCILAMA": "dipertahankan",
  "PlanList": [{"Plan": "UJI PLAN", "PlanID": "P1", "Name": "UJI BIZ", "Benefit": "UJI MANFAAT", "RIRATE": "UJI RATE",
    "RIRATEID": "R1", "pxObjClass": "ASM-FW-GISFW-Data-Plan"}],
  "UnderwritingLimitList": [{"Description": "UJI", "Medical": "FCL", "MinAge": 22, "MaxAge": "70",
    "MinInsured": "0", "MaxInsured": 1175000000.123456789012}],
  "FinancialUnderwritingList": [{"MinInsured": "1", "MaxInsured": "2", "Employee": "Y", "Non_Employee": "N"}],
  "DocumentClaim": [{"Document": "UJI DOK"}],
  "LienClause": [{"Usia": "1", "Manfaat": "50"}],
  "OutwardList": [{"REINSTYPEID": "10200", "REINSTYPENAME": "OR", "TRANSACTIONYEAR": "2026", "TREATYCONTRACTID": "1000001",
    "UNDERWRITINGYEAR": "2026"}],
  "CommentList": [{"Date": "20260301T120000.123 GMT", "OperatorName": "UJI-OP", "Suggest": "UJI"}]
}`

const jsonInwardUji = `{"ID": "UJI-01", "PRODUCTID": "UJI-01", "BEGIN": "01/03/2026", "MATURE": "28/02/2027",
  "STNC": "26/03/2026", "CEDINGLIMIT": 150000000, "POLICYHODER": "UJI-ORG-1", "POLICYHODERNAME": "UJI PEMEGANG",
  "MAXEXPIREDCLAIM": "180", "MAXDATARECEIVE": "90", "CURRENCY": "IDR", "CURRENCYID": "1", "EXPIRYAGE": "75"}`

func TestUraiProdukMembacaKunciPega(t *testing.T) {
	p, err := UraiProduk("UJI-01", jsonUmumUji, "UJI-01", jsonInwardUji)
	if err != nil {
		t.Fatal(err)
	}
	u := p.Umum
	if u.PolicyHolder != "UJI-ORG-1" || u.PolicyHolderName != "UJI PEMEGANG" {
		t.Errorf("POLICYHODER (sic) tidak terbaca: %+v", u)
	}
	if u.RIComm != "0" || !u.IsORS || u.TypeBasicRider != "1" || u.CedingID != "L0UJI" {
		t.Errorf("skalar umum: %+v", u)
	}
	if len(p.PlanList) != 1 || p.PlanList[0].RIRateID != "R1" || !strings.Contains(p.PlanList[0].Asli, "pxObjClass") {
		t.Errorf("plan + kunci asing: %+v", p.PlanList)
	}
	if got := p.UnderwritingLimit[0].MaxInsured; got != "1175000000.123456789012" {
		t.Errorf("uang dari angka JSON berubah: %q", got)
	}
	if p.UnderwritingLimit[0].MinAge != "22" || p.FinancialUnderwriting[0].NonEmployee != "N" {
		t.Errorf("baris: %+v %+v", p.UnderwritingLimit, p.FinancialUnderwriting)
	}
	if p.CommentList[0].Date != "20260301T120000.123 GMT" {
		t.Errorf("tanggal komentar Pega harus apa adanya: %q", p.CommentList[0].Date)
	}
	in := p.Inward
	if in.Begin != "2026-03-01" || in.Mature != "2027-02-28" || in.STNC != "2026-03-26" {
		t.Errorf("tanggal dd/MM/yyyy → API: %+v", in)
	}
	if in.CedingLimit != "150000000" || in.ProductID != "UJI-01" || in.ID != "UJI-01" || in.ExpiryAge != "75" {
		t.Errorf("inward: %+v", in)
	}
	if p.OutwardList[0].TransactionYear != "2026" || p.LienClause[0].Manfaat != "50" || p.DocumentClaim[0].Document != "UJI DOK" {
		t.Errorf("daftar lain: %+v %+v %+v", p.OutwardList, p.LienClause, p.DocumentClaim)
	}
}

func TestUraiProdukTanpaInwardDanDaftarKosong(t *testing.T) {
	p, err := UraiProduk("UJI-02", `{"PRODUCTNAME":"X"}`, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if p.PlanList == nil || len(p.PlanList) != 0 || p.Inward.ID != "" {
		t.Errorf("daftar absen = larik kosong, inward absen = kosong: %+v", p)
	}
}

func TestUraiProdukJSONRusakGagalTerang(t *testing.T) {
	if _, err := UraiProduk("UJI-03", `{"PRODUCTNAME":`, "", ""); !errors.Is(err, ErrJSONRusak) {
		t.Errorf("JSON rusak harus ErrJSONRusak, dapat %v", err)
	}
	if _, err := UraiProduk("UJI-03", `{}`, "UJI-03", `[`); !errors.Is(err, ErrJSONRusak) {
		t.Errorf("JSON inward rusak harus ErrJSONRusak, dapat %v", err)
	}
}

func TestRingkasanDariKunciGrid(t *testing.T) {
	r, err := RingkasanDari("UJI-01", jsonUmumUji)
	if err != nil {
		t.Fatal(err)
	}
	if r.Ceding != "UJI CEDING" || r.TreatyNumber != "UJI/001" || r.InwardName != "UJI PRODUK UJI PEMEGANG" ||
		r.CreateOp != "UJI-OP" || r.UpdateOp != "UJI-OP" {
		t.Errorf("ringkasan: %+v", r)
	}
}

func TestTanggalPegaBolakBalik(t *testing.T) {
	if TanggalKeAPI("01/03/2026") != "2026-03-01" || TanggalKePega("2026-03-01") != "01/03/2026" {
		t.Error("dd/MM/yyyy ↔ YYYY-MM-DD")
	}
	if TanggalKeAPI("bukan tanggal") != "bukan tanggal" || TanggalKePega("") != "" {
		t.Error("teks bukan tanggal apa adanya; kosong tetap kosong")
	}
}

func TestPilihInwardPRODUCTIDLaluID(t *testing.T) {
	baris := []barisJSON{
		{id: "100001", isi: `{"PRODUCTID":"100002"}`},
		{id: "100002", isi: `{"PRODUCTID":"100009"}`},
		{id: "100003", isi: `{"PRODUCTID":"100002"}`},
	}
	if b, ada := pilihInward("100002", baris); !ada || b.id != "100003" {
		t.Errorf("PRODUCTID cocok, terakhir menang: %+v %v", b, ada)
	}
	// Audit 02-10-2026: baris ber-ID = PRODUCTID = produk (yang dibaca pembaca hilir `WHERE ID = produk`)
	// menang atas baris lain ber-PRODUCTID sama, walau ID-nya lebih kecil.
	sendiri := []barisJSON{
		{id: "100002", isi: `{"PRODUCTID":"100002"}`},
		{id: "100005", isi: `{"PRODUCTID":"100002"}`},
	}
	if b, ada := pilihInward("100002", sendiri); !ada || b.id != "100002" {
		t.Errorf("baris milik sendiri ber-ID produk didahulukan: %+v %v", b, ada)
	}
	// Baris ber-ID = produk tetapi PRODUCTID-nya produk LAIN (ID sequence inward warisan): bukan milik produk ini.
	if b, ada := pilihInward("100002", baris[1:2]); ada {
		t.Errorf("baris milik produk 100009 tidak boleh dipilih untuk 100002: %+v", b)
	}
	if pid := inwardMilikLain("100002", baris[1:2]); pid != "100009" {
		t.Errorf("pemilik baris ber-ID sama: %q", pid)
	}
	kosong := []barisJSON{{id: "100002", isi: `{"ID":"100002"}`}}
	if b, ada := pilihInward("100002", kosong); !ada || b.id != "100002" || inwardMilikLain("100002", kosong) != "" {
		t.Errorf("tanpa PRODUCTID, ID = produk: %+v %v", b, ada)
	}
	if _, ada := pilihInward("100005", baris); ada {
		t.Error("tidak ada yang cocok")
	}
}
