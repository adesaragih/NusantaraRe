# Struktur Data Komite Sistem Lama — pohon halaman kerja sampai turunan terdalam

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Komite Claim Non Prop` — **59 berkas XML**, disapu seluruhnya (bukan cuplikan), termasuk `excludeXML/`.
> Dibangkitkan `alat/buat-pohon-komite.py` — sapuan **21 September 2026**.
>
> **Berkas ini memerikan SISTEM LAMA, bukan rancangan sistem baru.** Nama properti ditulis apa adanya seperti di XML, termasuk salah ejanya (`ComiteeClaim`, `Comitee`, `CurencyAdjustment`, `KomiteAproval`, `ReponseCode`, `CircumtansesCouseOfLoss`). Penerjemahan ke istilah `CONTEXT.md` terjadi di berkas lain, bukan di sini.
>
> **TURUNAN.** Ketiga CSV di §1.4 adalah pendamping mesin-baca berkas ini. Bila CSV berbeda dari sumber XML-nya, **alatnya yang salah**.

Pendamping mesin-baca: **tiga CSV**, lihat §1.4.
Melengkapi [`PENGETAHUAN.md`](../PENGETAHUAN.md) §7.1 dan §7.2, yang berhenti di **daftar properti tingkat satu** dan **sepuluh sasaran tulis-balik**. Berkas ini meneruskannya sampai **kedalaman 5** dan memberi kamus per simpul.

Sejajar dengan [`struktur-claimdata-lama.md`](../../claim-non-prop/4-erd-dan-tabel-datar/struktur-claimdata-lama.md) sisi Klaim. Keduanya memerikan pohon yang **bersinggungan tetapi tidak sama**: yang di sana berakar di `.ClaimData`, yang di sini berakar di halaman kerja Komite dan **melihat** `.ClaimData` dari tiga arah sekaligus (§5).

---

## 1. Cara baca

### 1.1 Notasi

| Tulisan | Arti |
|---|---|
| `.Prop` | properti tunggal (skalar) |
| `Prop[]` | **Page List** — baris berulang, dirujuk dengan indeks `(...)` di XML |
| `Prop{}` | **Page** — satu halaman bersarang, tanpa indeks |
| `→ KELAS` | kelas Pega yang menaungi halaman, diambil dari `pyPagesAndClassesPage`/`Class` di berkas Activity |

### 1.2 Tingkat keyakinan

Dua tingkat, dan bedanya penting.

| Nilai | Artinya |
|---|---|
| `jalur-penuh` | jalurnya terbaca utuh di XML, mis. `TempMainWork.ClaimData.AdjustmentList(n).ComiteeClaim(n).KomiteAproval`. Seluruh isi CSV bertingkat ini |
| `rujukan-relatif` | hanya terbaca sebagai `.Prop` di dalam repeat Section atau iterasi Activity; induknya disimpulkan dari sumber repeat itu. **Tidak masuk CSV**; disebut di badan berkas ini dengan penanda tegas |

Empat anggota `KomiteList[]` hanya ada di tingkat kedua — lihat §4.2. Itu satu-satunya tempat berkas ini memakai tingkat kedua.

### 1.3 Tipe data

Pega tidak mengekspor `Rule-Obj-Property`, dan `Rule-Obj-FieldValue` juga tidak ikut
(`INVENTARIS-BUKTI.md` §1, kolom "Apa yang TIDAK diliput"). **Tidak ada satu pun tipe yang
dinyatakan langsung oleh sistem lama**, dan berkas ini karena itu **tidak memuat kolom tipe
sama sekali** — berbeda dari berkas sejajarnya di sisi Klaim, yang memuatnya sebagai dugaan.

Yang dapat dipegang dari berkas ini adalah **bentuk pohon, nama, dan seberapa hidup tiap
simpul**. Presisi angka ada di `ADR-0003`; arti nilai ada di `CONTEXT.md`.

### 1.4 Tiga CSV pendamping

| Berkas | Baris | Isi | Kolom |
|---|---:|---|---|
| [`datar-komite-lama.csv`](datar-komite-lama.csv) | 688 | **pohon properti** — satu baris per simpul, seluruh halaman akar | `JALUR` · `AKAR` · `KEDALAMAN` · `NAMA` · `BENTUK` · `KELAS_PEGA` · `CACAH_ANAK` · `REF` · `CACAH_BERKAS` · `BERKAS_CONTOH` |
| [`datar-komite-kelas.csv`](datar-komite-kelas.csv) | 112 | **peta kelas** — satu baris per pasangan halaman–kelas; 84 halaman berbeda | `HALAMAN` · `KELAS_PEGA` · `CACAH_DEKLARASI` |
| [`datar-komite-tulis-balik.csv`](datar-komite-tulis-balik.csv) | 42 | **kontrak tulis-balik** — satu baris per pasangan jalur induk–berkas | `JALUR_INDUK` · `BERKAS` · `REF` |

Catatan pemakaian:

- `JALUR` memakai **nama kanonik tanpa indeks**. `TempMainWork.ClaimData.AdjustmentList.AcceptedNo` berarti `…AdjustmentList(n).AcceptedNo` di XML; `BENTUK = Page List` pada simpul induknya yang menandai perulangan.
- `REF` adalah cacah rujukan mentah, `CACAH_BERKAS` cacah berkas berbeda yang merujuknya. Simpul akar ber-`REF = 0` wajar: nama halamannya sendiri jarang ditulis tanpa properti di belakangnya.
- `BERKAS_CONTOH` memuat maksimal **tiga** berkas. Untuk daftar penuh, jalankan ulang alatnya.
- Kolom `KELAS_PEGA` **kosong untuk sebagian besar simpul**, dan itu bukan kelalaian: Pega hanya mendeklarasikan kelas untuk **halaman**, bukan untuk properti bersarang di dalamnya.

---

## 2. Ringkasan angka

| | |
|---|---:|
| Berkas XML disapu | **59** |
| Simpul dalam seluruh pohon | **688** |
| Halaman akar berbeda | **59** |
| Kedalaman maksimum | **5** tingkat |
| Simpul di bawah `pyWorkPage` (halaman kerja Komite) | **66** |
| Simpul di bawah `pyWorkCover` (klaim induk, jalur baca) | **64** |
| Simpul di bawah `TempMainWork` (klaim induk, jalur tulis) | **39** |
| Page List di seluruh pohon | **31** |
| Halaman yang dideklarasikan kelasnya | **84** |
| Kelas Pega berbeda | **32** |

**Simpul paling hidup**, sebagai penanda seberapa sering sesuatu disentuh — bukan penanda penting:

| Jalur | Ref | Berkas |
|---|---:|---:|
| `pyWorkCover.ClaimData` | 89 | 9 |
| `TempMainWork.ClaimData` | 54 | 4 |
| `pyWorkPage.AcceptStatus` | 54 | 4 |
| `pyWorkPage.Adjustment` | 42 | 6 |
| `TempKasir.CARI20` | 39 | 1 |

---

## 3. Peta kelas

**Tiga puluh dua kelas** muncul di 59 berkas ini, tetapi hanya **enam** yang menyusun data komite itu sendiri. Sisanya kelas integrasi, kelas pinjaman dari domain lain, dan kelas bawaan Pega.

| Kelas Pega | Dipakai sebagai | Catatan |
|---|---|---|
| `ASM-FW-GCNMFW-Work-KomiteTreatyNonProp` | `pyWorkPage` | **akar modul ini** — berkas sirkulasi |
| `ASM-FW-GCNMFW-Work-ClaimTreatyNonProp` | `TempMainWork`, `TempOpenPage` | klaim induk, dibuka untuk ditulis |
| `ASM-FW-GCNMFW-Data-Comitee` | `pyWorkPage.Komite` | **satu-satunya halaman komite yang berkelas sendiri**; `KomiteList[]` berkelas sama tetapi tidak dideklarasikan |
| `ASM-FW-GCNMFW-Data-Adjustment` | `Adjustment`, `Primary` | usulan yang diedarkan |
| `ASM-FW-GCNMFW-Data-ClaimData` | `ClaimData` | isi klaim |
| `ASM-FW-GCNMFW-Data-osAkseptasi` | `InputParamOs` | wadah muatan keluar akseptasi |

Sembilan kelas integrasi dan pinjaman: `Data-PaymentKasir` · `Int-DOCUMENT_CLAIM` · `Int-V_POLIS` · `Int-os_akseptasi_klaim` · `Int-BANKACCOUNT` · `Int-M_LINK_SERVICE` · `Int-TREATY_IN` · `Int-T_STORAGE_IMAGE` · `Int-policyjson`.
Enam kelas bawaan Pega: `Code-Pega-List` · `Code-Pega-PDF` · `Code-Pega-Process` · `Data-Admin-Operator-ID` · `Data-EmailAttachments` · `Data-WorkAttach-File` · `Link-Attachment`.

> **Konsekuensi untuk pemetaan tabel.** `ASM-FW-GCNMFW-Data-Comitee` dideklarasikan **hanya sekali**, untuk `pyWorkPage.Komite` — sebuah halaman berisi konteks kerugian (`CircumtansesCouseOfLoss`, `ExtentOfLoss`, `LegalLiability`, `Occupation`, `Remarks`), **bukan** berisi jenjang. Kelas yang sama menaungi `KomiteList[]`, `ComiteeClaim[]`, dan `ClaimComitee[]` di tempat lain, dan **di sanalah jenjang tinggal**. Satu kelas, dua arti yang tidak berhubungan. Memetakannya jadi satu tabel akan menggabungkan konteks kerugian dengan daftar pemutus.

---

## 4. Pohon `pyWorkPage` — berkas sirkulasi

**66 simpul, 18 skalar langsung, 6 kontainer, kedalaman 5.**

```
pyWorkPage{}                                   → GCNMFW-Work-KomiteTreatyNonProp
│
├─ 18 properti skalar (keputusan, penanda, jejak Pega)  ..................  §4.1
│
├─ KomiteList[]                                → (kelas tak dideklarasikan)   n=20  ★ daftar jenjang
│    ├─ KomiteID · KomiteAproval · KomiteComment · DateApproval · DateApprove
│    └─ KomiteEmail · IDKomite · KomitePost · Initial      ◄── rujukan-relatif, §4.2
│
├─ Komite{}                                    → GCNMFW-Data-Comitee          n=8
│    └─ CircumtansesCouseOfLoss · ExtentOfLoss · LegalLiability · Occupation · Remarks
│
├─ Adjustment{}                                → GCNMFW-Data-Adjustment       n=42  ★ usulan
│    ├─ AcceptedNo · PaymentType · IndexObject · CNPIndexInterim · XOLID
│    ├─ CNPAccNoAdjustF · CNPAccNoOtherF · CNPAccNoSalvage · CNPAccNoReinstate
│    ├─ IsProposeClose
│    └─ CNPLayerList · SpreadingRisk · AlokasiXOLPaid   ← disentuh sebagai list utuh, isinya tidak
│
├─ ClaimData{}                                 → GCNMFW-Data-ClaimData        n=26  ← SALINAN, §5
│    ├─ NoClaim · InsuredName · InsuredRelationship · Payable
│    ├─ TotalListClaimAmount · TotalListClaimAmountIDR · InterestList
│    ├─ PolicyData{} └─ EndDateTime
│    ├─ AdjustmentList[]
│    │    └─ CurencyAdjustment[]                                      ◄── KEDALAMAN 4
│    │         └─ CurrencyID
│    └─ ObjectList[]
│         └─ ObjectItemList[]
│              └─ Adjustment[]                                        ◄── KEDALAMAN 5
│                   └─ CurrencyID
│
├─ QuotationData{}                             → (tak dideklarasikan)         n=2
│    └─ InsuredName
└─ TreatyInMaster{}                            → (tak dideklarasikan)         n=2
     └─ ReportingStart
