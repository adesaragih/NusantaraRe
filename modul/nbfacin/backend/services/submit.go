package services

// Submit tahap marketing kasus FIRE NB (tiket 57): Email Type Accept / Reject / Banding / Decline + Comment. Bukti
// `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\`: `Section\EmailSection.xml`, `FlowAction\InwardFacultative.xml`
// (post-processing `SetValidateDate_PostAct` + DT `AddToListSuggestOfferFacIn_DT`, yang APPEND ViewSuggest-nya
// pyDisabled -> SATU baris riwayat per Submit), `Flow\InputInwardFacultativeOffer.xml` (Decision3: confirm / Banding /
// decline; konektor Decision29 -> Assignment14 TEAM LEADER), `Activity\SetBanding_ACT.xml`, `Activity\SaveViewSuggest.xml`,
// `Activity\InsertHistoryAkseptasiPega.xml`, `Activity\InsertFacInPerformance_Act.xml`, `Activity\SetCODRate.xml`,
// `DataTransform\SetCorrespondenceReadOnly.xml`. Pop-up Cause of Decline = port `frontend/components/sebabDecline.ts`.
//
// Keputusan work owner 07-10-2026 (diteruskan sesi 9d): Reject = case TETAP di marketing; Banding ikut Pega; semua proteksi
// Accept dibangun (yang datanya belum ada dilaporkan); riwayat ke TABEL LAMA; JSON_OFFER / FACINOFFER DITUNDA (nol CALL,
// ADR-0043); hanya tahap marketing; empat opsi. NB_STATUS (keputusan 9d dari korpus): Accept "NB IS IN TEAM LEADER'S
// INBOX" (konektor Decision29), Reject / Decline = SetValidateDate_PostAct 6.9 / 6.8 verbatim, Banding = tanpa nama orang.
// Tanpa email (SendOfferEmail langkah 1 Exit-Activity).
//
// Tiket 58 - Submit tahap TEAM LEADER (Assignment14, flow action yang sama `InwardFacultative`), keputusan work owner
// 07-10-2026 ("setuju" atas 6 saran + `DDL\EmailTypeTL.xml`): Email Type TL 1 Accept / 2 Reject / 3 Ask / 7 Decline; Accept
// -> Decision9 confirm -> `DataTransform\IsJuniorUW_B.xml` (posisi ReasFacInJuniorUnderwriting -> ReasFacInJuniorUnderwritingB ADR-0032, DateToUW, SubmitToUW bila
// kosong, NB status); Reject / Ask -> Utility8 SendEmailPolicy (TANPA email, butir 2) -> Transition170 kembali ke marketing;
// Decline -> End5 (Resolved-Decline, butir 5). Isian Deductible Acceptance + proteksi Protection_Act 9-12 (butir 4).

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

var (
	// ErrMasukanSubmit - isian Submit tidak sah. 400.
	ErrMasukanSubmit = errors.New("services: isian Submit tidak sah")
	// ErrSubmitBukanTahapMarketing - kasus bukan tahap marketing Offer (posisi / FlagOnGoingPolicy), atau posisinya berubah
	// selama Submit. 409. = ErrKasusBacaSaja (satu sumber aturan, bacasaja.go).
	ErrSubmitBukanTahapMarketing = ErrKasusBacaSaja
	// ErrSubmitTanpaDatabase - penyimpan Submit / kasus tidak terpasang. 503.
	ErrSubmitTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, Submit tidak dapat dijalankan")
	// ErrProteksiSubmit - proteksi Accept gagal; pesannya di GalatProteksi. 422.
	ErrProteksiSubmit = errors.New("services: proteksi Accept gagal")
)

// GalatProteksi - 422 + pesan proteksi verbatim; tidak ada yang ditulis.
type GalatProteksi struct{ Pesan []string }

func (g GalatProteksi) Error() string {
	return ErrProteksiSubmit.Error() + ": " + strings.Join(g.Pesan, "; ")
}

// Is - errors.Is(err, ErrProteksiSubmit).
func (g GalatProteksi) Is(target error) bool { return target == ErrProteksiSubmit }

const (
	// posisiTeamLeader - Assignment14 `pyRouterProp.pyCallParams.Workbasket`.
	posisiTeamLeader = "ReasFacInTeamLeader"
	// posisiJuniorUWB - Accept Team Leader: `IsJuniorUW_B` PositionNote (Assignment13 JUNIOR UNDERWITER B) "ReasFacInJuniorUnderwriting",
	// dibetulkan menjadi kode baru (keputusan work owner 07-10-2026, ADR-0032; tabel pemetaan issues/58).
	posisiJuniorUWB = "ReasFacInJuniorUnderwritingB"
	// lebarCommentDeductible - T_GENERAL_POLIS.COMMENT_DEDUCTIBLE VARCHAR2(500) (rancangan, migrasi 210).
	lebarCommentDeductible = 500
	// posisiDecline / statusDecline - Decline (End5; korpus tidak mengubah pyStatusWork): POSITION memakai kosakata
	// T_WORK_POLIS.POSITION yang sudah ada ("Decline") `[dugaan]`; STATUS_WORK "Resolved-Decline" = keputusan work owner
	// 07-10-2026.
	posisiDecline = "Decline"
	statusDecline = "Resolved-Decline"
	// posisiCedingConfirmOffer - `.IsCedingConfirm` tahap Offer (HISTORYAKSEPTASIPRODUCTION.POSISI; SetBanding juga "Offer").
	posisiCedingConfirmOffer = "Offer"
	// batasKomentar - SaveViewSuggest CARI10 `@substring(.CommentSuggest,0,3990)`; KETERANGAN VARCHAR2(4000) bita.
	batasKomentar     = 3990
	batasKomentarBita = 4000
	// bentukTglInp - InsertViewSuggest_SQL `To_date(CARI5,'DD/MM/YYYY HH24:MI:SS')`; Pega memformat "dd/MM/yyyy hh:mm:ss"
	// (jam 12 tanpa AM/PM - cacat, dibetulkan jam 24).
	bentukTglInp = "02/01/2006 15:04:05"
)

