# PROMPT — lanjutan 1 modul Claim Life: sesudah tiket 03 (`7aa0e94`), mulai tiket 06

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief modul **`PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE.md`** §1–§11 berlaku **seluruhnya** — aturan
> per modul *(§1.1)*, **XML menang atas tiket** *(§1.2)*, cara membaca XML *(§2)*, urutan *(§3)*,
> rencana per tiket *(§4)*, keputusan *(§5)*, katalog *(§6)*, gaya *(§8)*, persetujuan *(§9)*,
> telemetri *(§10)*, dan verifikasi tiket 03 *(§11)*. Berkas ini hanya memuat **keadaan sesudah
> `7aa0e94`** dan yang baru diketahui untuk tiket 06, 04, 05, 15, 07. Brief induk
> `PROMPT-IMPLEMENTASI-GO-REACT.md` tetap berlaku.
>
> **Urutan baca:** brief modul §1–§3 → §11 → berkas ini → tiket 06 utuh → XML tiket 06 *(§4-06 di bawah)*.
>
> **SESI INI:** Langkah 0 → **06 → 04 → 05 → 15 → 07** → *(bila masih ada ruang)* 08 → 09 → 10 → 12 →
> 11 → 13, satu commit per tiket, berhenti hanya sesudah commit hijau.

---

## 0. KEADAAN AWAL — 26 September 2026 malam, sesudah `7aa0e94`

| | Keadaan |
| --- | --- |
| `HEAD` | `7aa0e94` *claim-life: tiket 03 — spreading, gerbang dokumen, pewarisan delapan kolom*; sebelumnya `436eb39` *(docs Langkah 0)*, `4c10059` *(prasyarat 03)* |
| Working tree | **dua berkas docs belum di-commit**: `PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE.md` *(bab 11 baru)* dan `.scratch/claim-life/KATALOG-TABEL-PESERTA-DAN-TREATY.md` *(ralat kedua `GetJsonProductLife`)*, ditambah berkas ini. Ketiganya masuk **Langkah 0** |
| Uji yang lulus | vet · vet db · gofmt nol · build · **144 PASS · 0 FAIL · 24 SKIP** *("ORACLE_DSN belum dikonfigurasi")* · `tsc` · **5** JS · **88** modul |
| Penjaga yang dikunci test | `CREATE` = **20** *(`010` hanya `ALTER`; kolom `ALTER` terlihat lewat `KolomAlterTambah`)* · pemanggil `skemauji.Buka()` = **8** · `kolomSalin` = **24** · penyebut `EDMSTATUS` = **2** |
| Tiket 03 | `claimed` **11/28** — 14 AC sisa menunggu **ag** *(JSON produk)*, `[data DBA]` `IDR`/`USD`, Oracle, dan izin `011` *(**aj**)*. Bukan pekerjaan sesi ini kecuali §1 menyatakan lain |
| Verifikasi independen `7aa0e94` | brief modul **§11** — seluruh angka tereproduksi; kode precondition `2`/`3` dikorroborasi dari dua langkah lain; satu catatan kode kecil *(pemeriksaan urutan kaskade memakai kapasitas mata uang klaim, sumbernya `TO_NUMBER(IDR)`)* untuk saat 03 disentuh lagi |

**Tiket lain:** 01 `claimed` 4/14 · 14 `claimed` 40/53 · 02 `claimed` 7/26 · 06, 04, 05, 15, 07, 08,
09, 10, 12, 11, 13 `ready-for-agent`.

---

## 1. KEPUTUSAN — keadaan saat berkas ini ditulis

| | Keputusan | Keadaan | Bila `[DIPUTUSKAN]` sebelum sesi |
| ---: | --- | --- | --- |
| **aj** | langkah migrasi **`011`**: `BRANCH_OF_BANK`, `SWIFT_CODE`, `PAYABLE_TO` pada `T_CLAIMLF_ADJUSTMENT` *(layar Pega merujuk enam field bank)* | `[USULAN]` — rekomendasi setujui | dikerjakan sebagai **commit tambahan tiket 03** *(`claim-life: tiket 03 lanjutan — kolom bank`)* **sesudah tiket 06 selesai**, bukan menyela 06; STRUKTUR diralat bertanggal; satu baris di tiket 14 |
| **ag** | sumber `OUTWARDRATEID` per plan dan daftar treaty-year tanpa JSON produk — modul Master Product Name Life | `[USULAN]` ag1 | bila tabel plan relasionalnya **sudah ada** di `APP_RNM`, 14 AC sisa tiket 03 dikerjakan sebagai commit tambahan tiket 03 di **akhir** sesi |
| **ah**, **ak** | rate ganda; mata uang selain IDR/USD | `[terbuka — Product+UW]` | — |
| **af** | tabel Komite dibuat di modul ini bila tiket 00 Komite belum jalan | `[USULAN]` | dipakai tiket 10 |
| **o1–o3** | penomoran | `[USULAN]` | `PenomorCounter` dari `SUMBER-PENOMORAN-DBA.md`, commit tambahan tiket 02 |
| **ac**, **ad**, s′, v2, y, j, k, l | tetap sebagaimana brief modul §5 | | |

