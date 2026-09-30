package penjaga

// Penjaga SELURUH APLIKASI - refactor bentuk B paket 8 (30-09-2026).
//
// Aturan di berkas ini berlaku untuk setiap modul sekaligus, jadi rumahnya
// `inti`, bukan satu modul: alamat layanan (ADR-U-0013), nama orang, nama
// tabel telanjang (ADR-U-0033), arah handlers -> services -> repository, dan
// nol tabel baru Treaty Contract Out (tco4). Masing-masing dipindah APA ADANYA
// dari modul tempat ia lahir (Claim Life services/repository, Treaty
// repository); akarnya kini akar aplikasi dari `inti/penjaga`, dan cacah log
// berangkanya dibandingkan dengan sebelum pindah.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// buangKomentar menghapus baris komentar sebelum pencocokan, supaya prosa yang
// MENYEBUT sebuah pemanggilan tidak dituduh sebagai pemanggilannya.
// Salinan `modul/claimlife/services/jejak_statik_test.go` (paket berbeda).
func buangKomentar(isi string) string {
	var b strings.Builder
	for _, baris := range strings.Split(isi, "\n") {
		if strings.HasPrefix(strings.TrimSpace(baris), "//") {
			continue
		}
		b.WriteString(baris)
		b.WriteString("\n")
	}
	return b.String()
}

// lewatiFolderPindai - folder yang bukan sumber Go aplikasi.
// Salinan `modul/claimlife/services/pindai_bantu_test.go`.
func lewatiFolderPindai(nama string) bool {
	switch nama {
	case "frontend", "node_modules", "bin", ".git", "dist", "unggahan":
		return true
	}
	return false
}

// berkasGoSelainTest mengumpulkan seluruh berkas .go yang BUKAN test, berkunci
// jalur bergaya-garis-miring dari folder paket ini.
// Salinan `modul/claimlife/repository/batasanpemakaian_test.go`.
func berkasGoSelainTest(t *testing.T) map[string]string {
	t.Helper()
	hasil := map[string]string{}
	err := filepath.Walk(filepath.FromSlash(akarAplikasi), func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// frontend dan hasil bangun tidak memuat kode Go yang relevan.
			switch info.Name() {
			case "frontend", "node_modules", "bin", ".git":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}
		isi, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		hasil[filepath.ToSlash(p)] = string(isi)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hasil) == 0 {
		t.Fatal("nol berkas .go terbaca; pembacanya yang rusak, bukan kodenya")
	}
	return hasil
}

// konstantaTeksDiSumber membaca nilai konstanta/variabel teks bernama dari
// sumber Go bukan-uji sebuah folder paket, TANPA mengimpor paketnya - `inti`
// tidak mengimpor modul (`impor_lintas_modul_test.go`). Nama yang tidak
// ditemukan menggagalkan uji pemanggilnya: daftar yang basi harus berbunyi.
// ⛔ Dipanggil DARI uji, bukan dari inisialisasi paket - kegagalan di
// inisialisasi menjatuhkan SELURUH biner uji penjaga (temuan /code-review).
func konstantaTeksDiSumber(t *testing.T, dir string, nama ...string) []string {
	t.Helper()
	folder := filepath.Join(akarAplikasi, filepath.FromSlash(dir))
	isi, err := os.ReadDir(folder)
	if err != nil {
		t.Fatal(err)
	}
	nilai := map[string]string{}
	fset := token.NewFileSet()
	for _, e := range isi {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(folder, e.Name()), nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range f.Decls {
			g, ok := d.(*ast.GenDecl)
			if !ok || (g.Tok != token.CONST && g.Tok != token.VAR) {
				continue
			}
			for _, s := range g.Specs {
				v := s.(*ast.ValueSpec)
				for i, n := range v.Names {
					if i >= len(v.Values) {
						continue
					}
					if lit, ok := v.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
						if teks, err := strconv.Unquote(lit.Value); err == nil {
							nilai[n.Name] = teks
						}
					}
				}
			}
		}
	}
	hasil := make([]string, 0, len(nama))
	for _, n := range nama {
		teks, ada := nilai[n]
		if !ada {
			t.Fatalf("konstanta teks %s tidak ada di %s", n, dir)
		}
		hasil = append(hasil, teks)
	}
	return hasil
}