```

### 4.1 Skalar langsung di `pyWorkPage` (18)

**Keputusan komite** (4) — inilah satu-satunya masukan yang benar-benar milik modul ini
`AcceptStatus` · `Comment` · `IsSubjectivity` · `SubjectivityNote`

**Penentu giliran** (2) — dua penyimpan untuk satu fakta
`KomiteCount` · `KomiteLoop`

**Penanda jalur** (3)
`IsCloseFile` · `IsReject` · `IsPrevious`

**Jejak Pega** (8)
`pxCoverInsKey` · `pxCreateDateTime` · `pxCreateOpName` · `pxCreateOperator` · `pxLockHandle` · `pyID` · `pyWorkIDPrefix` · `pzInsKey`

**Sisa** (1)
`CLMNO`

> `pxCoverInsKey` (n=3) adalah **satu-satunya tali** dari berkas sirkulasi ke klaim induknya. Nomor sirkulasi yang ditulis balik ke klaim bukan tali — ia salinan, dan `ADR-0031` menetapkan tali yang sah dibaca dari sisi ini.

### 4.2 `KomiteList[]` — daftar jenjang, dan empat nama yang hanya terbaca relatif

Page List paling hidup di modul ini (n=20, 4 berkas). Satu baris = satu jenjang.

| Anggota | Tingkat keyakinan | Terbaca di |
|---|---|---|
| `KomiteID` | `jalur-penuh` | `KomitePostAdjustment`, `KomitePostAdjustmentCWP`, `KomiteRouter`, `GenerateAccCNP_act`, `SetKomiteList_Act` |
| `KomiteAproval` | `jalur-penuh` | `KomitePostAdjustment`, `KomitePostAdjustmentCWP`, `KomiteRouter`, `ShowTransfer` |
| `KomiteComment` | `jalur-penuh` | `KomitePostAdjustment`, `KomitePostAdjustmentCWP`, `ShowTransfer` |
| `DateApproval` | `jalur-penuh` | `KomitePostAdjustment`, `KomitePostAdjustmentCWP` |
| `DateApprove` | `jalur-penuh` | `KomitePostAdjustment`, `KomitePostAdjustmentCWP`, `GenerateAccCNP_act` |
| **`KomiteEmail`** | **`rujukan-relatif`** | `SendEmailKlaim_KMT` |
| **`IDKomite`** | **`rujukan-relatif`** | `SendEmailKlaim_KMT`, `SetKomiteList_Act` |
| **`KomitePost`** | **`rujukan-relatif`** | `SetKomiteList_Act`, `ShowTransfer` |
| **`Initial`** | **`rujukan-relatif`** | `SetKomiteList_Act`, `ShowTransfer` |

Sembilan anggota, sesuai `PENGETAHUAN.md` §7.1. Empat terakhir **tidak pernah ditulis dengan jalur penuh** di 59 berkas — hanya sebagai `.Prop` di dalam iterasi. Induknya pasti karena repeat-nya bersumber pada `KomiteList`, tetapi tingkat keyakinannya berbeda dan tidak diratakan di sini.

> **`DateApproval` dan `DateApprove` adalah dua kolom untuk satu fakta.** Keduanya terbaca dengan jalur penuh, di dua berkas yang sama, di modul yang sama. Bukan salah baca.

### 4.3 `Adjustment{}` — usulan yang sedang diedarkan (n=42, 6 berkas)

Salinan usulan dari klaim induk, dipakai layar dan activity pasca-keputusan. Sebelas properti terbaca dengan jalur penuh; tiga di antaranya (`CNPLayerList`, `SpreadingRisk`, `AlokasiXOLPaid`) **disentuh sebagai list utuh dan isinya tidak pernah dirujuk dari modul ini** — bentuk isinya hanya terbaca di sisi Klaim.

Empat properti bernama `CNPAccNo…` adalah **penanda nomor akseptasi per komponen** — biaya penilaian, biaya lain, salvage, pemulihan. Empat penanda, satu nomor.

---

## 5. Tiga nama untuk satu klaim — dan ini jebakan pemetaan yang paling mahal

Klaim induk muncul di modul ini lewat **tiga halaman berbeda**, dan ketiganya menunjuk berkas yang sama.

| Halaman | Simpul | Kelas | Perannya |
|---|---:|---|---|
| `pyWorkCover.ClaimData` | 64 | `Work-ClaimTreatyNonProp` | **jalur baca** — layar dan surat membacanya; ref 89 di 9 berkas, tertinggi di modul |
| `TempMainWork.ClaimData` | 39 | `Work-ClaimTreatyNonProp` | **jalur tulis** — dibuka, disunting, disimpan; §6 |
| `pyWorkPage.ClaimData` | 26 | `Data-ClaimData` | **salinan** di dalam berkas sirkulasi sendiri |

Ketiganya **berisi properti yang tumpang-tindih tetapi tidak identik**. `NoClaim` ada di ketiganya. `PolicyData.PolicyNo` ada di dua. `SpreadingRisk[]` hanya di `pyWorkCover`. `AdjustmentList[].ComiteeClaim[]` hanya di `TempMainWork`.

**Akibatnya untuk rancangan:** yang terlihat seperti "tiga tabel klaim" adalah **satu klaim dilihat dari tiga arah pada tiga saat berbeda** — saat dibaca, saat ditulis, dan saat difoto ke dalam berkas sirkulasi. Memetakan ketiganya menjadi tiga tabel akan melipatgandakan satu entitas; memetakannya menjadi satu tanpa mencatat **kapan** masing-masing dipakai akan menghapus fakta bahwa salinannya beku.

Hal yang sama berlaku pada usulan, yang muncul di **empat** tempat: `pyWorkPage.Adjustment{}`, halaman langkah `Adjustment{}`, alias `Primary{}`, dan `TempMainWork.ClaimData.AdjustmentList[]`. Ketiga yang pertama berkelas `Data-Adjustment` dan menunjuk hal yang sama; yang keempat adalah rumahnya.

`Primary{}` membawa tiga simpul yang **tidak muncul di mana pun lagi**: `StatusServiceKasir{}` beserta `ReponseCode` dan `ResponseMsg` — salah ejanya ada di sumbernya. Ia jejak integrasi kasir, dan hanya hidup di cabang ini.

---

## 6. Kontrak tulis-balik ke klaim induk

`TempMainWork` adalah klaim induk yang **dibuka untuk ditulis**. Inilah kontrak yang harus dipertahankan sistem baru — atau digantikan satu transaksi tunggal.

```
TempMainWork{}                                 → GCNMFW-Work-ClaimTreatyNonProp
│
├─ CNPStatusCase · IsCloseFile · stsReject · NamePIC       ← keadaan & penanda pada klaim
├─ pxCreateOperator · pxObjClass · pyID · pzInsKey          ← jejak Pega
│
├─ ClaimData{}                                                             n=54
│    ├─ NoClaim · InsuredName · DateOfLoss · ReportDate · Amount
│    ├─ DeductibleValue · ShareCeding · ListClaimAmount
│    ├─ IsSubjectivity · IsFInalAccXOL · ClaimComitee
│    ├─ PolicyData{} └─ PolicyNo · StartDateTime
│    └─ AdjustmentList[]                                                   n=29  ★ sasaran utama
│         ├─ AcceptanceStatus · AcceptedNo · AcceptedDate
│         ├─ IsKomite · IsSubjectivity · SubjectivityNote
│         └─ ComiteeClaim[]                                                n=9
│              └─ KomiteAproval · KomiteComment · DateApproval · DateApprove
│
└─ OfferFacIn{}
     └─ QuotationData{} └─ BusinessOldId
