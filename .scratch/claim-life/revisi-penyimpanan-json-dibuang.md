# Revisi Penyimpanan Claim Life — JSON DIBUANG, ClaimData → tabel baru relasional

Tanggal: 2026-09-15 (diperbarui 2026-09-16 setelah audit Claude + penutupan 5 OQ work owner)
Sumber: sesi grilling ulang dengan work owner + sensus property (context-gatherer) + **audit ulang
Claude (skill grilling)** yang menemukan 6 masalah pemblokir pada skema kandidat 5-tabel.
Status: `[keputusan work owner]` — merevisi bagian penyimpanan spec Claim Life (JSON → relasional).

> Konteks: spec Claim Life & 13 tiket sudah terbit. Revisi ini menyentuh **cara penyimpanan** —
> perlu diselaraskan ke spec + tiket yang menyentuh persistence.

---

## Keputusan inti `[keputusan work owner]`

1. **SEMUA JSON dibuang** di Claim Life. Tidak ada blob JSON di sistem baru.
   - `JSON_KLAIM` / serialisasi `ClaimData` → dibuang; data disimpan **relasional/flat**.
   - Konversi ke Arasapas/produksi → **TIDAK pakai JSON**; sistem tujuan **SELECT langsung dari flat
     table baru** (bukan payload JSON). `convertJsonNusareToProductionClaimLife` + serialisasi JSON
     keluar = dead, tidak dimigrasikan.
2. **Seluruh isi `ClaimData` → tabel BARU relasional** (`T_GENERAL_CLAIM` dkk). Yang dibuang: **JSON saja**
   (`JSON_KLAIM`/serialisasi ClaimData).
   ⚠️ **KOREKSI 2026-09-16 `[keputusan work owner — Tafsir A]`:** `OS_AKSEPTASI_KLAIM_LIFE` **TETAP
   di-INSERT** (flat), **berdampingan** dengan tabel relasional baru — BUKAN dibuang. Sistem baru
   menulis DUA tempat: tabel relasional baru + INSERT flat ke `OS_AKSEPTASI_KLAIM_LIFE` (karena hilir/
   Arasapas/produksi masih baca dari situ). Yang dibuang HANYA JSON. `M_LIFE_PREMIUM_DETAIL` = dibaca
   sebagai snapshot. → **Setiap AC/klaim yang berbunyi "OS_AKSEPTASI tidak ditulis / test menolaknya"
   adalah KELIRU — harus dibalik jadi: OS_AKSEPTASI tetap ditulis, JSON dibuang.**
3. **Kolom = property `ClaimData` yang DIPAKAI** (sensus context-gatherer + audit ulang Claude: di-set
   activity / ditampilkan section aktif). Nol property dikecualikan (tak ada yang di balik visibilitas
   mati `1=2`/`never` — sudah dicek per-field; uji tingkat-berkas atas `StsKatastrofe`/`KatastrofeNote`/
   `IsKPR` dibuang karena tak sahih).
4. **FK antar tabel klaim = `ON DELETE CASCADE`** + popup konfirmasi (pola Master Contract Retro Life
   / Treaty Contract Out). Hapus klaim → seluruh anak ikut, setelah konfirmasi Ya/Batal.
5. Uang/share/premi = **NUMBER** (decimal, ADR-0003); tanggal = **DATE**; semua kolom **nullable**
   (wajib-isi ditegakkan di Go). ID via sequence (ADR-0006).

⚠️ **Nama rule menyesatkan (OQ-066):** `InsertJsonKlaimLife_sql` `[terverifikasi]` isinya
`INSERT INTO OS_AKSEPTASI_KLAIM_LIFE` (51 kolom — 49 bind `CARI1`–`CARI49` + 2 hardcode) — **BUKAN**
JSON, INSERT flat. Nama "Json" hanya label. Jangan tiru sebagai blob.

---

## Hasil audit ulang Claude — 6 masalah pemblokir pada skema kandidat 5-tabel

Skema kandidat 5-tabel (adjustment digantung ke header) **cacat**. Audit Claude (skill grilling)
menemukan enam masalah, dua di antaranya pemblokir. Kiro sudah verifikasi dua temuan pemblokir ke
korpus — **LULUS**:

