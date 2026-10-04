# Pemetaan `M_TREATY_IN2` → empat tab, kolom demi kolom

**3 Oktober 2026.** Empat tab — **Limits · Share · Event Limits · RNM Share** — dibaca dari satu
tabel warisan. **Nol tabel baru, nol pemuatan.**

---

## 0 · ⛔ Nama kolomnya BUKAN nama properti Pega

Rancangan ronde ini meminta pemetaan diturunkan dengan *"mengambil `pyValue` tiap medan di dalam
blok tab, lalu mencocokkannya ke nama kolom"*. **Cara itu hanya menghasilkan 8 dari 41.**

Diukur: ke-126 `pyValue` di `TreatyInTabsProportional.xml` dan ke-258 di
`TreatyInTabsNonProportional.xml` (sesudah `<pyIncludedRuleXML>` dibuang) menamai medannya menurut
**properti Pega** — `.Limit`, `.MDP`, `.Deductible`, `.AggregateLimit`, `.NetPremi`, `.SharePct`.
`M_TREATY_IN2` menamai kolomnya menurut **ekstraksinya sendiri** — `LIMIT_100`, `MDP_RNM_100`,
`EPI100`, `RIOGR`, `QSOR`, `LIABILITYQSRI`. Hanya delapan nama yang kebetulan bertemu:
`TREATYTYPE`, `LAYERTYPE`, `LAYER`, `CURRENCY`, `MDP`, `BROKERAGEPERCENTP`, `EARTHQUAKE`,
`RNMSHARE`.

Jadi pemetaannya diturunkan dari **DUA** garis bukti yang berdiri sendiri:

| Bukti | Caranya |
| --- | --- |
| **A — wilayah tab** | Tab ditandai `<pyTitle>X</pyTitle><pyHeaderType>TABBED</pyHeaderType>`. Ke-11 penanda di ekspor non-proporsional memotong berkas menjadi wilayah; seluruh `pyValue` berbentuk properti di dalam satu wilayah adalah medan tab itu. |
| **B — kecocokan NILAI** | Untuk tiap baris `M_TREATY_IN2`, elemen `Limits[]`/`Share[]` dengan `Layer` yang sama dicari di dokumen kontrak yang sama, lalu tiap kolom dibandingkan nilainya dengan tiap kunci. 7.281 baris × 41 kolom. |

⚠️ **Bukti B tidak pernah 100% untuk kolom berpasangan mata uang.** Dokumen menyimpan `Limit` dan
`Limit2` (IDR dan USD) dalam satu elemen, sementara tabel menyimpan **satu baris per mata uang** —
jadi `LIMIT_100` cocok ke `Limits.Limit` pada 32% baris dan ke `Limits.Limit2` pada 18%. Itu bukan
kelemahan bukti; itu bentuk datanya.

### Wilayah tab yang dipakai (offset bita, berkas sesudah dibersihkan)

| Tab | Berkas | Dari | Sampai |
| --- | --- | ---: | ---: |
| Limits | `TreatyInTabsProportional.xml` | 526.981 | 651.816 |
| Share | `TreatyInTabsProportional.xml` | 651.816 | 1.031.030 |
| Event Limits | `TreatyInTabsNonProportional.xml` | 223.612 | 490.872 |
| Limits | `TreatyInTabsNonProportional.xml` | 886.054 | 1.695.720 |
| Share | `TreatyInTabsNonProportional.xml` | 1.695.720 | 2.093.710 |
| RNM Share | `TreatyInTabsNonProportional.xml` | 2.093.710 | 3.291.797 |

⚠️ `TreatyInTabsProportional.xml` memuat **nol** penanda `TABBED` — ke-11 tabnya ditandai `BAR`.
Offsetnya sudah tercatat per tab di `frontend/labels.ts` sejak ronde layar 1, dan dipakai apa
adanya di sini.

---

## 1 · Pemetaan — ke-41 kolom

