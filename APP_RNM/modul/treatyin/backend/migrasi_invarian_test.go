package backend

// Uji invarian atas TEKS DDL - bukan atas Oracle.
//
// ⛔ `L-3`: tidak ada instans Oracle yang terjangkau, sehingga tidak satu baris
// pun DDL modul ini pernah dijalankan di mana pun. Yang dapat diperiksa tanpa
// Oracle adalah BENTUKNYA, dan berkas ini memeriksanya - supaya invarian yang
// hilang dari DDL tertangkap di sini, bukan saat data pertama masuk.
//
// Yang berkas ini TIDAK dapat buktikan: bahwa Oracle menerima teksnya, dan
// bahwa constraint-nya sungguh menolak. Itu menunggu instans.

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
)

// tanpaKomentar membuang baris komentar `--`, menyisakan PERNYATAAN SQL saja.
//
// ⛔ Perlu, dan sebabnya layak dibaca: berkas migrasi modul ini MENJELASKAN
// keputusannya di dalam komentar - "TANPA \"ON DELETE\"", "nol INSERT". Sapuan
// yang memindai teks mentah menemukan kata-kata itu di prosanya sendiri dan
// menuduh DDL melakukan hal yang justru dinyatakan TIDAK dilakukan.
func tanpaKomentar(sql string) string {
	var b strings.Builder
	for _, baris := range strings.Split(sql, "\n") {
		if strings.HasPrefix(strings.TrimSpace(baris), "--") {
			continue
		}
		b.WriteString(baris)
		b.WriteByte('\n')
	}
	return b.String()
}

func seluruhMigrasi(t *testing.T) map[string]string {
	t.Helper()
	entri, err := fs.ReadDir(berkasMigrasi, "migrations")
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, e := range entri {
		isi, err := fs.ReadFile(berkasMigrasi, "migrations/"+e.Name())
		if err != nil {
			t.Fatal(err)
		}
		out[e.Name()] = string(isi)
	}
	if len(out) == 0 {
		t.Fatal("nol berkas migrasi terbaca; pembacanya yang rusak")
	}
	return out
}

func majuSaja(t *testing.T) map[string]string {
	t.Helper()
	out := map[string]string{}
	for nama, isi := range seluruhMigrasi(t) {
		if !strings.HasSuffix(nama, "_down.sql") {
			out[nama] = isi
		}
	}
	return out
}

func gabungan(t *testing.T) string {
	t.Helper()
	var b strings.Builder
	for _, isi := range majuSaja(t) {
		b.WriteString(isi)
	}
	return b.String()
}

// INV-01: setiap tabel berkunci utama.
func TestSetiapTabelBerkunciUtama(t *testing.T) {
	sql := gabungan(t)
	tabel := []string{
		// tiket 15 dan 14
		"MATA_UANG", "JENIS_POTONGAN", "KELAS_BISNIS", "KELOMPOK_TREATY", "BAHAYA",
		"JENIS_REASURANSI", "KONTRAK", "VERSI_KONTRAK",
		// tiket 20, 22-31, 39
		"MATA_UANG_KONTRAK", "RETENSI_CEDANT", "EGNPI", "PORTOFOLIO", "PERIODE_PELAPORAN",
		"PERIODE_AKUMULASI", "TERMIN", "SKALA_KOASURANSI", "BATAS_PER_BAHAYA",
		"DOKUMEN_KONTRAK", "LAYER", "NILAI_MDP", "NILAI_MDP_MINIMUM", "PEMULIHAN_LIMIT",
		"JEJAK_PERUBAHAN",
		// tiket 33, 34, 37
		"BAGIAN", "NILAI_PREMI_BRUTO", "NILAI_PREMI_BRUTO_MINIMUM",
		"DETAIL_PROPORSIONAL", "NILAI_CADANGAN_PREMI", "POTONGAN",
	}
	for _, n := range tabel {
		if !strings.Contains(sql, "CREATE TABLE {skema}."+n+" (") {
			t.Errorf("tabel %s tidak dibuat migrasi mana pun", n)
			continue
		}
		if !strings.Contains(sql, "CONSTRAINT PK_"+n+" PRIMARY KEY") {
			t.Errorf("INV-01: tabel %s tanpa kunci utama bernama PK_%s", n, n)
		}
	}
}

