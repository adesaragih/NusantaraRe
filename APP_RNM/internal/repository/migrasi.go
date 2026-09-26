package repository

// Pelari migrasi basis data.
//
// Untuk apa berkas ini: menjalankan berkas SQL di folder `migrations/` secara
// berurutan, sekali saja masing-masing, dan menyediakan jalur mundur.
//
// Dibaca sesudah: repository.go (yang memperkenalkan DB dan Qualify).
//
// Istilah yang dipakai di sini, sekali dijelaskan:
//   - migrasi    : satu langkah perubahan bentuk basis data, disimpan sebagai
//                  satu berkas .sql bernomor. Nomornya menentukan urutan.
//   - idempoten  : aman dijalankan berulang kali. Menjalankan dua kali memberi
//                  hasil yang sama dengan menjalankan sekali.
//   - embed      : menanam isi berkas ke dalam biner Go saat dibangun, sehingga
//                  program tidak perlu mencari berkas .sql di disk saat jalan.
//   - jalur mundur: berkas berpasangan (_down.sql) yang membatalkan langkahnya.
//
// Aturan yang dijaga:
//   ADR-U-0029  nol COMMIT di teks SQL. Perintah DDL Oracle memang menutup
//               transaksinya sendiri, dan itu sifat Oracle - bukan alasan untuk
//               menulis COMMIT di teks SQL. PeriksaSQL menolaknya.
//   ADR-U-0033  nama skema disebut eksplisit; berkas .sql memakai penanda
//               {skema} yang diganti saat dijalankan.
//   ADR-U-0005  menolak berjalan bila lingkungan menunjuk produksi Pega.

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
)

//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// ErrMigrasiDiProduksi menolak perubahan bentuk basis data di lingkungan
// produksi Pega. Migrasi membuat dan membuang tabel; ia tidak pernah boleh
// berjalan di sana tanpa keputusan manusia.
var ErrMigrasiDiProduksi = errors.New("repository: menolak migrasi saat IS_PEGA_PROD=true")

// namaTabelMigrasi mencatat langkah mana yang sudah dijalankan.
//
// [usulan] Nama dan bentuk tabel pencatat ini belum pernah diputuskan lewat
// ADR; ia boleh diganti lewat keputusan tertulis.
const namaTabelMigrasi = "T_MIGRASI"

// Migrasi adalah satu langkah, sudah dipecah menjadi pernyataan-pernyataan.
type Migrasi struct {
	Nama       string
	Pernyataan []string
}

// LaporanMigrasi menceritakan apa yang benar-benar terjadi.
type LaporanMigrasi struct {
	Dijalankan []string // langkah yang baru dijalankan kali ini
	Dilewati   []string // langkah yang memang sudah pernah dijalankan
	// Pernyataan mencacah yang benar-benar DIEKSEKUSI. Pernyataan yang
	// dilewati karena objeknya sudah ada tidak ikut dihitung; ia muncul
	// di ObjekSudahAda, dan memang tidak dikirim ke Oracle.
	Pernyataan int
	// ObjekSudahAda mencatat pernyataan CREATE yang dilewati karena objeknya
	// ternyata sudah berdiri. Ia TIDAK kosong hanya pada keadaan tidak normal -
	// lihat JalankanMigrasi - sehingga pemanggil yang menemukannya berisi tahu
	// bahwa ada langkah yang dulu gagal separuh jalan.
	ObjekSudahAda []string
}

// pecahPernyataan memecah isi berkas .sql menjadi pernyataan terpisah.
//
// Pemisahnya adalah baris yang HANYA berisi tanda garis miring - konvensi
// SQL*Plus. Memakai titik koma sebagai pemisah tidak aman karena titik koma
// juga muncul di dalam blok PL/SQL.
func pecahPernyataan(isi string) []string {
	var out []string
	var sekarang []string
	for _, baris := range strings.Split(isi, "\n") {
		if strings.TrimSpace(baris) == "/" {
			if p := gabung(sekarang); p != "" {
				out = append(out, p)
			}
			sekarang = nil
			continue
		}
		sekarang = append(sekarang, baris)
	}
	if p := gabung(sekarang); p != "" {
		out = append(out, p)
	}
	return out
}

