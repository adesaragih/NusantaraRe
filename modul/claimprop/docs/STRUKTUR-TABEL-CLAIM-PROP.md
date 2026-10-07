# Struktur Tabel — Claim Prop (Klaim Treaty Inward Proporsional)

> ⭐ `[keputusan work owner]` **2026-09-20.** **Awalan tabel klaim non-life yang lama — berhuruf `P` sesudah `CLAIM` — diganti menjadi `T_CLAIM_`**, sebab tabelnya kini **dipakai bersama lini FAC dan PROP** sehingga huruf **P** pada awalan menyesatkan. ⭐ **NONPROP nanti ikut tabel yang sama.**
>
> ⛔ **Nama awalan lamanya sengaja TIDAK dikutip harfiah** di berkas mana pun di luar lini Life — ⭐ supaya pencarian atas awalan lama itu **hanya** menemukan berkas yang memang belum diselaraskan, bukan kalimat yang menerangkan penggantiannya. ⭐ **Tujuh tabel dipakai bersama; lima di antaranya berkunci asing GANDA.** ⛔ **Tabel lini Life tidak disentuh** — awalannya berbeda dan tidak ikut terganti.
>
> ⚠️ Akibatnya **dua berkas lini Life masih menyebut nama tabel penyesuaian yang lama**. ⭐ Itu **dicatat sebagai `[terbuka]`**, ⛔ **bukan diperbaiki**.

> ⚠️ **RALAT 07-10-2026** — kalimat lama: *"Tujuh tabel dipakai bersama; lima di antaranya berkunci asing GANDA."* Tabel bersama berjumlah **lima**, sesuai badan berkas ini: `T_WORK_CLAIM`, `T_GENERAL_CLAIM`,
> `T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST`, `T_VIEW_SUGGEST`. "Tujuh" di kalimat ini basi.


**Tanggal:** 2026-09-19 · **Modul korpus:** `Claim Prop` (329 berkas) · Korpus **READ-ONLY**

Acuan bentuk tabel untuk aplikasi Go. Berkas ini menggambarkan **BENTUK**, bukan alasan — alasannya
ada di `spec.md`, `issues/*.md`, dan berkas grilling.

> ⚠️ **Berkas ini adalah acuan TUNGGAL nama kolom Claim Prop.** Seluruh nama kolom **snake_case
> huruf besar**, dibuat **sesuai isinya** — **AC 106**. Ejaan yang muncul di `spec.md` dan berkas
> grilling adalah **ejaan korpus**: itu **bukti asal kolom**, **bukan** nama kolom.
>
> ⛔ **Tidak ada `CREATE TABLE` di berkas ini.** Bentuk kolom ditulis dalam kalimat.
> ⛔ **Nol nomor baris XML dikutip** — bukti memakai path berkas + nama rule + nomor langkah Pega.
> ⛔ **Tidak satu pun butir `[terbuka]` dinyatakan tertutup.**

**Tipe ditulis sebagai kategori logis:** teks · angka desimal · bilangan bulat · DATE.

**Sumber tiap kolom** salah satu dari dua, tidak boleh kosong:

- **korpus** — nama rule + nomor langkah Pega (atau nama layar + selnya)
- **keputusan** — `[keputusan work owner]` + tanggalnya

⚠️ Seluruh kolom **nullable** kecuali PK dan yang disebut **NOT NULL**; wajib-isi ditegakkan di Go.

---

## TUGAS 1 — Nama tabel, dan dari mana asalnya

### ⚠️ Fakta yang wajib diketahui sebelum membaca berkas ini

`[keputusan work owner]` Kesepuluh nama tabel `T_CLAIM_*` di bawah **hanya ada di diagram milik
work owner** (`Diagram-Skema-Tabel-NusantaraRe.xlsx`, sheet *Claim Prop*). ⛔ Nama-nama itu
**tidak ada** di `spec.md`, **tidak ada** di `issues/*.md`, dan **tidak ada** di berkas grilling
mana pun. **Berkas inilah yang pertama kali menuliskannya ke dalam dokumen.**

⛔ **Nol nama dikarang sendiri. Nol nama diganti.** Dipakai apa adanya.

### Sepuluh tabel milik Claim Prop

| Tabel | Asal halaman Pega | Catatan |
| --- | --- | --- |
| `T_CLAIM_ESTIMATION` | `ClaimData.EstimationList` | |
| `T_CLAIM_INTEREST` | `ClaimData.InterestList` | objek pertanggungan |
| `T_CLAIM_CLAIM_AMOUNT` | `ClaimData.ListClaimAmount` | |
| `T_CLAIM_SPREADING` | `ClaimData.SpreadingClaim` | **+ `TREATY_ID`** |
| `T_CLAIM_BREAK_QS` | `ClaimData.SpreadingBreakQS` | **+ `TREATY_ID`** · ⚠️ **SEJAJAR**, bukan anak `_SPREADING` |
| `T_CLAIM_FAC_RETRO` | `ClaimData.FacRetroList` | |
| `T_CLAIM_ADJUSTMENT` | `ClaimData.AdjustmentList` **+ `DataCommitteeTreaty`** | 7 kolom `KOMITE_*` — lihat K7 |
| `T_CLAIM_ADJ_SPREADING` | `AdjustmentList(n).SpreadingAdjustment` | **+ `TREATY_ID`** |
| `T_CLAIM_ADJ_QUOTA_SHARE` | `AdjustmentList(n).SpreadingQuotaShare` | **+ `TREATY_ID`** · ⚠️ **SEJAJAR** |
| `T_CLAIM_ADJ_LOSS_ALLOCATION` | `AdjustmentList(n).LossAllocation` | |

### ⚠️ Keberatan nama — **SATU**

> `[terbuka]` **`T_CLAIM_BREAK_QS`** — namanya menyebut *"break QS"*, sedangkan halaman asalnya
> `SpreadingBreakQS` berisi **baris spreading quota-share pada sisi klaim**, sejajar dengan
> `T_CLAIM_SPREADING`. Pasangannya di sisi penyesuaian dinamai lengkap
> `T_CLAIM_ADJ_QUOTA_SHARE`. ⚠️ Dua tabel berperan sama diberi dua gaya nama yang berbeda —
> pembaca baru mudah mengira `_BREAK_QS` sesuatu yang lain.
> ⛔ **Nama TIDAK saya ganti.** Diserahkan ke work owner.

⛔ Sembilan nama lain: **tidak ada keberatan**.

### Lima tabel dipakai bersama — ⛔ TIDAK didefinisikan ulang di sini

#### T3 · `T_WORK_CLAIM` — aturan yang mengikat lini PROP

| Aturan | Isi | Sumber |
| --- | --- | --- |
| `ID` baris **klaim** | berawalan **`CLMP-`** | `[keputusan work owner 2026-09-18]` |
| `ID` baris **komite** | berawalan **`TKMT-`** | `[keputusan work owner 2026-09-18]` |
| `COVER_KEY` pada baris komite | **WAJIB terisi** — `ID` baris klaim induknya | `[keputusan work owner 2026-09-18]` |
| `LINI` | **`PROP`** pada kedua baris | `[keputusan work owner 2026-09-18]` |

⚠️ `T_WORK_CLAIM` **lintas-lini**, jadi `ID`-nya menyebut **delapan awalan** — sepasang per lini:

| LINI | baris klaim | baris komite |
| --- | --- | --- |
| **FAC** | `CLM-` | `KMT-` |
| **PROP** | **`CLMP-`** | **`TKMT-`** |
| **NONPROP** | `CLMNP-` | `KMTNP-` |
| **LIFE** | `CLMLF-` | `KMTLF-` |

#### T4 · `T_GENERAL_CLAIM` — **shared PK** dengan `T_WORK_CLAIM`

⛔ Definisi penuhnya **bukan milik berkas ini**. Yang mengikat lini PROP hanya: `ID`-nya **sama
persis** dengan `T_WORK_CLAIM.ID` baris klaim (`CLMP-…`), tanpa kolom penyambung.

⚠️ **Kolom khas satu lini WAJIB nullable.** Yang **khas PROP** dan karena itu wajib nullable:

| Kolom khas PROP | Isinya | Sumber |
| --- | --- | --- |
| `TREATY_GROUP_ID` | kelompok treaty — kunci tunggal klasifikasi lini bisnis | korpus — `Activity/SetValueToClaim_Act.xml` langkah **8** |
| `SHARE_CEDING` | persen share ceding klaim | korpus — `Activity/AddListClaimAmount.xml` langkah **1** |
| `NET_DEDUCTIBLE_VALUE` | nilai deductible bersih tingkat klaim | korpus — `Activity/SetFormat_Act.xml`, catatan pengembang *"add pyWorkPage.ClaimData.NetDeductibleValue"* |
| `DLA_NO_CEDING` | nomor DLA pihak ceding | korpus — `Activity/SetDLACedingSOB.xml` langkah **1** |
| `DLA_NO_SOB` | nomor DLA SOB | korpus — `Activity/SetDLACedingSOB.xml` langkah **2** |
| `PERIOD_POLICY_TBA` | penanda periode polis **belum pasti** | korpus — `Activity/AddAdjustment_Act.xml`; dipakai `Section/InputAcceptation.xml` dan `Section/OutstandingClaim.xml` |
| `IS_CLOSE_FILE` | penanda usul tutup berkas, ditulis modul Komite | korpus — `Komite Claim Prop/Activity/KomitePostAdjustment.xml` langkah **11**, tanpa gerbang |
| `IS_RESERVED_CLAIM` | penanda usul klaim dicadangkan, ditulis modul Komite | korpus — **KomitePostAdjustment** langkah **11**, tanpa gerbang |
| ⛔ ~~`FLAG_ON_GOING_COMMITTEE`~~ **DIBUANG** | ~~penanda komite sedang berjalan~~ — **tidak dibawa ke sistem baru**; layar akseptasi **menurunkannya sendiri** dari kasus komite yang berjalan, sehingga satu fakta tidak punya dua sumber kebenaran. `[keputusan work owner]` **2026-09-19** · lihat `spec.md` **14c** dan **AC 129**. ⚠️ **Nilai lamanya tidak dimigrasikan.** | korpus — `Activity/AddKomiteTreatyChild_ACT.xml` langkah **3**; ditulis juga **KomitePostAdjustment** langkah **26.1** dan **27** — **jejak bahwa kolom ini pernah ada di Pega sengaja dibiarkan terbaca** |
| `IS_ANY_ACCEPTATION` | penanda sudah ada akseptasi | korpus — **KomitePostAdjustment** langkah **14** |
| `IS_SUBJECTIVITY` | penanda persetujuan bersyarat | korpus — **KomitePostAdjustment** langkah **24** |

