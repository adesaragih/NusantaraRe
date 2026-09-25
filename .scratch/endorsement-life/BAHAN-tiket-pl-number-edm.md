> ⚠️ **BAHAN, bukan tiket aktif.** Berkas ini semula `07-endorsement-life-pl-number-edm.md` di
> `.scratch/premiumlist-life/issues/`. Dipindah 2026-09-15 setelah `[keputusan work owner]`
> menetapkan **Endorsement Life sebagai konteks terpisah**. Isinya dipertahankan apa adanya sebagai
> bahan untuk spec dan tiket konteks Endorsement Life; **jangan dikerjakan sebagai tiket** sebelum
> grilling konteks ini selesai dan spec-nya terbit.

# 07: Endorsement Life — `PL_NUMBER_EDM` sejajar dan penulisan detail endorsement

**Status:** (dipindah — bahan konteks Endorsement Life)

**Blocked by:** 05a (jalur simpan summary — endorsement memakai jalur yang sama), 03 (bentuk premium
list detail)

## Hasil & nilai pengguna

Sebagai **inputor Life**, saya ingin membuat endorsement atas polis berjalan dan memperoleh
**nomor endorsement tersendiri** yang tidak menimpa nomor new business, supaya perubahan polis punya
jejak sendiri dan premium list asalnya tetap utuh dan dapat dirujuk. *(User story 25–28 di spec)*

## Area codebase

`internal/handlers` (endpoint buat endorsement + submit), `internal/services` (penomoran EDM
idempoten; penentuan field yang boleh diubah per form), `internal/repository` (penulisan detail +
JSON polis jalur endorsement), `frontend/` (layar Input EDM, penandaan baris baru/diubah).

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `InputEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE` / `RULE-OBJ-FLOW` | `Endorsement Life/Flow/InputEDMLife.xml` | titik masuk flow |
| `CreateCaseEMDL` | `DATA-PORTAL` / `CREATECASEEMDL` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/CreateCaseEMDL.xml` | buat case |
| `GenerateNoEDM_Life` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `GENERATENOEDM_LIFE` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/GenerateNoEDM_Life.xml` | **penomoran EDM, 7 langkah** |
| `Generate_NoEndorsmentLife` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!GENERATE_NOENDORSMENTLIFE` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/Generate_NoEndorsmentLife.xml` | rakit nomor |
| `GetProdKeOldData_SQL` | `ASM-FW-GISFW-INT-OFFERJSON` / `ASM!GETPRODKEOLDDATA_SQL` / `RULE-CONNECT-SQL` | *(dirujuk `GenerateNoEDM_Life` step 3)* | ambil `PRODKE` terakhir |
| `InsertJsonPolisLife_Act` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml` (444.634 byte, **16 langkah**) | orkestrator simpan |
| `SaveMasterLPDet` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/SaveMasterLPDet.xml` | `INSERT INTO POOLDATA.M_LIFE_PREMIUM_DETAIL` + `COMMIT;` (252) |
| `InsertJsonPolisEDM` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONPOLISEDM` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/InsertJsonPolisEDM.xml` | `INSERT INTO POOLDATA.JSON_POLIS (…, NOENDORS, PRODKE, …)` + `COMMIT;` (107) |
| `GetEdmTypeLife`, `GetPL_NumberLife`, `GetProdkeNopolis` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!…` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/` | baca data polis lama dari `JSON_POLIS` |
| `MappingEDMLife`, `GetOldDetail_EDM`, `SetPremi_EDM`, `SaveCSVEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `…` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/` | salin & ubah detail lama |

### `PL_NUMBER_EDM` bukan dari `PROC_GENERATE_SEQUENCE_NUMBER`

⚠️ `[terverifikasi]` **Ini berbeda dari `PL_NUMBER`.** `GenerateNoEDM_Life` merakitnya sendiri:

| Step | Baris | Isi |
| ---: | ---: | --- |
| 3 | 680 | `RDB-List` "Get ProdKe dari JSON_POLIS" → `GetProdKeOldData_SQL` |
| 4 | 872 | "Set Prodke" — `Local.Prodke+1` (946) → `InputData.CARI14` (966) |
| **5** | 1068 | `RDB-List` "Generate No Endorsement" → `Generate_NoEndorsmentLife` |
| **6** | 1264 | "Set No Endorsement" — `PL_NUMBER_EDM = Local.Nopolis` (1337) |
| 7 | 1430 | `Obj-Save` |

`[terverifikasi]` SQL `Generate_NoEndorsmentLife`:

