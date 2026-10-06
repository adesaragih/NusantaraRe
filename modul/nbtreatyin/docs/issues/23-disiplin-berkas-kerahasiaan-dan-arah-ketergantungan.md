# 23: Disiplin berkas, kerahasiaan, dan arah ketergantungan

**Status:** selesai — AC 63 📄 (migrasi dijalankan work owner), AC 64 ✅ lawan diagram F10 *(putaran 2, konsolidasi P10 04-10-2026 — rincian `docs/HASIL-IMPLEMENTASI.md` bab 9; semula: sebagian, implementasi 2026-10-03; awalnya ready-for-agent)*
**Blocked by:** — *(dapat mulai segera)*
**Menutup:** NB AC **60–66** *(7 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-1 · ID-21 · ID-25 · Bab 9

## Hasil & nilai pengguna

Pekerjaan penyimpanan tidak membocorkan data nasabah ke dalam repositori, dan susunan lapisannya
tidak terbalik.

⭐ **Tiket ini tanpa penahan dan dapat dikerjakan sejajar dengan tiket mana pun** — sebagian besar
isinya penjaga otomatis, bukan fitur.

## Yang dibangun

Penjaga otomatis yang gagal bila dilanggar:

| Penjaga | Yang dicegahnya |
| --- | --- |
| nol nama orang di berkas proyek | nilai produksi tersalin ke spec, test, atau rancangan |
| nol nomor polis harfiah | pengenal nasabah masuk repositori |
| test berjalan di atas **data buatan** | cuplikan produksi menjadi fixture |
| perubahan skema menuntut persetujuan manusia | perubahan bentuk tabel berjalan otomatis |
| arah ketergantungan tidak terbalik | lapisan penyimpanan memanggil lapisan layanan |

Ditambah dua pernyataan bentuk yang diuji langsung: cacah medan tabel data umum, dan penerimaan
bahwa bentuk gabungan ceding **boleh tidak sinkron** dengan barisnya.

## Batas — yang TIDAK termasuk

⛔ Isi tabel — tiket **19**.
⚠️ Cacah medan yang diuji AC 64 adalah **batas bawah**, bukan total.

## Cara mengujinya

Sebagian lewat seam `repository`, sebagian lewat pemeriksaan repositori otomatis.

## Acceptance criteria

- [x] **AC 60** — nol nama orang tersalin ke berkas rancangan, spec, maupun test
- [x] **AC 61** — nol nomor polis ditulis apa adanya
- [x] **AC 62** — test berjalan di atas **data buatan**
- [x] 📄 **AC 63** — perubahan skema memerlukan **persetujuan manusia** *(P10: butir dokumen — migrasi dijalankan work owner, MODUL.md)*
- [x] ✅ **AC 64** — tabel data umum menampung cacah medan yang ditetapkan spec *(P10: ✅ lawan diagram F10 — `docs/PERBANDINGAN-KOLOM-DIAGRAM.md` bab 1, 10)*
- [x] **AC 65** — bentuk gabungan ceding **boleh tidak sinkron** dengan barisnya, tanpa penjaga
- [x] **AC 66** — arah ketergantungan **tidak pernah dibalik**