1. 🔴 **`T_CLAIMLF_ADJUSTMENT` salah induk.** Di korpus, `AdjustmentList` menggantung pada
   **`PremiumListDetail`**, bukan pada klaim. Bukti: `Claim Life/Activity/SavePesertaClaim.xml`
   (`ASM-FW-GCNMFW-WORK-CLAIMLIFE!SAVEPESERTACLAIM`) menulis 8 field ke
   `…PremiumListSummary.PremiumListDetail(<LAST>).AdjustmentList(<LAST>).*`. Setiap peserta punya
   daftar adjustment sendiri. **Perbaikan: FK adjustment → `T_CLAIMLF_PREMIUMLIST_DETAIL.ID`.**
2. 🔴 **`T_CLAIMLF_PREMIUMLIST_DETAIL` kehilangan 9 field** — termasuk `IsCheck` (penanda peserta
   dipilih untuk diklaim; **inilah penyimpan aturan "PremiumListSummary = hanya peserta yang
   diklaim"**) dan tiga tanggal per-peserta.
3. 🟠 `T_CLAIMLF_ADJUSTMENT` kehilangan 4 field: `CLAIM_AMOUNT`, `STS_REJECT`, `ACCEPTEDNO`,
   `ACCEPTATION_DATE`.
4. 🟠 `T_CLAIM_POLICY` kehilangan 8 field (3 `Tanggal*` + `ProductNameID`/`ProductName`/`TeamGroup`/
   `DateReceived`/`BranchName`/`BranchCode`).
5. 🟠 **PremiumListSummary** 1:1 dengan klaim (tak pernah ber-subscript) → **dilipat utuh ke header**.
   `PL_NUMBER` milik summary (`PolicyDataLife.PremiumListSummary.PL_NUMBER`), pindah ke header, bukan
   `T_CLAIM_POLICY`.
6. 🟠 **Sub-object dokumen per peserta** (`SaveOutStandingLife_Act` beriterasi `.DocumentList`, gerbang
   simpan "The document hasn't been uploaded person number") → **tabel baru**.

Verdict: lima tabel jadi **enam**, kedalaman tiga tingkat.

---

## Penutupan 5 OQ audit `[keputusan work owner 2026-09-16]`

| OQ | Pertanyaan | Keputusan |
| --- | --- | --- |
| **OQ-1** | Tanggal per peserta atau per klaim? | **Per peserta.** `CONFIRMATION_DATE`, `COMPLETE_DATE`, `CLAIM_RECEIVED_DATE` di `T_CLAIMLF_PREMIUMLIST_DETAIL`, **bukan** header. |
| **OQ-2** | Dokumen per peserta: tabel sendiri atau lampiran bawaan? | **Tabel sendiri: `DOCUMENT_CLAIM`.** `SELECT * FROM document_claim`. Tolak opsi lampiran bawaan (Google Storage / `DATA-WORKATTACH-FILE`). |
| **OQ-3** | Status tangga klaim (pyPosition/AcceptStatus/SendtoAdmin/SendtoMedical/Type) disimpan di mana? | **Tabel baru terpisah `T_WORK_CLAIM`** — dipakai bersama **semua klaim Life dan Non-Life** (lintas-lini). Keadaan tangga BUKAN di `T_GENERAL_CLAIM`. |
| **OQ-4** | Hardcode existing `ACCEPTATION_DATE=SYSDATE` & `STS_REJECT='0'` — ditiru atau nilai sebenarnya? | **Nilai sebenarnya.** `STS_REJECT` mengikuti aksi: Admin insert OS → `0`, Admin reject langsung → `2`, SPV tambah OS → `0`. `ACCEPTATION_DATE` diisi tanggal akseptasi asli saat baris benar-benar diaksep, **bukan** distempel SYSDATE tiap insert. |
| **OQ-5A** | `PEGA_JSON_KLAIM_PNC` (output BRANCH_NAME/CODE + COMMIT) — dibuang atau lookup? | **Tidak baca JSON — dibuang.** Sistem baru tidak memanggil procedure JSON ini. Bila butuh nama/kode cabang, ambil relasional biasa. |
| **OQ-5B** | `GetJsonProductLife` baca `m_product_life.JSONDATA` (lintas konteks) | **Setuju: baca dari tabel relasional `product_life`** (hasil migrasi Master Product Name Life), bukan JSONDATA. |

---

## Skema tabel klaim BARU — FINAL (6 tabel, 3 tingkat) + `T_WORK_CLAIM`

> ⚠️ **RALAT 2026-09-18 — DUA TABEL DIHAPUS, SKEMA JADI ENAM TABEL.** Judul di atas
> ("6 tabel, 3 tingkat"), catatan "**8 tabel klaim**" di bawah pohon, dan pohon itu sendiri
> **SUDAH TIDAK BERLAKU**. Teks lama dibiarkan utuh sebagai jejak keputusan; yang mengikat
> adalah ralat ini. Pohon dan tabel relasi yang berlaku ada di **`spec.md` §2b RALAT D dan E** —
> tidak diulang di sini supaya tidak ada dua sumber yang bisa berselisih.

`[keputusan work owner]` **`T_CLAIM_POLICY` dan `T_CLAIM_MARKETING` DIHAPUS.** Bukan diganti nama —
tabelnya **memang tidak ada**. Skema klaim Life turun dari **8 tabel menjadi 6**: lima tabel klaim
(`T_GENERAL_CLAIM`, `T_CLAIMLF_PREMIUMLIST_DETAIL`, `T_CLAIMLF_ADJUSTMENT`,
`T_CLAIMLF_ADJUSTMENT_SPREADING`, `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO`) ditambah
**`DOCUMENT_CLAIM`**. Data polis dan marketing **dibaca dari tabel polis** (`T_PREMIUM_LIST` dkk,
modul **PremiumList Life**) lewat **penunjuk di `T_GENERAL_CLAIM`**.

⛔ Jangan membuat `T_CLAIMLF_POLICY` maupun `T_CLAIMLF_MARKETING`.

**Mengapa boleh dihapus** `[terverifikasi]`:

| Tabel | Bukti |
| --- | --- |
| `T_CLAIM_POLICY` | **27 dari 32 kolomnya tersedia di `T_PREMIUM_LIST`** (32 = 23 kolom dasar + 9 kolom tambahan audit butir 4; `ID`/`CLAIM_ID` tidak dihitung). Lima sisanya kehilangan rumah — daftar di bawah. |
| `T_CLAIM_MARKETING` | **Terbukti turunan, nol nilai asli milik klaim.** `Claim Life/Activity/SetMOClaim_Act.xml` mengisi `MarketingData` dari `PolicyDataLife` (`MOID` · `MarketingCode` · `MarketingName` · `TeamGroup` · `BranchCode` · `BranchName`) **atau** dari master marketing officer (`ASM-FW-GISFW-Int-marketingofficer`, lewat `Claim Life/ReportDefinition/BrowseMarketingOfficer_RD.xml`). |

⚠️ **Penyimpangan sadar — potret berubah menjadi baca hidup, dan ini MEMBALIK keputusan di berkas
ini.** `[keputusan work owner]` Butir "PolicyDataLife = **snapshot** … *disalin, bukan baca live*"
di bagian `T_CLAIM_POLICY` **dicabut**. Akibatnya nyata dan diterima: **polis yang berubah sesudah
klaim dibuat akan mengubah tampilan klaim lama.** Untuk polis ber-endorsement itu **pasti terjadi**,
karena endorsement membuat versi baru.

⚠️ `[terbuka]` **Lima kolom kehilangan rumah — JANGAN DITEBAK.** Ini 5 dari 32 kolom
`T_CLAIM_POLICY` yang **tidak** tersedia di `T_PREMIUM_LIST`:

| Kolom | Keadaan |
| --- | --- |
| `TEAM_GROUP` | Tidak ada di `T_PREMIUM_LIST`; **ada di master marketing officer** — diambil lewat `MO_ID` |
| `BUSINESS_ID` | Tidak ada di `T_PREMIUM_LIST`; dipakai sebagai **parameter** `Claim Life/RDBList/Generate_NoAccept_Life.xml`. **Sumbernya wajib ditetapkan** |
| `TANGGAL_RESPON` · `TANGGAL_REALISASI` · `TANGGAL_KONFIRMASI_BALIK` | **Tanggal proses klaim, bukan atribut polis.** Tidak ada di `T_PREMIUM_LIST`, dan **tidak lagi ada di klaim**. Rumahnya `T_GENERAL_CLAIM` atau `T_WORK_CLAIM` — **belum diputuskan** |

⚠️ **Kolom yang bertambah dan yang pindah** (rincian di `spec.md` §2b RALAT B):
`T_GENERAL_CLAIM` **+**`CASEID_POLICY`/`POLICY_NO`/`ENDORSMENT_NO` (ketiganya ⚠️ `[terbuka]`),
**−**`PL_NUMBER` (diganti nama jadi `POLICY_NO`), **−**`CREATE_OP`/`CREATE_OP_NAME`/`TGL_UPDATE`
(**pindah** ke `T_WORK_CLAIM`, yang juga **+**`LINI` ⚠️ `[terbuka]`).

⚠️ `[terbuka]` **`T_WORK_CLAIM` naik menjadi AKAR pohon**, dan relasi
`T_WORK_CLAIM` → `T_GENERAL_CLAIM` **tidak punya kolom di sisi mana pun** — **memblokir tiket 14**.
Lihat `spec.md` §2b RALAT C4. **Jangan mengarang kolomnya.**


```
T_GENERAL_CLAIM (PK ID)                              header — PremiumListSummary dilipat 1:1
  ├─ T_CLAIM_POLICY               1:1  FK CLAIM_ID → T_GENERAL_CLAIM.ID   ON DELETE CASCADE
  ├─ T_CLAIM_MARKETING            1:1  FK CLAIM_ID → T_GENERAL_CLAIM.ID   ON DELETE CASCADE
  └─ T_CLAIMLF_PREMIUMLIST_DETAIL  1:N  FK CLAIM_ID → T_GENERAL_CLAIM.ID   ON DELETE CASCADE
        ├─ T_CLAIMLF_ADJUSTMENT     1:N  FK PREMIUM_LIST_DETAIL_ID → T_CLAIMLF_PREMIUMLIST_DETAIL.ID  ON DELETE CASCADE  ⬅ INDUK = DETAIL
        │     └─ T_CLAIMLF_ADJUSTMENT_SPREADING       1:N  FK ADJUSTMENT_ID → T_CLAIMLF_ADJUSTMENT.ID  ON DELETE CASCADE  ⬅ TABEL BARU (temuan 2026-09-16)
        │            └─ T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO  1:N  FK SPREADING_ID  ON DELETE CASCADE  ⬅ TABEL BARU
        └─ DOCUMENT_CLAIM         1:N  FK PREMIUM_LIST_DETAIL_ID → T_CLAIMLF_PREMIUMLIST_DETAIL.ID  ON DELETE CASCADE  ⬅ TABEL BARU (gerbang simpan)

T_WORK_CLAIM (PK ID)                             ⬅ TABEL BARU, TERPISAH, LINTAS-LINI (Life + Non-Life)
                                                    status/posisi tangga klaim — bukan bagian ClaimData
```

> **8 tabel klaim + `T_WORK_CLAIM`.** Kedalaman kini **empat tingkat** di sisi peserta:
> peserta → adjustment → spreading → spreading_retro.

### T_GENERAL_CLAIM (header) — skalar ClaimData + PremiumListSummary (1:1) dilipat
`ID` (PK), `CLAIM_NO`, `PY_ID`, `STS_KATASTROFE`, `KATASTROFE_NOTE`, `IS_KPR`, `STNC_CLAIM`,
`ACCEPTED_NO`, `STS_REJECT`, `RISLIPRNM`, `PL_NUMBER` (dari summary), `BUSINESS_NAME` (dari summary),
`CASEID`, `CLAIM_RETRO`, `CREATE_OP`, `CREATE_OP_NAME`, `TGL_UPDATE` DATE.

> ⚠️ Tiga tanggal per-peserta (CONFIRMATION/COMPLETE/CLAIM_RECEIVED) **TIDAK** di header (OQ-1).
> `ACCEPTED_NO`/`STS_REJECT` di header = **cerminan** baris adjustment terakhir (ADR-0011).
> `CREATE_OP`/`CREATE_OP_NAME`/`TGL_UPDATE` + `CASEID` = dari work object & INSERT existing, bukan
> property ClaimData murni — aturan "kolom = property ClaimData" diperlebar agar jujur (audit butir 9).

> ⚠️ **RALAT 2026-09-18 — TABEL INI DIHAPUS.** Seluruh bagian di bawah tidak berlaku, termasuk
> catatan **snapshot** *"disalin, bukan baca live"* yang kini **DIBALIK** menjadi **baca hidup**
> ⚠️ **penyimpangan sadar**. **27 dari 32** kolomnya ada di `T_PREMIUM_LIST`; **lima** kehilangan
> rumah ⚠️ `[terbuka]` — ketiga `TANGGAL_*`, `TEAM_GROUP`, dan `BUSINESS_ID`.

### T_CLAIM_POLICY (PolicyDataLife, 1:1)
`ID` (PK), `CLAIM_ID` (FK), `BUSINESS_CODE`, `BUSINESS_ID`, `BUSINESS_NAME`, `CEDING_CO`,
`CEDING_CO_NAME`, `DESCRIPTION`, `JENIS_ASURANSI`, `MARKETING_CODE`, `MARKETING_NAME`, `MO_ID`,
`NO_OFFER`, `POLICY_HOLDER`, `POLICY_HOLDER_NAME`, `PRO_RATE_TYPE`, `RETRO_ID`, `RETRO_NAME`,
`SECURITY_REINSURER`, `SECURITY_REINSURER_ID`, `SOB_NAME`, `SOURCE_OF_BUSINESS`, `TYPE`,
`TYPE_CEDING`, `WPC`,
**+8 (audit butir 4):** `TANGGAL_RESPON` DATE, `TANGGAL_REALISASI` DATE, `TANGGAL_KONFIRMASI_BALIK`
DATE, `PRODUCT_NAME_ID`, `PRODUCT_NAME`, `TEAM_GROUP`, `DATE_RECEIVED` DATE, `BRANCH_NAME`,
`BRANCH_CODE`.
**−`PL_NUMBER`** (pindah ke header — milik summary, bukan polis).

> `[keputusan work owner]` PolicyDataLife = **snapshot** dari data polis life / PremiumList NB/EDM
> saat klaim dibuat (disalin, bukan baca live). `PRODUCT_NAME`/`PRODUCT_NAME_ID` diisi dari
> **`product_life` relasional** (OQ-5B), bukan JSON.

> ⚠️ **RALAT 2026-09-18 — TABEL INI DIHAPUS.** Seluruh kolom di bawah **terbukti turunan**
> (`SetMOClaim_Act.xml`): nol nilai asli milik klaim. `TEAM_GROUP` termasuk lima kolom yang
> ⚠️ `[terbuka]` kehilangan rumah.

### T_CLAIM_MARKETING (MarketingData, 1:1)
`ID` (PK), `CLAIM_ID` (FK), `MARKETING_ID`, `CLIENT_ID`, `CLIENT_NAME`, `BRANCH_DETAIL_ID`,
`BRANCH_DETAIL_NAME`, `TEAM_GROUP`.

### T_CLAIMLF_PREMIUMLIST_DETAIL (1:N, per peserta diklaim)
`ID` (PK), `CLAIM_ID` (FK), `PL_NUMBER`, `POLICY_NO`, `POLICY_HOLDER`, `CERTIFICATE_NO`,
`NAME_OF_INSURED`, `DOB` DATE, `AGE`, `SEX`, `PLAN`, `DISEASE`, `ICD_CODE`, `DESCRIPTION`, `NOTES`,
`KETERANGAN`, `DATE_OF_LOSS` DATE, `RECEIVED_DATE` DATE, `BEGIN_DATE` DATE, `EFFECTIVE_DATE` DATE,
`EXPIRED_DATE` DATE, `LAPSE_DATE` DATE, `GROSS_VALUATION_BEGIN_DATE` DATE,
`GROSS_VALUATION_EXPIRED_DATE` DATE, `RETROCESSION_VALUATION_BEGIN_DATE` DATE,
`RETROCESSION_VALUATION_EXPIRED_DATE` DATE, `CURRENCY`, `SUM_INSURED` NUMBER, `SUM_REASURED` NUMBER,
`GROSS_PREMIUM` NUMBER, `NET_PREMIUM` NUMBER, `CLAIM_AMOUNT` NUMBER, `EM_PERCENT` NUMBER,
`SHARE_NUSANTARA_RE` NUMBER, `SHARE_RETRO` NUMBER, `RETROCEDED_SHARE` NUMBER, `WPC`, `STNC_TREATY`,
**+9 (audit butir 2):** `IS_CHECK` (penanda peserta dipilih diklaim), `STATUS`, `RECOMMENDATION`,
`STS_REJECT`, `SOURCE_ID`, `CONFIRMATION_DATE` DATE, `COMPLETE_DATE` DATE, `CLAIM_RECEIVED_DATE` DATE,
`CEDING_RETENTION` NUMBER.

> Helper `idx`/`IndexPremium` = variabel loop, BUKAN kolom.
> ⚠️ `STS_REJECT` tingkat detail = **cerminan** baris adjustment terakhir (ADR-0011).
> ⚠️ Jebakan `EDMSTATUS` (peserta hidup) ada di sisi `M_LIFE_PREMIUM_DETAIL` (PremiumList), tak rusak
> oleh pindah tabel; justru karena itu `IS_CHECK` wajib ada (snapshot merekam siapa yang dipilih).

### T_CLAIMLF_ADJUSTMENT (1:N) — INDUK = T_CLAIMLF_PREMIUMLIST_DETAIL
`ID` (PK), **`PREMIUM_LIST_DETAIL_ID` (FK → T_CLAIMLF_PREMIUMLIST_DETAIL.ID)** ⬅ bukan CLAIM_ID,
`SHARE_NUSANTARA_RE` NUMBER, `CEDING_RETENTION` NUMBER, `SUM_REASURED` NUMBER, `SUM_INSURED` NUMBER,
`SHARE_RETRO` NUMBER, `RETROCEDED_SHARE` NUMBER, `CURRENCY_ID`, `CURRENCY`,
**+4 (audit butir 3):** `CLAIM_AMOUNT` NUMBER, `STS_REJECT`, `ACCEPTEDNO`, `ACCEPTATION_DATE` DATE,
**+3 (data bank — terlewat audit, ditemukan verifikasi Kiro 2026-09-16):** `NAME_OF_BANK` (←
`NameOfBank`), `ID_BANK` (← `IDOfBank`), `ACCOUNT_NO` (← `NoAccount`).
**+1 (referensi komite — 2026-09-16):** `KOMITE_ID` (referensi ke kasus/roster komite di konteks
Komite Life; NULL bila belum dikirim). Roster+keputusan komite TIDAK disimpan di sini — di-view lewat
join ke tabel Komite Life.

> `STS_REJECT` di sini = **unit keputusan** (ADR-0011): `0`/`1`/`2`. `ACCEPTATION_DATE` diisi nilai
> sebenarnya saat aksep, bukan SYSDATE hardcode (OQ-4).
> ⚠️ **Tiga field bank** milik class `ASM-FW-GISFW-Data-AdjustmentLife` (baris adjustment, bukan
> policy/header). `[terverifikasi]` sumber: `Claim Life/Section/AdjustmentDetail_Section.xml`
> menampilkan `.NameOfBank`/`.ACCOUNTNO`; kolom fisik `NAME_OF_BANK`/`IDBANK`/`ACCOUNTNO` di
> `UpdateOsAkseptasiClaimLife_sql` (nama "Update" menyesatkan — isinya **INSERT** 55 kolom).
> ⚠️ **Gerbang wajib-isi terverifikasi:** `Claim Life/Activity/GetListKomiteLife.xml`
> (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `GETLISTKOMITELIFE`) precondition
> `.NameOfBank==""||.NoAccount==""||.IDOfBank==""`, pesan `"Name of bank cannot be empty"` (memo rule
> "add proteksi name of bank") → ketiga field **wajib terisi sebelum penyerahan ke Komite**. AC di
> tiket penyerahan Komite (10) + tiket adjustment (03).
> `[terverifikasi + keputusan work owner 2026-09-16]` **KAPAN wajib = PARITAS.** Bank **boleh kosong
> saat Save Outstanding**, wajib **hanya sebelum penyerahan Komite** — `SaveOutStandingLife_Act.xml`
> **nol** gerbang bank (dicek); gerbang hanya di `GetListKomiteLife`. BUKAN penyimpangan, tak
> ditandai ⚠️.

### T_CLAIMLF_ADJUSTMENT_SPREADING (1:N) — TABEL BARU, INDUK = T_CLAIMLF_ADJUSTMENT `[temuan 2026-09-16]`
`ID` (PK), **`ADJUSTMENT_ID` (FK → T_CLAIMLF_ADJUSTMENT.ID)**, + kolom per treaty-year spreading:
`TREATY_TYPE_ID`, `TREATY_TYPE_NAME`, `TREATY_YEAR_LIFE`, `RATE` NUMBER, dan turunan (mengikuti pola
`T_PREMIUM_LIST_SPREADING`; kolom final dari sensus `SpreadingClaimLife_Act`).

> ⚠️ **Terlewat di audit awal — ditemukan verifikasi Kiro 2026-09-16.** `AdjustmentList` (class
> `ASM-FW-GISFW-Data-AdjustmentLife`) punya anak `SpreadingList`. `[terverifikasi]`
> `Claim Life/Section/AdjustmentDetail_Section.xml` grid `.SpreadingList` (class AdjustmentLife);
> `Claim Life/Activity/SpreadingClaimLife_Act.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` /
> `SPREADINGCLAIMLIFE_ACT`) MENGISI `.SpreadingList` + `.RetroLifeList` dengan perhitungan.
> `[keputusan work owner]` **Disimpan** (dibekukan, auditable) — konsisten dgn spreading
> PremiumList/Endorsement.

