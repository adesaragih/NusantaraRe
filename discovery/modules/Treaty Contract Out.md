# Modul — Treaty Contract Out

STEP D3. Sintesis dari `../inventory/Treaty Contract Out.md` (D1) +
`../flows/Treaty Contract Out.md` (D2 Tahap 5) + `../flows/_SUMMARY-treaty-master.md`.
**303 file**. Modul **tanpa rule `Flow`** (OQ-005).
Audit: `find "Treaty Contract Out" -name '*.xml' | wc -l`

> **Nama modul menyesatkan.** §2 menyajikan lima uji berbukti; jawabannya: modul ini **bukan**
> tentang treaty outward.

## 1. Peran modul

`[terverifikasi]` **Editor master *term / arrangement* kontrak treaty** — bukan modul treaty
outward. Ia mengelola `TREATYCONTRACT`, `M_TREATYYEAR`, `TREATYREINSURER`, `TREATYBUSINESS`,
`M_PROPORTIONALARRG` (+ child), `MTREATYSECURITY`, `TREATYEXCHANGE`, dengan **16 jenis klausul**
yang dapat diedit.

`[terverifikasi]` **Titik masuk**: `Harness/InboxTreatyContract.xml` → `RULE-HTML-HARNESS` /
`DATA-PORTAL` / `INBOXTREATYCONTRACT`, 128.846 byte. Dua harness lain:
`InboxTreatyContractDescription.xml` (502.598 B, meng-`<pySection>` `GridTreatyArrangementPortfolioList`,
`GridTreatyArrangementProfitCommision`, `NitipKurs`, `ViewDetailDescription`) dan
`InboxTreatyContractReinsType.xml` (85.622 B).

`[terverifikasi]` **Bentuk modul: editor master, bukan proses berjenjang.** Tidak ada
`Akseptasi_DT`, tidak ada `StatusAkseptasi` (0 file), tidak ada tombol Submit/Akseptasi/Decline.
Audit: `grep -rl "StatusAkseptasi" "Treaty Contract Out" --include="*.xml" | wc -l` → **0**.

## 2. UJI NAMA MODUL — outward atau inward?

Rincian lengkap: **`../flows/Treaty Contract Out.md`** §2 (26 rule ditelusur).

`[terverifikasi]` Lima uji terpisah, seluruhnya dapat diaudit ulang:

| Uji | Hasil |
| --- | --- |
| **Objek database** | **0** objek bernama `*_OUT*`. Yang ada: `M_PROPORTIONALARRG` (10 rule), `MTREATYSECURITY` (5), `TREATYREINSURER` (3), `TREATYBUSINESS` (3), `TREATYEXCHANGE`, `TREATYCONTRACT`, `PROPORTIONALARRG`, `M_TREATYYEAR` |
| **Class rule** | **0** class bernama `*OUT*`; justru **5 rule berclass `ASM-FW-GISFW-INT-TREATY_IN`** (`DELETE_ACT`, `DELETEATTACHMENT2_SQL`, `GETALLATTACHMENT2_SQL`, `GETATTACHMENT2_SQL`, `INSERTATATCHMENT_SQL`), dan **11 file** menyebut class itu |
| **Penamaan rule** | 81 activity `*TreatyArr*` vs **4** `*Out*` — keempatnya soal **lampiran** (`LoadAttachmentTreatyOut`, `TreatyOutDownloadAll_Act`, `TreatyOutDownloadOne`, `TreatyOutSaveAttachment`) |
| **Arah tulis** | **9 Connect-SQL MENULIS master** lewat `PEGA_TREATYCONTRACT`, `PEGA_TREATYYEAR`, `PEGA_TREATYREINSURER`, `PEGA_TREATYBUSINESS`, `PEGA_PROPORTIONALARRG`, `PEGA_M_PROPORTIONALARRG_CHILD`, `PROSESCOPY`, `PEGA_M_ATTACHMENT`, `GET_TOKEN_STORAGE` |
| **Pembanding korpus** | Modul yang benar-benar menyentuh objek treaty outward: `NB Treaty In` **8** file, `Claim Non Prop` **4**, `EDM Treaty In` **4**, `Treaty In Adjustment` 2, `Claim Prop` 1, **`Treaty Contract Out` 0** |

Audit:
```
for pat in M_TREATY_OUT TREATY_OUT TREATYOUTDETAIL FACOUTPRODUCTION; do
  echo "$pat: $(grep -rli "$pat" "Treaty Contract Out" --include='*.xml' | wc -l)"
done                                               # semuanya 0
awk -F'\t' '$1 ~ "^Treaty Contract Out/"{print $3}' all-rules.tsv | grep -i out   # kosong
```

`[terverifikasi]` Satu-satunya class treaty-outward di **seluruh korpus** adalah
`ASM-FW-GISFW-INT-TREATYOUTDETAIL` (5 rule), letaknya di `NB Treaty In` (3),
`Treaty In Adjustment` (1), `EDM Treaty In` (1) — **tidak satu pun di modul ini**.

