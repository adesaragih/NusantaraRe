package repository

// Test bentuk migrasi - TANPA Oracle.
//
// Untuk apa berkas ini: memeriksa isi berkas .sql yang ditanam ke biner. Ia
// tidak menyambung ke basis data sama sekali; yang diuji adalah apa yang akan
// dikirim ke Oracle. Karena itu ia tetap berjalan di mesin tanpa instance.
//
// Banyak acceptance criteria tiket 14 berbentuk "test yang menemukan X gagal".
// Justru itu yang dikerjakan di sini.

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
)

func seluruhSQL(t *testing.T, mundur bool) map[string]string {
	t.Helper()
	langkah, err := migrasi.Daftar(mundur, berkasMigrasi)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, m := range langkah {
		out[m.Nama] = strings.Join(m.Pernyataan, "\n")
	}
	return out
}

func gabungSemua(t *testing.T) string {
	t.Helper()
	var b strings.Builder
	for _, isi := range seluruhSQL(t, false) {
		b.WriteString(isi)
		b.WriteString("\n")
	}
	return strings.ToUpper(b.String())
}

// Tiket 14: enam tabel klaim + T_WORK_CLAIM harus dibuat.
func TestTujuhTabelDibuat(t *testing.T) {
	sql := gabungSemua(t)
	mau := []string{
		"T_WORK_CLAIM",
		"T_GENERAL_CLAIM",
		"T_CLAIMLF_PREMIUMLIST_DETAIL",
		"T_CLAIMLF_ADJUSTMENT",
		"T_CLAIMLF_ADJUSTMENT_SPREADING",
		"T_CLAIMLF_ADJ_SPREADING_RETRO",
		"T_CLAIMLF_DOCUMENT",
	}
	for _, tb := range mau {
		if !strings.Contains(sql, "CREATE TABLE {SKEMA}."+tb+" (") {
			t.Errorf("tabel %s tidak dibuat", tb)
		}
	}
}

// Header klaim tidak lagi memuat kolom yang pindah ke T_WORK_CLAIM,
// dan PL_NUMBER sudah berganti nama menjadi POLICY_NO.
func TestHeaderKlaimTidakMemuatKolomYangPindah(t *testing.T) {
	isi := ""
	for nama, teks := range seluruhSQL(t, false) {
		if strings.Contains(nama, "t_general_claim") {
			isi = strings.ToUpper(teks)
		}
	}
	if isi == "" {
		t.Fatal("berkas T_GENERAL_CLAIM tidak ketemu")
	}
	for _, kol := range []string{"CASEID ", "CREATE_OP", "CREATE_OP_NAME", "TGL_UPDATE", "PL_NUMBER"} {
		if strings.Contains(isi, kol) {
			t.Errorf("kolom %q masih ada di T_GENERAL_CLAIM - ia pindah atau berganti nama", kol)
		}
	}
	for _, kol := range []string{"CASEID_POLICY", "POLICY_NO", "ENDORSMENT_NO"} {
		if !strings.Contains(isi, kol) {
			t.Errorf("penunjuk polis %q tidak ada di T_GENERAL_CLAIM", kol)
		}
	}
}

// Relasi 2: shared primary key - T_GENERAL_CLAIM.ID sekaligus PK dan FK.
func TestSharedPrimaryKey(t *testing.T) {
	var isi string
	for nama, teks := range seluruhSQL(t, false) {
		if strings.Contains(nama, "t_general_claim") {
			isi = strings.ToUpper(teks)
		}
	}
	if !strings.Contains(isi, "PRIMARY KEY (ID)") {
		t.Error("T_GENERAL_CLAIM.ID bukan kunci utama")
	}
	if !strings.Contains(isi, "FOREIGN KEY (ID)") {
		t.Error("T_GENERAL_CLAIM.ID bukan kunci tamu ke T_WORK_CLAIM - shared PK tidak terbentuk")
	}
}

