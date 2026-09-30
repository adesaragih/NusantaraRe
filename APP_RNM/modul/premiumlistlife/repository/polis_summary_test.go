package repository

// Rekap uang dan salinan warisan - tiket 05a bagian 2. TANPA Oracle.

import (
	"database/sql"
	"os"
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/migrasi"
	"nusantarare/inti/utils"
	"nusantarare/modul/premiumlistlife/models"
)

// kolomMigrasi membaca kolom satu tabel dari satu berkas migrasi.
func kolomMigrasi(t *testing.T, berkas, tabel string) []string {
	t.Helper()
	isi, err := berkasMigrasi.ReadFile("migrations/" + berkas)
	if err != nil {
		t.Fatalf("membaca migrasi %s: %v", berkas, err)
	}
	var keluar []string
	for _, pernyataan := range strings.Split(string(isi), "\n/") {
		nama, kolom := migrasi.KolomCreateTable(pernyataan)
		nama = strings.ToUpper(nama)
		if i := strings.LastIndex(nama, "."); i >= 0 {
			nama = nama[i+1:]
		}
		if nama != tabel {
			continue
		}
		for _, k := range kolom {
			keluar = append(keluar, strings.ToUpper(k))
		}
	}
	if len(keluar) == 0 {
		t.Fatalf("nol kolom %s di %s; pembacanya yang rusak", tabel, berkas)
	}
	return keluar
}

// TestKolomRekapSamaDenganMigrasi055 - dua arah, supaya tidak ada kolom yatim.
//
// ⛔ Kolom rekap yang tidak ada di DDL gagal di Oracle (ORA-00904); kolom
// DDL yang tidak pernah diisi adalah rekap yang diam-diam NULL.
func TestKolomRekapSamaDenganMigrasi055(t *testing.T) {
	ddl := map[string]bool{}
	for _, k := range kolomMigrasi(t, "055_t_premium_list_summary.sql", "T_PREMIUM_LIST_SUMMARY") {
		ddl[k] = true
	}
	isi := map[string]bool{}
	for _, k := range kolomSisipRekap() {
		if isi[k] {
			t.Errorf("kolom rekap %s disisipkan dua kali", k)
		}
		isi[k] = true
		if !ddl[k] {
			t.Errorf("kolom rekap %s tidak ada di migrasi 055", k)
		}
	}
	for k := range ddl {
		if !isi[k] {
			t.Errorf("kolom 055 %s tidak pernah diisi rekap", k)
		}
	}
}

// TestKolomBacaSummaryAdaDiMigrasi052 - pembaca tidak meminta kolom hantu.
func TestKolomBacaSummaryAdaDiMigrasi052(t *testing.T) {
	ada := kolomTabelPeserta(t)
	for _, k := range models.KolomBacaSummary() {
		if !ada[k] {
			t.Errorf("kolom baca rekap %s tidak ada di migrasi 052", k)
		}
	}
	q := sqlBarisUangPolis("SKEMAUJI.T_PREMIUM_LIST_DETAIL")
	if !strings.Contains(q, "d.PREMIUM_LIST_ID = :1") {
		t.Errorf("pembaca uang tidak dikurung satu polis:\n%s", q)
	}
	// ⛔ Uang lewat TO_CHAR ber-NLS, tidak pernah lewat float.
	if n := strings.Count(q, "NLS_NUMERIC_CHARACTERS"); n != len(models.KolomBacaSummary()) {
		t.Errorf("%d kolom ber-TO_CHAR NLS, mau %d", n, len(models.KolomBacaSummary()))
	}
}

