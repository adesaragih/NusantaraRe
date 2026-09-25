# FINDING-001 — Ambang kewenangan Komite dibandingkan tanpa konversi mata uang

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

**Jenis**: laporan kondisi sistem lama (bukan keputusan migrasi)
**Status**: BELUM TERBUKTI — menunggu REQ-011
**Ruang lingkup bukti**: `D:\XML_NURE\Claim Non Prop` saja
**Tanggal**: 18 September 2026

Dokumen ini memaparkan rantai perhitungan apa adanya. Ia **tidak** menyatakan adanya pelanggaran, kelalaian, atau kerugian, dan **tidak** menunjuk pihak mana pun. Apakah kondisi ini benar-benar berdampak hanya dapat dijawab oleh data produksi (REQ-011), bukan oleh XML.

---

## 1. Rantai bukti

Seluruhnya berada dalam satu berkas: `Activity\CreateChildKomiteCNP_Act.xml`.

| Langkah | Aksi | Nilai |
|---|---|---|
| 8 | `Property-Set` | `Local.TotalValueAdjust = 0` |
| 10 | `Property-Set` | `Local.CekLimitPersen = pyWorkPage.TreatyInMaster.RNMShare` |
| 10 | `Property-Set` | `Local.LimitMax = 30000000.00` |
| 10 | `Property-Set` | `Local.LimitMaxDivHead = 50000000.00` |
| 10 | `Property-Set` | `Local.LimitPersenMax = 30.00` · `Local.LimitPersenMaxDivHead = 30.00` |
| 11 | *(loop)* atas `pyWorkPage.ClaimData.AdjustmentList`, `repeat=EMBEDDED` | — |
| 11.1 | `Property-Set` | `Local.TotalValueAdjust = .ValueAdjustment` |
| 12 | `Property-Set`, kondisi `Local.CekLimitPersen<=Local.LimitPersenMax` **dan** `Local.TotalValueAdjust<=Local.LimitMax` | `Local.Flagkomite = 1` |
| 13 | `Property-Set`, kondisi `Local.CekLimitPersen<=Local.LimitPersenMaxDivHead` **dan** `Local.TotalValueAdjust>Local.LimitMax && Local.TotalValueAdjust<=Local.LimitMaxDivHead` | `Local.Flagkomite = 2` |
| 26.4 | `Property-Set`, kondisi `Local.Flagkomite==1` | `Param.LIMIT_BOTTOM = 0` |
| 26.5 | `Property-Set`, kondisi `Local.Flagkomite==2` | `Param.LIMIT_BOTTOM = 25000001` |
| 26.7 | `Call pxRetrieveReportData` | menjalankan `FilterEmailKomiteWithLimit` dengan `Param.LIMIT_BOTTOM` |
| 26.13 | `Property-Set` | `ChildWorkPage.KomiteLoop = @SizeOfPropertyList(pyReportContentPage.pxResults)` |

Hasil kueri `FilterEmailKomiteWithLimit` menentukan **siapa saja anggota Komite** dan **berapa jenjang persetujuan** yang harus dilalui.

## 2. Titik di mana konversi tidak terjadi

Perbandingan terhadap ambang berlangsung di **langkah 12 dan 13**, atas `Local.TotalValueAdjust`.

Nilai itu berasal dari langkah 11.1, yaitu `.ValueAdjustment` apa adanya. Di sepanjang langkah 8 sampai 13 **tidak ada satu pun operasi konversi mata uang**: tidak ada pemanggilan kurs, tidak ada pembacaan `Currency`, tidak ada `@divide`/`@multiply` terhadap faktor konversi. Bandingkan dengan langkah 20 dan 23 pada berkas yang sama, yang justru memakai kurs secara eksplisit (`.KursIDR`, `@divide(...RNMShare,100,20)`) ketika menyusun `CNPLayerList`.

Konstanta `30000000.00` dan `50000000.00` tidak membawa penanda mata uang.

## 3. Temuan sampingan pada rantai yang sama

Langkah 11.1 adalah **penugasan** (`=`), bukan akumulasi (`+=`), di dalam loop `AdjustmentList`.

