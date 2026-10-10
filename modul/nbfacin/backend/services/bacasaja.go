package services

// Case baca-saja sesudah keluar dari tahap marketing. Keputusan work owner 07-10-2026 ("setuju, buatkan jangan halu", atas
// saran sesi 9d): Pega `finishAssignment` - sesudah Submit marketing tidak lagi memegang assignment kasus, jadi layar
// marketing tidak dapat mengubahnya. SATU sumber aturan untuk:
//   - `bolehUbah` / `alasanBacaSaja` GET kasus dan jawaban PUT …/general,
//   - `bolehSubmit` / `alasan` Submit (tiket 57),
//   - BolehUbah Scoring Risk (CallDualScoringRisk "Write": tahap saja, tanpa pelaku),
//   - penjaga setiap endpoint yang MENULIS data kasus (409 tahap, 403 workbasket), dicek sesudah 401 dan SEBELUM menulis.
// Tab IT-only (Loss Record, Correspondence) memakai syarat TAHAP saja: akses tab itu sudah JOB_POSITION "IT" (tiket 55 / 56) -
// menambah syarat workbasket marketing akan membuatnya hanya dapat diubah akun yang memegang keduanya `[keputusan agent]`
// (dipertahankan work owner 07-10-2026).
// Tiket 58 (keputusan work owner 07-10-2026 butir 3): Team Leader boleh mengubah isi tab dan Submit selama case di posisinya
// - tahap "teamLeader" (POSITION ReasFacInTeamLeader, flag 0) bersama "marketing"; syarat workbasket = workbasket POSISI saat
// itu (Pega: pemegang assignment).
// Tiket 64 lanjutan (work owner 08-10-2026: "inikan lagi dimarketing harusnya msh bisa di edit, kecuali sudah di klik binding,
// baru smua tidak bisa diedit lagi"): tahap "binding" (POSITION ReasFacInMarketing, flag 2) dapat DIUBAH pemegang workbasket
// marketing selama CONFIRM_BINDING != "1"; sesudah tombol Binding ditekan seluruhnya baca-saja (409).

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services/rules"
)

// ErrKasusBacaSaja - kasus sudah tidak berada di tahap marketing Offer (POSITION / FLAG_ONGOING_POLICY). 409.
var ErrKasusBacaSaja = errors.New("services: case baca-saja: sudah tidak berada di tahap marketing")

// DenganTahap memasang pembaca tahap kasus. nil = penjaga baca-saja tidak aktif (layanan tanpa basis data / uji yang tidak
// memasangnya); DariDasar SELALU memasangnya (dijaga TestDariDasarMemasangTahap).
func (s *Service) DenganTahap(t repository.PembacaTahapKasus) *Service { s.tahap = t; return s }

// Tahap kasus yang dapat diubah / di-Submit (kontrak `TampilanSubmit.tahap`).
const (
	TahapMarketing  = "marketing"
	TahapTeamLeader = "teamLeader"
)

// workbasketTahap - workbasket pemegang tiap tahap (Assignment12 MARKETING / Assignment14 TEAM LEADER).
var workbasketTahap = map[string]string{TahapMarketing: repository.WorkbasketMarketing, TahapTeamLeader: posisiTeamLeader,
	TahapBinding: repository.WorkbasketMarketing}

// tahapMarketing - PositionNote ReasFacInMarketing (posisiMarketing) dan FlagOnGoingPolicy 0 (tahap Offer). Juga syarat
// "Write" Scoring Risk (CallDualScoringRisk - marketing saja).
func tahapMarketing(t models.TahapKasus) bool {
	return posisiMarketing[t.Position] && t.FlagOnGoingPolicy == flagScoringUbah
}

// awalanSelesai - STATUS_WORK case selesai ("Resolved-Decline", "Resolved-Completed"), seperti saringan kotak masuk.
const awalanSelesai = "Resolved-"

// kasusSelesai - case sudah selesai: baca-saja di tahap mana pun (tiket 69a-2: selesai Policy tidak mengubah POSITION).
func kasusSelesai(statusWork string) bool { return strings.HasPrefix(statusWork, awalanSelesai) }

// tahapDari - "marketing" / "teamLeader" / "binding" (sebelum tombol Binding) / "" (baca-saja).
func tahapDari(t models.TahapKasus) string {
	switch {
	case kasusSelesai(t.StatusWork):
		return ""
	case tahapMarketing(t):
		return TahapMarketing
	case t.Position == posisiTeamLeader && t.FlagOnGoingPolicy == flagScoringUbah:
		return TahapTeamLeader
	case tahapSubmitBinding(t) && t.ConfirmBinding != confirmBindingYa:
		return TahapBinding
	}
	return ""
}

