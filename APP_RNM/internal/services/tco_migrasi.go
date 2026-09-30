package services

// Migrasi data Treaty Contract Out - tiket 01 (tco2).
//
// Satu pintu dari composition root (`cmd/api -migrate-data-treaty-contract-out`)
// ke pemindah di repository. Nol aturan dagang: batas transaksinya di
// repository sebab ia membungkus pembacaan, penulisan, dan rekonsiliasi
// dalam satu transaksi yang sama.
//
// ⛔ TIDAK dijalankan di DEV oleh sesi modul; milik work owner sesudah
// penyatuan ke main dan `-migrate`.

import (
	"context"

	"nusantarare/internal/repository"
)

// PindahkanDataTreatyContractOut memindahkan enam tabel warisan ke T_*.
func (s *Service) PindahkanDataTreatyContractOut(ctx context.Context) (repository.LaporanMigrasiTCO, error) {
	if !s.PunyaDatabase() {
		return repository.LaporanMigrasiTCO{}, repository.ErrTanpaOracle
	}
	return repository.NewMigrasiTCO(s.db).Pindahkan(ctx)
}
