package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"nusantarare/modul/treatygroup/backend/models"
	"nusantarare/modul/treatygroup/backend/repository"
)

// cobaIDMaksimum - ID baru yang dicoba sebelum menyerah bila nomor sequence sudah terpakai baris lama.
const cobaIDMaksimum = 50

// wib - tampilan waktu ubah (TGLUPDATE disimpan GMT seperti Pega).
var wib = time.FixedZone("WIB", 7*60*60)

// FormatWaktuPega - bentuk TGLUPDATE yang ditulis Pega (`20220105T135630.318 GMT`).
func FormatWaktuPega(t time.Time) string { return t.UTC().Format("20060102T150405.000") + " GMT" }

// TampilWaktu - TGLUPDATE format Pega ke `DD-MM-YYYY HH:MM` WIB; bentuk lain dikembalikan apa adanya.
func TampilWaktu(s string) string {
	t, err := time.Parse("20060102T150405.000 MST", strings.TrimSpace(s))
	if err != nil {
		return s
	}
	return t.In(wib).Format("02-01-2006 15:04")
}

func lengkapi(g models.Grup) models.Grup {
	g.Diubah = TampilWaktu(g.TglUpdate)
	return g
}

// Daftar - baris bersaring kata dan OJK ("" = semua).
func (l *Layanan) Daftar(ctx context.Context, kata, ojk string) ([]models.Grup, error) {
	d, err := l.gudang.Daftar(ctx, kata, ojk)
	for i := range d {
		d[i] = lengkapi(d[i])
	}
	return d, err
}

// Pilihan - isi dropdown OJK Business.
func (l *Layanan) Pilihan(ctx context.Context) ([]models.Ojk, error) { return l.gudang.DaftarOjk(ctx) }

// Buka - satu grup dan grup bisnis anaknya.
func (l *Layanan) Buka(ctx context.Context, id string) (models.Detail, error) {
	id = strings.TrimSpace(id)
	g, err := l.gudang.Ambil(ctx, nil, id)
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Detail{}, ErrTidakAda
	}
	if err != nil {
		return models.Detail{}, err
	}
	anak, err := l.gudang.Anak(ctx, id)
	if err != nil {
		return models.Detail{}, err
	}
	return models.Detail{Grup: lengkapi(g), Anak: anak}, nil
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

// Rapikan - spasi tepi dibuang; nama dan nama SOA huruf besar (keputusan work owner 05-10-2026).
func Rapikan(i models.Isian) models.Isian {
	return models.Isian{
		ID:      strings.TrimSpace(i.ID),
		OjkID:   strings.TrimSpace(i.OjkID),
		Name:    strings.ToUpper(strings.TrimSpace(i.Name)),
		SoaName: strings.ToUpper(strings.TrimSpace(i.SoaName)),
	}
}

func periksa(i models.Isian) error {
	var pesan []string
	if i.OjkID == "" {
		pesan = append(pesan, "OJK Business is required")
	}
	if i.Name == "" {
		pesan = append(pesan, "Treaty Group Name is required")
	}
	for _, c := range []struct{ nama, nilai string }{{"Treaty Group Name", i.Name}, {"SOA Name", i.SoaName}} {
		if len(c.nilai) > models.BatasTeks {
			pesan = append(pesan, fmt.Sprintf("%s is longer than %d characters", c.nama, models.BatasTeks))
		}
	}
	if len(pesan) > 0 {
		return tolak("%s", strings.Join(pesan, "; "))
	}
	return nil
}

// salinOjk mengisi kolom OJK grup dari baris TREATYGROUPOJK; nilai yang tidak muat kolom TREATYGROUP ditolak.
func salinOjk(r *models.Grup, o models.Ojk) error {
	if len(o.Name) > models.BatasTeks || len(o.NameIDN) > models.BatasTeks || len(o.OrderNo) > models.BatasOrderNo {
		return tolak("OJK Business %s does not fit the treaty group table (Order No longer than %d characters)", o.ID, models.BatasOrderNo)
	}
	r.OjkID, r.OjkName, r.OjkNameIDN, r.OrderNo = o.ID, o.Name, o.NameIDN, o.OrderNo
	return nil
}

// Simpan - Add (ID kosong) atau Edit. Menjawab grup yang tersimpan.
func (l *Layanan) Simpan(ctx context.Context, a Aktor, isi models.Isian) (models.Detail, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Detail{}, err
	}
	isi = Rapikan(isi)
	if err := periksa(isi); err != nil {
		return models.Detail{}, err
	}
	id := isi.ID
	err := l.tx(ctx, func(tx *dbTx) error {
		var r models.Grup
		if id != "" {
			lama, err := l.gudang.Ambil(ctx, tx, id)
			if err != nil {
				return err
			}
			r = lama
		}
		// OJK tidak diganti (Edit): salinannya dan COAID dibiarkan ("jgn ada ubah data"). Selainnya disalin dari
		// TREATYGROUPOJK, dan COAID ikut OJK itu ("ubah pas pilih OJK Business") - sama dengan yang ditampilkan form.
		if id == "" || isi.OjkID != r.OjkID {
			o, err := l.gudang.AmbilOjk(ctx, tx, isi.OjkID)
			if errors.Is(err, repository.ErrTidakAda) {
				return tolak("OJK Business %s is not in the treaty group OJK list", isi.OjkID)
			}
			if err != nil {
				return err
			}
			if err := salinOjk(&r, o); err != nil {
				return err
			}
			if r.CoaID, err = l.gudang.CoaSeOjk(ctx, tx, r.OjkID); err != nil {
				return err
			}
		}
		ganda, err := l.gudang.PemakaiNama(ctx, tx, isi.Name, id)
		if err != nil {
			return err
		}
		if len(ganda) > 0 {
			return tolak("Treaty Group Name is already used by %s", ganda[0])
		}
		r.Name, r.SoaName = isi.Name, isi.SoaName
		r.TglUpdate, r.UserID = FormatWaktuPega(l.jam()), potongByte(a.AkunID, models.BatasUserID)
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
		return models.Detail{}, ErrTidakAda
	}
	if err != nil {
		return models.Detail{}, err
	}
	return l.Buka(ctx, id)
}
