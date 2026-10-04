package handlers

// Pintu HTTP kotak masuk - F0.4.
//
// `GET /api/klaim-life?tahap=<1..4>&halaman=&ukuran=`
//
// ⚠️ Handler TIPIS. Seluruh gerbang ada di `services.Inbox.Ambil`; di sini
// hanya penguraian parameter dan penerjemahan galat ke kode HTTP.
//
// Dibaca sesudah: handlers.go.

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/services"
)

// barisInboxJSON adalah satu baris antrian di kabel.
//
// ⛔ Nol medan uang (ADR-U-0003). Nilai klaim tidak ikut di daftar; layar yang
// memerlukannya membuka kasusnya.
type barisInboxJSON struct {
	// caseId - label `Case ID`, `InboxPremiumList.xml:721`.
	CaseID string `json:"caseId"`
	// id - pengenal work, dipakai membuka kasusnya.
	ID string `json:"id"`
	// tahap - nama assignment VERBATIM (butir at).
	Tahap string `json:"tahap"`
	// createOpName - label `Create Operator Name`, baris 751.
	CreateOpName string `json:"createOpName"`
	// tglCreate - label `Create Date/Time`, baris 736. Kosong bila belum ada.
	TglCreate string `json:"tglCreate"`
	// status - label `Work Status`, baris 765. KATA, bukan kode mentah.
	Status     string `json:"status"`
	NomorKlaim string `json:"nomorKlaim"`
	NomorPolis string `json:"nomorPolis"`
	NamaBisnis string `json:"namaBisnis"`
	MataUang   string `json:"mataUang"`
}

type halamanInboxJSON struct {
	Baris []barisInboxJSON `json:"baris"`
	// total adalah cacah SELURUH kasus pada tahap itu - lencana tab
	// memakainya, bukan panjang halaman.
	Total   int    `json:"total"`
	Tahap   int    `json:"tahap"`
	Halaman int    `json:"halaman"`
	Ukuran  int    `json:"ukuran"`
	Nama    string `json:"namaTahap"`
}

// bilanganKueri membaca satu parameter bilangan, dengan nilai bawaan.
//
// ⚠️ Parameter yang ADA tetapi bukan bilangan DITOLAK, tidak diam-diam jatuh
// ke bawaan: `?ukuran=limapuluh` yang dijawab 50 menyembunyikan salah ketik.
func bilanganKueri(r *http.Request, nama string, bawaan int) (int, bool) {
	teks := r.URL.Query().Get(nama)
	if teks == "" {
		return bawaan, true
	}
	n, err := strconv.Atoi(teks)
	if err != nil {
		return 0, false
	}
	return n, true
}

// kotakMasuk melayani GET /api/klaim-life.
func kotakMasuk(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		nomorTahap, sah := bilanganKueri(r, "tahap", int(models.TahapOutstanding))
		if !sah {
			galat.Tulis(w, http.StatusBadRequest, "parameter tahap harus bilangan 1..4")
			return
		}
		halaman, sah := bilanganKueri(r, "halaman", 1)
		if !sah {
			galat.Tulis(w, http.StatusBadRequest, "parameter halaman harus bilangan")
			return
		}
		ukuran, sah := bilanganKueri(r, "ukuran", 0)
		if !sah {
			galat.Tulis(w, http.StatusBadRequest, "parameter ukuran harus bilangan")
			return
		}
		if halaman < 1 {
			halaman = 1
		}
		tahap := models.Tahap(nomorTahap)
		// ⚠️ Offset diturunkan dari nomor halaman, bukan diterima mentah:
		// offset mentah membiarkan pemanggil melompati batas halaman dan
		// membuat lencana tidak cocok dengan isinya.
		offset := (halaman - 1) * services.BatasUkuran(ukuran)

		hal, err := svc.KotakMasuk().Ambil(r.Context(), inti.PelakuDari(r, stubPelaku),
			tahap, offset, ukuran)
		switch {
		case errors.Is(err, inti.ErrTanpaIdentitas):
			galat.Tulis(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
			return
		case errors.Is(err, services.ErrTahapTidakSah):
			galat.Tulis(w, http.StatusBadRequest, "tahap harus 1..4")
			return
		case errors.Is(err, inti.ErrTanpaWewenang):
			// ⛔ 403, BUKAN daftar kosong. Daftar kosong terbaca sebagai
			// "tidak ada pekerjaan".
			galat.Tulis(w, http.StatusForbidden, "peran tidak memegang tahap ini")
			return
		case err != nil:
			// ⚠️ Sebabnya DICATAT di log server (GILIRAN-12 paket 3): uji asap
			// DEV menjumpai satu 500 sesaat di sini yang tidak dapat
			// direproduksi, dan tanpa baris log ia tidak dapat ditelusuri.
			log.Printf("kotak masuk Claim Life: %v", err)
			galat.Tulis(w, http.StatusInternalServerError, "gagal membaca kotak masuk")
			return
		}

		keluar := halamanInboxJSON{
			Baris:   make([]barisInboxJSON, 0, len(hal.Baris)),
			Total:   hal.Total,
			Tahap:   int(hal.Tahap),
			Halaman: halaman,
			Ukuran:  hal.Ukuran,
			Nama:    hal.Tahap.String(),
		}
		for _, b := range hal.Baris {
			stempel := ""
			if !b.TglCreate.IsZero() {
				stempel = b.TglCreate.Format(time.RFC3339)
			}
			keluar.Baris = append(keluar.Baris, barisInboxJSON{
				CaseID: b.CaseID, ID: b.WorkID, Tahap: b.Tahap,
				CreateOpName: b.CreateOpName, TglCreate: stempel,
				// Status sebagai KATA - `models.StatusBaris.String()`.
				Status:     kontrak.StatusBarisDariKode(b.StatusKlaim).String(),
				NomorKlaim: b.NomorKlaim, NomorPolis: b.NomorPolis,
				NamaBisnis: b.NamaBisnis, MataUang: b.MataUang,
			})
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(keluar)
	}
}
