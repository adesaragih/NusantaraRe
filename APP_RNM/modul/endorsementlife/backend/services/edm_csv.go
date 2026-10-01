package services

// Tiket 07 - unggah CSV peserta baru (`Upload CSV` b8973, `Add CSV Data`
// b10405 → `SaveCSVEDMLife`).
//
// ⛔ TANPA BATAS BARIS (spec penyimpangan 6, AC 36): berkas dibaca BERTAHAP -
// satu baris di memori untuk validasi, satu kelompok `ukuranKelompokCSV` untuk
// penyisipan. Gerbang `>50000` `SetPremi_EDM` b3373 tidak ditiru.
//
// ⛔ SELURUH BARIS DIPERIKSA DULU (R24). Pega menambahkan baris sampai baris
// pertama yang `PLAN`/`POLICY_HOLDER`-nya beda lalu berhenti; di sini satu
// pelanggaran menolak seluruh berkas - lintasan pertama hanya memeriksa,
// lintasan kedua (di dalam transaksi) memeriksa ulang dan menyisip.
//
// ⛔ Sesudah `Add CSV Data` unggahan terkunci (`.EditInput1 = 1` b2899 →
// `pyDisabledWhen` b8965/b10403): kasus ber-baris `New` menolak unggahan
// berikutnya (409). Pembuangan baris `New` langkah 2.1 tetap dijalankan.

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/jejak"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

const (
	// ukuranKelompokCSV - baris per kelompok penyisipan.
	ukuranKelompokCSV = 1000
	// batasPesanCSV - penolakan yang dikembalikan; sisanya hanya dicacah.
	batasPesanCSV = 200
)

var (
	// ErrCSVBukanPerubahanData - baris `New` hanya lahir pada `EdmType=1`
	// (4.3 b2791; wadah unggah b8698 `.EdmType=1`) - AC 13.
	ErrCSVBukanPerubahanData = errors.New("services: CSV upload is only available for EDM Type Perubahan Data")
	// ErrCSVTerkunci - `Add CSV Data` sudah dijalankan (`.EditInput1 = 1`).
	ErrCSVTerkunci = errors.New("services: CSV data has already been added to this endorsement; Upload CSV is locked")
	// ErrCSVKosong - berkas tanpa baris data.
	ErrCSVKosong = errors.New("services: the CSV file has no data rows below the header")
	// ErrCSVRusak - berkas tidak dapat diurai sebagai CSV.
	ErrCSVRusak = errors.New("services: the file is not a readable CSV")
	// ErrCSVTanpaAcuan - kasus tanpa peserta lama: `PremiumListDetail(1)` tidak ada.
	ErrCSVTanpaAcuan = errors.New("services: the endorsement has no existing participant to compare PLAN and POLICY_HOLDER with")
)

// HasilPeriksaCSV - tinjauan unggahan: cacah baris, yang ditolak, dan pesannya.
type HasilPeriksaCSV struct {
	Total     int               `json:"total"`
	Ditolak   int               `json:"ditolak"`
	Pesan     []models.PesanCSV `json:"pesan"`
	Terpotong bool              `json:"terpotong"`
	// Diabaikan - judul kolom yang bukan properti 4.1: tidak disimpan, ditampilkan.
	Diabaikan []string `json:"diabaikan"`
}

// GalatCSV - berkas ditolak karena barisnya; nol baris tersimpan.
type GalatCSV struct{ HasilPeriksaCSV }

func (g GalatCSV) Error() string {
	return fmt.Sprintf("services: %d of %d CSV rows are rejected; no row is saved", g.Ditolak, g.Total)
}

// HasilTambahCSV - hasil `Add CSV Data`.
type HasilTambahCSV struct {
	Disimpan int                 `json:"disimpan"`
	Dibuang  int                 `json:"dibuang"`
	Rekap    []map[string]string `json:"rekap"`
}

// periksaKasusCSV - gerbang kasus untuk unggahan.
func periksaKasusCSV(k models.Kasus) error {
	switch {
	case !k.Terbuka():
		return ErrKasusTertutup
	case k.EdmType != models.EdmTypePerubahanData:
		return ErrCSVBukanPerubahanData
	case k.Cacah[models.StatusNew] > 0:
		return ErrCSVTerkunci
	}
	return nil
}

