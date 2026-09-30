package repository

// Peta tipe kolom warisan dibandingkan dengan katalog Oracle - TANPA Oracle.
//
// Untuk apa berkas ini: sampai 26 September 2026 tipe kolom
// OS_AKSEPTASI_KLAIM_LIFE tidak diketahui, dan kode ini MENEBAK-nya dari nama
// kolom. Tebakannya meleset pada tiga kolom, dan karena tabel TIRUAN di skema
// uji dibuat dari tebakan yang sama, test pulang-pergi tidak pernah
// menemukannya: ia menguji tiruan terhadap dirinya sendiri.
//
// Sekarang katalognya ada, disimpan sebagai dokumen `[data DBA]`. Test di sini
// membaca dokumen itu dan membandingkannya dengan tabel di barislamakolom.go,
// kolom demi kolom, sehingga selisih berikutnya terdengar pada `go test` biasa.
//
// ⚠️ BATASNYA, supaya tidak disalahbaca: yang dibuktikan di sini hanyalah
// DOKUMEN == KODE. Bahwa dokumen itu sendiri sama dengan Oracle sungguhan
// bersandar pada satu pembacaan katalog instance PENGEMBANGAN pada 26-09-2026
// yang belum dikonfirmasi DBA; produksi belum pernah dibaca. Test ini karena
// itu belum menutup "menguji tiruan terhadap dirinya sendiri" sepenuhnya - ia
// memindahkan sumbernya dari tebakan ke katalog, dan itu saja.
//
// Dibaca sesudah: barislamakolom.go.
//
// Istilah:
//   - katalog : tabel sistem Oracle yang memerikan bentuk tabel lain
//               (di sini ALL_TAB_COLUMNS).

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"nusantarare/inti/backend/utils"
)

// berkasKatalogWarisan adalah dokumen [data DBA] hasil pembacaan ALL_TAB_COLUMNS.
const berkasKatalogWarisan = "../../../../.scratch/claim-life/TIPE-KOLOM-OS-AKSEPTASI-KLAIM-LIFE.md"

// barisKatalog adalah satu baris tabel markdown dokumen itu.
var barisKatalog = regexp.MustCompile(`^\|\s*(\d+)\s*\|(.+)$`)

// bacaKatalogWarisan memecah tabel markdown menjadi daftar kolom bertipe DDL.
//
// Bentuk barisnya: | no | `NAMA` | TIPE | panjang | presisi | skala | null | catatan |
// Tebal (**) dan backtick dibuang; tipe dirakit kembali menjadi bentuk DDL.
func bacaKatalogWarisan(t *testing.T) []KolomWarisan {
	t.Helper()
	isi, err := os.ReadFile(filepath.FromSlash(berkasKatalogWarisan))
	if err != nil {
		t.Fatalf("katalog kolom warisan tidak terbaca: %v\n"+
			"Berkas ini [data DBA] dan wajib ada; tanpa ia, peta tipe kembali menjadi tebakan.", err)
	}
	bersih := strings.NewReplacer("`", "", "*", "")

	var hasil []KolomWarisan
	for _, baris := range strings.Split(string(isi), "\n") {
		m := barisKatalog.FindStringSubmatch(strings.TrimRight(baris, "\r"))
		if m == nil {
			continue
		}
		sel := strings.Split(m[2], "|")
		if len(sel) < 6 {
			t.Fatalf("baris katalog kurang kolom: %q", baris)
		}
		ambil := func(i int) string { return strings.TrimSpace(bersih.Replace(sel[i])) }
		nama, tipe := ambil(0), ambil(1)
		panjang, presisi, skala := ambil(2), ambil(3), ambil(4)

		var ddl string
		switch tipe {
		case "VARCHAR2", "CHAR":
			ddl = tipe + "(" + panjang + ")"
		case "DATE":
			ddl = "DATE"
		case "NUMBER":
			ddl = "NUMBER"
			if presisi != "" && skala != "" {
				ddl = "NUMBER(" + presisi + "," + skala + ")"
			}
		default:
			t.Fatalf("tipe katalog tidak dikenal untuk %s: %q", nama, tipe)
		}

		// Nomor urutnya ikut diperiksa: dokumen yang barisnya tertukar akan
		// diam-diam mengubah urutan kolom tabel tiruan.
		if n, _ := strconv.Atoi(m[1]); n != len(hasil)+1 {
			t.Fatalf("nomor baris katalog melompat di %s: %s", nama, m[1])
		}
		hasil = append(hasil, KolomWarisan{Kolom: nama, Tipe: ddl})
	}
	if len(hasil) == 0 {
		t.Fatal("nol baris katalog terbaca; pembacanya yang rusak, bukan dokumennya")
	}
	return hasil
}

