# 01: Pilih polis new business + lima gerbang kelayakan endorsement

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, **CL-01** (kerangka aplikasi + seam API — scaffolding lintas konteks; tidak dibuat di
sini)

## Hasil & nilai pengguna

Sebagai **inputor Life**, saya memasukkan **nomor polis** yang ingin saya endorse dan sistem langsung
memberi tahu apakah polis itu **boleh** di-endorse — lengkap dengan alasannya bila tidak — sehingga
saya tidak pernah membuat endorsement yang mustahil diselesaikan dan tidak ada case menggantung.
*(User story 1–8 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Permintaan kelayakan endorsement; hasil kelayakan beserta daftar alasan penolakan |
| `internal/repository` | Lookup polis di `JSON_POLIS`; pencarian endorsement berjalan; baca `EdmType` polis terakhir; **satu repository terpisah** untuk pembacaan Arasapas |
| `internal/services` | Mesin lima gerbang — berjalan berurutan, mengumpulkan seluruh alasan |
| `internal/handlers` | Endpoint cek kelayakan; endpoint pencarian polis |
| `frontend/` | Layar Input EDM: isian nomor polis, penelusuran polis, tampilan alasan penolakan |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SetErrorBatalEndorsement_Act` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETERRORBATALENDORSEMENT_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/SetErrorBatalEndorsement_Act.xml` (141.275 byte) | **gerbang kelayakan** |
| `GetPL_NumberLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETPL_NUMBERLIFE` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/GetPL_NumberLife.xml` | cari polis di `JSON_POLIS` |
| `FilterProteksiEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `FILTERPROTEKSIEDMLIFE` / `RULE-OBJ-REPORT-DEFINITION` | `Endorsement Life/ReportDefinition/FilterProteksiEDMLife.xml` | cari EDM berjalan |
| `GetEdmTypeLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETEDMTYPELIFE` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/GetEdmTypeLife.xml` | baca `EdmType` dari CLOB |
| `SearcStatusBayarArasaps_SQL` | `ASM-FW-GISFW-INT-POLICYJSON` / `ASM!SEARCSTATUSBAYARARASAPS_SQL` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/SearcStatusBayarArasaps_SQL.xml` | **cek pembayaran Arasapas** |
| `BrowsePremiumList_RD` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `BROWSEPREMIUMLIST_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Endorsement Life/ReportDefinition/BrowsePremiumList_RD.xml` | **alat bantu pencarian** |
| `InputEDMLife` (Section) | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE` / `RULE-OBJ-HTML-SECTION` | `Endorsement Life/Section/InputEDMLife.xml` (1.255.106 byte) | layar |

`[terverifikasi]` **Peta lima gerbang** — sub-langkah di bawah step 3 `SetErrorBatalEndorsement_Act`:

| Sub-step | Precondition | Arti |
| --- | --- | --- |
| 3.1 | `param.Nopolis==""` | **nomor polis kosong** |
| 3.2 → 3.3 | `OutData.pxResults(1).CARI1==""` (lewat `GetPL_NumberLife`, baris 895) | **polis tidak ada di `JSON_POLIS`** |
| 3.5 → 3.6 | `@LengthOfPageList(ListEdm.pxResults)>0` (RD `FilterProteksiEDMLife`, baris 1211-1212) | **sudah ada EDM belum resolve** |
| 3.7 → 3.8 | `@contains(OutData1.pxResults(1).CARI1,"3")` (lewat `GetEdmTypeLife`, baris 1766) | **sudah pernah EDM batal** |
| 3.9 → 3.10 | `@LengthOfPageList(ListPembayaran.pxResults)>0 && TempWork.EdmType=="3"` (baris 2090) | **sudah ada pembayaran** |
| ~~4~~ | `Obj-Save` | **REMARK** (`<pyStepsBlockName>//`, baris 2562) |

`[terverifikasi]` Kunci di seluruh rule endorsement adalah **`TempWork.PolicyNo`** (`NOPOLIS`) —
bukan `PL_NUMBER`. `[keputusan work owner]` `BrowsePremiumList_RD` **alat bantu pencarian, bukan
sumber kunci**.

`[terverifikasi]` Kueri pembayaran:
`select * from ARASAPAS.DETAIL_INVOICE where inv_inv_no = {InputData.CARI18} and IVD_JR_ID = '5'`.
`[keputusan work owner]` **`IVD_JR_ID = '5'` = pembayaran/pelunasan.**

## ADR terkait

**ADR-0001** (batas konteks), **ADR-0007** (jejak audit), **ADR-0009** (migrasi penuh).

## Acceptance criteria

- [ ] Kelayakan dinilai dengan **nomor polis** sebagai kunci; penelusuran premium list hanya membantu
      menemukan nomor, **tidak** menjadi kunci. *(AC 1–2 US; §3 spec)*
- [ ] Ditolak bila **nomor polis kosong**, dengan pesan yang menyebutnya. *(AC 1 spec)*
- [ ] Ditolak bila **polis tidak ditemukan** di rekam polis. *(AC 2 spec)*
- [ ] Ditolak bila polis itu **masih punya endorsement yang belum selesai** — **satu endorsement
      terbuka per polis**. Ini **aturan integritas inti** konteks ini. *(AC 3 spec)*
- [ ] Ditolak bila polis **sudah pernah dibatalkan** (`EdmType` polis terakhir = `3`). *(AC 4 spec)*
- [ ] Endorsement **Batal** ditolak bila polis **sudah dibayar**; endorsement **Perubahan Data** atas
      polis yang sudah dibayar **tetap boleh** — gerbang ini hanya berlaku bila `EdmType=3`.
      *(AC 5 spec)*
- [ ] Kelima alasan penolakan **terbaca pengguna**, dan bila lebih dari satu gerbang gagal,
      **semuanya** dilaporkan — bukan hanya yang pertama. *(AC 8 US)*
- [ ] ⚠️ **Penyimpangan sadar — Arasapas terkurung.** Pembacaan `ARASAPAS.DETAIL_INVOICE` berada di
      **satu repository** yang ditandai **batas lintas sistem**. Test yang menemukan kueri skema
      `ARASAPAS` di luar repository itu **gagal**. *(AC 47 spec; `[keputusan desain]`)*
- [ ] ⚠️ **Rule gerbang diganti nama.** Tidak ada padanan bernama "SetErrorBatal…" di kode baru — ia
      **gerbang kelayakan endorsement**, bukan pembatalan. *(`[keputusan work owner]`)*
- [ ] **Tidak ada case endorsement yang dibuat** oleh tiket ini; kelayakan murni pemeriksaan.
      *(AC 6 spec — pembuatan case milik tiket 02)*

## Catatan pengujian

`[keputusan desain]` **Pembacaan pembayaran Arasapas TIDAK di-fake** — ia **gerbang bisnis**, bukan
sekadar integrasi. Diuji terhadap **skema uji nyata**. Yang di-fake hanya **kiriman keluar** Arasapas
dan email (tiket 10).

## Blocker

**Tidak ada.** Menunggu scaffolding **CL-01** yang sudah `ready-for-agent` di konteks Claim Life.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