⛔ **Definisi Life tidak disalin.** Kesebelas kolom di atas adalah **tambahan khas PROP** pada tabel
lintas-lini, bukan definisi ulang tabelnya. ⚠️ **Sepuluh di antaranya dibawa; satu —
`FLAG_ON_GOING_COMMITTEE` — DIBUANG** `[keputusan work owner]` 2026-09-19.

> ⛔ **RALAT 2026-09-19** — kalimat lamanya **dikutip, tidak dihapus**: *"⛔ **Definisi Life tidak
> disalin.** Kesebelas kolom di atas adalah **tambahan khas PROP** pada tabel lintas-lini, bukan
> definisi ulang tabelnya."* Jumlah kolom **yang dibawa** kini **10**, bukan 11. ⛔ **Barisnya tidak
> dihapus** — ia ditandai DIBUANG supaya jejak bahwa kolom itu pernah ada di Pega tetap terbaca.

#### `T_VIEW_SUGGEST` · `T_GENERAL_KOMITE` · `T_KOMITE_KOMITELIST`

| Tabel | Aturan yang mengikat lini PROP |
| --- | --- |
| `T_VIEW_SUGGEST` | ⚠️ **tidak ketemu** di korpus `Claim Prop` maupun `Komite Claim Prop` — disisir seluruh berkas, penyaringan **tidak peka huruf besar-kecil** (aturan **I1**). `[terbuka]` apa perannya bagi lini PROP |
| `T_GENERAL_KOMITE` | ⛔ **sudah terkunci** di `STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md`. Yang mengikat sisi klaim: `ADJUSTMENT_ID`-nya menunjuk **`T_CLAIM_ADJUSTMENT.ID`** untuk lini PROP, **NOT NULL**, index **UNIK**, **tanpa `REFERENCES`** |
| `T_KOMITE_KOMITELIST` | ⛔ **sudah terkunci**. **Nol** hubungan langsung ke tabel Claim Prop — ia anak `T_GENERAL_KOMITE` |

> ⚠️ **RALAT 07-10-2026** — kalimat lama: *"`T_VIEW_SUGGEST` — tidak ketemu di korpus `Claim Prop` maupun `Komite Claim Prop` … `[terbuka]` apa perannya bagi lini PROP"* `T_VIEW_SUGGEST` **ketemu** sebagai halaman `SuggestList` (78 kemunculan di 4 berkas; RELASI §B6).
> Kolom `CLAIM_ID`-nya dipasang (migrasi `532`, keputusan work owner 07-10-2026 "1 tabel aja gabung life dan non life").

### T5 · ⚠️ Tabel kesebelas: jejak audit — **BELUM PUNYA NAMA**

**AC 4** menuntut: *"Jejak audit menjadi tabel milik Claim Prop sendiri, bukan menumpang struktur
modul lain."* ⛔ **Tabel itu belum ada di diagram mana pun.**

> ## ⛔ NAMANYA TIDAK SAYA TETAPKAN
>
> `[terbuka]` **Nama tabel jejak audit Claim Prop belum ditetapkan.** Menunggu **work owner**.
>
> **Usulan penamaan — *ini penamaan, bukan rancangan*:** **`T_CLAIM_AUDIT_TRAIL`**.
> Alasannya satu baris: mengikuti awalan `T_CLAIM_` yang sudah dipakai sepuluh tabel lain,
> dan kata terakhirnya menyebut isinya.
> ⛔ **Jangan dianggap ditetapkan.** Selama work owner belum menjawab, tabel ini **tidak punya
> nama** dan kolom-kolomnya di bawah menggantung bersamanya.

**Isinya menurut tiket 14 dan korpus** — kolom sudah punya sumber, **hanya namanya yang belum**:

| Kolom | Tipe | Null | Isinya | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | **tidak** | PK | `[keputusan work owner]` — struktural |
| `CLAIM_ID` | teks | **tidak** | FK → `T_WORK_CLAIM.ID` baris klaim | `[keputusan work owner]` — struktural |
| `ACTION` | teks | ya | jenis aksi | korpus — `DataTransform/InsertChronology_DT.xml` langkah **1.1.1**; ejaan Pega slot memo kronologi |
| `ACTOR_NAME` | teks | ya | **nama pelaku**, kolom terpisah | korpus — `Activity/SethistoryKlaimTreaty.xml`; ⚠️ **AC 78** — di Pega tercampur ke properti tingkat wewenang |
| `AUTHORITY_LEVEL` | teks *(enum tertutup)* | ya | tingkat wewenang | korpus — **InsertChronology_DT** langkah **1.1.4–1.1.7**; ⚠️ ejaan Pega properti berawalan `Is…` yang **bukan boolean** — **AC 78** |
| `EVENT_TIME` | **DATE** | ya | waktu aksi | korpus — **InsertChronology_DT**; ⚠️ **AC 82** — riwayat wajib diurutkan dari kolom ini, **bukan** urutan baris |
| `CLOSING_MEMO` | teks | ya | memo penutupan klaim | korpus — `Activity/CloseClaimProp.xml` langkah **5**; ⚠️ **AC 81** — di Pega ditulis ke slot yang **tidak pernah dibaca** |

⚠️ **Jenis peristiwa yang terbaca dari korpus** dan wajib tertampung: *"Add Claim Amount"* ·
*"Delete Loss Allocation"* (`Activity/RemoveLossAlloction_act.xml` langkah **2**) ·
*"Send to Commite"* (`AddKomiteTreatyChild_ACT` langkah **3**) · *"Add type and Change % Loss
Allocation"* dan *"Edit type and Change % Loss Allocation"* (`Activity/CountPersen_act.xml`
langkah **8** dan **9**) · *"Count Value Adjustment"* (`Activity/SetNameCurrency_Act.xml`
langkah **8**).

---

## TUGAS 2 — Kolom tiap tabel

> ⚠️ **RALAT 07-10-2026** — empat hal di bab ini basi:
> 1. **Induk FK anak**: setiap *"`CLAIM_ID` → `T_WORK_CLAIM.ID`"* di bawah = **`T_GENERAL_CLAIM.ID`** (diagram sheet
>    Claim Prop; RELASI relasi 3–10; migrasi 521–531).
> 2. **Ketelitian**: *"20 digit seluruhnya, 8 di antaranya di belakang koma"* = **`NUMBER(38,10)`** (keputusan work
>    owner 07-10-2026).
> 3. **Kurs per baris**: tabel spreading tanpa kolom kurs memang benar menurut XML — lihat RALAT AC 25 di spec.
> 4. `T_CLAIM_FAC_RETRO.TOTAL_ESTIMATION_REINS` dibuat seperti diagram, tetapi pada kasus baru **tetap kosong**: satu-
>    satunya penulis `ClaimData.FacRetroList.TotalEstimasiReas` adalah `AddKomiteTreatyChild_ACT` 8.1 (penyerahan komite =
>    OQ-CP-16); `PrintDLATreatyIn` 15.13 menulis salinan `AdjustmentList(n).FacRetroList`, bukan daftar ini. Kasus lama
>    dari `JSON_KLAIM` membawa nilainya lewat pemuat.
>
> Bentuk yang berlaku = **Lampiran pengikat penjaga** di akhir berkas (dibandingkan penjaga dengan DDL).

### Ketelitian angka — berlaku untuk seluruh tabel di bawah

`[keputusan work owner]`:

| Golongan | Bentuknya |
| --- | --- |
| **Uang · persen · kurs** | **angka desimal**, **20 digit seluruhnya, 8 di antaranya di belakang koma**. Dihitung di Go sampai **20 angka di belakang koma** dengan **nol pembulatan di tengah jalan**; **tampilan layar 4 angka di belakang koma**. **Tipe desimal, BUKAN bilangan pecahan biner** (**ADR-0003**) |
| **Pencacah** | **bilangan bulat biasa** — tidak ikut aturan di atas |

⚠️ **Pembulatan terjadi di batas penyimpanan, dan itu diterima sadar.**

### K5 · Mata uang **per baris** — **AC 25 / AC 27**

> ⛔ **Invariant "satu klaim satu mata uang" milik Claim Life JANGAN ikut disalin.**
> Di Claim Prop mata uang ditetapkan **per baris** — `Activity/CurencyEstimation_Act.xml`.
> **Setiap tabel yang memuat uang punya kolom mata uang dan kursnya sendiri.**

> ⚠️ **RALAT 07-10-2026** — kalimat lama: *"Setiap tabel yang memuat uang punya kolom mata uang dan kursnya sendiri."* Diputuskan dari XML: kurs per baris HANYA ditulis pada `EstimationList.KursValue`
> (`CurencyEstimation_Act` 5), `AdjustmentList.KursIDR` (`SetNameCurrency_Act` 6), `InterestList.KursObjectItem`
> (`SetCurencyInterest_act` 5), `ListClaimAmount.IDR` (`SetCurencyList_act` 5) dan `SpreadingRisk.PremiumSpreaded`
> (`SetCurrency_Act` 6; `LossAllocation` salinannya). Baris spreading — `SpreadingClaim`, `SpreadingBreakQS`,
> `AdjustmentList(n).SpreadingAdjustment` / `.SpreadingQuotaShare` — **tidak pernah** diberi kurs (hanya
> `CurrencyID` / `Currency`), maka keempat tabel spreading tanpa kolom kurs (migrasi 525, 526, 529, 530).

---

### 1 · `T_CLAIM_ESTIMATION` — estimasi klaim