### T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO (1:N) — TABEL BARU, INDUK = T_CLAIMLF_ADJUSTMENT_SPREADING
`ID` (PK), **`SPREADING_ID` (FK → T_CLAIMLF_ADJUSTMENT_SPREADING.ID)**, `REINSURER_NAME`
(← `.REINSURERNAME`), `PERCENT_SHARE` NUMBER (← `.PERCENTSHARE`), `AMOUNT` NUMBER (← `.Amount`),
`RATE` NUMBER, `PREMIUM_SPREADED_GROSS` NUMBER, `PREMIUM_SPREADED_NET` NUMBER, `TREATY_TYPE_ID`,
`TREATY_TYPE_NAME` (kolom final dari sensus).

> `[terverifikasi]` `Claim Life/Section/RetroDetailClaimLife.xml` tampil `.REINSURERNAME`,
> `.PERCENTSHARE`, `.Amount` (class `ASM-FW-GISFW-Int-RETROCESSIONLIFE`); `SpreadingClaimLife_Act`
> set `.Amount = AmountRetroShare * PERCENTSHARE/100`, `.PREMIUM_SPREADED_GROSS = Rate*(1+EMPercent)*Amount`.
> Sama class & pola dengan `T_PREMIUM_LIST_SPREADING_RETRO`.

