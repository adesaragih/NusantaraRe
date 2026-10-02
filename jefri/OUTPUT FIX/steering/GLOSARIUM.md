# Facultative Inward — Glosarium

Penerimaan risiko reasuransi fakultatif masuk. Satu bounded context; tiga siklus (New Business,
Renewal, Endorsement) berjalan di atas **satu basis rule**, dibedakan saat runtime.

Glosarium ini **hanya kosakata** — tanpa detail implementasi. Keputusan rancangan ada di `../adr/`,
keputusan lingkup di `../00-KEPUTUSAN-WORK-OWNER.md`.

Istilah yang artinya **belum diketahui dari korpus** tetap dicantumkan dengan penanda ☐ —
menghapusnya membuat kekurangan itu tidak terlihat.

---

## Siklus dan kasus

**Siklus**:
Salah satu dari tiga jalur bisnis — New Business, Renewal, Endorsement — yang dibedakan oleh nilai
`StatusBusiness` pada satu basis rule yang sama.
_Avoid_: modul, aplikasi, produk

**Penawaran (`OfferFacIn`)**:
Agregat utama yang menampung seluruh data bisnis sebuah kasus, diserialkan utuh ke kolom JSON polis.
_Avoid_: order, submission, deal

**Before-image**:
Salinan kerja nilai sebelum perubahan, dipakai menghitung selisih endorsement. Ada tiga mekanisme
berbeda yang tidak boleh disatukan.
_Avoid_: riwayat, history, snapshot

**Versi polis**:
Penomoran bertambah pada baris polis; riwayat endorsement tersimpan sebagai **baris baru**, bukan
pembaruan di tempat.
_Avoid_: revisi, before-image

---

## Tangga persetujuan

**Tangga akseptasi**:
Mesin persetujuan berjenjang berbasis limit wewenang. **Satu keputusan manusia = satu transisi**,
bukan satu proses yang menghitung seluruh rantai approver sekaligus.
_Avoid_: workflow approval, approval chain, loop persetujuan

**Antrean (`PositionNote`)**:
Peran yang **sedang** memegang kasus. Ruang nama tersendiri, mis. `ReasFacInUnderwriting`.
_Avoid_: jabatan, posisi, status

**Kode jabatan (`next_approver_position`)**:
Jabatan **tujuan berikutnya** dalam tangga. Ruang nama tersendiri, mis. `SENIORUW`. Disimpan di
properti warisan bernama `LetterNo`, yang **tidak** berisi nomor surat.
_Avoid_: nomor surat, LetterNo, antrean

**Hasil keputusan (`ProposalAcceptStatus`)**:
Hasil keputusan underwriting terakhir. Enam nilai, artinya terverifikasi dari korpus:
`1` Accept · `2` Reject · `3` Ask · `4` **Banding** · `7` Decline · `9` Revise.
_Avoid_: status kasus, state

**Banding**:
Jalur naik-banding setelah penolakan, memakai ambang kedua atau atasan penolak terakhir.
_Avoid_: appeal, eskalasi

**Nilai dasar akseptasi**:
Nilai yang dibandingkan dengan tabel limit. **New Business dan Renewal memakai TSI penuh;
Endorsement memakai selisih** terhadap before-image.
_Avoid_: TSI, nilai transaksi

**Tangga selesai**:
Keadaan ketika tidak ada jabatan tujuan yang cocok — wewenang sudah cukup. Ini **penyelesaian
normal, bukan galat**.
_Avoid_: gagal, tidak ditemukan, error

---

## Nilai dan satuan

**Uang (`Money`)**:
Nilai moneter beserta mata uangnya. Mata uang **wajib menyertai**, dan boleh berkeadaan
`Unknown` yang eksplisit.
_Avoid_: amount, nominal, angka

**Rasio (`Ratio`)**:
Rate atau persentase beserta **skalanya**. Bukan uang, dan tidak dapat dijumlahkan dengan uang.
_Avoid_: rate, persen, faktor

**Skala rasio**:
Satuan sebuah rasio — per mille (‰) atau persen (%). Ditentukan **lini bisnis**, bukan properti:
FIRE, PA, dan Layering memakai ‰; ANEKA, BONDING, GOLF, MARINE CARGO, dan MBU memakai %.
Dikunci oleh K-018.
_Avoid_: satuan rate, format

**Pembagi komposit**:
Pembagi gabungan pada satu operasi pembagian di rumus premi, yang merupakan **hasil kali** sumbangan
tiap faktor bersatuan — bukan satu satuan tunggal. Satuan sebuah faktor dibaca dari sumbangannya
sendiri, tidak pernah dari pembagi total.
_Avoid_: pembagi, skala rumus

**Lini bisnis (COB)**:
Kelompok produk yang menentukan cabang perhitungan dan skala rasio. Dikenali lewat predikat
`IsFire`, `IsPA`, `IsMBU`, `IsAneka`, `IsMarineCargo`, `IsGolfInsurance`, `IsBondingAndCustomBonds`.
_Avoid_: produk, jenis asuransi, class

