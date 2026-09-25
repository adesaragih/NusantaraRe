# Pertanyaan untuk Akuntansi

> **SELESAI 2026-09-18 — digantikan `REGISTER-RATIFIKASI.md`; butir 7 ditambahkan 18 Sep dan turun jadi verifikasi 19 Sep 2026.**
> Keenam pertanyaannya ditutup sebagai DECIDED-TEKNIS (delapan butir AK) dan **tidak ditanyakan ulang**; ratifikasi akuntansi menyusul tanpa menahan pekerjaan.
> Disimpan apa adanya karena ia memuat bukti berangka yang melahirkan putusan itu — tabel presisi, keterjangkauan `SetPPNPPH`, dan dua hipotesis yang terbukti gugur.

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.
> Setiap angka dan setiap pernyataan tentang perilaku sistem di dokumen ini membawa nama berkas dan nomor barisnya. Tidak ada yang berdiri tanpa rujukan.

**Enam pertanyaan.** Lima tentang kebijakan presisi dan tarif; satu berisi tiga butir lama yang belum pernah dijawab.

---

## 0. Definisi: apa yang dimaksud "nilai yang masuk jurnal"

Pertanyaan 1, 2, dan 3 memakai frasa ini, jadi isinya ditetapkan lebih dulu. Daftar ini **tertutup** — diturunkan dari penelusuran seluruh jalur yang membawa nilai keluar dari modul klaim non-proporsional.

### 0.1 Tiga jalur keluar, dan hanya satu yang membawa uang

| Jalur | Rule | Yang dibawa |
|---|---|---|
| **Kasir** | `Activity\HitServiceToKasir_Act.xml` → `ConnectREST\SendAcceptationToKasir.xml` | **nilai uang** — menjadi instruksi bayar |
| **Arasapas** | `Activity\KonversiKlaim_Act.xml` → `ConnectREST\KonversiKlaimNonLife.xml` | **hanya `CASEID`, `NOPOLIS`, `STS_REJECT`** — pemicu, bukan data nilai |
| **Persistensi** | `Activity\SaveDataToOSAksep_Act.xml` → `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` | seluruh dokumen `DATA_JSON` |

Jalur Arasapas **tidak membawa satu pun nilai uang**. Nilainya sudah lebih dulu tersimpan di `OS_AKSEPTASI_KLAIM.DATA_JSON` lewat jalur persistensi, dan Arasapas membacanya dari sana.

### 0.2 Daftar tertutup — nilai yang menjadi instruksi bayar

Dari `Activity\HitServiceToKasir_Act.xml`, muatan `TempKasir.CARIn`:

| Posisi | Isi | Sifat |
|---|---|---|
| `CARI9` | `.TotalClaim - .PremiumSpreaded` | **nilai uang** |
| `CARI8` | `.AdjustmentValue` | **nilai uang** |
| `CARI22` | `@toDecimal(TempKasir.CARI9) + .AdjustmentValue` | **nilai uang** |
| `CARI23` | `@toDecimal(TempKasir.CARI9) + .AdjusterFeeValue` | **nilai uang** |
| `CARI24` | `@toDecimal(TempKasir.CARI9) + .SalvageValue` | **nilai uang** |
| `CARI9` (cabang lain) | `.IndividualRiskRNM` | **nilai uang** |
| `CARI10` | `.IndividualRiskPercentage` | persentase |
| `CARI11`, `CARI12`, `CARI13` | `"NUSARE"`, `"D0031"`, `"100081"` | kode tetap |
| `CARI15`/`CARI17`/`CARI18` | `"100115"` pada satu cabang | kode tetap, berbeda per node |

**Jadi "nilai yang masuk jurnal" berarti tepat tujuh besaran**: `.TotalClaim`, `.PremiumSpreaded`, `.AdjustmentValue`, `.AdjusterFeeValue`, `.SalvageValue`, `.IndividualRiskRNM`, dan `.IndividualRiskPercentage`. Aturan "2 desimal, tanpa toleransi" berlaku untuk ketujuh itu, **tidak untuk nilai antara**.

