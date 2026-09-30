package services

// Jenis reasuransi untuk Treaty Contract Out - tiket 02.
//
// Untuk apa berkas ini: SATU tempat daftar jenis reasuransi non-life dibaca
// (AC tiket 02: dipakai layar kontrak maupun seluruh grid klausul, bukan
// disalin per layar). Saringannya hidup di repository sebagai SQL; tabel
// kebenarannya `repository.LolosSaringanNonLifeTCO`.
//
// ⛔ Master yang KOSONG atau tidak terbaca adalah KEGAGALAN yang terlihat
// (ADR-0015), bukan daftar kosong yang diam: pemilih tanpa pilihan terbaca
// "belum ada jenis reasuransi", padahal yang terjadi adalah master yang
// tidak terjangkau atau saringannya menyingkirkan seluruhnya.
//
// ⛔ Pembaca DISUNTIK (pola `penyuntikan_test.go`): bawaannya gagal terang,
// handler memasang implementasi Oracle. Uji murni memasang pembaca palsu.
//
// Dibaca sesudah: repository/tco_jenisreasuransi.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatycontractout/backend/repository"
)

var (
	// ErrMasterJenisReasuransiKosong - saringan non-life tidak meloloskan satu
	// pun baris master, atau masternya memang kosong.
	ErrMasterJenisReasuransiKosong = errors.New(
		"services: reinsurance type master REINSURANCETYPE is empty or has no active non-life type")
	// ErrPembacaJenisReasuransiBelumDisuntik - handler lupa memasang pembaca.
	ErrPembacaJenisReasuransiBelumDisuntik = errors.New(
		"services: reinsurance type master reader is not injected")
	// ErrPilihanAnakTreatyLimitKosong - master tidak memuat satu pun jenis
	// porsi aktif maupun induknya: pemilih anak Treaty Limit tanpa pilihan.
	ErrPilihanAnakTreatyLimitKosong = errors.New(
		"services: reinsurance type master REINSURANCETYPE has no active Treaty Limit child type")
)

// PembacaJenisReasuransiTCO membaca master tersaring.
type PembacaJenisReasuransiTCO interface {
	DaftarNonLife(ctx context.Context) ([]repository.JenisReasuransiTCO, error)
	// DaftarAnakTreatyLimit - pilihan ReinsType anak Treaty Limit di bawah
	// induk `induk` (`models.PilihanReinsAnakTreatyLimit`).
	DaftarAnakTreatyLimit(ctx context.Context, induk string) ([]repository.JenisReasuransiTCO, error)
}

// JenisReasuransi adalah satu pilihan untuk layar.
//
// ⚠️ Nama kuncinya mengikuti nama properti Pega (`.ID`, `.Note`, `.Type`)
// supaya pembanding layar tidak perlu menerjemahkan; `tipe` dipilih atas
// `type` sebab yang terakhir kata kunci di banyak bahasa.
type JenisReasuransi struct {
	ID   string `json:"id"`
	Note string `json:"note"`
	Tipe string `json:"tipe"`
}

type pembacaJenisReasuransiBelumDisuntik struct{}

func (pembacaJenisReasuransiBelumDisuntik) DaftarNonLife(context.Context) (
	[]repository.JenisReasuransiTCO, error) {
	return nil, ErrPembacaJenisReasuransiBelumDisuntik
}

func (pembacaJenisReasuransiBelumDisuntik) DaftarAnakTreatyLimit(context.Context, string) (
	[]repository.JenisReasuransiTCO, error) {
	return nil, ErrPembacaJenisReasuransiBelumDisuntik
}

type pembacaJenisReasuransiOracle struct{ svc *Service }

func (p pembacaJenisReasuransiOracle) DaftarNonLife(ctx context.Context) (
	[]repository.JenisReasuransiTCO, error) {
	if !p.svc.PunyaDatabase() {
		return nil, db.ErrTanpaOracle
	}
	return repository.NewMasterJenisReasuransi(p.svc.DB()).DaftarNonLife(ctx)
}

func (p pembacaJenisReasuransiOracle) DaftarAnakTreatyLimit(ctx context.Context, induk string) (
	[]repository.JenisReasuransiTCO, error) {
	if !p.svc.PunyaDatabase() {
		return nil, db.ErrTanpaOracle
	}
	return repository.NewMasterJenisReasuransi(p.svc.DB()).DaftarAnakTreatyLimit(ctx, induk)
}

// PembacaJenisReasuransiOracle adalah pembaca sungguhan, dipasang handler.
func PembacaJenisReasuransiOracle(svc *Service) PembacaJenisReasuransiTCO {
	return pembacaJenisReasuransiOracle{svc: svc}
}

// JenisReasuransiTreaty melayani daftar jenis reasuransi non-life.
type JenisReasuransiTreaty struct {
	svc     *Service
	pembaca PembacaJenisReasuransiTCO
}

// JenisReasuransiTreaty menyusun layanannya dengan pembaca yang GAGAL TERANG.
func (s *Service) JenisReasuransiTreaty() *JenisReasuransiTreaty {
	return &JenisReasuransiTreaty{svc: s, pembaca: pembacaJenisReasuransiBelumDisuntik{}}
}

// DenganPembaca memasang pembaca master.
func (j *JenisReasuransiTreaty) DenganPembaca(p PembacaJenisReasuransiTCO) *JenisReasuransiTreaty {
	salin := *j
	salin.pembaca = p
	return &salin
}

// Daftar mengembalikan jenis reasuransi non-life yang aktif, urutan `.Note`.
func (j *JenisReasuransiTreaty) Daftar(ctx context.Context, pelaku inti.Pelaku) ([]JenisReasuransi, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	baris, err := j.pembaca.DaftarNonLife(ctx)
	if err != nil {
		return nil, err
	}
	if len(baris) == 0 {
		return nil, ErrMasterJenisReasuransiKosong
	}
	out := make([]JenisReasuransi, 0, len(baris))
	for _, b := range baris {
		out = append(out, JenisReasuransi{ID: b.ID, Note: b.Note, Tipe: b.Tipe})
	}
	return out, nil
}

// DaftarAnakTreatyLimit mengembalikan pilihan ReinsType baris anak Treaty Limit
// di bawah induk `induk`: dua belas jenis porsi + induknya, Flag active,
// urutan `.Note` [keputusan work owner 30-09-2026].
func (j *JenisReasuransiTreaty) DaftarAnakTreatyLimit(ctx context.Context, pelaku inti.Pelaku, induk string) (
	[]JenisReasuransi, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	induk = strings.TrimSpace(induk)
	if induk == "" {
		return nil, fmt.Errorf("%w: parent ReinsTypeID is required", galat.ErrPermintaanTidakSah)
	}
	baris, err := j.pembaca.DaftarAnakTreatyLimit(ctx, induk)
	if err != nil {
		return nil, err
	}
	if len(baris) == 0 {
		return nil, ErrPilihanAnakTreatyLimitKosong
	}
	out := make([]JenisReasuransi, 0, len(baris))
	for _, b := range baris {
		out = append(out, JenisReasuransi{ID: b.ID, Note: b.Note, Tipe: b.Tipe})
	}
	return out, nil
}
