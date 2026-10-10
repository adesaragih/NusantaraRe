package models

// Untuk apa berkas ini: MESIN TATA LETAK. Section Pega menyatakan per sel: tampil (`pyVisible`/`pyCondition`),
// hanya-baca (`pyReadOnlyCondition`), nonaktif (`pyDisabledWhen`), wajib, dan aksi (`refresh` + Activity). Di sini
// ekspresi itu ditulis sebagai fungsi Go atas halaman TERSIMPAN, dievaluasi di server, lalu dikirim ke layar sebagai
// pohon `Tata` - layar hanya merender. Dari pohon yang sama server menurunkan:
//
//   - JALUR yang boleh diubah (`MedanTerbuka`) - kiriman layar untuk jalur lain DIABAIKAN (layar tidak pernah menulis
//     medan terkunci);
//   - AKSI yang terbuka (`AksiTerbuka`) - aksi tombol / sel yang tidak tampil atau nonaktif ditolak.
//
// Definisi section ada di `layar*.go`. Disalin dari `modul/claimnonprop/backend/models/tata.go`, bukan diimpor.
//
// Tambahan Claim Fac In: grid `masterDetail` / expandPane (ObjectList -> ObjectItemList -> Estimasi / Adjustment) memakai
// `Rincian` = prefiks PANEL. Baris ke-n grid berjalur `daftar` membuka panel berkunci `KunciPanel(prefiks, daftar, n)`
// = "prefiks:daftar(n)" - jalur baris itu sendiri, sehingga sarang berapa pun dalamnya beralamat tanpa protokol khusus.

import (
	"sort"
	"strconv"
	"strings"
)

// Jenis unsur tata letak.
const (
	JenisBagian = "bagian" // Layout ber-judul / blok
	JenisMedan  = "medan"  // sel properti
	JenisLabel  = "label"  // LABEL
	JenisTombol = "tombol" // TOMBOL / pxIcon
	JenisGrid   = "grid"   // GRID PageList
)

// Kendali medan.
const (
	KTeks     = "teks"          // pxTextInput
	KArea     = "area"          // pxTextArea
	KAngka    = "angka"         // pxNumber / Decimal
	KTanggal  = "tanggal"       // pxDateTime tanggal saja
	KWaktu    = "tanggal-waktu" // pxDateTime tanggal + jam
	KPilih    = "pilih"         // pxDropdown
	KCentang  = "centang"       // pxCheckbox
	KRadio    = "radio"         // pxRadioButtons
	KTampil   = "tampil"        // pxDisplayText / kontrol bawaan hanya-baca
	KOtomatis = "otomatis"      // pxAutoComplete
	KTelepon  = "telepon"       // nomor telepon: hanya angka, tanpa pemisah ribuan (keputusan work owner 08-10-2026)
)

// Kondisi atas halaman; nil = selalu benar (tampil) atau selalu salah (hanya-baca / nonaktif / wajib) menurut
// penggunaannya.
type Kondisi func(h *Halaman) bool

// KondisiBaris atas satu baris grid.
type KondisiBaris func(h *Halaman, b Baris) bool

