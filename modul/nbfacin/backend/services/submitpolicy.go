package services

// Jalur Policy (tiket 69a-1) - Team Leader dan jenjang UW tahap Policy (FlagOnGoingPolicy 1). Kontrak sesi 9d, keputusan work
// owner K69-1..7 dan keputusan agent A69-1..3 di docs/issues/69-jalur-policy.md. Bukti `[terverifikasi]`
// `D:\migrasi\RNM\NB FacIn\`:
//   - Flow\InputInwardFacultativeRISlip.xml: Start1 (FlagOnGoingPolicy 1, IsCedingConfirm "Policy"); Assignment6 TEAM LEADER
//     (flow action InwardFacultative) -> Decision24 confirm -> Decision21 `IsFlagUW` -> Utility8 GetLimitAkseptasi_JUW_UW ->
//     Decision33 (ToJUW_A / ToUW / ToSeniorUW / ToDepHeadUW; Else -> Decision4) / Else -> Utility1 SaveJsonPolicyFacIn_Act;
//     Assignment8 UNDERWRITING (Decision16) dan Assignment15 JUNIOR UNDERWITER A (Decision32) confirm -> Utility8; Assignment14
//     SENIOR (Decision22), Assignment10 DEP HEAD (Decision29), Assignment1 UW FINANCIAL (Decision2) confirm -> Decision4 FAC OUT?
//     (IsFacRetro -> Decision3 IsInputFacRetro -> Utility1 / Else SubProcess3 OfferFacRetro; Else -> Utility1).
//   - Decision33 ToSeniorUW / ToDepHeadUW men-SET NBStatus berisi NAMA ORANG - tidak disalin: teks jabatan statusInboxPosisi
//     (pola tiket 61 D7); ToUW / ToJUW_A verbatim.
//   - Flow\InputInwardFacultativeOffer.xml Transition69 (Decision14 reject): PositionNote ReasFacInMarketing, IsCedingConfirm
//     "Binding", FlagOnGoingPolicy 2, NBStatus "NB IS IN " + MarketingName + "'S INBOX".
//   - Activity\GetLimitAkseptasi_JUW_UW.xml langkah 7: PositionNote ReasFacInTeamLeader -> CARI5 "LEADER"; langkah 8 IsGroup ->
//     keluar (LetterNo "" -> Else).
//   - Activity\CekLimitSpreading_Act.xml langkah 3-4 + RDBList\GetCountAccumulation.xml (IsFlagUW, K69-5).
//   - Activity\SetValidateDate_PostAct.xml 3.5 / SetValidateDateUW_PostAct.xml 5.4 + SaveViewSuggest L793: riwayat POSISI
//     "Policy", PUTARAN "2"; InsertFacInPerformance_Act langkah 5 POSITION_DOC "POLICY", langkah 13 dilewati di flag 1.
//
// "Selesai Policy" (SaveJsonPolicyFacIn_Act) = tiket 69a-2 (selesaipolicy.go). Jalur FacRetro BELUM dibangun: Accept yang
// berujung ke sana DITOLAK 409 (tanpa status setengah jadi) sampai tiket 65.

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// Tahap jalur Policy (kontrak TampilanSubmit.tahap, A69-3).
const (
	TahapTeamLeaderPolicy = "teamLeaderPolicy"
	TahapUWPolicy         = "uwPolicy"
	// jabatanLeader - GetLimitAkseptasi_JUW_UW langkah 7 CARI5 untuk PositionNote ReasFacInTeamLeader.
	jabatanLeader = "LEADER"
	// treatyORS / treatyFACOUT / treatySPL - IsFlagUW K69-5: TreatyType memuat "10007" / "10015" atau TreatyName memuat "SPL".
	teksTreatySPL = "SPL"
	// bentukTanggalAkumulasi - CekLimitSpreading_Act 4.1.1.1 `@FormatDateTime(StartDateTime,"dd/MM/yyyy","Asia/Jakarta")`.
	bentukTanggalAkumulasi = "02/01/2006"
)

