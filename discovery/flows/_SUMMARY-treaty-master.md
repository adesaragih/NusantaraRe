# Sintesis Lintas-Modul — Treaty & Master (modul tanpa `Flow`)

STEP D2 Tahap 5 — **tahap terakhir D2**. Disusun 2026-09-13 setelah menelusur
**Treaty In**, **Treaty In Adjustment**, **Treaty Contract Out**, **Master Product Name Life**,
**Master Contract Retro Life**. Konvensi: `_METHOD.md` + `_METHOD-noflow.md`.

Dokumen ini **menyajikan bukti dan membandingkan**. Ia **tidak** menetapkan bounded context —
itu STEP D4, setelah gate manusia.

**Cakupan:** 1.191 file (329 + 379 + 303 + 114 + 66), **12,7 %** dari 9.369 file korpus.

---

## 1. Dua bentuk modul yang berbeda tajam

`[terverifikasi]` Kelimanya sama-sama tanpa rule `Flow`, tetapi terbelah menjadi **dua bentuk**:

| | **Bentuk A — proses berjenjang** | **Bentuk B — editor master** |
| --- | --- | --- |
| Modul | `Treaty In`, `Treaty In Adjustment` | `Treaty Contract Out`, `Master Product Name Life`, `Master Contract Retro Life` |
| Properti status proses | **`StatusAkseptasi`** (4 nilai) | **tidak ada** |
| Mesin transisi | **`Akseptasi_DT`** (10 `WHEN`, 18 `OTHERWISE_WHEN`, 64 `SET`) | tidak ada |
| Tombol Submit/Decline | **ada** | tidak ada |
| Tangga peran | **4 tingkat** (`Admin` → `SecHead` → `DeptHead` → `Director`) | tidak ada |
| `FlowAction` | 32 / 43 | 2 / 11 / 2 |
| Pola rule dominan | Section + Activity per tab data | `New…` / `Set…` / `Save…` / `Delete…` / `CancelActivity…` |

Perintah audit:
```
for m in "Treaty In" "Treaty In Adjustment" "Treaty Contract Out" \
         "Master Product Name Life" "Master Contract Retro Life"; do
  echo "$m  StatusAkseptasi=$(grep -rl 'StatusAkseptasi' "$m" --include='*.xml' | wc -l)  FlowAction=$(ls "$m/FlowAction" 2>/dev/null | wc -l)"
done
```
→ `Treaty In` 61 file / 32 FA; `Treaty In Adjustment` 61 / 43; `Treaty Contract Out` **0** / 2;
`Master Product Name Life` **0** / 11; `Master Contract Retro Life` **0** / 2.

**Kesimpulan berbukti:** "tanpa `Flow`" **bukan** berarti "tanpa proses". Dua modul treaty inward
punya proses persetujuan lengkap — ia hanya tidak ditulis sebagai graf. Tiga modul lainnya memang
editor master murni.

---

## 2. Mesin status treaty inward — satu rule, dua modul

`[terverifikasi]` `DataTransform/Akseptasi_DT.xml` (`DATA-PORTAL!AKSEPTASI_DT`) **identik**
(hash `58b8e650`) di `Treaty In` dan `Treaty In Adjustment`. Tangga lengkapnya di
`Treaty In.md` §2.1.

### 2.1 Ringkas transisi `[terverifikasi]`

| Posisi sekarang | Aksi pengguna | Posisi berikutnya | `StatusAkseptasi` |
| --- | --- | --- | --- |
| `""` / `ReasTreatyInAdmin` | — (dipilih oleh `pyTelephone`) | `ReasTreatyInSecHead` | `Accept` |
| `ReasTreatyInSecHead` | `Accept` | `ReasTreatyInDeptHead` | `Accept` |
| `ReasTreatyInDeptHead` | `Accept` | `ReasTreatyInDirector` | `Accept` |
| `ReasTreatyInGroupLeader` | `Accept` | `ReasTreatyInDirector` | `Accept` |
| `ReasTreatyInDirector` | `Accept` | `""` | **`Resolve Complete`** |
| mana pun | `Reject` | `ReasTreatyInAdmin` | `Reject` |
| mana pun | `Decline` | `""` | `Decline` |
| **jalur revisi** (`RevisionState == 1`) | `Accept` di `SecHead` | `""` | **`Resolve Complete`** |

