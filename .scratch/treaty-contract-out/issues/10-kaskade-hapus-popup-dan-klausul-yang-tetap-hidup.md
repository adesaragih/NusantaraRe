# 10: Kaskade hapus + popup konfirmasi — klausul **tetap hidup**

**Status:** selesai (29-09-2026)

**Blocked by:** 06 (security), 07 (business), 09 (pembungkus transaksi)

## Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin menghapus sebuah kontrak **beserta** reinsurer,
security, dan business-nya supaya tidak ada sisa yang menggantung — tetapi saya ingin **diberi
peringatan berisi jumlah baris yang akan ikut terhapus** lebih dulu, supaya saya dapat membatalkan.
*(User story 26–27 di spec)*

Dan sebagai **underwriter**, saya ingin **klausul bertahan** ketika kontrak dihapus, karena klausul
itu milik tahun/grup/jenis reasuransi, bukan milik satu kontrak. *(User story 23)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Kaskade hapus tiga anak + kontrak, **dalam satu transaksi** |
| `internal/services` | Hitung jumlah baris terdampak sebelum menghapus |
| `internal/handlers` | Endpoint pratinjau dampak + endpoint hapus |
| `frontend/` | Popup konfirmasi Ya/Batal dengan rincian jumlah |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `DeleteFromTREATYCONTRACT_SQL` | `ASM-FW-GISFW-INT` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteFromTREATYCONTRACT_SQL.xml` | kaskade empat tabel |
| `DeleteFromTreatyReinsurer_Act` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteFromTreatyReinsurer_Act.xml` | hapus security + reinsurer |
| `DeleteTreatyReins_Act` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `DELETETREATYREINS_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/DeleteTreatyReins_Act.xml` | orkestrator |
| `DeleteRowBusiness` | `@BASECLASS` / `DELETEROWBUSINESS` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/DeleteRowBusiness.xml` | hapus baris bisnis |
| `DeleteSecurityReinsurer` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteSecurityReinsurer.xml` | hapus security |
| `BrowseDeleteRowTreatyInContract` | `@BASECLASS` / `BROWSEDELETEROWTREATYINCONTRACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/BrowseDeleteRowTreatyInContract.xml` | ⚠️ bagian **salin**-nya tidak dimigrasikan |

`[terverifikasi]` Kaskade existing — empat `DELETE` berurut dalam SQL mentah, lalu `COMMIT`:

1. `DELETE FROM treatycontract WHERE id = … AND IDTREATYYEAR = …`
2. `DELETE FROM treatybusiness WHERE TREATYYEAR = … AND (TREATYYEARID = … OR TREATYYEARID IS NULL)
   AND TREATYGROUPID = … AND REINSTYPEID = …`
3. `DELETE FROM MTREATYSECURITY WHERE REAS_ID IN (SELECT id FROM TREATYREINSURER a WHERE …)`
4. `DELETE FROM TREATYREINSURER WHERE TreatyYear = … AND TreatyGroupID = … AND ReinsTypeID = …`

⚠️ **`PROPORTIONALARRG` tidak ikut** — dan itu **benar**.

⚠️ **Penyimpangan sadar 4 — kaskade + popup; klausul sengaja dikecualikan.**
`[fakta bisnis — work owner]` Klausul menggantung pada **(TreatyYear, TreatyGroupID, ReinsTypeID)**
dan **milik level tahun/grup/jenis**, dipakai bersama lintas kontrak. Ketidakikutannya adalah
**desain, bukan bug**.

## ADR terkait

**ADR-0007** (jejak audit setiap transisi), **ADR-0015** (kegagalan ditangani eksplisit).

## Acceptance criteria

- [ ] ⚠️ Menghapus kontrak **mengkaskade** ke `TREATYBUSINESS`, `MTREATYSECURITY`, dan
      `TREATYREINSURER`. *(AC 42 spec; User story 26; penyimpangan sadar 4)*
- [ ] ⚠️ Penghapusan didahului **popup konfirmasi Ya/Batal** yang **menyebut jumlah baris tiap
      jenis** yang akan ikut terhapus. Memilih **Batal** membatalkan seluruhnya dan **tidak
      mengubah apa pun**. *(AC 43 spec; User story 27; penyimpangan sadar 4)*