Kolom **Bukti A** menyebut properti Pega di wilayah tab yang bersangkutan. Kolom **Bukti B**
menyebut kunci JSON yang nilainya sama, dengan persentase baris terisi yang cocok.

### Tab **Limits** — 14 kolom

| Kolom | Golongan | Bukti A (`pyValue`) | Bukti B (nilai) |
| --- | --- | --- | --- |
| `LAYER` | teks | `.Layer` — NONPROP Limits | `Limits.Layer` **100%** |
| `LAYERTYPE` | teks | `.LayerType` — NONPROP Limits | `Limits.LayerType` **92%** |
| `BASIS_COVER` | teks | — | `Limits.Cover` **99%** |
| `TREATYTYPE` | teks | `.TreatyType` — **PROP Limits** | `Limits.TreatyType` **80%** |
| `CURRENCY` | teks | `.Currency` — NONPROP Limits | `Limits.Currency` 36% · `Currency2` 23% |
| `LIMIT_100` | **uang** | `.Limit` `.Limit2` — NONPROP Limits | `Limits.Limit` 32% · `Limit2` 18% |
| `CEDANT_RETENTION` | **uang** | `.Deductible` `.Deductible2` — NONPROP Limits | `Limits.Deductible` **51%** |
| `MDP` | **uang** | `.MDP` `.MDP2` — NONPROP Limits | — *(dihitung ekstrak)* |
| `MDP_RATIO` | persen | — | `Limits.MDPPct` **100%** |
| `ROL` | persen | `TreatyIn.TotalLimitsROL` — NONPROP Limits | `Limits.ROLPct` **88%** |
| `ADJ_RATE` | persen | — | `Limits.AdjRate` **90%** |
| `EARN_PREMIUM` | **uang** | — | — *(dihitung ekstrak)* |
| `CURRENCYRELATION` | teks | — | `Limits.CurrencyRelation` **100%** |
| `CURRENCYLIMIT` | teks | `.Currency` — NONPROP Limits | `Limits.Currency` 57% · `Currency2` 45% |

### Tab **Share** — 11 kolom

| Kolom | Golongan | Bukti A (`pyValue`) | Bukti B (nilai) |
| --- | --- | --- | --- |
| `LAYER` | teks | `.Layer` — NONPROP Share | `Share.Layer` **100%** |
| `CESSIONPCT` | persen share | `.SharePct` — PROP & NONPROP Share | — *(dihitung ekstrak)* |
| `SPREADINGTYPE` | teks | — | `Share.SpreadingTypeXOL` **60%** |
| `BROKERAGEPERCENTP` | persen share | `TreatyIn.BrokeragePercentP` — **PROP Share** · `TreatyIn.BrokeragePercent` — NONPROP Share | akar `BrokeragePercent` 51% · `BrokeragePercentP` 35% |
| `CESSION_TO_RI` | **uang** ⚠️ | — | — |
| `QSOR` | persen share | — | — *(dihitung ekstrak)* |
| `QSRI` | persen share | — | — *(dihitung ekstrak)* |
| `LIABILITYQSOR` | **uang** | — | — |
| `LIABILITYQSRI` | **uang** | — | — |
| `EPI100` | **uang** | — | — |
| `RIOGR` | **tidak tergolong** ⚠️ | — | — |

### Tab **Event Limits** — 3 kolom, dan itu TEMUAN

| Kolom | Golongan | Bukti A (`pyValue`) | Bukti B (nilai) |
| --- | --- | --- | --- |
| `LAYER` | teks | — | `Limits.Layer` 100% |
| `CURRENCY` | teks | `.Currency` — NONPROP Event Limits | `Limits.Currency` 36% |
| `EARTHQUAKE` | **uang** | `TreatyIn.Earthquake` — NONPROP Event Limits | akar `Earthquake` **36%** |

