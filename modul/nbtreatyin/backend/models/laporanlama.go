package models

// Untuk apa berkas ini: LAPORAN PEMUAT DOKUMEN LAMA - tiket 22.
//
// `[keputusan work owner]` K17 (PROMPT-NB-TREATY-IN-PUTARAN-2.md bab 2): medan
// dokumen tanpa kolom TIDAK ditampung di tabel (`T_POLIS_MEDAN_LAIN` dihapus -
// bukan tabel diagram grilling) melainkan di BERKAS CSV per jalankan, di
// folder keluaran yang ditentukan operator.
//
// `[keputusan work owner]` F3 (04-10-2026, PROMPT-NB-TREATY-IN-PUTARAN-3.md
// bab 2): setiap medan dokumen lama tanpa kolom diputuskan per medan - dibuang
// dengan alasan + bukti (`medan_abaikan_lama.json`) - dan berkas CSV itu
// menjadi ARSIP AUDIT PEMUATAN, bukan penampung: kolom `POLIS_ID`, `JALUR`,
// `NILAI`, `KEPUTUSAN` untuk SETIAP medan yang tidak ditulis ke kolom (yang
// dibuang menurut keputusan tertulis dan yang belum diputuskan). AC 59 (RALAT
// F3): yang wajib 0 sebelum pekerjaan dinyatakan selesai adalah medan BELUM
// DIPUTUSKAN. Laporan galat juga berkas, bukan tabel (AC 58; K15: tanggal
// ambigu tidak ditebak, jumlahnya dilaporkan).
//
// Murni: hanya menulis ke io.Writer yang diberikan pemanggil. Ringkasan yang
// dicetak hanya memuat cacah dan pola jalur - NILAI medan (dapat berupa nama
// orang) hanya ada di berkas CSV.

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// Kepala kedua berkas laporan.
var (
	KepalaArsipMedan   = []string{"POLIS_ID", "JALUR", "NILAI", "KEPUTUSAN"}
	KepalaLaporanGalat = []string{"IDPEGA", "NOPOLIS", "JALUR", "NILAI", "SEBAB"}
)

// Nilai kolom KEPUTUSAN arsip medan.
const (
	// KeputusanBelumDiputuskan - medan tanpa kolom dan tanpa keputusan
	// tertulis; menahan selesai (AC 59).
	KeputusanBelumDiputuskan = "BELUM DIPUTUSKAN"
	// awalanDibuang + kunci alasan `medan_abaikan_lama.json`.
	awalanDibuang = "dibuang: "
)

// ErrIDKasusDipakai - ID kasus (pyID) sudah dipakai baris T_GENERAL_POLIS
// lain ber-IDPEGA berbeda; dokumen tidak dimuat, tidak menimpa.
var ErrIDKasusDipakai = errors.New("models: ID kasus sudah dipakai kasus lain (IDPEGA berbeda)")

// ErrDokumenGanda - lebih dari satu baris JSON_POLIS generasi NB untuk kasus
// yang sama dalam satu jalankan. Hanya yang pertama dimuat; sisanya
// dilaporkan, tidak ditebak mana yang benar.
var ErrDokumenGanda = errors.New("models: dokumen generasi NB ganda untuk kasus yang sama")

// JenisGalatSimpan - galat di luar daftar jenis: penyimpanan ditolak basis
// data atau penjaga repository.
const JenisGalatSimpan = "penyimpanan ditolak (basis data / penjaga repository)"

// jenisGalat - label ringkasan per galat penanda.
var jenisGalat = []struct {
	err  error
	nama string
}{
	{ErrDokumenRusak, "dokumen tak terurai"},
	{ErrProdKe, "PRODKE tak terbaca"},
	{ErrIDPega, "IDPEGA tak terbaca"},
	{ErrTanggalAmbigu, "tanggal ambigu - tidak ditebak (K15)"},
	{ErrFormatTanggal, "format tanggal di luar YYYYMMDD / GMT"},
	{ErrNilaiKolom, "nilai tidak sesuai tipe kolom"},
	{ErrNoPolisKosong, "NOPOLIS kosong"},
	{ErrNoPolisBeda, "PolicyNo berbeda dari NOPOLIS"},
	{ErrBentukTidakSah, "bentuk tidak sah untuk jenis proporsinya (AC 31, 33)"},
	{ErrIDKasusDipakai, "ID kasus dipakai kasus lain"},
	{ErrDokumenGanda, "dokumen ganda untuk kasus yang sama"},
}

