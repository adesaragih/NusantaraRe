# FINDING-006 — Kurs yang tidak ditemukan dikembalikan sebagai 1

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.
> Berlaku hanya untuk keadaan sistem pada kedua ekspor itu.

**Jenis**: laporan kondisi sistem lama
**Status**: TERBUKTI DARI DDL — besarannya belum terukur
**Ruang lingkup bukti**: `pengetahuan/DDL_Script_ClaimNonProp.xls`
**Tanggal**: 18 September 2026

Dipisahkan dari `FINDING-001` atas permintaan Round 6. Uraian penuh rantai buktinya tetap di `FINDING-001` bagian 7; dokumen ini memuat temuannya sendiri supaya dapat dirujuk terpisah.

---

## 1. Mekanismenya

`POOLDATA.GETCURRENCYSTANDARD` berakhir dengan:

```sql
EXCEPTION
   WHEN NO_DATA_FOUND
   THEN
      RETURN 1;
END;
```

Bila tidak ada baris kurs untuk suatu mata uang, fungsi mengembalikan **`1`** — yang secara aritmetika berarti konversi satu banding satu.

## 2. Akibatnya

Nilai 1.000 USD diperlakukan sebagai 1.000 IDR. Tanpa galat, tanpa pesan, tanpa penanda.

Ini pilihan yang disengaja — ada penangan eksepsi khusus untuknya, bukan kelalaian. Konsekuensinya tetap sama: **kegagalan pencarian kurs tidak dapat dibedakan dari kurs yang memang bernilai 1.**

Rupiah terhadap rupiah memang berkurs 1. Jadi baris hasil tidak menyimpan informasi apa pun tentang apakah konversinya terjadi atau gagal.

## 3. Pola yang sama muncul tiga kali di sesi ini

Ini bukan kejadian tunggal, melainkan satu kebiasaan yang berulang: **kegagalan yang menyamar menjadi hasil yang sah.**

| Mekanisme | Kegagalan | Terlihat sebagai |
|---|---|---|
| `GETCURRENCYSTANDARD` | kurs tidak ada | kurs bernilai 1 |
| Database link ke HRD | tautan putus | daftar pengguna kosong |
| `pyStepsPreCondParamsWhen` yang tidak cocok | kondisi gagal | langkah dilewati tanpa pesan |

Ketiganya tidak menghasilkan galat. Ketiganya menghasilkan jawaban yang bentuknya benar dan isinya salah.

## 4. Batas klaim

1. **Belum diketahui apakah ada mata uang tanpa baris kurs di produksi.** Bila `m_currencystandard` lengkap untuk semua mata uang yang dipakai, jalur `NO_DATA_FOUND` tidak pernah tersentuh dan temuan ini nol dampaknya.
2. **Tidak dapat diukur dari hasilnya.** Karena kurs 1 sah untuk IDR, mencari "baris berkurs 1" akan mencampur kegagalan dengan kasus normal. Pengukurannya harus dari sisi lain: cari mata uang **bukan IDR** yang kursnya tersimpan `1` — itu hampir pasti hasil jalur kegagalan. → **REQ-022**.
3. Berlaku juga untuk `result_cache` pada fungsi itu: hasil `1` yang pernah dikembalikan dapat tersimpan di cache.

## 5. Hubungan dengan sistem baru

Sudah diputuskan terpisah di **ADR-0014**: konversi tanpa kurs tidak menghasilkan angka. Nilai IDR dibiarkan kosong dan ditandai menunggu kurs, sehingga tidak ada angka salah yang ikut ke perhitungan lain.