// labelApproval - per Email Type: riwayat = SaveViewSuggest CARI9 (HISTORYAKSEPTASIPRODUCTION.APPROVAL), pega =
// InsertHistoryAkseptasiPega CARI5 (HISTORYAKSEPTASIPEGA.STATUS).
var labelApproval = map[string]struct{ riwayat, pega string }{
	models.EmailTypeAccept:  {"Accept", "ACCEPT"},
	models.EmailTypeReject:  {"Reject", "REJECT"},
	models.EmailTypeAsk:     {"Ask", "ASK"},
	models.EmailTypeBanding: {"Banding", "BANDING"},
	models.EmailTypeDecline: {"Decline", "DECLINE"},
	// Tiket 64 - tahap Binding: SaveViewSuggest 1.2.2 "Revise"; InsertHistoryAkseptasiPega langkah 5 REVISE.
	models.EmailTypeRevise: {"Revise", "REVISE"},
}

// tujuanBanding - SetBanding_ACT langkah 3 (StatusBusiness bukan 3): LetterNo (= Banding to) -> PositionNote dan teks NB
// status. Teks Pega yang memuat NAMA ORANG tidak disalin (CLAUDE.md §4 butir 10) - diganti jabatannya (keputusan 9d
// 07-10-2026); yang tanpa nama dipakai verbatim. MANAGERTEKNIK (3.7) dan Financial (3.6) ber-label `//` - tidak ada.
var tujuanBanding = map[string]struct{ posisi, label string }{
	"JUW_A":              {"ReasFacInJuniorUnderwritingA", "JUNIOR UNDERWRITER (A)"},
	"JUW_B":              {"ReasFacInJuniorUnderwritingB", "JUNIOR UNDERWRITER (B)"},
	"UNDERWRITER":        {"ReasFacInUnderwriting", "UNDERWRITER"},
	"SENIORUW":           {"ReasFacInSeniorUnderwriting", "SENIORUW"},
	"DEPHEADUNDERWRITER": {"ReasFacInDepHeadUnderwriting", "DEPHEADUNDERWRITER"},
	"KADIVFACULTATIVE":   {"ReasFacInFacultativeDivHead", "KADIVFACULTATIVE"},
	"KADIVTEKNIK":        {"ReasFacInGroupLeader", "KADIVTEKNIK"},
	"DIREKTURMARKETING":  {"ReasFacInMarketingDirector", "DIREKTURMARKETING"},
	"DIREKTURTEKNIK":     {"ReasFacInTechnicalDirector", "DIREKTURTEKNIK"},
}

// emailTypeTahap - Email Type sah per tahap: marketing 1 / 2 / 4 / 7 (keputusan tiket 57 butir 7), Team Leader 1 / 2 / 3 / 7
// (`DDL\EmailTypeTL.xml`, tiket 58 butir 1).
var emailTypeTahap = map[string][]string{
	TahapMarketing:  {models.EmailTypeAccept, models.EmailTypeReject, models.EmailTypeBanding, models.EmailTypeDecline},
	TahapTeamLeader: {models.EmailTypeAccept, models.EmailTypeReject, models.EmailTypeAsk, models.EmailTypeDecline},
	// Tiket 61 - Underwriting `DDL\EmailTypeUW.xml` (ASM-FW-GISFW-WORK!EMAILTYPEUW) 1 / 2 / 3 / 7.
	TahapUW: {models.EmailTypeAccept, models.EmailTypeReject, models.EmailTypeAsk, models.EmailTypeDecline},
	// Tiket 64 - Binding `DDL\EmailTypeBinding.xml` (ASM-FW-GISFW-WORK!EMAILTYPEBINDING) 1 / 4 / 7 / 9.
	TahapBinding: {models.EmailTypeAccept, models.EmailTypeBanding, models.EmailTypeDecline, models.EmailTypeRevise},
}

// DenganSubmit memasang penyimpan Submit (tiket 57).
func (s *Service) DenganSubmit(p repository.PenyimpanSubmit) *Service { s.submit = p; return s }

// TampilanSubmit - GET …/submit.
type TampilanSubmit struct {
	BolehSubmit bool
	Alasan      string
	Riwayat     []models.RiwayatAkseptasi
	OpsiBanding []string
	SebabTc     []models.BarisSebabTc
	SebabNonTc  []string
	// PctRnmRate / PctRnmRIComm - SetCODRate: rate dari `TempTotal.pxResults(1).TotalNetRate` - halaman itu diisi
	// SumTSIPremiSpreadedRNM_FIRE_Act (per mata uang, lokasi Top Risk, Premium/TSI×1000) yang dipanggil pemeriksa spreading;
	// pemeriksa itu keluar untuk tahap marketing, jadi isi halaman saat popup marketing belum terverifikasi -> selalu ""
	// `[pertanyaan terbuka]`; RI Commission = `.Policy.Payment.RICommision` (PAY_RI_COMMISION).
	PctRnmRate, PctRnmRIComm string
	// Tahap - "marketing" / "teamLeader" / "" (tiket 58). DeductibleBedaMataUang - Pega `IsBedaCurr` (hanya dihitung di tahap
	// Team Leader); DeductibleAcceptance / CommentDeductible - isian TL tersimpan.
	Tahap                                        string
	DeductibleBedaMataUang, DeductibleAcceptance bool
	CommentDeductible                            string
	// Tiket 61 - tahap Underwriting: Validity Of Acceptance (hari) dan dua tanggal kabel DD-MM-YYYY; "" tahap lain.
	DaysValidity, DateValidity, WaitingBindDate string
	// Tiket 64 - tahap Binding: ConfirmBinding ("1") / ReceivedRiSlip ("true") tersimpan (migrasi 213); false tahap lain.
	ConfirmBinding, ReceivedRiSlip bool
}

