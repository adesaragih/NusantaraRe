# Prompt — implement NB saja, secepatnya (sesi `nusantarare-c3`)

> Perintah work owner, 01-10-2026. **Menggantikan bagian "Jatah" dan "Urutan"** di
> `PROMPT-LANJUT-IMPLEMENT-TANPA-TANYA-NB-RNW.md`. Bagian "Aturan kerja" (butir 1–10) berkas itu tetap
> berlaku penuh: langsung implement, tanpa bertanya kecuali 4 hal, korpus diam → `belum terverifikasi`,
> tiket kurang → tulis dari XML, CLAUDE.md §4/§4a/§7, jangan commit, satu laporan di akhir.

## Jatah: HANYA `modul/nbfacin`

- ⛔ **Berhenti mengerjakan `modul/rnwfacin`.** Biarkan apa adanya; pastikan baris `Status:` R01–R08
  mencerminkan keadaan terakhir, lalu jangan disentuh lagi.
- ⛔ Jangan menyentuh `modul/endorsmentfacin` atau modul lain.
- `inti/backend/kontrak/facin.go` dan generator `rules/bangkit`: hanya bila perlu untuk NB atau bila ada
  bug. Perubahan kontrak tetap dicatat sebagai butir yang perlu disetujui tim inti.

## Urutan

1. **Koreksi hasil verifikasi** (pesan dari sesi `nusantarare-0f`, 01-10-2026) — catat di register bab
   Ralat, sebut cara mana yang keliru.
2. **Tuntaskan tiket NB yang masih bisa dikerjakan tanpa pihak luar** (mis. sisa tiket 07).
3. **Jadikan NB modul yang berjalan** — ikuti `docs/bersama/PANDUAN-TIM-PER-MODUL.md` bab 4,
   `modul/_templat/`, dan pola modul yang sudah dimigrasi (`claimlife`, `premiumlistlife`):
   `backend/modul.go` (`Pendaftaran()`), `models/`, `repository/` (termasuk pemuat yang dibutuhkan
   tiket 17), `handlers/`, lalu `frontend/` (`menu.ts`, `rute.tsx`, `pages/`, `labels.ts`, `api.ts`).
   Tiket untuk lapisan ini belum ada → tulis ringkas dulu dari spec/XML (aturan 5).
   - **Skema Oracle hanya dari DDL** yang tersedia (`D:\migrasi\RNM\DDL\`), jangan menebak kolom/tipe.
     Kolom yang tidak ada di DDL → `belum terverifikasi`, tanyakan di laporan akhir.
   - Berkas migrasi boleh ditulis di rentang `180-219`; **jangan dijalankan** ke Oracle nyata.
   - ⛔ Jangan mengubah nilai rentang migrasi atau slot menu di `MODUL.md` (dibaca penjaga; perubahan
     lewat tim inti).
   - Penjaga (`go test ./inti/backend/penjaga/...`, `inti/frontend/lapisan.guard.test.ts`) dan tes
     frontend harus hijau.
4. Tiap kelompok selesai: full suite Go + tes frontend, code review dua sumbu, perbaiki temuan.

## Tetap tertahan — lewati, sebut di laporan

Tiket 04 (data DBA), kasus MBU #5 (pemilik Pega), butir yang menunggu work owner (A31–A36, A38, A39,
A42, A43, A45), persetujuan tim inti atas `facin.go`.
