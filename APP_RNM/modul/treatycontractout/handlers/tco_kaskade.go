package handlers

// Pintu HTTP kaskade hapus + popup - tiket 10.
//
//	GET    /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/dampak-hapus                    isi popup Ya/Batal
//	DELETE /api/treaty-contract-out/tahun/{id}/kontrak/{kid}?reinsurer=&security=&business=&bersama=  Delete b11809 (Ya)
//	GET    /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/reinsurer/{rid}/dampak-hapus      isi popup reinsurer
//	DELETE /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/reinsurer/{rid}?security=         Delete b4936 (Ya)
//
// Nomor baris = `Section/InputTreatyContractReinsType.xml` /
// `Section/ViewDetailTreatyReinsurerGrid1.xml`. Jumlah di kueri DELETE = jumlah
// yang pemakai lihat di popup; server menolak (409) bila angkanya sudah lain.
//
// Dibaca sesudah: services/tco_kaskade.go.

import (
	"net/http"
	"strconv"

	"nusantarare/inti"
	"nusantarare/inti/galat"
	"nusantarare/modul/treatycontractout/services"
)

func layananKaskadeTCO(svc *services.Service) *services.KaskadeTCO {
	return svc.KaskadeTCO().
		DenganKaskade(services.KaskadeOracle(svc)).
		DenganKontrak(services.PemegangKontrakOracle(svc)).
		DenganTahun(services.GudangTahunTreatyOracle(svc)).
		DenganReinsurer(services.GudangReinsurerOracle(svc))
}

const jalurKontrakKaskadeTCO = "/api/treaty-contract-out/tahun/{id}/kontrak/{kid}"

func daftarkanRuteKaskadeTCO(mux *http.ServeMux, svc *services.Service, stub bool) {
	mux.HandleFunc("GET "+jalurKontrakKaskadeTCO+"/dampak-hapus", dampakKontrakTCO(svc, stub))
	mux.HandleFunc("DELETE "+jalurKontrakKaskadeTCO, hapusKontrakTCO(svc, stub))
	mux.HandleFunc("GET "+jalurKontrakKaskadeTCO+"/reinsurer/{rid}/dampak-hapus", dampakReinsurerTCO(svc, stub))
	mux.HandleFunc("DELETE "+jalurKontrakKaskadeTCO+"/reinsurer/{rid}", hapusReinsurerKaskadeTCO(svc, stub))
}

// konfirmasiDariKueri membaca jumlah yang dikonfirmasi pemakai; wajib.
func konfirmasiDariKueri(r *http.Request, kunci ...string) (services.KonfirmasiHapus, bool) {
	var k services.KonfirmasiHapus
	isi := map[string]*int64{"reinsurer": &k.Reinsurer, "security": &k.Security, "business": &k.Business, "bersama": &k.Bersama}
	for _, n := range kunci {
		v, err := strconv.ParseInt(r.URL.Query().Get(n), 10, 64)
		if err != nil || v < 0 {
			return services.KonfirmasiHapus{}, false
		}
		*isi[n] = v
	}
	return k, true
}

type jawabanHapusKaskade struct {
	Pesan string `json:"pesan"`
}

func dampakKontrakTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		d, err := layananKaskadeTCO(svc).DampakHapusKontrak(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"), r.PathValue("kid"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, d)
	}
}

func hapusKontrakTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		k, ok := konfirmasiDariKueri(r, "reinsurer", "security", "business", "bersama")
		if !ok {
			galat.Tulis(w, http.StatusBadRequest, "confirmed counts of reinsurer, security, business and other contracts are required (see dampak-hapus)")
			return
		}
		pesan, err := layananKaskadeTCO(svc).HapusKontrak(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"), r.PathValue("kid"), k)
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, jawabanHapusKaskade{Pesan: pesan})
	}
}

func dampakReinsurerTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		d, err := layananKaskadeTCO(svc).DampakHapusReinsurer(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"),
			r.PathValue("kid"), r.PathValue("rid"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, d)
	}
}

func hapusReinsurerKaskadeTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		k, ok := konfirmasiDariKueri(r, "security")
		if !ok {
			galat.Tulis(w, http.StatusBadRequest, "confirmed security count is required (see dampak-hapus)")
			return
		}
		k.Reinsurer = 1
		pesan, err := layananKaskadeTCO(svc).HapusReinsurer(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"),
			r.PathValue("kid"), r.PathValue("rid"), k)
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, jawabanHapusKaskade{Pesan: pesan})
	}
}
