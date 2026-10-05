package services

// Save - versi Go `SaveAggregate_Act`: periksa setiap baris (pesan Pega apa adanya, seluruhnya dikumpulkan), lalu
// simpan seluruh baris dalam SATU transaksi (Pega menyimpan per baris dengan WriteNow; di sini gagal satu baris =
// tidak ada yang tersimpan). Baris "Total :" dilewati.

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/aggregate/backend/models"
)

// PesanTersimpan - `pyPortal.Status` sesudah Save berhasil.
const PesanTersimpan = "Successfully Saved"

// batasCobaNomor - nomor SEQ_AGGREGATE yang sudah terpakai dilewati paling banyak sekian kali.
const batasCobaNomor = 1000

// PermintaanSimpan - grid pratinjau apa adanya (termasuk baris "Total :").
type PermintaanSimpan struct {
	Baris []models.Baris `json:"baris"`
}

// HasilSimpan - jumlah baris tersimpan dan pesan layar.
type HasilSimpan struct {
	Disimpan int    `json:"disimpan"`
	Pesan    string `json:"pesan"`
}

var kolomDikenal = func() map[string]models.Kolom {
	m := map[string]models.Kolom{}
	for _, k := range models.KolomGrid {
		m[k.Nama] = k
	}
	return m
}()

// PeriksaBaris - `SaveAggregate_Act` langkah 3: pesan Pega untuk setiap baris data (nomor = urutan di grid).
func PeriksaBaris(baris []models.Baris) []string {
	var pesan []string
	for i, b := range baris {
		if b["ASSESMENT_ZONE"] == models.ZonaTotal {
			continue
		}
		n := i + 1
		if b["ASSESMENT_ZONE"] == "" {
			pesan = append(pesan, fmt.Sprintf("Assesment Zone in list %d Not Found", n))
		}
		if b["TO_USD"] == "" {
			pesan = append(pesan, fmt.Sprintf("To USD in list %d cannot be empty", n))
		}
		if !slices.Contains(models.TreatyTypeSah, b["TREATY_TYPE"]) {
			pesan = append(pesan, fmt.Sprintf("Treaty Type in List %d Can Only be Filled With OR, QS or SURPLUS", n))
		}
		if b["RNM_SHARE"] == "" {
			pesan = append(pesan, fmt.Sprintf("RNM Share in List %d Can't Null, Please Check Your Data", n))
		}
		if b["CEDING_NAME"] == "" {
			pesan = append(pesan, fmt.Sprintf("Ceding Name in List %d Not Found, Please Check CedingID", n))
		}
	}
	return pesan
}

// isiSah - nilai baris dapat ditulis ke AGGREGATE: angka terurai, tanggal sah, teks muat lebar kolom.
func isiSah(b models.Baris) bool {
	for _, k := range models.KolomGrid {
		v := strings.TrimSpace(b[k.Nama])
		if v == "" {
			continue
		}
		switch k.Jenis {
		case models.Angka:
			if desimal(v) == nil {
				return false
			}
		case models.Tanggal:
			if _, err := time.Parse("02-01-2006", v); err != nil {
				return false
			}
		default:
			if len(v) > k.Lebar {
				return false
			}
		}
	}
	return true
}

// Simpan - Save.
func (l *Layanan) Simpan(ctx context.Context, p inti.Pelaku, req PermintaanSimpan) (HasilSimpan, error) {
	if err := periksaPelaku(p); err != nil {
		return HasilSimpan{}, err
	}
	for i, b := range req.Baris {
		for nama := range b {
			if _, ada := kolomDikenal[nama]; !ada {
				return HasilSimpan{}, tolak("row %d has an unknown column %s", i+1, nama)
			}
		}
	}
	if pesan := PeriksaBaris(req.Baris); len(pesan) > 0 {
		return HasilSimpan{}, GalatSimpan{Pesan: pesan}
	}
	var data []int
	for i, b := range req.Baris {
		if b["ASSESMENT_ZONE"] == models.ZonaTotal {
			continue
		}
		if !isiSah(b) {
			return HasilSimpan{}, GalatSimpan{Pesan: []string{fmt.Sprintf("Failed to Save, Error in row %d", i+1)}}
		}
		data = append(data, i)
	}
	if len(data) == 0 {
		return HasilSimpan{}, tolak("there is no data row to save; upload a CSV file first")
	}
	err := l.tx(ctx, func(tx *dbTx) error {
		sekarang, err := l.gudang.Sekarang(ctx, tx)
		if err != nil {
			return err
		}
		for _, i := range data {
			id, err := l.idBerikut(ctx, tx)
			if err != nil {
				return err
			}
			if err := l.gudang.Sisip(ctx, tx, id, sekarang, p.AkunID, req.Baris[i]); err != nil {
				return fmt.Errorf("baris %d: %w", i+1, err)
			}
		}
		return nil
	})
	if err != nil {
		return HasilSimpan{}, err
	}
	return HasilSimpan{Disimpan: len(data), Pesan: PesanTersimpan}, nil
}

// idBerikut - `AGG-<SEQ_AGGREGATE>`; nomor yang sudah terpakai (mis. dibuat Pega) dilewati.
func (l *Layanan) idBerikut(ctx context.Context, tx *dbTx) (string, error) {
	for range batasCobaNomor {
		n, err := l.gudang.NomorBerikut(ctx, tx)
		if err != nil {
			return "", err
		}
		id := models.AwalanID + strconv.FormatInt(n, 10)
		ada, err := l.gudang.AdaID(ctx, tx, id)
		if err != nil {
			return "", err
		}
		if !ada {
			return id, nil
		}
	}
	return "", fmt.Errorf("services: %d nomor SEQ_AGGREGATE berturut-turut sudah terpakai", batasCobaNomor)
}