- [ ] ⚠️ **`PROPORTIONALARRG` (klausul) TIDAK ikut terhapus** dan tetap dapat dibaca setelah
      kontrak dihapus. Test yang menemukan klausul ikut terhapus **gagal**. *(AC 44 spec;
      User story 23; `[fakta bisnis — work owner]` — desain, bukan bug)*
- [ ] Popup **menyebut secara eksplisit** bahwa klausul **tidak** akan terhapus, supaya pengguna
      tidak menyangka sebaliknya. *(turunan AC 43, 44)*
- [ ] Seluruh kaskade berjalan dalam **satu transaksi**; kegagalan di langkah mana pun
      **membatalkan seluruhnya**. *(AC 45 spec; tiket 09)*
- [ ] Jumlah yang ditampilkan popup **sama** dengan jumlah yang benar-benar terhapus — dihitung dari
      kombinasi yang sama, bukan dari perkiraan.
- [ ] ⚠️ Kaskade menyentuh **satu keluarga tabel saja** — tidak ada tabel kembar JSON yang ikut
      atau tertinggal. *(AC 63, 64 spec; penyimpangan sadar 1)*
- [ ] Penghapusan mencatat **jejak audit** — siapa, kapan, dan berapa baris tiap jenis.
      *(**ADR-0007**)*
- [ ] Penghapusan yang gagal menghasilkan kegagalan **terang-terangan**. *(AC 38 spec;
      **ADR-0015**)*

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Existing tidak konsisten terhadap tabel kembar.** `[terverifikasi]` Kaskade ini menghapus
`treatycontract` tetapi **bukan** `M_TREATYCONTRACT`, sedangkan
`RDBList/DeleteRowBusinessList.xml` menghapus **keduanya**. Inkonsistensi itu **lenyap dengan
sendirinya** begitu penyimpangan sadar 1 diterapkan (tiket 01) — jangan mencoba "memperbaikinya"
dengan menambahkan penghapusan tabel JSON.

⚠️ **`DeleteFromTreatyReinsurer_Act` tanpa `COMMIT`** `[terverifikasi]` — konsisten dengan temuan
bahwa modul ini menyerahkan batas transaksi kepada pemanggil. Di sistem baru, pemanggil itu adalah
pembungkus transaksi tiket 09.

⚠️ **`InputData.CARI19` dipakai untuk dua arti berbeda di existing** `[terverifikasi]`:
`ReinsTypeID` pada kaskade hapus ini, tetapi `userId` pada jalur salin
(`RDBList/SaveMasterCopyData_SQL.xml`). Jebakan bagi migrasi yang meniru nama slot — di sistem baru
setiap nilai punya namanya sendiri. Jalur salin sendiri **tidak dimigrasikan** (tiket 03).

⚠️ **`TREATYYEARID` boleh NULL** `[terverifikasi]` — kaskade existing menulis
`(TREATYYEARID = … OR TREATYYEARID IS NULL)`. Kaskade baru harus **tahan terhadap NULL** pada
kolom itu, bukan mengandaikannya selalu terisi; bila tidak, baris bisnis lama akan tertinggal.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — kaskade yang **mengecualikan satu tabel**
adalah pernyataan tentang basis data; memalsukannya berarti tidak mengujinya. Uji dua arah: tiga
anak **hilang**, dan klausul **masih ada**.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Pembacaan ulang XML — 29-09-2026 (sesi modul, lanjutan 1)

Nomor baris = baris mentah berkas korpus `Treaty Contract Out/`.

