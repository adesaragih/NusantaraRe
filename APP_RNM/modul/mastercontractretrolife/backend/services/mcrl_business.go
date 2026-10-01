package services

// Business (paket 6, tiket 07/08) - padanan `SaveBusinessLife_Act` (`Save`
// b5739, aksi `pyBehaviors` lama) dan `SaveBusinessToAllLife_Act` (`Copy to all
// Reinstype` b12020) TANPA procedure `INSERTBUSINESS_LIFE`:
//
//	langkah 3 b589   wajib BIZCODE, BIZNAME, RIRATEID, RIRATE -> "All value cannot be empty." (b293)
//	BUSINESS NAME    autocomplete `BrowseBusinessLife_RD` mengisi BIZCODE; nama = master NOTE
//	R/I RATE         `RIRATE` = nama tabel rate, TEKS apa adanya (R7); sumbernya menunggu OQ-MCRL-13
//	salin-semua 1    CARI1 <- TREATYYEARID (tahun kontrak asal)
//	salin-semua 2    `GetTreatyContract_life`: kontrak `idtreatyyear =` tahun itu
//	salin-semua 3.1  `.REINSTYPEID == Param.REINSTYPEID` WhenTrue 3 = LEWATI -> sasaran = jenis BERBEDA (R2)
//	salin-semua 3.2  `SaveTreatyBusinessAll_Life_SQL`: INSERT baris baru per sasaran (TREATYBUSINESS_LIFE_SEQ)
//	salin-semua 4    "Copied to all reins types." (b1258)
//
// ⛔ Penyimpangan 5 tetap: pratinjau + konfirmasi. Konfirmasi membawa DAFTAR
// sasaran yang dilihat pengguna; yang tidak lagi sama dengan data = 409, nol
// baris lahir. Satu transaksi (R4): semua sasaran atau tidak sama sekali.
// ⛔ K6: USERID/TGLUPDATE tiap baris + satu baris log berisi cacah.
// ⛔ Tanpa penjaga dobel, seperti Pega (OQ-MCRL-06). `TREATYYEAR` ditulis (K4)
// walau SQL Pega tidak menulisnya.

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/mastercontractretrolife/backend/models"
	"nusantarare/modul/mastercontractretrolife/backend/repository"
)

// ErrDampakBerubah - data berubah sejak pratinjau/popup (409); nol baris disentuh.
var ErrDampakBerubah = errors.New("services: the data changed since the confirmation was shown; " +
	"review it again - nothing was changed")

// BusinessMasuk adalah badan permintaan simpan business (form `Business List`).
type BusinessMasuk struct {
	ID       string `json:"id"`
	BizCode  string `json:"bizCode"`  // BUSINESS CODE (diisi autocomplete)
	BizName  string `json:"bizName"`  // BUSINESS NAME - nama akhir dari master
	RIRateID string `json:"riRateId"` // diisi autocomplete R/I RATE
	RIRate   string `json:"riRate"`   // R/I RATE - teks apa adanya
}

func (m BusinessMasuk) keModel() (models.Business, error) {
	var w wajib
	b := models.Business{BizCode: w.teks("BUSINESS CODE", m.BizCode)}
	w.teks("BUSINESS NAME", m.BizName)
	// ⛔ R7: RIRATEID/RIRATE disimpan apa adanya - pemangkasan hanya untuk
	// memeriksa kosong, tidak untuk nilai yang disimpan.
	w.teks("RIRATEID", m.RIRateID)
	w.teks("R/I RATE", m.RIRate)
	b.RIRateID, b.RIRate = m.RIRateID, m.RIRate
	if err := w.galat(PesanKosongSemua); err != nil {
		return b, err
	}
	if err := muat("RIRATEID", b.RIRateID, lebarKode); err != nil {
		return b, err
	}
	return b, muat("R/I RATE", b.RIRate, lebarTeks)
}

