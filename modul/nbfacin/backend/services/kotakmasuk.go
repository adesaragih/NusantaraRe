package services

// Kotak masuk Beranda NB Fac In (permintaan work owner 07-10-2026 "ini kenapa tidak ada ya, padahal ini login teamleader";
// pola nbtreatyin, keputusan 06-10-2026 "beranda menunjukkan berapa banyak case yang masuk di akun dia, mengikuti
// workbasket"). Antrean = workbasket tempat case NB Fac In dapat berada di sistem baru (antreanKasusNB) yang DIPEGANG akun
// (inti.Pelaku.Peran). Marketing dan Team Leader: hanya case se-team akun (keputusan work owner 08-10-2026 "yg muncul sesuai
// workbasketnya dan team groupnya" - penyimpangan sadar, Pega hanya per workbasket; akun tanpa team = antrean kosong + log);
// workbasket UW / lainnya: seluruh case ber-POSITION workbasket itu. Satu saringan untuk hitungan dan daftar.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// ErrKotakMasukTanpaDatabase - pembaca kotak masuk tidak terpasang. 503.
var ErrKotakMasukTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, kotak masuk tidak terbaca")

// antreanKasusNB - urutan tampil: marketing, Team Leader (Accept, tiket 57), lalu tujuan Banding SetBanding_ACT
// (tujuanBanding) dan jenjang lanjutan tiket 62 (UW Financial, Kadiv Financial, Manager Teknik); tanpa saringan team. Posisi
// lain (mis. "Decline") tidak menunggu siapa pun.
var antreanKasusNB = []string{repository.WorkbasketMarketing, posisiTeamLeader,
	"ReasFacInJuniorUnderwritingA", "ReasFacInJuniorUnderwritingB", "ReasFacInUnderwriting", "ReasFacInSeniorUnderwriting",
	"ReasFacInDepHeadUnderwriting", "ReasFacInUnderwritingFinancial", "ReasFacInFacultativeDivHead", "ReasFacInGroupLeader",
	"ReasFacInFinDivHead", "ReasFacInMarketingDirector", "ReasFacInTechnicalDirector", "ReasFacInManagerTeknik"}

// posisiMarketingKotakMasuk - nilai POSITION tahap marketing (posisiMarketing), dihitung ke workbasket ReasFacInMarketing.
var posisiMarketingKotakMasuk = []string{repository.PosisiAwalCaseNB, repository.WorkbasketMarketing}

// AntreanKotakMasuk - satu baris kotak masuk Beranda.
type AntreanKotakMasuk struct {
	Workbasket, Nama string
	Jumlah           int
}

// DenganKotakMasuk memasang pembaca kotak masuk.
func (s *Service) DenganKotakMasuk(p repository.PembacaKotakMasuk) *Service {
	s.kotakMasuk = p
	return s
}

// saringanKotakMasuk - workbasket `pilih` (yang dipegang pelaku) -> saringan repository: marketing (POSITION tahap marketing)
// dan Team Leader bersaring team akun (TeamKotakMasuk, dibaca hanya bila perlu); lainnya tanpa team. Akun tanpa team untuk
// antrean itu -> antrean kosong + log (keputusan work owner 08-10-2026 butir 3d).
func (s *Service) saringanKotakMasuk(ctx context.Context, pelaku inti.Pelaku, pilih []string) (repository.SaringanKotakMasuk, error) {
	var out repository.SaringanKotakMasuk
	var teamMkt, teamTL []string
	if dalam(pilih, repository.WorkbasketMarketing) || dalam(pilih, posisiTeamLeader) {
		var err error
		if teamMkt, teamTL, err = s.kotakMasuk.TeamKotakMasuk(ctx, pelaku.AkunID); err != nil {
			return out, err
		}
	}
	tim := func(w string, posisi, team []string, sumber string) {
		if len(team) == 0 {
			log.Printf("nbfacin: kotak masuk akun %s: tanpa team group (%s) - antrean %s kosong", pelaku.AkunID, sumber, w)
		}
		out.Tim = append(out.Tim, repository.AntreanTim{Posisi: posisi, Team: team})
	}
	for _, w := range pilih {
		switch w {
		case repository.WorkbasketMarketing:
			tim(w, posisiMarketingKotakMasuk, teamMkt, "MARKETINGOFFICER.AKSES_LOGIN")
		case posisiTeamLeader:
			tim(w, []string{posisiTeamLeader}, teamTL, "M_LIMIT_PROPERTYY.LOGIN")
		default:
			out.Lain = append(out.Lain, w)
		}
	}
	return out, nil
}

// dipegang - antreanKasusNB yang dipegang pelaku, urut antrean.
func dipegang(pelaku inti.Pelaku) []string {
	var out []string
	for _, w := range antreanKasusNB {
		if pelaku.PunyaPeran(w) {
			out = append(out, w)
		}
	}
	return out
}

// KotakMasuk - GET /api/nbfacin/kotak-masuk: satu baris per workbasket antrean yang dipegang (jumlah boleh 0); nama =
// M_WORKBASKET.NAME, kosong -> kode. Tanpa workbasket antrean -> [].
func (s *Service) KotakMasuk(ctx context.Context, pelaku inti.Pelaku) ([]AntreanKotakMasuk, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	if s.kotakMasuk == nil {
		return nil, ErrKotakMasukTanpaDatabase
	}
	pegang := dipegang(pelaku)
	out := []AntreanKotakMasuk{}
	if len(pegang) == 0 {
		return out, nil
	}
	saring, err := s.saringanKotakMasuk(ctx, pelaku, pegang)
	if err != nil {
		return nil, err
	}
	cacah, err := s.kotakMasuk.HitungKotakMasuk(ctx, saring)
	if err != nil {
		return nil, err
	}
	nama, err := s.kotakMasuk.NamaWorkbasket(ctx, pegang)
	if err != nil {
		return nil, err
	}
	for _, w := range pegang {
		n := cacah[w]
		if w == repository.WorkbasketMarketing {
			n = 0
			for _, p := range posisiMarketingKotakMasuk {
				n += cacah[p]
			}
		}
		label := strings.TrimSpace(nama[w])
		if label == "" {
			label = w
		}
		out = append(out, AntreanKotakMasuk{Workbasket: w, Nama: label, Jumlah: n})
	}
	return out, nil
}

