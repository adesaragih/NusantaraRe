package services

// Kesiapan Confirm Premium List Detail (keputusan work owner 06-10-2026):
// "pastikan tidak ada error, atau pesan yang menyatakan apapun sebelum bisa
// generate PL Number dan Resolved-Completed".
//
// Untuk apa berkas ini: SATU daftar pemeriksaan yang dipakai dua pintu -
//   - GET .../kesiapan-confirm: layar mengunci Confirm dan menyebut alasannya;
//   - Confirm di tahap Input Premium Detail: ditolak 409 (nol tulisan, nomor tidak
//     terbit, kasus tidak ditutup) selama daftar ini tidak kosong.
// Semua yang kurang dilaporkan sekaligus. MEMBACA saja.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/premiumlistlife/backend/models"
	"nusantarare/modul/premiumlistlife/backend/repository"
)

// ErrBelumSiapConfirm - masih ada yang kurang sebelum Confirm (409, berpesan).
var ErrBelumSiapConfirm = errors.New("services: polis belum siap di-Confirm")

// Pesan kesiapan Confirm - tampil di layar Decision.
const (
	PesanSiapDataPolis = "Policy data is incomplete (%s) - fill it in and press Save Data."
	PesanSiapWPC       = "WPC is empty - press Save Data."
	PesanSiapPeserta   = "No participants yet - upload the CSV and press Calculate CSV."
	PesanSiapSummary   = "Summary is not calculated - press Calculate CSV."
	PesanSiapSpreading = "%d participant(s) have no retro spreading - press Calculate CSV again."
	PesanSiapPeriode   = "Production period cannot be determined (%s)."
)

// JawabanKesiapanConfirm - jawaban GET .../kesiapan-confirm.
type JawabanKesiapanConfirm struct {
	Siap   bool     `json:"siap"`
	Alasan []string `json:"alasan"`
}

// KesiapanConfirm - pintu baca (dengan pagar identitas / kasus).
func (p *Penawaran) KesiapanConfirm(ctx context.Context, pelaku inti.Pelaku, polisID string) (JawabanKesiapanConfirm, error) {
	keadaan, err := p.pagari(ctx, pelaku, polisID)
	if err != nil {
		return JawabanKesiapanConfirm{}, err
	}
	alasan, err := p.alasanBelumSiap(ctx, pelaku, keadaan)
	if err != nil {
		return JawabanKesiapanConfirm{}, err
	}
	return JawabanKesiapanConfirm{Siap: len(alasan) == 0, Alasan: alasan}, nil
}

// periksaSiapConfirm - gerbang Confirm: galat ErrBelumSiapConfirm berisi seluruh alasan.
func (p *Penawaran) periksaSiapConfirm(ctx context.Context, pelaku inti.Pelaku, keadaan repository.KeadaanPolis) error {
	alasan, err := p.alasanBelumSiap(ctx, pelaku, keadaan)
	if err != nil {
		return err
	}
	if len(alasan) > 0 {
		return fmt.Errorf("%w: %s", ErrBelumSiapConfirm, strings.Join(alasan, " "))
	}
	return nil
}

// alasanBelumSiap - seluruh yang kurang; kosong = siap. Hanya tahap Input Premium Detail
// (tahap lain tidak menerbitkan nomor lewat Confirm).
func (p *Penawaran) alasanBelumSiap(ctx context.Context, pelaku inti.Pelaku, keadaan repository.KeadaanPolis) ([]string, error) {
	alasan := []string{}
	if keadaan.Status != models.TahapPolisDetail {
		return alasan, nil
	}
	// 1-2. Data polis lengkap (aturan Save Data yang sama) dan WPC sudah dihitung.
	d, err := repository.NewPenawaran(p.svc.DB()).BacaDataPolis(ctx, keadaan.ID)
	if errors.Is(err, repository.ErrHeaderPolisTidakAda) {
		return nil, fmt.Errorf("%w: %s", ErrPolisTakDitemukan, keadaan.ID)
	}
	if err != nil {
		return nil, err
	}
	if _, err := models.SusunDataPolis(d.IsianDataPolis); err != nil {
		sebab := strings.TrimPrefix(err.Error(), models.ErrDataPolisBelumLengkap.Error()+": ")
		sebab = strings.TrimPrefix(sebab, models.ErrDataPolisTidakSah.Error()+": ")
		alasan = append(alasan, fmt.Sprintf(PesanSiapDataPolis, sebab))
	}
	if d.WPC == nil || d.WPC.IsZero() {
		alasan = append(alasan, PesanSiapWPC)
	}
	// 3-5. Peserta, summary, spreading retro (Type QR).
	c, err := repository.NewKesiapanPolis(p.svc.DB()).Cacah(ctx, keadaan.ID)
	if err != nil {
		return nil, err
	}
	if c.Peserta == 0 {
		alasan = append(alasan, PesanSiapPeserta)
	}
	if c.Peserta > 0 && c.Summary == 0 {
		alasan = append(alasan, PesanSiapSummary)
	}
	if models.HitungPesertaPerTipe(d.Type) && c.TanpaSpreading > 0 {
		alasan = append(alasan, fmt.Sprintf(PesanSiapSpreading, c.TanpaSpreading))
	}
	// 6. Periode produksi dapat ditentukan (TANGGAL_CLOSING) - bahan PL Number dan WPC.
	if _, err := p.svc.Periode().Sekarang(ctx, pelaku); err != nil {
		if errors.Is(err, db.ErrTanpaOracle) {
			return nil, err
		}
		alasan = append(alasan, fmt.Sprintf(PesanSiapPeriode, err.Error()))
	}
	return alasan, nil
}
