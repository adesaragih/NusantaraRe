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
)

// berkasKatalogWarisan adalah dokumen [data DBA] hasil pembacaan ALL_TAB_COLUMNS.
const berkasKatalogWarisan = "../../../.scratch/claim-life/TIPE-KOLOM-OS-AKSEPTASI-KLAIM-LIFE.md"

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
// ⚠️ APA YANG SEBENARNYA TERKUNCI DI SINI - dikatakan terus terang, sebab
// laporan ronde 5 sempat mengklaim lebih:
//
//   - CLAIM_RETRO TIDAK diuji apa pun oleh test ini. BarisLamaDari tidak
//     pernah mengisinya, jadi nilainya selalu kosong dan selalu lolos lewat
//     cabang "kosong itu sah". Jaminan yang diminta brief belum ada.
//   - STS_REJECT diisi dari models KodeStatus, yang SENGAJA teks bebas dan
//     dikunci begitu oleh TestKodeStatusTidakPernahJadiBilangan. Test ini
//     lulus hanya karena fixture memakai "1"; kode status non-angka akan
//     tetap sampai ke Oracle dan dijawab ORA-01722.
//
// ⛔ [terbuka] Itu pertentangan nyata antara ADR-U-0022 (kode tetap teks) dan
// katalog (STS_REJECT NUMBER(38,0)), dan executor TIDAK memutuskannya:
// menambah penolakan di jalur tulis berarti diam-diam memihak katalog dan
// melanggar ADR; membiarkannya berarti menunggu ORA-01722 di lapangan.
// Pemiliknya work owner. Yang dilakukan test ini sampai itu dijawab: menjaga
// agar fixture tidak menambah kasus baru yang pasti gagal.
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
