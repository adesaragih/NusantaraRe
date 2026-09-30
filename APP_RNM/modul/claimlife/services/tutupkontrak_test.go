package services

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Penjaga kontrak butir bb: SETIAP rute pengubah menolak kasus yang tertutup.
//
// ⛔ Kenapa penjaga statik dan bukan tujuh uji perilaku. Aturan yang ditulis
// ulang di tujuh berkas adalah aturan yang suatu hari hanya ada di enam, dan
// yang ketujuh tidak akan berbunyi: tiap berkas hijau sendirian. Itu bentuk
// cacat yang sudah terjadi tiga kali di modul ini (envelope `galat`, rute
// tanpa pemanggil, penanda `IsCheck`). Penjaga ini menagih pemanggilannya
// dari daftar, sehingga layanan pengubah KEDELAPAN pun akan tertagih - asal
// namanya ditambahkan ke sini, dan menambahkannya adalah pekerjaan sadar.

// layananPengubah memetakan berkas layanan ke fungsi masuknya - SETIAP
// fungsi yang melayani rute pengubah, bukan satu per berkas.
//
// ⛔ DIPERKETAT 28-09-2026 (temuan /code-review sensus §3.1). Ronde
// sebelumnya memeriksa panggilan di BERKAS, bukan di fungsinya: ketika
// `dol.go` memperoleh fungsi masuk kedua (`SetTanggalKlaim`), menghapus
// pemeriksaan kasus tertutup dari fungsi itu tetap hijau - panggilan milik
// `Set` di berkas yang sama membungkam penjaga. Kini tubuh tiap fungsi yang
// disebut diperiksa sendiri-sendiri.
//
// ⚠️ Daftar ini DITULIS TANGAN dan itu disengaja. Menurunkannya otomatis dari
// "berkas yang memanggil DalamTransaksi" akan membuat penjaga ini menjaga
// apa pun yang kebetulan ada, bukan apa yang kita putuskan harus dijaga -
// dan daftar yang menyesuaikan diri tidak pernah gagal.
var layananPengubah = map[string][]string{
	// Komite Claim Life tiket 02 - keputusan tingkat; klaim induk yang
	// tertutup tidak menerima keputusan komite lagi.
	"komite_keputusan.go": {"Putuskan"},
	"akseptasi.go":        {"SimpanAdjustment"},
	// Dua rute dialog Edit Date: DOL (tiket 06) dan tiga tanggal (§3.1).
	"dol.go":         {"Set", "SetTanggalKlaim"},
	"hapus.go":       {"Hapus"},
	"hasilkomite.go": {"Tambah"},
	"diagnosa.go":    {"pagari"},
	"unggahan.go":    {"pagari"},
	"komite.go":      {"Serahkan"},
	"tahap.go":       {"Pindah"},
	// GILIRAN-11 paket 1 - Save to RNM (SaveOutStandingLife_Act).
	"simpanrnm.go":   {"Simpan"},
	"statusbaris.go": {"ubah"},
	// GILIRAN-17 (OQ-M6) - tombol `DELETE` `InputOSClaimLife` b17865.
	"cabutpeserta.go": {"CabutPeserta"},
}

// tubuhFungsi mengembalikan teks satu fungsi (atau metode) bernama, dari
// kepalanya sampai sebelum `func` berikutnya di awal baris. Kosong bila tidak ada.
func tubuhFungsi(teks, nama string) string {
	pola := regexp.MustCompile(`(?m)^func (\([^)]*\) )?` + regexp.QuoteMeta(nama) + `\(ctx context\.Context`)
	loc := pola.FindStringIndex(teks)
	if loc == nil {
		return ""
	}
	sisa := teks[loc[1]:]
	if j := strings.Index(sisa, "\nfunc "); j >= 0 {
		sisa = sisa[:j]
	}
	return teks[loc[0]:loc[1]] + sisa
}

