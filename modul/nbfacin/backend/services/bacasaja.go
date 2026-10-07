package services

// Case baca-saja sesudah keluar dari tahap marketing. Keputusan work owner 07-10-2026 ("setuju, buatkan jangan halu", atas
// saran sesi 9d): Pega `finishAssignment` - sesudah Submit marketing tidak lagi memegang assignment kasus, jadi layar
// marketing tidak dapat mengubahnya. SATU sumber aturan untuk:
//   - `bolehUbah` / `alasanBacaSaja` GET kasus dan jawaban PUT …/general,
//   - `bolehSubmit` / `alasan` Submit (tiket 57),
//   - BolehUbah Scoring Risk (CallDualScoringRisk "Write": tahap saja, tanpa pelaku),
//   - penjaga setiap endpoint yang MENULIS data kasus (409 tahap, 403 workbasket), dicek sesudah 401 dan SEBELUM menulis.
// Tab IT-only (Loss Record, Correspondence) memakai syarat TAHAP saja: akses tab itu sudah JOB_POSITION "IT" (tiket 55 / 56) -
// menambah syarat workbasket marketing akan membuatnya hanya dapat diubah akun yang memegang keduanya `[keputusan agent]`.

import (
	"context"
	"errors"
	"fmt"
	"log"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// ErrKasusBacaSaja - kasus sudah tidak berada di tahap marketing Offer (POSITION / FLAG_ONGOING_POLICY). 409.
var ErrKasusBacaSaja = errors.New("services: case baca-saja: sudah tidak berada di tahap marketing")

// DenganTahap memasang pembaca tahap kasus. nil = penjaga baca-saja tidak aktif (layanan tanpa basis data / uji yang tidak
// memasangnya); DariDasar SELALU memasangnya (dijaga TestDariDasarMemasangTahap).
func (s *Service) DenganTahap(t repository.PembacaTahapKasus) *Service { s.tahap = t; return s }

// tahapMarketing - PositionNote ReasFacInMarketing (posisiMarketing) dan FlagOnGoingPolicy 0 (tahap Offer).
func tahapMarketing(t models.TahapKasus) bool {
	return posisiMarketing[t.Position] && t.FlagOnGoingPolicy == flagScoringUbah
}

// alasanUbah - kosong = boleh. Workbasket (403) diperiksa lebih dulu daripada tahap (409) bila wajibWorkbasket.
func alasanUbah(t models.TahapKasus, pelaku inti.Pelaku, wajibWorkbasket bool) (string, error) {
	if wajibWorkbasket && !pelaku.PunyaPeran(repository.WorkbasketMarketing) {
		return "Akun tidak memegang workbasket " + repository.WorkbasketMarketing + ".",
			fmt.Errorf("%w: kasus tahap marketing hanya dapat diubah workbasket %s", inti.ErrTanpaWewenang, repository.WorkbasketMarketing)
	}
	if !tahapMarketing(t) {
		return fmt.Sprintf("Case baca-saja: sudah tidak berada di tahap marketing Offer (posisi %q).", t.Position), ErrKasusBacaSaja
	}
	return "", nil
}

// tahapKasus - (tahap, terpasang, galat); 404 bila kasus tidak ada.
func (s *Service) tahapKasus(ctx context.Context, id string) (models.TahapKasus, bool, error) {
	if s.tahap == nil {
		return models.TahapKasus{}, false, nil
	}
	if !idKasusSah(id) {
		return models.TahapKasus{}, true, ErrKasusTidakAda
	}
	t, err := s.tahap.TahapKasus(ctx, id)
	if err != nil {
		return t, true, galatKasus(err)
	}
	return t, true, nil
}

// wajibBolehUbah - penjaga endpoint tulis: workbasket ReasFacInMarketing (403) dan tahap marketing (409).
func (s *Service) wajibBolehUbah(ctx context.Context, pelaku inti.Pelaku, id string) error {
	t, ada, err := s.tahapKasus(ctx, id)
	if err != nil || !ada {
		return err
	}
	_, err = alasanUbah(t, pelaku, true)
	return err
}

// wajibTahapUbah - penjaga tab IT-only (Loss Record, Correspondence): tahap marketing saja (409).
func (s *Service) wajibTahapUbah(ctx context.Context, id string) error {
	t, ada, err := s.tahapKasus(ctx, id)
	if err != nil || !ada {
		return err
	}
	_, err = alasanUbah(t, inti.Pelaku{}, false)
	return err
}

// AksesUbahKasus - `bolehUbah` / `alasanBacaSaja` GET kasus. LUNAK: galat baca tahap -> baca-saja (dicatat), supaya layar
// kasus tidak gagal; tanpa pembaca tahap -> boleh (penjaga tidak aktif).
func (s *Service) AksesUbahKasus(ctx context.Context, pelaku inti.Pelaku, id string) (bool, string) {
	t, ada, err := s.tahapKasus(ctx, id)
	if !ada {
		return true, ""
	}
	if err != nil {
		log.Printf("nbfacin: tahap kasus %s tidak terbaca, dianggap baca-saja: %v", id, err)
		return false, "Tahap case tidak terbaca."
	}
	alasan, _ := alasanUbah(t, pelaku, true)
	return alasan == "", alasan
}
