# 10: Kaskade hapus + popup konfirmasi — klausul **tetap hidup**

**Status:** ready-for-agent

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