// TestNilaiSisipRekapSejajarDanTepat - argumen sejajar, teks uang apa adanya.
//
// ⛔ Literal 975 dari uji rumus (`TestBalanceKeempatCabangDariLiteral`)
// dibawa sampai ke argumen SQL: yang dikirim `975.0000`, bukan `975` dan
// bukan float.
func TestNilaiSisipRekapSejajarDanTepat(t *testing.T) {
	baris := []models.BarisUang{{}}
	for k, v := range map[string]string{
		"GROSS_PREMIUM": "1000", "DEDUCTION": "10", "RI_ADMIN_FEE": "1",
		"BROKERAGE_FEE": "2", "TAX": "3", "PROF_COMM": "4", "CLAIM": "5",
		"COMM": "-0.00004",
	} {
		d, err := utils.ParseDecimal(v)
		if err != nil {
			t.Fatal(err)
		}
		baris[0][k] = d
	}
	rekap, err := models.RekapPerMataUang("QR", baris, []string{"IDR"})
	if err != nil {
		t.Fatal(err)
	}
	arg := nilaiSisipRekap("POLIS1", rekap[0])
	kolom := kolomSisipRekap()
	if len(arg) != len(kolom) {
		t.Fatalf("%d argumen, %d kolom", len(arg), len(kolom))
	}
	if n := strings.Count(sqlSisipRekap("S.T"), ":"); n != len(kolom) {
		t.Errorf("%d penanda, %d kolom", n, len(kolom))
	}
	nilai := map[string]any{}
	for i, k := range kolom {
		nilai[k] = arg[i]
	}
	for k, mau := range map[string]string{
		"BALANCE": "975.0000", "PREMIUM": "1000.0000", "DEDUCTION": "10.0000",
		"CURRENCY": "IDR", "PREMIUM_LIST_ID": "POLIS1",
		// ⚠️ Negatif kecil dibulatkan, tanda dipertahankan - bukan dibuang.
		"COMMISSION": "0.0000",
	} {
		if s, _ := nilai[k].(string); s != mau && !(k == "COMMISSION" && s == "-0.0000") {
			t.Errorf("%s = %v, mau %s", k, nilai[k], mau)
		}
	}
	if nilai["ID"] != PengenalRekap("POLIS1", "IDR") {
		t.Errorf("ID rekap %v bukan PengenalRekap", nilai["ID"])
	}
}

// TestPengenalRekapTepat32DanDeterministik - kolom ID `VARCHAR2(32)`.
func TestPengenalRekapTepat32DanDeterministik(t *testing.T) {
	a, b := PengenalRekap("P", "IDR"), PengenalRekap("P", "IDR")
	if len(a) != 32 || a != b {
		t.Errorf("pengenal %q/%q", a, b)
	}
	if PengenalRekap("P", "USD") == a {
		t.Error("dua mata uang satu polis berpengenal sama")
	}
	if PengenalRekap("P", "IDR") == PengenalPesertaUnggah("P", 1) {
		t.Error("pengenal rekap bertabrakan dengan pengenal peserta")
	}
}

// TestNolCommitDiQueryRekapDanWarisan - ADR-U-0029.
//
// ⛔ `SaveMasterLPDet` baris 252 memuat `COMMIT;`. Salinan kami tidak.
func TestNolCommitDiQueryRekapDanWarisan(t *testing.T) {
	for nama, q := range map[string]string{
		"baca uang":     sqlBarisUangPolis("S.D"),
		"hapus rekap":   sqlHapusRekap("S.S"),
		"sisip rekap":   sqlSisipRekap("S.S"),
		"sumber":        sqlSumberWarisan("S.D", "S.P"),
		"hapus warisan": sqlHapusPesertaWarisan("S.M"),
		"sisip warisan": sqlSisipPesertaWarisan("S.M", "S.Q"),
	} {
		if strings.Contains(strings.ToUpper(q), "COMMIT") {
			t.Errorf("query %q memuat COMMIT:\n%s", nama, q)
		}
	}
}

