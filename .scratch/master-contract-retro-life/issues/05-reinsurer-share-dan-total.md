# 05: Reinsurer — share, komisi, dan total yang ditampilkan tanpa memblokir

**Status:** ready-for-agent

**Blocked by:** 02 (reinsurer lahir di bawah kontrak), 04 (jenis reasuransi diwarisi dari kontrak)

## Hasil & nilai pengguna

Sebagai **admin master**, saya menambahkan **reinsurer** ke sebuah kontrak beserta persentase share,
komisi, dan overriding commission-nya; dan sebagai **underwriter** saya melihat **total share** yang
sudah teralokasi — tanpa dihalangi menyimpan ketika kontrak masih saya susun bertahap.
*(User story 12–17 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas reinsurer; share/komisi sebagai desimal |
| `internal/repository` | Pemanggilan procedure penulis reinsurer; kueri total share per kontrak |
| `internal/services` | Validasi 0–100 per baris; penjumlahan total share |
| `internal/handlers` | Endpoint CRUD reinsurer; total share pada respons kontrak |
| `frontend/` | Grid reinsurer; **total share mencolok** + tanda "belum 100%" |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterTreatyReinsurer_Life_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!SAVEMASTERTREATYREINSURER_LIFE_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/SaveMasterTreatyReinsurer_Life_SQL.xml` | `POOLDATA.INSERTREINSURER_LIFE` |
| `SaveSecurityLife_Act` | `ASM-FW-GISFW-…` / `SAVESECURITYLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveSecurityLife_Act.xml` | orkestrator simpan |
| `NewInputSecurityLife_Act` | `ASM-FW-GISFW-…` / `NEWINPUTSECURITYLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/NewInputSecurityLife_Act.xml` | baris baru |
| `SetSecurityLife_Act` | `ASM-FW-GISFW-…` / `SETSECURITYLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SetSecurityLife_Act.xml` | isi form |
| `SetErrorMessageReinsurer` | `@BASECLASS` / `SETERRORMESSAGEREINSURER` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SetErrorMessageReinsurer.xml` | **validasi 0–100** |
| `CountingPercentShare_Act` | `@BASECLASS` / `COUNTINGPERCENTSHARE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/CountingPercentShare_Act.xml` | ⚠️ **penjumlah**, bukan pembagi |
| `GetMasterReinsurerLifeList_SQl` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!GETMASTERREINSURERLIFELIST_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/GetMasterReinsurerLifeList_SQl.xml` | sumber total share |
| `BrowseDetailTreatyReisurerLife_RD` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `BROWSEDETAILTREATYREISURERLIFE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseDetailTreatyReisurerLife_RD.xml` | daftar |

`[terverifikasi]` **Validasi per baris** (`SetErrorMessageReinsurer`):

```
@toDecimal(PctShare) > 100 || @toDecimal(PctShare) < 0   → galat
@toDecimal(Ricomm)   > 100 || @toDecimal(Ricomm)   < 0   → galat
```

`[terverifikasi]` **Total share** dihitung `CountingPercentShare_Act`:
`Local.TotalShare += @toDecimal(.PCTSHARE)` atas hasil
`SELECT PCTSHARE FROM POOLDATA.TREATYREINSURER_LIFE WHERE TREATYYEARID = … AND TREATYCONTRACTID = …`,
lalu hasilnya ditaruh di field **`STDRATING`** untuk **ditampilkan**.
⚠️ **Tidak ada perbandingan terhadap 100** di seluruh modul, dan `[data DBA]` **tidak ada pula di
procedure**.

`[data DBA]` Kolom: `PCTSHARE`, `COMMISION`, `OVR_COMM` bertipe **`NUMBER`** — desimal, bukan teks.

## ADR terkait

**ADR-0003** (uang & persentase non-float), **ADR-0006** (identitas dari basis data), **ADR-0007**.

## Acceptance criteria

- [ ] Reinsurer lahir **di bawah** satu kontrak; tanpa kontrak induk **ditolak**. *(AC 15 spec)*
- [ ] `PCTSHARE` dan komisi di luar rentang **0–100 ditolak**, dengan pesan yang **menyebut
      kolomnya**. *(AC 16 spec)*
- [ ] Total share seluruh reinsurer pada satu kontrak **dihitung dan ditampilkan**. *(AC 17 spec)*
- [ ] ⚠️ Total share **TIDAK memblokir penyimpanan** — kontrak dengan total ≠ 100% **tetap
      tersimpan**. *(AC 18 spec; `[keputusan work owner]`)*
- [ ] Kontrak dengan total ≠ 100% **ditandai mencolok** di layar. *(AC 19 spec)*
- [ ] Share dan komisi diperlakukan sebagai **desimal presisi arbitrer**; **tidak** melewati `float`.
      *(AC 49 spec; **ADR-0003**)*
- [ ] `REINSTYPEID` **tidak ditulis** pada baris reinsurer — diwarisi dari kontrak. *(AC 41 spec;
      tiket 04)*
- [ ] Menyimpan reinsurer yang sudah ada = **upsert**, bukan baris kedua; identitas baru dibuat basis
      data. *(AC 45, 4 spec)*
- [ ] `o_message` diperiksa; kegagalan ditampilkan. *(AC 46 spec — HTML dan konformansi di
      **tiket 10**)*

## Blocker

**Tidak ada pemblokir.**

⚠️ **Keputusan tiket yang belum diambil — anti-dobel logis.** `[data DBA]` Upsert dikunci `ID`
mencegah dua baris ber-`ID` sama, **tidak** mencegah **dua reinsurer yang sama pada satu kontrak**.
Basis data tidak menjaganya (**nol `UNIQUE`**), dan Pega pun tidak. Bila anti-dobel dikehendaki, itu
aturan Go dan idealnya `UNIQUE` di basis data. **Putuskan sebelum tiket ini ditutup.**

## Catatan

⚠️ **Nama menyesatkan — dua di tiket ini.** `[terverifikasi]`

| Yang tertulis | Yang sebenarnya |
| --- | --- |
| `CountingPercentShare_Act` "rumus pembagian share retro" | hanya **menjumlahkan** `PCTSHARE` |
| field `STDRATING` "standard rating" | menyimpan **total share** |

Di sistem baru, **beri nama yang jujur** — jangan bawa `STDRATING` sebagai tempat total share.

⚠️ `COMMISION` dieja demikian di basis data (satu `S`). Bawa nama kolomnya apa adanya di lapisan
repository; gunakan ejaan yang benar di lapisan domain.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```
