package services

// Tiket 60 (permintaan work owner 08-10-2026 "di uw juga bisa ubah acumulation"): Accumulation Code / Address coverage dapat
// diubah di tahap Underwriting lewat PUT /api/nbfacin/kasus/{caseId}/akumulasi; field coverage lain tetap baca-saja, PUT
// …/objek tetap 409.
//
// `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\Section\CoverageItem_IsUW.xml`: `.AccumulationCode` / `.AccumulationDescription`
// read-only - nilainya hanya berubah lewat tombol Choose / Copy Accumulation, yang kontainernya tampil bila (L7247)
// PositionNote ReasFacInUnderwriting / ReasFacInSeniorUnderwriting / ReasFacInJuniorUnderwriting / ReasFacInJuniorUnderwritingA.
// K60-2 penyimpangan sadar [keputusan work owner 08-10-2026 "DepHead dan UnderwritingFinancial juga bisa dong ubah
// accumulation"]: tahap = tahapUW (When IsUWSpread, enam posisi termasuk DepHead / UnderwritingFinancial; JUW dibaca B,
// ADR-0032), workbasket wajib = POSITION - aturan yang sama dengan tab Spreading.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

var (
	// ErrAkumulasiHanyaUW - PUT …/akumulasi di luar tahap uw. 409.
	ErrAkumulasiHanyaUW = errors.New("Ubah Accumulation Code lewat tab Coverage hanya tersedia di tahap Underwriting - " +
		"di tahap marketing / Team Leader pakai Save tab Object")
	// ErrObjekBerubah - bentuk objek / item / coverage badan PUT …/akumulasi tidak sama dengan tersimpan. 409.
	ErrObjekBerubah = errors.New("Objek berubah sejak tab Coverage dimuat - muat ulang tab")
)

// DenganAkumulasiCoverage memasang penulis akumulasi coverage (tiket 60).
func (s *Service) DenganAkumulasiCoverage(p repository.PenyimpanAkumulasiCoverage) *Service {
	s.akumulasiCoverage = p
	return s
}

// alasanUbahAkumulasi - di luar tahap uw (tahapUW, K60-2) -> 409 ErrAkumulasiHanyaUW; tahap Policy (flag 1) -> 409
// ErrTahapPolicyBacaSaja (K69-8: "Tidak, hanya Spreading"); tanpa workbasket POSITION -> 403.
func alasanUbahAkumulasi(t models.TahapKasus, pelaku inti.Pelaku) error {
	if t.FlagOnGoingPolicy == flagPolicy && !kasusSelesai(t.StatusWork) {
		return ErrTahapPolicyBacaSaja
	}
	if !tahapUW(t) {
		return ErrAkumulasiHanyaUW
	}
	if !pelaku.PunyaPeran(t.Position) {
		return fmt.Errorf("%w: akumulasi kasus di posisi %s hanya dapat diubah workbasket itu", inti.ErrTanpaWewenang, t.Position)
	}
	return nil
}

// AksesAkumulasi - `ubahAkumulasi` GET kasus: tahap uw + pemegang workbasket POSITION. LUNAK: galat baca tahap atau tanpa
// pembaca tahap -> false (tanpa pembaca tahap tidak ada tahap uw).
func (s *Service) AksesAkumulasi(ctx context.Context, pelaku inti.Pelaku, id string) bool {
	t, ada, err := s.tahapKasus(ctx, id)
	if !ada {
		return false
	}
	if err != nil {
		log.Printf("nbfacin: tahap kasus %s tidak terbaca, akumulasi dianggap baca-saja: %v", id, err)
		return false
	}
	return alasanUbahAkumulasi(t, pelaku) == nil
}

// SimpanAkumulasi - PUT …/akumulasi badan daftar objek (bentuk PUT …/objek). Bentuk lokasi / item / coverage + oldId wajib
// sama dengan tersimpan (409); HANYA accumulationCode / accumulationDescription coverage yang ditulis, medan lain diabaikan
// (tanpa hitung ulang premi / coverage / spreading). Validasi dua medan itu = aturan PUT …/objek (lebarCoverage 50 / 500
// byte; PUT objek tidak memeriksa kode terhadap master). Tanpa pembaca tahap = penjaga tidak aktif. Jawab baca ulang objek.
func (s *Service) SimpanAkumulasi(ctx context.Context, pelaku inti.Pelaku, id string, baris []models.ObjekFire) ([]models.ObjekFire, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	t, ada, err := s.tahapKasus(ctx, id)
	if err != nil {
		return nil, err
	}
	if ada {
		if err := alasanUbahAkumulasi(t, pelaku); err != nil {
			return nil, err
		}
	}
	if m := periksaIsianAkumulasi(baris); len(m) > 0 {
		return nil, fmt.Errorf("%w: %s", ErrMasukanObjek, strings.Join(m, "; "))
	}
	if s.objek == nil || s.spreading == nil || s.akumulasiCoverage == nil || s.transaksi == nil {
		return nil, ErrObjekTanpaDatabase
	}
	if !idKasusSah(id) {
		return nil, ErrKasusTidakAda
	}
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		// Pohon coverage tersimpan (ID baris; urut lokasi / item / coverage SEQ_NO = posisi GET …/objek, seperti K48-7).
		k, err := s.spreading.BacaSpreading(ctx, tx, id)
		if err != nil {
			return err
		}
		tulis, ok := akumulasiMenurutPosisi(k, baris)
		if !ok {
			return ErrObjekBerubah
		}
		return s.akumulasiCoverage.UbahAkumulasiCoverage(ctx, tx, id, tulis)
	})
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	if err != nil {
		return nil, err
	}
	return s.BacaObjek(ctx, id)
}

// periksaIsianAkumulasi - lebarCoverage untuk accumulationCode / accumulationDescription saja.
func periksaIsianAkumulasi(baris []models.ObjekFire) []string {
	var masalah []string
	for i, o := range baris {
		for j, it := range o.Items {
			for m, c := range it.Coverages {
				for _, l := range lebarCoverage {
					if (l.nama == "accumulationCode" || l.nama == "accumulationDescription") && len(l.nilai(c)) > l.n {
						masalah = append(masalah, fmt.Sprintf("baris[%d].items[%d].coverages[%d].%s paling banyak %d byte", i, j, m, l.nama, l.n))
					}
				}
			}
		}
	}
	return masalah
}

// akumulasiMenurutPosisi - pasangan (ID coverage tersimpan, akumulasi badan) menurut posisi; ok=false bila jumlah lokasi /
// item / coverage atau oldId coverage berbeda.
func akumulasiMenurutPosisi(k models.KasusSpreading, baris []models.ObjekFire) ([]models.AkumulasiCoverage, bool) {
	if len(k.Lokasi) != len(baris) {
		return nil, false
	}
	var out []models.AkumulasiCoverage
	for i, l := range k.Lokasi {
		if len(l.Items) != len(baris[i].Items) {
			return nil, false
		}
		for j, it := range l.Items {
			badan := baris[i].Items[j].Coverages
			if len(it.Coverages) != len(badan) {
				return nil, false
			}
			for m, c := range it.Coverages {
				if c.OldID != badan[m].OldID {
					return nil, false
				}
				out = append(out, models.AkumulasiCoverage{ID: c.ID, AccumulationCode: badan[m].AccumulationCode,
					AccumulationDescription: badan[m].AccumulationDescription})
			}
		}
	}
	return out, true
}