⛔ **Wilayah tab Event Limits menyebut EMPAT batas**, bukan satu: `TreatyIn.Earthquake`,
`TreatyIn.FloodJab`, `TreatyIn.FloodNation`, `TreatyIn.RSMDLimit`, masing-masing dengan mata
uangnya (`TreatyIn.CurrencyEarthquake`, `CurrencyFloodJab`, `CurrencyFloodNat`, `CurrencyRSMD`).
**`M_TREATY_IN2` hanya punya `EARTHQUAKE`** — dan kolom itu terisi pada **759 dari 7.281** baris.
Ketiga batas lain ada di akar `JSONDATA` dan **tidak punya sumber tabel**.

⚠️ Jadi pernyataan *"`M_TREATY_IN2` melayani Event Limits"* benar untuk **seperempat** tab itu.

### Tab **RNM Share** — 7 kolom

| Kolom | Golongan | Bukti A (`pyValue`) | Bukti B (nilai) |
| --- | --- | --- | --- |
| `LAYER` | teks | `.Layer` — NONPROP RNM Share | `Limits.Layer` 100% |
| `RNMSHARE` | persen share | `.RNMShare` — NONPROP RNM Share · `TreatyIn.RNMShareP` — PROP Share | akar `RNMShare` **64%** · `RNMShareP` 45% |
| `LIABILITY_RNM` | **uang** | — | — |
| `MDP_RNM_100` | **uang** | `.MDP` `.MDP2` — NONPROP RNM Share | — *(dihitung ekstrak)* |
| `EPIRNMQS100` | **uang** | — | — |
| `RNM_RETAINED_PREMI` | **uang** | `.NetPremi` — NONPROP RNM Share | — *(dihitung ekstrak)* |
| `RNM_QS_PREMI` | **uang** | `.NetPremi2` — NONPROP RNM Share | — *(dihitung ekstrak)* |

### Sepuluh kolom KEPALA — tidak satu tab pun menampilkannya

`MASTERID` · `PROPORTIONTYPE` · `CEDINGID` · `CEDING` · `SOBID` · `SOB` · `TREATYGROUP` ·
`TREATYCONTRACTNAME` · `COMMENCEMENT` · `TERMINATION`

Kesepuluhnya **kepala kontrak yang berulang pada setiap baris layer**, dan form sudah
menampilkannya di bagian atas dari `TREATY_IN`. Mengulangnya di dalam grid akan menuliskan nilai
yang sama pada setiap baris tanpa menambah satu pun keterangan.

Bukti B untuk kesepuluhnya **100%** terhadap akar dokumen — dan itu berguna, sebab ia
**membuktikan `MASTERID` benar-benar kontrak yang sama**:

| Kolom | Kunci akar | Cocok |
| --- | --- | ---: |
| `MASTERID` | `ID` | 100% |
| `PROPORTIONTYPE` | `ProportionType` | 100% |
| `CEDINGID` | `CedingID` | 100% |
| `CEDING` | `Ceding` | 100% |
| `SOBID` | `LeadingReinsSourceID` | 100% |
| `SOB` | `LeadingReinsSource` | 100% |
| `TREATYCONTRACTNAME` | `TreatyContractName` | 100% |
| `TREATYGROUP` | — | nol — nama kelompok, bukan salinan nilai dokumen |
| `COMMENCEMENT` · `TERMINATION` | — | nol — **bentuknya berbeda**: `01/01/2023` di tabel, `20230101` di dokumen |

### ⛔ Penjumlahannya HABIS

**31 kolom di tab + 10 kolom kepala = 41.** Nol kolom hilang, nol kolom terhitung dua kali.
`TestKe41KolomTerhitungHabis` di `frontend/layar.test.ts` menjaganya — kolom yang lenyap dari
layar tidak menimbulkan satu pun galat, hanya sebuah angka yang tidak ada lagi di mana pun.

---

## 2 · Golongan angka — diturunkan dari SEBARAN NILAI, bukan dari nama

