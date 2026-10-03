package services_test

// Pembantu test before-image - TANPA Oracle.
//
// Seluruh angka dan kode di berkas test modul ini SINTETIS. Tidak ada nomor
// polis, nama, atau nilai produksi (CLAUDE.md §4.10).

import (
	"testing"
	"time"

	"nusantarare/inti/backend/uang"
	"nusantarare/modul/endorsmentfacin/backend/models"
)

const mataUangUji = "UJI-MU"

func duit(t *testing.T, s string) uang.Money {
	t.Helper()
	m, err := uang.NewMoney(s, mataUangUji)
	if err != nil {
		t.Fatalf("duit %q: %v", s, err)
	}
	return m
}

// kosong - nilai uang tak terisi, tetapi mata uang barisnya diketahui.
func kosong() uang.Money { return uang.Money{Currency: mataUangUji} }

func rasio(t *testing.T, s string) uang.Ratio {
	t.Helper()
	r, err := uang.NewRatio(s, 6)
	if err != nil {
		t.Fatalf("rasio %q: %v", s, err)
	}
	return r
}

func tanggal(t *testing.T, s string) time.Time {
	t.Helper()
	w, err := time.Parse("2006-01-02 15:04", s)
	if err != nil {
		t.Fatalf("tanggal %q: %v", s, err)
	}
	return w
}

// samaUang - sama nilai DAN mata uang; kosong hanya sama dengan kosong.
func samaUang(t *testing.T, label string, got, want uang.Money) {
	t.Helper()
	if got.Kosong() != want.Kosong() {
		t.Errorf("%s: kosong=%v, mau kosong=%v (%q)", label, got.Kosong(), want.Kosong(), got.String())
		return
	}
	if got.Kosong() {
		return
	}
	if got.Amount.Cmp(want.Amount) != 0 || got.Currency != want.Currency {
		t.Errorf("%s: %q, mau %q", label, got.String(), want.String())
	}
}

func samaRasio(t *testing.T, label string, got, want uang.Ratio) {
	t.Helper()
	if got.Kosong() != want.Kosong() {
		t.Errorf("%s: kosong=%v, mau kosong=%v (%q)", label, got.Kosong(), want.Kosong(), got.String())
		return
	}
	if got.Kosong() {
		return
	}
	if got.Value.Cmp(want.Value) != 0 || got.Scale != want.Scale {
		t.Errorf("%s: %s (skala %d), mau %s (skala %d)", label, got.String(), got.Scale, want.String(), want.Scale)
	}
}

// barisMU - satu baris mata uang terisi penuh.
func barisMU(t *testing.T, tsi, premi, rate string) models.BarisMataUang {
	return models.BarisMataUang{TSI: duit(t, tsi), Premium: duit(t, premi), Rate: rasio(t, rate)}
}

// spread - satu baris spreading.
func spread(t *testing.T, treaty, tsi, premi, share string) models.Spreading {
	return models.Spreading{
		TreatyType:      treaty,
		TSISpreaded:     duit(t, tsi),
		PremiumSpreaded: duit(t, premi),
		SharePercentage: rasio(t, share),
	}
}

// polisLamaFire - dokumen polis lama lini kebakaran: satu lokasi, satu item,
// satu coverage dengan dua baris spreading berjenis treaty BERBEDA.
func polisLamaFire(t *testing.T) models.OfferFacIn {
	t.Helper()
	return models.OfferFacIn{
		QuotationData: models.Quotation{
			BusinessCode:   "UJI-BC",
			BusinessType:   "UJI-BT",
			MarketingName:  "UJI-MKT",
			StatusBusiness: "1",
		},
		PolicyData: models.PolicyData{
			StartDateTime: tanggal(t, "2026-01-01 00:00"),
			EndDateTime:   tanggal(t, "2026-01-05 00:00"),
			Payment:       models.Pembayaran{Installment: "UJI-INST", RICommision: "UJI-RIC", PctBrokerageFee: "UJI-PBF"},
		},
		Currency:         mataUangUji,
		CurrencyList:     []models.BarisMataUang{barisMU(t, "1000", "10", "1")},
		IsB2B:            "UJI-B2B",
		PPnCheck:         "UJI-PPN",
		CedingCedantList: []models.Cedant{{CurrencyList: []models.BarisMataUang{barisMU(t, "400", "4", "1")}}},
		LocationList: []models.Lokasi{{
			Property: models.PropertiLokasi{
				PropertyItemList: []models.ItemProperti{{
					TSIObjectItem:           duit(t, "1000"),
					TotalGrossPremi:         duit(t, "12"),
					TotalPremiumNusantaraRe: duit(t, "6"),
					CoverageList: []models.Coverage{{
						EDM: "UJI-EDM",
						SpreadingList: []models.Spreading{
							spread(t, "UJI-T1", "600", "6", "60"),
							spread(t, "UJI-T2", "400", "4", "40"),
						},
					}},
				}},
				TotalTSIList:           []models.BarisMataUang{barisMU(t, "1000", "12", "1.2")},
				TotalTSIPremiGrossList: []models.BarisMataUang{barisMU(t, "1000", "12", "1.2")},
			},
		}},
	}
}

// kasusBaru - objek kerja endorsement sesudah langkah 1-13 (milik E03):
// tanggal endorsement sudah terisi, sisanya belum.
func kasusBaru(t *testing.T, tanggalEdm string) models.KasusEndorsement {
	t.Helper()
	return models.KasusEndorsement{
		OfferFacIn: models.OfferFacIn{
			QuotationData: models.Quotation{EdmDate: tanggal(t, tanggalEdm), StatusBusiness: "3", Type: "1"},
		},
	}
}
