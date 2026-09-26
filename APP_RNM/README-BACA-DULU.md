# APP_RNM — baca dulu sebelum menyentuh kode

Folder ini memuat **seluruh aplikasi** hasil migrasi Pega → Go + React + Oracle. Dokumen
(spec, tiket, ADR, discovery) tetap berada di folder induk `OUTPUT_HASIL_RNM\`. Ditulis
25 September 2026 untuk pembaca yang **baru mengenal Go dan React**.

## 1. Peta folder — satu kalimat per folder

| Folder / berkas | Isinya, dalam bahasa sehari-hari |
| --- | --- |
| `cmd/api/main.go` | Pintu masuk backend. Membaca pengaturan, menyambung Oracle, mendaftarkan alamat-alamat HTTP, lalu menunggu permintaan |
| `internal/config/` | Membaca **env var** (variabel lingkungan) seperti `ORACLE_DSN`. Tidak ada alamat atau kata sandi yang ditulis di kode |
| `internal/handlers/` | Penerima permintaan HTTP. Tugasnya sempit: baca permintaan, panggil *service*, tulis jawaban JSON |
| `internal/services/` | Aturan dagang dan perakitan data. Di sinilah "klaim punya peserta, peserta punya baris" disusun |
| `internal/repository/` | Satu-satunya lapisan yang berbicara ke Oracle. Seluruh SQL ada di sini, dan **hanya** di sini |
| `internal/repository/migrations/` | Berkas `.sql` bernomor yang membentuk tabel. Satu berkas = satu langkah, dan tiap langkah punya pasangan `_down.sql` untuk membatalkannya. Ditanam ke biner, jadi tidak perlu dicari di disk saat program jalan |
| `internal/repository/barislamakolom.go` | Daftar **62 kolom** tabel datar warisan beserta tipenya dari katalog Oracle, ditulis SEKALI; **55** di antaranya yang ditulis rule Pega. Pembaca, penulis fixture, dan tabel tiruan mengambil daftar yang sama, sehingga urutan `SELECT` dan urutan `Scan` tidak mungkin berselisih |
| `internal/repository/skemauji/` | Menyiapkan skema uji Oracle untuk test bertag `db`: menjalankan migrasi yang sama dengan aplikasi, mengisi fixture buatan, lalu membongkarnya |
| `internal/models/` | Bentuk data (struct): `Klaim`, `Peserta`, `BarisAdjustment`, `Money`, `Ratio` |
| `pkg/utils/` | Alat bantu umum: konversi teks ↔ desimal, format tanggal |
| `frontend/` | Tampilan React (TypeScript, berkas `.tsx`). Dijalankan Vite |
| `Makefile` | Daftar perintah: jalankan, bangun, uji. Setiap target adalah satu-dua perintah biasa |
| `bin/` | Hasil `go build` — tidak masuk git |

**Arah ketergantungan** — selalu satu arah, tidak pernah memotong:

```
handlers  →  services  →  repository  →  Oracle
```

`handlers` tidak boleh menyentuh `repository` langsung. Kalau Anda menulis SQL di luar
`repository`, atau memanggil `repository` dari `handlers`, itu salah tempat.

## 2. Satu permintaan, dari layar sampai Oracle dan kembali

Contoh: pengguna mengetik `UJI-KLAIM-1` lalu menekan **Buka**.

| Langkah | Berkas | Yang terjadi |
| ---: | --- | --- |
| 1 | `frontend/src/pages/KlaimLife.tsx` | Fungsi `cari` dipanggil; halaman masuk keadaan "memuat" |
| 2 | `frontend/src/services/api.ts` | `ambilKlaimLife(id)` mengirim `GET /api/klaim-life/UJI-KLAIM-1` |
| 3 | `internal/handlers/klaimlife.go` | Backend menerima, mengambil `{id}` dari alamat, memanggil *service* |
| 4 | `internal/services/klaimlife.go` | `Ambil` meminta header, peserta, dan baris ke *repository*, lalu merakitnya |
| 5 | `internal/repository/klaimlife.go` | Tiga query SQL ke Oracle. Uang diminta sebagai **teks** lewat `TO_CHAR` |
| 6 | `internal/models/klaimlife.go` | `MarshalJSON` menentukan bentuk JSON yang dikirim balik; kode status diterjemahkan menjadi kata |
| 7 | `KlaimLife.tsx` | Jawaban disimpan lewat `setKlaim`; React menggambar ulang tabel |

Nama medan JSON di langkah 6 (`nomorKlaim`, `jumlahKlaim`, …) **harus sama persis** dengan
`interface` di `api.ts`. Kalau salah satu diubah, yang lain ikut diubah.

## 3. Urutan membaca yang disarankan

**Go** — dari luar ke dalam: `cmd/api/main.go` → `internal/config/config.go` →
`internal/handlers/handlers.go` → `internal/handlers/klaimlife.go` →
`internal/services/klaimlife.go` → `internal/repository/klaimlife.go` →
`internal/models/klaimlife.go` → `internal/models/money.go`.

**React** — dari titik masuk: `frontend/src/main.tsx` → `App.tsx` → `pages/KlaimLife.tsx` →
`services/api.ts` → `store/index.ts`.

Setiap berkas punya komentar kepala yang menjelaskan **untuk apa berkas ini** dan istilah
yang dipakainya.

## 4. Istilah yang sering muncul

| Istilah | Artinya di proyek ini |
| --- | --- |
| **env var** | Nilai pengaturan yang dibaca dari lingkungan proses, bukan dari kode. Contohnya di `.env.example` |
| **handler** | Fungsi Go yang menjawab satu alamat HTTP |
| **service** | Fungsi Go yang memuat aturan dagang; boleh memanggil banyak *repository* |
| **repository** | Fungsi Go yang menjalankan SQL. Tidak tahu apa-apa soal HTTP |
| **model** | Struct Go yang menggambarkan satu jenis data |
| **seam** | Titik tempat test menggerakkan sistem: lewat HTTP, lewat *repository*, atau lewat *service* murni |
| **fixture** | Data contoh buatan untuk test — selalu berawalan `UJI-`, tidak pernah data sungguhan |
| **skema uji** | Tabel sementara yang dibuat test lalu dibuang lagi (`internal/repository/skemauji/`) |
| **test statik** | Test yang membaca kode sumber sebagai teks, bukan menjalankannya. Dipakai untuk aturan yang tidak dapat dijaga kompilator — misalnya "`Hapus` tidak boleh dipanggil di luar test" (`batasanpemakaian_test.go`) |
| **komponen** (React) | Fungsi yang mengembalikan tampilan |
| **`useState`** | Kotak penyimpan nilai di dalam komponen; mengganti isinya menggambar ulang tampilan |
| **`interface`** (TypeScript) | Bentuk data; hanya ada saat pemeriksaan tipe, tidak ikut ke browser |

## 5. Menjalankan — dari folder `APP_RNM\` ini

Bila `go`, `node`, atau `npm` tidak dikenal di terminal, tambahkan dulu ke PATH sesi
(rinciannya di `..\PANDUAN-RANTAI-ALAT.md` §1):

```powershell
$env:Path = 'C:\Program Files\Go\bin;C:\Program Files\nodejs;' + $env:Path
```

| Tujuan | Perintah | Yang diharapkan |
| --- | --- | --- |
| Uji backend tanpa Oracle | `go vet ./...` lalu `go test ./...` | `ok` di `internal/models` dan `internal/repository` |
| Bentuk tabel ke Oracle uji | `go run ./cmd/api -migrate` | menjalankan berkas di `internal/repository/migrations/` sekali masing-masing; aman diulang, dan **menolak** berjalan bila `IS_PEGA_PROD=true` |
| Uji backend **dengan** Oracle uji | `go test -tags=db ./internal/...` | perlu `ORACLE_DSN` + `ORACLE_SCHEMA` + `ORACLE_SKEMA_UJI=true`; tanpa `ORACLE_DSN` test **melewati** dengan pesan, bukan lulus diam-diam |
| Jalankan backend | `go run ./cmd/api` | log `http: mendengarkan di :8080` |
| Periksa tipe frontend | `cd frontend` lalu `npm run typecheck` | tidak mencetak galat |
| Uji frontend | `npm test` | `5 passed` |
| Jalankan frontend | `npm run dev` | buka `http://localhost:5173/` |
| Bangun frontend | `npm run build` | menjalankan `tsc` dulu, lalu Vite membuat `dist/` |