### DOCUMENT_CLAIM (1:N) — TABEL BARU, per peserta
`ID` (PK), **`PREMIUM_LIST_DETAIL_ID` (FK → T_CLAIMLF_PREMIUMLIST_DETAIL.ID)**, + kolom dokumen
(nama/kategori/kelengkapan) — rincian menyusul dari sensus `.DocumentList` di `SaveOutStandingLife_Act`.

> Gerbang simpan Outstanding: "The document hasn't been uploaded person number …" +
> "Documents are incomplete, please complete the documents". Dokumen **per peserta**, jadi gerbang.

### T_WORK_CLAIM (PK ID) — TABEL BARU, TERPISAH, LINTAS-LINI
Menyimpan **keadaan tangga/posisi** klaim: `pyPosition`, `AcceptStatus`, `SendtoAdmin`,
`SendtoMedical`, `Type`, dst. Dipakai bersama **Life dan Non-Life** — bukan bagian `ClaimData`, jadi
tidak dilipat ke `T_GENERAL_CLAIM`. Bentuk kolom final menyusul saat penyelarasan lintas-lini.

> ⚠️ Karena lintas-lini, `T_WORK_CLAIM` **bukan** anak `T_GENERAL_CLAIM` dengan cascade — ia tabel work
> mandiri yang mereferensikan klaim (Life/Non-Life) via kunci klaim + lini. Relasi & cascade-nya
> ditetapkan saat konteks Non-Life digarap; untuk Claim Life ia menampung posisi klaim life.