**Catatan yang perlu diketahui**: ketujuh nilai itu disusun dengan `@toDecimal(...)` **tanpa pembersihan teks** — lihat pertanyaan 4.

---

## 1. Berapa desimal untuk nilai antara?

**Sudah disepakati**: nilai yang masuk jurnal (daftar 0.2) dibulatkan ke **2 desimal**. Yang belum: berapa desimal untuk perhitungan di tengah jalan.

### 1.1 Keadaan sekarang

| Fakta | Sumber |
|---|---|
| Dari 120 kolom `NUMBER` di basis data, **88 tidak menyatakan presisi maupun scale sama sekali** | `pengetahuan/ddl/*.sql`, terhitung di `pengetahuan/SCHEMA-ACTUAL.csv` |
| Satu-satunya yang menyatakan: `NUMBER(20,4)` pada 23 kolom nilai | `pengetahuan/ddl/TABLE_TREATYINPRODUCTION.sql` |
| Nilai klaim **tidak punya kolom sama sekali** — tersimpan sebagai teks di dalam JSON | `pengetahuan/ddl/VIEW_CLAIMXOL.sql`, seluruh kolom nilai bertipe `varchar2` |
| 77% ekspresi aritmetika di kode tidak menyatakan skala | `pengetahuan/arithmetic-inventory.tsv`, 673 ekspresi |
| Skala pajak `8` | `Activity\SetPPNPPH.xml` baris 822, 890, 911 — `@divide(2.5,100,8)`, `@divide(2,100,8)`, `@divide(2.2,100,8)` |

**Basis data tidak memaksakan pembulatan apa pun.** Apa pun yang dihitung, tersimpan apa adanya.

### 1.2 Di mana angkanya benar-benar pecah — contoh berangka

Rumus premi reinstatement. **Ada dua varian di dalam sistem, dan keduanya berbeda tepat pada hal yang sedang ditanyakan:**

| Berkas & baris | Ekspresi | Skala |
|---|---|---|
| `Activity\CountReinstatement_Act.xml` **baris 624** | `((((.ClaimEstimation+.AdjusterFee) - (.Salvage*100/RNMShare)) / .CNPLimit) * .CNPMDP) * (.CNPPctReinstate/100)` | **tidak dinyatakan sama sekali** |
| `Activity\AdjClaimCNP_Act.xml` **baris 3109** | `@divide(.TotalClaim, .CNPLimit, 20) * .CNPMDP * @divide(.CNPPctReinstate, 100, 20)` | **20** |

Varian pertama menyerahkan presisinya kepada perilaku bawaan Pega, yang tidak terdokumentasi di dalam ekspor ini. Varian kedua menetapkan 20 desimal secara tegas. Keduanya juga memakai **pembilang yang berbeda** — yang satu `ClaimEstimation + AdjusterFee − Salvage×100/RNMShare`, yang lain `TotalClaim` — sehingga bukan sekadar dua penulisan untuk rumus yang sama.

Bentuknya sama: satu pembagian, lalu dua perkalian. Pembulatan pada hasil bagi **diperbesar** oleh dua pengali sesudahnya.

Contoh: pembilang Rp 4.490.977.654,68 · CNPLimit Rp 12.500.000.000 · CNPMDP Rp 1.750.000.000 · Reinstate 100%

| Desimal pada hasil bagi | Hasil bagi | Premi reinstatement | **Selisih** |
|---|---|---|---|
| **2** | 0,36 | Rp 630.000.000,00 | **+Rp 1.263.128** |
| **4** | 0,3593 | Rp 628.775.000,00 | **+Rp 38.128** |
| **6** | 0,359278 | Rp 628.736.500,00 | −Rp 372 |
| **8** | 0,35927821 | Rp 628.736.867,50 | −Rp 4 |
| tanpa pembulatan | 0,35927821237… | Rp 628.736.871,66 | — |

**Pada satu klaim.** Dua desimal meleset Rp 1,26 juta; empat desimal masih meleset Rp 38 ribu.

