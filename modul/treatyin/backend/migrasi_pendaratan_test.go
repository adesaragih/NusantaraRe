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
	"T_TREATY_REPORTING_PERIOD": 4548,
	"T_TREATY_PORTFOLIO":        1925,
	"T_TREATY_ACCUMULATION":     60,
	"T_TREATY_EGNPI":            2298,
	"T_TREATY_RETENTION":        2511,
	"T_TREATY_INSTALLMENT":      796,
	"T_TREATY_INSTALLMENT_ITEM": 3033,
	"T_VIEW_COMMENT":            11365,
	// Tabel kesembilan - migrasi 432, tab Co-Ins Scale.
	"M_TREATYIN_COINSCALE": 702,

	// ⭐ Tiga belas anak dari migrasi 437 — cacah disapu dari 1.854 dokumen.
	"T_TREATY_LIMITS":              4210,
	"T_TREATY_SHARE":               2923,
	"T_TREATY_RETRO_SHARE":         9,
	"T_TREATY_FAC_SHARE":           12,
	"T_TREATY_FAC_REINSURER":       15,
	"T_TREATY_LIMIT_DETAIL":        2868,
	"T_TREATY_LIMIT_GROUP":         8433,
	"T_TREATY_SHARE_SPREADING":     5550,
	"T_TREATY_SHARE_DEDUCTION":     2419,
	"T_TREATY_FAC_SHARE_DEDUCTION": 12,
	"T_TREATY_LIMIT_COB":           6547,
	"T_TREATY_LIMIT_ACHIEVEMENT":   2031,
	"T_TREATY_LIMIT_GROUP_COB":     15746,

	// ⭐ Tiga tabel NILAI dari migrasi 438 — cacah disapu atas 1.855 dokumen.
	"T_TREATY_LIMIT_AMOUNT":     11475,
	"T_TREATY_SHARE_AMOUNT":     6097,
	"T_TREATY_FAC_SHARE_AMOUNT": 24,

	// ⭐ Lima tabel dari migrasi 439 — yang menutup larangan JSONDATA.
	//
	// ⭐ Ketiga NOL DIISI 6 Oktober 2026. Sapuan atas SELURUH 1.855 dokumen
	// `M_TREATY_IN`, tanpa `ROWNUM` — penyebut yang SAMA dengan ke-27 entri
	// lain di tabel ini, dan sapuan itu mengukur ulang ke-27-nya untuk
	// membuktikannya (ke-27 cocok persis). Angka korpus `M_TREATY_IN_EDM`
	// disebut di `repository/pendaratan_peta.go`, tidak dijumlahkan ke sini.
	"T_TREATY_REVISION":      1855,
	"T_TREATY_LIMIT_MEASURE": 8849,
	"T_TREATY_LIMIT_SUMMARY": 2855,
	"T_TREATY_TOTAL":         6395,

	// ⭐ Migrasi 446 (parkir di folder `treatyinadjustment`), 7 Oktober
	// 2026 — tabel akar KEDUA, satu baris per dokumen seperti REVISION.
	"T_TREATY_HAZARD_LIMIT": 1855,

	// ⭐ Migrasi 448 (parkir), 7 Oktober 2026 — ringkasan tab Share
	// Non-Prop. Baru: nol baris korpus; diisi tombol Save.
	"T_TREATY_SHARE_SUMMARY": 0,

	// ⭐ Migrasi 449 (parkir), 8 Oktober 2026 — `Detail.SpreadingList` tab
	// Share Prop. Baru: nol baris; diisi tombol Save.
	"T_TREATY_LIMIT_SPREADING": 0,

	// ⭐ Migrasi 450 (parkir), 8 Oktober 2026 — Deduction, Parameter
	// Achievement, Reinstatement. Baru: nol baris; diisi tombol Save.
	"T_TREATY_LIMIT_DEDUCTION":     0,
	"T_TREATY_LIMIT_ACH_PARAM":     0,
	"T_TREATY_LIMIT_REINSTATEMENT": 0,
}

