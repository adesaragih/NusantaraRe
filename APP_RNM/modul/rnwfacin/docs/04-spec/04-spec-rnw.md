# Spesifikasi — Siklus Renewal (RNW)

> **Lingkup:** siklus **Renewal** (`RNW Fac In\`). Pendamping `03-spec-modul-terverifikasi.md` (NB) —
> **spec NB tidak diubah oleh dokumen ini**.
>
> Dasar keputusan: **K-030 … K-039**. Label mengikuti `CLAUDE.md` §3. Kosakata mengikuti
> `../steering/GLOSARIUM.md`.
>
> ⛔ **Tidak ada seam baru.** Ketiga seam yang disetujui untuk NB — `premium.Calculate` ·
> `acceptance.Next` · `rules.Eval` — **diwarisi apa adanya**. Menambah seam untuk renewal berarti
> mengarang pemisahan yang tidak ada di sistem lama.

---

## Problem Statement

Renewal berjalan di atas **basis kode yang sama dengan New Business**. `[terverifikasi]` **1.907 dari
1.927** berkas RNW **identik byte-per-byte** dengan berkas bernama sama di NB — nol berbeda, metadata
ekspor termasuk (K-031).

Yang membedakan renewal hanyalah **cara sebuah kasus masuk dan ditampilkan**.

Tanpa spec ini, dua kesalahan berlawanan sama-sama mungkin, dan keduanya mahal:

1. **Membangun ulang perhitungan yang sudah ada** — menghabiskan waktu untuk modul yang secara harfiah
   berkas yang sama, lalu melahirkan dua implementasi yang bisa menyimpang.
2. **Melewatkan 14 berkas yang memang khas renewal** — sehingga siklus renewal tidak punya pintu masuk
   dan layar sendiri.

---

## Solution

**Pakai ulang seluruh modul NB tanpa perubahan.** Bangun **14 berkas delta** sebagai lapisan **masuk
dan tampilan** di atasnya.

| Yang diwarisi utuh dari NB | Yang dibangun untuk renewal |
| --- | --- |
| `pkg/money` · `pkg/ratio` · `internal/rules` · `services/acceptance` · `services/premium` · `services/spreading` | **14 berkas delta**: alur masuk · layar input & periode · kelompok bisnis · daftar & konversi |
| `03-spec-modul-terverifikasi.md` + **16 tiket** di `../05-tickets/` | tiket tersendiri, menyusul |

`[terverifikasi]` **Renewal tidak punya mesin hitung ulang.** Hanya empat activity RNW yang
menggerbangi `StatusBusiness = 2` — masa berlaku tanggal, validasi tanggal, proteksi spreading, input
pembayaran — dan **tidak satu pun perhitungan premi**. Keempatnya berkas bersama yang identik dengan NB.

📌 Ini **menjelaskan K-015 dari sisi mekanisme**: renewal dinilai atas nilai pertanggungan **penuh**
bukan karena ada kebijakan yang memilih demikian, melainkan karena **tidak ada perhitungan selisih
sama sekali** — yang ada hanya perhitungan NB di atas data kasus baru.

**Delta = 14 berkas**, hasil `20 − 2 (K-036) − 4 (K-038)`.

---

## User Stories

### Kelompok 1 — Alur masuk (3 berkas)

`Flow\InputRenewalFacultativeIn` · `FlowAction\Renewal_FlowAct` · `FlowAction\Renewal_FlowAct_IsUW`

1. Sebagai underwriter, saya ingin memulai renewal dengan memasukkan **No. Polis + Renewal Date +
   Note** lalu menekan **OK**, agar kasus renewal tercipta tanpa saya mengetik ulang seluruh data polis.
2. Sebagai underwriter, saya ingin **data polis lama tersalin sebagai nilai awal** pada kasus baru,
   agar saya cukup menyunting yang berubah.
3. Sebagai pengembang, saya ingin titik masuk renewal **ditetapkan eksplisit di sistem baru**, karena
   penentunya di sistem lama ada di konfigurasi yang tidak terekspor.
4. Sebagai penguji, saya ingin alur masuk renewal berakhir pada kasus yang **dapat diproses tangga
   akseptasi NB tanpa penyesuaian**.

`[terverifikasi]` `RNW Fac In\FlowAction\Renewal_FlowAct.xml`:

| Tag | Nilai |
| --- | --- |
| `pySectionReference` | `InputRenewal` |
| `pyPreProcessingActivity` | `InputOfferFacInEngineer_preACT` |
| `pyLocalActionActivity` | `SetValidateDate_PostAct` |
| `pyActionTransformRule` | `AddToListSuggestOfferFacIn_DT` |

`[terverifikasi]` `Flow\InputRenewalFacultativeIn` memuat **75 shape** dan merujuk **19 rule `When`** —
seluruhnya predikat routing tangga akseptasi yang sama dengan NB; **tidak ada predikat gerbang-masuk**
di antaranya.

Referensi layar: berkas gambar `DDL\halaman depan renewal.JPG` **ada** (30 KB).
⛔ Isinya **tidak dibaca** sebagai fakta — hanya keberadaannya dirujuk.

### Kelompok 2 — Layar input & periode (7 berkas)

`Section\InputRenewal` · `InputRenewal_IsUW` · `InputRenewalDtl` · `InputRenewalDtl_IsUW` ·
`PeriodeRenewal` · `PeriodeRenewal_IsUW` · `SFAPortal_Renewal`

5. Sebagai underwriter, saya ingin layar renewal menampilkan **detail penawaran yang sama seperti NB**,
   agar saya tidak perlu mempelajari layar kedua untuk pekerjaan yang sama.
6. Sebagai pengembang, saya ingin layar renewal **menyisipkan komponen NB**, bukan menduplikasinya,
   agar perbaikan pada satu komponen berlaku di kedua siklus.
7. Sebagai underwriter, saya ingin **periode polis baru** beserta `OldPolicyNo` dan tanggal renewal
   terlihat di satu layar.
8. Sebagai underwriter, saya ingin melihat **daftar kandidat renewal** dari portal, agar saya tahu
   polis mana yang mendekati jatuh tempo.

`[terverifikasi]` **Layar renewal adalah wadah, bukan duplikat.** Section yang disisipkan (dibaca dari
`<pyInclude>` **dan** `<pySection>`):

| Layar | Jumlah section disisipkan |
| --- | ---: |
| `InputRenewalDtl_IsUW` | **17** |
| `InputRenewalDtl` | **14** |
| `InputRenewal` | **7** |
| `InputRenewal_IsUW` | **6** |
| `PeriodeRenewal` | 4 |
| `PeriodeRenewal_IsUW` | 2 |
| `SFAPortal_Renewal` | 0 |

Contoh `[terverifikasi]` — `InputRenewal` menyisipkan `AllSummarySection` · `PeriodeRenewal` ·
`InputRenewalDtl` · `InputDtlObject_FacIn` · `InputInwardFacultativeSuggest` · `EmailSection` ·
`EmailSectionCeding`.

`[terverifikasi]` `InputRenewal` dan `InputRenewal_IsUW` masing-masing hanya punya **2 properti
sendiri** (`.IsShowDetail` + template). Isinya seluruhnya datang dari section yang disisipkan.

`[terverifikasi]` `Section\PeriodeRenewal` mengikat **22 properti**, termasuk
`.QuotationData.OldPolicyNo` · `.QuotationData.RNWDate` · `.QuotationData.StatusBusiness` ·
`.PolicyData.StartDateTime` · `.PolicyData.EndDateTime`.

`[terverifikasi]` `SFAPortal_Renewal` mengikat 9 properti — bentuknya sama dengan daftar kandidat
(`NBStatus`, `NBStatusNew`, tanggal polis, `OldPolicyNo`).

> ### ✅ Lima section yang disisipkan tetapi tidak ada di RNW — **TERJAWAB 18 September 2026**
>
> Dari **35** section unik yang disisipkan ketujuh layar, **30 ada** sebagai berkas RNW dan **5 tidak**.
> Kelimanya sudah diputuskan, dan **tidak satu pun merupakan rujukan menggantung**:
>
> | Section | Keputusan |
> | --- | --- |
> | `InputInwardFacultativeSuggest` · `InwardFacIn` · `OfferFacIn_NusaRe` · `OfferFacIn_NusaRe_IsUW` | **K-040 — DIWARISI dari NB.** `[terverifikasi]` keempatnya ada di `NB FacIn\Section\`. Konsekuensi wajar K-031: layar renewal menyisipkan section NB apa adanya. **Bukan delta baru** — delta tetap 14 |
> | `InputDtlObject` | **K-041 — USANG**, digantikan `InputDtlObject_FacIn`. `[terverifikasi]` tidak ada di NB maupun RNW; penggantinya ada di keduanya dan muncul **26×** di berkas yang sama. Empat sisipan tersisa (L2764 · L2944 · L3207 · L3390) = **sisa kelewat** |
>
> ⚠️ **`InputDtlObject` ≠ `InputDtlObject_FacIn`** — dua section berbeda, dibedakan hanya oleh sufiks.
> Jangan disatukan.
>
> ⛔ Keempat sisipan usang **diport apa adanya** (`CLAUDE.md` §1) — jangan dihapus diam-diam. Karena
> section tujuannya tidak ada, sisipan itu tidak menghasilkan apa pun dan menjadi tidak berdampak
> dengan sendirinya.

### Kelompok 3 — Kelompok bisnis (2 berkas)

`Activity\GetBusinessGroup_Act` · `RDBList\CariBusinessGID`

9. Sebagai sistem, saya ingin menurunkan **kelompok bisnis** dari kode bisnis penawaran, agar
   klasifikasi kasus renewal konsisten dengan sistem lama.

`[terverifikasi]` `GetBusinessGroup_Act`, 4 langkah kedalaman penuh:

| Langkah | Metode | Isi |
| ---: | --- | --- |
| 1 | `Property-Set` | `InputData.CARI2` ← `…QuotationData.BusinessCode` |
| 2 | `RDB-List` | → `CariBusinessGID` |
| 3 | `Property-Set` | `ParamBis.CARI10` ← `OutBis.pxResults(1).CARI3` |
| 4 | `Page-Remove` | bersihkan halaman |

`[terverifikasi]` `CariBusinessGID` `<pyBrowseSQL>`:

```sql
select ID, OLDID, NOTE as "Note", GROUPPANEL as "GroupPanel", BusinessGroupID as CARI3
from business where ID = {InputData.CARI2}
```

⚠️ **Blok `<pySaveSQL>` yang memanggil `POOLDATA.PROSESCOPY` DIBUANG** (K-037) — perubahan perilaku
yang sudah dicatat bernomor. **`<pyBrowseSQL>` tetap diport apa adanya.**

⚠️ Ini menurunkan **kelompok bisnis**, **bukan** lini bisnis (COB) yang menggerakkan skala rasio
K-018. COB tetap ditentukan predikat `IsFire`/`IsPA`/`IsMBU`/… yang diwarisi dari NB.

### Kelompok 4 — Daftar & konversi (2 berkas)

`ReportDefinition\RenewalList_RD` · `Activity\serviceInsertArasapasRNW_act`

10. Sebagai underwriter, saya ingin **daftar kandidat renewal** yang menampilkan nomor polis lama dan
    tanggal berakhirnya, agar saya dapat memilih polis yang akan diperpanjang.
11. Sebagai sistem, saya ingin kasus renewal yang disetujui **dikonversi ke produksi dengan jalur yang
    sama seperti NB**, agar tidak ada dua jalur simpan yang bisa menyimpang.

`[terverifikasi]` `RenewalList_RD` — kelas `ASM-FW-GISFW-Work-Renewal`, **11 kolom**:
`.pyID` · `.pzInsKey` · `.pxCreateDateTime` · `.pxCreateOperator` · `.Quotation.OldPolicyNo` ·
`.Quotation.InsuredName` · `.Quotation.MarketingName` · `.OfferFacIn.PolicyData.StartDateTime` ·
`.EndDateTime` · `.NBStatus` · `.NBStatusNew`.

⚠️ `Work-Renewal` adalah **irisan pelaporan**, bukan kelas kerja tersendiri: `[terverifikasi]`
`pyWorkClass` pada flow renewal adalah `ASM-FW-GISFW-Work` — sama dengan NB — dan `Work-Renewal`
tidak dipakai rule lain mana pun.

`[terverifikasi]` `serviceInsertArasapasRNW_act` — kelas `ASM-FW-GISFW-Work`, **10 langkah**:
`Call serviceInsertArasapas_act` → `Property-Set` → `RDB-List` → `Property-Set` →
**`Call …Int-M_LINK_SERVICE.GetLinkService`** → **`Connect-REST`** → `Property-Set` →
`Call InsertLogServiceProd` → `Page-Remove` → `RDB-List`.

Pola **endpoint dari `M_LINK_SERVICE` → `Connect-REST` → catat log** menguatkan `CLAUDE.md` §4.4
langsung dari korpus. Rancangan konversi renewal **mengikuti pola NB** (K-035).

---

## Implementation Decisions

**Tidak ada modul baru di lapisan perhitungan.** `services/premium`, `services/acceptance`,
`internal/rules`, `pkg/money`, `pkg/ratio`, `services/spreading` dipakai **apa adanya** dari spec NB.
Renewal tidak menambah satu pun cabang di dalamnya.

**Siklus dibedakan oleh `StatusBusiness`, bukan oleh basis kode terpisah.** `[terverifikasi]` empat
activity bersama menggerbanginya, dan tidak satu pun perhitungan premi.

**Layar renewal disusun sebagai komposisi, bukan salinan.** Komponen NB disisipkan; komponen khas
renewal hanya wadah dan field periode. Menduplikasi komponen NB akan melahirkan dua layar yang
menyimpang diam-diam.

**`.OldTSI` diport apa adanya dari XML** (`CLAUDE.md` §1 — reproduksi perilaku terekam).
`[terverifikasi]` ia ditulis activity **perhitungan**, dari nilai terhitung:

```
Activity\CountPremiumNet             L2792 · L5505 · L7325   .OldTSI = Local.TSISpreadTotal
Activity\CountPremiumNetElse         L1781 · L3916 · L5751   .OldTSI = Local.TSISpreadTotal
Activity\SumTSIPremiSpreadedRNM_Act  L8391                   .OldTSI = Local.TsiLoLSpread + .OldTSI
```

⛔ **Jangan menafsir ulang, jangan menyatukan dengan mekanisme lain, jangan "merapikan".** Ikuti isi
XML.

**Jalur simpan produksi renewal = jalur simpan NB.** `[terverifikasi]` `SaveJsonPolicyFacIn_Act`
identik byte-per-byte antara NB dan RNW (371.819 byte, SHA-256 sama). Renewal **tidak memerlukan
penghasil nomor polis tersendiri**.

**Ketertelusuran (`CLAUDE.md` §4.6)** berlaku sama seperti NB: setiap query, predikat, dan transisi
menyebut rule Pega asalnya dalam komentar.

---

## Testing Decisions

**Modul inti tidak diuji ulang.** Ketiga seam NB sudah punya kriteria penerimaannya sendiri di
`03-spec-modul-terverifikasi.md` dan 16 tiket NB. Menguji ulang perhitungan yang berkasnya identik
hanya menambah waktu tanpa menambah keyakinan.

**Yang diuji pada renewal adalah lapisan masuk dan tampilan:**

1. Alur masuk menghasilkan kasus yang **diterima tangga akseptasi NB tanpa penyesuaian**.
2. Data polis lama **tersalin sebagai nilai awal** pada kasus baru.
3. Layar renewal **menyisipkan** komponen NB — bukan menyalinnya. Diuji dengan memastikan perubahan
   pada komponen NB terlihat di layar renewal.
4. Kelompok bisnis diturunkan benar dari kode bisnis, dan **blok `pySaveSQL` yang dibuang tidak
   dieksekusi**.
5. Konversi produksi renewal memakai **jalur yang sama** dengan NB.

**Rekonsiliasi eksak tetap berlaku** (ADR-0001): karena perhitungan renewal adalah perhitungan NB,
kasus renewal dari berkas fixture ter-de-identifikasi harus cocok **sampai digit terakhir** dengan
sistem lama — memakai kerangka rekonsiliasi tiket 16 NB, tanpa pembanding baru.

---

## Out of Scope

| Butir | Alasan / menunggu |
| --- | --- |
| **Modul inti** — perhitungan · registry predikat · tangga akseptasi · spreading | **Diwarisi dari NB (K-031)**, tidak dispec ulang. Tidak menunggu apa pun |
| **`EditMarketing`** (Section + FlowAction) | **Dibuang** (K-036); pemilihan marketing mengikuti NB. Temuan keamanan `.Password` **gugur** karena komponennya tidak diport |
| **Fitur pilih-tertanggung** — `Harness\ChooseInsured` · `Section\ChooseInsuredDtl` · `ReportDefinition\BrowseAccountInsuredEDM` · `Activity\SetDataInsuredEDM_Act` | **Dibuang** (K-038). Renewal mengambil tertanggung dari **polis lama**, bukan pemilihan CRM/SFA. Dua di antaranya berkelas SFA |
| **Gerbang masuk teknis** (work type / portal Pega) | **Tidak terekspor.** Sistem baru menetapkannya eksplisit sesuai K-034 — **keputusan bisnis, bukan porting** |
| **Rule konversi asli kelas `Work`** | **Celah K-004 tetap terbuka.** Arah rancangan = pola NB; rule aslinya belum ada di korpus mana pun |
| **Blok `pySaveSQL` `PROSESCOPY`** | **Dibuang** (K-037) — perubahan perilaku yang sudah dicatat bernomor |
| **Isi `InputRenewalDtl` per-field** (4,65 MB bersama `_IsUW`) | **Ditunda ke tahap tiket.** Strukturnya sudah terpetakan; rincian field belum diperlukan untuk spec |
| **`Activity\SetOldData` dan properti ber-sufiks `TSIOld`** | ⏸ **Ditunda ke fase Endorsement (EDM).** `[dugaan]` kemungkinan milik fase EDM; **tidak dianalisis lebih jauh** dan **tidak masuk lingkup delta renewal**. ⚠️ `TSIOld` (sufiks) **≠** `.OldTSI` (properti perhitungan) — dua hal berbeda |
| **`Section\InputDtlObject`** (tanpa sufiks) | ⏸ **Usang** (K-041), digantikan `InputDtlObject_FacIn`. Empat sisipan tersisa **diport apa adanya**, tidak dibangun sebagai komponen baru |
| ~~Lima section yang disisipkan tetapi tidak ada di RNW~~ | ✅ **TERJAWAB** (K-040 · K-041) — empat **diwarisi dari NB**, satu **usang**. Bukan lagi pertanyaan terbuka, dan **tidak menambah delta** |

---

## Further Notes

### Mengapa renewal begitu tipis — dan mengapa itu bukan kebetulan

Tiga pengukuran independen menunjuk hal yang sama:

1. **Berkas:** 1.907 dari 1.927 identik byte-per-byte.
2. **Perhitungan:** nol activity renewal yang menggerbangi perhitungan premi lewat `StatusBusiness`.
3. **Layar:** dua layar masuk utamanya hanya punya **2 properti sendiri**; sisanya section NB yang
   disisipkan.

Renewal bukan siklus yang mirip NB — ia **NB yang dimasuki lewat pintu lain**.

### Tiga jebakan penamaan yang sudah terbukti di siklus ini

Ketiganya menghasilkan kesimpulan salah sebelum dikoreksi, dan ketiganya bentuknya berbeda:

| Nama | Jebakannya |
| --- | --- |
| `isApproved` | **properti** vs DecisionTable — beda **tipe rule** |
| `GetInsuredID` | **Activity** (ada di RNW) vs **RDBList** (NB-only) — beda **tipe rule**, nama sama persis |
| `TSIOld` vs `.OldTSI` | beda **posisi kata** — sufiks vs prefiks, dua properti berlainan |
| `InputDtlObject` vs `InputDtlObject_FacIn` | beda **ada-tidaknya sufiks** — dua section berlainan; yang satu usang, yang lain dipakai luas (K-041) |

**Konsekuensi kerja:** tipe rule dan tag pembawa wajib diperiksa; pencocokan nama tidak pernah cukup.

### Yang berubah bila ekspor RNW yang lebih lengkap tiba

Lima section pada catatan Kelompok 2 dan celah K-004 adalah satu-satunya tempat spec ini bergantung
pada kelengkapan ekspor. Sisanya bersandar pada berkas yang **ada dan terbaca**.

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan. Berkas gambar dirujuk keberadaannya saja.*