// JenisGalat - label ringkasan sebuah galat dokumen.
func JenisGalat(err error) string {
	for _, j := range jenisGalat {
		if errors.Is(err, j.err) {
			return j.nama
		}
	}
	return JenisGalatSimpan
}

// RingkasanPemuat - cacah satu jalankan pemuat.
type RingkasanPemuat struct {
	// Tulis - `-jalankan` (true) atau uji-kering (false).
	Tulis bool
	// BarisDibaca - baris JSON_POLIS ber-PRODKE 0 / kosong yang dibaca.
	BarisDibaca int
	// BarisProdKeLain - baris JSON_POLIS ber-PRODKE lain (seluruh lini), tidak dibaca.
	BarisProdKeLain    int
	BukanTreatyIn      int
	GenerasiEndorsemen int
	// Dimuat - ditulis (`-jalankan`) atau siap ditulis (uji-kering).
	Dimuat       int
	SudahDimuat  int
	DokumenGagal int
	// GalatPerJenis - jenis -> cacah sebab (satu dokumen dapat membawa beberapa).
	GalatPerJenis map[string]int
	// TanggalAmbigu - cacah NILAI tanggal ambigu (K15).
	TanggalAmbigu int
	// DiabaikanPerAlasan - teks alasan keputusan tertulis -> cacah medan dibuang.
	DiabaikanPerAlasan map[string]int
	// MedanBelumDiputuskan - baris arsip berkeputusan BELUM DIPUTUSKAN (AC 59
	// RALAT F3: wajib 0).
	MedanBelumDiputuskan   int
	BelumDiputuskanPerPola map[string]int
	// UsulanDisalin - baris SuggestList dokumen lama yang ditulis (atau siap
	// ditulis) ke POOLDATA.HISTORYAKSEPTASIPRODUCTION (F3).
	UsulanDisalin int
	// UsulanTanpaAksesLogin, UsulanTanpaPIC - baris yang disalin dengan
	// AKSES_LOGIN / PIC kosong (NULL): isi baris dokumen tidak memuatnya dan
	// nilainya tidak dikarang (F3).
	UsulanTanpaAksesLogin int
	UsulanTanpaPIC        int
	// UsulanDokumenSudahAda - dokumen yang IDPEGA-nya sudah punya baris riwayat
	// produksi: salinan SuggestList-nya dilewati (penjaga dobel, F3).
	UsulanDokumenSudahAda int
	// AngkaKasusTerbesar - nomor terbesar pyID berawalan `NB-` yang dimuat:
	// SEQ_WORK_POLIS (`IDKasusBerikut`) wajib dimajukan melewatinya.
	AngkaKasusTerbesar int64
}

// Selesai - nol dokumen gagal dan nol medan BELUM DIPUTUSKAN (AC 59 RALAT
// F3). Medan yang dibuang menurut keputusan tertulis tidak menahan selesai.
func (r RingkasanPemuat) Selesai() bool { return r.DokumenGagal == 0 && r.MedanBelumDiputuskan == 0 }

func tulisPeta(b *strings.Builder, m map[string]int) {
	k := make([]string, 0, len(m))
	for n := range m {
		k = append(k, n)
	}
	sort.Strings(k)
	for _, n := range k {
		fmt.Fprintf(b, "    %-90s %d\n", n, m[n])
	}
}