// gabung menyatukan baris menjadi satu pernyataan, membuang baris komentar
// murni dan spasi di tepi. Mengembalikan teks kosong bila tidak ada isi.
func gabung(baris []string) string {
	var isi []string
	for _, b := range baris {
		// U+FEFF (byte order mark) menempel di awal berkas bila penyunting
		// menulisnya sebagai UTF-8 ber-BOM. TrimSpace TIDAK membuangnya,
		// sehingga tanpa baris ini satu baris komentar dapat lolos menjadi
		// bagian pernyataan SQL dan Oracle menolaknya. Pernah terjadi di
		// berkas 004 - karena itu dijaga di sini, bukan hanya dibersihkan
		// sekali di berkasnya.
		b = strings.TrimPrefix(b, "\ufeff")
		if strings.HasPrefix(strings.TrimSpace(b), "--") {
			continue
		}
		isi = append(isi, b)
	}
	return strings.TrimSpace(strings.Join(isi, "\n"))
}

// daftarMigrasi membaca seluruh langkah maju, terurut menurut namanya.
func daftarMigrasi(mundur bool) ([]Migrasi, error) {
	entri, err := berkasMigrasi.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("repository: membaca folder migrasi: %w", err)
	}
	var nama []string
	for _, e := range entri {
		n := e.Name()
		if !strings.HasSuffix(n, ".sql") {
			continue
		}
		if strings.HasSuffix(n, "_down.sql") != mundur {
			continue
		}
		nama = append(nama, n)
	}
	sort.Strings(nama)
	if mundur {
		// Jalur mundur berjalan MENURUN supaya anak hilang sebelum induknya.
		for i, j := 0, len(nama)-1; i < j; i, j = i+1, j-1 {
			nama[i], nama[j] = nama[j], nama[i]
		}
	}

	out := make([]Migrasi, 0, len(nama))
	for _, n := range nama {
		isi, err := berkasMigrasi.ReadFile(path.Join("migrations", n))
		if err != nil {
			return nil, fmt.Errorf("repository: membaca %s: %w", n, err)
		}
		out = append(out, Migrasi{Nama: n, Pernyataan: pecahPernyataan(string(isi))})
	}
	return out, nil
}

// PernyataanLangkah mengembalikan pernyataan SQL satu langkah migrasi, apa
// adanya, dengan penanda {skema} MASIH utuh.
//
// Dipakai test db supaya ia dapat menjalankan pernyataan yang BENAR-BENAR ada
// di langkah itu - bukan pernyataan tiruan yang bentuknya berbeda. Menirukan
// kegagalan separuh jalan dengan objek berbentuk lain akan menguji hal lain.
func PernyataanLangkah(nama string) ([]string, error) {
	isi, err := berkasMigrasi.ReadFile(path.Join("migrations", nama))
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", nama, err)
	}
	return pecahPernyataan(string(isi)), nil
}

// kunciLangkah menyamakan nama berkas maju dan mundur menjadi satu kunci,
// supaya jalur mundur tahu langkah mana yang dibatalkannya.
func kunciLangkah(nama string) string {
	return strings.TrimSuffix(strings.TrimSuffix(nama, ".sql"), "_down")
}

// ---------------------------------------------------------------------------
// Pra-terbang bentuk tabel
//
// Ronde 3 menutup separuh masalah: sebuah CREATE yang dilewati kini dibuktikan
// KEBERADAAN objeknya di katalog. Yang tidak diperiksa adalah BENTUKNYA. Pada
// 26-09-2026 kelemahan itu terbukti nyata: DOCUMENT_CLAIM sudah ada di skema
// warisan dengan empat belas kolom yang sama sekali berbeda dari DDL 007.
// Tanpa pemeriksaan bentuk, migrasi akan melewatinya, mencatat langkahnya
// sukses, dan aplikasi berjalan di atas tabel yang kolomnya bukan miliknya.
//
// `[keputusan work owner 26-09-2026, butir x]`: bentuk diperiksa SEBELUM satu
// pernyataan pun dikirim. Yang dibandingkan hanya NAMA kolom, tanpa memandang
// urutan maupun huruf besar-kecil - tipe dan panjang sengaja tidak, sebab
// selisih tipe belum tentu salah dan akan menghasilkan penolakan palsu.
// ---------------------------------------------------------------------------

