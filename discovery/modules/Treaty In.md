# Modul — Treaty In

STEP D3. Sintesis dari `../inventory/Treaty In.md` (D1) + `../flows/Treaty In.md` (D2 Tahap 5) +
`../flows/_SUMMARY-treaty-master.md`. **329 file**. Modul **tanpa rule `Flow`** (OQ-005).
Audit: `find "Treaty In" -name '*.xml' | wc -l`

## 1. Peran modul

`[terverifikasi]` **Master treaty inward** — layar input dan tangga akseptasi atas kontrak treaty
inward. Class terbanyak **`DATA-PORTAL`** (146 rule), lalu
`ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` (31), `ASM-FW-GISFW-INT-TREATY_IN` (26),
`ASM-FW-GISFW-DATA-TREATYINLIMITS` (23), `ASM-FW-GISFW-DATA-TREATYINSHARE` (18).

`[terverifikasi]` **Tanpa rule `Flow`, tetapi BUKAN tanpa proses**: ia punya properti status
(`StatusAkseptasi`, 4 nilai) dan tangga persetujuan 4 tingkat — hanya tidak ditulis sebagai graf
(§2.2).

## 2. Proses / fitur utama

Rincian: **`../flows/Treaty In.md`** (24 rule ditelusur). Metode titik masuk tanpa `Flow`:
**`../flows/_METHOD-noflow.md`**.

`[terverifikasi]` **Titik masuk**: `Treaty In/Harness/InputTreatyInOffer.xml` → `RULE-HTML-HARNESS`
/ `DATA-PORTAL` / `INPUTTREATYINOFFER`, **8.225.095 byte — file terbesar di korpus**. Dibaca
**hanya lewat grep**.

### 2.1 Panel tombol = daftar transisi

`[terverifikasi]` `Treaty In/Section/TreatyInActionButtons.xml`
(`RULE-HTML-SECTION` / `DATA-PORTAL` / `TREATYINACTIONBUTTONS`, 369.789 B) memuat seluruh tombol:
`Submit`, `Submit Revision`, `Save`, `Save ROL Profile`, `Fix Soa Upload RIComm`, `Actions`,
`Close`, ditambah **10 tombol ber-label `(dev)`** termasuk `Force Resolve Complete(dev)` dan
`Force Edit (dev)` → **OQ-050**.

Audit: `grep -oE "<pyLabel>[^<]{1,45}" … | grep -ci "(dev)"` → 10.

`[terverifikasi]` Rantai `Submit`: `Section/TreatyInfoSubmit` → `Activity/TreatyInSubmit.xml`
(`DATA-PORTAL!TREATYINSUBMIT`, 62.246 B, 7 langkah) → `Call TreatyInCheckID` →
`call TreatyInCheckError` → precondition `OutputParam.ERRMSG == ""` →
**`Apply-DataTransform Akseptasi_DT`** → `Call AddCommentList_Act` → `Call SaveTreatyIn_Act` →
`Call TreatyInInputVis`.

### 2.2 Mesin status — `Akseptasi_DT`

`[terverifikasi]` `Treaty In/DataTransform/Akseptasi_DT.xml`
(`RULE-OBJ-MODEL` / `DATA-PORTAL` / `AKSEPTASI_DT`) — **10 `WHEN`, 18 `OTHERWISE_WHEN`, 64 `SET`**.
Hash `58b8e650`, **identik di `Treaty In Adjustment`**.

Ringkas transisi:

| Posisi sekarang | Aksi | Posisi berikutnya | `StatusAkseptasi` |
| --- | --- | --- | --- |
| `""` / `ReasTreatyInAdmin` | (dipilih `pyTelephone`) | `ReasTreatyInSecHead` | `Accept` |
| `ReasTreatyInSecHead` | `Accept` | `ReasTreatyInDeptHead` | `Accept` |
| `ReasTreatyInDeptHead` | `Accept` | `ReasTreatyInDirector` | `Accept` |
| `ReasTreatyInGroupLeader` | `Accept` | `ReasTreatyInDirector` | `Accept` |
| `ReasTreatyInDirector` | `Accept` | `""` | **`Resolve Complete`** |
| mana pun | `Reject` | `ReasTreatyInAdmin` | `Reject` |
| mana pun | `Decline` | `""` | `Decline` |
| **jalur revisi** (`RevisionState == 1`) | `Accept` di SecHead | `""` | **`Resolve Complete`** |

