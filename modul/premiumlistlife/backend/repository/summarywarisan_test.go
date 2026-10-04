package repository

// PL-09 (GILIRAN-18) - `M_LIFE_PREMIUM_SUMMARY` seperti prosedur
// `PEGA_M_LIFE_PREMIUM_SUMMARY`. TANPA Oracle.

import (
	"context"
	"errors"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/premiumlistlife/backend/models"
)

// berkasProsedurSummary adalah badan prosedur yang dibaca asisten dari
// `ALL_SOURCE` DEV (commit 526fc93) - satu-satunya tempat nama kolomnya.
const berkasProsedurSummary = "../../docs/dba-procedure-PEGA_M_LIFE_PREMIUM_SUMMARY.md"

// badanProsedurSummary membaca blok SQL berkas itu tanpa nomor baris.
func badanProsedurSummary(t *testing.T) string {
	t.Helper()
	isi, err := os.ReadFile(berkasProsedurSummary)
	if err != nil {
		t.Fatalf("membaca badan prosedur: %v", err)
	}
	m := regexp.MustCompile("(?s)```sql\n(.*?)```").FindStringSubmatch(string(isi))
	if m == nil {
		t.Fatal("blok SQL badan prosedur tidak terbaca; pembacanya yang rusak")
	}
	nomor := regexp.MustCompile(`^\s*\d+ ?`)
	var baris []string
	for _, b := range strings.Split(m[1], "\n") {
		baris = append(baris, nomor.ReplaceAllString(b, ""))
	}
	return strings.Join(baris, "\n")
}

// pisahKoma memecah daftar SQL dipisah koma menjadi butir yang dipangkas.
func pisahKoma(s string) []string {
	var keluar []string
	for _, k := range strings.Split(s, ",") {
		if k = strings.TrimSpace(k); k != "" {
			keluar = append(keluar, k)
		}
	}
	return keluar
}

// TestKolomSummaryWarisanVERBATIMDariProsedur - dua cara, dari badan prosedur.
//
// Cara A: daftar kolom `INSERT INTO M_LIFE_PREMIUM_SUMMARY ( … )` - 38, `ID`
// pertama. Cara B: butir `VALUES` - 38, `ID` dari sequence lalu `P_<kolom>`
// untuk SETIAP kolom lain. Cara B yang membuktikan tiap parameter masuk ke
// kolom BERNAMA SAMA; cara A yang mengunci urutan kami posisi demi posisi.
func TestKolomSummaryWarisanVERBATIMDariProsedur(t *testing.T) {
	badan := badanProsedurSummary(t)
	m := regexp.MustCompile(`(?s)INSERT INTO M_LIFE_PREMIUM_SUMMARY\s*\((.*?)\)\s*VALUES\s*\((.*?)\);`).
		FindStringSubmatch(badan)
	if m == nil {
		t.Fatal("INSERT prosedur tidak terbaca; pembacanya yang rusak")
	}
	kolom, nilai := pisahKoma(m[1]), pisahKoma(m[2])
	if len(kolom) != 38 || len(nilai) != 38 {
		t.Fatalf("kolom %d, nilai %d; mau keduanya 38 (katalog DEV: 38 kolom)", len(kolom), len(nilai))
	}
	if kolom[0] != "ID" || nilai[0] != "id_LIFE_PREMIUM_SUMMARY_ins" {
		t.Errorf("butir pertama %s <- %s; mau ID <- sequence", kolom[0], nilai[0])
	}
	if !strings.Contains(badan, "TO_CHAR("+namaUrutanSummaryWarisan+".nextval)") {
		t.Errorf("badan prosedur tidak mengisi ID dari %s", namaUrutanSummaryWarisan)
	}
	for i := 1; i < len(kolom); i++ {
		if nilai[i] != "P_"+kolom[i] {
			t.Errorf("posisi %d: kolom %s diisi %s, bukan parameter bernama sama", i, kolom[i], nilai[i])
		}
	}
	kami := append([]string{"ID"}, kolomSummaryWarisan...)
	if len(kami) != len(kolom) {
		t.Fatalf("pemetaan kami %d kolom, prosedur %d", len(kami), len(kolom))
	}
	for i := range kolom {
		if kolom[i] != kami[i] {
			t.Errorf("posisi %d: prosedur %s, kami %s", i, kolom[i], kami[i])
		}
	}
	// ⛔ Satu-satunya pernyataan DML prosedur itu INSERT - tidak ada cabang
	// UPDATE/MERGE yang harus ditiru. Bila kelak muncul, uji ini merah dulu.
	for _, kata := range []string{"UPDATE ", "MERGE ", "DELETE "} {
		if strings.Contains(strings.ToUpper(badan), kata) {
			t.Errorf("badan prosedur memuat %q; tiruan kami hanya punya cabang INSERT", kata)
		}
	}
}

