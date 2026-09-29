# 07: Business + penonaktifan

**Status:** selesai (29-09-2026)

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

---

## Pembacaan ulang XML — 29-09-2026 (sesi modul, lanjutan 1)

Nomor baris = baris mentah berkas korpus; langkah aktivitas dibaca lengkap.

| Unsur | Bukti | Dibawa sebagai |
| --- | --- | --- |
| jalan masuk | `InputTreatyContractReinsType.xml` b10842 `Business List` → `BrowseTreatyBusinessList_Act` b10859 (`CARI1 = TreatyYear`, `CARI2 = TreatyGroupID`, `CARI3 = ReinsTypeID`) + `SetTreatyBusinessList_Act`; panel `ViewDetailTreatyBusinessGrid` b14064 | tombol `Business List` membuka `PanelBusinessKombinasi` |
| daftar | `RDBList/GetMasterBusinessList.xml` per kombinasi; grid `ReportDefinition/BrowseTreatyBusiness_RD.xml` saringan `A AND B AND C AND D AND E` (b590) dengan **E = `.IsActive = 1`** (b670–b675) | seluruh baris kombinasi, aktif maupun nonaktif |
| form | `ViewDetailTreatyBusinessGrid.xml` b6241 `Business Name` (pemilih `BrowseFilterBusiness_RD` b6301: `.ID` → BizCode b6318, `.Note` → BIZNAME), b6499 `Active` (radio wajib b6542, nilai dari daftar properti — tidak diekspor), b6680 `Business Code` (tersembunyi `1=2` b6831); ID, TreatyYear, TreatyGroupID, ReinsTypeID tersembunyi (b7855–b8472) | form `Business Name` + `Active`; `Business Code` baca-saja |
| master bisnis | `BrowseFilterBusiness_RD.xml` kelas `ASM-FW-GISFW-Int-BUSINESS`, INNER JOIN `BUSINESSGROUP` (b528–b536), saringan Note/ID/OLDID/GRUPBIS tanpa nilai dari layar, urut `.Note ASC` (b861–b867); parameter `TREATYGROUP` yang dikirim layar (b6363) **tidak dipakai saringan mana pun**; kolom fisik `ID`, `NOTE`, `BUSINESSGROUPID` terbukti di SQL korpus modul lain | `GET /api/treaty-contract-out/business-master` (tidak tersaring grup treaty) |
| simpan | `SaveTreatyBusinessDetail_Act`: langkah 1 `TreatyYearID = InputTreatyContract.IDTreatyYear`; langkah 4 RDB `SaveMasterTreatyBusiness_SQL` tanpa prasyarat; langkah 6 `ERRMSG4 = "Data sudah pernah di Input"` tanpa prasyarat, ditampilkan `Information` b10064 hanya bila `ERRMSG != ''` (b10215) | satu transaksi {kunci kontrak, dobel kode per kombinasi 409, tulis SELURUH medan, jejak} |
| hapus | `DeleteRowBusiness.xml` langkah 2 memanggil `DeleteRowBusinessList` dengan kelas **`ASM-FW-GISFW-Int-TREATYBUSINESS_LIFE`**; `RDBList/DeleteRowBusinessList.xml` memegang dua SQL (tabel utama + tabel dokumen kembar); langkah 3 `"Data Dengan ID" + " " + ID + " " + "Berhasil di Hapus"` | `DELETE .../business/{bid}` satu tabel, berbatas kombinasi, jejak, pesan VERBATIM |
| grid | b3422 `Treaty Group` (`.TreatyGroupName`), b3566 `Business ID` (**`.ID` baris** b4228), b3710 `Business Name` (`.BIZNAME`); `Edit` b4547, `Delete` b4826; `Close List` b10889 → `CancelActivity` | grid + kolom status `Active` |

### Ralat bertanggal 29-09-2026

1. **Grid Pega menyembunyikan baris nonaktif** (`BrowseTreatyBusiness_RD` saringan E `.IsActive = 1`). AC 22 menuntut
   baris nonaktif tetap terbaca — daftar baru memuat seluruh baris, dengan kolom `Active`.
2. **Nilai `IsActive`**: `1` terbukti (hilir `isactive='1'`, saringan grid); nilai nonaktif `0` adalah dugaan kuat —
   daftar pilihan radio milik properti tidak diekspor. **OQ-TCO-13**.
3. **Pemilih bisnis tidak tersaring grup treaty** — parameter `TREATYGROUP` dikirim tetapi tidak dipakai. Dibawa apa
   adanya. `INNER JOIN BUSINESSGROUP` diganti syarat `BUSINESSGROUPID IS NOT NULL` karena tabel `BUSINESSGROUP` tidak
   terbukti di SQL korpus (OQ-TCO-13).
4. **Anti-dobel kode bisnis per kombinasi** `[dugaan kuat]` — pesan `Data sudah pernah di Input` hanya tampil saat
   prosedur menjawab galat; kondisi prosedurnya milik DBA. 409 menyebut baris lain.
5. **Hapus Pega mungkin tidak pernah berjalan**: aktivitas memanggil rule dengan kelas `TREATYBUSINESS_LIFE`, sedangkan
   rule-nya berkelas `TREATYBUSINESS`. Sistem baru menghapus sesuai AC (satu tabel).
6. **`BIZNAME` dari master**, tidak diterima dari klien.

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_business.go` (+uji) | `PeriksaBusinessTCO`, `BusinessAktifTCO` |
| repository | `tco_business.go` (+uji) | daftar tanpa saringan aktif, perbarui SELURUH medan (penjaga AC 23 dibuktikan menggigit dengan mutasi), hapus satu tabel, master `BUSINESS` |
| services | `tco_business.go` (+uji) | `BusinessTCO` |
| handlers | `tco_business.go` (+uji, +uji `db`) | 5 rute |
| frontend | `PanelBusinessKombinasi.tsx` (+uji), `BUSINESS_TCO` (13 baris diuji ke korpus), `api.ts` (+4), tombol `Business List` hidup | |

**Status:** selesai 29-09-2026 — commit `treaty-contract-out: tiket 07 — business + nonaktif`.

## Keputusan work owner 29-09-2026

- **OQ-TCO-13 — ditutup.** Jawaban: *"0 berarti nonaktif"*. `models.BusinessNonaktif = "0"`,
  `models.BusinessAktif = "1"` dikunci `TestNilaiIsActiveKeputusanWorkOwner` `[keputusan work owner 29-09-2026]`.
