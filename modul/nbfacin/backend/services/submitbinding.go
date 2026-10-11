package services

// Form tahap MARKETING (BINDING) (tiket 64). `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\`:
//   - layar `Section\EmailSectionCeding.xml` (tampil `.IsCedingConfirm = 'Binding' || .FlagOnGoingPolicy=='2'`): Email Type
//     `.EmailTypeBinding` (`DDL\EmailTypeBinding.xml` 1 Accept / 4 Banding / 7 Decline / 9 Revise), tombol Binding (disabled bila
//     workbasket bukan posisi, `.OfferFacIn.ConfirmBinding = 1`, atau proteksi), checkbox Received RI Slip (disabled bila
//     ConfirmBinding != 1);
//   - tombol Binding `Activity\UpdateTglBindJsonOffer_Act.xml`: JSON_OFFER TGL_BIND = sysdate + jsondata, ConfirmBinding = 1,
//     riwayat "Offer has been bind.", FACINPERFORMANCE BindStatus "Binding" (D4: TIDAK di sini, sekali saat Submit Accept);
//   - Submit -> post-activity `SetValidateDate_PostAct` (riwayat PUTARAN "2", HISTORYAKSEPTASIPEGA, FACINPERFORMANCE) ->
//     Decision13: confirm -> SubProcess1 POLICY (`Flow\InputInwardFacultativeRISlip.xml` Start: FlagOnGoingPolicy 1,
//     IsCedingConfirm "Policy"; Decision27 IsTBonding -> Assignment1 UNDERWRITING FINANCIAL; Else -> Utility5 CekLimitSpreading
//     (D3: ditunda ke tiket Policy) -> Assignment6 TEAM LEADER); banding -> SetBanding_ACT (jenjang UW tahap Offer, flag 0);
//     revise -> Assignment12 MARKETING Offer (flag 0, `SetReviseProposal`: ConfirmBinding 0, ReceivedRiSlip false); decline
//     -> End3.
// Penyimpangan sadar (K64-2): status riwayat (POSISI) sesuai tahap tujuan, bukan "Accepted" untuk semua; NB status sesuai
// tujuan (Pega: Decline tidak terpasang); FACINPERFORMANCE BINDING/ACCEPT SEKALI (Pega: tombol + Submit). Layar Policy (TL / UW
// Financial flag 1) baca-saja sampai tiket jalur Policy.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// TahapBinding - tahap Submit MARKETING (BINDING): POSITION ReasFacInMarketing + FLAG_ONGOING_POLICY 2 (kontrak
// `TampilanSubmit.tahap`). Tab case TETAP baca-saja (tahapDari tidak berubah).
const TahapBinding = "binding"

const (
	// flagPolicy / flagOffer - FlagOnGoingPolicy tujuan: SubProcess1 POLICY Start "1"; SetBanding_ACT / Transition132 "0".
	flagPolicy = "1"
	flagOffer  = "0"
	// posisiCedingConfirmPolicy / posisiCedingConfirmBinding - `.IsCedingConfirm` tahap tujuan (HISTORYAKSEPTASIPRODUCTION.POSISI).
	posisiCedingConfirmPolicy  = "Policy"
	posisiCedingConfirmBinding = "Binding"
	// putaranBinding - SaveViewSuggest: PUTARAN "2" untuk Binding / Accepted / Policy.
	putaranBinding = "2"
	// keteranganTekanBinding - UpdateTglBindJsonOffer_Act langkah 4 CommentSuggest verbatim.
	keteranganTekanBinding = "Offer has been bind."
	// teamGroupBonding - When IsTBonding baris A `.OfferFacIn.QuotationData.TeamGroup = 5`.
	teamGroupBonding = "5"
	// pesanBelumBinding - Accept tanpa ConfirmBinding (tombol Send Email Pega disabled `.OfferFacIn.ConfirmBinding != 1`): teks
	// SISTEM BARU (bukan korpus).
	pesanBelumBinding = "Tekan tombol Binding terlebih dahulu sebelum Accept."
	// pesanCedingCoKosong / pesanSOBKosong - Protection_Act langkah 33 `Local.errmsgceding` "Ceding Co cannot be empty"
	// (verbatim); langkah 34 (SOB kosong) Pega memakai pesan yang SAMA - diganti "Source Of Business cannot be empty"
	// (keputusan work owner 08-10-2026, penyimpangan sadar).
	pesanCedingCoKosong = "Ceding Co cannot be empty"
	pesanSOBKosong      = "Source Of Business cannot be empty"
)