// statusCedingConfirm - kolom inbox "Status" = padanan `.IsCedingConfirm` (keputusan work owner 07-10-2026 "setuju , buatin
// jangan halu"): `Activity\InputOfferFacInEngineer_preACT.xml` langkah 3 IsCedingConfirm "Offer" bila FlagOnGoingPolicy 0;
// konektor flow men-set "Binding" bersama FlagOnGoingPolicy 2 (Decision29 / Decision35). Selain itu "" - "Policy" (SetBanding
// jalur EDM) dan "Accepted" (hanya baris riwayat SetValidateDate 3.2) belum punya padanan tersimpan `[pertanyaan terbuka]`.
var statusCedingConfirm = map[string]string{"0": "Offer", "2": "Binding"}

// BarisInbox - satu baris inbox "New Business" (kontrak api.ts BarisKotakMasukNbFacIn).
type BarisInbox struct {
	CaseID, InsuredName, BusinessName, SourceOfBusiness, CedingCompany, InceptionDate, ShareRNM, Marketing, Status string
	CreatedDate, OfferingDate, ValidityDate, StatusInbox, AksiTerakhir, TglUpdate, SubmitToUW                      string
}

// keBarisInbox - format tanggal kabel DD-MM-YYYY, NB status kosong -> "NEW NB", aksi terakhir = Approval + " by " + PIC
// (`Activity\SetNBStatus_Act` langkah 6: StsApproval "Accept by " … + ViewSuggest(<LAST>).PICSuggest; label APPROVAL
// tersimpan dipakai apa adanya). validityDate = "": `.OfferFacIn.DateValidity` (GetDateValidity_ACT, dipanggil pemeriksa
// spreading tahap UW) tidak dihitung / disimpan sistem baru `belum terverifikasi`.
func keBarisInbox(b models.BarisKotakMasuk) BarisInbox {
	out := BarisInbox{CaseID: b.CaseID, InsuredName: b.InsuredName, BusinessName: b.BusinessName, SourceOfBusiness: b.SourceOfBusiness,
		CedingCompany: b.CedingCompany, InceptionDate: waktuKeKabel(b.StartDateTime), ShareRNM: b.ShareRNM, Marketing: b.Marketing,
		Status: statusCedingConfirm[strings.TrimSpace(b.FlagOnGoingPolicy)], CreatedDate: b.CreatedDate, OfferingDate: tanggalKeKabel(b.OfferingDate), StatusInbox: b.NBStatus,
		TglUpdate: b.TglUpdate, SubmitToUW: waktuUWKeKabel(b.SubmitToUW)}
	if strings.TrimSpace(out.StatusInbox) == "" {
		out.StatusInbox = nbStatusBaru
	}
	if b.ApprovalTerakhir != "" {
		out.AksiTerakhir = b.ApprovalTerakhir + " by " + b.PICAkhir
	}
	return out
}

// bentukWaktuInbox - kabel tglUpdate / submitToUw (sama dengan TO_CHAR 'YYYY-MM-DD HH24:MI:SS').
const bentukWaktuInbox = "2006-01-02 15:04:05"

// waktuUWKeKabel - kolom inbox UW "Time in Underwriting" (tiket 58): SUBMIT_TO_UW teks Pega DateTime UTC (ditulis Accept Team
// Leader, `IsJuniorUW_B`) -> YYYY-MM-DD HH:MM:SS WIB. Kosong / teks yang tidak cocok BentukWaktuPega -> "" (layar tidak
// dapat menghitung lama dari teks mentah).
func waktuUWKeKabel(s string) string {
	t, err := time.Parse(BentukWaktuPega, s)
	if err != nil {
		return ""
	}
	return t.In(WIB).Format(bentukWaktuInbox)
}

// DaftarMenunggu - GET /api/nbfacin/kotak-masuk/kasus?workbasket=: case yang menunggu pelaku di `workbasket` (kosong =
// semua yang dipegang); workbasket yang tidak dipegang / bukan antrean NB Fac In -> 403.
func (s *Service) DaftarMenunggu(ctx context.Context, pelaku inti.Pelaku, workbasket string) ([]BarisInbox, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	if s.kotakMasuk == nil {
		return nil, ErrKotakMasukTanpaDatabase
	}
	pilih := dipegang(pelaku)
	if workbasket = strings.TrimSpace(workbasket); workbasket != "" {
		if !dalam(pilih, workbasket) {
			return nil, fmt.Errorf("%w: kotak masuk %q bukan workbasket NB Fac In yang dipegang akun", inti.ErrTanpaWewenang, workbasket)
		}
		pilih = []string{workbasket}
	}
	out := []BarisInbox{}
	if len(pilih) == 0 {
		return out, nil
	}
	saring, err := s.saringanKotakMasuk(ctx, pelaku, pilih)
	if err != nil {
		return nil, err
	}
	baris, err := s.kotakMasuk.DaftarMenunggu(ctx, saring)
	if err != nil {
		return nil, err
	}
	for _, b := range baris {
		out = append(out, keBarisInbox(b))
	}
	return out, nil
}