// INV-02 dan INV-03: pengenal dari SEQUENCE, dan sequence tanpa CYCLE.
//
// Uji POSITIF ikut di sini: kedelapan sequence HARUS ada. Uji yang hanya
// memeriksa "NOCYCLE pada yang ada" lulus walau tujuh di antaranya hilang.
func TestSequenceAdaDanTanpaCycle(t *testing.T) {
	sql := gabungan(t)
	seq := []string{
		"SEQ_TRIN_KONTRAK", "SEQ_TRIN_VERSI_KONTRAK", "SEQ_TRIN_MATA_UANG",
		"SEQ_TRIN_JENIS_POTONGAN", "SEQ_TRIN_JENIS_REASURANSI", "SEQ_TRIN_BAHAYA",
		"SEQ_TRIN_KELOMPOK_TREATY", "SEQ_TRIN_KELAS_BISNIS",
		"SEQ_TRIN_MATA_UANG_KONTRAK", "SEQ_TRIN_RETENSI_CEDANT", "SEQ_TRIN_EGNPI",
		"SEQ_TRIN_PORTOFOLIO", "SEQ_TRIN_PERIODE_PELAPORAN", "SEQ_TRIN_PERIODE_AKUMULASI",
		"SEQ_TRIN_TERMIN", "SEQ_TRIN_SKALA_KOASURANSI", "SEQ_TRIN_BATAS_PER_BAHAYA",
		"SEQ_TRIN_DOKUMEN_KONTRAK", "SEQ_TRIN_LAYER", "SEQ_TRIN_NILAI_MDP",
		"SEQ_TRIN_NILAI_MDP_MINIMUM", "SEQ_TRIN_PEMULIHAN_LIMIT", "SEQ_TRIN_JEJAK_PERUBAHAN",
		"SEQ_TRIN_BAGIAN", "SEQ_TRIN_NILAI_PREMI_BRUTO", "SEQ_TRIN_NILAI_PB_MINIMUM",
		"SEQ_TRIN_DETAIL_PROPORSIONAL", "SEQ_TRIN_NILAI_CADANGAN_PREMI", "SEQ_TRIN_POTONGAN",
	}
	for _, n := range seq {
		awal := "CREATE SEQUENCE {skema}." + n + " "
		i := strings.Index(sql, awal)
		if i < 0 {
			t.Errorf("INV-02: sequence %s tidak ada", n)
			continue
		}
		baris := sql[i:]
		if j := strings.IndexByte(baris, '\n'); j >= 0 {
			baris = baris[:j]
		}
		if !strings.Contains(baris, "NOCYCLE") {
			t.Errorf("INV-03: sequence %s tanpa NOCYCLE: %q", n, baris)
		}
	}
	// Panjang nama diperiksa TestNamaObjekDiBawahTigaPuluhBita, yang
	// membacanya dari teks DDL - bukan dari daftar literal di atas.
}

// INV-04: NOMOR_URUT_VERSI unik DI DALAM satu KONTRAK - bukan unik sendirian.
//
// ⛔ Uji ini menyebut KEDUA kolomnya dengan sengaja. `UNIQUE (NOMOR_URUT_VERSI)`
// saja menolak hal yang salah: ia melarang dua kontrak berbeda sama-sama punya
// versi nomor 1, dan itu MENOLAK DATA YANG SAH.
func TestNomorUrutVersiUnikDiDalamKontrak(t *testing.T) {
	sql := gabungan(t)
	if !strings.Contains(sql, "CONSTRAINT UQ_VERSI_KONTRAK UNIQUE (ID_KONTRAK, NOMOR_URUT_VERSI)") {
		t.Error("INV-04: UQ_VERSI_KONTRAK bukan (ID_KONTRAK, NOMOR_URUT_VERSI)")
	}
}