// Relasi 4: adjustment menggantung pada PESERTA, bukan pada header klaim.
func TestAdjustmentMenggantungPadaPeserta(t *testing.T) {
	var isi string
	for nama, teks := range seluruhSQL(t, false) {
		if strings.HasPrefix(nama, "004_") {
			isi = strings.ToUpper(teks)
		}
	}
	if !strings.Contains(isi, "REFERENCES {SKEMA}.T_CLAIMLF_PREMIUMLIST_DETAIL (ID)") {
		t.Error("FK adjustment tidak menunjuk T_CLAIMLF_PREMIUMLIST_DETAIL")
	}
	if strings.Contains(isi, "REFERENCES {SKEMA}.T_GENERAL_CLAIM") {
		t.Error("adjustment menggantung pada header klaim - seharusnya pada peserta")
	}
}

// milikClaimLife menjawab apakah berkas migrasi itu milik modul Claim Life.
//
// ⛔ Batasnya NOMOR, dan itu keputusan yang tercatat: Claim Life memakai
// 001-049, PremiumList Life mulai 050. Memisahkan lewat nama tabel akan
// gagal pada tabel yang namanya tidak menyebut modulnya.
func milikClaimLife(nama string) bool {
	return nama < "050_"
}

// Setiap kunci tamu ber-index, dan KOMITE_ID ber-index UNIK.
func TestKunciTamuBerIndex(t *testing.T) {
	// Pernyataan CREATE INDEX boleh ditulis pada beberapa baris, jadi spasinya
	// diratakan dulu. Tanpa ini, test-nya sendiri yang keliru melaporkan index
	// yang sebenarnya ada.
	sql := regexp.MustCompile(`\s+`).ReplaceAllString(gabungSemua(t), " ")
	for _, kol := range []string{
		"COVER_KEY", "CLAIM_ID", "PREMIUM_LIST_DETAIL_ID", "ADJUSTMENT_ID", "SPREADING_ID",
	} {
		if !regexp.MustCompile(`CREATE INDEX [^;]*?\(` + kol + `\)`).MatchString(sql) {
			t.Errorf("kunci tamu %s tidak ber-index", kol)
		}
	}
	if !regexp.MustCompile(`CREATE UNIQUE INDEX [^;]*?\(KOMITE_ID\)`).MatchString(sql) {
		t.Error("KOMITE_ID tidak ber-index UNIK")
	}
}

// ADR-U-0006: identitas seluruh tabel T_CLAIMLF_* dari sequence.
func TestSequenceUntukTabelYangMemakainya(t *testing.T) {
	sql := gabungSemua(t)
	for _, s := range []string{
		"SEQ_CLAIMLF_PLD", "SEQ_CLAIMLF_ADJ", "SEQ_CLAIMLF_SPR",
		"SEQ_CLAIMLF_SPR_RETRO", "SEQ_T_CLAIMLF_DOCUMENT",
	} {
		if !strings.Contains(sql, "CREATE SEQUENCE {SKEMA}."+s) {
			t.Errorf("sequence %s tidak dibuat", s)
		}
	}
}

// Identitas T_WORK_CLAIM adalah teks berformat, bukan angka sequence.
func TestIdentitasWorkClaimBerupaTeks(t *testing.T) {
	var isi string
	for nama, teks := range seluruhSQL(t, false) {
		if strings.HasPrefix(nama, "001_") {
			isi = strings.ToUpper(teks)
		}
	}
	if !regexp.MustCompile(`ID\s+VARCHAR2\(\d+\) NOT NULL`).MatchString(isi) {
		t.Error("T_WORK_CLAIM.ID bukan teks - tiket 14 menetapkan teks berformat CLM-/KMT-")
	}
	if regexp.MustCompile(`ID\s+NUMBER`).MatchString(isi) {
		t.Error("T_WORK_CLAIM.ID bertipe numerik")
	}
}
<<<<<<< HEAD

