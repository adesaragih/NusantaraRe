package handlers

// Pintu HTTP simpan utuh satu kontrak - tiket 09.
//
//	POST /api/treaty-contract-out/tahun/{id}/kontrak-utuh          kontrak BARU + reinsurer/security/business/klausul
//	PUT  /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/utuh    kontrak yang ADA + seluruh perubahannya
//
// Satu permintaan = satu transaksi (services.SimpanUtuhTCO). Kode HTTP
// mengikuti galat baris yang gagal (409, 422, ...); PESANNYA menyebut baris
// itu (AC 38) - kecuali galat server, yang isinya tidak dibocorkan.
//
// Dibaca sesudah: services/tco_simpan_utuh.go.

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"nusantarare/internal/services"
)

func layananSimpanUtuhTCO(svc *services.Service) *services.SimpanUtuhTCO {
	return svc.SimpanUtuhTCO().
		DenganKontrak(layananKontrakTCO(svc)).
		DenganReinsurer(layananReinsurerTCO(svc)).
		DenganSecurity(layananSecurityTCO(svc)).
		DenganBusiness(layananBusinessTCO(svc)).
		DenganKlausul(layananKlausulTCO(svc)).
		DenganPenetapIdentitas(services.PenetapIdentitasOracle(svc))
}

func daftarkanRuteSimpanUtuhTCO(mux *http.ServeMux, svc *services.Service, stub bool) {
	mux.HandleFunc("POST /api/treaty-contract-out/tahun/{id}/kontrak-utuh", simpanUtuhTCO(svc, stub, false))
	mux.HandleFunc("PUT /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/utuh", simpanUtuhTCO(svc, stub, true))
}

// batasBadanUtuhTCO - badan simpan utuh lebih besar dari satu baris (<= 500 baris).
const batasBadanUtuhTCO = 4 << 20

func simpanUtuhTCO(svc *services.Service, stub, perbarui bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		var masuk services.KontrakUtuhMasuk
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanUtuhTCO)).Decode(&masuk); err != nil {
			galat(w, http.StatusBadRequest, "badan permintaan bukan JSON simpan utuh yang sah")
			return
		}
		id := strings.TrimSpace(masuk.Kontrak.ID)
		switch {
		case !perbarui && id != "":
			galat(w, http.StatusBadRequest, "identitas kontrak dibuat server; POST kontrak-utuh tidak boleh membawa id")
			return
		case perbarui && id != "" && id != r.PathValue("kid"):
			galat(w, http.StatusBadRequest, "id kontrak di badan berbeda dari id di jalur")
			return
		case perbarui:
			masuk.Kontrak.ID = r.PathValue("kid")
		}
		h, err := layananSimpanUtuhTCO(svc).Simpan(r.Context(), pelakuDari(r, stub), r.PathValue("id"), masuk)
		if jawabGalatUtuhTCO(w, err) {
			return
		}
		tulisJSONPolis(w, h)
	}
}

// penangkapKodeTCO menangkap kode HTTP pemetaan galat; badannya dibuang -
// pesan yang menyebut baris ditulis sesudahnya.
type penangkapKodeTCO struct {
	http.ResponseWriter
	kode int
}

func (p *penangkapKodeTCO) WriteHeader(kode int)        { p.kode = kode }
func (p *penangkapKodeTCO) Write(b []byte) (int, error) { return len(b), nil }

// jawabGalatUtuhTCO - kode dari galat baris; pesan menyebut baris yang gagal.
func jawabGalatUtuhTCO(w http.ResponseWriter, err error) bool {
	var g services.GalatSimpanUtuhTCO
	if err == nil || !errors.As(err, &g) {
		if errors.Is(err, services.ErrSimpanUtuhTidakSah) {
			galat(w, http.StatusBadRequest, err.Error())
			return true
		}
		return jawabGalatTreatyContractOut(w, err)
	}
	p := &penangkapKodeTCO{ResponseWriter: w}
	jawabGalatTreatyContractOut(p, g.Galat)
	// 500 tidak membocorkan isi galat; 503 SENGAJA menyebut master yang belum
	// siap (ADR-0015) - pesannya dipertahankan (temuan /code-review).
	if p.kode == http.StatusInternalServerError {
		galat(w, p.kode, "simpan utuh dibatalkan seluruhnya - gagal pada "+g.Bagian+" (galat server)")
		return true
	}
	galat(w, p.kode, g.Error())
	return true
}
