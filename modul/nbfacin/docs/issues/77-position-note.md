# 77: T_WORK_POLIS.POSITION_NOTE — workbasket dan tahap dipisah

Status: selesai di backend (belum di-commit; migrasi 254 ditulis, belum dijalankan; konversi data ditulis, belum dijalankan).

## Asal

Work owner 12-10-2026 (lewat sesi 9d): "kalo tambahin kolom baru untuk positionnote gimana? jadi di position itu tetap offer,
policy, binding" → "Ya, kerjakan"; bug DEV NB-151039 ("itu knapa positionnya offer"). Pega punya dua properti: `PositionNote`
(workbasket = kotak masuk) dan `IsCedingConfirm` (tahap Offer / Binding / Policy); sistem baru memakai satu kolom
`T_WORK_POLIS.POSITION` untuk keduanya ('Offer' = antrean marketing tahap Offer, butir 77).

Kepemilikan: T_WORK_POLIS milik premiumlistlife. **Opsi 1 dipilih work owner**: migrasi 254 di nbfacin + SATU baris
`POSITION_NOTE` di `modul/premiumlistlife/docs/STRUKTUR-TABEL-PREMIUMLIST-LIFE.md` (izin eksplisit work owner 12-10-2026 untuk
baris itu saja; work owner yang mengabari pemegang PremiumList Life). Kolom juga dipakai endorsmentfacin (LINI FACEDM, migrasi
endorsmentfacin 263) — nama / lebar jangan diubah tanpa mengabari mereka.

## Sensus pemakai `POSITION` `[terverifikasi]` — dua cara, cocok: 14 baris

(A) pengurai literal string Go di luar komentar, (B) grep kata utuh `POSITION` (bukan JOB_POSITION / POSITION_NOTE) tanpa baris
komentar — keduanya 14 baris (angka "28+" permintaan awal keliru, diakui sesi 9d):

| Berkas | Baris | Pemakaian |
| --- | --- | --- |
| repository/kasus.go | 60, 152 / 97 | baca kasus / sqlTahapKasus |
| repository/submit.go | 116 / 206-207 | sqlKonteksSubmit / sqlPosisiSubmit (SET + penjaga 409) |
| repository/kotakmasuk.go | 91, 95 / 102 | saringan antrean + team / GROUP BY |
| repository/scoring.go | 73 | kepala Scoring (tahap marketing) |
| repository/binding.go | 259 | penjaga tombol Binding |
| repository/casenb.go | 76 | sisip case |
| services/loader/amandemen.go, aturan.go | 54, 179 | pemetaan loader data lama — **tidak diubah** |

## Rancangan (diterima sesi 9d)

- **Kosakata services TIDAK berubah** `[keputusan agent]`. Repository menerjemahkan penyimpanan (`repository/posisinote.go`):
  - tulis: `POSITION_NOTE` = PosisiBaru, kecuali 'Offer' → 'ReasFacInMarketing' (PositionNote awal Pega) dan 'Decline' → NULL
    (keputusan work owner: Decline mengosongkan antrean); `POSITION` = tahap dari FLAG_ONGOING_POLICY (`TahapDariFlag`: '0' Offer,
    '2' Binding, lainnya Policy), **disinkronkan di akhir transaksi** Submit (`sqlSinkronTahap`, sesudah FlagBaru dan Transition38
    tulisBinding yang men-set flag 2 di pernyataan terpisah); buat case langsung 'Offer' + 'ReasFacInMarketing'.
  - baca, posisi efektif (`sqlPosisiEfektif`, padanan Go `PosisiEfektif`):
    1. NOTE NULL + Resolved + POSITION tahap (Decline biner baru) → 'Decline' bila Resolved-Decline, selain itu NULL;
    2. x = NVL(NOTE, POSITION) — baris belum dikonversi dibaca POSITION lama; x 'ReasFacInMarketing' + flag 0 → 'Offer'.
  - Setara aturan kontrak "NVL(NOTE, POSITION lama) hanya untuk case non-Resolved", dinyatakan dalam kosakata lama. Normalisasi
    'ReasFacInMarketing' + flag 0 → 'Offer' aman: services menyamakan keduanya (posisiMarketing, posisiMarketingKotakMasuk dijumlah,
    wewenang tahap marketing = workbasketTahap, tahap Binding = 'ReasFacInMarketing' + flag 2) — dan membuat pembacaan SAMA sebelum /
    sesudah konversi.