// Teks - ringkasan yang dicetak alat pemuat.
func (r RingkasanPemuat) Teks() string {
	var b strings.Builder
	mode, dimuat, disalin := "UJI-KERING (nol tulisan)", "Siap dimuat", "siap disalin"
	if r.Tulis {
		mode, dimuat, disalin = "JALANKAN (satu transaksi per dokumen)", "Dimuat", "disalin"
	}
	fmt.Fprintf(&b, "Pemuat dokumen lama NB Treaty In - %s\n", mode)
	fmt.Fprintf(&b, "Baris JSON_POLIS PRODKE 0/kosong dibaca          : %d\n", r.BarisDibaca)
	fmt.Fprintf(&b, "Baris JSON_POLIS PRODKE lain, tidak dibaca       : %d  (generasi endorsemen - pemuat EDM tiket 10)\n", r.BarisProdKeLain)
	fmt.Fprintf(&b, "Bukan dokumen Treaty In (lini lain), dilewati    : %d\n", r.BukanTreatyIn)
	if r.GenerasiEndorsemen > 0 {
		fmt.Fprintf(&b, "Generasi endorsemen, dilewati                    : %d\n", r.GenerasiEndorsemen)
	}
	fmt.Fprintf(&b, "%-49s: %d\n", dimuat, r.Dimuat)
	fmt.Fprintf(&b, "Sudah dimuat sebelumnya, dilewati                : %d\n", r.SudahDimuat)
	fmt.Fprintf(&b, "Dokumen gagal (berkas laporan galat)             : %d\n", r.DokumenGagal)
	tulisPeta(&b, r.GalatPerJenis)
	fmt.Fprintf(&b, "Tanggal ambigu, tidak ditebak (K15)              : %d nilai\n", r.TanggalAmbigu)
	fmt.Fprintf(&b, "Catatan SuggestList %-29s: %d baris -> POOLDATA.HISTORYAKSEPTASIPRODUCTION (F3)\n", disalin, r.UsulanDisalin)
	fmt.Fprintf(&b, "    AKSES_LOGIN kosong (baris dokumen tanpa OperatorID, ditulis NULL): %d\n", r.UsulanTanpaAksesLogin)
	fmt.Fprintf(&b, "    PIC kosong (baris dokumen tanpa OperatorName, ditulis NULL)      : %d\n", r.UsulanTanpaPIC)
	fmt.Fprintf(&b, "    dokumen dilewati - IDPEGA sudah punya baris riwayat produksi    : %d\n", r.UsulanDokumenSudahAda)
	fmt.Fprintf(&b, "Medan dibuang menurut keputusan tertulis (arsip CSV, KEPUTUSAN dibuang):\n")
	tulisPeta(&b, r.DiabaikanPerAlasan)
	fmt.Fprintf(&b, "Medan BELUM DIPUTUSKAN (arsip CSV, KEPUTUSAN %s): %d  - WAJIB 0 sebelum pekerjaan dinyatakan selesai (AC 59, F3)\n",
		KeputusanBelumDiputuskan, r.MedanBelumDiputuskan)
	tulisPeta(&b, r.BelumDiputuskanPerPola)
	if r.AngkaKasusTerbesar > 0 {
		fmt.Fprintf(&b, "Nomor kasus terbesar yang dimuat: %s%d - SEQ_WORK_POLIS wajib dimajukan melewatinya sebelum kasus baru dibuat\n",
			AwalanKasus, r.AngkaKasusTerbesar)
	}
	if r.Selesai() {
		b.WriteString("STATUS: SELESAI - nol dokumen gagal, nol medan belum diputuskan\n")
	} else {
		b.WriteString("STATUS: BELUM SELESAI - periksa berkas laporan galat dan medan BELUM DIPUTUSKAN di arsip\n")
	}
	return b.String()
}

// LaporanPemuat - dua berkas CSV satu jalankan beserta ringkasannya.
type LaporanPemuat struct {
	arsip, galat *csv.Writer
	r            RingkasanPemuat
}

// LaporanPemuatBaru menulis kepala kedua berkas. `tulis` = mode `-jalankan`.
func LaporanPemuatBaru(arsip, galat io.Writer, tulis bool) (*LaporanPemuat, error) {
	l := &LaporanPemuat{arsip: csv.NewWriter(arsip), galat: csv.NewWriter(galat), r: RingkasanPemuat{
		Tulis: tulis, GalatPerJenis: map[string]int{}, DiabaikanPerAlasan: map[string]int{}, BelumDiputuskanPerPola: map[string]int{},
	}}
	if err := l.arsip.Write(KepalaArsipMedan); err != nil {
		return nil, err
	}
	if err := l.galat.Write(KepalaLaporanGalat); err != nil {
		return nil, err
	}
	return l, nil
}