| Kolom | Tipe | Null | Isinya | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | **tidak** | PK | `[keputusan work owner]` — struktural |
| `CLAIM_ID` | teks | **tidak** | FK → `T_WORK_CLAIM.ID` | `[keputusan work owner]` — struktural |
| `ESTIMATION_DATE` | **DATE** | ya | tanggal estimasi | korpus — `Activity/AddEstimation_Act.xml` langkah **5.3** *(Pega `EstimationDate`)* |
| `TYPE_LOSS_ID` | teks | ya | kode jenis kerugian | korpus — **AddEstimation_Act** langkah **5.3** *(Pega `TypeLossID`)* |
| `TYPE_LOSS_NAME` | teks | ya | nama jenis kerugian | korpus — **AddEstimation_Act** langkah **5.3** *(Pega `TypeLoss`)* |
| `CURRENCY_ID` | teks | ya | kode mata uang **baris ini** | korpus — **AddEstimation_Act** langkah **5.3** *(Pega `CurrencyID`)* |
| `CURRENCY_NAME` | teks | ya | nama mata uang | korpus — **AddEstimation_Act** langkah **5.3** *(Pega `Currency`)* |
| `KURS` | angka desimal | ya | kurs **baris ini** | korpus — **AddEstimation_Act** langkah **5.3** *(Pega `KursValue`)* |
| `ESTIMATION_VALUE` | angka desimal | ya | nilai estimasi | korpus — `Activity/CopyOldataCurr_act.xml` langkah **8.1** *(Pega `EstimationValue`)* |
| `ESTIMATION_VALUE_IDR` | angka desimal | ya | nilai estimasi dalam IDR | korpus — **CopyOldataCurr_act** langkah **8.1** *(Pega `ConvertValue`)* |
| `GROSS_ESTIMATION_PCT` | angka desimal | ya | persen estimasi kotor | korpus — **AddEstimation_Act** langkah **5.3** *(Pega `GrossEstimationPct`)* |
| `GROSS_ESTIMATION_IDR` | angka desimal | ya | estimasi kotor dalam IDR | korpus — **CopyOldataCurr_act** langkah **8.1** *(Pega `ConvertGrossEstimasi`)* |
| `RESERVE_VALUE` | angka desimal | ya | nilai cadangan | korpus — **CopyOldataCurr_act** langkah **8.1** · ⚠️ ejaan Pega `EstimastionReserve` **salah ketik**, dibetulkan di sini |
| `PERSEN_RNM` | angka desimal | ya | persen bagian RNM | korpus — `Activity/CountEstimation_Act.xml` langkah **3.1** *(Pega `PersenRNM`)* |
| `NO_PLA` | teks | ya | nomor PLA yang terbit dari baris ini | korpus — `Activity/TryMakePLA_Act.xml` langkah **18.1** *(Pega `NoPLA`)* |
| `IS_PRINT_FACE_CLAIM` | teks *(penanda)* | ya | penanda sudah dicetak | korpus — `Activity/ChangeData_Act.xml` langkah **2.1** *(Pega `PrintFaceClaim`)* |

⚠️ **Empat properti turunan tidak dijadikan kolom** — `TotalEstimasi` · `TotalGrossEstimasi` ·
`TotalGrossEstimasiIDR` · `IDR`/`Value`. Keempatnya **jumlah berjalan** yang ditulis balik ke baris
oleh `CountEstimation_Act` langkah **8.1** dan `CountTotalInsterest_Act` langkah **15.2.2**; di Go
mereka **dihitung, bukan disimpan** — sejalan dengan keputusan *"nilai turunan dihitung sekali
jalan dari basis asli"*. ⛔ Tidak ditulis sebagai kolom karena bukan data, melainkan hasil.

**Kunci tamu:** `CLAIM_ID` → `T_WORK_CLAIM.ID` · **ON DELETE CASCADE** · index **biasa**.

---

### 2 · `T_CLAIM_INTEREST` — objek pertanggungan

| Kolom | Tipe | Null | Isinya | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | **tidak** | PK | `[keputusan work owner]` — struktural |
| `CLAIM_ID` | teks | **tidak** | FK → `T_WORK_CLAIM.ID` | `[keputusan work owner]` — struktural |
| `OBJECT_NAME` | teks | ya | nama objek pertanggungan | korpus — `Activity/AddInterest_act.xml` langkah **1** *(Pega `ObjectName`)* |
| `CURRENCY_ID` | teks | ya | kode mata uang **objek ini** | korpus — `Activity/SetCurencyInterest_act.xml` langkah **1** *(Pega `CurrencyID`)* |
| `CURRENCY_NAME` | teks | ya | nama mata uang | korpus — **SetCurencyInterest_act** langkah **3.1** *(Pega `Currency`)* |
| `KURS` | angka desimal | ya | kurs **objek ini** — **AC 27** | korpus — **SetCurencyInterest_act** langkah **5** *(Pega `KursObjectItem`)* |
| `TSI_VALUE` | angka desimal | ya | nilai pertanggungan objek | korpus — `Activity/CountTotalInsterest_Act.xml` langkah **11.1** *(Pega `TSIPerObject`)* |
| `TSI_VALUE_IDR` | angka desimal | ya | nilai pertanggungan dalam IDR | korpus — **CountTotalInsterest_Act** langkah **7.1** *(Pega `TSIPerObjectIDR`)* |
| `IS_ADJ_VALUE` | teks *(penanda)* | ya | penanda ikut nilai penyesuaian | korpus — `Activity/SaveOutstanding_Act.xml` langkah **28.1** *(Pega `IsAdjVal`)* |

**Kunci tamu:** `CLAIM_ID` → `T_WORK_CLAIM.ID` · **ON DELETE CASCADE** · index **biasa**.

---

### 3 · `T_CLAIM_CLAIM_AMOUNT` — nilai klaim per mata uang

| Kolom | Tipe | Null | Isinya | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | **tidak** | PK | `[keputusan work owner]` — struktural |
| `CLAIM_ID` | teks | **tidak** | FK → `T_WORK_CLAIM.ID` | `[keputusan work owner]` — struktural |
| `CURRENCY_ID` | teks | ya | kode mata uang **baris ini** | korpus — `Activity/AddListClaimAmount.xml` langkah **6.1** *(Pega `CurrencyID`)* |
| `CURRENCY_NAME` | teks | ya | nama mata uang | korpus — **AddListClaimAmount** langkah **6.1** *(Pega `Currency`)* |
| `KURS` | angka desimal | ya | kurs **baris ini** | korpus — **AddListClaimAmount** langkah **6.1** *(Pega `KursObjectItem`, disalin dari objek pertanggungan)* |
| `CLAIM_AMOUNT` | angka desimal | ya | nilai klaim | korpus — **AddListClaimAmount** langkah **6.1** *(Pega `ClaimAmount`)* |
| `NET_DEDUCTIBLE_VALUE` | angka desimal | ya | nilai setelah deductible | korpus — **AddListClaimAmount** langkah **7.1** *(Pega `NetDeductibleValue`)* |
| `VALUE` | angka desimal | ya | nilai bagian RNM sesudah share ceding | korpus — **AddListClaimAmount** langkah **7.1** *(Pega `Value`)* |
| `VALUE_IDR` | angka desimal | ya | nilai dalam IDR | korpus — **AddListClaimAmount** langkah **7.1** · ⚠️ ejaan Pega **`USD`** — **nama berbohong**, isinya IDR; dibetulkan di sini |
| `KURS_IDR` | angka desimal | ya | kurs ke IDR | korpus — **AddListClaimAmount** langkah **6.1** *(Pega `IDR`)* |
| `NOTE` | teks | ya | catatan baris | korpus — `Activity/SaveOutstanding_Act.xml` langkah **24.1** *(Pega `Note`)* |

⚠️ **`USD` adalah alias yang berbohong** — nilainya `.Value * .IDR`, yaitu konversi ke **IDR**,
bukan USD. Sejalan dengan **AC 106**, namanya **dibetulkan**.

**Kunci tamu:** `CLAIM_ID` → `T_WORK_CLAIM.ID` · **ON DELETE CASCADE** · index **biasa**.

---

### 4 · `T_CLAIM_SPREADING` — spreading klaim per treaty

| Kolom | Tipe | Null | Isinya | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | **tidak** | PK | `[keputusan work owner]` — struktural |
| `CLAIM_ID` | teks | **tidak** | FK → `T_WORK_CLAIM.ID` | `[keputusan work owner]` — struktural |
| **`TREATY_ID`** | teks | ya | ⭐ penunjuk treaty — **tambahan diagram work owner** | `[keputusan work owner]` |
| `TREATY_NAME` | teks | ya | nama treaty | korpus — `Activity/SetTreatyNameSpreading_Act.xml` langkah **6.1** *(Pega `TreatyName`)* |
| `SHARE_PERCENTAGE` | angka desimal | ya | persen share reinsurer | korpus — `Activity/CountSpreading_Act.xml` *(Pega `SharePercentage`)* |
| `CLAIM_SPREADED` | angka desimal | ya | nilai klaim tersebar | korpus — `Activity/CopyOldataCurr_act.xml` langkah **9.1** *(Pega `ClaimSpreaded`)* |
| `CURRENCY_ID` | teks | ya | kode mata uang **baris ini** | korpus — **CopyOldataCurr_act** langkah **9.1** *(Pega `CurrencyID`)* |
| `CURRENCY_NAME` | teks | ya | nama mata uang | korpus — **CopyOldataCurr_act** langkah **9.1** *(Pega `Currency`)* |
| `TOTAL_SPREAD` | angka desimal | ya | total tersebar | korpus — `Activity/AddKomiteTreatyChild_ACT.xml` langkah **7.2** *(Pega `TotalSpread`)* |
| `IS_OLD_DATA` | teks *(penanda)* | ya | penanda baris warisan | korpus — `Activity/SaveOutstanding_Act.xml` langkah **26.1** *(Pega `IsOldData`)* |

⚠️ **Kurs tidak punya kolom sendiri di tabel ini** — tidak ketemu di `pyParamArray`,
`pyStepsCallParams`, maupun `pyStepsJavaSource` seluruh Activity. Sesuai aturan **I2** kalimatnya
dibatasi pada medan yang disisir. `[terbuka]` — lihat penutup.

**Kunci tamu:** `CLAIM_ID` → `T_WORK_CLAIM.ID` · **ON DELETE CASCADE** · index **biasa**.
`TREATY_ID` → master treaty · **di Go** · index **biasa**.

---

### 5 · `T_CLAIM_BREAK_QS` — spreading quota-share klaim ⚠️ **SEJAJAR**

> ⚠️ **Tabel ini SEJAJAR dengan `T_CLAIM_SPREADING`, bukan anaknya.** Keduanya anak langsung
> baris klaim. `[keputusan work owner]` — dan didukung korpus: `CopyOldataCurr_act` langkah **9**
> dan **10** mengulang **dua daftar berbeda** yang **sejajar**, bukan bersarang.

| Kolom | Tipe | Null | Isinya | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | **tidak** | PK | `[keputusan work owner]` — struktural |
| `CLAIM_ID` | teks | **tidak** | FK → `T_WORK_CLAIM.ID` | `[keputusan work owner]` — struktural |
| **`TREATY_ID`** | teks | ya | ⭐ penunjuk treaty — **tambahan diagram work owner** | `[keputusan work owner]` |
| `SHARE_PERCENTAGE` | angka desimal | ya | persen share | korpus — `Activity/CopyOldataCurr_act.xml` langkah **10.1** *(Pega `SharePercentage`)* |
| `CLAIM_SPREADED` | angka desimal | ya | nilai klaim tersebar | korpus — **CopyOldataCurr_act** langkah **10.1** *(Pega `ClaimSpreaded`)* |
| `CURRENCY_ID` | teks | ya | kode mata uang **baris ini** | korpus — **CopyOldataCurr_act** langkah **10.1** *(Pega `CurrencyID`)* |
| `CURRENCY_NAME` | teks | ya | nama mata uang | korpus — **CopyOldataCurr_act** langkah **10.1** *(Pega `Currency`)* |
| `TSI_SPREADED` | angka desimal | ya | nilai pertanggungan tersebar | korpus — `Activity/AddKomiteTreatyChild_ACT.xml` langkah **9.2** *(Pega `TSISpreaded`)* |
| `IS_OLD_DATA` | teks *(penanda)* | ya | penanda baris warisan | korpus — `Activity/SaveOutstanding_Act.xml` langkah **27.1** *(Pega `IsOldData`)* |

