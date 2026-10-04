package handlers

// Layar Inward Facultative tahap 2 (tiket 31): baca case, simpan General, pilihan
// Marketing Name. Kontrak kabel di bawah = jawaban ke frontend sesi 0f.

import (
	"encoding/json"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/services"
)

// generalKabel - objek `general`; empat medan terakhir tampil-saja.
type generalKabel struct {
	ReffNumber         string `json:"reffNumber"`
	QQName             string `json:"qqName"`
	BeginDate          string `json:"beginDate"`
	OfferingDate       string `json:"offeringDate"`
	EndDate            string `json:"endDate"`
	PolicyType         string `json:"policyType"`
	MarketingID        string `json:"marketingId"`
	Day                string `json:"day"`
	TypeFacultative    string `json:"typeFacultative"`
	SourceOfBusinessID string `json:"sourceOfBusinessId"`
	SourceOfBusiness   string `json:"sourceOfBusiness"`
	CedingCoName       string `json:"cedingCoName"`
	GroupName          string `json:"groupName"`
	OldPolicyNumber    string `json:"oldPolicyNumber"`
	// CedingList - daftar Ceding Co urut pilih (tiket 34); selalu larik.
	CedingList []barisCeding `json:"cedingList"`
}

// kasusKabel - jawaban GET /api/nbfacin/kasus/{caseId} dan PUT .../general.
type kasusKabel struct {
	CaseID      string           `json:"caseId"`
	Position    string           `json:"position"`
	StatusWork  string           `json:"statusWork"`
	Opportunity isianOpportunity `json:"opportunity"`
	InsuredName string           `json:"insuredName"`
	General     generalKabel     `json:"general"`
}

// isianGeneral - badan PUT .../general (sepuluh medan; tanpa medan tampil-saja).
type isianGeneral struct {
	ReffNumber      string `json:"reffNumber"`
	QQName          string `json:"qqName"`
	BeginDate       string `json:"beginDate"`
	OfferingDate    string `json:"offeringDate"`
	EndDate         string `json:"endDate"`
	PolicyType      string `json:"policyType"`
	MarketingID     string `json:"marketingId"`
	Day             string `json:"day"`
	TypeFacultative string `json:"typeFacultative"`
	// SourceOfBusinessID - kode SOB (tiket 33); nama diambil server dari AGENT.
	SourceOfBusinessID string `json:"sourceOfBusinessId"`
	// CedingIDs - kode Ceding Co urut pilih (tiket 34); kosong/tidak ada = daftar dikosongkan.
	CedingIDs []string `json:"cedingIds"`
}

// barisCeding - satu baris daftar Ceding Co.
type barisCeding struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// batasBadanGeneral - sepuluh medan, terlebar 500 bita.
const batasBadanGeneral = 16 << 10

func keKasusKabel(k services.KasusKabel) kasusKabel {
	o, g := k.Opportunity, k.General
	return kasusKabel{CaseID: k.CaseID, Position: k.Position, StatusWork: k.StatusWork, InsuredName: k.InsuredName,
		Opportunity: isianOpportunity{EstimatedClosingDate: o.EstimatedClosingDate, BusinessProspectName: o.BusinessProspectName,
			AccountID: o.AccountID, InsuredID: o.InsuredID, GroupBusinessID: o.GroupBusinessID, GroupBusiness: o.GroupBusiness,
			ClassOfBusiness: o.ClassOfBusiness, TypeOfInward: o.TypeOfInward, TypeOfFacultative: o.TypeOfFacultative,
			Phase: o.Phase, Stage: o.Stage, OpportunitySource: o.OpportunitySource, BusinessStatus: o.BusinessStatus,
			Description: o.Description},
		General: generalKabel{ReffNumber: g.ReffNumber, QQName: g.QQName, BeginDate: g.BeginDate,
			OfferingDate: g.OfferingDate, EndDate: g.EndDate, PolicyType: g.PolicyType, MarketingID: g.MarketingID,
			Day: g.Day, TypeFacultative: g.TypeFacultative, SourceOfBusinessID: g.SourceOfBusinessID,
			SourceOfBusiness: g.SourceOfBusiness,
			CedingCoName:     g.CedingCoName, GroupName: g.GroupName, OldPolicyNumber: g.OldPolicyNumber, CedingList: keBarisCeding(g.CedingList)}}
}

// bacaKasus - GET /api/nbfacin/kasus/{caseId}.
func bacaKasus(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		k, err := svc.BacaKasus(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keKasusKabel(k))
	}
}

// simpanGeneral - PUT /api/nbfacin/kasus/{caseId}/general (tombol "Save for later").
func simpanGeneral(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b isianGeneral
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanGeneral)).Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON isian General yang sah")
			return
		}
		k, err := svc.SimpanGeneral(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), services.IsianGeneral{
			ReffNumber: b.ReffNumber, QQName: b.QQName, BeginDate: b.BeginDate, OfferingDate: b.OfferingDate,
			EndDate: b.EndDate, PolicyType: b.PolicyType, MarketingID: b.MarketingID, Day: b.Day,
			TypeFacultative: b.TypeFacultative, SourceOfBusinessID: b.SourceOfBusinessID, CedingIDs: b.CedingIDs,
		})
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keKasusKabel(k))
	}
}

// barisMarketing - satu pilihan Marketing Name: id disimpan ke MOID, nama tampil.
type barisMarketing struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}

// daftarMarketing - GET /api/nbfacin/marketing-officer.
func daftarMarketing(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hasil, err := svc.DaftarMarketingOfficer(r.Context())
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]barisMarketing, 0, len(hasil))
		for _, m := range hasil {
			baris = append(baris, barisMarketing{ID: m.ID, Nama: m.Nama})
		}
		galat.TulisJSON(w, struct {
			Baris []barisMarketing `json:"baris"`
		}{baris})
	}
}

func keBarisCeding(d []services.CedingKabel) []barisCeding {
	hasil := make([]barisCeding, 0, len(d))
	for _, c := range d {
		hasil = append(hasil, barisCeding{ID: c.ID, Name: c.Name})
	}
	return hasil
}
