package tiruan

// Tiruan keputusan kasus (repository/edm_putusan.go). `Panggilan` mencatat
// urutan metode tulis - uji urutan `services.Putuskan` (AC 42).

import (
	"context"
	"fmt"
	"sort"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

func (g *Gudang) catat(nama string) { g.Panggilan = append(g.Panggilan, nama) }

// Riwayat meniru `sqlRiwayat` (`NO DESC`).
func (g *Gudang) Riwayat(_ context.Context, _ *db.Tx, kasusID string) ([]models.BarisRiwayat, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return nil, g.Galat
	}
	r := append([]models.BarisRiwayat{}, g.RiwayatKasus[kasusID]...)
	sort.SliceStable(r, func(i, j int) bool { return r[i].No > r[j].No })
	return r, nil
}

// SisipRiwayat meniru `sqlSisipRiwayat`.
func (g *Gudang) SisipRiwayat(_ context.Context, _ *db.Tx, r models.RiwayatTulis) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return g.Galat
	}
	g.catat("SisipRiwayat")
	if g.RiwayatKasus == nil {
		g.RiwayatKasus = map[string][]models.BarisRiwayat{}
	}
	no := 0
	for _, b := range g.RiwayatKasus[r.KasusID] {
		no = max(no, b.No)
	}
	g.RiwayatKasus[r.KasusID] = append(g.RiwayatKasus[r.KasusID], models.BarisRiwayat{
		No: no + 1, Tanggal: r.Waktu, PIC: r.PIC, Status: r.Status, Komentar: r.Komentar})
	return nil
}

// AdaVersiResmi meniru `sqlAdaVersiResmi`.
func (g *Gudang) AdaVersiResmi(_ context.Context, _ *db.Tx, nomorPolis string, prodKe int) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return false, g.Galat
	}
	g.catat("AdaVersiResmi")
	for _, p := range g.Polis {
		if p.NoPolis == nomorPolis && p.ProdKe == prodKe {
			return true, nil
		}
	}
	return false, nil
}

// Resmikan meniru `sqlResmikanKepala` + `sqlResmikanPeserta`.
func (g *Gudang) Resmikan(_ context.Context, _ *db.Tx, r models.ResmiKasus) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return 0, g.Galat
	}
	g.catat("Resmikan")
	p, ada := g.Polis[r.ID]
	if !ada || p.Status != "" {
		return 0, fmt.Errorf("tiruan: peresmian kasus %s menyentuh 0 baris", r.ID)
	}
	p.NoPolis, p.ProdKe, p.NoEndors, p.PLNumberEDM, p.Status = r.NomorPolis, r.ProdKe, r.Nomor, r.Nomor, models.StatusKasusSelesai
	c := 0
	for _, d := range g.Peserta {
		if d.PolisID == r.ID {
			d.Nilai["PL_NUMBER_EDM"] = r.Nomor
			d.Nilai["STATUS_OLD"] = models.StatusLama(d.EdmStatus)
			d.Nilai["STATUS"] = r.StatusJenis
			c++
		}
	}
	return c, nil
}

// Tolak meniru `sqlTolak`.
func (g *Gudang) Tolak(_ context.Context, _ *db.Tx, kasusID string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return g.Galat
	}
	g.catat("Tolak")
	p, ada := g.Polis[kasusID]
	if !ada || p.Status != "" {
		return fmt.Errorf("tiruan: penolakan kasus %s menyentuh 0 baris", kasusID)
	}
	p.Status = models.StatusKasusDitolak
	return nil
}

// TulisRekapWarisan meniru `sqlHapusRekapWarisan` + `sqlSisipRekapWarisan`.
func (g *Gudang) TulisRekapWarisan(_ context.Context, _ *db.Tx, r repository.RekapWarisanTulis) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return 0, g.Galat
	}
	g.catat("TulisRekapWarisan")
	sisa := g.RekapWarisan[:0]
	for _, b := range g.RekapWarisan {
		if b["PL_NUMBER"] != r.NomorPolis || b["IDPEGA"] != r.KasusID {
			sisa = append(sisa, b)
		}
	}
	g.RekapWarisan = sisa
	for _, rk := range g.RekapData[r.KasusID] {
		b := map[string]string{"COB": r.COB, "PL_NUMBER": r.NomorPolis, "PL_NUMBER_EDM": r.Nomor, "IDPEGA": r.KasusID}
		for _, k := range repository.KolomRekapWarisan {
			if _, sudah := b[k]; sudah {
				continue
			}
			v := rk[k]
			if v == "" && k != "CURRENCY" {
				v = "0"
			}
			b[k] = v
		}
		g.RekapWarisan = append(g.RekapWarisan, b)
	}
	return len(g.RekapData[r.KasusID]), nil
}

// TulisProduksiWarisan meniru `sqlSisipProduksiWarisan`: satu baris dari kepala kasus.
func (g *Gudang) TulisProduksiWarisan(_ context.Context, _ *db.Tx, r repository.ProduksiWarisanTulis) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Galat != nil {
		return 0, g.Galat
	}
	g.catat("TulisProduksiWarisan")
	p, ada := g.Polis[r.KasusID]
	if !ada {
		return 0, fmt.Errorf("tiruan: produksi warisan kasus %s menyentuh 0 baris", r.KasusID)
	}
	kp := func(k string) string { return p.Kepala[k] }
	g.ProduksiWarisan = append(g.ProduksiWarisan, map[string]string{
		"IDPEGA": r.KasusID, "NOPOLIS": r.NomorPolis, "NOENDORS": r.Nomor, "BUSINESSCODE": kp("BUSINESS_CODE"),
		"BUSINESSNAME": kp("BUSINESS_NAME"), "CEDINGCO": kp("CEDING_CO"), "CEDINGCONAME": kp("CEDING_CO_NAME"),
		"DATERECEIVED": kp("DATE_RECEIVED"), "MARKETINGCODE": kp("MARKETING_CODE"), "MARKETINGNAME": kp("MARKETING_NAME"),
		"POLICYHOLDER": kp("POLICY_HOLDER"), "POLICYHOLDERNAME": kp("POLICY_HOLDER_NAME"), "PRORATETYPE": kp("PRO_RATE_TYPE"),
		"CREATEOPNAME": p.Pembuat, "SOB": kp("SOB"), "SOBNAME": kp("SOB_NAME"), "TYPE": kp("TYPE"),
		"TYPECEDING": kp("TYPE_CEDING"), "MOID": kp("MO_ID"), "NOOFFER": kp("NO_OFFER"), "RISLIPRNM": kp("RI_SLIP_RNM"),
		"RETROID": kp("RETRO_ID"), "RETRONAME": kp("RETRO_NAME"), "TYPECEDINGNAME": models.NamaJenisCeding[kp("TYPE_CEDING")],
		"SECURITYREINSURERID": kp("SECURITY_REINSURER_ID"), "SECURITYREINSURER": kp("SECURITY_REINSURER"),
		"TGL_INPUT": "SYSDATE", // cap waktu basis data - tidak ditiru
	})
	return 1, nil
}