// TestPenulisWarisanTanpaKueriBaca - batas yang membuat penjaga 66,8 juta
// baris tetap tajam.
//
// ⛔ polis_warisan.go menyebut `namaTabelPeserta`, jadi SETIAP kueri baca di
// dalamnya ditagih penyaring ber-index dan batas hasil. Bahan salinannya
// dibaca dari tabel kami di polis_summary.go. Bila kelak ada kueri baca
// ditambahkan ke sana, uji ini merah lebih dulu - dan yang menambahkannya
// harus memutuskan dengan sadar, bukan kebetulan lolos.
func TestPenulisWarisanTanpaKueriBaca(t *testing.T) {
	isi, err := os.ReadFile("polis_warisan.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.ToUpper(string(isi)), "SELECT") {
		t.Error("polis_warisan.go memuat SELECT; ia hanya boleh menulis")
	}
	sumber, err := os.ReadFile("polis_summary.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(sumber), "namaTabelPeserta") {
		t.Error("polis_summary.go menyebut tabel warisan; pembaca sumber harus membaca tabel kami")
	}
}

// TestHapusWarisanDikurungNomorDanWork - idempoten pl2 tanpa memakan EDM.
func TestHapusWarisanDikurungNomorDanWork(t *testing.T) {
	q := sqlHapusPesertaWarisan("S.M")
	if !strings.Contains(q, "WHERE PL_NUMBER = :1 AND IDPEGA = :2") {
		t.Errorf("penghapusan warisan tidak dikurung nomor + work:\n%s", q)
	}
}

// TestKolomWarisanVERBATIMDariSaveMasterLPDet - dibaca dari korpus, dua cara.
//
// Cara A: daftar kolom `INSERT INTO … ( … )`. Cara B: cacah butir `VALUES`.
// Keduanya harus 80, dan cara A harus sama URUT dengan `ID` +
// `kolomPesertaWarisan`. Kolom yang tergeser satu posisi menyimpan nilai
// sebuah kolom ke kolom tetangganya - dan tipenya sering sama.
func TestKolomWarisanVERBATIMDariSaveMasterLPDet(t *testing.T) {
	const letak = `D:\XML\RNM_BRD\PremiumList Life\RDBList\SaveMasterLPDet.xml`
	isi, err := os.ReadFile(letak)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v); pemetaan tidak terperiksa", err)
	}
	teks := string(isi)
	m := regexp.MustCompile(`(?s)INSERT INTO POOLDATA\.M_LIFE_PREMIUM_DETAIL\s*\((.*?)\)\s*VALUES\s*\((.*?)\);`).
		FindStringSubmatch(teks)
	if m == nil {
		t.Fatal("INSERT SaveMasterLPDet tidak terbaca; pembacanya yang rusak")
	}
	var caraA []string
	for _, k := range strings.Split(m[1], ",") {
		if k = strings.TrimSpace(k); k != "" {
			caraA = append(caraA, k)
		}
	}
	// Cara B: butir VALUES dipisah koma di tingkat teratas - koma di dalam
	// `To_date(…, 'DD/MM/YYYY')` dan `to_Char(…)` tidak dihitung.
	tingkat, caraB := 0, 1
	for _, c := range m[2] {
		switch c {
		case '(':
			tingkat++
		case ')':
			tingkat--
		case ',':
			if tingkat == 0 {
				caraB++
			}
		}
	}
	if n := strings.Count(strings.ToLower(m[2]), "to_date("); n != 11 {
		t.Errorf("korpus memuat %d To_date, mau 11", n)
	}
	if len(caraA) != 80 || caraB != 80 {
		t.Errorf("kolom %d (cara A), nilai %d (cara B); mau keduanya 80", len(caraA), caraB)
	}
	kami := []string{"ID"}
	for _, k := range kolomPesertaWarisan {
		kami = append(kami, k.Kolom)
	}
	if len(kami) != len(caraA) {
		t.Fatalf("pemetaan kami %d kolom, korpus %d", len(kami), len(caraA))
	}
	for i := range caraA {
		if !strings.EqualFold(caraA[i], kami[i]) {
			t.Errorf("posisi %d: korpus %s, kami %s", i, caraA[i], kami[i])
		}
	}
}

