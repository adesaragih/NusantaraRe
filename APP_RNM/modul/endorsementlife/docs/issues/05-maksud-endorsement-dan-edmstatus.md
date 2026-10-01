# 05: Maksud endorsement (`EdmType`) dan status per baris (`EDMStatus`)

> **Ralat 01-10-2026** (gelombang 2 brief, `../RALAT-DEV-01-10-2026.md` — ralat mengalahkan isi di bawah). Teks lama yang tidak berlaku:
> - **R12** — sub-jenis `EDM Type Perubahan Data` / `EDM Type Batal` → opsinya tidak ada di korpus dan tak berkolom — tidak dibangun, OQ-EDM-005.
> - **R13** — `Description` (`DesBatal`) → disimpan di `EDM_NOTE`.

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 02 (baris peserta harus sudah tersalin dan bertanda `Old`)

## Hasil & nilai pengguna

Sebagai **inputor Life**, saya menyatakan **maksud** endorsement — mengubah data atau membatalkan
polis — dan pada Perubahan Data saya dapat **menambah** peserta baru serta **menandai** peserta yang
keluar, sehingga perubahan keanggotaan tercatat per baris dan dapat ditelusuri.
*(User story 13–17 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | `EdmType` pada endorsement; `EDMStatus` pada tiap baris detail |
| `internal/repository` | Baca `EdmType` polis dari atribut di dalam CLOB |
| `internal/services` | Mesin `EDMStatus`; aturan "baris baru hanya pada `EdmType=1`" |
| `internal/handlers` | Endpoint pilih maksud; endpoint tambah/tandai-hapus peserta |
| `frontend/` | Pilihan maksud endorsement; grid detail dengan penanda status per baris |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `GetEdmTypeLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETEDMTYPELIFE` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/GetEdmTypeLife.xml` | `SELECT A.DATA_JSON.EdmType … FROM POOLDATA.JSON_POLIS A` |
| `MappingEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `MAPPINGEDMLIFE` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/MappingEDMLife.xml` | menulis `"Old"` (baris 2608); membuang baris `"Delete"` (step 11.2) |
| `SetPremi_EDM` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETPREMI_EDM` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/SetPremi_EDM.xml` (345.689 byte) | menulis `"Delete"` (1384) dan `"Batal"` (1506) |
| `SaveCSVEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SAVECSVEDMLIFE` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/SaveCSVEDMLife.xml` | menulis `"New"` (2714) |
| `EditDetail_Section` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `EDITDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION` | `Endorsement Life/Section/EditDetail_Section.xml` | layar sunting detail |

### `EdmType` — dua maksud

`[keputusan work owner]` **`1` = Perubahan Data, `3` = Batal.** Nilai `2` dan `4` **tidak dipakai** —
enum efektif `{1, 3}`.

`[terverifikasi]` Nilai `3` terbukti **dari kode**: precondition `SetPremi_EDM` step 2.1
`.EdmBatal=="True" || pyWorkPage.EdmType==3` (baris 1273).

`[terverifikasi]` `EdmType` tersimpan **di dalam CLOB `DATA_JSON`**, **bukan** sebagai kolom.

### `EDMStatus` — empat nilai, bedanya **cakupan minus**

`[keputusan work owner]`

| Nilai | Kapan | Akibat pada nilai uang |
| --- | --- | --- |
| `Old` | warisan polis new business | tidak diubah |
| `New` | peserta ditambah — **hanya pada `EdmType=1`** | positif, baris baru |
| `Delete` | peserta dihapus dalam Perubahan Data | **diminuskan — selektif per peserta** (tiket 06) |
| `Batal` | lewat `EdmType=3` | **seluruh peserta — diminuskan menyeluruh** (tiket 06) |

⚠️ **`Delete` bukan sekadar penanda** — ia **juga** menghasilkan nilai negatif. Perhitungannya
dikerjakan tiket **06**; tiket ini menegakkan **penandaannya**.

## ADR terkait

**ADR-0011** (mesin status per baris — pola sama dengan `AdjustmentList` di Claim Life),
**ADR-0007** (jejak audit), **ADR-0003** (uang non-float).

## Acceptance criteria

- [ ] `EdmType` hanya menerima **`1`** dan **`3`**; nilai lain **ditolak**. *(AC 11 spec)*
- [ ] `EDMStatus` hanya menerima **`Old`**, **`New`**, **`Delete`**, **`Batal`**. *(AC 12 spec)*
- [ ] Baris **`New`** hanya dapat lahir pada `EdmType=1`; pada `EdmType=3` penambahan peserta
      **ditolak**. *(AC 13 spec)*
- [ ] `EdmType=3` **otomatis** menandai **seluruh** peserta polis sebagai `Batal` — tanpa penandaan
      satu per satu. *(AC 14 spec)*
- [ ] `Delete` menandai **hanya peserta yang dipilih**. *(AC 15 spec)*
- [ ] Baris `Old` yang tidak disentuh **tetap `Old`** dan nilainya tidak berubah.
- [ ] Maksud endorsement **terkunci** setelah case dibuat — percobaan mengubahnya ditolak **di sisi
      server**, bukan hanya di layar. *(AC 33 spec; penegakan layar di tiket 08)*
- [ ] `EdmType` dibaca dari atribut di dalam JSON polis, **bukan** dari kolom tersendiri; test yang
      mengasumsikan kolom **gagal**.
- [ ] Penanda status tiap baris **terlihat pengguna** di grid detail.

### Status sebagai turunan ⚠️ BARU 2026-09-16 — spec §16

- [ ] ⚠️ Status peserta adalah **turunan** dari `PARENT_ID` + aksi: hasil salin yang diubah →
      `Change`; peserta baru (`PARENT_ID` `NULL`) → `New`; baris salin yang ditandai keluar →
      `Delete`; pembatalan polis → `Batal`. *(AC 61, 66, 67 spec; penyimpangan sadar 11)*
- [ ] Peserta baru **tidak punya pengurang** — tidak ada baris lama untuk diselisih. *(AC 61 spec)*
- [ ] ⚠️ Selisih `new − old` dihitung dengan menyandingkan baris pada **`PARENT_ID`**-nya, bukan pada
      indeks. *(AC 65 spec; penyimpangan sadar 10)*

## Blocker

**Tidak ada.**

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
