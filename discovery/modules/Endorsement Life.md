# Modul — Endorsement Life

STEP D3. Sintesis dari `../inventory/Endorsement Life.md` (D1) + `../flows/Endorsement Life.md`
(D2 Tahap 1). **75 file — modul ber-`Flow` terkecil di korpus**.
Audit: `find "Endorsement Life" -name '*.xml' | wc -l`

> Ejaan nama modul (`Endorsement`, dengan `e`) berbeda dari `Endorsment Fac In` — **ditulis apa
> adanya** sesuai korpus.

## 1. Peran modul

`[terverifikasi]` Menangani **endorsement (EDM) untuk lini life**. Class kerja
`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` (38 rule). Titik masuk
`Endorsement Life/Flow/InputEDMLife.xml` → `RULE-OBJ-FLOW` /
`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE`, `<pyStartActivity>Start1`, 68.374 byte.

## 2. Proses / fitur utama

Rincian: **`../flows/Endorsement Life.md`** (14 rule ditelusur).

`[terverifikasi]` **Flow paling sederhana di korpus**: 2 Assignment, 1 Decision, 1 Utility,
4 connector.

```
Start1 ─(Always)─> ASSIGNMENT63 "InputEDMLife"  [WorkList]
                        │ FlowAction: InputEDMLife
                        v
                    Decision1 "Accept"  (DecisionTable: IsLifeAccepted)
                        ├─ Status=Decline ─> End1
                        └─ Status=Confirm ─> Utility1  (Activity InsertJsonPolisLife_Act)
                                                 └─> END52
```

`[terverifikasi]` **Kedua Assignment memakai `WorkList`**, bukan `WorkBasket` — kontras dengan
NB Treaty In dan EDM Treaty In yang seluruhnya `WorkBasket` → **OQ-028**.

`[pertanyaan terbuka]` `Assignment1` "Input EDM Summary" **tidak dituju connector mana pun** —
hanya punya connector keluar. Ini pola kedua setelah NB Treaty In → **OQ-023**.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 17 Connect-SQL, 16 Activity, 14 Section, 8 Harness,
6 ReportDefinition, 6 FlowAction, 2 When, 2 DataTransform, 1 Flow, 1 DecisionTable, 1 ConnectREST,
1 SystemSettings.

| Objek | Rule perujuk |
| --- | ---: |
| `POOLDATA.JSON_POLIS` | 6 |
| `POOLDATA.M_TEMPUPLOADLIFE` | 2 |
| `V_TEMPRETENSICEDING`, `V_MIN`, `V_COUNT`, `VCOUNT`, `RICOMM_LIFE`, `RATE_LIFE`, `PRODUCT_LIFE`, `POOLDATA.M_LIFE_PREMIUM_DETAIL` | 1 masing-masing |

`[terverifikasi]` Seluruh proses berpusat pada `Activity/InsertJsonPolisLife_Act.xml`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT`, **22 langkah**).

### 3.1 Kode yang menggerbangi logika

`[terverifikasi]` Dari precondition `InsertJsonPolisLife_Act`:

| Properti / kondisi | Nilai literal | Arti |
| --- | --- | --- |
| `Local.currentdate` | diuji `@toDecimal(...) > 25` | **ambang tanggal 25** — sama dengan default `TglProd=25` di EDM Treaty In → **OQ-030** |
| `PremiumListSummary.PL_NUMBER_EDM` | `== ""` | nomor premium list kosong |
| `pyWorkPage.Type` | **`"QR"`, `"QP"`** (step 12) | **arti belum terverifikasi** → OQ-020 |
| `pyWorkPage.Type` | **`"TR"`, `"TP"`** (step 13) | **arti belum terverifikasi** → OQ-020 |

`[terverifikasi]` Kode `Type` yang sama (`QP`/`QR` vs `TP`/`TR`) dipakai
`Komite Claim Life/…/Generate_NoAccept_KMT_Life*` untuk memilih rule penomoran — **satu enumerasi,
dua pemakai**.

`[terverifikasi]` Modul ini memberi **bukti konteks untuk `EdmType`** → OQ-014
(`../flows/Endorsement Life.md`).

## 4. Integrasi eksternal

`[terverifikasi]` Satu `RULE-CONNECT-REST`: `convertJsonNusareToProduction`
(`SETTING` → `LinkService!LinkService`, tanpa URL literal).

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **PremiumList Life** | `PremiumListSummary.PL_NUMBER_EDM`; `POOLDATA.M_LIFE_PREMIUM_DETAIL`; keduanya memakai `DecisionTable/IsLifeAccepted.xml` | `[terverifikasi]` |
| **Master Product Name Life** | membaca `PRODUCT_LIFE`, `RATE_LIFE` | `[terverifikasi]` |
| **Komite Claim Life** | berbagi enumerasi `Type` (`QP`/`QR`/`TP`/`TR`) dan class `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `[terverifikasi]` |
| **EDM Treaty In** | ambang tanggal 25 muncul di kedua modul | `[terverifikasi]` kesamaan nilai; `[dugaan]` hari tutup buku — OQ-030 |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **Dua stored procedure**, isinya tidak ada di korpus (OQ-002):
`POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY`, `DBMS_LOB.CREATETEMPORARY`.

`[terverifikasi]` `DecisionTable/IsLifeAccepted.xml` menggerbangi satu-satunya Decision di flow ini
— **baris keputusannya tidak ikut terekspor**, seperti seluruh 49 `DecisionTable` korpus
→ **OQ-043**. Ia juga ada sebagai dua tipe rule → OQ-026.

`[terverifikasi]` `POOLDATA.JSON_POLIS` (6 rule) — struktur JSON tidak ada di korpus → OQ-012.
Telusur Tahap 1 membuat **49 path JSON sebagian terbaca** dari modul ini — kemajuan parsial
atas OQ-012.

`[terverifikasi]` `IsPEGAPROD` ada; kondisinya tidak terbaca → OQ-029.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **0 file**,
`OperatorID.pyTelephone` **0 file** — modul tanpa jejak guard identitas.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-012**, **OQ-014**, **OQ-020**, **OQ-023**, **OQ-026**, **OQ-028**, **OQ-029**,
**OQ-030**, **OQ-043**.