// ErrFacRetroBelumTersedia - Accept berujung SubProcess3 OfferFacRetro (IsFacRetro, bukan IsInputFacRetro) - tiket 65. 409.
var ErrFacRetroBelumTersedia = errors.New("Case FacRetro belum dapat diselesaikan di sistem baru (tiket 65).")

// posisiUWPolicy - posisi UW tahap Policy (Assignment8 / 14 / 10 / 1 / 15); jenjang = Accept naik lewat Utility8.
var posisiUWPolicy = map[string]bool{posisiUWUtama: true, posisiSeniorUW: true, posisiDepHeadUW: true, posisiUWFinancial: true,
	posisiJuniorUWA: true}

// jabatanJenjangPolicy - posisi yang Accept-nya menjalankan GetLimitAkseptasi_JUW_UW (Utility8) di tahap Policy.
var jabatanJenjangPolicy = map[string]string{posisiTeamLeader: jabatanLeader, posisiUWUtama: "UNDERWRITER", posisiJuniorUWA: "JUW_A"}

// emailTypeTahap jalur Policy: TL `DDL\EmailTypeTL.xml` 1 / 2 / 3 / 7 (K69-2); UW `DDL\EmailTypeUWPolicy.xml` tanpa Banding
// (K69-6) 1 / 7.
func init() {
	emailTypeTahap[TahapTeamLeaderPolicy] = []string{models.EmailTypeAccept, models.EmailTypeReject, models.EmailTypeAsk,
		models.EmailTypeDecline}
	emailTypeTahap[TahapUWPolicy] = []string{models.EmailTypeAccept, models.EmailTypeDecline}
	workbasketTahap[TahapTeamLeaderPolicy] = posisiTeamLeader
}

// tahapTLPolicy - Assignment6 TEAM LEADER tahap Policy.
func tahapTLPolicy(t models.TahapKasus) bool {
	return t.Position == posisiTeamLeader && t.FlagOnGoingPolicy == flagPolicy && !kasusSelesai(t.StatusWork)
}

// tahapSubmitUWPolicy - posisi UW tahap Policy + flag 1.
func tahapSubmitUWPolicy(t models.TahapKasus) bool {
	return posisiUWPolicy[t.Position] && t.FlagOnGoingPolicy == flagPolicy && !kasusSelesai(t.StatusWork)
}

// tahapPolicy - salah satu tahap jalur Policy.
func tahapPolicy(tahap string) bool { return tahap == TahapTeamLeaderPolicy || tahap == TahapUWPolicy }

// tahapTL - Team Leader tahap Offer atau Policy (form InwardFacultative yang sama: Deductible Acceptance, Cause of Decline).
func tahapTL(tahap string) bool { return tahap == TahapTeamLeader || tahap == TahapTeamLeaderPolicy }

// tahapUWForm - form EmailSection_IsUW (Comment wajib, Validity Of Acceptance, tanpa Cause of Decline / Banding).
func tahapUWForm(tahap string) bool { return tahap == TahapUW || tahap == TahapUWPolicy }

// tanggalAkumulasi - StartDateTime teks Pega -> dd/MM/yyyy WIB; tak terbaca / kosong -> "".
func tanggalAkumulasi(s string) string {
	t, err := time.Parse(BentukWaktuPega, strings.TrimSpace(s))
	if err != nil {
		return ""
	}
	return t.In(WIB).Format(bentukTanggalAkumulasi)
}