var (
	// polaCreateTabel memisahkan nama tabel dari badan CREATE TABLE.
	polaCreateTabel = regexp.MustCompile(`(?is)CREATE\s+TABLE\s+\{skema\}\.(\w+)\s*\((.*?)\n\)`)
	// polaKolomDDL mengenali satu baris deklarasi kolom di dalam badan itu.
	polaKolomDDL = regexp.MustCompile(`^([A-Z][A-Z0-9_]*)\s+\S`)
)

// tampakCreateTabel mengenali pernyataan CREATE TABLE secara longgar.
//
// Dipakai HANYA untuk menolak pernyataan yang terlihat seperti CREATE TABLE
// tetapi tidak dapat diurai KolomCreateTable. Tanpa penjaga ini, satu tabel
// yang ditulis dengan bentuk lain - misalnya seluruhnya dalam satu baris -
// akan melewati pemeriksaan bentuk tanpa ada yang menyadarinya.
var polaCreateLonggar = regexp.MustCompile(`(?is)CREATE\s+TABLE\s`)

func tampakCreateTabel(pernyataan string) bool {
	return polaCreateLonggar.MatchString(pernyataan)
}

// KolomCreateTable membaca nama tabel dan daftar kolomnya dari satu pernyataan.
//
// Nama tabel kosong berarti pernyataan itu bukan CREATE TABLE. Baris CONSTRAINT
// dan REFERENCES dilewati: keduanya bukan kolom.
func KolomCreateTable(pernyataan string) (string, []string) {
	m := polaCreateTabel.FindStringSubmatch(pernyataan)
	if m == nil {
		return "", nil
	}
	var kolom []string
	for _, b := range strings.Split(m[2], "\n") {
		atas := strings.ToUpper(strings.TrimSpace(b))
		if atas == "" || strings.HasPrefix(atas, "CONSTRAINT") ||
			strings.HasPrefix(atas, "REFERENCES") {
			continue
		}
		if k := polaKolomDDL.FindStringSubmatch(atas); k != nil {
			kolom = append(kolom, k[1])
		}
	}
	return strings.ToUpper(m[1]), kolom
}

// SelisihKolom menyebut kolom yang diminta DDL tetapi tidak ada di katalog, dan
// sebaliknya.
//
// Perbandingannya tanpa urutan dan tanpa huruf besar-kecil: Oracle menyimpan
// pengenal tanpa kutip dalam huruf besar, dan urutan kolom katalog mengikuti
// urutan pembuatan, bukan urutan berkas migrasi.
func SelisihKolom(ddl, katalog []string) (kurang, lebih []string) {
	adaDi := func(daftar []string) map[string]bool {
		m := make(map[string]bool, len(daftar))
		for _, n := range daftar {
			m[strings.ToUpper(strings.TrimSpace(n))] = true
		}
		return m
	}
	diKatalog, diDDL := adaDi(katalog), adaDi(ddl)
	for _, n := range ddl {
		if !diKatalog[strings.ToUpper(strings.TrimSpace(n))] {
			kurang = append(kurang, strings.ToUpper(strings.TrimSpace(n)))
		}
	}
	for _, n := range katalog {
		if !diDDL[strings.ToUpper(strings.TrimSpace(n))] {
			lebih = append(lebih, strings.ToUpper(strings.TrimSpace(n)))
		}
	}
	sort.Strings(kurang)
	sort.Strings(lebih)
	return kurang, lebih
}

// kolomKatalog membaca nama kolom sebuah tabel dari katalog Oracle.
func (d *DB) kolomKatalog(ctx context.Context, tabel string) ([]string, error) {
	q := `SELECT COLUMN_NAME FROM SYS.ALL_TAB_COLUMNS
	        WHERE UPPER(OWNER) = UPPER(:1) AND UPPER(TABLE_NAME) = UPPER(:2)
	        ORDER BY COLUMN_ID`
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := d.sql.QueryContext(ctx, q, d.skema, tabel)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca kolom %s: %w", tabel, err)
	}
	defer func() { _ = baris.Close() }()
	var out []string
	for baris.Next() {
		var n string
		if err := baris.Scan(&n); err != nil {
			return nil, fmt.Errorf("repository: membaca kolom %s: %w", tabel, err)
		}
		out = append(out, n)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca kolom %s: %w", tabel, err)
	}
	return out, nil
}

