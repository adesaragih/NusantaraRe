# Discovery Renewal — Penutup

> **Cakupan:** `D:\migrasi\RNM\RNW Fac In\` (1.927 berkas). `Endorsment Fac In\` dan korpus Treaty
> **tidak dibaca** (K-030, K-005). Empat putaran, metode parser XML dengan pemeriksaan tag pembawa.
>
> ⚠️ **Dokumen ini membedakan dua tingkat bukti secara eksplisit:**
> **`[terverifikasi]`** = sudah **diperiksa ulang work owner** ke korpus ·
> **`[temuan Claude]`** = diukur dan dilaporkan, **belum** diverifikasi ulang.
> Tidak ada butir yang dinaikkan tingkatnya tanpa verifikasi.

---

## 1. Hasil pokok

### `[terverifikasi]` RNW **adalah** NB, minus 176 berkas, plus 20

| Ukuran | Nilai |
| --- | ---: |
| Berkas bersama NB↔RNW | **1.907** |
| — identik **byte-per-byte** | **1.907** (100 %) |
| — berbeda | **0** |
| Hanya di RNW | **20** |
| Hanya di NB | **176** |

Seluruh delta renewal muat dalam **20 berkas**. Pertanyaan "pakai ulang modul NB atau implementasi
tersendiri" terjawab untuk 1.907 berkas sekaligus: **pakai ulang**.

### Kedua puluh berkas delta — terpetakan seluruhnya

| Kelompok | Isi |
| --- | --- |
| Alur masuk | `Flow\InputRenewalFacultativeIn` (75 shape, 19 rule `When`) · `FlowAction\Renewal_FlowAct` → Section `InputRenewal` · `Renewal_FlowAct_IsUW` |
| Layar input | `Section\InputRenewal` · `InputRenewal_IsUW` · `InputRenewalDtl` (2,25 MB) · `InputRenewalDtl_IsUW` (2,40 MB) |
| Periode & polis lama | `Section\PeriodeRenewal` (22 properti) · `PeriodeRenewal_IsUW` |
| Portal | `Section\SFAPortal_Renewal` (kelas `Data-Portal`) |
| Pemilih tertanggung | `Harness\ChooseInsured` · `Section\ChooseInsuredDtl` (3 properti) · `ReportDefinition\BrowseAccountInsuredEDM` · `Activity\SetDataInsuredEDM_Act` |
| Kelompok bisnis | `Activity\GetBusinessGroup_Act` (4 langkah) · `RDBList\CariBusinessGID` |
| ~~Sunting marketing~~ | ~~`FlowAction\EditMarketing` · `Section\EditMarketing`~~ ⛔ **DIKELUARKAN dari lingkup port (K-036, 18 Sept 2026)** — pemilihan marketing mengikuti NB; temuan keamanan `.Password` **gugur** karena komponennya tidak diport |
| Daftar renewal | `ReportDefinition\RenewalList_RD` (11 kolom, kelas `Work-Renewal`) |
| Konversi produksi | `Activity\serviceInsertArasapasRNW_act` (10 langkah, kelas `Work`) |

---

## 2. ⛔ TIGA rujukan menggantung di RNW — bukan satu

Putaran 3 menemukan satu; pemeriksaan per-nama atas seluruh 40 nama sisa menemukan **dua lagi**.
Setiap vonis berdasar **tag pembawa**, bukan nama.

| Rule | Tipe | Mekanisme rujukan | Perujuk di RNW | Status |
| --- | --- | --- | --- | :-: |
| **`SumTreatyCapacity_Act`** | Activity | `<pyStepsActivityName>Call …` | `CountASMShareTotal_ACT` L444 · `SetValidateDate_Act` L1129 | `[terverifikasi]` |
| **`GenerateNoPolicy`** | RDBList | `<RequestType>` | `GenerateNopolis_Act` L3179 · `SaveJsonPolicyFacIn_Act` L2204 | `[temuan Claude]` |
| **`SearchJobID`** | RDBList | `<RequestType>` | `UploadCSVPerson_PostAct` L2404 | `[temuan Claude]` |

Ketiganya: berkasnya **ada di NB, tidak ada di RNW**, sementara perujuknya **ada di RNW**.

⚠️ **`GenerateNoPolicy` yang paling berdampak** — ia menghasilkan **nomor polis**, dan dirujuk dari
jalur simpan polis (`SaveJsonPolicyFacIn_Act`). Bila benar menggantung di renewal, jalur produksi
renewal kehilangan penghasil nomor polis.

⚠️ **Bukan vonis.** `CLAUDE.md` §4.5: **"hilang dari ekspor" ≠ "usang"** — dan di sini bobotnya lebih
besar dari biasanya, karena kedua folder terbukti diekspor dari keadaan yang sama (1.907/1.907
identik). Satu berkas yang hadir di satu sisi saja **lebih** mencurigakan, bukan kurang.

### Yang TERNYATA BUKAN menggantung — dan mengapa

Empat kasus yang tampak menggantung dari namanya, ternyata tidak:

| Nama | Kenyataan |
| --- | --- |
| `isApproved` | **properti** (`Rule-Obj-Property`) pada dua kelas, bukan DecisionTable — `[terverifikasi]` |
| `IsUWAccepted` | **DecisionTable** bernama sama **ada** di RNW; hanya rule `When`-nya NB-only |
| `GetLimitAkseptasi_ActFlow` | hanya di `<pzOriginalInstanceKey>` — **metadata *save-as***, bukan panggilan — `[terverifikasi]` |
| `GetInsuredID` | **`Activity\GetInsuredID` ADA di RNW**; yang NB-only adalah `RDBList\GetInsuredID`. Panggilan `Call GetInsuredID` terselesaikan ke Activity |

Dan `GetKurs`: `<RequestType>GetKurs</RequestType>` → **0 berkas**. Kemunculannya hanya sebagai **nama
halaman**, bukan pemanggilan rule.

📌 **Pelajaran metodologis yang dipaksakan data ini:** dari 40 nama, **tidak ada satu aturan umum**
yang memisahkan rujukan sungguhan dari tabrakan nama. Empat mekanisme berbeda muncul — `pyStepsActivityName`,
`RequestType`, `Rule-Obj-Property`, `pzOriginalInstanceKey` — dan dua nama yang sama persis
(`GetInsuredID`) berperilaku berbeda menurut **tipe rule**-nya. Pemeriksaan **per nama** tidak dapat
digantikan.

**Sisa 33 nama lain** dari 40: seluruhnya **bukan panggilan** — muncul sebagai `pySection`,
`pyInclude`, `pyLocalAction`, `pyActivity`, nama halaman, atau metadata. Tidak ada yang menggantung.

---

## 3. Status R-1 … R-12

| # | Pertanyaan | Jawaban | Tingkat bukti |
| --- | --- | --- | --- |
| R-1 | `isApproved` menggantung? | **Tidak** — properti, bukan DecisionTable. `"1"`=Accept / `"0"`=Reject, **ruang nilai berbeda** dari `ProposalAcceptStatus` (K-021) — jangan disatukan | `[terverifikasi]` |
| R-2 | Gerbang masuk renewal | `IsOfferFacIn` **0 rujukan**; flow merujuk 19 rule `When`, semuanya routing → menjadi R-8 | `[temuan Claude]` |
| R-3 | `Work-Renewal` antrean tersendiri? | **Bukan** → dikonfirmasi ulang di R-12 | `[temuan Claude]` |
| R-4 | Penentuan kelompok bisnis | `BusinessCode` → tabel `business` → `BusinessGroupID`. ⚠️ Ini **kelompok bisnis**, bukan COB penggerak skala K-018 | `[temuan Claude]` |
| R-5 | Ambil polis lama & hitung ulang? | → dijawab bersama R-9 | — |
| R-6 | 176 NB-only wajar tak dipakai? | **132 bersih**; dari 44 sisa → **3 menggantung**, 41 tabrakan/bukan-panggilan | sebagian `[terverifikasi]` |
| R-7 | Tipe rujukan 40 nama | **Selesai seluruhnya** — lihat §2 | campuran |
| R-8 | Gerbang masuk di luar korpus | **Ya** — konfigurasi work type/portal tidak terekspor → **keputusan bisnis, bukan porting** | `[temuan Claude]` |
| R-9 | Di mana hitung ulang renewal? | **Tidak ada hitung ulang.** Renewal menjalankan perhitungan NB yang sama; kaitan polis lama hanya **rujukan** `OldPolicyNo`. Hanya 4 activity bergerbang `StatusBusiness=2`, tak satu pun premi | `[temuan Claude]` |
| R-10 | `PROSESCOPY` | Blok PL/SQL **berargumen literal keras**, memanggil prosedur yang tidak ada di DB, keluaran ke `dbms_output`, tidak berhubungan dengan browse SQL di rule yang sama. **Artinya tidak disimpulkan** | `[temuan Claude]` |
| R-11 | Celah K-004 terisi? | **TIDAK.** Yang hilang adalah `serviceInsertArasapas_act` **kelas `Work`** (tanpa sufiks). `serviceInsertArasapasRNW_act` adalah **saudara**, bukan rule itu | `[terverifikasi]` |
| R-12 | `Work-Renewal` | **Irisan pelaporan** — dipakai hanya `RenewalList_RD`; `pyWorkClass` flow = `ASM-FW-GISFW-Work` | `[temuan Claude]` |

---

## 4. Temuan yang menyentuh keputusan terkunci

### ⛔ Menyentuh **K-006** — `SumTreatyCapacity_Act`

`[terverifikasi]` Cabang 4 K-006 ("kapasitas treaty") **putus satu hop lebih awal di renewal**:

```
NB : CountASMShareTotal_ACT / SetValidateDate_Act → SumTreatyCapacity_Act (ADA) → CountTreatyCapacity_Act (tiba lewat DDL\)
RNW: CountASMShareTotal_ACT / SetValidateDate_Act → SumTreatyCapacity_Act (HILANG) ✗
```

Amandemen K-006 (17 September 2026) memulihkan keenam cabang **berdasarkan korpus NB**. Temuan ini
**tidak membatalkannya**, tetapi menunjukkan keadaan RNW berbeda untuk cabang 4.

### ⛔ Menyentuh **K-004** — celah tetap terbuka

`[terverifikasi]` Peta lengkap activity konversi:

| Activity | NB | RNW | Kelas | Langkah |
| --- | :-: | :-: | --- | ---: |
| `serviceInsertArasapas_act` | ada | ada | `Data-PolicyTreatyIn` | 1 (stub) |
| `serviceInsertArasapasEDM_act` | ada | ada | `ASM-FW-GISFW-Work` | 20 |
| `serviceInsertArasapasRNW_act` | — | ada | `ASM-FW-GISFW-Work` | 10 |
| **`serviceInsertArasapas_act` kelas `Work`** | **tidak ada** | **tidak ada** | sasaran delegasi stub | — |

Yang bernilai: dua saudara terbaca memperlihatkan **polanya** —
`GetLinkService` (dari `M_LINK_SERVICE`) → `Connect-REST` → catat log layanan. `[terverifikasi]` Pola
ini menguatkan `CLAUDE.md` §4.4 langsung dari korpus.

### ⚠️ Temuan keamanan — `Section\EditMarketing`

`[temuan Claude]` Mengikat empat properti: `.ID` · **`.Password`** · `.QuotationData.MOID` ·
template, pada kelas `Data-OfferFacIn`. Otorisasi tertanam di lapisan data; apa yang divalidasi dan
terhadap apa **belum ditelusuri**. Perlu tinjauan keamanan sebelum diport (`CLAUDE.md` §6).

⛔ Tidak ada nilai kredensial disalin ke dokumen mana pun — hanya nama propertinya.

---

## 5. `Struktur_InputRenewalFacultativeIn.xlsx` — **bukan data pelanggan**

`[temuan Claude]` Dugaan awal bahwa berkas ini memuat data pelanggan **tidak terbukti**. Baris 1
dibaca — **hanya baris 1**, sesuai batasan:

| Kolom | Isi header |
| ---: | --- |
| 1 | `No` |
| 2–20 | (kosong pada baris 1) |
| 21–28 | `Nama Rule` · `Jenis Rule` · `Bagian` · `Shape ID` · `Dipanggil dari` · `XML` · `Class` · `Status XML` |

Ini **peta struktur rule** buatan tim — inventaris rule alur renewal, bukan data transaksi.
81.379 baris × 28 kolom, sheet tunggal `Struktur`, format OOXML **sah**.

⚠️ **Satu kehati-hatian sebelum isinya dibaca:** kolom **`XML`** kemungkinan memuat cuplikan XML rule
Pega, dan XML rule Pega **memuat nama operator** (`pxCreateOpName`) — terbukti berulang kali di korpus.
Pembacaan isinya wajib menyaring kolom itu.

📌 Berpotensi bernilai tinggi: kolom `Dipanggil dari` dan `Status XML` dapat **mengonfirmasi silang**
ketiga rujukan menggantung di §2 dari sumber independen.

---

## 6. Kandidat keputusan work owner

| # | Kandidat | Menyentuh |
| ---: | --- | --- |
| 1 | `SumTreatyCapacity_Act` menggantung di RNW — cabang kapasitas treaty tidak berlaku di renewal, atau ekspor RNW tidak lengkap? | **K-006** |
| 2 | `GenerateNoPolicy` menggantung — penghasil **nomor polis** pada jalur simpan polis renewal | jalur produksi |
| 3 | `SearchJobID` menggantung pada jalur unggah CSV | — |
| 4 | Bolehkah pola `serviceInsertArasapasEDM/RNW_act` dipakai merancang jalur konversi, meski bukan rule aslinya? | **K-004** |
| 5 | Gerbang masuk renewal tidak terekam → sistem baru **wajib menetapkannya eksplisit** | keputusan bisnis |
| 6 | `Section\EditMarketing` meminta ID + kata sandi — perlu tinjauan keamanan | `CLAUDE.md` §6 |
| 7 | `PROSESCOPY` — blok literal memanggil prosedur yang tidak ada; diport apa adanya atau ditinjau? | — |
| 8 | Baca isi `Struktur_…xlsx` (81.379 baris) sebagai konfirmasi silang, dengan penyaringan kolom `XML` | — |

---

## 7. Yang tidak dikerjakan, dan mengapa

| Butir | Alasan |
| --- | --- |
| Isi `Section\InputRenewalDtl` + `_IsUW` (4,65 MB) | Belum dibedah per-field; strukturnya sudah terpetakan sebagai layar input renewal |
| Isi `Struktur_…xlsx` (81.379 baris) | Volume memerlukan putaran tersendiri; kolom `XML` perlu penyaringan nama operator |
| `Endorsment Fac In\` | Di luar fokus (K-030) |
| Korpus Treaty tersendiri | Tidak dibaca (K-005) |

---

## Peta dokumen discovery RNW

| Putaran | Dokumen | Isi |
| ---: | --- | --- |
| 1 | `01-inventaris-dan-delta-rnw.md` | Inventaris + 1.907/1.907 identik + 20/176 delta |
| 2 | `02-bedah-delta-dan-jawaban-r1-r6.md` | Bedah 20 berkas + R-1…R-6 |
| 3 | `03-jawaban-r7-r12-dan-sisa-delta.md` | R-7…R-12 + sisa berkas |
| 4 | `04-penutup-discovery-rnw.md` | **Dokumen ini** — penutup |

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan, tanpa nilai kredensial.*
