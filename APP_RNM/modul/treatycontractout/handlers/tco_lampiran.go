package handlers

// Pintu HTTP lampiran tahun treaty - tiket 12 Treaty Contract Out.
//
//	GET    /api/treaty-contract-out/kategori-lampiran
//	GET    /api/treaty-contract-out/tahun/{id}/lampiran                 Refresh b1023
//	POST   /api/treaty-contract-out/tahun/{id}/lampiran                 Add attachment b578 (multipart)
//	GET    /api/treaty-contract-out/tahun/{id}/lampiran/semua           Download All b2659 (zip)
//	GET    /api/treaty-contract-out/tahun/{id}/lampiran/selaras         keselarasan rekam-berkas (AC 61)
//	GET    /api/treaty-contract-out/tahun/{id}/lampiran/{lid}/isi       tautan nama berkas b3470
//	POST   /api/treaty-contract-out/tahun/{id}/lampiran/{lid}/ulangi    unggah ulang (AC 58)
//	DELETE /api/treaty-contract-out/tahun/{id}/lampiran/{lid}           Delete b3897
//
// Nomor baris = `Section/GridTreatyArrangementAttachment.xml`.
//
// ⛔ Setiap jalur menyebut TAHUN treaty-nya, termasuk unduhan: lampiran
// dibaca dengan batas tahun, bukan dengan ID lampiran saja. Unduhan
// bergerbang identitas seperti rute lain; layar mengambilnya lewat `fetch`
// berheader identitas, bukan tautan biasa.
//
// Dibaca sesudah: rute_treaty_contract_out.go, services/tco_lampiran.go.

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	"nusantarare/inti"
	"nusantarare/inti/galat"
	"nusantarare/inti/unggah"
	"nusantarare/modul/treatycontractout/services"
)

// batasPermintaanLampiran membatasi SELURUH badan multipart: berkas + formulir.
//
// ⚠️ Batas berkasnya sendiri tetap `services.BatasUkuranUnggahan`, ditegakkan
// saat menyalin; yang di sini supaya permintaan raksasa berhenti di pintu.
const batasPermintaanLampiran = unggah.BatasUkuranUnggahan + galat.BatasFormulir

// LayananLampiranTCO memasang seluruh implementasi nyata - dipakai rute di
// berkas ini DAN pekerja latar di `cmd/api` (OQ-TCO-09).
//
// Penyimpanannya dipilih `PenyimpananLampiranTCO` menurut `PELAKSANA_STORAGE`
// (OQ-TCO-08, keputusan work owner 29-09-2026): bawaan STUB LOKAL; `nyata`
// hanya bila work owner menyetelnya.
func LayananLampiranTCO(svc *services.Service) *services.LampiranTahunTCO {
	return svc.LampiranTahunTCO().
		DenganGudang(services.GudangLampiranOracle(svc)).
		DenganKategori(services.KategoriLampiranOracle(svc)).
		DenganAntrean(services.AntreanLampiranOracle(svc)).
		DenganPenyimpanan(services.PenyimpananLampiranTCO(svc)).
		DenganTahun(services.GudangTahunTreatyOracle(svc))
}

func daftarkanRuteLampiranTCO(mux *http.ServeMux, svc *services.Service, stub bool) {
	mux.HandleFunc("GET /api/treaty-contract-out/kategori-lampiran", kategoriLampiranTCO(svc, stub))
	mux.HandleFunc("GET /api/treaty-contract-out/tahun/{id}/lampiran", daftarLampiranTCO(svc, stub))
	mux.HandleFunc("POST /api/treaty-contract-out/tahun/{id}/lampiran", unggahLampiranTCO(svc, stub))
	mux.HandleFunc("GET /api/treaty-contract-out/tahun/{id}/lampiran/semua", unduhSemuaLampiranTCO(svc, stub))
	mux.HandleFunc("GET /api/treaty-contract-out/tahun/{id}/lampiran/selaras", selarasLampiranTCO(svc, stub))
	mux.HandleFunc("GET /api/treaty-contract-out/tahun/{id}/lampiran/{lid}/isi", unduhLampiranTCO(svc, stub))
	mux.HandleFunc("POST /api/treaty-contract-out/tahun/{id}/lampiran/{lid}/ulangi", ulangiLampiranTCO(svc, stub))
	mux.HandleFunc("DELETE /api/treaty-contract-out/tahun/{id}/lampiran/{lid}", hapusLampiranTCO(svc, stub))
}

// punyaDBTCO menjawab 503 bila Oracle belum dikonfigurasi.
func punyaDBTCO(w http.ResponseWriter, svc *services.Service) bool {
	if !svc.PunyaDatabase() {
		galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
		return false
	}
	return true
}

type jawabanKategoriLampiran struct {
	Daftar []string `json:"daftar"`
	Total  int      `json:"total"`
}

func kategoriLampiranTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		d, err := LayananLampiranTCO(svc).Kategori(r.Context(), inti.PelakuDari(r, stub))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, jawabanKategoriLampiran{Daftar: d, Total: len(d)})
	}
}

type jawabanDaftarLampiran struct {
	Daftar []services.LampiranTampil `json:"daftar"`
	Total  int                       `json:"total"`
}

func daftarLampiranTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		d, err := LayananLampiranTCO(svc).Daftar(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, jawabanDaftarLampiran{Daftar: d, Total: len(d)})
	}
}

func unggahLampiranTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, batasPermintaanLampiran)
		if err := r.ParseMultipartForm(galat.BatasFormulir); err != nil {
			var besar *http.MaxBytesError
			if errors.As(err, &besar) {
				galat.Tulis(w, http.StatusRequestEntityTooLarge, "berkas melebihi batas ukuran")
				return
			}
			galat.Tulis(w, http.StatusBadRequest, "permintaan bukan multipart yang sah")
			return
		}
		berkas, kepala, err := r.FormFile("berkas")
		if err != nil {
			// VERBATIM `TreatyOutSaveAttachment.xml` b376.
			galat.Tulis(w, http.StatusBadRequest, services.PesanTanpaBerkasTCO)
			return
		}
		defer func() { _ = berkas.Close() }()
		hasil, err := LayananLampiranTCO(svc).Unggah(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"),
			unggah.BerkasMasuk{NamaFile: kepala.Filename, Mime: kepala.Header.Get("Content-Type"),
				Kategori: r.FormValue("kategori"), Isi: berkas})
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(hasil)
	}
}

// kepalaUnduhan menyetel header unduhan yang aman.
//
// ⛔ `attachment`, bukan `inline`: isi berkas datang dari pengunggah, dan
// yang dirender di asal kita berjalan di asal kita. Nama berkas dikodekan
// `mime.FormatMediaType` (RFC 2231) - nama beraksara non-ASCII tidak memecah
// header.
func kepalaUnduhan(w http.ResponseWriter, jenis, nama string) {
	if jenis == "" {
		jenis = "application/octet-stream"
	}
	w.Header().Set("Content-Type", jenis)
	disposisi := mime.FormatMediaType("attachment", map[string]string{"filename": nama})
	if disposisi == "" {
		disposisi = "attachment"
	}
	w.Header().Set("Content-Disposition", disposisi)
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func unduhLampiranTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		meta, isi, err := LayananLampiranTCO(svc).Unduh(r.Context(), inti.PelakuDari(r, stub),
			r.PathValue("id"), r.PathValue("lid"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		defer func() { _ = isi.Close() }()
		kepalaUnduhan(w, meta.FileMimeType, meta.FileName)
		_, _ = io.Copy(w, isi)
	}
}

func unduhSemuaLampiranTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		tahunID := r.PathValue("id")
		// ⚠️ Arsip dibuka pada entri PERTAMA: kegagalan sebelum itu - identitas,
		// tahun tidak ada, rekam tanpa berkas - masih dapat dijawab sebagai
		// galat JSON biasa.
		var arsip *zip.Writer
		err := LayananLampiranTCO(svc).UnduhSemua(r.Context(), inti.PelakuDari(r, stub), tahunID,
			func(nama string, isi io.Reader) error {
				if arsip == nil {
					kepalaUnduhan(w, "application/zip", "lampiran-tahun-treaty-"+tahunID+".zip")
					arsip = zip.NewWriter(w)
				}
				f, err := arsip.Create(nama)
				if err != nil {
					return err
				}
				_, err = io.Copy(f, isi)
				return err
			})
		if arsip == nil {
			jawabGalatTreatyContractOut(w, err)
			return
		}
		// Galat di tengah aliran: arsipnya dibiarkan TIDAK tertutup - zip tanpa
		// direktori pusat ditolak pembukanya, bukan diterima sebagai arsip
		// yang diam-diam kurang isi.
		if err == nil {
			_ = arsip.Close()
		}
	}
}

type jawabanSelaras struct {
	Temuan []services.TemuanSelarasTCO `json:"temuan"`
	Total  int                         `json:"total"`
}

func selarasLampiranTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		t, err := LayananLampiranTCO(svc).PeriksaSelaras(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, jawabanSelaras{Temuan: t, Total: len(t)})
	}
}

func ulangiLampiranTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		hasil, err := LayananLampiranTCO(svc).Ulangi(r.Context(), inti.PelakuDari(r, stub),
			r.PathValue("id"), r.PathValue("lid"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

type jawabanHapusLampiran struct {
	Peringatan string `json:"peringatan"`
}

func hapusLampiranTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		p, err := LayananLampiranTCO(svc).Hapus(r.Context(), inti.PelakuDari(r, stub),
			r.PathValue("id"), r.PathValue("lid"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, jawabanHapusLampiran{Peringatan: p})
	}
}
