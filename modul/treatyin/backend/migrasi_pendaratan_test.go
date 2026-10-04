package backend

// Penjaga kedelapan tabel PENDARATAN tab Treaty In — migrasi 430 dan 431.
//
// ⛔ Tabel pendaratan punya aturan yang BERBEDA dari tabel model baru, dan
// perbedaannya tidak terlihat dari namanya. Berkas ini menuliskan perbedaan
// itu sebagai uji, sebab aturan yang hanya hidup di komentar kepala migrasi
// berhenti berlaku pada orang berikutnya yang menambah kolom.
//
// Ukurannya dari sapuan 3 Oktober 2026 atas SELURUH 1.854 dokumen
// `POOLDATA.M_TREATY_IN.JSONDATA`, diurai utuh sebagai JSON — nol dokumen
// gagal urai. Bukan contoh 300, bukan contoh 66.

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/repository"
)

// Kesembilan tabel, dan cacah BARIS yang sapuan temukan untuk masing-masing.
//
// ⚠️ Delapan sampai 3 Oktober 2026 siang; `M_TREATYIN_COINSCALE` menyusul
// sore harinya lewat migrasi 432. Angka di bawah naik bersama tabelnya.
// Cacah itu ditulis di sini supaya rekonsiliasi pemuat punya angka untuk
// diadu — bukan sekadar "ada isinya".
var tabelPendaratan = map[string]int{
	"M_TREATYIN_REPORTINGPERIOD": 4548,
	"M_TREATYIN_PORTFOLIO":       1925,
	"M_TREATYIN_ACCUMULATION":    60,
	"M_TREATYIN_EGNPI":           2298,
	"M_TREATYIN_RETENTION":       2511,
	"M_TREATYIN_INSTALLMENT":     796,
	"M_TREATYIN_INSTALLMENTITEM": 3033,
	"M_TREATYIN_COMMENT":         11365,
	// Tabel kesembilan - migrasi 432, tab Co-Ins Scale.
	"M_TREATYIN_COINSCALE": 702,
}

func TestSembilanTabelPendaratanAdaDanBerkunciUtama(t *testing.T) {
	sql := gabungan(t)
	for nama := range tabelPendaratan {
		if !strings.Contains(sql, "CREATE TABLE {skema}."+nama+" (") {
			t.Errorf("tabel pendaratan %s tidak dibuat migrasi mana pun", nama)
			continue
		}
		pendek := strings.TrimPrefix(nama, "M_TREATYIN_")
		if !strings.Contains(sql, "CONSTRAINT PK_MTI_"+pendek+" PRIMARY KEY (ID)") {
			t.Errorf("INV-01: %s tanpa kunci utama PK_MTI_%s", nama, pendek)
		}
		if !strings.Contains(sql, "CREATE SEQUENCE {skema}.SEQ_MTI_"+pendek+" ") {
			t.Errorf("INV-02: %s tanpa sequence SEQ_MTI_%s", nama, pendek)
		}
	}
}

// ⛔ YANG PALING MUDAH DILANGGAR DIAM-DIAM, dan sebabnya masuk akal: setiap
// tabel anak di modul ini merujuk induknya dengan kunci asing, jadi memasang
// satu lagi terasa seperti mengikuti pola.
//
// Ia TIDAK DAPAT dipasang. Diukur di POOLDATA 3 Oktober 2026:
//
//	all_constraints  TREATY_IN, tipe P atau U  -> NOL BARIS
//	all_indexes      INDEX_ID (ID)             -> NONUNIQUE
//	all_tab_columns  ID VARCHAR2(100)          -> NULLABLE = Y
//
// Oracle menolaknya dengan ORA-02270, dan satu-satunya cara memperbaikinya
// adalah DDL terhadap tabel warisan — yang dilarang. Uji ini gagal di sini,
// di komputer orang yang menulisnya, alih-alih gagal di Oracle pada orang
// lain beberapa hari kemudian.
func TestPendaratanTidakMerujukTabelWarisanDenganKunciAsing(t *testing.T) {
	warisan := []string{"TREATY_IN", "M_TREATY_IN", "M_TREATY_IN2", "TREATYEXCHANGEYEARLY", "M_TREATY_IN_DETAIL"}
	pola := regexp.MustCompile(`(?is)REFERENCES\s+\{skema\}\.(\w+)`)
	sql := tanpaKomentar(gabungan(t))
	for _, m := range pola.FindAllStringSubmatch(sql, -1) {
		for _, w := range warisan {
			if strings.EqualFold(m[1], w) {
				t.Errorf("kunci asing merujuk tabel WARISAN %s; Oracle menolaknya (ORA-02270) "+
					"selama kolomnya tidak memimpin kunci utama maupun UNIQUE, dan memberinya "+
					"satu berarti DDL terhadap tabel warisan", w)
			}
		}
	}
}