### 1.3 Pertanyaan

Berapa desimal untuk nilai antara?

| | Nilai antara | Akibat menurut tabel 1.2 |
|---|---|---|
| **A** | 4 desimal | masih meleset ±Rp 38 ribu per klaim |
| **B** | 6 desimal | meleset ±Rp 372 per klaim |
| **C** | 8 desimal | meleset ±Rp 4 per klaim |

**Usulan: B, enam desimal** — dengan catatan bahwa rantai perhitungan yang memuat pembagian diikuti perkalian besar sebaiknya memakai delapan.

**Dua hal yang perlu disampaikan terus terang:**

1. Dugaan bahwa **alokasi berlapis** menumpuk galat antar layer (`Activity\CountLossAllocation_act.xml`) **tidak terbukti** pada pengujian: dengan 7 layer, hasilnya identik pada 2, 4, 6, maupun 20 desimal. Kerusakannya ada di rumus reinstatement, bukan di perulangan alokasi.
2. Dugaan bahwa **pembagian-lalu-perkalian-kembali** dengan persentase share merusak angka (`AdjusterFeeValue = AdjusterFee ÷ (RNMShare÷100)`, lalu `TotalXOLRNM = TotalClaim × (ClaimPercentage÷100)`) juga **tidak terbukti**: pembulatan akhir ke 2 desimal menyerap galatnya, bahkan bila nilai antara hanya 2 desimal. Diuji pada share 37,5%, 33,33%, 12,35%, 7,77%, dan 23,45%.

Jadi yang membenarkan enam desimal hanyalah tabel 1.2. Itu satu alasan, tetapi alasan yang kuat.

### 1.4 Pertanyaan sejarah, bukan penentu

**Dari mana `NUMBER(20,4)` pada `TREATYINPRODUCTION` berasal?** Bila ia ketentuan resmi, ia ketentuan untuk **nilai tersimpan di tabel itu**, bukan untuk presisi perhitungan di tengah jalan. Jawabannya berguna, tetapi tidak menentukan jawaban 1.3.

**Menghambat**: `ADR-0003` (draft sejak awal).

---

## 2a. Bolehkah pembulatan mengubah angka lama?

Bila data produksi menyimpan lebih dari 2 desimal, menetapkan pembulatan 2 desimal **akan mengubah angka yang selama ini diterima akuntansi**.

**Pertanyaan**: ikuti presisi lama apa adanya, atau bulatkan dan terima selisihnya?

**Menghambat**: `ADR-0003`, dan kriteria lulus shadow-run di `ADR-0005`.

---

## 2b. Bila perhitungan lama ternyata **keliru** — dipertahankan atau diperbaiki?

Ini pertanyaan yang berbeda sifatnya dari 2a, dan tidak boleh dijawab bersamaan. 2a soal presisi; 2b soal nilai yang hilang sama sekali.

### 2b.1 Contoh yang sudah terverifikasi

`Activity\CountLossAllocation_act.xml`, penetapan `Local.TotalUR`:

```
@if(Local.Currency=="IDR",  (@divide(.Deductible,Local.Kurs,10)*Local.ProrateClaim/100),
@if(Local.Currency==.Currency, (.Deductible * 0 * Local.ProrateClaim/100),
                               ... ))
```

Pada cabang **`Local.Currency == .Currency`** — yaitu ketika mata uang klaim sama dengan mata uang yang sedang diproses — perhitungannya dikalikan **`0`**. Hasilnya **selalu nol**.

Bandingkan dengan penetapan `Local.UR` pada berkas yang sama, di cabang yang sama persis:

```
@if(Local.Currency==.Currency, (.Deductible * Local.ProrateClaim/100), ... )
```

Tanpa `* 0`.

**Jadi `TotalUR ≠ UR` untuk kasus mata uang sama, dan `TotalUR` bernilai nol.** Ini bukan soal desimal — ini nilai yang tidak pernah terbentuk.