// praTerbangBentuk memeriksa SELURUH langkah sebelum satu pun dijalankan.
//
// Mengembalikan himpunan nama tabel yang sudah ada DAN bentuknya cocok -
// itulah yang boleh dilewati. Tabel yang bentuknya berbeda menghentikan
// migrasi di sini, sebelum ada yang ditulis ke mana pun.
func (d *DB) praTerbangBentuk(ctx context.Context, langkah []Migrasi,
	selesai map[string]bool) (map[string]bool, error) {
	cocok := map[string]bool{}
	for _, m := range langkah {
		// ⛔ Langkah yang SUDAH tercatat tidak diperiksa, dan itu bukan sekadar
		// penghematan. Tabel yang dibuat langkah selesai adalah milik basis
		// data sejak saat itu: DBA boleh menambahkan kolom audit padanya, dan
		// itu sah. Memeriksanya di sini membuat satu kolom tambahan
		// MENGGAGALKAN seluruh migrasi berikutnya - termasuk langkah baru yang
		// tidak ada hubungannya - dan satu-satunya pemulihan adalah membongkar
		// seluruh skema. Yang dijaga pra-terbang hanyalah langkah yang HENDAK
		// dijalankan.
		if selesai[kunciLangkah(m.Nama)] {
			continue
		}
		for _, p := range m.Pernyataan {
			nama, kolomDDL := KolomCreateTable(p)
			if nama == "" {
				// ⛔ Pernyataan yang TAMPAK CREATE TABLE tetapi tidak terurai
				// tidak boleh lewat diam-diam: ia akan jatuh ke pemeriksaan
				// keberadaan saja, dan bentuknya tidak pernah dibandingkan -
				// persis lubang yang pra-terbang ini dibuat untuk menutup.
				if tampakCreateTabel(p) {
					return nil, fmt.Errorf(
						"repository: migrasi %s: ada CREATE TABLE yang tidak dapat diurai "+
							"nama dan kolomnya, sehingga bentuknya tidak dapat diperiksa. "+
							"Tulis ulang pernyataannya mengikuti bentuk berkas migrasi lain", m.Nama)
				}
				continue
			}
			ada, err := d.objekAda(ctx, nama)
			if err != nil {
				return nil, fmt.Errorf("repository: migrasi %s: %w", m.Nama, err)
			}
			if !ada {
				continue // akan dibuat; tidak ada bentuk untuk dibandingkan
			}
			kolomKat, err := d.kolomKatalog(ctx, nama)
			if err != nil {
				return nil, fmt.Errorf("repository: migrasi %s: %w", m.Nama, err)
			}
			kurang, lebih := SelisihKolom(kolomDDL, kolomKat)
			if len(kurang) > 0 || len(lebih) > 0 {
				return nil, fmt.Errorf(
					"repository: migrasi %s: tabel %s sudah ada di skema %s tetapi BENTUKNYA BERBEDA "+
						"- kolom yang diminta migrasi tetapi tidak ada: %v; kolom yang ada tetapi "+
						"tidak diminta: %v. Migrasi dihentikan sebelum satu pernyataan pun dikirim; "+
						"tidak ada yang diubah",
					m.Nama, nama, d.skema, kurang, lebih)
			}
			cocok[nama] = true
		}
	}
	return cocok, nil
}

// siapkanTabelMigrasi membuat tabel pencatat bila belum ada.
func (d *DB) siapkanTabelMigrasi(ctx context.Context) error {
	tabel, err := d.Qualify(namaTabelMigrasi)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`CREATE TABLE %s (
		NAMA            VARCHAR2(128) NOT NULL,
		DIJALANKAN_PADA DATE,
		CONSTRAINT PK_%s PRIMARY KEY (NAMA))`, tabel, namaTabelMigrasi)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := d.sql.ExecContext(ctx, q); err != nil {
		if sudahAda(err) {
			return nil
		}
		return fmt.Errorf("repository: membuat tabel migrasi: %w", err)
	}
	return nil
}

// sudahAda mengenali SATU galat Oracle: nama yang hendak dipakai sudah dipakai
// objek lain - tabel, index, atau sequence (ORA-00955). Itu yang membuat
// pembuatan objek menjadi idempoten.
//
// ⛔ ORA-02264 sengaja TIDAK termasuk, dan ini ralat atas ronde 2.
// ORA-02264 berarti "nama itu sudah dipakai constraint lain", dan Oracle baru
// memeriksanya ketika TABELNYA BELUM ADA - bila tabelnya sudah ada, ia menjawab
// ORA-00955 lebih dulu. Jadi ORA-02264 pada sebuah CREATE TABLE berarti tabel
// itu JUSTRU TIDAK terbuat. Menelannya berarti mencatat langkah sebagai sukses
// atas tabel yang tidak pernah ada - kegagalan yang baru ketahuan jauh di hilir.
func sudahAda(err error) bool {
	return err != nil && strings.Contains(err.Error(), "ORA-00955")
}