// TestSumberWarisanAdaDiMigrasi - `d.X` di 052, `p.X` di 051.
func TestSumberWarisanAdaDiMigrasi(t *testing.T) {
	detail := kolomTabelPeserta(t)
	header := map[string]bool{}
	for _, k := range kolomMigrasi(t, "051_t_premium_list.sql", "T_PREMIUM_LIST") {
		header[k] = true
	}
	for _, k := range kolomPesertaWarisan {
		switch {
		case strings.HasPrefix(k.Sumber, "d."):
			if !detail[strings.TrimPrefix(k.Sumber, "d.")] {
				t.Errorf("%s bersumber %s yang tidak ada di 052", k.Kolom, k.Sumber)
			}
		case strings.HasPrefix(k.Sumber, "p."):
			if !header[strings.TrimPrefix(k.Sumber, "p.")] {
				t.Errorf("%s bersumber %s yang tidak ada di 051", k.Kolom, k.Sumber)
			}
		case k.Sumber == "":
			switch k.Kolom {
			case "PL_NUMBER", "IDPEGA", "STATUSOLD", "STATUS":
			default:
				t.Errorf("%s tanpa sumber dan tanpa perakit", k.Kolom)
			}
		default:
			t.Errorf("%s bersumber %q yang tidak berawalan d./p.", k.Kolom, k.Sumber)
		}
	}
}

// TestNilaiSalinWarisanSejajarDanTertanam - STATUS/STATUSOLD VERBATIM.
func TestNilaiSalinWarisanSejajarDanTertanam(t *testing.T) {
	b := BarisWarisan{Tipe: "TP", Nilai: map[string]sql.NullString{
		"GROSS_PREMIUM": {String: "1.5", Valid: true},
		"SEX":           {String: "  ", Valid: true},
	}}
	arg := nilaiSalinWarisan("RNML-UJI", "UJI-WORK-1", b)
	if len(arg) != len(kolomPesertaWarisan) {
		t.Fatalf("%d argumen, %d kolom", len(arg), len(kolomPesertaWarisan))
	}
	nilai := map[string]any{}
	for i, k := range kolomPesertaWarisan {
		nilai[k.Kolom] = arg[i]
	}
	if nilai["PL_NUMBER"] != "RNML-UJI" || nilai["IDPEGA"] != "UJI-WORK-1" ||
		nilai["STATUSOLD"] != "0" ||
		nilai["STATUS"] != "1" || nilai["GROSS_PREMIUM"] != "1.5" {
		t.Errorf("nilai salin %v", nilai)
	}
	// Kosong menjadi NULL - bukan teks kosong.
	if nilai["SEX"] != nil || nilai["DOB"] != nil {
		t.Errorf("kosong dikirim %v/%v, mau nil", nilai["SEX"], nilai["DOB"])
	}
	for tipe, mau := range map[string]string{"QR": "0", "QP": "0", "TP": "1", "TR": "1", "": "1"} {
		if got := StatusSalinWarisan(tipe); got != mau {
			t.Errorf("STATUS tipe %q = %s, mau %s", tipe, got, mau)
		}
	}
	// Tanggal dan teks-tanggal dibungkus TO_DATE berpola, sisanya tidak.
	q := sqlSisipPesertaWarisan("S.M", "S.Q")
	// 9 kolom DATE + STNC + WPC = 11, sama dengan cacah `To_date` korpus
	// (ditagih juga di `TestKolomWarisanVERBATIMDariSaveMasterLPDet`).
	if n := strings.Count(q, "TO_DATE("); n != 11 {
		t.Errorf("%d TO_DATE, mau 11 (9 tanggal + STNC + WPC)", n)
	}
	if !strings.Contains(q, "TO_CHAR(S.Q.NEXTVAL)") {
		t.Errorf("ID tidak dari sequence:\n%s", q)
	}
}
