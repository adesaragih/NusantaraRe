package services

// Daftar kasus renewal (tiket R05) - port RenewalList_RD:
// `D:\migrasi\RNM\RNW Fac In\ReportDefinition\RenewalList_RD.xml`
// (ASM-FW-GISFW-WORK-RENEWAL!RENEWALLIST_RD), dibaca dari `pyContent` (kueri saat
// jalan). `[terverifikasi]` wadah `pyUI` memuat empat filter yang sama, dihitung dua
// cara: urai ElementTree per wadah (4 + 4) dan `grep -c "<pyFilterName>"` = 8.
//
// ⚠️ `[keputusan agent]` A40, menunggu konfirmasi: daftar ini BUKAN daftar kandidat
// jatuh tempo. `[terverifikasi]` filternya hanya pembuat, team group, dan status
// kerja; tidak ada filter tanggal berakhir maupun OldPolicyNo - keduanya kolom
// tampil. Yang terekam = kasus renewal milik pembuatnya, di team group yang sama,
// kecuali berstatus Resolved-Completed dan Resolved-Rejected (status Resolved-*
// lain tetap lolos).

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"nusantarare/modul/rnwfacin/backend/models"
	"nusantarare/modul/rnwfacin/backend/repository"
)

// batasBaris - `pyContent/pyMaxRecords` = 500. `[terverifikasi]` satu-satunya di rule
// itu: `grep -c "<pyMaxRecords>"` = 1, dan urai menemukannya di `/pyContent` saja.
const batasBaris = 500

var (
	// ErrSaringanKosong - parameter filter kosong. ⚠️ Penyimpangan sadar, bukan port:
	// `pyUseNullIfEmpty` keempat filter tidak diisi, dan `[dugaan]` Pega lalu MELEWATI
	// filternya - tanpa filter A, daftar memuat kasus pengguna lain. `[pertanyaan
	// terbuka]` perilaku sebenarnya; sampai itu ditolak.
	ErrSaringanKosong = errors.New("services: parameter daftar renewal kosong")
	// ErrStatusKerjaKosong - baris sumber tanpa status kerja. `[dugaan]` filter B/D
	// menjadi SQL `<>`, dan Oracle memperlakukan '' sebagai NULL sehingga baris itu
	// terbuang; pembanding Go akan meloloskannya. Belum terverifikasi → ditolak.
	ErrStatusKerjaKosong = errors.New("services: baris daftar renewal tanpa status kerja")
)

// Saringan - parameter RenewalList_RD.
type Saringan struct {
	// PembuatID - `Param.UserNameID` (filter A, `.pxCreateOperator =`): login
	// pengguna yang membuka daftar. Penyedianya pemanggil (sesi); layanan ini tidak
	// menyentuh autentikasi.
	PembuatID string
	// TeamGroup - `Param.TeamGroup` (filter C, `.Quotation.TeamGroup =`).
	TeamGroup string
}

// HasilDaftar - baris terurut; Terpotong = lebih dari batasBaris yang lolos saringan
// (tambahan sistem baru: Pega memotong tanpa tanda yang terekam di rule).
type HasilDaftar struct {
	Baris     []models.BarisDaftarRenewal
	Terpotong bool
}

// ServiceDaftar - layanan daftar kasus renewal.
type ServiceDaftar struct {
	sumber repository.SumberDaftarRenewal
}

// NewServiceDaftar merakit layanan daftar.
func NewServiceDaftar(sumber repository.SumberDaftarRenewal) *ServiceDaftar {
	return &ServiceDaftar{sumber: sumber}
}

// Daftar - `pyFilterLogic` `B AND A AND C AND D`, urut `pxCreateDateTime` DESC
// (pySortOrder 1) lalu `pyID` DESC (2), 500 baris teratas. Semua pembanding persis
// (`pyIgnoreCase` tidak diisi). `[dugaan]` pyID dibandingkan sebagai teks, seperti
// ORDER BY kolom karakter.
func (s *ServiceDaftar) Daftar(ctx context.Context, sr Saringan) (HasilDaftar, error) {
	// Pesan galat tanpa nilai parameter: PembuatID adalah login (CLAUDE.md §4.10).
	if sr.PembuatID == "" {
		return HasilDaftar{}, fmt.Errorf("%w: pembuat", ErrSaringanKosong)
	}
	if sr.TeamGroup == "" {
		return HasilDaftar{}, fmt.Errorf("%w: team group", ErrSaringanKosong)
	}
	semua, err := s.sumber.KasusRenewal(ctx)
	if err != nil {
		return HasilDaftar{}, err
	}
	var lolos []models.BarisDaftarRenewal
	for _, b := range semua {
		if b.StatusKerja == "" {
			return HasilDaftar{}, fmt.Errorf("%w: kasus %s", ErrStatusKerjaKosong, b.IDKasus)
		}
		// Urutan pyFilterLogic dipertahankan: B (bukan selesai) AND A (pembuat) AND
		// C (team group) AND D (bukan ditolak).
		if b.StatusKerja != models.StatusKerjaSelesai && b.PembuatID == sr.PembuatID &&
			b.TeamGroup == sr.TeamGroup && b.StatusKerja != models.StatusKerjaDitolak {
			lolos = append(lolos, b)
		}
	}
	sort.SliceStable(lolos, func(i, j int) bool {
		if !lolos[i].DibuatPada.Equal(lolos[j].DibuatPada) {
			return lolos[i].DibuatPada.After(lolos[j].DibuatPada)
		}
		return lolos[i].IDKasus > lolos[j].IDKasus
	})
	if len(lolos) > batasBaris {
		return HasilDaftar{Baris: lolos[:batasBaris], Terpotong: true}, nil
	}
	return HasilDaftar{Baris: lolos}, nil
}
