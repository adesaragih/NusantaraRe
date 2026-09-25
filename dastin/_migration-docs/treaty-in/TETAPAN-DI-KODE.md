# Tetapan yang ditulis langsung di kode — sapuan 24 September 2026

**Aturan yang melahirkannya:** *"Setiap tetapan yang ditulis langsung di kode adalah **calon**."*
Pola **cacat bersembunyi di balik nilai bawaan** sudah muncul empat kali di modul ini, dan pemilik
proses menaikkannya dari pengamatan menjadi **alat cari**.

Perkakas: `alat/sapu-tetapan-di-kode.py`. Ia melaporkan apa yang ditolaknya, sesuai aturan perkakas.

```
berkas disapu                                    : 708
Property-Set bernilai LITERAL ke properti bisnis : 449
--- yang DITOLAK, dan sebabnya ---
  sasaran bukan properti bisnis (Local/Param/Temp) : 843
  nilai kosong / benar-salah                       : 265
  nilainya bukan literal (rujukan atau ekspresi)   : 6.028
```

Dari 449, sebagian besar adalah perancah: `pxObjClass`, `IDPEGA`, pesan galat, penanda pencarian
layar, `"Select..."`. **Yang menyentuh arti bisnis ada di bawah**, dan dua di antaranya serius.

---

## 1. PERSENTASE RETRO DITULIS TETAP DI KODE — `FetchQSfromMasterXOL.xml`

```
Primary.SpreadingListXOL(<LAST>).ReinsTypeID       = "10007"
Primary.SpreadingListXOL(<LAST>).ReinsTypeName     = "ORS"
Primary.SpreadingListXOL(<LAST>).ParentReinsTypeID = "00"
Primary.SpreadingListXOL(<LAST>).Pct               = "15.00"
```

| | |
|---|---|
| **Langkah hidup?** | **56 dari 56 langkah hidup**, nol mati |
| **Terjangkau?** | ya — dipanggil `TreatyInXOLAddSpreading`, `…Detail`, `…DetailActual`, dan **`Section/Share.xml`** |
| **Apa artinya** | sebuah **bagian retro 15%** ke jenis reasuransi `ORS` (pengenal `10007`) **tidak datang dari tabel master** — ia ditulis di dalam aturan |

**Akibatnya tiga, dan ketiganya nyata:**

1. **Mengubah 15% itu menuntut mengubah kode.** Tabel acuan susunan retro tidak berdaya atasnya.
2. **INV-50 menjumlahkannya.** Persen penyebaran harus berjumlah 100; salah satu sukunya berasal
   dari kode, bukan dari data. Sebuah *constraint* yang menjumlahkan angka yang sebagiannya tidak
   dapat dilihat dari basis data adalah constraint yang tidak dapat ditelusuri ketika gagal.
3. **Ia menyamarkan perbedaan, persis seperti keempat instans sebelumnya.** Selama master memang
   menyisakan 15% untuk `ORS`, hasil "dari master" dan hasil "dari kode" **identik** — dan tidak ada
   yang pernah tahu mana yang berlaku.

> **Di sistem baru ia menjadi baris data di susunan retro, bukan baris kode.** Yang perlu dipastikan
> bisnis: apakah 15% ke `ORS` itu ketentuan yang berlaku umum, atau tambalan untuk keadaan tertentu.
> **Pertanyaan wawancara**, satu kalimat.

---

## 2. DUA NOMOR KONTRAK DITULIS TETAP SEBAGAI KONDISI

```
pyStepsPreCondParamsWhen = TreatyIn.ID=="1000951"      14 kemunculan
pyStepsPreCondParamsWhen = TreatyIn.ID=="1000069"       2 kemunculan
```

| | |
|---|---|
| **Label cocok dengan kondisi?** | **ya** — `pyStepsDescription` dan `pyStepsPreCondParamsWhen` berbunyi sama. Ini kasus langka di modul ini: labelnya jujur |
| **Langkah hidup?** | **ya** — `pyStepsBlockName` kosong pada seluruhnya |
| **Kapan dibuat** | 16 Februari 2023 |
| **Tercatat atas nama** | `Ade Samuel Tua Saragih` — orang yang sama dapat menjawab **kenapa**, dan itu jalan tercepat |

**Dua kontrak diperlakukan berbeda dari seluruh kontrak lain, di dalam aturan yang hidup.** Selama
itu berlaku, setiap kalimat berbentuk *"sistem lama melakukan X"* punya **dua pengecualian yang
tidak disebutkan di mana pun**.

**Yang dituntutnya:**

| Siapa | Apa |
|---|---|
| bisnis | apa yang berbeda pada kedua kontrak itu, dan apakah perbedaannya masih berlaku |
| migrasi | keputusan tertulis: kedua kontrak dipindahkan **apa adanya**, atau dinormalkan. Tidak boleh diputuskan diam-diam oleh skrip migrasi |
| sesi tiket | tidak ada tiket yang menyalin pengecualian ini. **Pengecualian per kontrak tidak dibawa** — bila perbedaannya nyata, ia harus menjadi **data**, bukan kondisi |

