package services

// Tiket 05 (maksud endorsement dan `EDMStatus`) dan 06 (jurnal balik) -
// `Save` b37202 (`VIS .EditInput=1`, `pyDisabledWhen` b37200 `.IsJsonPolis=1`)
// → `SetPremi_EDM`.
//
//	EdmType 1 (Perubahan Data): peserta `Old` yang dicentang (`.EdmBatal`,
//	  2.2 b1441) → `Delete` + 32 kolom dibalik (2.1 b1273).
//	EdmType 3 (Batal): SELURUH peserta `Old` → `Batal` + dibalik (2.3 b1583),
//	  tanpa penandaan satu per satu (spec user story 16).
//	Lalu rekap mata uang dihitung ulang (6 b4574 `AppendCurrencySummary_DT`)
//	dan kasus berstatus "sudah simpan" (8 b5600 `.IsJsonPolis = 1`).
//
// ⛔ Simpan kedua ditolak 409 - tombolnya mati di Pega sesudah simpan pertama
// (b37200), dan server tidak boleh hanya bergantung pada layar.

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/jejak"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

// ErrKasusTertutup - kasus sudah diputuskan; tidak dapat disunting lagi (AC 29).
var ErrKasusTertutup = errors.New("services: the endorsement case is already decided and cannot change")

// ErrSudahDisimpan - `Save` sudah ditekan (`IsJsonPolis = 1`, b37200).
var ErrSudahDisimpan = errors.New("services: the endorsement is already saved; Save is disabled once the summary exists")

// ErrTanpaPeserta - kasus tanpa peserta tidak dapat disimpan (rekap kosong).
var ErrTanpaPeserta = errors.New("services: the endorsement case has no participants to save")

// HasilSimpan - cacah peserta yang berubah status dan rekap mata uang baru.
type HasilSimpan struct {
	Ditandai int                 `json:"ditandai"`
	Status   string              `json:"status"`
	Rekap    []map[string]string `json:"rekap"`
}

// Simpan - `SetPremi_EDM`, satu transaksi.
func (l *Layanan) Simpan(ctx context.Context, p inti.Pelaku, id string, pilihan models.PilihanHapus) (HasilSimpan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilSimpan{}, err
	}
	if !models.KasusEDM(id) {
		return HasilSimpan{}, fmt.Errorf("%w: %q", ErrKasusTidakAda, id)
	}
	var hasil HasilSimpan
	err := l.tx(ctx, func(tx *db.Tx) error {
		k, err := l.gudang.AmbilKasus(ctx, tx, id, true)
		if errors.Is(err, repository.ErrTidakAda) {
			return fmt.Errorf("%w: %q", ErrKasusTidakAda, id)
		}
		if err != nil {
			return err
		}
		if !k.Terbuka() {
			return ErrKasusTertutup
		}
		sudah, err := l.gudang.AdaRekap(ctx, tx, id)
		if err != nil {
			return err
		}
		if sudah {
			return ErrSudahDisimpan
		}
		cacah, err := l.gudang.CacahPeserta(ctx, tx, id)
		if err != nil {
			return err
		}
		jumlah := 0
		for _, c := range cacah {
			jumlah += c
		}
		if jumlah == 0 {
			return ErrTanpaPeserta
		}
		switch k.EdmType {
		case models.EdmTypeBatal:
			// 2.3: seluruh peserta - pilihan layar tidak berlaku (grid Batal b17500 tanpa kotak centang).
			hasil.Status = models.StatusBatal
			hasil.Ditandai, err = l.gudang.Tandai(ctx, tx, id, models.StatusBatal, models.PilihanHapus{Semua: true})
		case models.EdmTypePerubahanData:
			hasil.Status = models.StatusDelete
			if pilihan.Semua || len(pilihan.Pilih) > 0 {
				hasil.Ditandai, err = l.gudang.Tandai(ctx, tx, id, models.StatusDelete, pilihan)
			}
		default:
			return fmt.Errorf("%w: EDM Type %q", ErrMasukanTidakSah, k.EdmType)
		}
		if errors.Is(err, repository.ErrDaftarKecualiTerlaluPanjang) {
			return fmt.Errorf("%w: too many unchecked rows after DELETE ALL", ErrMasukanTidakSah)
		}
		if err != nil {
			return err
		}
		if _, err := l.gudang.HitungRekap(ctx, tx, id, k.Kepala["TYPE"]); err != nil {
			return err
		}
		if hasil.Rekap, err = l.gudang.RekapKasus(ctx, tx, id); err != nil {
			return err
		}
		return l.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
			KlaimID: id, Dari: models.TahapInputEDMLife, Ke: models.TahapInputEDMLife, AkunID: p.AkunID, Waktu: l.jam(),
			Komentar: models.AksiSimpan,
		})
	})
	if err != nil {
		return HasilSimpan{}, err
	}
	if hasil.Rekap == nil {
		hasil.Rekap = []map[string]string{}
	}
	return hasil, nil
}
