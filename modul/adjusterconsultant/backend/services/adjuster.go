package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"nusantarare/modul/adjusterconsultant/backend/models"
	"nusantarare/modul/adjusterconsultant/backend/repository"
)

// cobaIDMaksimum - ID baru yang dicoba sebelum menyerah bila nomor sequence sudah terpakai baris lama.
const cobaIDMaksimum = 50

// Daftar - baris bersaring: kata (ID, nama, alamat, telepon) dan status (`models.Status*`).
func (l *Layanan) Daftar(ctx context.Context, kata, status string) ([]models.Adjuster, error) {
	bendera := ""
	switch strings.ToLower(strings.TrimSpace(status)) {
	case models.StatusSemua:
	case models.StatusAktif:
		bendera = models.BenderaAktif
	case models.StatusNonaktif:
		bendera = models.BenderaNonaktif
	default:
		return nil, tolak("Status must be active, inactive, or empty")
	}
	return l.gudang.Daftar(ctx, kata, bendera)
}

// Buka - satu baris.
func (l *Layanan) Buka(ctx context.Context, id string) (models.Adjuster, error) {
	a, err := l.gudang.Ambil(ctx, nil, strings.TrimSpace(id))
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Adjuster{}, ErrTidakAda
	}
	return a, err
}

// potongByte memotong teks ke n byte tanpa memotong satu karakter.
func potongByte(s string, n int) string {
	if len(s) <= n {
		return s
	}
	s = s[:n]
	for !utf8.ValidString(s) {
		s = s[:len(s)-1]
	}
	return s
}

// Rapikan - isian seperti `SaveAdjusterConsultant_Act`: spasi tepi dibuang, NAME dan ADDRESS huruf besar, TELPNO apa
// adanya.
func Rapikan(i models.Isian) models.Isian {
	return models.Isian{
		ID:      strings.TrimSpace(i.ID),
		Name:    strings.ToUpper(strings.TrimSpace(i.Name)),
		Address: strings.ToUpper(strings.TrimSpace(i.Address)),
		TelpNo:  strings.TrimSpace(i.TelpNo),
	}
}

// periksa menilai isian yang sudah dirapikan; seluruh kesalahan dikumpulkan.
func periksa(i models.Isian) error {
	var pesan []string
	if i.Name == "" {
		pesan = append(pesan, "Name is required")
	}
	for _, c := range []struct {
		nama, nilai string
		batas       int
	}{{"Name", i.Name, models.BatasNama}, {"Address", i.Address, models.BatasAlamat}, {"Telp No", i.TelpNo, models.BatasTelp}} {
		if len(c.nilai) > c.batas {
			pesan = append(pesan, fmt.Sprintf("%s is longer than %d characters", c.nama, c.batas))
		}
	}
	if len(pesan) > 0 {
		return tolak("%s", strings.Join(pesan, "; "))
	}
	return nil
}

// Simpan - Add (ID kosong) atau Edit. Menjawab baris yang tersimpan.
func (l *Layanan) Simpan(ctx context.Context, a Aktor, isi models.Isian) (models.Adjuster, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Adjuster{}, err
	}
	isi = Rapikan(isi)
	if err := periksa(isi); err != nil {
		return models.Adjuster{}, err
	}
	username := potongByte(a.AkunID, models.BatasUsername)
	id := isi.ID
	err := l.tx(ctx, func(tx *dbTx) error {
		if id != "" {
			if _, err := l.gudang.Ambil(ctx, tx, id); err != nil {
				return err
			}
		}
		ganda, err := l.gudang.PemakaiNama(ctx, tx, isi.Name, id)
		if err != nil {
			return err
		}
		if len(ganda) > 0 {
			g := ganda[0]
			if !g.Active {
				return tolak("Name is already used by %s, which is inactive - activate it instead of adding a new one", g.ID)
			}
			return tolak("Name is already used by %s", g.ID)
		}
		baris := models.Adjuster{ID: id, Name: isi.Name, Address: isi.Address, TelpNo: isi.TelpNo}
		if id != "" {
			return l.gudang.Ubah(ctx, tx, baris, username)
		}
		for i := 0; ; i++ {
			if i >= cobaIDMaksimum {
				return fmt.Errorf("services: tidak menemukan ID kosong sesudah %d nomor sequence", cobaIDMaksimum)
			}
			baru, err := l.gudang.IDBaru(ctx, tx)
			if err != nil {
				return err
			}
			ada, err := l.gudang.AdaID(ctx, tx, baru)
			if err != nil {
				return err
			}
			if !ada {
				id, baris.ID = baru, baru
				break
			}
		}
		return l.gudang.Sisip(ctx, tx, baris, username)
	})
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Adjuster{}, ErrTidakAda
	}
	if err != nil {
		return models.Adjuster{}, err
	}
	return l.Buka(ctx, id)
}

// SetelAktif - Activate / Deactivate (pengganti hapus, keputusan work owner 05-10-2026).
func (l *Layanan) SetelAktif(ctx context.Context, a Aktor, id string, aktif bool) (models.Adjuster, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Adjuster{}, err
	}
	id = strings.TrimSpace(id)
	err := l.tx(ctx, func(tx *dbTx) error {
		return l.gudang.SetelAktif(ctx, tx, id, aktif, potongByte(a.AkunID, models.BatasUsername))
	})
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Adjuster{}, ErrTidakAda
	}
	if err != nil {
		return models.Adjuster{}, err
	}
	return l.Buka(ctx, id)
}