---

## 2b. `9989998` — KEGAGALAN YANG MENYAMAR SEBAGAI PERSENTASE

**Ditambahkan 24 September 2026, putaran penutup to-spec.** Ia tidak ditemukan sapuan literal di §1,
melainkan saat menelusuri penulis `ROLPct` untuk `TDA-17` — dan itu keterangan yang berguna:
**tetapan di dalam sebuah ekspresi tidak terlihat oleh sapuan yang mencari nilai literal polos.**

```
Activity/DetailCalculationROL.xml:
  .ROLPct = @toDecimal(@if(Local.TotalLimit==0, 9989998,
                           @divide(local.TotalPremi, Local.TotalLimit, 8))) * 100
```

| | |
|---|---|
| **Apa artinya** | ketika limitnya **nol**, sistem lama tidak menolak dan tidak mengosongkan — ia **menulis `9989998` sebagai persentase** |
| **Golongannya** | **ADR-0035** *"kegagalan bukan nilai"*. Ini instansnya yang paling telanjang di modul ini: angkanya bahkan tidak menyerupai persentase |
| **Akibat pada migrasi** | setiap baris warisan berlimit nol membawa `ROL_PCT = 9989998`. Migrasi yang memperlakukannya sebagai persentase **memindahkan kegagalan sebagai fakta**, dan rata-rata apa pun yang menghitungnya menjadi tak berarti |
| **Akibat pada DDL** | ia **tidak muat** di `NUMBER(11,8)` — tipe yang kelompok `P2` usulkan, yang hanya menampung tiga angka di depan koma. Salah satu dasar `KTV-A` menyeragamkan presisi angka |
| **Berapa barisnya** | belum diketahui. **`Uji AQ`** menghitungnya, dan melaporkannya **di baris tersendiri**, bukan sebagai penyimpangan |

> **Kolom yang dituju sedang dipertikaikan.** `PERSEN_ROL` berdiri di §10.3 atas putusan `TDA-17`,
> dan `F-15` membantahnya dengan rumus di atas. Apa pun keputusannya, **`9989998` tidak dibawa**:
> bila kolomnya bertahan, baris berlimit nol menjadi **kosong**, bukan sembilan juta.

---

## 3. Yang lain, diadili singkat

| Tetapan | Berkas | Putusan |
|---|---|---|
| `Reinstatement_List(…).ReinstatementPct = "100"` · `AdditionalPct = "100"` | `SetReinstatementPct` | **instans keempat** pola nilai bawaan; melahirkan entitas `PEMULIHAN_LIMIT` dan satu kemampuan **BARU** — `SPEC-MODEL-DATA.md` §10.3b |
| `DeductionList(1).Comment = "Comm to NuRe"` | `CalculateSharePctToRetro` | nama potongan sebagai teks — di model baru ia **rujukan tabel acuan** (§14.2). Tetapan ini hilang dengan sendirinya |
| `DeductionList(<APPEND>).Comment = "Brokerage fee"` / `"Facultative Brokerage fee"` | `TreatyEDMDifferenceDeduction` | sama |
| `ReinsTypeName = "ORS"` | `FetchQSfromMaster`, `…XOL` | nama disalin dari master — INV-59 sudah melarangnya; ia tidak disimpan |
| `ShareFacultativeReinsurers(…).Layer = "ALL LAYER"` | `GetSpreadingRetro` | **GEL-2**; dicatat agar tidak hilang saat gelombang itu dikerjakan |
| `CommentList(<LAST>).IsApproved = "Accept"` · `Suggest = "Create Revision"` | `AddCommentList_Act`, `SetTreatyIn_Act` | nilai enumerasi keadaan — sah sebagai tetapan, ia memang konstanta program |
| `SaveData.DEDUCTION1 = 0` · `DEDUCTION2 = 0` · `SPREAD_RNM_SHARE_VALUE = 0` | `SaveTreatyInDetail_Act` dan pasangan EDM-nya | penulisan **nol** ke kolom datar saat menyimpan. Sekeluarga dengan ADR-0035 *"kegagalan bukan nilai"*; tabel datar di model baru **diturunkan**, tidak ditulis aplikasi |
| `GetAchievement` menulis nol ke delapan ruas | `GetAchievement` | **GEL-3**, dan ia instans ADR-0035 yang **sudah** tercatat |

---

## 4. Cara menjalankan ulang

```
python alat/sapu-tetapan-di-kode.py
```

Keluarannya mencetak hitungan yang ditolak lebih dulu, lalu seluruh temuan per berkas. **Saringan
manual tetap diperlukan** — perkakas ini memisahkan literal dari ekspresi, ia tidak dapat memisahkan
konstanta program dari parameter bisnis. Itu pekerjaan yang membaca.
