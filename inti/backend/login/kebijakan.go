package login

import "time"

// Kebijakan login - seluruhnya disetujui work owner 01-10-2026.
const (
	// BatasGagal - sandi salah beruntun sebelum akun dikunci.
	BatasGagal = 5
	// LamaKunci - lama akun terkunci sesudah BatasGagal.
	LamaKunci = 15 * time.Minute
	// BatasDiam - sesi berakhir bila tidak dipakai selama ini.
	BatasDiam = 30 * time.Minute
	// BatasSesi - sesi berakhir paling lambat sekian sejak login, walau terus dipakai.
	BatasSesi = 10 * time.Hour
	// PanjangMinSandi - panjang sandi minimal, dalam KARAKTER.
	PanjangMinSandi = 10
	// PanjangMaksSandi - bcrypt hanya membaca 72 BYTE pertama; sisanya akan
	// diabaikan diam-diam, jadi ditolak terang.
	PanjangMaksSandi = 72
	// BiayaBcrypt - faktor kerja bcrypt.
	BiayaBcrypt = 12
	// NamaCookie - cookie sesi (HttpOnly, SameSite=Strict).
	NamaCookie = "rnm_sesi"
)
