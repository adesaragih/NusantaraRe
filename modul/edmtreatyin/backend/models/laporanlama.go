package models

// Untuk apa berkas ini: LAPORAN PEMUAT DOKUMEN LAMA ENDORSEMEN - tiket EDM 10. Asal: adaptasi
// `modul/nbtreatyin/backend/models/laporanlama.go` (tiket NB 22).
//
// `[keputusan work owner]` K17 / F3 (NB, 04-10-2026): medan dokumen tanpa kolom TIDAK ditampung di tabel melainkan
// di BERKAS CSV per jalankan - ARSIP AUDIT PEMUATAN (`POLIS_ID`, `JALUR`, `NILAI`, `KEPUTUSAN`) untuk SETIAP medan
// yang tidak ditulis ke kolom, yang dibuang menurut keputusan tertulis maupun yang BELUM DIPUTUSKAN (wajib 0
// sebelum pekerjaan dinyatakan selesai). Laporan galat juga berkas, bukan tabel; tanggal ambigu tidak ditebak (K15).
// EDM menambah kolom PRODKE di berkas galat (generasi dimuat menurut nomornya) dan cacah penanda migrasi (tiket 09).
//
// Murni: hanya menulis ke io.Writer pemanggil. Ringkasan yang dicetak hanya memuat cacah dan pola jalur - NILAI
// medan (dapat berupa nama orang) hanya ada di berkas CSV.

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
	KepalaLaporanGalat = []string{"IDPEGA", "NOPOLIS", "PRODKE", "JALUR", "NILAI", "SEBAB"}
)

// KeputusanBelumDiputuskan - nilai kolom KEPUTUSAN arsip medan tanpa kolom dan tanpa keputusan tertulis.
const KeputusanBelumDiputuskan = "BELUM DIPUTUSKAN"

// pmAwalanDibuang + kunci alasan `medan_abaikan_lama.json`.
const pmAwalanDibuang = "dibuang: "

// JenisGalatSimpan - galat di luar daftar jenis: penyimpanan ditolak basis data atau penjaga repository.
const JenisGalatSimpan = "penyimpanan ditolak (basis data / penjaga repository)"

// pmJenisGalat - label ringkasan per galat penanda.
var pmJenisGalat = []struct {
	err  error
	nama string
}{
	{ErrDokumenRusak, "dokumen tak terurai"},
	{ErrProdKe, "PRODKE tak terbaca"},
	{ErrIDPega, "IDPEGA bukan kasus EndorsementTreaty EDMT-"},
	{ErrTanggalAmbigu, "tanggal ambigu - tidak ditebak (K15)"},
	{ErrFormatTanggal, "format tanggal di luar YYYYMMDD / GMT"},
	{ErrNilaiKolom, "nilai tidak sesuai tipe kolom"},
	{ErrNoPolisKosong, "NOPOLIS kosong"},
	{ErrNoPolisBeda, "PolicyNo berbeda dari NOPOLIS"},
	{ErrProdKeBeda, "ProdKe dokumen berbeda dari PRODKE"},
	{ErrEDMNoKosong, "nomor endorsemen kosong"},
	{ErrEDMNoBeda, "EDMNo dokumen berbeda dari NOENDORS"},
	{ErrBentukTidakSah, "bentuk tidak sah untuk jenis proporsinya (AC 31, 33)"},
	{ErrDokumenGanda, "dokumen ganda untuk kasus yang sama"},
	{ErrGenerasiSebelumnyaTidakAda, "generasi sebelumnya belum dimuat (OLD_POLIS_ID tidak ditebak)"},
	{ErrPercabangan, "percabangan ditolak - UNIQUE OLD_POLIS_ID (ID-10)"},
	{ErrKeutuhan, "keutuhan ditolak - baris spreading generasi sebelumnya hilang (ID-15)"},
	{ErrOldDataTakSesuai, "OldData.EDMNo tidak sesuai generasi sebelumnya (ID-36)"},
	{ErrIDKasusBentrok, "ID kasus dipakai generasi lain"},
}

