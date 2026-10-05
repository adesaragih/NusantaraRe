package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"nusantarare/modul/treatyexchangeyearly/backend/models"
	"nusantarare/modul/treatyexchangeyearly/backend/repository"
)

// cobaIDMaksimum - ID baru yang dicoba sebelum menyerah bila nomor sequence sudah terpakai baris lama.
const cobaIDMaksimum = 50

var (
	polaTahun = regexp.MustCompile(`^[0-9]{4}$`)
	// polaAngka - kurs: angka dengan titik desimal, tanpa pemisah ribuan (bentuk data warisan `16560.00`, `14500`).
	polaAngka = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
	// polaHariPega - delapan digit tanggal di depan teks Pega (juga baris warisan `20190801T00000.000 GMT`).
	polaHariPega = regexp.MustCompile(`^([0-9]{8})`)
	wib          = time.FixedZone("WIB", 7*60*60)
)

// FormatWaktuPega - bentuk DATEIU / DATEIN yang ditulis Pega (`20230713T104129.285 GMT`).
func FormatWaktuPega(t time.Time) string { return t.UTC().Format("20060102T150405.000") + " GMT" }

// TanggalPega - tanggal layar `YYYY-MM-DD` ke format Pega yang benar `YYYYMMDDT000000.000 GMT`.
func TanggalPega(s string) (string, bool) {
	t, err := time.Parse("2006-01-02", strings.TrimSpace(s))
	if err != nil {
		return "", false
	}
	return t.Format("20060102") + "T000000.000 GMT", true
}

// HariPega - delapan digit tanggal teks Pega (`YYYYMMDD`); "" bila bentuknya lain.
func HariPega(s string) string {
	m := polaHariPega.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return ""
	}
	if _, err := time.Parse("20060102", m[1]); err != nil {
		return ""
	}
	return m[1]
}

// TampilTanggal - teks Pega ke `DD-MM-YYYY`; bentuk lain apa adanya.
func TampilTanggal(s string) string {
	h := HariPega(s)
	if h == "" {
		return s
	}
	return h[6:8] + "-" + h[4:6] + "-" + h[0:4]
}

// TampilWaktu - teks Pega ke `DD-MM-YYYY HH:MM` WIB; bentuk lain apa adanya.
func TampilWaktu(s string) string {
	t, err := time.Parse("20060102T150405.000 MST", strings.TrimSpace(s))
	if err != nil {
		return s
	}
	return t.In(wib).Format("02-01-2006 15:04")
}

func lengkapi(k models.Kurs) models.Kurs {
	k.Mulai, k.Akhir, k.Diubah = TampilTanggal(k.StartDate), TampilTanggal(k.EndDate), TampilWaktu(k.DateIU)
	return k
}

// Daftar - baris bersaring kata dan tahun treaty ("" = semua).
func (l *Layanan) Daftar(ctx context.Context, kata, tahun string) ([]models.Kurs, error) {
	d, err := l.gudang.Daftar(ctx, kata, tahun)
	for i := range d {
		d[i] = lengkapi(d[i])
	}
	return d, err
}

// Pilihan - isi dropdown Currency dan saringan tahun.
func (l *Layanan) Pilihan(ctx context.Context) ([]models.MataUang, []string, error) {
	m, err := l.gudang.DaftarMataUang(ctx)
	if err != nil {
		return nil, nil, err
	}
	t, err := l.gudang.DaftarTahun(ctx)
	return m, t, err
}

// potongByte memotong teks ke n byte tanpa memotong satu karakter.
func potongByte(s string, n int) string {
	if len(s) <= n {
		return s
	}
	s = s[:n]
	for !utf8.ValidString(s) {
		s = s[:len(s)-1]
	}
	return s
}

// Rapikan - spasi tepi dibuang; Quarter kosong = 0 (tahunan).
func Rapikan(i models.Isian) models.Isian {
	r := models.Isian{
		Kunci:      strings.TrimSpace(i.Kunci),
		TreatyYear: strings.TrimSpace(i.TreatyYear),
		IDCurrency: strings.TrimSpace(i.IDCurrency),
		StartDate:  strings.TrimSpace(i.StartDate),
		EndDate:    strings.TrimSpace(i.EndDate),
		ToIDR:      strings.TrimSpace(i.ToIDR),
		ToUSD:      strings.TrimSpace(i.ToUSD),
		Quarter:    strings.TrimSpace(i.Quarter),
	}
	if r.Quarter == "" {
		r.Quarter = "0"
	}
	return r
}