Yang masih `[USULAN]` saat sesi berjalan **tidak ditebak**. Sesi ini **tidak** membutuhkan satu pun
keputusan di atas untuk tiket 06, 04, 05, 15, 07.

---

## 2. LANGKAH 0

`git add PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE.md .scratch/claim-life/KATALOG-TABEL-PESERTA-DAN-TREATY.md
PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE-LANJUTAN-1.md` → commit `docs: brief modul bab 11, katalog ralat
kedua, brief lanjutan 1` → `git status --porcelain` **kosong** → uji tanpa Oracle hijau *(144 · 24 SKIP ·
5 JS · 88 modul)* → SHA = titik tetap **tiket 06**.

---

## 3. YANG SUDAH ADA DAN WAJIB DIPAKAI ULANG — jangan menulis kembar

| Ada di kode *(sesudah `7aa0e94`)* | Dipakai oleh |
| --- | --- |
| `models.StatusBaris` *(`TidakDiketahui`, `Outstanding`, `Aksep`, `Ditolak`)*, `String()`, `Diketahui()`, `StatusBarisDariKode`, konstanta `KodeOutstanding`/`KodeAksep`/`KodeDitolak` *(teks `"0"`/`"1"`/`"2"`; `"4"` sengaja tanpa nama)* | 04, 05, 11 |
| `services.TandaiOutstanding(pohon)` — menulis `"0"` **menurut aksi** Save ke Outstanding, bergerbang `PrintFaceClaim` | 04 *(sebagai transisi masuk)*, 06 *(jalur simpan)* |
| `services.PeranSimpanOutstanding` — `ReasLifeSPV` dari XML, **belum punya pemanggil** | 07 *(dan jalur simpan 06)* |
| `services.WajibPeran`, `PunyaPeran`, `Pelaku{AkunID, Peran}`, `handlers.pelakuDari` + `config.AuthStub` *(`AUTH_STUB=true`, ditolak saat `IS_PEGA_PROD`)* | 05, 07, 08 |
| `services.PeriksaDokumenAda`, `PeriksaDokumenLengkap`, `PesertaDokumenTidakLengkap`, `KategoriWajibBelumDiketahui` | 06 *(jalur simpan)* |
| `services.TambahBaris`, `WarisiKolom`, `KolomDiwarisi()` | 11 |
| `services.HitungSpreading`, `PilihRate`, `TahunPolis` | 03 lanjutan |
| `repository.PohonKlaim.Simpan`, `.Hapus(ctx, tx, id, caseID)`, `.AmbilSpreading`, `KlaimLife.AmbilHeader/AmbilPeserta/AmbilBaris` | 15, 04 |
| `models.Peserta` — tanggal valuasi dan WPC sebagai **teks** berformat `utils.TanggalWaktu` *(`2006-01-02 15:04:05`)*; `utils.ParseTanggal`, `FormatTanggal`, `FormatTanggalWaktu` | 06 |
| Rute yang ada: `GET /healthz`, `GET /api/klaim-life/{id}`, `POST /api/klaim-life`, `GET /api/peserta-life` | semua |
| `CONTEXT.md` § *"Daftar BusinessCode → produk → jenis klaim"* — `L1`…`L21` → `DEATH`/`HEALTH`/`CI`/`TPD`/`TI` *(daftar penuh dari work owner; XML hanya menguji `L1`…`L11`)* | 06 |

---

## 4. RENCANA PER TIKET — tambahan atas brief modul §4

### Tiket 06 — Validasi DOL per `Type` + `ContentNote` *(0/5)*

**XML `[terverifikasi]` `Activity\ValidasiDOL_Act.xml`** *(berkas pecahan 1.408 baris; nomor baris
di bawah dari pecahan itu)*:

