package models

// Untuk apa berkas ini: PEMBANTU DEFINISI LAYAR - kondisi, pembentuk unsur tata, format layout, mode layar, dan daftar
// pilihan kode `associated`. Disalin dari `modul/claimnonprop/backend/models/layar_dasar.go` (asal Claim Prop), bukan
// diimpor. Definisi section Claim Fac In ada di `layar_*.go`.

import (
	"strings"
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
	isEstimasi     = sama(JalurIsEstimation, "1")
	isAdjustment   = sama(JalurIsAdjustment, "1")
	adaAkseptasi   = sama(JalurIsAnyAccept, "1")
	isCFS          = sama(JalurIsCFS, "1")
	nonCatastrophe = sama(CD+"StsKatastrofe", "Non-Catastrophe")
	editCatastrope = sama(CD+"EditCatastrope", "true")
	selalu         = func(*Halaman) bool { return true }
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

var hSelalu = func(*Halaman, Baris) bool { return true }

// ---------------------------------------------------------------- pembentuk

func medan(jalur, label, kendali string) Unsur {
	return Unsur{Jenis: JenisMedan, Jalur: jalur, Label: label, Kendali: kendali}
}

// Pembungkus kondisi MENGGABUNG dengan kondisi yang sudah terpasang, tidak menimpanya: tampil = DAN (pembungkus luar =
// container ber-visible-when atas unsur yang punya visible-when sendiri), hanya-baca / nonaktif = ATAU (tombolOQ tetap
// nonaktif walau dibungkus naJika). Temuan cek layar 10-10-2026: menimpa membuat "Yes" Reject Claim aktif dan grid
// coverage Marine tampil di semua lini.
func ro(u Unsur) Unsur                { u.HanyaBaca = selalu; return u }
func roJika(u Unsur, k Kondisi) Unsur { u.HanyaBaca = gabungAtau(u.HanyaBaca, k); return u }
func naJika(u Unsur, k Kondisi) Unsur { u.Nonaktif = gabungAtau(u.Nonaktif, k); return u }
func tampil(u Unsur, k Kondisi) Unsur { u.Tampil = gabungDan(u.Tampil, k); return u }

func gabungAtau(lama, k Kondisi) Kondisi {
	if lama == nil {
		return k
	}
	if k == nil {
		return lama
	}
	return atau(lama, k)
}

func gabungDan(lama, k Kondisi) Kondisi {
	if lama == nil {
		return k
	}
	if k == nil {
		return lama
	}
	return dan(lama, k)
}
func wajibU(u Unsur) Unsur            { u.Wajib = selalu; return u }
func wajibJ(u Unsur, k Kondisi) Unsur { u.Wajib = k; return u }
func aksi(u Unsur, a string) Unsur    { u.Aksi = a; return u }
func sumber(u Unsur, s string) Unsur  { u.Sumber = s; return u }
func tampilIsi(u Unsur) Unsur         { return tampil(u, terisi(u.Jalur)) }
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

// ModeLayar - penanda MODE layar Pega yang hidup di clipboard (halaman requestor / properti tanpa kolom) tetapi tidak
// disimpan (pola Claim Prop): ikon Edit / Save Catastrophe (`EditCatastrope`), lokasi sementara (`TempLocaion.CARI1`,
// SetTempLocation_Act), pop-up Choose Polis (`SearchType`, `SearchName`, `InputParam.CARI4` = detail polis tampil),
// pop-up Print PLA / DLA (`ProtectPrint.CARI50` kirim surel, `Email.Message`), penanda SendPICProtect_Act
// (`Protect.CARI1/CARI2` - hanya menampilkan tombol; penyerahan menghitung ulang di server), medan tampilan pop-up
// Comittee / ClaimComiteeReject (`TempCommiteClaim.*`, `TreatyExchangeYearly.UserName`), dan pilihan Close Without
// Payment (`TempCommiteClaim.AllocationShareSalvage`). `CekLimit.CARI1` BUKAN mode layar:
// halaman requestor itu tak pernah bernilai 1 (OQ-CFI-24), jadi tidak diterima dari layar. Server mengirimnya di `Layar.Mode`, layar mengembalikannya di setiap
// aksi, server memasangnya SEBELUM tata dihitung. Pilihan calon objek layar Register
// (`TempClaimData.ClaimData.ObjectList(n).Selected`) dikirim sebagai isian grid dan dibawa `BawaSementara`.
var ModeLayar = []string{CD + "EditCatastrope", JalurLokasiSementara, JalurJenisCari, JalurTeksCari, JalurDetailPolis,
	JalurKirimSurel, JalurPesanSurel, JalurProtect1, JalurProtect2, JalurTanggalKomite, JalurInisialKomite,
	JalurAlokasiSalvage, JalurTKInisial, JalurTKTipe}

// Jalur pop-up Choose Polis (ViewPolis).
const (
	JalurJenisCari   = "SearchType"
	JalurTeksCari    = "SearchName"
	JalurDetailPolis = "InputParam.CARI4" // SetInputParam_Act / CopyNB_Act: InputInwardFacultativeDtl tampil
)

// modeSah - nilai kiriman layar yang diterima per penanda ModeLayar.
var modeSah = map[string]func(string) bool{
	CD + "EditCatastrope": benarSalah,
	JalurLokasiSementara:  teksPendek(4000),
	JalurJenisCari:        func(v string) bool { return v == CariNoPolis || v == CariCeding || v == CariInsured || v == CariQQ },
	JalurTeksCari:         teksPendek(200),
	JalurDetailPolis:      benarSalah,
	JalurKirimSurel:       benarSalah,
	JalurPesanSurel:       teksPendek(20000),
	JalurProtect1:         modeNolSatu,
	JalurProtect2:         modeNolSatu,
	JalurTanggalKomite:    teksPendek(32),
	JalurInisialKomite:    teksPendek(128),
	JalurAlokasiSalvage:   benarSalah,
	JalurTKInisial:        teksPendek(128),
	JalurTKTipe:           func(v string) bool { return v == "1" || v == "5" },
}

func benarSalah(v string) bool { return v == "true" || v == "false" }

func modeNolSatu(v string) bool { return v == "0" || v == "1" }

func teksPendek(n int) func(string) bool {
	return func(v string) bool { return utf8.RuneCountInString(v) <= n }
}

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

// BawaSementara memindahkan keadaan layar sementara dari halaman aksi ke halaman yang dimuat ulang sesudahnya:
// penanda ModeLayar dan pilihan calon objek (menurut `KunciPolis`).
func BawaSementara(dari, ke *Halaman) {
	for _, k := range ModeLayar {
		if v := dari.Ambil(k); v != "" {
			ke.Setel(k, v)
		}
	}
	pilih := map[string]string{}
	for _, b := range dari.AmbilDaftar(DaftarCalon) {
		pilih[b[PropKunciPolis]] = b[PropDipilih]
	}
	for _, b := range ke.AmbilDaftar(DaftarCalon) {
		if v, ada := pilih[b[PropKunciPolis]]; ada {
			b[PropDipilih] = v
		}
	}
}

// AksiLihatPolis - tombol / tautan yang membuka pop-up halaman polis tanpa aksi server.
const AksiLihatPolis = "LihatPolis"

func tombol(id, lbl, aksiNama string) Unsur {
	return Unsur{Jenis: JenisTombol, ID: id, Label: lbl, Aksi: aksiNama}
}

// tombolOQ - tombol yang tampil sesuai section tetapi nonaktif (rule tidak diekspor / layanan belum disetujui).
func tombolOQ(id, lbl, oq string) Unsur {
	return Unsur{Jenis: JenisTombol, ID: id, Label: lbl, Nonaktif: selalu, Catatan: oq}
}

func kol(prop, judul, kendali string) Unsur {
	return Unsur{Jenis: JenisMedan, Jalur: prop, Label: judul, Kendali: kendali}
}

func kRO(u Unsur) Unsur                  { u.HanyaBacaB = hSelalu; return u }
func kROJ(u Unsur, k KondisiBaris) Unsur { u.HanyaBacaB = gabungAtauB(u.HanyaBacaB, k); return u }
func kNA(u Unsur, k KondisiBaris) Unsur  { u.NonaktifB = gabungAtauB(u.NonaktifB, k); return u }
func kTampil(u Unsur, k KondisiBaris) Unsur {
	if u.TampilB != nil && k != nil {
		lama := u.TampilB
		u.TampilB = func(h *Halaman, b Baris) bool { return lama(h, b) && k(h, b) }
		return u
	}
	if k != nil {
		u.TampilB = k
	}
	return u
}

// gabungAtauB - pasangan gabungAtau untuk kondisi baris grid.
func gabungAtauB(lama, k KondisiBaris) KondisiBaris {
	if lama == nil {
		return k
	}
	if k == nil {
		return lama
	}
	return func(h *Halaman, b Baris) bool { return lama(h, b) || k(h, b) }
}
func kAksi(u Unsur, a string) Unsur   { u.Aksi = a; return u }
func kSumber(u Unsur, s string) Unsur { u.Sumber = s; return u }
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