func periksa(i models.Isian) error {
	var pesan []string
	if !polaTahun.MatchString(i.TreatyYear) {
		pesan = append(pesan, "Treaty Year must be 4 digits")
	}
	if i.IDCurrency == "" {
		pesan = append(pesan, "Currency is required")
	}
	mulai, okMulai := TanggalPega(i.StartDate)
	akhir, okAkhir := TanggalPega(i.EndDate)
	if !okMulai {
		pesan = append(pesan, "Start Date is required (a valid date)")
	}
	if !okAkhir {
		pesan = append(pesan, "End Date is required (a valid date)")
	}
	if okMulai && okAkhir && mulai > akhir {
		pesan = append(pesan, "End Date must not be before Start Date")
	}
	if !polaAngka.MatchString(i.ToIDR) {
		pesan = append(pesan, "To IDR must be a number (use a dot for decimals, no thousand separators)")
	}
	if i.ToUSD != "" && !polaAngka.MatchString(i.ToUSD) {
		pesan = append(pesan, "To USD must be a number (use a dot for decimals, no thousand separators)")
	}
	if !slices.Contains(models.Quarter, i.Quarter) {
		pesan = append(pesan, "Quarter must be 0 (yearly), 1, 2, 3, or 4")
	}
	if len(i.ToIDR) > models.BatasTeks || len(i.ToUSD) > models.BatasTeks {
		pesan = append(pesan, fmt.Sprintf("rates are longer than %d characters", models.BatasTeks))
	}
	if len(pesan) > 0 {
		return tolak("%s", strings.Join(pesan, "; "))
	}
	return nil
}

// tanggalSimpan - tanggal yang ditulis: Edit dengan hari yang sama = teks lama apa adanya; selainnya format Pega benar.
func tanggalSimpan(isian, lama string, edit bool) string {
	baru, _ := TanggalPega(isian)
	if edit && HariPega(lama) != "" && HariPega(lama) == baru[:8] {
		return lama
	}
	return baru
}

// Simpan - Add (Kunci kosong) atau Edit. Menjawab baris yang tersimpan.
func (l *Layanan) Simpan(ctx context.Context, a Aktor, isi models.Isian) (models.Kurs, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Kurs{}, err
	}
	isi = Rapikan(isi)
	if err := periksa(isi); err != nil {
		return models.Kurs{}, err
	}
	edit := isi.Kunci != ""
	var hasilID string
	err := l.tx(ctx, func(tx *dbTx) error {
		var k models.Kurs
		if edit {
			lama, err := l.gudang.AmbilKunci(ctx, tx, isi.Kunci)
			if err != nil {
				return err
			}
			k = lama
		}
		// Currency tidak diganti (Edit): kode tersimpan dibiarkan. Selainnya disalin dari view CURRENCY.
		if !edit || isi.IDCurrency != k.IDCurrency {
			m, err := l.gudang.AmbilMataUang(ctx, tx, isi.IDCurrency)
			if errors.Is(err, repository.ErrTidakAda) {
				return tolak("Currency %s is not in the currency list", isi.IDCurrency)
			}
			if err != nil {
				return err
			}
			k.IDCurrency, k.Currency, k.CurrencyName = m.ID, m.Kode, m.Nama
		}
		ganda, err := l.gudang.Kembar(ctx, tx, isi.TreatyYear, k.IDCurrency, isi.Quarter, isi.Kunci)
		if err != nil {
			return err
		}
		if len(ganda) > 0 {
			return tolak("Treaty Year %s, %s, Quarter %s already exists (ID %s)", isi.TreatyYear, k.Currency, isi.Quarter, ganda[0])
		}
		sekarang := FormatWaktuPega(l.jam())
		k.TreatyYear, k.Quarter, k.ToIDR, k.ToUSD = isi.TreatyYear, isi.Quarter, isi.ToIDR, isi.ToUSD
		k.StartDate, k.EndDate = tanggalSimpan(isi.StartDate, k.StartDate, edit), tanggalSimpan(isi.EndDate, k.EndDate, edit)
		k.UserID, k.DateIU = potongByte(a.AkunID, models.BatasTeks), sekarang
		if edit {
			hasilID = k.ID
			return l.gudang.Ubah(ctx, tx, k)
		}
		k.DateIn = sekarang
		for i := 0; ; i++ {
			if i >= cobaIDMaksimum {
				return fmt.Errorf("services: tidak menemukan ID kosong sesudah %d nomor sequence", cobaIDMaksimum)
			}
			baru, err := l.gudang.IDBaru(ctx, tx)
			if err != nil {
				return err
			}
			ada, err := l.gudang.AdaID(ctx, tx, baru)
			if err != nil {
				return err
			}
			if !ada {
				k.ID, hasilID = baru, baru
				break
			}
		}
		return l.gudang.Sisip(ctx, tx, k)
	})
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Kurs{}, ErrTidakAda
	}
	if err != nil {
		return models.Kurs{}, err
	}
	var k models.Kurs
	if edit {
		k, err = l.gudang.AmbilKunci(ctx, nil, isi.Kunci)
	} else {
		k, err = l.gudang.AmbilID(ctx, nil, hasilID)
	}
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Kurs{}, ErrTidakAda
	}
	return lengkapi(k), err
}