// Berkas migrasi tidak boleh berawalan byte order mark.
//
// Ini bukan kerewelan gaya. BOM membuat baris komentar pertama lolos menjadi
// bagian pernyataan SQL, dan Oracle menolaknya - sementara seluruh test tanpa
// basis data tetap hijau. Sekali terjadi di berkas 004; sejak itu dijaga.
func TestBerkasMigrasiTanpaBOM(t *testing.T) {
	entri, err := berkasMigrasi.ReadDir("migrations")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entri {
		isi, err := berkasMigrasi.ReadFile("migrations/" + e.Name())
		if err != nil {
			t.Fatal(err)
		}
		if len(isi) >= 3 && isi[0] == 0xEF && isi[1] == 0xBB && isi[2] == 0xBF {
			t.Errorf("%s berawalan BOM UTF-8", e.Name())
		}
	}
}

// Berkas migrasi tidak boleh memuat byte carriage return.
//
// Sepupu dekat jebakan BOM di atas. Alat Windows - PowerShell, penyunting yang
// disetel salah, git tanpa .gitattributes - menulis akhiran baris CRLF. CR yang
// terbawa masuk ke teks pernyataan yang dikirim ke Oracle, dan sekali lagi
// seluruh test tanpa basis data tetap hijau sementara instance menolaknya.
// Berkas .gitattributes di akar repositori menjaga sisi git; test ini menjaga
// sisi berkas.
func TestBerkasMigrasiTanpaCR(t *testing.T) {
	entri, err := berkasMigrasi.ReadDir("migrations")
	if err != nil {
		t.Fatal(err)
	}
	if len(entri) == 0 {
		t.Fatal("nol berkas migrasi terbaca; pembacanya yang rusak")
	}
	for _, e := range entri {
		isi, err := berkasMigrasi.ReadFile("migrations/" + e.Name())
		if err != nil {
			t.Fatal(err)
		}
		if n := strings.Count(string(isi), "\r"); n > 0 {
			t.Errorf("%s memuat %d byte CR; akhiran barisnya harus LF", e.Name(), n)
		}
	}
}

// Pernyataan pertama setiap berkas harus benar-benar mulai dengan kata perintah
// SQL - bukan dengan sisa komentar.
func TestPernyataanMulaiDenganPerintah(t *testing.T) {
	for _, mundur := range []bool{false, true} {
		langkah, err := migrasi.Daftar(mundur, berkasMigrasi)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range langkah {
			for i, p := range m.Pernyataan {
				kata := strings.ToUpper(strings.Fields(p)[0])
				switch kata {
				case "CREATE", "ALTER", "DROP", "INSERT", "UPDATE", "DELETE", "SELECT":
				default:
					t.Errorf("%s pernyataan %d mulai dengan %q, bukan perintah SQL",
						m.Nama, i, kata)
				}
			}
		}
	}
}

