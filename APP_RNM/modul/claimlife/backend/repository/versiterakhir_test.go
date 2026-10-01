package repository

// Peserta VERSI TERAKHIR - keputusan work owner 01-10-2026 (OQ-N14, brief
// `PROMPT-PERBAIKAN-CLAIM-LIFE-PESERTA-VERSI-TERAKHIR.md`): sertifikat yang versi polis terakhirnya
// `Delete` atau `Batal` lewat endorsement TIDAK tampil dan TIDAK dapat diklaim. Penyimpangan sadar dari
// `Claim Life/RDBList/GetPesertaClaim_sql1.xml` (nol saringan versi maupun `EDMSTATUS`).
//
// Inventori pembaca `M_LIFE_PREMIUM_DETAIL` di modul ini (paket 1):
//
//	pesertapolis.go  Cari             WAJIB - daftar `Find Insured`: menentukan sertifikat yang dapat dipilih
//	pesertapolis.go  AmbilUntukKlaim  WAJIB - baris yang disalin ke klaim saat pendaftaran
//	gandawarisan.go  DOBSumberKosong, StatusWarisanTerakhir, AdaWarisanSamaDOL, IsiTertanggungCermin
//	                                  TIDAK - berkunci `ID` baris sumber (`SOURCE_ID`) yang SUDAH dipilih
//	                                  AmbilUntukKlaim saat pendaftaran; membaca riwayat baris itu, bukan memutuskan
//	                                  dapat-tidaknya diklaim
//	kolompeserta.go, models/klaimlife.go, services/simpanrnm.go
//	                                  TIDAK - hanya komentar; nol SQL ke tabel itu
//
// Tabel kasus di bawah dipakai DUA kali: oleh acuan aturan Go di berkas ini (selalu jalan) dan oleh uji
// `db` atas SQL sungguhan (`versiterakhir_db_test.go`, MELEWATI tanpa ORACLE_DSN).

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// barisUjiVersi - satu baris peserta tiruan (`UJI-`, nol data orang).
type barisUjiVersi struct {
	ID, PLNumberEDM, EDMStatus, TglInput, Nama string // TglInput `YYYY-MM-DD HH24:MI:SS`, kosong = NULL
}

// kasusVersi - satu sertifikat di satu PL_NUMBER dan jawaban yang dituntut.
type kasusVersi struct {
	Kode, Nama, PL, Sertifikat string
	Baris                      []barisUjiVersi
	Tampil                     bool   // tampil di Cari (dengan CariNama) - AmbilUntukKlaim mengikuti PilihAcuan
	IDTerpilih                 string // ID baris yang dipakai bila tampil
	CariNama                   string // kotak nama Find Insured (kosong = tanpa saringan nama)
}

