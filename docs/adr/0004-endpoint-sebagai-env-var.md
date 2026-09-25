---
status: superseded
digantikan-oleh: ADR-0013
tanggal: 2026-09-14
sumber: grilling Ronde 1 Q2, keputusan work owner — **dibatalkan 2026-09-14** setelah OQ-047 ditutup
---

> ⛔ **ADR INI SUDAH TIDAK BERLAKU.** Ia digantikan **ADR-0013**, yang memutuskan sebaliknya:
> alamat layanan keluar **wajib** di-resolve runtime dari `M_LINK_SERVICE`, dan **dilarang**
> ditanam sebagai env var. ADR-0004 ditulis saat isi `M_LINK_SERVICE` belum diketahui (OQ-047
> masih terbuka); `[data DBA]` isinya kini diketahui — 19 baris, dipakai bersama banyak modul.
> Teks di bawah dipertahankan untuk jejak audit.


# Alamat layanan keluar menjadi env var, bukan replikasi lookup `M_LINK_SERVICE`

Di Pega, alamat layanan keluar diambil **saat runtime dari tabel Oracle** `M_LINK_SERVICE`.
Di sistem baru, alamat menjadi **konfigurasi lingkungan (env var)**; mekanisme lookup itu
**tidak direplikasi**.

Keputusan ini terbatas pada **dari mana alamat berasal**. Apakah efek keluar dijalankan sama sekali
adalah keputusan terpisah — lihat **ADR-0005**.

## Keadaan sekarang `[terverifikasi]`

Dua rule `RULE-CONNECT-REST` di `Claim Life`, keduanya `pyBaseURLSelectionType = SETTING` dengan
`pyBaseURLSetting = LinkService!LinkService` — **tanpa URL literal**:

| Berkas | `pyServiceName` |
| --- | --- |
| `Claim Life/ConnectREST/ServiceGoogle.xml` | `ServiceGoogle` |
| `Claim Life/ConnectREST/convertJsonNusareToProductionClaimLife.xml` | `convertJsonNusareToProductionClaimLife` |

Alamat sesungguhnya diresolusi lewat `Activity/GetLinkService.xml`
(`ASM-FW-GISFW-INT-M_LINK_SERVICE!GETLINKSERVICE` / `RULE-OBJ-ACTIVITY`), yang melakukan
`Obj-Browse` atas class `ASM-FW-GISFW-Int-M_LINK_SERVICE` dengan kunci **`KATEGORI_1`** dan
**`KATEGORI_2`** — pola yang sama di seluruh korpus (`discovery/understanding-report.md` §3.1).

**Isi tabel `M_LINK_SERVICE` tidak ada di korpus** → **OQ-047**. Artinya daftar endpoint yang
sesungguhnya **tidak diketahui**, dari korpus maupun dari ADR ini.

## Considered Options

- **Env var per layanan** — dipilih
- Replikasi tabel `M_LINK_SERVICE` + lookup runtime — ditolak: memindahkan konfigurasi ke dalam
  data operasional membuat alamat tidak terbaca dari konfigurasi deployment, dan **isi tabelnya
  sendiri tidak diketahui** sehingga tidak dapat direplikasi dengan setia
- URL literal di kode — ditolak oleh aturan proyek (`CLAUDE.md` §10)

## Consequences

- Daftar endpoint menjadi bagian konfigurasi deployment, terbaca dan dapat berbeda per lingkungan.
- **Migrasi memerlukan isi `M_LINK_SERVICE` terlebih dahulu** (OQ-047) — tanpa itu, tidak diketahui
  berapa endpoint yang harus disediakan, atau apa arti `KATEGORI_1`/`KATEGORI_2`.
- Aturan proyek tetap berlaku: endpoint & host internal **bukan literal** di kode, dan tidak masuk
  prompt, tiket, atau test.

## OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-047** | Isi tabel `M_LINK_SERVICE` — daftar endpoint, arti `KATEGORI_1`/`KATEGORI_2`. Pemilik: DBA + Platform |
| **OQ-018** (terjawab untuk Claim — Life) | `jboss1073` = production, `jboss117` = dev; lingkungan `pega-nusre` **belum dinyatakan** |