`[terverifikasi]` **`ChooseStatusAkseptasi` punya 3 nilai** (`Accept`/`Reject`/`Decline`, dipilih
pengguna) tetapi **`StatusAkseptasi` punya 4** — `Resolve Complete` hanya dihasilkan mesin.
**Beda `Reject` vs `Decline` tidak dijelaskan korpus** → **OQ-052**. Arti keempat nilai
**belum terverifikasi** → OQ-020.

`[terverifikasi]` `TreatyIn.ProportionType` bernilai `"Proportional"` / `"NonProportional"` —
pembeda lini yang memisahkan layar, tab, dan perhitungan limit.

### 2.3 Otorisasi — tiga lapis, seluruhnya di luar model peran

| Lapis | Mekanisme | OQ |
| --- | --- | --- |
| Masuk tangga | `OperatorID.pyTelephone` = `"TREATY1"`, **`"TREATY2"`**, `"SPVTREATY1"`, **`"SPVTREATY2"`** | OQ-027 (dua nilai baru) |
| Visibilitas tombol | **indeks tetap** `OperatorID.pyWorkBasketList(2).pyWorkBasketName` (36 kemunculan) | **OQ-051** |
| Pemilik tugas berikutnya | `TreatyIn.PositionUsername` diisi **5 identitas orang ter-hardcode** | **OQ-053** |

`[terverifikasi]` Lapis 3 **lebih berat dari OQ-021**: di sana identitas orang hanya *diuji*; di sini
identitas orang *ditetapkan sebagai tujuan rute berikutnya*. Nilai nama orang **tidak disalin**.

Audit (menghitung tanpa menampilkan nilai):
```
awk '/<pyPropertiesName>/{n=$0;gsub(/<[^>]*>/,"",n)}
     /<pyPropertiesValue>/{v=$0;gsub(/<[^>]*>/,"",v);
        if(n=="TreatyIn.PositionUsername" && v!="" && v !~ /\./){gsub(/"/,"",v); print v}; n=""}' \
  "Treaty In/DataTransform/Akseptasi_DT.xml" | sort -u | wc -l      # -> 5
```

`[pertanyaan terbuka]` Dua guard visibilitas **tidak dapat bernilai benar** menurut ekspresi yang
terbaca (`FALSE && …`, `… && 1=2`) — dicatat apa adanya, paralel OQ-023.

## 3. Entitas & tabel data

`[terverifikasi]` Distribusi tipe rule: 149 Activity, 50 Section, 40 Connect-SQL,
35 DataTransform, 32 FlowAction, 16 ReportDefinition, 3 Harness, 1 When, 1 DecisionTable,
1 ConnectREST, 1 SystemSettings. **Nol rule `Flow`.**

| Objek | Rule perujuk |
| --- | ---: |
| `M_ATTACHMENTTREATY_2` | 4 |
| `T_STORAGE_IMAGE`, `POOLDATA.TREATYINPRODUCTION`, `POOLDATA.ACHIEVEMENT` | 3 masing-masing |
| `TREATY_IN_EDM`, **`POOLDATA.OS_AKSEPTASI_KLAIM`**, `POOLDATA.M_KATEGORIMASTERTREATY`, `CATEGORY_ATTACH_REAS` | 2 masing-masing |
| `TREATYINDETAILEDM`, `TREATYINDETAIL`, `TREATYEXCHANGEYEARLY`, `PROPORTIONALARRG`, `POOLDATA.TREATYINOFFER`, **`POOLDATA.M_TREATY_IN`**, `POOLDATA.M_TREATY_IN_EDM`, `POOLDATA.T_FOLDER_IMAGE` | 1 masing-masing |

`[terverifikasi]` 32 `FlowAction` membentuk struktur data treaty: limit & layer (`DetailLimits`,
`TotalLimits`, `Layers`, `MaxRetention`, `LimitProportional`, `LimitFacRetro`), share/retrosesi
(`Share`, `ShareRetro`, `DetailShare`, `TreatyRetroList`, `TreatyinChooseRetro`), spreading
(`SpreadingTPDtl`, `SpreadingTXOLDtl`), premi & angsuran (`Installments`, `DetailEGNPI`,
`AchievementCombine`), aksi status (`TreatyInAction`, `TreatyInDeclineConfirmation`).

