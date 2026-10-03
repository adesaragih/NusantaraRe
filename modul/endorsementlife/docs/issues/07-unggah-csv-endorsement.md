# 07: Unggah CSV endorsement — tanpa batas baris, mengganti bukan menumpuk

> **Ralat 01-10-2026** (gelombang 2 brief, `../RALAT-DEV-01-10-2026.md` — ralat mengalahkan isi di bawah). Teks lama yang tidak berlaku:
> - **R05** — *"`call ASM-FW-GISFW-Work-LIFE.Calculate1_Act` | ⚠️ lintas class — **tidak direplikasi**"* → tetap tidak dijalankan, padahal langkahnya hidup (b629 `·`) — dicatat sebagai OQ-EDM-003.
> - **R24** — validasi `PLAN`/`POLICY_HOLDER` per baris → seluruh baris diperiksa lebih dulu; sesudah `Add CSV Data` tombol unggah mati (`.EditInput1=1`).
> - **R31** — *"Unggah ulang mengganti"* (AC 35) → unggahan kedua sesudah `Add CSV Data` ditolak 409 (kunci `.EditInput1`), pembuangan baris `New` 2.1 tetap berjalan; CSV sesudah `Save` menghitung ulang rekap.
> - **R32** — judul di luar 4.1 diabaikan dan ditampilkan (bukan ditolak); `UW_STATUS`, `SUM_AT_RISK`, `REMAINING_PERIOD` tidak disimpan — OQ-EDM-018.

**Status:** done 01-10-2026 — e5621e6

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 05 (baris `New` menuntut mesin `EDMStatus` sudah berjalan)

## Hasil & nilai pengguna