**Pasangan nilai-sesudah / selisih**:
Cara tabel produksi menyimpan perubahan endorsement: **nilai sesudah dan delta berdampingan**, bukan
satu nilai.
_Avoid_: nilai baru, perubahan

---

## Spreading dan kapasitas

**Spreading**:
Pembagian risiko ke beberapa penerima. Terdapat **tiga mesin berbeda** dengan presisi berbeda; tidak
boleh disatukan.
_Avoid_: alokasi, distribusi, sesi

**Proteksi kapasitas**:
Pemeriksaan ambang kapasitas. `[terverifikasi]` Akibat pelanggarannya adalah **tombol kirim mati**,
**bukan** eskalasi approver.
_Avoid_: limit check, validasi limit

**Skor risiko**:
Hasil penilaian risiko. `[terverifikasi]` **Tidak menggerakkan alur sama sekali** — hanya ditampilkan;
satu-satunya pemakaian non-tampilan memeriksa kelengkapannya, bukan nilainya.
_Avoid_: rating, grade

---

## Enumerasi — **artinya terverifikasi, seluruhnya tertutup** (17 September 2026)

> 📌 **Otoritasnya K-029** di `../00-KEPUTUSAN-WORK-OWNER.md`. Bagian ini adalah kosakatanya; bila
> keduanya berbeda, **K-029 yang berlaku**.
>
> `[terverifikasi]` Keenam berkas property Pega di `D:\migrasi\RNM\DDL\` memuat `<pyPromptTableList>`
> berisi pasangan `<pyStandardValue>` (kode) dan `<pyLocalizedValue>` (arti). Nilai **dan** artinya
> terbaca langsung dari ekspor. Tidak ada yang ditebak di bawah ini.

**`QuotationData.Type`** — jenis penyesuaian · sumber `DDL\Type.xml` (ruleset 01-01-71)

| Kode | Arti | | Kode | Arti |
| ---: | --- | --- | ---: | --- |
| `0` | Adjustment Reff. Number | | `8` | Adjustment Deduction |
| `1` | Extend Period | | `9` | Adjustment Insured Name |
| `2` | Adjustment TSI / Add Object / Rate / Premium | | `11` | Adjustment Share Cedant |
| `4` | Adjustment Spreading | | `12` | Adjustment PPN/PPH |
| `6` | Adjustment Period | | | |
| `7` | Adjustment Currency | | | |

**`QuotationData.EdmType`** — jenis endorsement · sumber `DDL\EdmType.xml` (ruleset 01-01-90)

| Kode | Arti |
| ---: | --- |
| `1` | Batal Sejak Semula |
| `2` | Batal Prorata |
| `4` | Penambahan / Pengurangan / Perubahan |

**`TypeDeductible`** · sumber `DDL\TypeDeductible.xml` (ruleset 01-01-60)

| Kode | Arti | | Kode | Arti |
| ---: | --- | --- | ---: | --- |
| `0` | NIL | | `4` | % Of Approved Loss Value |
| `1` | % Of Claim | | `5` | % Of Recoverable Claim Amount |
| `2` | % Of Loss | | `6` | % Of Recoverable Amount |
| `3` | % Of TSI | | `7` | In Amount |

**`TypeDeductible2`** · sumber `DDL\TypeDeductible2.xml` (ruleset 01-01-60)

| Kode | Arti |
| ---: | --- |
| `0` | NIL · `1` % Of Claim · `2` % Of Loss · `3` % Of TSI · `4` % Of TSI Whichever Is Higher |

⚠️ **`TypeDeductible` dan `TypeDeductible2` adalah dua property berbeda**, bukan satu. Kode `4`
berarti hal yang **berlainan** di masing-masing (*Approved Loss Value* vs *TSI Whichever Is Higher*).
Menyatukannya menghasilkan salah tafsir. Ini menjelaskan sebagian catatan lama tentang "dua ruang
nilai" pada deductible.

**`TeamGroup`** · sumber `DDL\TeamGroup.xml` (ruleset 01-01-53)

| Kode | Arti |
| ---: | --- |
| `1`–`4` | Group 1 … Group 4 |
| `5` | **Bonding** — memutus pola "Group N"; jangan diturunkan dari nomornya |

**`ProRateType`** · sumber `DDL\ProRateType.xml` (ruleset 01-01-54)

| Kode | Arti |
| ---: | --- |
| `1` | Based on age and period (At Once) |
| `2` | Based on period (At Once) |
| `3` | Based on age (Yearly) |
| `4` | Monthly |

---

**`IsB2B`** — **flag penanda**, bukan boolean meski namanya berawalan `Is` · sumber `DDL\B2B.xml`

| Nilai sah |
| --- |
| `ASM` · `KBRU` · `BDX` · `SRB` |

✅ **Ditutup K-029.** Tidak ada arti bisnis yang lebih dalam — **nilainya sendiri adalah isinya**,
diperlakukan apa adanya sebagai kode string.
_Avoid_: boolean, penanda ASM, flag kemitraan

---

## Kode usang — masih tertulis di rule, **tetap diport**

✅ Ditutup **K-029**. Ketiganya dahulu ada lalu dihapus; jejaknya sisa yang kelewat. Pola sama dengan
`IsFacout` (K-019) dan `JUW_A` (K-024).

⚠️ **Sikap portingnya berbeda — jangan disamakan:**

| Kode | Jejaknya di mana | Sikap |
| --- | --- | --- |
| `Type = 3` · `Type = 5` | kondisi rule `When` | Diport apa adanya. `[terverifikasi]` efektif setara `IsEdmAdjTSI` karena tidak pernah ditawarkan UI — menghapusnya tidak mengubah apa pun, tetapi tetap tidak dihapus |
| **`EdmType = 3`** | **logika perhitungan premi yang aktif** | ⛔ **WAJIB diport apa adanya.** `[terverifikasi]` cabang `EdmType==1\|\|==2\|\|==3` (selisih persentase × −1) ada di **7 berkas folder NB**. Menghapus `==3` **mengubah perilaku record lama** |

Status usang **dicatat, tidak dieksekusi**: tanpa data baru bernilai `3`, cabangnya menjadi jalur mati
dengan sendirinya — tanpa satu baris kode pun diubah.

---

## Istilah yang artinya **masih belum terverifikasi**

Nilai literalnya terbaca; artinya tidak. **Jangan ditebak** (`CLAUDE.md` §3 butir 4).

| Istilah | Nilai yang terbaca | Arti | Pemilik |
| --- | --- | :-: | --- |
| `BusinessCode` | 98 kode | ☐ | Product |
| `BusinessOldId` | 87 kode | ☐ | Product |

**Hanya dua ini yang tersisa.** Seluruh enumerasi lain sudah tertutup K-029.

---

## Kosakata yang dihindari

`[terverifikasi]` Nama warisan yang **tidak menggambarkan isinya**:

| Jangan pakai | Pakai | Alasan |
| --- | --- | --- |
| "nomor surat" untuk `LetterNo` | **kode jabatan** | isinya jabatan, bukan dokumen |
| satu "status" untuk kasus | **antrean · kode jabatan · hasil keputusan** | tiga field state terpisah |
| "riwayat" untuk before-image | **before-image** | riwayat ada di versi polis |
| `IsFacout` sebagai penanda fac out | — | `[terverifikasi]` isinya menguji **banding** |
| "Insert" untuk rule bernama `Insert…` | **hapus-lalu-sisip-ulang** | dua di antaranya berisi `DELETE` |
| **"prorata" untuk `Local.prorate` pada jalur non-EDM** | **konstanta skala** | `[terverifikasi]` nilainya **selalu `100`** di jalur itu — tidak ada pembagian waktu |

### `Local.prorate` di jalur non-EDM bukan prorata

`[terverifikasi]` **Bukti — kedua varian `CountRateRetroCov`, alamat langkah dicantumkan:**

| Berkas | Alamat | Isi |
| --- | --- | --- |
| `D:\migrasi\RNM\DDL\CountRateRetroCov.xml` (FIRE, `ASM-FW-GISFW-Data-PropertyItem`) | `RH_2.pySteps(1)` `[3]` | `Local.prorate := 1` — menimpa dua rumus sebelumnya di langkah yang sama |
| idem | `RH_2.pySteps(4)` | gerbang `IsEDM` **T=3 / F=2** → berjalan saat **bukan** EDM; `Local.prorate := 100` |
| `D:\migrasi\RNM\DDL\CountRateRetroCov(ANEKA).xml` (ANEKA, `ASM-FW-GISFW-Data-Aneka`) | `RH_1.pySteps(1)` `[3]` · `RH_1.pySteps(4)` | **identik** |

Di jalur **NB/RNW (non-EDM)**, `Local.prorate` keluar bernilai **`100`** dan masuk ke
`@Math.divide((Share × Rate × Local.prorate), <pembagi>, 20)`. Angka `100` itu **menyatu dengan
pembagi** — `100000 ÷ 100 = 1.000` (‰, FIRE) dan `10000 ÷ 100 = 100` (%, ANEKA), persis skala rasio
per lini bisnis yang dikunci **K-018**. Menyebutnya "prorata" menyiratkan pembagian waktu yang
**tidak terjadi**, dan mendorong "perbaikan" berupa perhitungan periode yang akan **mengubah angka**.

📌 Di jalur **EDM** namanya baru berarti: `pySteps(2)` mengisi
`(ProrateEDMEnd + ProrateStartEDM) × 100`. Jadi properti yang sama **berperan berbeda per jalur** —
itu sebabnya satu nama tidak cukup.
