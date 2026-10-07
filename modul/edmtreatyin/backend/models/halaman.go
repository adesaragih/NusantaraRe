// Package models memegang bentuk data dan aturan murni modul NB Treaty In
// (`nbtreatyin`): halaman kerja, penggolong, penanda persetujuan, rumus uang,
// tangga persetujuan, dan nomor polis. Tidak satu pun berkas di sini menyentuh
// basis data atau HTTP - inilah seam 3 `spec.md` §6.2 ("fungsi murni penggolong
// dan penghitung").
//
// Untuk apa berkas ini: HALAMAN KERJA. Pega menyimpan seluruh data realisasi di
// halaman `pyWorkPage` beserta halaman anaknya (`PolicyTreatyIn`, `Quotation`,
// `TreatyIn`), dan setiap aktivitas membaca/menulis properti lewat jalurnya -
// `.PolicyTreatyIn.PremiOgp`, `pyWorkPage.Quotation.BusinessOldId`. Halaman di
// sini meniru bentuk itu: nilai TEKS per jalur relatif terhadap `pyWorkPage`,
// ditambah daftar baris untuk PageList (`SpreadingRiskList`, `ListInstallment`,
// `SuggestList`, ...). Dengan begitu setiap port aktivitas dapat dibaca
// berdampingan dengan langkah aslinya di `docs/INVENTARIS-XML.md` bab 6.
//
// ⛔ Nilai disimpan sebagai TEKS, persis seperti di clipboard Pega. Konversi ke
// desimal terjadi di titik pakai (`Angka`), dan kembali ke teks lewat
// `utils.FormatDecimal` - tidak pernah lewat float (ADR-U-0003, spec §5.6).
// Konversi ke tipe kolom Oracle terjadi sekali, di repository (ADR-U-0022).
package models

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// Jalur halaman anak di bawah pyWorkPage.
const (
	// HalamanPolis adalah `pyWorkPage.PolicyTreatyIn` - data realisasi.
	HalamanPolis = "PolicyTreatyIn"
	// HalamanQuotation adalah `pyWorkPage.Quotation`.
	HalamanQuotation = "Quotation"
	// HalamanMaster adalah `pyWorkPage.TreatyIn` - nilai kontrak treaty
	// induk. ⭐ Di sistem baru diisi dari view relasional
	// `POOLDATA.TREATYINDETAILJOINEDM`, bukan dari kolom dokumen JSONDATA
	// (spec §5.1, keputusan work owner P29).
	HalamanMaster = "TreatyIn"
)

// Halaman adalah satu halaman kerja: nilai teks per jalur dan daftar baris.
//
// Jalur ditulis RELATIF terhadap `pyWorkPage` tanpa titik di depan:
// `PolicyTreatyIn.PremiOgp`, `Quotation.BusinessOldId`, `NBStatus`.
type Halaman struct {
	Nilai  map[string]string   `json:"nilai"`
	Daftar map[string][]Baris  `json:"daftar"`
	Pesan  map[string][]string `json:"pesan,omitempty"`
}

// Baris adalah satu anggota PageList: nilai teks per nama properti anggota.
type Baris map[string]string

// HalamanBaru membuat halaman kosong yang siap ditulis.
func HalamanBaru() *Halaman {
	return &Halaman{Nilai: map[string]string{}, Daftar: map[string][]Baris{}, Pesan: map[string][]string{}}
}

// pastikan menyiapkan peta yang nil - halaman hasil urai JSON dapat datang
// tanpa salah satunya.
func (h *Halaman) pastikan() {
	if h.Nilai == nil {
		h.Nilai = map[string]string{}
	}
	if h.Daftar == nil {
		h.Daftar = map[string][]Baris{}
	}
	if h.Pesan == nil {
		h.Pesan = map[string][]string{}
	}
}

// Ambil membaca nilai teks satu jalur. Jalur yang tidak pernah ditulis
// bernilai "" - sama dengan properti Pega yang belum diisi.
func (h *Halaman) Ambil(jalur string) string {
	if h == nil || h.Nilai == nil {
		return ""
	}
	return h.Nilai[jalur]
}

// Setel menulis nilai teks satu jalur.
func (h *Halaman) Setel(jalur, nilai string) {
	h.pastikan()
	h.Nilai[jalur] = nilai
}

// Hapus meniru `Property-Remove`: jalurnya kembali tak bernilai.
func (h *Halaman) Hapus(jalur string) {
	h.pastikan()
	delete(h.Nilai, jalur)
}