Sebagai **inputor Life**, saya mengunggah CSV berisi peserta baru untuk endorsement Perubahan Data,
**tanpa dibatasi jumlah baris**, dan bila berkasnya salah saya dapat mengunggah ulang tanpa
menggandakan apa pun. *(User story 31–36 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Baris unggahan; hasil validasi per baris |
| `internal/repository` | Penyimpanan baris hasil unggahan; pembuangan baris `New` sebelumnya |
| `internal/services` | Parsing uang; validasi konsistensi; penandaan `New`; **pemrosesan bertahap** |
| `internal/handlers` | Endpoint unggah; endpoint tinjau hasil |
| `frontend/` | Form unggah; tabel hasil validasi berlabel baris dan kolom |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `UploadCSVEDMLifePremium_Act` | `@BASECLASS` / `UPLOADCSVEDMLIFEPREMIUM_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/UploadCSVEDMLifePremium_Act.xml` (119.366 byte) | impor berkas — **rule base-class** |
| `SaveCSVEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SAVECSVEDMLIFE` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/SaveCSVEDMLife.xml` (230.581 byte) | **penyimpan, 5 langkah, nol remark** |
| `ViewCSVResult_LifeEDM` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `VIEWCSVRESULT_LIFEEDM` / `RULE-OBJ-HTML-SECTION` | `Endorsement Life/Section/ViewCSVResult_LifeEDM.xml` | layar tinjau |

`[terverifikasi]` **`UploadCSVEDMLifePremium_Act` — hanya tiga langkah pertama hidup:**

| Step | Langkah | Status |
| ---: | --- | --- |
| 1–2 | `Page-Remove`, `Page-New` "deklarasi pyWorkPage" | aktif |
| **3** | **`Call pxUploadCSVResults`** "import data life premium" | aktif |
| ~~4–6~~ | `Page-New`, `Property-Set` (+5.2, 5.4), `Obj-Save` | **REMARK** (645, 791, 2120, 2400, 2525) |

`[terverifikasi]` **`SaveCSVEDMLife` — 5 langkah:**

| Step | Langkah | Precondition |
| --- | --- | --- |
| **2.1** | `Property-Remove` **"Remove EdmStatus \"New\""** | **`.EDMStatus=="New"`** (547) |
| 3 | `call ASM-FW-GISFW-Work-LIFE.Calculate1_Act` | ⚠️ lintas class — **tidak direplikasi**, lihat catatan |
| **4.1** | "Set property PremiumListDetail" | `.PLAN = …PremiumListDetail(1).PLAN` (2466); `.POLICY_HOLDER = …PremiumListDetail(1).POLICY_HOLDER` (2495) |
| 4.2 | `Page-Set-Messages` "Set message error" | `local.errmsg==""` (2649) |
| **4.3** | **"Set \"New\" untuk detail baru"** | **`pyWorkPage.EdmType==1`** (2791) |
| 5 | `Property-Set` | `.EditInput1 = 1` (2899) |

⚠️ `[keputusan work owner]` **Batas 50.000 baris dibuang.** Gerbang Pega
`@SizeOfPropertyList(TempWorkPage.ListLifePremiumDetailUpload)>50000` (`SetPremi_EDM` baris 3373)
**tidak direplikasi**.

## ADR terkait

**ADR-0003** (uang non-float — parsing CSV adalah titik masuk uang), **ADR-0010** (penyimpanan berkas
tetap Google Storage untuk lampiran), **ADR-0007** (jejak audit unggahan).

## Acceptance criteria

- [ ] ⚠️ **Penyimpangan sadar — tanpa batas baris.** Unggahan **tidak ditolak karena jumlah baris**.
      Test memuat berkas **di atas 50.000 baris** dan berhasil. *(AC 36 spec;
      `[keputusan work owner]`)*
- [ ] Unggahan sangat besar diproses **bertahap** (streaming/batch internal) tanpa memuat seluruh
      berkas sekaligus, dan **tanpa** menolak.
- [ ] **Unggah ulang mengganti, bukan menumpuk** — seluruh baris `EDMStatus == "New"` sebelumnya
      dibuang lebih dulu. Dibuktikan dengan mengunggah dua kali lalu menghitung baris.
      *(AC 35 spec)*
- [ ] Baris `New` **hanya lahir pada `EdmType=1`**; pada endorsement **Batal**, unggahan **tidak**
      menambah peserta. *(AC 13 spec)*
- [ ] Baris dengan **`PLAN`** atau **`POLICY_HOLDER`** berbeda dari **baris pertama** ditolak, dengan
      pesan yang menyebut **nomor baris dan nama kolom**. *(AC 37 spec)*
- [ ] Nilai uang di-parse dengan **format pemisah yang dinyatakan eksplisit**, langsung ke desimal
      presisi arbitrer — **tidak** lewat `float`. Test memuat kasus **pemisah ribuan**.
      *(AC 38 spec; **ADR-0003**)*
- [ ] Pengguna dapat **meninjau** hasil unggahan — baris lolos dan baris ditolak — sebelum menyimpan.
- [ ] Baris `Old` dan baris bertanda `Delete` **tidak tersentuh** oleh unggah ulang.

## Catatan — `Calculate1_Act` tidak direplikasi

⚠️ `[keputusan work owner]` Meski `SaveCSVEDMLife` step 3 memanggil `Calculate1_Act`
(`ASM-FW-GISFW-WORK-LIFE` / `CALCULATE1_ACT` / `RULE-OBJ-ACTIVITY`, 302.597 byte) lintas class, dan
berkasnya **identik** di kedua folder modul (diff ternormalisasi **nol baris**), **di jalur
endorsement ia TIDAK dijalankan.** Ia milik **PremiumList Life (new business)** saja. Perhitungan
endorsement dikerjakan **`SetPremi_EDM`** (tiket 06).

⚠️ **Pelajaran, seiring OQ-066: berkas identik ≠ dipakai.** Sebagaimana `<pyStepsBlockName>` tidak
boleh dipakai sendirian untuk menyimpulkan hidup/mati, **kesamaan berkas antar folder tidak boleh
dipakai sendirian untuk menyimpulkan "mesin bersama"** — dan rujukan `Call` hanya menunjukkan
**kemungkinan** pemanggilan.

## Blocker

**Tidak ada.**

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Status 01-10-2026 — K2 keputusan work owner 01-10-2026 (OQ-EDM-003) — ⏸️ terhenti, OQ-EDM-021

Keputusan *"jalankan `Calculate1_Act` saat simpan CSV"* dicatat. `Activity/Calculate1_Act.xml` dibaca langkah demi langkah
(`pyStepsBlockName` kosong di seluruh langkah; langkah 7.12 b3626 berlabel `AA` b3636 = sasaran lompatan, bukan `//`):

| Langkah | Baris | Isi | Arah prakondisi (1 lompat · 2 lanjut · 3 lewati · 6 keluar activity) |
| --- | --- | --- | --- |
| 1 | b416 | tanpa metode — gerbang | `.IsCalculationSystem` b476 T=2 **F=6**; `.Type=="QR"` b499 T=2 **F=6** |
| 2–3 | b539, b678 | buang `InData`/`DataProduct`; `InData.CARI1/2/21/26` dari kepala polis | — |
| 4 | b924 | `GetProductDtlPL`: `PRODUCT_LIFE` × `PRODUCTINWARD_LIFE` × `rate_life_summary` menurut ceding, jenis ceding, policy holder, `PRODUCTTYPE` | — |
| 5 | b1101 | retensi/limit ceding, `ShareRNM` dari produk | — |
| 6, 8 | b1361, b5575 | `DeleteTempUploadDataLife` (`delete from M_TEMPUPLOADLIFE where idpega = …`) | — |
| 7 | b1538 | ulang `TempWorkPage.ListLifePremiumDetailUpload` (b5554) | — |
| 7.2 / 7.3 | b1796 / b1978 | QS / SP: `CEDING_RETENTION`, `SUM_REASURED`, `SHARE_NUSANTARA_RE` | `TypeCeding=="1"` / `"2"` |
| 7.5–7.11 | b2480–b3452 | `CekDoubleInsured` + `InsertDataUploadLife` (temp `M_TEMPUPLOADLIFE`) | 7.6 b2706 `TypeCeding=="2"` **F=1 → lompat `AA`** |
| 7.13–7.20 | b3707–b5123 | prorata, `SelectComm_SQL`, `GetRateLifePM` (`RATE_LIFE`), `GROSS_PREMIUM`, `COMM`, `BROKERAGE_FEE`, `NET_PREMIUM` | 7.14 b4072/b4103, 7.15 b4325, 7.16 b4487 |

**Mengapa terhenti (brief §4 K2):** gerbang langkah 1 membaca `.IsCalculationSystem` — `Rule-Obj-Property` kelas
`ASM-FW-GISFW-Work-LIFE` (b6722/b6725) yang tidak diekspor; nol rule di korpus mengisinya; satu-satunya sumber yang mungkin,
`pyDefault` (model awal flow `InputEDMLife` b140), juga tidak diekspor. Tanpa nilai, Pega keluar di langkah 1 (F=6) — sehingga
perilaku yang terbukti dari XML sama dengan bawaan yang sudah dibangun: nilai uang CSV dipakai apa adanya. Membangun langkah 2–8
di belakang gerbang yang tidak pernah terbuka = kode mati. Pertanyaan ke pemilik ekspor Pega: **OQ-EDM-021**.

**Temuan tambahan bila gerbang ternyata terbuka** (untuk pembangun berikutnya): `InsertDataUploadLife` b58 menyisip
`M_TEMPUPLOADLIFE` **tanpa** kolom `IDPEGA`, sedangkan `DeleteTempUploadDataLife` b78 dan `SUM(CEDING_RETENTION)` di
`CekDoubleInsured` menyaring `idpega = …` — keduanya tidak pernah mengena baris yang baru disisip (kecuali ada bawaan/trigger
kolom, data DBA); `COUNT`/`MIN(NO)` `CekDoubleInsured` tidak menyaring `idpega` sama sekali. `SelectComm_SQL` (7.18) mengisi
`HasilComm` yang tidak dibaca langkah mana pun (7.20 memakai `HasilProductName.RICOMM`).

**Selisih `Calculate1_Act` Endorsement vs PremiumList:** NOL selisih logika. Keduanya rule yang SAMA
(`pzInsKey` `RULE-OBJ-ACTIVITY ASM-FW-GISFW-WORK-LIFE CALCULATE1_ACT #20210907T082034.453 GMT`, versi `01-01-60`,
`pxUpdateDateTime` `20210927T081024.226 GMT`) yang diekspor dua kali (cap ekspor 07-09 vs 08-09-2026). Yang berbeda hanya cap
waktu ekspor, urutan elemen parameter langkah 4/6/7.5/7.11/7.18/7.19/8, dan karenanya nomor baris definisi ulang langkah 7
(b5554 Endorsement, b5555 PremiumList) serta memo (b13 vs b12). 24 langkah `pyStepsActivityName` di keduanya.
