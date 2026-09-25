# Modul — PremiumList Life

STEP D3. Sintesis dari `../inventory/PremiumList Life.md` (D1) + `../flows/PremiumList Life.md`
(D2 Tahap 1). **124 file**.
Audit: `find "PremiumList Life" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` Menangani **penawaran dan daftar premi (premium list) untuk lini life** —
dari input offer hingga detail premi. Class kerja `ASM-FW-GISFW-WORK-LIFE` (64 rule).

`[terverifikasi]` **Titik masuk**: `PremiumList Life/InputPolicyHolder.xml` → `RULE-OBJ-FLOW` /
`ASM-FW-GISFW-WORK-LIFE` / `INPUTPOLICYHOLDER`, `<pyStartActivity>Start1`.

`[terverifikasi]` **Letak file tidak lazim**: rule `Flow` ini berada **langsung di root modul**,
bukan di folder `Flow/`. Tipenya diverifikasi dari `<pzOriginalInstanceKey>` **dan** `<pxObjClass>`
(`Rule-Obj-Flow`) → **OQ-004 (terjawab)**. Ini yang mengoreksi daftar OQ-005 dari 6 modul menjadi 5.

## 2. Proses / fitur utama

Rincian: **`../flows/PremiumList Life.md`** (13 rule ditelusur).

`[terverifikasi]` Graf: 3 Assignment, 3 Decision, 1 Utility, 11 connector.

```
Start1 ─(Always)─> Assignment2 "Input Offer"  [WorkList]
                        │ FlowAction: InputDataOfferLife
                        v
                    Decision1 "Accept"  (DecisionTable: IsLifeAccepted)
                        ├─ Status=Decline ─> End1
                        └─ Status=Confirm ─> Decision3 "FlagOnGoingPolicy"
                                              (DecisionTable: IsFlagOnGoingPolicy)
                              ├─ Status=Offer   ─> END52   (selesai di tahap penawaran)
                              └─ Status=Premium ─> ASSIGNMENT63 "Input Premium List Detail" [WorkList]
                                                      │ FlowAction: ShowLifePremiumDetail
                                                      v
                                                  Decision2 "Accept" (IsLifeAccepted)
                                                      ├─ Status=Decline ─> End1
                                                      ├─ Status=Reject  ─> Assignment2  (balik ke Input Offer)
                                                      └─ Status=Confirm ─> Utility1
                                                                             │ Activity: InsertJsonPolisLife_Act
                                                                             └─> END52
```

`[terverifikasi]` Nilai status connector: `Confirm`, `Decline`, `Reject`, `Offer`, `Premium`.
**Arti masing-masing belum terverifikasi** selain dari posisinya di graf → OQ-020.

`[terverifikasi]` Ketiga Assignment memakai **`WorkList`** → OQ-028.

`[pertanyaan terbuka]` `Assignment1` "Input Premium List Summary" **tidak dituju connector mana
pun** — kejadian ketiga pola yang sama (NB Treaty In, Endorsement Life, modul ini) → **OQ-023**.
Layar itu jelas bagian dari proses (ada Section dan FlowAction-nya) tetapi **tidak terhubung ke
graf**.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 37 Activity, 25 Connect-SQL, 19 Section, 10 FlowAction,
9 ReportDefinition, 9 Harness, 7 DataTransform, 2 When, 2 DecisionTable, **1 HTMLRule**, 1 Flow,
1 ConnectREST.

| Objek | Rule perujuk |
| --- | ---: |
| `POOLDATA.M_TEMPUPLOADLIFE`, **`POOLDATA.JSON_OFFER_LIFE`**, `JSON_POLIS`, `CATEGORY_ATTACH_REAS` | 2 masing-masing |
| `V_TEMPRETENSICEDING`, `V_MIN`, `V_COUNT`, `VCOUNT`, `RICOMM_LIFE`, `RATE_LIFE` | 1 masing-masing |

`[terverifikasi]` **Kode literal yang menggerbangi logika** (dari `../flows/PremiumList Life.md`):
empat ID `1000032`–`1000035` menggerbangi `Property-Set` berbeda — **mewakili apa belum
terverifikasi** → **OQ-031**. Ambang **tanggal 25** juga muncul di modul ini → **OQ-030**.

## 4. Integrasi eksternal

`[terverifikasi]` Satu `RULE-CONNECT-REST`: `convertJsonNusareToProduction`
(`SETTING` → `LinkService!LinkService`, tanpa URL literal).

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **Endorsement Life** | `Activity/InsertJsonPolisLife_Act.xml` dipakai kedua modul; `DecisionTable/IsLifeAccepted.xml` dipakai kedua modul; `PremiumListSummary.PL_NUMBER_EDM` dirujuk dari Endorsement Life | `[terverifikasi]` |
| **Master Product Name Life** | class `ASM-FW-GISFW-INT-PRODUCT_LIFE` (6 rule); `RATE_LIFE`, `RICOMM_LIFE` | `[terverifikasi]` |
| **Claim Life / Komite Claim Life** | class `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` (11 rule) dipakai bersama | `[terverifikasi]` |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **Tujuh stored procedure**, isinya tidak ada di korpus (OQ-002):
**`POOLDATA.GETQUARTER`**, **`POOLDATA.GETQUARTERRETRO`**, **`POOLDATA.INSERTJSONOFFERLIFE`**,
**`POOLDATA.INSERTJSONPOLISLIFE`**, `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY`,
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`, `DBMS_LOB.CREATETEMPORARY`. Empat yang pertama
**hanya muncul di modul ini**.

`[terverifikasi]` **Dua `DecisionTable`** menggerbangi seluruh percabangan flow
(`IsLifeAccepted`, `IsFlagOnGoingPolicy`) — **baris keputusannya tidak ikut terekspor**, seperti
seluruh 49 `DecisionTable` korpus → **OQ-043**. Jadi apa yang menyebabkan hasil `Offer` vs
`Premium`, atau `Confirm` vs `Reject` vs `Decline`, **tidak dapat dinyatakan**.

`[terverifikasi]` `IsLifeAccepted` juga ada sebagai **dua tipe rule** di korpus → OQ-026.

`[terverifikasi]` `POOLDATA.JSON_OFFER_LIFE` dan `JSON_POLIS` — struktur JSON tidak ada di korpus
→ OQ-012.

`[terverifikasi]` `IsPEGAPROD` ada; kondisinya tidak terbaca → OQ-029.

**Guard identitas** `[terverifikasi]`: `OperatorID.pyUserIdentifier` **0 file**,
`OperatorID.pyTelephone` **0 file**.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-002**, **OQ-004** (terjawab), **OQ-012**, **OQ-018**, **OQ-020**, **OQ-023**, **OQ-026**,
**OQ-028**, **OQ-029**, **OQ-030**, **OQ-031**, **OQ-043**.