| Unsur | Bukti | Dibawa sebagai |
| --- | --- | --- |
| tombol hapus kontrak | `Section/InputTreatyContractReinsType.xml` b11809 `Delete` → b11833 `BrowseDeleteRowTreatyInContract` dengan `IDTreatyYear` b11848, `TreatyYear` b11854, `TreatyGroupID` b11860, `ReinsTypeID` b11866, `ID` b11872; **tanpa konfirmasi** | `Delete` → popup Ya/Batal (penyimpangan sadar 4) |
| orkestrator | `Activity/BrowseDeleteRowTreatyInContract.xml` b371–b482 (`CARI16..CARI20`), b615 `DeleteFromTREATYCONTRACT_SQL`, b762 `"Data Berhasil di Hapus"`; langkah salin b879/b922 di-remark | pesan VERBATIM; jalur salin tidak dibawa |
| kaskade | `RDBList/DeleteFromTREATYCONTRACT_SQL.xml` b80–b82 kontrak (`id` + `IDTREATYYEAR`), b83–b87 business (`(TREATYYEARID = … OR TREATYYEARID IS NULL)` b85), b88–b89 security (`REAS_ID IN (SELECT id FROM TREATYREINSURER …)`), b90–b93 reinsurer, b94 `COMMIT;` | empat DELETE urutan sama di tabel `T_`, satu transaksi, nol COMMIT di teks SQL; klausul TIDAK disentuh |
| hapus reinsurer | `Section/ViewDetailTreatyReinsurerGrid1.xml` b4936 `Delete` → b4953 `DeleteTreatyReins_Act` (b248 `TempCariTreatyInsurer.ID = .ID`, b408 → `DeleteFromTreatyReinsurer_Act` b60–b64: security lalu reinsurer); tanpa konfirmasi, tanpa pesan | popup (jumlah security) → hapus satu transaksi |
| panel rinci exclusion | `FlowAction/DetailTreatyExclustion.xml` b92 → `Section/DetailTreatyExclustion_Sec.xml` (`Add` b881 → `NewTreatyArrExclutionTreaty`, `Edit` b2032 → `SetTreatyArrExclustionTreatyOccupation_Act`), dipakai sebagai `pyEditAction` grid exclusion Occupation (`GridTreatyArrangementExclutionTreatyOccupation.xml` b8801, `expandPane`) | baris klausul exclusion — ikut TETAP HIDUP; di sistem baru form baris di `PanelJenisKlausul` |

### Ralat bertanggal 29-09-2026

1. **Pega tidak punya popup konfirmasi** untuk hapus kontrak maupun reinsurer — popup Ya/Batal adalah penyimpangan sadar 4,
   seluruh teksnya `[tidak ada di korpus]` kecuali pesan sukses b762.
2. **Angka popup = angka terhapus** dijamin di server: `DELETE` membawa jumlah yang dilihat pemakai; di dalam transaksi
   kontrak dikunci, dampak dihitung ULANG dengan saringan yang sama (satu sumber `WHERE`), berbeda → 409 tanpa menghapus,
   dan jumlah baris yang benar-benar terhapus dibandingkan lagi (berbeda → transaksi dibatalkan).
3. **Klausul yang "tetap hidup" dihitung** untuk popup: baris `T_PROPORTIONALARRG` tahun itu yang `REINSTYPEID` atau
   `PARENTREINSTYPEID`-nya jenis reasuransi kontrak `[keputusan kami]` (**OQ-TCO-20**).
4. **Hapus reinsurer** memakai penghapusan eksplisit security lalu reinsurer (jumlahnya dicatat), walau FK `ON DELETE
   CASCADE` migrasi 303 juga ada.
5. `DetailTreatyExclustion` → `DetailTreatyExclustion_Sec` (brief) ternyata panel rinci baris exclusion, bukan popup
   hapus; dicatat di paritas sebagai bagian klausul yang tidak disentuh kaskade.

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| repository | `tco_kaskade.go` (+uji) | `KaskadeTCO`: dampak + kaskade, satu sumber saringan, langkah kaskade diuji langsung (nol langkah klausul) |
| services | `tco_kaskade.go` (+uji) | `KaskadeTCO`: pratinjau, hapus satu transaksi dengan hitung ulang + konfirmasi, jejak berjumlah |
| handlers | `tco_kaskade.go` (+uji, +uji `db` dua arah) | 4 rute (`dampak-hapus` + `DELETE` kontrak/reinsurer) |
| frontend | `KonfirmasiHapusTCO.tsx` (+uji), `HAPUS_TCO`, `api.ts` (+4), tombol `Delete` kontrak & reinsurer hidup | popup memakai `Modal` bersama (X, backdrop, Batal, Escape) |

**Status:** selesai 29-09-2026 — commit `treaty-contract-out: tiket 10 — kaskade hapus, popup, klausul yang tetap hidup`.
