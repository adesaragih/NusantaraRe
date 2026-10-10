// Package models memegang bentuk data dan aturan murni modul Claim Prop (`claimprop`, klaim treaty inward
// proporsional): halaman kerja, katalog kolom, port activity Pega (rumus uang, validasi, penomoran), dan tata letak
// layar. Tidak satu pun berkas di sini menyentuh basis data atau HTTP.
//
// Untuk apa berkas ini: HALAMAN KERJA. Pega menyimpan seluruh isi klaim di `pyWorkPage` beserta halaman anaknya
// (`ClaimData`, `TreatyInMaster`, `OfferFacIn`), dan setiap activity membaca/menulis properti lewat jalurnya -
// `.ClaimData.DateOfLoss`, `pyWorkPage.ClaimData.EstimationList(1).GrossEstimationPct`. Halaman di sini meniru bentuk
// itu: nilai TEKS per jalur relatif terhadap `pyWorkPage`, ditambah daftar baris untuk PageList. Dengan begitu setiap
// port activity dapat dibaca berdampingan dengan langkah aslinya di `Activity/*.xml`.
//
// Pola ini DISALIN dari modul NB Treaty In (`models/halaman.go`), bukan diimpor - modul tidak pernah mengimpor
// modul lain.
//
// ⛔ Nilai disimpan sebagai TEKS, persis seperti di clipboard Pega. Konversi ke desimal terjadi di titik pakai
// (`Angka`), dan kembali ke teks lewat `utils.FormatDecimal` - tidak pernah lewat float (ADR-0003, AC 18).
// Konversi ke tipe kolom Oracle terjadi sekali, di repository.
package models

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// Halaman adalah satu halaman kerja: nilai teks per jalur dan daftar baris.
//
// Jalur ditulis RELATIF terhadap `pyWorkPage` tanpa titik di depan: `ClaimData.DateOfLoss`, `IsOutstanding`,
// `TreatyInMaster.RNMShareP`. Daftar bersarang memakai indeks Pega berbasis satu:
// `ClaimData.AdjustmentList(2).SpreadingAdjustment` (`JalurAnak`).
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

// pastikan menyiapkan peta yang nil - halaman hasil urai JSON dapat datang tanpa salah satunya.
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

// Ambil membaca nilai teks satu jalur. Jalur yang tidak pernah ditulis bernilai "" - sama dengan properti Pega yang
// belum diisi.
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

// HapusAwalan meniru `Page-Remove` atas halaman tunggal: seluruh jalur berawalan `awalan.` dibuang.
func (h *Halaman) HapusAwalan(awalan string) {
	h.pastikan()
	for k := range h.Nilai {
		if strings.HasPrefix(k, awalan+".") {
			delete(h.Nilai, k)
		}
	}
	for k := range h.Daftar {
		if strings.HasPrefix(k, awalan+".") {
			delete(h.Daftar, k)
		}
	}
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
	if len(b) == 0 {
		delete(h.Daftar, jalur)
		return
	}
	h.Daftar[jalur] = b
}

// TambahBaris meniru append ke PageList (`Page-New` + `.(<APPEND>)`), mengembalikan indeks Pega baris baru (1..n).
func (h *Halaman) TambahBaris(jalur string, b Baris) int {
	h.pastikan()
	h.Daftar[jalur] = append(h.Daftar[jalur], b)
	return len(h.Daftar[jalur])
}

// HapusBaris meniru `Page-Remove` atas `jalur(n)` (n berbasis satu). Daftar bersarang baris-baris sesudahnya ikut
// bergeser satu indeks, sama dengan PageList Pega.
func (h *Halaman) HapusBaris(jalur string, n int) {
	d := h.AmbilDaftar(jalur)
	if n < 1 || n > len(d) {
		return
	}
	awalan := jalur + "("
	geser := map[string][]Baris{}
	for k, v := range h.Daftar {
		if !strings.HasPrefix(k, awalan) {
			continue
		}
		i, sisa, ok := indeksAnak(k[len(awalan):])
		if !ok {
			continue
		}
		delete(h.Daftar, k)
		switch {
		case i < n:
			geser[k] = v
		case i > n:
			geser[JalurAnak(jalur, i-1, sisa)] = v
		}
	}
	for k, v := range geser {
		h.Daftar[k] = v
	}
	h.SetelDaftar(jalur, append(append([]Baris{}, d[:n-1]...), d[n:]...))
}

