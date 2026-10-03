# 20: Transaksi tunggal dan skema eksplisit

**Status:** sebagian — dapat dikerjakan: uji bertag `db` kegagalan menulis **tabel anak** membatalkan induk (AC 45) belum ditulis; penahan pihak luar: K11 (AC 45, 46) *(putaran 2, konsolidasi P10 04-10-2026 — rincian `docs/HASIL-IMPLEMENTASI.md` bab 9; semula: sebagian, implementasi 2026-10-03; awalnya ready-for-agent)*
**Blocked by:** **16** · **19**
**Menutup:** NB AC **45–48** *(4 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-32..ID-35

## Hasil & nilai pengguna

Satu polis tersimpan **seluruhnya atau tidak sama sekali**. Sistem lama meng-`COMMIT` di dalam
empat program basis data yang berbeda, sehingga kegagalan di tengah meninggalkan data separuh yang
permanen.

⚠️ `[penyimpangan sadar]` **Ini lebih ketat daripada sistem lama, dan disengaja.**

## Yang dibangun

Pembungkusan seluruh urutan penyimpanan satu polis — induk dan seluruh anaknya — dalam **satu
transaksi**. Dan dua aturan penulisan query:

| Aturan | Sebab |
| --- | --- |
| skema basis data ditulis **eksplisit** pada setiap pernyataan | sistem lama tidak konsisten; satu skema kedua terbukti disentuh |
| ⛔ **nol pemanggilan program tersimpan** dari aplikasi | mandat proyek; logikanya ditulis ulang |

## Batas — yang TIDAK termasuk

⛔ Transaksi yang mencakup tabel proyeksi selisih — tiket **26**.
⛔ Pemuatan massal — tiket **22**.

## Cara mengujinya

Lewat seam `repository`. Uji utama: kegagalan disuntikkan pada penulisan tabel anak terakhir —
tidak boleh ada satu baris pun tertinggal, termasuk baris induknya.

## Acceptance criteria

- [ ] 🟡 **AC 45** — kegagalan menulis tabel anak mana pun **membatalkan seluruh** penyimpanan
- [ ] 🟡 **AC 46** — seluruh penyimpanan satu polis terjadi dalam **satu transaksi**
- [x] **AC 47** — setiap pernyataan menyebut skema **eksplisit**
- [x] **AC 48** — penyimpanan **tidak memanggil** satu pun program tersimpan

## ⭐ Putaran 2 — paket penyimpanan (03-10-2026)

Dasar: PROMPT-NB-TREATY-IN-PUTARAN-2 bab 0 butir 11–12, bab 2 K4/K16/K17; rincian kolom `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`.

- `HISTORYAKSEPTASIPRODUCTION` (K4) ditulis di transaksi tunggal submit, skema eksplisit, nol `COMMIT`
  (`TestSQLRiwayatProduksiMengikutiInsertViewSuggest`), nol procedure — sama dengan AC 47, 48.
- Urutan simpan disesuaikan FK diagram: anak dihapus **sebelum** `T_POLIS_QUOTATION` ditulis ulang (ceding
  menunjuk quotation), lalu anak disisip — tetap satu transaksi.
- AC 45 dan 46 tetap 🟡 (uji `-tags db` belum dijalankan — K11).