// Nama yang DIBUAT migrasi 430/432, untuk tiap nama yang dipakai hari ini.
//
// ⛔ Migrasi 430 dan 432 TIDAK disunting ketika nama tabelnya berganti; yang
// mengganti migrasi 436, dan riwayat yang disunting berbohong tentang apa
// yang pernah dijalankan. Akibatnya penjaga ini harus MENGIKUTI pergantian
// nama, bukan mencari nama hari ini di dalam DDL yang membuatnya.
var namaDDLPendaratan = map[string]string{
	"T_TREATY_REPORTING_PERIOD": "M_TREATYIN_REPORTINGPERIOD",
	"T_TREATY_PORTFOLIO":        "M_TREATYIN_PORTFOLIO",
	"T_TREATY_ACCUMULATION":     "M_TREATYIN_ACCUMULATION",
	"T_TREATY_EGNPI":            "M_TREATYIN_EGNPI",
	"T_TREATY_RETENTION":        "M_TREATYIN_RETENTION",
	"T_TREATY_INSTALLMENT":      "M_TREATYIN_INSTALLMENT",
	"T_TREATY_INSTALLMENT_ITEM": "M_TREATYIN_INSTALLMENTITEM",
	"T_VIEW_COMMENT":            "M_TREATYIN_COMMENT",
	// ⚠️ Kesembilan TIDAK berganti nama: `Diagram-Skema-Tabel-TreatyIn-dan-
	// EDM-v2.xlsx` nol padanan untuk Co-Ins Scale, disapu keenam lembarnya.
	"M_TREATYIN_COINSCALE": "M_TREATYIN_COINSCALE",

	// 437 membuatnya langsung dengan nama ini — nol pergantian nama.
	"T_TREATY_LIMITS":              "T_TREATY_LIMITS",
	"T_TREATY_SHARE":               "T_TREATY_SHARE",
	"T_TREATY_RETRO_SHARE":         "T_TREATY_RETRO_SHARE",
	"T_TREATY_FAC_SHARE":           "T_TREATY_FAC_SHARE",
	"T_TREATY_FAC_REINSURER":       "T_TREATY_FAC_REINSURER",
	"T_TREATY_LIMIT_DETAIL":        "T_TREATY_LIMIT_DETAIL",
	"T_TREATY_LIMIT_GROUP":         "T_TREATY_LIMIT_GROUP",
	"T_TREATY_SHARE_SPREADING":     "T_TREATY_SHARE_SPREADING",
	"T_TREATY_SHARE_DEDUCTION":     "T_TREATY_SHARE_DEDUCTION",
	"T_TREATY_FAC_SHARE_DEDUCTION": "T_TREATY_FAC_SHARE_DEDUCTION",
	"T_TREATY_LIMIT_COB":           "T_TREATY_LIMIT_COB",
	"T_TREATY_LIMIT_ACHIEVEMENT":   "T_TREATY_LIMIT_ACHIEVEMENT",
	"T_TREATY_LIMIT_GROUP_COB":     "T_TREATY_LIMIT_GROUP_COB",

	// 438 dan 439 juga membuatnya langsung dengan nama ini.
	"T_TREATY_LIMIT_AMOUNT":     "T_TREATY_LIMIT_AMOUNT",
	"T_TREATY_SHARE_AMOUNT":     "T_TREATY_SHARE_AMOUNT",
	"T_TREATY_FAC_SHARE_AMOUNT": "T_TREATY_FAC_SHARE_AMOUNT",
	"T_TREATY_REVISION":         "T_TREATY_REVISION",
	"T_TREATY_LIMIT_MEASURE":    "T_TREATY_LIMIT_MEASURE",
	"T_TREATY_LIMIT_SUMMARY":    "T_TREATY_LIMIT_SUMMARY",
	"T_TREATY_TOTAL":            "T_TREATY_TOTAL",

	// 446 membuatnya langsung dengan nama ini, di folder parkir.
	"T_TREATY_HAZARD_LIMIT": "T_TREATY_HAZARD_LIMIT",
	// 448 juga, di folder parkir.
	"T_TREATY_SHARE_SUMMARY": "T_TREATY_SHARE_SUMMARY",
	// 449 juga, di folder parkir.
	"T_TREATY_LIMIT_SPREADING": "T_TREATY_LIMIT_SPREADING",
	// 450 juga, di folder parkir.
	"T_TREATY_LIMIT_DEDUCTION":     "T_TREATY_LIMIT_DEDUCTION",
	"T_TREATY_LIMIT_ACH_PARAM":     "T_TREATY_LIMIT_ACH_PARAM",
	"T_TREATY_LIMIT_REINSTATEMENT": "T_TREATY_LIMIT_REINSTATEMENT",
}