### 2b.2 Nilai nol itu ikut tersimpan

`Activity\CountLossAllocation_act.xml` baris 7015–7120 menulis baris Retensi Cedant:

```
SpreadingRisk(<LAST>).TotalClaim = Local.TotalUR
```

Jadi nol itu **masuk ke baris Retensi Cedant**, yang pada sistem baru menjadi tabel tersendiri (`ADR-0010`). Bila perilaku lama dipertahankan, tabel baru itu akan memuat kolom yang selalu nol untuk seluruh klaim bermata uang sama.

### 2b.3 Pertanyaan

Bila ditemukan perhitungan lama yang keliru — bukan kurang presisi, melainkan salah — apakah sistem baru **mempertahankan kesalahannya** demi kesesuaian angka dengan data lama, atau **memperbaikinya** dan menerima selisih terhadap data lama?

Kedua jawaban sah. Akibatnya sangat berbeda, dan itu sebabnya dipisah dari 2a.

**Menghambat**: `ADR-0010`, dan definisi baseline shadow-run di `ADR-0005`.

---

## 3. Toleransi shadow-run untuk nilai antara

**Sudah disepakati**: nilai yang masuk jurnal (daftar 0.2) harus **sama persis sampai 2 desimal, tanpa toleransi**.

**Pertanyaan**: berapa toleransi untuk nilai antara?

**Konteks besaran** — agar angka toleransi dapat dinilai: contoh di tabel 1.2 memakai `CNPLimit` Rp 12,5 miliar dan `CNPMDP` Rp 1,75 miliar, yang merupakan besaran wajar untuk satu layer XOL. Nilai klaim tertinggi dan rata-rata per tahun **belum tersedia** — itu menunggu profil data (REQ-005), dan bila diperlukan sebelum pertemuan, angkanya dapat ditarik lebih dulu.

**Usulan: dua batas berjalan bersama, yang dilanggar duluan yang berlaku.**

| | Batas |
|---|---|
| Relatif | 0,01% dari nilai baris |
| Mutlak | Rp 1.000 per baris |

Alasannya: batas relatif saja terlalu longgar untuk klaim besar — 0,1% dari Rp 50 miliar adalah Rp 50 juta. Batas mutlak saja terlalu ketat untuk klaim besar. Dua batas bersama menutup kedua ujungnya.

Setiap selisih di atas ambang wajib dijelaskan satu per satu, tidak diabaikan sebagai derau.

**Menghambat**: `ADR-0005`.

---

## 4. Ambang nilai untuk menghentikan migrasi

### 4.1 Ini fakta, bukan kemungkinan

Nilai uang di sistem lama tersimpan sebagai **teks**. Dan sistem lama **sudah pernah menemui teks yang tidak dapat diurai**, lalu menanganinya:

```
@toDecimal(@replaceAll(.Deductible2, ",", "."))
```
`Activity\CountLossAllocation_act.xml` — mengganti koma menjadi titik sebelum mengubah ke angka.

Tidak ada yang menulis pembersihan semacam itu untuk masalah yang tidak pernah terjadi. **Data dengan koma sebagai pemisah desimal ada di produksi.**

### 4.2 Penanganannya tidak seragam — dan sangat timpang

Sapuan seluruh 279 berkas atas `toDecimal`:

| | Jumlah |
|---|---|
| Seluruh pemanggilan `toDecimal` | **127** |
| Didahului pembersihan `replaceAll` | **9** — seluruhnya atas `.Deductible2` |
| **Tanpa pembersihan apa pun** | **118** |

Properti yang paling sering diubah ke angka **tanpa dibersihkan**: `TempDla.CARI34` (12×), `pyWorkPage.TreatyInMaster.Share` (9×), `OutSpreading.pxResults` (7×), `.KursValue` (6×), `TempKasir.CARI20` dan `CARI21` (6× masing-masing), `.TSIPerObject` (4×), `.BalanceDueTo` (3×), `.PersenRNM` (3×).