// namaBusiness membaca nama business dari master; pilihan BARU wajib lolos saringan autocomplete
// `BrowseBusinessLife_RD` b651 (`.OLDID StartsWith "L"`).
func (l *Layanan) namaBusiness(ctx context.Context, kode string, pilihanBaru bool) (string, error) {
	m, ada, err := l.gudang.AmbilMasterBusiness(ctx, kode)
	if err != nil {
		return "", err
	}
	if !ada {
		return "", fmt.Errorf("%w: BUSINESS CODE %q is not in the business master", ErrMasukanTidakSah, kode)
	}
	if pilihanBaru && !m.Life() {
		return "", fmt.Errorf("%w: BUSINESS CODE %q is not a life business", ErrMasukanTidakSah, kode)
	}
	return m.Note, nil
}

// salinanInduk menulis salinan kontrak dan tahun ke baris business (K4).
func salinanInduk(b *models.Business, k models.Kontrak, th models.TahunTreaty) {
	b.TreatyYearID, b.TreatyYear, b.TreatyContractID = k.IDTreatyYear, th.TreatyYear, k.ID
	b.ReinsTypeID, b.ReinsTypeName = k.ReinsTypeID, k.ReinsTypeName
}

// SimpanBusiness menyimpan business baru (`m.ID` kosong) di bawah kontrakID,
// atau mengubah business m.ID (kontrakID boleh kosong; bila terisi harus cocok).
func (l *Layanan) SimpanBusiness(ctx context.Context, p inti.Pelaku, kontrakID string, m BusinessMasuk) (models.Business, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.Business{}, err
	}
	var hasil models.Business
	if err := pelakuMuat(p, lebarKode); err != nil {
		return models.Business{}, err
	}
	err := l.tx(ctx, func(tx *db.Tx) error {
		ubah := m.ID != ""
		var kodeLama string
		if ubah {
			lama, err := l.gudang.AmbilBusiness(ctx, tx, m.ID)
			if err != nil {
				return tidakAda(err, ErrBusinessTidakAda, m.ID)
			}
			if kontrakID != "" && lama.TreatyContractID != kontrakID {
				return fmt.Errorf("%w: %s in treaty contract %s", ErrBusinessTidakAda, m.ID, kontrakID)
			}
			kontrakID, kodeLama = lama.TreatyContractID, lama.BizCode
		}
		// ⛔ Induk dikunci lebih dulu - kaskade hapus kontrak yang bersamaan menunggu (K2 tanpa FK).
		if err := l.kunci(ctx, tx, HapusKontrak, kontrakID); err != nil {
			return err
		}
		k, err := l.ambilKontrak(ctx, tx, kontrakID)
		if err != nil {
			return err
		}
		th, err := l.ambilTahun(ctx, tx, k.IDTreatyYear)
		if err != nil {
			return err
		}
		b, err := m.keModel()
		if err != nil {
			return err
		}
		if b.BizName, err = l.namaBusiness(ctx, b.BizCode, !ubah || b.BizCode != kodeLama); err != nil {
			return err
		}
		salinanInduk(&b, k, th)
		b.UserID = p.AkunID
		if !ubah {
			if b.ID, err = l.gudang.SisipBusiness(ctx, tx, b); err != nil {
				return err
			}
		} else {
			b.ID = m.ID
			if err := l.gudang.PerbaruiBusiness(ctx, tx, b); err != nil {
				return tidakAda(err, ErrBusinessTidakAda, b.ID)
			}
		}
		hasil, err = l.gudang.AmbilBusiness(ctx, tx, b.ID)
		return tidakAda(err, ErrBusinessTidakAda, b.ID)
	})
	return hasil, err
}

// PratinjauSalin - apa yang akan ditulis `Copy to all Reinstype`.
type PratinjauSalin struct {
	Business    models.Business  `json:"business"`
	ReinsTypeID string           `json:"reinsTypeId"` // dasar pemilihan: jenis kontrak asal
	Sasaran     []models.Kontrak `json:"sasaran"`
}

// HasilSalin - jawaban eksekusi `Copy to all Reinstype`.
type HasilSalin struct {
	Pesan  string            `json:"pesan"`
	Jumlah int               `json:"jumlah"`
	Baru   []models.Business `json:"baru"`
}

