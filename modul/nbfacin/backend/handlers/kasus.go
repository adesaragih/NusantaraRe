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
	// RiskScoring - tiket 54: `{finalScore, noteFinalScore}` tersimpan; "" bila belum ada.
	RiskScoring riskScoringKabel `json:"riskScoring"`
	// AksesLossRecord - tiket 55: akun login ber-JOB_POSITION "IT" (lunak: galat baca = false).
	AksesLossRecord bool `json:"aksesLossRecord"`
	// AksesKorespondensi - tiket 56: aturan sama dengan AksesLossRecord (satu pembacaan JOB_POSITION).
	AksesKorespondensi bool `json:"aksesKorespondensi"`
	// BolehUbah / AlasanBacaSaja - case baca-saja sesudah keluar dari tahap marketing (aturan bersama Submit, bacasaja.go).
	BolehUbah      bool   `json:"bolehUbah"`
	AlasanBacaSaja string `json:"alasanBacaSaja"`
	// BolehUbahSpreading / UbahBarisSpreading - tiket 59: tab Spreading dapat diubah (marketing / Team Leader seperti
	// bolehUbah, ATAU tahap uw = When IsUWSpread + workbasket posisi); grid baris spreading (Add / Delete / Treaty Type /
	// % Share, Copy Spreading) HANYA tahap uw.
	BolehUbahSpreading bool `json:"bolehUbahSpreading"`
	UbahBarisSpreading bool `json:"ubahBarisSpreading"`
	// UbahAkumulasi - tiket 60: Accumulation Code / Address coverage dapat diubah (PUT …/akumulasi) - tahap uw (When
	// IsUWSpread, K60-2) + workbasket posisi.
	UbahAkumulasi bool `json:"ubahAkumulasi"`
	// KunciPeriode - tiket 64d K1: Begin / End date terkunci (tahap Binding).
	KunciPeriode bool `json:"kunciPeriode"`
	// BolehFollowing - tiket 70 (A70-4): tombol Search Following Policy (boleh ubah + ConfirmBinding != 1).
	BolehFollowing bool `json:"bolehFollowing"`
	// BusinessType / Lini - tiket 73 P0: BusinessType_DeT atas BUSINESS GROUPPANEL + OLDID dan predikat When korpus
	// (gerbang tab menurut Class of Business, bukan groupBusiness).
	BusinessType string    `json:"businessType"`
	Lini         liniKabel `json:"lini"`
	// TampilanAneka - tiket 73 P1: kontainer form objek Aneka + sakelar (services.TampilanAneka).
	TampilanAneka tampilanAnekaKabel `json:"tampilanAneka"`
	// JumlahLokasi - tiket 75 F3: jumlah lokasi objek FIRE (layar memilih tab Object per halaman); null bila tak terbaca.
	JumlahLokasi *int `json:"jumlahLokasi"`
	// JenisCsv - tiket 73 P6: template CSV objek ('fire' | 'cargo' | 'aneka' | 'marinehull' | 'cpm' | ''), services.LiniKasus.
	JenisCsv string `json:"jenisCsv"`
	// MarketingBawaan - saran awal Marketing Name (ID MARKETINGOFFICER pembuat case) bila MOID belum tersimpan; "" selain itu
	// (keputusan work owner 12-10-2026).
	MarketingBawaan string `json:"marketingBawaan"`
}

// tampilanAnekaKabel - services.TampilanAneka; kontainer selalu larik ("C1".."C9", [] di luar gerbang tab Object Aneka).
type tampilanAnekaKabel struct {
	Kontainer        []string `json:"kontainer"`
	KolomTahunJumlah bool     `json:"kolomTahunJumlah"`
	HE               bool     `json:"he"`
	Golf             bool     `json:"golf"`
	TanamanHasil     bool     `json:"tanamanHasil"`
	Engineering      bool     `json:"engineering"`
	Bonding          bool     `json:"bonding"`
	Briguna          bool     `json:"briguna"`
	TopRisk          bool     `json:"topRisk"`
	SectionLol       bool     `json:"sectionLol"`
	JumlahSatuan     bool     `json:"jumlahSatuan"`
	// Tiket 73 P2 - form coverage Aneka.
	CurrencyCoverage bool `json:"currencyCoverage"`
	MBD              bool `json:"mbd"`
}