Akibatnya `Local.TotalValueAdjust` berisi nilai `.ValueAdjustment` dari **baris terakhir** yang dikunjungi loop, bukan jumlah seluruh Adjustment — meskipun namanya "Total". Ini observasi terpisah dari soal mata uang, dan mengoreksi pernyataan saya sebelumnya yang menyebut nilai ini sebagai `Σ AdjustmentList[].ValueAdjustment`.

## 4. Batas klaim — apa yang TIDAK dapat saya simpulkan dari XML

Bagian ini menentukan apakah §2 berarti sesuatu atau tidak. **Ketiganya harus diuji lebih dulu sebelum temuan ini diperlakukan sebagai cacat.**

1. **Mata uang `.ValueAdjustment` tidak diketahui.**
   `.ValueAdjustment` muncul di **tepat satu berkas** di seluruh folder ini — `CreateChildKomiteCNP_Act.xml` — dan di situ ia hanya **dibaca**. Tidak ada satu pun rule di folder ini yang menulisinya, dan ia **bukan field UI** di section atau harness mana pun di folder ini. Asal-usul nilainya: **TIDAK DITEMUKAN DI XML**.
   Konsekuensinya: sangat mungkin nilai itu sudah dikonversi ke IDR di hulu oleh rule di luar folder ini, atau oleh layar milik modul lain. **Bila `.ValueAdjustment` selalu IDR, tidak ada cacat apa pun di sini.**

2. **Keberadaan Adjustment non-IDR di praktik tidak diketahui.**
   Struktur data jelas mendukung banyak mata uang, tetapi apakah Adjustment non-IDR benar-benar pernah dibuat adalah pertanyaan data, bukan pertanyaan kode.

3. **Perilaku bila `.ValueAdjustment` kosong tidak diketahui.**
   Bila properti itu tidak pernah terisi, `Local.TotalValueAdjust` tetap `0` dari langkah 8, sehingga `Flagkomite` selalu bernilai 1 dan `LIMIT_BOTTOM` selalu 0 — mekanisme penjenjangan efektif tidak berjalan sama sekali. Ini kemungkinan ketiga yang sama masuk akalnya, dan berbeda konsekuensinya.

## 5. Yang dibutuhkan untuk menutup temuan ini

**REQ-011** — kuantifikasi dari Oracle. Query ada di `pengetahuan/RECON.sql` BAGIAN 10.

| Pertanyaan | Menentukan |
|---|---|
| Berapa Adjustment bermata uang selain IDR? | Bila nol, temuan ini gugur seluruhnya |
| Berapa di antaranya bernilai setara di atas 30.000.000 IDR? | Besaran paparan |
| Jenjang persetujuan mana yang benar-benar menanganinya | Apakah penjenjangan berjalan sebagaimana dimaksud |
| Apakah `.ValueAdjustment` tersimpan dan dalam mata uang apa | Menjawab batas klaim §4.1 |

## 6. Hubungan dengan sistem baru

Tidak ada. Keputusan untuk sistem baru sudah ditetapkan terpisah di [ADR-0007](./docs/adr/0007-nilai-uang-berpasangan-dan-ambang-idr.md): ambang dibandingkan terhadap nilai IDR hasil konversi, dengan kurs dan tanggal kurs tercatat. Keputusan itu berlaku apa pun hasil investigasi ini.

---

## 7. Tambahan 18 September 2026 — isi `GETCURRENCYSTANDARD` terbaca

Sumber: `pengetahuan/DDL_Script_ClaimNonProp.xls`. Ini menjawab pertanyaan yang sempat saya ajukan sebagai T6 lalu dicabut karena bukan pertanyaan untuk pengguna.

```sql
CREATE OR REPLACE EDITIONABLE FUNCTION "POOLDATA"."GETCURRENCYSTANDARD" (
   i_kurs_id    IN   m_currencystandard.ID%TYPE,
   i_tgl_kurs   IN   m_currencystandard.CurrencyDate%TYPE
)  RETURN NUMBER result_cache
IS
   vhasil   m_currencystandard.CurrencyValue%TYPE;
BEGIN
   SELECT CurrencyValue INTO vhasil
     FROM (SELECT CurrencyValue FROM m_currencystandard
            WHERE ID = i_kurs_id
              AND TRUNC (CurrencyDate) <= TRUNC (sysdate)
         ORDER BY CurrencyDate DESC)
    WHERE ROWNUM < 2;
   RETURN vhasil;
EXCEPTION
   WHEN NO_DATA_FOUND THEN RETURN 1;
END;
```

