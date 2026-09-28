package services

// Gerbang dokumen sebelum Save ke Outstanding - tiket 03.
//
// Untuk apa berkas ini: `SaveOutStandingLife_Act` menolak menyimpan bila
// dokumen peserta belum memenuhi syarat. Rule itu MEMUAT dua gerbang, memakai
// dua daftar yang berbeda, dan berkas ini memisahkan keduanya seperti aslinya.
//
// ⛔ RALAT 28-09-2026 (GILIRAN-11, temuan /code-review): hanya gerbang PERTAMA
// yang hidup. Langkah 12 - gerbang kedua - ber-`pyStepsBlockName = //`
// (b6178): ter-remark, tidak pernah jalan. Save to RNM (simpanrnm.go) karena
// itu tidak memanggil PeriksaDokumenLengkap; lihat komentarnya.
//
// Dibaca sesudah: adjustment.go.
//
// ⛔ Kedua gerbangnya MURNI. Yang tidak murni hanya SumberKategoriWajib, dan
// ia sengaja antarmuka: isinya belum ada di korpus mana pun.

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"nusantarare/internal/models"

	"nusantarare/internal/repository"
)

var (
	// ErrDokumenBelumDiunggah - gerbang pertama, langkah 3 + 3.1 + 4.
	ErrDokumenBelumDiunggah = errors.New("services: ada peserta yang belum mengunggah dokumen")
	// ErrDokumenTidakLengkap - gerbang kedua, langkah 12.
	//
	// ⚠️ Teksnya PERSIS pesan `Local.Err4` di XML, dan sengaja TIDAK menyebut
	// peserta mana: rule aslinya memang tidak menyebutnya. Peserta yang
	// bermasalah dikembalikan terpisah lewat PesertaDokumenTidakLengkap,
	// supaya pemanggil tetap dapat menunjukkannya tanpa kami mengarang teks
	// pesan yang tidak ada di sistem lama.
	ErrDokumenTidakLengkap = errors.New(
		"Documents are incomplete, please complete the documents")
	// ErrKategoriWajibBelumDiketahui - daftar kategori wajibnya belum ada.
	ErrKategoriWajibBelumDiketahui = errors.New(
		"services: daftar kategori dokumen wajib belum diketahui")
)

// tipeTanpaGerbangDokumen adalah Type yang DILEWATI gerbang pertama.
//
// `[terverifikasi]` precondition langkah 3.1: `pyWorkPage.Type=="TP" ||
// pyWorkPage.Type=="TR"` dengan WhenTrue=3 - benar berarti LEWATI langkah.
// Jadi klaim treaty tidak dituntut berdokumen.
var tipeTanpaGerbangDokumen = map[string]bool{"TP": true, "TR": true}

// SumberKategoriWajib memberi daftar kategori dokumen yang wajib ada.
//
// ⛔ Kenapa ini antarmuka dan bukan daftar tetap: MEKANISMEnya terbaca dari
// XML - cacah kategori berbeda yang terunggah harus sama dengan cacah baris
// yang `GetCategoryLife_SQL` kembalikan - tetapi rule itu **tidak ada di
// korpus**; seluruh 29 berkas `Claim Life/RDBList/` sudah dicacah. Menuliskan
// daftarnya sendiri berarti mengarang kebijakan dokumen sebuah perusahaan
// reasuransi.
//
// Sama dengan Penomor pada tiket 02: tempatnya dipisah, dan yang belum
// diketahui terlihat sebagai satu galat terang.
type SumberKategoriWajib interface {
	KategoriWajib(ctx context.Context) ([]string, error)
}

// KategoriWajibBelumDiketahui adalah implementasi bawaan; ia selalu gagal.
type KategoriWajibBelumDiketahui struct{}

