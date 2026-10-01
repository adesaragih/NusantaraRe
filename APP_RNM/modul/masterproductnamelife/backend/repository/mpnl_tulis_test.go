package repository

// Penulis sisi umum (paket 3) - identitas dan perakit JSONDATA, tanpa Oracle.

import (
	"errors"
	"strings"
	"testing"

	go_ora "github.com/sijms/go-ora/v2"

	"nusantarare/modul/masterproductnamelife/backend/models"
)

func TestFormatIdentitasLimaDigitSepertiProsedur(t *testing.T) {
	// `dba-procedures-and-ddl.md` §1: concat('1', lpad(M_PRODUCT_LIFE_SEQ.nextval, 5, '0')).
	for n, mau := range map[int64]string{1: "100001", 421: "100421", 99999: "199999"} {
		if got, err := FormatIdentitas(n); err != nil || got != mau {
			t.Errorf("FormatIdentitas(%d) = %q, %v; mau %q", n, got, err, mau)
		}
	}
	for _, n := range []int64{100000, -1} {
		if _, err := FormatIdentitas(n); !errors.Is(err, ErrIdentitasMelampauiLebar) {
			t.Errorf("FormatIdentitas(%d): nomor tak muat 5 digit harus gagal terang, dapat %v", n, err)
		}
	}
}

func produkUji() models.Produk {
	return models.Produk{
		ID: "100044",
		Umum: models.ProdukUmum{ProductName: "UJI <PRODUK> & CO", Ceding: "UJI CEDING", CedingID: "L0UJI",
			SOBName: "UJI SOB", SOBID: "L0SOB", RIComm: "12.3456789012345678901234567890", RIRisk: "UJI RISK",
			RIRiskID: "1000117", InwardName: "UJI PRODUK UJI PEMEGANG", TreatyNumber: "UJI/001", Cause: "ANY CAUSE",
			CauseID: "100004", PolicyHolder: "UJI-ORG-1", PolicyHolderName: "UJI PEMEGANG", CreateOp: "UJI-A",
			UpdateOp: "UJI-B", IsORS: true, Comment: "UJI komentar"},
		Inward:                models.ProdukInward{Begin: "2026-03-01"},
		PlanList:              []models.BarisPlan{{Plan: "UJI PLAN", RIRate: "UJI RATE", RIRateID: "R1", Asli: `{"pxObjClass":"ASM-FW-GISFW-Data-Plan"}`}},
		FinancialUnderwriting: []models.BarisFinUW{{MinInsured: "1", NonEmployee: "N"}},
		DocumentClaim:         []models.BarisDokumen{{Document: "UJI DOK"}},
		OutwardList:           []models.BarisOutward{{ReinsTypeID: "10200"}},
	}
}

func TestRakitUmumKunciPegaDanKunciView(t *testing.T) {
	teks, err := RakitUmum(produkUji(), "", true)
	if err != nil {
		t.Fatal(err)
	}
	obj, err := uraiObjek(teks)
	if err != nil {
		t.Fatalf("JSON hasil rakit tidak sah: %v\n%s", err, teks)
	}
	for _, k := range KunciViewProduk {
		if _, ada := obj[k]; !ada {
			t.Errorf("kunci view PRODUCT_LIFE %q tidak ada di JSON hasil simpan", k)
		}
	}
	for _, w := range []string{`"POLICYHODER":"UJI-ORG-1"`, `"Non_Employee":"N"`, `"RICOMM":"12.3456789012345678901234567890"`,
		`"PRODUCTNAME":"UJI <PRODUK> & CO"`, `"IsORS":"true"`, `"ID":"100044"`, `"pxObjClass":"ASM-FW-GISFW-Data-Plan"`,
		`"OVR_COMM":""`, `"Document":"UJI DOK"`, `"CommentList":[]`} {
		if !strings.Contains(teks, w) {
			t.Errorf("JSON tanpa %s:\n%s", w, teks)
		}
	}
	if strings.Contains(teks, "u003c") || strings.Contains(teks, "u0026") {
		t.Errorf("teks diloloskan HTML - Pega menulis apa adanya:\n%s", teks)
	}
	if strings.Contains(teks, `"IsView"`) {
		t.Error("produk baru Pega tidak membawa IsView (NewProductLife tidak mengisinya)")
	}
}