func kasusVersiTerakhir() []kasusVersi {
	return []kasusVersi{
		{Kode: "a", Nama: "NB saja", PL: "UJI-PLV-A", Sertifikat: "UJI-C1",
			Baris:  []barisUjiVersi{{ID: "101", EDMStatus: ""}},
			Tampil: true, IDTerpilih: "101"},
		{Kode: "b", Nama: "NB + /01 Old: tampil sekali, nilai dari /01", PL: "UJI-PLV-B", Sertifikat: "UJI-C1",
			Baris:  []barisUjiVersi{{ID: "101"}, {ID: "102", PLNumberEDM: "UJI-PLV-B/01", EDMStatus: "Old"}},
			Tampil: true, IDTerpilih: "102"},
		{Kode: "c", Nama: "NB + /01 Delete: tidak tampil", PL: "UJI-PLV-C", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101"}, {ID: "102", PLNumberEDM: "UJI-PLV-C/01", EDMStatus: "Delete"}}},
		{Kode: "d", Nama: "NB + /01 Batal: tidak tampil", PL: "UJI-PLV-D", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101"}, {ID: "102", PLNumberEDM: "UJI-PLV-D/01", EDMStatus: "Batal"}}},
		{Kode: "e", Nama: "NB + /01 Delete + /02 New (masuk lagi): tampil dari /02", PL: "UJI-PLV-E", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101"}, {ID: "102", PLNumberEDM: "UJI-PLV-E/01", EDMStatus: "Delete"},
				{ID: "103", PLNumberEDM: "UJI-PLV-E/02", EDMStatus: "New"}},
			Tampil: true, IDTerpilih: "103"},
		{Kode: "f", Nama: `EDMSTATUS berspasi " Delete ": tidak tampil`, PL: "UJI-PLV-F", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101"}, {ID: "102", PLNumberEDM: "UJI-PLV-F/01", EDMStatus: " Delete "}}},
		{Kode: "g", Nama: "dua baris versi sama: TGL_INPUT terbaru menang (bukan ID terbesar)", PL: "UJI-PLV-G", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101"},
				{ID: "102", PLNumberEDM: "UJI-PLV-G/01", EDMStatus: "Old", TglInput: "2026-02-01 08:00:00"},
				{ID: "103", PLNumberEDM: "UJI-PLV-G/01", EDMStatus: "Delete", TglInput: "2026-01-01 08:00:00"}},
			Tampil: true, IDTerpilih: "102"},
		{Kode: "h", Nama: "PL_NUMBER berakhiran angka, baris NB ber-PL_NUMBER_EDM = PL_NUMBER: versi NB tetap 0", PL: "UJI-PLV-2024", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101", PLNumberEDM: "UJI-PLV-2024"},
				{ID: "102", PLNumberEDM: "UJI-PLV-2024/01", EDMStatus: "Delete"}}},
		{Kode: "i", Nama: "seri versi dan TGL_INPUT kosong: ID terbesar secara ANGKA (100 > 99)", PL: "UJI-PLV-I", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "99", PLNumberEDM: "UJI-PLV-I/01", EDMStatus: "Delete"},
				{ID: "100", PLNumberEDM: "UJI-PLV-I/01", EDMStatus: "Old"}},
			Tampil: true, IDTerpilih: "100"},
		{Kode: "j", Nama: "saringan nama Find Insured atas versi TERAKHIR, bukan versi lama", PL: "UJI-PLV-J", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101", Nama: "UJI-NAMA-LAMA"},
				{ID: "102", PLNumberEDM: "UJI-PLV-J/01", EDMStatus: "Old", Nama: "UJI-NAMA-BARU"}},
			CariNama: "lama"},
		{Kode: "k", Nama: "versi sama, TGL_INPUT terisi lawan NULL: yang terisi menang (NULLS LAST), bukan ID terbesar", PL: "UJI-PLV-K", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101"},
				{ID: "102", PLNumberEDM: "UJI-PLV-K/01", EDMStatus: "Old", TglInput: "2026-01-01 08:00:00"},
				{ID: "103", PLNumberEDM: "UJI-PLV-K/01", EDMStatus: "Delete"}},
			Tampil: true, IDTerpilih: "102"},
		{Kode: "l", Nama: "versi dibanding ANGKA: /100 menang atas /99 (teks: '99' > '100')", PL: "UJI-PLV-L", Sertifikat: "UJI-C1",
			Baris: []barisUjiVersi{{ID: "101", PLNumberEDM: "UJI-PLV-L/100", EDMStatus: "Old"},
				{ID: "102", PLNumberEDM: "UJI-PLV-L/99", EDMStatus: "Delete"}},
			Tampil: true, IDTerpilih: "101"},
	}
}

// versiAcuan - ACUAN Go aturan §1 butir 1 (bukan kode produksi): angka sesudah `<PL_NUMBER>/` di
// `PL_NUMBER_EDM`; selain itu 0 (new business). Diperketat dari `REGEXP_SUBSTR(PL_NUMBER_EDM, '[0-9]+$')`
// brief: angka di akhir PL_NUMBER sendiri tidak pernah terbaca sebagai versi (kasus h).
func versiAcuan(pl, edm string) int {
	e := strings.TrimSpace(edm)
	if !strings.HasPrefix(e, pl+"/") {
		return 0
	}
	n, err := strconv.Atoi(e[len(pl)+1:])
	if err != nil || strings.Trim(e[len(pl)+1:], "0123456789") != "" {
		return 0
	}
	return n
}