func TestSetiapLayananPengubahMemeriksaKasusTerbuka(t *testing.T) {
	const panggilan = "PastikanKasusTerbuka(ctx, klaimID)"
	for berkas, daftar := range layananPengubah {
		isi, err := bacaBerkasLayanan(berkas)
		if err != nil {
			t.Errorf("membaca %s: %v", berkas, err)
			continue
		}
		for _, fungsi := range daftar {
			tubuh := tubuhFungsi(string(isi), fungsi)
			if tubuh == "" {
				t.Errorf("%s: fungsi masuk %q tidak ditemukan lagi - daftar "+
					"layananPengubah harus disesuaikan, bukan penjaga ini dimatikan",
					berkas, fungsi)
				continue
			}
			if !strings.Contains(tubuh, panggilan) {
				t.Errorf("%s (%s) tidak memanggil %s.\n"+
					"Butir bb: sesudah kasus ditutup, SETIAP rute pengubah harus "+
					"menolak. Rute yang lupa memeriksanya tidak akan membuat satu "+
					"pun uji lain merah - ia hanya akan mengubah kasus yang sudah "+
					"selesai.", berkas, fungsi, panggilan)
			}
		}
	}
}

// TestDaftarLayananPengubahMencakupSeluruhRutePengubah menagih arah
// sebaliknya: rute pengubah yang BARU harus masuk daftar di atas.
//
// ⛔ Tanpa uji ini, daftar itu hanya menjaga dirinya sendiri. Rute kedelapan
// dapat lahir lengkap dengan uji perilakunya yang hijau, dan penjaga di atas
// akan tetap diam - sebab ia hanya memeriksa nama yang sudah tertulis.
func TestDaftarLayananPengubahMencakupSeluruhRutePengubah(t *testing.T) {
	// Berkas layanan yang MENULIS ke basis data adalah yang memanggil
	// DalamTransaksi. Setiap satu di antaranya harus ada di daftar, atau
	// dinyatakan alasannya di pengecualian bernama di bawah.
	dikecualikan := map[string]string{
		// Pendaftaran MELAHIRKAN kasus; belum ada kasus untuk ditutup.
		"pendaftaran.go": "melahirkan kasus, bukan mengubah kasus yang ada",
		// GILIRAN-13 butir bn: tombol portal Input Offer/Input Premium juga
		// MELAHIRKAN kasus - kali ini kasus polis (T_WORK_POLIS).
		"polis_kasus.go": "melahirkan kasus polis, bukan mengubah kasus yang ada",
		// PremiumList Life tiket 01. Ia MEMERIKSA kasus tertutup - lihat
		// `Penawaran.pagari` - tetapi lewat `T_WORK_POLIS` dan
		// `models.KasusPolisTertutup`, BUKAN lewat `PastikanKasusTerbuka`
		// yang membaca `T_WORK_CLAIM`. Dua modul, dua tabel kerja: layanan
		// polis yang menanyai tabel klaim akan selalu menjawab "tidak ada".
		"polis_penawaran.go": "modul PremiumList - memeriksa T_WORK_POLIS lewat KasusPolisTertutup",
		// PremiumList Life tiket 03. Sebab yang SAMA dengan di atas: ia
		// memeriksa T_WORK_POLIS, bukan T_WORK_CLAIM. Penjaganya disebut
		// namanya di `penjagaTutup` di bawah, jadi ini bukan kelonggaran -
		// ia tetap dituntut penjaga, hanya penjaga yang lain.
		"polis_nomor.go": "modul PremiumList - memeriksa T_WORK_POLIS lewat KasusPolisTertutup",
		// PremiumList Life tiket 04. Sebab yang SAMA: T_WORK_POLIS, bukan
		// T_WORK_CLAIM. Penjaganya disebut namanya di `penjagaTutup`.
		"polis_unggah.go": "modul PremiumList - memeriksa T_WORK_POLIS lewat KasusPolisTertutup",
		// PremiumList Life tiket 05a bagian 2. Sebab yang SAMA: T_WORK_POLIS.
		// Penjaganya disebut namanya di `penjagaTutup`.
		"polis_summary.go": "modul PremiumList - memeriksa T_WORK_POLIS lewat KasusPolisTertutup",
		// Outbox dan pelaksana efek bekerja atas baris antrean, bukan atas
		// kasus - dan efek yang sudah terlanjur diantre tetap harus selesai
		// walau kasusnya kemudian ditutup.
		"antrean.go": "menjalankan antrean efek, bukan mengubah kasus",
		// Tutup itu sendiri: ia yang MENUTUP, dan ia memeriksa dengan
		// pembacaan status kerjanya sendiri.
		"tutup.go": "layanan penutupnya sendiri",
		// Penerbit nomor akseptasi dipanggil DARI akseptasi.go, yang sudah
		// diperiksa di pintunya.
		"spreading.go": "dipanggil dari akseptasi.go yang sudah memeriksa",
		// Tolak MENYALURKAN ke statusbaris.go `ubah`, yang memeriksanya.
		// Penjaga di dua pintu menuju satu ruang adalah dua tempat untuk lupa.
		"tolak.go": "menyalurkan ke statusbaris.go ubah yang sudah memeriksa",
		// services.go adalah tempat DalamTransaksi DIDEFINISIKAN, bukan
		// pemakainya. Ia cocok dengan pencariannya karena namanya sendiri.
		"services.go": "mendefinisikan DalamTransaksi, bukan memakainya",
	}
	masuk, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("membaca direktori layanan: %v", err)
	}
	diperiksa := 0
	for _, e := range masuk {
		nama := e.Name()
		if e.IsDir() || !strings.HasSuffix(nama, ".go") ||
			strings.HasSuffix(nama, "_test.go") {
			continue
		}
		isi, err := os.ReadFile(filepath.Join(".", nama))
		if err != nil {
			t.Fatalf("membaca %s: %v", nama, err)
		}
		if !strings.Contains(string(isi), "DalamTransaksi(ctx") {
			continue
		}
		diperiksa++
		if _, ada := layananPengubah[nama]; ada {
			continue
		}
		if alasan, ada := dikecualikan[nama]; ada {
			if alasan == "" {
				t.Errorf("%s dikecualikan tanpa alasan tertulis", nama)
			}
			continue
		}
		t.Errorf("%s menulis ke basis data tetapi tidak ada di layananPengubah "+
			"maupun di daftar pengecualian.\n"+
			"Butir bb: tiap layanan pengubah harus menolak kasus yang sudah "+
			"ditutup, atau menyatakan alasannya dengan nama.", nama)
	}
	if diperiksa == 0 {
		t.Fatal("nol layanan bertransaksi ditemukan - penjaga ini tidak menjaga apa pun")
	}
}