// IsiPeriksa - badan POST …/submit/periksa (tiket 58: isian TL ikut, proteksi Deductible; tiket 61: daysValidity tahap uw).
type IsiPeriksa struct {
	EmailType            string
	DeductibleAcceptance bool
	CommentDeductible    string
	DaysValidity         string
}

// HasilSubmit - jawaban periksa / Submit.
type HasilSubmit struct {
	Pesan              []string
	Position, NBStatus string
	// PesanBinding - tiket 64b: periksa Accept tahap Binding = penahan tombol "Binding" (flagTombolBinding); [] selain itu.
	PesanBinding []string
	// NomorPolis - tiket 69a-2 (work owner 09-10-2026): nomor polis bila Submit menyelesaikan Policy; "" selain itu.
	NomorPolis string
}

// idPega - pzInsKey kasus.
func idPega(id string) string { return repository.AwalanIDPega + id }

// awalanNBStatus - `@replaceAll(pyWorkIDPrefix,"-","")` = "NB".
func awalanNBStatus() string { return strings.ReplaceAll(repository.AwalanCaseNB, "-", "") }

func (s *Service) siapSubmit(id string) error {
	if s.submit == nil {
		return ErrSubmitTanpaDatabase
	}
	if !idKasusSah(id) {
		return ErrKasusTidakAda
	}
	return nil
}

// konteksSubmit - baca kasus (404).
func (s *Service) konteksSubmit(ctx context.Context, id string) (models.KonteksSubmit, error) {
	k, err := s.submit.KonteksSubmit(ctx, id)
	if err != nil {
		return k, galatKasus(err)
	}
	return k, nil
}

// alasanTidakBoleh - tombol Submit Pega aktif bila `TempEmail.CARI28 == .PositionNote` (workbasket operator = posisi
// kasus) dan Email Type tampil (FlagOnGoingPolicy 0) = aturan baca-saja bersama (alasanUbah, bacasaja.go). Tiket 61: tahap uw
// (When IsUWSpread, FlagOnGoingPolicy 0) - workbasket POSITION (alasanUbahSpreading). Kosong = boleh.
func alasanTidakBoleh(k models.KonteksSubmit, pelaku inti.Pelaku) (string, error) {
	t := tahapKasusKonteks(k)
	if kasusSelesai(t.StatusWork) {
		return alasanUbah(t, pelaku, true)
	}
	// K69-8: TL tahap Policy tetap boleh Submit (form saja), walau tab baca-saja.
	if tahapSubmitUW(t) || tahapSubmitBinding(t) || tahapSubmitUWPolicy(t) || tahapTLPolicy(t) {
		if !pelaku.PunyaPeran(t.Position) {
			return "Akun tidak memegang workbasket " + t.Position + ".",
				fmt.Errorf("%w: kasus di posisi %s hanya dapat di-Submit workbasket itu", inti.ErrTanpaWewenang, t.Position)
		}
		return "", nil
	}
	return alasanUbah(t, pelaku, true)
}

// tahapSubmitUW - Submit tahap uw: (When IsUWSpread ATAU jenjang lanjutan tiket 62) + FlagOnGoingPolicy 0 (tiket 61, D7).
func tahapSubmitUW(t models.TahapKasus) bool {
	return (tahapUW(t) || posisiJenjangLanjut[t.Position]) && t.FlagOnGoingPolicy == flagScoringUbah
}

// unik - urutan pertama dipertahankan (M_LIMIT_PROPERTYY memuat jabatan kembar per TEAM_GROUP).
func unik(xs []string) []string {
	hasil := []string{}
	lihat := map[string]bool{}
	for _, x := range xs {
		if x = strings.TrimSpace(x); x != "" && !lihat[x] {
			lihat[x] = true
			hasil = append(hasil, x)
		}
	}
	return hasil
}

// TampilanSubmit - GET …/submit (401 tanpa sesi).
func (s *Service) TampilanSubmit(ctx context.Context, pelaku inti.Pelaku, id string) (TampilanSubmit, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TampilanSubmit{}, err
	}
	if err := s.siapSubmit(id); err != nil {
		return TampilanSubmit{}, err
	}
	k, err := s.konteksSubmit(ctx, id)
	if err != nil {
		return TampilanSubmit{}, err
	}
	t := TampilanSubmit{Tahap: tahapKonteks(k), DeductibleAcceptance: k.DeductibleAcceptance == "true",
		CommentDeductible: k.CommentDeductible}
	t.Alasan, _ = alasanTidakBoleh(k, pelaku)
	t.BolehSubmit = t.Alasan == ""
	if tahapTL(t.Tahap) {
		objek, err := s.BacaObjek(ctx, id)
		if err != nil {
			return TampilanSubmit{}, err
		}
		t.DeductibleBedaMataUang = deductibleBedaMataUang(objek)
	}
	if tahapUWForm(t.Tahap) {
		t.DaysValidity, t.DateValidity, t.WaitingBindDate = validitasTampilan(k)
	}
	if t.Tahap == TahapBinding {
		t.ConfirmBinding, t.ReceivedRiSlip = k.ConfirmBinding == "1", k.ReceivedRiSlip == "true"
	}
	if t.Riwayat, err = s.submit.RiwayatAkseptasi(ctx, idPega(id)); err != nil {
		return TampilanSubmit{}, err
	}
	banding, err := s.submit.OpsiBanding(ctx, idPega(id), k.TeamGroup)
	if err != nil {
		return TampilanSubmit{}, err
	}
	t.OpsiBanding = bandingAktif(banding)
	baris, err := s.submit.SebabDecline(ctx)
	if err != nil {
		return TampilanSubmit{}, err
	}
	t.SebabTc, t.SebabNonTc = susunSebabDecline(baris)
	if k.RICommision != nil {
		t.PctRnmRIComm = utils.FormatDecimal(k.RICommision)
	}
	return t, nil
}

