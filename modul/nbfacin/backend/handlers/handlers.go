// Package handlers memuat rute HTTP modul NB Fac In (tiket 20). Handler hanya
// mengurai permintaan, memanggil services, dan memetakan galat ke kode HTTP; ia
// tidak mengimpor repository (TestHandlersTidakMengimporRepository).
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/services"
)

// DaftarkanRute memasang rute modul ini. stubPelaku = config AuthStub (penunda
// identitas X-Pelaku, inti/backend/pelaku_http.go); sesi login selalu didahulukan.
func DaftarkanRute(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	mux.HandleFunc("POST /api/nbfacin/premi", hitungPremi(svc))
	mux.HandleFunc("POST /api/nbfacin/akseptasi/langkah", langkahAkseptasi(svc))
	mux.HandleFunc("GET /api/nbfacin/account", cariAkun(svc))
	mux.HandleFunc("GET /api/nbfacin/class-of-business", kelasBisnis(svc))
	mux.HandleFunc("POST /api/nbfacin/opportunity", buatOpportunity(svc, stubPelaku))
	mux.HandleFunc("GET /api/nbfacin/kasus/{caseId}", bacaKasus(svc))
	mux.HandleFunc("PUT /api/nbfacin/kasus/{caseId}/general", simpanGeneral(svc, stubPelaku))
	mux.HandleFunc("GET /api/nbfacin/marketing-officer", daftarMarketing(svc))
	mux.HandleFunc("GET /api/nbfacin/opportunity", cariPortal(svc, stubPelaku))
	mux.HandleFunc("GET /api/nbfacin/sob", cariSOB(svc))
	mux.HandleFunc("GET /api/nbfacin/kasus/{caseId}/objek", bacaObjek(svc))
	mux.HandleFunc("PUT /api/nbfacin/kasus/{caseId}/objek", simpanObjek(svc, stubPelaku))
	mux.HandleFunc("GET /api/nbfacin/risk-address", cariRisk(svc))
	mux.HandleFunc("POST /api/nbfacin/risk-address", simpanAlamat(svc, stubPelaku))
	mux.HandleFunc("GET /api/nbfacin/rw", cariRW(svc))
	mux.HandleFunc("GET /api/nbfacin/occupation", cariOccupation(svc))
	mux.HandleFunc("GET /api/nbfacin/jenis-item-objek", daftarJenisItem(svc))
	mux.HandleFunc("GET /api/nbfacin/mata-uang", daftarMataUang(svc))
	mux.HandleFunc("GET /api/nbfacin/kasus/{caseId}/table-of-limit", cariTableOfLimit(svc))
	mux.HandleFunc("GET /api/nbfacin/coverage", cariCoverage(svc))
	mux.HandleFunc("GET /api/nbfacin/coverage-otomatis", coverageOtomatis(svc))
	mux.HandleFunc("POST /api/nbfacin/kasus/{caseId}/hitung-coverage", hitungCoverage(svc))
	mux.HandleFunc("POST /api/nbfacin/kasus/{caseId}/hitung-net-rate", hitungNetRate(svc))
	mux.HandleFunc("GET /api/nbfacin/akumulasi", cariAkumulasi(svc))
	mux.HandleFunc("GET /api/nbfacin/akumulasi/saran/{jenis}", saranAkumulasi(svc))
	mux.HandleFunc("POST /api/nbfacin/akumulasi", tambahAkumulasi(svc, stubPelaku))
	mux.HandleFunc("GET /api/nbfacin/akumulasi/czone", czoneZip(svc))
	mux.HandleFunc("GET /api/nbfacin/akumulasi/zipcode", cariZipAkumulasi(svc))
	mux.HandleFunc("GET /api/nbfacin/kasus/{caseId}/klausa", bacaKlausa(svc))
	mux.HandleFunc("PUT /api/nbfacin/kasus/{caseId}/klausa", simpanKlausa(svc, stubPelaku))
	mux.HandleFunc("GET /api/nbfacin/klausa", cariKlausa(svc))
	mux.HandleFunc("GET /api/nbfacin/klausa/{id}/argumen", argumenKlausa(svc))
	mux.HandleFunc("GET /api/nbfacin/kasus/{caseId}/spreading", bacaSpreading(svc))
	mux.HandleFunc("PUT /api/nbfacin/kasus/{caseId}/spreading", simpanSpreading(svc, stubPelaku))
	mux.HandleFunc("POST /api/nbfacin/kasus/{caseId}/spreading/hitung-share", hitungShareSpreading(svc))
	mux.HandleFunc("POST /api/nbfacin/kasus/{caseId}/spreading/salin", salinSpreading(svc))
	mux.HandleFunc("GET /api/nbfacin/kasus/{caseId}/cedant", bacaCedant(svc))
	mux.HandleFunc("PUT /api/nbfacin/kasus/{caseId}/cedant", simpanCedant(svc, stubPelaku))
}