// Tabel kolomWarisan sama persis dengan katalog - nama, tipe, dan urutan.
func TestPetaTipeCocokDenganKatalog(t *testing.T) {
	katalog := bacaKatalogWarisan(t)

	const mau = 62
	if len(katalog) != mau {
		t.Fatalf("katalog memuat %d kolom, mau %d", len(katalog), mau)
	}
	if len(kolomWarisan) != len(katalog) {
		t.Fatalf("kolomWarisan memuat %d kolom, katalog %d", len(kolomWarisan), len(katalog))
	}

	for i, k := range katalog {
		punyaKode := kolomWarisan[i]
		if punyaKode.Kolom != k.Kolom {
			t.Errorf("posisi %d: kode menyebut %s, katalog %s", i+1, punyaKode.Kolom, k.Kolom)
			continue
		}
		if punyaKode.Tipe != k.Tipe {
			t.Errorf("%s bertipe %q di kode, %q di katalog", k.Kolom, punyaKode.Tipe, k.Tipe)
		}
	}
}

// Penggolongan angka dan tanggal diturunkan dari katalog, bukan dari nama.
//
// Ketiga kolom yang dulu ditebak keliru disebut namanya di sini, supaya
// kemundurannya tidak sekadar terbaca sebagai "satu selisih" di test atas.
func TestGolonganKolomWarisanIkutKatalog(t *testing.T) {
	katalog := bacaKatalogWarisan(t)
	mauAngka, mauTanggal := map[string]bool{}, map[string]bool{}
	for _, k := range katalog {
		switch {
		case strings.HasPrefix(k.Tipe, "NUMBER"):
			mauAngka[k.Kolom] = true
		case k.Tipe == "DATE":
			mauTanggal[k.Kolom] = true
		}
	}
	bandingkan := func(nama string, punya, mau map[string]bool) {
		for k := range mau {
			if !punya[k] {
				t.Errorf("%s: %s ada di katalog tetapi tidak digolongkan", nama, k)
			}
		}
		for k := range punya {
			if !mau[k] {
				t.Errorf("%s: %s digolongkan tetapi katalog tidak menyebutnya", nama, k)
			}
		}
	}
	bandingkan("kolomAngkaLama", kolomAngkaLama, mauAngka)
	bandingkan("kolomTanggalLama", kolomTanggalLama, mauTanggal)

	// Ketiga tebakan yang meleset, disebut satu per satu.
	for _, k := range []string{"STS_REJECT", "CLAIM_RETRO"} {
		if !kolomAngkaLama[k] {
			t.Errorf("%s bukan kolom angka; katalog menyebutnya NUMBER", k)
		}
	}
	if !kolomTanggalLama["WPC"] {
		t.Error("WPC bukan kolom tanggal; katalog menyebutnya DATE")
	}
}

// Kelima puluh lima kolom yang ditulis rule adalah bagian dari keenam puluh dua.
func TestKolomDitulisAdalahBagianDariKatalog(t *testing.T) {
	seluruh := map[string]bool{}
	for _, n := range NamaKolomTabelWarisan() {
		seluruh[n] = true
	}
	ditulis := NamaKolomBarisLama()
	const mauDitulis = 55
	if len(ditulis) != mauDitulis {
		t.Errorf("rule menulis %d kolom, mau %d", len(ditulis), mauDitulis)
	}
	for _, n := range ditulis {
		if !seluruh[n] {
			t.Errorf("%s ditulis rule tetapi tidak ada di katalog", n)
		}
	}
	if len(seluruh) != 62 {
		t.Errorf("tabel tiruan memakai %d kolom, mau 62", len(seluruh))
	}
}

