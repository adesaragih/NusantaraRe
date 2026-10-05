package db

import go_ora "github.com/sijms/go-ora/v2"

// Biner membungkus isi berkas sebagai nilai bind kolom BLOB (Template Manager, migrasi inti 912).
//
// ⛔ Bukan []byte polos: go-ora mengikat []byte sebagai RAW, yang berubah menjadi LONG RAW di atas batas RAW.
// `go_ora.Blob` diikat sebagai LOB, aman berapa pun ukurannya.
func Biner(isi []byte) any { return go_ora.Blob{Data: isi} }

// PindaiBiner adalah tujuan Scan kolom BLOB.
type PindaiBiner struct{ b go_ora.Blob }

// Scan memenuhi sql.Scanner.
func (p *PindaiBiner) Scan(v any) error { return p.b.Scan(v) }

// Isi - isi kolom; nil bila NULL.
func (p *PindaiBiner) Isi() []byte { return p.b.Data }

// PindaiTeksPanjang adalah tujuan Scan kolom CLOB (mis. JSON Pega lama yang hanya DIBACA alat Copy Old).
type PindaiTeksPanjang struct{ c go_ora.Clob }

// Scan memenuhi sql.Scanner.
func (p *PindaiTeksPanjang) Scan(v any) error { return p.c.Scan(v) }

// Teks - isi kolom; "" bila NULL.
func (p *PindaiTeksPanjang) Teks() string { return p.c.String }