// bandingAktif - pilihan GetAksepBanding_SQL yang punya tujuan aktif di SetBanding_ACT (keputusan work owner 07-10-2026:
// MANAGERTEKNIK, pemetaan 3.7 ber-`//`, disembunyikan - disaring generik), kembar sekali.
func bandingAktif(xs []string) []string {
	hasil := []string{}
	for _, x := range unik(xs) {
		if tujuanBanding[x].posisi != "" {
			hasil = append(hasil, x)
		}
	}
	return hasil
}

// teksBoolPega - boolean Pega disimpan teks "true" / "false" (butir 68.1).
func teksBoolPega(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// tahapKonteks - tahap kasus Submit (bacasaja.go tahapDari); tiket 61: When IsUWSpread + FlagOnGoingPolicy 0 = "uw" (Pega
// menampilkan EmailTypeUW juga pada flag 2 = tahap Binding - di luar cakupan, D7).
func tahapKonteks(k models.KonteksSubmit) string {
	t := tahapKasusKonteks(k)
	if tahapSubmitUW(t) {
		return TahapUW
	}
	if tahapSubmitBinding(t) {
		return TahapBinding
	}
	if tahapSubmitUWPolicy(t) {
		return TahapUWPolicy
	}
	if tahapTLPolicy(t) {
		return TahapTeamLeaderPolicy
	}
	return tahapDari(t)
}

// tahapKasusKonteks - TahapKasus dari konteks Submit (ConfirmBinding tidak dipakai penentu tahap Submit).
func tahapKasusKonteks(k models.KonteksSubmit) models.TahapKasus {
	return models.TahapKasus{Position: k.Position, FlagOnGoingPolicy: k.FlagOnGoingPolicy, StatusWork: k.StatusWork}
}

// emailTypeSah - Email Type sah tahap itu; pesan 400 menyebut daftarnya.
func emailTypeSah(tahap, e string) error {
	if dalam(emailTypeTahap[tahap], e) {
		return nil
	}
	return fmt.Errorf("%w: emailType tahap %s wajib salah satu %s", ErrMasukanSubmit, tahap, strings.Join(emailTypeTahap[tahap], " / "))
}

// deductibleTL - isian Team Leader untuk proteksi Accept (nil = tahap marketing, Protection_Act 9-12 dilewati).
func deductibleTL(tahap string, terima bool, komentar string) *isianDeductibleTL {
	if !tahapTL(tahap) {
		return nil
	}
	return &isianDeductibleTL{terima: terima, komentar: strings.TrimSpace(komentar)}
}

// PeriksaSubmit - POST …/submit/periksa: Accept -> proteksi (TL: + Deductible Acceptance); Reject / Ask / Banding /
// Decline -> kosong (Protection_Act 7-8). Kasus baca-saja (tahap "") -> 409.
func (s *Service) PeriksaSubmit(ctx context.Context, pelaku inti.Pelaku, id string, isi IsiPeriksa) (HasilSubmit, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HasilSubmit{}, err
	}
	if err := s.siapSubmit(id); err != nil {
		return HasilSubmit{}, err
	}
	k, err := s.konteksSubmit(ctx, id)
	if err != nil {
		return HasilSubmit{}, err
	}
	// Penjaga sama dengan Submit / POST …/validitas: 403 workbasket / 409 tahap (review tiket 61).
	if _, err := alasanTidakBoleh(k, pelaku); err != nil {
		return HasilSubmit{}, err
	}
	tahap := tahapKonteks(k)
	if err := emailTypeSah(tahap, isi.EmailType); err != nil {
		return HasilSubmit{}, err
	}
	if tahapUWForm(tahap) {
		if _, err := uraiHariValiditas(isi.DaysValidity); err != nil {
			return HasilSubmit{}, err
		}
	}
	pesan := []string{}
	switch {
	case tahap == TahapBinding:
		var pesanBinding []string
		if pesan, pesanBinding, err = s.periksaBinding(ctx, id, k, isi.EmailType); err != nil {
			return HasilSubmit{}, err
		}
		return HasilSubmit{Pesan: pesan, PesanBinding: pesanBinding}, nil
	case tahap == TahapMarketing && isi.EmailType == models.EmailTypeBanding:
		if pesan, err = s.proteksiBanding(ctx, id, flagBandingMarketing); err != nil {
			return HasilSubmit{}, err
		}
	case isi.EmailType != models.EmailTypeAccept:
	case tahapPolicy(tahap):
		// Tiket 69a: selesai Policy / FacRetro -> 409 juga di periksa.
		if pesan, _, err = s.periksaAcceptPolicy(ctx, id, k, pelaku, tahap,
			deductibleTL(tahap, isi.DeductibleAcceptance, isi.CommentDeductible)); err != nil {
			return HasilSubmit{}, err
		}
	case tahap == TahapUW:
		if pesan, _, err = s.periksaAcceptUW(ctx, id, k, pelaku); err != nil {
			return HasilSubmit{}, err
		}
	default:
		if pesan, err = s.ProteksiAccept(ctx, id, deductibleTL(tahap, isi.DeductibleAcceptance, isi.CommentDeductible)); err != nil {
			return HasilSubmit{}, err
		}
	}
	return HasilSubmit{Pesan: pesan}, nil
}