// Setiap pernyataan CREATE di berkas migrasi harus dapat dibaca namanya.
//
// Kalau ada satu saja yang tidak terbaca, jalur "dilewati lalu dibuktikan"
// menolak dengan galat - dan lebih baik test ini yang menemukannya lebih dulu,
// di mesin tanpa Oracle.
//
// ⚠️ Versi pertama test ini memecah ulang teks yang sudah disambung seluruhSQL,
// sehingga pemisahnya tidak pernah memisah apa pun dan yang diperiksa hanya
// SATU pernyataan per berkas - 8 dari 19. Penjaga yang lebih lemah dari
// namanya. Sekarang pernyataannya diambil dari daftarMigrasi apa adanya.
func TestSeluruhCreateDapatDibacaNamanya(t *testing.T) {
	langkah, err := migrasi.Daftar(false, berkasMigrasi)
	if err != nil {
		t.Fatal(err)
	}
	diperiksa := 0
	for _, m := range langkah {
		for _, p := range m.Pernyataan {
			if !migrasi.PernyataanBuat(p) {
				continue
			}
			diperiksa++
			if migrasi.NamaObjekDibuat(p) == "" {
				t.Errorf("%s: nama objek tidak terbaca dari %q", m.Nama, migrasi.RingkasPernyataan(p))
			}
		}
	}
	// Angkanya dikunci: kalau pemisah pernyataan rusak lagi, cacahnya anjlok
	// dan test ini gagal alih-alih diam-diam memeriksa lebih sedikit.
	// Dua puluh sejak 26-09-2026: langkah 009 menambah SEQ_WORK_CLAIM
	// (butir aa, tiket 02).
	//
	// ⛔ Angkanya DIPERBARUI, bukan dilonggarkan - A1 27-09-2026 menambah
	// butir am: 1 tabel (T_CLAIMLF_JEJAK) + 1 sequence (SEQ_CLAIMLF_JEJAK) +
	// 2 index = 4 pernyataan CREATE baru.
	//
	// ⛔ Diperbarui LAGI - A1 menambah butir af (2 tabel + 1 sequence +
	// 2 index) dan temuan audit A0 (BUSINESS_CODE, ALTER - tidak dihitung).
	//
	// ⛔ Diperbarui LAGI - A2 menambah butir aq: 1 tabel (T_LOG_SERVICE_RNM) +
	// 1 sequence + 2 index = 4 pernyataan CREATE baru.
	//
	// ⛔ Diperbarui LAGI - tiket 00 PremiumList Life menambah TUJUH tabel
	// (T_WORK_POLIS, T_PREMIUM_LIST, _DETAIL, _SPREADING, _SPREADING_RETRO,
	// _SUMMARY, T_VIEW_SUGGEST) dan ENAM index FK (dua pada _DETAIL, satu
	// pada masing-masing tabel anak lainnya). Nol sequence: pengenalnya
	// dirakit di repository, pola PengenalWorkBerikut (butir pl3).
	//
	// ⛔ Diperbarui LAGI - butir bd menambah T_CLAIMLF_DIAGNOSE:
	// 1 tabel + 1 sequence (SEQ_CLAIMLF_DIAGNOSE) + 1 index (FK peserta).
	// ALTER pada DISEASE tidak dihitung - ia bukan CREATE.
	//
	// 11+7+1 = 19 tabel + 10 sequence + 13+6+1 = 20 index = 49.
	// +1 tabel dari 019 (butir be, kartu berkas unggahan) = 50.
	//
<<<<<<< HEAD
	// ⛔ Diperbarui LAGI - butir bn (GILIRAN-13): +1 sequence SEQ_WORK_POLIS
	// (057). RALAT atas catatan tiket 00 di atas: "nol sequence" untuk polis
	// keliru - pl3 memutuskan SEQ_WORK_POLIS, dan ia terlewat. ALTER kolom
	// FLAG_ONGOING_POLICY tidak dihitung.
	const mau = 51
=======
	// ⛔ Diperbarui LAGI - tiket 01 Treaty Contract Out (300-306) menambah
	// TUJUH tabel (T_TREATYYEAR, T_TREATYCONTRACT, T_TREATYREINSURER,
	// T_MTREATYSECURITY, T_TREATYBUSINESS, T_PROPORTIONALARRG,
	// T_TREATYCO_JEJAK) + 7 sequence + 7 index = 21 pernyataan CREATE = 71.
<<<<<<< HEAD
	const mau = 71
>>>>>>> 1872d26 (treaty-contract-out: tiket 01 — skema relasional + migrasi + tipe dirapikan)
=======
	// +3 dari 307 (tiket 12 Treaty Contract Out, 29-09-2026): tabel
	// T_TREATYYEAR_LAMPIRAN + index + sequence = 74.
<<<<<<< HEAD
	const mau = 74
>>>>>>> 7b1db9b (treaty-contract-out: tiket 12 — lampiran di tahun treaty)
=======
	// Penyatuan 29-09-2026: 50 (dasar) + 1 (057, butir bn) + 24 (300-307 Treaty Contract Out) = 75.
	// ⛔ tco4 (keputusan work owner 29-09-2026): Treaty Contract Out NOL tabel
	// baru - 300-307 dibuang, kembali ke 50 + 1 = 51.
<<<<<<< HEAD
	const mau = 51
>>>>>>> e6905ed (treaty-contract-out: tco4 — nol tabel baru, migrasi 300-307 dibuang)
=======
	// ⛔ OQ-PL-15 (GILIRAN-15): 058 membuat ULANG SEQ_WORK_POLIS - DROP (tidak
	// dihitung) lalu CREATE SEQUENCE ... START WITH 22374 (+1) = 52.
	const mau = 52
>>>>>>> c31eb12 (premiumlist-life: PL-15 — SEQ_WORK_POLIS mulai 22374)
	if diperiksa != mau {
		t.Errorf("pernyataan CREATE diperiksa %d, mau %d", diperiksa, mau)
	}
}