// polaURL mencocokkan alamat layanan yang ter-hardcode.
//
// ⚠️ `localhost` dan `127.0.0.1` TIDAK dikecualikan. Alamat pengembangan yang
// ter-hardcode adalah persis bentuk percabangan dev-prod yang ADR-U-0013
// larang, dan ia yang paling sering lolos karena terasa tidak berbahaya.
var polaURL = regexp.MustCompile(`://`)

// polaDSNOracle adalah SATU-SATUNYA skema `://` yang bukan alamat layanan.
//
// ⚠️ ADR-U-0013 mengatur ALAMAT LAYANAN KELUAR - endpoint bisnis yang
// di-resolve lewat `M_LINK_SERVICE`. DSN basis data bukan itu: ia diatur
// aturan `ORACLE_DSN` tersendiri (tidak pernah menunjuk produksi). Membiarkan
// penjaga ini menuduhnya berarti penjaga yang menuduh hal yang bukan
// urusannya - dan penjaga seperti itu akhirnya dimatikan orang.
var polaDSNOracle = regexp.MustCompile(`oracle://`)

// kunciEnvBolehBeralamat adalah kunci env yang boleh memuat `://`, beserta
// alasannya. Jumlahnya dikunci di bawah.
var kunciEnvBolehBeralamat = map[string]string{
	"ORACLE_DSN":       "DSN basis data, bukan alamat layanan keluar; diatur aturan ORACLE_DSN",
	"DEV_PROXY_TARGET": "sasaran proxy Vite ke backend sendiri saat pengembangan; tidak pernah ada di produksi",
}

// polaEnvApaPun mencocokkan SETIAP pembacaan env var di luar `inti/config`.
//
// ⛔ Ronde pertama hanya mencocokkan nama yang memuat URL/ENDPOINT/HOST - dan
// dielakkan oleh `os.Getenv("ARASAPAS_TUJUAN")`, yang tidak memuat satu pun
// kata itu. Menambah kata ke daftar tidak menolong: nama env var tak terbatas.
//
// Yang menggantikannya aturan ARSITEKTUR, bukan tebakan nama: env dibaca di
// SATU tempat, `inti/config`, dan tidak di mana pun lagi. Aturan itu dapat
// ditegakkan utuh, sedangkan daftar kata tidak pernah bisa.
var polaEnvApaPun = regexp.MustCompile(`os\.(Getenv|LookupEnv)\(`)

// polaPayloadKeluar mencocokkan pengiriman payload JSON ke luar.
// ⛔ Dari pemeriksaan NAMA menjadi pemeriksaan PERILAKU. Ronde pertama
// mencocokkan tiga pengenal harfiah - dan `json.Marshal` diikuti sebuah POST
// lolos utuh, sebab namanya tidak termasuk ketiganya. Nama dapat diganti;
// yang tidak dapat disembunyikan adalah KLIEN HTTP-nya.
//
// Hari ini modul ini punya NOL klien HTTP keluar, dan itulah yang dikunci:
// hilir membaca langsung dari tabel klaim, tidak ada payload yang dikirim.
// Ketika Arasapas kelak benar-benar dihubungkan (menuntut persetujuan
// manusia), penjaga ini akan merah - dan itu memang yang diinginkan: ia
// memaksa keputusan itu dilihat, bukan menyelinap.
var polaKlienHTTP = regexp.MustCompile(
	`http\.(Post|Get|Head|NewRequest|DefaultClient)|http\.Client\{|\.Do\(req`)