Dua hal yang terbaca langsung dari badan fungsi ini, dan keduanya berdiri sendiri dari soal ambang di bagian 2.

### 7.1 Parameter tanggal dideklarasikan tetapi tidak pernah dipakai

`i_tgl_kurs` muncul di daftar parameter dan **tidak muncul lagi di mana pun** di dalam badan fungsi. Penyaringnya memakai `TRUNC(CurrencyDate) <= TRUNC(sysdate)` — **tanggal hari ini**, bukan tanggal yang dikirim pemanggil.

Akibatnya: apa pun tanggal yang dikirim, yang dikembalikan selalu **kurs terbaru sampai hari ini**. Seluruh pemanggilan di XML memang mengirim `sysdate` (`RDBList\CurrencyStandard.xml`), jadi untuk pemakaian yang ada sekarang hasilnya sama saja. Yang tertutup adalah kemungkinan lain: **kurs historis tidak dapat diambil lewat fungsi ini**, meskipun tanda tangannya menjanjikan bisa.

Ini berhubungan langsung dengan ADR-0007, yang mewajibkan **tanggal dan sumber kurs** ikut tersimpan. Sistem lama tidak dapat memenuhi syarat itu lewat fungsi ini: nilai yang dikembalikan tidak dapat direproduksi ulang untuk tanggal tertentu di kemudian hari, karena hasilnya bergeser mengikuti `sysdate`. Rekonsiliasi angka lama tidak akan menghasilkan angka yang sama.

### 7.2 Kurs yang tidak ditemukan menghasilkan `1`, bukan galat

`WHEN NO_DATA_FOUND THEN RETURN 1`.

Bila tidak ada baris kurs untuk suatu mata uang, fungsi mengembalikan **`1`** — yang secara aritmetika berarti **konversi satu banding satu**. Nilai 1.000 USD akan diperlakukan sebagai 1.000 IDR, tanpa galat, tanpa pesan, tanpa penanda.

Ini pilihan yang disengaja (ada penangan eksepsi khusus untuk itu), bukan kelalaian. Tetapi konsekuensinya: **kegagalan pencarian kurs tidak dapat dibedakan dari kurs yang memang bernilai 1.**

### 7.3 Objek baru: `m_currencystandard`

Tabel kurs yang sebenarnya, dengan kolom `ID`, `CurrencyDate`, `CurrencyValue`, `UserID`, `InputDate`. **Tidak pernah muncul di 279 berkas XML** — Pega hanya melihat fungsi dan view `CURRENCYSTANDARD` di atasnya. Masuk `pengetahuan/PULL-LIST.csv` sebagai REQ-017, prioritas BLOCKER.

### 7.4 Apa yang ini ubah dan tidak ubah pada temuan ini

**Tidak mengubah** inti FINDING-001: perbandingan `Local.TotalValueAdjust` terhadap `30000000.00` di langkah 12–13 tetap terjadi tanpa konversi apa pun, dan batas klaim di bagian 4 tetap utuh. `.ValueAdjustment` masih bisa saja sudah IDR sejak hulu.

**Menambah** dua hal yang harus diperhitungkan saat REQ-011 dijalankan:

1. Angka kurs yang dipakai sistem lama untuk suatu transaksi **tidak dapat direkonstruksi** dari tanggal transaksi. Jadi kuantifikasi "berapa yang setara di atas ambang" hanya dapat memakai kurs hari ini, dan itu harus dinyatakan sebagai batasan hasilnya.
2. Bila ada mata uang tanpa baris kurs, nilainya akan terbaca seolah sudah IDR. Query REQ-011 harus **memeriksa keberadaan baris kurs secara terpisah**, tidak boleh bersandar pada hasil fungsi ini.