**Dua di antaranya adalah muatan ke kasir** — `TempKasir.CARI20` dan `CARI21` — yaitu nilai yang masuk jurnal menurut daftar 0.2.

Satu tempat membersihkan; 118 tempat tidak. Setiap tempat yang membersihkan menandai satu jenis kotoran yang pernah ditemui seseorang; setiap tempat yang tidak membersihkan adalah calon kegagalan.

### 4.3 Pertanyaan

**Ambangnya berbasis nilai, bukan cacah baris.** Satu baris senilai lima miliar lebih berat daripada lima ratus baris senilai seratus ribu.

Dua angka:

1. **Nilai total** yang belum terselesaikan yang masih dapat diterima sebelum migrasi dihentikan dan data dibereskan lebih dulu.
2. **Nilai satu baris** yang, bila dilampaui, langsung menghentikan migrasi tanpa menunggu total.

Setiap baris yang tidak dapat diurai tetap diselesaikan satu per satu, berapa pun jumlahnya. Tidak ada baris yang dibuang karena "cuma sedikit".

**Menghambat**: `ADR-0014`, dan rencana pembersihan data sebelum cutover.

---

## 5. Tarif pajak dan brokerage

### 5.1 Keterjangkauan — diverifikasi lebih dulu

Sempat diragukan apakah `SetPPNPPH` termasuk lingkup klaim non-proporsional, karena rule itu berada pada class `ASM-FW-GISFW-Data-PolicyTreatyIn` — kerangka polis, bukan kerangka klaim.

**Rantainya sampai, dan seluruh pangkalnya di class klaim:**

```
Section\OutstandingClaim(1).xml                      (layar klaim)
  -> Harness\ViewPolisNonProp.xml                    class ASM-FW-GCNMFW-Work-ClaimTreatyNonProp
     -> Section\ViewDetailDeptHeadTreatyIn_UW.xml    <pyInclude>
        -> <pyActivity>CountNetPremi_act
           -> Call SetPPNPPH
```

Jadi pertanyaan ini **sah** dan tetap dibawa. Sifatnya perlu dicatat: jalur itu adalah layar **peninjauan polis** yang dibuka petugas klaim, jadi tarifnya memengaruhi angka premi yang **dilihat**, bukan langsung nilai pembayaran klaim.

### 5.2 Tarif, beserta letaknya

Seluruhnya di `Activity\SetPPNPPH.xml`:

| Baris | Properti | Ekspresi persis | Tarif |
|---|---|---|---|
| 821–822 | `.BrokerageFee` | `@divide(2.5,100,8) * (.PremiOgp + .PremiOnp)` | **brokerage 2,5%** |
| 868–869 | `.BrokerageFeeSebenarnya` | `@if(.TypeTax=="Inclusive", @divide(.Deduction1, @divide(102.2,100,8), 8), .Deduction1)` | **faktor 102,2** |
| 889–890 | `.PPHValue` | `.BrokerageFeeSebenarnya * @divide(2,100,8)` | **PPh 2%** |
| 910–911 | `.PPNValue` | `.BrokerageFeeSebenarnya * @divide(2.2,100,8)` | **PPN 2,2%** |

**Angka keempat yang belum pernah dilaporkan**: `102.2` pada baris 869 adalah faktor gross-up untuk pajak *inclusive* — yaitu `100 + 2,2`. Ia **turunan dari tarif PPN**, tetapi ditulis terpisah. Bila PPN berubah, angka itu harus ikut diubah di tempat lain, dan tidak ada apa pun yang memaksanya.

**Satu hal yang tidak saya simpulkan**: `.BrokerageFee` dihitung 2,5% dari premi, tetapi `.BrokerageFeeSebenarnya` — yang menjadi dasar PPh dan PPN — diturunkan dari `.Deduction1`, **bukan** dari `.BrokerageFee`. Apakah `.BrokerageFee` karena itu hanya untuk tampilan, tidak dapat saya pastikan dari XML.

### 5.3 Pertanyaan, tiga bagian