var (
	// ErrSudahBinding - tombol Binding ditekan lagi (ConfirmBinding sudah 1). 409.
	ErrSudahBinding = errors.New("services: penawaran sudah di-Binding")
	// ErrPenawaranBelumAda - JSON_OFFER case belum ada (Accept jalur akhir tiket 63 belum terjadi). 409.
	ErrPenawaranBelumAda = errors.New("services: penawaran (JSON_OFFER) case belum tersimpan")
)

// tahapSubmitBinding - Assignment9 MARKETING (BINDING): POSITION ReasFacInMarketing + flag 2.
func tahapSubmitBinding(t models.TahapKasus) bool {
	return t.Position == repository.WorkbasketMarketing && t.FlagOnGoingPolicy == flagBinding
}

// kasusBonding - When IsTBonding (`When\IsTBonding.xml`, "A OR B OR C OR D"): TeamGroup = 5 ATAU MarketingName salah satu
// literal nama marketing (NBFACIN_BONDING_MARKETINGNAME; literal TIDAK di kode). Nama dibandingkan tak peka huruf, dipangkas
// `[keputusan agent]` (Pega "=").
func (g KonfigGrup) kasusBonding(k models.KonteksSubmit) bool {
	if strings.TrimSpace(k.TeamGroup) == teamGroupBonding {
		return true
	}
	nama := strings.TrimSpace(k.MarketingName)
	for _, n := range g.BondingMarketingName {
		if nama != "" && strings.EqualFold(n, nama) {
			return true
		}
	}
	return false
}

// proteksiTahapBinding - Protection_Act di tahap Binding bukan Banding (tiket 64b): port tiket 57 + deductible 9.2.30.6.6 / 13
// (flag 2) + langkah 31 / 33 / 34 (CedingCoName / SourceOfBusiness kosong di flag 2 -> ProtectSpreading.CARI1; teks SOB =
// keputusan work owner D5). Pesan + flag; penahan dipilih pemanggil (Accept / tombol Binding).
// Langkah 25 CheckRISLIP (hanya PolicyStatus, tidak memblokir) dan survey report (langkah 32, data tidak ada) TIDAK diport (D5).
func (s *Service) proteksiTahapBinding(ctx context.Context, id string) (*pengumpulPesan, error) {
	p, err := s.proteksiBerflag(ctx, id, nil, true)
	if err != nil {
		return nil, err
	}
	if err := s.proteksiSpreadingBinding(ctx, id, p); err != nil {
		return nil, err
	}
	kasus, err := s.kasus.BacaKasus(ctx, id)
	if err != nil {
		return nil, galatKasus(err)
	}
	if strings.TrimSpace(kasus.General.CedingCoName) == "" {
		p.tambahFlag(pesanCedingCoKosong, flagSpreading1)
	}
	if strings.TrimSpace(kasus.General.SourceOfBusinessID) == "" {
		p.tambahFlag(pesanSOBKosong, flagSpreading1)
	}
	return p, nil
}

// proteksiSpreadingBinding - K73-29 (keputusan work owner 09-10-2026, SEMUA COB termasuk FIRE): Protection_Act langkah 22
// CheckSpreadingProtect_ACT berjalan penuh di tahap Binding (langkah 2 WHEN `pyWorkPage.IsCedingConfirm=="Binding"` T=5, L1026)
// - pesan yang sama dengan tahap uw (proteksiSpreadingUW), ber-flag ProtectSpreading.CARI1 (flag pesan CheckSpreadingProtectFire
// / AnekaGolf; untuk M4 / M5 `[dugaan]` flag sama).
func (s *Service) proteksiSpreadingBinding(ctx context.Context, id string, p *pengumpulPesan) error {
	pesan, err := s.proteksiSpreadingUW(ctx, id)
	if err != nil {
		return err
	}
	for _, m := range pesan {
		p.tambahFlag(m, flagSpreading1)
	}
	return nil
}