// potongKomentar - 3990 karakter (substring Pega), lalu ≤ 4000 bita (KETERANGAN VARCHAR2(4000) BYTE).
func potongKomentar(s string) string {
	if utf8.RuneCountInString(s) > batasKomentar {
		s = string([]rune(s)[:batasKomentar])
	}
	for len(s) > batasKomentarBita {
		_, n := utf8.DecodeLastRuneInString(s)
		s = s[:len(s)-n]
	}
	return s
}

// rapikanIsiSubmit - isian 400 (emailType per tahap, comment, bandingTo, Cause of Decline, Comment Deductible TL). Decline:
// Comment disusun ulang server.
func rapikanIsiSubmit(tahap string, isi models.IsiSubmit, opsiBanding []string, tc []models.BarisSebabTc, nonTc []string) (models.IsiSubmit, error) {
	out := isi
	out.Comment, out.BandingTo = strings.TrimSpace(isi.Comment), strings.TrimSpace(isi.BandingTo)
	out.CommentDeductible = strings.TrimSpace(isi.CommentDeductible)
	if err := emailTypeSah(tahap, out.EmailType); err != nil {
		return out, err
	}
	var masalah []string
	if len(out.CommentDeductible) > lebarCommentDeductible {
		masalah = append(masalah, fmt.Sprintf("commentDeductible paling banyak %d bita", lebarCommentDeductible))
	}
	if tahapUWForm(tahap) {
		// Tiket 61: tanpa Cause of Decline / Banding; Comment wajib untuk semua (EmailSection_IsUW), ke COMMENT_TEXT 500 bita.
		if _, err := uraiHariValiditas(out.DaysValidity); err != nil {
			return out, err
		}
		if out.Comment == "" {
			masalah = append(masalah, "comment wajib diisi")
		}
		if len(out.Comment) > lebarCommentUW {
			masalah = append(masalah, fmt.Sprintf("comment paling banyak %d bita", lebarCommentUW))
		}
		if len(masalah) > 0 {
			return out, fmt.Errorf("%w: %s", ErrMasukanSubmit, strings.Join(masalah, "; "))
		}
		return out, nil
	}
	switch out.EmailType {
	case models.EmailTypeBanding:
		ada := false
		for _, o := range opsiBanding {
			ada = ada || o == out.BandingTo
		}
		switch {
		case out.BandingTo == "":
			masalah = append(masalah, "bandingTo wajib diisi")
		case !ada:
			masalah = append(masalah, fmt.Sprintf("bandingTo %q bukan pilihan Banding to kasus ini", out.BandingTo))
		case tujuanBanding[out.BandingTo].posisi == "":
			masalah = append(masalah, fmt.Sprintf("jabatan %q tidak punya tujuan Banding (SetBanding_ACT)", out.BandingTo))
		}
	case models.EmailTypeDecline:
		pilihan, err := periksaSebab(out.SebabTc, out.SebabNonTc, tc, nonTc)
		if err != nil {
			return out, err
		}
		out.SebabTc = pilihan
		out.Comment = komentarDecline(pilihan, out.SebabNonTc)
	}
	if out.EmailType != models.EmailTypeDecline && out.Comment == "" {
		masalah = append(masalah, "comment wajib diisi")
	}
	if len(masalah) > 0 {
		return out, fmt.Errorf("%w: %s", ErrMasukanSubmit, strings.Join(masalah, "; "))
	}
	out.Comment = potongKomentar(out.Comment)
	return out, nil
}

// nbStatusSubmit - NB_STATUS dan posisi baru per tahap dan Email Type.
func nbStatusSubmit(tahap string, isi models.IsiSubmit, k models.KonteksSubmit, akun models.AkunSubmit) (posisi, status string) {
	a := awalanNBStatus()
	if tahapPolicy(tahap) {
		return nbStatusPolicy(tahap, isi, k, akun)
	}
	if tahap == TahapBinding {
		// Accept: posisi / NB status dari lengkapiTulisanBinding (TL / UW Financial). Banding: SetBanding_ACT (sama dengan
		// marketing). Revise: Assignment12 marketing. Decline: pola UW (dua spasi) - Pega tidak memasangnya (K64-2).
		switch isi.EmailType {
		case models.EmailTypeAccept:
			return "", ""
		case models.EmailTypeBanding:
			t := tujuanBanding[isi.BandingTo]
			return t.posisi, a + " IS IN " + t.label + "'S INBOX"
		case models.EmailTypeRevise:
			return repository.WorkbasketMarketing, a + " IS IN " + k.MarketingName + "'S INBOX"
		}
		return posisiDecline, a + " WAS DECLINED BY  " + strings.ToUpper(akun.Nama)
	}
	if tahap == TahapUW {
		// Accept: posisi / NB status dari jenjang (acceptUW, tulisanSubmit). Reject / Ask: Utility8 -> Transition170 (sama
		// dengan Team Leader). Decline: SetValidateDateUW_PostAct L3474 (dua spasi).
		switch isi.EmailType {
		case models.EmailTypeAccept:
			return "", "" // diisi lengkapiTulisanUW dari jenjang
		case models.EmailTypeReject, models.EmailTypeAsk:
			return repository.WorkbasketMarketing, a + " IS IN " + k.MarketingName + "'S INBOX"
		}
		return posisiDecline, a + " WAS DECLINED BY  " + strings.ToUpper(akun.Nama)
	}
	if tahap == TahapTeamLeader {
		switch isi.EmailType {
		case models.EmailTypeAccept:
			// IsJuniorUW_B verbatim.
			return posisiJuniorUWB, a + " IS IN JUNIOR UNDERWRITER (B)'S INBOX"
		case models.EmailTypeReject, models.EmailTypeAsk:
			// Transition170: MarketingName APA ADANYA (tanpa huruf besar, beda dengan Reject marketing 6.9).
			return repository.WorkbasketMarketing, a + " IS IN " + k.MarketingName + "'S INBOX"
		}
		// Decline TL: SetValidateDate 6.8 hanya untuk PositionNote Admin / Marketing (NextUser kosong -> "NB IS IN  'S
		// INBOX" di Pega) - teks Decline marketing dipakai `[keputusan agent]`.
		return posisiDecline, a + " WAS DECLINED BY  " + strings.ToUpper(akun.Nama)
	}
	switch isi.EmailType {
	case models.EmailTypeAccept:
		return posisiTeamLeader, a + " IS IN TEAM LEADER'S INBOX"
	case models.EmailTypeReject:
		return k.Position, a + " IS IN " + strings.ToUpper(strings.ReplaceAll(k.MarketingName, " -", "")) + "'S OPPORTUNITY"
	case models.EmailTypeBanding:
		t := tujuanBanding[isi.BandingTo]
		return t.posisi, a + " IS IN " + t.label + "'S INBOX"
	}
	return posisiDecline, a + " WAS DECLINED BY  " + strings.ToUpper(akun.Nama)
}