// pesanTahapPolicy - alasan baca-saja FlagOnGoingPolicy 1 (K69-8).
const pesanTahapPolicy = "Case baca-saja: tahap Policy (sudah di-bind)."

// ErrTahapPolicyBacaSaja - K69-8 (work owner 09-10-2026 "jika flagongoingpolicy = 1, ini data2nya tidak bisa dirubah lagi.
// karena itu sudah di bind, jadi data gaboleh berubah"): seluruh tab baca-saja di flag 1 (TL / UW tahap Policy), kecuali form
// Submit dan - RALAT K69-8 - tab Spreading tahap uw (alasanUbahSpreading, tahapUW). 409 (membungkus ErrKasusBacaSaja).
var ErrTahapPolicyBacaSaja = fmt.Errorf("%w: tahap Policy (sudah di-bind)", ErrKasusBacaSaja)

// confirmBindingYa - CONFIRM_BINDING sesudah tombol Binding (teks Pega "1").
const confirmBindingYa = "1"

// alasanUbah - kosong = boleh. Tahap di luar marketing / Team Leader -> 409; akun tanpa workbasket posisi saat itu -> 403
// (bila wajibWorkbasket).
func alasanUbah(t models.TahapKasus, pelaku inti.Pelaku, wajibWorkbasket bool) (string, error) {
	if kasusSelesai(t.StatusWork) {
		return "Case baca-saja: case sudah selesai (" + t.StatusWork + ").", ErrKasusBacaSaja
	}
	if t.FlagOnGoingPolicy == flagPolicy {
		return pesanTahapPolicy, ErrTahapPolicyBacaSaja
	}
	tahap := tahapDari(t)
	if tahap == "" && tahapSubmitBinding(t) {
		return "Case baca-saja: tombol Binding sudah ditekan.", ErrKasusBacaSaja
	}
	if tahap == "" {
		return fmt.Sprintf("Case baca-saja: sudah tidak berada di tahap marketing / Team Leader (posisi %q).", t.Position), ErrKasusBacaSaja
	}
	if wb := workbasketTahap[tahap]; wajibWorkbasket && !pelaku.PunyaPeran(wb) {
		return "Akun tidak memegang workbasket " + wb + ".",
			fmt.Errorf("%w: kasus di posisi %s hanya dapat diubah workbasket %s", inti.ErrTanpaWewenang, t.Position, wb)
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

// ErrPeriodeTerkunci - tiket 64d K1 (work owner 08-10-2026 "Ikut Pega"): Begin / End date tidak dapat diubah di tahap
// Binding (`Section\Periode.xml` `.PolicyData.StartDateTime` / `.EndDateTime` `.ConfirmBinding = 1||pyWorkPage.FlagOnGoingPolicy=2`).
// 409 (membungkus ErrKasusBacaSaja).
var ErrPeriodeTerkunci = fmt.Errorf("%w: Begin Date dan End Date tidak dapat diubah di tahap Binding", ErrKasusBacaSaja)

// KunciPeriode - `kunciPeriode` GET kasus: case di tahap Binding (TampilanSubmit.tahap "binding": POSITION ReasFacInMarketing,
// flag 2 - sebelum maupun sesudah tombol Binding) ATAU jalur Policy (tiket 69a: TL / UW flag 1 - ConfirmBinding sudah 1, Periode.xml
// `.ConfirmBinding = 1||FlagOnGoingPolicy=2`). LUNAK: tanpa pembaca tahap / galat baca -> false.
func (s *Service) KunciPeriode(ctx context.Context, id string) bool {
	t, ada, err := s.tahapKasus(ctx, id)
	return ada && err == nil && (tahapSubmitBinding(t) || tahapTLPolicy(t) || tahapSubmitUWPolicy(t))
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

// ---- Tahap Underwriting: tab Spreading (tiket 59) -----------------------------------------------------------------
//
// Keputusan work owner 07-10-2026 ("di pega ada when IsUWSpread, itu kondisi untuk ubah spreading ini, jadi di uw bisa ubah
// % Share RNM dan add delete ubah spreading") dan 08-10-2026 ("spreading hanya bisa di copy atau add yg masuk ke when IsUW, di
// marketing tidak bisa add spreading"; meralat K48 S-1). Tahap "uw" = When `IsUWSpread` (`NB FacIn\When\IsUWSpread.xml`,
// dievaluasi dari registry bangkitan - ReasFacInJuniorUnderwriting dibaca B, ADR-0032): HANYA tab Spreading yang terbuka;
// penjaga tulis tab lain (wajibBolehUbah) tetap 409. Syarat workbasket = POSITION saat itu (pola Team Leader).
//   - PUT …/spreading: marketing / Team Leader / uw (+ workbasket posisi); baris badan hanya dipakai di uw;
//   - POST …/spreading/salin dan …/hitung-baris: HANYA uw (+ workbasket posisi);
//   - POST …/spreading/hitung-share: tidak menulis, tidak dijaga (seperti sebelum tiket 59).

// TahapUW - tahap Underwriting (When IsUWSpread).
const TahapUW = "uw"

// ErrSpreadingHanyaUW - Copy Spreading / hitung baris spreading di luar tahap Underwriting. 409.
var ErrSpreadingHanyaUW = errors.New("Copy Spreading dan ubah baris spreading hanya tersedia di tahap Underwriting")

// halamanPosisi - rules.Kasus dengan satu properti `pyWorkPage.PositionNote` (satu-satunya yang diuji IsUWSpread).
type halamanPosisi string

// Nilai - rules.Kasus: hanya `pyWorkPage.PositionNote`; jalur lain tidak ada (teks kosong).
func (p halamanPosisi) Nilai(jalur string) (string, bool) {
	if jalur == "pyWorkPage.PositionNote" {
		return string(p), true
	}
	return "", false
}

// tahapUW - When IsUWSpread atas POSITION; FlagOnGoingPolicy tidak diuji (rule-nya hanya PositionNote). Galat Eval
// (ErrTafsirBerbeda - hanya mungkin bila POSITION berupa angka) = bukan uw, dicatat.
func tahapUW(t models.TahapKasus) bool {
	if kasusSelesai(t.StatusWork) {
		return false
	}
	ok, err := rules.Eval("ISUWSPREAD", halamanPosisi(t.Position))
	if err != nil {
		log.Printf("nbfacin: When IsUWSpread atas posisi %q gagal dinilai, dianggap bukan uw: %v", t.Position, err)
		return false
	}
	return ok
}

// alasanUbahSpreading - (tahap, alasan, galat) tab Spreading: uw + workbasket posisi (403); selain itu aturan alasanUbah.
func alasanUbahSpreading(t models.TahapKasus, pelaku inti.Pelaku) (string, string, error) {
	if tahapUW(t) {
		if !pelaku.PunyaPeran(t.Position) {
			return TahapUW, "Akun tidak memegang workbasket " + t.Position + ".",
				fmt.Errorf("%w: spreading kasus di posisi %s hanya dapat diubah workbasket itu", inti.ErrTanpaWewenang, t.Position)
		}
		return TahapUW, "", nil
	}
	alasan, err := alasanUbah(t, pelaku, true)
	return tahapDari(t), alasan, err
}

// penjagaPutSpreading - penjaga PUT …/spreading (409 / 403 seperti alasanUbahSpreading). barisBadan = baris spreading badan
// dipakai: tahap uw, atau penjaga tidak aktif (tanpa pembaca tahap, perilaku tiket 48); marketing / Team Leader -> false.
// uw = tahap uw nyata (langkah 11 CountPremiAndTSI ikut ditulis); penjaga tidak aktif -> false.
func (s *Service) penjagaPutSpreading(ctx context.Context, pelaku inti.Pelaku, id string) (barisBadan, uw bool, err error) {
	t, ada, err := s.tahapKasus(ctx, id)
	if err != nil || !ada {
		return !ada, false, err
	}
	tahap, _, err := alasanUbahSpreading(t, pelaku)
	return tahap == TahapUW, tahap == TahapUW, err
}

// wajibTahapUW - penjaga salin / hitung-baris: di luar tahap uw -> 409 ErrSpreadingHanyaUW; uw tanpa workbasket posisi -> 403.
func (s *Service) wajibTahapUW(ctx context.Context, pelaku inti.Pelaku, id string) error {
	t, ada, err := s.tahapKasus(ctx, id)
	if err != nil || !ada {
		return err
	}
	if !tahapUW(t) {
		return ErrSpreadingHanyaUW
	}
	_, _, err = alasanUbahSpreading(t, pelaku)
	return err
}

// AksesSpreading - `bolehUbahSpreading` / `ubahBarisSpreading` GET kasus. LUNAK seperti AksesUbahKasus: galat baca tahap ->
// keduanya false (dicatat); tanpa pembaca tahap -> keduanya true (penjaga tidak aktif, endpoint pun tidak menolak).
func (s *Service) AksesSpreading(ctx context.Context, pelaku inti.Pelaku, id string) (bool, bool) {
	t, ada, err := s.tahapKasus(ctx, id)
	if !ada {
		return true, true
	}
	if err != nil {
		log.Printf("nbfacin: tahap kasus %s tidak terbaca, spreading dianggap baca-saja: %v", id, err)
		return false, false
	}
	tahap, alasan, _ := alasanUbahSpreading(t, pelaku)
	boleh := tahap != "" && alasan == ""
	return boleh, boleh && tahap == TahapUW
}