// rutePengubah memetakan SETIAP rute non-GET ke berkas layanan yang
// melayaninya.
//
// ⛔ LUBANG YANG INI TUTUP. Penjaga di atas menemukan calonnya lewat
// `strings.Contains(…, "DalamTransaksi(ctx")`. `hapus.go` MENGUBAH
// (`DELETE /api/klaim-life/{id}`) tetapi tidak memanggilnya sama sekali,
// sehingga ia lolos penemuan dan hanya aman karena kebetulan tertulis tangan
// di `layananPengubah`. Layanan pengubah BARU yang berbentuk seperti
// `hapus.go` akan lolos DARI KEDUA penjaga itu - tanpa satu pun uji merah.
//
// ⚠️ Yang dipakai di sini bukan lagi bentuk KODEnya melainkan TABEL RUTEnya:
// apa pun yang terdaftar dengan metode selain GET adalah pengubah, titik.
// Itu definisi yang tidak dapat diakali dengan menulis kode berbeda bentuk.
var rutePengubah = map[string]string{
	"POST /api/klaim-life":                                                    "pendaftaran.go",
	"DELETE /api/klaim-life/{id}":                                             "hapus.go",
	"POST /api/klaim-life/{id}/peserta/{pesertaId}/akseptasi":                 "akseptasi.go",
	"POST /api/klaim-life/{id}/peserta/{pesertaId}/putaran":                   "hasilkomite.go",
	"POST /api/klaim-life/{id}/peserta/{pesertaId}/adjustment/{adjId}/komite": "komite.go",
	"POST /api/klaim-life/{id}/tahap/{tujuan}":                                "tahap.go",
	"POST /api/klaim-life/{id}/outstanding":                                   "simpanrnm.go",
	"POST /api/klaim-life/{id}/tutup":                                         "tutup.go",
	"POST /api/klaim-life/{id}/adjustment/{adjId}/tolak":                      "statusbaris.go",
	"PUT /api/klaim-life/{id}/peserta/{pesertaId}/tanggal-kejadian":           "dol.go",
	"PUT /api/klaim-life/{id}/peserta/{pesertaId}/tanggal-klaim":              "dol.go",
	"POST /api/klaim-life/{id}/peserta/{pesertaId}/cabut":                     "cabutpeserta.go",
	"POST /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa":                  "diagnosa.go",
	"PUT /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa/{diagId}":          "diagnosa.go",
	"DELETE /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa/{diagId}":       "diagnosa.go",
	"POST /api/klaim-life/{id}/peserta/{pesertaId}/dokumen":                   "unggahan.go",
	"DELETE /api/klaim-life/{id}/dokumen/{dokId}":                             "unggahan.go",
	// Modul PremiumList Life (tiket 01). Keduanya memeriksa kasus tertutup
	// lewat T_WORK_POLIS - lihat pengecualian bernama "polis_penawaran.go"
	// di TestDaftarLayananPengubahMencakupSeluruhRutePengubah.
	// GILIRAN-13 butir bn - tombol portal membuat kasus polis.
	"POST /api/polis-life":                "polis_kasus.go",
	"POST /api/polis-life/{id}/keputusan": "polis_penawaran.go",
	// Tiket 03. Penerbitan nomor MENULIS - ia menaikkan baris penghitung dan
	// menuliskan nomornya ke baris peserta - jadi ia pengubah, dan kasus yang
	// sudah ditutup tidak boleh memperoleh nomor baru.
	"POST /api/polis-life/{id}/nomor": "polis_nomor.go",
	// Tiket 04. `simpan` MENGGANTI isi tabel peserta - pengubah sepenuhnya.
	// `tinjau` POST hanya karena ia MENGIRIM berkas; ia tidak menyentuh apa
	// pun, dan itu dinyatakan di daftar pengecualian di bawah.
	"POST /api/polis-life/{id}/unggah/simpan": "polis_unggah.go",
	"POST /api/polis-life/{id}/unggah/tinjau": "polis_unggah_tinjau.go",
	// Tiket 05a bagian 2. Submit summary MENULIS - nomor, rekap, dan salinan
	// peserta warisan dalam satu transaksi - jadi ia pengubah.
	"POST /api/polis-life/{id}/summary": "polis_summary.go",
	// Komite Claim Life tiket 02 - keputusan satu tingkat MENULIS anak tangga
	// dan kepala kasus.
	"POST /api/komite/{id}/keputusan": "komite_keputusan.go",
	// Tiket 03 - eskalasi MENULIS anak tangga dan kepala kasus.
	"POST /api/komite/{id}/eskalasi": "komite_keputusan.go",
}

