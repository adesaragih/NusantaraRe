package services

// Kasus polis baru - tombol portal `Input Offer` / `Input Premium`
// (GILIRAN-13 paket 1, pl3 + bn).
//
// Untuk apa berkas ini: meniru `Activity/CreateInputLife.xml` - melahirkan
// satu work object `ASM-FW-GISFW-Work-LIFE` (`svcAddWorkObject` b444),
// menyimpan benderanya (b618), dan menyimpannya (b726) - di SATU transaksi
// bersama jejaknya. Keadaan awalnya dari models.SusunKasusPolisBaru.
//
// ⛔ Tanpa gerbang peran, sama dengan seluruh layanan PremiumList: korpus
// tidak menggerbangi tombol portal ini, dan mengarang gerbang berarti
// memutuskan siapa boleh membuat penawaran.
//
// Dibaca sesudah: models/polis_kasus.go, repository/polis_kasus.go.

import (
	"context"
	"fmt"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/jejak"
	"nusantarare/modul/premiumlistlife/backend/models"
	"nusantarare/modul/premiumlistlife/backend/repository"
)

// KasusPolis adalah layanan pembuat kasus polis.
type KasusPolis struct {
	svc   *Service
	jejak jejak.Jejak
}

// KasusPolis menyusun layanannya dengan jejak bawaan yang gagal terang.
func (s *Service) KasusPolis() *KasusPolis {
	return &KasusPolis{svc: s, jejak: jejak.JejakBelumDiputuskan{}}
}

// DenganJejak mengganti perekamnya.
func (k *KasusPolis) DenganJejak(j jejak.Jejak) *KasusPolis {
	salin := *k
	salin.jejak = j
	return &salin
}

// HasilKasusPolisBaru adalah kasus yang baru lahir - cukup untuk membukanya.
type HasilKasusPolisBaru struct {
	CaseID     string `json:"caseId"`
	StatusWork string `json:"statusWork"`
	Position   string `json:"position"`
	Flag       string `json:"flag"`
}

// namaTombolFlag menyebut tombol pembuatnya di jejak - KATA, bukan kode.
var namaTombolFlag = map[string]string{
	models.FlagPolisPenawaran: "Input Offer",
	models.FlagPolisPremium:   "Input Premium",
}

// Buat melahirkan satu kasus polis.
func (k *KasusPolis) Buat(ctx context.Context, pelaku inti.Pelaku, flag string,
	saat time.Time) (HasilKasusPolisBaru, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HasilKasusPolisBaru{}, err
	}
	awal, err := models.SusunKasusPolisBaru(flag)
	if err != nil {
		return HasilKasusPolisBaru{}, fmt.Errorf("%w: %w", galat.ErrPermintaanTidakSah, err)
	}
	if k == nil || k.svc == nil || !k.svc.PunyaDatabase() {
		return HasilKasusPolisBaru{}, db.ErrTanpaOracle
	}
	kerja := repository.NewWorkPolis(k.svc.DB())
	var id string
	err = k.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		var err error
		if id, err = kerja.PengenalBerikut(ctx, tx); err != nil {
			return err
		}
		if err := kerja.SisipKasusBaru(ctx, tx, id, awal, pelaku.AkunID, saat); err != nil {
			return err
		}
		// ⛔ Kelahiran kasus ikut terekam (ADR-0007): tanpa baris ini jejak
		// polis dimulai dari perpindahan pertamanya, dan siapa yang membuatnya
		// tidak tercatat di mana pun.
		return k.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
			KlaimID: id,
			Dari:    "",
			Ke:      awal.Status + " (" + namaTombolFlag[flag] + ")",
			AkunID:  pelaku.AkunID,
			Waktu:   saat,
		})
	})
	if err != nil {
		return HasilKasusPolisBaru{}, err
	}
	return HasilKasusPolisBaru{CaseID: id, StatusWork: awal.Status,
		Position: awal.Posisi, Flag: awal.Flag}, nil
}