// polaObjekDibuat menangkap nama objek yang dibuat sebuah pernyataan CREATE.
var polaObjekDibuat = regexp.MustCompile(
	`(?is)CREATE\s+(?:UNIQUE\s+)?(?:TABLE|INDEX|SEQUENCE)\s+\{skema\}\.(\w+)`)

// namaObjekDibuat membaca nama objek dari pernyataan CREATE, atau teks kosong
// bila pernyataannya bukan CREATE yang dikenali.
func namaObjekDibuat(pernyataan string) string {
	m := polaObjekDibuat.FindStringSubmatch(pernyataan)
	if m == nil {
		return ""
	}
	return strings.ToUpper(m[1])
}

// objekAda bertanya ke katalog Oracle apakah sebuah objek benar-benar berdiri.
//
// Dipakai sesudah sebuah CREATE dilewati karena menjawab ORA-00955. Oracle
// hanya bilang "nama sudah dipakai"; ia tidak bilang objek APA. Tanpa
// pemeriksaan ini, nama yang dipakai objek jenis lain akan lolos sebagai
// "sudah ada".
//
// ⚠️ Yang diperiksa keberadaan, BUKAN bentuk. Kolom yang berbeda tidak
// tertangkap di sini; yang menjaga hal itu jalur mundur - bongkar dulu, baru
// pasang lagi.
//
// SYS.ALL_OBJECTS ditulis berawalan skema, bukan telanjang (ADR-U-0033). Ia
// dipilih di atas USER_OBJECTS karena skema sasaran belum tentu sama dengan
// pengguna sambungan, dan di atas DBA_OBJECTS karena yang terakhir menuntut hak
// istimewa yang tidak perlu. ALL_OBJECTS memperlihatkan objek yang terlihat oleh
// sesi - dan sesi ini baru saja mencoba membuat objek di skema itu, jadi ia
// memang punya akses ke sana.
func (d *DB) objekAda(ctx context.Context, nama string) (bool, error) {
	// UPPER di kedua sisi: Oracle menyimpan pengenal tanpa kutip dalam huruf
	// besar, tetapi ORACLE_SCHEMA datang dari env var dan boleh ditulis
	// bagaimana saja. Tanpa ini, skema yang ditulis huruf kecil membuat
	// pembuktian menjawab "tidak ada" dan MENGGAGALKAN migrasi yang sehat.
	q := `SELECT COUNT(*) FROM SYS.ALL_OBJECTS
	        WHERE UPPER(OWNER) = UPPER(:1) AND UPPER(OBJECT_NAME) = UPPER(:2)`
	if err := PeriksaSQL(q); err != nil {
		return false, err
	}
	var n int
	if err := d.sql.QueryRowContext(ctx, q, d.skema, nama).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: memeriksa keberadaan %s: %w", nama, err)
	}
	return n > 0, nil
}

// pernyataanBuat menjawab apakah sebuah pernyataan SQL membuat objek baru.
//
// Hanya pernyataan CREATE yang boleh dilewati saat objeknya sudah ada. ALTER,
// INSERT, atau DROP yang menjawab galat serupa BUKAN keadaan yang sama, dan
// menelannya akan menyembunyikan kerusakan sungguhan.
func pernyataanBuat(q string) bool {
	return strings.HasPrefix(strings.ToUpper(strings.TrimSpace(q)), "CREATE")
}

// ringkasPernyataan mengambil beberapa kata pertama sebuah pernyataan, supaya
// laporan menyebut objek yang dimaksud tanpa menyalin seluruh DDL-nya.
func ringkasPernyataan(q string) string {
	kata := strings.Fields(q)
	if len(kata) > 4 {
		kata = kata[:4]
	}
	return strings.Join(kata, " ")
}