// letakTabelRute - berkas yang mendaftarkan rute berpenjaga kasus tertutup.
var letakTabelRute = []string{
	filepath.Join("..", "handlers", "handlers.go"),
	filepath.Join("..", "..", "..", "modul", "premiumlistlife", "handlers", "rute_premiumlist.go"),
	filepath.Join("..", "..", "..", "modul", "komiteclaimlife", "handlers", "rute_komite.go"),
}

// bacaBerkasLayanan membaca berkas layanan pelayan rute, di modul mana pun.
//
// Refactor bentuk B (30-09-2026): dulu selalu di folder ini; layanan modul
// yang sudah pindah tinggal di modul/<nama>/services.
func bacaBerkasLayanan(nama string) ([]byte, error) {
	if b, err := os.ReadFile(nama); err == nil {
		return b, nil
	}
	cocok, _ := filepath.Glob(filepath.Join("..", "..", "..", "modul", "*", "services", nama))
	if len(cocok) != 1 {
		return nil, fmt.Errorf("%s ditemukan di %d folder layanan modul", nama, len(cocok))
	}
	return os.ReadFile(cocok[0])
}

var polaRute = regexp.MustCompile(`mux\.HandleFunc\(\s*\n?\s*"([A-Z]+) ([^"]+)"`)

// penjagaTutupBawaan adalah panggilan yang membuktikan sebuah layanan
// menolak kasus yang sudah ditutup.
const penjagaTutupBawaan = "PastikanKasusTerbuka(ctx, klaimID)"

// penjagaTutup menyebut penjaga LAIN bagi layanan yang tabel kerjanya
// bukan `T_WORK_CLAIM`.
//
// ⛔ DIPERSEMPIT 28-09-2026, dan sebabnya prinsip yang sama dengan
// penjaga nama tabel telanjang: penjaga yang menuduh hal yang BENAR akan
// dilonggarkan orang, bukan dipatuhi. Ronde pertama menuntut
// `PastikanKasusTerbuka` dari SETIAP layanan pengubah - kalimat yang benar
// selama hanya ada satu modul. `polis_penawaran.go` memang memeriksa kasus
// tertutup, tetapi lewat `T_WORK_POLIS`: layanan polis yang menanyai tabel
// klaim akan selalu menjawab "tidak ada".
//
// ⚠️ Ini BUKAN pengecualian. Berkas yang tidak ada di peta ini tetap
// dituntut penjaga bawaan, dan yang ada di sini tetap dituntut penjaga yang
// DISEBUT namanya. Yang nol penjaga tetap gagal.
var penjagaTutup = map[string]string{
	"polis_penawaran.go": "models.KasusPolisTertutup(k.Status)",
	"polis_nomor.go":     "models.KasusPolisTertutup(keadaanKerja.Status)",
	"polis_unggah.go":    "models.KasusPolisTertutup(keadaanKerja.Status)",
	"polis_summary.go":   "models.KasusPolisTertutup(keadaan.Status)",
}