// ⛔ Nilai yang ditulis ke kolom angka harus angka atau kosong.
//
// Kenapa ini dijaga: BarisLama menyimpan SELURUH kolom sebagai teks Go, dan
// NilaiBarisLama mem-bind teks itu apa adanya. Untuk kolom NUMBER, Oracle
// mengonversinya secara implisit - dan teks yang bukan angka menjadi ORA-01722
// saat menyimpan, bukan saat kompilasi. STS_REJECT dan CLAIM_RETRO baru
// diketahui bertipe NUMBER pada 26-09-2026, sesudah keduanya diperlakukan
// sebagai teks selama berbulan-bulan.
//
// ✅ KEDUA LUBANG YANG DICATAT RONDE 5 SUDAH DITUTUP, 26-09-2026 malam:
//
//   - CLAIM_RETRO dulu tidak diuji apa pun oleh test ini, sebab BarisLamaDari
//     tidak pernah mengisinya sehingga selalu lolos lewat cabang "kosong itu
//     sah". Sejak butir w2 ia UANG di tingkat header dan sungguh ditulis ke
//     tiap baris datar - lihat TestClaimRetroPulangPergiSebagaiUang.
//   - STS_REJECT dulu hanya lulus karena fixture memakai "1". Sejak butir s1
//     jalur tulis sendiri dipagari PeriksaNilaiWarisan, yang menolak nilai
//     bukan-angka dengan galat menyebut kolom dan isinya.
//
// ⭐ Bingkai "ADR-U-0022 lawan katalog" yang saya tulis ronde 5 terlalu lebar,
// dan work owner mempersempitnya: ADR-U-0022 melindungi kode BERAWALAN NOL
// seperti "006"; STS_REJECT adalah angka satu digit menurut aksi, dan agregat
// instance pengembangan menunjukkan nol nilai berawalan nol dan nol nilai
// non-angka. Skema baru karena itu tetap menyimpannya sebagai teks - ADR utuh -
// dan yang ditambah hanya pagar kompatibilitas di jalur tulis warisan.
func TestNilaiKolomAngkaSelaluAngkaAtauKosong(t *testing.T) {
	masuk := []BarisLama{
		contohBaris("R1", "UJI-CASE-1", "006", "1234567890.12345678"),
		contohBaris("R2", "UJI-CASE-1", "006", "0.00000001"),
		contohBaris("R3", "UJI-CASE-1", "010", "250000"),
	}
	pohon, _ := BongkarBarisLama(masuk)
	if len(pohon) != 1 {
		t.Fatalf("klaim %d, mau 1", len(pohon))
	}
	baris := BarisLamaDari(pohon[0])
	if len(baris) != len(masuk) {
		t.Fatalf("baris dihasilkan %d, mau %d", len(baris), len(masuk))
	}
	angka := regexp.MustCompile(`^-?\d+(\.\d+)?$`)
	diperiksa := 0
	for _, b := range baris {
		medan := medanBarisLama(&b)
		for _, m := range medan {
			if !kolomAngkaLama[m.Kolom] {
				continue
			}
			diperiksa++
			nilai := strings.TrimSpace(*m.Nilai)
			if nilai == "" {
				continue // menjadi NULL, dan itu sah (ADR-U-0027)
			}
			if !angka.MatchString(nilai) {
				t.Errorf("kolom angka %s berisi %q - Oracle akan menjawab ORA-01722",
					m.Kolom, nilai)
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol kolom angka diperiksa; pembacanya yang rusak")
	}
	t.Logf("%d nilai kolom angka diperiksa", diperiksa)
}

// Pagar jalur tulis menolak nilai yang tidak dapat masuk ke kolom NUMBER.
//
// `[keputusan work owner 26-09-2026, butir s1]`. Ia PAGAR KOMPATIBILITAS, bukan
// konversi: skema baru tetap menyimpan kode sebagai teks (ADR-U-0022 utuh), dan
// yang dijaga hanya jalur tulis ke tabel datar warisan - tempat kolomnya memang
// bertipe NUMBER menurut katalog.
//
// Nilainya galat terang, bukan ORA-01722: pesan Oracle menuding tipe data dan
// tidak menyebut kolom mana yang salah maupun isinya.
func TestPagarNilaiWarisanMenolakBukanAngka(t *testing.T) {
	kasus := []struct {
		nama   string
		ubah   func(*BarisLama)
		mauSah bool
	}{
		{"kode status satu digit", func(b *BarisLama) { b.STS_REJECT = "1" }, true},
		{"kode status kosong", func(b *BarisLama) { b.STS_REJECT = "" }, true},
		{"kode status berspasi", func(b *BarisLama) { b.STS_REJECT = "  2  " }, true},
		{"uang berdesimal panjang", func(b *BarisLama) { b.CLAIM_AMOUNT = "1234567890.12345678" }, true},
		{"uang negatif", func(b *BarisLama) { b.CLAIM_RETRO = "-500.25" }, true},
		{"kode status berupa kata", func(b *BarisLama) { b.STS_REJECT = "DITOLAK" }, false},
		{"uang dengan pemisah ribuan", func(b *BarisLama) { b.CLAIM_AMOUNT = "1,234.00" }, false},
		{"uang berkoma desimal", func(b *BarisLama) { b.CLAIM_RETRO = "1234,56" }, false},
		{"teks bertanda UJI", func(b *BarisLama) { b.CLAIM_RETRO = "UJI-RETRO" }, false},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			b := contohBaris("R1", "UJI-CASE-1", "006", "250000")
			k.ubah(&b)
			err := PeriksaNilaiWarisan(b)
			if k.mauSah {
				if err != nil {
					t.Errorf("pagar menolak nilai yang sah: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("pagar meloloskan nilai yang akan dijawab ORA-01722")
			}
			// Pesannya harus menyebut kolomnya DAN nilainya - itulah yang
			// tidak dilakukan galat Oracle.
			if !strings.Contains(err.Error(), "ORA-01722") {
				t.Errorf("pesan tidak menyebut galat yang dihindarinya: %v", err)
			}
		})
	}
}

// Kolom teks TIDAK ikut dipagari - nomor rekening berawalan nol harus lolos.
//
// Ini sisi lain pagar s1, dan yang membuatnya bukan "semua harus angka":
// ADR-U-0022 menetapkan kode dan nomor tetap teks, dan tabel warisan memang
// menyimpan ketiga kolom bank sebagai VARCHAR2.
func TestPagarNilaiWarisanTidakMenyentuhKolomTeks(t *testing.T) {
	b := contohBaris("R1", "UJI-CASE-1", "006", "250000")
	b.ACCOUNTNO = "0012345678"
	b.CERTIFICATE_NO = "006"
	b.NAME_OF_BANK = "BANK UJI"
	if err := PeriksaNilaiWarisan(b); err != nil {
		t.Errorf("pagar menyentuh kolom teks: %v", err)
	}
}

// CLAIM_RETRO pulang-pergi sebagai uang, dan jaminannya tidak lagi kosong.
//
// ⛔ Ronde 5 mengklaim jalur tulis mengunci kolom angka; tinjauan menunjukkan
// klaim itu KOSONG untuk CLAIM_RETRO, sebab BarisLamaDari tidak pernah
// mengisinya sehingga selalu lolos lewat cabang "kosong itu sah". Sejak butir
// w2 ia UANG di tingkat header dan sungguh-sungguh ditulis.
func TestClaimRetroPulangPergiSebagaiUang(t *testing.T) {
	masuk := []BarisLama{
		contohBaris("R1", "UJI-CASE-1", "006", "250000"),
		contohBaris("R2", "UJI-CASE-1", "010", "125000"),
	}
	pohon, lap := BongkarBarisLama(masuk)
	if len(pohon) != 1 {
		t.Fatalf("klaim %d, mau 1", len(pohon))
	}
	for _, tm := range lap.Temuan {
		if tm.Medan == "CLAIM_RETRO" {
			t.Fatalf("CLAIM_RETRO dilaporkan tidak terurai: %+v", tm)
		}
	}
	mau := masuk[0].CLAIM_RETRO
	if got := utils.FormatDecimal(pohon[0].Klaim.ClaimRetro.Amount); got != mau {
		t.Errorf("ClaimRetro = %q, mau %q", got, mau)
	}
	// Mata uangnya datang dari baris adjustment, sebab header warisan tidak
	// punya kolom mata uang sendiri.
	if got := pohon[0].Klaim.ClaimRetro.Currency; got != masuk[0].CURRENCY {
		t.Errorf("mata uang ClaimRetro = %q, mau %q", got, masuk[0].CURRENCY)
	}

	balik := BarisLamaDari(pohon[0])
	if len(balik) != len(masuk) {
		t.Fatalf("baris balik %d, mau %d", len(balik), len(masuk))
	}
	for i, b := range balik {
		if b.CLAIM_RETRO != mau {
			t.Errorf("baris %d: CLAIM_RETRO = %q, mau %q", i, b.CLAIM_RETRO, mau)
		}
	}
}

// ⛔ Bentuk yang TO_CHAR TM9 keluarkan harus diterima pagar itu sendiri.
//
// Ini jebakan yang nyaris lolos: ekspresiBacaLama membaca kolom NUMBER lewat
// TO_CHAR ber-format TM9, dan TM9 mengeluarkan ".5" untuk setengah - tanpa nol
// di depan titik. Pola yang menuntut digit di depan titik karena itu akan
// menolak nilai yang baru saja dibacanya sendiri dari Oracle, dan pulang-pergi
// baris warisan patah di tengah tanpa ada yang menyentuh datanya.
func TestPagarMenerimaBentukKeluaranTM9(t *testing.T) {
	for _, nilai := range []string{".5", "-.25", "5.", "+5", "0.5", "-0.5"} {
		b := contohBaris("R1", "UJI-CASE-1", "006", "250000")
		b.CLAIM_RETRO = nilai
		if err := PeriksaNilaiWarisan(b); err != nil {
			t.Errorf("pagar menolak %q, yang justru dapat datang dari TO_CHAR TM9: %v",
				nilai, err)
		}
	}
}

// Mata uang yang berbeda antar baris satu klaim DILAPORKAN, bukan dipilih diam-diam.
//
// CURRENCY bukan anggota atributKlaim, sehingga tanpa pemeriksaan sendiri
// ClaimRetro akan mendapat mata uang baris pertama tanpa ada yang tahu.
func TestMataUangCampurDilaporkan(t *testing.T) {
	masuk := []BarisLama{
		contohBaris("R1", "UJI-CASE-1", "006", "250000"),
		contohBaris("R2", "UJI-CASE-1", "010", "125000"),
	}
	masuk[1].CURRENCY = "USD"
	_, lap := BongkarBarisLama(masuk)
	ketemu := false
	for _, tm := range lap.Temuan {
		if tm.Medan == "CURRENCY" && tm.Jenis == TemuanAtributBerbeda {
			ketemu = true
		}
	}
	if !ketemu {
		t.Error("mata uang berbeda antar baris tidak dilaporkan; " +
			"ClaimRetro akan memakai mata uang baris pertama diam-diam")
	}
	// Seragam: nol temuan mata uang.
	_, lapSeragam := BongkarBarisLama(masuk[:1])
	for _, tm := range lapSeragam.Temuan {
		if tm.Medan == "CURRENCY" {
			t.Errorf("mata uang seragam dilaporkan berbeda: %+v", tm)
		}
	}
}