---

## Yang DIBACA (bukan tabel milik klaim)
- Data polis life / PremiumList NB/EDM → sumber snapshot `T_CLAIM_POLICY` + `T_CLAIMLF_PREMIUMLIST_DETAIL`.
  ⚠️ **RALAT 2026-09-18** — baris ini DICABUT: tidak ada lagi `T_CLAIM_POLICY`, dan tidak ada
  lagi snapshot. Data polis **dibaca hidup** dari `T_PREMIUM_LIST` dkk lewat penunjuk di
  `T_GENERAL_CLAIM` ⚠️ **penyimpangan sadar**.
- **`product_life` relasional** (OQ-5B) → `PRODUCT_NAME`/`PRODUCT_NAME_ID` — bukan `JSONDATA`.
- Master (retro, security reinsurer, currency, dll) → dibaca saat isi klaim.

## OQ tersisa — KOSONG untuk penyimpanan Claim Life
Seluruh 5 OQ audit ditutup 2026-09-16. `PEGA_JSON_KLAIM_PNC` tidak dipakai (tidak baca JSON).
Satu-satunya sisa: `[terbuka — Non-Life]` relasi & cascade formal `T_WORK_CLAIM` — **tidak
memblokir** Claim Life (untuk Life cukup: hapus klaim → baris work ikut, AC di tiket 15).

