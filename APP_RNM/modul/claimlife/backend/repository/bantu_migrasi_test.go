package repository

// Folder migrasi Claim Life untuk uji paket ini - TANPA Oracle.
//
// Refactor bentuk B (30-09-2026): dulu `berkasMigrasi` adalah variabel
// `//go:embed` di paket ini. Penanaman kini tinggal di `modul/claimlife`
// (modul.go), dan uji di paket repository tidak boleh mengimpor paket perakit
// modulnya sendiri (siklus impor). Maka uji membaca folder yang SAMA dari
// disk - isi berkasnya identik dengan yang ditanam.

import (
	"io/fs"
	"os"
)

// berkasMigrasiDisk membuka folder `migrations/` modul ini sebagai fs.FS.
type berkasMigrasiDisk struct{ fs.FS }

// ReadFile membaca satu berkas, jalur `migrations/<nama>`.
func (b berkasMigrasiDisk) ReadFile(nama string) ([]byte, error) { return fs.ReadFile(b.FS, nama) }

// ReadDir membaca isi folder, jalur `migrations`.
func (b berkasMigrasiDisk) ReadDir(nama string) ([]fs.DirEntry, error) { return fs.ReadDir(b.FS, nama) }

// berkasMigrasi - folder `migrations/` modul ini (001-029).
var berkasMigrasi = berkasMigrasiDisk{os.DirFS("..")}
