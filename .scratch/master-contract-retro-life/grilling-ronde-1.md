# Grilling — Master Contract Retro Life — Ronde 1

Tanggal: 2026-09-15
Konteks: `master-contract-retro-life` — **konteks/menu sendiri** `[keputusan work owner]`
Modul: **Master Contract Retro Life** (66 berkas — **modul terkecil di korpus**)
Skill: `/mattpocock-skills:grilling`
Sumber: korpus `D:\XML\RNM_BRD\` (READ-ONLY), `discovery/modules|flows/Master Contract Retro Life.md`,
`discovery/context-map.md` §2.9, `CONTEXT.md`, `docs/adr/`

> **Konvensi.** `[terverifikasi]` = terbukti korpus dengan **class + nama + path**;
> `[keputusan work owner]`; `[data DBA]`; `[dugaan]`; `[terbuka]` = OQ.

> `[keputusan work owner]` **Master Product Name Life adalah konteks TERPISAH** — tidak digarap di
> sini. Pelajaran dari NB/EDM: menggabung menyembunyikan perbedaan.

---

## Bagian A — Tiga rumus yang belum pernah dibaca: **sudah saya buka**

Batas telusur D3 menyebut tiga activity `@BASECLASS` sebagai "belum dibaca". Ketiganya kini terbaca —
dan **dua di antaranya ternyata bukan rumus**.

### A1. `CountingPercentShare_Act` — **penjumlah share, bukan pembagi**

`[terverifikasi]` `@BASECLASS` / `COUNTINGPERCENTSHARE_ACT` / `RULE-OBJ-ACTIVITY`
(`Master Contract Retro Life/Activity/CountingPercentShare_Act.xml`, 59.429 byte, 4 langkah):

| Step | Isi |
| ---: | --- |
| 1 | `TempSearch.CARI1 = Param.TREATYYEARID`; `TempSearch.CARI2 = Param.TREATYCONTRACTID`; `InputSecurityLife.STDRATING = ""` |
| 2 | `RDB-List` → `GetMasterReinsurerLifeList_SQl` — "Count TotalShare" |
| 3 | `Local.TotalShare = 0`; lalu **`Local.TotalShare = Local.TotalShare + @toDecimal(.PCTSHARE)`** |
| 4 | **`InputSecurityLife.STDRATING = Local.TotalShare`** |

`[terverifikasi]` Kueri sumbernya (`ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` /
`ASM!GETMASTERREINSURERLIFELIST_SQL` / `RULE-CONNECT-SQL`):

```sql
SELECT PCTSHARE FROM POOLDATA.TREATYREINSURER_LIFE
WHERE TREATYYEARID = {TempSearch.CARI1} AND TREATYCONTRACTID = {TempSearch.CARI2}
```

**Ia hanya MENJUMLAHKAN** `PCTSHARE` seluruh reinsurer pada satu (tahun treaty, kontrak). **Tidak ada
pembagian share, tidak ada alokasi.** Nama "CountingPercentShare" akurat; dugaan D3 bahwa ini "rumus
pembagian share retro" **tidak terbukti**.

⚠️ **Dua kejanggalan:**

1. Hasil jumlah ditaruh di field bernama **`STDRATING`** ("standard rating") — field dipakai ulang
   untuk menyimpan **total share**. Nama tidak mencerminkan isi.
2. `[terverifikasi]` **Tidak ada satu pun perbandingan terhadap 100** di seluruh modul. Sensus:
   `TotalShare` hanya muncul pada tiga ekspresi di atas, dan `STDRATING` tidak pernah diuji. Total
   dihitung dan **ditampilkan**, tidak divalidasi. → **Q4**

### A2. `TreatyLimit_TypeProtect` — **lookup, bukan rumus**

`[terverifikasi]` `@BASECLASS` / `TREATYLIMIT_TYPEPROTECT` / `RULE-OBJ-ACTIVITY` (54.893 byte,
4 langkah):

| Step | Isi |
| ---: | --- |
| 1–2 | `TempInputData.CARI1 = InputTreatyContract.REINSTYPEID`; `Param.pyReportName = "BrowseReinsuranceTypeLimit_RD"`; `Param.pyReportClass = "ASM-FW-GISFW-Int-REINSURANCETYPE"`; `Param.pyPageName = "ReinsuranceType"` |
| 3 | `Call Rule-Obj-Report-Definition.pxRetrieveReportData` |
| **4** | *(loop)* precondition **`TempInputData.CARI1 == .ID`** → **`InputTreatyContract.REINSTYPENAME = .Note`** |

**Ia menerjemahkan `REINSTYPEID` menjadi nama**, dengan namanya diambil dari kolom **`.Note`**.
Tidak ada logika "jenis proteksi limit" — dugaan D3 **tidak terbukti**.

### ⚠️ A2.1 Ini **menyempitkan OQ-057** — sumber enumerasinya ditemukan

`[terverifikasi]` `Master Contract Retro Life/ReportDefinition/BrowseReinsuranceTypeLimit_RD.xml` →
`ASM-FW-GISFW-INT-REINSURANCETYPE` / `BROWSEREINSURANCETYPELIMIT_RD` / `RULE-OBJ-REPORT-DEFINITION`.

| Hal | Temuan |
| --- | --- |
| Class penyimpan enumerasi | **`ASM-FW-GISFW-INT-REINSURANCETYPE`** |
| Kolom tersedia | `.ID`, `.Code`, `.Note`, `.Type`, `.Flag`, `.SOANote`, `.UserID` |
| **Penyaring** | **`.Flag = 1`** (baris 581/585 dan salinan indeks 887/893) |
| Nama yang ditampilkan | **`.Note`** — bukan `.Code`, bukan `.Type` |

`[terverifikasi]` RD serupa ada di **sembilan modul lain** (`BrowseReinsuranceType_RD`, class yang
sama) — jadi enumerasi ini **dipakai lintas lini**, bukan khas life. RD versi modul lain difilter
lewat `Param.ID` / `Param.Name` / `Param.Note` / `Param.Type`, **bukan** `.Flag`.

**Yang masih belum terbaca:** **nilai-nilainya** — ia tabel master di basis data. → **Q7**

### A3. `SetValueRetroLimit_TreatyYearLife` — **penyalur konteks, bukan rumus**

`[terverifikasi]` `@BASECLASS` / `SETVALUERETROLIMIT_TREATYYEARLIFE` / `RULE-OBJ-ACTIVITY`
(63.678 byte, 2 langkah). Ia **menyalin konteks tahun treaty terpilih** ke tiga halaman input:

| Sasaran | Isi yang disalin |
| --- | --- |
| `InputTreatyContract` | `IDTREATYYEAR`, `TREATYSTARTDATE`, `TREATYENDDATE` |
| `InputBusinessLife` | `TREATYCONTRACTID`, `TREATYYEAR`, `TREATYYEARID`, `REINSTYPEID`, `REINSTYPENAME` |
| `InputSecurityReinsurer` | `TREATYYEARID`, `TREATYCONTRACTID`, `TREATYREINSURERID`, `REINSURERNAME`, `PCTSHARE` |

Ditambah `OutputParam.STSSAVE` dan `OutputParam.DATASHOW` — deskripsi langkahnya sendiri:
*"OutputParam.DATASHOW - hide inputan, OutputParam.STSSAVE - hide status save/delete"*, yakni
**kendali tampilan**.

**Tidak ada perhitungan limit apa pun.** Dugaan D3 **tidak terbukti**.

> **Kesimpulan Bagian A:** ketiga "rumus" itu **bukan rumus**. Satu penjumlah, satu lookup, satu
> penyalur konteks. **Perhitungan retro life yang sesungguhnya — bila ada — tidak berada di korpus.**
> Kandidat tunggalnya adalah kelima stored procedure penulis. → **Q9**

---

## Bagian B — Hierarki entitas: **terbaca penuh dari tanda tangan procedure**

`[terverifikasi]` Kelima Connect-SQL penulis beserta tanda tangannya:

| Rule | Class / Nama | Procedure + parameter |
| --- | --- | --- |
| `SaveMasterTreatyYear_Life_SQL` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` / `ASM!SAVEMASTERTREATYYEAR_LIFE_SQL` | `INSERTTREATYYEAR_LIFE(p_ID, p_TREATYYEAR, p_UNDERWRITINGYEAR, p_USERID, p_TGLUPDATE, p_STARTDATE, p_ENDDATE, HASIL1 out)` |
| `SaveMasterTreatyContract_Life_SQL` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` / `ASM!SAVEMASTERTREATYCONTRACT_LIFE_SQL` | `INSERTTREATYCONTRACT_LIFE(p_ID, **p_IDTREATYYEAR**, p_REINSTYPEID, p_REINSTYPENAME, p_TREATYSTARTDATE, p_TREATYENDDATE, p_USERID, p_TGLUPDATE, **p_IDR, p_USD, p_B_IDR, p_B_USD, p_IDR_SELISIH, p_USD_SELISIH**, HASIL1 out)` |
| `SaveMasterTreatyReinsurer_Life_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!SAVEMASTERTREATYREINSURER_LIFE_SQL` | `INSERTREINSURER_LIFE(p_ID, **p_TREATYYEARID, p_TREATYCONTRACTID**, p_REINSTYPEID, p_REINSTYPENAME, p_REINSURERID, p_REINSURERNAME, **p_PCTSHARE, p_COMMISION, p_OVR_COMM**, p_USERID, p_TGLUPDATE, HASIL1 out)` |
| `SaveMasterTreatySecurityReinsurer_Life_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!SAVEMASTERTREATYSECURITYREINSURER_LIFE_SQL` | `INSERTSECURITYREINSURER_LIFE(p_ID, p_TREATYYEARID, p_TREATYCONTRACTID, **p_TREATYREINSURERID**, p_REINSURERID, p_REINSURERNAME, **p_PCTSHARE**, p_USERID, p_TGLUPDATE, HASIL1 out)` |
| `SaveMasterTreatyBusiness_Life_SQL` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` / `ASM!SAVEMASTERTREATYBUSINESS_LIFE_SQL` | `INSERTBUSINESS_LIFE(p_ID, p_TREATYYEARID, p_TREATYYEAR, p_TREATYCONTRACTID, p_REINSTYPEID, p_REINSTYPENAME, **p_BIZCODE, p_BIZNAME, p_RIRATEID, p_RIRATE**, p_USERID, p_TGLUPDATE, HASIL1 out)` |

