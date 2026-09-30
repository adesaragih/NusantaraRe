package services

// Kurs USD -> IDR - tiket 11 Treaty Contract Out.
//
// Untuk apa berkas ini: padanan aktivitas `testingKurs` (jalur yang benar-benar
// dipakai; nama barunya jujur, penyimpangan sadar 8): kurs yang berlaku pada
// `StartDate` tahun treaty (`InputTreatyArrangementDesc.StartDate`,
// `Harness/InboxTreatyContractDescription.xml` b6153), dan konversi Rp <-> Usd
// yang di Pega dikerjakan `HitungRpUsd_depan` / `CalculateTSIExcludeTreaty`.
//
// ⛔ Pengenal USD dibaca dari master mata uang lewat kodenya (AC 47); `QUARTER`
// diikuti apa adanya (AC 48).
//
// Dibaca sesudah: models/tco_kurs.go, repository/tco_kurs.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/treatycontractout/models"
	"nusantarare/modul/treatycontractout/repository"
)

var (
	// ErrGudangKursBelumDisuntik - pembaca kurs belum dipasang.
	ErrGudangKursBelumDisuntik = errors.New("services: exchange-rate reader is not injected")
	// ErrKursTidakAda - tidak ada kurs berlaku pada tanggal mulai tahun (422).
	ErrKursTidakAda = models.ErrKursTidakAda
	// ErrMataUangTidakDikenal - sentinel `inti/backend/db` dibuka untuk handler, yang
	// dilarang mengimpor lapisan repository (penjaga lintas aplikasi).
	ErrMataUangTidakDikenal = db.ErrMataUangTidakDikenal
	// ErrMasterKursRusak - master kurs/mata uang tidak dapat dipakai (503).
	ErrMasterKursRusak = errors.New("services: exchange-rate or currency master cannot be used")
	// ErrKonversiKursTidakSah - permintaan konversi di luar arah/skala yang dikenal (400).
	ErrKonversiKursTidakSah = errors.New("services: invalid exchange-rate conversion request")
)

// PembacaKursTCO memberi kurs berlaku untuk satu tahun treaty.
type PembacaKursTCO interface {
	Berlaku(ctx context.Context, tahun models.TahunTreaty) (models.KursTCO, error)
}

// PembacaMasterKursTCO membaca baris kurs yang berlaku pada satu tanggal -
// tanggalnya diurai dan dibandingkan Oracle (dibaca saja).
type PembacaMasterKursTCO interface {
	BacaBerlaku(ctx context.Context, idCurrency, quarter string, tanggal time.Time) (models.HasilMasterKursTCO, error)
}

// PembacaMataUangTCO menerjemahkan kode mata uang menjadi pengenalnya.
type PembacaMataUangTCO interface {
	Pengenal(ctx context.Context, kode string) (string, error)
}

type kursBelumDisuntik struct{}

func (kursBelumDisuntik) Berlaku(context.Context, models.TahunTreaty) (models.KursTCO, error) {
	return models.KursTCO{}, ErrGudangKursBelumDisuntik
}
func (kursBelumDisuntik) BacaBerlaku(context.Context, string, string, time.Time) (models.HasilMasterKursTCO, error) {
	return models.HasilMasterKursTCO{}, ErrGudangKursBelumDisuntik
}
func (kursBelumDisuntik) Pengenal(context.Context, string) (string, error) {
	return "", ErrGudangKursBelumDisuntik
}

// KursTampil adalah kurs berlaku seperti dikirim ke layar (`NitipKurs`).
type KursTampil struct {
	Kurs       string `json:"kurs"`
	Tanggal    string `json:"tanggal"`
	Mulai      string `json:"mulai"`
	Akhir      string `json:"akhir"`
	Currency   string `json:"currency"`
	IDCurrency string `json:"idCurrency"`
	Quarter    string `json:"quarter"`
	TreatyYear string `json:"treatyYear"`
	// BarisMasterDitolak - cacah baris master yang tanggalnya ditolak Oracle;
	// dilaporkan, tidak dipakai (lanjutan 6).
	BarisMasterDitolak int `json:"barisMasterDitolak"`
	// BarisMasterKembar - cacah baris berlaku lain yang TOIDR-nya sama (master
	// memuat baris kembar) [keputusan work owner 29-09-2026]; dilaporkan.
	BarisMasterKembar int `json:"barisMasterKembar"`
}