`[terverifikasi]` **Jalur revisi memotong dua tingkat** — dari Sec Head langsung selesai, tanpa
Dept Head dan Direktur.

`[terverifikasi]` **`ChooseStatusAkseptasi` punya 3 nilai** (`Accept`, `Reject`, `Decline`) tetapi
**`StatusAkseptasi` punya 4** — `Resolve Complete` hanya dihasilkan mesin, tidak pernah dipilih.

**Arti keempat nilai belum terverifikasi** (OQ-020), dan **beda `Reject` vs `Decline` tidak
dijelaskan korpus** (OQ-052).

### 2.2 Otorisasi: tiga lapis, seluruhnya di luar model peran `[terverifikasi]`

| Lapis | Mekanisme | Contoh |
| --- | --- | --- |
| 1. Masuk tangga | `OperatorID.pyTelephone` menyimpan kode peran | `"TREATY1"`, `"TREATY2"`, `"SPVTREATY1"`, `"SPVTREATY2"` |
| 2. Visibilitas tombol | **indeks tetap** daftar workbasket | `OperatorID.pyWorkBasketList(2).pyWorkBasketName = 'ReasTreatyInAdmin'` (36×) |
| 3. Pemilik tugas berikutnya | **identitas orang ter-hardcode** | `TreatyIn.PositionUsername := "<orang>"` — **5 identitas berbeda** |

`[terverifikasi]` Lapis 1 memperluas **OQ-027** (dua nilai baru: `TREATY2`, `SPVTREATY2`).
Lapis 2 adalah **OQ-051** (baru). Lapis 3 adalah **OQ-053** (baru) — dan **lebih berat dari
OQ-021**: di OQ-021 identitas orang hanya *diuji*; di sini identitas orang *ditetapkan* sebagai
pemilik tugas berikutnya.

Nilai nama orang **tidak disalin** ke artefak mana pun.

---

## 3. `Treaty In` vs `Treaty In Adjustment` — bukti untuk OQ-010

`[terverifikasi]`

| Ukuran | Hasil |
| --- | ---: |
| File bernama sama | **323** |
| — identik (normalisasi 18 tag) | **277** (85,8 %) |
| — berbeda | 46 |
| — berbeda setelah normalisasi **21 tag** | **43** |
| — **konflik semu** | **3** (keluarga `T_STORAGE_IMAGE`) |
| Hanya di `Treaty In` | 6 |
| Hanya di `Treaty In Adjustment` | **56** |

`[terverifikasi]` Yang **eksklusif** di `Adjustment`: 23 layar *data lama* (9 FlowAction +
14 Section `*OldData*`), mesin penomoran revisi (`TreatyInRevisi_post`, `GetTreatyRevisionID`,
`TreatyInSetAddendumToHistory`, `TreatyRevisionCopyAttachment`), `TreatyCreateEDM`,
`PickerTreatyInMaster(Revisi)`, dan satu-satunya `RULE-OBJ-MENU` di Tahap 5.

`[terverifikasi]` **Pembedanya terbaca**: `Param.type` bernilai `"revision"` / `"adjustment"`,
yang men-set `TreatyIn.EDMState` `"1"` / `"3"`. Ditambah `RevisionState` yang memendekkan tangga
persetujuan (§2.1).

`[terverifikasi]` Untuk 43 rule yang tetap berbeda, **struktur langkah dan Section yang di-include
justru identik** (diperiksa pada `TreatyInSubmit` dan `Harness/InputTreatyInOffer`). Selisihnya ada
pada tingkat yang tidak terbaca dari tag struktural — **tidak dikarakterisasi**, batas cakupan
telusur.

**Pola ini sejajar dengan OQ-015** (facultative NB/RNW/EDM): satu basis rule, dibedakan saat
runtime oleh nilai properti. **Penetapan konteks bukan di sini.**

---

## 4. `Treaty Contract Out` — nama modul menyesatkan (OQ-022 terjawab)

`[terverifikasi]` Lima uji terpisah, seluruhnya dapat diaudit ulang (`Treaty Contract Out.md` §2):