**Jawaban `[terverifikasi]`: BUKAN outward.** → **OQ-022 terjawab sebagian**; yang tersisa adalah
*mengapa* dinamai "Out" dan ke bounded context mana 303 rule ini ditempatkan (**D4**).

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 168 Activity, 49 Section, 37 Connect-SQL,
30 ReportDefinition, 10 DataTransform, 3 Harness, **2 FlowAction** (paling sedikit di korpus untuk
modul sebesar ini), 2 DataPage, 1 ConnectREST, 1 SystemSettings. **Nol rule `Flow`, nol `When`.**

`[terverifikasi]` Pola rule simetris per jenis klausul: 27 `SaveTreatyArr*_Act`,
`New…`/`Set…`/`Browse…`, dan **16 `CancelActivity*`** yang menamai jenis klausul:
`BordereAux`, `CashLossLimit`, `ClaimCoorperation`, `Epi`, `ExGratia`, `ExGratiaLimitChild`,
`FacIn`, `FacInList`, `PLA`, `ProfitCommision`, `Ricomm`, `TreatyContract`, `TreatyLimit`,
`TreatyLimitChild`, `pPortfolio`, `pTerrLimit`.
Total activity bernama `*TreatyArr*`: **81**.
`EPI`, `PLA`, `Ricomm`, `LOL`, `MB` **kepanjangan belum terverifikasi**.

`[terverifikasi]` `ReinsTypeID` diuji `== ""` **114×**; diisi dari `InputData.CARIREINS…`.
`Harness/InboxTreatyContractReinsType.xml` adalah layar tersendiri → "jenis reasuransi" adalah
**dimensi master tersendiri**. **Arti belum terverifikasi** → OQ-020.

## 4. Integrasi eksternal

`[terverifikasi]` Satu `RULE-CONNECT-REST`: `ServiceGoogle` (`SETTING` → `LinkService!LinkService`,
tanpa URL literal). Penyimpanan berkas: `T_STORAGE_IMAGE` (2 rule),
`POOLDATA.GET_TOKEN_STORAGE`, `POOLDATA.PEGA_M_ATTACHMENT`, dan 7 rule berclass
`ASM-FW-GISFW-INT-T_STORAGE_IMAGE`.

**Tidak ada Arasapas / Kasir / Konversi / Gemini AI.**

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **Domain treaty inward** | 5 rule berclass `ASM-FW-GISFW-INT-TREATY_IN`; 11 file menyebut class itu; 12 rule bernama FacIn/TreatyIn (`BrowseDeleteRowTreatyInContract`, `SaveTreatyArrFacIn_Act`, `SetCategoryAttachTreatyin`, …) | `[terverifikasi]` |
| **Claim Prop / Komite Claim Prop** | modul-modul itu **membaca** `TREATYREINSURER`, `TREATYBUSINESS`, `PROPORTIONALARRG` yang **ditulis** di sini | `[terverifikasi]` arah tulis/baca |
| **Claim Fac In** | membaca `TREATYCONTRACT`, `TREATYBUSINESS`, `PROPORTIONALARRG` | `[terverifikasi]` |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **11 stored procedure — terbanyak dari kelima modul Tahap 5**, isinya tidak ada
di korpus (OQ-002). Karena **seluruh jalur tulis master melewati procedure**, aturan penyimpanan
master arrangement treaty — validasi, versi, kunci, kaskade ke child — **seluruhnya berada di sisi
database**.

`[terverifikasi]` **`DBMS_OUTPUT.PUT_LINE` muncul di modul ini — satu-satunya kemunculannya di
korpus** (`awk -F'\t' '$8 ~ /DBMS_OUTPUT/{print $1}' all-rules.tsv`). `POOLDATA.PROSESCOPY` juga
hanya di sini.

`[terverifikasi]` **202 dari 303 rule (66,7 %) berada di `@BASECLASS`** — proporsi tertinggi di
korpus. Akibatnya **pemetaan rule → entitas domain tidak dapat diturunkan dari class** untuk dua
pertiga modul → **OQ-009**.

Audit:
```
awk -F'\t' '$1 ~ "^Treaty Contract Out/"{print $3}' all-rules.tsv | sort | uniq -c | sort -rn | head
```

`[terverifikasi]` **Tidak ada objek db-link** — OQ-017 tidak berlaku.
`IsPEGAPROD` **tidak ada** — OQ-029 tidak berlaku di modul ini.

`[terverifikasi]` Folder **`Claude outputs`** berisi `Struktur_InboxTreatyContract.xlsx`
(306.979 B) — berkas non-Pega di dalam korpus READ-ONLY; **tidak dibaca** → **OQ-054**.

`[terverifikasi]` **Yang tidak dapat direkonstruksi** (`../flows/_METHOD-noflow.md` §3.5): urutan
wajib antar layar, mekanisme assignment, SLA, titik akhir proses, dan siapa yang boleh mengedit
klausul — bahkan ekspresi visibilitas ber-workbasket pun **tidak ditemukan** di modul ini,
berbeda dari `Treaty In`.

`[terverifikasi]` 168 Activity, 49 Section, 30 ReportDefinition **belum ditelusur isinya** — batas
**cakupan telusur**.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **2 file**,
`OperatorID.pyTelephone` **0**.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-005**, **OQ-009**, **OQ-016**, **OQ-018**, **OQ-020**, **OQ-021**,
**OQ-022** (terjawab sebagian), **OQ-042**, **OQ-047**, **OQ-054**.
