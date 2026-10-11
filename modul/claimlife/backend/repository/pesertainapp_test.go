package repository

// Aturan in-app Register (10-10-2026) - bentuk SQL dan aturan status tanpa Oracle; perilaku sungguhan di uji -tags db
// services (TestDaftarPesertaAktifDiKlaimLain).

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
)

func TestSqlPesertaAktifDiKlaimLain(t *testing.T) {
	q := sqlPesertaAktifDiKlaimLain("S.T_CLAIMLF_PREMIUMLIST_DETAIL", "S.T_GENERAL_CLAIM", "S.T_CLAIMLF_ADJUSTMENT", 3)
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	s := strings.Join(strings.Fields(q), " ")
	for _, wajib := range []string{
		"FROM S.T_CLAIMLF_PREMIUMLIST_DETAIL p JOIN S.T_GENERAL_CLAIM h ON h.ID = p.CLAIM_ID",
		// baris adjustment putaran: yang bukan ditolak (:1) dan yang aksep (:2)
		"(SELECT COUNT(*) FROM S.T_CLAIMLF_ADJUSTMENT a WHERE a.PREMIUM_LIST_DETAIL_ID = p.ID AND (a.STS_REJECT IS NULL OR TRIM(a.STS_REJECT) <> :1))",
		"(SELECT COUNT(*) FROM S.T_CLAIMLF_ADJUSTMENT a WHERE a.PREMIUM_LIST_DETAIL_ID = p.ID AND TRIM(a.STS_REJECT) = :2)",
		// kunci = PL + sertifikat dipangkas, satu kueri berdaftar IN
		"WHERE TRIM(p.PL_NUMBER) = :3 AND TRIM(p.CERTIFICATE_NO) IN (:4, :5, :6)",
		// baris tercabut tidak dihitung
		"AND p.STS_HAPUS IS NULL",
	} {
		if !strings.Contains(s, wajib) {
			t.Errorf("SQL tanpa %q:\n%s", wajib, s)
		}
	}
	// header TIDAK menentukan; tidak ada literal kode; bukan tabel warisan.
	for _, larang := range []string{"h.STS_REJECT", "STATUS_WORK", "'0'", "'1'", "'2'", "SOURCE_ID", "OS_AKSEPTASI"} {
		if strings.Contains(s, larang) {
			t.Errorf("SQL memuat %q:\n%s", larang, s)
		}
	}
}

// Aturan status tingkat peserta (keputusan work owner 10-10-2026 + model putaran Add).
func TestAturanAktifInApp(t *testing.T) {
	for _, k := range []struct {
		nama                   string
		kode                   string
		barisAktif, barisAksep int
		aktif                  bool
		kuat                   string
	}{
		{"terdaftar belum diputuskan (NULL)", "", 1, 0, true, ""},
		{"Outstanding", kontrak.KodeOutstanding, 1, 0, true, kontrak.KodeOutstanding},
		{"Accept", kontrak.KodeAksep, 1, 1, true, kontrak.KodeAksep},
		{"Rejected di semua tingkat", kontrak.KodeDitolak, 0, 0, false, ""},
		{"Rejected lalu putaran Add baru (baris kosong)", kontrak.KodeDitolak, 1, 0, true, ""},
		{"Rejected lalu putaran baru diaksep", kontrak.KodeDitolak, 1, 1, true, kontrak.KodeAksep},
		{"klaim lama tanpa baris adjustment", "", 0, 0, true, ""},
	} {
		aktif, kuat := aktifInApp(k.kode, k.barisAktif, k.barisAksep)
		if aktif != k.aktif || kuat != k.kuat {
			t.Errorf("%s: aktif %v kode %q, mau %v %q", k.nama, aktif, kuat, k.aktif, k.kuat)
		}
	}
}

// TestIndeksPesertaAktifCocokDenganKueri - migrasi 024 (keputusan work owner 10-10-2026): ekspresi indeks fungsi
// IX_PLD_TRIM_PL_SERT harus PERSIS ekspresi saringan kueri, dalam urutan yang sama - indeks fungsi hanya terpakai
// bila ekspresinya sama. Mengubah salah satunya tanpa yang lain membuat penjaga ini merah.
func TestIndeksPesertaAktifCocokDenganKueri(t *testing.T) {
	baca := func(nama string) string {
		b, err := os.ReadFile(filepath.Join("..", "migrations", nama))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	const nama = "IX_PLD_TRIM_PL_SERT"
	if len(nama) > 30 {
		t.Fatalf("nama indeks %d byte", len(nama))
	}
	maju := baca("024_ix_peserta_aktif_klaim.sql")
	awal := "CREATE INDEX {skema}." + nama + " ON {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL ("
	i := strings.Index(maju, awal)
	if i < 0 {
		t.Fatal("CREATE INDEX 024 tidak ditemukan")
	}
	sisa := maju[i+len(awal):]
	// sisa = "TRIM(PL_NUMBER), TRIM(CERTIFICATE_NO))';" - kurung luar ditutup tepat sebelum kutip penutup.
	ekspresi := strings.Split(sisa[:strings.Index(sisa, ")'")], ", ")
	q := strings.Join(strings.Fields(sqlPesertaAktifDiKlaimLain("S.P", "S.H", "S.A", 2)), " ")
	if len(ekspresi) != 2 || ekspresi[0] != "TRIM(PL_NUMBER)" || ekspresi[1] != "TRIM(CERTIFICATE_NO)" {
		t.Fatalf("ekspresi indeks %q", ekspresi)
	}
	letakPL := strings.Index(q, "TRIM(p.PL_NUMBER) = :3")
	letakSert := strings.Index(q, "TRIM(p.CERTIFICATE_NO) IN (")
	if letakPL < 0 || letakSert < 0 || letakSert < letakPL {
		t.Errorf("kueri tidak lagi menyaring dengan ekspresi indeks yang sama:\n%s", q)
	}
	if !strings.Contains(maju, "AND INDEX_NAME = '"+nama+"'") || strings.Contains(strings.ToUpper(buangKomentarSumber("024.sql", maju)), "COMMIT") {
		t.Error("024 maju: pelindung ALL_INDEXES hilang atau memuat COMMIT")
	}
	mundur := baca("024_ix_peserta_aktif_klaim_down.sql")
	if !strings.Contains(mundur, "DROP INDEX {skema}."+nama) || !strings.Contains(mundur, "AND INDEX_NAME = '"+nama+"'") {
		t.Error("024 mundur tidak membuang indeks yang sama secara berpelindung")
	}
}