## REVISI LANJUTAN (2026-09-16) — spreading adjustment + referensi Komite `[keputusan work owner]`

**Temuan Kiro (diverifikasi korpus):** `AdjustmentList` punya anak `SpreadingList` → `RetroLifeList`
yang **terlewat** di skema 6-tabel. Skema klaim jadi **8 tabel**, kedalaman 4 tingkat di sisi peserta:
- `T_CLAIMLF_ADJUSTMENT_SPREADING` (per treaty-year) + `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` (per reinsurer)
- **Disimpan** (dibekukan, auditable) — konsisten dgn PremiumList/Endorsement. ⚠️ penyimpangan sadar.

**Komite → `T_CLAIMLF_ADJUSTMENT` simpan `KOMITE_ID` saja `[keputusan work owner]`:**
`T_CLAIMLF_ADJUSTMENT` menyimpan **`KOMITE_ID`** (referensi ke kasus/roster komite), NULL bila belum
dikirim. **Roster + keputusan komite (KomiteAproval/KomiteComment/DateApprove per anggota) BUKAN di
Claim Life** — milik **konteks Komite Life (tabel terpisah)**. Adjustment **view** keputusan lewat
join `KOMITE_ID` → tabel komite. TIDAK ada `T_CLAIMLF_ADJUSTMENT_KOMITE`.

