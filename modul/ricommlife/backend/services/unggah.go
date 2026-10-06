package services

// Upload CSV - "View Upload" (pratinjau, tanpa menulis) dan "Simpan Upload" (urai ULANG di server, tolak seluruhnya
// bila ada galat, lalu satu transaksi). Rule Pega `UploadCSV_RICOMM`, `ViewCSVResult_RIComm`, `SubmitRIComm_Act` tidak
// tersedia: dirancang dari label `Format excel : USEDBY, CONTRACT, YEAR, COMM` b5569, gaya riratelife (butir 7).

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"nusantarare/modul/ricommlife/backend/models"
)

// MaksKalimatGalat - kalimat galat per baris paling banyak di pesan 422 Simpan Upload.
const MaksKalimatGalat = 50

// PermintaanUnggah - isi berkas CSV (teks).
type PermintaanUnggah struct {
	CSV string `json:"csv"`
}

// RingkasanUnggah - satu R/I COMM NAME di berkas: ringkasan yang cocok (ID) atau baru.
type RingkasanUnggah struct {
	UsedBy string `json:"usedby"`
	// ID - ringkasan yang cocok; kosong = dibuat saat Simpan Upload (Pratinjau) / ID barunya (Simpan Upload).
	ID     string `json:"id"`
	Baru   bool   `json:"baru"`
	Jumlah int    `json:"jumlah"`
}

// HasilUnggah - View Upload: baris sah, galat per baris, dan ringkasan tujuan.
type HasilUnggah struct {
	Baris     []models.BarisCSV   `json:"baris"`
	Galat     []models.GalatBaris `json:"galat"`
	Ringkasan []RingkasanUnggah   `json:"ringkasan"`
	// Sah - nol galat: Simpan Upload akan diterima (selama data tidak berubah).
	Sah bool `json:"sah"`
}

// HasilSimpanUnggah - Simpan Upload.
type HasilSimpanUnggah struct {
	Disimpan      int               `json:"disimpan"`
	RingkasanBaru int               `json:"ringkasanBaru"`
	Ringkasan     []RingkasanUnggah `json:"ringkasan"`
}

// kelompok - baris berkas per R/I COMM NAME (urutan kemunculan pertama).
type kelompok struct {
	nama  string
	baris []models.BarisCSV
	// cocok - ringkasan yang sudah ada (nil = baru).
	cocok *models.Ringkasan
}

// periksaUnggah - urai, cocokkan nama ke ringkasan, dan periksa kembar (CONTRACT, YEAR) terhadap rincian yang ada.
func (l *Layanan) periksaUnggah(ctx context.Context, tx *dbTx, teks string) (HasilUnggah, []*kelompok, error) {
	baris, galat, err := models.UraiCSV(teks)
	if errors.Is(err, models.ErrCSV) {
		return HasilUnggah{}, nil, tolak("%s", strings.TrimPrefix(err.Error(), models.ErrCSV.Error()+": "))
	}
	if err != nil {
		return HasilUnggah{}, nil, err
	}
	var urut []*kelompok
	per := map[string]*kelompok{}
	for _, b := range baris {
		k := models.KunciNama(b.UsedBy)
		if per[k] == nil {
			per[k] = &kelompok{nama: b.UsedBy}
			urut = append(urut, per[k])
		}
		per[k].baris = append(per[k].baris, b)
	}
	if len(urut) > models.MaksRingkasanCSV {
		return HasilUnggah{}, nil, tolak("the CSV file has %d different USEDBY values; at most %d can be uploaded at once",
			len(urut), models.MaksRingkasanCSV)
	}
	var ids []string
	sah := urut[:0:0]
	for _, k := range urut {
		ada, err := l.gudang.PemakaiNama(ctx, tx, k.nama, "")
		if err != nil {
			return HasilUnggah{}, nil, err
		}
		switch len(ada) {
		case 0:
		case 1:
			r := ada[0]
			k.cocok, k.nama = &r, r.UsedBy
			ids = append(ids, r.ID)
		default:
			var s []string
			for _, r := range ada {
				s = append(s, r.ID)
			}
			galat = append(galat, models.GalatBaris{Baris: k.baris[0].Baris, Pesan: fmt.Sprintf(
				"USEDBY %s matches more than one R/I comm summary (IDs %s); rename them first", k.nama, strings.Join(s, ", "))})
			continue
		}
		sah = append(sah, k)
	}
	lama, err := l.gudang.KomisiDari(ctx, tx, ids)
	if err != nil {
		return HasilUnggah{}, nil, err
	}
	sudah := map[string]string{}
	for _, r := range lama {
		sudah[strings.TrimSpace(r.IDUsedBy)+"\x01"+models.KunciKomisi(r.Contract, r.Year)] = r.ID
	}
	h := HasilUnggah{Baris: []models.BarisCSV{}}
	for _, k := range sah {
		ru := RingkasanUnggah{UsedBy: k.nama, Baru: k.cocok == nil}
		if k.cocok != nil {
			ru.ID = k.cocok.ID
		}
		var tetap []models.BarisCSV
		for _, b := range k.baris {
			b.UsedBy = k.nama
			if k.cocok != nil {
				if idLama, ada := sudah[k.cocok.ID+"\x01"+models.KunciKomisi(b.Contract, b.Year)]; ada {
					galat = append(galat, models.GalatBaris{Baris: b.Baris, Pesan: fmt.Sprintf(
						"CONTRACT %s, YEAR %s already exists in R/I COMM NAME %s (ID %s)", b.Contract, b.Year, k.nama, idLama)})
					continue
				}
			}
			tetap = append(tetap, b)
		}
		k.baris = tetap
		ru.Jumlah = len(tetap)
		h.Baris = append(h.Baris, tetap...)
		h.Ringkasan = append(h.Ringkasan, ru)
	}
	sort.SliceStable(galat, func(a, b int) bool { return galat[a].Baris < galat[b].Baris })
	sort.SliceStable(h.Baris, func(a, b int) bool { return h.Baris[a].Baris < h.Baris[b].Baris })
	if galat == nil {
		galat = []models.GalatBaris{}
	}
	if h.Ringkasan == nil {
		h.Ringkasan = []RingkasanUnggah{}
	}
	h.Galat, h.Sah = galat, len(galat) == 0
	return h, sah, nil
}

