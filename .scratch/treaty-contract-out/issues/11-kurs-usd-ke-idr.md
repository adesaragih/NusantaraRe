# 11: Kurs USD → IDR

**Status:** ready-for-agent

**Blocked by:** 08 (kurs melekat pada baris klausul)

## Hasil & nilai pengguna

Sebagai **underwriter**, saya ingin nilai dalam **USD** dapat dilihat padanannya dalam **IDR**
menurut kurs yang berlaku **pada periode kontrak**; dan sebagai **Finance**, saya ingin nilai
**Rp** dan **USD** tetap tercatat sebagai **dua nilai terpisah**, karena keduanya memang dua angka
yang berbeda. *(User story 28–30 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Baca kurs per periode dari master (read-only) |
| `internal/services` | Konversi USD → IDR; rujukan mata uang |
| `internal/handlers` | Endpoint kurs berlaku |
| `frontend/` | Tampilan padanan IDR pada baris klausul |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `testingKurs` | `@BASECLASS` / `TESTINGKURS` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/testingKurs.xml` | ⚠️ **jalur yang benar-benar dipakai** — nol langkah di-remark |
| `RefreshKurs` | `@BASECLASS` / `REFRESHKURS` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/RefreshKurs.xml` | muat ulang kurs |
| `GetMasterKursList` | `ASM-FW-GISFW-INT` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetMasterKursList.xml` | baca kurs per periode |
| `NitipKurs` | `RULE-HTML-SECTION` | `Treaty Contract Out/Section/NitipKurs.xml` | tampilan kurs |
| `SetTreatyArrangementDesc_Act` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `SETTREATYARRANGEMENTDESC_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetTreatyArrangementDesc_Act.xml` | ⚠️ **MATI** — 5 dari 6 langkah di-remark |

`[terverifikasi]` Kueri kurs existing:

```
select TOIDR from treatyexchange
 where Quarter = '0'
   and to_date({tanggal},'YYYYMMDD') BETWEEN trunc(TO_TIMESTAMP_TZ(STARTDATE, 'YYYYMMDD"T"HH24MISS.FF3 TZR'))
                                         AND trunc(TO_TIMESTAMP_TZ(ENDDATE,   'YYYYMMDD"T"HH24MISS.FF3 TZR'))
   and IDCURRENCY = '10001'
```

`[data DBA]` **Nama tabel sebenarnya `TREATYEXCHANGEYEARLY`** (bukan `TREATYEXCHANGE`); kolom
`TOIDR`, `TOUSD`, `IDCURRENCY`, `CURRENCY`, `QUARTER`, `STARTDATE`, `ENDDATE` — **seluruhnya
`VARCHAR2`**, termasuk kedua tanggal.

⚠️ **Penyimpangan sadar 7 — `IDCURRENCY='10001'` jadi rujukan master, bukan literal.**
`[keputusan work owner]` `IDCURRENCY = '10001'` **berarti USD**, dan konversi **USD → IDR** memang
aturan bisnis. Yang tidak boleh adalah identitas mata uangnya ditanam sebagai konstanta program.

⚠️ **Penyimpangan sadar 8 — nama jujur.** `[keputusan work owner]` Komponen di sistem baru
**tidak** mengandung kata `"testing"`.

## ADR terkait

**ADR-0003** (uang non-float), **ADR-0015** (kegagalan ditangani eksplisit).

## Acceptance criteria

- [ ] Padanan **IDR** untuk nilai **USD** dihitung dari kurs yang **berlaku pada periode** kontrak —
      tanggal berada di antara tanggal mulai dan tanggal akhir baris kurs. *(AC 46 spec;
      User story 28)*
- [ ] ⚠️ Identitas mata uang USD (`IDCURRENCY = '10001'`) adalah **rujukan ke master mata uang**,
      **bukan literal di kode**. Test yang menemukan `'10001'` sebagai konstanta program **gagal**.
      *(AC 47 spec; penyimpangan sadar 7)*
- [ ] `QUARTER = '0'` **diikuti apa adanya**; artinya dicatat sebagai `[terbuka]` dan **tidak
      ditebak**. *(AC 48 spec)*
- [ ] ⚠️ Nilai **Rp** dan **Usd** tetap **dua nilai terpisah** pada baris klausul — bukan satu nilai
      dengan kode mata uang. *(AC 49 spec; User story 30)*
- [ ] ⚠️ Jalur kurs yang dimigrasikan adalah yang **benar-benar dipakai** (`testingKurs`), bukan
      yang namanya lebih wajar tetapi **di-remark**. Nama komponen barunya **tidak mengandung kata
      "testing"**. *(AC 50 spec; penyimpangan sadar 8; **OQ-066**)*
- [ ] Tanggal mulai/akhir baris kurs diperlakukan sebagai **tanggal**, bukan teks yang di-parse
      setiap kueri. *(AC 53 spec; penyimpangan sadar 6)*
- [ ] Nilai kurs diperlakukan sebagai **desimal presisi arbitrer**; **tidak** melewati `float`.
      *(AC 51 spec; **ADR-0003**)*
- [ ] Master kurs **tidak ditulis** oleh konteks ini. Test yang menemukan tulisan ke
      `TREATYEXCHANGEYEARLY` **gagal**.
- [ ] Periode yang **tidak punya baris kurs** menghasilkan kegagalan yang **terlihat**, bukan nilai
      nol atau kosong yang diam. *(**ADR-0015**)*

## Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka]` **Arti `QUARTER = '0'`** — OQ kecil, **tidak memblokir**: AC di atas mewajibkan
nilainya diikuti apa adanya, bukan dipahami.

## Catatan

⚠️ **OQ-066 terbalik di sini.** `[terverifikasi]` Activity bernama **`testingKurs`**
(`@BASECLASS!TESTINGKURS`) **nol langkah di-remark** — ia hidup dan dipakai. Sebaliknya
`SetTreatyArrangementDesc_Act` (`ASM-FW-GISFW-INT-PROPORTIONALARRG`), yang namanya paling wajar,
**lima dari enam langkahnya di-remark** — ia mati. Biasanya nama "testing" menandai yang mati; di
sini justru sebaliknya. **Baca kodenya, jangan namanya.**

⚠️ **Nama tabel di korpus dan di basis data berbeda.** Kueri Pega menyebut `treatyexchange`;
`[data DBA]` menyebut tabel sebenarnya **`TREATYEXCHANGEYEARLY`**. Migrasi (tiket 01) memakai nama
DBA.

⚠️ **Dua kolom uang dipertahankan dengan sengaja.** `[keputusan work owner]` `Rp` dan `Usd` bukan
satu nilai yang ditampilkan dua cara — keduanya **dua angka berbeda** yang disimpan berdampingan
(dan untuk baris anak, keduanya **turunan** dari `Pct` × nilai induk — lihat tiket 08). Jangan
"merapikannya" jadi satu kolom + kode mata uang.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — pemilihan baris kurs menurut periode
adalah pertanyaan rentang tanggal di basis data.

```
go test ./internal/...
cd frontend && npm test
make check
```
