package services

// Cabut peserta - OQ-M6 DITUTUP 29-09-2026 (GILIRAN-17) `[keputusan work owner]`:
// mencabut peserta = PENANDA, layar menyembunyikannya.
//
// `[terverifikasi]` `Section/InputOSClaimLife.xml`: tombol `DELETE` b17865 di grid
// `PremiumListDetail` menjalankan `deleteRow` b17874 (tanpa konfirmasi b18021),
// lalu `DeletePesertaClaimLife` b17909; tampil hanya bila `CLAIM_NO == ''`
// (b18082). Di Pega peserta DILEPAS dari halaman kasus; di sini ditandai
// (`STS_HAPUS`, migrasi 022, ADR-U-0031), dan setiap pembaca menyaringnya.
//
// ⛔ Baris cermin warisan `OS_AKSEPTASI_KLAIM_LIFE` peserta itu (ditulis saat
// pendaftaran, status NULL; sejak OQ-N2 bernama/ber-DOB) DIBUANG di transaksi
// yang sama (temuan /code-review GILIRAN-17): di Pega baris cermin baru lahir
// saat Save Outstanding, jadi peserta yang dilepas sebelumnya tidak pernah
// terlihat hilir. Hanya baris berstatus NULL milik `CASEID` klaim ini.
//
// Dibaca sesudah: dol.go (gerbang tahap), statusbaris.go (jejak).

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/jejak"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

// ErrPesertaTidakDapatDicabut - tombol `DELETE` hanya ada di layar Outstanding
// dan hanya sebelum Save to RNM pertama (b18082).
var ErrPesertaTidakDapatDicabut = errors.New(
	"services: peserta hanya dapat dicabut di tahap Outstanding Claim sebelum Save to RNM")

// CabutPeserta menandai seorang peserta dicabut dari klaimnya.
//
// Urutan gerbangnya sama dengan dialog Edit Date: identitas, peran Admin,
// pengenal, kasus terbuka, tahap + penanda Save to RNM, lalu kepemilikan.
func (st *Status) CabutPeserta(ctx context.Context, pelaku inti.Pelaku,
	klaimID, pesertaID string, saat time.Time) error {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if err := inti.WajibPeran(pelaku, inti.PeranAdmin); err != nil {
		return err
	}
	if strings.TrimSpace(klaimID) == "" || strings.TrimSpace(pesertaID) == "" {
		return fmt.Errorf("%w: pengenal klaim dan peserta wajib diisi", galat.ErrPermintaanTidakSah)
	}
	if !st.svc.PunyaDatabase() {
		return db.ErrTanpaOracle
	}
	if err := st.svc.PastikanKasusTerbuka(ctx, klaimID); err != nil {
		return err
	}

	baca := repository.NewKlaimLife(st.svc.DB())
	tahap, err := tahapKasus(ctx, baca, klaimID)
	if err != nil {
		return err
	}
	sudah, err := baca.SudahSaveRNM(ctx, klaimID)
	if err != nil {
		return err
	}
	if !models.BolehCabutPeserta(tahap, sudah) {
		return fmt.Errorf("%w: tahap %s, sudah Save to RNM %v", ErrPesertaTidakDapatDicabut, tahap, sudah)
	}
	peserta, err := baca.AmbilPeserta(ctx, klaimID)
	if err != nil {
		return err
	}
	milik := false
	for _, p := range peserta {
		if p.ID == pesertaID {
			milik = true
			break
		}
	}
	if !milik {
		return fmt.Errorf("%w: peserta %q bukan milik klaim %q",
			galat.ErrPermintaanTidakSah, pesertaID, klaimID)
	}
	perBaris, err := baca.AmbilBaris(ctx, klaimID)
	if err != nil {
		return err
	}
	caseID, err := baca.CaseIDKlaim(ctx, klaimID)
	if err != nil {
		return err
	}
	// ⛔ Penanda, cermin, dan jejaknya dalam SATU transaksi (ADR-U-0007).
	return st.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		if err := baca.CabutPeserta(ctx, tx, klaimID, pesertaID); err != nil {
			return err
		}
		for _, b := range perBaris[pesertaID] {
			if err := baca.HapusCerminBelumDisimpan(ctx, tx, b.ID, caseID); err != nil {
				return err
			}
		}
		return st.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
			KlaimID: klaimID,
			Dari:    "peserta " + pesertaID,
			Ke:      "dicabut",
			AkunID:  pelaku.AkunID,
			Waktu:   saat,
		})
	})
}
