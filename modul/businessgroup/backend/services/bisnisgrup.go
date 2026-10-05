package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/modul/businessgroup/backend/models"
	"nusantarare/modul/businessgroup/backend/repository"
)

// cobaIDMaksimum - ID baru yang dicoba sebelum menyerah bila nomor sequence sudah terpakai baris lama.
const cobaIDMaksimum = 50

// Syariah - nama berakhiran SYARIAH (tanpa beda huruf dan spasi tepi): tidak dikelola modul ini.
func Syariah(nama string) bool {
	return strings.HasSuffix(strings.ToUpper(strings.TrimSpace(nama)), models.AkhiranSyariah)
}

// Daftar - baris tanpa SYARIAH, bersaring kata dan Treaty Group ("" = semua).
func (l *Layanan) Daftar(ctx context.Context, kata, topID string) ([]models.BisnisGrup, error) {
	return l.gudang.Daftar(ctx, kata, topID)
}

// Pilihan - isi dropdown Treaty Group.
func (l *Layanan) Pilihan(ctx context.Context) ([]models.TreatyGroup, error) {
	return l.gudang.DaftarTreaty(ctx)
}

// ambil - satu baris; yang berakhiran SYARIAH dianggap tidak ada.
func (l *Layanan) ambil(ctx context.Context, tx *dbTx, id string) (models.BisnisGrup, error) {
	b, err := l.gudang.Ambil(ctx, tx, id)
	if err == nil && Syariah(b.Name) {
		return models.BisnisGrup{}, repository.ErrTidakAda
	}
	return b, err
}

// Buka - satu baris.
func (l *Layanan) Buka(ctx context.Context, id string) (models.BisnisGrup, error) {
	b, err := l.ambil(ctx, nil, strings.TrimSpace(id))
	if errors.Is(err, repository.ErrTidakAda) {
		return models.BisnisGrup{}, ErrTidakAda
	}
	return b, err
}

// Rapikan - spasi tepi dibuang; Name dan Alias Name huruf besar (keputusan work owner 05-10-2026); Alias kosong = Name.
func Rapikan(i models.Isian) models.Isian {
	r := models.Isian{
		ID:    strings.TrimSpace(i.ID),
		TopID: strings.TrimSpace(i.TopID),
		Name:  strings.ToUpper(strings.TrimSpace(i.Name)),
		Alias: strings.ToUpper(strings.TrimSpace(i.Alias)),
	}
	if r.Alias == "" {
		r.Alias = r.Name
	}
	return r
}

func periksa(i models.Isian) error {
	var pesan []string
	if i.TopID == "" {
		pesan = append(pesan, "Treaty Group is required")
	}
	if i.Name == "" {
		pesan = append(pesan, "Name is required")
	}
	if Syariah(i.Name) {
		pesan = append(pesan, "Business groups ending with SYARIAH are not managed here")
	}
	for _, c := range []struct{ nama, nilai string }{{"Name", i.Name}, {"Alias Name", i.Alias}} {
		if len(c.nilai) > models.BatasTeks {
			pesan = append(pesan, fmt.Sprintf("%s is longer than %d characters", c.nama, models.BatasTeks))
		}
	}
	if len(pesan) > 0 {
		return tolak("%s", strings.Join(pesan, "; "))
	}
	return nil
}

// Simpan - Add (ID kosong) atau Edit. Menjawab baris yang tersimpan.
func (l *Layanan) Simpan(ctx context.Context, a Aktor, isi models.Isian) (models.BisnisGrup, error) {
	if err := wajibPenuh(a); err != nil {
		return models.BisnisGrup{}, err
	}
	isi = Rapikan(isi)
	if err := periksa(isi); err != nil {
		return models.BisnisGrup{}, err
	}
	id := isi.ID
	err := l.tx(ctx, func(tx *dbTx) error {
		var b models.BisnisGrup
		if id != "" {
			lama, err := l.ambil(ctx, tx, id)
			if err != nil {
				return err
			}
			b = lama
		}
		// Treaty Group tidak diganti (Edit): TREATYNAME dibiarkan ("jgn ada ubah data"). Selainnya disalin dari TREATYGROUP.
		if id == "" || isi.TopID != b.TopID {
			t, err := l.gudang.AmbilTreaty(ctx, tx, isi.TopID)
			if errors.Is(err, repository.ErrTidakAda) {
				return tolak("Treaty Group %s is not in the treaty group list", isi.TopID)
			}
			if err != nil {
				return err
			}
			if len(t.Name) > models.BatasTeks {
				return tolak("the name of Treaty Group %s is longer than %d characters", t.ID, models.BatasTeks)
			}
			b.TopID, b.TreatyName = t.ID, t.Name
		}
		ganda, err := l.gudang.PemakaiNama(ctx, tx, isi.Name, id)
		if err != nil {
			return err
		}
		if len(ganda) > 0 {
			return tolak("Name is already used by %s", ganda[0])
		}
		b.Name, b.Alias = isi.Name, isi.Alias
		if id != "" {
			return l.gudang.Ubah(ctx, tx, b)
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
				id, b.ID = baru, baru
				break
			}
		}
		return l.gudang.Sisip(ctx, tx, b)
	})
	if errors.Is(err, repository.ErrTidakAda) {
		return models.BisnisGrup{}, ErrTidakAda
	}
	if err != nil {
		return models.BisnisGrup{}, err
	}
	return l.Buka(ctx, id)
}
