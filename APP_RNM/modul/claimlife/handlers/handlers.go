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
	"net/http"

	"nusantarare/modul/claimlife/services"
)

// Router menyusun rute modul ini SAJA, di mux sendiri - dipakai uji HTTP.
//
// Refactor bentuk B (30-09-2026): dulu `Router` menyusun SELURUH rute
// aplikasi. Kini `cmd/api` menyusun mux bersama dari modul yang AKTIF
// (MODUL_AKTIF), dan `GET /healthz` - milik aplikasi, bukan modul - tinggal
// di sana.
func Router(svc *services.Service, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, stubPelaku)
	return mux
}

// DaftarkanRute mendaftarkan seluruh rute modul Claim Life.
//
// stubPelaku datang dari konfigurasi dan disetel SEKALI di sini. Ia bukan
// autentikasi - lihat pelaku.go - dan konfigurasi sudah menolaknya bila
// lingkungan menunjuk produksi.
func DaftarkanRute(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
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
}