> `[terverifikasi]` Korpus menaruh `KomiteList` sebagai anak `AdjustmentList`
> (`KomitePostAdjustment` tulis `…AdjustmentList(idx).KomiteList(k).KomiteAproval/KomiteComment/DateApprove`)
> + `CreateKMTLife_Act` set `childPageKomite.KomiteList`/`IndexPremiumList`/`IndexAdjustment`
> (subscript posisi). **Sistem baru menyimpang sadar:** roster+keputusan dipindah ke tabel Komite Life,
> adjustment cuma pegang `KOMITE_ID`; referensi pakai **ID stabil**, bukan index posisi Pega. ⚠️
> ⚠️ **Lintas konteks:** tabel komite yang menampung roster+keputusan ini ada di **Komite Claim Life**
> (Task 1). Perlu diselaraskan — apakah tabel komite di sana sudah punya kolom `ADJUSTMENT_ID`/keputusan
> per baris, atau perlu ditambah. `[terbuka — selaraskan Komite Claim Life]`.

## STATUS FINAL (2026-09-16)
Revisi penyimpanan Claim Life terverifikasi korpus. Artefak yang berlaku:
- `spec.md` §2b + AC 31–57, US 1–60 (bagian penyimpanan direvisi; perilaku/peran/status TIDAK diubah)
- Tiket: **14** (PREFACTOR skema+migrasi), **15** (hapus klaim+popup+baris work) — keduanya baru;
  diselaraskan: 02, 03, 04, 05, 10, 12, 13. Tak disentuh: 01, 06, 07, 08, 09, 11.
