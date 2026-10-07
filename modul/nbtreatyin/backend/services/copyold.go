package services

// Untuk apa berkas ini: COPY OLD - perintah work owner 07-10-2026 ("nb ttreatyin tobol copy untuk data lama mana?" ->
// "langsung anda kerjakan!"; dikerjakan sesi EDM TREATY IN dengan izin WO). Tombol di samping Create membuka popup berisi
// dokumen polis NB Treaty In lama (`POOLDATA.JSON_POLIS` generasi 0) yang BELUM ada di tabel flat; yang dicentang
// disalin lewat `Process Copy`.
//
//   - aturan salin = pemuat dokumen lama (`pemuat.go` `muat`, tiket 22): satu transaksi per dokumen lewat antarmuka
//     penyimpanan jalur biasa; satu dokumen yang gagal tidak membatalkan yang lain;
//   - hanya superadmin (pemegang menu Kelola User) dengan menu NB Treaty In ber-hak PENUH - pola Copy Old Data
//     Bordereaux (View only berlaku juga bagi superadmin, 05-10-2026) dan Copy Old EDM Treaty In;
//   - JSON_POLIS hanya DIBACA. Pesan layar dari `models.AlasanSalinLama` - tanpa nilai dokumen.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbtreatyin/backend/models"
)

// MaksSalinLama - ID paling banyak dalam satu `Process Copy`.
const MaksSalinLama = 500

// ErrBukanSuperadmin - Copy Old dibuka akun yang bukan superadmin, atau menu NB Treaty In-nya View only.
var ErrBukanSuperadmin = fmt.Errorf("%w: Copy Old is only for super admin with full access to NB Treaty In", inti.ErrTanpaWewenang)

// Hak - hak layar portal yang bergantung pada akun.
type Hak struct {
	// CopyOld - tombol Copy Old tampil.
	CopyOld bool `json:"copyOld"`
}

// HakPortal - hak layar portal. `superadmin` = pemegang Kelola User dengan menu NB Treaty In ber-hak penuh (dinilai
// handlers dari menu sesi). Tanpa pembaca JSON_POLIS (gudang tiruan) tombol tidak tampil.
func (l *Layanan) HakPortal(p inti.Pelaku, superadmin bool) Hak {
	return Hak{CopyOld: l.wajibCopyOld(p, superadmin) == nil}
}

func (l *Layanan) wajibCopyOld(p inti.Pelaku, superadmin bool) error {
	if err := inti.WajibIdentitas(p); err != nil {
		return err
	}
	if !superadmin {
		return ErrBukanSuperadmin
	}
	if l.pemuat == nil {
		return ErrTanpaOracle
	}
	return nil
}

// DaftarDokumenLama - isi popup Copy Old.
func (l *Layanan) DaftarDokumenLama(ctx context.Context, p inti.Pelaku, superadmin bool) ([]models.DokumenLama, error) {
	if err := l.wajibCopyOld(p, superadmin); err != nil {
		return nil, err
	}
	d, _, err := l.pemuat.siapkanLama(ctx)
	return d, err
}

// SalinDokumenLama - `Process Copy`: ID yang dicentang disalin satu per satu menurut urutan daftar.
func (l *Layanan) SalinDokumenLama(ctx context.Context, p inti.Pelaku, superadmin bool, ids []string) (models.JawabanSalinLama, error) {
	if err := l.wajibCopyOld(p, superadmin); err != nil {
		return models.JawabanSalinLama{Hasil: []models.HasilSalinLama{}}, err
	}
	return l.pemuat.salinLama(ctx, ids)
}

// dokumenSiap - dokumen yang boleh disalin: baris JSON_POLIS (IDPEGA untuk SuggestList) dan hasil pecahnya.
type dokumenSiap struct {
	b models.BarisJSONPolis
	h models.HasilPecah
}