// Unsur - satu definisi sel / layout / tombol / grid.
type Unsur struct {
	Jenis   string
	ID      string // pengenal stabil tombol (dipakai permintaan aksi)
	Label   string
	Jalur   string // medan: jalur halaman (atau nama properti anggota untuk kolom grid); grid: jalur daftar
	Kendali string
	Sumber  string // kunci daftar pilihan
	// Tampilan - jalur teks yang DITAMPILKAN untuk medan ber-sumber; nilai medan (Jalur) tetap yang disimpan. Consultant /
	// Adjuster: pilih dan tampil nama, simpan ID (work owner 08-10-2026).
	Tampilan string
	Aksi     string // aksi saat ubah / klik ("" = postValue saja)
	Catatan  string // OQ / alasan nonaktif

	Tampil    Kondisi
	HanyaBaca Kondisi
	Nonaktif  Kondisi
	Wajib     Kondisi

	// Kolom grid: kondisi per baris.
	TampilB    KondisiBaris
	HanyaBacaB KondisiBaris
	NonaktifB  KondisiBaris

	Anak   []Unsur // bagian
	Kolom  []Unsur // grid: templat sel baris
	Kaki   []Unsur // grid: area kaki (total)
	Tambah *Unsur  // grid: tombol "Add" di header
	// Bernomor - grid bernomor baris.
	Bernomor bool
	// Letak - format layout Pega (`pyLayoutOtherFormat` / layout group) untuk layar; "" = Stacked with labels left.
	Letak string
	// Ikon - gambar / kelas ikon tombol tanpa label (pi-plus, pi-trash, pi-pencil, pi-check, pyWorkActionsAddWork).
	Ikon string
	// PerHalaman - paging grid (`pyGridPaginator`); 0 = tanpa paging.
	PerHalaman int
	// Rincian - grid masterDetail: prefiks panel baris (`KunciPanel`); "" = baris tidak dapat dibuka.
	Rincian string
}

// Tata - unsur sesudah dievaluasi, dikirim ke layar.
type Tata struct {
	Jenis     string      `json:"jenis"`
	ID        string      `json:"id,omitempty"`
	Label     string      `json:"label,omitempty"`
	Jalur     string      `json:"jalur,omitempty"`
	Kendali   string      `json:"kendali,omitempty"`
	Sumber    string      `json:"sumber,omitempty"`
	Tampilan  string      `json:"tampilan,omitempty"`
	Aksi      string      `json:"aksi,omitempty"`
	Catatan   string      `json:"catatan,omitempty"`
	HanyaBaca bool        `json:"hanyaBaca,omitempty"`
	Nonaktif  bool        `json:"nonaktif,omitempty"`
	Wajib     bool        `json:"wajib,omitempty"`
	Anak      []Tata      `json:"anak,omitempty"`
	Kolom     []Tata      `json:"kolom,omitempty"`
	Baris     [][]SelTata `json:"baris,omitempty"`
	Kaki      []Tata      `json:"kaki,omitempty"`
	Tambah    *Tata       `json:"tambah,omitempty"`
	Bernomor  bool        `json:"bernomor,omitempty"`
	Letak     string      `json:"letak,omitempty"`
	Ikon      string      `json:"ikon,omitempty"`
	// PerHalaman - jumlah baris per halaman grid.
	PerHalaman int `json:"perHalaman,omitempty"`
	// Rincian - prefiks panel baris grid masterDetail (layar membuka `Layar.Panel[KunciPanel(...)]`).
	Rincian string `json:"rincian,omitempty"`
}

// KunciPanel - kunci panel baris ke-n grid `daftar` ("est:ClaimData.ObjectList(1)").
func KunciPanel(prefiks, daftar string, n int) string {
	return prefiks + ":" + JalurBaris(daftar, n)
}

// JalurBaris - jalur halaman baris ke-n sebuah daftar ("ClaimData.ObjectList(1)").
func JalurBaris(daftar string, n int) string { return daftar + "(" + strconv.Itoa(n) + ")" }

// SelTata - keadaan satu sel grid pada satu baris (sejajar `Tata.Kolom`).
type SelTata struct {
	Tampil    bool `json:"tampil"`
	HanyaBaca bool `json:"hanyaBaca,omitempty"`
	Nonaktif  bool `json:"nonaktif,omitempty"`
}

func ya(k Kondisi, h *Halaman, kosong bool) bool {
	if k == nil {
		return kosong
	}
	return k(h)
}

func yaB(k KondisiBaris, h *Halaman, b Baris, kosong bool) bool {
	if k == nil {
		return kosong
	}
	return k(h, b)
}

