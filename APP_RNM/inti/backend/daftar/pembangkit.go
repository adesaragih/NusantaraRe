package daftar

// Direktif pembangkit daftar modul - di berkas TERSENDIRI yang namanya berurut
// SESUDAH `modul_*_gen.go`, dan itu disengaja.
//
// `go generate` membaca berkas paket satu per satu menurut abjad dan menjalankan
// setiap direktif yang ditemuinya. Pembangkit MEMBUANG berkas bangkitan modul
// yang foldernya sudah tidak punya `backend/modul.go`. Bila direktifnya di
// berkas yang berurut lebih dulu (dulu `daftar.go`), `go generate` masih hendak
// membaca berkas yang baru saja dibuang itu dan berhenti bergalat - padahal
// hasilnya sudah benar (ditemukan uji coba bab 5, 30-09-2026). Di sini setiap
// `modul_*_gen.go` sudah terbaca sebelum pembangkit berjalan. Letaknya dijaga
// `bangkit/main_test.go`.

//go:generate go run ./bangkit
