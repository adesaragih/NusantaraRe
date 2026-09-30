# 08: Lampiran — opsional, terpisah dari metadata, dengan status yang terlihat

**Status:** ready-for-agent

**Blocked by:** 02 (lampiran melekat pada produk yang sudah tersimpan)

## Hasil & nilai pengguna

Sebagai **admin master**, saya melampirkan berkas pada produk, mengunduhnya kembali, dan
menghapusnya — dan bila unggahan gagal, **produk saya tetap tersimpan** dengan status lampiran yang
jelas, sehingga gangguan penyimpanan berkas tidak menghentikan pekerjaan saya.
*(User story 24–27 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Metadata lampiran + statusnya |
| `internal/repository` | Rekam lampiran; **klien penyimpanan berkas terpisah** |
| `internal/services` | **Pemisahan metadata dan lampiran**; status per lampiran |
| `internal/handlers` | Endpoint unggah, unduh, hapus, daftar lampiran |
| `frontend/` | Panel lampiran dengan status per berkas |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `ProductNameSaveAttachment` | `ASM-FW-GISFW-…` / `PRODUCTNAMESAVEATTACHMENT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/ProductNameSaveAttachment.xml` | simpan lampiran |
| `LoadAttachmentProdName` | `ASM-FW-GISFW-…` / `LOADATTACHMENTPRODNAME` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/LoadAttachmentProdName.xml` | muat daftar |
| `DownloadAttProdName_Act` | `ASM-FW-GISFW-…` / `DOWNLOADATTPRODNAME_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/DownloadAttProdName_Act.xml` | unduh satu |
| `DownloadAll_Act` | `ASM-FW-GISFW-…` / `DOWNLOADALL_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/DownloadAll_Act.xml` | unduh seluruhnya |
| `DeleteAttacProdName_act` | `ASM-FW-GISFW-…` / `DELETEATTACPRODNAME_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/DeleteAttacProdName_act.xml` | hapus |
| `InsertAttachProdName_Sql` | `ASM-FW-GISFW-…` / `ASM!INSERTATTACHPRODNAME_SQL` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/InsertAttachProdName_Sql.xml` | rekam lampiran |
| `GetAttachmentProdName_Sql` | `ASM-FW-GISFW-…` / `ASM!GETATTACHMENTPRODNAME_SQL` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/GetAttachmentProdName_Sql.xml` | baca lampiran |
| `DeleteAttachProdName_Sql` | `ASM-FW-GISFW-…` / `ASM!DELETEATTACHPRODNAME_SQL` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/DeleteAttachProdName_Sql.xml` | hapus rekam |
| `GetMimeType` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `GETMIMETYPE` / `RULE-OBJ-DECISIONTABLE` | `Master Product Name Life/DecisionTable/GetMimeType.xml` | jenis berkas |
| `SetCategory_act` | `ASM-FW-GISFW-…` / `SETCATEGORY_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/SetCategory_act.xml` | kategori lampiran |

⚠️ **Penyimpangan sadar 4 — lampiran opsional dan terpisah.** `[keputusan work owner]` Lampiran
adalah **pelengkap** master produk, **bukan syarat sahnya**. Memblokir penyimpanan produk karena
penyimpanan berkas sedang tak dapat dihubungi akan menghentikan pekerjaan tanpa perlu.

## ADR terkait

**ADR-0010** (penyimpanan berkas tetap di Google Storage), **ADR-0015** (efek keluar — kegagalan
tidak boleh diam-diam), **ADR-0007** (jejak audit).

## Acceptance criteria

- [ ] ⚠️ **Lampiran OPSIONAL**: produk tersimpan **tanpa** lampiran tanpa keluhan. *(AC 33 spec;
      `[keputusan work owner]` — penyimpangan sadar 4)*
- [ ] ⚠️ **Lampiran TERPISAH dari metadata**: kegagalan unggah **tidak** membatalkan produk yang
      sudah tersimpan. Produk tersimpan **lebih dulu**; lampiran menyusul. *(AC 34 spec)*
- [ ] ⚠️ Kegagalan unggah **tercatat dan terlihat** — bukan senyap — dan **dapat diulang**.
      *(AC 35 spec; **ADR-0015**)*
- [ ] **Status tiap lampiran** terbaca lewat API: **terunggah**, **gagal**, atau **belum**.
      *(AC 36 spec)*
- [ ] Lampiran dapat **diunduh** satu per satu **dan** seluruhnya sekaligus. *(AC 38 spec)*
- [ ] Lampiran dapat **dihapus**; penghapusan mencakup rekam **dan** berkasnya.
- [ ] Jenis berkas ditentukan dari isinya, dan berkas yang jenisnya tidak didukung **ditolak** dengan
      pesan yang menyebut jenisnya.
- [ ] Menghapus **produk** tidak meninggalkan lampiran yatim.
- [ ] Mengunggah berkas bernama sama **tidak** menimpa diam-diam — pengguna diberi tahu.

## Blocker

**Tidak ada.** Resolusi alamat penyimpanan, token, dan pengulangan kegagalan ditangani tiket **09**.

## Catatan

⚠️ **Urutannya mengikat: produk dulu, lampiran menyusul.** `[keputusan work owner]` Bukan sekadar
preferensi — ia yang membuat AC 34 dapat diuji: menyuntikkan kegagalan unggah lalu memastikan produk
**tetap ada**.

⚠️ **Berbeda dari Master Contract Retro Life.** Konteks itu **tanpa integrasi luar sama sekali**;
konteks ini punya efek keluar nyata, sehingga **ADR-0010** dan **ADR-0015** berlaku di sini.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**; **penyimpanan berkas di-fake** di balik
interface — yang diperiksa adalah **efeknya**: unggah terjadi atau tidak, status lampiran, kegagalan
tercatat dan dapat diulang.

```
go test ./internal/...
cd frontend && npm test
make check
```