| Langkah | Isi | Baris |
| --- | --- | ---: |
| 2 | precondition `pyWorkPage.PolicyDataLife.Type=="QR" \|\| =="QP"` *(`WhenTrue=2`, `WhenFalse=3` → jalan hanya untuk QR/QP)*; `local.errmsg = "Invalid DOL"`; `local.DOL = @addCalendar(.DATE_OF_LOSS,0,0,0,0,0,0,0)`; `local.Begin = @CompareDates(local.DOL, .GROSS_VALUATION_BEGIN_DATE)`; `local.Expired = @CompareDates(local.DOL, .GROSS_VALUATION_EXPIRED_DATE)` | 410–584 |
| 3 | precondition `Type=="TR" \|\| =="TP"`; sama, dengan `@addCalendar(.DATE_OF_LOSS,0,0,0,1,0,0,0)` dan `.RETROCESSION_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | 650–804 |
| 4 | jalur galat: precondition `local.Begin==false \|\| local.Expired==true` *(`WhenFalse=3` → jalan bila benar)* — baca isi langkahnya *(pesan/`Page-Set-Messages`)* dan siapa pemanggil activity ini *(cari `ValidasiDOL_Act` di `Flow\`, `Section\`, `Activity\` — di mana ia dipanggil menentukan **di mana** Go memasangnya: pendaftaran atau simpan ke Outstanding)* | 963–968 |

Dua hal yang **bukan** fakta korpus melainkan pustaka Pega, dan karena itu `[dugaan]` yang wajib
ditulis di kode dan tiket sebagai asumsi bernama: **(i)** `@CompareDates(a, b)` bernilai benar bila
`a` **sesudah** `b` *(ketat)*; **(ii)** argumen keempat `@addCalendar` adalah **jam**. Bila keduanya
benar: QP/QR sah ⇔ `BEGIN < DOL ≤ EXPIRED`; TP/TR sah ⇔ `BEGIN < DOL + 1 jam ≤ EXPIRED` *(untuk tanggal
tanpa jam: `BEGIN ≤ DOL < EXPIRED`)*. Tulis konstanta `pergeseranTPTR = time.Hour` dan
`bandingKetat = true` dengan komentar asal, test batas untuk **tiap** `Type` pada `BEGIN`, `BEGIN+1h`,
`EXPIRED`, `EXPIRED−1h`, dan satu kalimat di tiket: *kasus mana yang berbalik bila Product+UW
memutuskan lain*. Satuan dan keketatan tetap `[terbuka — Product+UW]`.

**Sumber nilai:** jendela valuasi dibaca dari **peserta klaim** *(`models.Peserta.ValuasiGross*`,
`ValuasiRetro*` — teks `TanggalWaktu`; urai lewat `utils.ParseTanggal` di **satu** tempat; teks kosong
= galat *"tanggal valuasi peserta kosong"*, bukan waktu nol)*. Nama properti Pega
`RETROCESSION_VALUATION_*` = kolom `RETRO_VALUATION_*` `[terverifikasi katalog]`. `Type` dari
`T_WORK_CLAIM.TYPE` *(nilai `QP`/`QR`/`TP`/`TR`; selain itu **galat**, bukan lolos)*. `DATE_OF_LOSS` —
tentukan dari XML pemanggil di mana ia diisi *(baris adjustment atau klaim)* dan pakai kolom yang
sudah ada di `003`/`004`; jangan membuat kolom baru tanpa bukti.

**`ContentNote`:** XML `SaveOutStandingLife_Act` menguji `L1`…`L11` → `"DEATH"` `[terverifikasi]`;
sisanya dari daftar work owner di `CONTEXT.md` `[keputusan work owner]`. Implementasi: **satu tabel
data** `map[string]string` 21 baris, nol `if` bercabang; kode di luar daftar → galat; test 21 baris +
kode asing. Bila XML menurunkan `ContentNote` untuk kode lain dengan cara berbeda dari tabel
`CONTEXT.md`, **XML menang** *(§1.2)* — ralat bertanggal di tiket dan di `CONTEXT.md` **tidak**
disentuh *(ia dokumen konteks; catat selisihnya di tiket, pemilik work owner)*.

**Seam:** `services.ValidasiDOL(tipe string, dol time.Time, p models.Peserta) error` murni;
`services.ContentNoteDari(kodeBisnis string) (string, error)` murni; keduanya dipasang pada jalur yang
XML tunjukkan *(pemanggil `ValidasiDOL_Act`)*. Pesan ke layar **persis** `"Invalid DOL"` + pengenal
peserta terpisah *(pola gerbang dokumen tiket 03: kalimat XML utuh, data tambahan di medan sendiri)*.
Frontend: pesan pada formulir yang bersangkutan; `ContentNote` tampil sebagai kata.

### Tiket 04 — Mesin status per baris *(0/8)*

Pakai `models.StatusBaris` dan konstanta kode yang **sudah ada**; jangan lahirkan literal `"0"`/`"1"`/`"2"`
baru *(penjaga statik: literal kode status hanya di `models`)*. Tambah `Klaim.StatusTurunan()` *(brief
modul §4-04)*, `services.Transisi(baris, ke)` *(hanya dari Outstanding; final tidak berubah; `"4"` →
`TidakDiketahui`, tidak pernah ditulis)*, `repository.PohonKlaim.PerbaruiStatusBaris` *(pencerminan
dua tingkat dalam satu transaksi)*. **XML:** `SetSTS_Reject.xml` — baca arah pencerminan
*(`.STS_REJECT = Primary.STS_REJECT`: dari klaim ke baris, atau sebaliknya)* dan **aksi
precondition-nya**; `KomitePostAdjustment.xml` *(penulis `1`/`2`, gerbang `KomiteCount == KomiteLoop`)*;
`RejectOSClaimLife_Act.xml`. `TandaiOutstanding` yang ada adalah transisi **masuk** ke Outstanding —
`Transisi` tidak menggantikannya, keduanya memakai konstanta yang sama.

### Tiket 05 — Reject Outstanding oleh Admin *(0/7)*

Brief modul §4-05 berlaku. Gerbang XML `pyPosition=='ReasLifeAdmin' && CLAIM_NO!='' && .STS_REJECT==0`
`[terverifikasi]`. ⚠️ Sampai **o** diputuskan, `PenomorBelumDiputuskan` membuat **tidak ada** klaim
bernomor: jalur HTTP nyata akan menjawab 422 *"klaim belum bernomor"* untuk semua klaim — nyatakan itu
di tiket sebagai akibat **o**, uji `services` memakai `penomorUji` yang sudah ada. Pintu
`POST /api/klaim-life/{id}/adjustment/{adjId}/tolak`; pelaku lewat `pelakuDari` *(**ab**)*.

### Tiket 15 — Hapus klaim *(0/12)*

`PohonKlaim.Hapus(ctx, tx, id, caseID)` **sudah ada** dari tiket 14 beserta
`TestHapusMengkaskadeSampaiCicit` *(db, SKIP)* — perluas ke **lima tingkat** + dokumen + baris
`T_WORK_CLAIM`, satu transaksi; `services.Dampak(ctx, id)` menghitung per jenis **sebelum** menghapus;
tolak bila `KOMITE_ID` terisi *(kecuali XML/tiket berkata lain)*. **XML:** cari rule penghapus kasus di
`Claim Life\` *(`Delete`, `Hapus`, `pxDelete`, `Obj-Delete`)* — tulis hasilnya walau nol; popup React
Ya/Batal dengan rincian angka.

### Tiket 07 — Penegakan peran *(0/6)*

`PeranSimpanOutstanding` mendapat pemanggilnya di sini. Tiga gerbang Section `[terverifikasi tiket 03]`:
Save to Outstanding `ReasLifeSPV`; jalur Komite `ReasLifeSPV || Type=='TP' || Type=='TR'`; Reject
`ReasLifeAdmin && CLAIM_NO!='' && STS_REJECT==0`. Penjaga statik: **setiap** fungsi layanan yang
mengubah status/tahap memanggil `WajibPeran`. Periksa konsistensi tiket 03 *(ralat 1: SPV yang
menyimpan)*, 05 *(Admin yang menolak)*, dan bab "Hasil" tiket 07; bila saling bertentangan, XML menang.

---

## 5. TELEMETRI EKSEKUSI — bab wajib di laporan akhir

Persis brief modul §10, ditambah satu baris: **keputusan §1 yang berubah status sejak brief ini**
*(mana yang `[DIPUTUSKAN]` sebelum sesi, dan commit tambahan mana yang lahir karenanya)*.

---

## 6. PERSETUJUAN MANUSIA

Brief modul §9. Tambahan: langkah migrasi **`011`** hanya bila **aj** `[DIPUTUSKAN]`; tabel jejak
audit *(tiket 09)* dan tabel Komite *(tiket 10, **af**)* tetap memerlukan pernyataan terang di tiket 14.

---

*Disusun 26 September 2026 malam sesudah verifikasi independen `7aa0e94` (brief modul §11):
`ValidasiDOL_Act` dipecah dan dibaca (langkah 2–4 beserta aksi precondition), nama fungsi dan tipe
yang sudah ada dibaca dari kode, letak tabel `BusinessCode` dibaca dari `CONTEXT.md`.*
