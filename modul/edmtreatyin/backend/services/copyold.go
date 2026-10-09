package services

// Untuk apa berkas ini: COPY OLD - perintah work owner 07-10-2026 ("BUATKAN TOMBOL COPY OLD SAMA SEPERTI MASTER
// PRODUCTNAME LIFE, KHUSUS BUAT SUPERUSER"). Tombol di samping Create membuka popup berisi dokumen endorsemen Treaty In
// lama (`POOLDATA.JSON_POLIS` ber-PRODKE > 0) yang BELUM ada di tabel flat; yang dicentang disalin lewat `Process Copy`.
//
//   - aturan salin = pemuat dokumen lama (`pemuat.go` `muat`, tiket 09/10): satu transaksi per dokumen lewat antarmuka
//     penyimpanan jalur biasa - penjaga percabangan, keutuhan, dan proyeksi selisih 'PEGA' ikut berjalan; satu
//     dokumen yang gagal tidak membatalkan yang lain;
//   - urutan GENERASI (NOPOLIS lalu PRODKE), bukan urutan centang - generasi 1 tersalin sebelum generasi 2;
//   - hanya superadmin (pemegang menu Kelola User) dengan menu EDM Treaty In ber-hak PENUH - pola Copy Old Data
//     Bordereaux (04-10-2026; View only berlaku juga bagi superadmin, 05-10-2026);
//   - JSON_POLIS hanya DIBACA. Pesan layar dari `models.AlasanSalinLama` - tanpa nomor polis / nilai dokumen.
//
// Pola tampilan: Copy Old `modul/masterproductnamelife` (03-10-2026).

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/edmtreatyin/backend/models"
)

// MaksSalinLama - ID paling banyak dalam satu `Process Copy`.
const MaksSalinLama = 500

// ErrBukanSuperadmin - Copy Old dibuka akun yang bukan superadmin, atau menu EDM Treaty In-nya View only.
var ErrBukanSuperadmin = fmt.Errorf("%w: Copy Old is only for super admin with full access to EDM Treaty In", inti.ErrTanpaWewenang)

// Hak - hak layar portal yang bergantung pada akun.
type Hak struct {
	// CopyOld - tombol Copy Old tampil.
	CopyOld bool `json:"copyOld"`
}

// HakPortal - hak layar portal. `superadmin` = pemegang Kelola User dengan menu EDM Treaty In ber-hak penuh
// (dinilai handlers dari menu sesi). Tanpa pembaca JSON_POLIS (gudang tiruan) tombol tidak tampil.
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

// SalinDokumenLama - `Process Copy`: ID yang dicentang disalin satu per satu menurut urutan generasi.
func (l *Layanan) SalinDokumenLama(ctx context.Context, p inti.Pelaku, superadmin bool, ids []string) (models.JawabanSalinLama, error) {
	if err := l.wajibCopyOld(p, superadmin); err != nil {
		return models.JawabanSalinLama{Hasil: []models.HasilSalinLama{}}, err
	}
	return l.pemuat.salinLama(ctx, ids)
}

// kunciGen - kunci generasi (NOPOLIS, PRODKE).
func kunciGen(nopolis string, prodke int) string { return nopolis + "\x00" + strconv.Itoa(prodke) }

