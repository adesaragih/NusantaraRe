// Package handlers memuat HTTP controller.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak
// mengimpor repository, dan tidak pernah memegang koneksi database: seluruh
// pertanyaan ke Oracle - termasuk pemeriksaan kesehatan - lewat services.
// Nol aturan dagang di sini.
//
// [usulan] Router memakai net/http bawaan Go 1.22 (pola "GET /path"),
// sehingga nol ketergantungan pihak ketiga.
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"nusantarare/internal/services"
)

// Router menyusun seluruh rute.
//
// stubPelaku datang dari konfigurasi dan disetel SEKALI di sini. Ia bukan
// autentikasi - lihat pelaku.go - dan konfigurasi sudah menolaknya bila
// lingkungan menunjuk produksi.
func Router(svc *services.Service, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz(svc))
	mux.HandleFunc("GET /api/klaim-life/{id}", klaimLife(svc))
	// Kotak masuk per tahap - F0.4. Rutenya GET pada koleksi yang sama
	// dengan POST pendaftaran: satu sumber daya, dua metode.
	mux.HandleFunc("GET /api/klaim-life", kotakMasuk(svc, stubPelaku))
	mux.HandleFunc("POST /api/klaim-life", daftarKlaim(svc, stubPelaku))
	mux.HandleFunc("GET /api/peserta-life", cariPeserta(svc))
	// Pencarian diagnosa - kelompok Medis. Tabelnya 97.586 baris, dan
	// batasnya dari rule (pyMaxRecords 500), bukan dari klien.
	mux.HandleFunc("GET /api/penyakit-life", cariPenyakit(svc, stubPelaku))
	mux.HandleFunc("GET /api/klaim-life/{id}/dampak-hapus", dampakHapus(svc, stubPelaku))
	mux.HandleFunc("DELETE /api/klaim-life/{id}", hapusKlaim(svc, stubPelaku))
	mux.HandleFunc("POST /api/klaim-life/{id}/peserta/{pesertaId}/akseptasi",
		simpanAdjustment(svc, stubPelaku))
	mux.HandleFunc("POST /api/klaim-life/{id}/peserta/{pesertaId}/putaran",
		tambahPutaran(svc, stubPelaku))
	mux.HandleFunc(
		"POST /api/klaim-life/{id}/peserta/{pesertaId}/adjustment/{adjId}/komite",
		serahkanKomite(svc, stubPelaku))
	// Perpindahan tahap - butir aw. Dua tombol layar Outstanding
	// memanggilnya: `Send Back to Register` (InputOSClaimLife.xml
	// 21404) dan `Send to Medical Check` (21839). Tujuannya KATA,
	// bukan angka: jalur berisi /tahap/2 tidak terbaca siapa pun, dan
	// angka yang bergeser memindahkan kasus ke tempat yang salah.
	mux.HandleFunc("POST /api/klaim-life/{id}/tahap/{tujuan}", pindahTahap(svc, stubPelaku))
	// Gerbang Close Claim: memeriksa, belum menyelesaikan penugasan.
	mux.HandleFunc("GET /api/klaim-life/{id}/boleh-tutup", bolehTutup(svc))
	// Penutupan kasus - butir bb. POST, sebab ia MENGUBAH: STATUS_WORK
	// diisi Resolved-Completed dan TAHAP dikosongkan.
	mux.HandleFunc("POST /api/klaim-life/{id}/tutup", tutupKlaim(svc, stubPelaku))
	mux.HandleFunc("POST /api/klaim-life/{id}/adjustment/{adjId}/tolak",
		tolakBaris(svc, stubPelaku))
	mux.HandleFunc("PUT /api/klaim-life/{id}/peserta/{pesertaId}/tanggal-kejadian",
		setTanggalKejadian(svc, stubPelaku))
	return mux
}

type jawabanSehat struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// healthz menjawab tanpa menyentuh aturan dagang mana pun.
//
// Ia tetap menjawab 200 ketika Oracle belum dikonfigurasi: Fase 0 harus dapat
// dijalankan tanpa instance, dan keadaan database dilaporkan apa adanya di
// dalam badan jawaban, bukan disembunyikan.
func healthz(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jawab := jawabanSehat{Status: "sehat", Database: "tidak dikonfigurasi"}

		if svc.PunyaDatabase() {
			ctx, batal := context.WithTimeout(r.Context(), 3*time.Second)
			defer batal()
			if err := svc.CekKesehatan(ctx); err != nil {
				jawab.Database = "tidak terjangkau"
			} else {
				jawab.Database = "terjangkau"
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(jawab)
	}
}