// lebihBaruAcuan - urutan §1: versi DESC, TGL_INPUT DESC (kosong terakhir), ID DESC secara angka.
func lebihBaruAcuan(pl string, a, b barisUjiVersi) bool {
	if va, vb := versiAcuan(pl, a.PLNumberEDM), versiAcuan(pl, b.PLNumberEDM); va != vb {
		return va > vb
	}
	if a.TglInput != b.TglInput {
		if a.TglInput == "" || b.TglInput == "" {
			return b.TglInput == ""
		}
		return a.TglInput > b.TglInput
	}
	if len(a.ID) != len(b.ID) {
		return len(a.ID) > len(b.ID)
	}
	return a.ID > b.ID
}

// pilihAcuan - baris versi terakhir satu kasus menurut acuan, dan apakah ia hidup (PesertaHidup PRODUKSI).
// Tanpa saringan Nama: itulah jawaban AmbilUntukKlaim.
func pilihAcuan(k kasusVersi) (barisUjiVersi, bool) {
	terpilih := k.Baris[0]
	for _, b := range k.Baris[1:] {
		if lebihBaruAcuan(k.PL, b, terpilih) {
			terpilih = b
		}
	}
	return terpilih, PesertaHidup(terpilih.EDMStatus)
}

// TestAturanVersiTerakhirAcuan - jawaban kasus (a)-(l) menurut acuan aturan; PesertaHidup PRODUKSI menilai
// baris terpilih. Tabel kasusnya juga yang dijalankan uji `db` atas SQL sungguhan.
func TestAturanVersiTerakhirAcuan(t *testing.T) {
	for _, k := range kasusVersiTerakhir() {
		t.Run(k.Kode, func(t *testing.T) {
			terpilih, tampil := pilihAcuan(k)
			if k.CariNama != "" && !strings.Contains(strings.ToUpper(terpilih.Nama), strings.ToUpper(k.CariNama)) {
				tampil = false
			}
			if tampil != k.Tampil || (tampil && terpilih.ID != k.IDTerpilih) {
				t.Errorf("%s: tampil %v ID %s, mau %v ID %s", k.Nama, tampil, terpilih.ID, k.Tampil, k.IDTerpilih)
			}
		})
	}
}

// rataSQL - SQL dirapikan satu spasi supaya pencocokan tidak bergantung indentasi.
func rataSQL(q string) string { return strings.Join(strings.Fields(q), " ") }

// Potongan SQL yang DITUNTUT aturan - teks harfiah, bukan konstanta produksi: mengubah konstanta produksi
// tidak boleh ikut mengubah tuntutannya.
const (
	tuntutPeringkat = "ROW_NUMBER() OVER (PARTITION BY CERTIFICATE_NO ORDER BY"
	tuntutVersi     = "NVL(TO_NUMBER(CASE WHEN SUBSTR(TRIM(PL_NUMBER_EDM), 1, LENGTH(PL_NUMBER) + 1) = PL_NUMBER || '/' " +
		"THEN REGEXP_SUBSTR(SUBSTR(TRIM(PL_NUMBER_EDM), LENGTH(PL_NUMBER) + 2), '^[0-9]+$') END), 0) DESC"
	tuntutSeri  = "TGL_INPUT DESC NULLS LAST, LENGTH(ID) DESC NULLS LAST, ID DESC NULLS LAST)"
	tuntutPilih = "RN_VERSI = 1 AND (EDMSTATUS IS NULL OR TRIM(EDMSTATUS) NOT IN ('Batal','Delete'))"
)