// Pratinjau - View Upload (`ViewCSVResult_RIComm`): urai dan periksa tanpa menulis.
func (l *Layanan) Pratinjau(ctx context.Context, a Aktor, p PermintaanUnggah) (HasilUnggah, error) {
	if err := wajibPenuh(a); err != nil {
		return HasilUnggah{}, err
	}
	h, _, err := l.periksaUnggah(ctx, nil, p.CSV)
	return h, err
}

// KalimatGalat - pesan 422 Simpan Upload: satu kalimat per baris, paling banyak MaksKalimatGalat.
func KalimatGalat(g []models.GalatBaris) string {
	var s []string
	for i, x := range g {
		if i == MaksKalimatGalat {
			s = append(s, fmt.Sprintf("And %d more rows.", len(g)-MaksKalimatGalat))
			break
		}
		s = append(s, x.Kalimat())
	}
	return "The upload was not saved. " + strings.Join(s, " ")
}

// SimpanUnggah - Simpan Upload (`SubmitRIComm_Act`): urai ulang di server; satu galat = tidak ada yang tersimpan.
// Per R/I COMM NAME: ringkasan bernama sama dipakai (OPERATORID / MODIFIEDDATE-nya diperbarui) atau dibuat baru; baris
// rincian disisipkan ke tabel flat (IDUSEDBY = ID ringkasan). Satu transaksi.
func (l *Layanan) SimpanUnggah(ctx context.Context, a Aktor, p PermintaanUnggah) (HasilSimpanUnggah, error) {
	if err := wajibPenuh(a); err != nil {
		return HasilSimpanUnggah{}, err
	}
	var hasil HasilSimpanUnggah
	err := l.tx(ctx, func(tx *dbTx) error {
		h, kel, err := l.periksaUnggah(ctx, tx, p.CSV)
		if err != nil {
			return err
		}
		if !h.Sah {
			return tolak("%s", KalimatGalat(h.Galat))
		}
		operator, waktu := potongByte(a.AkunID, models.BatasNama), FormatWaktuPega(l.jam())
		idRingkasan := &pemberiID{l: l, ringkasan: true}
		idKomisi := &pemberiID{l: l}
		hasil = HasilSimpanUnggah{Ringkasan: []RingkasanUnggah{}}
		for _, k := range kel {
			r := models.Ringkasan{UsedBy: k.nama, OperatorID: operator, ModifiedDate: waktu}
			ru := RingkasanUnggah{UsedBy: k.nama, Baru: k.cocok == nil, Jumlah: len(k.baris)}
			if k.cocok != nil {
				r.ID = k.cocok.ID
				if err := l.gudang.UbahRingkasan(ctx, tx, r); err != nil {
					return err
				}
			} else {
				if r.ID, err = idRingkasan.baru(ctx, tx); err != nil {
					return err
				}
				if err := l.gudang.SisipRingkasan(ctx, tx, r); err != nil {
					return err
				}
				hasil.RingkasanBaru++
			}
			ru.ID = r.ID
			for _, b := range k.baris {
				id, err := idKomisi.baru(ctx, tx)
				if err != nil {
					return err
				}
				if err := l.gudang.SisipKomisi(ctx, tx, models.Komisi{ID: id, IDUsedBy: r.ID, UsedBy: k.nama,
					Contract: b.Contract, Year: b.Year, Comm: b.Comm}); err != nil {
					return err
				}
				hasil.Disimpan++
			}
			hasil.Ringkasan = append(hasil.Ringkasan, ru)
		}
		return nil
	})
	if err != nil {
		return HasilSimpanUnggah{}, petaTidakAda(err)
	}
	return hasil, nil
}