// berkasAlamatDikecualikan berkunci JALUR PENUH relatif akar modul.
//
// ⚠️ Daftar ini pendek dan tiap barisnya beralasan. Jumlahnya dikunci di
// bawah: pengecualian yang tumbuh tanpa alasan adalah cara penjaga mati
// diam-diam.
var berkasAlamatDikecualikan = map[string]string{
	"inti/penjaga/lintasaplikasi_test.go": "berkas penjaga ini sendiri - polanya harus tertulis",
}

// berkasKlienHTTPDisetujui dikecualikan dari pemeriksaan KLIEN HTTP SAJA -
// `://` dan env tetap diperiksa. Tiap baris menyebut persetujuan manusianya;
// jumlahnya dikunci di bawah.
var berkasKlienHTTPDisetujui = map[string]string{
	"modul/treatycontractout/services/tco_pengirim_storage.go": "transport penyimpanan lampiran Treaty Contract Out - " +
		"[keputusan work owner 29-09-2026, OQ-TCO-08]; alamat dari M_LINK_SERVICE saat jalan, " +
		"hanya aktif bila PELAKSANA_STORAGE=nyata",
}

// bolehBacaEnv menyatakan sebuah berkas berhak membaca env var.
//
// ⚠️ Hanya `inti/config` dan berkas test. Daftar ini aturan arsitektur,
// bukan pengecualian per berkas: ia tidak tumbuh seiring berkas bertambah.
func bolehBacaEnv(jalur string) bool {
	rel := filepath.ToSlash(jalur)
	if strings.Contains(rel, "inti/config/") || strings.HasSuffix(rel, "_test.go") {
		return true
	}
	// ⚠️ SATU pengecualian bernama, dan ia mendahului tiket 12: pagar
	// penghancuran skema uji membaca `ORACLE_SKEMA_UJI` SENDIRI, tidak lewat
	// Config. Itu disengaja - pagar yang hanya berlaku bila seseorang
	// menyusun Config dengan benar bukan pagar. Ia membaca pengakuan
	// keselamatan, bukan alamat layanan.
	return strings.HasSuffix(rel, "uji/skemauji/skemauji.go")
}

