package services_test

// Penjaga statik tiket 12 - alamat layanan dan payload keluar.
//
// Pemilik: tiket 12. Dibaca sesudah: efekkeluar.go.
//
// Dua aturan dijaga, keduanya diminta tiket dalam bentuk test:
//
//  1. ADR-U-0013: nol URL sebagai literal, konstanta, MAUPUN env var. Yang
//     boleh jadi konstanta hanyalah kunci kategori. Pemisahan dev-prod
//     terjadi lewat isi tabel per-database, bukan percabangan di kode.
//  2. `[keputusan work owner 2026-09-16]`: konversi ke produksi tidak lagi
//     mengirim payload JSON; hilir membaca langsung dari tabel klaim.
//
// ⛔ Ketiga kelemahan penjaga tiket 10 ditutup SEJAK AWAL di sini:
// pengecualian berkunci jalur penuh (bukan nama berkas), akar telusurnya akar
// modul (bukan `internal/`), dan pencacahnya naik SEBELUM pengecualian
// sehingga swa-periksa tidak dapat dikelabui pengecualian yang terlalu lebar.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

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

// polaEnvApaPun mencocokkan SETIAP pembacaan env var di luar `internal/config`.
//
// ⛔ Ronde pertama hanya mencocokkan nama yang memuat URL/ENDPOINT/HOST - dan
// dielakkan oleh `os.Getenv("ARASAPAS_TUJUAN")`, yang tidak memuat satu pun
// kata itu. Menambah kata ke daftar tidak menolong: nama env var tak terbatas.
//
// Yang menggantikannya aturan ARSITEKTUR, bukan tebakan nama: env dibaca di
// SATU tempat, `internal/config`, dan tidak di mana pun lagi. Aturan itu dapat
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
	"internal/services/efekkeluar_statik_test.go": "berkas penjaga ini sendiri - polanya harus tertulis",
}

// bolehBacaEnv menyatakan sebuah berkas berhak membaca env var.
//
// ⚠️ Hanya `internal/config` dan berkas test. Daftar ini aturan arsitektur,
// bukan pengecualian per berkas: ia tidak tumbuh seiring berkas bertambah.
func bolehBacaEnv(jalur string) bool {
	rel := filepath.ToSlash(jalur)
	if strings.Contains(rel, "internal/config/") || strings.HasSuffix(rel, "_test.go") {
		return true
	}
	// ⚠️ SATU pengecualian bernama, dan ia mendahului tiket 12: pagar
	// penghancuran skema uji membaca `ORACLE_SKEMA_UJI` SENDIRI, tidak lewat
	// Config. Itu disengaja - pagar yang hanya berlaku bila seseorang
	// menyusun Config dengan benar bukan pagar. Ia membaca pengakuan
	// keselamatan, bukan alamat layanan.
	return strings.HasSuffix(rel, "internal/repository/skemauji/skemauji.go")
}

