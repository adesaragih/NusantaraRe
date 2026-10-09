// Package config membaca seluruh konfigurasi dari environment variable.
//
// ⛔ Nol literal host, endpoint, kata sandi, atau DSN di berkas ini maupun di
// berkas lain mana pun (ADR-U-0013, CLAUDE.md bab 10). Daftar endpoint
// sesungguhnya ada di tabel Oracle M_LINK_SERVICE yang isinya belum ada di
// korpus (OQ-047) - sampai ia ada, alamat datang dari env.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Prefix env var untuk alamat layanan luar: SERVICE_<NAMA>.
const prefixLayanan = "SERVICE_"

// Config adalah seluruh konfigurasi proses.
type Config struct {
	// HTTPAddr alamat dengar server. Hanya port secara bawaan; tidak ada
	// nama host yang ditanam.
	HTTPAddr string

	// OracleDSN kosong berarti backend berjalan tanpa Oracle. Fase 0
	// mengizinkannya supaya GET /healthz dapat dijawab tanpa instance; tiket
	// penyimpanan mana pun akan menuntutnya terisi.
	OracleDSN string

	// OracleSchema wajib terisi bila OracleDSN terisi: setiap query menyebut
	// nama skema secara eksplisit (ADR-U-0033).
	OracleSchema string

	// IsPegaProd menandai lingkungan yang menunjuk data produksi Pega
	// (ADR-U-0005). Dipakai sebagai gerbang perilaku, bukan sebagai hiasan.
	IsPegaProd bool

	// StorageTokenSalt adalah garam token penyimpanan berkas (butir an).
	//
	// ⛔ RAHASIA. Ia TIDAK PERNAH disalin dari procedure ke repo, tidak
	// dicetak, dan tidak masuk artefak mana pun. Diminta ke DBA lewat jalur
	// rahasia - bukan lewat repo maupun percakapan.
	//
	// Kosong berarti token penyimpanan GAGAL TERANG; itu keadaan yang benar
	// sampai garamnya diberikan.
	StorageTokenSalt string

	// SkemaUjiDiakui adalah pengakuan sadar bahwa skema yang ditunjuk
	// OracleSchema boleh DIHAPUS isinya. Dibaca dari ORACLE_SKEMA_UJI.
	//
	// Hanya operasi merusak yang memakainya - test bertag db dan
	// `-migrate-down`. `-migrate` tidak: ia hanya membuat objek.
	SkemaUjiDiakui bool

	// AuthStub menyalakan pembacaan pelaku dari header HTTP.
	//
	// ⛔ INI BUKAN AUTENTIKASI. Header dapat ditulis siapa saja yang dapat
	// mengirim permintaan, jadi ia tidak membuktikan apa pun. Ia penunda
	// sampai tiket 07 / IAM memasang sumber peran yang sebenarnya
	// (ADR-U-0030), supaya jalur wewenang dapat dibangun dan diuji lebih
	// dulu. Mati secara bawaan, dan ditolak keras saat IS_PEGA_PROD=true.
	AuthStub bool

	// SesiRahasia adalah kunci HMAC cookie sesi login (SESI_RAHASIA),
	// minimal 32 karakter. Kosong = login menjawab 503 (stub tetap dapat
	// dipakai di pengembangan). Keputusan work owner 01-10-2026.
	//
	// ⛔ Nilainya tidak pernah ditulis di repo, di log, atau di percakapan.
	SesiRahasia string

	// SesiCookieAman adalah atribut Secure cookie sesi (SESI_COOKIE_SECURE,
	// bawaan true). `false` hanya untuk server http tanpa TLS di jaringan
	// uji - dan ditolak saat IS_PEGA_PROD=true.
	SesiCookieAman bool

	// UnggahanDir adalah folder LOKAL tempat berkas unggahan disimpan
	// (butir be).
	//
	// ⛔ BUKAN korpus, dan bukan folder mana pun di bawah `D:/XML/RNM_BRD`
	// selain `OUTPUT_HASIL_RNM`. Korpus dibaca sebagai bukti; ia tidak
	// pernah ditulisi.
	//
	// Kosong berarti unggahan GAGAL TERANG - itu keadaan yang benar sampai
	// seseorang memilih foldernya. Bawaan diam-diam ke `./unggahan` akan
	// membuat berkas pelanggan mendarat di folder kerja siapa pun yang
	// kebetulan menjalankan server.
	UnggahanDir string

	// ModulAktif - modul yang dipasang proses ini (refactor bentuk B), dari
	// MODUL_AKTIF: nama dipisah koma, huruf kecil, tanpa ganda. KOSONG
	// berarti SEMUA modul terdaftar aktif - bawaan yang sama dengan sebelum
	// variabel ini ada. Namanya diperiksa `cmd/api` terhadap daftar modul;
	// paket ini tidak mengenal nama modul mana pun.
	//
	// ⛔ Migrasi TIDAK membaca ini: `-migrate` selalu menjalankan migrasi
	// SEMUA modul terdaftar, sehingga skema selalu utuh.
	ModulAktif []string

	// Layanan memetakan nama layanan luar ke alamatnya, seluruhnya dari env.
	Layanan map[string]string
}