// indeksAnak membaca "<n>).<anak>" menjadi n dan anak.
func indeksAnak(s string) (int, string, bool) {
	tutup := strings.Index(s, ").")
	if tutup < 1 {
		return 0, "", false
	}
	var n int
	if _, err := fmt.Sscanf(s[:tutup], "%d", &n); err != nil {
		return 0, "", false
	}
	return n, s[tutup+2:], true
}

// JalurAnak membentuk jalur daftar bersarang `induk(n).anak` (n berbasis satu, sama dengan Pega).
func JalurAnak(induk string, n int, anak string) string {
	return fmt.Sprintf("%s(%d).%s", induk, n, anak)
}

// TambahPesan meniru `Property-Set-Messages` / `Page-Set-Messages`: pesan galat ditempel pada satu jalur (atau ""
// untuk pesan halaman).
func (h *Halaman) TambahPesan(jalur, pesan string) {
	h.pastikan()
	for _, p := range h.Pesan[jalur] {
		if p == pesan {
			return
		}
	}
	h.Pesan[jalur] = append(h.Pesan[jalur], pesan)
}

// BersihkanPesan meniru `Page-Clear-Messages`.
func (h *Halaman) BersihkanPesan() {
	h.pastikan()
	h.Pesan = map[string][]string{}
}

// AdaPesan - ada setidaknya satu pesan galat terpasang (`@hasMessages(pyWorkPage)`).
func (h *Halaman) AdaPesan() bool {
	if h == nil {
		return false
	}
	for _, p := range h.Pesan {
		if len(p) > 0 {
			return true
		}
	}
	return false
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
		out = append(out, h.Pesan[k]...)
	}
	return out
}

// Salin membuat salinan dalam - activity boleh menulis tanpa mengubah halaman pemanggil.
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
			baru[i] = b.Salin()
		}
		s.Daftar[k] = baru
	}
	for k, p := range h.Pesan {
		s.Pesan[k] = append([]string(nil), p...)
	}
	return s
}

// Salin menyalin satu baris.
func (b Baris) Salin() Baris {
	nb := Baris{}
	for k, v := range b {
		nb[k] = v
	}
	return nb
}

// SalinDaftar menyalin satu PageList baris demi baris (`Page-Copy` daftar).
func SalinDaftar(d []Baris) []Baris {
	out := make([]Baris, len(d))
	for i, b := range d {
		out[i] = b.Salin()
	}
	return out
}

// ---------------------------------------------------------------- angka

// ErrBukanAngka - nilai properti bukan bilangan desimal. Pega melempar galat ekspresi di titik yang sama; sistem
// baru mengembalikannya sebagai galat yang menyebut jalurnya.
var ErrBukanAngka = errors.New("models: nilai properti bukan angka")

// AngkaTeks membaca teks properti sebagai desimal dengan semantik ekspresi Pega: properti Decimal KOSONG bernilai 0
// di dalam ekspresi aritmetika (`@toDecimal("")` = 0).
//
// ⚠️ Kosong-sebagai-nol berlaku untuk ARITMETIKA saja. Pembandingan teks (`.SharePercentage == ""`) tetap memakai
// `Ambil`, sebab di sana kosong dan "0" dibedakan.
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

// Angka membaca satu properti baris sebagai desimal (kosong = 0).
func (b Baris) Angka(prop string) (*apd.Decimal, error) {
	return AngkaTeks("."+prop, b[prop])
}

// AdalahDesimal meniru `@String.isDouble(x)`: teks itu bilangan. Teks kosong BUKAN bilangan.
func AdalahDesimal(teks string) bool {
	s := strings.TrimSpace(teks)
	if s == "" {
		return false
	}
	_, err := utils.ParseDecimal(s)
	return err == nil
}
