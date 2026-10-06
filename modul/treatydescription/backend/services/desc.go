package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/modul/treatydescription/backend/models"
	"nusantarare/modul/treatydescription/backend/repository"
)

// cobaIDMaksimum - ID baru yang dicoba sebelum menyerah bila nomor sequence sudah terpakai baris lama.
const cobaIDMaksimum = 50

// saringan - nilai saringan "0"/"1" diteruskan; selain itu tanpa saringan.
func saringan(s string) string {
	if s = strings.TrimSpace(s); s == "0" || s == "1" {
		return s
	}
	return ""
}

// Daftar - baris bersaring kata (ID atau nama), jenis ("0" Non XOL / "1" XOL), dan status ("1" Active / "0"
// Inactive); "" atau nilai lain = semua.
func (l *Layanan) Daftar(ctx context.Context, kata, xol, status string) ([]models.Desc, error) {
	return l.gudang.Daftar(ctx, kata, saringan(xol), saringan(status))
}

// Buka - satu baris dan jumlah pemakaiannya.
func (l *Layanan) Buka(ctx context.Context, id string) (models.Desc, error) {
	d, err := l.gudang.Ambil(ctx, nil, strings.TrimSpace(id))
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Desc{}, ErrTidakAda
	}
	return d, err
}

// Rapikan - spasi tepi dibuang; nama huruf besar (seluruh nama DEV huruf besar).
func Rapikan(i models.Isian) models.Isian {
	return models.Isian{
		ID:          strings.TrimSpace(i.ID),
		DescName:    strings.ToUpper(strings.TrimSpace(i.DescName)),
		IsXOL:       strings.TrimSpace(i.IsXOL),
		StatusAktif: strings.TrimSpace(i.StatusAktif),
	}
}

func periksa(i models.Isian) error {
	var pesan []string
	if i.DescName == "" {
		pesan = append(pesan, "Description Name is required")
	}
	if len(i.DescName) > models.BatasNama {
		pesan = append(pesan, fmt.Sprintf("Description Name is longer than %d characters", models.BatasNama))
	}
	if i.IsXOL != models.NonXOL && i.IsXOL != models.XOL {
		pesan = append(pesan, "Type must be Non XOL or XOL")
	}
	if i.StatusAktif != models.Aktif && i.StatusAktif != models.Nonaktif {
		pesan = append(pesan, "Status must be Active or Inactive")
	}
	if len(pesan) > 0 {
		return tolak("%s", strings.Join(pesan, "; "))
	}
	return nil
}

// Simpan - Add (ID kosong) atau Edit. Menjawab baris yang tersimpan.
func (l *Layanan) Simpan(ctx context.Context, a Aktor, isi models.Isian) (models.Desc, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Desc{}, err
	}
	isi = Rapikan(isi)
	if err := periksa(isi); err != nil {
		return models.Desc{}, err
	}
	id := isi.ID
	err := l.tx(ctx, func(tx *dbTx) error {
		if id != "" {
			if _, err := l.gudang.Ambil(ctx, tx, id); err != nil {
				return err
			}
		}
		ganda, err := l.gudang.PemakaiNama(ctx, tx, isi.DescName, id)
		if err != nil {
			return err
		}
		if len(ganda) > 0 {
			return tolak("Description Name is already used by %s", ganda[0])
		}
		r := models.Desc{ID: id, DescName: isi.DescName, IsXOL: isi.IsXOL, StatusAktif: isi.StatusAktif}
		if id != "" {
			return l.gudang.Ubah(ctx, tx, r)
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
				id, r.ID = baru, baru
				break
			}
		}
		return l.gudang.Sisip(ctx, tx, r)
	})
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Desc{}, ErrTidakAda
	}
	if err != nil {
		return models.Desc{}, err
	}
	return l.Buka(ctx, id)
}