// INV-68: KODE unik di KEENAM tabel acuan - diperiksa satu per satu, bukan
// disimpulkan dari satu.
func TestKodeUnikDiKeenamTabelAcuan(t *testing.T) {
	sql := gabungan(t)
	for _, n := range []string{"MATA_UANG", "JENIS_POTONGAN", "KELAS_BISNIS", "KELOMPOK_TREATY", "BAHAYA", "JENIS_REASURANSI"} {
		if !strings.Contains(sql, "CONSTRAINT UQ_"+n+" UNIQUE (KODE)") {
			t.Errorf("INV-68: %s tanpa UNIQUE (KODE)", n)
		}
	}
}

// INV-62 dan ADR-0038: nol CHECK berisi DAFTAR NILAI - bukan nol CHECK.
//
// ⛔ Uji ini pernah melarang SETIAP "CHECK (", dan itu terlalu luas.
// `ddl-usulan/` sendiri memuat dua CHECK - `CK_POTONGAN_INDUK` dan
// `CK_PENYEBARAN_INDUK` (tiket 38) - dan keduanya SAH: keduanya menyatakan
// BENTUK baris (induk polimorfik, tepat satu kolom induk terisi), bukan nilai
// apa yang boleh masuk sebuah kolom.
//
// Yang dilarang adalah CHECK yang MENGENUMERASI nilai, sebab itulah bentuk
// yang ADR-0038 buang: himpunan yang dapat bertambah disimpan sebagai DATA,
// bukan sebagai CHECK. Bahaya kesembilan tidak boleh menuntut perubahan skema.
//
// Larangan atas TRIGGER dan PROCEDURE berdiri terpisah, di
// TestNolTriggerProcedureDanCommit - ADR-0056 melarang keduanya, bukan CHECK.
func TestNolCheckDaftarNilai(t *testing.T) {
	// CHECK yang SAH disebut satu per satu beserta sebabnya. Yang tidak
	// terdaftar ditolak: daftar putih yang tumbuh diam-diam bukan daftar putih.
	sah := map[string]string{
		"CK_POTONGAN_INDUK": "KTV-B - induk polimorfik, tepat satu dari ID_BAGIAN dan " +
			"ID_DETAIL_PROPORSIONAL terisi; menyatakan BENTUK baris, bukan daftar nilai",
	}
	bernama := regexp.MustCompile(`(?is)CONSTRAINT\s+(\w+)\s+CHECK\s*\(`)
	semua := regexp.MustCompile(`(?is)CHECK\s*\(`)
	for nama, isi := range majuSaja(t) {
		bersih := tanpaKomentar(isi)
		cocok := bernama.FindAllStringSubmatch(bersih, -1)
		for _, m := range cocok {
			if sebab, ok := sah[m[1]]; ok {
				t.Logf("%s: CHECK %s sah - %s", nama, m[1], sebab)
				continue
			}
			t.Errorf("%s: CHECK %s tidak terdaftar sebagai CHECK yang sah. "+
				"CHECK berisi DAFTAR NILAI dilarang ADR-0038; CHECK yang menyatakan "+
				"BENTUK baris boleh, tetapi didaftar di uji ini beserta sebabnya.",
				nama, m[1])
		}
		// CHECK tanpa CONSTRAINT bernama tidak dapat diadili di sini, dan
		// tidak dapat dicabut dengan tepat kelak.
		if n := len(semua.FindAllString(bersih, -1)); n != len(cocok) {
			t.Errorf("%s memuat %d CHECK tetapi hanya %d yang bernama; "+
				"CHECK tanpa nama tidak dapat diadili", nama, n, len(cocok))
		}
	}
}