// JenisGalat - label ringkasan sebuah galat dokumen.
func JenisGalat(err error) string {
	for _, j := range pmJenisGalat {
		if errors.Is(err, j.err) {
			return j.nama
		}
	}
	return JenisGalatSimpan
}

// RingkasanPenanda - cacah penanda migrasi baris selisih SUMBER 'PEGA' (tiket 09).
type RingkasanPenanda struct {
	// Baris* - baris selisih yang ditandai (setiap baris PEGA terisi "0" / "1").
	BarisSpreading, BarisAngsuran, BarisLapisan int
	// Bergeser* - baris PASANGAN_BERGESER = 1: anomali sungguhan (ID-18, AC 13).
	BergeserSpreading, BergeserAngsuran, BergeserLapisan int
	// Berlapis - baris RUMUS_BERLAPIS = 1 (spreading + angsuran).
	Berlapis int
}

// RingkasanPemuat - cacah satu jalankan pemuat.
type RingkasanPemuat struct {
	// Tulis - `-jalankan` (true) atau uji-kering (false).
	Tulis bool
	// BarisDibaca - baris JSON_POLIS ber-PRODKE bukan 0 yang dibaca.
	BarisDibaca   int
	BukanTreatyIn int
	GenerasiNB    int
	// BarisAplikasiBaru - baris json_polis tulisan Utility1 aplikasi baru (ErrBarisAplikasiBaru), dilewati.
	BarisAplikasiBaru int
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
	// MedanBelumDiputuskan - baris arsip berkeputusan BELUM DIPUTUSKAN (F3: wajib 0).
	MedanBelumDiputuskan   int
	BelumDiputuskanPerPola map[string]int
	Usulan                 RingkasanUsulan
	Penanda                RingkasanPenanda
	// AngkaKasusTerbesar - nomor terbesar pyID `EDMT-` yang dimuat: SEQ_WORK_POLIS (`IDKasusBerikut`) wajib
	// dimajukan melewatinya.
	AngkaKasusTerbesar int64
	// GenerasiTerbesar - PRODKE terbesar yang dimuat.
	GenerasiTerbesar int
}

// RingkasanUsulan - cacah salinan SuggestList dokumen lama ke POOLDATA.HISTORYAKSEPTASIPRODUCTION.
type RingkasanUsulan struct {
	// Disalin - baris yang ditulis (atau siap ditulis). AKSES_LOGIN setiap baris NULL (`usulanlama.go`).
	Disalin int
	// TanpaPIC - baris dengan PIC kosong (NULL): baris dokumen tanpa `OperatorName`; tidak dikarang.
	TanpaPIC int
	// DokumenSudahAda - dokumen yang IDPEGA-nya sudah punya baris riwayat produksi: salinan dilewati (penjaga dobel).
	DokumenSudahAda int
}

// NasibUsulan - akibat salinan SuggestList satu dokumen.
type NasibUsulan int

const (
	// UsulanTanpaBaris - dokumen tanpa baris SuggestList yang disalin.
	UsulanTanpaBaris NasibUsulan = iota
	// UsulanDisalin - baris ditulis (atau, uji-kering, siap ditulis).
	UsulanDisalin
	// UsulanDilewati - penjaga dobel: IDPEGA sudah punya baris riwayat produksi.
	UsulanDilewati
)

// NasibUsulanUjiKering - uji-kering tidak menulis apa pun: baris yang ada siap disalin.
func NasibUsulanUjiKering(h HasilPecahEDM) NasibUsulan {
	if len(h.Usulan) == 0 {
		return UsulanTanpaBaris
	}
	return UsulanDisalin
}

// Selesai - nol dokumen gagal dan nol medan BELUM DIPUTUSKAN (F3). Medan yang dibuang menurut keputusan tertulis
// tidak menahan selesai.
func (r RingkasanPemuat) Selesai() bool { return r.DokumenGagal == 0 && r.MedanBelumDiputuskan == 0 }

