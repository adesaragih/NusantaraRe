package services_test

// Aksi tambahan dari XML (paket 9): `Copy`, `On Retention` → `OutwardList`, `Generate`.

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/services"
	"nusantarare/modul/masterproductnamelife/backend/tiruan"
)

const jsonAsal = `{"ID":"100007","PRODUCTNAME":"UJI ASAL","TYPE":"1","GRUP":"2","PRODUCTCODE":"UJI-KODE",` +
	`"CREATEOP":"UJI-ASAL","IsORS":true,` +
	`"CommentList":[{"Date":"20260101T000000.000 GMT","OperatorName":"UJI-ASAL","IsApproved":"","Suggest":"UJI KOMENTAR ASAL"}],` +
	`"OutwardList":[{"REINSTYPEID":"10200","REINSTYPENAME":"UJI OR","TRANSACTIONYEAR":"2025","TREATYCONTRACTID":"","UNDERWRITINGYEAR":"2025","OVR_COMM":""}]}`

const jsonAsalInward = `{"ID":"100007","PRODUCTID":"100007","MONTHS":"UJI-BULAN","LIENCLAUSE":"UJI-LIEN"}`

func layananSalin() (*services.Layanan, *tiruan.Gudang) {
	l, g := layananMaster()
	g.Umum["100007"] = jsonAsal
	g.Inward["100007"] = jsonAsalInward
	return l, g
}

func TestCopyProdukBaruMewarisiMedanServerPembuatPelaku(t *testing.T) {
	l, g := layananSalin()
	m := produkMasuk()
	m.SalinanDari = "100007"
	m.Umum.IsORS = true
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != "100044" || p.Inward.ID != "100044" || p.Inward.ProductID != "100044" {
		t.Errorf("salinan = produk baru, kedua ID dari sequence (CopyProduct b144/b173): %q %q %q", p.ID, p.Inward.ID, p.Inward.ProductID)
	}
	if p.Umum.CreateOp != "UJI-PELAKU" {
		t.Errorf("CREATEOP salinan = pelaku (OQ-MPNL-13): %q", p.Umum.CreateOp)
	}
	if p.Umum.TypeBasicRider != "1" || p.Umum.Grup != "2" || p.Umum.ProductCode != "UJI-KODE" ||
		p.Inward.Months != "UJI-BULAN" || p.Inward.LienClause != "UJI-LIEN" {
		t.Errorf("medan mati ikut tersalin seperti halaman Pega: %+v %+v", p.Umum, p.Inward)
	}
	if len(p.CommentList) != 2 || p.CommentList[0].Suggest != "UJI KOMENTAR ASAL" {
		t.Errorf("riwayat komentar asal + satu baris simpan: %+v", p.CommentList)
	}
	if len(p.OutwardList) != 1 || p.OutwardList[0].ReinsTypeName != "UJI OR" {
		t.Errorf("OutwardList asal ikut tersalin (checkbox tidak diubah): %+v", p.OutwardList)
	}
	if g.MintaOR != [2]string{} {
		t.Errorf("salinan tanpa perubahan checkbox tidak membaca kontrak: %v", g.MintaOR)
	}
	if g.Umum["100007"] != jsonAsal {
		t.Error("produk asal tidak tersentuh")
	}
	if strings.Contains(g.Umum["100044"], "salinanDari") || strings.Contains(g.Umum["100044"], "hitungOutward") {
		t.Error("penanda permintaan tidak pernah disimpan")
	}
}

func TestCopyGalat(t *testing.T) {
	l, g := layananSalin()
	m := produkMasuk()
	m.SalinanDari = "100999"
	if _, err := l.SimpanProduk(context.Background(), pelakuUji, m, true); !errors.Is(err, services.ErrSalinanTidakAda) ||
		!strings.Contains(err.Error(), "100999") {
		t.Errorf("produk asal tidak ada: %v", err)
	}
	m.SalinanDari, m.ID = "100007", "100007"
	if _, err := l.SimpanProduk(context.Background(), pelakuUji, m, false); !errors.Is(err, services.ErrSalinanPadaUbah) {
		t.Errorf("salinanDari pada ubah ditolak: %v", err)
	}
	if g.Komit != 0 {
		t.Errorf("nol transaksi sukses: %d", g.Komit)
	}
}