**Kunci tamu:** sama dengan tabel 4.

---

### 6 · `T_CLAIM_FAC_RETRO` — fakultatif retro

| Kolom | Tipe | Null | Isinya | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | **tidak** | PK | `[keputusan work owner]` — struktural |
| `CLAIM_ID` | teks | **tidak** | FK → `T_WORK_CLAIM.ID` | `[keputusan work owner]` — struktural |
| `REINSURER_ID` | teks | ya | kode reinsurer retro | korpus — `Activity/SaveAcceptationTreaty_Act.xml` langkah **9.3** *(Pega `ReinsurerID`)* |
| `REINSURER_NAME` | teks | ya | nama reinsurer retro | korpus — **SaveAcceptationTreaty_Act** langkah **9.3** *(Pega `ReinsurerName`)* |
| `SHARE_PCT` | angka desimal | ya | persen share retro | korpus — **SaveAcceptationTreaty_Act** langkah **9.3** *(Pega `PctShareAllObj`)* |
| `RI_COMMISSION_PCT` | angka desimal | ya | persen komisi reasuransi | korpus — **SaveAcceptationTreaty_Act** langkah **9.3** *(Pega `RiCommAllObj`)* |
| `ADDITIONAL_INFO` | teks | ya | keterangan tambahan | korpus — `Activity/SetTreatyNameSpreading_Act.xml` langkah **17.1** *(Pega `AdditionalInfo`)* |
| `TOTAL_ESTIMATION_REINS` | angka desimal | ya | total estimasi bagian retro | korpus — `Activity/AddKomiteTreatyChild_ACT.xml` langkah **8.1** *(Pega `TotalEstimasiReas`)* |
| `TSI_ALL_OBJECT` | angka desimal | ya | nilai pertanggungan seluruh objek | korpus — **AddKomiteTreatyChild_ACT** langkah **10.1** *(Pega `TFAllObj`)* |

⚠️ `TOTAL_ESTIMATION_REINS` juga ditulis `Activity/PrintDLATreatyIn.xml` langkah **15.13** dari
`local.hasilakhir`, **yang rumusnya belum ditelusuri**. `[terbuka]` — lihat penutup.

⚠️ **Mata uang dan kurs tidak ketemu** di medan yang disisir untuk tabel ini. `[terbuka]`.

**Kunci tamu:** `CLAIM_ID` → `T_WORK_CLAIM.ID` · **ON DELETE CASCADE** · index **biasa**.

⚠️ `[terbuka]` Di Pega `FacRetroList` muncul **dua tempat**: sebagai `ClaimData.FacRetroList`
*(tabel ini)* **dan** sebagai `AdjustmentList(n).FacRetroList` *(bukan tabel — lihat J1)*.

---

### 7 · `T_CLAIM_ADJUSTMENT` — baris penyesuaian ⭐ tabel terbesar

| Kolom | Tipe | Null | Isinya | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | **tidak** | PK | `[keputusan work owner]` — struktural |
| `CLAIM_ID` | teks | **tidak** | FK → `T_WORK_CLAIM.ID` | `[keputusan work owner]` — struktural |
| `KOMITE_ID` | teks | ya | ⭐ penunjuk balik ke baris komite · index **UNIK** | `[keputusan work owner]` — ⛔ **tidak terbukti korpus**; di Pega hanya ada penanda boleh/tidak |
| `ADJUSTMENT_TYPE` | teks | ya | jenis penyesuaian | korpus — `Activity/AddAdjustment_Act.xml` langkah **5** *(Pega `Type`)* |
| `CURRENCY_ID` | teks | ya | kode mata uang **baris ini** | korpus — `Activity/PrintFileAcceptance.xml` langkah **10** *(Pega `CurrencyID`)* |
| `CURRENCY_NAME` | teks | ya | nama mata uang | korpus — **PrintFileAcceptance** langkah **10** *(Pega `Currency`)* |
| `KURS` | angka desimal | ya | kurs **baris ini**, **dikunci saat pengisian pertama** | korpus — `Activity/SetNameCurrency_Act.xml` langkah **6** *(Pega `KursIDR`)* |
| `CURRENCY_ADJUSTMENT` | teks | ya | mata uang penyesuaian | korpus — **AddAdjustment_Act** langkah **5** *(Pega `CurencyAdjustment`, salah ketik di Pega)* |
| `GROSS_ADJUSTMENT` | angka desimal | ya | nilai kotor penyesuaian | korpus — `Activity/PrintDLATreatyIn.xml` langkah **15.12** *(Pega `GrossAdjustment`)* |
| `GROSS_VALUE` | angka desimal | ya | nilai kotor bagian RNM | korpus — `Activity/CountGrossAdjTreaty_Act.xml` langkah **3** dan **4** *(Pega `GrossValue`)* |
| `ADJUSTMENT_VALUE` | angka desimal | ya | nilai penyesuaian | korpus — `Activity/CountValueADJTreaty_Act.xml` langkah **12** *(Pega `AdjustmentValue`)* |
| `VALUE_ADJUSTMENT_IDR` | angka desimal | ya | nilai penyesuaian dalam IDR | korpus — `Activity/SetNameCurrency_Act.xml` langkah **8** *(Pega `ValueAdjustment`)* |
| `PROPOSE_ADJUSTMENT_VALUE` | angka desimal | ya | nilai penyesuaian yang diusulkan | korpus — **CountValueADJTreaty_Act** langkah **12** *(Pega `ProposeAdjustmentValue`)* |
| `PERSEN_RNM` | angka desimal | ya | persen bagian RNM | korpus — **AddAdjustment_Act** langkah **5** *(Pega `PersenRNM`)* |
| `INDIVIDUAL_RISK_TYPE` | teks | ya | jenis risiko perorangan | korpus — **AddAdjustment_Act** langkah **6** *(Pega `IndividualRiskType`)* |
| `INDIVIDUAL_RISK_PCT` | angka desimal | ya | persen risiko perorangan | korpus — **AddAdjustment_Act** langkah **6** *(Pega `IndividualRiskPercentage`)* |
| `INDIVIDUAL_RISK_VALUE` | angka desimal | ya | nilai risiko perorangan | korpus — **AddAdjustment_Act** langkah **6** *(Pega `IndividualRiskValue`)* |
| `INDIVIDUAL_RISK_RNM` | angka desimal | ya | nilai risiko perorangan bagian RNM | korpus — **CountValueADJTreaty_Act** langkah **12** *(Pega `IndividualRiskRNM`)* |
| `ADJUSTER_FEE_VALUE` | angka desimal | ya | nilai fee adjuster | korpus — **PrintFileAcceptance** langkah **10** *(Pega `AdjusterFeeValue`)* |
| `TOTAL_ESTIMATION_VALUE` | angka desimal | ya | total nilai estimasi baris ini | korpus — **AddAdjustment_Act** langkah **5** *(Pega `TotalEstimasiValue`)* |
| `PAYMENT_TYPE` | teks | ya | jenis pembayaran | korpus — **PrintFileAcceptance** langkah **10** *(Pega `PaymentType`)* |
| `IS_PAYABLE` | teks *(penanda)* | ya | dapat dibayar | korpus — **PrintFileAcceptance** langkah **10** *(Pega `Payable`)* |
| `PAYABLE_TO` | teks | ya | dibayarkan kepada | korpus — **PrintFileAcceptance** langkah **10** *(Pega `PayableTo`)* |
| `BANK_NAME` | teks | ya | nama bank | korpus — **PrintFileAcceptance** langkah **10** *(Pega `NameOfBank`)* |
| `BANK_BRANCH` | teks | ya | cabang bank | korpus — **PrintFileAcceptance** langkah **10** *(Pega `BranchOfBank`)* |
| `BANK_ACCOUNT_NO` | teks | ya | nomor rekening | korpus — **PrintFileAcceptance** langkah **10** *(Pega `NoAccount`)* |
| `SWIFT_CODE` | teks | ya | kode swift | korpus — `Komite Claim Prop/Section/ShowTransfer.xml`, sel kode swift *(Pega `SwiftCode`)* |
| `IS_DIRECT_TO_KASIR` | teks *(penanda)* | ya | penanda kirim langsung ke Kasir | korpus — **AddAdjustment_Act** langkah **5** *(Pega `DirectToKasir`)* |
| `ACCEPTED_NO` | teks | ya | nomor akseptasi | korpus — `Komite Claim Prop/Activity/KomitePostAdjustment.xml` langkah **16.9** *(Pega `AcceptedNo`)* |
| `ACCEPTED_DATE` | **DATE** | ya | tanggal akseptasi | korpus — **KomitePostAdjustment** langkah **16.9** *(Pega `AcceptedDate`)* |
| `ACCEPTANCE_STATUS` | teks | ya | status akseptasi · `1` setuju · `2` tolak | korpus — **KomitePostAdjustment** langkah **16.9** dan **25** *(Pega `AcceptanceStatus`)* |
| `IS_APPROVED` | teks *(penanda)* | ya | penanda disetujui komite | korpus — **KomitePostAdjustment** langkah **15** *(Pega `IsApproved`)* |
| `IS_KOMITE` | teks *(penanda)* | ya | penanda sedang/pernah di komite | korpus — `Activity/AddKomiteTreatyChild_ACT.xml` langkah **1** *(Pega `IsKomite`)* |
| `IS_SUBJECTIVITY` | teks *(penanda)* | ya | persetujuan bersyarat | korpus — **KomitePostAdjustment** langkah **24** *(Pega `IsSubjectivity`)* |
| `SUBJECTIVITY_NOTE` | teks | ya | catatan syarat | korpus — **KomitePostAdjustment** langkah **24** *(Pega `SubjectivityNote`)* |
| `NOTES` | teks | ya | catatan komite | korpus — **KomitePostAdjustment** langkah **27** *(Pega `Notes`)* |
| `DLA_NO` | teks | ya | nomor DLA | korpus — `Activity/PrintDLATreatyIn.xml` langkah **12** *(Pega `DLA_No`)* |
| `IS_PRINT_ACCEPT` | teks *(penanda)* | ya | penanda cetak akseptasi | korpus — **PrintFileAcceptance** langkah **3** *(Pega `IsPrintAccept`)* |
| `SALVAGE_VALUE` | angka desimal | ya | nilai salvage | korpus — `Komite Claim Prop/Section/ShowTransfer.xml` *(Pega `SalvageValue`)* · ⚠️ pembacanya di `PrintDLATreatyIn` langkah **15.11.2** **ber-remark** |
| `CREATED_BY` | teks | ya | operator pembuat | korpus — **AddAdjustment_Act** langkah **4** *(Pega `pxCreateOperator`)* |
| `CREATED_BY_NAME` | teks | ya | nama operator pembuat | korpus — **AddAdjustment_Act** langkah **4** *(Pega `pxCreateOpName`)* |
| `CREATED_AT` | **DATE** | ya | waktu pembuatan | korpus — **AddAdjustment_Act** langkah **4** *(Pega `pxCreateDateTime`)* |

