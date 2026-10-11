package services

// Copy To All Spreading per lokasi (tiket 73 P5, K73-40 keputusan work owner 10-10-2026: SEMUA COB, FIRE juga). Port
// `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\Activity\CopyToAllLocSpreading_ACT.xml` (`ASM-FW-GISFW-DATA!COPYTOALLLOCSPREADING_ACT`,
// tombol "Copy To All Spreading" `Section\PropertyItemListCoverageSpreading_IsUW.xml` sel L4444-4809, pyVisible ALWAYS, kontainer
// (IsUWSpread) L659):
//   - langkah 5 FIRE (lokasi Param.IdxLoc): SpreadingList tiap coverage DIGANTI template; TSISpreaded / PremiumSpreaded /
//     TSIGrossSpreaded menurut share (L2433-2504) = salinTemplate dibatasi satu lokasi;
//   - langkah 6 Aneka (IsAneka ∨ Bonding) / 7 Golf: sama tanpa TSIGross (L3677 / L3729, L4937 / L4989) = salinTemplateAneka dibatasi
//     satu lokasi.
//
// Tidak menyimpan (seperti salin case-level); tidak memanggil cekSpreadingFactIn / SumTSI (Pega) -> tanpa pesan jenis treaty.
// K73-50: efek IsFlagUW langkah 2-3 TIDAK diport (tidak disimpan). K73-51: flagSpreading (Type Facultative) wajib bila treatyType
// "10015" (`pyRequiredWhen=.TreatyType=="10015"` L3020), terisi = salah satu opsiFlagSpreading, tidak disimpan. Langkah 4 (OfferFacIn.FlagSpreading) tidak disimpan.

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

// pesanFlagSpreadingWajib - K73-51 (teks agent).
const pesanFlagSpreadingWajib = "flagSpreading wajib untuk treaty 10015"

// opsiFlagSpreading - `D:\migrasi\RNM\DDL\FlagSpreading.xml` (`ASM-FW-GISFW-DATA-SPREADINGRISK!FLAGSPREADING`, diekspor work owner
// 10-10-2026): nilai lokal pyStandardValue = pyLocalizedValue L186-192 (sensus dua cara: grep 2, urai XML 2). Di luar daftar -> 400.
var opsiFlagSpreading = map[string]bool{"SF-MR 2023": true, "Open Fac Market": true}

// periksaTemplateLokasi - penjaga + validasi salin-lokasi = salin case-level (tahap uw + workbasket, % Share, template ≥ 1 baris,
// periksaBarisShare, proteksiTemplate) + flagSpreading (K73-51) + idxLokasi 1..n.
func (s *Service) periksaTemplateLokasi(ctx context.Context, pelaku inti.Pelaku, id string, ps *apd.Decimal, templat []models.TemplateSpreading) error {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if err := s.wajibTahapUW(ctx, pelaku, id); err != nil {
		return err
	}
	if err := periksaShareKasus(ps); err != nil {
		return err
	}
	if len(templat) == 0 {
		return fmt.Errorf("%w: template wajib berisi paling sedikit satu baris", ErrMasukanSpreading)
	}
	jenis, nama, share := make([]string, len(templat)), make([]string, len(templat)), make([]*apd.Decimal, len(templat))
	var masalah []string
	for i, t := range templat {
		jenis[i], nama[i], share[i] = t.TreatyType, t.TreatyName, t.SharePercentage
		f := strings.TrimSpace(t.FlagSpreading)
		if strings.TrimSpace(t.TreatyType) == treatyFACOUT && f == "" {
			masalah = append(masalah, fmt.Sprintf("template[%d].%s", i, pesanFlagSpreadingWajib))
		}
		if f != "" && !opsiFlagSpreading[f] {
			masalah = append(masalah, fmt.Sprintf("template[%d].flagSpreading %q bukan pilihan (SF-MR 2023 / Open Fac Market)", i, f))
		}
	}
	masalah = append(periksaBarisShare("template", jenis, nama, share), masalah...)
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanSpreading, strings.Join(masalah, "; "))
	}
	if s.spreading == nil {
		return ErrSpreadingTanpaDatabase
	}
	m, err := s.proteksiTemplate(ctx, templat)
	if err != nil {
		return err
	}
	if len(m) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanSpreading, strings.Join(m, "; "))
	}
	return nil
}

// idxLokasiSah - 1-based (Param.IdxLoc = .pxListSubscript) -> indeks 0-based; di luar rentang -> 400.
func idxLokasiSah(idx, n int) (int, error) {
	if idx < 1 || idx > n {
		return 0, fmt.Errorf("%w: idxLokasi %d di luar 1..%d", ErrMasukanSpreading, idx, n)
	}
	return idx - 1, nil
}

// SalinSpreadingLokasi - POST …/spreading/salin-lokasi (FIRE, langkah 5).
func (s *Service) SalinSpreadingLokasi(ctx context.Context, pelaku inti.Pelaku, id string, ps *apd.Decimal, idxLokasi int,
	templat []models.TemplateSpreading) (TampilanSpreading, error) {
	if err := s.periksaTemplateLokasi(ctx, pelaku, id, ps, templat); err != nil {
		return TampilanSpreading{}, err
	}
	k, err := s.bacaSpreading(ctx, nil, id)
	if err != nil {
		return TampilanSpreading{}, err
	}
	i, err := idxLokasiSah(idxLokasi, len(k.Lokasi))
	if err != nil {
		return TampilanSpreading{}, err
	}
	k.PercentShare = ps
	hitungNR(&k)
	satuLokasi := models.KasusSpreading{PercentShare: ps, Lokasi: k.Lokasi[i : i+1]} // berbagi larik: coverage k ikut berubah
	salinTemplate(&satuLokasi, templat)
	return s.tampilan(ctx, id, k, nil, false)
}

// SalinSpreadingAnekaLokasi - POST …/spreading-aneka/salin-lokasi (Aneka langkah 6 / Golf langkah 7).
func (s *Service) SalinSpreadingAnekaLokasi(ctx context.Context, pelaku inti.Pelaku, id string, ps *apd.Decimal, idxLokasi int,
	templat []models.TemplateSpreading) (TampilanSpreadingAneka, error) {
	if err := s.periksaTemplateLokasi(ctx, pelaku, id, ps, templat); err != nil {
		return TampilanSpreadingAneka{}, err
	}
	d, err := s.bacaDataSpreadingAneka(ctx, id)
	if err != nil {
		return TampilanSpreadingAneka{}, err
	}
	i, err := idxLokasiSah(idxLokasi, len(d.lokasi))
	if err != nil {
		return TampilanSpreadingAneka{}, err
	}
	d.kasus.PercentShare = ps
	hitungNRAneka(d.lokasi, ps)
	salinTemplateAneka(d.lokasi[i:i+1], ps, templat, d.lini) // berbagi larik: objek d.lokasi ikut berubah
	return s.tampilanSpreadingAneka(ctx, id, d, nil, false)
}
