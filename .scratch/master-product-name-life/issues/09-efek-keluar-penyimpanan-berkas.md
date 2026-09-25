# 09: Efek keluar penyimpanan berkas — resolusi alamat, token, dan pengulangan

**Status:** ready-for-agent

**Blocked by:** 08 (lampiran harus ada lebih dulu — tiket ini mengeraskan jalur keluarnya)

## Hasil & nilai pengguna

Sebagai **tim operasi**, saya ingin alamat penyimpanan berkas **dibaca dari konfigurasi saat
dijalankan** — bukan ditanam di kode — dan saya ingin unggahan yang gagal **tercatat dan dapat
diulang**, sehingga perpindahan lingkungan tidak menuntut penempelan ulang dan tidak ada berkas yang
hilang diam-diam. *(User story 28 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Resolusi alamat penyimpanan; pengambilan dan **cache token** |
| `internal/clients` | Klien penyimpanan berkas **di balik interface** |
| `internal/services` | Orkestrasi efek keluar; penandaan kegagalan; pengulangan |
| `internal/handlers` | Status efek keluar terbaca API |
| `frontend/` | Penanda "terkirim" / "tertunda" pada lampiran |

## Rule Pega sumber

`[terverifikasi]` Seluruh rantai berkas ber-class **`ASM-FW-GISFW-INT-T_STORAGE_IMAGE`**:

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `ServiceGoogle` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `SERVICEGOOGLE` / `RULE-CONNECT-REST` | `Master Product Name Life/ConnectREST/ServiceGoogle.xml` | **satu-satunya ConnectREST** modul ini |
| `GetTokenStorage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!GETTOKENSTORAGE_SQL` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/GetTokenStorage_SQL.xml` | ambil token |
| `GetLinkStorage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!GETLINKSTORAGE_SQL` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/GetLinkStorage_SQL.xml` | ambil tautan |
| `Insert_T_Storage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!INSERT_T_STORAGE_SQL` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/Insert_T_Storage_SQL.xml` | rekam berkas |
| `Update_T_Storage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!UPDATE_T_STORAGE_SQL` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/Update_T_Storage_SQL.xml` | perbarui rekam |
| `GetLinkService` | `ASM-FW-GISFW-…` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/GetLinkService.xml` | **resolusi alamat** |
| `InsertGoogleStorage_Act` | `ASM-FW-GISFW-…` / `INSERTGOOGLESTORAGE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/InsertGoogleStorage_Act.xml` | unggah |
| `GetUrlGoogleStorage_Act` | `ASM-FW-GISFW-…` / `GETURLGOOGLESTORAGE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/GetUrlGoogleStorage_Act.xml` | ambil URL |
| `DeleteGoogleStorage_Act` | `ASM-FW-GISFW-…` / `DELETEGOOGLESTORAGE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/DeleteGoogleStorage_Act.xml` | hapus berkas |

`[terverifikasi]` `ServiceGoogle` **tanpa URL literal** — alamat datang dari `M_LINK_SERVICE`.

`[data DBA]` Token mengikuti pola `GET_TOKEN_STORAGE`: **MD5**, **di-cache**, **berlaku 1 menit**.

## ADR terkait

**ADR-0013** (alamat di-resolve **runtime**; **dilarang** sebagai literal, konstanta, **maupun env
var** — menggantikan ADR-0004), **ADR-0015** (efek keluar — kegagalan tidak boleh diam-diam),
**ADR-0010** (penyimpanan berkas tetap Google Storage), **ADR-0005** (flag lingkungan bila
diperlukan).

## Acceptance criteria

- [ ] Alamat penyimpanan di-resolve **runtime** dari konfigurasi. Test yang memindai kode untuk URL
      sebagai **literal, konstanta, atau pembacaan env var** **gagal** bila menemukannya.
      *(AC 37 spec; **ADR-0013**)*
- [ ] Token **di-cache** dan **diperbarui sebelum kedaluwarsa**; kedaluwarsanya **tidak** menggagalkan
      permintaan pengguna secara langsung. *(AC 39 spec; `[data DBA]` berlaku 1 menit)*
- [ ] Token yang **gagal diambil** menghasilkan kegagalan yang **terlihat**, bukan unggahan yang diam
      saja tidak terjadi.
- [ ] ⚠️ Kegagalan unggah **tercatat dan dapat diulang**; pengulangan **tidak** menggandakan berkas.
      *(AC 35 spec; **ADR-0015**)*
- [ ] Klien penyimpanan berkas berada **di balik interface** dan **di-fake** di test; yang diperiksa
      adalah **efeknya**.
- [ ] ⚠️ Kegagalan efek keluar **tidak** membatalkan produk maupun rekam lampirannya — hanya status
      lampiran yang berubah. *(AC 34 spec)*
- [ ] Rekam berkas di basis data dan berkas di penyimpanan **tetap sejalan**: rekam tanpa berkas
      terdeteksi dan dapat diperbaiki.
- [ ] Penghapusan berkas yang **sudah tidak ada** di penyimpanan **tidak** menggagalkan penghapusan
      rekamnya.

## Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka]` **OQ-047 — tidak memblokir:** alamat fisik Google Storage tersimpan di
`M_LINK_SERVICE`, dan isinya belum dibaca. **ADR-0013** sudah mengatur **cara** meresolusinya, jadi
tiket ini dapat selesai tanpa mengetahui alamatnya — yang mengikat adalah **alamat tidak ditanam**.

## Catatan

⚠️ **Berbeda dari Master Contract Retro Life.** Konteks itu **tanpa `ConnectREST` sama sekali**,
sehingga ADR-0013 dan ADR-0015 tidak berlaku di sana. Di sini **berlaku keduanya** — inilah
satu-satunya efek keluar konteks ini.

⚠️ **Jangan samakan dengan ADR-0015 versi Komite.** Di Komite Claim Life, efek keluar **wajib
berhasil** (transactional outbox). Di sini `[keputusan work owner]` lampiran **opsional** — yang
wajib adalah **kegagalannya terlihat dan dapat diulang**, bukan berhasil. Taruhannya berbeda.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**; **penyimpanan berkas di-fake** di balik
interface.

```
go test ./internal/...
cd frontend && npm test
make check
```