#### K7 · ⚠️ Tujuh kolom `KOMITE_*` dari `DataCommitteeTreaty`

`[terverifikasi]` **Tujuh, dihitung ulang** dari sel layar `Section/ComiteeClaimTreaty.xml`,
`Harness/CommitteeTreaty.xml`, dan `Section/DtlDataCommitte.xml`:

| Kolom | Tipe | Null | Isinya | Sumber |
| --- | --- | --- | --- | --- |
| `KOMITE_CIRCUM_CAUSE_OF_LOSS` | teks | ya | uraian sebab kerugian | korpus — `Activity/ProteksiInitialandDate_Act.xml` langkah **1** *(Pega `CircumCauseOfLoss`)* |
| `KOMITE_ADJUSTER_FEE` | angka desimal | ya | fee adjuster menurut komite | korpus — `Activity/AddKomiteTreatyChild_ACT.xml` langkah **15** *(Pega `AdjusterFee`)* |
| `KOMITE_REMARKS` | teks | ya | catatan komite | korpus — **ProteksiInitialandDate_Act** langkah **1** *(Pega `Remarks`)* |
| `KOMITE_SALVAGE` | angka desimal | ya | nilai salvage menurut komite | korpus — **AddKomiteTreatyChild_ACT** langkah **15** *(Pega `Salvage`)* |
| `KOMITE_LEGAL_LIABILITY` | teks | ya | tanggung jawab hukum | korpus — **AddKomiteTreatyChild_ACT** langkah **15** *(Pega `LegalLiability`)* |
| `KOMITE_EXTENT_OF_LOSS` | teks | ya | luas kerugian | korpus — **AddKomiteTreatyChild_ACT** langkah **15** *(Pega `ExtentOfLoss`)* |
| `KOMITE_OCCUPATION` | teks | ya | okupasi tertanggung | korpus — sel layar `Section/ComiteeClaimTreaty.xml` dan `Harness/CommitteeTreaty.xml` *(Pega `Occupation`)* |

> ## ⚠️ JANGAN DIKELIRUKAN
>
> **Ketujuh kolom `KOMITE_*` di atas ada di `T_CLAIM_ADJUSTMENT` — tabel milik Claim Prop, berisi
> data komite yang MELEKAT pada baris penyesuaian. Ia BUKAN dan TIDAK ADA HUBUNGAN STRUKTUR dengan
> tujuh kolom `T_GENERAL_KOMITE`, yang milik modul Komite, sudah terkunci di berkas terpisah, dan
> berisi header kasus komite (`ID` · `ADJUSTMENT_ID` · `KOMITE_LOOP` · `KOMITE_COUNT` ·
> `ACCEPT_STATUS` · `CLOSE_FILE` · `RESERVED_CLAIM`). Dua tabel berbeda, dua modul berbeda,
> kebetulan sama-sama tujuh.**

**Kunci tamu `T_CLAIM_ADJUSTMENT`:**

| Kunci | Ke mana | ON DELETE | Index |
| --- | --- | --- | --- |
| `CLAIM_ID` | `T_WORK_CLAIM.ID` baris klaim | **CASCADE** | **biasa** |
| `KOMITE_ID` | `T_WORK_CLAIM.ID` baris komite | **di Go** *(penunjuk)* | ⭐ **UNIK**, nullable |

---

### 8 · `T_CLAIM_ADJ_SPREADING` — spreading penyesuaian

| Kolom | Tipe | Null | Isinya | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | **tidak** | PK | `[keputusan work owner]` — struktural |
| `ADJUSTMENT_ID` | teks | **tidak** | FK → `T_CLAIM_ADJUSTMENT.ID` | `[keputusan work owner]` — struktural |
| **`TREATY_ID`** | teks | ya | ⭐ penunjuk treaty — **tambahan diagram work owner** | `[keputusan work owner]` |
| `TREATY_NAME` | teks | ya | nama treaty | korpus — `Activity/SendEmailKlaim.xml` langkah **8.2.1** *(Pega `TreatyName`)* |
| `TREATY_TYPE` | teks | ya | jenis treaty | korpus — `Activity/SaveAcceptationTreaty_Act.xml` langkah **7.1** *(Pega `TreatyType`)* |
| `SHARE_PERCENTAGE` | angka desimal | ya | persen share reinsurer | korpus — `Activity/CountSpreadingADJ_Act.xml` langkah **6.1** *(Pega `SharePercentage`)* |
| `CLAIM_SPREADED` | angka desimal | ya | nilai penyesuaian tersebar | korpus — **CountSpreadingADJ_Act** langkah **6.2** *(Pega `ClaimSpreaded`)* |
| `TOTAL_SPREAD` | angka desimal | ya | total tersebar | korpus — **CountSpreadingADJ_Act** langkah **6.3** *(Pega `TotalSpread`)* |
| `CURRENCY_ID` | teks | ya | kode mata uang **baris ini** | korpus — `Activity/HitServiceToKasir_Act.xml` langkah **9.3** *(Pega `CurrencyID`)* |
| `CURRENCY_NAME` | teks | ya | nama mata uang | korpus — **SendEmailKlaim** langkah **8.2.1** *(Pega `Currency`)* |
| `PREMIUM_SPREADED` | angka desimal | ya | premi tersebar | korpus — **HitServiceToKasir_Act** langkah **9.3** *(Pega `PremiumSpreaded`)* |
| `BANK_ACCOUNT_NO` | teks | ya | nomor rekening penerima | korpus — **HitServiceToKasir_Act** langkah **9.3** *(Pega `NoAccount`)* |
| `BANK_ID` | teks | ya | kode bank penerima | korpus — **HitServiceToKasir_Act** langkah **9.3** *(Pega `IDOfBank`)* |

⚠️ **Kurs tidak punya kolom sendiri di tabel ini** — tidak ketemu di medan yang disisir.
`[terbuka]`.

**Kunci tamu:** `ADJUSTMENT_ID` → `T_CLAIM_ADJUSTMENT.ID` · **ON DELETE CASCADE** · index
**biasa**. `TREATY_ID` → master treaty · **di Go** · index **biasa**.

---

### 9 · `T_CLAIM_ADJ_QUOTA_SHARE` — quota-share penyesuaian ⚠️ **SEJAJAR**

> ⚠️ **SEJAJAR dengan `T_CLAIM_ADJ_SPREADING`, bukan anaknya.** Keduanya anak langsung baris
> penyesuaian. Didukung korpus: `CountSpreadingADJ_Act` langkah **6** dan **7** mengulang **dua
> daftar berbeda yang sejajar**.

| Kolom | Tipe | Null | Isinya | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | **tidak** | PK | `[keputusan work owner]` — struktural |
| `ADJUSTMENT_ID` | teks | **tidak** | FK → `T_CLAIM_ADJUSTMENT.ID` | `[keputusan work owner]` — struktural |
| **`TREATY_ID`** | teks | ya | ⭐ penunjuk treaty — **tambahan diagram work owner** | `[keputusan work owner]` |
| `SHARE_PERCENTAGE` | angka desimal | ya | persen share | korpus — `Activity/CountSpreadingADJ_Act.xml` langkah **7.1** *(Pega `SharePercentage`)* |
| `CLAIM_SPREADED` | angka desimal | ya | nilai tersebar | korpus — **CountSpreadingADJ_Act** langkah **7.2** *(Pega `ClaimSpreaded`)* |
| `TOTAL_SPREAD` | angka desimal | ya | total tersebar | korpus — **CountSpreadingADJ_Act** langkah **7.3** *(Pega `TotalSpread`)* |

⚠️ **Mata uang dan kurs tidak ketemu** di medan yang disisir untuk tabel ini. `[terbuka]`.

**Kunci tamu:** sama dengan tabel 8.

---

### 10 · `T_CLAIM_ADJ_LOSS_ALLOCATION` — alokasi kerugian penyesuaian

| Kolom | Tipe | Null | Isinya | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | **tidak** | PK | `[keputusan work owner]` — struktural |
| `ADJUSTMENT_ID` | teks | **tidak** | FK → `T_CLAIM_ADJUSTMENT.ID` | `[keputusan work owner]` — struktural |
| `SHARE_PERCENTAGE` | angka desimal | ya | persen alokasi | korpus — `Activity/CountGrossAdjTreaty_Act.xml` langkah **2.1** *(Pega `SharePercentage`)* |
| `CLAIM_SPREADED` | angka desimal | ya | nilai tersebar | korpus — **CountGrossAdjTreaty_Act** langkah **5.2** *(Pega `ClaimSpreaded`)* |
| `CLAIM_ESTIMATION` | angka desimal | ya | estimasi klaim baris ini | korpus — **CountGrossAdjTreaty_Act** langkah **5.2** *(Pega `ClaimEstimation`)* |

⚠️ **Tabel paling tipis — hanya tiga kolom bermuatan.** Mata uang, kurs, dan penunjuk treaty
**tidak ketemu** di medan yang disisir. ⚠️ Halaman asalnya `LossAllocation` juga **diisi dengan
menyalin seluruh halaman** (`SetNameCurrency_Act` langkah **11** menyalin `TempSpreadingRisk`
ke `.LossAllocation`), jadi kolom sebenarnya mungkin lebih banyak. `[terbuka]`.

**Kunci tamu:** `ADJUSTMENT_ID` → `T_CLAIM_ADJUSTMENT.ID` · **ON DELETE CASCADE** · index
**biasa**.

---

## TUGAS 3 — Yang bukan tabel, dan tiga jebakan migrasi

### J1 · Halaman Pega yang **BUKAN** tabel

| # | Halaman Pega | Sebabnya | Penggantinya di Go |
| --- | --- | --- | --- |
| 1 | `ClaimData.SpreadingRisk` | **bacaan master treaty**, bukan data klaim | dibaca dari master treaty saat dibutuhkan |
| 2 | `PaymentData` + `AcceptationList` + `CurrencyList` + `DetailPayment` | milik **sistem kasir**, bukan milik klaim | dibaca dari layanan Kasir |
| 3 | `ClaimData.SpreadingAdjustment` | ⚠️ **halaman kerja sementara** tingkat klaim — salinan, bukan sumber | dihitung dari `T_CLAIM_ADJ_SPREADING` |
| 4 | `ClaimData.SpreadingAdjustmentQS` | idem, sisi quota-share | dihitung dari `T_CLAIM_ADJ_QUOTA_SHARE` |
| 5 | `AdjustmentList(n).FacRetroList` | **salinan** `ClaimData.FacRetroList` ke tingkat penyesuaian | dibaca dari `T_CLAIM_FAC_RETRO` |
| 6 | `InterestListDtl` | **salinan** `InterestList` | dibaca dari `T_CLAIM_INTEREST` |
| 7 | `ListTotalEstimation` + `TotalInterestInsured` | **hasil penjumlahan**, bukan data | dihitung saat dibaca |
| 8 | `ObjectList` + `ObjectItemList` | halaman kerja penyusun objek | dihitung dari `T_CLAIM_INTEREST` |