func kontrakUji(g *tiruan.Gudang) {
	tgl := func(s string) time.Time { v, _ := time.Parse("2006-01-02", s); return v }
	g.KontrakOR = []tiruan.KontrakOR{
		{ReinsTypeID: "10200", ReinsTypeName: "UJI OR 2025", TreatyYear: "2025", UnderwritingYear: "2025",
			Mulai: tgl("2025-01-01"), Akhir: tgl("2036-12-31")},
		{ReinsTypeID: "10196", ReinsTypeName: "UJI QS", TreatyYear: "2025", UnderwritingYear: "2025",
			Mulai: tgl("2025-01-01"), Akhir: tgl("2036-12-31")},
		{ReinsTypeID: "10200", ReinsTypeName: "UJI OR 2018", TreatyYear: "2018", UnderwritingYear: "2018",
			Mulai: tgl("2018-01-01"), Akhir: tgl("2024-12-31")},
	}
}

func TestOnRetentionDiubahMenggantiOutwardList(t *testing.T) {
	l, g := layananSalin()
	kontrakUji(g)
	ctx := context.Background()
	m, err := l.AmbilProduk(ctx, pelakuUji, "100007")
	if err != nil {
		t.Fatal(err)
	}
	isi := produkMasuk()
	isi.ID, isi.Inward.Begin, isi.Inward.Mature = "100007", "2026-03-01", "2027-02-28"
	isi.OutwardList = []models.BarisOutward{{ReinsTypeID: "KLIEN"}}
	// Checkbox tidak diubah: OutwardList tersimpan dipertahankan, kiriman klien diabaikan.
	p, err := l.SimpanProduk(ctx, pelakuUji, isi, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.OutwardList) != 1 || p.OutwardList[0].ReinsTypeName != m.OutwardList[0].ReinsTypeName {
		t.Errorf("tanpa perubahan checkbox: %+v", p.OutwardList)
	}
	// Checkbox diubah - JUGA saat dilepas: keempat prakondisi PRE=false (`GetReinsTypeOR_Life` b236, b534, b731).
	isi.HitungOutward, isi.Umum.IsORS = true, false
	p, err = l.SimpanProduk(ctx, pelakuUji, isi, false)
	if err != nil {
		t.Fatal(err)
	}
	if g.MintaOR != [2]string{"01/03/2026", "28/02/2027"} {
		t.Errorf("Temp.CARI1/2 = BEGIN/MATURE dd/MM/yyyy (`GetReinsTypeOR_Life` 2 b383): %v", g.MintaOR)
	}
	want := models.BarisOutward{ReinsTypeID: "10200", ReinsTypeName: "UJI OR 2025", TransactionYear: "2025", UnderwritingYear: "2025"}
	if len(p.OutwardList) != 1 || p.OutwardList[0].ReinsTypeID != want.ReinsTypeID || p.OutwardList[0].ReinsTypeName != want.ReinsTypeName ||
		p.OutwardList[0].TransactionYear != want.TransactionYear || p.OutwardList[0].UnderwritingYear != want.UnderwritingYear ||
		p.OutwardList[0].TreatyContractID != "" || p.OutwardList[0].OvrComm != "" {
		t.Errorf("OutwardList dari BrowseReinstypeOR_SQL (`GetReinsTypeOR_Life` 4.1 b770): %+v", p.OutwardList)
	}
	if !strings.Contains(g.Umum["100007"], `"OutwardList":[{"OVR_COMM":"","REINSTYPEID":"10200"`) {
		t.Errorf("kunci Pega + OVR_COMM di JSON: %s", g.Umum["100007"])
	}
	// Tanggal kosong: TO_DATE NULL - nol baris.
	isi.Inward.Mature = ""
	if p, err = l.SimpanProduk(ctx, pelakuUji, isi, false); err != nil || len(p.OutwardList) != 0 {
		t.Errorf("MATURE kosong = OutwardList kosong: %+v %v", p.OutwardList, err)
	}
}

func TestProdukBaruOnRetentionDihitung(t *testing.T) {
	l, g := layananMaster()
	kontrakUji(g)
	m := produkMasuk()
	m.Umum.IsORS, m.Inward.Mature = true, "2027-02-28"
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.OutwardList) != 1 || !p.Umum.IsORS {
		t.Errorf("produk baru bercentang: %+v", p.OutwardList)
	}
	g.GagalMaster = errors.New("UJI kontrak tak terbaca")
	if _, err := l.SimpanProduk(context.Background(), pelakuUji, m, true); err == nil {
		t.Error("kontrak tak terbaca = simpan gagal terang, bukan OutwardList kosong diam-diam")
	}
}