func TestSembilanTabelPendaratanAdaDanBerkunciUtama(t *testing.T) {
	sql := gabunganDenganParkir(t)
	for nama := range tabelPendaratan {
		ddl, ok := namaDDLPendaratan[nama]
		if !ok {
			t.Errorf("tabel pendaratan %s tidak punya nama DDL-nya di penjaga ini", nama)
			continue
		}
		if !strings.Contains(sql, "CREATE TABLE {skema}."+ddl+" (") {
			t.Errorf("tabel pendaratan %s (DDL %s) tidak dibuat migrasi mana pun", nama, ddl)
			continue
		}
		// ⛔ Yang berganti nama WAJIB punya pernyataan penggantinya. Tanpa
		// pemeriksaan ini, salah ketik di peta di atas lolos tanpa suara.
		if ddl != nama && !strings.Contains(sql, "RENAME TO "+nama) {
			t.Errorf("%s dibuat sebagai %s tetapi nol migrasi menggantinya namanya", nama, ddl)
		}
		// ⛔ Nama batasannya DIBACA dari DDL, tidak lagi disusun dengan
		// menanggalkan awalan. Sejak migrasi 437 ada DUA keluarga nama -
		// `PK_MTI_*` untuk kedelapan tabel migrasi 430/432, `PK_TT_*` untuk
		// ketiga belas migrasi 437, dan pendeknya bukan potongan nama tabel
		// melainkan singkatan tersendiri (`LIMIT_ACHIEVE`, `FAC_SHARE_DED`).
		// Yang diuji keberadaan dan BENTUKNYA, bukan konvensi namanya.
		badan := badanCreateTable(t, nama)
		if !strings.Contains(badan, "PRIMARY KEY (ID)") {
			t.Errorf("INV-01: %s tanpa kunci utama atas kolom ID", nama)
		}
		mPK := regexp.MustCompile(`CONSTRAINT\s+(\w+)\s+PRIMARY KEY \(ID\)`).FindStringSubmatch(badan)
		mSeq := regexp.MustCompile(`Seq:\s*"(\w+)"`)
		_ = mSeq
		if mPK == nil {
			t.Errorf("INV-01: %s kunci utamanya tanpa nama", nama)
			continue
		}
		// Tiap tabel pendaratan WAJIB punya sequence-nya sendiri; namanya
		// mengikuti kunci utamanya dengan awalan yang sama.
		seq := strings.Replace(mPK[1], "PK_", "SEQ_", 1)
		if !strings.Contains(sql, "CREATE SEQUENCE {skema}."+seq+" ") {
			t.Errorf("INV-02: %s tanpa sequence %s", nama, seq)
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
	sql := tanpaKomentar(gabunganDenganParkir(t))
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
			// ⭐ LIMA kolom teks panjang — migrasi 439. Keduanya tetap teks
			// apa adanya; yang berbeda hanya WADAHNYA, sebab `VARCHAR2`
			// Oracle berhenti di 4.000 bita sementara teks tab
			// `Special Conditions` terukur 23.453 aksara. Kolom yang terlalu
			// pendek MEMOTONG tanpa bersuara.
			if teksPanjangPendaratan[nama+"."+kolom] {
				if !strings.HasPrefix(tipe, "CLOB") {
					t.Errorf("%s.%s mau CLOB (teks panjang), dapat %s", nama, kolom, tipe)
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
		// ⛔ Keberadaan `UNIQUE`-nya yang diuji, bukan namanya. Dua keluarga
		// nama hidup berdampingan sejak migrasi 437, dan menuntut satu pola
		// berarti penjaga ini menolak tabel yang sebenarnya benar.
		if !strings.Contains(badan, " UNIQUE (") {
			t.Errorf("%s tanpa batasan UNIQUE; pemuat idempoten menyandar padanya, dan "+
				"idempotensi yang hanya dijaga kode pemanggil benar sampai dua pemuat "+
				"berjalan bersamaan", nama)
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

// teksPanjangPendaratan - kolom pendaratan yang SENGAJA `CLOB`, bukan
// `VARCHAR2`. Daftar tertutup: kolom teks panjang baru harus disebut di sini
// beserta sebabnya, dan kolom yang tidak disebut tetap wajib teks pendek.
var teksPanjangPendaratan = map[string]bool{
	"T_TREATY_REVISION.EXCLUSIONS":         true,
	"T_TREATY_REVISION.EXCLUSIONSP":        true,
	"T_TREATY_REVISION.SPECIALCONDITIONS":  true,
	"T_TREATY_REVISION.SPECIALCONDITIONSP": true,
	// ⚠️ Ejaan KETIGA, huruf `p` kecil di dokumen (`SpecialConditionsp`),
	// berisi di 292 dokumen dan isinya BERBEDA dari saudaranya. Oracle tidak
	// membedakan besar-kecil nama kolom, jadi pembedanya pindah ke akhiran.
	"T_TREATY_REVISION.SPECIALCONDITIONSLC": true,
}

// Nama kunci JSON `Date` mendarat sebagai `TANGGAL`, dan kolom `DATE`
// tidak pernah lahir. `TestNolKataCadanganOracleSebagaiKolom` sudah
// menjaganya untuk seluruh repo; yang dijaga DI SINI pasangannya — bahwa
// penggantinya benar-benar ada, sehingga kuncinya tidak hilang diam-diam
// alih-alih berganti nama.
func TestKunciDateMendaratSebagaiTanggal(t *testing.T) {
	badan := badanCreateTable(t, "T_VIEW_COMMENT")
	if !strings.Contains(badan, "\n  TANGGAL ") {
		t.Error("T_VIEW_COMMENT tanpa kolom TANGGAL; kunci `Date` " +
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

	dibuang := kolomDibuang(t)
	ditambah := kolomDitambah(t)

	for _, p := range repository.PetaPendaratan {
		diDDL := map[string]bool{}
		// ⭐ Kolom yang migrasi BERIKUTNYA TAMBAHKAN lewat `ALTER … ADD`
		// tidak pernah muncul di `CREATE TABLE` mana pun — tetapi ia ADA di
		// basis data. Tanpa ini penjaga menuduh peta menyebut kolom yang
		// "tidak ada di DDL", padahal pemuatnyalah yang benar dan
		// penjaganya yang buta.
		for k := range ditambah[namaDDLPendaratan[p.Tabel]] {
			diDDL[k] = true
		}
		for _, baris := range strings.Split(badanCreateTable(t, p.Tabel), "\n") {
			ruas := strings.Fields(strings.TrimSpace(baris))
			if len(ruas) < 2 || ruas[0] == "CONSTRAINT" || ruas[0] == "REFERENCES" {
				continue
			}
			// ⛔ Kolom yang migrasi BERIKUTNYA buang tidak lagi ada di basis
			// data, walau `CREATE TABLE` yang melahirkannya masih menyebutnya.
			// Membacanya sebagai "ada" membuat penjaga ini menuntut peta
			// mengisi kolom yang sudah tidak ada.
			if !struktur[ruas[0]] && !dibuang[namaDDLPendaratan[p.Tabel]+"."+ruas[0]] {
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
			// ⭐ `JENIS` diisi PEMUAT, bukan dibaca dari dokumen — ia menandai
			// larik asal baris pada tabel yang menggabungkan beberapa larik
			// (migrasi 438 dan 439). Ia tidak boleh ada di `Kunci`.
			if k == "JENIS" {
				if len(p.LarikGabung) == 0 {
					t.Errorf("%s punya kolom JENIS tetapi petanya tidak menggabungkan larik mana pun", p.Tabel)
				}
				continue
			}
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
// ⭐ RALAT 6 Oktober 2026: `CurrencyList` dan `Limits` KINI di daftar ini.
// Keputusan pemilik proses melarang keras menarik nilai dari `JSONDATA`,
// jadi keduanya tidak lagi diurai — `Limits` dari tabel pendaratan
// (`pendaratan_layer.go`), `CurrencyList` dari tabel yang belum ada,
// sehingga gridnya kosong sampai tabel itu berdiri.
func TestLarikYangSudahPunyaTabelTidakDiuraiLagi(t *testing.T) {
	const berkas = "repository/warisan_kontrak.go"
	isi, err := os.ReadFile(berkas)
	if err != nil {
		t.Fatalf("membaca %s: %v", berkas, err)
	}
	teks := tanpaKomentarGo(string(isi))
	for _, larik := range []string{"ReportingPeriodList", "Portfolio", "AccumulationList",
		"CurrencyList", "Limits"} {
		tag := `json:"` + larik + `"`
		if strings.Contains(teks, tag) {
			t.Errorf("%s masih mengurai %s dari CLOB; ia sudah punya tabelnya sejak migrasi 430, "+
				"dan dua sumber untuk satu tab berarti salah satunya akan basi tanpa suara", berkas, larik)
		}
	}
	// ⭐ NOL medan tersisa sejak 6 Oktober 2026: berkas itu tidak lagi
	// mengurai dokumen sama sekali. Yang dijaga kini KEBALIKANNYA — ia tidak
	// boleh kembali mengurai apa pun.
	if strings.Contains(teks, "json.Unmarshal") {
		t.Error("jalur baca kontrak warisan mengurai dokumen lagi; " +
			"nilai dari JSONDATA dilarang keras sejak 6 Oktober 2026")
	}
	// ⛔ Dan penjaganya tetap MENGGIGIT: polanya dibuktikan atas teks tiruan.
	if !strings.Contains(`json:"ReportingPeriodList"`, "ReportingPeriodList") {
		t.Error("pembanding tag rusak; penjaga ini tidak menjaga apa pun")
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
		{"BacaPeriodePelaporan", "T_TREATY_REPORTING_PERIOD"},
		{"BacaPortofolio", "T_TREATY_PORTFOLIO"},
		{"BacaAkumulasi", "T_TREATY_ACCUMULATION"},
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
	sql := gabunganDenganParkir(t)
	// ⛔ DDL membuatnya dengan nama LAMA; migrasi 436 yang menggantinya.
	// Mencari nama hari ini di dalam `CREATE TABLE` karena itu selalu gagal
	// untuk kedelapan tabel yang berganti nama.
	if ddl, ok := namaDDLPendaratan[tabel]; ok {
		tabel = ddl
	}
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

// kolomDibuang mengumpulkan kolom yang migrasi mana pun CABUT, berkunci
// `TABEL.KOLOM` dengan nama tabel sebagaimana DDL menyebutnya.
//
// ⛔ Dihitung dari teks migrasinya, bukan didaftar tangan. Daftar tangan
// benar sampai kolom KEDUA dibuang, dan yang membuangnya tidak akan tahu
// penjaga ini ada.
func kolomDibuang(t *testing.T) map[string]bool {
	t.Helper()
	keluar := map[string]bool{}
	for _, baris := range strings.Split(gabungan(t), "\n") {
		ruas := strings.Fields(strings.TrimSpace(baris))
		// ALTER TABLE {skema}.NAMA DROP COLUMN KOLOM
		if len(ruas) < 6 || ruas[0] != "ALTER" || ruas[1] != "TABLE" ||
			ruas[3] != "DROP" || ruas[4] != "COLUMN" {
			continue
		}
		tabel := strings.TrimPrefix(ruas[2], "{skema}.")
		keluar[tabel+"."+ruas[5]] = true
	}
	return keluar
}

// kolomDitambah mengumpulkan kolom yang migrasi mana pun TAMBAHKAN lewat
// `ALTER TABLE {skema}.NAMA ADD ( … )`, berkunci nama tabel.
//
// ⛔ KEMBARAN `kolomDibuang`, dan ketiadaannya adalah lubang: penjaga ini
// sudah tahu kolom dapat DIBUANG sesudah `CREATE TABLE`, tetapi tidak tahu
// kolom dapat DITAMBAHKAN. Migrasi `444` yang pertama menambahkan, dan
// tanpa pembaca ini keenam belas kolomnya terbaca sebagai "peta menyebut
// kolom yang tidak ada di DDL" — penjaga merah yang menuduh kode yang benar.
//
// ⚠️ Bentuk yang dibaca BERKURUNG BANYAK BARIS:
//
//	ALTER TABLE {skema}.NAMA ADD (
//	  KOLOM  TIPE,
//	  KOLOM  TIPE
//	)
//
// Bentuk satu baris (`ADD KOLOM TIPE`) TIDAK dibaca, dan itu disengaja: nol
// migrasi memakainya hari ini, dan pembaca yang menebak dua bentuk lebih
// mudah salah diam-diam daripada yang menolak bentuk tak dikenal.
func kolomDitambah(t *testing.T) map[string]map[string]bool {
	t.Helper()
	masuk := map[string]map[string]bool{}
	var tabel string
	for _, baris := range strings.Split(gabungan(t)+"\n"+migrasiTetanggaParkir(t), "\n") {
		teks := strings.TrimSpace(baris)
		ruas := strings.Fields(teks)
		if len(ruas) >= 5 && ruas[0] == "ALTER" && ruas[1] == "TABLE" &&
			ruas[3] == "ADD" && ruas[4] == "(" {
			tabel = strings.TrimPrefix(ruas[2], "{skema}.")
			if masuk[tabel] == nil {
				masuk[tabel] = map[string]bool{}
			}
			continue
		}
		if tabel == "" {
			continue
		}
		if strings.HasPrefix(teks, ")") {
			tabel = ""
			continue
		}
		// `KOLOM TIPE…` — kolomnya ruas pertama, koma di ujung dibuang.
		if len(ruas) >= 2 && !strings.HasPrefix(teks, "--") {
			masuk[tabel][strings.TrimSuffix(ruas[0], ",")] = true
		}
	}
	return masuk
}

// migrasiTetanggaParkir membaca migrasi modul `treatyinadjustment` — hanya
// untuk `kolomDitambah`.
//
// ⛔⛔ INI KOMPROMI, DAN IA DICATAT SUPAYA TIDAK TERBACA SEBAGAI RANCANGAN.
//
// Rentang migrasi `treatyin` adalah `400-439` dan ia PENUH TANPA CELAH —
// keempat puluh nomornya terpakai. Migrasi `444`, yang menambah enam belas
// kolom ke `T_TREATY_REVISION` (tabel MILIK modul ini), karena itu terpaksa
// berdiri di folder `treatyinadjustment`: satu-satunya rentang sah yang
// memuat nomor `444` adalah `440-479` milik modul itu.
//
// Akibatnya `TestPetaPendaratanCocokDenganDDL` kehilangan DDL-nya — ia
// membaca `berkasMigrasi` modul ini saja, lalu menuduh peta menyebut enam
// belas kolom yang "tidak ada di DDL", padahal kolomnya ADA di basis data
// dan pemuat mengisinya.
//
// ⚠️ Ketiga pilihan, dan mengapa yang ketiga diambil:
//
//	biarkan merah          penjaga yang selalu merah berhenti dibaca orang
//	longgarkan penjaganya  ia berhenti menangkap kolom yang SUNGGUH hilang
//	baca folder tetangga   penjaga tetap benar; harganya satu jalur relatif
//
// ⛔ YANG SEHARUSNYA MENGGANTIKAN INI: jatah rentang baru untuk `treatyin`
// dari tim inti — prosedur yang `MODUL.md` sebut sendiri. Begitu ada, `444`
// dinomori ulang ke dalamnya, dipindahkan kembali, dan fungsi ini DIHAPUS.
//
// ⚠️ Ia membaca BERKAS, bukan mengimpor paket, jadi
// `TestModulTidakMengimporModulLain` tidak dilanggar. Folder yang tidak ada
// dijawab teks kosong — modul tetangga boleh tidak terpasang di pohon yang
// sedang diuji, dan itu bukan kegagalan penjaga ini.
// gabunganDenganParkir — migrasi modul ini DITAMBAH yang diparkir di folder
// `treatyinadjustment`. Dipakai pemeriksaan yang mencari `CREATE TABLE`
// tabel pendaratan: sejak `446` (`T_TREATY_HAZARD_LIMIT`) satu tabel
// pendaratan LAHIR di folder parkir, bukan hanya ditambah kolom. Uji
// invarian lain tetap membaca `gabungan` saja — folder tetangga punya
// penjaganya sendiri.
func gabunganDenganParkir(t *testing.T) string {
	t.Helper()
	return gabungan(t) + "\n" + migrasiTetanggaParkir(t)
}

func migrasiTetanggaParkir(t *testing.T) string {
	t.Helper()
	const dir = "../../treatyinadjustment/backend/migrations"
	entri, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	var b strings.Builder
	for _, e := range entri {
		if !strings.HasSuffix(e.Name(), ".sql") || strings.HasSuffix(e.Name(), "_down.sql") {
			continue
		}
		isi, err := os.ReadFile(dir + "/" + e.Name())
		if err != nil {
			t.Fatalf("membaca %s: %v", e.Name(), err)
		}
		b.Write(isi)
		b.WriteByte('\n')
	}
	return b.String()
}