1. Apakah ketiga tarif masih benar per hari ini?
2. Pernahkah berubah sejak sistem berjalan? Bila ya, **klaim lama dihitung dengan tarif lama atau tarif baru** — itu menentukan apakah shadow-run harus menyimpan tarif historis.
3. Apakah tarif dapat berbeda per jenis bisnis, per cedant, atau per mata uang? Kode sekarang hanya mendukung satu nilai untuk semua.

**Yang diusulkan apa pun jawabannya**: tarif pindah ke konfigurasi bertanggal berlaku, bukan tertanam di kode — termasuk faktor `102,2` yang harus diturunkan dari tarif PPN, bukan ditulis terpisah. Itu tidak perlu persetujuan akuntansi; yang perlu adalah **nilai dan tanggal berlakunya**.

**Menghambat**: perhitungan pajak di sistem baru, dan kriteria shadow-run untuk nilai yang masuk jurnal.

---

## 6. Tiga butir lama yang belum pernah dijawab

Terbuka sejak 17 September, tercatat di `MEMORI_PEMAHAMAN.MD` §11. Dibawa sekalian karena dua di antaranya akan sangat mahal bila baru ditanyakan setelah sistem baru jadi.

### 6.1 Ambang otoritas Komite — 15 atau 30?

`MEMORI_PEMAHAMAN.MD` §11 butir 13. `Activity\CreateChildKomiteCNP_Act.xml` langkah 10 menetapkan `Local.LimitPersenMax = 30.00` **dan** `Local.LimitPersenMaxDivHead = 30.00` — dua ambang yang berbeda perannya diberi nilai yang sama.

**Bila salah satunya seharusnya 15, maka satu jenjang persetujuan tidak pernah berjalan selama ini.**

### 6.2 `TotalUR` yang selalu nol — perilaku benar atau bukan?

`MEMORI_PEMAHAMAN.MD` §11 butir 14. Ini butir **2b** di atas, sudah terverifikasi dari XML.

### 6.3 Aturan cut-off tanggal produksi — masih berlaku?

`MEMORI_PEMAHAMAN.MD` §11 butir 17. Aturannya ada di dua tempat, dan keduanya memakai tanggal 25:

| Tempat | Bentuk |
|---|---|
| `Activity\HitServiceToKasir_Act.xml` | `@if(TempKasir.CARI19 > 25, @toDecimal(TempKasir.CARI20) + 1, …)` — bila tanggal lewat 25, bulan digeser ke depan |
| `pengetahuan/ddl/PROCEDURE_PROC_GENERATE_SEQUENCE_NUMBER.sql` | membaca `POOLDATA.TANGGAL_CLOSING`, lalu `ADD_MONTHS(v_now, CASE WHEN TO_NUMBER(TO_CHAR(v_now,'DD')) > v_day_closing THEN 1 …)` |

Yang di kode Pega **menuliskan 25 langsung**; yang di basis data **membacanya dari tabel**. Dua sumber kebenaran untuk satu aturan.

**Pertanyaan**: tanggal 25 masih berlaku? Dan bila tanggal di `TANGGAL_CLOSING` diubah, apakah disadari bahwa angka 25 di kode Pega tidak ikut berubah?

---

## 7. Nilai rupiah dari biaya penilaian dan salvage — **butir baru, 18 September 2026**

> Butir ini ditambahkan sesudah berkas ini diarsipkan. Keenam butir di atas sudah tertutup lewat delapan butir AK.
>
> **Kedudukan butir 7 berubah 19 September 2026: ia turun dari penahan menjadi verifikasi, dan berkas ini kembali arsip.**

### 7.1 Apa yang ditemukan

Sapuan S4 atas 279 berkas ekspor, empat lapisan, mencari setiap pengenal bersufiks `IDR` dan membandingkannya terhadap daftar properti bernilai uang. Penamaan alternatif — `Rupiah`, `Idr`, `_IDR`, `Rp`, `InRp` — ikut disapu dan mengembalikan **nol**.

**Tujuh baris, sembilan nama, tanpa padanan IDR sama sekali:**

