# 01: Kerangka aplikasi + seam API + tracer "buka satu klaim Life"

**Status:** ready-for-agent

**Blocked by:** None (can start immediately)

## Hasil & nilai pengguna

Seorang pengguna dapat membuka satu klaim Life dan melihat baris-baris `AdjustmentList`-nya di
layar. Itu perilaku yang tipis, tetapi ia menembus **seluruh lapisan** — React → HTTP → handler →
service → repository → Oracle — sehingga sekaligus **menegakkan seam** yang dipakai dua belas tiket
sesudahnya.

Nilainya bukan fiturnya, melainkan **jalurnya**: setelah tiket ini selesai, setiap tiket berikutnya
menambah perilaku di jalur yang sudah terbukti hidup, bukan membangun jalur baru.

## Area codebase

Seluruhnya **di dalam `OUTPUT_HASIL_RNM/`** — hari ini belum ada satu pun berkas kode di sana.

| Bagian | Yang dibuat |
| --- | --- |
| akar | `go.mod`, `Makefile`, `.env.example` |
| `cmd/api` | entry point HTTP |
| `internal/config` | pembacaan env var + koneksi Oracle |
| `internal/models` | agregat klaim Life + koleksi baris adjustment |
| `internal/repository` | pembacaan klaim + baris dari `POOLDATA` |
| `internal/services` | perakitan agregat |
| `internal/handlers` | satu endpoint REST baca-saja |
| `frontend/` | Vite + React; satu halaman yang menampilkan klaim dan barisnya |

Arah dependency **`handlers` → `services` → `repository`** (`CLAUDE.md` §5). Tidak boleh terbalik,
tidak boleh memotong lapisan.

## Rule Pega sumber

Tiket ini **tidak meniru perilaku bisnis** — ia hanya membaca. Bentuk data diambil dari kolom yang
sudah terbaca:

| Rule | Identitas | Dipakai untuk |
| --- | --- | --- |
| `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `[terverifikasi]` 55 nama kolom `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` terbaca langsung dari blok PL/SQL-nya |
| `Claim Life/Section/AdjustmentDetail_Section.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION` | `[terverifikasi]` bentuk baris `AdjustmentList` dan kolom yang ditampilkan |

## ADR terkait

**ADR-0003** (uang non-float — presisi desimal eksak; kolom Oracle `NUMBER` tanpa presisi),
**ADR-0013** (koneksi database dari env var; **alamat layanan keluar TIDAK** — ia di-lookup runtime
dari `M_LINK_SERVICE`), **ADR-0011** (agregat memuat **koleksi baris**, bukan satu status).

## Acceptance criteria

- [ ] `GET` satu klaim Life mengembalikan klaim beserta **seluruh** baris `AdjustmentList`-nya,
      bukan hanya baris terakhir.
- [ ] Halaman React menampilkan klaim dan daftar barisnya, masing-masing dengan statusnya.
- [ ] Status ditampilkan sebagai kata (Outstanding / Aksep / Ditolak), **bukan** sebagai nama field
      `STS_REJECT` dan bukan sebagai angka. *(AC 26 spec)*
- [ ] Tidak ada nilai uang yang direpresentasikan sebagai *binary floating point* di lapisan mana
      pun maupun di JSON respons. *(AC 22 spec)*
- [ ] Tidak ada host, endpoint, atau kredensial sebagai literal di kode — semuanya env var.
- [ ] Ada satu test ujung-ke-ujung yang menggerakkan sistem lewat HTTP dan memeriksa hasilnya lewat
      HTTP, terhadap skema uji Oracle — bukan mock repository.
- [ ] `handlers` tidak memanggil `repository` langsung; `repository` tidak memuat business logic.

## Skema uji

`[terbuka]` **OQ-001** (pemilik **DBA**) — tidak ada DDL di korpus, sehingga **tipe dan presisi**
kolom tidak diketahui. Untuk tiket ini dipakai **skema uji provisional**: nama kolom dari rule di
atas, tipe numerik berpresisi longgar, tipe teks berpanjang longgar.

Ini **sah untuk skema uji** dan **tidak sah untuk DDL produksi maupun migrasi** — keduanya tetap
menunggu OQ-001 dan ditangani di tiket **13**.

## Perintah verifikasi

Toolchain belum ada; tiket ini yang membuatnya. Setelah selesai, ketiga perintah berikut **harus
berjalan hijau** dan dicatat di `Makefile`:

```
go build ./...
go test ./internal/...
cd frontend && npm test
```

Tambahkan target gabungan `make check` yang menjalankan ketiganya.