func (l *Layanan) sasaranSalin(ctx context.Context, tx *db.Tx, businessID string) (PratinjauSalin, error) {
	b, err := l.gudang.AmbilBusiness(ctx, tx, businessID)
	if err != nil {
		return PratinjauSalin{}, tidakAda(err, ErrBusinessTidakAda, businessID)
	}
	semua, err := l.gudang.DaftarKontrak(ctx, tx, b.TreatyYearID)
	if err != nil {
		return PratinjauSalin{}, err
	}
	p := PratinjauSalin{Business: b, ReinsTypeID: b.ReinsTypeID, Sasaran: []models.Kontrak{}}
	for _, k := range semua {
		if k.ReinsTypeID != b.ReinsTypeID {
			p.Sasaran = append(p.Sasaran, k)
		}
	}
	return p, nil
}

// PratinjauSalinSemua - langkah pratinjau (penyimpangan 5): nol tulisan.
func (l *Layanan) PratinjauSalinSemua(ctx context.Context, p inti.Pelaku, businessID string) (PratinjauSalin, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return PratinjauSalin{}, err
	}
	return l.sasaranSalin(ctx, nil, businessID)
}

func samaHimpunan(a, b []string) bool {
	x, y := append([]string{}, a...), append([]string{}, b...)
	sort.Strings(x)
	sort.Strings(y)
	return strings.Join(x, "\x00") == strings.Join(y, "\x00")
}

// SalinSemua menjalankan `Copy to all Reinstype` untuk sasaran yang sudah
// dikonfirmasi. Sasaran dihitung ulang DI DALAM transaksi; berbeda = 409.
func (l *Layanan) SalinSemua(ctx context.Context, p inti.Pelaku, businessID string, dikonfirmasi []string) (HasilSalin, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilSalin{}, err
	}
	if err := pelakuMuat(p, lebarKode); err != nil {
		return HasilSalin{}, err
	}
	hasil := HasilSalin{Pesan: PesanSalinSemua, Baru: []models.Business{}}
	err := l.tx(ctx, func(tx *db.Tx) error {
		pr, err := l.sasaranSalin(ctx, tx, businessID)
		if err != nil {
			return err
		}
		ada := make([]string, len(pr.Sasaran))
		for i, k := range pr.Sasaran {
			ada[i] = k.ID
		}
		if !samaHimpunan(ada, dikonfirmasi) {
			return fmt.Errorf("%w (confirmed %d target contracts, now %d)", ErrDampakBerubah, len(dikonfirmasi), len(ada))
		}
		th, err := l.ambilTahun(ctx, tx, pr.Business.TreatyYearID)
		if err != nil {
			return err
		}
		for _, k := range pr.Sasaran {
			// ⛔ Sasaran dikunci - kontrak yang terhapus bersamaan = keadaan berubah, nol baris ditulis.
			if err := l.gudang.KunciBaris(ctx, tx, string(HapusKontrak), k.ID); err != nil {
				if errors.Is(err, repository.ErrTidakAda) {
					return fmt.Errorf("%w (target treaty contract %s was deleted)", ErrDampakBerubah, k.ID)
				}
				return err
			}
			b := models.Business{BizCode: pr.Business.BizCode, BizName: pr.Business.BizName,
				RIRateID: pr.Business.RIRateID, RIRate: pr.Business.RIRate, UserID: p.AkunID}
			salinanInduk(&b, k, th)
			id, err := l.gudang.SisipBusiness(ctx, tx, b)
			if err != nil {
				return err
			}
			b.ID = id
			hasil.Baru = append(hasil.Baru, b)
		}
		hasil.Jumlah = len(hasil.Baru)
		return nil
	})
	if err != nil {
		return HasilSalin{}, err
	}
	l.catat(fmt.Sprintf("master contract retro life: business %s copied to %d treaty contracts%s", businessID, hasil.Jumlah, oleh(p)))
	return hasil, nil
}
