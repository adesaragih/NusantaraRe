# 11: Kurs USD → IDR

**Status:** selesai (29-09-2026); **diralat lanjutan 6** (29-09-2026) — tanggal master kurs diurai Oracle seperti `GetMasterKursList`

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
      setiap kueri. *(AC 53 spec; penyimpangan sadar 6)* ⚠️ **Diralat lanjutan 6 (29-09-2026):** tanggal diurai
      **Oracle** di setiap kueri, seperti Pega; Go menerima `DATE` — bab bertanggal di akhir tiket.
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

## Keputusan work owner 29-09-2026

- **OQ-TCO-18 — ditutup.** Jawaban: *"setuju"*. `KURS` diisi kurs yang dipakai menghitung `Usd` tujuh induk berkurs,
  skala 8 untuk ketujuh pembagian (exclusion IDR → USD tetap 4 VERBATIM), dan dua baris kurs berlaku = master rusak (503)
  `[keputusan work owner 29-09-2026]`.

## ⛔ Ralat bertanggal — 29-09-2026 (lanjutan 6: kurs dibaca seperti `GetMasterKursList`)

Laporan work owner: `List Description` → `Show TreatyDesc` pada tahun `1000682` menampilkan *"Backend tidak terhubung"*
padahal backend hidup. Diagnosis asisten (brief lanjutan 6 §0): `GET …/tahun/1000682/kurs` dan `…/kurs/konversi` menjawab
**503** `models: nilai master kurs tidak dapat diurai: STARTDATE "20190801T00000.000 GMT" bukan bentuk …`.
`[data DEV — brief lanjutan 6 §0, agregat bentuk]` `STARTDATE` 129 baris `99999999T999999.999 GMT` dan 11 baris
`99999999T99999.999 GMT` (jam lima angka); `ENDDATE` 140/140 normal.

1. **Pengurai teks Go lebih ketat dari Oracle — dibuang.** Pega tidak mengurai tanggal ini di Java: `GetMasterKursList`
   b85–b86 menyerahkannya ke Oracle, dan Oracle menerima jam lima angka. `models.UraiTanggalKursTCO` menolaknya, sehingga
   satu baris mematikan kurs seluruh tahun. Kini `repository.sqlBerlakuKursTCO` membaca
   `trunc(TO_TIMESTAMP_TZ(STARTDATE|ENDDATE, 'YYYYMMDD"T"HH24MISS.FF3 TZR'))` dengan topeng VERBATIM dan membandingkan
   `to_date(:1,'YYYYMMDD') BETWEEN …` di SQL; Go menerima `DATE` dan mengambil harinya apa adanya. AC 53 diralat:
   penyimpangan sadar 6 dicabut untuk kolom ini.
2. **Satu baris cacat tidak mematikan yang lain.** Dua beda yang disengaja dari RDB: `DEFAULT NULL ON CONVERSION ERROR`
   (Oracle 12.2+; DEV 12.2.0.1 menurut `claim-life/SUMBER-PENOMORAN-DBA.md`) dan `BETWEEN` sebagai kolom `BERLAKU`, bukan
   saringan. Baris yang tanggalnya Oracle tolak (atau kosong) disaring dan **dicacah** — `barisMasterDitolak` di jawaban
   `GET …/kurs`. Tanpa baris berlaku tetapi ada yang ditolak → **503** berkata-kata (`… N baris master tanggalnya ditolak
   Oracle (bentuk …), mis. STARTDATE "…"`), bukan "Tidak ada Nilai Kurs": salah satunya mungkin baris yang dicari. Di Pega
   satu baris cacat menggagalkan seluruh kueri.
3. **Tabel.** RDB menyebut `treatyexchange`; yang dipakai tetap `TREATYEXCHANGEYEARLY`: `[data DBA]` `TREATYEXCHANGE` tidak
   ada (`dba-procedures.md` b83; bab "Rule Pega sumber" di atas); 12 RDB korpus di enam modul lain (Endorsment Fac In,
   NB FacIn, NB Treaty In, RNW Fac In, Treaty In, Treaty In Adjustment) membaca `treatyexchangeyearly`, sedangkan
   `treatyexchange` hanya disebut `GetMasterKursList`; agregat §0 dan galat 503 DEV pun berasal dari
   `TREATYEXCHANGEYEARLY`. `[belum diverifikasi executor di katalog DEV]` — perintah pemeriksaan baca-saja:
   `SELECT OWNER, OBJECT_NAME, OBJECT_TYPE FROM ALL_OBJECTS WHERE OBJECT_NAME IN ('TREATYEXCHANGE','TREATYEXCHANGEYEARLY')`.
4. **Uji.** Tanpa Oracle: `TestSQLKursTCO` (ekspresi VERBATIM + klausa `DEFAULT`, bind urut kemunculan, satu-satunya
   literal = topeng), `TestTambahBarisKursTCO` (hari dari `DATE` apa adanya — `.UTC()` menggeser 1 Januari WIB; NULL =
   ditolak, dicacah), `TestPilihKursBerlakuTCO`, `TestKursBarisTanggalDitolakOracle`; tujuh mutasi merah. Uji `db`
   `TestKursTanggalDiuraiOracle` (jam lima angka diterima; huruf, bulan 13, panjang salah, dan kosong ditolak lalu
   dicacah; 503 tanpa baris berlaku) — **SKIP di mesin executor** (tanpa skema uji), belum pernah dijalankan.
5. Klien yang menyebut 503 ini "Backend tidak terhubung" diperbaiki terpisah (perbaikan 2 lanjutan 6,
   `lib/keadaanGalat.ts`).

## ⛔ Keputusan work owner bertanggal — 29-09-2026 (penyisiran layar: baris kurs KEMBAR)

*Temuan.* Sesudah lanjutan 6, penyisiran `GET` seluruh rute Treaty Contract Out (182 tahun treaty) mendapati `…/kurs`
menjawab 503 *"lebih dari satu kurs berlaku pada tanggal yang sama: 2 baris"* pada **23 tahun** (tanggal mulai 2019-08-01,
2025-07-01, 2026-06-01). `[data DEV 29-09-2026 — dibaca executor, SELECT baca-saja atas izin work owner]`: master USD `QUARTER='0'` memuat 12 baris, di antaranya **dua pasang baris
kembar persis** — `14500.00` `20190801T00000.000 GMT`–`20200630T000000.000 GMT` (dua kali) dan `16500.00`
`20250701T140000.000 GMT`–`20260630T140000.000 GMT` (dua kali); nol tanggal ditolak Oracle.

*Keputusan (jawaban: "Kembar identik = satu kurs").* **OQ-TCO-18 dipersempit**: dua baris berlaku atau lebih yang
TOIDR-nya **sama menurut angka** adalah satu kurs — Pega "terakhir menang" memberi nilai yang sama, jadi tidak ada yang
ditebak; yang dipakai baris yang mulainya paling akhir (periode `Mulai`/`Akhir` di layar tidak bergantung urutan baca).
TOIDR **berbeda** tetap master rusak (503), kini dengan kedua nilainya disebut. Cacah baris kembar dilaporkan
(`barisMasterKembar` di `GET …/kurs`) dan **disebut di layar** bersama `barisMasterDitolak`
(`PanelJenisKlausul.catatanMasterKurs`). Uji: `TestPilihKursBerlakuTCOBarisKembar` (kembar persis, kembar menurut angka
dengan periode berbeda, TOIDR berbeda, TOIDR rusak); uji lama yang menuntut galat atas `{berlaku, berlaku}` diganti.
Data master TIDAK disentuh; menghapus baris kembar tetap urusan DBA.