// KategoriWajib selalu gagal, dengan pesan yang menyebut apa yang ditunggu.
func (KategoriWajibBelumDiketahui) KategoriWajib(context.Context) ([]string, error) {
	return nil, fmt.Errorf("%w: rule GetCategoryLife_SQL tidak ada di korpus, dan "+
		"korpus tidak pernah mengenumerasi nilai KATEGORI_1/KATEGORI_2; daftarnya "+
		"diperlukan dari DBA atau work owner lebih dulu", ErrKategoriWajibBelumDiketahui)
}

// PeriksaDokumenAda adalah gerbang PERTAMA - langkah 3, 3.1, dan 4.
//
// Meniru ketiga precondition-nya beserta aksinya:
//
//	Type "TP" atau "TR"          -> seluruh gerbang dilewati
//	peserta tidak dipilih        -> peserta itu dilewati
//	peserta nol dokumen          -> pesan bertambah, menyebut NOMOR URUTnya
//
// ⚠️ "person number" di pesan aslinya adalah `local.idx2 = .pxListSubscript`,
// yaitu NOMOR URUT peserta dalam daftar - bukan pengenalnya. Nomor urut 1-based
// dipertahankan supaya pesannya dapat dibandingkan dengan sistem lama.
//
// ⚠️ `[terbuka]` Rule aslinya menyaring `.IsAccept=="true"`; tabel kita tidak
// punya kolom itu. `IS_CHECK` - "dipilih untuk diklaim" - adalah padanan
// terdekat yang ada, dan bukan hal yang sama. Dicatat, tidak didiamkan.
func PeriksaDokumenAda(tipe string, peserta []models.Peserta) error {
	if tipeTanpaGerbangDokumen[strings.ToUpper(strings.TrimSpace(tipe))] {
		return nil
	}
	kurang := nomorPesertaYangGagal(peserta, func(p models.Peserta) bool {
		return len(p.Dokumen) == 0
	})
	if len(kurang) == 0 {
		return nil
	}
	// Teks ini meniru `local.Errmsg10` di XML, yang memang menyebut nomor urut.
	return fmt.Errorf("%w: person number %s",
		ErrDokumenBelumDiunggah, gabungNomor(kurang))
}

// PeriksaDokumenLengkap adalah gerbang KEDUA - langkah 12 sampai 12.4.
//
// Aturannya `[terverifikasi]`: kategori dokumen peserta dikumpulkan,
// DI-DEDUP (langkah Java atas `.CARI1`), lalu cacahnya dibandingkan dengan
// cacah kategori wajib. Pega membandingkan CACAHnya, bukan himpunannya - dan
// peniruan yang jujur mengikuti itu, termasuk kelemahannya: dua dokumen
// berkategori salah dengan cacah yang kebetulan pas akan lolos.
//
// ⛔ Perbandingan cacah dipertahankan APA ADANYA, bukan "diperbaiki" menjadi
// perbandingan himpunan. Memperbaiki diam-diam berarti sistem baru menolak
// klaim yang sistem lama terima, tanpa seorang pun memutuskannya.
//
// ⛔ TANPA PEMANGGIL PRODUKSI, dan itu disengaja: langkah 12 ter-remark
// (b6178), jadi sistem lama TIDAK PERNAH menolak simpan karena dokumen tidak
// lengkap. Memasangnya adalah penyimpangan baru yang harus diputuskan work
// owner (OQ-N6), bukan paritas - ia dipertahankan hanya sebagai aturan siap
// pakai bila keputusan itu jatuh; bila tidak, ia dibuang.
func PeriksaDokumenLengkap(peserta []models.Peserta, wajib []string) error {
	if len(wajib) == 0 {
		return fmt.Errorf("%w: daftar kategori wajib kosong", ErrKategoriWajibBelumDiketahui)
	}
	if len(PesertaDokumenTidakLengkap(peserta, wajib)) == 0 {
		return nil
	}
	// ⛔ Pesannya berhenti di sini, tanpa nomor peserta. XML tidak menyebutnya,
	// dan pesan galat adalah logika bisnis: menambahinya berarti sistem baru
	// berbicara dengan kalimat yang tidak pernah ada.
	return ErrDokumenTidakLengkap
}

