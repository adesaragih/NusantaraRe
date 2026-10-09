package models

// Untuk apa berkas ini: PEMBANTU DEFINISI LAYAR - kondisi, pembentuk unsur tata, format layout, mode layar, dan daftar
// pilihan kode `associated`. Disalin dari `modul/claimprop/backend/models/layar.go`, bukan diimpor. Definisi section Non
// Prop ada di `layar_*.go`.

import (
	"strings"
	"time"
	"unicode/utf8"
)

// ---------------------------------------------------------------- kondisi

func sama(j, v string) Kondisi   { return func(h *Halaman) bool { return h.Ambil(j) == v } }
func beda(j, v string) Kondisi   { return func(h *Halaman) bool { return h.Ambil(j) != v } }
func terisi(j string) Kondisi    { return func(h *Halaman) bool { return h.Ambil(j) != "" } }
func atau(ks ...Kondisi) Kondisi { return func(h *Halaman) bool { return anyK(h, ks) } }

func anyK(h *Halaman, ks []Kondisi) bool {
	for _, k := range ks {
		if k(h) {
			return true
		}
	}
	return false
}

func dan(ks ...Kondisi) Kondisi {
	return func(h *Halaman) bool {
		for _, k := range ks {
			if !k(h) {
				return false
			}
		}
		return true
	}
}

var (
	isOutstanding  = sama("IsOutstanding", "1")
	isAcceptation  = sama("IsAcceptation", "1")
	isEstimation   = sama("isEstimation", "1")
	isCFS          = sama("IsCFS", "1")
	nonCatastrophe = sama(CD+"StsKatastrofe", "Non-Catastrophe")
	editCatastrope = sama(CD+"EditCatastrope", "true")
	selalu         = func(*Halaman) bool { return true }
	pyNoteKosong   = sama("pyNote", "")
)

func bSama(p, v string) KondisiBaris { return func(_ *Halaman, b Baris) bool { return b[p] == v } }

// bTerisi - medan baris terisi.
func bTerisi(p string) KondisiBaris { return func(_ *Halaman, b Baris) bool { return b[p] != "" } }

// bAtau - salah satu kondisi baris benar.
func bAtau(ks ...KondisiBaris) KondisiBaris {
	return func(h *Halaman, b Baris) bool {
		for _, k := range ks {
			if k(h, b) {
				return true
			}
		}
		return false
	}
}

var (
	bNote     = bSama("Note", "Yes")
	bTerkunci = bSama("CNPFlagOuts", "1")
	hSelalu   = func(*Halaman, Baris) bool { return true }
)

// ---------------------------------------------------------------- pembentuk

func medan(jalur, label, kendali string) Unsur {
	return Unsur{Jenis: JenisMedan, Jalur: jalur, Label: label, Kendali: kendali}
}

func ro(u Unsur) Unsur                { u.HanyaBaca = selalu; return u }
func roJika(u Unsur, k Kondisi) Unsur { u.HanyaBaca = k; return u }
func naJika(u Unsur, k Kondisi) Unsur { u.Nonaktif = k; return u }
func tampil(u Unsur, k Kondisi) Unsur { u.Tampil = k; return u }
func wajibU(u Unsur) Unsur            { u.Wajib = selalu; return u }
func wajibJ(u Unsur, k Kondisi) Unsur { u.Wajib = k; return u }
func aksi(u Unsur, a string) Unsur    { u.Aksi = a; return u }
func sumber(u Unsur, s string) Unsur  { u.Sumber = s; return u }
func tampilIsi(u Unsur) Unsur         { u.Tampil = terisi(u.Jalur); return u }
func label(t string) Unsur            { return Unsur{Jenis: JenisLabel, Label: t} }
func catatan(u Unsur, c string) Unsur { u.Catatan = c; return u }

func bagian(judul string, anak ...Unsur) Unsur {
	return Unsur{Jenis: JenisBagian, Label: judul, Anak: anak}
}

// Letak layout Pega (`pyLayoutOtherFormat`, layout group) - dirender layar; "" = Stacked with labels left.
const (
	LetakDua     = "dua"     // Inline grid double: setiap anak satu sel, dua sel per baris
	LetakSebaris = "sebaris" // Inline / Inline labels left: anak sebaris; label bagian = label baris
	LetakTab     = "tab"     // layout group Tab: anak = bagian berjudul (satu tab per bagian)
	LetakJudul   = "judul"   // kepala layar: Inline grid triple dengan label di sel tengah
	LetakTabel   = "tabel"   // layout bebas berkolom (tanpa format): anak = baris (bagian), baris pertama = judul kolom
)

// Ikon tombol (kelas ikon / gambar tombol di XML); label tetap dikirim sebagai keterangan.
const (
	IkonTambah = "tambah" // pi-plus, webwb/pyWorkActionsAddWork.png
	IkonHapus  = "hapus"  // pi-trash
	IkonUbah   = "ubah"   // pi-pencil, webwb/pyEditIcon.png
	IkonSimpan = "simpan" // pi-check
)