⛔ **`go test -tags=db` MENGHAPUS tabel di skema yang ditunjuk `ORACLE_SCHEMA`** — termasuk
`OS_AKSEPTASI_KLAIM_LIFE`, tabel datar warisan. Karena itu ia menolak berjalan kecuali
`ORACLE_SKEMA_UJI=true` **dan** `ORACLE_SCHEMA` bukan `POOLDATA`; penolakannya berupa **galat**,
bukan lewat. Tunjuk hanya skema uji kosong yang dibuat DBA khusus untuk itu, tidak pernah skema
mana pun yang memuat tabel warisan sungguhan — sekalipun di instance pengembangan. `-migrate`
tidak ikut dipagari: ia hanya `CREATE`, tidak pernah `DROP`.

Dengan `make`: `make test` · `make typecheck` · `make check` · `make test-db` · `make build`.

## 6. Tiga aturan yang paling sering ditegur reviewer

1. **Uang tidak pernah `float`, dan tidak pernah `Number` di JavaScript.** Oracle → Go → JSON →
   React membawanya sebagai **teks** desimal. Float membulatkan diam-diam. *(ADR-U-0003, ADR-U-0016)*
2. **Kode dan nomor tetap teks.** `"006"` bukan `6`. Mengubahnya menjadi bilangan menghilangkan
   nol di depan dan memecahkan penggolong. *(ADR-U-0022)*
3. **Akhiran baris berkas `.sql` dan `.go` selalu LF, tanpa BOM.** Alat Windows menulis CRLF dan BOM diam-diam; keduanya lolos ke teks SQL dan ditolak Oracle, sementara seluruh test tanpa basis data tetap hijau. Dijaga `.gitattributes` di akar dan dua test penjaga.
4. **Nol alamat host, kata sandi, atau nomor polis sungguhan di kode dan test.** Alamat dari
   env var; data uji berawalan `UJI-`. *(ADR-U-0004)*

Rujukan `ADR-U-nnnn` menunjuk `..\docs\adr\`; `ADR-D-<modul>-nnnn` ke `..\dastin\...\docs\adr\`;
`ADR-F-nnnn` ke `..\jefri\OUTPUT FIX\adr\`.

## 7. Dari mana pekerjaan datang

Setiap baris kode menunjuk satu **tiket**. Tiket Claim Life ada di
`..\.scratch\claim-life\issues\`; spec-nya di `..\.scratch\claim-life\spec.md`. Brief kerja
untuk sesi implementasi: `..\PROMPT-IMPLEMENTASI-GO-REACT.md`.
