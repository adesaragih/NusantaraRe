package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"nusantarare/modul/reinsurancetype/backend/models"
	"nusantarare/modul/reinsurancetype/backend/repository"
)

// cobaIDMaksimum - ID baru yang dicoba sebelum menyerah bila nomor sequence sudah terpakai baris lama.
const cobaIDMaksimum = 50

var (
	polaAngka = regexp.MustCompile(`^[0-9]+$`)
	wib       = time.FixedZone("WIB", 7*60*60)
)

// FormatWaktuPega - bentuk TGLUPDATE yang ditulis Pega (`20261005T102343.866 GMT`).
func FormatWaktuPega(t time.Time) string { return t.UTC().Format("20060102T150405.000") + " GMT" }

// TampilWaktu - teks Pega ke `DD-MM-YYYY HH:MM` WIB; bentuk lain apa adanya.
func TampilWaktu(s string) string {
	t, err := time.Parse("20060102T150405.000 MST", strings.TrimSpace(s))
	if err != nil {
		return s
	}
	return t.In(wib).Format("02-01-2006 15:04")
}

func lengkapi(j models.Jenis) models.Jenis {
	j.Diubah = TampilWaktu(j.TglUpdate)
	return j
}

// Daftar - baris bersaring kata, Type, dan Flag ("" = semua).
func (l *Layanan) Daftar(ctx context.Context, kata, tipe, flag string) ([]models.Jenis, error) {
	d, err := l.gudang.Daftar(ctx, kata, tipe, flag)
	for i := range d {
		d[i] = lengkapi(d[i])
	}
	return d, err
}

// Buka - satu baris.
func (l *Layanan) Buka(ctx context.Context, id string) (models.Jenis, error) {
	j, err := l.gudang.Ambil(ctx, nil, strings.TrimSpace(id))
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Jenis{}, ErrTidakAda
	}
	return lengkapi(j), err
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

// Rapikan - spasi tepi dibuang; Name, SOA Name, Group Type huruf besar; Flag huruf kecil; Code kosong = `00`.
func Rapikan(i models.Isian) models.Isian {
	r := models.Isian{
		ID:        strings.TrimSpace(i.ID),
		Name:      strings.ToUpper(strings.TrimSpace(i.Name)),
		Type:      strings.TrimSpace(i.Type),
		SoaName:   strings.ToUpper(strings.TrimSpace(i.SoaName)),
		Code:      strings.TrimSpace(i.Code),
		Flag:      strings.ToLower(strings.TrimSpace(i.Flag)),
		NoUrut:    strings.TrimSpace(i.NoUrut),
		GroupType: strings.ToUpper(strings.TrimSpace(i.GroupType)),
	}
	if r.Code == "" {
		r.Code = models.CodeBawaan
	}
	return r
}

// sahAtauLama - nilai baru harus salah satu pilihan; nilai warisan lain boleh tetap selama tidak diubah (Edit).
func sahAtauLama(nilai string, pilihan []string, lama *models.Jenis, ambil func(models.Jenis) string) bool {
	return slices.Contains(pilihan, nilai) || (lama != nil && nilai == ambil(*lama))
}

func periksa(i models.Isian, lama *models.Jenis) error {
	var pesan []string
	if i.Name == "" {
		pesan = append(pesan, "Name is required")
	}
	if !sahAtauLama(i.Type, models.Type, lama, func(j models.Jenis) string { return j.Type }) {
		pesan = append(pesan, "Type must be 1 (Own Retention), 2 (Treaty Out), 3 (Facultative), or 4 (Treaty In)")
	}
	if !sahAtauLama(i.Flag, models.Flag, lama, func(j models.Jenis) string { return j.Flag }) {
		pesan = append(pesan, "Flag must be active or inactive")
	}
	if i.GroupType != "" && !sahAtauLama(i.GroupType, models.GroupType, lama, func(j models.Jenis) string { return j.GroupType }) {
		pesan = append(pesan, "Group Type must be OR, QS, RI, SPL, or empty")
	}
	if !polaAngka.MatchString(i.Code) {
		pesan = append(pesan, "Code must be a number")
	}
	if i.NoUrut != "" && !polaAngka.MatchString(i.NoUrut) {
		pesan = append(pesan, "No Urut must be a number")
	}
	for _, c := range []struct{ nama, nilai string }{{"Name", i.Name}, {"SOA Name", i.SoaName}, {"Code", i.Code}, {"No Urut", i.NoUrut}} {
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
func (l *Layanan) Simpan(ctx context.Context, a Aktor, isi models.Isian) (models.Jenis, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Jenis{}, err
	}
	isi = Rapikan(isi)
	id := isi.ID
	err := l.tx(ctx, func(tx *dbTx) error {
		var lama *models.Jenis
		if id != "" {
			j, err := l.gudang.Ambil(ctx, tx, id)
			if err != nil {
				return err
			}
			lama = &j
		}
		if err := periksa(isi, lama); err != nil {
			return err
		}
		ganda, err := l.gudang.PemakaiNama(ctx, tx, isi.Name, id)
		if err != nil {
			return err
		}
		if len(ganda) > 0 {
			return tolak("Name %s is already used by %s", isi.Name, ganda[0].ID)
		}
		j := models.Jenis{ID: id, Name: isi.Name, Type: isi.Type, SoaName: isi.SoaName, Code: isi.Code, Flag: isi.Flag,
			NoUrut: isi.NoUrut, GroupType: isi.GroupType, UserID: potongByte(a.AkunID, models.BatasTeks),
			TglUpdate: FormatWaktuPega(l.jam())}
		if id != "" {
			return l.gudang.Ubah(ctx, tx, j)
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
				id, j.ID = baru, baru
				break
			}
		}
		return l.gudang.Sisip(ctx, tx, j)
	})
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Jenis{}, ErrTidakAda
	}
	if errors.Is(err, repository.ErrKembar) {
		return models.Jenis{}, tolak("Name %s is already used", isi.Name)
	}
	if err != nil {
		return models.Jenis{}, err
	}
	return l.Buka(ctx, id)
}