// KonversiTampil adalah hasil satu konversi Rp <-> Usd.
type KonversiTampil struct {
	Rp   string `json:"rp"`
	Usd  string `json:"usd"`
	Kurs string `json:"kurs"`
}

// KursTCO melayani kurs berlaku dan konversinya.
type KursTCO struct {
	svc      *Service
	tahun    PemeriksaTahunTCO
	master   PembacaMasterKursTCO
	mataUang PembacaMataUangTCO
}

// KursTCO menyusun layanannya; bawaannya gagal terang.
func (s *Service) KursTCO() *KursTCO {
	return &KursTCO{svc: s, tahun: gudangTahunTreatyBelumDisuntik{}, master: kursBelumDisuntik{},
		mataUang: kursBelumDisuntik{}}
}

func (l *KursTCO) salin() *KursTCO { s := *l; return &s }

// DenganTahun memasang pemeriksa tahun treaty.
func (l *KursTCO) DenganTahun(t PemeriksaTahunTCO) *KursTCO { s := l.salin(); s.tahun = t; return s }

// DenganMaster memasang pembaca master kurs.
func (l *KursTCO) DenganMaster(m PembacaMasterKursTCO) *KursTCO {
	s := l.salin()
	s.master = m
	return s
}

// DenganMataUang memasang penerjemah kode mata uang.
func (l *KursTCO) DenganMataUang(m PembacaMataUangTCO) *KursTCO {
	s := l.salin()
	s.mataUang = m
	return s
}

// MasterKursOracle menyusun pembaca master `TREATYEXCHANGEYEARLY`.
func MasterKursOracle(svc *Service) PembacaMasterKursTCO {
	return repository.NewMasterKursTCO(svc.DB())
}

// MataUangOracle menyusun penerjemah kode mata uang (`CURRENCY`, kode bersama).
func MataUangOracle(svc *Service) PembacaMataUangTCO { return db.NewMataUang(svc.DB()) }

// PembacaKursOracle menyusun pembaca kurs di atas Oracle - dipakai layanan klausul.
func PembacaKursOracle(svc *Service) PembacaKursTCO {
	return svc.KursTCO().DenganMaster(MasterKursOracle(svc)).DenganMataUang(MataUangOracle(svc))
}

// Berlaku - `testingKurs`: kurs `QUARTER = '0'` mata uang USD yang rentangnya
// memuat `StartDate` tahun. Tanpa baris -> "Tidak ada Nilai Kurs di Tahun : <tahun>".
func (l *KursTCO) Berlaku(ctx context.Context, tahun models.TahunTreaty) (models.KursTCO, error) {
	if tahun.StartDate.IsZero() {
		return models.KursTCO{}, models.GalatKursTidakAda{TreatyYear: tahun.TreatyYear}
	}
	id, err := l.mataUang.Pengenal(ctx, models.KodeMataUangAsalKursTCO)
	if err != nil {
		if errors.Is(err, db.ErrMataUangTidakDikenal) {
			return models.KursTCO{}, fmt.Errorf("%w: %v", ErrMasterKursRusak, err)
		}
		return models.KursTCO{}, err
	}
	h, err := l.master.BacaBerlaku(ctx, id, models.QuarterKursTahunanTCO, tahun.StartDate)
	if err != nil {
		return models.KursTCO{}, err
	}
	k, err := models.PilihKursBerlakuTCO(h, tahun.StartDate)
	switch {
	case errors.Is(err, models.ErrKursTidakAda):
		return models.KursTCO{}, models.GalatKursTidakAda{TreatyYear: tahun.TreatyYear, Tanggal: tahun.StartDate}
	case errors.Is(err, models.ErrKursGanda), errors.Is(err, models.ErrKursTakTerurai):
		return models.KursTCO{}, fmt.Errorf("%w: %v", ErrMasterKursRusak, err)
	case err != nil:
		return models.KursTCO{}, err
	}
	if k, err = models.LengkapiKursTCO(k); err != nil {
		return models.KursTCO{}, fmt.Errorf("%w: %v", ErrMasterKursRusak, err)
	}
	if k.Currency == "" {
		k.Currency = models.KodeMataUangAsalKursTCO
	}
	return k, nil
}