// KolomCreateTable membaca nama dan kolom dari setiap CREATE TABLE migrasi.
//
// Cacahnya dikunci: delapan CREATE TABLE. Pernyataan yang BUKAN CREATE TABLE -
// CREATE INDEX dan CREATE SEQUENCE - harus mengembalikan nama kosong, kalau
// tidak pra-terbang akan mencari "bentuk" sebuah sequence.
func TestKolomCreateTableMembacaSeluruhTabel(t *testing.T) {
	langkah, err := migrasi.Daftar(false, berkasMigrasi)
	if err != nil {
		t.Fatal(err)
	}
	tabel, bukanTabel := 0, 0
	for _, m := range langkah {
		for _, p := range m.Pernyataan {
			nama, kolom := migrasi.KolomCreateTable(p)
			if nama == "" {
				bukanTabel++
				continue
			}
			tabel++
			if len(kolom) == 0 {
				t.Errorf("%s: tabel %s terbaca tanpa satu pun kolom", m.Nama, nama)
			}
			for _, k := range kolom {
				if strings.HasPrefix(k, "CONSTRAINT") || strings.HasPrefix(k, "REFERENCES") {
					t.Errorf("%s: %s menganggap %q sebagai kolom", m.Nama, nama, k)
				}
			}
		}
	}
	// Tujuh, bukan delapan: T_MIGRASI dibuat siapkanTabelMigrasi, di luar
	// berkas migrasi. Sesudah migrasi, katalog memang memuat delapan tabel.
	// 11 tabel Claim Life + 7 tabel PremiumList Life (tiket 00) = 18.
	// 20 tabel + 10 sequence + 20 index = 50 pernyataan CREATE, cocok dengan
	// cacah yang dikunci TestSeluruhCreateDapatDibacaNamanya.
	// +1 tabel, +1 sequence, +1 index dari 018 (butir bd, diagnosa).
	// +1 tabel dari 019 (butir be) - TANPA sequence dan TANPA index:
	// identitasnya cap waktu `models.IDDokumenBaru`, bukan nomor kita,
	// dan PK-nya sudah berindeks sendiri.
	// +7 tabel dari 300-306 (tiket 01 Treaty Contract Out) = 27.
	// +1 tabel dari 307 (tiket 12 Treaty Contract Out, lampiran) = 28.
	// ⛔ tco4 (29-09-2026): 300-307 dibuang - kembali ke 20.
	const mauTabel = 20
	if tabel != mauTabel {
		t.Errorf("CREATE TABLE terbaca %d, mau %d", tabel, mauTabel)
	}
	if bukanTabel == 0 {
		t.Error("nol pernyataan bukan-tabel; CREATE INDEX dan SEQUENCE seharusnya ada")
	}
}
=======
>>>>>>> 9b478cc (refactor(bentuk-B) paket 3: modul/premiumlist — polis_*, migrasi 050-058, kontrak PembacaPolis)
