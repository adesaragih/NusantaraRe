# Peta asal dokumen `rnwfacin/docs/`

Dipindah 30 September 2026 dari `jefri/OUTPUT FIX/` dengan `git mv`. **Isi setiap berkas tidak diubah**;
tautan relatif di dalamnya masih menunjuk tata letak lama. Peta lengkap tautan keluar — ADR, steering,
register keputusan, Fac Out — ada di `../../nbfacin/docs/PETA-ASAL.md`.

## Berkas di folder ini

| Asal (`jefri/OUTPUT FIX/…`) | Sekarang (`rnwfacin/docs/…`) |
| --- | --- |
| `04-spec/04-spec-rnw.md` | `04-spec/` — nama sama |
| `06-rnw/01…05-*.md` (5) | `06-rnw/` — nama sama |
| `05-tickets/rnw/00-INDEKS-RNW.md` + `R01…R08-*.md` (9) | **`issues/`** — nama sama |

## Tautan ke New Business

| Tautan di dokumen | Sekarang di |
| --- | --- |
| `..\..\04-spec\…` (spec NB), `NB-08`, `NB-11`, `NB-16`, `..\0x-…md` | `APP_RNM/modul/nbfacin/docs/` (`04-spec/`, `issues/`) |

⛔ **Tiket di sini ter-block tiket `nbfacin`** — R01 menunggu NB-08 dan NB-11, R08 menunggu NB-16.
Karena `modul/rnwfacin` tidak boleh mengimpor `modul/nbfacin`, mesin NB yang dipakai ulang
(registry predikat, tangga akseptasi, hitung premi) harus sampai ke modul ini lewat
`inti/backend/kontrak` (`docs/bersama/PANDUAN-TIM-PER-MODUL.md` bab 7).