// TestNolAlamatLayananDiKode menegakkan ADR-U-0013.
func TestNolAlamatLayananDiKode(t *testing.T) {
	diperiksa, dikecualikan, klienDisetujui := 0, 0, 0
	akar := filepath.Join("..", "..")
	err := filepath.Walk(akar, func(jalur string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// node_modules dan dist bukan sumber kita.
			if n := info.Name(); n == "node_modules" || n == "dist" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(jalur, ".go") {
			return nil
		}
		diperiksa++
		// Jalur relatif akar APP_RNM - `internal/...` maupun `inti/...`.
		rel := filepath.ToSlash(jalur)
		for strings.HasPrefix(rel, "../") {
			rel = strings.TrimPrefix(rel, "../")
		}
		if alasan, boleh := berkasAlamatDikecualikan[rel]; boleh {
			t.Logf("dikecualikan: %s (%s)", rel, alasan)
			dikecualikan++
			return nil
		}
		isi, err := os.ReadFile(jalur)
		if err != nil {
			return err
		}
		kode := buangKomentar(string(isi))
		if polaURL.MatchString(polaDSNOracle.ReplaceAllString(kode, "")) {
			t.Errorf("%s: memuat `://` di luar komentar; ADR-U-0013 menuntut "+
				"alamat di-resolve lewat M_LINK_SERVICE saat jalan. Yang dicari "+
				"`://`, bukan `http://` - ronde pertama dielakkan oleh "+
				"penyambungan (`skema + \"://\" + inang`) dan oleh "+
				"`Sprintf(\"%%s://%%s\", …)`", filepath.ToSlash(jalur))
		}
		if m := polaEnvApaPun.FindString(kode); m != "" && !bolehBacaEnv(jalur) {
			t.Errorf("%s: membaca env var (%q) di luar inti/config; env "+
				"dibaca di SATU tempat, dan alamat layanan tidak pernah datang "+
				"dari env sama sekali (ADR-U-0013)", filepath.ToSlash(jalur), m)
		}
		// ⚠️ Berkas test dikecualikan DARI PEMERIKSAAN INI SAJA: test HTTP
		// memanggil server KITA SENDIRI lewat `httptest`, bukan layanan bisnis
		// luar. Batasnya dinyatakan: sebuah panggilan keluar yang disembunyikan
		// di dalam berkas test tidak akan tertangkap - tetapi ia juga tidak
		// pernah berjalan di produksi.
		m := polaKlienHTTP.FindString(kode)
		if alasan, boleh := berkasKlienHTTPDisetujui[rel]; boleh && m != "" {
			t.Logf("klien HTTP disetujui: %s (%s)", rel, alasan)
			klienDisetujui++
			m = ""
		}
		if m != "" && !strings.HasSuffix(filepath.ToSlash(jalur), "_test.go") {
			t.Errorf("%s: memuat klien HTTP keluar (%q); konversi ke produksi "+
				"TIDAK lagi mengirim payload - hilir membaca langsung dari tabel "+
				"klaim (`[keputusan work owner 2026-09-16]`), dan menghubungkan "+
				"layanan luar menuntut persetujuan manusia",
				filepath.ToSlash(jalur), m)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// ⛔ Berkas `.env*` IKUT dibaca. Ronde pertama hanya membaca `.go`, dan
	// `APP_RNM/.env`, `.env.example`, serta `frontend/.env` tidak pernah
	// tersentuh - padahal memindahkan URL ke sana persis elakan yang paling
	// wajar dilakukan orang.
	envDiperiksa, envDikecualikan := 0, 0
	for _, pola := range []string{".env", ".env.*", "frontend/.env", "frontend/.env.*"} {
		cocok, _ := filepath.Glob(filepath.Join(akar, pola))
		for _, berkas := range cocok {
			isi, err := os.ReadFile(berkas)
			if err != nil {
				continue
			}
			envDiperiksa++
			for _, baris := range strings.Split(string(isi), "\n") {
				if !polaURL.MatchString(baris) {
					continue
				}
				kunci := strings.TrimSpace(strings.SplitN(baris, "=", 2)[0])
				if _, boleh := kunciEnvBolehBeralamat[kunci]; boleh {
					envDikecualikan++
					continue
				}
				// ⛔ Nilainya TIDAK dicetak. Berkas env memuat kredensial
				// pengembangan, dan keluaran test terbaca siapa pun yang
				// menjalankannya; yang perlu diketahui hanya KUNCI mana.
				t.Errorf("%s: kunci %q memuat `://`; alamat layanan keluar "+
					"di-resolve lewat M_LINK_SERVICE, tidak pernah dari berkas "+
					"env (ADR-U-0013)", filepath.ToSlash(berkas), kunci)
			}
		}
	}
	t.Logf("berkas .go dibaca: %d; berkas .env dibaca: %d; kunci env dikecualikan: %d",
		diperiksa, envDiperiksa, envDikecualikan)
	// ⛔ Pengecualian env dikunci: tiap kunci di petanya harus benar-benar
	// terpakai. Pengecualian yang tidak terpakai adalah pengecualian yang
	// tidak ada yang periksa lagi.
	if envDikecualikan != len(kunciEnvBolehBeralamat) {
		t.Errorf("%d kunci env dikecualikan, sedangkan petanya memuat %d",
			envDikecualikan, len(kunciEnvBolehBeralamat))
	}
	if diperiksa < 40 {
		t.Fatalf("hanya %d berkas terbaca; pembacanya yang rusak, bukan kodenya",
			diperiksa)
	}
	if dikecualikan != len(berkasAlamatDikecualikan) {
		t.Errorf("%d berkas dikecualikan, sedangkan petanya memuat %d",
			dikecualikan, len(berkasAlamatDikecualikan))
	}
	// ⛔ Pengecualian klien HTTP juga dikunci: baris yang tidak lagi memuat
	// klien HTTP harus dihapus dari petanya.
	if klienDisetujui != len(berkasKlienHTTPDisetujui) {
		t.Errorf("%d berkas klien HTTP disetujui terpakai, sedangkan petanya memuat %d",
			klienDisetujui, len(berkasKlienHTTPDisetujui))
	}
}

// polaNamaOrangTetap mencocokkan pemberian TEKS TETAP ke medan yang namanya
// menandakan nama orang - termasuk teks berkutip-balik.
//
// ⛔ STRUKTURAL, bukan berdaftar-nama. Percobaan pertama memuat daftar nama
// operator nyata dari korpus, dan itu sendiri melanggar pagar keamanan brief
// induk ("nilai nama orang tidak disalin ke artefak mana pun").
//
// ⚠️ Ia HEURISTIK, dan batasnya dinyatakan di tiket: nama yang masuk lewat
// konstanta perantara, lewat medan yang tidak terdaftar, atau lewat
// perbandingan `AkunID == "..."` tidak tertangkap. Penjaga yang menyebut
// batasnya lebih berguna daripada penjaga yang mengaku sempurna.
var polaNamaOrangTetap = regexp.MustCompile(
	"(?i)(OpName|NamaOrang|PolicyHolder|NameOfInsured|Tertanggung)" +
		"\\s*[:=]+\\s*[\"`][^\"`]+[\"`]")

// pesanVerbatimYangSah adalah NILAI yang cocok dengan pola di atas tetapi
// BUKAN nama orang - masing-masing beserta alasannya.
//
// ⛔ DIPERSEMPIT 28-09-2026 dengan daftar bernama, BUKAN dengan melonggarkan
// polanya. Penjaga yang menuduh hal yang benar akan dilonggarkan orang, bukan
// dipatuhi - pelajaran yang sudah dibayar dua kali di repo ini (nama tabel
// telanjang, ambang tutup buku).
//
// Yang menuduh di sini: konstanta PESAN VALIDASI yang disalin VERBATIM dari
// `ValidasiUploadPL_act.xml`. Namanya memuat "PolicyHolder"/"Tertanggung"
// sebab itulah KOLOM yang divalidasi, dan nilainya kalimat galat huruf besar
// - bukan nama siapa pun.
//
// ⛔ KUNCINYA DIAMBIL DARI `models`, TIDAK DIKETIK ULANG. Refactor bentuk B
// paket 8: penjaga ini kini di `inti/penjaga`, dan `inti` tidak mengimpor
// modul - nilainya karena itu DIBACA dari sumber konstantanya
// (`konstantaTeksDiSumber`), tetap tanpa diketik ulang. Menuliskan
// kalimatnya harfiah di sini akan membuat penjaga ini menuduh DIRINYA
// SENDIRI - dan itu persis yang terjadi pada ronde pertama penyempitan ini.
// Mengambilnya dari konstantanya juga berarti daftar ini ikut basi begitu
// pesannya berubah, alih-alih diam-diam tetap mengecualikan teks lama.
// ⛔ ALASANNYA DI KOMENTAR, BUKAN DI NILAI. Menuliskannya sebagai nilai teks
// membuat BARIS DAFTAR INI SENDIRI cocok dengan polanya - `...Tertanggung:
// "alasan"` - dan penjaga ini menuduh dirinya sendiri. Sudah terjadi, dua
// kali, saat penyempitan ini ditulis.
func pesanVerbatimYangSah(t *testing.T) []string {
	return konstantaTeksDiSumber(t, "modul/premiumlistlife/models",
		// ValidasiUploadPL_act `local.err3` - pesan kolom NAME_OF_INSURED.
		"PesanNamaTertanggung",
		// ValidasiUploadPL_act `local.err17` - pesan rujukan master POLICY HOLDER.
		"PesanPolicyHolder",
	)
}

// pesanVerbatimDiterima menjawab apakah sebuah nilai ada di daftar itu.
func pesanVerbatimDiterima(sah []string, nilai string) bool {
	for _, p := range sah {
		if p == nilai {
			return true
		}
	}
	return false
}

// polaNilaiTerkutip mengambil nilai di dalam tanda kutip sebuah baris cocok.
var polaNilaiTerkutip = regexp.MustCompile("[\"`]([^\"`]+)[\"`]")

func TestNolNamaOrangDiKode(t *testing.T) {
	sah := pesanVerbatimYangSah(t)
	diperiksa := 0
	err := filepath.Walk(akarAplikasi, func(jalur string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() && lewatiFolderPindai(info.Name()) {
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(jalur, ".go") {
			return nil
		}
		isi, err := os.ReadFile(jalur)
		if err != nil {
			return err
		}
		diperiksa++
		for _, baris := range strings.Split(string(isi), "\n") {
			if strings.HasPrefix(strings.TrimSpace(baris), "//") {
				continue
			}
			cocok := polaNamaOrangTetap.FindString(baris)
			if cocok == "" {
				continue
			}
			// Nilai sintetis BERAWALAN UJI- memang bentuk yang pagar keamanan
			// brief tuntut untuk fixture. Diperiksa di AWAL teksnya, bukan di
			// mana saja - "Budi UJI-1" bukan nilai sintetis.
			if awalanUjiSintetis(cocok) {
				continue
			}
			// Pesan validasi VERBATIM - didaftar satu per satu beserta
			// alasannya, lihat `pesanVerbatimYangSah`. Yang dicocokkan
			// NILAINYA, bukan seluruh barisnya: baris yang sama dapat ditulis
			// dengan spasi berbeda, dan daftar yang mencocokkan spasi akan
			// lolos begitu seseorang menjalankan gofmt.
			if m := polaNilaiTerkutip.FindStringSubmatch(cocok); m != nil &&
				pesanVerbatimDiterima(sah, m[1]) {
				continue
			}
			t.Errorf("%s memberi nilai tetap ke medan bernama-orang: %s",
				filepath.ToSlash(jalur), strings.TrimSpace(baris))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if diperiksa < 10 {
		t.Fatalf("hanya %d berkas terbaca; pembacanya yang rusak", diperiksa)
	}
}

// awalanUjiSintetis menyatakan teks yang dikutip dimulai dengan UJI-.
func awalanUjiSintetis(cocok string) bool {
	for _, kutip := range []string{"\"", "`"} {
		i := strings.Index(cocok, kutip)
		if i < 0 {
			continue
		}
		return strings.HasPrefix(cocok[i+1:], "UJI-")
	}
	return false
}

// Lapisan handlers tidak boleh menyentuh repository langsung.
//
// Arah ketergantungan proyek ini satu arah: handlers -> services -> repository.
// Memotongnya membuat aturan dagang tersebar ke lapisan yang tugasnya hanya
// menerima permintaan HTTP.
func TestHandlersTidakMengimporRepository(t *testing.T) {
	diperiksa := 0
	// Refactor bentuk B (30-09-2026): handlers SETIAP modul, dan repository
	// SETIAP modul. `inti/db` ikut dilarang: isinya dulu kepala paket
	// repository (koneksi dan transaksi), jadi larangan lama tetap utuh.
	polaRepository := regexp.MustCompile(`"nusantarare/(internal/repository|modul/[^/"]+/repository|inti/db)"`)
	for nama, isi := range berkasGoSelainTest(t) {
		if !strings.Contains(nama, "/handlers/") {
			continue
		}
		diperiksa++
		if polaRepository.MatchString(isi) {
			t.Errorf("%s mengimpor repository; seharusnya lewat services", nama)
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol berkas handlers terbaca; pembacanya yang rusak")
	}
}

// ADR-U-0033: nol nama tabel telanjang di dalam teks query.
//
// "Telanjang" berarti tanpa awalan skema. Query begitu benar hanya selama sesi
// kebetulan menunjuk skema yang tepat, dan salahnya baru muncul saat pindah
// lingkungan - jauh dari orang yang menulisnya. ADR-U-0033 Akibat 3 menuntut
// test yang menemukannya gagal; sampai ronde 3 test itu tidak pernah ada, dan
// SYS.ALL_OBJECTS sempat lolos sebagai ALL_OBJECTS telanjang.
//
// Yang dianggap SAH sesudah FROM / INTO / UPDATE / JOIN:
//   - "%s"            nama yang sudah dilewatkan Qualify
//   - "A.B"           sudah berawalan skema, termasuk SYS.
//   - "{skema}.B"     penanda di berkas migrasi
//   - "DUAL"          tabel semu milik Oracle, tidak punya skema
//
// Dan `FOR UPDATE` DILEWATI: ia klausa penguncian baris, bukan pernyataan
// UPDATE, sehingga kata sesudahnya (`SKIP`, `NOWAIT`, atau tidak ada) bukan
// nama tabel. Penyempitan ini dibuktikan MASIH MENGGIGIT sebelum dipakai.
// ⛔ KOMENTAR DIBUANG SEBELUM PENCOCOKAN, sejak 27-09-2026. Tanpa itu
// prosa yang MENERANGKAN sebuah query - hal biasa di repositori ini -
// dituduh sebagai query-nya. Yang menyalakannya satu kalimat di
// `diagnosa.go`: "menirunya dengan N UPDATE berarti daftar berlubang", dan
// penjaga membaca `UPDATE berarti` sebagai `UPDATE <nama tabel>`.
//
// ⚠️ Penjaga yang menuduh hal yang BENAR akan dilonggarkan orang,
// bukan dipatuhi. Jadi ia dipersempit sekarang - dan dibuktikan MASIH

// polaTabelTelanjang mencari kata sesudah FROM / INTO / UPDATE / JOIN.
var polaTabelTelanjang = regexp.MustCompile(
	`(?i)(\bFOR\s+)?\b(FROM|INTO|UPDATE|JOIN)\s+([A-Za-z_{%][\w{}%.]*)`)

func TestNolNamaTabelTelanjangDiQuery(t *testing.T) {
	diperiksa := 0
	for nama, isi := range berkasGoSelainTest(t) {
		// Refactor bentuk B (30-09-2026): SQL kini juga tinggal di
		// `modul/*/repository` dan `inti/*`. Dulu hanya `/internal/repository/`
		// - saringan itu diam-diam menyempit begitu kode pindah (cacah log
		// dasar 275 rujukan; sesudah paket 1-2 tanpa perbaikan ini, 168).
		// Paket 5: skema uji pindah dari internal/repository/skemauji ke
		// uji/skemauji - SQL tiruannya tetap dalam cakupan.
		if !strings.Contains(nama, "/repository/") && !strings.Contains(nama, "/inti/") &&
			!strings.Contains(nama, "/uji/skemauji/") {
			continue
		}
		// Baris komentar dibuang lebih dulu (paket 8: `buangKomentar` berkas
		// ini, pengganti `polaKomentarBaris` yang setara).
		isi = buangKomentar(isi)
		for _, m := range polaTabelTelanjang.FindAllStringSubmatch(isi, -1) {
			if m[1] != "" { // klausa `FOR UPDATE`, bukan pernyataan
				continue
			}
			objek := m[3]
			diperiksa++
			switch {
			case strings.Contains(objek, "."): // berawalan skema atau {skema}
			case strings.Contains(objek, "%s"): // datang dari Qualify
			case strings.EqualFold(objek, "DUAL"): // tabel semu Oracle
			default:
				t.Errorf("%s: %s %s - nama tabel telanjang (ADR-U-0033)",
					nama, strings.ToUpper(m[2]), objek)
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol rujukan tabel terbaca; pembacanya yang rusak")
	}
	t.Logf("%d rujukan tabel diperiksa", diperiksa)
}

// tco4 (keputusan work owner 29-09-2026): NOL tabel baru untuk modul ini -
// rentang migrasinya 300-319 kosong, dan tidak satu pun berkas migrasi mana
// pun membuat tabel atau sequence bernama Treaty Contract Out.
func TestTCONolTabelBaru(t *testing.T) {
	pola := regexp.MustCompile(`(?i)CREATE\s+(TABLE|SEQUENCE)\s+\{skema\}\.(\w+)`)
	nama := regexp.MustCompile(`(?i)^(SEQ_)?(T_)?(M?TREATY|PROPORTIONAL)`)
	berkas := 0
	for n, teks := range seluruhSQL(t, false) {
		berkas++
		if n >= "300_" && n < "320_" {
			t.Errorf("%s: berkas migrasi di rentang Treaty Contract Out 300-319 - tco4 menolak tabel baru", n)
		}
		for _, m := range pola.FindAllStringSubmatch(teks, -1) {
			if nama.MatchString(m[2]) {
				t.Errorf("%s membuat %s %s - tco4: modul ini memakai tabel warisan", n, m[1], m[2])
			}
		}
	}
	if berkas == 0 {
		t.Fatal("nol berkas migrasi terbaca; pembacanya yang rusak")
	}
}

// tco4: nol nama tabel baru modul (T_ + TREATY…/MTREATY…/PROPORTIONAL…) di KODE
// Go dan frontend. Komentar - catatan sejarah dan ralat - dibuang lebih dulu.
// Polanya dirakit dari potongan supaya berkas ini sendiri tidak memuatnya.
func TestTCONolNamaTabelBaruDiKode(t *testing.T) {
	pola := regexp.MustCompile("T" + "_" + `(TREATY|MTREATY|PROPORTIONAL)\w*`) // tanpa \b: SEQ_ + nama ikut
	if !pola.MatchString("SELECT ID FROM S."+"T"+"_TREATYYEAR") || !pola.MatchString("S.SEQ_"+"T"+"_TREATYYEAR") ||
		pola.MatchString("SELECT ID FROM S.TREATYYEAR") {
		t.Fatal("pola penjaga tidak menggigit atau menuduh nama warisan")
	}
	ekor := regexp.MustCompile(`(^|\s)//.*$`)
	blok := regexp.MustCompile(`(?s)/\*.*?\*/`)
	berkas := 0
	for _, akar := range []string{akarAplikasi + "/inti", akarAplikasi + "/modul", akarAplikasi + "/cmd", akarAplikasi + "/frontend/src"} {
		err := filepath.Walk(filepath.FromSlash(akar), func(jalur string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if info.Name() == "node_modules" || info.Name() == "dist" {
					return filepath.SkipDir
				}
				return nil
			}
			if ext := filepath.Ext(jalur); ext != ".go" && ext != ".ts" && ext != ".tsx" {
				return nil
			}
			isi, err := os.ReadFile(jalur)
			if err != nil {
				return err
			}
			berkas++
			kode := blok.ReplaceAllString(string(isi), "")
			for i, baris := range strings.Split(kode, "\n") {
				baris = ekor.ReplaceAllString(baris, "")
				if filepath.Ext(jalur) != ".go" && strings.HasPrefix(strings.TrimSpace(baris), "*") {
					continue // badan JSDoc
				}
				if m := pola.FindString(baris); m != "" {
					t.Errorf("%s:%d menyebut %s - tco4: modul ini memakai tabel warisan", filepath.ToSlash(jalur), i+1, m)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if berkas < 100 {
		t.Fatalf("hanya %d berkas terbaca; pembacanya yang rusak", berkas)
	}
}