// TestSumberSummaryWarisanDariKorpus - parameter posisional `InsertPLSummary`
// dan `TempInputData.CARIn` `InsertJsonPolisLife_Act` langkah 8.1 (b2296…).
//
// ⛔ Yang dibuktikan: parameter `P_X` menerima `.X` baris `CurrencyList` - nama
// yang SAMA dengan kolomnya - sehingga `nilaiSummaryWarisan` cukup mengambil
// kolom rekap bernama sama. Tiga turunan (`PREMIUM`, `COMMISSION`, `BALANCE`)
// dan 29 jumlah; `CURRENCY` dari baris itu; empat sisanya dari halaman kerja.
func TestSumberSummaryWarisanDariKorpus(t *testing.T) {
	const rdb = `D:\XML\RNM_BRD\PremiumList Life\RDBList\InsertPLSummary.xml`
	const act = `D:\XML\RNM_BRD\PremiumList Life\Activity\InsertJsonPolisLife_Act.xml`
	isiRDB, err := os.ReadFile(rdb)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v); sumber tidak terperiksa", err)
	}
	isiAct, err := os.ReadFile(act)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v); sumber tidak terperiksa", err)
	}
	m := regexp.MustCompile(`(?s)PEGA_M_LIFE_PREMIUM_SUMMARY\((.*?)\{InputParam\.ERRMSG`).FindStringSubmatch(string(isiRDB))
	if m == nil {
		t.Fatal("panggilan InsertPLSummary tidak terbaca; pembacanya yang rusak")
	}
	var arg []string
	for _, a := range pisahKoma(m[1]) {
		arg = append(arg, strings.Trim(a, "{}"))
	}
	param := regexp.MustCompile(`(P_\w+) IN VARCHAR2`).FindAllStringSubmatch(badanProsedurSummary(t), -1)
	if len(arg) != 37 || len(param) != 37 {
		t.Fatalf("argumen %d, parameter IN %d; mau keduanya 37", len(arg), len(param))
	}
	halaman := map[string]string{
		"P_COB":           "pyWorkPage.BusinessName",
		"P_PL_NUMBER":     "pyWorkPage.PremiumListSummary.PL_NUMBER",
		"P_PL_NUMBER_EDM": "pyWorkPage.PremiumListSummary.PL_NUMBER_EDM",
		"P_IDPEGA":        "pyWorkPage.pzInsKey",
	}
	cari := regexp.MustCompile(`<PropertiesName>TempInputData\.(CARI\d+)</PropertiesName>\s*<PropertiesValue>([^<]*)</PropertiesValue>`).
		FindAllStringSubmatch(string(isiAct), -1)
	sumberCari := map[string]string{}
	for _, c := range cari {
		sumberCari[c[1]] = c[2]
	}
	jumlah := map[string]bool{}
	for _, k := range models.KolomJumlahSummary {
		jumlah[k] = true
	}
	for i, p := range param {
		nama, a := p[1], arg[i]
		if mau, ada := halaman[nama]; ada {
			if a != mau {
				t.Errorf("%s menerima %s, mau %s", nama, a, mau)
			}
			continue
		}
		kode, ok := strings.CutPrefix(a, "TempInputData.")
		if !ok {
			t.Errorf("%s menerima %s, bukan TempInputData.CARIn", nama, a)
			continue
		}
		kolom := strings.TrimPrefix(nama, "P_")
		if got := sumberCari[kode]; got != "."+kolom {
			t.Errorf("%s <- %s = %q, mau .%s", nama, kode, got, kolom)
		}
		switch kolom {
		case "CURRENCY", "PREMIUM", "COMMISSION", "BALANCE":
		default:
			if !jumlah[kolom] {
				t.Errorf("%s bukan kolom KolomJumlahSummary; rekap tidak menghitungnya", kolom)
			}
		}
	}
}

