package services

// "IS ada Ask ?" (tiket 57, keputusan work owner 11-10-2026). Sumber = rule Pega LIVE kiriman work owner di `D:\migrasi\RNM\DDL\`
// (BUKAN korpus - flow korpus 01-01-94 / 20260512 tidak memuatnya) `[terverifikasi]`:
//   - `DDL\InputInwardFacultativeOffer.xml` Decision43 "IS ada Ask ?" (sesudah Utility9 SAVE JSON_OFFER): When IsFlagAsk ->
//     Utility3 SetBanding_ACT (Transition173); Else -> Decision35 (jalur lama). When IsFlagAsk (layar, RS 01-01-95):
//     `pyWorkPage.OfferFacIn.IsFlagAsk = "true"`.
//   - `DDL\Protection_Act.xml` (01-01-96): L1651 IsFlagAsk = "" lalu langkah 21 "Set FlagASk JIka aksep terakhir ask" IsFlagAsk =
//     "true" berprakondisi `ViewSuggest(<LAST>).Approval==3` - entri TERAKHIR riwayat akseptasi, SIAPA PUN pembuatnya, ber-Approval
//     Ask. Sistem baru: baris HISTORYAKSEPTASIPRODUCTION terbaru (RiwayatAkseptasi, SaveViewSuggest yang sama) ber-APPROVAL "Ask"
//     (labelApproval EmailTypeAsk "3").
//   - Tujuan: Pega SetBanding_ACT (LetterNo -> PositionNote + NB status); LetterNo tidak diisi saat Ask (kode jabatan approver yang
//     terakhir dirutekan) -> tujuan `[keputusan agent]` = posisi penanya ("balik ke siapa yg ask terakhir"): WORKBASKET catatan
//     HISTORYAKSEPTASIPEGA terbaru non-marketing (wajib berSTATUS ASK; tidak -> jalur lama + log). NB status = label tujuanBanding
//     (SetBanding_ACT) bila posisinya ada di peta itu; Team Leader / posisi uw lain di luar peta: teks jalur yang sudah ada.
// Proteksi Accept, riwayat, dan kunci posisi sama dengan Accept biasa - hanya tujuan + NB status yang berbeda.

import (
	"context"
	"log"
	"strings"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// posisiJuniorUWLama - workbasket JUW sebelum ADR-0032 (catatan Pega lama) -> posisiJuniorUWB (migrasi 211, issues/58).
const posisiJuniorUWLama = "ReasFacInJuniorUnderwriting"

// statusInboxAsk - NB status (sesudah awalan "NB") tujuan kembali-ke-penanya: label tujuanBanding (SetBanding_ACT) bila posisinya
// ada di peta itu; Team Leader = teks Accept marketing; posisi uw lain = statusInboxPosisi (tanpa nama orang).
func statusInboxAsk(posisi string) (string, bool) {
	for _, t := range tujuanBanding {
		if t.posisi == posisi {
			return " IS IN " + t.label + "'S INBOX", true
		}
	}
	if posisi == posisiTeamLeader {
		return " IS IN TEAM LEADER'S INBOX", true
	}
	s, ada := statusInboxPosisi[posisi]
	return s, ada
}

// tujuanAskMarketing - IsFlagAsk (entri riwayat terakhir = Ask) -> posisi / NB status kembali-ke-penanya; ok = false -> jalur lama
// (Team Leader). Catatan Pega penanya tidak ber-ASK atau workbasket di luar posisi yang dikenal -> jalur lama + log `[keputusan
// agent]` (tidak menjatuhkan case ke antrean tanpa pemilik).
func (s *Service) tujuanAskMarketing(ctx context.Context, id string) (posisi, status string, ok bool, err error) {
	riwayat, err := s.submit.RiwayatAkseptasi(ctx, idPega(id))
	if err != nil || len(riwayat) == 0 ||
		!strings.EqualFold(strings.TrimSpace(riwayat[0].Approval), labelApproval[models.EmailTypeAsk].riwayat) {
		return "", "", false, err
	}
	st, wb, ada, err := s.submit.CatatanPegaNonMarketing(ctx, idPega(id))
	if err != nil {
		return "", "", false, err
	}
	if !ada || !strings.EqualFold(strings.TrimSpace(st), labelApproval[models.EmailTypeAsk].pega) {
		log.Printf("nbfacin: Submit kasus %s: riwayat terakhir Ask tetapi catatan Pega penanya tidak ber-ASK (%q) - jalur lama ke Team Leader", id, st)
		return "", "", false, nil
	}
	posisi = strings.TrimSpace(wb)
	if posisi == posisiJuniorUWLama {
		posisi = posisiJuniorUWB // ADR-0032
	}
	teks, dikenal := statusInboxAsk(posisi)
	if !dikenal || posisi == repository.WorkbasketMarketing {
		log.Printf("nbfacin: Submit kasus %s: Ask terakhir dari workbasket %q tidak dikenal - jalur lama ke Team Leader", id, wb)
		return "", "", false, nil
	}
	return posisi, awalanNBStatus() + teks, true, nil
}