| Uji | Hasil |
| --- | --- |
| Objek database | **0** objek bernama `*_OUT*`; yang ada `M_PROPORTIONALARRG` (10), `MTREATYSECURITY` (5), `TREATYREINSURER` (3), `TREATYBUSINESS` (3), `TREATYCONTRACT`, `M_TREATYYEAR`, `TREATYEXCHANGE` |
| Class rule | **0** class bernama `*OUT*`; justru **5 rule berclass `ASM-FW-GISFW-INT-TREATY_IN`** |
| Penamaan rule | 81 `*TreatyArr*` vs **4 `*Out*`** (seluruhnya lampiran) |
| Arah tulis | **9 Connect-SQL menulis master** lewat `PEGA_TREATYCONTRACT`, `PEGA_TREATYYEAR`, `PEGA_TREATYREINSURER`, `PEGA_TREATYBUSINESS`, `PEGA_PROPORTIONALARRG`, `PEGA_M_PROPORTIONALARRG_CHILD`, `PROSESCOPY` |
| Siapa yang menyentuh outward | `NB Treaty In` **8** file, `Claim Non Prop` **4**, `EDM Treaty In` **4**, `Treaty In Adjustment` 2, `Claim Prop` 1, **`Treaty Contract Out` 0** |

**Jawaban `[terverifikasi]`: modul ini adalah editor master *term / arrangement* kontrak treaty —
bukan modul treaty outward.** Ia memiliki 16 jenis klausul yang dapat diedit (`CancelActivity*`:
BordereAux, CashLossLimit, ClaimCoorperation, EPI, ExGratia, FacIn, PLA, ProfitCommision, Ricomm,
TreatyContract, TreatyLimit, Portfolio, TerrLimit, …).

`[terverifikasi]` Satu-satunya class treaty-outward di **seluruh korpus** adalah
`ASM-FW-GISFW-INT-TREATYOUTDETAIL` (5 rule), dan letaknya di `NB Treaty In` (3),
`Treaty In Adjustment` (1), `EDM Treaty In` (1) — **bukan** di modul ini.

**Yang masih terbuka:** *mengapa* dinamai "Out", dan ke konteks mana 303 rule ini ditempatkan (D4).

---

## 5. Master Life — dua modul, dua sifat berbeda

`[terverifikasi]`

| | `Master Product Name Life` (114 file) | `Master Contract Retro Life` (66 file) |
| --- | --- | --- |
| Harness | **1** (class entitas `…INT-PRODUCT_LIFE`) | **4** (seluruhnya `DATA-PORTAL`) |
| FlowAction | 11 (7 pemilih master + 2 konfirmasi + rate + lampiran) | 2 (keduanya tampilan rate) |
| Objek Oracle | 13, **termasuk `POOLDATA.M_TREATY_IN`** dan `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` | **4**, seluruhnya `*_LIFE` |
| Stored procedure | 5 (termasuk **`PEGA_TREATY_IN`**) | 5 (seluruhnya `INSERT*_LIFE`, **eksklusif modul ini**) |
| `ConnectREST` | 1 (`ServiceGoogle`) | **tidak ada** |
| Kolom JSON | ada (`JSONDATA` via `BrowseTreatyIn`) | **tidak ada** |
| Guard identitas | `pyUserIdentifier` 3 file | **0 — tanpa jejak sama sekali** |
| Entri register OQ-011 | 22 | **5** (paling sedikit di korpus) |

### 5.1 `Master Product Name Life` membawa jalur tulis master treaty inward `[terverifikasi]`

**35 dari 114 file (30,7 %) bernama sama dengan file di `Treaty In`**: 16 identik, 19 berbeda
(12 setelah normalisasi 21 tag; **7 konflik semu**, termasuk `SystemSettings/LinkService.xml`).

Yang tersalin adalah **infrastruktur bersama** (lampiran, penyimpanan berkas, komentar) **ditambah
tiga rule treaty inward nyata**:

| Rule | Sifat |
| --- | --- |
| `RDBList/SaveTreatyIn.xml` (`ASM!SAVETREATYIN`) | memanggil `POOLDATA.PEGA_TREATY_IN`; **identik di 6 modul** (`e371c194`) |
| `RDBList/BrowseTreatyIn.xml` (`ASM!BROWSETREATYIN`) | membaca `M_TREATY_IN`, `M_TREATY_IN_EDM`; identik di 6 modul (`196d49b5`) |
| `Activity/SetTreatyIn_Act.xml` | 15 langkah, **urutan identik** dengan `Treaty In`, isi berbeda 47 byte |

**Jawaban berbukti:** ya — modul master produk life membawa jalur tulis ke `M_TREATY_IN`.
**Nuansanya:** rule tulisnya **bukan salinan menyimpang** (identik di 6 modul, tidak terdaftar di
register OQ-011), dan **tangga persetujuan treaty tidak ikut tersalin** (`Akseptasi_DT`,
`TreatyInSubmit`, `TreatyInActionButtons` tidak ada di sini). Apakah jalur itu benar-benar
dieksekusi dari modul ini **tidak dapat dipastikan dari korpus** → **OQ-056**.

### 5.2 `Master Contract Retro Life` — kontrak retro Life, apa adanya

`[terverifikasi]` Empat grid master (`RetroLifeReinsurersList`, `RetroLimitReinsurers`,
`BusinessLifeReinsurers`, `SecurityReinsurerLife`) di atas 4 tabel `POOLDATA.*_LIFE`, dengan
5 procedure `INSERT*_LIFE` yang **hanya muncul di modul ini**.

`[terverifikasi]` **`REINSTYPEID` muncul 129×**, tetapi **tidak pernah sebagai nilai literal** —
selalu diisi dari slot generik `TempInputData.CARI3/4/5` lalu diteruskan sebagai `p_REINSTYPEID`
ke tiga procedure. Daftar nilainya **tidak dapat dinyatakan** → OQ-020, **OQ-057**.

`[terverifikasi]` **Kode `OR`** (pada `GetReinsTypeOR_Life`, `BrowseReinstypeOR_SQL` di
`Master Product Name Life`) juga **tidak pernah muncul sebagai nilai literal** —
`grep -rhoE "\"OR[A-Z0-9]*\"|'OR[A-Z0-9]*'"` → kosong. Hanya nama rule yang memuatnya, dan
`_METHOD.md` §2.2 melarang memakai nama sebagai bukti → **OQ-057**.

`[terverifikasi]` Pola slot generik `CARI1`…`CARI5` juga dipakai di `NB Treaty In`
(`InputData.CARI20`, `CARI21`, `CARI3`) dan `Treaty In Adjustment` (`InputData.CARI1`)
→ **OQ-059**.

---

## 6. Batas pengetahuan kelompok ini

| # | Batas | Modul | Sifat | OQ |
| ---: | --- | --- | --- | --- |
| 1 | Aturan simpan master treaty inward ada di 5 SP `PEGA_M_TREATY_IN*` | Treaty In, Adjustment | korpus | OQ-002 |
| 2 | Aturan simpan master arrangement ada di 11 SP; **seluruh jalur tulis lewat SP** | Treaty Contract Out | korpus | OQ-002 |
| 3 | Aturan simpan kontrak retro life ada di 5 SP `INSERT*_LIFE` | Master Contract Retro Life | korpus | OQ-002 |
| 4 | Struktur `JSONDATA` pada `M_TREATY_IN` / `M_TREATY_IN_EDM` | Treaty In, Adjustment, Master Product Name Life | korpus | OQ-012 |
| 5 | Rumus limit/layer/spreading/ROL — 31 activity, terbesar 825 KB | Treaty In, Adjustment | **cakupan telusur** | — |
| 6 | Rumus share retro & tipe proteksi — 3 activity `@BASECLASS` | Master Contract Retro Life | **cakupan telusur** | — |
| 7 | Selisih isi 43 + 12 rule yang tetap berbeda | Adjustment, Master Product Name Life | **cakupan telusur** | OQ-010, OQ-011 |
| 8 | Kondisi `IsTreatyUser` (6 kondisi kosong), `recordEvent` | Treaty In Adjustment | korpus | OQ-029 |
| 9 | Nilai `REINSTYPEID`, kode `OR`, slot `CARI*` | Master Life | korpus | OQ-020, OQ-057, OQ-059 |
| 10 | Harness 8,2 MB hanya di-grep | Treaty In, Adjustment | **cakupan telusur** | — |
| 11 | 202/303 rule di `@BASECLASS` | Treaty Contract Out | korpus | OQ-009 |