| Properti uang | Yang ada sebagai gantinya |
|---|---|
| `AdjusterFee`, `AdjusterFeeValue` | `AdjusterFeeRNM` |
| `Salvage`, `SalvageValue` | `SalvageRNM` |
| `CNPOthersFee` | `CNPOthersFeeRNM` |
| `OthersFee` | `OthersFeeRNM` |
| `AdjustmentValue` | — |
| `TotalClaim` | `TotalClaimKurs`, `TotalClaimRNM` |
| `PremiumSpreaded` | — |

Sufiks `RNM` bukan padanan mata uang — ia porsi pihak. `TotalClaimKurs` adalah kursnya, bukan nilai terkonversinya.

Lima dari sembilan masuk hitungan instruksi bayar di `HitServiceToKasir_Act`. Rantai `.AdjusterFeeValue` ditelusuri dari lahir sampai keluar (S10) dan **tidak melewati satu pun titik konversi**: ia masuk instruksi bayar dalam mata uang aslinya.

### 7.2 Pertanyaannya

**Apakah nilai rupiah dari biaya penilaian dan salvage pernah dibutuhkan — dan bila ya, selama ini diambil dari mana?**

Pertanyaannya bukan "kolomnya kurang atau tidak"; itu rancangan, dan rancangan menunggu jawaban ini.

### 7.3 Mengapa kolom tidak dapat menjawabnya

Kolom hanya menunjukkan bahwa padanannya tidak ada. Ia tidak menunjukkan apakah ketiadaan itu **kesengajaan** — karena besaran itu memang tidak pernah masuk jurnal rupiah — atau **kekurangan** yang selama ini ditambal di luar sistem.

### 7.4 Sampai terjawab — **dicabut 19 September 2026**

Rumusan lama berbunyi: *"berlaku disiplin yang sama seperti H1/H2 — tidak boleh ada rancangan yang mengandaikan salah satu dari dua pembacaan."* **Disiplin itu dicabut, dan kedua cabangnya ikut dicabut dari badan ADR-0029.**

Sebabnya bukan jawabannya datang, melainkan **kedua cabang tidak berbiaya setara**:

| Bila ternyata… | dan kolomnya disediakan | dan tidak disediakan |
|---|---|---|
| rupiahnya tidak pernah dibutuhkan | kolom kosong — tidak ada yang rugi | benar, kebetulan |
| rupiahnya dibutuhkan | benar | **migrasi kedua** |

**ADR-0029 karena itu naik ke `accepted` dan berlaku seragam** untuk seluruh nilai uang tanpa pengecualian — termasuk kesembilan nama di 7.1. Kolomnya disediakan; bila memang tidak pernah ada isinya, ia kosong (**kosong, bukan nol** — ADR-0019).

### 7.5 Yang tersisa untuk akuntansi

**Butir 7 tetap berdiri, sebagai verifikasi**, dan pertanyaannya tidak berubah: *apakah nilai rupiah dari biaya penilaian dan salvage pernah dibutuhkan, dan bila ya, selama ini diambil dari mana?*

Yang berubah adalah akibat jawabannya. Sekarang ia menentukan **apakah kolom itu akan terisi**, bukan apakah kolom itu akan ada. **Ia tidak menahan ADR-0029, tidak menahan `SPEC-MODEL-DATA.md`, dan tidak menahan tiket mana pun.**

Butir 5 (tarif brokerage / PPh / PPN) berkedudukan sama: **verifikasi nilai**, bukan penahan — ADR-0025 menjadikan tarif data bertanggal berlaku, sehingga memperbaikinya adalah pekerjaan entri.

Temuan terkait yang memperberat, bukan menjawab: **D23** — `ClaimAmountIDR` diisi oleh dua rule, dan salah satunya (`SetActualPremium_ACT` baris 681) menyalin `TotalClaim` apa adanya tanpa menguji mata uang dan tanpa mengalikan kurs. Kolom berlabel IDR karena itu dapat memuat angka yang bukan rupiah.
