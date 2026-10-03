# 07: Business + penonaktifan

**Status:** ready-for-agent

**Blocked by:** 04 (kombinasi tahun/grup/jenis dibuka oleh kontrak)

## Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin mencatat **jenis bisnis** yang ditanggung sebuah
kontrak dengan kode dan namanya, dan **menonaktifkan** satu baris tanpa menghapusnya — supaya
riwayatnya tetap ada. *(User story 15–16 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas baris bisnis |
| `internal/repository` | Upsert dikunci `ID`; baca daftar per kombinasi |
| `internal/services` | Penonaktifan; pembaruan seluruh field |
| `internal/handlers` | Endpoint daftar + simpan bisnis |
| `frontend/` | Grid bisnis di layar kontrak |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveTreatyBusinessDetail_Act` | `@BASECLASS` / `SAVETREATYBUSINESSDETAIL_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveTreatyBusinessDetail_Act.xml` | orkestrator simpan |
| `NewTreatyBusinessDetail_Act` | `@BASECLASS` / `NEWTREATYBUSINESSDETAIL_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/NewTreatyBusinessDetail_Act.xml` | baris baru |
| `SetUbahTreatyBusinessList_Act` | `@BASECLASS` / `SETUBAHTREATYBUSINESSLIST_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetUbahTreatyBusinessList_Act.xml` | muat untuk diubah |
| `SaveMasterTreatyBusiness_SQL` | `ASM-FW-GISFW-INT-TREATYBUSINESS` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyBusiness_SQL.xml` | → `POOLDATA.PEGA_TREATYBUSINESS` |
| `GetMasterBusinessList` | `ASM-FW-GISFW-INT-TREATYBUSINESS` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetMasterBusinessList.xml` | baca daftar per kombinasi |
| `DeleteRowBusinessList` | `ASM-FW-GISFW-INT-TREATYBUSINESS` / `ASM!DELETEROWBUSINESSLIST` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteRowBusinessList.xml` | ⚠️ memegang **kedua** tabel kembar |
| `BrowseTreatyBusiness_RD`, `BrowseFilterBusiness_RD` | `ASM-FW-GISFW-INT-TREATYBUSINESS` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/` | daftar & saringan |

`[terverifikasi]` Parameter `PEGA_TREATYBUSINESS` (12 in + 2 out): `ID`, `IsActive`, `TreatyYear`,
`TreatyYearID`, `TreatyGroupID`, `TreatyGroupName`, `ReinsTypeID`, `ReinsTypeName`, `BizCode`,
`BIZNAME`, `UserID`, `TglUpdate`.

`[data DBA]` Upsert dikunci `ID`; identitas `'1' || lpad(TREATY_BUSINESS_SEQ.nextval, 6, '0')`;
**tidak `COMMIT` sendiri**; `StsSimpan` **1 = sukses / 0 = gagal**.
⚠️ **Saat UPDATE, procedure existing hanya mengisi `ISACTIVE`, `BIZCODE`, `BIZNAME`, `USERID`,
`TGLUPDATE`** — `REINSTYPEID` dan kawan-kawan hanya diisi saat INSERT.

## ADR terkait

**ADR-0006** (identitas lewat sequence), **ADR-0007** (jejak audit), **ADR-0015** (kegagalan
eksplisit).

## Acceptance criteria

- [ ] Jenis bisnis dicatat dengan **kode** dan **nama**-nya pada sebuah kontrak. *(AC 21 spec;
      User story 15)*
- [ ] Satu baris bisnis dapat **dinonaktifkan** tanpa dihapus, dan tetap terbaca sebagai baris
      nonaktif. *(AC 22 spec; User story 16)*
- [ ] ⚠️ Memperbarui baris bisnis **memperbarui seluruh field yang dikirim** — **bukan** hanya lima
      kolom seperti procedure existing. Test yang menemukan field terkirim yang tidak tersimpan
      **gagal**. *(AC 23 spec; `[keputusan work owner]` — perbaikan sadar atas perilaku existing)*
- [ ] Identitas baris bisnis **tidak pernah diketik pengguna**; berbentuk `'1'` + 6 digit.
      *(AC 5, 6 spec; **ADR-0006**)*
- [ ] Menyimpan baris yang **sudah ada** memperbaruinya, **bukan** menambah baris baru.
      *(AC 8 spec)*
- [ ] Daftar bisnis disaring **per kombinasi** (tahun, grup, jenis reasuransi), bukan per `ID`
      kontrak.
- [ ] ⚠️ Penghapusan baris bisnis menyentuh **satu tabel saja** — tidak ada tabel kembar JSON yang
      ikut dihapus. *(AC 63, 64 spec; penyimpangan sadar 1)*
- [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan**. *(AC 38, 39 spec;
      **ADR-0015**)*
- [ ] Setiap penyimpanan mencatat **jejak audit**. *(AC 41 spec; **ADR-0007**)*

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Rule ini adalah bukti paling telanjang dari dualitas tabel kembar.** `[terverifikasi]`
`RDBList/DeleteRowBusinessList.xml` memegang dua SQL sekaligus:

```
<pyBrowseSQL>  delete from treatybusiness   where id = {…}
<pyDeleteSQL>  delete from m_treatybusiness where id = {…}; commit;
```

Satu rule, dua tabel, dua ejaan. Setelah penyimpangan sadar 1 diterapkan (tiket 01), hanya yang
pertama tersisa.

⚠️ **`TREATYYEARID` boleh NULL di data lama.** `[terverifikasi]` Kaskade hapus existing
(`RDBList/DeleteFromTREATYCONTRACT_SQL.xml`) menulis
`(TREATYYEARID = {…} OR TREATYYEARID IS NULL)` — pengakuan bahwa sebagian baris bisnis lama tidak
punya nilai itu. Migrasi (tiket 01) dan pembacaan di sini harus **tahan terhadap NULL**, bukan
mengandaikannya selalu terisi.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```
