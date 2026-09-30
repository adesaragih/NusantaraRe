package services

// Kotak masuk klaim - F0.4.
//
// Untuk apa berkas ini: gerbang antrian kerja. Siapa boleh melihat antrian
// mana, dan antrian mana yang hanya berisi kasus MILIKNYA.
//
// ⛔ DUA BENTUK ANTRIAN, dan XML memisahkannya dengan terang:
//
//	`[terverifikasi]` `Claim Life/Flow/Register_Flow.xml`
//	  Input Register    `WorkList`   1511 + `Current operator` 1508
//	  Outstanding Claim `WorkList`   1300 + `ToWorklist`        1351
//	  Medical Check     `WorkBasket`  993 + `<Workbasket>ReasLifeMedicalAdvisor` 1072
//	  Claim Analis      `WorkBasket` 1119 + `<Workbasket>ReasLifeSPV`            1198
//
// WORKLIST = kotak masuk PRIBADI: hanya kasus yang pelakunya sendiri buat.
// WORKBASKET = kotak masuk BERSAMA: setiap pemegang peran itu melihat semua.
//
// Menyamakan keduanya bukan penyederhanaan melainkan kebocoran ke satu arah
// dan kebutaan ke arah lain: Admin akan melihat kasus rekannya, dan Medical
// Advisor tidak akan melihat kasus yang menunggu perannya.
//
// Dibaca sesudah: tahap.go.

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

var (
	// ErrTahapTidakSah - nomor tahap di luar keempat yang ada.
	ErrTahapTidakSah = errors.New("services: tahap tidak sah")
)

// AntrianPribadi menyatakan tahap ini WORKLIST, bukan workbasket.
//
// ⛔ Daftar, bukan tebakan dari perannya. Kedua tahap Admin kebetulan
// keduanya worklist, tetapi itu FAKTA XML - bukan aturan yang dapat
// disimpulkan dari "Admin selalu worklist".
func AntrianPribadi(t models.Tahap) bool {
	switch t {
	case models.TahapInputRegister, models.TahapOutstanding:
		return true
	default:
		return false
	}
}

// HalamanInbox adalah satu halaman antrian beserta totalnya.
type HalamanInbox struct {
	Baris []repository.BarisInbox
	// Total adalah cacah SELURUH kasus pada tahap itu, bukan panjang halaman.
	// Lencana tab memakainya.
	Total  int
	Tahap  models.Tahap
	Offset int
	Ukuran int
}

// Inbox melayani pembacaan antrian kerja.
type Inbox struct{ svc *Service }

// KotakMasuk menyusun layanannya.
func (s *Service) KotakMasuk() *Inbox { return &Inbox{svc: s} }

// BatasUkuran memangkas ukuran halaman ke rentang yang XML tetapkan.
//
// ⚠️ Nol atau negatif jatuh ke bawaan, bukan menjadi galat: pemanggil yang
// tidak menyebut ukuran meminta ukuran bawaan.
func BatasUkuran(n int) int {
	switch {
	case n <= 0:
		return repository.UkuranHalamanBawaan
	case n > repository.UkuranHalamanMaksimum:
		return repository.UkuranHalamanMaksimum
	default:
		return n
	}
}

// Ambil membaca satu halaman antrian sebuah tahap.
//
// ⛔ Pelaku KOSONG ditolak, bukan dijawab daftar kosong. Daftar kosong
// terbaca sebagai "tidak ada pekerjaan" - dan itu kalimat yang berbeda
// artinya dari "saya tidak tahu siapa Anda".
func (in *Inbox) Ambil(ctx context.Context, pelaku inti.Pelaku, tahap models.Tahap,
	offset, ukuran int) (HalamanInbox, error) {

	var kosong HalamanInbox
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return kosong, err
	}
	if !tahap.Diketahui() {
		return kosong, fmt.Errorf("%w: %d", ErrTahapTidakSah, int(tahap))
	}
	peran, ada := models.PeranPemegangTahap(tahap)
	if !ada {
		return kosong, fmt.Errorf("%w: tahap %q", ErrPeranTahapBelumDiputuskan, tahap)
	}
	// ⛔ Gerbangnya SAMA untuk kedua bentuk: pemegang tahap itu. Yang berbeda
	// hanya apa yang terlihat sesudah gerbangnya terbuka.
	if err := inti.WajibPeran(pelaku, peran); err != nil {
		return kosong, err
	}
	if !in.svc.PunyaDatabase() {
		return kosong, db.ErrTanpaOracle
	}
	if offset < 0 {
		offset = 0
	}
	ukuran = BatasUkuran(ukuran)

	saring := repository.SaringInbox{
		Tahap:      tahap.String(),
		PeranTahap: peran,
		Offset:     offset,
		Ukuran:     ukuran,
	}
	// ⛔ Worklist disaring ke kasus MILIK pelaku; workbasket tidak.
	// `[terverifikasi]` `Register_Flow.xml` 1508 `Current operator`, 1351
	// `ToWorklist`.
	if AntrianPribadi(tahap) {
		saring.AkunID = pelaku.AkunID
	}

	baris, total, err := repository.NewKlaimLife(in.svc.DB()).AmbilInbox(ctx, saring)
	if err != nil {
		return kosong, err
	}
	return HalamanInbox{
		Baris: baris, Total: total, Tahap: tahap,
		Offset: offset, Ukuran: ukuran,
	}, nil
}

// TahapTerlihat menyebut tahap yang pelaku ini berhak lihat.
//
// ⚠️ Pelaku BERPERAN GANDA melihat gabungan - bukan salah satu. `pelakuDari`
// memecah `X-Peran` pada koma justru supaya itu mungkin, dan layar yang
// memaksa memilih satu akan menyembunyikan separuh pekerjaan orang itu.
func TahapTerlihat(pelaku inti.Pelaku) []models.Tahap {
	urut := []models.Tahap{
		models.TahapInputRegister,
		models.TahapOutstanding,
		models.TahapMedicalCheck,
		models.TahapClaimAnalis,
	}
	var terlihat []models.Tahap
	for _, t := range urut {
		peran, ada := models.PeranPemegangTahap(t)
		if !ada {
			continue
		}
		if inti.WajibPeran(pelaku, peran) == nil {
			terlihat = append(terlihat, t)
		}
	}
	return terlihat
}