func TestRakitUmumMempertahankanKunciLama(t *testing.T) {
	lama := `{"ID":"100044","KUNCILAMA":{"a":1},"IsORS":false,"IsView":"true","PRODUCTNAME":"LAMA"}`
	teks, err := RakitUmum(produkUji(), lama, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(teks, `"KUNCILAMA":{"a":1}`) || !strings.Contains(teks, `"PRODUCTNAME":"UJI <PRODUK> & CO"`) {
		t.Errorf("kunci tak dikelola dipertahankan, kunci dikelola ditimpa:\n%s", teks)
	}
	if !strings.Contains(teks, `"IsORS":true`) {
		t.Errorf("IsORS lama boolean → ditulis boolean:\n%s", teks)
	}
	if !strings.Contains(teks, `"IsView":"false"`) {
		t.Errorf("simpan sesudah Edit: IsView = \"false\" (SetViewEdit b145):\n%s", teks)
	}
}

func TestRakitUmumBolakBalik(t *testing.T) {
	p := produkUji()
	teks, err := RakitUmum(p, "", true)
	if err != nil {
		t.Fatal(err)
	}
	q, err := UraiProduk(p.ID, teks, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if q.Umum != p.Umum {
		t.Errorf("umum berubah menyeberang JSON:\n%+v\n%+v", q.Umum, p.Umum)
	}
	if q.PlanList[0] != p.PlanList[0] || q.FinancialUnderwriting[0] != p.FinancialUnderwriting[0] {
		t.Errorf("baris berubah: %+v %+v", q.PlanList, q.FinancialUnderwriting)
	}
}

func TestArgSimpanUmumMengikatClobDanKolomDatar(t *testing.T) {
	p := produkUji()
	args := argSisipUmum(p, `{"ID":"100044"}`)
	if len(args) != 6 {
		t.Fatalf("6 bind: %v", args)
	}
	if c, ok := args[1].(go_ora.Clob); !ok || !c.Valid || c.String != `{"ID":"100044"}` {
		t.Errorf("JSONDATA wajib diikat sebagai CLOB: %#v", args[1])
	}
	if args[2] != "1000117" || args[3] != "UJI RISK" || args[4] != "UJI <PRODUK> & CO" || args[5] != "01/03/2026" {
		t.Errorf("kolom datar RIRISKID, RIRISK, PRODUCTNAME, BEGIN_DATE: %v", args)
	}
	q := rata(sqlSisipUmum("S.M_PRODUCT_LIFE"))
	for _, w := range []string{"INSERT INTO S.M_PRODUCT_LIFE (ID, JSONDATA, RIRISKID, RIRISK, PRODUCTNAME, BEGIN_DATE)",
		"TO_DATE(:6, 'DD/MM/YYYY')"} {
		if !strings.Contains(q, w) {
			t.Errorf("SQL tanpa %q: %s", w, q)
		}
	}
	// ⛔ go-ora mengikat menurut URUTAN KEMUNCULAN: nomor placeholder = urutan argumen.
	u := rata(sqlPerbaruiUmum("S.M_PRODUCT_LIFE"))
	if !strings.Contains(u, "SET JSONDATA = :1, RIRISKID = :2, RIRISK = :3, PRODUCTNAME = :4, BEGIN_DATE = TO_DATE(:5, 'DD/MM/YYYY') WHERE ID = :6") {
		t.Errorf("UPDATE dikunci ID: %s", u)
	}
	ua := argPerbaruiUmum(p, `{}`)
	if ua[5] != "100044" {
		t.Errorf("ID adalah argumen terakhir UPDATE: %v", ua)
	}
}