// isianOpportunity - badan POST /api/nbfacin/opportunity, kontrak frontend
// `IsianOpportunity` (modul/nbfacin/frontend/api.ts); semua teks, tanggal DD-MM-YYYY.
type isianOpportunity struct {
	EstimatedClosingDate string `json:"estimatedClosingDate"`
	BusinessProspectName string `json:"businessProspectName"`
	AccountID            string `json:"accountId"`
	InsuredID            string `json:"insuredId"`
	GroupBusinessID      string `json:"groupBusinessId"`
	GroupBusiness        string `json:"groupBusiness"`
	ClassOfBusiness      string `json:"classOfBusiness"`
	TypeOfInward         string `json:"typeOfInward"`
	TypeOfFacultative    string `json:"typeOfFacultative"`
	Phase                string `json:"phase"`
	Stage                string `json:"stage"`
	OpportunitySource    string `json:"opportunitySource"`
	BusinessStatus       string `json:"businessStatus"`
	Description          string `json:"description"`
}

// batasBadanOpportunity - badan terbesar yang diterima: 14 medan, terlebar 4000 bita.
const batasBadanOpportunity = 64 << 10

// buatOpportunity - POST /api/nbfacin/opportunity (tiket 29): 201 {"caseId":"NB-<n>"}.
func buatOpportunity(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b isianOpportunity
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanOpportunity)).Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON isian opportunity yang sah")
			return
		}
		id, err := svc.BuatOpportunity(r.Context(), inti.PelakuDari(r, stubPelaku), services.IsianOpportunity{
			EstimatedClosingDate: b.EstimatedClosingDate, BusinessProspectName: b.BusinessProspectName,
			AccountID: b.AccountID, InsuredID: b.InsuredID, GroupBusinessID: b.GroupBusinessID,
			GroupBusiness: b.GroupBusiness, ClassOfBusiness: b.ClassOfBusiness, TypeOfInward: b.TypeOfInward,
			TypeOfFacultative: b.TypeOfFacultative, Phase: b.Phase, Stage: b.Stage,
			OpportunitySource: b.OpportunitySource, BusinessStatus: b.BusinessStatus, Description: b.Description,
		})
		if err != nil {
			tulisGalat(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(struct {
			CaseID string `json:"caseId"`
		}{id})
	}
}

// barisKelasBisnis - satu pilihan Class Of Business (tiket 28); teks apa adanya, NULL = "".
type barisKelasBisnis struct {
	ID   string `json:"id"`
	Note string `json:"note"`
}

// kelasBisnis - GET /api/nbfacin/class-of-business?groupBusinessId= (tiket 28): semua
// pilihan, tanpa paging.
func kelasBisnis(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hasil, err := svc.KelasBisnis(r.Context(), r.URL.Query().Get("groupBusinessId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]barisKelasBisnis, 0, len(hasil))
		for _, k := range hasil {
			baris = append(baris, barisKelasBisnis{ID: k.ID, Note: k.Note})
		}
		galat.TulisJSON(w, struct {
			Baris []barisKelasBisnis `json:"baris"`
		}{baris})
	}
}

// barisAkun - satu baris jawaban lookup akun (tiket 27); teks apa adanya, NULL = "".
type barisAkun struct {
	ID              string `json:"id"`
	InsuredID       string `json:"insuredId"`
	InsuredName     string `json:"insuredName"`
	GroupBusinessID string `json:"groupBusinessId"`
	GroupBusiness   string `json:"groupBusiness"`
}

// cariAkun - GET /api/nbfacin/account?cari=&halaman= (tiket 27, popup ChooseAccount).
func cariAkun(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		halaman := 1
		if h := r.URL.Query().Get("halaman"); h != "" {
			n, err := strconv.Atoi(h)
			if err != nil {
				galat.Tulis(w, http.StatusBadRequest, services.ErrMasukanAkun.Error())
				return
			}
			halaman = n
		}
		hasil, err := svc.CariAkun(r.Context(), r.URL.Query().Get("cari"), halaman)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]barisAkun, 0, len(hasil.Baris))
		for _, a := range hasil.Baris {
			baris = append(baris, barisAkun{ID: a.ID, InsuredID: a.InsuredID, InsuredName: a.InsuredName,
				GroupBusinessID: a.GroupBusinessID, GroupBusiness: a.GroupBusiness})
		}
		galat.TulisJSON(w, struct {
			Baris   []barisAkun `json:"baris"`
			Total   int         `json:"total"`
			Halaman int         `json:"halaman"`
			Ukuran  int         `json:"ukuran"`
		}{baris, hasil.Total, hasil.Halaman, hasil.Ukuran})
	}
}

