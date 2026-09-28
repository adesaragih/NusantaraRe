# 06: Efek keluar — kiriman Arasapas, log panggilan, dan alarm kegagalan simpan

**Status:** sebagian — log per panggilan (berhasil + isi respons), pekerja pengulang modul, dan panggilan nyata Arasapas/email (OQ-PL-11) belum ada

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**, 05a (efek keluar berjalan setelah penyimpanan satu-transaksi selesai dan keadaannya diketahui)

> ⚠️ Dikoreksi 2026-09-16: sebelumnya menunjuk 05b, yang kini **wontfix** (JSON dibuang). Simpan polis
> kini satu transaksi utuh di **05a**, jadi efek keluar bergantung pada 05a.

## Hasil & nilai pengguna

Sebagai **tim operasi**, saya ingin premium list yang sudah tersimpan diteruskan ke Arasapas dengan
alamat yang **selalu dibaca runtime**, setiap panggilan **tercatat** beserta responsnya, dan saya
diberi tahu ketika penyimpanan ternyata **gagal** — supaya kegagalan integrasi tidak pernah menjadi
temuan tutup buku. *(User story 37–38, 40 di spec)*

## Area codebase

`internal/services` (orkestrasi efek keluar; gerbang lingkungan), `internal/repository` (resolusi
alamat dari `M_LINK_SERVICE`; tulis log panggilan), `internal/clients` (klien Arasapas + pengirim
email di balik interface), `internal/handlers` (status efek keluar terbaca API).

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `serviceInsertArasapasLife_act` | `ASM-FW-GISFW-WORK-LIFE` / `SERVICEINSERTARASAPASLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/serviceInsertArasapasLife_act.xml` (104.302 byte, **9 langkah**) | efek bisnis |
| `GetLinkService` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/GetLinkService.xml` | resolusi alamat runtime |
| `InsertLogServiceProd` (Activity) | `ASM-FW-GISFW-WORK` / `INSERTLOGSERVICEPROD` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/InsertLogServiceProd.xml` | tulis log |
| `InsertLogServiceProd` (SQL) | `ASM-FW-GISFW-WORK` / `RNM!INSERTLOGSERVICEPROD` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/InsertLogServiceProd.xml` | `INSERT INTO pooldata.MONITORING_PROD_LOG(IDPEGA, NOPOLIS, PARAMETER, JN_SERVICE, STS_MESSAGE, RESPON_MESSAGE)` |
| `SendEmailNotification` | `@BASECLASS` / `SENDEMAILNOTIFICATION` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/SendEmailNotification.xml` | **jalur alarm** |
| `GetPolicyNoByCaseId` | `ASM-FW-GISFW-…` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GetPolicyNoByCaseId.xml` | nomor polis dari case id |
| `serviceInsertArasapasLife_act` (Endorsement) | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SERVICEINSERTARASAPASLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/serviceInsertArasapasLife_act.xml` | **rule berbeda** — class berbeda, isi berbeda. ⚠️ **konteks lain**, disebut hanya sebagai pembanding |

`[terverifikasi]` **Peta 9 langkah `serviceInsertArasapasLife_act`**:

| Step | Baris | Langkah |
| ---: | ---: | --- |
| 1 | 304 | ambil `CaseId` |
| 2 | 438 | `RDB-List` → `GetPolicyNoByCaseId` |
| 3 | 615 | set `PolicyNo` + sysdate |
| **4** | 799 | **`Call ASM-FW-GISFW-Int-M_LINK_SERVICE.GetLinkService`** — "GET LINK SERVICE" |
| **5** | 915 | **`Connect-REST`** — "Hit service Arasapas" |
| 6 | 1066 | set param insert log |
| **7** | 1305 | **`Call InsertLogServiceProd`** — "Insert Log Service" |
| 8 | 1470 | `Page-Remove` |
| ~~9~~ | 1585 | ~~`RDB-List` "Set err_note sts_konversi 9 (UPDATE JSON_POLIS)"~~ — **REMARK** (1597) |

