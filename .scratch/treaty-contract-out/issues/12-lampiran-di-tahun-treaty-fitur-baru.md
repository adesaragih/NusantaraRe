# 12: Lampiran di tahun treaty — **FITUR BARU**

**Status:** ready-for-agent

**Blocked by:** 03 (tahun treaty sebagai induk lampiran)

⚠️ **Ini fitur baru, bukan migrasi.** `[keputusan work owner]` Jalur lampiran di Pega
(`M_ATTACHMENTTREATY_2`) **belum rampung di-develop**. Di sistem baru fungsinya **dilengkapi** dan
melekat pada **tahun treaty**.

## Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin **melampirkan berkas pada sebuah tahun treaty** beserta
kategorinya, supaya dokumen kontraknya tersimpan bersama datanya — **tanpa** itu menjadi syarat
tersimpannya data tahun treaty; dan sebagai **tim operasi**, saya ingin **kegagalan unggah terlihat
dan dapat diulang**, serta alamat penyimpanan **dibaca dari konfigurasi saat dijalankan**.
*(User story 31–36 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Rekam lampiran; resolusi alamat penyimpanan; **cache token** |
| `internal/clients` | Klien penyimpanan berkas **di balik interface** |
| `internal/services` | Orkestrasi efek keluar; penandaan kegagalan; pengulangan |
| `internal/handlers` | Endpoint unggah / unduh / hapus / status |
| `frontend/` | Panel lampiran di layar tahun treaty; penanda "terkirim" / "tertunda" |

## Rule Pega sumber

`[terverifikasi]` Rantai berkas yang ada sekarang — pola yang sama dengan Master Product Name Life:

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `ServiceGoogle` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `SERVICEGOOGLE` / `RULE-CONNECT-REST` | `Treaty Contract Out/ConnectREST/ServiceGoogle.xml` | **satu-satunya** ConnectREST modul ini |
| `LinkService` | `LINKSERVICE` / `LINKSERVICE` / `RULE-ADMIN-SYSTEM-SETTINGS` | `Treaty Contract Out/SystemSettings/LinkService.xml` | sumber alamat |
| `GetTokenStorage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetTokenStorage_SQL.xml` | → `POOLDATA.GET_TOKEN_STORAGE` |
| `GetLinkStorage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetLinkStorage_SQL.xml` | ambil tautan berkas |
| `Update_T_Storage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/Update_T_Storage_SQL.xml` | perbarui rekam |
| `DeleteStorage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteStorage_SQL.xml` | hapus rekam |
| `InsertAtatchment_Sql` | `ASM-FW-GISFW-INT-TREATY_IN` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/InsertAtatchment_Sql.xml` | → `POOLDATA.PEGA_M_ATTACHMENT` (CLOB) |
| `GetAllAttachment2_Sql` | `ASM-FW-GISFW-INT-TREATY_IN` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetAllAttachment2_Sql.xml` | ⚠️ berkunci `TreatyIn.ID` |
| `CategoryAttach_SQL` | `ASM-FW-GISFW-INT` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/CategoryAttach_SQL.xml` | → `CATEGORY_ATTACH_REAS` |
| `GetLinkService` | `@BASECLASS` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/GetLinkService.xml` | resolusi alamat |
| `TreatyOutSaveAttachment`, `TreatyOutDownloadAll_Act`, `TreatyOutDownloadOne`, `LoadAttachmentTreatyOut` | `Treaty Contract Out/Activity/` | jalur berkas yang belum rampung |

`[terverifikasi]` `ServiceGoogle` **tanpa URL literal** — alamat datang dari `LinkService`.
⚠️ `pzOriginalInstanceKey`-nya menunjuk `…GOOGLESTORAGE_UPLOAD` — nama berkas ≠ nama asal
(**OQ-066** lagi).

⚠️ **Kebocoran batas yang tidak dibawa.** `[terverifikasi]` `GetAllAttachment2_Sql` berclass
**`ASM-FW-GISFW-INT-TREATY_IN`** dan membaca `M_ATTACHMENTTREATY_2 WHERE treatyid = {TreatyIn.ID}` —
lampiran existing dikunci ke **ID treaty inward**, bukan ke tahun treaty. `[keputusan work owner]`
Di sistem baru lampiran melekat pada **tahun treaty**; kepemilikan `M_ATTACHMENTTREATY_2` tetap di
konteks treaty inward.

⚠️ **Penyimpangan sadar 9 — fitur baru.** `[keputusan work owner]`

## ADR terkait

**ADR-0013** (alamat di-resolve **runtime**; **dilarang** sebagai literal, konstanta, **maupun env
var** — menggantikan ADR-0004), **ADR-0015** (efek keluar — kegagalan tidak boleh diam-diam),
**ADR-0010** (penyimpanan berkas tetap Google Storage), **ADR-0007** (jejak audit).

## Acceptance criteria

- [ ] ⚠️ Berkas dapat **dilampirkan pada sebuah tahun treaty** — fitur yang **belum ada** di Pega.
      *(AC 54 spec; User story 31; penyimpangan sadar 9)*
- [ ] ⚠️ Lampiran **opsional**: kegagalan unggah **tidak** membatalkan tersimpannya tahun treaty —
      hanya status lampiran yang berubah. *(AC 55 spec; User story 32)*
- [ ] Tiap lampiran dapat diberi **kategori** dari master kategori. *(AC 56 spec; User story 33)*
- [ ] Lampiran dapat **diunduh** dan **dihapus**. *(AC 57 spec; User story 34)*
- [ ] ⚠️ Kegagalan unggah **tercatat dan dapat diulang**; pengulangan **tidak** menggandakan berkas.
      *(AC 58 spec; User story 35; **ADR-0015**)*
- [ ] Alamat penyimpanan di-resolve **runtime** dari konfigurasi. Test yang memindai kode untuk URL
      sebagai **literal, konstanta, atau pembacaan env var** **gagal** bila menemukannya.
      *(AC 59 spec; User story 36; **ADR-0013**)*
- [ ] Token **di-cache** dan **diperbarui sebelum kedaluwarsa**; token yang **gagal diambil**
      menghasilkan kegagalan yang **terlihat**, bukan unggahan yang diam saja tidak terjadi.
      *(AC 60 spec)*
- [ ] Rekam berkas di basis data dan berkas di penyimpanan **tetap sejalan**: rekam tanpa berkas
      terdeteksi dan dapat diperbaiki. *(AC 61 spec)*
- [ ] Klien penyimpanan berkas berada **di balik interface** dan **di-fake** di test; yang diperiksa
      adalah **efeknya**. *(AC 62 spec)*
- [ ] Penghapusan berkas yang **sudah tidak ada** di penyimpanan **tidak** menggagalkan penghapusan
      rekamnya.
- [ ] Lampiran melekat pada **tahun treaty**, bukan pada ID treaty inward. Test yang menemukan
      rujukan ke ID treaty inward **gagal**.
- [ ] Setiap unggah dan penghapusan mencatat **jejak audit**. *(**ADR-0007**)*

## Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka]` **OQ-047 — tidak memblokir:** alamat fisik Google Storage tersimpan di
`M_LINK_SERVICE`, dan isinya belum dibaca. **ADR-0013** sudah mengatur **cara** meresolusinya, jadi
tiket ini dapat selesai tanpa mengetahui alamatnya — yang mengikat adalah **alamat tidak ditanam**.

## Catatan

⚠️ **Taruhannya berbeda dari Komite Claim Life.** Di sana efek keluar **wajib berhasil**
(transactional outbox, **ADR-0015** versi Komite). Di sini `[keputusan work owner]` lampiran
**opsional** — yang wajib adalah **kegagalannya terlihat dan dapat diulang**, bukan berhasil.
Sama dengan perlakuan di Master Product Name Life.

⚠️ **Fitur ini tidak lengkap di Pega** `[keputusan work owner]` — jadi korpus bukan sumber
kebenaran perilaku di sini, hanya sumber **rantai teknis** (token, alamat, rekam berkas, kategori).
Perilaku yang tidak ada di korpus **tidak ditebak**; ia ditetapkan oleh AC di atas.

⚠️ **Satu-satunya efek keluar konteks ini.** `[terverifikasi]` Modul ini punya **satu**
`RULE-CONNECT-REST`, dan tidak ada integrasi Arasapas / Kasir / Konversi / Gemini AI. Berbeda dari
Master Contract Retro Life yang **tanpa `ConnectREST` sama sekali** — di sana ADR-0013 dan ADR-0015
tidak berlaku; di sini berlaku keduanya.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**; **penyimpanan berkas di-fake** di balik
interface — ia satu-satunya yang dipalsukan di seluruh konteks ini.

```
go test ./internal/...
cd frontend && npm test
make check
```