- Draf lama `spec-penyimpanan-relasional.md` = **superseded** (ditandai, tidak dihapus).
- Penyimpangan sadar (⚠️); gerbang bank = paritas (bukan ⚠️).
- ⚠️ **PERLU revisi lanjutan spec + tiket 14** (skema):
  1. Tambah `T_CLAIMLF_ADJUSTMENT_SPREADING` + `_SPREADING_RETRO` (jadi **8 tabel**). Kolom final dari
     sensus `SpreadingClaimLife_Act`/`RetroDetailClaimLife`.
  2. Tambah kolom `KOMITE_ID` di `T_CLAIMLF_ADJUSTMENT` (referensi; roster+keputusan komite di tabel
     Komite Life terpisah, di-view lewat join). TIDAK ada `T_CLAIMLF_ADJUSTMENT_KOMITE`.
  3. Selaraskan **Komite Claim Life** (Task 1) `[keputusan work owner: setuju]`: tabel komite di sana
     menyimpan keputusan per baris dengan **`ADJUSTMENT_ID` (FK balik → `T_CLAIMLF_ADJUSTMENT.ID`)** +
     `KOMITE_APROVAL`/`KOMITE_COMMENT`/`DATE_APPROVE`/roster. Adjustment view lewat `KOMITE_ID` /
     join `ADJUSTMENT_ID`. Referensi ID stabil (bukan index posisi Pega `IndexPremiumList`/`IndexAdjustment`).

## Dampak ke artefak (langkah selanjutnya)
- Revisi **spec Claim Life** bagian penyimpanan (§ persistence): JSON dibuang, ClaimData → **6 tabel
  baru + `T_WORK_CLAIM`**, existing tak dipakai untuk simpan.
- **Slice skema/migrasi = PREFACTOR (tiket pertama)** karena bentuk skema berubah total.
- Tiket Claim Life yang menyentuh simpan (register, outstanding, akseptasi, migrasi) diselaraskan —
  revisi terarah, bukan regenerasi.
- Penyimpangan sadar (⚠️): buang JSON, adjustment induk = detail (perbaikan relasi Pega), `IS_CHECK`
  eksplisit, dokumen jadi tabel, status tangga jadi `T_WORK_CLAIM` lintas-lini, nilai STS_REJECT/
  ACCEPTATION_DATE sebenarnya (bukan hardcode).