// telusuriCSV membaca berkas baris demi baris; `sisip` menerima kelompok baris
// lolos SELAMA belum ada penolakan (nil = periksa saja).
func telusuriCSV(berkas io.Reader, acuan models.AcuanCSV, sisip func([]models.BarisCSV) error) (HasilPeriksaCSV, error) {
	c := csv.NewReader(berkas)
	c.FieldsPerRecord = -1
	c.TrimLeadingSpace = true
	c.ReuseRecord = true
	mentah, err := c.Read()
	if errors.Is(err, io.EOF) {
		return HasilPeriksaCSV{}, ErrCSVKosong
	}
	if err != nil {
		return HasilPeriksaCSV{}, fmt.Errorf("%w: %w", ErrCSVRusak, err)
	}
	judul, asing, err := models.JudulCSV(append([]string{}, mentah...))
	if err != nil {
		return HasilPeriksaCSV{}, fmt.Errorf("%w: %v", ErrMasukanTidakSah, err)
	}
	h := HasilPeriksaCSV{Pesan: []models.PesanCSV{}, Diabaikan: asing}
	if h.Diabaikan == nil {
		h.Diabaikan = []string{}
	}
	tolak := func(p ...models.PesanCSV) {
		for _, x := range p {
			if len(h.Pesan) < batasPesanCSV {
				h.Pesan = append(h.Pesan, x)
			} else {
				h.Terpotong = true
			}
		}
	}
	var kelompok []models.BarisCSV
	kirim := func() error {
		if sisip == nil || h.Ditolak > 0 || len(kelompok) == 0 {
			kelompok = kelompok[:0]
			return nil
		}
		err := sisip(kelompok)
		kelompok = nil
		return err
	}
	for {
		rec, err := c.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return h, fmt.Errorf("%w: %w", ErrCSVRusak, err)
		}
		if strings.TrimSpace(strings.Join(rec, "")) == "" {
			continue // baris kosong - bukan peserta
		}
		h.Total++
		baris := make(map[string]string, len(judul))
		for i, v := range rec {
			if i < len(judul) {
				baris[judul[i]] = v
			}
		}
		nilai, pesan := models.RapikanBarisCSV(h.Total, baris, acuan)
		if len(rec) > len(judul) {
			pesan = append(pesan, models.PesanCSV{Baris: h.Total, Pesan: fmt.Sprintf("row has %d fields, header has %d", len(rec), len(judul))})
		}
		if len(pesan) > 0 {
			h.Ditolak++
			tolak(pesan...)
			continue
		}
		if sisip != nil && h.Ditolak == 0 {
			kelompok = append(kelompok, models.BarisCSV{Nomor: h.Total, Nilai: nilai})
			if len(kelompok) >= ukuranKelompokCSV {
				if err := kirim(); err != nil {
					return h, err
				}
			}
		}
	}
	if h.Total == 0 {
		return h, ErrCSVKosong
	}
	return h, kirim()
}

// acuanCSV membaca `PremiumListDetail(1)`.
func (l *Layanan) acuanCSV(ctx context.Context, tx *db.Tx, id string) (models.AcuanCSV, error) {
	a, ada, err := l.gudang.AcuanCSV(ctx, tx, id)
	if err != nil {
		return models.AcuanCSV{}, err
	}
	if !ada {
		return models.AcuanCSV{}, ErrCSVTanpaAcuan
	}
	return a, nil
}

// PeriksaCSV - `Upload CSV`: mengurai dan memvalidasi, nol penulisan (tinjauan
// sebelum `Add CSV Data`).
func (l *Layanan) PeriksaCSV(ctx context.Context, p inti.Pelaku, id string, berkas io.Reader) (HasilPeriksaCSV, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilPeriksaCSV{}, err
	}
	k, err := l.muatKasus(ctx, id)
	if err != nil {
		return HasilPeriksaCSV{}, err
	}
	if err := periksaKasusCSV(k); err != nil {
		return HasilPeriksaCSV{}, err
	}
	acuan, err := l.acuanCSV(ctx, nil, id)
	if err != nil {
		return HasilPeriksaCSV{}, err
	}
	return telusuriCSV(berkas, acuan, nil)
}

// TambahCSV - `Add CSV Data` (`SaveCSVEDMLife`), satu transaksi: baris `New`
// lama dibuang (2.1), seluruh baris diperiksa ulang lalu disisip (4.1/4.3),
// rekap dihitung ulang bila kasus sudah disimpan (R31), jejak.
func (l *Layanan) TambahCSV(ctx context.Context, p inti.Pelaku, id string, berkas io.ReadSeeker) (HasilTambahCSV, error) {
	cek, err := l.PeriksaCSV(ctx, p, id, berkas)
	if err != nil {
		return HasilTambahCSV{}, err
	}
	if cek.Ditolak > 0 {
		return HasilTambahCSV{}, GalatCSV{cek}
	}
	if _, err := berkas.Seek(0, io.SeekStart); err != nil {
		return HasilTambahCSV{}, fmt.Errorf("services: rewinding the CSV file: %w", err)
	}
	var hasil HasilTambahCSV
	err = l.tx(ctx, func(tx *db.Tx) error {
		k, err := l.gudang.AmbilKasus(ctx, tx, id, true)
		if errors.Is(err, repository.ErrTidakAda) {
			return fmt.Errorf("%w: %q", ErrKasusTidakAda, id)
		}
		if err != nil {
			return err
		}
		if k.Cacah, err = l.gudang.CacahPeserta(ctx, tx, id); err != nil {
			return err
		}
		if err := periksaKasusCSV(k); err != nil {
			return err
		}
		acuan, err := l.acuanCSV(ctx, tx, id)
		if err != nil {
			return err
		}
		if hasil.Dibuang, err = l.gudang.HapusPesertaBaru(ctx, tx, id); err != nil {
			return err
		}
		h, err := telusuriCSV(berkas, acuan, func(b []models.BarisCSV) error {
			n, err := l.gudang.SisipPesertaCSV(ctx, tx, id, k.PLNumber, b)
			hasil.Disimpan += n
			return err
		})
		if err != nil {
			return err
		}
		if h.Ditolak > 0 {
			return GalatCSV{h} // berkas berubah di antara kedua lintasan
		}
		sudah, err := l.gudang.AdaRekap(ctx, tx, id)
		if err != nil {
			return err
		}
		if sudah {
			if _, err := l.gudang.HitungRekap(ctx, tx, id, k.Kepala["TYPE"]); err != nil {
				return err
			}
		}
		if hasil.Rekap, err = l.gudang.RekapKasus(ctx, tx, id); err != nil {
			return err
		}
		return l.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
			KlaimID: id, Dari: models.TahapInputEDMLife, Ke: models.TahapInputEDMLife, AkunID: p.AkunID, Waktu: l.jam(),
			Komentar: models.AksiTambahCSV,
		})
	})
	if err != nil {
		return HasilTambahCSV{}, err
	}
	if hasil.Rekap == nil {
		hasil.Rekap = []map[string]string{}
	}
	return hasil, nil
}
