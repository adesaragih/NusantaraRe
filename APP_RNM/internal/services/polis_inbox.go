package services

// Kotak masuk PremiumList Life - tiket 01 bagian 2.
//
// Untuk apa berkas ini: gerbang identitas atas pembacaan kotak masuk, dan
// penjepitan halaman. Aturan bentuknya ada di repository; yang di sini
// siapa yang boleh membacanya.
//
// ⚠️ Nol gerbang PERAN. `InboxPremiumList.xml` tidak menyaring per peran -
// ia daftar kerja seluruh inputor Life, dan menyempitkannya ke satu peran
// berarti menutup daftar yang di sistem lama terbuka. Yang membatasi
// nantinya `KomiteRouter`-nya sendiri di modul Komite, bukan layar ini.
//
// Dibaca sesudah: repository/polis_inbox.go.

import (
	"context"
	"time"

	"nusantarare/internal/repository"
	"nusantarare/inti"
	"nusantarare/inti/db"
)

// BarisInboxPolis adalah satu baris kotak masuk, sebagaimana LAYAR
// membutuhkannya.
//
// ⛔ Tipe MILIK services, bukan tipe repository yang diteruskan. Handlers
// tidak boleh mengimpor repository (arah ketergantungan, tiket 01 AC-7), dan
// meneruskan tipe repository lewat services membuat aturan itu benar hanya
// pada impornya - tidak pada kenyataannya.
//
// ⚠️ Tanggal menyeberang sebagai TEKS `YYYY-MM-DD`, sebagaimana kotak
// masuk Claim Life: layar tidak mengurai ulang tanggal, dan `null` tetap
// dapat dibedakan dari tanggal nol lewat teks kosong.
type BarisInboxPolis struct {
	CaseID           string `json:"caseId"`
	TglCreate        string `json:"tglCreate"`
	CreateOpName     string `json:"createOpName"`
	StatusWork       string `json:"statusWork"`
	Position         string `json:"position"`
	CedingCoName     string `json:"cedingCoName"`
	PolicyHolderName string `json:"policyHolderName"`
	PLNumber         string `json:"plNumber"`
	RISlipRNM        string `json:"riSlipRnm"`
	Type             string `json:"type"`
	MarketingName    string `json:"marketingName"`
	SobName          string `json:"sobName"`
	DateReceived     string `json:"dateReceived"`
}

// HalamanInboxPolis adalah satu halaman beserta cacah totalnya.
type HalamanInboxPolis struct {
	Baris []BarisInboxPolis `json:"baris"`
	// Total mencacah SELURUH baris yang cocok, bukan yang di halaman ini.
	Total   int `json:"total"`
	Halaman int `json:"halaman"`
	Ukuran  int `json:"ukuran"`
}

// InboxPolis melayani pembacaan kotak masuk polis.
type InboxPolis struct{ svc *Service }

// InboxPolis menyusun layanannya.
func (s *Service) InboxPolis() *InboxPolis { return &InboxPolis{svc: s} }

// Ambil membaca satu halaman kotak masuk.
//
// `posisi` kosong berarti seluruh posisi.
func (i *InboxPolis) Ambil(ctx context.Context, pelaku inti.Pelaku,
	posisi string, halaman, ukuran int) (HalamanInboxPolis, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HalamanInboxPolis{}, err
	}
	if i == nil || i.svc == nil || !i.svc.PunyaDatabase() {
		return HalamanInboxPolis{}, db.ErrTanpaOracle
	}
	if halaman < 1 {
		halaman = 1
	}
	ukuran = repository.BatasUkuranHalamanPolis(ukuran)
	hasil, err := repository.NewInboxPolis(i.svc.DB()).Ambil(ctx, posisi, halaman, ukuran)
	if err != nil {
		return HalamanInboxPolis{}, err
	}
	keluar := HalamanInboxPolis{
		Baris: make([]BarisInboxPolis, 0, len(hasil.Baris)),
		Total: hasil.Total, Halaman: halaman, Ukuran: ukuran,
	}
	for _, b := range hasil.Baris {
		keluar.Baris = append(keluar.Baris, BarisInboxPolis{
			CaseID: b.CaseID, TglCreate: tanggalTeks(b.TglCreate),
			CreateOpName: b.CreateOpName, StatusWork: b.StatusWork,
			Position: b.Position, CedingCoName: b.CedingCoName,
			PolicyHolderName: b.PolicyHolderName, PLNumber: b.PLNumber,
			RISlipRNM: b.RISlipRNM, Type: b.Type,
			MarketingName: b.MarketingName, SobName: b.SobName,
			DateReceived: tanggalPenunjukTeks(b.DateReceived),
		})
	}
	return keluar, nil
}

// tanggalTeks menulis tanggal sebagai `YYYY-MM-DD`; nol menjadi kosong.
func tanggalTeks(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// tanggalPenunjukTeks menulis tanggal nullable; nil menjadi kosong.
func tanggalPenunjukTeks(t *time.Time) string {
	if t == nil {
		return ""
	}
	return tanggalTeks(*t)
}