// periksaBinding - (pesan, pesanBinding) tahap Binding. Langkah 8 melewati Revise / Decline. Banding = proteksiBanding
// (EmailSectionCeding, tiket 57b). Accept: pesan = penahan tombol Send Email Accept (flagAcceptBinding) didahului
// pesanBelumBinding bila ConfirmBinding bukan "1"; pesanBinding = penahan tombol Binding (flagTombolBinding).
func (s *Service) periksaBinding(ctx context.Context, id string, k models.KonteksSubmit, emailType string) ([]string, []string, error) {
	switch emailType {
	case models.EmailTypeBanding:
		p, err := s.proteksiBerflag(ctx, id, nil, false)
		if err != nil {
			return nil, nil, err
		}
		if err := s.proteksiSpreadingBinding(ctx, id, p); err != nil { // K73-29
			return nil, nil, err
		}
		return p.penahan(flagBandingBinding), []string{}, nil
	case models.EmailTypeAccept:
	default:
		return []string{}, []string{}, nil
	}
	p, err := s.proteksiTahapBinding(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	pesan := p.penahan(flagAcceptBinding)
	if k.ConfirmBinding != "1" {
		pesan = append([]string{pesanBelumBinding}, pesan...)
	}
	return pesan, p.penahan(flagTombolBinding), nil
}

// lengkapiTulisanBinding - tulisan Submit tahap Binding per pilihan (D4 / D6 / K64-3).
func (s *Service) lengkapiTulisanBinding(w *models.TulisanSubmit, id string, pelaku inti.Pelaku, isi models.IsiSubmit,
	k models.KonteksSubmit) {
	a := awalanNBStatus()
	w.Riwayat.Putaran = putaranBinding
	switch isi.EmailType {
	case models.EmailTypeAccept:
		// SubProcess1 POLICY: Decision27 IsTBonding -> UW Financial; Else -> Team Leader (Transition2 NBStatus).
		w.PosisiBaru, w.NBStatus = posisiTeamLeader, a+" IS IN TEAM LEADER'S INBOX"
		if s.grup.kasusBonding(k) {
			w.PosisiBaru, w.NBStatus = posisiUWFinancial, a+" IS IN UNDERWRITER FINANCIAL'S INBOX"
		}
		w.FlagBaru, w.Riwayat.Posisi = flagPolicy, posisiCedingConfirmPolicy
		ri := teksBoolPega(isi.ReceivedRiSlip)
		w.StatusBinding = &models.StatusBinding{ConfirmBinding: "1", ReceivedRiSlip: ri}
		w.Performance = performanceBinding("ACCEPT", pelaku, k, ri, id)
	case models.EmailTypeBanding:
		// SetBanding_ACT langkah 3: FlagOnGoingPolicy "0", IsCedingConfirm "Offer" (posisi / NB status dari nbStatusSubmit).
		w.FlagBaru, w.Riwayat.Posisi = flagOffer, posisiCedingConfirmOffer
	case models.EmailTypeRevise:
		// Transition132 + SetReviseProposal (K64-3): Offer flag 0, ConfirmBinding 0, ReceivedRiSlip false; TGL_BIND dibiarkan.
		w.FlagBaru, w.Riwayat.Posisi = flagOffer, posisiCedingConfirmOffer
		w.StatusBinding = &models.StatusBinding{ConfirmBinding: "0", ReceivedRiSlip: "false"}
	case models.EmailTypeDecline:
		// End3: case berhenti di tahap Binding; FACINPERFORMANCE BINDING/DECLINE (post-activity langkah 10).
		w.Riwayat.Posisi = posisiCedingConfirmBinding
		w.Performance = performanceBinding("DECLINE", pelaku, k, k.ReceivedRiSlip, id)
	}
}

// performanceBinding - InsertFacInPerformance_Act tahap Binding: POSITIONDOC "BINDING" (FlagOnGoingPolicy 2), RISLIPRECEIVED
// = ReceivedRiSlip; dilewati bila kolom wajib kosong (langkah 12). POLICY/OUTSTANDING (langkah 13) TIDAK ditulis: USERNAME
// Pega = LimitAkseptasi.pxResults(1).CARI5 dan halaman itu tidak diisi di jalur Binding (GetLimitAkseptasi_Act 4-5 nonaktif)
// -> `[dugaan]` langkah 12 melewati insert (keputusan work owner D4 (i): ikut Pega).
func performanceBinding(status string, pelaku inti.Pelaku, k models.KonteksSubmit, ri, id string) *models.PerformanceFacIn {
	f := models.PerformanceFacIn{PositionDoc: "BINDING", StatusDoc: status, Username: pelaku.AkunID, InsuredName: k.InsuredName,
		SOB: k.SobName, MOLeader: k.MOLeader, COB: k.ClassOfBusiness, RISlipReceived: ri}
	if f.Username == "" || f.InsuredName == "" || f.SOB == "" || f.MOLeader == "" || f.COB == "" {
		log.Printf("nbfacin: Submit Binding kasus %s: FACINPERFORMANCE BINDING/%s dilewati (kolom wajib kosong)", id, status)
		return nil
	}
	return &f
}

// TekanBinding - POST …/submit/binding (tombol Binding). Urutan: 401 -> 503 -> 404 -> 403 workbasket -> 409 bukan tahap
// Binding -> 409 sudah Binding -> 422 proteksi Accept -> satu transaksi (409 posisi berubah / sudah Binding / JSON_OFFER belum
// ada) -> TampilanSubmit terbaru.
func (s *Service) TekanBinding(ctx context.Context, pelaku inti.Pelaku, id string) (TampilanSubmit, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TampilanSubmit{}, err
	}
	if err := s.siapSubmit(id); err != nil {
		return TampilanSubmit{}, err
	}
	if s.transaksi == nil {
		return TampilanSubmit{}, ErrSubmitTanpaDatabase
	}
	k, err := s.konteksSubmit(ctx, id)
	if err != nil {
		return TampilanSubmit{}, err
	}
	if _, err := alasanTidakBoleh(k, pelaku); err != nil {
		return TampilanSubmit{}, err
	}
	if tahapKonteks(k) != TahapBinding {
		return TampilanSubmit{}, fmt.Errorf("%w: tombol Binding hanya di tahap MARKETING (BINDING)", ErrSubmitBukanTahapMarketing)
	}
	if k.ConfirmBinding == "1" {
		return TampilanSubmit{}, ErrSudahBinding
	}
	p, err := s.proteksiTahapBinding(ctx, id)
	if err != nil {
		return TampilanSubmit{}, err
	}
	if pesan := p.penahan(flagTombolBinding); len(pesan) > 0 {
		return TampilanSubmit{}, GalatProteksi{Pesan: pesan}
	}
	akun, err := s.submit.AkunSubmit(ctx, pelaku.AkunID)
	if err != nil {
		return TampilanSubmit{}, err
	}
	nama := akun.Nama
	if nama == "" {
		nama = pelaku.AkunID
	}
	doc, err := json.Marshal(map[string]string{"IDNewBisnis": id})
	if err != nil {
		return TampilanSubmit{}, err
	}
	b := models.TekanBinding{CaseID: id, IDPega: idPega(id), Position: k.Position, Flag: k.FlagOnGoingPolicy, JSONData: string(doc),
		Riwayat: models.RiwayatProduksi{TypePolis: awalanNBStatus(), Posisi: posisiCedingConfirmBinding, PIC: nama,
			TglInp: s.sekarang().In(WIB).Format(bentukTglInp), Div: akun.Divisi, Type: "F", Putaran: putaranBinding,
			Approval: labelApproval[models.EmailTypeAccept].riwayat, Keterangan: keteranganTekanBinding, AksesLogin: pelaku.AkunID,
			PercentRNM: k.PercentShare}}
	err = s.transaksi(ctx, func(tx *db.Tx) error { return s.submit.TulisTekanBinding(ctx, tx, b) })
	switch {
	case errors.Is(err, repository.ErrPosisiBerubah):
		return TampilanSubmit{}, ErrSubmitBukanTahapMarketing
	case errors.Is(err, repository.ErrSudahBinding):
		return TampilanSubmit{}, ErrSudahBinding
	case errors.Is(err, repository.ErrPenawaranTidakAda):
		return TampilanSubmit{}, ErrPenawaranBelumAda
	case err != nil:
		return TampilanSubmit{}, galatKasus(err)
	}
	log.Printf("nbfacin: Binding kasus %s oleh %s", id, pelaku.AkunID)
	return s.TampilanSubmit(ctx, pelaku, id)
}
