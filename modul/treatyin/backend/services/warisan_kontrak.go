package services

// Baca satu kontrak WARISAN untuk form — keputusan pemilik proses.
//
// ⛔ Terjemahan tampilnya MEMAKAI ULANG `TanggalTampil` dan
// `SifatProporsiTampil` yang sudah hidup di `warisan_daftar.go`. Dua
// penerjemah untuk satu aturan adalah dua tempat untuk salah, dan yang kedua
// selalu yang basi.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/repository"
)

var (
	// ErrWarisanTidakAda - pengenalnya tidak menunjuk kontrak mana pun.
	ErrWarisanTidakAda = repository.ErrWarisanTidakAda
	// ErrJSONWarisanRusak - dokumennya ada tetapi tidak dapat diurai.
	ErrJSONWarisanRusak = repository.ErrJSONWarisanRusak
)

// BacaKontrakWarisan membaca satu kontrak warisan, siap tampil.
func (l *Layanan) BacaKontrakWarisan(ctx context.Context, p inti.Pelaku, id string) (models.KontrakWarisan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.KontrakWarisan{}, err
	}
	// ⛔ Pengenal KOSONG ditolak di sini, bukan diteruskan: kueri dengan
	// pengenal kosong mengembalikan nol baris, dan "tidak ada" adalah
	// jawaban yang berbeda dari "tidak ditanyakan".
	if strings.TrimSpace(id) == "" {
		return models.KontrakWarisan{}, fmt.Errorf("%w: pengenal kontrak kosong", ErrMasukanTidakSah)
	}
	k, err := l.gudang.BacaKontrakWarisan(ctx, id)
	if err != nil {
		return models.KontrakWarisan{}, err
	}
	// ⭐ TIGA TAB dari tabel pendaratan, bukan lagi dari CLOB. Satu
	// mekanisme untuk prop dan non-prop: nol cabang menurut
	// `PROPORTIONTYPE` di sini, dan tab yang kosong mengembalikan nol baris
	// alih-alih galat.
	if k.PeriodePelaporan, err = l.gudang.BacaPeriodePelaporan(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	if k.Portofolio, err = l.gudang.BacaPortofolio(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	if k.Akumulasi, err = l.gudang.BacaAkumulasi(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	// Empat tab berikutnya - dimuat migrasi 430/431, dan sampai 3 Oktober
	// 2026 nol layar membacanya. Tabel terisi yang tidak dibaca siapa pun
	// adalah pekerjaan yang terlihat selesai dan tidak sampai ke pemakai.
	if k.Egnpi, err = l.gudang.BacaEgnpi(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	if k.Retensi, err = l.gudang.BacaRetensi(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	if k.Angsuran, err = l.gudang.BacaAngsuran(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	if k.Catatan, err = l.gudang.BacaCatatan(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}

	// ⭐ Empat tab dari `M_TREATY_IN2`, satu pembacaan.
	//
	// ⚠️ Kosong di sini punya DUA arti, dan services tidak membedakannya —
	// yang membedakan layar, lewat petunjuk kosong. Terukur: 510 kontrak
	// punya `Limits[]` berisi di dokumennya tetapi nol baris di tabel ini.
	if k.Layer, err = l.gudang.BacaLayerWarisan(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}
	if k.SkalaKoasuransi, err = l.gudang.BacaSkalaKoasuransi(ctx, id); err != nil {
		return models.KontrakWarisan{}, err
	}

	// ⭐ Dua tab TEKS, dipilih menurut cabang — Jalan B, keputusan §15.
	// Dijalankan SESUDAH `SifatProporsiAsli` terisi, sebab cabangnya yang
	// menentukan ejaan mana yang dipakai.
	isiTabTeks(&k)

	k.SifatProporsi = SifatProporsiTampil(k.SifatProporsiAsli)
	k.TanggalMulai = TanggalTampil(k.TanggalMulaiAsli)
	k.TanggalBerakhir = TanggalTampil(k.TanggalBerakhirAsli)

	// ⭐ TANGGAL DI DALAM LARIK ikut diterjemahkan.
	//
	// ⛔ Terlewat sampai 3 Oktober 2026: ketiga medan kepala sudah memakai
	// `TanggalTampil`, sementara tanggal di dalam grid tidak - di layar
	// terbaca `20250118`, bentuk simpanan, bukan bentuk baca. Satu aturan
	// yang berlaku separuh lebih buruk daripada yang tidak berlaku: yang
	// melihatnya mengira dua medan itu memang bertipe berbeda.
	//
	// Nilai TERSIMPAN tidak berubah - tabel pendaratan tetap memegang
	// `YYYYMMDD`, dan yang diterjemahkan hanya jawaban API. `TanggalTampil`
	// mengembalikan apa adanya untuk yang bukan delapan angka, jadi nilai
	// kosong dan nilai aneh lewat tanpa dikarang.
	for i := range k.Kurs {
		k.Kurs[i].BerlakuDari = TanggalTampil(k.Kurs[i].BerlakuDari)
		k.Kurs[i].BerlakuSampai = TanggalTampil(k.Kurs[i].BerlakuSampai)
	}
	for i := range k.PeriodePelaporan {
		p := &k.PeriodePelaporan[i]
		p.TanggalAwal = TanggalTampil(p.TanggalAwal)
		p.JatuhTempoKirim = TanggalTampil(p.JatuhTempoKirim)
		p.JatuhTempoKonfir = TanggalTampil(p.JatuhTempoKonfir)
		p.JatuhTempoBayar = TanggalTampil(p.JatuhTempoBayar)
	}
	for i := range k.Angsuran {
		g := &k.Angsuran[i]
		g.JatuhTempo = TanggalTampil(g.JatuhTempo)
		g.TanggalBayar = TanggalTampil(g.TanggalBayar)
	}
	for i := range k.Egnpi {
		k.Egnpi[i].PerTanggal = TanggalTampil(k.Egnpi[i].PerTanggal)
	}
	for i := range k.Catatan {
		k.Catatan[i].Tanggal = TanggalTampil(k.Catatan[i].Tanggal)
	}
	for i := range k.Akumulasi {
		a := &k.Akumulasi[i]
		a.TanggalLapor = TanggalTampil(a.TanggalLapor)
		a.JatuhTempoKirim = TanggalTampil(a.JatuhTempoKirim)
	}
	if k.AdaDiJSON == nil {
		k.AdaDiJSON = map[string]bool{}
	}
	return k, nil
}

// WarisanTidakAda menjawab apakah galatnya "kontraknya tidak ada".
//
// Dipisah supaya handler tidak perlu mengimpor `repository` - arah
// ketergantungan `handlers -> services -> repository` dijaga penjaga inti.
func WarisanTidakAda(err error) bool { return errors.Is(err, ErrWarisanTidakAda) }

// WarisanJSONRusak menjawab apakah galatnya "dokumennya tidak dapat diurai".
func WarisanJSONRusak(err error) bool { return errors.Is(err, ErrJSONWarisanRusak) }