// PesertaDokumenTidakLengkap mengembalikan NOMOR URUT peserta yang gagal
// gerbang kedua - 1-based, seperti `.pxListSubscript`.
func PesertaDokumenTidakLengkap(peserta []models.Peserta, wajib []string) []int {
	return nomorPesertaYangGagal(peserta, func(p models.Peserta) bool {
		return len(KategoriBerbeda(p)) != len(wajib)
	})
}

// nomorPesertaYangGagal memakai satu kerangka untuk kedua gerbang: lewati
// peserta yang tidak dipilih, kumpulkan NOMOR URUT yang gagal syaratnya.
func nomorPesertaYangGagal(peserta []models.Peserta, gagal func(models.Peserta) bool) []int {
	var out []int
	for i, p := range peserta {
		if !dipilihUntukDiklaim(p) {
			continue
		}
		if gagal(p) {
			out = append(out, i+1)
		}
	}
	return out
}

// gabungNomor menulis daftar nomor urut untuk pesan galat.
func gabungNomor(nomor []int) string {
	teks := make([]string, 0, len(nomor))
	for _, n := range nomor {
		teks = append(teks, fmt.Sprintf("%d", n))
	}
	return strings.Join(teks, ", ")
}

// KategoriBerbeda mengumpulkan kategori dokumen seorang peserta, tanpa kembar.
//
// ⚠️ `[terbuka - work owner]` Kolom mana yang memegang kategori pembanding
// belum diputuskan. Gerbang Pega membacanya dari daftar LAMPIRAN bawaan
// (`.pyCategory`), yang di skema relasional tidak punya padanan langsung.
// `KATEGORI_2` dipakai di sini karena itulah kategori yang `Section/
// DocumentLife.xml` tampilkan kepada manusia - alasan yang dinyatakan, bukan
// tebakan yang didiamkan.
func KategoriBerbeda(p models.Peserta) []string {
	lihat := map[string]bool{}
	var out []string
	for _, d := range p.Dokumen {
		k := strings.TrimSpace(d.Kategori2)
		if k == "" || lihat[k] {
			continue
		}
		lihat[k] = true
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// dipilihUntukDiklaim membaca penanda IS_CHECK apa adanya.
//
// Nilainya teks warisan, dan korpus memakai "true" maupun "'true'" - keduanya
// diterima, apa pun yang tidak dikenal dianggap TIDAK dipilih.
func dipilihUntukDiklaim(p models.Peserta) bool {
	v := strings.Trim(strings.TrimSpace(p.IsCheck), "'\"")
	return strings.EqualFold(v, "true") || v == "1" || strings.EqualFold(v, "Y")
}

// kategoriOracle membaca daftar kategori wajib - butir ar1, A2.
type kategoriOracle struct{ pohon *repository.PohonKlaim }

// KategoriWajibOracle menyusun pembaca kategori yang memakai Oracle.
//
// ⚠️ `ar1` `[USULAN yang disahkan]`, bukan `[terverifikasi]`: rule
// `GetCategoryLife_SQL` tidak ada di korpus, sehingga tabel mana yang Pega
// baca tidak dapat dipastikan. Kandidat kedua dicatat di repository dan tidak
// dipakai.
func KategoriWajibOracle(svc *Service) SumberKategoriWajib {
	return kategoriOracle{pohon: repository.NewPohonKlaim(svc.db)}
}

// KategoriWajib membaca daftar kategori untuk sebuah kode bisnis.
func (k kategoriOracle) KategoriWajib(ctx context.Context) ([]string, error) {
	// ⛔ Kode bisnis belum tersimpan di model relasional - temuan audit A0,
	// kolomnya lahir di `014`. Sampai pembacanya ada, penyaring `BISNIS`
	// memakai daftar `ALL` saja, dan itu DINYATAKAN, bukan disembunyikan:
	// keenam baris DEV memang seluruhnya `ALL`.
	return k.pohon.AmbilKategoriWajib(ctx, "")
}