// siapkanLama - baris popup (urutan IDPEGA) dan dokumen yang boleh disalin. SELECT saja.
//
//	sumber       JSON_POLIS generasi 0 yang IDPEGA-nya ada di tabel kerja Pega DAN di TREATYINPRODUCTION
//	             (perintah WO 07-10-2026, `repository.sqlKunciJSONPolisCopyOld`)
//	lewat        dokumen lini lain / generasi endorsemen (milik EDM) / baris json_polis tulisan aplikasi baru
//	tidak tampil kasus sudah di tabel flat (`AdaKasus`)
//	ditolak      galat pemecah, dokumen ganda
func (pm *Pemuat) siapkanLama(ctx context.Context) ([]models.DokumenLama, map[string]dokumenSiap, error) {
	kunci, err := pm.g.KunciJSONPolisCopyOld(ctx)
	if err != nil {
		return nil, nil, err
	}
	out := []models.DokumenLama{}
	siap := map[string]dokumenSiap{}
	terlihat := map[string]bool{}
	for _, k := range kunci {
		if err := ctx.Err(); err != nil {
			return nil, nil, err
		}
		b, err := pm.g.BacaJSONPolis(ctx, k)
		if err != nil {
			return nil, nil, err
		}
		h, err := models.PecahDokumenLama(b)
		if errors.Is(err, models.ErrBukanTreatyIn) || errors.Is(err, models.ErrGenerasiEndorsemen) ||
			errors.Is(err, models.ErrBarisAplikasiBaru) {
			continue
		}
		d := models.DokumenLama{NoPolis: strings.TrimSpace(b.NoPolis), TglProd: b.TglProd, Alasan: []string{}}
		d.ID, _ = models.IDKasusDariIDPega(b.IDPega)
		if err != nil {
			d.Alasan = append(d.Alasan, teksAlasan(err))
			out = append(out, d)
			continue
		}
		d.ID, d.NoPolis = h.ID, h.NoPolis
		d.NoOffer = h.Halaman.Ambil(models.HalamanPolis + ".NoOffer")
		d.InsuredName = models.NilaiQuotation(h.Halaman, "InsuredName")
		if d.InsuredName == "" {
			d.InsuredName = h.Halaman.Ambil(models.HalamanPolis + ".InsuredName")
		}
		d.BusinessName = models.NilaiQuotation(h.Halaman, "BusinessName")
		d.SOBName = h.Halaman.Ambil(models.HalamanPolis + ".SOBName")
		d.CedingCoName = h.Halaman.Ambil(models.HalamanPolis + ".CedingCoName")
		for _, g := range h.Galat {
			d.Alasan = append(d.Alasan, g.Jalur+": "+teksAlasan(g.Err))
		}
		if terlihat[h.ID] {
			d.Alasan = append(d.Alasan, models.AlasanSalinLama(models.ErrDokumenGanda))
		}
		terlihat[h.ID] = true
		if len(d.Alasan) == 0 {
			ada, err := pm.g.AdaKasus(ctx, nil, h.ID)
			if err != nil {
				return nil, nil, err
			}
			if ada {
				continue // sudah di tabel flat
			}
		}
		d.BolehDisalin = len(d.Alasan) == 0
		if d.BolehDisalin {
			siap[h.ID] = dokumenSiap{b: b, h: h}
		}
		out = append(out, d)
	}
	return out, siap, nil
}

// teksAlasan - teks layar galat dokumen; galat lain tanpa rincian (teksnya dapat memuat nilai dokumen).
func teksAlasan(err error) string {
	if t := models.AlasanSalinLama(err); t != "" {
		return t
	}
	return "the old document cannot be copied"
}

// salinLama - `Process Copy`. Galat basis data satu dokumen dicatat log (nomor kasus saja) dan dilaporkan `gagal`
// untuk dokumen itu saja.
func (pm *Pemuat) salinLama(ctx context.Context, ids []string) (models.JawabanSalinLama, error) {
	j := models.JawabanSalinLama{Hasil: []models.HasilSalinLama{}}
	minta := map[string]bool{}
	var unik []string
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" || minta[id] {
			continue
		}
		minta[id] = true
		unik = append(unik, id)
	}
	switch {
	case len(unik) == 0:
		return j, &GalatValidasi{Pesan: []string{"select at least one old document to copy"}}
	case len(unik) > MaksSalinLama:
		return j, &GalatValidasi{Pesan: []string{fmt.Sprintf("at most %d old documents can be copied at once", MaksSalinLama)}}
	}
	daftar, siap, err := pm.siapkanLama(ctx)
	if err != nil {
		return j, err
	}
	diproses := map[string]bool{}
	for _, d := range daftar {
		if d.ID == "" || !minta[d.ID] || diproses[d.ID] {
			continue
		}
		diproses[d.ID] = true
		h := models.HasilSalinLama{ID: d.ID, Pesan: []string{}}
		if !d.BolehDisalin {
			h.Status, h.Pesan = models.SalinDitolak, d.Alasan
			j.Hasil = append(j.Hasil, h)
			continue
		}
		s := siap[d.ID]
		_, err := pm.muat(ctx, s.b, s.h)
		switch {
		case err == nil:
			h.Status = models.SalinDisalin
			j.Disalin++
		case errors.Is(err, errSudahDimuat):
			h.Status, h.Pesan = models.SalinSudahAda, []string{"already in the new tables"}
		case errors.Is(err, ErrNomorPolisSudahAda):
			h.Status, h.Pesan = models.SalinDitolak, []string{"this policy number is already used by another NB case in the new tables"}
		default:
			log.Printf("nbtreatyin: Copy Old %s gagal: %v", d.ID, err)
			h.Status, h.Pesan = models.SalinGagal, []string{"database error - nothing was written for this document"}
		}
		j.Hasil = append(j.Hasil, h)
	}
	for _, id := range unik {
		if !diproses[id] {
			j.Hasil = append(j.Hasil, models.HasilSalinLama{ID: id, Status: models.SalinSudahAda,
				Pesan: []string{"already in the new tables, or not an old NB Treaty In document"}})
		}
	}
	log.Printf("nbtreatyin: Copy Old menyalin %d dari %d dokumen lama", j.Disalin, len(unik))
	return j, nil
}