// rekapUjiPL09 - satu rekap USD, sebagian jumlah absen (nil).
func rekapUjiPL09() models.RekapMataUang {
	return models.RekapMataUang{
		Currency:   "USD",
		Premium:    apd.New(9750000, -4),
		Commission: apd.New(1250, -4),
		Balance:    apd.New(-8970000, -4),
		Jumlah: map[string]*apd.Decimal{
			"CLAIM":     apd.New(10, 0),
			"DEDUCTION": apd.New(333, -2),
		},
	}
}

// TestNilaiSummaryWarisanSejajarDanNolBilaKosong - PL-10 di tabel summary.
//
// ⛔ Uang kosong = "0" (`@toDecimal`-nya Pega, keputusan PL-10), BUKAN NULL -
// sama dengan salinan detail. Teks kosong (`PL_NUMBER_EDM` new business) =
// NULL, sama dengan `”` di Oracle.
func TestNilaiSummaryWarisanSejajarDanNolBilaKosong(t *testing.T) {
	kepala := KepalaSummaryWarisan{NomorPL: "UJI-PL-0918", COB: "UJI-COB", IDPega: "NBLF-918"}
	arg := nilaiSummaryWarisan(kepala, rekapUjiPL09())
	if len(arg) != len(kolomSummaryWarisan) {
		t.Fatalf("argumen %d, kolom %d", len(arg), len(kolomSummaryWarisan))
	}
	nilai := map[string]any{}
	for i, k := range kolomSummaryWarisan {
		nilai[k] = arg[i]
	}
	mau := map[string]any{
		"COB": "UJI-COB", "CURRENCY": "USD", "PL_NUMBER": "UJI-PL-0918",
		"PL_NUMBER_EDM": nil, "IDPEGA": "NBLF-918",
		"PREMIUM": "975.0000", "COMMISSION": "0.1250", "BALANCE": "-897.0000",
		"CLAIM": "10", "DEDUCTION": "3.33",
		"BROKERAGE_FEE": "0", "RI_ADMIN_FEE": "0", "NET_PREMIUM_REFUND_RETRO": "0",
	}
	for k, v := range mau {
		if nilai[k] != v {
			t.Errorf("%s = %#v, mau %#v", k, nilai[k], v)
		}
	}
	for _, k := range kolomSummaryWarisan {
		if nilai[k] == nil && k != "PL_NUMBER_EDM" {
			t.Errorf("%s NULL; hanya PL_NUMBER_EDM yang boleh kosong", k)
		}
	}
	kepala.NomorEDM = "UJI-EDM-1"
	if got := nilaiSummaryWarisan(kepala, rekapUjiPL09()); got[indeksKolomSummary(t, "PL_NUMBER_EDM")] != "UJI-EDM-1" {
		t.Errorf("PL_NUMBER_EDM terisi tidak tersalin: %#v", got)
	}
}

func indeksKolomSummary(t *testing.T, nama string) int {
	t.Helper()
	for i, k := range kolomSummaryWarisan {
		if k == nama {
			return i
		}
	}
	t.Fatalf("kolom %s tidak ada", nama)
	return -1
}

// TestSisipSummaryWarisanBerurutanDanBerSequence - ID dari sequence warisan.
func TestSisipSummaryWarisanBerurutanDanBerSequence(t *testing.T) {
	q := sqlSisipSummaryWarisan("S.M", "S.Q")
	if !strings.HasPrefix(q, "INSERT INTO S.M (ID, BALANCE, BROKERAGE_FEE, CLAIM, OVR_COMM, COB, ") {
		t.Errorf("kepala INSERT tidak urut prosedur:\n%s", q)
	}
	if !strings.Contains(q, "VALUES (TO_CHAR(S.Q.NEXTVAL), :1, :2, ") || !strings.HasSuffix(q, ", :37)") {
		t.Errorf("VALUES bukan sequence + 37 penanda:\n%s", q)
	}
}

