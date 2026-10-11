package services

// Tombol `Save` layar Medical Check - 11-10-2026 (work owner: "RECOMMENDATION per peserta di bawah grid diagnosa,
// dapat diisi Medical Advisor di tahap Medical Check, disimpan tombol Save Medical Check").
//
// Korpus: `kelvin/InboxClaimLife (Done)/FlowAction/MedicalCheckClaimLife.xml` (Save · Send Back to Admin · Send to
// Claim Analyst) dan `Section/ClaimLifeDetailGCNM.xml` b7316/b7323 `.RECOMMENDATION` (pyDisabledWhen
// `.STS_REJECT=='1' || .STS_REJECT=='2'`). Kolom `T_CLAIMLF_PREMIUMLIST_DETAIL.RECOMMENDATION VARCHAR2(1000)` sudah
// ada sejak migrasi 003 - tanpa migrasi baru.
//
// ⛔ Diagnosa TIDAK lewat sini: grid diagnosa menulis per aksi (Add / Choose / Delete) lewat diagnosa.go, gerbang
// yang sama (tahap Medical Check + Medical Advisor + STS_REJECT).
//
// Dibaca sesudah: isianoutstanding.go, diagnosa.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

// BatasRekomendasi - T_CLAIMLF_PREMIUMLIST_DETAIL.RECOMMENDATION VARCHAR2(1000), migrasi 003. Byte, diperiksa SEBELUM
// bind supaya tidak pernah jatuh ke ORA-12899.
const BatasRekomendasi = 1000

var (
	// ErrMedisBukanMedicalCheck - tombol Save Medical Check hanya ada di tahap Medical Check.
	ErrMedisBukanMedicalCheck = errors.New("services: Save Medical Check hanya tersedia pada tahap Medical Check")
	// ErrRekomendasiTerkunci - Recommendation dibekukan bila peserta sudah diputus (b7323 pyDisabledWhen STS_REJECT).
	ErrRekomendasiTerkunci = errors.New("services: Recommendation tidak dapat diubah untuk peserta yang sudah diputus")
	// ErrRekomendasiKepanjangan - melebihi lebar kolom.
	ErrRekomendasiKepanjangan = errors.New("services: Recommendation melebihi 1000 byte")
)

// IsianMedis - isian tombol Save Medical Check: Recommendation per pengenal peserta.
type IsianMedis struct {
	Rekomendasi map[string]string `json:"rekomendasi"`
}

// PeriksaIsianMedis - bentuk isian, MURNI. `aktif` = pengenal peserta aktif klaim.
func PeriksaIsianMedis(is IsianMedis, aktif map[string]bool) error {
	if len(is.Rekomendasi) == 0 {
		return fmt.Errorf("%w: tidak ada isian untuk disimpan", galat.ErrPermintaanTidakSah)
	}
	for id, v := range is.Rekomendasi {
		if !aktif[id] {
			return fmt.Errorf("%w: peserta %q bukan peserta aktif klaim", galat.ErrPermintaanTidakSah, id)
		}
		if len(v) > BatasRekomendasi {
			return fmt.Errorf("%w: peserta %q (%d byte)", ErrRekomendasiKepanjangan, id, len(v))
		}
	}
	return nil
}

// SimpanMedis adalah layanan tombol Save layar Medical Check.
type SimpanMedis struct{ svc *Service }

// SimpanMedis menyusun layanannya.
func (s *Service) SimpanMedis() *SimpanMedis { return &SimpanMedis{svc: s} }

// Simpan menulis Recommendation per peserta dalam SATU transaksi.
//
// Gerbang: identitas, pengenal, kasus terbuka, tahap Medical Check + pemegangnya (ReasLifeMedicalAdvisor), bentuk
// isian, lalu tulis (syarat STS_REJECT di DALAM pernyataan). ⛔ Status baris, nomor, dan tahap TIDAK disentuh.
func (x *SimpanMedis) Simpan(ctx context.Context, pelaku inti.Pelaku, klaimID string, is IsianMedis) error {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if strings.TrimSpace(klaimID) == "" {
		return fmt.Errorf("%w: pengenal klaim wajib diisi", galat.ErrPermintaanTidakSah)
	}
	if x == nil || x.svc == nil || !x.svc.PunyaDatabase() {
		return db.ErrTanpaOracle
	}
	if err := x.svc.PastikanKasusTerbuka(ctx, klaimID); err != nil {
		if errors.Is(err, repository.ErrWorkTidakAda) {
			return fmt.Errorf("%w: %w", ErrKlaimTidakAda, err)
		}
		return err
	}
	baca := repository.NewKlaimLife(x.svc.DB())
	tahap, err := tahapKasus(ctx, baca, klaimID)
	if err != nil {
		return err
	}
	if tahap != models.TahapMedicalCheck {
		return fmt.Errorf("%w: tahap %s", ErrMedisBukanMedicalCheck, tahap)
	}
	if err := WajibPemegangTahap(pelaku, tahap); err != nil {
		return err
	}
	peserta, err := baca.AmbilPeserta(ctx, klaimID)
	if err != nil {
		return err
	}
	aktif := map[string]bool{}
	for _, p := range peserta {
		aktif[p.ID] = true
	}
	if err := PeriksaIsianMedis(is, aktif); err != nil {
		return err
	}
	return x.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		for id, v := range is.Rekomendasi {
			ok, err := baca.SetelRekomendasi(ctx, tx, klaimID, id, v)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("%w: peserta %q", ErrRekomendasiTerkunci, id)
			}
		}
		return nil
	})
}