// permintaanPremi - masukan satu coverage; angka berupa TEKS desimal bertitik
// (sama dengan premium.Input), tidak pernah bilangan JSON - tanpa float.
type permintaanPremi struct {
	LiniBisnis             string `json:"liniBisnis"`
	CalculateMethod        string `json:"calculateMethod"`
	MataUang               string `json:"mataUang"`
	TSI                    string `json:"tsi"`
	Rate                   string `json:"rate"`
	ProRatePercent         string `json:"proRatePercent"`
	PctShortPeriod         string `json:"pctShortPeriod"`
	Discount               string `json:"discount"`
	DiscountType           string `json:"discountType"`
	DiscountPercentage     string `json:"discountPercentage"`
	PremiSebelumnya        string `json:"premiSebelumnya"`
	Loading                string `json:"loading"`
	ProRatePercentCoverage string `json:"proRatePercentCoverage"`
	IndemnityPercentage    string `json:"indemnityPercentage"`
	LossLimit              string `json:"lossLimit"`
	PctAdjustment          string `json:"pctAdjustment"`
	CoverageBasis          string `json:"coverageBasis"`
	NetRate                string `json:"netRate"`
	FirstScale             string `json:"firstScale"`
	MBD                    bool   `json:"mbd"`
	MasterPolicy           bool   `json:"masterPolicy"`
}

func (p permintaanPremi) masukan() kontrak.MasukanPremiFacIn {
	return kontrak.MasukanPremiFacIn{LiniBisnis: p.LiniBisnis, CalculateMethod: p.CalculateMethod, MataUang: p.MataUang,
		TSI: p.TSI, Rate: p.Rate, ProRatePercent: p.ProRatePercent, PctShortPeriod: p.PctShortPeriod, Discount: p.Discount,
		DiscountType: p.DiscountType, DiscountPercentage: p.DiscountPercentage, PremiSebelumnya: p.PremiSebelumnya,
		Loading: p.Loading, ProRatePercentCoverage: p.ProRatePercentCoverage, IndemnityPercentage: p.IndemnityPercentage,
		LossLimit: p.LossLimit, PctAdjustment: p.PctAdjustment, CoverageBasis: p.CoverageBasis, NetRate: p.NetRate,
		FirstScale: p.FirstScale, MBD: p.MBD, MasterPolicy: p.MasterPolicy}
}

func hitungPremi(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p permintaanPremi
		if !urai(w, r, &p) {
			return
		}
		h, err := svc.HitungPremi(p.masukan())
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, struct {
			Premi     string `json:"premi"`
			MataUang  string `json:"mataUang"`
			AsalRumus string `json:"asalRumus"`
		}{utils.FormatDecimal(h.Premi.Amount), h.Premi.Currency, h.AsalRumus})
	}
}

// permintaanAkseptasi - halaman kerja kasus sebagai jalur properti Pega → nilai,
// dan pengguna yang baru memutuskan.
//
// ⚠️ Jabatan DARI ISIAN PERMINTAAN - keputusan work owner 02-10-2026, butir 58.
// TIDAK AMAN UNTUK PRODUKSI: klien dapat mengaku jabatan mana pun. Diganti saat
// model peran (OQ RBAC) diputuskan; jawabannya membawa peringatan ini.
type permintaanAkseptasi struct {
	Kasus       map[string]string `json:"kasus"`
	Jabatan     string            `json:"jabatan"`
	AnggotaGrup bool              `json:"anggotaGrup"`
}

const peringatanJabatan = "jabatan pengguna diambil dari isian permintaan (butir 58) - tidak aman untuk produksi"

type kasusPeta map[string]string

func (k kasusPeta) Nilai(jalur string) (string, bool) { v, ada := k[jalur]; return v, ada }

