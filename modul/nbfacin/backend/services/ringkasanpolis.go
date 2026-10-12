package services

// Total kasus dan DATE_OFFER T_GENERAL_POLIS (migrasi 252) - tiket 76, keputusan work owner 12-10-2026:
//   - kolom hitungan "SIMPAN dari hitungan Go yang sudah ada": TOTAL_TSI_NUSA_RE / TOTAL_TSI_TOP_RISK (totalTSIUW =
//     CountTotalTSIPremiNusaRe_Act langkah 14) dan TOTAL_TSI_NUSA_RE_SPREADING / TOTAL_PREMI_NUSA_RE (ringkasanRNM + totalKeIDR =
//     Total RNM IDR tab Cedant, MENYIMPANG K49-11 per mata uang -> IDR). Titik tulis `[keputusan agent]`: SETIAP Submit -
//     setter Pega (CountPremiAndTSINusantaraRe_ACT / CountTotalTSIPremiNusaRe_Act) berjalan tiap form detail dibuka / diubah dan
//     di GetLimitAkseptasi*, jadi nilai halaman saat Submit = nilai terakhir; Go menghitung saat dipakai, bukan per ubahan isian.
//     SUM_TOTAL_TSI / TREATY_CAPACITY TIDAK ditulis: Go belum punya hitungannya (SumTSIObjItem_act / SumTreatyCapacity_Act).
//   - hitung gagal / kurs tidak ada -> log dan total itu NULL; Submit tidak pernah gagal karena total ini.
//   - DATE_OFFER: Submit tahap marketing, hanya bila masih NULL; ekspresi ViewOfferFacInUW_PreDT baris 2
//     `@DateTime.CurrentDate("MM/dd/yyyy hh:mm a","GMT")` (pendekatan work owner; Pega mengisinya tiap form dibuka).

import (
	"context"
	"log"
	"sync"

	"nusantarare/modul/nbfacin/backend/models"
)

// bentukDateOffer - "MM/dd/yyyy hh:mm a" Java (jam 12-an, AM / PM) dalam GMT.
const bentukDateOffer = "01/02/2006 03:04 PM"

// dateOfferSubmit - DATE_OFFER Submit tahap marketing; "" tahap lain.
func (s *Service) dateOfferSubmit(tahap string) string {
	if tahap != TahapMarketing {
		return ""
	}
	return s.sekarang().UTC().Format(bentukDateOffer)
}

// ringkasanTotalPolis - lihat kepala berkas; nil bila tanpa pembaca spreading. Galat baca -> log + nil (bukan galat Submit).
func (s *Service) ringkasanTotalPolis(ctx context.Context, id string) *models.RingkasanTotalPolis {
	if s.spreading == nil {
		return nil
	}
	g := &models.RingkasanTotalPolis{}
	// Dua cabang independen PARALEL (laporan work owner 12-10-2026: pohon spreading tidak lagi dibaca berurutan dua kali); tiap cabang
	// menulis medannya sendiri dan menoleransi galatnya sendiri seperti alur berurutan lama - tidak saling membatalkan.
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		d := dasarJenjang{}
		var err error
		if d.sp, err = s.bacaSpreading(ctx, nil, id); err == nil {
			if d.aneka, err = s.bacaJenjangAneka(ctx, id); err == nil {
				if d.cargo, err = s.bacaJenjangCargo(ctx, id); err == nil {
					err = s.totalTSIUW(ctx, &d)
				}
			}
		}
		switch {
		case err != nil:
			log.Printf("nbfacin: total TSI kasus %s tidak dihitung: %v", id, err)
		case d.pesanKurs != "":
			log.Printf("nbfacin: total TSI kasus %s tidak dihitung: %s", id, d.pesanKurs)
		default:
			g.TotalTSINusaRe, g.TotalTSITopRisk = simpan8(d.totalNR), simpan8(d.totalTop)
		}
	}()
	go func() {
		defer wg.Done()
		sp, err := s.bacaSpreadingRingkas(ctx, nil, id)
		if err != nil {
			log.Printf("nbfacin: total RNM kasus %s tidak dihitung: %v", id, err)
			return
		}
		_, jumlah, _ := ringkasanRNM(sp)
		tsi, premi, tidakAda, err := totalKeIDR(ctx, s.spreading, sp.StartDateTime, jumlah)
		if err != nil || len(tidakAda) > 0 {
			log.Printf("nbfacin: total RNM kasus %s tidak dihitung: %v kurs tidak ada %v", id, err, tidakAda)
			return
		}
		g.TotalTSINusaReSpreading, g.TotalPremiNusaRe = simpan8(tsi), simpan8(premi)
	}()
	wg.Wait()
	return g
}
