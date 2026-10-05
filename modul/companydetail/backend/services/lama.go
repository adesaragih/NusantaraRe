package services

// Copy Old - perintah work owner 04-10-2026: "buatkan fungsinya copy old seperti pada productname; copy old hanya
// muncul untuk superadmin". Popup berisi organisasi dokumen Pega (`M_CLIENT`) yang datanya belum pindah ke tabel
// datar; yang dicentang disalin lewat `Process Copy`, SATU TRANSAKSI PER ORGANISASI - satu yang gagal tidak
// membatalkan yang lain. Aturan salin SAMA dengan alat pindah. Dokumen hanya DIBACA.
//
// Superadmin = akun pemegang menu Kelola User (`kelolauser`, admin aplikasi): aplikasi tidak punya peran bernama
// superadmin (DEV 04-10-2026: M_WORKBASKET hanya Reas*Admin per lini). Ditegakkan di sini, bukan hanya di layar.

import (
	"context"
	"errors"
	"log"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/menu"
	"nusantarare/modul/companydetail/backend/models"
	"nusantarare/modul/companydetail/backend/repository"
)

// ErrBukanSuperadmin - Copy Old diminta akun yang bukan pemegang Kelola User (403).
var ErrBukanSuperadmin = errors.New("copy old is only for super admin")

// MaksSalinLama - ID paling banyak dalam satu permintaan `Process Copy` (layar mengirim per 20).
const MaksSalinLama = 1000

// Hak - tombol yang boleh tampil bagi akun ini.
type Hak struct {
	CopyOld bool `json:"copyOld"`
}

// superadmin - akun hasil login yang memegang menu Kelola User.
func superadmin(ctx context.Context) bool {
	kode, ada := inti.AksesMenuDari(ctx)
	return ada && inti.PunyaMenu(kode, menu.KodeKelolaUser)
}

// HakAkun - hak tombol akun yang meminta.
func (l *Layanan) HakAkun(ctx context.Context) Hak { return Hak{CopyOld: superadmin(ctx)} }

func wajibSuperadmin(ctx context.Context, p inti.Pelaku) error {
	if strings.TrimSpace(p.AkunID) == "" {
		return ErrTanpaPelaku
	}
	if !superadmin(ctx) {
		return ErrBukanSuperadmin
	}
	return nil
}

// DaftarLama - isi popup Copy Old.
func (l *Layanan) DaftarLama(ctx context.Context, p inti.Pelaku) ([]models.OrgLama, error) {
	if err := wajibSuperadmin(ctx, p); err != nil {
		return nil, err
	}
	return l.gudang.SiapkanOrgLama(ctx)
}

// SalinLama - `Process Copy`: ID yang dicentang disalin satu per satu. Galat basis data satu organisasi dicatat log
// dan dilaporkan `gagal` untuk organisasi itu saja.
func (l *Layanan) SalinLama(ctx context.Context, p inti.Pelaku, ids []string) (models.JawabanSalinLama, error) {
	j := models.JawabanSalinLama{Hasil: []models.HasilSalinLama{}}
	if err := wajibSuperadmin(ctx, p); err != nil {
		return j, err
	}
	var unik []string
	dilihat := map[string]bool{}
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" || dilihat[id] {
			continue
		}
		dilihat[id] = true
		unik = append(unik, id)
	}
	switch {
	case len(unik) == 0:
		return j, tolak("select at least one old organization to copy")
	case len(unik) > MaksSalinLama:
		return j, tolak("at most %d old organizations can be copied at once", MaksSalinLama)
	}
	// Title dan negara dibaca SEKALI per Process Copy (dulu sekali per organisasi - Select all ratusan baris melewati
	// batas waktu layar 30 detik: "signal is aborted without reason", 04-10-2026). Layar mengirim per 20 ID.
	r, err := l.gudang.ReferensiPindah(ctx)
	if err != nil {
		return j, err
	}
	for _, id := range unik {
		h := models.HasilSalinLama{ID: id, Pesan: []string{}}
		var ditulis bool
		err := l.tx(ctx, func(tx *dbTx) error {
			var err error
			ditulis, err = l.gudang.SalinOrgLama(ctx, tx, id, r)
			return err
		})
		switch {
		case errors.Is(err, repository.ErrTidakAda):
			h.Status, h.Pesan = models.SalinDitolak, []string{"not an old organization document"}
		case errors.Is(err, repository.ErrBelumDimigrasi):
			return j, err
		case err != nil:
			log.Printf("company detail: copy old %s: %v", id, err)
			h.Status, h.Pesan = models.SalinGagal, []string{"database error; details in the server log"}
		case ditulis:
			h.Status = models.SalinDisalin
			j.Disalin++
		default:
			h.Status, h.Pesan = models.SalinSudahAda, []string{"already in the new tables"}
		}
		j.Hasil = append(j.Hasil, h)
	}
	return j, nil
}