// AjukanSubmit - POST …/submit. Urutan: 401 -> 503 -> 404 -> 403 workbasket -> 409 posisi -> 400 isian -> 422 proteksi
// Accept (tanpa tulis) -> satu transaksi (409 bila posisi berubah serentak).
func (s *Service) AjukanSubmit(ctx context.Context, pelaku inti.Pelaku, id string, isi models.IsiSubmit) (HasilSubmit, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HasilSubmit{}, err
	}
	if err := s.siapSubmit(id); err != nil {
		return HasilSubmit{}, err
	}
	if s.transaksi == nil {
		return HasilSubmit{}, ErrSubmitTanpaDatabase
	}
	k, err := s.konteksSubmit(ctx, id)
	if err != nil {
		return HasilSubmit{}, err
	}
	if _, err := alasanTidakBoleh(k, pelaku); err != nil {
		return HasilSubmit{}, err
	}
	tahap := tahapKonteks(k)
	var opsiBanding, nonTc []string
	var tc []models.BarisSebabTc
	switch {
	case tahapUWForm(tahap): // tanpa Banding / Cause of Decline
	case isi.EmailType == models.EmailTypeBanding:
		if opsiBanding, err = s.submit.OpsiBanding(ctx, idPega(id), k.TeamGroup); err != nil {
			return HasilSubmit{}, err
		}
	case isi.EmailType == models.EmailTypeDecline:
		baris, err := s.submit.SebabDecline(ctx)
		if err != nil {
			return HasilSubmit{}, err
		}
		tc, nonTc = susunSebabDecline(baris)
	}
	isi, err = rapikanIsiSubmit(tahap, isi, opsiBanding, tc, nonTc)
	if err != nil {
		return HasilSubmit{}, err
	}
	var tujuan tujuanUW
	if tahap == TahapBinding {
		pesan, _, err := s.periksaBinding(ctx, id, k, isi.EmailType)
		if err != nil {
			return HasilSubmit{}, err
		}
		if len(pesan) > 0 {
			return HasilSubmit{}, GalatProteksi{Pesan: pesan}
		}
	} else if isi.EmailType == models.EmailTypeAccept || (tahap == TahapMarketing && isi.EmailType == models.EmailTypeBanding) {
		var pesan []string
		switch {
		case isi.EmailType == models.EmailTypeBanding: // tiket 57b
			pesan, err = s.proteksiBanding(ctx, id, flagBandingMarketing)
		case tahapPolicy(tahap): // tiket 69a
			pesan, tujuan, err = s.periksaAcceptPolicy(ctx, id, k, pelaku, tahap,
				deductibleTL(tahap, isi.DeductibleAcceptance, isi.CommentDeductible))
		case tahap == TahapUW:
			pesan, tujuan, err = s.periksaAcceptUW(ctx, id, k, pelaku)
		default:
			pesan, err = s.ProteksiAccept(ctx, id, deductibleTL(tahap, isi.DeductibleAcceptance, isi.CommentDeductible))
		}
		if err != nil {
			return HasilSubmit{}, err
		}
		if len(pesan) > 0 {
			return HasilSubmit{}, GalatProteksi{Pesan: pesan}
		}
	}
	akun, err := s.submit.AkunSubmit(ctx, pelaku.AkunID)
	if err != nil {
		return HasilSubmit{}, err
	}
	w := s.tulisanSubmit(tahap, id, pelaku, isi, k, akun)
	if tahap == TahapBinding {
		s.lengkapiTulisanBinding(&w, id, pelaku, isi, k)
	}
	if tahap == TahapUW {
		if err := s.lengkapiTulisanUW(ctx, &w, id, pelaku, isi, k, tujuan); err != nil {
			return HasilSubmit{}, err
		}
	}
	if tahapPolicy(tahap) {
		if err := s.lengkapiTulisanPolicy(ctx, &w, tahap, id, pelaku, isi, k, tujuan); err != nil {
			return HasilSubmit{}, err
		}
	}
	if err := s.transaksi(ctx, func(tx *db.Tx) error { return s.submit.TulisSubmit(ctx, tx, w) }); err != nil {
		if errors.Is(err, repository.ErrPosisiBerubah) {
			return HasilSubmit{}, ErrSubmitBukanTahapMarketing
		}
		return HasilSubmit{}, galatKasus(err)
	}
	log.Printf("nbfacin: Submit kasus %s %s oleh %s -> %s", id, labelApproval[isi.EmailType].riwayat, pelaku.AkunID, w.PosisiBaru)
	if w.Binding != nil {
		s.kirimSurelBinding(ctx, id, k, w.NBStatus)
	}
	h := HasilSubmit{Pesan: []string{}, Position: w.PosisiBaru, NBStatus: w.NBStatus}
	if w.Policy != nil {
		h.NomorPolis = w.Policy.NomorPolis
	}
	return h, nil
}