// ErrBukanSkemaUji menolak operasi merusak di luar skema uji.
var ErrBukanSkemaUji = errors.New("konfigurasi: menolak operasi merusak di luar skema uji")

// EnvSkemaUji adalah nama env var pengakuan itu.
const EnvSkemaUji = "ORACLE_SKEMA_UJI"

// NamaSkemaWarisan adalah skema Pega yang memuat tabel warisan sungguhan.
const NamaSkemaWarisan = "POOLDATA"

// PagarSkemaUji memutuskan apakah operasi MERUSAK boleh menyentuh skema ini.
//
// Dipisah dari environment supaya dapat diuji tanpa Oracle dan tanpa env var.
// Dua syarat, KEDUANYA harus benar:
//  1. env ORACLE_SKEMA_UJI bernilai "true"
//  2. nama skema tidak memuat POOLDATA
//
// Syarat kedua tidak dapat ditutupi oleh syarat pertama: menyetel env tidak
// membuat skema warisan menjadi skema uji.
//
// Tempatnya di sini, bukan di paket skema uji, sebab DUA pemanggil
// membutuhkannya: test bertag db dan flag `-migrate-down` di cmd/api - dan
// cmd/api tidak boleh mengimpor paket penunjang test.
func PagarSkemaUji(skema, diakui string) error {
	if strings.TrimSpace(strings.ToLower(diakui)) != "true" {
		return fmt.Errorf("%w: env %s belum bernilai true; operasi ini MENGHAPUS tabel "+
			"di skema yang ditunjuk ORACLE_SCHEMA. Setel %s=true hanya bila skema itu "+
			"memang skema uji kosong dari DBA", ErrBukanSkemaUji, EnvSkemaUji, EnvSkemaUji)
	}
	// MEMUAT, bukan sama persis: POOLDATA_DEV dan POOLDATA2 adalah skema
	// warisan juga. Untuk operasi yang MENGHAPUS, menolak terlalu banyak jauh
	// lebih murah daripada meloloskan satu yang salah.
	if strings.Contains(strings.ToUpper(strings.TrimSpace(skema)), NamaSkemaWarisan) {
		return fmt.Errorf("%w: ORACLE_SCHEMA menunjuk %s, yang memuat nama skema warisan Pega %s; "+
			"tabel sungguhan ada di sana. Pakai skema uji kosong, dan %s=true tidak mengubah hal ini",
			ErrBukanSkemaUji, skema, NamaSkemaWarisan, EnvSkemaUji)
	}
	return nil
}

// PastikanSkemaUji menjalankan pagar itu atas konfigurasi yang terbaca.
func (c Config) PastikanSkemaUji() error {
	diakui := "false"
	if c.SkemaUjiDiakui {
		diakui = "true"
	}
	return PagarSkemaUji(c.OracleSchema, diakui)
}

// ErrKonfigurasi membungkus seluruh kegagalan pembacaan konfigurasi.
var ErrKonfigurasi = errors.New("konfigurasi")

// bacaDaftarModul mengurai MODUL_AKTIF: dipisah koma, spasi dibuang, huruf
// kecil, butir kosong dan ganda dilewati. Kosong seluruhnya = nil (semua).
func bacaDaftarModul(raw string) []string {
	var hasil []string
	sudah := map[string]bool{}
	for _, n := range strings.Split(raw, ",") {
		n = strings.ToLower(strings.TrimSpace(n))
		if n == "" || sudah[n] {
			continue
		}
		sudah[n] = true
		hasil = append(hasil, n)
	}
	return hasil
}

