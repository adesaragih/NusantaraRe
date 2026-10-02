// Package login - masuk dengan akun dan sandi (`M_LOGIN_GO`), keputusan work
// owner 01-10-2026.
//
// Untuk apa paket ini: memeriksa akun dan sandi, mengunci akun sesudah salah
// berulang, menerbitkan dan membaca cookie sesi bertanda tangan, dan mengisi
// `inti.Pelaku` dari sesi itu - peran pelaku adalah workbasket-nya
// (`M_LOGIN_GO_WORKBASKET`, `WORKBASKET_ID` = nama peran).
//
// ⛔ Sandi tidak pernah disimpan, dicatat, atau dikirim balik - hanya hash
// bcrypt-nya. Tanpa tabel sesi: cookie memuat akun, versi sesi, dan waktunya,
// ditandatangani HMAC; `SESSION_VERSION` yang naik mencabut semua cookie lama
// akun itu.
//
// Dibaca sesudah: inti/backend/pelaku_http.go, migrations/902_m_login_go.sql.
package login