// TestHapusSummaryWarisanDikurungNomorDanWork - idempoten pl2 tanpa memakan
// baris EDM milik work lain yang ber-`PL_NUMBER` sama.
//
// ⚠️ PENYIMPANGAN SADAR dari prosedur, yang hanya INSERT: submit ulang di
// Pega menumpuk baris summary kembar. pl2 menuntut idempoten per `PL_NUMBER`.
func TestHapusSummaryWarisanDikurungNomorDanWork(t *testing.T) {
	q := sqlHapusSummaryWarisan("S.M")
	if q != "DELETE FROM S.M WHERE PL_NUMBER = :1 AND IDPEGA = :2" {
		t.Errorf("penghapusan summary warisan tidak dikurung nomor + work:\n%s", q)
	}
}

// TestKepalaSummaryWarisanAdaDiMigrasi051 - `COB` <- `BusinessName`.
func TestKepalaSummaryWarisanAdaDiMigrasi051(t *testing.T) {
	q := sqlKepalaSummaryWarisan("S.P")
	if q != "SELECT BUSINESS_NAME, PL_NUMBER_EDM FROM S.P WHERE ID = :1" {
		t.Errorf("pembaca kepala berubah:\n%s", q)
	}
	ada := map[string]bool{}
	for _, k := range kolomMigrasi(t, "051_t_premium_list.sql", "T_PREMIUM_LIST") {
		ada[k] = true
	}
	for _, k := range []string{"BUSINESS_NAME", "PL_NUMBER_EDM"} {
		if !ada[k] {
			t.Errorf("%s tidak ada di migrasi 051", k)
		}
	}
}

// TestGantiSummaryWarisanMenolakKepalaKosong - sebelum menyentuh basis data.
func TestGantiSummaryWarisanMenolakKepalaKosong(t *testing.T) {
	ctx := context.Background()
	r := NewSummaryWarisan(nil)
	rekap := []models.RekapMataUang{rekapUjiPL09()}
	for nama, k := range map[string]KepalaSummaryWarisan{
		"tanpa nomor": {IDPega: "NBLF-1"},
		"tanpa work":  {NomorPL: "UJI-PL-1"},
	} {
		if _, _, err := r.Ganti(ctx, nil, k, rekap); err == nil {
			t.Errorf("%s: diterima", nama)
		}
	}
	if _, _, err := r.Ganti(ctx, nil, KepalaSummaryWarisan{NomorPL: "UJI-PL-1", IDPega: "NBLF-1"}, nil); !errors.Is(err, ErrRekapKosong) {
		t.Errorf("rekap kosong: %v, mau ErrRekapKosong", err)
	}
}

// TestSetiapKolomUangSummaryPunyaSumberRekap - TANPA korpus (temuan
// /code-review GILIRAN-18).
//
// ⛔ Cabang bawaan `nilaiSummaryWarisan` membaca `r.Jumlah[kolom]`, dan kolom
// yang tidak dijumlah rekap menjadi "0" TANPA galat - "nol diam-diam" yang
// diperingatkan `models`. `TestSumberSummaryWarisanDariKorpus` menagihnya
// juga, tetapi ia SKIP di mesin tanpa korpus; uji ini tidak.
func TestSetiapKolomUangSummaryPunyaSumberRekap(t *testing.T) {
	jumlah := map[string]bool{}
	for _, k := range models.KolomJumlahSummary {
		jumlah[k] = true
	}
	uang, teks := 0, 0
	for _, k := range kolomSummaryWarisan {
		if !KolomUangSummaryWarisan(k) {
			teks++
			continue
		}
		uang++
		switch k {
		case "PREMIUM", "COMMISSION", "BALANCE":
			continue
		}
		if !jumlah[k] {
			t.Errorf("kolom uang %s tidak dijumlah models.KolomJumlahSummary; ia akan selalu 0", k)
		}
	}
	if uang != 32 || teks != 5 {
		t.Errorf("uang %d, teks %d; mau 32 + 5 (badan prosedur)", uang, teks)
	}
	// Kelima kolom teks dipetakan cabang tersendiri, bukan jatuh ke cabang uang.
	kepala := KepalaSummaryWarisan{NomorPL: "UJI-PL-1", NomorEDM: "UJI-EDM-1", COB: "UJI-COB", IDPega: "UJI-W-1"}
	arg := nilaiSummaryWarisan(kepala, rekapUjiPL09())
	for i, k := range kolomSummaryWarisan {
		if !KolomUangSummaryWarisan(k) && arg[i] == "0" {
			t.Errorf("kolom teks %s jatuh ke cabang uang", k)
		}
	}
}
