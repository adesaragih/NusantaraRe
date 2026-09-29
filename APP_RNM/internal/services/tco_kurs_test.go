package services_test

// Uji layanan kurs - TANPA Oracle (tiket 11).

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/utils"
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
	ditolak     []models.BarisKursDitolakTCO
	err         error
	id, quarter string
	tanggal     time.Time
}

// BacaBerlaku meniru saringan Oracle: baris yang rentangnya memuat tanggal.
func (m *masterKursUji) BacaBerlaku(_ context.Context, id, quarter string, tanggal time.Time) (models.HasilMasterKursTCO, error) {
	m.id, m.quarter, m.tanggal = id, quarter, tanggal
	h := models.HasilMasterKursTCO{Ditolak: m.ditolak}
	for _, k := range m.baris {
		if !tanggal.Before(k.Mulai) && !tanggal.After(k.Akhir) {
			h.Berlaku = append(h.Berlaku, k)
		}
	}
	return h, m.err
}

type mataUangUji map[string]string

func (m mataUangUji) Pengenal(_ context.Context, kode string) (string, error) {
	if id, ada := m[kode]; ada {
		return id, nil
	}
	return "", fmt.Errorf("%w: %q", db.ErrMataUangTidakDikenal, kode)
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
	if m.id != "UJI-USD" || m.quarter != "0" || !m.tanggal.Equal(mulaiTahunKurs) || k.BarisMasterDitolak != 0 {
		t.Errorf("saringan master: id %q quarter %q tanggal %v ditolak %d", m.id, m.quarter, m.tanggal, k.BarisMasterDitolak)
	}
}

// Lanjutan 6: baris yang tanggalnya ditolak Oracle TIDAK mematikan kurs yang
// berlaku - cacahnya dilaporkan ke layar; tanpa kurs berlaku, penolakan itu
// master rusak (503) yang menyebut teksnya, bukan "Tidak ada Nilai Kurs".
func TestKursBarisTanggalDitolakOracle(t *testing.T) {
	ditolak := []models.BarisKursDitolakTCO{{Kolom: "STARTDATE", Teks: "2019A801T000000.000 GMT"},
		{Kolom: "ENDDATE", Teks: "2019T000000.000 GMT"}}
	k, err := layananKurs(&masterKursUji{baris: barisKurs2026(), ditolak: ditolak}).KursTahun(context.Background(), pelakuUjiTCO, "1000001")
	if err != nil || k.Kurs != "15500.25" || k.BarisMasterDitolak != 2 {
		t.Errorf("berlaku + dua ditolak: %+v %v", k, err)
	}
	_, err = layananKurs(&masterKursUji{baris: barisKurs2026()[1:], ditolak: ditolak}).KursTahun(context.Background(), pelakuUjiTCO, "1000001")
	if !errors.Is(err, services.ErrMasterKursRusak) || errors.Is(err, services.ErrKursTidakAda) ||
		!strings.Contains(err.Error(), `STARTDATE "2019A801T000000.000 GMT"`) {
		t.Errorf("hanya ditolak: %v", err)
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
	// Galat Oracle diteruskan apa adanya - bukan disamarkan jadi "tidak ada kurs".
	gagalBaca := errors.New("ORA-00942: table or view does not exist")
	if _, err := layananKurs(&masterKursUji{err: gagalBaca}).KursTahun(context.Background(), pelakuUjiTCO, "1000001"); !errors.Is(err, gagalBaca) {
		t.Errorf("galat baca: %v", err)
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
	if _, err := layananKurs(&masterKursUji{}).KursTahun(context.Background(), inti.Pelaku{}, "1000001"); !errors.Is(err, inti.ErrTanpaIdentitas) {
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

// Temuan /code-review: TOIDR rusak di baris periode LAIN tidak menggagalkan
// pencarian; TOIDR rusak di baris yang berlaku adalah master rusak.
func TestKursBarisLainRusakTidakMenggagalkan(t *testing.T) {
	baris := []models.KursTCO{
		{TeksToIDR: "15500.25", Mulai: mulaiTahunKurs, Akhir: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC), Quarter: "0"},
		{TeksToIDR: "", Mulai: time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC), Akhir: time.Date(2012, 12, 31, 0, 0, 0, 0, time.UTC)},
	}
	k, err := layananKurs(&masterKursUji{baris: baris}).KursTahun(context.Background(), pelakuUjiTCO, "1000001")
	if err != nil || k.Kurs != "15500.25" {
		t.Errorf("baris lama rusak menggagalkan: %+v %v", k, err)
	}
	baris[0].TeksToIDR = "0"
	if _, err := layananKurs(&masterKursUji{baris: baris}).KursTahun(context.Background(), pelakuUjiTCO, "1000001"); !errors.Is(err, services.ErrMasterKursRusak) {
		t.Errorf("baris berlaku rusak: %v", err)
	}
}