// tulisanSubmit - seluruh baris tulis (SaveViewSuggest, InsertHistoryAkseptasiPega, InsertFacInPerformance_Act; TL: isian
// Deductible, IsJuniorUW_B).
func (s *Service) tulisanSubmit(tahap, id string, pelaku inti.Pelaku, isi models.IsiSubmit, k models.KonteksSubmit, akun models.AkunSubmit) models.TulisanSubmit {
	nama := akun.Nama
	if nama == "" {
		nama = pelaku.AkunID // OperatorID.pyUserName tak tersedia `[dugaan]`
	}
	label := labelApproval[isi.EmailType]
	posisi, status := nbStatusSubmit(tahap, isi, k, akun)
	w := models.TulisanSubmit{CaseID: id, IDPega: idPega(id), PosisiLama: k.Position, PosisiBaru: posisi, NBStatus: status,
		Riwayat: models.RiwayatProduksi{TypePolis: awalanNBStatus(), Posisi: posisiCedingConfirmOffer, PIC: nama,
			TglInp: s.sekarang().In(WIB).Format(bentukTglInp), Div: akun.Divisi, Type: "F", Putaran: "1", Approval: label.riwayat,
			Keterangan: isi.Comment, AksesLogin: pelaku.AkunID, PercentRNM: k.PercentShare},
		Pega: models.RiwayatPega{Status: label.pega, Username: nama, Workbasket: workbasketTahap[tahap]}}
	if tahapTL(tahap) {
		w.Deductible = &models.DeductibleTL{Acceptance: teksBoolPega(isi.DeductibleAcceptance), Comment: isi.CommentDeductible}
		if tahap == TahapTeamLeader && isi.EmailType == models.EmailTypeAccept {
			w.KeUW = s.sekarang().UTC().Format(BentukWaktuPega)
		}
	}
	if isi.EmailType == models.EmailTypeDecline {
		w.StatusWorkBaru = statusDecline
	}
	// FACINPERFORMANCE hanya Decline MARKETING: InsertFacInPerformance_Act tidak membaca EmailTypeTL `[terverifikasi]`.
	if isi.EmailType == models.EmailTypeDecline && tahap == TahapMarketing {
		f := models.PerformanceFacIn{PositionDoc: "OFFER", StatusDoc: "DECLINE", Username: pelaku.AkunID, InsuredName: k.InsuredName,
			SOB: k.SobName, MOLeader: k.MOLeader, COB: k.ClassOfBusiness}
		// Langkah 12: dilewati bila salah satu kolom NOT NULL kosong (seperti penjaga Pega).
		if f.Username != "" && f.InsuredName != "" && f.SOB != "" && f.MOLeader != "" && f.COB != "" {
			w.Performance = &f
		} else {
			log.Printf("nbfacin: Submit Decline kasus %s: FACINPERFORMANCE dilewati (kolom wajib kosong)", id)
		}
	}
	return w
}

// ---- Cause of Decline (port frontend/components/sebabDecline.ts) -------------------------------------------------

// Syarat kolom popup (`labels.ts` SEBAB_DECLINE.term).
var (
	termOption       = []string{"Period Of Insurance (POI)", "Extension of coverage", "Risk Exposure"}
	termRemarksWajib = []string{"Occupancy", "TSI (in IDR)", "Lack of data for acceptance", "Others T/C"}
)

const (
	termRate      = "Rate"
	termRIComm    = "R/I Comm, Deduction"
	termOJK       = "OJK Discount"
	termLossRatio = "Loss Ratio (net)"
	termLainnya   = "Others"
	// pesan* - CekCauseOfDecline_ACT verbatim (spasi depan ikut); pesanSebabKosong = TEKS_SEBAB_DECLINE.kosong (K57-5).
	pesanPctRnmSob   = " PctRNM or PctSOB Can't be NULL"
	pesanPctRnm      = " PctRNM Can't be NULL"
	pesanRemaks      = " Remaks Can't be NULL"
	pesanOption      = "Option Can't be NULL"
	pesanRemarks     = "Remarks Can't be NULL"
	pesanSebabKosong = "Pilih minimal satu Cause of Decline."
)

func dalam(daftar []string, x string) bool {
	for _, d := range daftar {
		if d == x {
			return true
		}
	}
	return false
}

// validasiSebab - validasiSebab sebabDecline.ts: pesan verbatim, unik, berurutan.
func validasiSebab(daftar []models.PilihanSebabTc, nonTc string) []string {
	var p pengumpulPesan
	for _, b := range daftar {
		if b.Prioritas == "" {
			continue
		}
		t := b.TermCondition
		if (t == termRate || t == termRIComm || t == termOJK) && (b.PctRNM == "" || b.PctSOB == "") {
			p.tambah(pesanPctRnmSob)
		}
		if t == termLossRatio && b.PctRNM == "" {
			p.tambah(pesanPctRnm)
		}
		if dalam(termRemarksWajib, t) && b.Description2 == "" {
			p.tambah(pesanRemaks)
		}
		if dalam(termOption, t) && b.Description1 == "" {
			p.tambah(pesanOption)
		}
		if (t == "Extension of coverage" || t == "Risk Exposure") && b.Description1 == termLainnya && b.Description2 == "" {
			p.tambah(pesanRemarks)
		}
	}
	satu := false
	for _, b := range daftar {
		satu = satu || b.Prioritas == "I"
	}
	if nonTc == "" && !satu {
		p.tambah(pesanSebabKosong)
	}
	return p.pesan
}

// teksSebabTc - teksSebabTc sebabDecline.ts (SetCauseOfDecline_ACT 5.4).
func teksSebabTc(b models.PilihanSebabTc) string {
	s := b.TermCondition + " " + b.Description1 + " " + b.Description2 + " " + b.PctRNM
	if b.PctRNM != "" {
		s += " % RNM "
	}
	s += b.PctSOB
	if b.PctSOB != "" {
		s += " % SOB "
	}
	return s
}