`[terverifikasi]` `EGNPI` muncul **1.047 kali** di modul ini. **Kepanjangan belum terverifikasi.**

## 4. Integrasi eksternal

`[terverifikasi]` Satu `RULE-CONNECT-REST`: `ServiceGoogle` (`SETTING` → `LinkService!LinkService`,
tanpa URL literal). Penyimpanan berkas lewat `T_STORAGE_IMAGE` / `POOLDATA.T_FOLDER_IMAGE` +
`POOLDATA.GET_TOKEN_STORAGE`; alamat dari tabel `M_LINK_SERVICE` lewat `Activity/GetLinkService.xml`
→ OQ-047.

`[terverifikasi]` **Tidak ada Arasapas, Kasir, Konversi, maupun Gemini AI di modul ini**
(`ls "Treaty In/Activity/" | grep -icE "arasapas|kasir|konversi|gemini"` → 0).

## 5. Ketergantungan ke modul lain

| Ketergantungan | Bukti | Label |
| --- | --- | --- |
| **Treaty In Adjustment** | **277 dari 323** file bernama sama **identik** (85,8 %); `Akseptasi_DT` identik (`58b8e650`) | `[terverifikasi]` — OQ-010 |
| **NB Treaty In (realisasi)** | `RDBList/SaveTreatyIn.xml` (`ASM!SAVETREATYIN`, `e371c194`) & `BrowseTreatyIn.xml` (`196d49b5`) — identik di 6 modul; `POOLDATA.TREATYINPRODUCTION` | `[terverifikasi]` |
| **Domain Claim/Komite** | `POOLDATA.OS_AKSEPTASI_KLAIM` (2 rule) — tabel yang sama yang disentuh Claim & Komite | `[terverifikasi]` titik temu |
| **Master Product Name Life** | modul itu membawa `ASM!SAVETREATYIN` yang sama → penulis `M_TREATY_IN` kedua | `[terverifikasi]` — OQ-056 |

## 6. Batasan & batas pengetahuan

`[terverifikasi]` **Lima stored procedure**, isinya tidak ada di korpus (OQ-002):
`POOLDATA.PEGA_TREATY_IN`, `POOLDATA.PEGA_M_TREATY_IN_EDM`, `POOLDATA.PEGA_M_TREATY_IN_DETAIL`,
`POOLDATA.PEGA_M_TREATY_IN_DETAIL_EDM`, `POOLDATA.GET_TOKEN_STORAGE`.
**Empat dari lima adalah penulis master treaty inward** — cara master disimpan (validasi, versi,
kunci) **berada di sisi database**.

`[terverifikasi]` **Tidak ada objek db-link** di modul ini — OQ-017 tidak berlaku.

`[terverifikasi]` **31 activity limit/layer/spreading/retention/EGNPI/ROL belum habis dibaca** —
terbesar `TreatyInNPSetTotal.xml` **825.279 B**, lalu `SaveTreatyInDetailEdm_Act.xml` (764.703),
`SaveTreatyInDetail_Act.xml` (655.709), `GetAchievement.xml` (518.157),
`SetSpreadingXOL.xml` (509.259). Batas **cakupan telusur**.
`ROL`, `QS`, `XOL`, `TP`, `TXOL` **kepanjangan belum terverifikasi**.

`[terverifikasi]` `BrowseTreatyIn` menyeleksi kolom `JSONDATA` dari `POOLDATA.M_TREATY_IN` dan
`M_TREATY_IN_EDM` — **ini temuan asal OQ-012**; strukturnya tidak ada di korpus.

`[terverifikasi]` Nomor revisi disimpan di dalam string `TreatyIn.ID` pada posisi karakter 10–12
(mekanismenya di `Treaty In Adjustment`) → OQ-055.

## 7. OQ yang menyentuh modul

Dari register `../open-questions.md`:
**OQ-001**, **OQ-002**, **OQ-005**, **OQ-007**, **OQ-009**, **OQ-010**, **OQ-011**, **OQ-012**,
**OQ-018**, **OQ-020**, **OQ-021**, **OQ-023**, **OQ-024**, **OQ-027**, **OQ-047**, **OQ-050**,
**OQ-051**, **OQ-052**, **OQ-053**, **OQ-055**, **OQ-056**, **OQ-059**.