// siapkanLama - baris popup (urutan generasi) dan hasil pecah dokumen yang boleh disalin. SELECT saja.
//
//	sumber       JSON_POLIS generasi endorsemen yang IDPEGA-nya ada di tabel kerja Pega DAN di TREATYINPRODUCTION
//	             (perintah WO 07-10-2026, `repository.sqlPmKunciJSONPolisEDMCopyOld`)
//	lewat        dokumen lini lain / generasi NB (milik pemuat NB) / baris json_polis tulisan aplikasi baru - tidak
//	             tampil
//	tidak tampil generasi sudah di tabel flat dengan kunci sama (`KunciGenerasiLama`)
//	ditolak      galat pemecah, dokumen ganda, ID kasus bentrok, generasi bernomor sama sudah ada (percabangan),
//	             generasi sebelumnya tidak ada di tabel flat dan tidak di antara dokumen yang boleh disalin
func (pm *Pemuat) siapkanLama(ctx context.Context) ([]models.DokumenLama, map[string]models.HasilPecahEDM, error) {
	kunci, err := pm.g.KunciJSONPolisEDMCopyOld(ctx)
	if err != nil {
		return nil, nil, err
	}
	models.UrutKunciGenerasi(kunci)
	out := []models.DokumenLama{}
	siap := map[string]models.HasilPecahEDM{}
	terlihat := map[string]bool{}
	calon := map[string]bool{} // generasi yang boleh disalin di daftar ini
	for _, k := range kunci {
		if err := ctx.Err(); err != nil {
			return nil, nil, err
		}
		b, err := pm.g.BacaJSONPolisEDM(ctx, k.Kunci)
		if err != nil {
			return nil, nil, err
		}
		h, err := models.PecahDokumenEDM(b)
		if errors.Is(err, models.ErrBukanTreatyIn) || errors.Is(err, models.ErrBukanGenerasiEndorsemen) ||
			errors.Is(err, models.ErrBarisAplikasiBaru) {
			continue
		}
		d := models.DokumenLama{NoPolis: strings.TrimSpace(b.NoPolis), EDMNo: strings.TrimSpace(b.NoEndors), TglProd: b.TglProd,
			Alasan: []string{}}
		d.ID = strings.TrimSpace(b.IDPega) // IDPEGA utuh (WO 07-10-2026)
		d.ProdKe, _ = strconv.Atoi(strings.TrimSpace(b.ProdKe))
		if err != nil {
			out = append(out, tolakLama(d, err))
			continue
		}
		d.ID, d.NoPolis, d.EDMNo, d.ProdKe, d.EDMType = h.ID, h.NoPolis, h.EDMNo, h.ProdKe, h.EDMType
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
			kg, ada, err := pm.g.KunciGenerasiLama(ctx, nil, h.ID)
			if err != nil {
				return nil, nil, err
			}
			if ada && kg.NoPolis == h.NoPolis && kg.ProdKe == h.ProdKe && kg.NoEndors == h.EDMNo {
				continue // sudah di tabel flat
			}
			if ada {
				d.Alasan = append(d.Alasan, models.AlasanSalinLama(models.ErrIDKasusBentrok))
			}
		}
		if len(d.Alasan) == 0 {
			if alasan, err := pm.alasanRantai(ctx, h, calon); err != nil {
				return nil, nil, err
			} else if alasan != "" {
				d.Alasan = append(d.Alasan, alasan)
			}
		}
		d.BolehDisalin = len(d.Alasan) == 0
		if d.BolehDisalin {
			siap[h.ID] = h
			calon[kunciGen(h.NoPolis, h.ProdKe)] = true
		}
		out = append(out, d)
	}
	return out, siap, nil
}

// alasanRantai - percabangan (generasi bernomor sama sudah di tabel flat atau sudah diambil dokumen lain di daftar)
// dan generasi sebelumnya (di tabel flat, atau dokumen sebelumnya di daftar yang boleh disalin).
func (pm *Pemuat) alasanRantai(ctx context.Context, h models.HasilPecahEDM, calon map[string]bool) (string, error) {
	sama, err := pm.g.GenerasiSebelumnya(ctx, nil, models.Kasus{ProdKe: h.ProdKe + 1}, h.NoPolis)
	if err != nil {
		return "", err
	}
	if sama != "" || calon[kunciGen(h.NoPolis, h.ProdKe)] {
		return models.AlasanSalinLama(models.ErrPercabangan), nil
	}
	if calon[kunciGen(h.NoPolis, h.ProdKe-1)] {
		return "", nil
	}
	lama, err := pm.g.GenerasiSebelumnya(ctx, nil, models.Kasus{ProdKe: h.ProdKe}, h.NoPolis)
	if err != nil {
		return "", err
	}
	if lama == "" {
		return models.AlasanSalinLama(models.ErrGenerasiSebelumnyaTidakAda), nil
	}
	return "", nil
}

// teksAlasan - teks layar galat dokumen; galat lain tanpa rincian (teksnya dapat memuat nilai dokumen).
func teksAlasan(err error) string {
	if t := models.AlasanSalinLama(err); t != "" {
		return t
	}
	return "the old document cannot be copied"
}

func tolakLama(d models.DokumenLama, err error) models.DokumenLama {
	d.Alasan = append(d.Alasan, teksAlasan(err))
	return d
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
		_, _, err := pm.muat(ctx, siap[d.ID])
		switch {
		case err == nil:
			h.Status = models.SalinDisalin
			j.Disalin++
		case errors.Is(err, pmErrSudahDimuat):
			h.Status, h.Pesan = models.SalinSudahAda, []string{"already in the new tables"}
		case models.AlasanSalinLama(err) != "":
			h.Status, h.Pesan = models.SalinDitolak, []string{models.AlasanSalinLama(err)}
		default:
			log.Printf("edmtreatyin: Copy Old %s gagal: %v", d.ID, err)
			h.Status, h.Pesan = models.SalinGagal, []string{"database error - nothing was written for this document"}
		}
		j.Hasil = append(j.Hasil, h)
	}
	for _, id := range unik {
		if !diproses[id] {
			j.Hasil = append(j.Hasil, models.HasilSalinLama{ID: id, Status: models.SalinSudahAda,
				Pesan: []string{"already in the new tables, or not an old EDM Treaty In document"}})
		}
	}
	log.Printf("edmtreatyin: Copy Old menyalin %d dari %d dokumen lama", j.Disalin, len(unik))
	return j, nil
}