// tidakAda mengenali galat Oracle "tabel atau view tidak ada" (ORA-00942)
// dan "sequence tidak ada" (ORA-02289). Jalur mundur memakainya supaya
// membongkar yang memang belum ada bukan kegagalan.
func tidakAda(err error) bool {
	if err == nil {
		return false
	}
	p := err.Error()
	return strings.Contains(p, "ORA-00942") || strings.Contains(p, "ORA-02289")
}

// sudahDijalankan membaca daftar langkah yang tercatat.
func (d *DB) sudahDijalankan(ctx context.Context) (map[string]bool, error) {
	tabel, err := d.Qualify(namaTabelMigrasi)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT NAMA FROM %s`, tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := d.sql.QueryContext(ctx, q)
	if err != nil {
		if tidakAda(err) {
			return map[string]bool{}, nil
		}
		return nil, fmt.Errorf("repository: membaca catatan migrasi: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := map[string]bool{}
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		out[n] = true
	}
	return out, rows.Err()
}

// JalankanMigrasi menjalankan seluruh langkah yang belum pernah dijalankan.
//
// Idempoten: langkah yang sudah tercatat dilewati, sehingga memanggilnya dua
// kali berturut-turut aman.
func (d *DB) JalankanMigrasi(ctx context.Context) (LaporanMigrasi, error) {
	var lap LaporanMigrasi
	if d == nil || d.sql == nil {
		return lap, ErrTanpaOracle
	}
	if d.isPegaProd {
		return lap, ErrMigrasiDiProduksi
	}
	if err := d.siapkanTabelMigrasi(ctx); err != nil {
		return lap, err
	}
	selesai, err := d.sudahDijalankan(ctx)
	if err != nil {
		return lap, err
	}
	langkah, err := daftarMigrasi(false)
	if err != nil {
		return lap, err
	}
	tabel, err := d.Qualify(namaTabelMigrasi)
	if err != nil {
		return lap, err
	}

	// Pra-terbang: bentuk tiap tabel yang sudah ada diperiksa lebih dulu, dan
	// migrasi berhenti di sini bila ada yang berbeda (butir x).
	bentukCocok, err := d.praTerbangBentuk(ctx, langkah, selesai)
	if err != nil {
		return lap, err
	}

	for _, m := range langkah {
		kunci := kunciLangkah(m.Nama)
		if selesai[kunci] {
			lap.Dilewati = append(lap.Dilewati, kunci)
			continue
		}
		for _, p := range m.Pernyataan {
			q := strings.ReplaceAll(p, "{skema}", d.skema)
			if err := PeriksaSQL(q); err != nil {
				return lap, fmt.Errorf("%s: %w", m.Nama, err)
			}
			if _, err := d.sql.ExecContext(ctx, q); err != nil {
				// DDL Oracle menutup transaksinya sendiri. Bila sebuah langkah
				// dulu gagal di pernyataan kedua, pernyataan pertamanya sudah
				// terlanjur jadi dan T_MIGRASI tidak sempat mencatat apa pun -
				// sehingga percobaan berikutnya mengulang langkah itu dari awal
				// dan mati di ORA-00955. Objek yang sudah berdiri dilewati,
				// dicatat, lalu langkahnya diteruskan sampai tuntas.
				// ✅ Risiko bentuk DITUTUP 26-09-2026 (butir x): pra-terbang
				// di atas sudah membandingkan kolom tiap tabel yang ada dengan
				// katalog, dan migrasi tidak sampai ke sini bila ada yang
				// berbeda. Yang boleh dilewati hanyalah tabel yang namanya
				// tercatat di bentukCocok.
				if pernyataanBuat(q) && sudahAda(err) {
					// ⛔ Dilewati TIDAK cukup - keberadaannya dibuktikan.
					// Oracle hanya bilang namanya terpakai, bukan bahwa objek
					// yang kita maksud ada. Bila ternyata tidak ada, langkah
					// ini GAGAL dan tidak dicatat di T_MIGRASI.
					nama := namaObjekDibuat(p)
					if nama == "" {
						return lap, fmt.Errorf(
							"repository: migrasi %s: dilaporkan sudah ada, nama objeknya tidak terbaca: %w",
							m.Nama, err)
					}
					// Tabel: cukup dibuktikan lewat pra-terbang, yang sudah
					// memeriksa keberadaan DAN bentuk. Nama yang tidak ada di
					// sana berarti tabelnya muncul sesudah pra-terbang - dan
					// bentuknya belum pernah diperiksa siapa pun.
					if namaTabel, _ := KolomCreateTable(p); namaTabel != "" {
						if !bentukCocok[namaTabel] {
							return lap, fmt.Errorf(
								"repository: migrasi %s: tabel %s dilaporkan sudah ada, "+
									"tetapi bentuknya belum pernah diperiksa - ia muncul sesudah "+
									"pra-terbang. Migrasi dihentikan: %w", m.Nama, namaTabel, err)
						}
						lap.ObjekSudahAda = append(lap.ObjekSudahAda, nama)
						continue
					}
					ada, errPeriksa := d.objekAda(ctx, nama)
					if errPeriksa != nil {
						return lap, fmt.Errorf("repository: migrasi %s: %w",
							m.Nama, errPeriksa)
					}
					if !ada {
						return lap, fmt.Errorf(
							"repository: migrasi %s: %s dilaporkan sudah ada, tidak ditemukan di katalog: %w",
							m.Nama, nama, err)
					}
					lap.ObjekSudahAda = append(lap.ObjekSudahAda,
						fmt.Sprintf("%s: %s", m.Nama, ringkasPernyataan(q)))
					continue
				}
				return lap, fmt.Errorf("repository: migrasi %s: %w", m.Nama, err)
			}
			lap.Pernyataan++
		}
		qCatat := fmt.Sprintf(`INSERT INTO %s (NAMA, DIJALANKAN_PADA) VALUES (:1, :2)`, tabel)
		if err := PeriksaSQL(qCatat); err != nil {
			return lap, err
		}
		if _, err := d.sql.ExecContext(ctx, qCatat, kunci, time.Now()); err != nil {
			return lap, fmt.Errorf("repository: mencatat migrasi %s: %w", kunci, err)
		}
		lap.Dijalankan = append(lap.Dijalankan, kunci)
	}
	return lap, nil
}

// BongkarMigrasi menjalankan jalur mundur untuk langkah yang sudah dijalankan.
//
// Urutannya menurun, sehingga tabel anak dibongkar sebelum induknya. Langkah
// yang objeknya memang sudah tidak ada dilewati tanpa galat.
func (d *DB) BongkarMigrasi(ctx context.Context) (LaporanMigrasi, error) {
	var lap LaporanMigrasi
	if d == nil || d.sql == nil {
		return lap, ErrTanpaOracle
	}
	if d.isPegaProd {
		return lap, ErrMigrasiDiProduksi
	}
	langkah, err := daftarMigrasi(true)
	if err != nil {
		return lap, err
	}
	tabel, err := d.Qualify(namaTabelMigrasi)
	if err != nil {
		return lap, err
	}
	// Hanya langkah yang TERCATAT pernah dijalankan yang dibongkar. Tanpa
	// pembacaan ini, laporan akan menyebut langkah yang tidak berbuat apa-apa
	// sebagai "dijalankan".
	selesai, err := d.sudahDijalankan(ctx)
	if err != nil {
		return lap, err
	}

	for _, m := range langkah {
		kunci := kunciLangkah(m.Nama)
		if !selesai[kunci] {
			lap.Dilewati = append(lap.Dilewati, kunci)
			continue
		}
		for _, p := range m.Pernyataan {
			q := strings.ReplaceAll(p, "{skema}", d.skema)
			if err := PeriksaSQL(q); err != nil {
				return lap, fmt.Errorf("%s: %w", m.Nama, err)
			}
			if _, err := d.sql.ExecContext(ctx, q); err != nil {
				if tidakAda(err) {
					continue
				}
				return lap, fmt.Errorf("repository: bongkar %s: %w", m.Nama, err)
			}
			lap.Pernyataan++
		}
		// DELETE ini pembukuan migrasi, BUKAN jalur pengguna - ADR-U-0031
		// melarang DELETE pada jalur pengguna, bukan pada catatan langkah.
		qHapus := fmt.Sprintf(`DELETE FROM %s WHERE NAMA = :1`, tabel)
		if err := PeriksaSQL(qHapus); err != nil {
			return lap, err
		}
		if _, err := d.sql.ExecContext(ctx, qHapus, kunci); err != nil {
			if !tidakAda(err) {
				return lap, fmt.Errorf("repository: menghapus catatan %s: %w", kunci, err)
			}
		}
		lap.Dijalankan = append(lap.Dijalankan, kunci)
	}
	return lap, nil
}