⛔ **Delapan baris, sepuluh halaman** — sesuai daftar diagram work owner. Diperiksa ulang terhadap
`spec.md`: **cocok**.

> ### ⭐ KANDIDAT KESEBELAS — **TEMUAN**
>
> `[terverifikasi]` **`.CashLossList`** diulang sebagai halaman langkah (`pyStepsObjectName`) di
> **tiga** activity: `Activity/CountEstimation_Act.xml`, `Activity/DeleteEstimation_Act.xml`, dan
> `Activity/SaveAcceptationTreaty_Act.xml`. ⛔ Ia **tidak ada** di daftar sepuluh tabel, **tidak
> ada** di daftar bukan-tabel, dan **tidak pernah disebut** di `spec.md`.
>
> ⚠️ **Nol properti ditulis padanya** di medan yang disisir (`PropertiesName`, `PropertiesValue`,
> `pyStepsObjectName`, `pyValue`) — sesuai aturan **I2** kalimatnya dibatasi pada medan itu saja.
>
> `[terbuka]` **Apakah `.CashLossList` tabel kesebelas, halaman kerja, atau sisa yang mati —
> belum dapat diputuskan.** ⛔ Tidak ditebak.

### J2 · Tiga jebakan migrasi dari `DATA_JSON` warisan

> ## ⛔ PERINGATAN MIGRASI — ditulis apa adanya
>
> | Halaman | Perintah | Sebabnya |
> | --- | --- | --- |
> | **`InterestListDtl`** | ⛔ **JANGAN diimpor** | ia **salinan** `InterestList` — mengimpornya menghasilkan **baris GANDA** |
> | **`SpreadingRisk`** | ⛔ **JANGAN diimpor** | ia **bacaan master treaty**, bukan data klaim |
> | **`PaymentData`** | ⛔ **JANGAN diimpor** | ia **milik sistem kasir** |

**Dan tiga aturan parser JSON warisan:**

1. ⚠️ **Buang seluruh kunci berawalan `px`, `py`, `pz`** — termasuk **`pyExpanded`** yang menyimpan
   status tampilan UI. **AC 10** dan **AC 11**. ⚠️ `pyExpanded` terbukti ada pada baris penyesuaian
   — `Activity/AddAdjustment_Act.xml` langkah **5**.
2. ⚠️ **Seluruh nilai di JSON lama bertipe teks** — angka, tanggal, dan penanda semuanya string.
   Konversi wajib eksplisit.
3. ⚠️ `[data DBA]` **Tanggal memakai DUA format berdampingan dalam satu dokumen**. **AC 12** —
   pembaca wajib menerima keduanya.

⚠️ **Parser tidak boleh mengandalkan sebuah kunci selalu hadir** — **AC 10**.

### J3 · Enam properti spreading = **enam peran, bukan satu**

> ## ⛔ CLASS YANG SAMA BUKAN BUKTI PERAN YANG SAMA
>
> **AC 3.** Di Pega keenam properti di bawah **berbagi satu definisi class**. ⛔ Itu **tidak punya
> akibat perilaku apa pun**, dan **bukan** alasan menyatukannya menjadi satu tabel.

| # | Properti Pega | Tabel tujuannya |
| --- | --- | --- |
| 1 | `ClaimData.SpreadingClaim` | **`T_CLAIM_SPREADING`** |
| 2 | `ClaimData.SpreadingBreakQS` | **`T_CLAIM_BREAK_QS`** |
| 3 | `AdjustmentList(n).SpreadingAdjustment` | **`T_CLAIM_ADJ_SPREADING`** |
| 4 | `AdjustmentList(n).SpreadingQuotaShare` | **`T_CLAIM_ADJ_QUOTA_SHARE`** |
| 5 | `ClaimData.SpreadingAdjustment` | ⛔ **BUKAN TABEL** — halaman kerja *(J1 no. 3)* |
| 6 | `ClaimData.SpreadingRisk` | ⛔ **BUKAN TABEL** — bacaan master treaty *(J1 no. 1)* |

⭐ **Empat menjadi tabel berbeda, dua tidak menjadi tabel sama sekali.** Menyatukannya karena
class-nya sama akan menggabungkan **spreading klaim**, **quota-share**, **spreading penyesuaian**,
dan **bacaan master treaty** ke dalam satu tabel yang tidak punya arti.

---

## TUGAS 4 — Pohon relasi dan penutup

### P1 · Pohon relasi — **lima tingkat**, berakar di `T_WORK_CLAIM`

```
T_WORK_CLAIM  (lintas-lini, satu baris per work object)
│   baris klaim   ID = CLMP-xxxxxx   COVER_KEY kosong    LINI = PROP
│   baris komite  ID = TKMT-xxxxxx   COVER_KEY = CLMP-…  LINI = PROP
│
├── T_GENERAL_CLAIM              ID = CLMP-xxxxxx   (shared PK, tanpa kolom penyambung)
│
├── T_CLAIM_ESTIMATION          CLAIM_ID → T_WORK_CLAIM.ID        1:N   CASCADE
├── T_CLAIM_INTEREST            CLAIM_ID → T_WORK_CLAIM.ID        1:N   CASCADE
├── T_CLAIM_CLAIM_AMOUNT        CLAIM_ID → T_WORK_CLAIM.ID        1:N   CASCADE
├── T_CLAIM_SPREADING           CLAIM_ID → T_WORK_CLAIM.ID        1:N   CASCADE   + TREATY_ID
├── T_CLAIM_BREAK_QS            CLAIM_ID → T_WORK_CLAIM.ID        1:N   CASCADE   + TREATY_ID
│                                  ⚠️ SEJAJAR dengan _SPREADING, bukan anaknya
├── T_CLAIM_FAC_RETRO           CLAIM_ID → T_WORK_CLAIM.ID        1:N   CASCADE
│
└── T_CLAIM_ADJUSTMENT          CLAIM_ID → T_WORK_CLAIM.ID        1:N   CASCADE
    │   KOMITE_ID → T_WORK_CLAIM.ID baris komite   1:1  UNIK, nullable, di Go
    │   + 7 kolom KOMITE_*  (dari DataCommitteeTreaty)
    │
    ├── T_CLAIM_ADJ_SPREADING       ADJUSTMENT_ID → …ADJUSTMENT.ID  1:N  CASCADE  + TREATY_ID
    ├── T_CLAIM_ADJ_QUOTA_SHARE     ADJUSTMENT_ID → …ADJUSTMENT.ID  1:N  CASCADE  + TREATY_ID
    │                                  ⚠️ SEJAJAR dengan _ADJ_SPREADING
    └── T_CLAIM_ADJ_LOSS_ALLOCATION  ADJUSTMENT_ID → …ADJUSTMENT.ID  1:N  CASCADE

    ┌─────────────── SISI KOMITE — ⛔ TIDAK didefinisikan ulang di sini ───────────────┐
    │  T_GENERAL_KOMITE.ADJUSTMENT_ID  ──1:1──►  T_CLAIM_ADJUSTMENT.ID                │
    │       NOT NULL · index UNIK · tanpa REFERENCES · ditegakkan di Go                │
    │       (definisinya terkunci di STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md)              │
    └──────────────────────────────────────────────────────────────────────────────────┘

    DIBACA dari luar, BUKAN anak, TIDAK ikut cascade:
       master treaty (SpreadingRisk)  ·  sistem kasir (PaymentData)  ·  DOCUMENT_CLAIM
```

**Tingkatnya lima:** `T_WORK_CLAIM` → `T_CLAIM_ADJUSTMENT` → `T_CLAIM_ADJ_SPREADING` →
*(baris di dalamnya)* → nilai. Jalur terdalam melewati **tiga tabel**.

⛔ **Tabel komite TIDAK didefinisikan ulang** — hanya digambar sambungannya.

### P2 · Butir `[terbuka]` berkas ini — ⭐ **SEBELAS**

⛔ **Tidak satu pun dinyatakan tertutup.**

| # | Butir | Menunggu |
| --- | --- | --- |
| **1** | ⭐ **Nama tabel jejak audit belum ditetapkan.** Usulan `T_CLAIM_AUDIT_TRAIL` — *penamaan, bukan rancangan* | **work owner** |
| **2** | **Rumus `local.hasilakhir`** yang mengisi `TOTAL_ESTIMATION_REINS` lewat `PrintDLATreatyIn` langkah **15.13** — **belum ditelusuri** | korpus |
| **3** | **Rule `RDB-List` mana yang mengisi `SpreadingRisk`** — sebelas activity menyebut keduanya, pengisinya belum diisolasi | korpus |
| **4** | **Induk `DOCUMENT_CLAIM` untuk lini PROP**, dan kolomnya `[data DBA]` | **DBA** |
| **5** | **Apakah `KOMITE_ID` dan `COVER_KEY` dipasangi `REFERENCES`** | **DBA** |
| **6** | ⚠️ **Keberatan nama `T_CLAIM_BREAK_QS`** — gaya nama tidak sejajar dengan pasangannya | **work owner** |
| **7** | ⭐ **`.CashLossList`** — tabel kesebelas, halaman kerja, atau sisa mati? | **work owner** |
| **8** | **Kurs tidak punya kolom** di `T_CLAIM_SPREADING`, `T_CLAIM_ADJ_SPREADING`, `T_CLAIM_ADJ_QUOTA_SHARE`, `T_CLAIM_ADJ_LOSS_ALLOCATION` — padahal **AC 25** menuntut *(nilai, mata uang, kurs)* per baris | **work owner** |
| **9** | **`T_CLAIM_ADJ_LOSS_ALLOCATION` hanya tiga kolom bermuatan** — halaman asalnya diisi dengan menyalin seluruh halaman, jadi kolomnya mungkin lebih banyak | korpus |
| **10** | **`T_VIEW_SUGGEST` tidak ketemu** di kedua modul — apa perannya bagi lini PROP? | **work owner** |
| **11** | **`FacRetroList` muncul di dua tempat** — tingkat klaim *(tabel 6)* dan tingkat penyesuaian *(bukan tabel)*; apakah keduanya memang data yang sama | **work owner** |

⭐ **Kelima butir yang brief sebut wajib muncul, muncul semua** — butir **1 · 2 · 3 · 4 · 5**.