// liniKabel - hasil predikat When per lini (services.LiniKasus).
type liniKabel struct {
	Fire            bool `json:"fire"`
	Aneka           bool `json:"aneka"`
	AnekaNotBonding bool `json:"anekaNotBonding"`
	Liability       bool `json:"liability"`
	Bonding         bool `json:"bonding"`
	Golf            bool `json:"golf"`
	MarineCargo     bool `json:"marineCargo"` // tiket 74 P0
}

// denganAksesTab - aksesLossRecord / aksesKorespondensi dari satu pembacaan lunak JOB_POSITION; bolehUbah / alasanBacaSaja
// dari tahap kasus (lunak).
func denganAksesTab(out kasusKabel, svc *services.Service, r *http.Request, stubPelaku bool) kasusKabel {
	pelaku := inti.PelakuDari(r, stubPelaku)
	akses := svc.AksesTabIT(r.Context(), pelaku)
	out.AksesLossRecord, out.AksesKorespondensi = akses, akses
	out.BolehUbah, out.AlasanBacaSaja = svc.AksesUbahKasus(r.Context(), pelaku, r.PathValue("caseId"))
	out.BolehUbahSpreading, out.UbahBarisSpreading = svc.AksesSpreading(r.Context(), pelaku, r.PathValue("caseId"))
	out.UbahAkumulasi = svc.AksesAkumulasi(r.Context(), pelaku, r.PathValue("caseId"))
	out.KunciPeriode = svc.KunciPeriode(r.Context(), r.PathValue("caseId"))
	out.BolehFollowing = svc.AksesFollowing(r.Context(), pelaku, r.PathValue("caseId"))
	l := svc.LiniKasus(r.Context(), r.PathValue("caseId"))
	out.BusinessType = l.BusinessType
	out.Lini = liniKabel{Fire: l.Fire, Aneka: l.Aneka, AnekaNotBonding: l.AnekaNotBonding, Liability: l.Liability, Bonding: l.Bonding,
		Golf: l.Golf, MarineCargo: l.MarineCargo}
	a := l.TampilanAneka
	if a.Kontainer == nil {
		a.Kontainer = []string{}
	}
	out.TampilanAneka = tampilanAnekaKabel{Kontainer: a.Kontainer, KolomTahunJumlah: a.KolomTahunJumlah, HE: a.HE, Golf: a.Golf,
		TanamanHasil: a.TanamanHasil, Engineering: a.Engineering, Bonding: a.Bonding, Briguna: a.Briguna, TopRisk: a.TopRisk,
		SectionLol: a.SectionLol, JumlahSatuan: a.JumlahSatuan,
		CurrencyCoverage: a.CurrencyCoverage, MBD: a.MBD}
	out.JumlahLokasi = svc.JumlahLokasi(r.Context(), r.PathValue("caseId"))
	out.JenisCsv = l.JenisCsv
	return out
}

type riskScoringKabel struct {
	FinalScore     string `json:"finalScore"`
	NoteFinalScore string `json:"noteFinalScore"`
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
		RiskScoring: riskScoringKabel{teks(k.RiskScoring.FinalScore), k.RiskScoring.NoteFinalScore}, MarketingBawaan: k.MarketingBawaan,
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
func bacaKasus(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		k, err := svc.BacaKasus(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, denganAksesTab(keKasusKabel(k), svc, r, stubPelaku))
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
		galat.TulisJSON(w, denganAksesTab(keKasusKabel(k), svc, r, stubPelaku))
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