// SetelAngka menulis desimal sebagai teks bertitik (`utils.FormatDecimal`).
func (h *Halaman) SetelAngka(jalur string, d *apd.Decimal) {
	h.Setel(jalur, utils.FormatDecimal(d))
}

// AmbilDaftar membaca PageList; daftar yang belum ada bernilai nil.
func (h *Halaman) AmbilDaftar(jalur string) []Baris {
	if h == nil || h.Daftar == nil {
		return nil
	}
	return h.Daftar[jalur]
}

// SetelDaftar mengganti seluruh PageList.
func (h *Halaman) SetelDaftar(jalur string, b []Baris) {
	h.pastikan()
	h.Daftar[jalur] = b
}

// TambahPesan meniru `Property-Set-Messages` / `Page-Set-Messages`: pesan galat
// ditempel pada satu jalur (atau "" untuk pesan halaman).
func (h *Halaman) TambahPesan(jalur, pesan string) {
	h.pastikan()
	h.Pesan[jalur] = append(h.Pesan[jalur], pesan)
}

// BersihkanPesan meniru `Page-Clear-Messages`.
func (h *Halaman) BersihkanPesan() {
	h.pastikan()
	h.Pesan = map[string][]string{}
}

// SemuaPesan merangkai pesan galat, urut jalur, untuk dikembalikan ke layar.
func (h *Halaman) SemuaPesan() []string {
	if h == nil {
		return nil
	}
	kunci := make([]string, 0, len(h.Pesan))
	for k := range h.Pesan {
		kunci = append(kunci, k)
	}
	sort.Strings(kunci)
	var out []string
	for _, k := range kunci {
		for _, p := range h.Pesan[k] {
			if k == "" {
				out = append(out, p)
			} else {
				out = append(out, k+": "+p)
			}
		}
	}
	return out
}

// Salin membuat salinan dalam - aktivitas boleh menulis tanpa mengubah
// halaman pemanggil (dipakai uji dan oleh services saat menghitung pratinjau).
func (h *Halaman) Salin() *Halaman {
	s := HalamanBaru()
	if h == nil {
		return s
	}
	for k, v := range h.Nilai {
		s.Nilai[k] = v
	}
	for k, d := range h.Daftar {
		baru := make([]Baris, len(d))
		for i, b := range d {
			nb := Baris{}
			for kk, vv := range b {
				nb[kk] = vv
			}
			baru[i] = nb
		}
		s.Daftar[k] = baru
	}
	for k, p := range h.Pesan {
		s.Pesan[k] = append([]string(nil), p...)
	}
	return s
}

// ---------------------------------------------------------------- angka

// ErrBukanAngka - nilai properti bukan bilangan desimal. Pega melempar galat
// ekspresi di titik yang sama; sistem baru mengembalikannya sebagai galat
// yang menyebut jalurnya.
var ErrBukanAngka = errors.New("models: nilai properti bukan angka")

// AngkaTeks membaca teks properti sebagai desimal dengan semantik ekspresi
// Pega: properti Decimal KOSONG bernilai 0 di dalam ekspresi aritmetika.
//
// ⚠️ Kosong-sebagai-nol berlaku untuk ARITMETIKA saja. Pembandingan teks
// (`.PremiOgp == ""`) tetap memakai `Ambil`, sebab di sana kosong dan "0"
// dibedakan - lihat `CountOGPONP_Act` langkah 1-2.
func AngkaTeks(jalur, teks string) (*apd.Decimal, error) {
	s := strings.TrimSpace(teks)
	if s == "" {
		return apd.New(0, 0), nil
	}
	d, err := utils.ParseDecimal(s)
	if err != nil {
		return nil, fmt.Errorf("%w: %s = %q", ErrBukanAngka, jalur, teks)
	}
	return d, nil
}

// Angka membaca satu jalur halaman sebagai desimal (kosong = 0).
func (h *Halaman) Angka(jalur string) (*apd.Decimal, error) {
	return AngkaTeks(jalur, h.Ambil(jalur))
}

// AdalahDesimal meniru `@String.isDouble(x)`: teks itu bilangan.
//
// ⚠️ Teks kosong BUKAN bilangan di sini - itulah sebabnya rule aslinya
// menulis `@String.isDouble(.RiCommOgp) || .RiCommOgp==""` (CountResult1_Act
// langkah 2).
func AdalahDesimal(teks string) bool {
	s := strings.TrimSpace(teks)
	if s == "" {
		return false
	}
	_, err := utils.ParseDecimal(s)
	return err == nil
}
