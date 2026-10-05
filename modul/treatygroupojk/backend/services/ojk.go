package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/modul/treatygroupojk/backend/models"
	"nusantarare/modul/treatygroupojk/backend/repository"
)

// Daftar - baris bersaring kata (ID, Name, Name IDN), urut Order No.
func (l *Layanan) Daftar(ctx context.Context, kata string) ([]models.Ojk, error) {
	return l.gudang.Daftar(ctx, kata)
}

// Buka - satu baris.
func (l *Layanan) Buka(ctx context.Context, id string) (models.Ojk, error) {
	o, err := l.gudang.Ambil(ctx, nil, strings.TrimSpace(id))
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Ojk{}, ErrTidakAda
	}
	return o, err
}

// Rapikan - spasi tepi dibuang; Name dan Name (IDN) huruf besar (keputusan work owner 05-10-2026).
func Rapikan(i models.Isian) models.Isian {
	return models.Isian{
		ID:      strings.TrimSpace(i.ID),
		Name:    strings.ToUpper(strings.TrimSpace(i.Name)),
		NameIDN: strings.ToUpper(strings.TrimSpace(i.NameIDN)),
	}
}

// periksa menilai isian yang sudah dirapikan; seluruh kesalahan dikumpulkan.
func periksa(i models.Isian) error {
	var pesan []string
	for _, c := range []struct{ nama, nilai string }{{"Name", i.Name}, {"Name (IDN)", i.NameIDN}} {
		if c.nilai == "" {
			pesan = append(pesan, c.nama+" is required")
		}
	}
	for _, c := range []struct {
		nama, nilai string
		batas       int
	}{{"Name", i.Name, models.BatasNama}, {"Name (IDN)", i.NameIDN, models.BatasNama}} {
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
func (l *Layanan) Simpan(ctx context.Context, a Aktor, isi models.Isian) (models.Ojk, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Ojk{}, err
	}
	isi = Rapikan(isi)
	if err := periksa(isi); err != nil {
		return models.Ojk{}, err
	}
	id := isi.ID
	err := l.tx(ctx, func(tx *dbTx) error {
		var order string
		if id != "" {
			lama, err := l.gudang.Ambil(ctx, tx, id)
			if err != nil {
				return err
			}
			order = lama.OrderNo
		} else {
			// Kunci dulu, baru periksa kembar: Add bersamaan menunggu giliran dan melihat baris yang baru tersimpan.
			n, o, err := l.gudang.KunciNomorTertinggi(ctx, tx)
			if err != nil {
				return err
			}
			id, order = fmt.Sprintf("%02d", n+1), fmt.Sprint(o+1)
			if len(id) > models.BatasID || len(order) > models.BatasOrder {
				return fmt.Errorf("services: ID %s atau Order No %s baru melewati batas kolom", id, order)
			}
		}
		if ganda, err := l.gudang.PemakaiNama(ctx, tx, isi.Name, isi.ID); err != nil {
			return err
		} else if len(ganda) > 0 {
			return tolak("Name is already used by %s", ganda[0].ID)
		}
		baris := models.Ojk{ID: id, Name: isi.Name, NameIDN: isi.NameIDN, OrderNo: order}
		if isi.ID != "" {
			return l.gudang.Ubah(ctx, tx, baris)
		}
		return l.gudang.Sisip(ctx, tx, baris)
	})
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Ojk{}, ErrTidakAda
	}
	if err != nil {
		return models.Ojk{}, err
	}
	return l.Buka(ctx, id)
}