| Golongan | Kolom | Dasar |
| --- | --- | --- |
| **uang**, 4 desimal | `LIMIT_100` `CEDANT_RETENTION` `MDP` `EARN_PREMIUM` `CESSION_TO_RI` `EPI100` `EARTHQUAKE` `LIABILITY_RNM` `MDP_RNM_100` `LIABILITYQSOR` `LIABILITYQSRI` `EPIRNMQS100` `RNM_RETAINED_PREMI` `RNM_QS_PREMI` | magnitudonya mencapai 1,8×10¹² — bukan persentase |
| **persen**, 2 desimal | `MDP_RATIO` (0–100) · `ADJ_RATE` (0,0028–109,6) · `ROL` (0–199,4) | ⚠️ dua yang terakhir **melampaui 100**, jadi laju, bukan bagian dari jumlah yang harus 100 |
| **persen share**, desimal apa adanya | `CESSIONPCT` (0–100) · `RNMSHARE` (0–100) · `QSOR` (0–40) · `QSRI` (0–80) · `BROKERAGEPERCENTP` (0–12,5) | jumlah beberapa baris harus **tepat 100** — lihat `PERTANYAAN-TERBUKA-PERSEN-SHARE.md` |
| **teks / pengenal**, tidak diformat | `LAYER` `LAYERTYPE` `BASIS_COVER` `TREATYTYPE` `CURRENCY` `CURRENCYLIMIT` `CURRENCYRELATION` `SPREADINGTYPE` + kesepuluh kolom kepala | `LAYER` memuat `1A` pada 222 dari 4.209 baris — pemformat angka akan merusaknya |

### ⚠️ Dua koreksi atas penggolongan yang diusulkan

| Kolom | Diusulkan | Terukur | Golongan yang berlaku |
| --- | --- | --- | --- |
| `CESSION_TO_RI` | persen | 0 – **1,434×10¹²**, 4.211 nol | **uang** — nilai sebesar itu bukan persentase |
| `RIOGR` | persen | 0 – **3.250**, 4.212 nol, nilai tersering `35` | **tidak tergolong** — melampaui 100, jadi bukan share; tidak cukup bukti untuk menyebutnya uang maupun laju |

⛔ `RIOGR` diformat sebagai **bilangan biasa** (titik ribuan, tanpa tanda `%`), dan itu pilihan
yang paling sedikit mengaku: tanda `%` akan menyatakan arti yang belum terbukti, dan membiarkannya
mentah membuat `3250` tak terbaca. **Ditagih ke pemilik proses.**

---

## 3 · ⚠️ Jangkauan tabel ini TIDAK penuh — dan layar harus mengatakannya

| | |
| --- | ---: |
| `M_TREATY_IN2` | 7.281 baris · **1.340 kontrak** |
| `(MASTERID, LAYER)` berbeda | 2.540 |
| dokumen dengan `Limits[]` berisi | **1.850 kontrak** · 4.210 elemen |
| **`Limits[]` berisi TAPI nol baris di `M_TREATY_IN2`** | **510 kontrak · 1.210 elemen** |
| baris di `M_TREATY_IN2` tapi `Limits[]` kosong | 0 |
| dokumen dengan `Share[]` berisi | 845 kontrak · 2.923 elemen |
| `Share[]` berisi tapi nol baris di `M_TREATY_IN2` | 257 kontrak |

⛔ **Grid kosong di keempat tab ini TIDAK berarti "kontrak ini memang tidak punya limit".** Pada
**510 kontrak** artinya *"sumber yang tabel ini ambil tidak mencakupnya"* — dan datanya **ada**, di
`Limits[]` dokumen. Petunjuk kosongnya (`FORM_KONTRAK.petunjukLayer`) mengatakan persis itu, dan
`layar.test.ts` menolak kalimat yang mengaku tahu sebabnya.

**Melepaskan keempat tab dari batas ini menuntut tabel pendaratan untuk `Limits[]` dan `Share[]`** —
yang ronde ini larang. Keputusannya milik pemilik proses.