// TestSQLCariPesertaVersiTerakhir - bentuk SQL `Find Insured` per kasus: setiap kasus menuntut potongan yang
// memutuskannya. Menghapus ROW_NUMBER membuat (b), (c), (d), (e), (g), (i) merah.
func TestSQLCariPesertaVersiTerakhir(t *testing.T) {
	q, _ := sqlCariPeserta("S.M", "UJI-PLV", "", "", 10)
	qn, _ := sqlCariPeserta("S.M", "UJI-PLV", "UJI-C", "lama", 10)
	for _, g := range galatCariVersi(q, qn) {
		t.Error(g)
	}
}

// galatCariVersi - pelanggaran aturan versi terakhir pada SQL Find Insured (q tanpa saringan, qn dengan
// saringan sertifikat + nama), berawalan kode kasus. Dipisah supaya uji gigit menjalankan PEMERIKSA YANG SAMA
// atas SQL termutasi.
func galatCariVersi(q, qn string) []string {
	q, qn = rataSQL(q), rataSQL(qn)
	var galat []string
	tuntut := map[string][]string{
		"a": {tuntutPilih},
		"b": {tuntutPeringkat, tuntutPilih},
		"c": {tuntutPeringkat, tuntutVersi, tuntutPilih},
		"d": {tuntutPeringkat, tuntutVersi, tuntutPilih},
		"e": {tuntutPeringkat, tuntutVersi, tuntutPilih},
		"f": {"TRIM(EDMSTATUS) NOT IN ('Batal','Delete')"},
		"g": {tuntutPeringkat, tuntutSeri},
		"h": {tuntutVersi},
		"i": {tuntutSeri},
		"k": {tuntutSeri},
		"l": {tuntutVersi},
	}
	for kode, potongan := range tuntut {
		for _, p := range potongan {
			if !strings.Contains(q, p) {
				galat = append(galat, "("+kode+") SQL tanpa "+strconv.Quote(p)+":\n"+q)
			}
		}
	}
	// (j) saringan nama di LUAR jendela (atas versi terakhir); saringan sertifikat boleh di dalam (seluruh
	// kelompok satu sertifikat ikut atau tidak, peringkat di dalamnya tidak berubah).
	if i := strings.Index(qn, "RN_VERSI = 1"); i < 0 || !strings.Contains(qn[i:], "UPPER(NAME_OF_INSURED) LIKE") {
		galat = append(galat, "(j) saringan nama tidak berada sesudah peringkat versi:\n"+qn)
	}
	for _, x := range []string{q, qn} {
		galat = append(galat, galatJendela(x)...)
	}
	// §1 butir 4: jendela hanya menyentuh SATU polis - PL_NUMBER = :1 di DALAM subkueri peringkat.
	if !strings.Contains(segmenJendela(q), "WHERE PL_NUMBER = :1") {
		galat = append(galat, "(pemindaian penuh) jendela peringkat tidak dikurung PL_NUMBER:\n"+q)
	}
	galat = append(galat, galatPenampung(q, 1)...)
	galat = append(galat, galatPenampung(qn, 3)...)
	return galat
}

// segmenJendela - subkueri peringkat: dari `FROM (` sampai `WHERE` luar yang memilih `RN_VERSI = 1`.
func segmenJendela(q string) string {
	a, b := strings.Index(q, "FROM ("), strings.Index(q, "RN_VERSI = 1")
	if a < 0 || b <= a {
		return ""
	}
	return q[a:b]
}

// galatJendela - penyaring status atau nama di DALAM jendela: baris yang tersaring sebelum diperingkat membuat
// versi lamanya naik menjadi "terakhir" - `/01 Delete` disaring, NB-nya tampil lagi (c); nama lama tetap
// menemukan pesertanya (j). Kehadiran potongan yang benar di luar tidak menebus kehadirannya di dalam.
func galatJendela(q string) []string {
	var galat []string
	j := segmenJendela(q)
	for kode, terlarang := range map[string][]string{
		"c": {"EDMSTATUS IS NULL", "NOT IN ('Batal','Delete')"},
		"j": {"NAME_OF_INSURED) LIKE"},
	} {
		for _, x := range terlarang {
			if strings.Contains(j, x) {
				galat = append(galat, "("+kode+") "+strconv.Quote(x)+" di DALAM jendela peringkat:\n"+q)
			}
		}
	}
	return galat
}