// judulGenerate - `GenerateUpload_Act` b1340 `CSVPropHeaders` VERBATIM.
const judulGenerate = "TYPE,TYPE_CEDING,CEDING,GRUP,PRODUCTNAME,PRODUCTCODE,PRODUCTTYPE,RIRATE,RICOMM,RIRISK,INWARDNAME," +
	"BENEFIT,CAUSE,POLICYHODERNAME,INSURED,ADDENDUMWORD,AMANDEMENTSCHD,BEGIN,MATURE,BIRTHDAY,CURRENCY,EXTRAMORTALITY," +
	"MAXCONTRACT,CEDINGRETENTIONNUM,CEDINGLIMIT,BROKERAGE,MINAGE,MAXAGE,EXTRAPREMI,MONTHS,MINSUMINSURED,MAXSUMINSURED," +
	"MAXSUMREASURED,RNMSHARE,RNMLIMITNUM,PAYMENT,PROPORTIONALTABLE,SUBJECTTO,TREATYNUMBER"

func generate(t *testing.T, m models.Produk) map[string]string {
	t.Helper()
	l, _ := layananUji()
	var buf bytes.Buffer
	if err := l.Generate(context.Background(), pelakuUji, m, &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "\r\n") {
		t.Error("CSV ber-CRLF")
	}
	r, err := csv.NewReader(&buf).ReadAll()
	if err != nil || len(r) != 2 {
		t.Fatalf("satu judul + satu baris: %v %v", r, err)
	}
	if strings.Join(r[0], ",") != judulGenerate {
		t.Errorf("judul VERBATIM b1340:\n%s", strings.Join(r[0], ","))
	}
	hasil := map[string]string{}
	for i, j := range r[0] {
		hasil[j] = r[1][i]
	}
	return hasil
}

func TestGenerateSeeDetailMenurutNamaJudul(t *testing.T) {
	m := produkMasuk()
	m.Umum.TypeBasicRider, m.Umum.TypeCeding, m.Umum.Grup, m.Umum.Benefit = "1", "1", "3", "UJI BENEFIT"
	m.Inward.Payment, m.Inward.LienClause, m.Inward.TreatyNumber, m.Inward.Mature = "2", "UJI-LIEN", "UJI/IN", "2027-02-28"
	m.Inward.MaxSumReasured = "1234567890.12"
	v := generate(t, m)
	cek := map[string]string{"TYPE": "Basic", "TYPE_CEDING": "QS", "GRUP": "PA", "PAYMENT": "Semi Annual",
		"CEDING": "UJI CEDING", "BENEFIT": "UJI BENEFIT", "POLICYHODERNAME": "UJI PEMEGANG", "BEGIN": "01/03/2026",
		"MATURE": "28/02/2027", "TREATYNUMBER": "UJI/IN", "RICOMM": "12,5", "MAXSUMREASURED": "1234567890.12"}
	for k, w := range cek {
		if v[k] != w {
			t.Errorf("%s = %q, mau %q", k, v[k], w)
		}
	}
	for _, x := range v {
		if x == "UJI-LIEN" {
			t.Error("CARI36 LIENCLAUSE tanpa judul tidak ikut (OQ-MPNL-06)")
		}
	}
	// `@if` lainnya: bukan 1 → Rider / SUPRLUS / HEALTH / Single.
	m.Umum.TypeBasicRider, m.Umum.TypeCeding, m.Umum.Grup, m.Inward.Payment = "", "2", "", "9"
	v = generate(t, m)
	if v["TYPE"] != "Rider" || v["TYPE_CEDING"] != "SUPRLUS" || v["GRUP"] != "HEALTH" || v["PAYMENT"] != "Single" {
		t.Errorf("cabang lainnya @if: %v", v)
	}
	l, _ := layananUji()
	if err := l.Generate(context.Background(), inti.Pelaku{}, m, &bytes.Buffer{}); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Error("Generate menuntut identitas")
	}
}

// TestJudulGenerateSamaDenganKorpus - b1340 dibaca dari korpus READ-ONLY; korpus
// tidak terjangkau = DILEWATI.
func TestJudulGenerateSamaDenganKorpus(t *testing.T) {
	berkas := filepath.Join(`D:\XML\RNM_BRD\Master Product Name Life`, "Activity", "GenerateUpload_Act.xml")
	isi, err := os.ReadFile(berkas)
	if err != nil {
		t.Skipf("korpus tidak terjangkau: %v", err)
	}
	baris := strings.Split(strings.ReplaceAll(string(isi), "><", ">\n<"), "\n")
	if len(baris) < 1343 {
		t.Fatalf("korpus pendek: %d baris", len(baris))
	}
	m := regexp.MustCompile(`<CSVPropHeaders>([^<]*)</CSVPropHeaders>`).FindStringSubmatch(baris[1339])
	if m == nil || m[1] != judulGenerate {
		t.Errorf("b1340 korpus: %q", baris[1339])
	}
	if !strings.Contains(baris[1334], "<FileName>SeeDetail</FileName>") {
		t.Errorf("b1335 korpus: %q", baris[1334])
	}
}
