# Laporan implementasi EDM — putaran 11-10-2026 (irisan I0–I2: modul terdaftar, layar pertama)

> Atas keputusan work owner 11-10-2026 (`KEPUTUSAN-WORK-OWNER-EDM.md` W-EDM-01…04). Backend dikerjakan sesi
> asisten `nusantarare-0b` atas kontrak sesi utama; frontend, `MODUL.md`, tiket, dan laporan oleh sesi utama.
> **Belum di-commit.** Korpus `D:\migrasi\RNM\Endorsment Fac In\` (READ-ONLY). Label mengikuti `CLAUDE.md` §4.
> Nol nomor polis, nol nama orang di berkas ini dan di kode modul.

## 1. Yang lahir

| Irisan | Isi | Bukti |
| --- | --- | --- |
| **I0** pendaftaran | `backend/modul.go` (tanpa HakLihat, Bangun tanpa Oracle), `inti/backend/daftar/modul_endorsmentfacin_gen.go` (bangkitan), migrasi `966_menu_endorsmentfacin.sql` (+ `_down`), `frontend/{menu.ts,rute.tsx}`, `MODUL.md` Status `dimigrasi`, prefix `/api/endorsmentfacin` | `TestDaftarBangkitanSamaDenganFolder` ok; `daftar.modulAktif.test.ts` ok |
| **I1** form | Portal `SFAPortalEndorsement` (judul `Addendum`, `Create Addendum`; **tanpa grid** — `pyGridRDName` kosong, sumber tidak ada di korpus) dan form *Create Addendum* = `Section/WorkPrimaryDetails.xml` + `EdmType1/2/3` + `crmNewHarnessButtons` (urutan medan, label, dan syarat tampil apa adanya) | `frontend/pages/{PortalEndorsement,FormAddendum}.tsx`; label diuji ke korpus |
| **I2** periksa | `POST /api/endorsmentfacin/periksa` (NOL tulis): 6 gerbang `SetErrorBatalEndorsement_Act`, jenis bisnis langkah 7–18, `CheckEDMPolisDate`; migrasi `261_m_edm_pengecualian_polis.sql` (tabel konfigurasi W-EDM-03) | `backend/{models/periksa.go, services/{periksa,layanan}.go, repository/{gudang,arasapas}.go, handlers/rute.go, tiruan/}`; `docs/STRUKTUR-TABEL-ENDORSMENTFACIN.md` |

Tombol **OK** sengaja mati: `SetValueToEDMWork` (buat kasus `EDM-`) adalah irisan I3 (W-EDM-02).

## 2. Keputusan agent — menunggu konfirmasi

A24–A33 dari sesi backend (diperiksa sesi utama atas berkasnya); A34–A36 dari sesi utama. A11–A23 tetap
ditunda (W-EDM-04).

| # | Keputusan | Bukti |
| --- | --- | --- |
| A24 | PK tabel 261 = (`NOPOLIS`, `JENIS`) + `CHECK JENIS IN ('TANGGAL','BATAL')` — satu polis dapat butuh dua jenis; CHECK menangkap salah ketik DBA | `migrations/261_…sql`; literal pengganti: `CheckEDMPolisDate` langkah 2 = 1, `SetValueToEDMWork` langkah 5 = 2 (dua cara: urai XML / cacah baris `pyStepsPreCondParamsWhen`, sepakat) |
| A25 | ✅ **DIPUTUSKAN W-EDM-06 — jam DATE dibaca WIB** (konvensi NB FacIn yang menulis kolom itu: `modul/nbfacin/backend/services/binding.go:54` `tanggalFacInOffer` = `t.In(WIB)`). Baris tulisan NB berjam 12:00 WIB → hasil tidak berubah; baris 00:00 WIB = 17:00 GMT hari sebelumnya → batas mundur sehari (K-046, `validasitanggal.go` tidak diubah). `tanggal.mulai/akhir` = batas yang DIPAKAI pembanding (kalender GMT). Uji: `periksa_test.go` (12:00 WIB, `TestPeriksaTanggalTengahMalamWIB`); mutan "kembali ke UTC" tertangkap | `services/periksa.go` `uraiDateOracle` |

> ⛔ **RALAT A25 (11-10-2026).** Versi pertama membaca jam DATE sebagai **GMT** dengan tafsir "JVM Pega ber-GMT",
> disimpulkan dari bentuk rule (`CheckEDMPolisDate` langkah 4 mengganti `Thh`→`T05` tanpa menggeser tanggal). Cara
> itu keliru: ia menebak zona dari pembaca, padahal zona ditentukan oleh **penulis** kolom, dan penulisnya (NB FacIn)
> terbukti menulis WIB. Untuk baris 12:00 WIB kedua tafsir memberi tanggal yang sama, sehingga kekeliruan itu tidak
> tertangkap uji yang hanya memakai jam tengah hari.
| A26 | `GetStartDate` tanpa baris / tanggal NULL → `tanggal.hasil = "tolak"`, pesan medan EndorsementDate "Policy period not found in facinproduction" (**teks bukan korpus**) — menutup OK tanpa 422 supaya gerbang tetap tampil | `services/periksa.go` |
| A27 | Notasi titik `b.DATA_JSON.QuotationData.EdmType` ditulis `JSON_VALUE(… NULL ON ERROR)`: constraint VALID_JSON `JSON_POLIS` berstatus DISABLE | `DDL\JSON_POLIS.txt`; `repository/gudang.go` |
| A28 | `ARASAPAS.DETAIL_INVOICE` hanya ditanya bila EdmType 1/2 — Pega selalu menjalankannya tetapi hanya membaca hasilnya untuk 1/2 (keluaran identik) | `SetErrorBatalEndorsement_Act` 5.15–5.18 |
| A29 | Hanya 7 predikat lini langkah 10–16 yang dinilai (bukan `PredikatDari`, yang juga menilai IsTravel/IsEDM dan butuh properti kasus) | langkah 10–16 |
| A30 | ✅ **DITUTUP W-EDM-05** ("life tidak termasuk") — penolakan di bawah ini perilaku yang dimaksud. Semula `[pertanyaan terbuka, berdampak ke layar]`: `IsLife` EDM membaca `pyWorkPage.Quotation.BusinessOldId`, dan **nol rule EDM** mengisinya di kasus portal → polis berpola RNML: jenis bisnis "Life" tetapi pemeriksaan bisnis **tolak** ("Invalid business, please contact IT!") | `When/IsLife.xml`; langkah 8 |
| A31 | `tanggal.catatan` (opsional) ditambahkan, mis. tabel 261 belum dibuat → daftar kosong + catatan; frontend menampilkannya | `models/periksa.go`, `api.ts` |
| A32 | `policyNo` dipakai apa adanya (teks bebas); hanya kosong-sesudah-trim yang 422. `retro` tidak divalidasi — "1.0" → `ErrTafsirKodeBerbeda` → 500 (sikap `samaKode`) | `handlers/rute.go` |
| A33 | `edmType` "" diterima (dinilai bukan-batal); "3" (`EdmUsang3`, K-029) → 422 | `services/layanan.go` |
| A34 | `[dugaan]` Teks pilihan **EDM Type** (`.Quotation.EdmTypeNew`) = Penambahan / Pengurangan / Perubahan / RI Slip. Prompt list properti itu **tidak ada** di korpus; nilai 1/2/3 terbukti dari visibilitas `EdmType1/2/3` (daftar `TypePenambahan/Pengurangan/Perubahan`), nilai 4 dari bypass "edm ri slip" (`CARI41=="4"`) | `frontend/labels.ts` `JENIS_EDM_BARU` |
| A35 | Detail EDM Type Penambahan memuat `12 Adjustment PPN/PPH` **sekali**. `[terverifikasi]` Korpus menambahkannya dua kali: langkah Pengurangan menulis ke `TypePenambahan` (K-046) — dropdown Pega menampilkan entri kembar; di sini kembaran dibuang (tampilan saja, nilai sama) | `Activity/SetEdmType.xml`; uji `endorsmentfacin.test.tsx` mencacah kembaran = 2 sebagai uji instrumen |
| A36 | ⚠️ Sebagian dijawab: Endorsement Status = W-EDM-13, Internal Retro = W-EDM-14; Endorsement Source tetap teks bebas. Semula: **Endorsement Source** sementara teks bebas: sumber pilihannya `D_EnumerationList` (`Type=sumberedm`, `Language=id-ID`, `Status=1`) — isi tabel enumerasi Pega tidak ada di korpus; autocomplete Pega pun mengizinkan isian bebas. **Endorsement Status** belum dirender (prompt list tidak ada di korpus). **Endorsement Internal Retro** menampilkan kode `1` / `0` (teks tidak ada di korpus) | `WorkPrimaryDetails` L2011-2290, L8836-8995, L9896-10213 |

**Tambahan W-EDM-07** — nomor A40/A41 (A37 dan A39 sudah dipakai modul ini untuk keputusan register NB yang dikutip kode: A37 = HALF_UP `@Math.divide`, A39 = rekonsiliasi EDM dengan rumus NB, ditahan). ⛔ Ralat: draf pertama bagian ini menomori klep A37/A38 tanpa memeriksa nomor terpakai - bentrok dengan A37 (temuan sesi `nusantarare-0b`, `grep -rnw A37` = 14 rujukan di kode dan dokumen):

| # | Keputusan | Bukti |
| --- | --- | --- |
| A40 | `[dugaan]` Saringan **F1 `.IDPEGA = Param.Idpega` TIDAK diport**: satu-satunya pemanggil (`SetErrorBatalEndorsement_Act` 5.3/5.12/5.16/5.20) tidak pernah mengisi `Param.Idpega`; filter tanpa penanda "null bila kosong" (grep `IfEmpty`/`NullIf` = 0). Bacaan "saringan berparameter kosong diabaikan" dan "produksi tanpa F1" (F1 dibuat di DEV `pegadevnusare2` 03-12-2025, C/D di produksi) memberi hasil sama; hanya bacaan "kosong = IDPEGA IS NULL" yang membuat klep tak pernah terbuka — `belum terverifikasi` | `DDL\BrowseOpenProteksiEdm_RD.xml` L593-664 |
| A41 | Keempat klep (Type 1–4) **selalu dibaca** berurutan, juga bila tak ada gerbang yang menolak (fidelitas urutan Pega); RI slip → nol klep. Arti `STS_AKSEP = "2"` `belum terverifikasi` (OQ-020) — diport sebagai saringan saja | `services/layanan.go` `bacaKlep` |

## 3. Penghalang — butuh work owner / DBA

1. **DBA (persetujuan manusia):** jalankan **261** lalu **966**; isi `M_EDM_PENGECUALIAN_POLIS` dari literal
   korpus (nilai tidak disalin agent). Tanpa 261 `/periksa` tetap jalan (catatan). Pastikan hak baca
   `ARASAPAS.DETAIL_INVOICE` dari skema aplikasi.
2. **Gerbang 2 `GetListEdm`**: body report tidak ada di korpus (kelas `ASM-FW-GISFW-Work-Endorsement`, saringan
   `FilterTermForEndorsement` = PolicyNo) → selalu `belum-diperiksa`, `lengkap = false`. Perlu definisi report
   atau keputusan sumber (sesudah I3 kasus EDM baru tinggal di `T_WORK_POLIS`; kasus EDM era Pega di tabel Pega).
3. ✅ **Klep `OPENPROTEKSI_EDM` — DIPASANG (W-EDM-07).** Work owner mengirim `DDL\BrowseOpenProteksiEdm_RD.xml`
   (`ASM-FW-GISFW-INT-OPENPROTEKSI_EDM!BROWSEOPENPROTEKSIEDM_RD`, logika `A AND B AND (C OR D) AND F1`: `.POLICY_NO =
   Param.Nopolis`, `.TYPE = Param.Type`, `.STS_AKSEP != "2"`, `.STS_AKSEP IS NULL`, `.IDPEGA = Param.Idpega`) dan
   `DDL\OPENPROTEKSI_EDM.txt` (kolom teks). Diport di `backend/repository/proteksi.go` (`AdaProteksi`); gerbang
   berklep yang terbuka = `diklep` + catatan "Opened by OPENPROTEKSI_EDM", tanpa pesan; tabel tak ada → `tolak` +
   catatan. Isi baris tetap milik DBA / menu lain. Diverifikasi sesi utama atas berkas DDL (L593-664) dan uji
   modul. Teks semula, untuk riwayat: bentuk query `BrowseOpenProteksiEdm_RD` tidak ada → dianggap kosong. Petunjuk
   (belum dipakai, belum terverifikasi untuk EDM): `modul/claimfacin/backend/repository/acuan.go:877`
   (`AdaProteksiPremi`, asal `CekProteksiKlaim`) membaca tabel itu dengan `TYPE = '5' AND STS_AKSEP = '1' AND
   POLICY_NO` — nama kolomnya terlihat; apakah `BrowseOpenProteksiEdm_RD` juga menyaring `STS_AKSEP` tetap
   `[pertanyaan terbuka]` (DBA / work owner).
4. **Prompt list** `EdmTypeNew`, `EndorsementStatus`, `EndorsementInternalRetro` dan isi enumerasi `sumberedm`
   (A34, A36).
5. ~~**A25** zona JVM Pega produksi; **A30** pengisi `Quotation.BusinessOldId` (Life).~~ Keduanya dijawab
   (W-EDM-05, W-EDM-06).
6. Grid portal: sumber daftar kasus endorsement ("Refresh EDM Grid") tidak ada di korpus.

## 4. Uji (dijalankan sesi utama)

| | Hasil |
| --- | --- |
| `go generate ./inti/backend/daftar` · `go build ./...` · `go build -tags db ./...` · `go vet ./modul/endorsmentfacin/...` · `gofmt -l modul/endorsmentfacin` | bersih |
| `go test -count=1 -short ./modul/endorsmentfacin/...` | ok (handlers, repository, services, medis, predikat, logika) |
| `go test ./inti/backend/daftar/...` | ok |
| `go test ./inti/backend/penjaga/... ./cmd/api/...` | gagal **di luar modul**, sama dengan keadaan awal: `TestNolAlamatLayananDiKode`, `TestHandlersTidakMengimporRepository`, `TestTCONolTabelBaru`, `TestTCONolNamaTabelBaruDiKode`, `TestMenuBersihDuaPuluhBarisSatuPerModul`, `TestPanggilanLintasModulTerdaftar`; `grep -ci endorsmentfacin` atas keluaran = 0 |
| `npx tsc --noEmit -p .` | 12 galat, nol di modul ini — `inti/frontend/modul.ts` (commit `19da2b1a`) tanpa `LayarModul` / `tanpaMenu` yang dipakai `frontend/daftar.ts` dan `komiteclaim*/frontend/layar.ts` |
| `npx vitest run modul/endorsmentfacin` + `daftar.modulAktif`, `cssmodul`, `lapisan`, `envelopegalat` | 27/27 |
| `npx vitest run frontend/ inti/frontend` | 28 gagal di luar modul (komite `tanpa baris M_NAV_MENU`, premiumlistlife, treatycontractout); nol menyebut endorsmentfacin |
| Uji mutasi backend (sesi 0b) | 35/36 tertangkap; 1 ekuivalen (tanggal layar diurai UTC vs Jakarta — validasi membaca tanggal kalender) |

## 5. Irisan I3 — kasus `EDM-` lahir (sore 11-10-2026, W-EDM-08…12)

**Yang lahir.**
- Backend (sesi 0b):
  - migrasi `260_seq_work_polis_edm.sql` dan `262_t_edm_facin_kasus.sql` (15 kolom VARCHAR2/DATE + indeks `OLD_POLICY_NO`);
  - `models/kasus.go`, `repository/kasus.go`, `services/kasus.go`;
  - `POST /api/endorsmentfacin/kasus`, yang menjalankan berurutan:
    1. periksa ulang;
    2. OpenCase kering;
    3. SATU transaksi: NEXTVAL → OpenCase → `T_WORK_POLIS` → `T_EDM_FACIN_KASUS`.

  Jawaban: 201 `{caseId, pesan}`; **200 `{caseId:"", periksa}` = ditolak (hasil bisnis, spec 07 §7.2)**; 403 tanpa akun;
  422 masukan; 503 tabel/sequence belum ada.
- G2 kini memeriksa kasus sistem baru (W-EDM-12) → `lengkap` dapat `true`.
- Frontend (sesi utama): tombol OK aktif bila medan wajib terisi ∧ pemeriksaan terakhir lolos ∧ lengkap ∧ tanggal
  diperiksa (padanan `CARIGROUP=='true'` keduanya). Kasus lahir → portal menampilkan "Endorsement case EDM-xxxxx created.";
  ditolak server → hasil pemeriksaan diganti rincian penolakan.

**Keputusan agent — menunggu konfirmasi** (nomor dicek `grep -rnw`, A42–A49 tidak dipakai sebelumnya):

| # | Keputusan | Bukti |
| --- | --- | --- |
| A42 | ID = `EDM-` + 5 digit berpadding (teks work owner "EDM-00001"). ⚠️ Kasus EDM Pega TANPA padding dan sudah melewati 10000 (`EDM-13445` di fixture NB; korpus EDM sendiri 19 nomor unik 537–8989, dua cara sepakat) → mulai `EDM-10000` bentuknya sama dengan IDPEGA Pega. Nomor Pega terakhir hanya dapat dipastikan DBA | `models/kasus.go` `RakitIDKasusEDM`; migrasi 260 |
| A43 | Tautan `DataToEDM` (`EndorsementID`/`EDMHandle`/`EDMHandle2`) TIDAK ditulis: kasus portal bukan entitas tersimpan di sistem baru. `pzInsKey` OpenCase = ID (bentuk kunci instans Pega tidak ditebak, tidak disimpan) | `services/kasus.go` |
| A44 | Dua Create bersamaan atas polis yang sama **tidak dikunci** (G2 dapat lolos di keduanya). Kunci murah butuh mengunci baris tabel warisan Pega — perlu persetujuan | `[pertanyaan terbuka]` |
| A45 | Detail EDM Type diperiksa SESUDAH pemeriksaan (syaratnya membaca jenis bisnis); di luar EdmType 4 ∧ EdmTypeNew 1/2/3 nilainya wajib kosong → 422 bila terisi | `Section/EdmType1/2/3.xml` `pyRequiredWhen` |
| A46 | ✅ **DIJAWAB W-EDM-13** (0 Internal · 1 External · 2 Batal Internal · 3 Batal External; wajib bila bukan Life, bawaan Internal) - penyimpangan di bawah DICABUT. Semula ⚠️ penyimpangan dari korpus: `EndorsementStatus` DIWAJIBKAN korpus bila bukan Life (`WorkPrimaryDetails` L8906 `pyRequiredWhen OutData.pxResults(1).CARI2!='Life'`), tetapi di sini TIDAK wajib dan tidak dirender — prompt list-nya (`D_EnumerationList`, kelas `ASM-FW-GISFW-Data-Enumeration`) tidak ada di korpus. Kasus lahir dengan `ENDORSEMENT_STATUS` kosong | medan wajib dicacah dua cara (urai sel / cacah baris; selisih di L9920 diselesaikan dengan membaca langsung L9880-9975) |
| A47 | ✅ **DIJAWAB W-EDM-14** (1 Yes · 0 No, bawaan No). Semula: Nilai radio `EndorsementInternalRetro` tidak ada di korpus → hanya wajib tidak kosong, maks 10 byte | `WorkPrimaryDetails` L9896-10213 |
| A48 | Tabel/sequence hilang → 503 di KEDUA rute (sebelumnya `/periksa` 500 untuk tabel warisan yang hilang), pesan generik "Endorsment Fac In: a required table or sequence is not in this schema; details are in the server log" — nama objek dan teks ORA hanya di log (atas permintaan sesi utama: pesan semula menyebut migrasi 260/262 juga untuk tabel warisan seperti `JSON_POLIS`) | `handlers/rute.go:120` |
| A49 | Lebar kolom 262 (OLD_POLICY_NO 100 = `JSON_POLIS.NOPOLIS`; REASON 4000; SOURCE_NOTE 1000; SOURCE/STATUS/NB_STATUS 255; kode 10; checkbox 5; IS_CEDING_CONFIRM 20), diperiksa dalam BYTE sebelum INSERT → 422 | migrasi 262; `services/kasus.go` |

**Uji (dijalankan ulang sesi utama):**
- `go build ./...`, `go build -tags db ./...`, `go vet`, `gofmt`: bersih.
- `go test -short ./modul/endorsmentfacin/... ./inti/backend/daftar/...`: ok.
- `go test ./inti/backend/penjaga/... ./cmd/api/...`: 6 gagal yang sama di luar modul; nol menyebut endorsmentfacin / 26x.
- `tsc`: nol galat modul ini.
- vitest modul + penjaga `cssmodul`/`lapisan`/`envelopegalat`/`daftar.modulAktif`: 29/29.
- Uji mutasi I3 (sesi 0b): 23/23.

**Keadaan DEV (dibaca sesi 0b atas perintah work owner "JALANKAN SAJA Migrasi 260 → 261 → 262"):** keempat migrasi
sudah tercatat di `T_MIGRASI` sebelum perintah itu dijalankan (260, 261, 966 pukul 17:12; 262 pukul 17:25) → pelari
menjalankan **nol** pernyataan. Katalog: `SEQ_WORK_POLIS_EDM`, `M_EDM_PENGECUALIAN_POLIS` (+PK), `T_EDM_FACIN_KASUS`
(+PK, indeks); `M_NAV_MENU` endorsmentfacin `DIMIGRASI = 1`. Isi: pengecualian 0 baris, kasus 0 baris. Sisa untuk DBA:
isi `M_EDM_PENGECUALIAN_POLIS`; restart backend.

## 6. Form Create Addendum dirapikan (W-EDM-13…15)

- **Tata letak = tangkapan layar Pega** dari work owner:
  - kiri: Policy Number, Endorsement Source, Endorsement Note, Endorsement Type (+ EDM Type / Detail EDM Type);
  - kanan: Endorsement Date, Endorsement Status, Endorsement Internal Retro.

  Label di kiri medan; OK lalu Cancel di kanan bawah, di luar kotak form; alasan OK mati di kirinya; hasil
  pemeriksaan di bawah.
- **Tema** = salinan tema NB FacIn (token `--edt-*`, kelas `edmf-*`).
- **Tanggal** dd/mm/yyyy = salinan `TanggalDMY` NB, sebab input date bawaan tampil mm/dd/yyyy.
- **Endorsement Status** dirender: wajib bila bukan Life, bawaan Internal.
- **Retro** Yes/No vertikal, bawaan No.
- **A50** `[terverifikasi]` **Biaya (`.AdminCharge`) dan Survey (`.Survey`) DIBUANG dari layar dan dari badan
  `/kasus`.** Kontainernya `pyContainerVisibleWhen = never` (`WorkPrimaryDetails.xml` L9174), jadi Pega tidak pernah
  menampilkannya. Versi pertama form menampilkannya karena laporan pembaca menyebut sel itu `ALWAYS`. Cara itu keliru:
  yang dibaca hanya visibilitas sel, bukan kontainer induknya. Diuji: kontainer terdekat sebelum kedua caption =
  `never`.
  - Backend (sesi 0b): `MasukanKasus` tanpa kedua medan (mengirimnya → 400), kolom `ADMIN_CHARGE` / `SURVEY` ditulis
    NULL.
  - Wajib-isi Status/Retro dinilai SESUDAH pemeriksaan (bukan Life saja, urut Detail → Status → Retro).
  - Diverifikasi sesi utama: build, vet, gofmt bersih; modul 6/6.
- **Uji:** `tsc` nol galat modul; vitest modul + penjaga 31/31.
- **12-10-2026 (W-EDM-16):** panel "Policy checks" dibuang atas perintah work owner. Pesan penolakan tetap di bawah
  medannya; alasan OK mati satu kalimat di samping tombol. Label `PERIKSA` dan CSS panel yang tak terpakai ikut
  dibuang. Uji: `tsc` nol galat modul; vitest modul + penjaga 31/31.

## 7. Polis NB FacIn sistem baru dapat di-endors (12-10-2026, W-EDM-17)

**Sebab "Invalid policy no!"** `[terverifikasi DEV, sesi 0b]`:
- NB FacIn sistem baru menulis `JSON_POLIS.DATA_JSON` minimal `{"IDNewBisnis": …}`
  (`modul/nbfacin/backend/services/selesaipolicy.go:209`).
- `GetBusinessType_Sql` membaca `DATA_JSON.QuotationData.BusinessType`, sehingga hasilnya kosong → langkah 9.

**Yang lahir:**
- Kontrak `kontrak.JenisBisnisFacIn` disediakan sesi `nusantarare-c3`: `inti/backend/kontrak/facin.go` (berkas inti
  bersama, perlu tinjauan tim inti saat PR) dan subpaket nbfacin `services/jenisbisnis`.
- Modul ini memakainya lewat `Membutuhkan`. Diperiksa sesi utama: `cmd/api/rakit.go` merakit SEMUA modul terdaftar
  sebelum menyaring `MODUL_AKTIF`, jadi menonaktifkan nbfacin tidak menggagalkan start modul ini.

| # | Keputusan | Bukti |
| --- | --- | --- |
| A51 | Polis sistem baru = `DATA_JSON` baris PRODKE terakhir memuat `$.IDNewBisnis`. Bagi polis itu: G1 dibaca dari kasus `T_EDM_FACIN_KASUS` ⋈ `T_WORK_POLIS` yang sudah selesai berproduksi; jenis bisnis dari `FACINPRODUCTION.BUSINESSCODE` → `BUSINESS` (GROUPPANEL, OLDID) → kontrak, HANYA bila `DATA_JSON` tanpa BusinessType. Bukan tepat satu baris BUSINESS → UNKNOWN → "Invalid business, please contact IT!". Polis era Pega tetap setia korpus. **Penyimpangan sadar** akibat keputusan 09-10 "sistem baru tidak pakai data_json lgi", disetujui 12-10 | `services/periksa.go`, `repository/gudang.go`, `kasus.go` |
| A52 | Kontrak tidak tersedia (nil, mis. uji tanpa perakit) → gerbang bisnis `belum-diperiksa`, `lengkap = false`; `MasukanGerbang.BisnisBelumDiperiksa` aditif (bawaan = korpus) | `services/gerbangtolak.go` |
| A53 | `StatusEDMSelesaiProduksi = ""` sampai irisan produksi EDM menetapkannya `[pertanyaan terbuka]` → selama itu G1 polis sistem baru selalu "lolos" + catatan "No completed endorsement in the new system" (jujur: polis sistem baru hanya dapat di-endors lewat sistem ini, dan belum ada EDM yang berproduksi) | `models/kasus.go` |

**Uji (diulang sesi utama):**
- build, `build -tags db`, vet, gofmt: bersih.
- Modul 7/7 paket ok; `inti/backend/daftar` ok.
- Gagal di luar modul (nol menyebut endorsmentfacin): 5 penjaga lama, `cmd/api TestPanggilanLintasModulTerdaftar`, dan
  `inti/backend/config TestSetiapEnvTerdokumentasi` (baru, milik c3: env belum di `.env.example`).
- Mutasi (sesi 0b): 12/12.

## 8. POSITION = tahap, POSITION_NOTE = workbasket (12-10-2026, W-EDM-18)

**Yang lahir (sesi 0b; diverifikasi ulang sesi utama):**
- `repository/posisinote.go` = salinan pola `modul/nbfacin/backend/repository/posisinote.go`. Kasus `EDM-` lahir
  dengan `POSITION = 'Policy'` dan `POSITION_NOTE = 'ReasFacInMarketing'` bila kolomnya ada.
- Tanpa kolom (nbfacin 254 belum jalan): SQL persis bentuk W-EDM-11 + satu baris log.
- Migrasi `263_posisi_kasus_edm.sql` mengonversi kasus lama `LINI = 'FACEDM'`.

**Keadaan DEV 12-10 17:35** (dibaca sesi 0b): 254 belum jalan, kolom belum ada, kasus FACEDM 0.
**Pembaruan 12-10** (disampaikan sesi `nusantarare-9d`, laporan work owner): nbfacin `254_position_note` SUDAH jalan di
DEV, kolom `POSITION_NOTE VARCHAR2(255)` ada. **263 SUDAH tercatat di DEV** (`T_MIGRASI` 12-10 17:55, kasus FACEDM 0
→ nol baris berubah); perintah work owner "jalankan" mendapati 263 sudah dijalankan `-migrate` lain — pelari sesi 0b
0 dijalankan, 1 dilewati. ⛔ Ralat: sesi 0b sempat melaporkan "jam server DB ~1 jam di depan mesin" — keliru, karena
disimpulkan dari SATU titik data (catatan 254 pukul 18:47); catatan 263 pukul 17:55 cocok dengan jam mesin. Waktu
catatan 254 yang janggal `belum terverifikasi` (tidak berdampak pada modul ini).

| # | Keputusan | Bukti |
| --- | --- | --- |
| A54 | Deteksi kolom per proses, ketiadaan di-cache dan dicoba ulang tiap 1 menit — sama dengan nbfacin (`posisinote.go:114` `jedaUlangCatatan = time.Minute`). Galat selain ORA-00904 → bentuk lama + log, bukan gagal. ⛔ Ralat: laporan pertama sesi 0b menulis angka nbfacin "tidak dibaca" — dilaporkan tanpa dicek, lalu dicek dan sama | `repository/posisinote.go` |
| A55 | `263_down` hanya mengembalikan baris yang masih di keadaan awal (`Policy` + `ReasFacInMarketing`) | migrasi 263 |

⚠️ **Celah peralihan:** kasus yang lahir di antara 263 selesai dan cache 1 menit kedaluwarsa masih berbentuk lama.
Bila saat itu DEV sudah punya kasus FACEDM: ulangi `UPDATE` 263 secara manual, atau restart backend sesudah `-migrate`.

⚠️ **Celah uji:** penyambungan deteksi di `SisipWorkPolis` butuh Oracle (mutan `|| true` lolos, celah yang sama di
nbfacin; `skemauji.Buka` dikunci claimlife). Uji di DEV sesudah 254: buat satu kasus, cek `POSITION`/`POSITION_NOTE`.

**Uji:**
- Modul ok; build, `build -tags db`, vet, gofmt bersih.
- Penjaga + `cmd/api`: 6 gagal lama, nol menyebut modul.

*Tanpa nama orang, tanpa nomor polis, tanpa data nasabah.*
