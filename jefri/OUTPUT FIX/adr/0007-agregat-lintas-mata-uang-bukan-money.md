---
status: accepted
---

# Total yang menjumlahkan lintas mata uang bukan `Money`, melainkan tipe ketiga yang tidak dapat diaritmetikakan

Lima total skalar di akar penawaran — `TOTAL_TSI_NUSA_RE`, `TOTAL_PREMI_NUSA_RE`,
`TOTAL_TSI_NUSA_RE_SPREADING`, `TOTAL_TSI_TOP_RISK`, `SUM_TOTAL_TSI` — **tidak dimodelkan sebagai
`Money`**. Keduanya tidak dapat disatukan: **ADR-0004** mensyaratkan `Money` membawa mata uangnya,
dan di sini **tidak ada satu mata uang pun yang benar untuk dibawa**.

Kelimanya memakai **tipe ketiga** di samping `Money` dan `Ratio`. Sifatnya:

- menyimpan **angka apa adanya**, tanpa mata uang;
- ⛔ **tidak boleh dijumlahkan dengan `Money`**;
- ⛔ **tidak boleh dikonversi**;
- ⛔ **tidak boleh dibandingkan** dengan nilai bermata-uang.

Boleh dibaca, ditampilkan, disimpan, dan direkonsiliasi terhadap nilai lama. Tidak boleh masuk
aritmetika bisnis.

## Mengapa, dengan bukti

`[terverifikasi]` `DDL\CONTOH\NB-173649.xml` adalah penawaran **dua mata uang** — `CurrencyList`
memuat entri **USD** dan **IDR**, dan `PropertyItemList` membawa `Currency` sendiri per baris
(**3 IDR + 2 USD**).

Pada berkas itu, dua kesetaraan terukur **persis**:

```
TotalTSINusaReSpreading  =  jumlah 2 baris TotalTSIPremiSpreadRNM.TSISpreaded
                            TreatyName baris 1 = USD
                            TreatyName baris 2 = IDR
TotalPremiNusaRe         =  jumlah 17 baris CoverageList.PremiNusantaraRe
SumTotalTSI              =  satu nilai LocationList/Property.TotalTSI
```

⛔ **Baris USD dan baris IDR dijumlahkan langsung, tanpa kurs.** Angka yang keluar karena itu bukan
jumlah uang dalam mata uang mana pun — ia jumlah aritmetis atas dua satuan berbeda.

Dua sisanya tidak cocok agregat sederhana mana pun dan **tidak ditebak**:

| Skalar | Diuji terhadap | Hasil |
| --- | --- | :-: |
| `TotalTSINusaRe` | jumlah `CoverageList.TSINusantaraRe` | **BEDA** |
| `TotalTSITopRisk` | `LocationList/Property.TotalTSI` | **BEDA** (bernilai nol di berkas ini) |

Keduanya tetap masuk tipe ketiga: asal-usulnya belum terbaca, sehingga **tidak ada dasar** untuk
menempelkan mata uang pada keduanya.

⚠️ **Yang tidak pecah hanya kelima skalar akar ini.** Struktur lain sudah menangani multi mata uang
dengan benar — `CurrencyList`, `TotalTSIList`, `TotalTSIPremiGrossList` dan
`TotalTSIPremiSpreadRNM` masing-masing **2 baris** di berkas itu, satu per mata uang. Jadi masalahnya
**bukan** bahwa sistem lama tidak mengenal multi mata uang; masalahnya kelima skalar ini meratakannya.

## Considered Options

- **Modelkan sebagai `Money` dengan `Currency = Unknown` (ADR-0006).** ⛔ Ditolak, dan ini penolakan
  yang paling penting. `Unknown` di ADR-0006 berarti *"mata uangnya ada, kita belum tahu apa"* —
  keadaan pengetahuan yang bisa diperbaiki dengan data. Di sini keadaannya berbeda secara jenis:
  **mata uangnya memang tidak ada**, karena angkanya campuran. Memakai `Unknown` menyamarkan cacat
  struktural menjadi celah pengetahuan, dan membuka jalan bagi seseorang kelak "mengisi" mata uang
  yang benar — yang justru **mengubah angka**.
- **Modelkan sebagai `Money` ber-`Currency = IDR`.** ⛔ Ditolak. Menebak (`CLAUDE.md` §3 butir 4),
  dan tebakannya salah persis pada kasus yang paling mahal: penempatan valas.
- **Perbaiki sekarang — konversi ke satu mata uang sebelum dijumlahkan.** ⛔ Ditolak. Itu
  **perbaikan, bukan migrasi** (`CLAUDE.md` §1). Ia mengubah angka yang tersimpan, sehingga
  rekonsiliasi paralel run akan berbeda **dan menyembunyikan salah-port yang sesungguhnya**
  (ADR-0001).
- **`decimal` telanjang.** ⛔ Ditolak dengan alasan yang sama seperti ADR-0004 menolak persen sebagai
  `decimal` telanjang: menghapus justru proteksi yang paling dibutuhkan. Tanpa tipe, tidak ada yang
  mencegah `TOTAL_PREMI_NUSA_RE` dijumlahkan dengan `Money` pada layar ringkasan.

## Consequences

**Kegagalan dipindahkan ke waktu kompilasi.** Sejalan ADR-0004: `Money + Ratio` tidak dapat
dikompilasi, dan kini `Money + <tipe ketiga>` juga tidak. Cacat yang tadinya senyap — menjumlahkan
total campuran dengan nilai bermata-uang — menjadi **kegagalan kompilasi**.

**Kandidat perbaikan bisnis, dicatat bukan dieksekusi.** Bahwa sistem lama menjumlahkan USD dan IDR
tanpa konversi adalah temuan yang layak dibawa ke bisnis. Sampai bisnis memutuskan, kelimanya
**diport apa adanya**.

**Batas terhadap ADR-0006.** `Unknown` tetap dipakai untuk kasus yang dimaksud ADR-0006 — nilai uang
bermata-uang-satu yang informasinya hilang dari korpus. Tipe ketiga ini **bukan** perluasan
`Unknown`; ia tipe lain, untuk angka yang **memang tidak punya** mata uang tunggal.

### ⚠️ Akibat untuk butir 8 — yang MASIH TERBUKA dan tidak diputuskan di sini

`TotalTSINusaReSpreading` adalah **jumlah dari `TotalTSIPremiSpreadRNM`** — tabel yang **V-28**
rencanakan dibuang.

⛔ Bila tabel itu dibuang lalu nilainya dihitung ulang dari `T_SPREADINGLIST`, **hitung ulang itu
WAJIB ikut menjumlahkan lintas mata uang** agar angkanya sama. Bila tidak, selisihnya **bukan soal
pembulatan lagi melainkan beda hasil** — dan tidak ada ambang toleransi yang dapat menutupinya.

📌 `TreatyName` di bawah `TotalTSIPremiSpreadRNM` adalah **pembawa kode mata uang**. Membuang tabel
itu ikut membuang pembawanya. **Butir 8 dan butir 9 saling menyentuh**; butir 8 diputuskan terpisah.