`[terverifikasi]` **Ini bukti langsung ADR-0013**: alamat Arasapas **tidak** ada sebagai literal di
`Connect-REST`; ia di-resolve satu langkah sebelumnya lewat kelas `ASM-FW-GISFW-Int-M_LINK_SERVICE`.

`[terverifikasi]` **Gerbang lingkungan.** Di `InsertJsonPolisLife_Act` (PremiumList), tiga
precondition `IsPEGAPROD` (baris 5128, 5236, 5497) menjaga step **14** (email), step **15**
(Arasapas), dan step **17** (`Connect-REST`, sudah REMARK).

`[terverifikasi]` **Email adalah alarm, bukan notifikasi bisnis.** Step 14 dijaga **dua** syarat:
`OutDataLife.pxResults(1).PL_NUMBER==""` (5105) **dan** `IsPEGAPROD` — yakni ia menyala **hanya bila
pembacaan balik menunjukkan penyimpanan gagal**. `[keputusan work owner]` dikonfirmasi.

⚠️ `[terverifikasi]` **Catatan lintas konteks — jalur endorsement berbeda, dan itu BUKAN cakupan
tiket ini.** Di `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT`, 444.634 byte, **16 langkah**),
hanya dua langkah yang REMARK — dan keduanya justru **jalur deteksi dan alarm**: step **13** "Cek
sudah masuk atau blm datanya" (6399) dan step **15** `SendEmailNotification` (6752). Yang tersisa
aktif hanyalah step **16** `serviceInsertArasapasLife_act`, dijaga `IsPEGAPROD` (7230).
**Endorsement hari ini berjalan tanpa alarm.**

`[keputusan work owner]` **Endorsement Life adalah konteks terpisah.** Temuan ini dicatat di
`.scratch/endorsement-life/` sebagai bahan; keputusan apakah alarm diseragamkan **diambil di konteks
itu**, bukan di sini. Tiket ini hanya mengikat jalur **new business**.

`[terverifikasi]` `ConvertJsonNusareToProduction` (`Connect-REST`, step 17) **REMARK** — tidak
dimigrasikan. *(AC 30 spec)* ⛔ *Ralat 28-09-2026 (sensus remark):* step 17 `InsertJsonPolisLife_Act`
yang ter-remark memanggil `InsertLifePremiumDetail` (b5422). `convertJsonNusareToProduction` **hidup**
di `serviceInsertArasapasLife_act` langkah 5 (b963) — bagian efek Arasapas tiket ini, yang masih stub.
Premis AC 30 → **OQ-PL-14**.

## ADR terkait

**ADR-0013** (alamat di-resolve runtime dari `M_LINK_SERVICE`; **dilarang** sebagai literal,
konstanta, **maupun env var** — menggantikan **ADR-0004**), **ADR-0005** (`IsPEGAPROD` → flag
lingkungan), **ADR-0008** (efek keluar asinkron), **ADR-0007** (jejak audit).

## Acceptance criteria

- [x] Alamat Arasapas di-resolve **runtime** dari `M_LINK_SERVICE` pada setiap panggilan. Test yang
      memindai kode untuk URL sebagai literal, konstanta, **atau pembacaan env var** **gagal** bila
      menemukannya. *(AC 26 spec; **ADR-0013**)* — bukti: `services/polis_efekkeluar.go:EfekArasapasPolis.Jalankan` (lewat `AlamatLayanan`); uji `TestNolAlamatLayananDiKode`
- [x] Di lingkungan non-production, **kedua** efek keluar tidak berjalan — sementara penyimpanan
      premium list **tetap berjalan penuh**. *(AC 25 spec; **ADR-0005**)* — bukti: uji `TestBukanProduksiNolPanggilan`; `services/polis_penawaran.go:terapkan` (simpan tidak bergerbang)
