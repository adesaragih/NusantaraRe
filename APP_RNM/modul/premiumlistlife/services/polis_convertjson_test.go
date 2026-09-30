package services_test

// OQ-PL-14 DITUTUP (GILIRAN-17): `convertJsonNusareToProduction` ditiru lewat
// outbox, pelaksana STUB (OQ-PL-11 - tanpa panggilan nyata).

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/premiumlistlife/services"
)

func TestParameterConvertJsonVerbatim(t *testing.T) {
	saat := time.Date(2026, 9, 29, 3, 4, 5, 0, time.UTC)
	p, err := services.RakitParameterConvertJson("UJI-POLIS-1", "UJI-NOPOL-1", saat)
	if err != nil {
		t.Fatal(err)
	}
	j, _ := json.Marshal(p)
	for _, mau := range []string{`"noPolis":"UJI-NOPOL-1"`, `"caseId":"UJI-POLIS-1"`, `"tglInput":"20260929T030405.000 GMT"`} {
		if !strings.Contains(string(j), mau) {
			t.Errorf("parameter tanpa %s: %s", mau, j)
		}
	}
	if _, err := services.RakitParameterConvertJson("UJI-POLIS-1", "", saat); err == nil {
		t.Error("nomor polis kosong diterima")
	}
	isi, err := os.ReadFile(`D:\XML\RNM_BRD\PremiumList Life\ConnectREST\ConvertJsonNusareToProduction.xml`)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	baris := strings.Split(string(isi), "\n")
	for n, nama := range map[int]string{200: "noPolis", 209: "caseId", 221: "tglInput"} {
		if !strings.Contains(baris[n-1], "<pyParameterName>"+nama+"</pyParameterName>") {
			t.Errorf("b%d bukan parameter %s: %q", n, nama, baris[n-1])
		}
	}
}

type pembacaNoPolisUji struct {
	diminta []string
}

func (p *pembacaNoPolisUji) NomorPolisDariID(_ context.Context, id string) (string, error) {
	p.diminta = append(p.diminta, id)
	return "UJI-NOPOL-1", nil
}

type resolverUjiConvert struct{}

func (resolverUjiConvert) Resolve(context.Context, layanan.KunciLayanan) (string, error) {
	return "alamat-uji", nil
}

func barisPolis(jenis string) outbox.BarisEfekKeluar {
	return outbox.BarisEfekKeluar{ID: "1", Modul: services.ModulPremiumListLife, Jenis: jenis,
		Muatan: `{"klaim_id":"UJI-POLIS-1","adjustment_id":"","akun_id":"UJI-AKUN","waktu":"2026-09-29T10:00:00Z"}`}
}

func TestPelaksanaPremiumListStub(t *testing.T) {
	ctx := context.Background()
	baca := &pembacaNoPolisUji{}
	p := services.PelaksanaPremiumList{Lingkungan: inti.Produksi, Resolver: resolverUjiConvert{}, Pembaca: baca}
	if err := p.Laksanakan(ctx, nil, barisPolis(services.NamaEfekArasapasPolis)); !errors.Is(err, outbox.ErrArasapasBelumDisetujui) {
		t.Errorf("produksi: %v, mau ErrArasapasBelumDisetujui (stub, OQ-PL-11)", err)
	}
	if len(baca.diminta) != 1 || baca.diminta[0] != "UJI-POLIS-1" {
		t.Errorf("parameter noPolis tidak dirakit dari polis: %v", baca.diminta)
	}
	if err := p.Laksanakan(ctx, nil, barisPolis(services.NamaEfekAlarmPolis)); !errors.Is(err, outbox.ErrEmailBelumDisetujui) {
		t.Errorf("alarm: %v, mau ErrEmailBelumDisetujui", err)
	}
	nonProd := services.PelaksanaPremiumList{Lingkungan: inti.BukanProduksi, Resolver: resolverUjiConvert{}, Pembaca: baca}
	if err := nonProd.Laksanakan(ctx, nil, barisPolis(services.NamaEfekArasapasPolis)); !errors.Is(err, outbox.ErrPengirimStubNonProduksi) {
		t.Errorf("non-produksi: %v, mau ErrPengirimStubNonProduksi", err)
	}
	salah := barisPolis(services.NamaEfekArasapasPolis)
	// Modul LAIN - nilai kolom MODUL milik Claim Life. Refactor bentuk B:
	// literal, sebab uji modul ini tidak boleh mengimpor Claim Life.
	salah.Modul = "CLAIMLIFE"
	if err := p.Laksanakan(ctx, nil, salah); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Errorf("modul lain: %v, mau ErrPermintaanTidakSah", err)
	}
	asing := barisPolis("asing")
	if err := p.Laksanakan(ctx, nil, asing); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Errorf("jenis asing: %v, mau ErrPermintaanTidakSah", err)
	}
}
