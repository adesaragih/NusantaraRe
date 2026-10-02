# Peta asal dokumen `nbfacin/docs/`

Dipindah 30 September 2026 dari `jefri/OUTPUT FIX/` dengan `git mv`. **Isi setiap berkas tidak diubah.**
Karena itu tautan relatif di dalam dokumen masih menunjuk tata letak lama — pakai tabel di bawah
untuk menemukan tujuannya.

`[terverifikasi]` Isi ke-85 berkas yang dipindah (tiga modul) sama dengan salinan asalnya di
`D:\migrasi\RNM\OUTPUT\`; bedanya hanya akhir baris (CRLF di repo, LF di asal). Diuji dua cara:
`cmp` mentah (85 beda — semuanya CRLF) dan `cmp` setelah `tr -d '\r'` (85 sama).

## Berkas di folder ini

| Asal (`jefri/OUTPUT FIX/…`) | Sekarang (`nbfacin/docs/…`) |
| --- | --- |
| `01-activity/01-inventaris-activity-nb.md` | `01-activity/` — nama sama |
| `02-layar/01-skema-field-nb.md` | `02-layar/` — nama sama |
| `03-celah/01…03-*.md` (3) | `03-celah/` — nama sama |
| `03-keputusan/RINGKASAN-GRILLING.md` | `03-keputusan/` — nama sama |
| `04-kuesioner/01…04-*.md` (4) | `04-kuesioner/` — nama sama |
| `04-spec/01-modul-go.md`, `02-model-data.md`, `03-spec-modul-terverifikasi.md` | `04-spec/` — nama sama |
| `05-tickets/00-INDEKS.md` + `05-tickets/01…16-*.md` (17) | **`issues/`** — nama sama |

## Tautan di dalam dokumen yang menunjuk ke luar folder ini

| Tautan di dokumen | Sekarang di |
| --- | --- |
| `rnw/…`, `..\rnw\R0x…`, `04-spec\04-spec-rnw.md`, `06-rnw\` | `APP_RNM/modul/rnwfacin/docs/` (`issues/`, `04-spec/`, `06-rnw/`) |
| `edm/…`, `..\edm\E0x…`, `04-spec\05…10-spec-edm-*`, `07-edm\` | `APP_RNM/modul/endorsmentfacin/docs/` (`issues/`, `04-spec/`, `07-edm/`) |
| `facout/…`, `09-facout\`, `10-audit\01…07`, `10-audit\09` | `jefri/OUTPUT FIX/_usulan-modul-facout/docs/` (`issues/`, `09-facout/`, `10-audit/`) — modul baru `facout`, keputusan 30-09-2026; **diparkir** sampai PR tim inti (bab 5), lalu pindah ke `APP_RNM/modul/facout/` |
| `../adr/`, ADR-0001…0007 | tetap di `jefri/OUTPUT FIX/adr/` — dirujuk kode sebagai **`ADR-F-000n`** (mis. `inti/backend/uang`). ⚠️ Nomornya **bukan** `docs/bersama/adr/000n` |
| `steering/GLOSARIUM.md`, `steering/PANDUAN-KERJA.md` | tetap di `jefri/OUTPUT FIX/steering/` |
| `08-flat\`, `04-spec\11-spec-pemuatan-data-lama.md`, `10-audit\08`, `10`, `11` | tetap di `jefri/OUTPUT FIX/` (lintas siklus) |
| `00-KEPUTUSAN-WORK-OWNER.md` (K-001…), `00-RINGKASAN-NB.md`, `_EKSTRAKSI-PEGA-SELAGI-HIDUP.md`, `_PAKET-PERMINTAAN-DBA-IT-PRODUCT.md`, `_ARSIP-lintas-siklus/`, dan 7 berkas `_*.md` lain | **folder ini** (`nbfacin/docs/`) — **disalin** 30-09-2026 dari `D:\migrasi\RNM\OUTPUT\` (26 berkas, `cmp` 26/26 identik), karena tidak pernah ikut ke `jefri/OUTPUT FIX/` |

## Yang perlu diingat sebelum mengerjakan tiket

- ⭐ Keputusan 30-09-2026 (register, modul `facout`, tiket NB-01/02 mengikuti `inti/backend/uang`):
  `KEPUTUSAN-30-09-2026.md`.

- Tiket ditulis untuk tata letak **sebelum** bentuk B: `pkg/money`, `pkg/ratio`, `internal/rules`,
  `services/…`, dan menyatakan "nol berkas `.go`". Di repo ini sudah ada `inti/backend/uang`
  (`Money`, `Ratio`) yang dipakai `claimlife` dan `premiumlistlife`, dengan perilaku berbeda dari
  tiket 01/02 (lintas mata uang → `ErrMataUangBerbeda`, bukan `panic`; tanpa keadaan `Unknown`).
- Jalur korpus di dokumen adalah `D:\migrasi\RNM\NB FacIn\`; `CLAUDE.md` menunjuk `D:\XML\RNM_BRD\`.
- `backend/modul.go` **belum** dibuat: penjaga `TestMenuDimigrasiSamaDenganModulBackend`
  mewajibkan setiap modul ber-`modul.go` sudah menyalakan menunya. Ia lahir bersama layar pertama
  (`docs/bersama/PANDUAN-TIM-PER-MODUL.md` bab 4).
