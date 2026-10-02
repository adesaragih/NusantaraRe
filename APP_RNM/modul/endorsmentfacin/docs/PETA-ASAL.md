# Peta asal dokumen `endorsmentfacin/docs/`

Dipindah 30 September 2026 dari `jefri/OUTPUT FIX/` dengan `git mv`. **Isi setiap berkas tidak diubah**;
tautan relatif di dalamnya masih menunjuk tata letak lama. Peta lengkap tautan keluar — ADR, steering,
register keputusan, Fac Out — ada di `../../nbfacin/docs/PETA-ASAL.md`.

## Berkas di folder ini

| Asal (`jefri/OUTPUT FIX/…`) | Sekarang (`endorsmentfacin/docs/…`) |
| --- | --- |
| `04-spec/05…10-spec-edm-*.md` (6) | `04-spec/` — nama sama |
| `07-edm/01…11-*.md` (10) | `07-edm/` — nama sama |
| `04-kuesioner/_DITUNDA-fase-endorsement.md` | `04-kuesioner/` — nama sama |
| `05-tickets/edm/00-INDEKS-EDM.md` + `E01…E22-*.md` (23) | **`issues/`** — nama sama |

## Tautan ke siklus lain

| Tautan di dokumen | Sekarang di |
| --- | --- |
| `..\01-…md` … `..\16-…md` (tiket NB), spec NB `04-spec\01…03` | `APP_RNM/modul/nbfacin/docs/` (`issues/`, `04-spec/`) |
| `..\rnw\R0x…` | `APP_RNM/modul/rnwfacin/docs/issues/` |
| `..\facout\F1x…` | `jefri/OUTPUT FIX/_usulan-modul-facout/docs/issues/` (diparkir; kelak `APP_RNM/modul/facout/`) |

⛔ **Seluruh 22 tiket di sini ter-block tiket `nbfacin`** (lihat `issues/00-INDEKS-EDM.md` §0). Seam 1
(`rules.Eval`) dan Seam 3 (`premium.Calculate`) milik NB; modul ini tidak boleh mengimpor
`modul/nbfacin` — keduanya harus datang lewat `inti/backend/kontrak`.

📌 `docs/bersama/adr/0041-endorsemen-berbagi-tabel-dengan-new-business.md` relevan untuk modul ini.