// polaPenampung - penampung bind `:n`. `:MI`/`:SS` format TO_CHAR bukan angka, tidak tercocok.
var polaPenampung = regexp.MustCompile(`:([0-9]+)`)

// galatPenampung - godror mengikat menurut URUTAN KEMUNCULAN, bukan nomor: penampung harus muncul :1, :2, …, :n
// berurutan di teks, tepat n buah - sama dengan urutan nilai bind yang dirakit.
func galatPenampung(q string, n int) []string {
	var dapat []string
	for _, m := range polaPenampung.FindAllStringSubmatch(q, -1) {
		dapat = append(dapat, m[1])
	}
	var mau []string
	for i := 1; i <= n; i++ {
		mau = append(mau, strconv.Itoa(i))
	}
	if strings.Join(dapat, ",") != strings.Join(mau, ",") {
		return []string{"(bind) penampung muncul " + strings.Join(dapat, ",") + ", mau " + strings.Join(mau, ",") + ":\n" + q}
	}
	return nil
}

// TestSQLAmbilPesertaKlaimVersiTerakhir - SQL penyalin ke klaim: baris VERSI TERAKHIR yang hidup, per kasus.
func TestSQLAmbilPesertaKlaimVersiTerakhir(t *testing.T) {
	for _, g := range galatAmbilVersi(sqlAmbilPesertaKlaim("S.M")) {
		t.Error(g)
	}
}

// galatAmbilVersi - pelanggaran aturan versi terakhir pada SQL AmbilUntukKlaim (tabel uji `S.M`).
func galatAmbilVersi(q string) []string {
	q = rataSQL(q)
	var galat []string
	for kode, potongan := range map[string][]string{
		"b-e": {tuntutPeringkat, tuntutVersi, "WHERE " + tuntutPilih + " FETCH FIRST 1 ROWS ONLY"},
		"f":   {"TRIM(EDMSTATUS) NOT IN ('Batal','Delete')"},
		"g,i": {tuntutSeri},
		"h,l": {tuntutVersi},
	} {
		for _, p := range potongan {
			if !strings.Contains(q, p) {
				galat = append(galat, "("+kode+") SQL tanpa "+strconv.Quote(p)+":\n"+q)
			}
		}
	}
	// Subkueri peringkat dikurung PL_NUMBER + CERTIFICATE_NO (§1 butir 4).
	if !strings.Contains(segmenJendela(q), "FROM S.M WHERE PL_NUMBER = :1 AND CERTIFICATE_NO = :2)") {
		galat = append(galat, "(pemindaian penuh) subkueri peringkat tidak dikurung PL_NUMBER + CERTIFICATE_NO:\n"+q)
	}
	galat = append(galat, galatJendela(q)...)
	galat = append(galat, galatPenampung(q, 2)...)
	// Daftar pilih luar = kolomSalin (posisi dikunci salinKePeserta); subkueri membawa kolom mentahnya, bukan `*`.
	if !strings.HasPrefix(q, "SELECT ID, PL_NUMBER, POLICY_NO, CERTIFICATE_NO, CURRENCY, TO_CHAR(STNC,") {
		galat = append(galat, "daftar pilih kolomSalin berubah:\n"+q)
	}
	if !strings.Contains(q, "FROM (SELECT ID, PL_NUMBER, POLICY_NO, CERTIFICATE_NO, CURRENCY, STNC, ") ||
		!strings.Contains(q, ", CLAIM_AMOUNT, EDMSTATUS, ROW_NUMBER()") {
		galat = append(galat, "(kolom) subkueri tidak membawa kolom mentah kolomSalin + EDMSTATUS:\n"+q)
	}
	if strings.Contains(q, "*") {
		galat = append(galat, "(KTP) `*` membaca seluruh kolom, termasuk KTP:\n"+q)
	}
	return galat
}

