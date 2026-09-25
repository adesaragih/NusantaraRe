# 23: Disiplin berkas, kerahasiaan, dan arah ketergantungan

**Status:** ready-for-agent
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

- [ ] **AC 60** — nol nama orang tersalin ke berkas rancangan, spec, maupun test
- [ ] **AC 61** — nol nomor polis ditulis apa adanya
- [ ] **AC 62** — test berjalan di atas **data buatan**
- [ ] **AC 63** — perubahan skema memerlukan **persetujuan manusia**
- [ ] **AC 64** — tabel data umum menampung cacah medan yang ditetapkan spec
- [ ] **AC 65** — bentuk gabungan ceding **boleh tidak sinkron** dengan barisnya, tanpa penjaga
- [ ] **AC 66** — arah ketergantungan **tidak pernah dibalik**