// INV-18: perilaku hapus setiap kunci asing ditetapkan SADAR.
//
// Tidak ada ON DELETE di modul ini: bawaan Oracle MENOLAK, dan menolak memang
// yang dikehendaki. Uji ini menjaga keputusan itu tetap terbaca - kaskade yang
// disisipkan diam-diam akan menghilangkan riwayat persetujuan tiket 54.
func TestNolKaskadeHapus(t *testing.T) {
	for nama, isi := range majuSaja(t) {
		if strings.Contains(strings.ToUpper(tanpaKomentar(isi)), "ON DELETE") {
			t.Errorf("INV-18: %s memuat ON DELETE; kebijakan modul ini adalah MENOLAK, tanpa klausa", nama)
		}
	}
	// Dan kunci asingnya memang ada - tanpa ini uji di atas lulus pada modul
	// yang tidak punya kunci asing sama sekali.
	sql := gabungan(t)
	for _, fk := range []string{
		"FK_VERSI_KONTRAK_1", "FK_VERSI_KONTRAK_MATA_UANG", "FK_JENIS_REASURANSI_INDUK",
		"FK_MATA_UANG_KONTRAK_1", "FK_MATA_UANG_KONTRAK_2", "FK_RETENSI_CEDANT_1",
		"FK_RETENSI_CEDANT_2", "FK_EGNPI_1", "FK_EGNPI_2", "FK_EGNPI_3", "FK_PORTOFOLIO_1",
		"FK_PERIODE_PELAPORAN_1", "FK_PERIODE_AKUMULASI_1", "FK_TERMIN_1",
		"FK_SKALA_KOASURANSI_1", "FK_BATAS_PER_BAHAYA_1", "FK_BATAS_PER_BAHAYA_2",
		"FK_DOKUMEN_KONTRAK_1", "FK_LAYER_1", "FK_NILAI_MDP_1", "FK_NILAI_MDP_MINIMUM_1",
		"FK_PEMULIHAN_LIMIT_1", "FK_JEJAK_PERUBAHAN_1",
		"FK_BAGIAN_1", "FK_NILAI_PREMI_BRUTO_1", "FK_NILAI_PREMI_BRUTO_MIN_1",
		"FK_DETAIL_PROPORSIONAL_1", "FK_DETAIL_PROPORSIONAL_2", "FK_NILAI_CADANGAN_PREMI_1",
		"FK_POTONGAN_1", "FK_POTONGAN_2", "FK_POTONGAN_3",
	} {
		if !strings.Contains(sql, "CONSTRAINT "+fk+" FOREIGN KEY") {
			t.Errorf("kunci asing %s tidak ada", fk)
		}
	}
}

// ADR-0056 (K-4) dan ADR-U-0029: nol trigger, nol procedure, nol COMMIT.
func TestNolTriggerProcedureDanCommit(t *testing.T) {
	for nama, isi := range seluruhMigrasi(t) {
		atas := strings.ToUpper(tanpaKomentar(isi))
		for _, terlarang := range []string{"CREATE TRIGGER", "CREATE OR REPLACE TRIGGER", "CREATE PROCEDURE", "CREATE FUNCTION", "\nCOMMIT"} {
			if strings.Contains(atas, terlarang) {
				t.Errorf("%s memuat %q; aturan bisnis tidak turun ke basis data (ADR-0056)", nama, terlarang)
			}
		}
	}
}

