package models

// Tab Correspondence kasus Inward Facultative (tiket 56) - `.OfferFacIn.CorrespondenceList` (kelas Data-Correspondence).

// PenerimaKorespondensi - satu `.RecipientList(n)`: nama (akun terpilih atau teks bebas) dan alamat email.
type PenerimaKorespondensi struct {
	Recipient, RecipientMail string
}

// Korespondensi - satu email korespondensi. No = ID baris (teks). Tanggal / DikirimPada: teks Pega DateTime di batas
// services <-> repository, kabel `DD-MM-YYYY HH:MM` WIB di batas services <-> handler. Message = HTML yang SUDAH
// dibersihkan. Terkirim = `.IsEditable` "0" (terkunci).
type Korespondensi struct {
	No, Tanggal, Sender, SenderMail, Subject, Message string
	Penerima                                          []PenerimaKorespondensi
	Terkirim                                          bool
	DikirimPada                                       string
}

// IsiKorespondensi - medan yang dapat diubah (Save / Send Email).
type IsiKorespondensi struct {
	Subject, Message string
	Penerima         []PenerimaKorespondensi
}

// Surel - satu email keluar yang siap dikirim pengirim SMTP: To / Cc alamat polos (sudah diperiksa), HTML bersih.
type Surel struct {
	To, Cc           []string
	ReplyTo, Subject string
	HTML             string
}
