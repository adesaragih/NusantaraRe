package models

// Untuk apa berkas ini: KONTEKS SATU AKSI dan JEJAK KRONOLOGI (DataTransform `ChronologyInsertion_DT`, Activity
// `SetchronologyKlaimFacIn`). Pola konteks disalin dari `modul/claimnonprop/backend/models/konteks.go`, bukan diimpor.
//
// `[terverifikasi]` `Claim Fac In/DataTransform/ChronologyInsertion_DT.xml`: WHEN `OperatorID.pyPosition != "IT
// Developer"` -> APPEND `.ClaimData.Chronology`: `.pyNote <- Data.CARI12`, `.ASMUser <- OperatorID.pyUserName`,
// `.ASMDateTimeChronology <- CurrentDateTime()`, `.ASMNoteType <- label langkah alur terakhir` lalu dipetakan
// (Choose Surveyor -> Adjustment, "" -> Committee, Input Register -> Registration, Input Estimasi -> Estimation, View
// Polis -> View Policy), dan `.ASMUserID <- "Claim Admin"` kecuali Committee.
//
// ⚠️ `SetchronologyKlaimFacIn` langkah 1.1 memetakan TIGA NAMA ORANG tertulis mati ke jabatan - DIBUANG (bawaan c prompt
// §3, pola Claim Prop): tingkat wewenang dibaca dari roster EMAILKOMITE STS_KLAIM FACIN (kunci OPERATOR_ID / workbasket
// pelaku, label JABATAN); tidak ketemu -> `TingkatClaimAdmin`. Nol nama orang di kode. Pengecualian "IT Developer"
// dibuang: semua peran terekam.
//
// Penyimpanan: T_VIEW_SUGGEST (CLAIM_ID, tabel bersama PremiumList Life, keputusan work owner 07-10-2026 "1 tabel aja")
// - `.ASMDateTimeChronology` DATE_SUGGEST, `.ASMUser` PIC_SUGGEST, `.ASMUserID` IS_CEDING_CONFIRM, `.pyNote`
// COMMENT_SUGGEST, `.ASMNoteType` INITIAL_SUGGEST (kolom bersama yang ada; OQ-CFI kolom jenis catatan).

import (
	"context"
	"sort"
	"strings"
	"time"
)

// TingkatClaimAdmin - nilai dasar `.ASMUserID` (ChronologyInsertion_DT langkah 1.1.11.1).
const TingkatClaimAdmin = "Claim Admin"

// JenisCatatanKomite - `.ASMNoteType` langkah tanpa label (kasus komite, ChronologyInsertion_DT 1.1.7).
const JenisCatatanKomite = "Committee"

// jenisCatatan - pemetaan label langkah -> `.ASMNoteType` (ChronologyInsertion_DT 1.1.6-1.1.10).
var jenisCatatan = map[string]string{
	"Choose Surveyor": "Adjustment",
	"":                JenisCatatanKomite,
	"Input Register":  "Registration",
	"Input Estimasi":  "Estimation",
	"View Polis":      "View Policy",
}

// JenisCatatan - `.ASMNoteType` untuk label langkah alur `label` (label lain apa adanya).
func JenisCatatan(label string) string {
	if v, ada := jenisCatatan[label]; ada {
		return v
	}
	return label
}

// Konteks membawa apa yang satu aksi butuhkan dari luar halaman.
type Konteks struct {
	Ctx   context.Context
	Acuan Acuan
	// Pelaku - akun pelaku (`OperatorID.pyUserName`).
	Pelaku string
	// Tingkat - tingkat wewenang pelaku dari roster komite FACIN (`TingkatClaimAdmin` bila tidak ada).
	Tingkat string
	// Langkah - label langkah alur yang sedang dikerjakan (`pxSteps(<LAST>).pyLabel`): label assignment kasus.
	Langkah string
	// Sekarang - waktu aksi, zona Jakarta.
	Sekarang time.Time
	// Produksi - padanan `When/IsPEGAPROD` (konfigurasi IS_PEGA_PROD).
	Produksi bool
}

// Jakarta - zona waktu Asia/Jakarta (WIB, UTC+7).
var Jakarta = time.FixedZone("WIB", 7*3600)

// Hari - tanggal hari aksi ("2006-01-02") - `@CurrentDate(...)`.
func (k *Konteks) Hari() string { return k.Sekarang.In(Jakarta).Format("2006-01-02") }

// Waktu - stempel waktu aksi ("2006-01-02 15:04:05") - `@CurrentDateTime()`.
func (k *Konteks) Waktu() string { return k.Sekarang.In(Jakarta).Format("2006-01-02 15:04:05") }

// Ctxt - konteks Go (latar bila kosong).
func (k *Konteks) Ctxt() context.Context {
	if k.Ctx == nil {
		return context.Background()
	}
	return k.Ctx
}

// DaftarKronologi - `pyWorkPage.ClaimData.Chronology` (grid "Claim Status" ViewHistoryClaim).
const DaftarKronologi = CD + "Chronology"

// PropRiwayatBaru - penanda baris Chronology yang lahir di aksi ini (repository menyisipkannya; T_VIEW_SUGGEST hanya
// bertambah, tidak pernah ditulis ulang).
const PropRiwayatBaru = "Baru"

// Kronologi meniru `Data.CARI12 = teks` + `Apply-DataTransform ChronologyInsertion_DT` (+ pemetaan jabatan
// SetchronologyKlaimFacIn tanpa nama orang). SetchronologyKlaimFacIn 1.1 memetakan jabatan untuk SETIAP baris pelaku yang
// dikenalnya (bukan hanya baris Committee); di sini jabatan roster pelaku untuk setiap baris, selainnya "Claim Admin"
// (pola Claim Prop; temuan review 10-10-2026).
func (k *Konteks) Kronologi(h *Halaman, teks string) {
	jenis := JenisCatatan(k.Langkah)
	tingkat := TingkatClaimAdmin
	if strings.TrimSpace(k.Tingkat) != "" {
		tingkat = k.Tingkat
	}
	h.TambahBaris(DaftarKronologi, Baris{
		"pyNote":                teks,
		"ASMUser":               k.Pelaku,
		"ASMDateTimeChronology": k.Waktu(),
		"ASMNoteType":           jenis,
		"ASMUserID":             tingkat,
		PropRiwayatBaru:         "1",
	})
}

// TampilKronologi - grid "Claim Status" (ViewHistoryClaim, kolom Date Input | User | Status | Note): urut
// `.ASMDateTimeChronology` menurun (`SetchronologyKlaimFacIn` langkah 2 Obj-Sort).
func TampilKronologi(d []Baris) []Baris {
	out := SalinDaftar(d)
	sort.SliceStable(out, func(i, j int) bool { return out[i]["ASMDateTimeChronology"] > out[j]["ASMDateTimeChronology"] })
	return out
}