// Dibaca mencatat satu baris JSON_POLIS yang dibaca.
func (l *LaporanPemuat) Dibaca() { l.r.BarisDibaca++ }

// ProdKeLain mencatat cacah baris JSON_POLIS yang tidak dibaca.
func (l *LaporanPemuat) ProdKeLain(n int) { l.r.BarisProdKeLain = n }

// Lewat mencatat dokumen di luar lingkup (ErrBukanTreatyIn, ErrGenerasiEndorsemen).
func (l *LaporanPemuat) Lewat(err error) {
	if errors.Is(err, ErrGenerasiEndorsemen) {
		l.r.GenerasiEndorsemen++
		return
	}
	l.r.BukanTreatyIn++
}

// Gagal menulis setiap sebab satu dokumen ke berkas galat (AC 58).
func (l *LaporanPemuat) Gagal(b BarisJSONPolis, g []GalatDokumen) error {
	l.r.DokumenGagal++
	for _, x := range g {
		l.r.GalatPerJenis[JenisGalat(x.Err)]++
		if errors.Is(x.Err, ErrTanggalAmbigu) {
			l.r.TanggalAmbigu++
		}
		if err := l.galat.Write([]string{b.IDPega, b.NoPolis, x.Jalur, x.Nilai, x.Err.Error()}); err != nil {
			return err
		}
	}
	return nil
}

// Berhasil mencatat dokumen yang dimuat (atau siap dimuat), menulis setiap
// medannya yang tidak masuk kolom ke arsip CSV beserta keputusannya (F3, K17,
// AC 57), dan mencatat salinan SuggestList-nya. `usulanDisalin` false =
// penjaga dobel melewati salinan karena IDPEGA sudah punya baris riwayat
// produksi (`repository.SalinUsulanLama`).
func (l *LaporanPemuat) Berhasil(h HasilPecah, usulanDisalin bool) error {
	l.r.Dimuat++
	for a, n := range h.Diabaikan {
		l.r.DiabaikanPerAlasan[a] += n
	}
	for _, m := range h.Arsip {
		keputusan := KeputusanBelumDiputuskan
		if m.Kunci != "" {
			keputusan = awalanDibuang + m.Kunci
		} else {
			l.r.MedanBelumDiputuskan++
			l.r.BelumDiputuskanPerPola[m.Pola]++
		}
		if err := l.arsip.Write([]string{h.ID, m.Jalur, m.Nilai, keputusan}); err != nil {
			return err
		}
	}
	switch {
	case usulanDisalin:
		l.r.UsulanDisalin += len(h.Usulan)
		for _, u := range h.Usulan {
			if u.AksesLogin == "" {
				l.r.UsulanTanpaAksesLogin++
			}
			if u.PIC == "" {
				l.r.UsulanTanpaPIC++
			}
		}
	case len(h.Usulan) > 0:
		l.r.UsulanDokumenSudahAda++
	}
	if strings.HasPrefix(h.ID, AwalanKasus) {
		if n, err := strconv.ParseInt(strings.TrimPrefix(h.ID, AwalanKasus), 10, 64); err == nil && n > l.r.AngkaKasusTerbesar {
			l.r.AngkaKasusTerbesar = n
		}
	}
	return nil
}

// SudahDimuat mencatat dokumen yang barisnya sudah ada dari jalankan sebelumnya.
func (l *LaporanPemuat) SudahDimuat() { l.r.SudahDimuat++ }

// Tutup mengosongkan penyangga kedua berkas.
func (l *LaporanPemuat) Tutup() error {
	l.arsip.Flush()
	l.galat.Flush()
	return errors.Join(l.arsip.Error(), l.galat.Error())
}

// Ringkasan - salinan cacah saat ini.
func (l *LaporanPemuat) Ringkasan() RingkasanPemuat { return l.r }