### P3 · Sebelum relasi tabel Claim Prop bisa ditulis

| # | Yang kurang | Kenapa menahan |
| --- | --- | --- |
| **1** | **Jawaban butir 8** — kurs pada empat tabel spreading | menentukan bentuk kolom, bukan hanya relasi; **AC 25** menggantung padanya |
| **2** | **Nama tabel jejak audit** *(butir 1)* | relasinya ke `T_WORK_CLAIM` tidak bisa ditulis tanpa nama |
| **3** | **Vonis `.CashLossList`** *(butir 7)* | kalau ia tabel, pohon relasinya bertambah satu cabang |
| **4** | **`DOCUMENT_CLAIM` untuk lini PROP** *(butir 4)* | ia tabel bersama; relasinya ke klaim belum terbaca |
| **5** | **Keputusan `REFERENCES`** *(butir 5)* | menentukan apakah `ON DELETE` ditegakkan basis data atau Go |
| **6** | **Isi `T_CLAIM_ADJ_LOSS_ALLOCATION`** *(butir 9)* | tabel dengan tiga kolom sulit dipercaya lengkap |

### P4 · Pertanyaan BARU untuk work owner — **2**

#### Pertanyaan 1 — Empat tabel spreading tanpa kolom kurs: ditambahkan atau tidak?

**Apa yang ditanyakan.** **AC 25** menetapkan bentuk uang adalah **(nilai, mata uang, kurs) per
baris**. Empat tabel — `T_CLAIM_SPREADING`, `T_CLAIM_ADJ_SPREADING`, `T_CLAIM_ADJ_QUOTA_SHARE`,
`T_CLAIM_ADJ_LOSS_ALLOCATION` — memuat uang tetapi **kolom kursnya tidak ketemu di korpus**. Dua
di antaranya bahkan tidak punya kolom mata uang.

**Kenapa muncul.** Ketiadaannya baru kelihatan saat kolom disusun berdampingan. Di Pega kursnya
diambil dari halaman induk saat dibutuhkan, jadi tidak pernah perlu disimpan per baris.

**Bedanya kalau A atau B.** **A — ditambahkan:** keempat tabel memenuhi AC 25 sepenuhnya, dan nilai
historis tetap dapat direkonstruksi walau kurs induk berubah; harganya empat kolom baru **tanpa
sumber korpus**. **B — tidak ditambahkan:** mengikuti Pega apa adanya, kurs dibaca dari induk;
tetapi **AC 25 tidak terpenuhi** untuk keempat tabel itu, dan nilai historis ikut berubah bila
kurs induk dikoreksi.

**Apa yang tertahan.** Bentuk akhir keempat tabel, dan apakah **AC 25** perlu diberi pengecualian
tertulis.

#### Pertanyaan 2 — `.CashLossList`: tabel, halaman kerja, atau sisa yang mati?

**Apa yang ditanyakan.** `.CashLossList` diulang sebagai halaman langkah di tiga activity —
`CountEstimation_Act`, `DeleteEstimation_Act`, `SaveAcceptationTreaty_Act` — tetapi **nol properti
ditulis padanya** di medan yang disisir, dan ia **tidak pernah disebut** di `spec.md`, di daftar
sepuluh tabel, maupun di daftar bukan-tabel.

**Kenapa muncul.** Ia ketemu saat memeriksa ulang daftar bukan-tabel, persis seperti yang diminta.

**Bedanya kalau A atau B.** **A — ia tabel kesebelas:** diagram bertambah satu tabel, pohon relasi
bertambah satu cabang, dan kolomnya perlu disensus dari awal. **B — halaman kerja atau sisa mati:**
masuk daftar bukan-tabel sebagai butir kesebelas, dan migrasi **tidak boleh** mengimpornya.

**Apa yang tertahan.** Jumlah tabel `T_CLAIM_*` — sepuluh atau sebelas — dan isi daftar
bukan-tabel.

### P5 · ⭐ Yang seharusnya dikerjakan tetapi TIDAK diperintahkan blok ini

⛔ **Disebutkan, tidak dikerjakan.**

| # | Butir |
| --- | --- |
| **1** | **Menulis `RELASI-TABEL-CLAIM-PROP.md`** — berkas relasi sisi klaim, pasangan berkas relasi sisi komite yang sudah ada |
| **2** | **Menelusuri rumus `local.hasilakhir`** di `PrintDLATreatyIn` — ia mengisi nilai uang dan rumusnya tidak terbaca |
| **3** | **Menyensus kolom `.CashLossList`** sebelum memutuskan vonisnya |
| **4** | **Memeriksa apakah `T_CLAIM_ADJ_LOSS_ALLOCATION` mewarisi kolom `SpreadingRisk`** — `SetNameCurrency_Act` langkah **11** menyalin seluruh halaman ke `.LossAllocation` |
| **5** | **Menyelaraskan `STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md`** dengan nama kolom baru `CLOSE_FILE`/`RESERVED_CLAIM` — berkas itu masih menulis usulan lama |
| **6** | **Menetapkan generator nomor** untuk awalan `CLMP-` dan `TKMT-` pada `T_WORK_CLAIM.ID` |
| **7** | **Memeriksa apakah lini lain memakai `T_VIEW_SUGGEST`** — ia terdaftar sebagai tabel bersama tetapi tidak ketemu di dua modul |

---

## Lampiran pengikat penjaga — kolom DDL tabel milik Claim Prop (07-10-2026)

Bab di bawah **mengikat**: `inti/backend/penjaga` (`TestKolomDDLCocokDenganStruktur`,
`TestGolonganTipeDDLCocokDenganStruktur`) membandingkan setiap judul `## T_…` dengan DDL migrasi
`521`–`531`. TUGAS 2 di atas tetap catatan rancangan beserta alasan per kolom; bila keduanya berbeda,
bab di bawah inilah yang berlaku. Kolom Claim Prop di tabel bersama `T_GENERAL_CLAIM` (migrasi `520`)
dicatat di dokumen pemiliknya, bukan di sini, supaya satu tabel digambarkan satu dokumen: bab
`## T_GENERAL_CLAIM` di `modul/claimlife/docs/STRUKTUR-TABEL-CLAIM-LIFE.md`.

`T_VIEW_SUGGEST.CLAIM_ID` (migrasi `532`) dicatat di bab `## T_VIEW_SUGGEST`
`modul/premiumlistlife/docs/STRUKTUR-TABEL-PREMIUMLIST-LIFE.md` — keputusan work owner 07-10-2026 "1 tabel aja gabung life dan non life".
(Sempat dicabut pada hari yang sama karena uji PremiumList Life mengunci 243 kolom; uji itu kini menghitung 244.)

Uang, persen, share, dan kurs `NUMBER(38,10)` (keputusan work owner 07-10-2026).

## T_CLAIM_ESTIMATION

Migrasi `521`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `521` |
| `CLAIM_ID` | teks | tidak | FK | migrasi `521` — → `T_GENERAL_CLAIM.ID`, ON DELETE CASCADE |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `521` |
| `ESTIMATION_DATE` | DATE | ya |  | migrasi `521` |
| `ESTIMATION_TYPE` | teks | ya |  | migrasi `521` |
| `TREATY_TYPE_ID` | teks | ya |  | migrasi `521` |
| `TREATY_TYPE_NAME` | teks | ya |  | migrasi `521` |
| `CURRENCY_ID` | teks | ya |  | migrasi `521` |
| `CURRENCY_NAME` | teks | ya |  | migrasi `521` |
| `KURS` | angka desimal | ya |  | migrasi `521` |
| `GROSS_ESTIMATION_VALUE` | angka desimal | ya |  | migrasi `521` |
| `GROSS_ESTIMATION_IDR` | angka desimal | ya |  | migrasi `521` |
| `PERSEN_RNM` | angka desimal | ya |  | migrasi `521` |
| `ESTIMATION_VALUE` | angka desimal | ya |  | migrasi `521` |
| `ESTIMATION_VALUE_IDR` | angka desimal | ya |  | migrasi `521` |
| `NO_PLA` | teks | ya |  | migrasi `521` |
| `IS_PRINT_FACE_CLAIM` | teks | ya |  | migrasi `521` |

**Unik:** `(CLAIM_ID, NOURUT)`.

## T_CLAIM_INTEREST

Migrasi `522`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `522` |
| `CLAIM_ID` | teks | tidak | FK | migrasi `522` — → `T_GENERAL_CLAIM.ID`, ON DELETE CASCADE |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `522` |
| `OBJECT_NAME` | teks | ya |  | migrasi `522` |
| `CURRENCY_ID` | teks | ya |  | migrasi `522` |
| `CURRENCY_NAME` | teks | ya |  | migrasi `522` |
| `KURS` | angka desimal | ya |  | migrasi `522` |
| `TSI_VALUE` | angka desimal | ya |  | migrasi `522` |
| `TSI_VALUE_IDR` | angka desimal | ya |  | migrasi `522` |
| `IS_ADJ_VALUE` | teks | ya |  | migrasi `522` |

**Unik:** `(CLAIM_ID, NOURUT)`.

## T_CLAIM_CLAIM_AMOUNT

Migrasi `523`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `523` |
| `CLAIM_ID` | teks | tidak | FK | migrasi `523` — → `T_GENERAL_CLAIM.ID`, ON DELETE CASCADE |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `523` |
| `CURRENCY_ID` | teks | ya |  | migrasi `523` |
| `CURRENCY_NAME` | teks | ya |  | migrasi `523` |
| `KURS` | angka desimal | ya |  | migrasi `523` |
| `CLAIM_AMOUNT` | angka desimal | ya |  | migrasi `523` |
| `NET_DEDUCTIBLE_VALUE` | angka desimal | ya |  | migrasi `523` |
| `VALUE` | angka desimal | ya |  | migrasi `523` |
| `VALUE_IDR` | angka desimal | ya |  | migrasi `523` |
| `NOTE` | teks | ya |  | migrasi `523` |

**Unik:** `(CLAIM_ID, NOURUT)`.

## T_CLAIM_LOSS_ALLOCATION

Migrasi `524`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `524` |
| `CLAIM_ID` | teks | tidak | FK | migrasi `524` — → `T_GENERAL_CLAIM.ID`, ON DELETE CASCADE |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `524` |
| `CURRENCY_ID` | teks | ya |  | migrasi `524` |
| `CURRENCY_NAME` | teks | ya |  | migrasi `524` |
| `TREATY_TYPE_ID` | teks | ya |  | migrasi `524` |
| `TREATY_NAME` | teks | ya |  | migrasi `524` |
| `SHARE_PERCENTAGE` | angka desimal | ya |  | migrasi `524` |
| `CLAIM_SPREADED` | angka desimal | ya |  | migrasi `524` |
| `CLAIM_ESTIMATION` | angka desimal | ya |  | migrasi `524` |
| `KURS` | angka desimal | ya |  | migrasi `524` |
| `IS_OLD_DATA` | teks | ya |  | migrasi `524` |