// Setiap kolom ISI bertipe teks — INV "simpan apa adanya".
//
// ⚠️ Yang dijaga BUKAN selera. Nilai di dalam `JSONDATA` seluruhnya string
// JSON, dan dua di antaranya membuktikan mahalnya menafsirkan:
// `Installment.AmountTotal` terpanjang 61 aksara (NUMBER(38,8) akan
// membulatkannya diam-diam), dan `Retention.Currency` memuat `1/04/2023` —
// tanggal di dalam kolom mata uang. Keduanya harus mendarat apa adanya
// supaya dapat ditemukan.
//
// Kolom STRUKTUR — `ID`, `IDINDUK`, `URUTAN` — justru harus angka: ia milik
// tabel ini, bukan milik dokumen.
func TestKolomIsiPendaratanBertipeTeks(t *testing.T) {
	struktur := map[string]bool{"ID": true, "IDINDUK": true, "URUTAN": true}
	diperiksa := 0
	for nama := range tabelPendaratan {
		badan := badanCreateTable(t, nama)
		for _, baris := range strings.Split(badan, "\n") {
			baris = strings.TrimSpace(baris)
			// Baris lanjutan definisi kunci asing bukan kolom.
			if baris == "" || strings.HasPrefix(baris, "CONSTRAINT") || strings.HasPrefix(baris, "REFERENCES") {
				continue
			}
			ruas := strings.Fields(baris)
			if len(ruas) < 2 {
				continue
			}
			kolom, tipe := ruas[0], ruas[1]
			diperiksa++
			if struktur[kolom] {
				if !strings.HasPrefix(tipe, "NUMBER(") {
					t.Errorf("%s.%s kolom struktur bertipe %s, mau NUMBER", nama, kolom, tipe)
				}
				continue
			}
			if !strings.HasPrefix(tipe, "VARCHAR2(") {
				t.Errorf("%s.%s bertipe %s; nilai JSON mendarat APA ADANYA sebagai teks — "+
					"menafsirkannya saat memuat membulatkan angka dan menolak nilai yang "+
					"salah bentuk, dan keduanya menghilangkan bukti", nama, kolom, tipe)
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol kolom terbaca; pembacanya yang rusak")
	}
	t.Logf("%d kolom pendaratan diperiksa", diperiksa)
}

// Ketiga kolom struktur ADA di kedelapan tabel.
//
// `URUTAN` yang hilang membuat larik kehilangan urutannya, dan larik Pega
// BERURUT — `Installment` nomor 1, 2, 3 bukan himpunan. Kehilangan itu tidak
// terlihat sampai seseorang membandingkan layar dengan sistem lama.
func TestPendaratanPunyaMasteridDanUrutan(t *testing.T) {
	for nama := range tabelPendaratan {
		badan := badanCreateTable(t, nama)
		for _, wajib := range []string{"ID ", "MASTERID ", "URUTAN "} {
			if !strings.Contains(badan, "\n  "+wajib) {
				t.Errorf("%s tanpa kolom %s", nama, strings.TrimSpace(wajib))
			}
		}
		pendek := strings.TrimPrefix(nama, "M_TREATYIN_")
		if !strings.Contains(badan, "CONSTRAINT UQ_MTI_"+pendek+" UNIQUE") {
			t.Errorf("%s tanpa UQ_MTI_%s; pemuat idempoten menyandar padanya, dan "+
				"idempotensi yang hanya dijaga kode pemanggil benar sampai dua pemuat "+
				"berjalan bersamaan", nama, pendek)
		}
	}
}

// ⛔ `AchievementLists` TIDAK boleh diam-diam ditambahkan sebagai tabel
// pendaratan kesembilan.
//
// Rancangan ronde ini mendaftarnya bersama kedelapan yang lain. Sapuan
// menemukan ia TIDAK ADA sebagai kunci puncak — nol kali di 1.854 dokumen.
// Jalurnya `Limits[].Detail[].AchievementLists` (1.961 kemunculan, 859
// kosong) dan `RevisionHistory[].Limits[].Detail[].AchievementLists` (3).
// Butirnya milik satu baris `Limits[].Detail[]`, yang hidup di
// `M_TREATY_IN2` — tabel warisan yang BERBEDA. Tabel berinduk `MASTERID`
// saja tidak dapat menyatakan butir itu milik baris limit yang mana, dan
// yang kehilangan induknya bukan sekadar kurang rapi: angka pencapaian yang
// menempel pada limit yang salah terbaca benar.
func TestAchievementBukanTabelPendaratan(t *testing.T) {
	sql := gabungan(t)
	if strings.Contains(sql, "CREATE TABLE {skema}.M_TREATYIN_ACHIEVEMENT") {
		t.Error("M_TREATYIN_ACHIEVEMENT dibuat sebagai tabel pendaratan berinduk MASTERID; " +
			"`AchievementLists` berinduk `Limits[].Detail[]`, bukan kontrak — lihat kepala " +
			"migrasi 430 dan docs/STRUKTUR-TABEL-TREATY-IN.md")
	}
}

// Nama kunci JSON `Date` mendarat sebagai `TANGGAL`, dan kolom `DATE`
// tidak pernah lahir. `TestNolKataCadanganOracleSebagaiKolom` sudah
// menjaganya untuk seluruh repo; yang dijaga DI SINI pasangannya — bahwa
// penggantinya benar-benar ada, sehingga kuncinya tidak hilang diam-diam
// alih-alih berganti nama.
func TestKunciDateMendaratSebagaiTanggal(t *testing.T) {
	badan := badanCreateTable(t, "M_TREATYIN_COMMENT")
	if !strings.Contains(badan, "\n  TANGGAL ") {
		t.Error("M_TREATYIN_COMMENT tanpa kolom TANGGAL; kunci `Date` " +
			"(11.365 kemunculan) tidak punya rumah")
	}
}

// ⛔ PETA DAN DDL TIDAK BOLEH BERSELISIH.
//
// `repository.PetaPendaratan` membangkitkan seluruh SQL sisip. Satu kolom
// di DDL yang tidak ada di peta adalah medan yang tidak pernah terisi;
// satu kolom di peta yang tidak ada di DDL adalah ORA-00904 yang baru
// muncul saat pemuatan sungguhan berjalan. Keduanya tidak terlihat sampai
// terlambat, dan keduanya ditangkap di sini tanpa koneksi Oracle.
func TestPetaPendaratanCocokDenganDDL(t *testing.T) {
	// Kolom struktur tidak ada di peta: ia milik tabel, bukan dokumen.
	struktur := map[string]bool{"ID": true, "IDINDUK": true, "MASTERID": true, "URUTAN": true}

	for _, p := range repository.PetaPendaratan {
		diDDL := map[string]bool{}
		for _, baris := range strings.Split(badanCreateTable(t, p.Tabel), "\n") {
			ruas := strings.Fields(strings.TrimSpace(baris))
			if len(ruas) < 2 || ruas[0] == "CONSTRAINT" || ruas[0] == "REFERENCES" {
				continue
			}
			if !struktur[ruas[0]] {
				diDDL[ruas[0]] = true
			}
		}
		diPeta := map[string]bool{}
		for _, k := range p.Kolom {
			diPeta[k] = true
			if !diDDL[k] {
				t.Errorf("%s: peta menyebut kolom %s yang TIDAK ada di DDL — "+
					"sisipnya akan gagal dengan ORA-00904 saat pemuatan sungguhan", p.Tabel, k)
			}
		}
		for k := range diDDL {
			if !diPeta[k] {
				t.Errorf("%s: DDL punya kolom %s yang TIDAK ada di peta — "+
					"ia tidak akan pernah terisi, dan nol galat akan menyebutkannya", p.Tabel, k)
			}
		}
	}
}

// Kedelapan nama tabel di peta sama dengan kedelapan yang berkas ini jaga.
// Dua daftar yang menyebut hal yang sama akan berselisih suatu hari; yang
// ini membuat hari itu gagal di sini.
func TestPetaDanPenjagaMenyebutTabelYangSama(t *testing.T) {
	if len(repository.PetaPendaratan) != len(tabelPendaratan) {
		t.Fatalf("peta %d tabel, penjaga %d", len(repository.PetaPendaratan), len(tabelPendaratan))
	}
	for _, p := range repository.PetaPendaratan {
		cacah, ada := tabelPendaratan[p.Tabel]
		if !ada {
			t.Errorf("peta menyebut %s yang penjaga ini tidak kenal", p.Tabel)
			continue
		}
		if p.CacahTerukur != cacah {
			t.Errorf("%s: peta mencatat %d baris terukur, penjaga %d", p.Tabel, p.CacahTerukur, cacah)
		}
	}
}

// ⛔ LARIK YANG SUDAH PUNYA TABELNYA TIDAK BOLEH DIURAI LAGI DARI CLOB.
//
// Ini penjaga terhadap kembalinya jalur lama. Menambahkan satu medan ke
// `jsonWarisan` adalah dua baris kerja dan selalu terasa paling mudah —
// dan sejak saat itu ada DUA sumber untuk satu tab, keduanya terlihat
// benar, dan yang satu membeku pada bentuk dokumen sementara yang lain
// ikut berubah bersama pemuatnya.
//
// ⚠️ `CurrencyList` SENGAJA tidak di daftar ini: ia belum punya tabel, sebab
// ronde pemindahan melarang membuat tabel Rate of Exchange. Begitu tabel
// kesembilan diputuskan, namanya masuk ke sini.
func TestLarikYangSudahPunyaTabelTidakDiuraiLagi(t *testing.T) {
	const berkas = "repository/warisan_kontrak.go"
	isi, err := os.ReadFile(berkas)
	if err != nil {
		t.Fatalf("membaca %s: %v", berkas, err)
	}
	teks := tanpaKomentarGo(string(isi))
	for _, larik := range []string{"ReportingPeriodList", "Portfolio", "AccumulationList"} {
		tag := `json:"` + larik + `"`
		if strings.Contains(teks, tag) {
			t.Errorf("%s masih mengurai %s dari CLOB; ia sudah punya tabelnya sejak migrasi 430, "+
				"dan dua sumber untuk satu tab berarti salah satunya akan basi tanpa suara", berkas, larik)
		}
	}
	// ⛔ Penjaganya MENGGIGIT: polanya harus menemukan yang memang masih ada.
	if !strings.Contains(teks, `json:"CurrencyList"`) {
		t.Error("pola tidak menemukan CurrencyList yang masih diurai; pembacanya yang rusak, " +
			"atau CurrencyList sudah pindah dan daftar di atas perlu diperbarui")
	}
}

// Pasangannya: pembacanya BENAR-BENAR ada, dan membaca dari tabel yang
// migrasi 430 buat. Tanpa ini, menghapus medan dari `jsonWarisan` tanpa
// menuliskan penggantinya lulus uji di atas dengan sempurna — dan ketiga
// tab menjadi kosong selamanya.
func TestPembacaTabPendaratanAda(t *testing.T) {
	const berkas = "repository/pendaratan_baca.go"
	isi, err := os.ReadFile(berkas)
	if err != nil {
		t.Fatalf("membaca %s: %v", berkas, err)
	}
	teks := string(isi)
	for _, pasang := range []struct{ fungsi, tabel string }{
		{"BacaPeriodePelaporan", "M_TREATYIN_REPORTINGPERIOD"},
		{"BacaPortofolio", "M_TREATYIN_PORTFOLIO"},
		{"BacaAkumulasi", "M_TREATYIN_ACCUMULATION"},
	} {
		if !strings.Contains(teks, "func (g *Gudang) "+pasang.fungsi+"(") {
			t.Errorf("%s tidak ada di %s", pasang.fungsi, berkas)
		}
		if !strings.Contains(teks, `"`+pasang.tabel+`"`) {
			t.Errorf("%s tidak menyebut tabel %s", berkas, pasang.tabel)
		}
	}
	// ⛔ `ORDER BY URUTAN` — larik Pega BERURUT, dan Oracle tanpa ORDER BY
	// bebas mengembalikan baris dalam urutan apa pun.
	if !strings.Contains(teks, "ORDER BY URUTAN") {
		t.Error("pembacanya tanpa ORDER BY URUTAN; urutan larik hilang, dan hilangnya " +
			"tidak terlihat sampai seseorang membandingkan layar dengan sistem lama")
	}
}

// badanCreateTable memotong isi satu CREATE TABLE, dari kurung buka sampai
// kurung tutup terakhir sebelum pembatas pernyataan.
func badanCreateTable(t *testing.T, tabel string) string {
	t.Helper()
	sql := gabungan(t)
	awal := strings.Index(sql, "CREATE TABLE {skema}."+tabel+" (")
	if awal < 0 {
		t.Fatalf("CREATE TABLE %s tidak ditemukan", tabel)
	}
	sisa := sql[awal:]
	// ⛔ Kepala pernyataannya DIBUANG. Membiarkannya membuat `CREATE TABLE`
	// terbaca sebagai kolom bernama CREATE bertipe TABLE — yang lolos
	// sebagai kolom aneh alih-alih gagal sebagai pembaca yang rusak.
	buka := strings.Index(sisa, "(\n")
	if buka < 0 {
		t.Fatalf("kurung buka CREATE TABLE %s tidak ditemukan", tabel)
	}
	sisa = sisa[buka+1:]
	akhir := strings.Index(sisa, "\n)\n")
	if akhir < 0 {
		t.Fatalf("akhir CREATE TABLE %s tidak ditemukan", tabel)
	}
	return sisa[:akhir]
}

// ⛔ Ketiga berkas pertanyaan TIDAK boleh hilang, dan harus menyatakan
// sudah dijawab beserta nomor keputusannya.
//
// Sebabnya bukan kerapian. Ukuran di dalamnya — "303 dokumen, nol yang
// isinya identik", "2 kontrak dari 1.854", "23.453 aksara" — adalah DASAR
// ketiga keputusan. Menghapusnya meninggalkan keputusan yang tidak dapat
// diadu dengan apa pun, dan keputusan semacam itu akan dibatalkan oleh
// orang berikutnya yang tidak tahu sebabnya.
func TestBerkasPertanyaanTetapAdaDanMenunjukKeputusannya(t *testing.T) {
	for _, p := range []struct{ berkas, keputusan string }{
		{"../docs/PERTANYAAN-TERBUKA-PERSEN-SHARE.md", "§13"},
		{"../docs/PERTANYAAN-TERBUKA-RETRO.md", "§14"},
		{"../docs/PERTANYAAN-TERBUKA-TAB-TEKS.md", "§15"},
	} {
		isi, err := os.ReadFile(p.berkas)
		if err != nil {
			t.Errorf("%s hilang: %v — ukurannya adalah dasar keputusan %s", p.berkas, err, p.keputusan)
			continue
		}
		teks := string(isi)
		if !strings.Contains(teks, "SUDAH DIJAWAB") {
			t.Errorf("%s tidak menyatakan sudah dijawab", p.berkas)
		}
		if !strings.Contains(teks, p.keputusan) {
			t.Errorf("%s tidak menunjuk keputusan %s", p.berkas, p.keputusan)
		}
		if !strings.Contains(teks, "KEPUTUSAN-PENYELARASAN-REPO.md") {
			t.Errorf("%s tidak menunjuk berkas keputusannya", p.berkas)
		}
	}

	// Dan ketiga keputusannya benar-benar ADA, dengan syarat pembalikannya.
	kep, err := os.ReadFile("../docs/KEPUTUSAN-PENYELARASAN-REPO.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, judul := range []string{"## 13 ·", "## 14 ·", "## 15 ·"} {
		if !strings.Contains(string(kep), judul) {
			t.Errorf("keputusan %q tidak ada", judul)
		}
	}
	// ⛔ Syarat pembalikan WAJIB disebut — §12 menetapkan bentuknya, dan
	// keputusan tanpa syarat pembalikan adalah keputusan yang tidak dapat
	// dicabut tanpa membongkar ulang seluruh alasannya.
	if n := strings.Count(string(kep), "Pembalikan"); n < 4 {
		t.Errorf("hanya %d keputusan menyebut pembalikannya; §12 sampai §15 seluruhnya wajib", n)
	}
}