// TestAturanVersiTerakhirMenggigit - uji gigit §2 butir 4, PEMERIKSA YANG SAMA atas SQL yang sungguh dirakit
// lalu dimutasi: tanpa ROW_NUMBER kasus (c) merah; tanpa `PL_NUMBER = :1` di jendela, penjaga "nol pemindaian
// penuh" merah; status di dalam jendela (c) merah; penampung tak urut (bind) merah. SQL asli harus nol galat
// (pemeriksanya tidak asal menuduh).
func TestAturanVersiTerakhirMenggigit(t *testing.T) {
	q, _ := sqlCariPeserta("S.M", "UJI-PLV", "", "", 10)
	qn, _ := sqlCariPeserta("S.M", "UJI-PLV", "UJI-C", "lama", 10)
	if g := galatCariVersi(q, qn); len(g) != 0 {
		t.Fatalf("SQL asli sudah bergalat: %v", g)
	}
	ada := func(galat []string, awalan string) bool {
		for _, g := range galat {
			if strings.HasPrefix(g, awalan) {
				return true
			}
		}
		return false
	}
	ganti := func(lama, baru string) func(string) string {
		return func(x string) string { return strings.Replace(x, lama, baru, 1) }
	}
	for _, m := range []struct {
		nama     string
		ubah     func(string) string
		mauMerah string
	}{
		{"tanpa ROW_NUMBER", ganti(peringkatVersi, "1"), "(c)"},
		{"tanpa PL_NUMBER di jendela", ganti("WHERE PL_NUMBER = :1", "WHERE 1 = 1"), "(pemindaian penuh)"},
		{"status JUGA di dalam jendela", ganti("WHERE PL_NUMBER = :1", "WHERE PL_NUMBER = :1 AND "+penyaringHidup), "(c)"},
		{"penampung tak urut", strings.NewReplacer(":2", ":3", ":3", ":2").Replace, "(bind)"},
	} {
		mq, mqn := m.ubah(q), m.ubah(qn)
		if mq == q && mqn == qn {
			t.Fatalf("%s: mutasi tidak mengubah SQL", m.nama)
		}
		if g := galatCariVersi(mq, mqn); !ada(g, m.mauMerah) {
			t.Errorf("mutasi %s: %s tetap hijau - pemeriksanya tidak menggigit (%v)", m.nama, m.mauMerah, g)
		}
	}
	ambil := sqlAmbilPesertaKlaim("S.M")
	if g := galatAmbilVersi(ambil); len(g) != 0 {
		t.Fatalf("SQL AmbilUntukKlaim asli sudah bergalat: %v", g)
	}
	for _, m := range []struct {
		nama     string
		ubah     func(string) string
		mauMerah string
	}{
		{"tanpa ROW_NUMBER", ganti(peringkatVersi, "1"), "(b-e)"},
		{"tanpa PL_NUMBER di subkueri", ganti("WHERE PL_NUMBER = :1 AND ", "WHERE "), "(pemindaian penuh)"},
		{"status JUGA di subkueri", ganti("CERTIFICATE_NO = :2)", "CERTIFICATE_NO = :2 AND "+penyaringHidup+")"), "(c)"},
		{"subkueri SELECT *", ganti(kolomMentahSalin+", EDMSTATUS,", "T.*,"), "(KTP)"},
	} {
		mut := m.ubah(ambil)
		if mut == ambil {
			t.Fatalf("%s: mutasi tidak mengubah SQL AmbilUntukKlaim", m.nama)
		}
		if g := galatAmbilVersi(mut); !ada(g, m.mauMerah) {
			t.Errorf("mutasi AmbilUntukKlaim %s: %s tetap hijau (%v)", m.nama, m.mauMerah, g)
		}
	}
}