```sql
SELECT NOPOLIS||'/'||{InputData.CARI14} AS HASIL1
FROM POOLDATA.JSON_POLIS
WHERE NOPOLIS = {pyWorkPage.PolicyNo}
ORDER BY PRODKE DESC
```

Jadi **`PL_NUMBER_EDM` = `<nomor polis>` + `/` + `<PRODKE terakhir + 1>`** — diturunkan dari data
polis, bukan dari sequence terpusat. **ADR-0006 tidak berlaku pada nomor ini**; jangan memaksanya
lewat `PROC_GENERATE_SEQUENCE_NUMBER`.

`[terverifikasi]` **Idempoten di korpus**: step 5 dan step 6 keduanya dijaga precondition
`pyWorkPage.PremiumListSummary.PL_NUMBER_EDM==""` (baris 1216 dan 1390). Nomor **lahir sekali**.

`[keputusan work owner]` `PL_NUMBER` (NB) dan `PL_NUMBER_EDM` (EDM) adalah **dua jalur penomoran yang
sejajar**, dibuat terpisah agar tidak saling timpa — **bukan** relasi induk-anak, bukan pointer ke PL
asal.

`[keputusan work owner]` Field yang boleh diubah endorsement ditentukan kondisi `When` **per form**;
**tidak ada** daftar global.

`[terverifikasi]` `InsertPLSummary` yang dipakai endorsement (`<RequestType>` baris 6226) adalah
**rule yang sama persis** dengan yang dipakai PremiumList Life — lihat tiket **05a**.

## ADR terkait

**ADR-0006** (berlaku untuk `PL_NUMBER`, **tidak** untuk `PL_NUMBER_EDM` — perbedaan ini wajib
terlihat di kode), **ADR-0003** (uang non-float), **ADR-0015** (batas transaksi), **ADR-0007**
(jejak audit), **ADR-0001** (batas konteks).

## Acceptance criteria

- [ ] Endorsement memperoleh `PL_NUMBER_EDM`; `PL_NUMBER` new business **tidak berubah** dan premium
      list asal tetap dapat dibaca utuh. *(AC 12 spec)*
- [ ] `PL_NUMBER_EDM` dirakit sebagai `<nomor polis>/<prodke+1>`, dengan `prodke` diambil dari rekam
      polis terakhir — **bukan** dari `PROC_GENERATE_SEQUENCE_NUMBER`.
- [ ] Penomoran **idempoten**: memanggil ulang pada endorsement yang sudah bernomor **tidak**
      mengubah nomornya dan tidak menaikkan `prodke`.
- [ ] Dua endorsement berurutan atas polis yang sama memperoleh `prodke` berurutan.
- [ ] Kedua nomor tersimpan pada rekam summary yang sama sebagai **dua nilai terpisah**.
      *(AC 13 spec)*
- [ ] Baris detail endorsement ditulis ke `M_LIFE_PREMIUM_DETAIL` membawa **`PL_NUMBER` maupun
      `PL_NUMBER_EDM`**; kolom uangnya tidak melewati `float`. *(AC 14 spec; **ADR-0003**)*
- [ ] Penulisan detail dan JSON polis endorsement **commit sendiri** — diperlakukan sebagai titik
      potong, persis seperti tiket 05b, dan dapat diulang tanpa menggandakan baris.
- [ ] Field yang dapat diubah ditentukan **per form**; tidak ada daftar global yang ditanam di kode.
- [ ] Baris yang **baru** ditambahkan endorsement dibedakan dari baris yang **disalin** dari polis
      lama, dan pembedaan itu terlihat pengguna.
- [ ] Endorsement atas polis yang tidak ditemukan ditolak dengan pesan yang menyebut nomor polisnya.

## Blocker

**Tidak ada pemblokir.**

## Catatan — asimetri yang tidak dibawa

⚠️ `[terverifikasi]` Jalur endorsement di korpus **kehilangan alarmnya**: step **13** "Cek sudah masuk
atau blm datanya" (baris 6399) dan step **15** `SendEmailNotification` (6752) keduanya
`<pyStepsBlockName>//`. Hanya step **16** `serviceInsertArasapasLife_act` yang aktif (dijaga
`IsPEGAPROD`, baris 7230). Sistem baru **menyeragamkan** — lihat tiket **06**.

⚠️ `[terverifikasi]` Endorsement menulis `JSON_POLIS` lewat `INSERT` **langsung**
(`InsertJsonPolisEDM`), sedangkan new business lewat **upsert** procedure `INSERTJSONPOLISLIFE`.
Sifat "dapat diulang" karena itu **tidak gratis** di jalur endorsement: pengulangan menghasilkan
baris kedua kecuali ditangani eksplisit. **AC di atas mensyaratkan penanganan itu.**

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```
