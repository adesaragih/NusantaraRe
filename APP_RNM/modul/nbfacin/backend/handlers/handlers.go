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

	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/services"
)

// DaftarkanRute memasang rute modul ini.
func DaftarkanRute(mux *http.ServeMux, svc *services.Service) {
	mux.HandleFunc("POST /api/nbfacin/premi", hitungPremi(svc))
	mux.HandleFunc("POST /api/nbfacin/akseptasi/langkah", langkahAkseptasi(svc))
	mux.HandleFunc("GET /api/nbfacin/account", cariAkun(svc))
	mux.HandleFunc("GET /api/nbfacin/class-of-business", kelasBisnis(svc))
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
	case errors.Is(err, services.ErrMasukanAkun), errors.Is(err, services.ErrMasukanKelasBisnis):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrTidakDapatDiproses):
		galat.Tulis(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, services.ErrTanpaDatabase), errors.Is(err, services.ErrTabelLimitTakTersedia),
		errors.Is(err, services.ErrAkunTanpaDatabase), errors.Is(err, services.ErrKelasBisnisTanpaDatabase):
		galat.Tulis(w, http.StatusServiceUnavailable, err.Error())
	default:
		log.Printf("nbfacin: galat server: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "galat server")
	}
}