// komentarDecline - komentarDecline sebabDecline.ts: alasan I (atau Non - T/C) + "\n" + alasan II.
func komentarDecline(daftar []models.PilihanSebabTc, nonTc string) string {
	var satu, dua *models.PilihanSebabTc
	for i := range daftar {
		if daftar[i].Prioritas == "I" && satu == nil {
			satu = &daftar[i]
		}
		if daftar[i].Prioritas == "II" && dua == nil {
			dua = &daftar[i]
		}
	}
	alasan1, alasan2 := nonTc, ""
	if nonTc == "" && satu != nil {
		alasan1 = teksSebabTc(*satu)
	}
	if dua != nil {
		alasan2 = teksSebabTc(*dua)
	}
	return alasan1 + "\n" + alasan2
}

// susunSebabDecline - baris CAUSEOFDECLINE urut ID angka: T/C = GetCauseOfDecline (`DESCRIPTION IS NULL`); opsi T/C X =
// TERMCONDITION baris ber-DESCRIPTION X dan TERMCONDITION2 kosong (sumber Pega `Getdescription_ACT` tidak di korpus -
// `[dugaan]` kuat dari bentuk data; teks = TERMCONDITION, bukan DESCRIPTION1 - ID 24 ber-DESCRIPTION1 salah); Non - T/C =
// GetCauseOfDecline2 (`TERMCONDITION2 IS NOT NULL`). Teks dipangkas: ID 10 "Clauses" tersimpan berakhiran baris baru sehingga
// syarat Pega `.TERMCONDITION != 'Clauses'` tak pernah cocok - dipangkas = maksud Pega (pembetulan, keputusan 9d 07-10-2026).
func susunSebabDecline(baris []models.BarisCauseOfDecline) ([]models.BarisSebabTc, []string) {
	urut := append([]models.BarisCauseOfDecline{}, baris...)
	angka := func(id string) (int, bool) { n, err := strconv.Atoi(strings.TrimSpace(id)); return n, err == nil }
	sort.SliceStable(urut, func(i, j int) bool {
		a, okA := angka(urut[i].ID)
		b, okB := angka(urut[j].ID)
		if okA && okB {
			return a < b
		}
		return okA && !okB
	})
	tc := []models.BarisSebabTc{}
	nonTc := []string{}
	for _, b := range urut {
		term, desc, term2 := strings.TrimSpace(b.TermCondition), strings.TrimSpace(b.Description), strings.TrimSpace(b.TermCondition2)
		if term != "" && desc == "" {
			tc = append(tc, models.BarisSebabTc{TermCondition: term, Opsi: []string{}})
		}
		if term != "" && term2 != "" {
			nonTc = append(nonTc, term)
		}
	}
	for i := range tc {
		for _, b := range urut {
			if strings.TrimSpace(b.Description) == tc[i].TermCondition && strings.TrimSpace(b.TermCondition2) == "" {
				if o := strings.TrimSpace(b.TermCondition); o != "" {
					tc[i].Opsi = append(tc[i].Opsi, o)
				}
			}
		}
	}
	return tc, nonTc
}

// periksaSebab - baris T/C dicentang saja, TANPA dipangkas (Comment = teks yang sama dengan popup, komentarDecline TS);
// T/C dan Non - T/C wajib dari daftar server; Option (description1) T/C ber-opsi wajib salah satu opsinya; satu baris per
// prioritas; persen desimal sah; lalu validasiSebab. Selain validasiSebab = pertahanan server, bukan bagian sebabDecline.ts.
func periksaSebab(pilihan []models.PilihanSebabTc, nonTc string, tc []models.BarisSebabTc, daftarNonTc []string) ([]models.PilihanSebabTc, error) {
	opsi := map[string][]string{}
	for _, x := range tc {
		opsi[x.TermCondition] = x.Opsi
	}
	var masalah []string
	dipakai := []models.PilihanSebabTc{}
	prioritas := map[string]bool{}
	for i, b := range pilihan {
		switch {
		case b.Prioritas == "":
			continue
		case b.Prioritas != "I" && b.Prioritas != "II":
			masalah = append(masalah, fmt.Sprintf("sebabTc[%d].prioritas wajib I / II", i))
		case prioritas[b.Prioritas]:
			masalah = append(masalah, fmt.Sprintf("prioritas %s hanya untuk satu T/C", b.Prioritas))
		case opsi[b.TermCondition] == nil:
			masalah = append(masalah, fmt.Sprintf("sebabTc[%d].termCondition %q bukan pilihan T/C", i, b.TermCondition))
		case b.Description1 != "" && len(opsi[b.TermCondition]) > 0 && !dalam(opsi[b.TermCondition], b.Description1):
			masalah = append(masalah, fmt.Sprintf("sebabTc[%d].description1 %q bukan Option %s", i, b.Description1, b.TermCondition))
		}
		prioritas[b.Prioritas] = true
		UraiDesimalIsian(fmt.Sprintf("sebabTc[%d].pctRnm", i), b.PctRNM, &masalah)
		UraiDesimalIsian(fmt.Sprintf("sebabTc[%d].pctSob", i), b.PctSOB, &masalah)
		dipakai = append(dipakai, b)
	}
	if nonTc != "" && !dalam(daftarNonTc, nonTc) {
		masalah = append(masalah, fmt.Sprintf("sebabNonTc %q bukan pilihan Non - T/C", nonTc))
	}
	if len(masalah) == 0 {
		masalah = validasiSebab(dipakai, nonTc)
	}
	if len(masalah) > 0 {
		return nil, fmt.Errorf("%w: %s", ErrMasukanSubmit, strings.Join(masalah, "; "))
	}
	return dipakai, nil
}
