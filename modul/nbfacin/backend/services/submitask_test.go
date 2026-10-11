package services

import (
	"context"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestAcceptMarketingKembaliKePenanya - tiket 57 "IS ada Ask ?" (`DDL\Protection_Act.xml` langkah 21 + `DDL\InputInwardFacultativeOffer.xml`
// Decision43): IsFlagAsk = entri riwayat akseptasi TERAKHIR (siapa pun) ber-Approval Ask -> Accept marketing ke WORKBASKET catatan
// Pega penanya, NB status label tujuanBanding (SetBanding_ACT) / TL / statusInboxPosisi; selain itu Team Leader (jalur lama).
// Riwayat (produksi / Pega) sama dengan Accept biasa; nama workbasket lama tanpa B -> B (ADR-0032).
func TestAcceptMarketingKembaliKePenanya(t *testing.T) {
	ctx := context.Background()
	ask := []models.RiwayatAkseptasi{{Approval: "Ask"}, {Approval: "Accept"}}
	for nama, c := range map[string]struct {
		riwayat        []models.RiwayatAkseptasi
		pega           []string
		posisi, status string
	}{
		"Ask JUW_B":             {ask, []string{"ASK", "ReasFacInJuniorUnderwritingB"}, "ReasFacInJuniorUnderwritingB", "NB IS IN JUNIOR UNDERWRITER (B)'S INBOX"},
		"Ask TL":                {ask, []string{"ASK", "ReasFacInTeamLeader"}, "ReasFacInTeamLeader", "NB IS IN TEAM LEADER'S INBOX"},
		"Ask Senior UW":         {ask, []string{"ASK", "ReasFacInSeniorUnderwriting"}, "ReasFacInSeniorUnderwriting", "NB IS IN SENIORUW'S INBOX"},
		"Ask Manager Teknik":    {ask, []string{"ASK", "ReasFacInManagerTeknik"}, "ReasFacInManagerTeknik", "NB IS IN MANAGER TEKNIK'S INBOX"},
		"Ask nama lama tanpa B": {[]models.RiwayatAkseptasi{{Approval: " ask "}}, []string{"ask ", " ReasFacInJuniorUnderwriting"}, "ReasFacInJuniorUnderwritingB", "NB IS IN JUNIOR UNDERWRITER (B)'S INBOX"},
		// Ask JUW_B lalu Reject marketing: entri terakhir = Reject (marketing) -> flag mati -> TL.
		"Ask lalu Reject marketing":   {[]models.RiwayatAkseptasi{{Approval: "Reject"}, {Approval: "Ask"}}, []string{"ASK", "ReasFacInJuniorUnderwritingB"}, "ReasFacInTeamLeader", "NB IS IN TEAM LEADER'S INBOX"},
		"Reject UW":                   {[]models.RiwayatAkseptasi{{Approval: "Reject"}}, []string{"REJECT", "ReasFacInJuniorUnderwritingB"}, "ReasFacInTeamLeader", "NB IS IN TEAM LEADER'S INBOX"},
		"tanpa riwayat":               {nil, nil, "ReasFacInTeamLeader", "NB IS IN TEAM LEADER'S INBOX"},
		"Ask workbasket asing":        {ask, []string{"ASK", "ReasFacInEntahApa"}, "ReasFacInTeamLeader", "NB IS IN TEAM LEADER'S INBOX"},
		"Ask, catatan Pega bukan ASK": {ask, []string{"ACCEPT", "ReasFacInJuniorUnderwritingB"}, "ReasFacInTeamLeader", "NB IS IN TEAM LEADER'S INBOX"},
	} {
		st := submitUji()
		st.riwayat, st.pegaTerakhir = c.riwayat, c.pega
		h, err := layananSubmit(st).AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "ok"})
		if err != nil || h.Position != c.posisi || h.NBStatus != c.status || len(st.tulis) != 1 {
			t.Errorf("%s: %v %+v", nama, err, h)
			continue
		}
		w := st.tulis[0]
		if w.PosisiBaru != c.posisi || w.PosisiLama != "Offer" || w.Riwayat.Approval != "Accept" ||
			w.Pega != (models.RiwayatPega{Status: "ACCEPT", Username: "Uji Pengguna", Workbasket: "ReasFacInMarketing"}) {
			t.Errorf("%s tulisan: %+v", nama, w)
		}
	}
	// Ask hanya mengubah Accept: Reject marketing tetap ke Opportunity walau Ask terakhir ada.
	st := submitUji()
	st.riwayat, st.pegaTerakhir = ask, []string{"ASK", "ReasFacInJuniorUnderwritingB"}
	if h, err := layananSubmit(st).AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1", models.IsiSubmit{EmailType: "2", Comment: "tolak"}); err != nil ||
		h.Position != "Offer" {
		t.Errorf("Reject marketing: %v %+v", err, h)
	}
}