func pmTulisPeta(b *strings.Builder, m map[string]int) {
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
	fmt.Fprintf(&b, "Pemuat dokumen lama EDM Treaty In (generasi endorsemen) - %s\n", mode)
	fmt.Fprintf(&b, "Baris JSON_POLIS PRODKE bukan 0 dibaca           : %d  (urut NOPOLIS, PRODKE naik)\n", r.BarisDibaca)
	fmt.Fprintf(&b, "Bukan dokumen Treaty In (lini lain), dilewati    : %d\n", r.BukanTreatyIn)
	if r.GenerasiNB > 0 {
		fmt.Fprintf(&b, "Generasi NB (PRODKE 0), dilewati - pemuat NB     : %d\n", r.GenerasiNB)
	}
	if r.BarisAplikasiBaru > 0 {
		fmt.Fprintf(&b, "Baris json_polis aplikasi baru, dilewati         : %d\n", r.BarisAplikasiBaru)
	}
	fmt.Fprintf(&b, "%-49s: %d\n", dimuat, r.Dimuat)
	fmt.Fprintf(&b, "Sudah dimuat sebelumnya, dilewati                : %d\n", r.SudahDimuat)
	fmt.Fprintf(&b, "Dokumen gagal (berkas laporan galat)             : %d\n", r.DokumenGagal)
	pmTulisPeta(&b, r.GalatPerJenis)
	fmt.Fprintf(&b, "Tanggal ambigu, tidak ditebak (K15)              : %d nilai\n", r.TanggalAmbigu)
	p := r.Penanda
	fmt.Fprintf(&b, "Penanda migrasi baris selisih SUMBER 'PEGA' (tiket 09, angka tidak diubah):\n")
	fmt.Fprintf(&b, "    PASANGAN_BERGESER = 1 (anomali, ID-18) spreading %d / %d baris, angsuran %d / %d, lapisan XOL %d / %d\n",
		p.BergeserSpreading, p.BarisSpreading, p.BergeserAngsuran, p.BarisAngsuran, p.BergeserLapisan, p.BarisLapisan)
	fmt.Fprintf(&b, "    RUMUS_BERLAPIS = 1 (generasi sebelumnya ber-EDMNo, ID-36)       : %d baris\n", p.Berlapis)
	fmt.Fprintf(&b, "Catatan SuggestList %-29s: %d baris -> POOLDATA.HISTORYAKSEPTASIPRODUCTION\n", disalin, r.Usulan.Disalin)
	b.WriteString("    AKSES_LOGIN selalu NULL (anggota sumbernya tidak ada di dokumen lama)\n")
	fmt.Fprintf(&b, "    PIC kosong (baris dokumen tanpa OperatorName, ditulis NULL)      : %d\n", r.Usulan.TanpaPIC)
	fmt.Fprintf(&b, "    dokumen dilewati - IDPEGA sudah punya baris riwayat produksi    : %d\n", r.Usulan.DokumenSudahAda)
	fmt.Fprintf(&b, "Medan dibuang menurut keputusan tertulis (arsip CSV, KEPUTUSAN dibuang):\n")
	pmTulisPeta(&b, r.DiabaikanPerAlasan)
	fmt.Fprintf(&b, "Medan BELUM DIPUTUSKAN (arsip CSV, KEPUTUSAN %s): %d  - WAJIB 0 sebelum pekerjaan dinyatakan selesai (F3)\n",
		KeputusanBelumDiputuskan, r.MedanBelumDiputuskan)
	pmTulisPeta(&b, r.BelumDiputuskanPerPola)
	if r.GenerasiTerbesar > 0 {
		fmt.Fprintf(&b, "PRODKE terbesar yang dimuat: %d\n", r.GenerasiTerbesar)
	}
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

// Lewat mencatat dokumen di luar lingkup (ErrBukanTreatyIn, ErrBukanGenerasiEndorsemen, ErrBarisAplikasiBaru).
func (l *LaporanPemuat) Lewat(err error) {
	if errors.Is(err, ErrBukanGenerasiEndorsemen) {
		l.r.GenerasiNB++
		return
	}
	if errors.Is(err, ErrBarisAplikasiBaru) {
		l.r.BarisAplikasiBaru++
		return
	}
	l.r.BukanTreatyIn++
}

// Gagal menulis setiap sebab satu dokumen ke berkas galat.
func (l *LaporanPemuat) Gagal(b BarisJSONPolis, g []GalatDokumen) error {
	l.r.DokumenGagal++
	for _, x := range g {
		l.r.GalatPerJenis[JenisGalat(x.Err)]++
		if errors.Is(x.Err, ErrTanggalAmbigu) {
			l.r.TanggalAmbigu++
		}
		if err := l.galat.Write([]string{b.IDPega, b.NoPolis, b.ProdKe, x.Jalur, x.Nilai, x.Err.Error()}); err != nil {
			return err
		}
	}
	return nil
}

// Berhasil mencatat dokumen yang dimuat (atau siap dimuat): setiap medannya yang tidak masuk kolom ke arsip CSV
// beserta keputusannya, nasib salinan SuggestList-nya, dan cacah penanda migrasinya.
func (l *LaporanPemuat) Berhasil(h HasilPecahEDM, usulan NasibUsulan, p PenandaMigrasi) error {
	l.r.Dimuat++
	for a, n := range h.Diabaikan {
		l.r.DiabaikanPerAlasan[a] += n
	}
	for _, m := range h.Arsip {
		keputusan := KeputusanBelumDiputuskan
		if m.Kunci != "" {
			keputusan = pmAwalanDibuang + m.Kunci
		} else {
			l.r.MedanBelumDiputuskan++
			l.r.BelumDiputuskanPerPola[m.Pola]++
		}
		if err := l.arsip.Write([]string{h.ID, m.Jalur, m.Nilai, keputusan}); err != nil {
			return err
		}
	}
	switch usulan {
	case UsulanDisalin:
		l.r.Usulan.Disalin += len(h.Usulan)
		for _, u := range h.Usulan {
			if u.PIC == "" {
				l.r.Usulan.TanpaPIC++
			}
		}
	case UsulanDilewati:
		l.r.Usulan.DokumenSudahAda++
	}
	s, a, x, b := p.Cacah()
	l.r.Penanda.BarisSpreading += len(p.Spreading)
	l.r.Penanda.BarisAngsuran += len(p.Angsuran)
	l.r.Penanda.BarisLapisan += len(p.Lapisan)
	l.r.Penanda.BergeserSpreading += s
	l.r.Penanda.BergeserAngsuran += a
	l.r.Penanda.BergeserLapisan += x
	l.r.Penanda.Berlapis += b
	if n, err := strconv.ParseInt(strings.TrimPrefix(h.ID, AwalanKasus), 10, 64); err == nil &&
		strings.HasPrefix(h.ID, AwalanKasus) && n > l.r.AngkaKasusTerbesar {
		l.r.AngkaKasusTerbesar = n
	}
	if h.ProdKe > l.r.GenerasiTerbesar {
		l.r.GenerasiTerbesar = h.ProdKe
	}
	return nil
}

// SudahDimuat mencatat dokumen yang generasinya sudah ada dari jalankan sebelumnya.
func (l *LaporanPemuat) SudahDimuat() { l.r.SudahDimuat++ }

// Tutup mengosongkan penyangga kedua berkas.
func (l *LaporanPemuat) Tutup() error {
	l.arsip.Flush()
	l.galat.Flush()
	return errors.Join(l.arsip.Error(), l.galat.Error())
}

// Ringkasan - salinan cacah saat ini.
func (l *LaporanPemuat) Ringkasan() RingkasanPemuat { return l.r }
