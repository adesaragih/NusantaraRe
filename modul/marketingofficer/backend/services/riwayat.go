package services

// Log perubahan satu MO dari `MARKETINGOFFICER_LOG` (permintaan work owner 03-10-2026: "anggotanya ada log
// perubahan, ambil dari sini MARKETINGOFFICER_LOG").
//
// Setiap baris log = keadaan LAMA sebelum satu UPDATE (trigger warisan). Versi berurutan = baris log (tertua dulu)
// lalu baris sekarang; satu perubahan = selisih dua versi berurutan:
//   - waktu = LOG_TIME versi lama (saat UPDATE itu terjadi); baris lama tanpa LOG_TIME memakai TANGGAL versi
//     baru dan ditandai `perkiraan` (urutannya juga perkiraan);
//   - oleh  = USERUPDATE versi baru (Pega dan aplikasi Go sama-sama mengisinya saat menyimpan);
//   - Login Account hanya dibandingkan bila AKSES_LOGIN kedua versi diketahui (baris log bertanda `UPDATE-GO`, atau
//     baris sekarang).

import (
	"context"
	"errors"

	"nusantarare/modul/marketingofficer/backend/models"
	"nusantarare/modul/marketingofficer/backend/repository"
)

// RuasBerubah adalah satu kolom yang berubah. `Kolom` = nama kolom MARKETINGOFFICER; label layarnya di frontend.
type RuasBerubah struct {
	Kolom   string `json:"kolom"`
	Sebelum string `json:"sebelum"`
	Sesudah string `json:"sesudah"`
}

// Perubahan adalah satu UPDATE: kapan, oleh siapa, dan kolom apa saja yang berubah (bisa kosong - disimpan ulang
// tanpa perubahan).
type Perubahan struct {
	Waktu     string        `json:"waktu"`
	Oleh      string        `json:"oleh"`
	Perkiraan bool          `json:"perkiraan"`
	Ruas      []RuasBerubah `json:"ruas"`
}

// Riwayat adalah log perubahan satu MO, TERBARU lebih dulu.
type Riwayat struct {
	ID        string      `json:"id"`
	Perubahan []Perubahan `json:"perubahan"`
	JumlahLog int         `json:"jumlahLog"`
}

// versi adalah satu keadaan MO di riwayat.
type versi struct {
	m              models.MarketingOfficer
	logTime        string
	aksesDiketahui bool
}

// ruasDibandingkan - urutan kolom di layar log.
var ruasDibandingkan = []struct {
	kolom string
	nilai func(models.MarketingOfficer) string
	akses bool
}{
	{"CLIENTNAME", func(m models.MarketingOfficer) string { return m.ClientName }, false},
	{"CLIENTID", func(m models.MarketingOfficer) string { return m.ClientID }, false},
	{"AKSES_LOGIN", func(m models.MarketingOfficer) string { return m.AksesLogin }, true},
	{"CLIENTID2", func(m models.MarketingOfficer) string { return m.ClientID2 }, false},
	{"MOLEADER", func(m models.MarketingOfficer) string { return m.MOLeader }, false},
	{"BRANCHPARENT", func(m models.MarketingOfficer) string { return m.BranchParent }, false},
	{"BRANCHDETAILID", func(m models.MarketingOfficer) string { return m.BranchDetailID }, false},
	{"BRANCHDETAILNAME", func(m models.MarketingOfficer) string { return m.BranchDetailName }, false},
	{"TEAMGROUP", func(m models.MarketingOfficer) string { return m.TeamGroup }, false},
	{"MOSTATUS", func(m models.MarketingOfficer) string { return m.MOStatus }, false},
}

// SusunRiwayat menyusun perubahan dari log (tertua dulu) dan baris sekarang - fungsi murni.
func SusunRiwayat(kini models.MarketingOfficer, log []models.BarisLog) Riwayat {
	vs := make([]versi, 0, len(log)+1)
	for _, b := range log {
		vs = append(vs, versi{m: b.MarketingOfficer, logTime: b.LogTime, aksesDiketahui: b.AksesDiketahui()})
	}
	vs = append(vs, versi{m: kini, aksesDiketahui: true})
	r := Riwayat{ID: kini.ID, Perubahan: []Perubahan{}, JumlahLog: len(log)}
	for i := len(vs) - 2; i >= 0; i-- {
		lama, baru := vs[i], vs[i+1]
		p := Perubahan{Waktu: lama.logTime, Oleh: baru.m.UserUpdate, Ruas: []RuasBerubah{}}
		if p.Waktu == "" {
			p.Waktu, p.Perkiraan = baru.m.Tanggal, true
		}
		for _, ru := range ruasDibandingkan {
			if ru.akses && !(lama.aksesDiketahui && baru.aksesDiketahui) {
				continue
			}
			if a, b := ru.nilai(lama.m), ru.nilai(baru.m); a != b {
				p.Ruas = append(p.Ruas, RuasBerubah{Kolom: ru.kolom, Sebelum: a, Sesudah: b})
			}
		}
		r.Perubahan = append(r.Perubahan, p)
	}
	return r
}

// Riwayat membaca log perubahan satu MO.
func (l *Layanan) Riwayat(ctx context.Context, id string) (Riwayat, error) {
	var r Riwayat
	err := l.tx(ctx, func(tx *dbTx) error {
		kini, err := l.gudang.AmbilMO(ctx, tx, id)
		if errors.Is(err, repository.ErrTidakAda) {
			return ErrTidakAda
		}
		if err != nil {
			return err
		}
		log, err := l.gudang.LogMO(ctx, id)
		if err != nil {
			return err
		}
		r = SusunRiwayat(kini, log)
		return nil
	})
	return r, err
}