- [x] Gerbang lingkungan adalah **satu** flag yang dibaca di satu tempat, bukan pemeriksaan tersebar. — bukti: `services/polis_efekkeluar.go:NewPenyalurPolis` (gerbang milik `Penyalur`); uji `TestGerbangLingkunganBenarBenarMenggerbangi`
- [ ] Setiap panggilan keluar menulis satu baris log berisi identitas kasus, nomor polis, parameter,
      jenis layanan, status, dan **isi respons** — baik saat berhasil maupun gagal. — belum: hanya kegagalan yang tercatat (outbox); log keberhasilan dan isi respons belum ada — kedua efek masih stub
- [ ] Email terkirim **hanya** bila pembacaan balik menunjukkan penyimpanan **gagal** (`PL_NUMBER`
      kosong). Penyimpanan yang berhasil **tidak** mengirim email. *(AC 27 spec)* — belum: digantikan AC relasional (pemicu kini kegagalan Arasapas); pengirim email masih stub (`ErrEmailBelumDisetujui`)
- [ ] Kegagalan efek keluar **tidak** membatalkan premium list yang sudah tersimpan; ia tercatat dan
      terlihat, dan dapat diulang. (**ADR-0008**) — belum: tidak membatalkan simpan dan tercatat di outbox (uji `TestEfekBerjalanSesudahCommitBukanDiDalamnya`), tetapi baris `PREMIUMLISTLIFE` belum punya pekerja pengulang
- [x] Arasapas dan pengirim email berada **di balik interface** dan difake di test; yang diperiksa — bukti: interface `EfekKeluar`/`Antrean`; uji `TestAlarmHanyaBilaArasapasGagal`, `TestAlarmYangGagalIkutTerantre`

### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §12

- [ ] ⚠️ Efek keluar membaca data polis **dari tabel relasional**, bukan dari payload JSON. Test yang
      menemukan perakitan CLOB JSON untuk dikirim keluar **gagal**. *(AC 32, 34 spec; penyimpangan
      sadar 1)* — belum: kiriman masih stub dengan muatan pengenal saja (OQ-PL-11), dan uji penolak perakitan CLOB JSON belum ada
- [x] ⚠️ Pemicu alarm **tidak lagi** berupa "deteksi keadaan separuh" antar-procedure — keadaan itu
      mustahil karena polis ditulis **atomik**. Alarm kini menyala pada **kegagalan efek keluar**
      saja. *(AC 35 spec; spec §6)*
      test adalah **efeknya** (panggilan terjadi/tidak, log tertulis, email terpicu/tidak). — bukti: `services/polis_efekkeluar.go:PenyalurPolis.Salurkan`; uji `TestAlarmHanyaBilaArasapasGagal`

## Blocker

**Tidak ada.**

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

## Implementasi — 28-09-2026 (giliran 10): outbox + stub, alarm hanya pada kegagalan

### XML yang dibaca sebelum kode

| Rule | Yang diambil |
| --- | --- |
| `Activity/InsertJsonPolisLife_Act.xml` | langkah 14 `SendEmailNotification` (b4723, *"Kalau blm, email errornya"*), 15 `serviceInsertArasapasLife_act` (b5168) |
| `Activity/serviceInsertArasapasLife_act.xml` | langkah 4 `GetLinkService` — `Kategori_1 = "Production"`, `Kategori_2 = "convertJsonNusareToProduction"` (pecahan b844–845); langkah 5 `Connect-REST` `ServiceName convertJsonNusareToProduction`, `POST` |
| `ConnectREST/ConvertJsonNusareToProduction.xml` | pemetaan permintaan: `.pzInsKey`, `.OfferFacIn.PolicyData.PolicyNo`, `.pxCreateDateTime` — **tiga pengenal**, respons ke `.StatusService` |

### Yang dibangun — memakai ulang mesin Claim Life, tidak menyalinnya

