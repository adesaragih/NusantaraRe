package services

// Efek kasus Komite dan laporan harian "perlu intervensi" - tiket 08.
//
// Untuk apa berkas ini: kegagalan efek keluar yang tidak pulih TERLIHAT -
// menempel pada kasusnya (`KasusKomiteTampil.Efek`), dan dalam satu laporan
// harian yang tetap ada walau kosong (AC tiket 08).
//
// ⚠️ `[keputusan work owner]` siapa yang menindaklanjuti belum ditetapkan.
// Laporannya dibatasi `ReasLifeAdmin` (`[asumsi — OQ-007/OQ-021]`, sama dengan
// eskalasi tiket 03); pengirimannya setiap hari menunggu penjadwal (belum ada).

import (
	"context"
	"strings"
	"time"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/modul/komite/models"
	"nusantarare/modul/komite/repository"
)

// EfekTampil adalah satu efek untuk layar - KATA, bukan kode (km5).
type EfekTampil struct {
	Jenis     string `json:"jenis"`
	Keadaan   string `json:"keadaan"`
	Percobaan int    `json:"percobaan"`
	// Sejak - waktu terakhir baris berubah (atau lahir, bila belum dipungut).
	Sejak string `json:"sejak"`
}

// RingkasEfekKasus adalah seluruh efek satu kasus beserta keadaannya.
type RingkasEfekKasus struct {
	Keadaan string       `json:"keadaan"`
	Efek    []EfekTampil `json:"efek"`
}

func keEfekTampil(e repository.EfekKasusKomite) EfekTampil {
	sejak := e.Dibuat
	if e.Diperbarui.Valid {
		sejak = e.Diperbarui.Time
	}
	return EfekTampil{Jenis: e.Jenis, Keadaan: models.KataStatusEfek(e.Status),
		Percobaan: e.Percobaan, Sejak: sejak.Format("2006-01-02 15:04:05")}
}

// ringkasEfekKasus menyusun ringkasan - murni.
func ringkasEfekKasus(daftar []repository.EfekKasusKomite) RingkasEfekKasus {
	r := RingkasEfekKasus{Efek: []EfekTampil{}}
	var kode []string
	for _, e := range daftar {
		r.Efek = append(r.Efek, keEfekTampil(e))
		kode = append(kode, e.Status)
	}
	r.Keadaan = models.KeadaanEfekKasus(kode)
	return r
}

// BarisLaporanIntervensi adalah satu baris laporan harian.
type BarisLaporanIntervensi struct {
	KasusID string `json:"kasusId"`
	EfekTampil
}

// LaporanHarianIntervensi adalah laporan harian.
//
// ⛔ `Kosong` DINYATAKAN, tidak disimpulkan dari panjang daftar: laporan yang
// tidak datang harus dapat dibedakan dari laporan yang datang tanpa masalah.
type LaporanHarianIntervensi struct {
	Tanggal string                   `json:"tanggal"`
	Kosong  bool                     `json:"kosong"`
	Baris   []BarisLaporanIntervensi `json:"baris"`
}

// susunLaporan - murni.
func susunLaporan(saat time.Time, daftar []repository.EfekKasusKomite) LaporanHarianIntervensi {
	l := LaporanHarianIntervensi{Tanggal: saat.Format("2006-01-02"), Baris: []BarisLaporanIntervensi{}}
	for _, e := range daftar {
		kasus := e.KasusID
		if i := strings.Index(kasus, "#"); i >= 0 {
			kasus = kasus[:i] // rujukan email `kasus#Tn`
		}
		l.Baris = append(l.Baris, BarisLaporanIntervensi{KasusID: kasus, EfekTampil: keEfekTampil(e)})
	}
	l.Kosong = len(l.Baris) == 0
	return l
}

// LaporanHarian menyusun laporan "perlu intervensi" hari ini.
func (i *InboxKomite) LaporanHarian(ctx context.Context, pelaku inti.Pelaku, saat time.Time) (
	LaporanHarianIntervensi, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return LaporanHarianIntervensi{}, err
	}
	if !pelaku.PunyaPeran(inti.PeranAdmin) {
		return LaporanHarianIntervensi{}, inti.ErrTanpaWewenang
	}
	if i == nil || i.svc == nil || !i.svc.PunyaDatabase() {
		return LaporanHarianIntervensi{}, db.ErrTanpaOracle
	}
	daftar, err := repository.NewInboxKomite(i.svc.DB()).EfekPerluIntervensi(ctx)
	if err != nil {
		return LaporanHarianIntervensi{}, err
	}
	return susunLaporan(saat, daftar), nil
}