### Hierarki yang terbaca `[terverifikasi]`

```
TREATYYEAR_LIFE            (ID, TREATYYEAR, UNDERWRITINGYEAR, STARTDATE, ENDDATE)
  └─ TREATYCONTRACT_LIFE   (ID, IDTREATYYEAR → tahun, REINSTYPEID, tanggal, ENAM kolom limit)
       ├─ TREATYREINSURER_LIFE          (ID, TREATYYEARID + TREATYCONTRACTID, REINSURERID,
       │                                  PCTSHARE, COMMISION, OVR_COMM)
       │    └─ TREATYSECURITYREINSURER_LIFE (ID, …, TREATYREINSURERID → reinsurer, PCTSHARE)
       └─ TREATYBUSINESS_LIFE           (ID, TREATYYEARID + TREATYCONTRACTID, BIZCODE/BIZNAME,
                                          RIRATEID, RIRATE)
```

**Ya, berjenjang** — dan **lebih dalam dari dugaan D3**: ada **dua tingkat share**.
`TREATYSECURITYREINSURER_LIFE` menunjuk `TREATYREINSURERID`, yakni **share di bawah share**. → **Q3**

`[terverifikasi]` `REINSTYPEID`/`REINSTYPENAME` **diulang** di contract, reinsurer, dan business —
denormalisasi. Tahun (`TREATYYEAR`) juga diulang di business. → **Q1**

