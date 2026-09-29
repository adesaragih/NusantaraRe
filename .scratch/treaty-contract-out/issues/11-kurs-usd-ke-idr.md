# 11: Kurs USD → IDR

**Status:** selesai (29-09-2026)

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

## Pembacaan ulang XML — 29-09-2026 (sesi modul, lanjutan 1)

Nomor baris = baris mentah berkas korpus `Treaty Contract Out/` kecuali disebut lain; langkah aktivitas dibaca lengkap.

| Unsur | Bukti | Dibawa sebagai |
| --- | --- | --- |
| jalur hidup | `Activity/testingKurs.xml`: b273 `CARI1 = Param.StartDate`; b438 `RDB-List GetMasterKursList`; b579–b580 `InputTreatyArrangement.Kurs = .HASIL1` (untuk SETIAP hasil — yang terakhir menang) | `services.KursTCO.Berlaku` (nama tanpa "testing") |
| kueri | `RDBList/GetMasterKursList.xml` b85–b86: `Quarter='0'`, `to_date(CARI1,'YYYYMMDD') BETWEEN trunc(TO_TIMESTAMP_TZ(STARTDATE,…)) AND trunc(TO_TIMESTAMP_TZ(ENDDATE,…))`, `IDCURRENCY` literal | `TREATYEXCHANGEYEARLY` `[data DBA]` dibaca saja; `QUARTER`/`IDCURRENCY` di-bind; tanggal diurai sekali di Go |
| pemanggil | `Harness/InboxTreatyContractDescription.xml` `Show` Non XOL b6138 (`StartDate = InputTreatyArrangementDesc.StartDate` b6153); `Show` XOL b9256 (`TreatyYear = …TreatyYear` b9271) | `GET /tahun/{id}/kurs` memakai `StartDate` tahun untuk KEDUA grid |
| master mata uang | `Claim Life/RDBList/GetCurrencyID.xml` b85 `SELECT ID … FROM POOLDATA.CURRENCY WHERE CURRENCY = {kode}`; `NB FacIn/RDBList/GetCurrencyIDByName.xml` (`ID, OLDID, CURRENCY, CURRENCYSYMBOL`) | pengenal USD = `MataUang.Pengenal("USD")` (kode bersama, tidak diubah) |
| konversi induk | `Activity/HitungRpUsd_depan.xml` b383 `Usd = @divide(Rp, Kurs, 8)` (TreatyLimit), b405–b510 `Usd = Rp / Kurs` (PLA, CashLossLimit, FacIn, ExGratia, EPI, ClaimCoorp); dipanggil onchange Rp (mis. `GridTreatyArrangementEpi.xml` b3630); medan `Usd` hanya dibaca di ketujuh form (EPI b3775, PLA b3821, TreatyLimit b4062, FacIn b3653, CashLossLimit b3777, ClaimCoorp b3747, ExGratia b1911) | tujuh induk: `Usd` TURUNAN `Rp ÷ Kurs` skala 8 di server |
| konversi exclusion | `Activity/CalculateTSIExcludeTreaty.xml` b248 IDR → `Usd = @divide(Rp, Kurs, 4)`; b394 USD → `Rp = Usd * Kurs`; dipanggil `GridTreatyArrangementExclutionTreatyOccupation.xml` b3003/b3123/b3290/b3412 | pratinjau dua arah lewat `GET /tahun/{id}/kurs/konversi`; kedua nilai tetap masukan |
| gerbang | 14 aktivitas `NewTreatyArr*` (7 induk + 7 anak): `DATASHOW = @if(Kurs="","","1")` lalu `ERRMSG2 = "Tidak ada Nilai Kurs di Tahun : " + TreatyYear` (mis. `NewTreatyArrEpi.xml` b707, b870, b948) | jenis berkurs ditolak 422 dengan pesan VERBATIM; `Add` nonaktif di layar |
| tempat kurs | `Section/NitipKurs.xml` b512 (`InputTreatyArrangement.Kurs`, harness b3882); `SaveMasterProportionalArrg.xml` b98 menulis `InputTreatyArrTreatyLimit.Kurs` | `KURS` klausul = kurs yang dipakai (tujuh induk berkurs) |

### Ralat bertanggal 29-09-2026

1. **Grid XOL tidak pernah menemukan kurs di Pega**: `Show` XOL (b9256) mengirim parameter `TreatyYear` (b9271), padahal
   `testingKurs` membaca `Param.StartDate` (b273) → `to_date('')` → nol baris. Sistem baru memakai `StartDate` tahun
   untuk kedua grid.
2. **Kolom `KURS` warisan praktis selalu kosong**: prosedur menulis `InputTreatyArrTreatyLimit.Kurs` (b98), sedangkan
   satu-satunya penulis kurs mengisi `InputTreatyArrangement.Kurs` (`testingKurs` b579, `RefreshKurs` b244). Sistem baru
   menyimpan kurs yang DIPAKAI menghitung `Usd` pada tujuh induk berkurs `[keputusan kami]` (**OQ-TCO-18**).
3. **Tiket 08 diralat**: `Usd` tujuh induk (TreatyLimit, PLA, CashLossLimit, FacIn, ExGratia, EPI, ClaimCoorp) bukan
   masukan — ia hanya dibaca di form dan diturunkan `Rp ÷ Kurs`. Klien yang mengirimnya ditolak 422.
4. **Skala pembagian**: hanya TreatyLimit yang tersurat 8 desimal (`@divide(…, 8)`); enam form lain memakai `/` tanpa
   skala. Sistem baru memakai 8 (skala kolom `NUMBER(38,8)`) untuk ketujuhnya, setengah-ke-atas; exclusion IDR → USD
   tetap 4 VERBATIM (**OQ-TCO-18**).
5. **Dua baris kurs berlaku pada tanggal yang sama**: Pega menyimpan yang terakhir (urutan tak tentu). Sistem baru
   menolaknya sebagai master rusak (503) `[keputusan kami]`, bukan menebak (**OQ-TCO-18**).
6. **`QUARTER = '0'`** diikuti apa adanya (`QuarterKursTahunanTCO`); artinya tetap `[terbuka]`.
7. **Tanggal mulai tahun kosong** diperlakukan sebagai "tidak ada kurs" (pesan yang sama), bukan tanggal nol.

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_kurs.go` (+uji) | urai tanggal/nilai master sekali, `PilihKursBerlakuTCO`, `UsdDariRpTCO` / `RpDariUsdTCO` (apd), `GalatKursTidakAda` VERBATIM; aturan klausul `Berkurs`/`Konversi` |
| repository | `tco_kurs.go` (+uji) | `MasterKursTCO` baca-saja, saringan di-bind; tiruan `TREATYEXCHANGEYEARLY` + `CURRENCY` |
| services | `tco_kurs.go` (+uji) | `KursTCO` (berlaku, konversi); `KlausulTCO.DenganKurs` — Usd induk turunan, `KURS` tersimpan, gerbang anak |
| handlers | `tco_kurs.go` (+uji, +uji `db`) | 2 rute GET; uji `db` klausul diperbarui |
| frontend | `PanelJenisKlausul.tsx` (+uji), `KURS_TCO`, `api.ts` (+2) | kurs berlaku / pesan server, `Add` nonaktif tanpa kurs, pratinjau konversi dari server |

**Status:** selesai 29-09-2026 — commit `treaty-contract-out: tiket 11 — kurs USD ke IDR`.
