package services

// Gerbang dokumen sebelum Save ke Outstanding - tiket 03.
//
// Untuk apa berkas ini: `SaveOutStandingLife_Act` menolak menyimpan bila ada
// peserta yang belum mengunggah dokumen (langkah 3-4), dan berkas ini meniru
// gerbang itu; ia juga memegang daftar kategori dokumen (butir ar1) yang
// dipakai validasi unggahan.
//
// ⛔ BUTIR bl (GILIRAN-12, OQ-N6 ditutup `[DIPUTUSKAN; veto work owner]`):
// gerbang KEDUA rule itu - dokumen lengkap per kategori, langkah 12 - ter-remark
// (`pyStepsBlockName = //`, b6178) dan TIDAK PERNAH berlaku di sistem lama.
// `PeriksaDokumenLengkap`, `PesertaDokumenTidakLengkap`, `KategoriBerbeda`,
// dan `ErrDokumenTidakLengkap` karena itu DIBUANG beserta ujinya. Bila bisnis
// menghendaki gerbang itu, ia keputusan BARU, bukan replikasi.
//
// Dibaca sesudah: adjustment.go.
//
// ⛔ Gerbangnya MURNI. Yang tidak murni hanya SumberKategoriWajib, dan ia
// sengaja antarmuka: isinya belum ada di korpus mana pun.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

var (
	// ErrDokumenBelumDiunggah - langkah 3 + 3.1 + 4.
	ErrDokumenBelumDiunggah = errors.New("services: ada peserta yang belum mengunggah dokumen")
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

// SumberKategoriWajib memberi daftar kategori dokumen (butir ar1).
//
// Kini dipakai validasi unggahan saja: gerbang kelengkapan yang dulu
// membandingkan cacahnya (`GetCategoryLife_SQL`, langkah 12) ter-remark dan
// dibuang (butir bl).
//
// ⛔ Kenapa ini antarmuka dan bukan daftar tetap: rule `GetCategoryLife_SQL`
// **tidak ada di korpus**; seluruh 29 berkas `Claim Life/RDBList/` sudah
// dicacah. Menuliskan daftarnya sendiri berarti mengarang kebijakan dokumen
// sebuah perusahaan reasuransi.
//
// Sama dengan Penomor pada tiket 02: tempatnya dipisah, dan yang belum
// diketahui terlihat sebagai satu galat terang.
//
// ⭐ 11-10-2026 (keputusan work owner): berbatas KLAIM - daftarnya Document Claim produk polis klaim itu.
type SumberKategoriWajib interface {
	KategoriWajib(ctx context.Context, klaimID string) ([]string, error)
}

// KategoriWajibBelumDiketahui adalah implementasi bawaan; ia selalu gagal.
type KategoriWajibBelumDiketahui struct{}

// KategoriWajib selalu gagal, dengan pesan yang menyebut apa yang ditunggu.
func (KategoriWajibBelumDiketahui) KategoriWajib(context.Context, string) ([]string, error) {
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
// ✅ 10-10-2026 (migrasi 026): rule aslinya menyaring `.IsAccept=="true"` (`SaveOutStandingLife_Act` b1181,
// WhenTrue 2 / WhenFalse 3) - kini dibaca dari IS_ACCEPT, checkbox Accept layar Outstanding. Dulu `IS_CHECK` dipakai
// sebagai padanan terdekat; itu ralat.
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

// nomorPesertaYangGagal - lewati peserta yang tidak di-Accept, kumpulkan NOMOR
// URUT yang gagal syaratnya (dipakai pula gerbang 1 Save to RNM).
func nomorPesertaYangGagal(peserta []models.Peserta, gagal func(models.Peserta) bool) []int {
	var out []int
	for i, p := range peserta {
		if !diterimaDiOutstanding(p) {
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

// dipilihUntukDiklaim membaca penanda IS_CHECK apa adanya.
//
// Nilainya teks warisan, dan korpus memakai "true" maupun "'true'" - keduanya
// diterima, apa pun yang tidak dikenal dianggap TIDAK dipilih.
// diterimaDiOutstanding - `.IsAccept=="true"` (b1181), IS_ACCEPT migrasi 026. Bentuk bertanda kutip dan "1"/"Y"
// tetap dibaca seperti pembaca IS_CHECK lama: nilai warisan yang sama-sama berarti tercentang.
func diterimaDiOutstanding(p models.Peserta) bool {
	v := strings.Trim(strings.TrimSpace(p.IsAccept), "'\"")
	return strings.EqualFold(v, "true") || v == "1" || strings.EqualFold(v, "Y")
}

type kategoriOracle struct{ svc *Service }

// KategoriWajibOracle menyusun pembaca kategori dokumen klaim.
//
// ⭐ 11-10-2026 (keputusan work owner "category document didapat dari document claim yang ada di product name yang
// digunakan oleh claim"): sumbernya Document Claim Master Product Name Life (M_PRODUCTNAME_LIFE_DOCCLAIM) untuk
// ProductNameID polis klaim (header polis via IDPEGA peserta) - menggantikan CATEGORY_ATTACH_CLAIMLIFE (butir ar1, yang
// kini tidak dipakai jalur mana pun; gerbang dokumen Save to RNM hanya menghitung ADA-nya dokumen, bukan kategorinya).
func KategoriWajibOracle(svc *Service) SumberKategoriWajib {
	return kategoriOracle{svc: svc}
}

// KategoriWajib membaca Document Claim produk polis klaim itu, urutan master. Polis tanpa produk / produk tanpa daftar
// = kosong (layar menyatakannya; unggahan menolak 422).
func (k kategoriOracle) KategoriWajib(ctx context.Context, klaimID string) ([]string, error) {
	klaim, err := k.svc.KlaimLife().Ambil(ctx, klaimID)
	if err != nil {
		return nil, err
	}
	polis, _, err := RingkasPolisIdPega(ctx, k.svc.PembacaPolis(), IdPegaPeserta(klaim.Peserta))
	if errors.Is(err, kontrak.ErrPolisNomorTakDitemukan) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(polis.ProductNameID) == "" {
		return []string{}, nil
	}
	return repository.NewPohonKlaim(k.svc.DB()).DokumenKlaimProduk(ctx, strings.TrimSpace(polis.ProductNameID))
}