// TestNolAlamatLayananDiKode menegakkan ADR-U-0013.
func TestNolAlamatLayananDiKode(t *testing.T) {
	diperiksa, dikecualikan := 0, 0
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
		rel := filepath.ToSlash(jalur)
		if i := strings.Index(rel, "internal/"); i >= 0 {
			rel = rel[i:]
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
			t.Errorf("%s: membaca env var (%q) di luar internal/config; env "+
				"dibaca di SATU tempat, dan alamat layanan tidak pernah datang "+
				"dari env sama sekali (ADR-U-0013)", filepath.ToSlash(jalur), m)
		}
		// ⚠️ Berkas test dikecualikan DARI PEMERIKSAAN INI SAJA: test HTTP
		// memanggil server KITA SENDIRI lewat `httptest`, bukan layanan bisnis
		// luar. Batasnya dinyatakan: sebuah panggilan keluar yang disembunyikan
		// di dalam berkas test tidak akan tertangkap - tetapi ia juga tidak
		// pernah berjalan di produksi.
		if m := polaKlienHTTP.FindString(kode); m != "" &&
			!strings.HasSuffix(filepath.ToSlash(jalur), "_test.go") {
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
}

// TestFlagLingkunganTidakMenggerbangiPenyimpanan - penyimpangan sadar.
//
// ⛔ Flag lingkungan hanya menggerbangi EFEK KELUAR. Di Pega `IsPEGAPROD`
// menggerbangi simpan utama juga; bila kita menirunya, lingkungan
// non-produksi tidak dapat dipakai menguji sama sekali - klaim tidak akan
// pernah tersimpan di sana.
//
// ⛔ PENJAGA TEKSTUAL DICABUT - ia terbukti dapat dielakkan.
//
// Ronde pertama memeriksa teks sumber: `polaCabangLingkungan` menuntut
// `Produksi()` BERSEBELAHAN dengan `if`. Dua elakan dibangun, dikompilasi,
// dan dijalankan - keduanya hijau:
//
//	prod := p.penyalur.lingkungan.AdalahProduksi()
//	if !prod { return nil }              // menyimpan dilewati di non-produksi
//
//	var _ = "if !p.lingkungan.AdalahProduksi()" // memuaskan strings.Contains
//
// Elakan kedua paling telak: ia memuaskan pemeriksaan "gerbangnya masih ada"
// dengan sebuah LITERAL TEKS, sementara gerbang aslinya dihapus - efek keluar
// menyala di setiap lingkungan, email nyata dari lingkungan uji.
//
// Pelajarannya bukan "perkuat polanya". Pola atas teks sumber selalu dapat
// dielakkan oleh penulisan ulang yang setara. Yang menggantikannya adalah
// penjaga PERILAKU di `efekkeluar_test.go` dan di bawah: ia memanggil kodenya
// dan memeriksa apa yang TERJADI, sehingga penulisan ulang apa pun yang
// mengubah perilaku akan tertangkap apa pun bentuknya.
func TestFlagLingkunganTidakMenggerbangiPenyimpanan(t *testing.T) {
	// ⛔ PERILAKU, bukan teks. Layanan disusun di lingkungan BUKAN produksi -
	// bawaan `svc.Komite()` - lalu jalur penyimpanannya dipanggil TANPA Oracle.
	//
	// Yang benar: ia berjalan sampai menyentuh basis data, lalu berhenti di
	// `ErrTanpaOracle`. Bila seseorang menyisipkan gerbang lingkungan di depan
	// penyimpanan - dalam bentuk APA PUN, termasuk yang dinaikkan ke variabel
	// atau ditulis sebagai `switch` - jalur itu akan pulang lebih awal dan
	// mengembalikan nil, dan test ini merah. Penulisan ulang yang setara tidak
	// menolongnya: yang diperiksa akibatnya, bukan bentuknya.
	svc := services.New(nil)
	err := svc.Komite().Serahkan(context.Background(),
		services.Pelaku{AkunID: "UJI-AKUN", Peran: []string{services.PeranAdmin}},
		"CLM-1", "P-1", "A-1", saatUji)
	if err == nil {
		t.Fatal("penyimpanan pulang tanpa galat di lingkungan bukan produksi; " +
			"flag lingkungan menggerbangi SIMPAN, bukan hanya efek keluar - " +
			"tanpa Oracle jalur ini seharusnya berhenti di ErrTanpaOracle")
	}
	if !errors.Is(err, repository.ErrTanpaOracle) {
		t.Errorf("galat = %v, mau ErrTanpaOracle; jalur penyimpanan tidak "+
			"sampai menyentuh basis data", err)
	}
}

// TestGerbangLingkunganBenarBenarMenggerbangi - pasangan test di atas.
//
// ⛔ Test di atas memastikan gerbangnya TIDAK ada di penyimpanan. Test ini
// memastikan ia ADA di efek keluar. Tanpa keduanya, menghapus gerbang itu
// seluruhnya akan tetap hijau - dan itulah elakan ketiga yang terbukti
// berhasil terhadap ronde pertama.
func TestGerbangLingkunganBenarBenarMenggerbangi(t *testing.T) {
	efek := &efekUji{nama: "UJI-GERBANG"}
	antre := &antreanUji{}

	services.NewPenyalur(services.BukanProduksi, antre, efek).
		Salurkan(context.Background(), services.MuatanEfek{KlaimID: "CLM-1"})
	if efek.dipanggil != 0 {
		t.Fatalf("efek berjalan %d kali di BUKAN produksi; gerbangnya hilang - "+
			"email nyata akan terkirim dari lingkungan uji", efek.dipanggil)
	}

	services.NewPenyalur(services.Produksi, antre, efek).
		Salurkan(context.Background(), services.MuatanEfek{KlaimID: "CLM-1"})
	if efek.dipanggil != 1 {
		t.Errorf("efek berjalan %d kali di produksi, mau 1; gerbangnya menutup "+
			"terlalu rapat dan efek keluar tidak pernah berjalan di mana pun",
			efek.dipanggil)
	}
}