```

Empat hal yang perlu dibaca dengan teliti:

**`ComiteeClaim[]` adalah salinan `KomiteList[]` ke dalam usulan.** Nama berbeda, ejaan berbeda, isi sama — empat anggota yang ditulis balik adalah persis empat anggota keputusan. `CONTEXT.md` §6 sudah membedakannya dari `ClaimComitee`; berkas ini menambahkan yang ketiga: `KomiteList` di berkas sirkulasi, `ComiteeClaim` di usulan, `ClaimComitee` di klaim. **Tiga nama, dua ejaan salah, satu konsep.**

**`ClaimComitee` disentuh tanpa anak.** Ia ditulis sebagai list utuh (n=1); tidak ada satu pun anggotanya yang dirujuk dari modul ini. Bentuknya baru terlihat di sisi Klaim.

**`IsFInalAccXOL` — huruf besar di tengah kata ada di sumbernya.** Jangan dirapikan saat memetakan; nama itu yang ada di XML.

**`CNPStatusCase`, `IsCloseFile`, dan `stsReject` menempel di `TempMainWork` langsung**, bukan di `ClaimData`. Keadaan klaim dan isi klaim tinggal di dua tingkat berbeda.

---

## 7. Muatan keluar — bentuk berubah tiga kali

Pohon di atas adalah bentuk **di memori**. Saat keluar dari Pega, bentuknya berubah.

### 7.1 `InputParamOs` — akseptasi outstanding (17 field)

Kelas `ASM-FW-GCNMFW-Data-osAkseptasi`. Wadah yang diisi lalu dikirim ke prosedur akseptasi.

`NoClaim` · `IDMasterTreaty` · `CauseOfLoss` · `CauseOfLossID` · `TypeLoss` · `TypeLossID` · `Currency` · `CurrencyID` · `KursValue` · `GrossValue` · `Value` · `Adjusterfee` · `Salvage` · `CNPOthersFee` · `PersenRNM` · `Type` · `EstimationDate`

Tujuh belas field, sama persis dengan yang dicatat sisi Klaim. **Jauh lebih sempit dari sumbernya** — pohon 688 simpul menyempit jadi 17 saat keluar.

> `Adjusterfee` dieja dengan `f` kecil di sini, sementara di `SpreadingRisk[]` sisi Klaim ia `AdjusterFee`. Satu field, dua ejaan, tergantung wadahnya.

### 7.2 `TempKasir` — muatan kiriman pembayaran (25 field, seluruhnya bernama `CARI…`)

`CARI1` … `CARI24` ditambah `CARIDATETIME`. **Nama-nama ini tidak mengatakan apa-apa** — artinya hanya terbaca dari urutan pengisian di `HitServiceToKasirKMT_Act`, bukan dari namanya.

Dua yang paling hidup: `CARI20` (n=39) dan `CARI19` (n=30).

**Isi muatan ini berpagar** (`PG-04`; `PAGAR-05` dan `PAGAR-06` pada `SPEC-KOMITE-01.md`). Berkas ini mencatat **bentuknya** — 25 field bernama `CARI…` — dan berhenti di situ. Arti tiap field tidak disimpulkan, karena pembandingnya adalah nol baris `POOLDATA.DIRECTTOKASIR_LOG` yang belum pernah diambil (`INVENTARIS-BUKTI.md` §2.5 baris 1).

---

## 8. Apa yang tidak ada di berkas ini, dan sebabnya

**Tipe data.** Tidak ada, dan tidak akan ada sampai `Rule-Obj-Property` diekspor. Berkas sejajarnya di sisi Klaim memuat kolom tipe sebagai **dugaan**; berkas ini memilih tidak menebak.

**Arti nilai.** `.Type`, `TransferType`, `FlagProrate`, `KomiteAproval` — keempatnya terbukti dipakai, artinya tidak tertulis. `Rule-Obj-FieldValue` tidak ikut diekspor (`INVENTARIS-BUKTI.md` §1). `TransferType` terbaca di dua berkas, `KomiteRouter` dan `SendEmailKlaimRejectClose_KMT`.

**`DEGREE` dan `LIMIT_BOTTOM`.** **Nol kemunculan** di 59 berkas ini. Keduanya hidup di `FilterEmailKomiteWithLimit`, yang berada di folder **Claim**, bukan folder Komite (`INVENTARIS-BUKTI.md` §2.4). Pohon roster karena itu tidak dapat disusun dari ekspor ini saja.

**Isi tabel roster.** `EMAILKOMITE` DDL-nya dipegang, isinya nol baris (`INVENTARIS-BUKTI.md` §2.5 baris 3). Kueri pembukanya ada di `KUERI-DBA-01.md` nomor 3.

**Rancangan sistem baru.** Tidak ada satu baris pun. Model data sistem baru hidup di `SPEC-KOMITE-01.md` bagian *Model data*; pemetaan lama ke baru hidup di bagian *Ketertelusuran*. Berkas ini hanya memerikan yang lama.
