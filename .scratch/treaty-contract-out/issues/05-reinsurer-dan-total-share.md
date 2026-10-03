# 05: Reinsurer + total share

**Status:** ready-for-agent

**Blocked by:** 04 (kombinasi tahun/grup/jenis dibuka oleh kontrak)

## Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin mencatat **para reinsurer** pada sebuah kombinasi tahun,
grup, dan jenis reasuransi, masing-masing dengan **share** dan **komisi reasuransi**-nya; dan
sebagai **underwriter**, saya ingin melihat **total share** seluruh reinsurer, supaya saya tahu
apakah penempatan sudah penuh. *(User story 11–12 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas reinsurer pada kombinasi |
| `internal/repository` | Upsert dikunci `ID`; baca daftar per kombinasi |
| `internal/services` | Penjumlahan total share |
| `internal/handlers` | Endpoint daftar + simpan reinsurer |
| `frontend/` | Grid reinsurer dengan baris total |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveTreatyReinsurerDetail1_Act` | `@BASECLASS` / `SAVETREATYREINSURERDETAIL1_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveTreatyReinsurerDetail1_Act.xml` | orkestrator simpan |
| `NewTreatyReinsurerDetail_Act` | `@BASECLASS` / `NEWTREATYREINSURERDETAIL_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/NewTreatyReinsurerDetail_Act.xml` | baris baru |
| `SetUbahTreatyReinsurerList_Act` | `@BASECLASS` / `SETUBAHTREATYREINSURERLIST_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetUbahTreatyReinsurerList_Act.xml` | muat untuk diubah |
| `SetErrorMessageReinsurer` | `@BASECLASS` / `SETERRORMESSAGEREINSURER` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetErrorMessageReinsurer.xml` | ⚠️ normalisasi koma→titik |
| `SaveMasterTreatyReinsurer_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyReinsurer_SQL.xml` | → `POOLDATA.PEGA_TREATYREINSURER` |
| `GetMasterReinsurerList` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetMasterReinsurerList.xml` | baca daftar per kombinasi |
| `BrowseDetailTreatyReisurer_RD` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/BrowseDetailTreatyReisurer_RD.xml` | daftar |

`[terverifikasi]` `GetMasterReinsurerList` menyaring **`TreatyYear` + `TreatyGroupID` +
`ReinsTypeID`** — kombinasi, bukan `ID` kontrak.
`[terverifikasi]` Parameter `PEGA_TREATYREINSURER` (19 in + 2 out): `ID`, `TreatyYear`,
`TreatyGroupID`, `TreatyGroupName`, `ReinsTypeID`, `ReinsTypeName`, `ReinsurerID`, `CLIENTID`,
`NAME`, `Ricomm`, `PctShare`, `IUDate`, `UserId`, `StartDate`, `EndDate`, `StatusOn`, `StdRating`,
`OperatorName`, `TglUpdate`.

`[data DBA]` `RICOMM` dan `PCTSHARE` sudah **`NUMBER`** di existing; sisanya `VARCHAR2`.
Upsert dikunci `ID`; identitas `'1' || lpad(M_TREATYREINSURER_SEQ.nextval, 6, '0')`; **tidak
`COMMIT` sendiri**; `StsSimpan` **1 = sukses / 0 = gagal**.

## ADR terkait

**ADR-0003** (uang & persen non-float), **ADR-0006** (identitas lewat sequence),
**ADR-0007** (jejak audit), **ADR-0015** (kegagalan eksplisit).

## Acceptance criteria

- [ ] Reinsurer dicatat pada kombinasi **(tahun, grup, jenis reasuransi)**, masing-masing dengan
      **share** dan **komisi reasuransi**. *(AC 14 spec; User story 11)*
- [ ] **Total share** seluruh reinsurer pada satu kombinasi **terlihat** bagi pengguna.
      *(AC 15 spec; User story 12)*
- [ ] ⚠️ Share dan komisi diperlakukan sebagai **desimal presisi arbitrer**; **tidak** melewati
      `float` dan **tidak** dibulatkan ke bilangan bulat. *(AC 16 spec; **ADR-0003**; penyimpangan
      sadar 6)*
- [ ] Identitas reinsurer **tidak pernah diketik pengguna**; berbentuk `'1'` + 6 digit.
      *(AC 5, 6 spec; **ADR-0006**)*
- [ ] Menyimpan reinsurer yang **sudah ada** memperbaruinya, **bukan** menambah baris baru.
      *(AC 8 spec)*
- [ ] Daftar reinsurer disaring **per kombinasi**, bukan per `ID` kontrak.
- [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan**. *(AC 38, 39 spec;
      **ADR-0015**)*
- [ ] Setiap penyimpanan mencatat **jejak audit**. *(AC 41 spec; **ADR-0007**)*

## Blocker

**Tidak ada.**

⚠️ `[terbuka]` **Total share = 100% bukan gerbang di tiket ini.** Existing **tidak** menolak
kombinasi yang totalnya ≠ 100 pada tingkat reinsurer; yang ada hanyalah penegakan 100% pada baris
**anak klausul** (tiket 08). Apakah total share reinsurer wajib 100% adalah **fakta bisnis yang
belum ditetapkan** — AC di atas hanya mewajibkan totalnya **terlihat**, bukan ditegakkan. Bila
Product + UW menetapkan sebaliknya, itu AC tambahan, bukan perubahan tiket ini.

## Catatan

⚠️ **Angka hari ini disimpan sebagai teks ber-koma desimal.** `[terverifikasi]`
`Activity/SetErrorMessageReinsurer.xml` (`@BASECLASS!SETERRORMESSAGEREINSURER`) menjalankan
`@replaceAll(InputTreatyReinsurer.PctShare, ",", ".")` dan hal yang sama untuk `Ricomm` — bukti
bahwa nilai masuk dari layar dalam format berkoma. Di sistem baru **normalisasi dilakukan di batas
masukan**, sekali, bukan ditempel di tengah alur simpan.

⚠️ `[terverifikasi]` Kolom `STDRATING` adalah **field yang dipakai-ulang** — pola yang sama sudah
ditemukan di Master Contract Retro Life. Isi sebenarnya belum terverifikasi; bawa apa adanya dan
jangan menafsirkan.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — presisi desimal share dan komisi hanya
terbukti benar pada basis data sungguhan.

```
go test ./internal/...
cd frontend && npm test
make check
```