// Evaluasi mengubah definisi menjadi tata. `kunci` = layar terkunci seluruhnya (kasus tertutup, bukan pemegang,
// ReadOnly=Yes harness): semua medan hanya-baca, semua tombol nonaktif.
func Evaluasi(h *Halaman, defs []Unsur, kunci bool) []Tata {
	var out []Tata
	for _, u := range defs {
		if !ya(u.Tampil, h, true) {
			continue
		}
		t := Tata{Jenis: u.Jenis, ID: u.ID, Label: u.Label, Jalur: u.Jalur, Kendali: u.Kendali, Sumber: u.Sumber,
			Tampilan: u.Tampilan, Aksi: u.Aksi, Catatan: u.Catatan, Bernomor: u.Bernomor, Letak: u.Letak, Ikon: u.Ikon,
			PerHalaman: u.PerHalaman, Rincian: u.Rincian}
		t.HanyaBaca = kunci || ya(u.HanyaBaca, h, false)
		t.Nonaktif = kunci || ya(u.Nonaktif, h, false)
		t.Wajib = ya(u.Wajib, h, false)
		switch u.Jenis {
		case JenisBagian:
			t.Anak = Evaluasi(h, u.Anak, kunci)
		case JenisGrid:
			t.Kolom = Evaluasi(HalamanBaru(), kolomTanpaKondisi(u.Kolom), false)
			for i := range t.Kolom {
				t.Kolom[i].HanyaBaca, t.Kolom[i].Nonaktif = false, false
			}
			for _, b := range h.AmbilDaftar(u.Jalur) {
				var sel []SelTata
				for _, k := range u.Kolom {
					sel = append(sel, SelTata{
						Tampil:    yaB(k.TampilB, h, b, true),
						HanyaBaca: kunci || yaB(k.HanyaBacaB, h, b, false),
						Nonaktif:  kunci || yaB(k.NonaktifB, h, b, false),
					})
				}
				t.Baris = append(t.Baris, sel)
			}
			t.Kaki = Evaluasi(h, u.Kaki, kunci)
			if u.Tambah != nil {
				if tt := Evaluasi(h, []Unsur{*u.Tambah}, kunci); len(tt) == 1 {
					t.Tambah = &tt[0]
				}
			}
		}
		out = append(out, t)
	}
	return out
}

// kolomTanpaKondisi - templat kolom untuk kepala grid (kondisi baris dievaluasi per baris).
func kolomTanpaKondisi(k []Unsur) []Unsur {
	out := make([]Unsur, len(k))
	for i, u := range k {
		out[i] = Unsur{Jenis: u.Jenis, ID: u.ID, Label: u.Label, Jalur: u.Jalur, Kendali: u.Kendali, Sumber: u.Sumber,
			Aksi: u.Aksi, Catatan: u.Catatan, Ikon: u.Ikon}
	}
	return out
}

// MedanTerbuka - jalur yang boleh diubah layar: medan tampil, tidak hanya-baca, tidak nonaktif, berkendali isian;
// sel grid yang terbuka pada baris itu (`daftar(n).prop`).
func MedanTerbuka(ts []Tata) map[string]bool {
	out := map[string]bool{}
	var jalan func(ts []Tata)
	jalan = func(ts []Tata) {
		for _, t := range ts {
			switch t.Jenis {
			case JenisMedan:
				if !t.HanyaBaca && !t.Nonaktif && t.Kendali != KTampil && t.Jalur != "" {
					out[t.Jalur] = true
				}
			case JenisBagian:
				jalan(t.Anak)
			case JenisGrid:
				for i, sel := range t.Baris {
					for j, s := range sel {
						k := t.Kolom[j]
						if k.Jenis == JenisMedan && s.Tampil && !s.HanyaBaca && !s.Nonaktif && k.Kendali != KTampil {
							out[JalurAnak(t.Jalur, i+1, k.Jalur)] = true
						}
					}
				}
				jalan(t.Kaki)
			}
		}
	}
	jalan(ts)
	return out
}