func langkahAkseptasi(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p permintaanAkseptasi
		if !urai(w, r, &p) {
			return
		}
		if p.Jabatan == "" {
			galat.Tulis(w, http.StatusBadRequest, "jabatan pengguna wajib diisi")
			return
		}
		tr, err := svc.LangkahAkseptasi(r.Context(), kasusPeta(p.Kasus),
			kontrak.PenggunaFacIn{Jabatan: kontrak.JabatanFacIn(p.Jabatan), AnggotaGrup: p.AnggotaGrup})
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, struct {
			Selesai             bool   `json:"selesai"`
			JabatanTujuan       string `json:"jabatanTujuan"`
			Antrean             string `json:"antrean"`
			PositionNoteDitulis bool   `json:"positionNoteDitulis"`
			Peringatan          string `json:"peringatan"`
		}{tr.Selesai, string(tr.JabatanTujuan), string(tr.Antrean), tr.PositionNoteDitulis, peringatanJabatan})
	}
}

// batasBadan - badan permintaan terbesar yang dibaca (seperti treatycontractout).
const batasBadan = 1 << 16

// urai - badan JSON ke `tujuan`; medan tak dikenal ditolak (salah eja tidak boleh
// menjadi medan kosong diam-diam). Pesan ke klien tetap (seperti claimlife); rincian
// urai tidak dikembalikan.
func urai(w http.ResponseWriter, r *http.Request, tujuan any) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadan))
	dec.DisallowUnknownFields()
	if err := dec.Decode(tujuan); err != nil {
		galat.Tulis(w, http.StatusBadRequest, "permintaan tidak terurai")
		return false
	}
	return true
}

func tulisGalat(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrMasukanAkun), errors.Is(err, services.ErrMasukanKelasBisnis),
		errors.Is(err, services.ErrMasukanOpportunity), errors.Is(err, services.ErrMasukanGeneral),
		errors.Is(err, services.ErrMasukanPortal), errors.Is(err, services.ErrMasukanSOB),
		errors.Is(err, services.ErrMasukanObjek), errors.Is(err, services.ErrMasukanRisk),
		errors.Is(err, services.ErrMasukanRW), errors.Is(err, services.ErrMasukanAlamat),
		errors.Is(err, services.ErrMasukanOccupation), errors.Is(err, services.ErrMasukanTableOfLimit),
		errors.Is(err, services.ErrMasukanCoverage), errors.Is(err, services.ErrMasukanAkumulasi),
		errors.Is(err, services.ErrWajibAkumulasi), errors.Is(err, services.ErrMasukanKlausa),
		errors.Is(err, services.ErrMasukanSpreading), errors.Is(err, services.ErrMasukanCedant):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, services.ErrKasusTidakAda):
		galat.Tulis(w, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrTableOfLimitTidakSiap), errors.Is(err, services.ErrPeriodeKasus),
		errors.Is(err, services.ErrAkumulasiSudahAda), errors.Is(err, services.ErrSpreadingBerubah):
		galat.Tulis(w, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrTidakDapatDiproses):
		galat.Tulis(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, services.ErrTanpaDatabase), errors.Is(err, services.ErrTabelLimitTakTersedia),
		errors.Is(err, services.ErrAkunTanpaDatabase), errors.Is(err, services.ErrKelasBisnisTanpaDatabase),
		errors.Is(err, services.ErrOpportunityTanpaDatabase), errors.Is(err, services.ErrKasusTanpaDatabase),
		errors.Is(err, services.ErrMarketingTanpaDatabase), errors.Is(err, services.ErrPortalTanpaDatabase),
		errors.Is(err, services.ErrSOBTanpaDatabase), errors.Is(err, services.ErrObjekTanpaDatabase),
		errors.Is(err, services.ErrRiskTanpaDatabase), errors.Is(err, services.ErrRWTanpaDatabase),
		errors.Is(err, services.ErrOccupationTanpaDatabase), errors.Is(err, services.ErrPilihanItemTanpaDatabase),
		errors.Is(err, services.ErrTableOfLimitTanpaDatabase), errors.Is(err, services.ErrCoverageTanpaDatabase),
		errors.Is(err, services.ErrAkumulasiTanpaDatabase), errors.Is(err, services.ErrKlausaTanpaDatabase),
		errors.Is(err, services.ErrSpreadingTanpaDatabase), errors.Is(err, services.ErrCedantTanpaDatabase):
		galat.Tulis(w, http.StatusServiceUnavailable, err.Error())
	default:
		log.Printf("nbfacin: galat server: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "galat server")
	}
}
