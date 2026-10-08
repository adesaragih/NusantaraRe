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
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
	maju := majuSaja(t)
	nama := make([]string, 0, len(maju))
	for n := range maju {
		nama = append(nama, n)
	}
	// ⛔ DIURUTKAN. Peta Go tidak berurut, dan urutan menentukan siapa menimpa
	// siapa: migrasi 420 membongkar lalu membuat ulang kunci asing yang
	// 403-418 definisikan. Menggabungkan dalam urutan acak membuat uji
	// perilaku hapus hijau atau merah bergantung iterasi peta.
	sort.Strings(nama)
	var b strings.Builder
	for _, n := range nama {
		b.WriteString(maju[n])
		b.WriteByte('\n')
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
		"CK_PENYEBARAN_INDUK": "ERD.md §2.5 - induk polimorfik yang SAMA bentuknya dengan " +
			"CK_POTONGAN_INDUK: tepat satu dari ID_BAGIAN dan ID_DETAIL_PROPORSIONAL terisi. " +
			"Menyatakan BENTUK baris, bukan daftar nilai, jadi ADR-0038 tidak kena",
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

// INV-18 - perilaku hapus tiap kunci asing DITETAPKAN SADAR, dan di sini ia
// diadu dengan sumber yang menetapkannya.
//
// ⚠️ RALAT 2 Oktober 2026. Uji ini pernah bernama TestNolKaskadeHapus dan
// menuntut NOL `ON DELETE` di seluruh modul, dengan alasan "kebijakan modul
// ini MENOLAK". Itu membaca INV-18 terbalik: invariannya menuntut keputusan
// yang DINYATAKAN, bukan bawaan yang kebetulan cocok — dan keputusannya sudah
// ada, mengikat, di `4-erd-dan-tabel-datar/ERD.md` §2. Uji lama membekukan
// kekeliruan itu dan akan menolak setiap perbaikannya.
//
// Tabel di bawah DISALIN dari ERD.md §2, satu baris per kunci asing yang
// modul ini punya. Relasi yang tabelnya belum dibuat tidak ada di sini.
func TestPerilakuHapusSesuaiERD(t *testing.T) {
	// "" berarti TOLAK: diwujudkan dengan TIDAK menulis klausa ON DELETE.
	const tolak = ""
	mau := map[string]string{
		// §2.1 tulang punggung
		"FK_VERSI_KONTRAK_1":      tolak,
		"FK_KONTRAK_DISALIN_DARI": "ON DELETE SET NULL",
		// §2.3 anak langsung versi — sepuluh ikut hapus, dua tolak
		//
		// ⛔ `FK_MATA_UANG_KONTRAK_1` DICABUT 4 Oktober 2026 bersama tabelnya
		// (migrasi 434). Entri yang tertinggal di sini akan LULUS tanpa
		// berarti apa-apa: uji ini mencocokkan teks `CREATE` di migrasi 401
		// dan 403, dan `DROP` di 434 tidak terlihat olehnya. Penjaga yang
		// menagih kunci asing yang tidak ada adalah penjaga yang berhenti
		// menggigit tanpa menjadi merah.
		"FK_RETENSI_CEDANT_1":    "ON DELETE CASCADE",
		"FK_EGNPI_1":             "ON DELETE CASCADE",
		"FK_PORTOFOLIO_1":        "ON DELETE CASCADE",
		"FK_PERIODE_PELAPORAN_1": "ON DELETE CASCADE",
		"FK_PERIODE_AKUMULASI_1": "ON DELETE CASCADE",
		"FK_TERMIN_1":            "ON DELETE CASCADE",
		"FK_SKALA_KOASURANSI_1":  "ON DELETE CASCADE",
		"FK_BATAS_PER_BAHAYA_1":  "ON DELETE CASCADE",
		"FK_DOKUMEN_KONTRAK_1":   "ON DELETE CASCADE",
		"FK_JEJAK_PERUBAHAN_1":   tolak,
		// §2.4 cabang
		"FK_LAYER_1":               "ON DELETE CASCADE",
		"FK_DETAIL_PROPORSIONAL_1": "ON DELETE CASCADE",
		"FK_BAGIAN_1":              "ON DELETE CASCADE",
		// §2.5 potongan — dua pelekatan, keduanya ikut hapus
		"FK_POTONGAN_1": "ON DELETE CASCADE",
		"FK_POTONGAN_2": "ON DELETE CASCADE",
		// ----------------------------------------------------------------
		// Entitas rekonsiliasi ERD, 2 Oktober 2026 — migrasi 422–425.
		//
		// ⛔ Keenam relasi di bawah TIDAK ada di `ERD.md` §2: ia lahir sesudah
		// §2 ditulis, sama seperti delapan relasi §2.3c. Perilakunya TIDAK
		// dikarang — ia diambil dari DUA sumber yang sepakat:
		//
		//   - `ERD-TREATY-IN-DAN-EDM.html` (ACUAN struktur sejak keputusan
		//     pemilik proses 2 Okt 2026) baris relasi 4, 9, 10, 27, 39;
		//   - aturan yang kelompoknya sendiri sudah nyatakan di `ERD.md`:
		//     §2.5 anak cabang ikut hapus, §2.7 rujukan acuan tolak,
		//     §2.3 anak langsung baris versi ikut hapus.
		//
		// ⚠️ Kolom `ON DELETE` di ERD HTML menyatakan dirinya "USULAN
		// rancangan ... bukan perilaku sistem lama", jadi ia tidak dipakai
		// sendirian. Di keenam baris ini kedua sumber sepakat.
		"FK_KELAS_BISNIS_LAYER_1":    "ON DELETE CASCADE", // ERD baris 4 · ERD.md §2.5
		"FK_KELAS_BISNIS_LAYER_2":    tolak,               // rujukan acuan · ERD.md §2.7
		"FK_KELOMPOK_LAYER_1":        "ON DELETE CASCADE", // ERD baris 9 · ERD.md §2.5
		"FK_KELOMPOK_LAYER_2":        tolak,               // rujukan acuan · ERD.md §2.7
		"FK_KELAS_BISNIS_KELOMPOK_1": "ON DELETE CASCADE", // ERD baris 10
		"FK_KELAS_BISNIS_KELOMPOK_2": tolak,               // rujukan acuan · ERD.md §2.7
		"FK_RINCIAN_ANGSURAN_1":      "ON DELETE CASCADE", // ERD baris 27 · ERD.md §2.3
		// §2.9 — MENGIKAT, dan satu-satunya tempat ERD.md dan ERD HTML
		// BERSELISIH. §2.9 menulis `KONTRAK 1--o< PENCAPAIAN [hapus: tolak]`;
		// ERD HTML baris 6 menulis induk `LIMIT_DETAIL` + CASCADE dengan bukti
		// yang ia sendiri tandai DAUN-RELATIF (lemah). Sebabnya di kepala
		// migrasi 423.
		"FK_PENCAPAIAN_1": tolak,
		// ERD baris 39 menulis ON DELETE-nya "di Go" — ia TIDAK meresepkan
		// aturan tingkat basis data. Di sini diwujudkan `tolak`: arsip yang
		// lenyap bersama kontraknya berhenti menjadi arsip. Sebabnya di kepala
		// migrasi 425.
		"FK_ARSIP_MUATAN_KELUAR_1": tolak,
		// ⛔ `FK_NILAI_SELISIH_1` dan `FK_NILAI_SEBELUM_PRO_RATE_1` PINDAH
		// ke `modul/treatyinadjustment` 4 Oktober 2026 bersama tabelnya —
		// KEPUTUSAN §19, migrasi 435 mencabut dan 442 membangun ulang.
		//
		// Barisnya IKUT PINDAH, tidak digandakan: dua tempat yang
		// menyatakan perilaku hapus yang sama akan berselisih, dan yang
		// salah satunya basi tidak akan berbunyi. Keduanya kini ada di
		// `modul/treatyinadjustment/backend/migrasi_invarian_test.go`.
		//
		// ⚠️ Yang TIDAK ikut pindah: alasan rancangannya. Kepala migrasi
		// `426` tetap di modul ini — kardinalitas 1:N, penghalang
		// `BESARAN_DAPAT_DISESUAIKAN`, dan INV-58 — dan `442` menunjuk
		// balik ke sana alih-alih menyalinnya.
		// §2.3 — MENGIKAT, dan ia SATU-SATUNYA anak VERSI_KONTRAK di §2.3
		// yang `tolak`; sepuluh lainnya `ikut hapus`. Sebabnya dikutip utuh
		// di kepala migrasi 427: "jejak yang dapat dihapus bersama bendanya
		// bukan jejak."
		"FK_CATATAN_PERSETUJUAN_1": tolak,
		// Perkakas pemindahan — migrasi 428. TIDAK ada di `ERD.md` §2 sebab
		// ia potret sistem LAMA dan ketiganya perkakas sistem BARU.
		//
		// ⚠️ ERD HTML baris 38 menulis ON DELETE-nya "di Go" — ia tidak
		// meresepkan aturan tingkat basis data. Dipilih `tolak`, dengan
		// alasan yang sudah dipakai migrasi 425: catatan forensik yang
		// lenyap bersama induknya berhenti menjadi catatan forensik.
		// Ini MENGOREKSI tiket 73, yang menyebut "ikut hapus".
		"FK_MIGRASI_NILAI_DITOLAK_1": tolak,
		// §2.5 cabang penyebaran — migrasi 429. Empat `ikut hapus` dan dua
		// `tolak`, seluruhnya dikutip: §2.5 untuk rantai induknya, §2.7
		// untuk rujukan tabel acuan.
		"FK_PENYEBARAN_1":         "ON DELETE CASCADE", // BAGIAN, §2.5
		"FK_PENYEBARAN_2":         "ON DELETE CASCADE", // DETAIL_PROPORSIONAL, §2.5
		"FK_PENYEBARAN_3":         tolak,               // JENIS_REASURANSI, §2.7
		"FK_RINCIAN_PENYEBARAN_1": "ON DELETE CASCADE", // PENYEBARAN, §2.5
		"FK_RINCIAN_PENYEBARAN_2": tolak,               // JENIS_REASURANSI, §2.7
		"FK_NILAI_PENYEBARAN_1":   "ON DELETE CASCADE", // RINCIAN_PENYEBARAN, §2.5
		// §2.3b dan §2.3c — tabel anak paket uang, seluruhnya ikut hapus
		"FK_PEMULIHAN_LIMIT_1":       "ON DELETE CASCADE",
		"FK_NILAI_MDP_1":             "ON DELETE CASCADE",
		"FK_NILAI_MDP_MINIMUM_1":     "ON DELETE CASCADE",
		"FK_NILAI_PREMI_BRUTO_1":     "ON DELETE CASCADE",
		"FK_NILAI_PREMI_BRUTO_MIN_1": "ON DELETE CASCADE",
		"FK_NILAI_CADANGAN_PREMI_1":  "ON DELETE CASCADE",
		// §2.7 tabel acuan — seluruhnya tolak
		//
		// ⛔ `FK_MATA_UANG_KONTRAK_2` dan `FK_VERSI_KONTRAK_MATA_UANG`
		// DICABUT 4 Oktober 2026 bersama `MATA_UANG` dan
		// `MATA_UANG_KONTRAK` (migrasi 434) — kurs dan daftar mata uang kini
		// dibaca dari `TREATYEXCHANGEYEARLY`. Lihat sebabnya di kepala
		// migrasi itu, termasuk kewajiban tiket 57 yang ikut terbuka.
		"FK_RETENSI_CEDANT_2":      tolak,
		"FK_EGNPI_2":               tolak,
		"FK_EGNPI_3":               tolak,
		"FK_BATAS_PER_BAHAYA_2":    tolak,
		"FK_DETAIL_PROPORSIONAL_2": tolak,
		"FK_POTONGAN_3":            tolak,
		// Tidak ada di ERD §2 — diputuskan di migrasi 400, ditagih ke pemilik ERD.
		"FK_JENIS_REASURANSI_INDUK": tolak,
		// ----------------------------------------------------------------
		// Tabel PENDARATAN tab Treaty In — migrasi 430, 3 Oktober 2026.
		//
		// ⛔ SATU-SATUNYA kunci asing di kedelapan tabel itu, dan ia tidak
		// punya baris ERD. Kedelapannya mendaratkan larik di dalam
		// `M_TREATY_IN.JSONDATA`; ERD menggambar entitas sistem lama, bukan
		// tabel pendaratan, jadi mengutip baris ERD untuknya akan mengarang
		// sumber. Yang mengikat bentuk datanya: `InstallmentList` bersarang
		// DI DALAM elemen `Installment` — diukur pada ke-1.854 dokumen, 796
		// elemen induk dan 3.033 butir anak — dan butir angsuran tanpa
		// terminnya tidak berarti apa pun. Migrasi `424_` memodelkan dua
		// tingkat yang sama persis dengan alasan yang sama.
		//
		// ⚠️ Yang TIDAK ada di sini jauh lebih banyak: delapan `MASTERID`
		// yang rancangannya sebut sebagai kunci asing ke `TREATY_IN.ID`.
		// Kedelapannya TIDAK dibuat, sebab `POOLDATA.TREATY_IN` tidak punya
		// kunci utama maupun UNIQUE pada `ID` (ia bahkan NULLABLE), dan
		// Oracle menolak merujuk kolom semacam itu dengan ORA-02270.
		// `MODUL.md` bab kaskade dan `KEPUTUSAN-PENYELARASAN-REPO.md` §12.
		// ⚠ Nama DDL-nya, bukan nama hari ini: migrasi 436 menggantinya
		// menjadi `FK_TT_INSTALLMENT_ITEM_1`, dan penjaga ini membaca DDL.
		"FK_MTI_INSTALLMENTITEM_1": "ON DELETE CASCADE",

		// ⭐ DELAPAN kunci asing migrasi 437 — pohon anak `TREATY_IN` menurut
		// `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`, keputusan pemilik
		// proses 6 Oktober 2026 ("puncaknya tetap TREATY_IN, anak-anaknya
		// mengikuti xlsx"). Tiap satu menyebut sel asalnya.
		//
		// ⛔ Tautan ke AKAR tidak ada di sini, dan itu disengaja: xlsx menulis
		// `FK TREATY_IN_ID -> TREATY_IN.ID` pada tiap anak tingkat pertama,
		// tetapi `TREATY_IN` nol kunci utama dan nol kunci unik — `ORA-02270`.
		// Anak tingkat pertama karena itu menaut dengan `MASTERID` teks.
		// Kedelapan di bawah seluruhnya antar tabel MILIK modul ini, jadi
		// keduanya dapat dipasang sungguhan.
		"FK_TT_LIMIT_DETAIL":    "ON DELETE CASCADE", // K23  -> T_TREATY_LIMITS
		"FK_TT_LIMIT_GROUP":     "ON DELETE CASCADE", // K40  -> T_TREATY_LIMITS
		"FK_TT_LIMIT_COB":       "ON DELETE CASCADE", // O27  -> T_TREATY_LIMIT_DETAIL
		"FK_TT_LIMIT_ACHIEVE":   "ON DELETE CASCADE", // O35  -> T_TREATY_LIMIT_DETAIL
		"FK_TT_LIMIT_GRP_COB":   "ON DELETE CASCADE", // O44  -> T_TREATY_LIMIT_GROUP
		"FK_TT_SHARE_SPREADING": "ON DELETE CASCADE", // K54  -> T_TREATY_SHARE
		"FK_TT_SHARE_DEDUCTION": "ON DELETE CASCADE", // K67  -> T_TREATY_SHARE
		"FK_TT_FAC_SHARE_DED":   "ON DELETE CASCADE", // K80  -> T_TREATY_FAC_SHARE

		// ⭐ Migrasi `438` — tiga tabel NILAI, dan ketiganya menaut tabel
		// MILIK modul ini, jadi kunci asingnya nyata seperti kedelapan di
		// atas. Dibangun sesudah cacah elemennya diukur atas SELURUH 1.855
		// dokumen: 11.475 · 6.097 · 24.
		"FK_TT_LIMIT_AMOUNT":     "ON DELETE CASCADE", // -> T_TREATY_LIMIT_DETAIL
		"FK_TT_SHARE_AMOUNT":     "ON DELETE CASCADE", // -> T_TREATY_SHARE
		"FK_TT_FAC_SHARE_AMOUNT": "ON DELETE CASCADE", // -> T_TREATY_FAC_SHARE

		// ⭐ Migrasi 439 — satu-satunya kunci asing barunya. `MDPList`,
		// `PremiumEarnedList`, dan `EgnpiTotalList` hidup di dalam elemen
		// `Limits[]`, jadi induknya `T_TREATY_LIMITS`, bukan `_DETAIL`.
		// Alasan kaskadenya sama dengan ketiga di atas: besaran sebuah layer
		// tanpa layernya tidak berarti apa pun.
		"FK_TT_LIMIT_MEASURE": "ON DELETE CASCADE", // -> T_TREATY_LIMITS
	}

	sql := gabungan(t)
	// Potong tiap definisi kunci asing sampai akhir baris REFERENCES-nya.
	pola := regexp.MustCompile(`(?is)CONSTRAINT\s+(FK_\w+)\s+FOREIGN KEY\s*\([^)]*\)\s*REFERENCES\s+\{skema\}\.\w+\s*\([^)]*\)([^,\n]*)`)
	// ⛔ YANG TERAKHIR MENANG. Migrasi 420 membongkar lalu membuat ulang kunci
	// asing yang 403-418 definisikan, jadi bentuk yang BERLAKU adalah definisi
	// terakhir di urutan migrasi - bukan yang pertama ditemukan. Membaca yang
	// pertama membuat uji ini menilai bentuk yang sudah diganti.
	akhir := map[string]string{}
	urut := []string{}
	for _, m := range pola.FindAllStringSubmatch(tanpaKomentar(sql), -1) {
		nama := m[1]
		if _, pernah := akhir[nama]; !pernah {
			urut = append(urut, nama)
		}
		akhir[nama] = strings.ToUpper(strings.TrimSpace(m[2]))
	}
	// ⛔ YANG DICABUT TIDAK DINILAI. Migrasi 434 membuang `MATA_UANG` dan
	// `MATA_UANG_KONTRAK` beserta `FK_VERSI_KONTRAK_MATA_UANG`; teks
	// `CREATE`-nya tetap ada di migrasi 401 dan 403, sebab migrasi tidak
	// pernah disunting mundur. Tanpa langkah ini uji menagih kunci asing
	// yang tidak akan ada di basis data mana pun - dan sebelum 4 Oktober
	// 2026 ia LULUS menagihnya, yaitu berhenti menggigit tanpa menjadi
	// merah.
	for nama := range akhir {
		if dicabut(tanpaKomentar(sql), nama) {
			delete(akhir, nama)
		}
	}
	hidup := urut[:0:0]
	for _, nama := range urut {
		if _, ada := akhir[nama]; ada {
			hidup = append(hidup, nama)
		}
	}
	urut = hidup

	lihat := map[string]bool{}
	for _, nama := range urut {
		harap, dikenal := mau[nama]
		if !dikenal {
			t.Errorf("kunci asing %s tidak ada di tabel ERD uji ini; "+
				"tiap kunci asing BARU wajib menyebut keputusan ERD.md §2-nya", nama)
			continue
		}
		lihat[nama] = true
		if akhir[nama] != strings.ToUpper(harap) {
			t.Errorf("%s: perilaku hapus %q, ERD.md menuntut %q", nama, akhir[nama], harap)
		}
	}
	for nama := range mau {
		if !lihat[nama] {
			t.Errorf("kunci asing %s didaftar uji ini tetapi tidak ada di DDL", nama)
		}
	}
}

// Dua puluh lima relasi IKUT HAPUS, dan cacahnya dijaga.
//
// ERD.md §2 menyatakan 28 relasi `ikut hapus`; tujuh di antaranya menyentuh
// tabel yang modul ini BELUM buat (PENYEBARAN, RINCIAN_PENYEBARAN,
// NILAI_PENYEBARAN, NILAI_SELISIH, PERISTIWA_KONTRAK, RETRO_KELUAR, dan
// PENYEBARAN cabang kedua). Angka di bawah naik bersama tabelnya.
//
// 3 Oktober 2026: 27 -> 31. Cabang penyebaran masuk bersama migrasi 429 —
// empat relasi §2.5 — sehingga yang "tabelnya belum dibuat" turun dari lima
// menjadi satu (`RETRO_KELUAR`, gelombang 2). Sebelumnya 21 -> 27:
// `NILAI_SELISIH` (§2.6) dan `NILAI_SEBELUM_PRO_RATE` bersama migrasi 426,
// dan empat lagi bersama migrasi 422 dan 424,
// dan keempatnya BUKAN dari §2 melainkan dari `ERD-TREATY-IN-DAN-EDM.html`
// baris 4, 9, 10, 27 — lihat tabel di `TestPerilakuHapusSesuaiERD`. Penyebut
// 28 TIDAK ikut naik: ia cacah §2, dan §2 belum memuat keempatnya.
//
// 3 Oktober 2026, kedua kalinya: 31 -> 32. Yang ke-32 `FK_MTI_INSTALLMENTITEM_1`
// dari migrasi 430, dan ia tidak bersumber dari ERD sama sekali — tabel
// PENDARATAN tidak digambar di sana. Penyebut 28 tetap, dan pembilang
// "bukan dari §2" naik dari 4 menjadi 5.
func TestCacahKaskadeSesuaiTabelYangAda(t *testing.T) {
	// 6 Oktober 2026, ketiga kalinya: 40 -> 43. Ketiga yang baru dari migrasi
	// `438` (`LIMIT_AMOUNT`, `SHARE_AMOUNT`, `FAC_SHARE_AMOUNT`), dan seperti
	// `437` ia tidak bersumber dari ERD — tabel PENDARATAN tidak digambar di
	// sana. Penyebut 28 tetap; "bukan dari §2" naik dari 13 menjadi 16.
	//
	// 6 Oktober 2026, keempat kalinya: 43 -> 44. Yang ke-44
	// `FK_TT_LIMIT_MEASURE` dari migrasi `439`, sebab yang sama persis.
	// "bukan dari §2" naik dari 16 menjadi 17.
	n := strings.Count(strings.ToUpper(tanpaKomentar(gabungan(t))), "ON DELETE CASCADE")
	if n != 44 {
		t.Errorf("ON DELETE CASCADE ditemukan %d, mau 44 (ERD.md §2: 28 relasi ikut hapus, "+
			"1 di antaranya tabelnya belum dibuat; ditambah 4 dari ERD HTML baris 4, 9, 10, 27, "+
			"ditambah 1 tabel pendaratan migrasi 430 yang tidak ada di ERD mana pun)", n)
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
// rapatkanSpasi mengubah tiap deret spasi putih menjadi satu spasi tunggal,
// sehingga pemilih berbentuk teks tidak bergantung pada tata letak berkasnya.
func rapatkanSpasi(s string) string { return strings.Join(strings.Fields(s), " ") }

func TestTabelTanpaKunciAlamiTidakDiberiDiamDiam(t *testing.T) {
	// ⚠️ Spasi DIRAPATKAN lebih dulu. Bentuk harfiah "CONSTRAINT UQ_X UNIQUE"
	// dapat dihindari - tanpa sengaja maupun tidak - hanya dengan memenggal
	// barisnya di antara nama constraint dan kata UNIQUE, dan penjaga yang
	// dapat dihindari dengan satu baris baru tidak menjaga apa pun.
	// (2 Oktober 2026: itu sungguh terjadi, pada percobaan memasang
	// UQ_PEMULIHAN_LIMIT untuk tiket 32.)
	sql := rapatkanSpasi(gabungan(t))
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

// INV-61 - arsip muatan keluar TIDAK punya jalur baca, dan angkanya dicetak.
//
// Tiket 42 menuntut ujinya berbentuk SAPUAN: "sapu seluruh basis kode untuk
// kueri yang menyentuh isi arsip; hasilnya HARUS nol, dan angkanya dicetak."
// Angka yang dicetak itulah yang membuat uji ini berguna setahun lagi: nol
// yang tidak terlihat tidak dapat dibedakan dari sapuan yang rusak.
//
// ⛔ Yang dilarang membaca KOLOM MUATAN, bukan menyentuh tabelnya. Menghitung
// berapa arsip yang ada tidak membuat arsip menjadi sumber kedua; membaca
// isinya membuatnya begitu (ADR-0034).
//
// ⚠️ Komentar DIBUANG lebih dulu, dan berkas ini sendiri dikecualikan. Dua
// sebab, keduanya ditemukan saat uji ini pertama dijalankan: berkas migrasi
// dan berkas repository MENJELASKAN keputusannya dengan menyebut kata
// `MUATAN` dan `SELECT` di dalam prosa, dan uji ini sendiri harus menyebut
// keduanya untuk dapat melarangnya. Sapuan yang memindai teks mentah
// menemukan kata-katanya sendiri - pelajaran yang sama yang sudah membuat
// `tanpaKomentar` lahir di berkas ini.
func TestArsipTidakPunyaJalurBaca(t *testing.T) {
	var pembaca []string
	disapu := 0
	err := filepath.Walk(".", func(jalur string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(jalur, ".go") {
			return nil
		}
		// Berkas penjaganya sendiri - lihat alasan di kepala fungsi.
		if filepath.Base(jalur) == "migrasi_invarian_test.go" {
			return nil
		}
		isi, errBaca := os.ReadFile(jalur)
		if errBaca != nil {
			return errBaca
		}
		disapu++
		// ⚠️ DIPERTAJAM 3 Oktober 2026. Bentuk sebelumnya memotong 600 aksara
		// sesudah tiap `SELECT` dan mencari `MUATAN` di dalamnya. Jendela itu
		// MENYEBERANG ke literal berikutnya: `bukti_db_test.go` punya
		// `SELECT COUNT(*) FROM MIGRASI_KORELASI` yang 600 aksara sesudahnya
		// memuat konstanta `insPendaratan` ber-kolom `MUATAN` — dan penjaga
		// berbunyi untuk berkas yang nol membaca arsip.
		//
		// Yang diperiksa sekarang LITERAL TEKS Go satu per satu: sebuah
		// pernyataan SQL hidup di dalam SATU literal, jadi `SELECT` dan
		// `MUATAN` yang berada di literal BERBEDA memang bukan satu kueri.
		// Lebih tajam, bukan lebih longgar: kueri yang sungguh memilih
		// `MUATAN` tetap tertangkap, di literal mana pun ia ditulis.
		for _, lit := range literalTeksGo(tanpaKomentarGo(string(isi))) {
			u := strings.ToUpper(lit)
			if strings.Contains(u, "SELECT") && strings.Contains(u, "MUATAN") {
				pembaca = append(pembaca, jalur)
				break
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if disapu == 0 {
		t.Fatal("nol berkas Go tersapu; pembacanya yang rusak, bukan kodenya")
	}
	t.Logf("INV-61: %d berkas Go disapu, %d memuat SELECT atas kolom MUATAN", disapu, len(pembaca))
	if len(pembaca) != 0 {
		t.Errorf("ADR-0034 dan INV-61: arsip TIDAK punya jalur baca, tetapi %d berkas "+
			"memuat SELECT atas kolom MUATAN: %v. Menambahkannya membalikkan ADR-0034 "+
			"tanpa membukanya", len(pembaca), pembaca)
	}
}

// literalTeksGo mengembalikan isi tiap literal teks Go - backtick maupun
// tanda kutip ganda.
//
// Sengaja sederhana dan sengaja MELEBIH: ia tidak mengurai escape, jadi
// sebuah literal dapat terbaca lebih panjang daripada yang sebenarnya.
// Melebih di sini aman - ia membuat penjaga menangkap LEBIH banyak, bukan
// lebih sedikit.
func literalTeksGo(src string) []string {
	var keluar []string
	for i := 0; i < len(src); i++ {
		switch src[i] {
		case '`':
			if j := strings.IndexByte(src[i+1:], '`'); j >= 0 {
				keluar = append(keluar, src[i+1:i+1+j])
				i += j + 1
			}
		case '"':
			j := i + 1
			for j < len(src) && src[j] != '"' && src[j] != '\n' {
				if src[j] == '\\' {
					j++
				}
				j++
			}
			if j < len(src) && src[j] == '"' {
				keluar = append(keluar, src[i+1:j])
			}
			i = j
		}
	}
	return keluar
}

// tanpaKomentarGo membuang komentar `//` dan `/* */` dari sumber Go.
//
// Ia sengaja sederhana: tanda `//` di dalam literal teks (mis. sebuah URL)
// akan ikut terpotong. Itu diterima di sini sebab modul ini nol URL di dalam
// literal, dan sapuan yang terlalu rajin membuang BUKTI - bukan menambahkannya.
func tanpaKomentarGo(src string) string {
	var b strings.Builder
	for len(src) > 0 {
		i := strings.Index(src, "//")
		j := strings.Index(src, "/*")
		switch {
		case i < 0 && j < 0:
			b.WriteString(src)
			return b.String()
		case j < 0 || (i >= 0 && i < j):
			b.WriteString(src[:i])
			k := strings.IndexByte(src[i:], '\n')
			if k < 0 {
				return b.String()
			}
			src = src[i+k:]
		default:
			b.WriteString(src[:j])
			k := strings.Index(src[j:], "*/")
			if k < 0 {
				return b.String()
			}
			src = src[j+k+2:]
		}
	}
	return b.String()
}

// ⛔ `TREATY_IN` adalah tabel WARISAN, dan modul ini BUKAN pemiliknya.
//
// Ia memuat 1.854 baris produksi-bayangan. Layar daftar membacanya
// (keputusan pemilik proses 3 Oktober 2026), dan membaca itu seluruh izinnya.
//
// Dua hal dijaga di sini, dan keduanya pernah menjadi cara modul merusak
// tabel yang bukan miliknya:
//
//  1. NOL `INSERT`/`UPDATE`/`DELETE`/`MERGE` terhadapnya di berkas Go mana
//     pun milik modul ini;
//  2. NOL penyebutan di migrasi mana pun — menuliskannya di migrasi berarti
//     mengklaim kepemilikan, dan `migrate` berikutnya akan mencoba
//     membuatnya di atas tabel yang sudah berisi.
func TestWarisanHanyaDibaca(t *testing.T) {
	// ⚠️ DIPERLUAS 3 Oktober 2026 ke `M_TREATY_IN` — pasangan `TREATY_IN`
	// yang memegang dokumen aslinya, 1.854 baris `CLOB`. Form kontrak
	// membacanya, dan membaca itu seluruh izinnya.
	// ⚠️ DIPERLUAS lagi 3 Oktober 2026 ke tiga tabel warisan yang ronde tab
	// mulai dibaca: `M_TREATY_IN2` (7.281 baris, empat tab),
	// `TREATYEXCHANGEYEARLY` (140), dan `M_TREATY_IN_DETAIL` (27.617).
	// Ketiganya DIBACA, dan membaca itu seluruh izinnya.
	//
	// ⛔ `M_TREATY_IN2` TIDAK tertangkap oleh pola `M_TREATY_IN`: `` sesudah
	// `IN` tidak cocok di depan angka `2`. Itu sebabnya ia disebut sendiri,
	// bukan diandaikan ikut terjaga.
	// ⭐ DIPERLUAS 6 Oktober 2026 ke dua tabel warisan yang sudah dibaca modul
	// ini tetapi belum dijaga di sini: `M_ATTACHMENTTREATY_2` (panel
	// Attachment) dan `TREATYINPRODUCTION` (panel Existing Policy).
	//
	// ⚠️ Sebabnya ditemukan saat menyapu gerbang berangka PERSIS: keduanya
	// dijaga HANYA oleh cacah barisnya di uji `db`, dan cacah itu kini
	// dilonggarkan menjadi LANTAI sebab Pega sendiri menambah barisnya.
	// Melonggarkannya tanpa memasang penjaga tulis yang sesungguhnya akan
	// meninggalkan keduanya tanpa penjaga sama sekali.
	for _, tabel := range []string{
		"TREATY_IN",
		"TREATYEXCHANGEYEARLY", "M_TREATY_IN_DETAIL",
		// ⭐ `M_ATTACHMENTTREATY_2` DIKELUARKAN 8 Oktober 2026 — keputusan
		// pemakai: *"untuk upload file seharusnya kesini SELECT * FROM
		// M_ATTACHMENTTREATY_2"*. Tombol `Upload file` menulisnya
		// (`repository/lampiran_tulis.go`, `InsertAttachment2_Sql`).
		"TREATYINPRODUCTION",
		// ⭐ 6 Oktober 2026 — isi dropdown `Treaty Type` tab Limits.
		"REINSURANCETYPE", "TREATYGROUP", "CURRENCY", "TREATYBUSINESS", "ACHIEVEMENT",
		// ⭐ 7 Oktober 2026 — tab Share Non-Prop: susunan treaty master
		// (dropdown Spreading Type dan anaknya) dan autocomplete Reinsurer.
		"PROPORTIONALARRG", "TREATYYEAR", "AGENT",
		// ⭐ 7 Oktober 2026 — cadangan skalar akar Share (`SaveTreatyInDetail_Act`).
		"TREATYINDETAIL",
	} {
		t.Run(tabel, func(t *testing.T) { warisanHanyaDibaca(t, tabel) })
	}
}

// namaKolomJuga - tabel warisan yang namanya juga dipakai sebagai kolom.
var namaKolomJuga = map[string]bool{"TREATYGROUP": true, "CURRENCY": true, "TREATYYEAR": true}

func warisanHanyaDibaca(t *testing.T, tabel string) {
	t.Helper()
	// ⚠️ Kata yang dicari harus berdiri sebagai KATA. `VERSI_KONTRAK` dan
	// `TREATY_IN_ID` memuat potongan yang sama, dan sapuan yang menangkapnya
	// akan merah selamanya tanpa satu pun pelanggaran.
	pola := regexp.MustCompile(`(?i)\b(INSERT\s+INTO|UPDATE|DELETE\s+FROM|MERGE\s+INTO)\s+[^\s;]*\b` + tabel + `\b`)

	disapu := 0
	var pelanggar []string
	err := filepath.Walk(".", func(jalur string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(jalur, ".go") {
			return nil
		}
		if filepath.Base(jalur) == "migrasi_invarian_test.go" {
			return nil // berkas penjaganya sendiri; ia menyebut kata-katanya
		}
		isi, errBaca := os.ReadFile(jalur)
		if errBaca != nil {
			return errBaca
		}
		disapu++
		for _, lit := range literalTeksGo(tanpaKomentarGo(string(isi))) {
			if pola.MatchString(lit) {
				pelanggar = append(pelanggar, jalur)
				break
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if disapu == 0 {
		t.Fatal("nol berkas Go tersapu; pembacanya yang rusak")
	}
	t.Logf("%s: %d berkas Go disapu, %d menulis ke tabel warisan", tabel, disapu, len(pelanggar))
	if len(pelanggar) != 0 {
		t.Errorf("%s adalah tabel WARISAN dan modul ini bukan pemiliknya; "+
			"%d berkas menulis ke sana: %v", tabel, len(pelanggar), pelanggar)
	}

	// ⚠️ Aturannya MENGGIGIT, dan itu diperiksa di sini — sapuan yang tidak
	// pernah menemukan apa pun tidak dapat dibedakan dari sapuan yang rusak.
	for _, jahat := range []string{
		"INSERT INTO POOLDATA." + tabel + " (ID) VALUES ('x')",
		"delete from {skema}." + tabel + " where ID = 1",
		"MERGE INTO " + tabel + " t USING dual ON (1=1)",
	} {
		if !pola.MatchString(jahat) {
			t.Errorf("pola tidak menangkap penulisan yang nyata: %q", jahat)
		}
	}
	// Dan ia tidak menangkap yang SAH: pembacaan, dan tabel bernama mirip.
	for _, sah := range []string{
		"SELECT ID FROM POOLDATA." + tabel + " ORDER BY ID",
		"INSERT INTO {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK) VALUES (1)",
		"INSERT INTO {skema}.MIGRASI_KORELASI (ID_MIGRASI_KORELASI) VALUES (1)",
	} {
		if pola.MatchString(sah) {
			t.Errorf("pola menangkap pernyataan yang SAH: %q", sah)
		}
	}

	// Dan ia tidak boleh disebut di migrasi mana pun.
	//
	// ⚠️ `TREATYGROUP` dan `CURRENCY` juga NAMA KOLOM sah di migrasi modul
	// ini; untuk keduanya yang dicari hanya penyebutan sebagai TABEL.
	sebut := regexp.MustCompile(`\b` + tabel + `\b`)
	if namaKolomJuga[tabel] {
		sebut = regexp.MustCompile(`(?i)\b(TABLE|REFERENCES|INTO|UPDATE|FROM|JOIN)\s+[^\s;(]*\b` + tabel + `\b`)
	}
	for nama, isi := range seluruhMigrasi(t) {
		if sebut.MatchString(tanpaKomentar(isi)) {
			t.Errorf("%s menyebut %s; tabel warisan TIDAK dimiliki modul ini, dan "+
				"menuliskannya di migrasi berarti `migrate` berikutnya mencoba membuatnya", nama, tabel)
		}
	}
}

// dicabut menjawab apakah sebuah kunci asing dibuang migrasi berikutnya —
// langsung lewat `DROP CONSTRAINT`, atau ikut terbawa `DROP TABLE` atas
// tabel yang memuatnya.
//
// ⛔ YANG TERAKHIR MENANG, dan itu bukan hiasan: migrasi 420 MEMBONGKAR
// lalu MEMASANG ULANG dua puluh satu kunci asing. Membaca "ada DROP" saja
// membuang kedua puluh satunya — FK_EGNPI_1, FK_DOKUMEN_KONTRAK_1, dan
// seterusnya — padahal ketiganya hidup. Yang menentukan POSISI: bila DROP
// terakhir berada SESUDAH definisi terakhir, barulah ia tercabut.
//
// ⛔ SENGAJA TANPA REGEXP. Versi pertama memakai pola ber-backslash dan
// gagal diam-diam: satu aksara kendali ikut tersalin ke dalam literalnya,
// polanya tidak pernah cocok, dan uji tetap hijau sambil menilai kunci
// asing yang sudah tidak ada.
func dicabut(sql, fk string) bool {
	atas := strings.ToUpper(sql)
	F := strings.ToUpper(fk)
	buat := strings.LastIndex(atas, "CONSTRAINT "+F+" FOREIGN KEY")
	if buat < 0 {
		return false
	}
	if buang := strings.LastIndex(atas, "DROP CONSTRAINT "+F); buang > buat {
		return true
	}
	// Tabel pemiliknya: CREATE/ALTER TABLE terakhir sebelum definisi itu.
	awal := atas[:buat]
	potong := strings.LastIndex(awal, "CREATE TABLE {SKEMA}.")
	if j := strings.LastIndex(awal, "ALTER TABLE {SKEMA}."); j > potong {
		potong = j
	}
	if potong < 0 {
		return false
	}
	sisa := awal[potong+strings.Index(awal[potong:], "{SKEMA}.")+len("{SKEMA}."):]
	n := 0
	for n < len(sisa) && (sisa[n] == '_' || (sisa[n] >= 'A' && sisa[n] <= 'Z') || (sisa[n] >= '0' && sisa[n] <= '9')) {
		n++
	}
	tabel := sisa[:n]
	if tabel == "" {
		return false
	}
	// Tabelnya dibuang SESUDAH dibuat — bukan sekadar pernah disebut DROP.
	buatTbl := strings.LastIndex(atas, "CREATE TABLE {SKEMA}."+tabel)
	buangTbl := strings.LastIndex(atas, "DROP TABLE {SKEMA}."+tabel)
	return buangTbl > buatTbl
}
