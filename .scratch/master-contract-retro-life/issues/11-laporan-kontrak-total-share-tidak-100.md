# 11: Laporan kontrak dengan total share ≠ 100%

**Status:** ready-for-agent

**Blocked by:** 05 (total share dihitung di sana)

## Hasil & nilai pengguna

Sebagai **manajemen**, saya ingin melihat **seluruh kontrak yang total sharenya belum 100%** dalam
satu tempat — sehingga celah alokasi risiko tidak tersembunyi di antara ratusan kontrak, meski
sistem sengaja tidak memblokir penyimpanannya. *(User story 17 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Kueri agregat total share per kontrak |
| `internal/services` | Penyaringan kontrak yang totalnya ≠ 100% |
| `internal/handlers` | Endpoint laporan |
| `frontend/` | Layar laporan; tiap baris dapat dibuka ke kontraknya |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `CountingPercentShare_Act` | `@BASECLASS` / `COUNTINGPERCENTSHARE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/CountingPercentShare_Act.xml` | penjumlah share per kontrak |
| `GetMasterReinsurerLifeList_SQl` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!GETMASTERREINSURERLIFELIST_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/GetMasterReinsurerLifeList_SQl.xml` | `SELECT PCTSHARE … WHERE TREATYYEARID = … AND TREATYCONTRACTID = …` |
| `BrowseTreatyContract_Life_RD` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` / `BROWSETREATYCONTRACT_LIFE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseTreatyContract_Life_RD.xml` | daftar kontrak |

⚠️ `[terverifikasi]` **Laporan ini tidak ada di Pega.** Korpus hanya **menjumlahkan dan menampilkan**
total per kontrak yang sedang dibuka; tidak ada satu pun rule yang mencari kontrak bercelah secara
menyeluruh. Ini **kemampuan baru**, konsekuensi langsung dari keputusan **tidak memblokir**
penyimpanan (Q4).

## ADR terkait

**ADR-0003** (persentase non-float — perbandingan terhadap 100 harus tepat, tanpa galat pembulatan),
**ADR-0001** (celah alokasi berdampak pada konteks hilir).

## Acceptance criteria

- [ ] Tersedia **laporan kontrak dengan total share ≠ 100%**. *(AC 20 spec)*
- [ ] Laporan menampilkan, per kontrak: tahun treaty, jenis reasuransi, **total share terhitung**,
      dan **selisihnya terhadap 100%**.
- [ ] Kontrak yang totalnya **kurang dari 100%** dan yang **lebih dari 100%** keduanya muncul —
      laporan tidak hanya mencari yang kurang.
- [ ] Kontrak **tanpa reinsurer sama sekali** (total `0`) muncul di laporan, tidak terlewat karena
      dianggap tidak punya data.
- [ ] Perbandingan terhadap 100 memakai **desimal presisi arbitrer**; kontrak yang totalnya tepat
      `100` **tidak** muncul karena galat pembulatan. *(**ADR-0003**)*
- [ ] Tiap baris laporan dapat **dibuka ke kontraknya**, sehingga celah dapat langsung diperbaiki.
- [ ] Laporan dapat disaring setidaknya per **tahun treaty**.
- [ ] Laporan **tidak mengubah apa pun** — murni baca.

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Mengapa laporan ini perlu ada.** `[keputusan work owner]` Total share **sengaja tidak
memblokir** penyimpanan, supaya kontrak dapat disusun bertahap. Konsekuensinya: celah alokasi
**pasti akan ada** pada suatu waktu. Tanpa laporan ini, celah itu hanya terlihat oleh orang yang
kebetulan membuka kontraknya — dan itu membuat keputusan "tidak memblokir" menjadi berbahaya.
Laporan inilah yang membuatnya aman.

⚠️ **Nama menyesatkan** `[terverifikasi]`: di Pega, total share ditaruh di field **`STDRATING`**
("standard rating"). Jangan bawa nama itu ke laporan — beri nama yang jujur.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```
