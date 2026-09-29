package services_test

// Uji layanan kurs - TANPA Oracle (tiket 11).

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
	"nusantarare/pkg/utils"
)

var mulaiTahunKurs = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

type tahunKursUji struct{ mulai time.Time }

func (t tahunKursUji) Ambil(_ context.Context, id string) (models.TahunTreaty, error) {
	if id != "1000001" {
		return models.TahunTreaty{}, repository.ErrTahunTreatyTidakAda
	}
	return models.TahunTreaty{ID: id, TreatyYear: "2026", StartDate: t.mulai}, nil
}

type masterKursUji struct {
	baris       []models.KursTCO
	err         error
	id, quarter string
}

func (m *masterKursUji) Daftar(_ context.Context, id, quarter string) ([]models.KursTCO, error) {
	m.id, m.quarter = id, quarter
	return m.baris, m.err
}

type mataUangUji map[string]string

func (m mataUangUji) Pengenal(_ context.Context, kode string) (string, error) {
	if id, ada := m[kode]; ada {
		return id, nil
	}
	return "", fmt.Errorf("%w: %q", repository.ErrMataUangTidakDikenal, kode)
}

func barisKurs2026() []models.KursTCO {
	return []models.KursTCO{
		{ToIDR: apd.New(1550025, -2), Mulai: mulaiTahunKurs, Akhir: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
			IDCurrency: "UJI-USD", Currency: "USD", Quarter: "0"},
		{ToIDR: apd.New(15000, 0), Mulai: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Akhir: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC), IDCurrency: "UJI-USD", Currency: "USD", Quarter: "0"},
	}
}

func layananKurs(m *masterKursUji) *services.KursTCO {
	return services.New(nil).KursTCO().DenganTahun(tahunKursUji{mulai: mulaiTahunKurs}).DenganMaster(m).
		DenganMataUang(mataUangUji{"USD": "UJI-USD"})
}

// AC 46-48: kurs berlaku pada StartDate tahun; pengenal USD dari master mata
// uang; QUARTER '0' apa adanya.
func TestKursTahunBerlaku(t *testing.T) {
	m := &masterKursUji{baris: barisKurs2026()}
	k, err := layananKurs(m).KursTahun(context.Background(), pelakuUjiTCO, "1000001")
	if err != nil {
		t.Fatal(err)
	}
	if k.Kurs != "15500.25" || k.TreatyYear != "2026" || k.Currency != "USD" || k.Quarter != "0" ||
		k.Tanggal != utils.FormatTanggal(mulaiTahunKurs) {
		t.Errorf("kurs: %+v", k)
	}
	if m.id != "UJI-USD" || m.quarter != "0" {
		t.Errorf("saringan master: id %q quarter %q", m.id, m.quarter)
	}
}

// ADR-0015: kegagalan terlihat - periode tanpa kurs, master rusak, mata uang tak dikenal.
func TestKursGagalTerang(t *testing.T) {
	_, err := layananKurs(&masterKursUji{baris: barisKurs2026()[1:]}).KursTahun(context.Background(), pelakuUjiTCO, "1000001")
	if !errors.Is(err, services.ErrKursTidakAda) || err.Error() != "Tidak ada Nilai Kurs di Tahun : 2026" {
		t.Errorf("tanpa kurs: %v", err)
	}
	ganda := append(barisKurs2026(), models.KursTCO{ToIDR: apd.New(1, 0), Mulai: mulaiTahunKurs, Akhir: mulaiTahunKurs})
	if _, err := layananKurs(&masterKursUji{baris: ganda}).KursTahun(context.Background(), pelakuUjiTCO, "1000001"); !errors.Is(err, services.ErrMasterKursRusak) {
		t.Errorf("ganda: %v", err)
	}
	rusak := &masterKursUji{err: fmt.Errorf("%w: TOIDR %q", models.ErrKursTakTerurai, "x")}
	if _, err := layananKurs(rusak).KursTahun(context.Background(), pelakuUjiTCO, "1000001"); !errors.Is(err, services.ErrMasterKursRusak) {
		t.Errorf("tak terurai: %v", err)
	}
	tanpaUSD := services.New(nil).KursTCO().DenganTahun(tahunKursUji{mulai: mulaiTahunKurs}).
		DenganMaster(&masterKursUji{baris: barisKurs2026()}).DenganMataUang(mataUangUji{})
	if _, err := tanpaUSD.KursTahun(context.Background(), pelakuUjiTCO, "1000001"); !errors.Is(err, services.ErrMasterKursRusak) {
		t.Errorf("USD tidak di master mata uang: %v", err)
	}
	tanpaMulai := services.New(nil).KursTCO().DenganTahun(tahunKursUji{}).DenganMaster(&masterKursUji{baris: barisKurs2026()}).
		DenganMataUang(mataUangUji{"USD": "UJI-USD"})
	if _, err := tanpaMulai.KursTahun(context.Background(), pelakuUjiTCO, "1000001"); !errors.Is(err, services.ErrKursTidakAda) {
		t.Errorf("tahun tanpa StartDate: %v", err)
	}
	if _, err := services.New(nil).KursTCO().DenganTahun(tahunKursUji{mulai: mulaiTahunKurs}).KursTahun(context.Background(),
		pelakuUjiTCO, "1000001"); !errors.Is(err, services.ErrGudangKursBelumDisuntik) {
		t.Errorf("bawaan: %v", err)
	}
	if _, err := layananKurs(&masterKursUji{}).KursTahun(context.Background(), services.Pelaku{}, "1000001"); !errors.Is(err, services.ErrTanpaIdentitas) {
		t.Errorf("identitas: %v", err)
	}
}

// HitungRpUsd_depan / CalculateTSIExcludeTreaty di server - desimal persis.
func TestKonversiKurs(t *testing.T) {
	l := layananKurs(&masterKursUji{baris: barisKurs2026()})
	for _, k := range []struct{ dari, nilai, skala, rp, usd string }{
		{"Rp", "31000500", "", "31000500", "2000.00000000"},
		{"Rp", "1000", "4", "1000", "0.0645"},
		{"Usd", "2000", "", "31000500.00000000", "2000"},
		{"Rp", "1000,5", "8", "1000.5", "0.06454735"},
	} {
		h, err := l.Konversi(context.Background(), pelakuUjiTCO, "1000001", k.dari, k.nilai, k.skala)
		if err != nil || h.Rp != k.rp || h.Usd != k.usd || h.Kurs != "15500.25" {
			t.Errorf("%+v: %+v %v", k, h, err)
		}
	}
	for _, k := range []struct{ dari, skala string }{{"IDR", ""}, {"Rp", "2"}} {
		if _, err := l.Konversi(context.Background(), pelakuUjiTCO, "1000001", k.dari, "1", k.skala); !errors.Is(err, services.ErrKonversiKursTidakSah) {
			t.Errorf("%+v: %v", k, err)
		}
	}
}
