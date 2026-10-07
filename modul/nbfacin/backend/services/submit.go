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
	// awalanIDPega - pzInsKey kelas kerja NB Fac In (`InsertHistoryAkseptasiPega` langkah 8 "ASM-FW-GISFW-WORK "+idpega;
	// index DBA HISTOAKSPEGA_INDEX4 memotong awalan yang sama).
	awalanIDPega = "ASM-FW-GISFW-WORK "
	// posisiTeamLeader - Assignment14 `pyRouterProp.pyCallParams.Workbasket`.
	posisiTeamLeader = "ReasFacInTeamLeader"
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
	models.EmailTypeBanding: {"Banding", "BANDING"},
	models.EmailTypeDecline: {"Decline", "DECLINE"},
}

// tujuanBanding - SetBanding_ACT langkah 3 (StatusBusiness bukan 3): LetterNo (= Banding to) -> PositionNote dan teks NB
// status. Teks Pega yang memuat NAMA ORANG tidak disalin (CLAUDE.md §4 butir 10) - diganti jabatannya (keputusan 9d
// 07-10-2026); yang tanpa nama dipakai verbatim. MANAGERTEKNIK (3.7) dan Financial (3.6) ber-label `//` - tidak ada.
var tujuanBanding = map[string]struct{ posisi, label string }{
	"JUW_A":              {"ReasFacInJuniorUnderwritingA", "JUNIOR UNDERWRITER (A)"},
	"JUW_B":              {"ReasFacInJuniorUnderwriting", "JUNIOR UNDERWRITER (B)"},
	"UNDERWRITER":        {"ReasFacInUnderwriting", "UNDERWRITER"},
	"SENIORUW":           {"ReasFacInSeniorUnderwriting", "SENIORUW"},
	"DEPHEADUNDERWRITER": {"ReasFacInDepHeadUnderwriting", "DEPHEADUNDERWRITER"},
	"KADIVFACULTATIVE":   {"ReasFacInFacultativeDivHead", "KADIVFACULTATIVE"},
	"KADIVTEKNIK":        {"ReasFacInGroupLeader", "KADIVTEKNIK"},
	"DIREKTURMARKETING":  {"ReasFacInMarketingDirector", "DIREKTURMARKETING"},
	"DIREKTURTEKNIK":     {"ReasFacInTechnicalDirector", "DIREKTURTEKNIK"},
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
}

// HasilSubmit - jawaban periksa / Submit.
type HasilSubmit struct {
	Pesan              []string
	Position, NBStatus string
}