func letak(l string, anak ...Unsur) Unsur { return Unsur{Jenis: JenisBagian, Letak: l, Anak: anak} }

// sebaris - layout Inline; `lbl` = label baris (boleh kosong).
func sebaris(lbl string, anak ...Unsur) Unsur {
	u := letak(LetakSebaris, anak...)
	u.Label = lbl
	return u
}

// dua - layout Inline grid double; bagian tanpa judul di dalamnya = satu kolom Stacked with labels left.
func dua(sel ...Unsur) Unsur { return letak(LetakDua, sel...) }

func ikon(u Unsur, i string) Unsur { u.Ikon = i; return u }

// ModeLayar - penanda MODE layar Pega yang hidup di clipboard tetapi tidak punya kolom di tabel datar (pola Claim Prop):
// ikon Edit / Save Catastrophe (`EditCatastrope`, Catastrope_Sec) dan halaman sementara pop-up komite KomiteCLMNP
// selama terbuka - `CekError.Lolos` (gerbang tombol Send Claim to Committee, dihitung ulang di server saat dikirim),
// tanggal (`TempCommiteClaim.DateOfComitee`) dan inisial (`TreatyExchangeYearly.UserName`). Server mengirimnya di
// `Layar.Mode`, layar mengembalikannya di setiap aksi, server memasangnya SEBELUM tata dihitung.
var ModeLayar = []string{CD + "EditCatastrope", JalurLolosKomite, JalurTanggalKomite, JalurPICKomite}

// Jalur halaman sementara pop-up komite (KomiteCLMNP / CloseClaimNP).
const (
	JalurLolosKomite   = "CekError.Lolos"
	JalurTanggalKomite = "TempCommiteClaim.DateOfComitee"
	JalurPICKomite     = "TreatyExchangeYearly.UserName"
)

// modeSah - nilai kiriman layar yang diterima per penanda ModeLayar.
var modeSah = map[string]func(string) bool{
	CD + "EditCatastrope": benarSalah,
	JalurLolosKomite:      nolSatu,
	JalurTanggalKomite: func(v string) bool {
		_, err := time.Parse("2006-01-02", v)
		return err == nil
	},
	JalurPICKomite: func(v string) bool { return utf8.RuneCountInString(v) <= 200 },
}

func benarSalah(v string) bool { return v == "true" || v == "false" }
func nolSatu(v string) bool    { return v == "0" || v == "1" }

// PasangMode memasang penanda mode kiriman layar; hanya kunci ModeLayar dengan nilai yang sah (`modeSah`).
func PasangMode(h *Halaman, mode map[string]string) {
	for _, k := range ModeLayar {
		if v, ada := mode[k]; ada && modeSah[k](v) {
			h.Setel(k, v)
		}
	}
}

// AmbilMode - penanda mode halaman untuk `Layar.Mode`.
func AmbilMode(h *Halaman) map[string]string {
	out := map[string]string{}
	for _, k := range ModeLayar {
		if v := h.Ambil(k); v != "" {
			out[k] = v
		}
	}
	return out
}

// AksiLihatPolis - tombol View polis (harness ViewPolisNonProp / ViewDetailDeptHeadTreatyIn_UW: bawaan OQ-CNP-13 =
// jendela Modal penuh NB / EDM Treaty In, pola Claim Prop); tidak ada aksi server.
const AksiLihatPolis = "LihatPolis"

func tombol(id, lbl, aksiNama string) Unsur {
	return Unsur{Jenis: JenisTombol, ID: id, Label: lbl, Aksi: aksiNama}
}

// tombolOQ - tombol yang tampil sesuai section tetapi nonaktif (rule tidak diekspor).
func tombolOQ(id, lbl, oq string) Unsur {
	return Unsur{Jenis: JenisTombol, ID: id, Label: lbl, Nonaktif: selalu, Catatan: oq}
}

func kol(prop, judul, kendali string) Unsur {
	return Unsur{Jenis: JenisMedan, Jalur: prop, Label: judul, Kendali: kendali}
}

func kRO(u Unsur) Unsur                  { u.HanyaBacaB = hSelalu; return u }
func kROJ(u Unsur, k KondisiBaris) Unsur { u.HanyaBacaB = k; return u }
func kNA(u Unsur, k KondisiBaris) Unsur  { u.NonaktifB = k; return u }
func kAksi(u Unsur, a string) Unsur      { u.Aksi = a; return u }
func kSumber(u Unsur, s string) Unsur    { u.Sumber = s; return u }
func kTombol(id, lbl, a string, na KondisiBaris) Unsur {
	return Unsur{Jenis: JenisTombol, ID: id, Label: lbl, Aksi: a, NonaktifB: na}
}

func ikonK(u Unsur, i string) Unsur { u.Ikon = i; return u }

// AdaKode - kunci sumber kode `associated`.
func AdaKode(s string) (string, bool) {
	if !strings.HasPrefix(s, AwalanKode) {
		return "", false
	}
	p := strings.TrimPrefix(s, AwalanKode)
	_, ada := KodePilihan[p]
	return p, ada
}
