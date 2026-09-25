# 12: Efek keluar asinkron + antre-ulang + flag lingkungan

**Status:** ready-for-agent

**Blocked by:** 09 (jejak audit) — kegagalan efek keluar harus masuk jalur audit, bukan hanya log

## Hasil & nilai pengguna

Sebagai **ReasLifeAdmin**, saya dapat mengunggah dan mengunduh dokumen pendukung klaim; sebagai
**ReasLifeSPV**, anggota Komite menerima notifikasi saat kasus diserahkan. Dan sebagai **pengguna
mana pun**, alur klaim saya **tidak pernah tertahan** karena layanan luar sedang gagal — kegagalan
itu tercatat dan diantre ulang, bukan hilang diam-diam.
*(User story 32–36 di spec)*

## Area codebase

`internal/services` (interface efek keluar + orkestrasi asinkron + antre-ulang),
`internal/repository` (implementasi pemanggil), `internal/config` (alamat layanan + flag
lingkungan), `internal/handlers` (endpoint unggah/unduh), `frontend/` (kontrol dokumen).

## Rule Pega sumber

| Efek | Rule | Identitas |
| --- | --- | --- |
| Unggah berkas | `Claim Life/Activity/InsertGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `INSERTGOOGLESTORAGE_ACT` / `RULE-OBJ-ACTIVITY`, 160.027 byte |
| Ambil URL / hapus | `GetUrlGoogleStorage_Act.xml`, `DeleteGoogleStorage_Act.xml` | kelas sama |
| Token penyimpanan | `Claim Life/RDBList/GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!GETTOKENSTORAGE_SQL` / `RULE-CONNECT-SQL` |
| Email | `Claim Life/Activity/SendEmailKlaimLF.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `SENDEMAILKLAIMLF` / `RULE-OBJ-ACTIVITY` |
| Arasapas | `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` | `[terverifikasi]` **satu-satunya salinan di korpus** — OQ-035 |
| ⚠️ ~~Konversi ke produksi lewat payload JSON~~ — **DIBUANG 2026-09-16** | ~~`Claim Life/ConnectREST/convertJsonNusareToProductionClaimLife.xml`~~ | `[keputusan work owner]` hilir **`SELECT` langsung dari tabel klaim** (spec §2b, §12) |
| Pencatatan layanan | `Claim Life/RDBList/InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` / `RNM!INSERTLOGSERVICECLAIM` / `RULE-CONNECT-SQL` → `INSERT INTO pooldata.monitoring_klaim_log` |

`[terverifikasi]` Token penyimpanan **tidak** diterbitkan aplikasi — ia datang dari Oracle:
`pooldata.GET_TOKEN_STORAGE({UploadDoc.App}, {OperatorID.pyUserIdentifier}, {UploadDoc.Kodestring OUT}, {DocAPI.ResponseMsg OUT})`.

## ADR terkait

**ADR-0008** (asinkron, tidak memblokir, antre-ulang), **ADR-0010** (penyimpanan tetap Google
Storage), **ADR-0013** (**resolusi endpoint lewat `M_LINK_SERVICE`** — meralat **ADR-0004** untuk
alamat layanan keluar), **ADR-0005** (flag lingkungan), **ADR-0007** (kegagalan masuk jalur audit).

## Acceptance criteria

- [ ] ⚠️ **Diselaraskan 2026-09-16:** efek keluar tinggal **tiga** — berkas, email, Arasapas.
      Konversi ke produksi **tidak lagi mengirim payload JSON**; hilir membaca **langsung dari
      tabel klaim**. Test yang menemukan payload JSON dikirim keluar **gagal**. *(AC 55 spec;
      penyimpangan sadar 1)*
- [ ] Seluruh efek keluar berada **di balik interface**, sehingga dapat diganti dalam test.
- [ ] Kegagalan efek keluar mana pun **tidak menahan** transisi status klaim. *(AC 19 spec)*
- [ ] Kegagalan tercatat di **jalur audit**, bukan hanya di log layanan, dan **dapat diantre ulang**.
      *(AC 20 spec)*
- [ ] Di lingkungan non-production, klaim **tetap tersimpan**; keempat efek keluar tidak berjalan.
      *(AC 21 spec)*
- [ ] Tidak ada host, endpoint, atau kredensial sebagai literal di kode.
- [ ] Kegagalan konfigurasi dapat dibedakan dari kegagalan jaringan, agar antre-ulang tidak berputar
      sia-sia.
- [ ] Dokumen yang sudah diunggah dapat diunduh kembali.
- [ ] Alamat endpoint keluar di-resolve lewat **runtime lookup** ke `M_LINK_SERVICE` dengan kunci
      `(KATEGORI_1, KATEGORI_2)` — untuk Arasapas: `("Klaim", "insertClaimLife")`.
- [ ] **Tidak ada URL** sebagai literal, konstanta, **maupun env var** di kode. Yang boleh menjadi
      konstanta hanyalah **kunci kategori**. *(**ADR-0013**)*
- [ ] Pemisahan dev–prod terjadi lewat **isi tabel per-database**, bukan lewat percabangan di kode.
- [ ] Bila kunci kategori tidak ditemukan di `M_LINK_SERVICE`, kegagalan **terang-terangan** dan
      masuk jalur audit — bukan diam-diam melewati efek keluar.

## Catatan penutupan (2026-09-14)

**OQ-047 TERTUTUP** `[terverifikasi + data DBA]` — kontrak resolusi endpoint diketahui penuh:
`Claim Life/Activity/GetLinkService.xml` (`ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` /
`RULE-OBJ-ACTIVITY`) melakukan `Obj-Browse` atas `M_LINK_SERVICE` dengan
`.KATEGORI_1 = Param.Kategori_1 AND .KATEGORI_2 = Param.Kategori_2`, mengambil `.URL`, lalu
`Connect-REST`.

`[terverifikasi]` Kunci Claim — Life terbaca di
`Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml`:
`Kategori_1 = "Klaim"`, `Kategori_2 = "insertClaimLife"`.
`[data DBA]` Isi tabel 19 baris; endpoint Life →
`http://10.100.10.75:7315/Nusare-Integration-WS/resources1/restws/NusareClaim/insertClaimLife`.

**OQ-018 TERTUTUP untuk Claim — Life** `[terverifikasi]` — **nol** URL endpoint bisnis ter-hardcode
di modul ini; seluruh `http(s)://` yang ada hanyalah `pyHelpURI` ke `community.pega.com` (132+) dan
3 tautan penampil dokumen Office. Pembedaan lingkungan memakai
`Claim Life/When/IsPEGAPROD.xml` (`@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN`):
`pzProductionLevel = "5"`.

**Flag lingkungan — keputusan spesifikasi tetap berlaku:** flag hanya menggerbangi **efek keluar**,
tidak pernah **penyimpanan**. Ini penyimpangan sadar dari Pega (di sana `IsPEGAPROD` juga
menggerbangi simpan utama); tanpa itu lingkungan non-production tidak dapat dipakai menguji.

`[terbuka]` **OQ-035** — kepemilikan `serviceInsertArasapasClaimLife_act`. **Tidak memblokir.**

## Catatan

`[terverifikasi]` "Tidak memblokir" adalah **paritas**, bukan perubahan: `Claim Life` hanya punya
**pencatatan** (`InsertLogServiceClaim`), bukan gerbang keberhasilan seperti `IsSuccessHitService`
di konteks facultative. Yang **ditambahkan** adalah antre-ulang.

⚠️ Konsekuensi yang diterima (**ADR-0008**): sebuah klaim dapat mencapai keadaan akhir sementara
efek keluarnya masih tertunda.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
