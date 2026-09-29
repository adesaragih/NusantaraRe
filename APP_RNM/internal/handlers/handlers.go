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
//
// Refactor bentuk B (30-09-2026): rute modul yang sudah pindah ke
// `modul/<nama>/` TIDAK didaftarkan di sini - paket ini tidak boleh mengimpor
// modul lain. `cmd/api` menyerahkannya lewat `tambahan`, dan urutannya sama
// dengan sebelumnya: sesudah seluruh rute di bawah.
func Router(svc *services.Service, stubPelaku bool, tambahan ...func(*http.ServeMux)) http.Handler {
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
	// `Save to RNM` - tombol layar Outstanding (`InputOSClaimLife.xml`
	// b21102 -> `SaveOutStandingLife_Act`). POST, sebab ia MENULIS.
	mux.HandleFunc("POST /api/klaim-life/{id}/outstanding", simpanRNM(svc, stubPelaku))
	// Gerbang Close Claim: memeriksa, belum menyelesaikan penugasan.
	mux.HandleFunc("GET /api/klaim-life/{id}/boleh-tutup", bolehTutup(svc))
	// Penutupan kasus - butir bb. POST, sebab ia MENGUBAH: STATUS_WORK
	// diisi Resolved-Completed dan TAHAP dikosongkan.
	mux.HandleFunc("POST /api/klaim-life/{id}/tutup", tutupKlaim(svc, stubPelaku))
	mux.HandleFunc("POST /api/klaim-life/{id}/adjustment/{adjId}/tolak",
		tolakBaris(svc, stubPelaku))
	mux.HandleFunc("PUT /api/klaim-life/{id}/peserta/{pesertaId}/tanggal-kejadian",
		setTanggalKejadian(svc, stubPelaku))
	// Tiga tanggal klaim lain dialog Edit Date - `Save` b1910
	// (`EditDateClaimLife_Section.xml`). Satu rute, sebab satu tombol.
	mux.HandleFunc("PUT /api/klaim-life/{id}/peserta/{pesertaId}/tanggal-klaim",
		setTanggalKlaim(svc, stubPelaku))
	// Cabut peserta - OQ-M6 (GILIRAN-17), tombol `DELETE` `InputOSClaimLife`
	// b17865. POST, bukan DELETE: yang terjadi penandaan (ADR-U-0031).
	mux.HandleFunc("POST /api/klaim-life/{id}/peserta/{pesertaId}/cabut",
		cabutPeserta(svc, stubPelaku))
	// Grid diagnosa - butir bd. Tiga tombol, tiga rute, dan jalurnya
	// BERSARANG di bawah pesertanya: `SetDisease.xml` b389 menutup dengan
	// `Obj-Save pyWorkPage`, jadi diagnosa tidak punya hidup di luar peserta
	// yang memuatnya.
	mux.HandleFunc("POST /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa",
		tambahDiagnosa(svc, stubPelaku))
	mux.HandleFunc("PUT /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa/{diagId}",
		ubahDiagnosa(svc, stubPelaku))
	mux.HandleFunc("DELETE /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa/{diagId}",
		hapusDiagnosa(svc, stubPelaku))
	// Dokumen pendukung - butir be. Tiga tombol `DocumentLife.xml`.
	//
	// ⛔ Jalur UNDUH tidak bersarang di bawah klaimnya, dan itu meniru
	// aslinya: `URLPUBLIC` dicari dengan `imageid` saja
	// (`GetLinkStorage_SQL.xml` b91), dan nilai kolom itulah yang tersimpan
	// di kartu berkas. Batas klaimnya tetap ditegakkan services.
	mux.HandleFunc("POST /api/klaim-life/{id}/peserta/{pesertaId}/dokumen",
		unggahDokumen(svc, stubPelaku))
	mux.HandleFunc("GET /api/dokumen/{dokId}/isi", isiDokumen(svc, stubPelaku))
	mux.HandleFunc("DELETE /api/klaim-life/{id}/dokumen/{dokId}",
		hapusDokumen(svc, stubPelaku))
	// --- modul yang sudah pindah ke modul/<nama>/ (Treaty Contract Out,
	// PremiumList Life, Komite Claim Life) ---
	//
	// Rutenya hidup di modulnya masing-masing; `cmd/api` yang menyerahkannya.
	for _, daftarkan := range tambahan {
		daftarkan(mux)
	}
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
