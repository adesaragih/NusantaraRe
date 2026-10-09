package models

// Untuk apa berkas ini: KONTEKS SATU AKSI dan JEJAK AUDIT (DataTransform `InsertChronology_DT`, Activity
// `SethistoryKlaimTreaty`; tiket 14, AC 77-83).
//
// `[terverifikasi]` `DataTransform/InsertChronology_DT.xml`: WHEN `OperatorID.pyPosition != "IT Developer"` -> APPEND
// `.ClaimData.SuggestList`: `.CommentSuggest <- DataChronology.CARI1`, `.PICSuggest <- OperatorID.pyUserName`,
// `.DateSuggest <- CurrentDateTime()`, `.IsCedingConfirm <- "Claim Admin"`, lalu tiga cabang nama orang di-hardcode
// menimpa tingkatnya.
//
// ⚠️ PENYIMPANGAN SADAR (keputusan work owner, spec §7/§10, AC 54, 55, 78, 80):
//   - tingkat wewenang dibaca dari SATU sumber, roster EMAILKOMITE (kunci OPERATOR_ID, label JABATAN baris STS_KLAIM
//     PROP); tidak ketemu -> `TingkatClaimAdmin`. Nol nama orang di kode.
//   - SEMUA peran terekam - pengecualian "IT Developer" dibuang (AC 80).
//   - pelaku (PICSuggest) dan tingkat (IsCedingConfirm) kolom terpisah (T_VIEW_SUGGEST.PIC_SUGGEST / IS_CEDING_CONFIRM).

import (
	"context"
	"sort"
	"strings"
	"time"
)

// TingkatClaimAdmin - nilai dasar tingkat wewenang (InsertChronology_DT langkah 1.1.4).
const TingkatClaimAdmin = "Claim Admin"

// Konteks membawa apa yang satu aksi butuhkan dari luar halaman.
type Konteks struct {
	Ctx   context.Context
	Acuan Acuan
	// Pelaku - akun pelaku (`OperatorID.pyUserIdentifier` / `pyUserName`; nama tampilan belum ada di login).
	Pelaku string
	// Tingkat - tingkat wewenang pelaku dari roster komite (`TingkatClaimAdmin` bila tidak ada).
	Tingkat string
	// Sekarang - waktu aksi, zona Jakarta.
	Sekarang time.Time
	// Produksi - padanan `When/IsPEGAPROD` (konfigurasi IS_PEGA_PROD).
	Produksi bool
}

// Jakarta - zona waktu Asia/Jakarta (WIB, UTC+7).
var Jakarta = time.FixedZone("WIB", 7*3600)

// Hari - tanggal hari aksi ("2006-01-02") - `@CurrentDate("dd/MM/yyyy", …)`.
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

// PropRiwayatBaru - penanda baris SuggestList yang lahir di aksi ini (repository menyisipkannya; T_VIEW_SUGGEST hanya
// bertambah, tidak pernah ditulis ulang).
const PropRiwayatBaru = "Baru"

// Riwayat meniru `DataChronology.CARI1 = teks` + `Apply-DataTransform InsertChronology_DT`.
func (k *Konteks) Riwayat(h *Halaman, teks string) {
	tingkat := k.Tingkat
	if strings.TrimSpace(tingkat) == "" {
		tingkat = TingkatClaimAdmin
	}
	h.TambahBaris(DaftarRiwayat, Baris{
		"CommentSuggest":  teks,
		"PICSuggest":      k.Pelaku,
		"DateSuggest":     k.Waktu(),
		"IsCedingConfirm": tingkat,
		PropRiwayatBaru:   "1",
	})
}

// TampilRiwayat - grid "Claim History" (Section OutstandingClaim / InputAcceptation, kolom "Name" | "Date" | "Noted")
// sesudah `SethistoryKlaimTreaty`: tingkat Claim Admin tampil "Admin Claim <pelaku>" (langkah 1.1), lainnya label
// tingkatnya; urut `.DateSuggest` MENURUN (langkah 1.11 Obj-Sort Descending - AC 82: dari kolom tanggal, bukan
// urutan baris). Langkah 1.2-1.10 (nama orang -> jabatan) DIBUANG (AC 55): tingkat sudah datang dari roster.
func TampilRiwayat(d []Baris) []Baris {
	out := SalinDaftar(d)
	for _, b := range out {
		t := b["IsCedingConfirm"]
		if strings.Contains(t, "Admin") || t == "" {
			b["IsCedingConfirm"] = "Admin Claim " + b["PICSuggest"]
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i]["DateSuggest"] > out[j]["DateSuggest"] })
	return out
}