// Kunci alami tabel anak - INV-05, INV-07 s.d. INV-14, INV-66, INV-67.
//
// ⛔ Tiap baris di bawah menyebut KOLOMNYA, bukan sekadar "ada UNIQUE". Itu
// pokoknya: `Z00_KUNCI_ALAMI.sql` memperingatkan bahwa menulis UNIQUE TANPA
// mata uang pada INV-08, INV-09, dan INV-12 bukan jalan tengah - ia mengubah
// artinya menjadi "dilarang dua baris bermata uang berbeda", dan itu MENOLAK
// DATA YANG SAH. Uji yang hanya memeriksa keberadaan UNIQUE lulus atas
// kekeliruan itu.
func TestKunciAlamiTabelAnak(t *testing.T) {
	sql := gabungan(t)
	mau := []struct{ inv, constraint, kolom string }{
		{"INV-05", "UQ_LAYER", "(ID_VERSI_KONTRAK, NOMOR_LAYER, BAGIAN_LAYER)"},
		{"INV-07", "UQ_MATA_UANG_KONTRAK", "(ID_VERSI_KONTRAK, KODE_MATA_UANG)"},
		{"INV-08", "UQ_RETENSI_CEDANT", "(ID_VERSI_KONTRAK, ID_KELOMPOK_TREATY, KODE_MATA_UANG)"},
		{"INV-09", "UQ_EGNPI", "(ID_VERSI_KONTRAK, ID_KELOMPOK_TREATY, KODE_MATA_UANG)"},
		{"INV-10", "UQ_PERIODE_PELAPORAN", "(ID_VERSI_KONTRAK, PERIODE)"},
		{"INV-11", "UQ_PERIODE_AKUMULASI", "(ID_VERSI_KONTRAK, PERIODE)"},
		{"INV-12", "UQ_TERMIN", "(ID_VERSI_KONTRAK, NOMOR_TERMIN, KODE_MATA_UANG)"},
		{"INV-13", "UQ_SKALA_KOASURANSI", "(ID_VERSI_KONTRAK, PERSEN_LIMIT)"},
		{"INV-14", "UQ_BATAS_PER_BAHAYA", "(ID_VERSI_KONTRAK, ID_BAHAYA)"},
		{"INV-66", "UQ_PORTOFOLIO", "(ID_VERSI_KONTRAK, ARAH_PORTOFOLIO, JENIS_PORTOFOLIO)"},
		{"INV-67", "UQ_DOKUMEN_KONTRAK", "(ID_VERSI_KONTRAK, ID_DOKUMEN)"},
		{"INV-64", "UQ_BAGIAN", "(ID_LAYER)"},
		{"INV-06", "UQ_DETAIL_PROPORSIONAL", "(ID_LAYER, ID_KELOMPOK_TREATY)"},
		// INV-15 terpasang sebagai DUA UNIQUE, satu per pelekatan - bukan satu
		// yang mencampur keduanya. Daftar periksa tiket 37 menuntutnya begitu.
		{"INV-15", "UQ_POTONGAN", "(ID_BAGIAN, ID_JENIS_POTONGAN)"},
		{"INV-15", "UQ_POTONGAN_2", "(ID_DETAIL_PROPORSIONAL, ID_JENIS_POTONGAN)"},
	}
	for _, m := range mau {
		if !strings.Contains(sql, "CONSTRAINT "+m.constraint+" UNIQUE "+m.kolom) {
			t.Errorf("%s: %s bukan UNIQUE %s", m.inv, m.constraint, m.kolom)
		}
	}
}

// Tabel tanpa UNIQUE, dan SEBABNYA BERBEDA-BEDA.
//
// ⛔ Ketiganya dulu dijelaskan dengan satu sebab yang sama - "peristiwa yang
// sama dapat terjadi dua kali". Itu KELIRU untuk dua dari tiga, dan
// kekeliruannya dibekukan ke dalam uji ini sebelum diperbaiki. Z00 memberi
// sebab yang berbeda untuk masing-masing:
//
//	JEJAK_PERUBAHAN    TIDAK ADA kunci alami, dan itu keputusan - peristiwa
//	                   yang sama dapat terjadi dua kali pada versi yang sama.
//	PEMULIHAN_LIMIT    kunci alaminya BELUM BERNOMOR; entitasnya baru
//	                   diterima §10.23c.
//	NILAI_MDP          Z00 TIDAK menyebutnya sama sekali. `KAMUS-KOLOM.md`
//	NILAI_MDP_MINIMUM  menyatakan KODE_MATA_UANG "kunci alami di dalam
//	                   induknya - belum bernomor". Jadi ia PUNYA kunci alami
//	                   yang belum bernomor, bukan sengaja tidak punya.
//
// Yang uji ini jaga hanya satu hal: tidak ada yang MENYISIPKAN UNIQUE diam-diam
// sebelum nomor invariannya turun. Ia tidak menyatakan ketiadaannya benar.
func TestTabelTanpaKunciAlamiTidakDiberiDiamDiam(t *testing.T) {
	sql := gabungan(t)
	for _, tabel := range []string{"JEJAK_PERUBAHAN", "NILAI_MDP", "NILAI_MDP_MINIMUM", "PEMULIHAN_LIMIT"} {
		if strings.Contains(sql, "CONSTRAINT UQ_"+tabel+" UNIQUE") {
			t.Errorf("%s diberi kunci alami tanpa nomor invarian; sebabnya berbeda per tabel - "+
				"lihat Z00_KUNCI_ALAMI.sql dan KAMUS-KOLOM.md sebelum memasangnya", tabel)
		}
	}
}