`[terverifikasi]` **Tidak ada db-link `@ASMD.SINARMAS.CO.ID` di kelima modul** — OQ-017 tidak
berlaku di Tahap 5 (ia terlokalisasi di 3 modul facultative).

`[terverifikasi]` **`IsPEGAPROD` tidak ada di kelima modul** — OQ-029 hanya berlaku lewat
`IsTreatyUser`/`recordEvent`, bukan lewat keluarga `IsPEGAPROD`.

---

## 7. Integrasi eksternal kelompok ini

`[terverifikasi]` Pola seragam dan sederhana:

| Modul | `ConnectREST` | Google Storage | `M_LINK_SERVICE` | Arasapas/Kasir/Konversi/Gemini |
| --- | --- | --- | --- | --- |
| Treaty In | 1 (`ServiceGoogle`) | ya | ya | **tidak ada** |
| Treaty In Adjustment | 1 (`ServiceGoogle`) | ya | ya | **tidak ada** |
| Treaty Contract Out | 1 (`ServiceGoogle`) | ya | ya | **tidak ada** |
| Master Product Name Life | 1 (`ServiceGoogle`) | ya | ya | **tidak ada** |
| Master Contract Retro Life | **tidak ada** | tidak | tidak | **tidak ada** |

Seluruhnya `pyBaseURLSelectionType = SETTING`, `pyBaseURLSetting = LinkService!LinkService`,
**tanpa URL literal endpoint**. Alamat sesungguhnya ada di tabel Oracle `M_LINK_SERVICE`
(OQ-047).

`[terverifikasi]` `SystemSettings/LinkService.xml` **identik** di modul-modul ini setelah
normalisasi 21 tag — konfigurasi base URL tidak bercabang (konsisten dengan temuan Tahap 4 #324).

---

## 8. Temuan lintas-korpus yang muncul di Tahap 5

| Temuan | Bukti |
| --- | --- |
| **10 tombol `(dev)` di panel aksi**, termasuk `Force Resolve Complete` dan `Force Edit` | `Treaty In/Section/TreatyInActionButtons.xml`; jumlah sama di Adjustment → **OQ-050** |
| **Guard yang tak dapat bernilai benar** (`FALSE && …`, `… && 1=2`) | ekspresi visibilitas tombol — paralel OQ-023 |
| **Otorisasi lewat indeks tetap** `pyWorkBasketList(2)` | 36 kemunculan → **OQ-051** |
| **5 identitas orang ditetapkan sebagai pemilik tugas** | `Akseptasi_DT` → **OQ-053** |
| **Nomor revisi di dalam string ID**, posisi 10–12, ter-hardcode di rule **dan** SQL | `TreatyInRevisi_post`, `GetTreatyRevisionID` → **OQ-055** |
| **Tabel internal Pega diakses langsung lewat SQL** | `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` → **OQ-058** |
| **`DBMS_OUTPUT.PUT_LINE` di SQL produksi** — satu-satunya di korpus | `Treaty Contract Out` |
| **Folder `Claude outputs`** (`.xlsx`, `.diff`) di dalam ekspor 4 modul | → **OQ-054** |
| **Konflik semu keluarga `T_STORAGE_IMAGE`** muncul lagi | 3 di Adjustment, 7 di Master Product Name Life — konsisten dengan 9 temuan Tahap 4 |

---

## 9. Yang **tidak** diputuskan di sini

Sesuai `_METHOD.md` §0, dokumen ini **tidak** menetapkan:

- apakah `Treaty In` + `Treaty In Adjustment` menjadi satu bounded context atau dua (→ **D4**);
- ke konteks mana `Treaty Contract Out` ditempatkan (→ **D4**);
- apakah master Life menjadi konteks tersendiri (→ **D4**);
- arti `StatusAkseptasi`, `EDMState`, `ProportionType`, `REINSTYPEID`, kode `OR` (→ OQ terbuka);
- apakah pola apa pun di sini baik atau buruk (→ **FASE B**, setelah gate).