func TestSetiapRuteNonGETPunyaPenjagaKasusTertutup(t *testing.T) {
	// Refactor bentuk B (30-09-2026): tabel rute kini tersebar - Router Claim
	// Life dan pendaftar rute modul yang sudah pindah ke modul/<nama>/.
	var isi []byte
	for _, letak := range letakTabelRute {
		b, err := os.ReadFile(letak)
		if err != nil {
			t.Fatalf("membaca tabel rute %s: %v", letak, err)
		}
		isi = append(append(isi, b...), '\n')
	}
	cocok := polaRute.FindAllStringSubmatch(string(isi), -1)
	if len(cocok) == 0 {
		t.Fatal("nol rute terbaca; pembacanya yang rusak, bukan kodenya")
	}
	// Layanan yang sengaja TIDAK memeriksa, beserta alasan tertulis.
	dikecualikan := map[string]string{
		"pendaftaran.go": "melahirkan kasus; belum ada kasus untuk ditutup",
		"polis_kasus.go": "melahirkan kasus polis; belum ada kasus untuk ditutup",
		// Tutup MENUTUP; ia membaca STATUS_WORK sendiri lalu menolak lewat
		// ErrKasusSudahTertutup. Memanggil PastikanKasusTerbuka di sini
		// berarti membaca baris yang sama dua kali untuk satu jawaban.
		"tutup.go": "layanan penutupnya sendiri; memeriksa dengan pembacaan status kerjanya sendiri",
		// ⛔ Tinjauan unggahan MEMBACA SAJA: ia mengurai berkas, memvalidasi,
		// dan mengembalikan hasilnya. Nol tulisan, nol transaksi. Ia POST
		// semata karena GET tidak punya badan permintaan untuk membawa
		// berkas. Menuntutnya menolak kasus tertutup berarti melarang orang
		// MELIHAT apa yang salah dengan berkasnya pada kasus yang sudah
		// selesai - larangan yang tidak melindungi apa pun.
		"polis_unggah_tinjau.go": "tinjauan membaca saja; nol tulisan, nol transaksi",
	}
	cacahPengubah := 0
	for _, m := range cocok {
		metode, jalur := m[1], m[2]
		if metode == "GET" {
			continue
		}
		cacahPengubah++
		kunci := metode + " " + jalur
		berkas, terdaftar := rutePengubah[kunci]
		if !terdaftar {
			t.Errorf("rute pengubah %q tidak ada di rutePengubah.\n"+
				"Butir bb: setiap rute non-GET harus menolak kasus yang sudah "+
				"ditutup. Daftarkan berkas layanannya, atau tuliskan alasan "+
				"pengecualiannya - keduanya pekerjaan sadar.", kunci)
			continue
		}
		if alasan, ada := dikecualikan[berkas]; ada {
			if alasan == "" {
				t.Errorf("%s dikecualikan tanpa alasan tertulis", berkas)
			}
			continue
		}
		b, err := bacaBerkasLayanan(berkas)
		if err != nil {
			t.Errorf("rute %q menunjuk %s yang tidak terbaca: %v", kunci, berkas, err)
			continue
		}
		penjaga, punya := penjagaTutup[berkas]
		if !punya {
			penjaga = penjagaTutupBawaan
		}
		if !strings.Contains(string(b), penjaga) {
			t.Errorf("rute pengubah %q dilayani %s, dan %s tidak memanggil %s.",
				kunci, berkas, berkas, penjaga)
		}
	}
	if cacahPengubah == 0 {
		t.Fatal("nol rute non-GET ditemukan; pembacanya yang rusak")
	}
	// Dan arah sebaliknya: daftar tidak boleh memuat rute yang sudah tidak ada.
	hidup := map[string]bool{}
	for _, m := range cocok {
		hidup[m[1]+" "+m[2]] = true
	}
	for kunci := range rutePengubah {
		if !hidup[kunci] {
			t.Errorf("rutePengubah memuat %q yang tidak lagi terdaftar di Router; "+
				"daftar yang menjaga rute mati tidak menjaga apa pun", kunci)
		}
	}
}