`[terverifikasi]` Setiap procedure menerima **`p_USERID` dan `p_TGLUPDATE`** (kolom audit) dan
mengembalikan **`HASIL1 out`**.

---

## Bagian C — CRUD: hapus **tanpa kaskade sama sekali**

`[terverifikasi]` Keempat penghapus, seluruhnya **`DELETE … WHERE ID = …` datar**:

| Rule | Class / Nama | SQL |
| --- | --- | --- |
| `DeleteTreatyLimit_SQL` | `ASM-FW-GISFW-INT-RETROCESSIONLIFE` / `ASM!DELETETREATYLIMIT_SQL` | `delete from POOLDATA.treatycontract_life where ID = {InputTreatyContract.ID}` |
| `DeleteSecurityReinsurer_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!DELETESECURITYREINSURER_SQL` | `DELETE FROM POOLDATA.TREATYREINSURER_LIFE WHERE ID = {InputSecurityLife.ID}` |
| `DeleteSecurityReinsurerLife_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!DELETESECURITYREINSURERLIFE_SQL` | `DELETE FROM POOLDATA.TREATYSECURITYREINSURER_LIFE WHERE ID = {InputSecurityReinsurerLife.ID}` |
| `DeleteRowBusinessList` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` / `ASM!DELETEROWBUSINESSLIST` | `delete from pooldata.treatybusiness_life where id={InputData.HASIL12}` |

⚠️ **Tiga akibat yang harus diputuskan:**

1. **Nol kaskade.** Menghapus kontrak **tidak** menghapus reinsurer, security reinsurer, maupun
   business di bawahnya. Menghapus reinsurer **tidak** menghapus security reinsurer yang menunjuknya
   lewat `TREATYREINSURERID`. **Risiko yatim nyata dan tidak ditangani di korpus.** → **Q5**
2. **`TREATYYEAR_LIFE` tidak punya penghapus sama sekali** — lima penulis, **empat** penghapus.
   Tahun treaty tidak dapat dihapus dari modul ini. → **Q6**
3. `DeleteTreatyLimit_SQL` ber-class **`ASM-FW-GISFW-INT-RETROCESSIONLIFE`** tetapi menghapus
   `treatycontract_life` — **class tidak sejalan dengan tabel sasarannya**. Dicatat sebagai
   kejanggalan penamaan, bukan cacat perilaku.

### C1. Validasi yang **ada** di korpus

`[terverifikasi]` `SetErrorMessageReinsurer` (`@BASECLASS` / `SETERRORMESSAGEREINSURER` /
`RULE-OBJ-ACTIVITY`) — **batas per baris**:

```
@toDecimal(InputTreatyReinsurer.PctShare) > 100 || @toDecimal(InputTreatyReinsurer.PctShare) < 0
@toDecimal(InputTreatyReinsurer.Ricomm)   > 100 || @toDecimal(InputTreatyReinsurer.Ricomm)   < 0
```

`[terverifikasi]` **Field wajib** saat simpan:

| Activity | Wajib diisi |
| --- | --- |
| `SaveTreatyYearLife_Act` | `UNDERWRITINGYEAR`, `TREATYYEAR`, `STARTDATE`, `ENDDATE` |
| `SaveTreatyLimit_Act` (kontrak) | `REINSTYPEID`, `TREATYSTARTDATE`, `TREATYENDDATE`, **`B_IDR`, `IDR`, `B_USD`** |

⚠️ **`USD` TIDAK termasuk wajib**, sementara `IDR`, `B_IDR`, dan `B_USD` wajib. Asimetri yang
mencurigakan. → **Q2**

`[terverifikasi]` **Gerbang konsistensi tahun** — dipakai bersama oleh `SaveBusinessLife_Act`,
`SaveSecurityLife_Act`, dan `SaveSecurityReinsurerLife_Act`:

```
@substring(InputRetrocessionLife.TREATYSTARTDATE,6,10) <> InputRetrocessionLifeTreatyType.TREATYYEAR_LIFE
```

yakni **tahun di dalam tanggal mulai treaty harus sama dengan tahun treaty** — dibaca dengan
`substring` posisi 6–10, sehingga **bergantung pada format tanggal**. → **Q8**

`[terverifikasi]` `SaveBusinessToAllLife_Act` bergerbang `.REINSTYPEID == Param.REINSTYPEID` —
fitur **terapkan ke semua** baris berjenis reasuransi sama. → **Q10**

### ⚠️ C2. Keluaran galat procedure **diabaikan**

`[terverifikasi]` Kelima procedure mengembalikan `{OutputData.HASIL1 out}`. Sensus 27 Activity
modul ini: **hanya `DeleteRowBusiness.xml`** yang menyebut `HASIL1`. **Tidak satu pun activity
`Save*` membacanya.**

**Artinya galat yang dilaporkan procedure jatuh diam-diam** — penyimpanan tampak berhasil meski
procedure menolak. Sekeluarga dengan "fallback diam" yang sudah kita tolak di PremiumList Life.
→ **Q11**

---

## Bagian D — Batas pengetahuan yang tetap berdiri

| Hal | Status |
| --- | --- |
| Body kelima procedure penulis | `[terbuka]` **OQ-002** — aturan simpan (validasi, versi, kunci unik, kaskade) ada di sisi basis data |
| Nilai `REINSTYPEID` | `[terbuka]` **OQ-057 / OQ-020** — **sumbernya kini diketahui** (§A2.1), nilainya belum |
| DDL keempat tabel `_LIFE` | `[terbuka]` **OQ-001** — tipe/presisi kolom uang dan share |
| `Claude outputs/` (`Struktur_GridRetrocessionLife.xlsx`, `perubahan_skill.diff`) | `[terbuka]` **OQ-054** — berkas non-Pega di korpus; **tidak dibaca**, dicatat saja |

`[terverifikasi]` Modul ini **tanpa `ConnectREST`, tanpa `When`, tanpa `SystemSettings`, tanpa
`IsPEGAPROD`, tanpa kolom JSON, tanpa guard identitas**. Tidak ada efek keluar, tidak ada flag
lingkungan, tidak ada ADR-0013. **Seluruh interaksi luar = 12 Connect-SQL ke `POOLDATA`.**

---

# Frontier Ronde 1 — 11 pertanyaan

**Empat fakta bisnis (OQ tim — TIDAK saya tebak): Q2, Q3, Q6, Q7.**
**Tujuh keputusan desain / konfirmasi (saya beri rekomendasi): Q1, Q4, Q5, Q8, Q9, Q10, Q11.**

❓ **Q1** — **Denormalisasi `REINSTYPEID` dan `TREATYYEAR`: dibawa atau dinormalkan?**
`[terverifikasi]` `REINSTYPEID`/`REINSTYPENAME` tersimpan **tiga kali** — di contract, reinsurer, dan
business. `TREATYYEAR` tersimpan di year **dan** di business. Di sistem baru: tiru apa adanya
(kolom berulang), atau normalkan (anak mewarisi dari induk lewat kunci asing)?

➡️ **Rekomendasi: normalkan, dengan satu pengecualian sadar.** Simpan `REINSTYPEID` **hanya di
contract**; reinsurer dan business mewarisinya lewat `TREATYCONTRACTID`. Alasannya: tiga salinan
dapat menyimpang, dan korpus **tidak punya** mekanisme yang menjaganya tetap sama. Pengecualian:
pertahankan **`REINSTYPENAME` hasil lookup** sebagai nilai tampilan yang di-cache bila laporan
lama bergantung padanya — tetapi tandai jelas sebagai turunan, bukan sumber kebenaran.

---

❓ **Q2** — **Enam kolom limit pada kontrak: `IDR`, `USD`, `B_IDR`, `B_USD`, `IDR_SELISIH`,
`USD_SELISIH` — apa artinya masing-masing?** *(fakta bisnis)* `[terverifikasi]` Keenamnya parameter
`INSERTTREATYCONTRACT_LIFE`. `[dugaan]` `B_` mungkin "batas"/"bruto", `_SELISIH` = selisih — **tetapi
korpus tidak menyatakannya dan saya tidak menebak.**

⚠️ Ada petunjuk yang memperdalam pertanyaan: **`IDR`, `B_IDR`, dan `B_USD` WAJIB diisi, sementara
`USD` TIDAK** (`SaveTreatyLimit_Act`). Apakah itu aturan sah, atau kelalaian validasi?

➡️ **Tidak ada rekomendasi — ini fakta bisnis.** Bila tidak terjawab saya buka OQ, dan spec hanya
mengikat bahwa keenam kolom ada dan bertipe uang menurut **ADR-0003**.

---

❓ **Q3** — **Dua tingkat share: apa hubungan `TREATYREINSURER_LIFE.PCTSHARE` dengan
`TREATYSECURITYREINSURER_LIFE.PCTSHARE`?** *(fakta bisnis)* `[terverifikasi]` Security reinsurer
menunjuk **`TREATYREINSURERID`** — ia **anak dari reinsurer**, bukan saudara. Apakah share tingkat
kedua itu: (a) porsi **dari** share reinsurer induknya (retrosesi atas retrosesi), (b) porsi dari
keseluruhan treaty, atau (c) sesuatu yang lain — misalnya jaminan/*security* atas share reinsurer?

⚠️ Ini menentukan apakah `10%` di tingkat dua berarti 10% dari treaty atau 10% dari share induknya —
**perbedaan besar pada angka**.

➡️ **Tidak ada rekomendasi — ini fakta bisnis.** Menebak di sini akan salah menghitung eksposur.

---

❓ **Q4** — **Total share: divalidasi harus 100, atau hanya ditampilkan?** `[terverifikasi]` Korpus
**hanya menjumlahkan** (`CountingPercentShare_Act`) lalu menaruh hasilnya di `STDRATING` untuk
ditampilkan. **Tidak ada** perbandingan terhadap 100 di seluruh modul; yang divalidasi hanyalah
batas per baris `0..100` (`SetErrorMessageReinsurer`).

➡️ **Rekomendasi: tampilkan total secara mencolok, tetapi JANGAN blokir penyimpanan — kecuali work
owner menyatakan 100% wajib.** Alasan: master kontrak lazim disusun bertahap, dan memaksa 100% di
tiap penyimpanan akan menghalangi pekerjaan setengah jadi. Yang saya sarankan **ditambahkan**:
penanda visual "belum 100%" dan **laporan kontrak yang sharenya ≠ 100** — sehingga celah terlihat
tanpa memblokir. Bila work owner menghendaki penegakan keras, itu **penyimpangan sadar** dan saya
catat demikian.

---

❓ **Q5** — **Hapus berkaskade atau tolak-bila-punya-anak?** `[terverifikasi]` Korpus **tidak punya
kaskade sama sekali** — empat `DELETE … WHERE ID = …` datar. Menghapus kontrak meninggalkan
reinsurer, security reinsurer, dan business **yatim**.

Tiga pilihan: (a) tiru apa adanya; (b) **kaskade** — hapus anak bersama induk; (c) **tolak** —
induk tidak dapat dihapus selama punya anak.

➡️ **Rekomendasi: (c) tolak-bila-punya-anak, dengan pesan yang menyebut berapa anak menghalangi.**
Alasan: ini **master data** yang dirujuk Claim Life dan Master Product Name Life — kaskade diam-diam
dapat menghapus riwayat yang masih dipakai konteks lain, dan meniru apa adanya membiarkan data yatim
yang sudah pasti keliru. ⚠️ **Fakta pendukung yang belum ada:** apakah basis data sudah punya
*foreign key* berkaskade — itu **OQ-001/OQ-002**, perlu **DBA**.

---

❓ **Q6** — **Mengapa tahun treaty tidak dapat dihapus?** *(fakta bisnis)* `[terverifikasi]` Lima
penulis, **empat** penghapus — `TREATYYEAR_LIFE` tidak punya jalur hapus di modul ini. Apakah itu
(a) disengaja (tahun treaty abadi setelah dibuat), atau (b) kelalaian yang selama ini ditambal lewat
basis data langsung?

➡️ **Tidak ada rekomendasi — ini fakta bisnis.** Bila (a), saya jadikan aturan eksplisit di spec;
bila (b), ia menjadi keputusan desain baru.

---

❓ **Q7** — **Daftar nilai `REINSTYPEID` dan arti `Flag = 1`.** *(fakta bisnis — OQ-057)*
`[terverifikasi]` **Sumbernya kini ditemukan**: class `ASM-FW-GISFW-INT-REINSURANCETYPE`, dibaca RD
`BROWSEREINSURANCETYPELIMIT_RD` dengan penyaring **`.Flag = 1`**, dan nama yang ditampilkan diambil
dari kolom **`.Note`**. Kolom tersedia: `.ID`, `.Code`, `.Note`, `.Type`, `.Flag`, `.SOANote`,
`.UserID`.

Tiga hal yang saya butuhkan: (a) **daftar nilai** `ID` beserta artinya; (b) arti **`Flag = 1`** —
apakah "aktif", atau "berlaku untuk limit"; (c) beda **`.Code`**, **`.Note`**, dan **`.Type`** —
mengapa nama diambil dari `Note`, bukan dari `Code`?

➡️ **Tidak ada rekomendasi — ini fakta bisnis.** Ini dapat dijawab **DBA** (isi tabel) **atau**
Product+UW (artinya). Menutup ini juga menutup **OQ-057** untuk sembilan modul lain yang memakai
class sama.

---

❓ **Q8** — **Gerbang konsistensi tahun lewat `substring` posisi 6–10: dipertahankan?**
`[terverifikasi]` Tiga activity simpan berbagi gerbang
`@substring(TREATYSTARTDATE,6,10) <> TREATYYEAR_LIFE` — membandingkan **potongan teks tanggal**
dengan tahun treaty. Ia benar hanya bila format tanggalnya tetap.

➡️ **Rekomendasi: pertahankan aturannya, buang caranya.** Aturan "tahun mulai treaty harus sama
dengan tahun treaty induknya" **sah dan berguna** — ia mencegah kontrak nyasar tahun. Tetapi di
sistem baru bandingkan **tahun dari tanggal bertipe tanggal**, bukan potongan teks. Ini
**perbaikan sadar**: `substring` diam-diam salah bila format tanggal berubah.

---

❓ **Q9** — **Kelima procedure penulis: dipanggil apa adanya, atau ditulis ulang di Go?**
`[terverifikasi]` **Seluruh jalur tulis** modul ini melewati stored procedure yang **bodinya tidak
ada di korpus** (OQ-002). Bagian A membuktikan tidak ada perhitungan di sisi Pega — jadi **apa pun
aturan bisnis penyimpanan retro life, ia ada di dalam procedure itu.**

➡️ **Rekomendasi: panggil apa adanya, seperti tiga konteks Life sebelumnya.** Menulis ulang menuntut
kita mengetahui isinya — dan kita tidak. Konsisten dengan **ADR-0006** (penomoran lewat procedure)
dan pola PremiumList Life. ⚠️ **Tetapi** ini menjadikan **OQ-002 pemblokir spec**, bukan sekadar
pemblokir tiket migrasi: tanpa body, saya tidak dapat menulis AC tentang *apa yang terjadi saat
simpan* — hanya *bahwa* procedure dipanggil. **Butuh DBA.**

---

❓ **Q10** — **Fitur "terapkan ke semua" (`SaveBusinessToAllLife_Act`): dibawa?**
`[terverifikasi]` Activity ini menyimpan business ke **seluruh** baris berjenis reasuransi sama
(`.REINSTYPEID == Param.REINSTYPEID`), lewat Connect-SQL tersendiri
`SaveTreatyBusinessAll_Life_SQL`. Ia operasi massal tanpa konfirmasi yang terbaca di korpus.

➡️ **Rekomendasi: bawa, tetapi dengan konfirmasi dan pratinjau.** Fitur ini jelas menghemat kerja
pada master besar. Yang saya tambahkan: layar **memberi tahu berapa baris akan terpengaruh** sebelum
dijalankan, dan hasilnya tercatat di jejak audit. Tanpa itu, satu klik dapat mengubah puluhan baris
tanpa jejak.

---

❓ **Q11** — **Keluaran galat procedure (`HASIL1`) yang diabaikan: ditegakkan?**
`[terverifikasi]` Kelima procedure mengembalikan `HASIL1 out`; **tidak satu pun activity `Save*`
membacanya** (hanya `DeleteRowBusiness` menyebut `HASIL1`). Galat procedure **jatuh diam-diam**.

➡️ **Rekomendasi: tegakkan — `HASIL1` WAJIB diperiksa, dan kegagalan ditampilkan.** Ini
**penyimpangan sadar** yang sekeluarga dengan keputusan "gagal terang-terangan" pada ambang tutup
buku PremiumList Life. Penyimpanan yang ditolak procedure **tidak boleh** tampak berhasil.
⚠️ Arti nilai `HASIL1` sendiri belum diketahui — bagian dari **OQ-002**.

---

## Rekap frontier Ronde 1

| # | Jenis | Pemblokir spec? |
| --- | --- | --- |
| Q1 denormalisasi | desain | ya — bentuk skema |
| **Q2 enam kolom limit** | **fakta bisnis** | **ya** |
| **Q3 dua tingkat share** | **fakta bisnis** | **ya** — menentukan arti angka |
| Q4 total 100% | desain | ya |
| Q5 kaskade hapus | desain (+ butuh DBA) | ya |
| **Q6 tahun tak dapat dihapus** | **fakta bisnis** | tidak |
| **Q7 enumerasi `REINSTYPEID`** | **fakta bisnis** | sebagian — sumbernya sudah diketahui |
| Q8 gerbang tahun | desain | tidak |
| **Q9 procedure apa adanya** | desain (+ **butuh DBA**) | **ya — pemblokir utama** |
| Q10 terapkan-ke-semua | desain | tidak |
| Q11 galat procedure | desain | ya |

**Frontier belum kosong.** Ronde 2 terbuka setelah Q1–Q5, Q9, dan Q11 dijawab — terutama untuk
mendalami isi keempat Harness (444–575 KB, baru di-grep), kedua FlowAction tampilan (`ViewRate`,
`ViewRateTable`), dan `RIRATEID`/`RIRATE` pada business yang belum saya telusur.