- `services/polis_efekkeluar.go`: `KunciArasapasPremiumList` VERBATIM (≠ kunci Claim Life, dikunci
  `TestKunciArasapasPolisVERBATIM` yang membaca activity-nya); `EfekArasapasPolis` me-resolve alamat
  **sungguhan** lewat `M_LINK_SERVICE` lalu berhenti terang (`ErrArasapasBelumDisetujui`);
  `EfekAlarmEmailPolis` stub (`ErrEmailBelumDisetujui`).
- `PenyalurPolis`: dua `Penyalur` bersarang — kiriman, lalu alarm **hanya bila** kiriman gagal
  (`TestAlarmHanyaBilaArasapasGagal`). Satu `Penyalur` berisi dua efek akan mengirim alarm pada setiap
  simpan yang berhasil. Gerbang lingkungan tetap **satu tempat** — milik `Penyalur`
  (`TestBukanProduksiNolPanggilan`).
- Outbox `T_LOG_SERVICE_RNM` dipakai bersama; `antreanOracle` kini membawa `modul` —
  `AntreanEfekOracleModul(svc, "PREMIUMLISTLIFE")`. Claim Life tetap `CLAIMLIFE` (perubahan aditif pada
  berkas milik Claim Life, dilaporkan).
- `Penawaran.terapkan`: efek keluar dipanggil **sesudah** transaksi commit, **hanya** bila
  `SimpanPolis` (`TestEfekBerjalanSesudahCommitBukanDiDalamnya`). Hasilnya ringkasan
  (`HasilSubmitSummary.EfekKeluar`), bukan galat. Handler menyuntikkan
  `PenyalurPremiumListOracle` di ketiga rute (keputusan, penggolong, summary).
- Layar summary mengatakan hasilnya: dilewati (bukan produksi) / gagal + terantre / gagal tak terantre.

### ⛔ OQ-PL-11 `[terbuka — work owner / pemilik Arasapas]` — layanan hilir membaca JSON yang tidak lagi ada

Nama layanannya `convertJsonNusareToProduction`, dan permintaannya hanya **tiga pengenal** — bukan isi
polis. Layanan itu tampaknya membaca `JSON_POLIS` di sisinya sendiri; tabel itu **tidak lagi ditulis**
(pl1). Mengaktifkan kiriman ini tanpa jawaban berarti memanggil layanan yang akan mencari JSON yang
tidak ada. Stub tetap gagal terang sampai pemilik layanan menyatakan dari mana ia membaca polis.

### AC — keadaan

| AC | Keadaan |
| --- | --- |
| alamat di-resolve runtime; nol URL literal/konstanta/env | ✅ `KunciLayanan` + `AlamatLayanan`; penjaga `TestNolAlamatLayananDiKode` (milik Claim Life) ikut memindai berkas baru |
| non-produksi: kedua efek tidak berjalan, simpan tetap penuh | ✅ `Penyalur` bergerbang; simpan tidak bergerbang |
| gerbang satu flag satu tempat | ✅ `Service.lingkungan` → `Penyalur` |
| satu baris log per panggilan, berhasil maupun gagal | ⚠️ **sebagian** — yang tertulis kegagalan (outbox). Panggilan berhasil belum mungkin: kedua efek stub. Padanan `MONITORING_PROD_LOG` untuk keberhasilan dibangun bersama panggilan nyata |
| email hanya bila gagal | ✅ pemicunya bergeser ke kegagalan efek keluar (AC relasional) |
| kegagalan efek tidak membatalkan simpan; tercatat; dapat diulang | ✅ sesudah commit; outbox + `Backoff` Claim Life |
| di balik interface, difake di uji | ✅ `EfekKeluar`, `Antrean` |
| nol perakitan CLOB JSON | ✅ muatan outbox = pengenal + waktu |

### Angka

Go **540 PASS · 0 FAIL** tingkat atas; vet (+`-tags db`), gofmt bersih · vitest **348** · tsc bersih.
