# Peta asal dokumen `facout/docs/`

Dipindah 30 September 2026 dari `jefri/OUTPUT FIX/` (`mv`, **isi tidak diubah**); tautan relatif di
dalamnya masih menunjuk tata letak lama. Peta tautan lintas siklus: `../../nbfacin/docs/PETA-ASAL.md`.

`[terverifikasi]` Isi ke-24 berkas sama dengan salinan asalnya di `D:\migrasi\RNM\OUTPUT\` setelah akhir
baris diseragamkan (`cmp` sesudah `tr -d '\r'`: 24 sama, 0 beda).

## Berkas di folder ini

| Asal (`jefri/OUTPUT FIX/…`) | Sekarang (`facout/docs/…`) |
| --- | --- |
| `05-tickets/facout/00-INDEKS-FACOUT.md` + `F01…F14-*.md` (15) | **`issues/`** — nama sama |
| `09-facout/01-temuan-dan-rancangan-facout.md` | `09-facout/` — nama sama |
| `10-audit/01…07-*.md`, `10-audit/09-*.md` (8) | `10-audit/` — nama sama. Masing-masing menyebut dirinya sumber angka tiket F-xx atau K-057…K-062 |

Audit `08`, `10`, `11` **lintas siklus** — tetap di `jefri/OUTPUT FIX/10-audit/`.

## Tautan di dalam dokumen yang menunjuk ke luar folder ini

| Tautan di dokumen | Sekarang di |
| --- | --- |
| `..\01…16-*.md` (tiket NB), `..\edm\E21…` | `APP_RNM/modul/nbfacin/docs/issues/`, `APP_RNM/modul/endorsmentfacin/docs/issues/` |
| `..\..\00-KEPUTUSAN-WORK-OWNER.md` (K-053…K-062) | `APP_RNM/modul/nbfacin/docs/00-KEPUTUSAN-WORK-OWNER.md` |
| `..\..\10-audit\08…`, `10…`, `11…` | `jefri/OUTPUT FIX/10-audit/` |

⛔ **Seluruh 14 tiket ter-block tiket `nbfacin`** (NB-01, 03, 08, 09, 11, 15, 16) dan F11 juga oleh
`endorsmentfacin` E21. `modul/facout` tidak boleh mengimpor `modul/nbfacin` — registry predikat
(K-050: satu registry, bukan dua), tangga akseptasi, dan hitung premi datang lewat `inti/backend/kontrak`.