- **Peralihan**: kolom dideteksi sekali per tabel (`adaCatatanPosisi`, `SELECT POSITION_NOTE … WHERE ROWNUM < 1`; ORA-00904 di-cache,
  dicoba ulang tiap 1 menit). **Tanpa 254 SQL persis seperti sebelum tiket 77** — teks lama dikunci uji yang sudah ada.
- **Riwayat** (bug POSISI / PUTARAN NB-151039): `posisiRiwayatFlag` = `repository.TahapDariFlag` (satu fungsi tahap, butir 6).

## Konversi data — `docs/alat/konversi_position_note_254_pooldata.sql` (siap-jalan POOLDATA) dan `konversi_position_note_254.sql` ({skema})

Sesudah 254; TIDAK wajib (baris lama dibaca lewat posisi efektif). Satu UPDATE tanpa COMMIT:
`POSITION_NOTE = CASE POSITION WHEN 'Offer' THEN 'ReasFacInMarketing' WHEN 'Decline' THEN NULL ELSE POSITION END`,
`POSITION = tahap dari flag`. Cakupan: LINI 'FAC' ID 'NB-%' (FACEDM milik endorsmentfacin), NOTE masih NULL, bukan baris yang sudah
berbentuk baru — aman diulang. `[keputusan agent]` Resolved-Completed IKUT dikonversi (NOTE = workbasket terakhir, seperti biner
baru saat Selesai Policy). **Pemeriksa**: sebaran posisi efektif seluruh case NB FacIn di BAGIAN 0 dan BAGIAN 2 wajib identik.

## Uji

`TestTahapDanCatatanPosisi`, `TestPosisiEfektifBarisCampuran` (12 baris lama / baru / Decline), `TestSQLCatatanPosisi` (varian
hanya mengganti kolom posisi; bind tidak bergeser; sinkron; sisip case; uji kolom lolos PeriksaSQL), `TestKotakMasukCatatanPosisi`,
`TestTulisSubmitCatatanPosisi`, `TestPembacaMemakaiDeteksi` (9 pemanggilan deteksi di 6 berkas), `TestJalurTulisPulangPergi`
(12 jalur tulis), `TestKonversiMempertahankanPosisiEfektif` (11 kombinasi lama: efektif sebelum = sesudah, idempoten, teks skrip),
`TestPosisiPutaranRiwayat`. Mutasi 12 / 12 mati (termasuk penjaga DDL ↔ STRUKTUR premiumlistlife sebagai uji instrumen).
⚠️ Tanpa driver Oracle tiruan, perilaku SQL diuji lewat teks + padanan Go; kesetaraan SQL ↔ Go dijaga uji teks yang sama.

## Uji premiumlistlife yang terdampak (keputusan work owner 12-10-2026)

Kolom 254 membuat dua uji milik premiumlistlife gagal: `TestMigrasi059WorkPolisSeragamDenganWorkClaim` (bentuk akhir T_WORK_POLIS
dikunci 10 kolom) dan `TestMigrasi050Sampai056TipeNullFKIndexSesuaiStruktur` (jumlah kolom STRUKTUR di-hard-code 244, kini 245).
**Dibiarkan untuk pemilik modul premiumlistlife** — berkas modul itu tidak diubah nbfacin (izin hanya untuk satu baris STRUKTUR).
Verifikasi nbfacin memakai `./modul/nbfacin/...` saja.