func tampilKurs(tahun models.TahunTreaty, k models.KursTCO) KursTampil {
	return KursTampil{Kurs: utils.FormatDecimal(k.ToIDR), Tanggal: utils.FormatTanggal(tahun.StartDate),
		Mulai: utils.FormatTanggal(k.Mulai), Akhir: utils.FormatTanggal(k.Akhir), Currency: k.Currency,
		IDCurrency: k.IDCurrency, Quarter: k.Quarter, TreatyYear: tahun.TreatyYear, BarisMasterDitolak: k.BarisDitolak,
		BarisMasterKembar: k.BarisKembar}
}

// KursTahun membaca kurs berlaku satu tahun treaty.
func (l *KursTCO) KursTahun(ctx context.Context, pelaku inti.Pelaku, tahunID string) (KursTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return KursTampil{}, err
	}
	tahun, err := l.tahun.Ambil(ctx, tahunID)
	if err != nil {
		return KursTampil{}, err
	}
	k, err := l.Berlaku(ctx, tahun)
	if err != nil {
		return KursTampil{}, err
	}
	return tampilKurs(tahun, k), nil
}

// Konversi menghitung padanan satu nilai dengan kurs berlaku tahun itu:
// `dari = "Rp"` -> `Usd = Rp / Kurs` pada skala 8 (`HitungRpUsd_depan`) atau 4
// (`CalculateTSIExcludeTreaty`); `dari = "Usd"` -> `Rp = Usd * Kurs`.
func (l *KursTCO) Konversi(ctx context.Context, pelaku inti.Pelaku, tahunID, dari, nilai, skala string) (KonversiTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return KonversiTampil{}, err
	}
	var sk int32
	switch strings.TrimSpace(skala) {
	case "", "8":
		sk = models.SkalaUsdDariRpTCO
	case "4":
		sk = models.SkalaUsdExclusionTCO
	default:
		return KonversiTampil{}, fmt.Errorf("%w: scale %q (4 or 8)", ErrKonversiKursTidakSah, skala)
	}
	if dari != models.MedanRp && dari != models.MedanUsd {
		return KonversiTampil{}, fmt.Errorf("%w: from %q (Rp or Usd)", ErrKonversiKursTidakSah, dari)
	}
	d, err := models.UraiDesimalMasukTCO(dari, nilai)
	if err != nil {
		return KonversiTampil{}, err
	}
	tahun, err := l.tahun.Ambil(ctx, tahunID)
	if err != nil {
		return KonversiTampil{}, err
	}
	k, err := l.Berlaku(ctx, tahun)
	if err != nil {
		return KonversiTampil{}, err
	}
	h := KonversiTampil{Kurs: utils.FormatDecimal(k.ToIDR)}
	if dari == models.MedanRp {
		usd, err := models.UsdDariRpTCO(d, k.ToIDR, sk)
		if err != nil {
			return KonversiTampil{}, err
		}
		h.Rp, h.Usd = utils.FormatDecimal(d), utils.FormatDecimal(usd)
		return h, nil
	}
	rp, err := models.RpDariUsdTCO(d, k.ToIDR)
	if err != nil {
		return KonversiTampil{}, err
	}
	h.Rp, h.Usd = utils.FormatDecimal(rp), utils.FormatDecimal(d)
	return h, nil
}
