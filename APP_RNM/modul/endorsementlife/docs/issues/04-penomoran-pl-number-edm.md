# 04: Penomoran `PL_NUMBER_EDM` dan kenaikan `PRODKE`

> **Ralat 01-10-2026** (gelombang 2 brief, `../RALAT-DEV-01-10-2026.md` — ralat mengalahkan isi di bawah). Teks lama yang tidak berlaku:
> - **R01** — *"`<nomor polis>`"* → nomor polis = `PL_NUMBER` (`InsertJsonPolisEDM` b102).
> - **R23** — `PRODKE` NB warisan kosong → dianggap versi 1 (E1) → nomor pertama `<polis>/02`; Pega memberi `/01` — OQ-EDM-008.

**Status:** done 01-10-2026 — c83bf68 (nomor lahir di `Confirm`)

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 02 (nomor dirakit dari nomor polis lama + `PRODKE` yang dibaca saat pemetaan)

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin setiap endorsement memperoleh **nomornya sendiri** yang tidak
menimpa penomoran new business, lahir **sekali saja**, dan urutannya terbaca — sehingga satu polis
dengan beberapa endorsement tetap dapat ditelusuri. *(User story 23–26 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | `PL_NUMBER_EDM` dan `PRODKE` pada entitas endorsement |
| `internal/repository` | Pembacaan `PRODKE` **satu urutan**; perakitan nomor |
| `internal/services` | Gerbang "nomor lahir sekali" |
| `internal/handlers` | Nomor tampil pada respons |
| `frontend/` | Nomor endorsement terlihat setelah terbit |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `GenerateNoEDM_Life` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `GENERATENOEDM_LIFE` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/GenerateNoEDM_Life.xml` | orkestrator, **7 langkah** |
| `Generate_NoEndorsmentLife` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!GENERATE_NOENDORSMENTLIFE` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/Generate_NoEndorsmentLife.xml` | **merakit nomor** |
| `GetProdKeOldData_SQL` | `ASM-FW-GISFW-INT-OFFERJSON` / `ASM!GETPRODKEOLDDATA_SQL` / `RULE-CONNECT-SQL` | *(dirujuk `GenerateNoEDM_Life` step 3)* | baca `PRODKE` — ⚠️ urutan berbeda |
| `GetProdkeNopolis` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETPRODKENOPOLIS` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/GetProdkeNopolis.xml` | baca `IDPEGA` — urutan `PRODKE DESC` |

`[terverifikasi]` **Rantai penomoran `GenerateNoEDM_Life`:**

| Step | Baris | Isi |
| ---: | ---: | --- |
| 3 | 680 | `RDB-List` "Get ProdKe dari JSON_POLIS" → `GetProdKeOldData_SQL` |
| 4 | 872 | "Set Prodke" — `Local.Prodke + 1` (946) → `InputData.CARI14` (966) |
| **5** | 1068 | `RDB-List` "Generate No Endorsement" → `Generate_NoEndorsmentLife` — precondition `PL_NUMBER_EDM==""` (1216) |
| **6** | 1264 | "Set No Endorsement" — `PL_NUMBER_EDM = Local.Nopolis` (1337) — precondition `PL_NUMBER_EDM==""` (1390) |
| 7 | 1430 | `Obj-Save` |

`[terverifikasi]` SQL perakit:
`SELECT NOPOLIS||'/'||{InputData.CARI14} AS HASIL1 FROM POOLDATA.JSON_POLIS WHERE NOPOLIS =
{pyWorkPage.PolicyNo} ORDER BY PRODKE DESC` → **`PL_NUMBER_EDM` = `<nomor polis>/<PRODKE+1>`**.

⚠️ `[terverifikasi]` **Dua pembaca `PRODKE` dengan urutan berbeda di Pega:**

| Rule | `ORDER BY` |
| --- | --- |
| `GetProdkeNopolis` | **`PRODKE DESC`** |
| `GetProdKeOldData_SQL` | `TGL_INPUT desc` |

## ADR terkait

**ADR-0006** — ⚠️ **TIDAK berlaku pada nomor EDM.** `PL_NUMBER_EDM` **bukan** dari
`PROC_GENERATE_SEQUENCE_NUMBER`; jangan memaksanya lewat sequence terpusat. ADR-0006 tetap berlaku
untuk `PL_NUMBER` new business (**PL-03**). Juga **ADR-0007** (jejak audit), **ADR-0015**.

## Acceptance criteria

- [ ] `PL_NUMBER_EDM` dirakit **`<nomor polis>/<PRODKE + 1>`** — **bukan** dari
      `PROC_GENERATE_SEQUENCE_NUMBER`. Test yang menemukan pemanggilan sequence terpusat di jalur ini
      **gagal**. *(AC 23 spec; **ADR-0006** tidak berlaku di sini)*
- [ ] Nomor **lahir sekali**: memanggil ulang pada endorsement yang sudah bernomor **tidak** mengubah
      nomornya dan **tidak** menaikkan `PRODKE`. *(AC 24 spec)*
- [ ] Dua endorsement berurutan atas polis yang sama memperoleh **`PRODKE` berurutan**.
      *(AC 25 spec)*
- [ ] `PL_NUMBER` new business **tidak berubah** oleh endorsement apa pun. *(AC 26 spec)*
- [ ] ⚠️ **Penyimpangan sadar — satu urutan `PRODKE`.** `PRODKE` dibaca dengan
      **`ORDER BY PRODKE DESC`** di **setiap** tempat. Ada test yang **gagal** bila ada jalur yang
      memakai `TGL_INPUT` atau urutan lain. *(AC 27 spec; `[keputusan desain]`)*
- [ ] Dua pembuatan endorsement serentak atas polis yang sama tidak menghasilkan `PRODKE` kembar.
- [ ] Nomor yang terbit **terlihat pengguna** dan dapat dibaca kembali lewat API.

## Catatan — mengapa urutannya disatukan

`PRODKE` adalah **urutan produksi** yang menjadi dasar nomor endorsement; `TGL_INPUT` adalah **waktu
pencatatan**, yang dapat menyimpang karena entri susulan atau koreksi. Bila keduanya pernah tidak
sejalan, kedua pembaca Pega memberi `PRODKE` berbeda — dan **nomor endorsement ikut salah**. Sistem
baru menutup celah itu.

## Blocker

**Tidak ada.**

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