// idPega - pzInsKey kasus.
func idPega(id string) string { return awalanIDPega + id }

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
// kasus) dan Email Type tampil (FlagOnGoingPolicy 0) = aturan baca-saja bersama (alasanUbah, bacasaja.go). Kosong = boleh.
func alasanTidakBoleh(k models.KonteksSubmit, pelaku inti.Pelaku) (string, error) {
	return alasanUbah(models.TahapKasus{Position: k.Position, FlagOnGoingPolicy: k.FlagOnGoingPolicy}, pelaku, true)
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
	t := TampilanSubmit{}
	t.Alasan, _ = alasanTidakBoleh(k, pelaku)
	t.BolehSubmit = t.Alasan == ""
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

// emailTypeSah - empat opsi tahap marketing (keputusan butir 7).
func emailTypeSah(e string) bool { _, ada := labelApproval[e]; return ada }

// PeriksaSubmit - POST …/submit/periksa: Accept -> proteksi; Reject / Banding / Decline -> kosong (Protection_Act 7-8).
func (s *Service) PeriksaSubmit(ctx context.Context, pelaku inti.Pelaku, id, emailType string) (HasilSubmit, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HasilSubmit{}, err
	}
	if err := s.siapSubmit(id); err != nil {
		return HasilSubmit{}, err
	}
	if !emailTypeSah(emailType) {
		return HasilSubmit{}, fmt.Errorf("%w: emailType wajib 1 / 2 / 4 / 7", ErrMasukanSubmit)
	}
	if _, err := s.konteksSubmit(ctx, id); err != nil {
		return HasilSubmit{}, err
	}
	pesan := []string{}
	if emailType == models.EmailTypeAccept {
		var err error
		if pesan, err = s.ProteksiAccept(ctx, id); err != nil {
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

// rapikanIsiSubmit - isian 400 (emailType, comment, bandingTo, Cause of Decline). Decline: Comment disusun ulang server.
func rapikanIsiSubmit(isi models.IsiSubmit, opsiBanding []string, tc []models.BarisSebabTc, nonTc []string) (models.IsiSubmit, error) {
	out := isi
	out.Comment, out.BandingTo = strings.TrimSpace(isi.Comment), strings.TrimSpace(isi.BandingTo)
	if !emailTypeSah(out.EmailType) {
		return out, fmt.Errorf("%w: emailType wajib 1 / 2 / 4 / 7", ErrMasukanSubmit)
	}
	var masalah []string
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

// nbStatusSubmit - NB_STATUS dan posisi baru per Email Type.
func nbStatusSubmit(isi models.IsiSubmit, k models.KonteksSubmit, akun models.AkunSubmit) (posisi, status string) {
	a := awalanNBStatus()
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
	var opsiBanding, nonTc []string
	var tc []models.BarisSebabTc
	switch isi.EmailType {
	case models.EmailTypeBanding:
		if opsiBanding, err = s.submit.OpsiBanding(ctx, idPega(id), k.TeamGroup); err != nil {
			return HasilSubmit{}, err
		}
	case models.EmailTypeDecline:
		baris, err := s.submit.SebabDecline(ctx)
		if err != nil {
			return HasilSubmit{}, err
		}
		tc, nonTc = susunSebabDecline(baris)
	}
	isi, err = rapikanIsiSubmit(isi, opsiBanding, tc, nonTc)
	if err != nil {
		return HasilSubmit{}, err
	}
	if isi.EmailType == models.EmailTypeAccept {
		pesan, err := s.ProteksiAccept(ctx, id)
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
	w := s.tulisanSubmit(id, pelaku, isi, k, akun)
	if err := s.transaksi(ctx, func(tx *db.Tx) error { return s.submit.TulisSubmit(ctx, tx, w) }); err != nil {
		if errors.Is(err, repository.ErrPosisiBerubah) {
			return HasilSubmit{}, ErrSubmitBukanTahapMarketing
		}
		return HasilSubmit{}, galatKasus(err)
	}
	log.Printf("nbfacin: Submit kasus %s %s oleh %s -> %s", id, labelApproval[isi.EmailType].riwayat, pelaku.AkunID, w.PosisiBaru)
	return HasilSubmit{Pesan: []string{}, Position: w.PosisiBaru, NBStatus: w.NBStatus}, nil
}

// tulisanSubmit - seluruh baris tulis (SaveViewSuggest, InsertHistoryAkseptasiPega, InsertFacInPerformance_Act).
func (s *Service) tulisanSubmit(id string, pelaku inti.Pelaku, isi models.IsiSubmit, k models.KonteksSubmit, akun models.AkunSubmit) models.TulisanSubmit {
	nama := akun.Nama
	if nama == "" {
		nama = pelaku.AkunID // OperatorID.pyUserName tak tersedia `[dugaan]`
	}
	label := labelApproval[isi.EmailType]
	posisi, status := nbStatusSubmit(isi, k, akun)
	w := models.TulisanSubmit{CaseID: id, IDPega: idPega(id), PosisiLama: k.Position, PosisiBaru: posisi, NBStatus: status,
		Riwayat: models.RiwayatProduksi{TypePolis: awalanNBStatus(), Posisi: posisiCedingConfirmOffer, PIC: nama,
			TglInp: s.sekarang().In(WIB).Format(bentukTglInp), Div: akun.Divisi, Type: "F", Putaran: "1", Approval: label.riwayat,
			Keterangan: isi.Comment, AksesLogin: pelaku.AkunID, PercentRNM: k.PercentShare},
		Pega: models.RiwayatPega{Status: label.pega, Username: nama, Workbasket: repository.WorkbasketMarketing}}
	if isi.EmailType == models.EmailTypeDecline {
		w.StatusWorkBaru = statusDecline
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