// AksiTerbuka menjawab: aksi `aksi` (tombol ber-ID, aksi ubah medan, atau aksi sel grid baris `indeks`) tersedia di
// tata yang dievaluasi. Aksi ubah medan menuntut medannya terbuka.
func AksiTerbuka(ts []Tata, aksi string, indeks int) bool {
	var cari func(ts []Tata) bool
	cari = func(ts []Tata) bool {
		for _, t := range ts {
			switch t.Jenis {
			case JenisTombol:
				if (t.ID == aksi || t.Aksi == aksi) && !t.Nonaktif {
					return true
				}
			case JenisMedan:
				if t.Aksi == aksi && !t.HanyaBaca && !t.Nonaktif {
					return true
				}
			case JenisBagian:
				if cari(t.Anak) {
					return true
				}
			case JenisGrid:
				if t.Tambah != nil && (t.Tambah.ID == aksi || t.Tambah.Aksi == aksi) && !t.Tambah.Nonaktif {
					return true
				}
				if indeks >= 1 && indeks <= len(t.Baris) {
					sel := t.Baris[indeks-1]
					for j, k := range t.Kolom {
						if (k.Aksi == aksi || k.ID == aksi) && sel[j].Tampil && !sel[j].Nonaktif &&
							(k.Jenis == JenisTombol || !sel[j].HanyaBaca) {
							return true
						}
					}
				}
				if cari(t.Kaki) {
					return true
				}
			}
		}
		return false
	}
	return cari(ts)
}

// ---------------------------------------------------------------- jalur

// pecahJalurBaris memecah `daftar(n).prop` (indeks terakhir) - prop boleh bertitik (`DataCommitteeTreaty.Remarks`).
func pecahJalurBaris(j string) (daftar string, n int, prop string, ok bool) {
	i := strings.LastIndex(j, ").")
	if i < 0 {
		return "", 0, "", false
	}
	awal := j[:i]
	b := strings.LastIndex(awal, "(")
	if b < 0 {
		return "", 0, "", false
	}
	n, err := strconv.Atoi(awal[b+1:])
	if err != nil {
		return "", 0, "", false
	}
	return awal[:b], n, j[i+2:], true
}

// AmbilJalur membaca nilai jalur halaman ATAU sel baris (`daftar(n).prop`).
func AmbilJalur(h *Halaman, j string) string {
	if d, n, p, ok := pecahJalurBaris(j); ok {
		rows := h.AmbilDaftar(d)
		if n >= 1 && n <= len(rows) {
			return rows[n-1][p]
		}
		return ""
	}
	return h.Ambil(j)
}

// SetelJalur menulis nilai jalur halaman ATAU sel baris yang sudah ada.
func SetelJalur(h *Halaman, j, v string) bool {
	if d, n, p, ok := pecahJalurBaris(j); ok {
		rows := h.AmbilDaftar(d)
		if n < 1 || n > len(rows) {
			return false
		}
		rows[n-1][p] = v
		return true
	}
	h.Setel(j, v)
	return true
}

// GabungMasukan menulis kiriman layar ke halaman tersimpan HANYA untuk jalur yang terbuka. Mengembalikan jalur yang
// berubah, urut.
func GabungMasukan(h *Halaman, masuk map[string]string, terbuka map[string]bool) []string {
	var ubah []string
	for j, v := range masuk {
		if !terbuka[j] {
			continue
		}
		if AmbilJalur(h, j) == v {
			continue
		}
		if SetelJalur(h, j, v) {
			ubah = append(ubah, j)
		}
	}
	sort.Strings(ubah)
	return ubah
}

// WajibKosong - label medan wajib yang tampil, terbuka, dan kosong (validasi klien Pega `pyClientValidation=true`
// sebelum submit). Grid: sel wajib per baris.
func WajibKosong(h *Halaman, ts []Tata, defs []Unsur) []string {
	var out []string
	var jalan func(ts []Tata)
	jalan = func(ts []Tata) {
		for _, t := range ts {
			switch t.Jenis {
			case JenisMedan:
				if t.Wajib && !t.HanyaBaca && !t.Nonaktif && strings.TrimSpace(AmbilJalur(h, t.Jalur)) == "" {
					out = append(out, t.Label)
				}
			case JenisBagian:
				jalan(t.Anak)
			}
		}
	}
	jalan(ts)
	return out
}