// isFlagUW - K69-5 (CekLimitSpreading_Act langkah 4, Copy To All spreading): 1 bila (a) coverage ber-AccumulationCode yang
// GetCountAccumulation-nya > 0 (FACINPRODUCTION + FACINOFFER, case sendiri dikecualikan, ENDDATE >= StartDateTime polis), atau
// (b) baris spreading ber-TreatyType memuat "10007" / "10015" atau TreatyName memuat "SPL". StartDateTime kosong -> TO_DATE teks kosong
// NULL di Pega (hitung 0) -> kueri dilewati.
func (s *Service) isFlagUW(ctx context.Context, id string) (bool, error) {
	if aneka, ok, err := s.isFlagUWAneka(ctx, id); err != nil || ok {
		return aneka, err
	}
	sp, err := s.bacaSpreading(ctx, nil, id)
	if err != nil {
		return false, err
	}
	for _, l := range sp.Lokasi {
		for _, it := range l.Items {
			for _, c := range it.Coverages {
				for _, b := range c.Spreading {
					if strings.Contains(b.TreatyType, treatyORS) || strings.Contains(b.TreatyType, treatyFACOUT) ||
						strings.Contains(b.TreatyName, teksTreatySPL) {
						return true, nil
					}
				}
			}
		}
	}
	tgl := tanggalAkumulasi(sp.StartDateTime)
	if tgl == "" {
		return false, nil
	}
	objek, err := s.BacaObjek(ctx, id)
	if err != nil {
		return false, err
	}
	sudah := map[string]bool{}
	for _, o := range objek {
		for _, it := range o.Items {
			for _, c := range it.Coverages {
				kode := c.AccumulationCode
				if kode == "" || sudah[kode] {
					continue
				}
				sudah[kode] = true
				n, err := s.submit.HitungAkumulasiBerjalan(ctx, kode, tgl, idPega(id), id)
				if err != nil {
					return false, err
				}
				if n > 0 {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

// isFlagUWAneka - tiket 73 P5 (K73-37 / K73-41, keputusan work owner 10-10-2026): case keluarga Aneka (IsAneka ∨
// IsBondingAndCustomBonds ∨ isGolfInsurance; ok=true) - CekLimitSpreading_Act langkah 4 (akumulasi / 10007) hanya IsFire, langkah 5
// mengisi 1 bila flag kosong dan bukan polis deklarasi (L2105), polis deklarasi (When IsDeclarationPolicy `pyWorkPage.Quotation.
// PolicyType = 2`, L137) -> langkah 7 = 0 (L2407). IsFlagUW tidak disimpan sistem baru: flag dianggap kosong. ok=false (lini tak
// terbaca / bukan Aneka) -> perhitungan FIRE K69-5 apa adanya.
func (s *Service) isFlagUWAneka(ctx context.Context, id string) (bool, bool, error) {
	kb, ada := s.kasusBisnisDari(ctx, id)
	if !ada {
		return false, false, nil
	}
	if l := liniSpreadingDari(id, kb); !l.aneka && !l.golf {
		return false, false, nil
	}
	kasus, err := s.kasus.BacaKasus(ctx, id)
	if err != nil {
		return false, true, galatKasus(err)
	}
	return kodePolicyType(kasus.General.PolicyType) != policyTypeDeklarasi, true, nil
}

// decision4Policy - Decision4 FAC OUT? (`IsFacRetro`: varian kelas Work `pyWorkPage.OfferFacIn.IsFacRetro = 1`; varian
// Data-OfferFacIn ekspresi `.IsFacRetro = 1` - dasarJenjang.adaFacRetro) -> Decision3 IsInputFacRetro (selalu tidak di sistem
// baru) -> SubProcess3 FacRetro (tiket 65, 409); Else -> Utility1 selesai Policy (tiket 69a-2).
func decision4Policy(id string, d dasarJenjang) (tujuanUW, error) {
	if d.adaFacRetro() {
		log.Printf("nbfacin: Accept Policy kasus %s: Decision4 IsFacRetro - FacRetro belum tersedia", id)
		return tujuanUW{}, ErrFacRetroBelumTersedia
	}
	log.Printf("nbfacin: Accept Policy kasus %s: Decision4 Else - selesai Policy (SaveJsonPolicyFacIn_Act)", id)
	return tujuanUW{selesai: true}, nil
}

// acceptPolicy - tujuan Accept jalur Policy: posisi + NB status (flag tetap 1), selesai Policy (tiket 69a-2), pesan (ditahan
// 422: kurs tidak ada), atau 409 (FacRetro). TL: IsFlagUW 1 -> Utility8, 0 -> selesai. UW / JUW A: Utility8. Senior / DepHead / UW Financial:
// Decision4.
func (s *Service) acceptPolicy(ctx context.Context, id string, k models.KonteksSubmit, pelaku inti.Pelaku) (tujuanUW, error) {
	if k.Position == posisiTeamLeader {
		flag, err := s.isFlagUW(ctx, id)
		if err != nil {
			return tujuanUW{}, err
		}
		if !flag {
			log.Printf("nbfacin: Accept TL Policy kasus %s: IsFlagUW 0 (Decision21 Else) - selesai Policy", id)
			return tujuanUW{selesai: true}, nil
		}
	}
	d, err := s.bacaDasarJenjang(ctx, id)
	if err != nil {
		return tujuanUW{}, err
	}
	jabatan, jenjang := jabatanJenjangPolicy[k.Position]
	if !jenjang {
		return decision4Policy(id, d)
	}
	if d.pesanKurs != "" {
		return tujuanUW{pesan: d.pesanKurs}, nil
	}
	letter, alasan := s.letterJUWUW(k, pelaku, d, jabatan)
	if posisi, ok := tujuanAcceptUW[letter]; ok {
		return tujuanPosisi(posisi), nil
	}
	log.Printf("nbfacin: Accept Policy kasus %s: Decision33 LetterNo %q (%s) - Else -> Decision4", id, letter, alasan)
	return decision4Policy(id, d)
}

// periksaAcceptPolicy - proteksi Accept lalu tujuan. TL: Protection_Act + Deductible Acceptance (sama dengan TL tahap Offer);
// UW: CheckSpreadingProtect (sama dengan UW tahap Offer). Pesan proteksi lebih dulu (422); tanpa pesan -> tujuan (409 bila
// selesai Policy / FacRetro).
func (s *Service) periksaAcceptPolicy(ctx context.Context, id string, k models.KonteksSubmit, pelaku inti.Pelaku, tahap string,
	tl *isianDeductibleTL) ([]string, tujuanUW, error) {
	var pesan []string
	var err error
	if tahap == TahapTeamLeaderPolicy {
		pesan, err = s.ProteksiAccept(ctx, id, tl)
	} else {
		pesan, err = s.proteksiSpreadingUW(ctx, id)
	}
	if err != nil || len(pesan) > 0 {
		return pesan, tujuanUW{}, err
	}
	t, err := s.acceptPolicy(ctx, id, k, pelaku)
	if err != nil {
		return nil, tujuanUW{}, err
	}
	if t.pesan != "" {
		return []string{t.pesan}, t, nil
	}
	return []string{}, t, nil
}

// lengkapiTulisanPolicy - riwayat POSISI "Policy" PUTARAN "2" (SaveViewSuggest); TL: Reject / Ask -> Marketing tahap Binding flag 2
// (Transition69; ConfirmBinding TIDAK direset, K69-3); Accept -> posisi jenjang (flag tetap 1); TL tidak menulis DATE_TO_UW
// (SetValidateDate_PostAct 3.5 / Decision33 tanpa DateToUW). UW: Validity Of Acceptance + Comment (EmailSection_IsUW
// menampilkannya tanpa syarat flag) dan FACINPERFORMANCE POLICY (InsertFacInPerformance_Act langkah 5; ACCEPT `[dugaan kuat]`
// EmailTypeUW sisa Accept tahap Offer = "1", DECLINE EmailTypeUWPolicy 7; tanpa OUTSTANDING - langkah 13 dilewati flag 1).
func (s *Service) lengkapiTulisanPolicy(ctx context.Context, w *models.TulisanSubmit, tahap, id string, pelaku inti.Pelaku,
	isi models.IsiSubmit, k models.KonteksSubmit, tujuan tujuanUW) error {
	w.Riwayat.Posisi, w.Riwayat.Putaran = posisiCedingConfirmPolicy, putaranBinding
	if isi.EmailType == models.EmailTypeAccept {
		w.PosisiBaru, w.NBStatus = tujuan.posisi, tujuan.status
	}
	if tahap == TahapTeamLeaderPolicy {
		if isi.EmailType == models.EmailTypeReject || isi.EmailType == models.EmailTypeAsk {
			w.FlagBaru = flagBinding
		}
	} else {
		if err := s.lengkapiTulisanUW(ctx, w, id, pelaku, isi, k, tujuan); err != nil {
			return err
		}
		w.Riwayat.Posisi, w.Riwayat.Putaran = posisiCedingConfirmPolicy, putaranBinding
		w.Performance, w.PerformanceLanjut = performancePolicy(isi.EmailType, pelaku, k, id), nil
	}
	if isi.EmailType == models.EmailTypeAccept && tujuan.selesai {
		// Tiket 69a-2: Utility1 tanpa SET PositionNote / NBStatus -> POSITION dan NB_STATUS tetap; Resolved-Completed (K69-7).
		p, err := s.tulisanPolicy(ctx, id, k, pelaku)
		if err != nil {
			return err
		}
		w.PosisiBaru, w.NBStatus, w.NBStatusTetap, w.StatusWorkBaru, w.Policy = k.Position, "", true, statusSelesaiPolicy, p
	}
	return nil
}

// performancePolicy - InsertFacInPerformance_Act tahap Policy (POSITION_DOC "POLICY"): Accept ACCEPT, Decline DECLINE; kolom
// wajib kosong -> dilewati (langkah 12).
func performancePolicy(emailType string, pelaku inti.Pelaku, k models.KonteksSubmit, id string) *models.PerformanceFacIn {
	status := map[string]string{models.EmailTypeAccept: "ACCEPT", models.EmailTypeDecline: "DECLINE"}[emailType]
	if status == "" {
		return nil
	}
	f := models.PerformanceFacIn{PositionDoc: "POLICY", StatusDoc: status, Username: pelaku.AkunID, InsuredName: k.InsuredName,
		SOB: k.SobName, MOLeader: k.MOLeader, COB: k.ClassOfBusiness}
	if f.Username == "" || f.InsuredName == "" || f.SOB == "" || f.MOLeader == "" || f.COB == "" {
		log.Printf("nbfacin: Submit Policy kasus %s: FACINPERFORMANCE POLICY/%s dilewati (kolom wajib kosong)", id, status)
		return nil
	}
	return &f
}

// nbStatusPolicy - posisi / NB status jalur Policy selain Accept (Accept dari acceptPolicy).
func nbStatusPolicy(tahap string, isi models.IsiSubmit, k models.KonteksSubmit, akun models.AkunSubmit) (string, string) {
	a := awalanNBStatus()
	switch isi.EmailType {
	case models.EmailTypeAccept:
		return "", ""
	case models.EmailTypeReject, models.EmailTypeAsk:
		if tahap == TahapTeamLeaderPolicy {
			// Transition69: "NB IS IN " + MarketingName + "'S INBOX" (nama dari data, tanpa huruf besar).
			return repository.WorkbasketMarketing, a + " IS IN " + k.MarketingName + "'S INBOX"
		}
	}
	// Decline (K69-4 / A69-1): pola Decline TL / UW (SetValidateDateUW_PostAct 6.10, dua spasi).
	return posisiDecline, a + " WAS DECLINED BY  " + strings.ToUpper(akun.Nama)
}