**Unik:** `(CLAIM_ID, NOURUT)`.

## T_CLAIM_SPREADING

Migrasi `525`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `525` |
| `CLAIM_ID` | teks | tidak | FK | migrasi `525` — → `T_GENERAL_CLAIM.ID`, ON DELETE CASCADE |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `525` |
| `TREATY_ID` | teks | ya |  | migrasi `525` |
| `TREATY_NAME` | teks | ya |  | migrasi `525` |
| `SHARE_PERCENTAGE` | angka desimal | ya |  | migrasi `525` |
| `CLAIM_SPREADED` | angka desimal | ya |  | migrasi `525` |
| `CURRENCY_ID` | teks | ya |  | migrasi `525` |
| `CURRENCY_NAME` | teks | ya |  | migrasi `525` |
| `IS_OLD_DATA` | teks | ya |  | migrasi `525` |

**Unik:** `(CLAIM_ID, NOURUT)`.

## T_CLAIM_BREAK_QS

Migrasi `526`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `526` |
| `CLAIM_ID` | teks | tidak | FK | migrasi `526` — → `T_GENERAL_CLAIM.ID`, ON DELETE CASCADE |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `526` |
| `TREATY_ID` | teks | ya |  | migrasi `526` |
| `TREATY_NAME` | teks | ya |  | migrasi `526` |
| `SHARE_PERCENTAGE` | angka desimal | ya |  | migrasi `526` |
| `CLAIM_SPREADED` | angka desimal | ya |  | migrasi `526` |
| `CURRENCY_ID` | teks | ya |  | migrasi `526` |
| `CURRENCY_NAME` | teks | ya |  | migrasi `526` |
| `IS_OLD_DATA` | teks | ya |  | migrasi `526` |

**Unik:** `(CLAIM_ID, NOURUT)`.

## T_CLAIM_FAC_RETRO

Migrasi `527`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `527` |
| `CLAIM_ID` | teks | tidak | FK | migrasi `527` — → `T_GENERAL_CLAIM.ID`, ON DELETE CASCADE |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `527` |
| `REINSURER_ID` | teks | ya |  | migrasi `527` |
| `REINSURER_NAME` | teks | ya |  | migrasi `527` |
| `SHARE_PCT` | angka desimal | ya |  | migrasi `527` |
| `RI_COMMISSION_PCT` | angka desimal | ya |  | migrasi `527` |
| `ADDITIONAL_INFO` | teks | ya |  | migrasi `527` |
| `TOTAL_ESTIMATION_REINS` | angka desimal | ya |  | migrasi `527` |

**Unik:** `(CLAIM_ID, NOURUT)`.

## T_CLAIM_ADJUSTMENT

Migrasi `528`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `528` |
| `CLAIM_ID` | teks | tidak | FK | migrasi `528` — → `T_GENERAL_CLAIM.ID`, ON DELETE CASCADE |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `528` |
| `KOMITE_ID` | teks | ya |  | migrasi `528` |
| `ADJUSTMENT_TYPE` | teks | ya |  | migrasi `528` |
| `FORM_TYPE` | teks | ya |  | migrasi `528` |
| `PAYMENT_TYPE` | teks | ya |  | migrasi `528` |
| `CURRENCY_ID` | teks | ya |  | migrasi `528` |
| `CURRENCY_NAME` | teks | ya |  | migrasi `528` |
| `KURS` | angka desimal | ya |  | migrasi `528` |
| `PERSEN_RNM` | angka desimal | ya |  | migrasi `528` |
| `LOSS_ALLOCATION_NAME` | teks | ya |  | migrasi `528` |
| `LOSS_ALLOCATION_SHARE` | angka desimal | ya |  | migrasi `528` |
| `GROSS_ADJUSTMENT` | angka desimal | ya |  | migrasi `528` |
| `GROSS_ADJUSTMENT_IDR` | angka desimal | ya |  | migrasi `528` |
| `GROSS_VALUE` | angka desimal | ya |  | migrasi `528` |
| `INDIVIDUAL_RISK_TYPE` | teks | ya |  | migrasi `528` |
| `INDIVIDUAL_RISK_PCT` | angka desimal | ya |  | migrasi `528` |
| `INDIVIDUAL_RISK_VALUE` | angka desimal | ya |  | migrasi `528` |
| `INDIVIDUAL_RISK_RNM` | angka desimal | ya |  | migrasi `528` |
| `PROPOSE_ADJUSTMENT_VALUE` | angka desimal | ya |  | migrasi `528` |
| `ADJUSTMENT_VALUE` | angka desimal | ya |  | migrasi `528` |
| `VALUE_ADJUSTMENT_IDR` | angka desimal | ya |  | migrasi `528` |
| `TOTAL_ESTIMATION_VALUE` | angka desimal | ya |  | migrasi `528` |
| `ADJUSTER_FEE_VALUE` | angka desimal | ya |  | migrasi `528` |
| `SALVAGE_VALUE` | angka desimal | ya |  | migrasi `528` |
| `PAYABLE` | teks | ya |  | migrasi `528` |
| `PAYABLE_TO` | teks | ya |  | migrasi `528` |
| `BANK_NAME` | teks | ya |  | migrasi `528` |
| `BANK_BRANCH` | teks | ya |  | migrasi `528` |
| `BANK_ACCOUNT_NO` | teks | ya |  | migrasi `528` |
| `SWIFT_CODE` | teks | ya |  | migrasi `528` |
| `BANK_ID` | teks | ya |  | migrasi `528` |
| `DLA_NO_CEDING` | teks | ya |  | migrasi `528` |
| `DLA_NO_SOB` | teks | ya |  | migrasi `528` |
| `IS_DIRECT_TO_KASIR` | teks | ya |  | migrasi `528` |
| `STATUS_KASIR` | teks | ya |  | migrasi `528` |
| `ACCEPTED_NO` | teks | ya |  | migrasi `528` |
| `ACCEPTED_DATE` | DATE | ya |  | migrasi `528` |
| `ACCEPTANCE_STATUS` | teks | ya |  | migrasi `528` |
| `IS_APPROVED` | teks | ya |  | migrasi `528` |
| `IS_KOMITE` | teks | ya |  | migrasi `528` |
| `IS_SUBJECTIVITY` | teks | ya |  | migrasi `528` |
| `SUBJECTIVITY_NOTE` | teks | ya |  | migrasi `528` |
| `NOTES` | teks | ya |  | migrasi `528` |
| `DLA_NO` | teks | ya |  | migrasi `528` |
| `REMARKS_DLA` | teks | ya |  | migrasi `528` |
| `IS_FAC_RETRO` | teks | ya |  | migrasi `528` |
| `IS_PRINT_ACCEPT` | teks | ya |  | migrasi `528` |
| `KOMITE_CIRCUM_CAUSE_OF_LOSS` | teks | ya |  | migrasi `528` |
| `KOMITE_ADJUSTER_FEE` | teks | ya |  | migrasi `528` |
| `KOMITE_REMARKS` | teks | ya |  | migrasi `528` |
| `KOMITE_SALVAGE` | teks | ya |  | migrasi `528` |
| `KOMITE_LEGAL_LIABILITY` | teks | ya |  | migrasi `528` |
| `KOMITE_EXTENT_OF_LOSS` | teks | ya |  | migrasi `528` |
| `KOMITE_OCCUPATION` | teks | ya |  | migrasi `528` |
| `CREATED_BY` | teks | ya |  | migrasi `528` |
| `CREATED_BY_NAME` | teks | ya |  | migrasi `528` |
| `CREATED_AT` | DATE | ya |  | migrasi `528` |

**Unik:** `(CLAIM_ID, NOURUT)`; `(KOMITE_ID)`.

## T_CLAIM_ADJ_SPREADING

Migrasi `529`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `529` |
| `ADJUSTMENT_ID` | teks | tidak | FK | migrasi `529` — → `T_CLAIM_ADJUSTMENT.ID`, ON DELETE CASCADE |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `529` |
| `TREATY_ID` | teks | ya |  | migrasi `529` |
| `TREATY_NAME` | teks | ya |  | migrasi `529` |
| `SHARE_PERCENTAGE` | angka desimal | ya |  | migrasi `529` |
| `CLAIM_SPREADED` | angka desimal | ya |  | migrasi `529` |
| `CURRENCY_ID` | teks | ya |  | migrasi `529` |
| `CURRENCY_NAME` | teks | ya |  | migrasi `529` |
| `PREMIUM_SPREADED` | angka desimal | ya |  | migrasi `529` |
| `BANK_ACCOUNT_NO` | teks | ya |  | migrasi `529` |
| `BANK_ID` | teks | ya |  | migrasi `529` |

**Unik:** `(ADJUSTMENT_ID, NOURUT)`.

## T_CLAIM_ADJ_QUOTA_SHARE

Migrasi `530`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `530` |
| `ADJUSTMENT_ID` | teks | tidak | FK | migrasi `530` — → `T_CLAIM_ADJUSTMENT.ID`, ON DELETE CASCADE |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `530` |
| `TREATY_ID` | teks | ya |  | migrasi `530` |
| `TREATY_NAME` | teks | ya |  | migrasi `530` |
| `SHARE_PERCENTAGE` | angka desimal | ya |  | migrasi `530` |
| `CLAIM_SPREADED` | angka desimal | ya |  | migrasi `530` |
| `CURRENCY_ID` | teks | ya |  | migrasi `530` |
| `CURRENCY_NAME` | teks | ya |  | migrasi `530` |

**Unik:** `(ADJUSTMENT_ID, NOURUT)`.

## T_CLAIM_ADJ_LOSS_ALLOCATION

Migrasi `531`.

| Kolom | Tipe | Null | Kunci | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | migrasi `531` |
| `ADJUSTMENT_ID` | teks | tidak | FK | migrasi `531` — → `T_CLAIM_ADJUSTMENT.ID`, ON DELETE CASCADE |
| `NOURUT` | bilangan bulat | tidak |  | migrasi `531` |
| `CURRENCY_ID` | teks | ya |  | migrasi `531` |
| `CURRENCY_NAME` | teks | ya |  | migrasi `531` |
| `TREATY_TYPE_ID` | teks | ya |  | migrasi `531` |
| `TREATY_NAME` | teks | ya |  | migrasi `531` |
| `SHARE_PERCENTAGE` | angka desimal | ya |  | migrasi `531` |
| `CLAIM_SPREADED` | angka desimal | ya |  | migrasi `531` |
| `CLAIM_ESTIMATION` | angka desimal | ya |  | migrasi `531` |
| `KURS` | angka desimal | ya |  | migrasi `531` |

**Unik:** `(ADJUSTMENT_ID, NOURUT)`.