// Setiap kunci asing terlayani sebuah index.
//
// Oracle TIDAK membuat index untuk kunci asing dengan sendirinya. Kunci asing
// tanpa index membuat penghapusan baris INDUK memindai seluruh tabel anak, dan
// pada sebagian versi Oracle menguncinya.
//
// Kunci asing yang kolomnya MEMIMPIN sebuah UNIQUE sudah terlayani index milik
// UNIQUE itu - memasang index kedua di sana bukan keamanan, ia ongkos tulis.
// Yang didaftar di bawah adalah yang TIDAK demikian: kunci asing ke tabel
// acuan, dan anak yang tidak punya kunci alami.
func TestSetiapKunciAsingTerlayaniIndex(t *testing.T) {
	sql := gabungan(t)
	for _, ix := range []string{
		// kunci asing ke VERSI_KONTRAK/LAYER pada tabel tanpa kunci alami.
		// `IX_VERSI_KONTRAK_DASAR` TIDAK di sini: ia dipasang migrasi 440 modul
		// `treatyinadjustment`, dan uji ini hanya membaca migrasi modul ini.
		"IX_VERSI_KONTRAK_KONTRAK",
		"IX_NILAI_MDP_LAYER", "IX_NILAI_MDP_MIN_LAYER",
		"IX_PEMULIHAN_LIMIT_LAYER", "IX_JEJAK_PERUBAHAN_VERSI",
		// kunci asing ke tabel acuan - tidak pernah memimpin kunci alami
		"IX_MATA_UANG_KONTRAK_MU", "IX_RETENSI_CEDANT_KLP",
		"IX_EGNPI_KELOMPOK", "IX_EGNPI_KELAS_BISNIS", "IX_BATAS_PER_BAHAYA_BHY",
		"IX_DETAIL_PROP_KLP", "IX_POTONGAN_JENIS",
		"IX_NILAI_PB_BAGIAN", "IX_NILAI_PB_MIN_BAGIAN", "IX_NILAI_CAD_PREMI_DP",
	} {
		if !strings.Contains(sql, "CREATE INDEX {skema}."+ix+" ") {
			t.Errorf("index %s tidak ada; kunci asingnya tidak terlayani index mana pun", ix)
		}
	}
}

// Nama Oracle dibatasi 30 bita (§16).
//
// ⛔ Dibaca DARI TEKS DDL, bukan dari daftar literal di dalam uji. Uji yang
// mengukur literal yang ditulis tiga baris di atasnya tidak pernah dapat gagal;
// yang ini menangkap nama ke-24 yang ditulis orang berikutnya.
func TestNamaObjekDiBawahTigaPuluhBita(t *testing.T) {
	pola := regexp.MustCompile(`(?i)(?:CREATE\s+(?:TABLE|SEQUENCE|(?:UNIQUE\s+)?INDEX)\s+\{skema\}\.|CONSTRAINT\s+)(\w+)`)
	diperiksa := 0
	for nama, isi := range majuSaja(t) {
		for _, m := range pola.FindAllStringSubmatch(tanpaKomentar(isi), -1) {
			diperiksa++
			if len(m[1]) > 30 {
				t.Errorf("%s: pengenal %s %d bita, batas 30", nama, m[1], len(m[1]))
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol pengenal terbaca; pembacanya yang rusak")
	}
	t.Logf("%d pengenal diperiksa panjangnya", diperiksa)
}

// Yang TIDAK diuji di sini, sebab penjaga inti sudah menegakkannya atas SELURUH
// modul - menyalinnya ke sini berarti dua tempat yang harus sepakat:
//
//   jalur mundur tiap langkah   -> penjaga `TestSetiapLangkahPunyaJalurMundur`
//   nomor di rentang modulnya   -> penjaga `TestSetiapMigrasiDiRentangAtauSlotModulnya`
//   slot menu hanya menyentuh
//   baris menu modulnya sendiri -> penjaga `TestSlotMenuHanyaMenyentuhMenuModulnya`
//
// Yang DI SINI adalah yang khusus modul ini: invarian bernomor dan pernyataan
// keputusan yang penjaga umum tidak dapat mengetahuinya.
