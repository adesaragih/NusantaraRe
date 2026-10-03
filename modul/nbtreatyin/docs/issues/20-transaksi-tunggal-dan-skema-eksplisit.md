# 20: Transaksi tunggal dan skema eksplisit

**Status:** selesai — penahan tersisa hanya pihak luar: **K11** (skema uji Oracle — uji bertag `db` AC 45 dan AC 46 sudah ditulis, belum dijalankan) *(putaran 2, paket P11 04-10-2026; semula: sebagian — konsolidasi P10 04-10-2026; sebagian, implementasi 2026-10-03; awalnya ready-for-agent)*
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

- [ ] 🟡 **AC 45** — kegagalan menulis tabel anak mana pun **membatalkan seluruh** penyimpanan *(P11: `repository/penyimpanan_db_test.go` `TestGagalTulisTabelAnakMembatalkanInduk` — K11)*
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

## ⭐ Putaran 2 — paket P11 (04-10-2026): uji `db` AC 45

`backend/repository/penyimpanan_db_test.go` `TestGagalTulisTabelAnakMembatalkanInduk` (bertag `db`, **belum
dijalankan** — K11): kasus tersimpan (`PremiOgp` 100, `GroupPanel` 006, dua angsuran); penyimpanan kedua mengubah
induk (`PremiOgp` 999), quotation (`007`), dan anak — baris angsuran kedua ber-`InstallmentNo` `12345678901` lolos
pemeriksaan Go (bilangan bulat) tetapi ditolak **Oracle** (`INSTALLMENT_NO NUMBER(10)`, ORA-01438) SESUDAH induk
diperbarui, anak lama dihapus, quotation ditulis ulang, dan baris anak pertama disisipkan. Sesudah gagal, kolom
dibaca langsung: `PREMI_OGP` 100, `GROUP_PANEL` 006, angsuran lama `1:1 2:2` — tak satu pun tersisa dari
penyimpanan yang gagal; galatnya bukan `ErrPermintaanTidakSah` (bukan galat Go).

Status: **selesai** — sisa penahan hanya K11 (AC 45, 46 🟡).