// Load membaca konfigurasi dari environment.
func Load() (Config, error) {
	c := Config{
		HTTPAddr:     ambil("HTTP_ADDR", ":8080"),
		OracleDSN:    os.Getenv("ORACLE_DSN"),
		OracleSchema: strings.TrimSpace(os.Getenv("ORACLE_SCHEMA")),
		Layanan:      map[string]string{},
	}

	c.SkemaUjiDiakui = strings.EqualFold(strings.TrimSpace(os.Getenv(EnvSkemaUji)), "true")
	c.AuthStub = strings.EqualFold(strings.TrimSpace(os.Getenv("AUTH_STUB")), "true")
	c.SesiRahasia = strings.TrimSpace(os.Getenv("SESI_RAHASIA"))
	c.SesiCookieAman = !strings.EqualFold(strings.TrimSpace(os.Getenv("SESI_COOKIE_SECURE")), "false")
	// ⛔ Tidak di-TrimSpace dan tidak dinormalkan: garam adalah byte apa
	// adanya, dan membetulkannya diam-diam menghasilkan token yang berbeda
	// dari yang sistem lama terbitkan.
	c.StorageTokenSalt = os.Getenv("STORAGE_TOKEN_SALT")
	c.UnggahanDir = strings.TrimSpace(os.Getenv("UNGGAHAN_DIR"))
	c.ModulAktif = bacaDaftarModul(os.Getenv("MODUL_AKTIF"))

	raw := strings.TrimSpace(os.Getenv("IS_PEGA_PROD"))
	if raw != "" {
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return Config{}, fmt.Errorf("%w: IS_PEGA_PROD bukan boolean: %q", ErrKonfigurasi, raw)
		}
		c.IsPegaProd = b
	}

	// ⛔ Stub pelaku TIDAK PERNAH hidup di lingkungan produksi. Ini ditolak
	// saat memuat konfigurasi, bukan saat permintaan pertama tiba: proses yang
	// menyala dengan keduanya benar akan melayani permintaan sebelum ada yang
	// sempat menyadarinya.
	if c.AuthStub && c.IsPegaProd {
		return Config{}, fmt.Errorf(
			"%w: AUTH_STUB=true ditolak saat IS_PEGA_PROD=true; stub pelaku membaca peran "+
				"dari header HTTP dan tidak membuktikan apa pun (ADR-U-0005, ADR-U-0030)",
			ErrKonfigurasi)
	}

	// ⛔ Rahasia pendek dapat ditebak; cookie tanpa Secure di produksi dapat
	// disadap. Keduanya ditolak saat menyala, bukan saat login pertama.
	if c.SesiRahasia != "" && len(c.SesiRahasia) < PanjangMinSesiRahasia {
		return Config{}, fmt.Errorf("%w: SESI_RAHASIA minimal %d karakter", ErrKonfigurasi, PanjangMinSesiRahasia)
	}
	if !c.SesiCookieAman && c.IsPegaProd {
		return Config{}, fmt.Errorf("%w: SESI_COOKIE_SECURE=false ditolak saat IS_PEGA_PROD=true", ErrKonfigurasi)
	}

	if c.OracleDSN != "" && c.OracleSchema == "" {
		return Config{}, fmt.Errorf(
			"%w: ORACLE_SCHEMA wajib terisi bila ORACLE_DSN terisi (ADR-U-0033)", ErrKonfigurasi)
	}

	for _, kv := range os.Environ() {
		nama, nilai, ada := strings.Cut(kv, "=")
		if !ada || !strings.HasPrefix(nama, prefixLayanan) || nilai == "" {
			continue
		}
		c.Layanan[strings.TrimPrefix(nama, prefixLayanan)] = nilai
	}

	return c, nil
}

// PanjangMinSesiRahasia - panjang minimal SESI_RAHASIA.
const PanjangMinSesiRahasia = 32

// PunyaOracle menyatakan apakah proses dikonfigurasi menyentuh Oracle.
func (c Config) PunyaOracle() bool { return c.OracleDSN != "" }

func ambil(nama, bawaan string) string {
	if v := strings.TrimSpace(os.Getenv(nama)); v != "" {
		return v
	}
	return bawaan
}
