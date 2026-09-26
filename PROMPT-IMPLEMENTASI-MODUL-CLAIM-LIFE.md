# PROMPT — sesi implementasi **MODUL Claim Life**: seluruh tiket dalam satu sesi, XML Pega sebagai sumber kebenaran logika bisnis

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief induk **`PROMPT-IMPLEMENTASI-GO-REACT.md`** berlaku seluruhnya. Brief batch
> **`PROMPT-IMPLEMENTASI-SESI-CLAIM-LIFE-TIKET-02-06-BATCH.md`** §1–§10 tetap berlaku untuk rincian
> seam tiket 02, 03, 06, 04, 05 *(§4 di sana)* dan keputusan yang belum diubah di sini; brief ronde
> 5–7 tetap berlaku untuk tiket 14 dan 01 *(Langkah A bila Oracle ada)*.
>
> **Urutan baca sebelum satu baris kode:** brief induk §4–§8 → `APP_RNM\README-BACA-DULU.md` →
> §1–§3 berkas ini → tiket yang sedang dikerjakan **utuh** → §4 bagian tiket itu → **XML rule sumber
> tiket itu** *(§2 cara membacanya)*.
>
> **SESI INI:** modul **Claim Life** — 15 tiket di `.scratch/claim-life/issues/`, dikerjakan
> **berurutan dalam satu sesi** menurut §3, **bukan** satu tiket per sesi. Tiap tiket tetap satu
> commit. Sesi berhenti hanya karena §9, atau karena tiket yang tersisa seluruhnya terkunci gerbang
> yang tidak dapat dibuka executor *(Oracle, keputusan manusia, modul lain)*.

---

## 0. KEADAAN AWAL — 26 September 2026 malam, sesudah `4c10059`

| | Keadaan |
| --- | --- |
| `HEAD` | `4c10059` *claim-life: tiket 03 prasyarat — tiruan tabel peserta dan treaty di skema uji*; sebelumnya `f5cef2b` *(docs bab 10)*, `f439107` *(02-L)*, `753cef2` *(02)*. Working tree bersih |
| Uji yang lulus | `go vet`, `go vet -tags=db`, `gofmt` nol, `go build`; **112 PASS · 0 FAIL**; **24 SKIP** bertag `db` dengan pesan *"skemauji: ORACLE_DSN belum dikonfigurasi"*; `tsc --noEmit`; **5** test JS; `vite build` **88 modul** |
| Migrasi | **belum pernah dijalankan di Oracle mana pun** — DBA belum membuat user kosong *(G1)*. `POOLDATA` **tidak pernah** dipakai untuk `-migrate`, `-migrate-down`, maupun `go test -tags=db` |
| Penjaga yang dikunci test | `CREATE` di migrasi = **20** *(7 tabel · 6 sequence · 7 index)*; pemanggil `skemauji.Buka()` = **8**; `CREATE TABLE` = 7; `kolomSalin` = **24** ekspresi berurutan *(`TestUrutanKolomSalinDikunci`)*; berkas yang **menyebut** `EDMSTATUS` = **2** *(pembacanya dan skema uji)*, yang **menyaring** = 1 |

**Tiket modul Claim Life — keadaan per tiket:**

| Tiket | Judul singkat | Status | AC `[x]`/total | Blocked by *(kolom tiket)* |
| ---: | --- | --- | ---: | --- |
| 01 | Kerangka aplikasi + seam API | `claimed` | 4/14 | — |
| 14 | Skema relasional klaim + migrasi *(PREFACTOR)* | `claimed` | 40/53 | 01 |
| 02 | Register klaim + penomoran | `claimed` | 7/26 | 01, 14 |
| 03 | Baris `AdjustmentList` + Save ke Outstanding | `claimed` | 0/26 | 02, 14 |
| 06 | Validasi DOL per `Type` + `ContentNote` | `ready-for-agent` | 0/5 | 03 |
| 04 | Mesin status per baris | `ready-for-agent` | 0/8 | 03 |
| 05 | Reject Outstanding oleh Admin | `ready-for-agent` | 0/7 | 04 |
| 15 | Hapus klaim — popup, kaskade, baris work | `ready-for-agent` | 0/12 | 14, 03 |
| 07 | Penegakan peran + wewenang per `Type` | `ready-for-agent` | 0/6 | 05 |
| 08 | Tahap Medical Check & Claim Analis + jalur balik | `ready-for-agent` | 0/7 | 07 |
| 09 | Jejak audit setiap transisi | `ready-for-agent` | 0/6 | 08 |
| 10 | Kontrak Komite — penyerahan | `ready-for-agent` | 0/18 | 07 |
| 12 | Efek keluar asinkron + flag lingkungan | `ready-for-agent` | 0/12 | 09 |
| 11 | Terima & tampilkan hasil Komite | `ready-for-agent` | 0/8 | 10, **Komite 05** *(modul lain)* |
| 13 | Migrasi data penuh | `ready-for-agent` | 0/6 | 11, 14 |

---

## 1. ATURAN BARU WORK OWNER — 26 September 2026 malam `[DIPUTUSKAN]`

### 1.1 Pekerjaan per MODUL, bukan per tiket

1. Satu sesi executor menempuh **seluruh tiket modul** menurut urutan §3, tanpa menunggu brief baru
   di antara dua tiket. Yang tetap per tiket: `Status: claimed` → baca ulang XML *(§1.3)* → test dulu
   pada seam → kode → verifikasi penuh → bab `## Implementasi — <tanggal>` → `/code-review` atas titik
   tetap tiket itu → perbaiki → **verifikasi penuh lagi** → **satu commit per tiket**
   *(`claim-life: tiket NN — <judul>` + baris `Tiket:`)*. Nol commit gabungan.
2. **Blocker di dalam tiket tidak menghentikan sesi.** AC yang menunggu Oracle, keputusan manusia,
   atau modul lain diberi tanda di bab Implementasi, tiket tetap `claimed`, sesi **lanjut** ke tiket
   berikutnya — dinyatakan, bukan diam-diam.
3. **Berhenti total** hanya bila: tiket sebelumnya merah dan tidak dapat dipulihkan; `/code-review`
   menemukan cacat AC yang tidak dapat ditutup; ada hal §9; atau seluruh tiket yang tersisa terkunci
   gerbang yang bukan milik executor. Bila berhenti, berhenti **sesudah commit yang hijau**; tiket
   yang belum disentuh tetap `ready-for-agent`.
4. Verifikasi penuh *(vet · vet db · gofmt · build · test · test db · typecheck · test JS · build
   JS)* dijalankan **di akhir tiap tiket** sebelum commit-nya **dan sekali lagi sesudah perbaikan
   review**. Klaim "nol FAIL" di tiket ditulis dari keluaran perintah, bukan dari ingatan.
5. Laporan akhir sesi: satu bab per tiket *(AC ditutup / total, SHA titik tetap dan commit, AC
   terbuka dengan pemiliknya)* + bab telemetri §10.

### 1.2 XML Pega = sumber kebenaran logika bisnis

Kata work owner, 26 September 2026 malam: *"Logika bisnis yang benar ada di XML-nya. Jika ada
kekeliruan di tiket, ubah sesuai yang ada di XML-nya."* Akibatnya:

| # | Aturan |
| ---: | --- |
| a | Untuk **logika bisnis** *(rumus, syarat, urutan langkah, siapa menulis nilai apa, gerbang layar, cabang per `Type`/peran/tahap, pesan galat)*, **XML menang atas tiket, spec, dan brief**. Bila tiket berkata lain, tiket **diubah** — bukan kodenya yang mengikuti tiket yang salah. |
| b | Perubahan tiket ditulis **terang**: blok `### Ralat menurut XML — <tanggal>` di bab Implementasi tiket, satu baris per perbedaan: *teks tiket lama → teks baru*, **bukti** *(path XML + nomor baris di berkas hasil pecahan §2)*. Teks AC yang keliru **disunting di tempat** supaya AC yang dicentang adalah AC yang benar; teks lamanya tetap terbaca di blok ralat. Ini satu-satunya keadaan executor boleh mengubah teks AC. |
| c | Yang **bukan** logika bisnis dan **tidak** diubah: keputusan bentuk penyimpanan dan arsitektur yang spec/tiket tandai `[keputusan work owner]` atau *penyimpangan sadar* — relasional bukan JSON, adjustment menunjuk **peserta**, dokumen per peserta, shared PK, `COVER_KEY`, rujukan lewat ID bukan subscript, nol pemanggilan procedure *(keputusan **o**)*, nol pembacaan `m_product_life.JSONDATA`, spreading dibekukan, jejak audit. Bila sebuah penyimpangan sadar ternyata **bertabrakan dengan logika bisnis** di XML *(bukan sekadar bentuk)*, executor **tidak** memilih sendiri: tulis keduanya di tiket dengan bukti dan tandai `[terbuka — work owner]`. |
| d | Butir `[terbuka]` yang alasannya *"tidak terbaca dari korpus"* **dibaca ulang**. Bila XML ternyata **memutuskan** *(contoh: kedua cabang `PREMIUM_SPREADED_NET` punya precondition `local.Year==1` — cabang mana yang jalan ditentukan aksi precondition-nya)*, executor **mengikuti XML**, menulis buktinya, dan menandai butir itu `[ditutup oleh XML — <tanggal>]`. Yang benar-benar tidak ada di XML *(satuan, arti kode bisnis, kebijakan)* tetap `[terbuka]` dengan pemiliknya. |
| e | Label kebenaran tetap ketat: `[terverifikasi]` hanya untuk yang **dibaca sendiri** dari XML di sesi ini dengan path + baris; turunan yang masuk akal tetapi belum dibaca = `[dugaan]`; yang berasal dari katalog = `[data DBA]`. |
| f | XML **tidak** mengubah larangan §9 dan pagar keamanan brief induk *(nilai nama orang, nomor polis, baris data warisan, kredensial — tidak pernah masuk artefak)*. Yang disalin dari XML hanya **struktur rule**: SQL, ekspresi, nama properti, kondisi. |

### 1.3 Baca ulang XML sebelum setiap tiket — langkah wajib, keluarannya dicatat

Sebelum menulis test pertama sebuah tiket, executor membaca **seluruh rule di tabel "Rule Pega
sumber" tiket itu** *(dan rule yang dipanggilnya: Connect-SQL, activity anak, When, Data Transform,
Section yang memicu)* dengan cara §2, lalu menulis bab
`## Pembacaan ulang XML — <tanggal>` di tiket: satu baris per rule *(path, tag yang dibaca, yang
diambil darinya)*, daftar **perbedaan** dengan teks tiket *(dijadikan ralat §1.2-b)*, dan daftar
pertanyaan yang XML **tidak** jawab. Backend **dan** frontend disusun dari bacaan itu: backend dari
activity/Connect-SQL/When, frontend dari Section/Harness/Flow *(kolom yang tampil, gerbang tombol,
pesan)*.

---

## 2. CARA MEMBACA XML PEGA — metode yang terbukti di sesi ini

**Letak:** korpus **baca-saja** `D:\XML\RNM_BRD\Claim Life\` *(subfolder `Activity\`, `RDBList\`,
`Section\`, `Flow\`, `When\`, `Harness\`, `ReportDefinition\`, `ConnectREST\`)* dan
`D:\XML\RNM_BRD\Komite Claim Life\` untuk rule Komite yang tiket 04/10/11 rujuk. Berkas Pega
ditulis **satu baris raksasa**; pecah dulu ke berkas kerja di scratchpad *(bukan ke repositori)*:

```
sed -e 's/></>\n</g' "D:\XML\RNM_BRD\Claim Life\Activity\<Rule>.xml" \
  | sed -E 's/&lt;/</g; s/&gt;/>/g; s/&quot;/"/g; s/&amp;/\&/g' > <scratchpad>\<Rule>-split.xml
```

Nomor baris di berkas hasil pecahan itulah yang dikutip sebagai bukti *(tulis juga perintah
pecahannya sekali di bab pembacaan ulang, supaya orang lain mendapat nomor baris yang sama)*.

| Jenis rule | Tag yang membawa logika | Cara baca |
| --- | --- | --- |
| **Connect-SQL** *(`RDBList\`)* | `<pyBrowseSQL>`, `<pyOpenSQL>`, `<pySaveSQL>`, `<pyDeleteSQL>` — **bisa banyak baris**; `{Page.Prop}` = parameter bind; `{Page.Prop out}` = parameter keluar | `awk '/<pyBrowseSQL>/{f=1} f{print} /<\/pyBrowseSQL>/{f=0}'` — jangan `grep` satu baris saja; kutip SQL utuh ke tiket |
| **Activity** *(`Activity\`)* | urutan langkah `<pyStepPageReference>` *(`RH_1.pySteps(n)…`)*, `<pyStepMethod>`, `<RequestType>` *(Connect-SQL yang dipanggil)*, pasangan `<PropertiesName>`/`<PropertiesValue>` *(Property-Set)*, `<pyStepsPreCondParamsWhen>` *(precondition aktif)* **beserta tag aksinya di blok yang sama** *(lewati/lompat, saat benar/salah)*, `<pyLoopPropertyName>` | `grep -n` untuk peta langkah, lalu **baca bloknya utuh** dengan `sed -n 'A,Bp'`. ⚠️ `<pyExpression>` di dekatnya sering memuat **teks lama/tidak aktif** yang berbeda dari precondition aktif — jangan mengutip `pyExpression` sebagai perilaku |
| **When** *(`When\`)* | ekspresi kondisi | baca berkas utuh; kecil |
| **Section / Harness** *(`Section\`, `Harness\`)* | `<pyCondition>` *(gerbang tampil/aktif)*, `<pyActivity>` *(tombol memanggil activity)*, nama properti yang ditampilkan, label | frontend meniru **himpunan kolom dan gerbangnya**; label boleh diterjemahkan, kolom dan gerbang **tidak** dikarang |
| **Flow** *(`Flow\`)* | shape `Assignment*` + `pyPosition` *(peran per tahap)*, transisi | tahap dan peran per tahap |
| **ReportDefinition** | filter, parameter *(`Param.*`)*, kolom | pemilihan roster/daftar |

**Fungsi Pega yang muncul dan artinya** *(dari dokumentasi Pega; padanan Go wajib desimal eksak,
nol float)*: `@divide(a,b,n)` = a ÷ b dibulatkan **n** desimal; `@toDecimal(s)` = teks → desimal;
`@replaceAll(s,a,b)`; `@substring(s,i,j)` = potongan indeks **0-based** i..j-1; `@contains(s,t)`;
`@round(x)`; `@if(k,a,b)`; `@addCalendar(tgl, y,M,d,h,m,s,ms)` *(urutan argumen `[dugaan]` — pastikan
dari dokumentasi Pega, tulis di tiket)*; `.pxResults` = baris hasil Connect-SQL; `.pxListSubscript`
= indeks 1-based baris yang sedang diulang; `pyWorkPage` = halaman kasus.

**Yang dicatat ke tiket dari XML:** SQL utuh, rumus dengan pembulatannya, precondition **dengan
aksinya**, urutan langkah, nama kolom/properti. **Yang tidak pernah disalin:** nilai data apa pun
yang kebetulan tertanam *(nama, nomor, alamat, kredensial)*.

**Uji dari XML:** tiap rumus menjadi fungsi murni Go + test dengan **angka contoh yang dihitung
tangan dari rumus XML** *(ditulis di test beserta turunannya)*; tiap precondition menjadi kasus
test; tiap gerbang Section menjadi test `services` **dan** kontrol React.

---

## 3. URUTAN MODUL — rantai ketergantungan, dari kolom "Blocked by"

```
01 ✅claimed ─ 14 ✅claimed ─ 02 ✅claimed ─ 03 ─ 06 ─ 04 ─ 05 ─ 15 ─ 07 ─ 08 ─ 09 ─ 10 ─ 12 ─ 11 ─ 13
```

| Urut | Tiket | Sebab urutan | Yang mengunci sebagian *(tidak menghentikan)* |
| ---: | --- | --- | --- |
| 1 | **03** | prasyaratnya *(tiruan peserta dan treaty)* sudah hijau di `4c10059`; 0/26 AC | `RATE` dari produk per plan *(§4-03 — lintas modul Master Product Name Life)*; `IDR`/`USD` per treaty-year *(`[data DBA]`)*; Oracle |
| 2 | **06** | gerbang simpan di 03 | satuan `@addCalendar` |
| 3 | **04** | status diturunkan dari baris 03 | arti kode `4` *(work owner)* |
| 4 | **05** | memakai mesin 04 + `AUTH_STUB` *(ab)* | — |
| 5 | **15** | hanya butuh 14 + 03; kecil; menutup kaskade lima tingkat selagi pohon segar di kepala | Oracle *(test db)* |
| 6 | **07** | gerbang perlu tindakan nyata *(05)* | sumber peran ADR-U-0030 *(IAM)* — `AUTH_STUB` tetap penunda |
| 7 | **08** | tiap tahap milik peran *(07)* | arti `IsSendtoMedical` |
| 8 | **09** | seluruh transisi harus ada dulu | — |
| 9 | **10** | wewenang kirim-Komite *(07)* | tabel Komite milik tiket 00 Komite Claim Life *(§5 af)* |
| 10 | **12** | kegagalan efek keluar masuk jalur audit *(09)* | alamat layanan luar, kredensial storage *(§9)* |
| 11 | **11** | membaca hasil Komite *(10)*; penulisnya **Komite 05** di modul lain | tanpa penulis, yang dapat dibuat: pembaca + baris lanjutan + status turunan; ditandai |
| 12 | **13** | migrasi data penuh: butuh 11, 14, Oracle, DDL produksi, persetujuan §9 | skrip + rekonsiliasi ditulis dan diuji di skema uji; **tidak** dijalankan ke `POOLDATA` |

Tiket **01** dan **14** *(claimed)* tidak disentuh, kecuali **Langkah A brief ronde 6/7** bila G1
terbuka, dan satu baris di tiket 14 bab Implementasi setiap kali sesi ini menambah langkah migrasi.

---

## 4. RENCANA PER TIKET

Seam tetap tiga *(brief induk §5)*: `repository` ↔ skema uji *(db, SKIP tanpa Oracle)*; HTTP ↔ skema
uji; `services` murni. Uang `Money`/`Ratio` dari teks; nol float; kode tetap teks *(ADR-U-0022)*;
lapisan `handlers → services → repository`; `PeriksaSQL` menolak `COMMIT`; `Qualify` memaksa
`{skema}`. Rincian seam tiket 02, 03, 06, 04, 05 di brief batch §4 **tetap berlaku**; di bawah ini
yang **ditambah atau diralat** oleh pembacaan XML dan katalog sesi ini.

### Tiket 03 — Baris `AdjustmentList` + Save ke Outstanding *(0/26; mulai di sini)*

**XML wajib dibaca ulang** *(§1.3)*: `Activity\SaveOutStandingLife_Act.xml`,
`Section\AdjustmentDetail_Section.xml`, `Activity\SetIndexAdjustmentList.xml`,
`Activity\SpreadingClaimLife_Act.xml`, `RDBList\GetRateRetro.xml`, `RDBList\GetRetroLife_SQL.xml`,
`RDBList\GetProductLife.xml`, `RDBList\GetJsonProductLife.xml`, `Section\RetroDetailClaimLife.xml`,
`Section\DocumentLife.xml` *(dokumen per peserta)*, `Activity\InsertDocument_Act.xml`.

**Yang sudah terbaca sesi ini dari `SpreadingClaimLife_Act` `[terverifikasi]`** *(berkas pecahan
§2; nomor baris di berkas pecahan)* — **ini mengubah tiket 03, terapkan §1.2-b:**

| Langkah | Isi | Baris |
| --- | --- | ---: |
| 1 | `ParamData.CARI1 = pyWorkPage.PolicyDataLife.ProductNameID` | 330–331 |
| 2 | Connect-SQL **`GetJsonProductLife`** *(class `Int-TREATYYEAR_LIFE`)*; baris pertama `pyBrowseSQL`-nya berbunyi `select * from treatyyear_life` — **lanjutannya belum dibaca; baca utuh** | 496 |
| 3 | Connect-SQL **`GetProductLife`**; baris pertama `pyBrowseSQL`: `SELECT M_PRODUCT_LIFE.JSONDATA AS CARI1, PRODUCT_LIFE.RICOMM FROM PRODUCT_LIFE` … — **lanjutannya belum dibaca; baca utuh**. ⚠️ Jadi pembaca JSON produk di Pega adalah rule **ini**, bukan `GetJsonProductLife` — katalog dan brief batch §3 keliru menamai; diralat di katalog | 673 |
| 4 | JSON produk *(`CARI1`)* diurai ke `TempProduct` / `TempProduct1` | 818–911 |
| 5.1 | ulang `PlanList`; bila `pyWorkPage.BusinessName == TempProduct1.PlanList(idx).Name` → **`InputData.CARI3 = .OUTWARDRATEID`** — kunci rate adalah **`OUTWARDRATEID` milik plan yang namanya = `BusinessName` klaim**, dari dalam JSON produk | 1527–1604 |
| 6 | Connect-SQL **`GetRateRetro`**: `SELECT AGE AS CARI1, CONTRACT AS CARI2, GENDER AS CARI3, RATE AS CARI4 FROM POOLDATA.RATE_LIFE WHERE IDUSEDBY = {InputData.CARI3}` — **tanpa `ORDER BY`** | `GetRateRetro.xml` 84 |
| 7 | `.SpreadingList = OutwardList.pxResults` — satu baris spreading per baris hasil langkah 2 *(treaty-year)* | 1987–1988 |
| 8.1 | `local.EMPercent = .EM_PERCENT` **tanpa ÷100**; `TempInputData.AGE = .AGE`; `TempInputData.PERIOD_MM = @if(.PERIOD_MM>=12, @divide(.PERIOD_MM,12,3), 1)` lalu `@round(...)`; `local.Year = @toDecimal(@substring(.GROSS_VALUATION_BEGIN_DATE,6,10)) − @toDecimal(@substring(.BEGIN_DATE,6,10)) + 1` *(tahun polis ke-n; `@substring(…,6,10)` mengambil **tahun** dari teks tanggal berbentuk `DD/MM/YYYY` `[dugaan — pastikan bentuk teks tanggal di halaman Pega]`)* | 2204–2314 |
| 8.2 | `local.Currency = .CURRENCY`; `local.ClaimNet = .CLAIM_GROSS` | 2419–2466 |
| 8.2.1 | per baris `SpreadingList` *(treaty-year)*: `ParamData.CARI2 = .ID` → **`GetRetroLife_SQL`**: `select * from retrocessionlife where idtreatyyear_life ={ParamData.CARI2} order by id asc` → `.RetroLifeList = RetroLife.pxResults` | 2618–2924 |
| 8.2.1.4–7 | ⚠️ **KASKADE KAPASITAS per treaty-year — tidak ada di AC tiket 03; tambahkan.** Bila `local.Currency=="IDR"`: jika `local.ClaimNet <= .IDR` → `.RetrocadedShare = local.ClaimNet`, `local.ClaimNet = 0`; jika `> .IDR` → `local.ClaimNet = local.ClaimNet − .IDR`, `.RetrocadedShare = .IDR`. Sama untuk `"USD"` dengan `.USD`. Sisa klaim mengalir ke treaty-year berikutnya | 3054–3729 |
| 8.2.1.8 | `local.AmountRetroShare = .RetrocadedShare` | 3794–3795 |
| 8.2.1.9.1 | per baris `RetroLifeList` *(reinsurer)*, ulang hasil `GetRateRetro`: **`local.Rate = @toDecimal(@replaceAll(.CARI4, ",", "."))`** — koma diganti titik **sebelum** diurai; precondition yang tampak: `@contains(.CARI3,"U")`, `TempInputData.AGE==.CARI1`, `TempInputData.AGE==.CARI1 && TempInputData.SEX==.CARI3` — **aksi tiap precondition dan sumber `TempInputData.SEX` belum dibaca; baca bloknya utuh** *(`pyExpression` di 4096 memuat varian lama dengan `PERIOD_YY==.CARI2` — abaikan kecuali terbukti aktif)* | 4007–4299 |
| 8.2.1.9.2 | `.RATE = @divide(local.Rate,1,4)` **lalu** `local.Rate = @divide(@toDecimal(local.Rate),1000,10)` → **rate dibagi 1000 (per-mil) sebelum dipakai**; `local.Comm = @divide(@toDecimal(.OVR_COMM),100,5)`; `.Amount = local.AmountRetroShare × @divide(.PERCENTSHARE,100,4)`; `.PREMIUM_SPREADED_GROSS = local.Rate × (1+local.EMPercent) × .Amount`; `local.Comm = GROSS × local.Comm`; **`.PREMIUM_SPREADED_NET = GROSS − local.Comm`**; precondition `local.Year==1` | 4396–4599 |
| 8.2.1.9.3 | rumus yang sama sampai GROSS; lalu `local.Discount = @divide(@toDecimal(.COMMISION),100,5)`; `local.Discount = GROSS × local.Discount`; `local.Comm = (GROSS − local.Discount) × local.Comm`; **`.PREMIUM_SPREADED_NET = GROSS − local.Discount − local.Comm`**; precondition `local.Year==1` | 4664–4907 |

**Ralat tiket 03 yang wajib ditulis** *(§1.2-b)*: **(1)** rumus AC "spreading" ditulis lengkap
dengan `÷1000` dan pembulatan `@divide` *(4, 10, 5, 4 desimal)* — tiket 14 sudah mencatat
"`local.Rate` sudah dibagi 1000", tiket 03 belum; **(2)** kaskade kapasitas `IDR`/`USD` per
treaty-year ditambahkan sebagai AC baru; **(3)** `PREMIUM_SPREADED_NET`: kedua cabang punya
precondition `local.Year==1` — executor membaca aksi precondition-nya; bila XML memutuskan cabang per
tahun polis, AC `[terbuka — Product+UW]` ditutup **oleh XML** *(§1.2-d)* dan kolom diisi; bila tidak
terbaca, tetap `NULL` dan `[terbuka]`; **(4)** `EM_PERCENT` dipakai sebagai **pecahan langsung**
*(tanpa ÷100)* — nyatakan di tiket dan di `models` *(`Ratio`)*; **(5)** kunci rate = `OUTWARDRATEID`
plan yang namanya = `BusinessName`, dicocokkan ke `RATE_LIFE` lewat `IDUSEDBY`, lalu `AGE` *(dan
`SEX` kecuali `GENDER` mengandung `U`)* — bukan "rate treaty".

**Fakta katalog `[data DBA]` yang mengubah kode** *(rinci di
`.scratch/claim-life/KATALOG-TABEL-PESERTA-DAN-TREATY.md`)*: `RATE_LIFE` adalah **VIEW** atas
`M_RATE_LIFE.JSONDATA`, **8 kolom** `VARCHAR2(4000)` *(`ID` `VARCHAR2(10)`)*: `ID, IDUSEDBY, USEDBY,
TYPE, GENDER, CONTRACT, AGE, RATE`; **98.305** baris; `RATE` **90.436 baris berkoma** *(desimal
Indonesia)*, 3.226 bertitik, 4.113 bernilai `0`, nol NULL, **556 baris berspasi tepi** dan **564
baris tidak polos** *(bukan `angka[,.]angka`; 2 memuat koma **dan** titik)* → pembaca `TrimSpace`
dulu, ganti koma, urai, dan **laporkan** sisanya dengan `ID`; `GENDER` `U` 98.004 · `M` 201 · `F` 99;
`AGE` teks angka 0–120; `CONTRACT` teks angka 0–120, NULL 10.203; kombinasi
`(IDUSEDBY, GENDER, AGE, CONTRACT)` **ganda pada 2.693 kombinasi** dan `GetRateRetro` tanpa
`ORDER BY` → **rate di Pega tidak deterministik** bila lebih dari satu baris cocok — pembaca Go
**mengurutkan** dan **melaporkan** ambiguitas, bukan memilih diam-diam; pilihan aturannya
`[terbuka — Product+UW]`. `RETROCESSIONLIFE` di DEV **13 baris**, `PERCENTSHARE` 5–100 dan berjumlah
100 per treaty-year, tanggal `DD/MM/YYYY`, dan **nol** `IDTREATYYEAR_LIFE`-nya cocok dengan
`TREATYYEAR_LIFE.ID` *(2 baris: 2018 dan 2025)* → di DEV join itu kosong; fixture skema uji harus
memakai ID yang **cocok**. `TREATYYEAR_LIFE` di `POOLDATA` **tidak punya kolom `IDR`/`USD`** →
sumber `.IDR`/`.USD` langkah 8.2.1.4 `[data DBA]`: objek `treatyyear_life` yang dilihat koneksi Pega,
dan SQL utuh `GetJsonProductLife`. `PRODUCT_LIFE` = VIEW 36 kolom atas `M_PRODUCT_LIFE.JSONDATA`
*(memuat `OUTWARDRATEID`, `RIRATEID`, `RICOMM`, `OUTWARDCOMM`, `OVR_COMM`, `TREATYNUMBER` tingkat
produk — **bukan** per plan)*.

**Backend** *(menambah brief batch §4-03)*:

| Bagian | Isi |
| --- | --- |
| `repository` | `Treaty.TahunTreaty(ctx)` *(baris `treatyyear_life` — kolomnya mengikuti SQL utuh `GetJsonProductLife`)*, `Treaty.Retrosesi(ctx, idTahun)` *(urut `ID`)*, `Treaty.Rate(ctx, idUsedBy)` *(urut tetap, misalnya `AGE, CONTRACT, GENDER, ID`)* — semua kolom dibaca sebagai teks, diurai `ParseDecimal` **sesudah** `strings.ReplaceAll(",", ".")` *(meniru langkah 9.1.1)*, kegagalan **dilaporkan** dengan nama kolom dan `ID` baris; tiruan skema uji sudah ada untuk `RETROCESSIONLIFE` dan `TREATYYEAR_LIFE` — tambahkan tiruan **`RATE_LIFE`** *(8 kolom `VARCHAR2(4000)`, fixture `UJI-*` dengan **koma** dan titik, satu baris `U`, satu `M`)* dan fixture treaty yang **saling menunjuk** |
| kunci rate | `OUTWARDRATEID` per plan berasal dari JSON produk yang **dilarang dibaca** *(§1.2-c)*. Executor memeriksa `.scratch/master-product-name-life/` *(modul yang memigrasikan produk ke relasional)*: bila tabel plan relasional dengan `OUTWARDRATEID` sudah ditentukan, baca dari sana; bila belum, `services.SumberRate` menjadi **antarmuka** *(`RateUntuk(ctx, idUsedBy, umur, jenisKelamin, kontrak) (Ratio, error)`)* dengan implementasi `SumberRateBelumTersedia` yang menggagalkan spreading secara terang — **bukan** menebak — dan tiket menandai `[terbuka — lintas modul Master Product Name Life]` *(§5 ag)* |
| `services.Spreading.Hitung` | murni; masukan: baris adjustment *(`CLAIM_GROSS`, `CURRENCY`, `EM_PERCENT`, `AGE`, `SEX`, `PERIOD_MM`, `BEGIN_DATE`, `GROSS_VALUATION_BEGIN_DATE`)*, daftar treaty-year *(dengan `IDR`/`USD`)*, retrosesi per treaty-year, rate; keluaran: baris spreading + retro dengan **seluruh** pembulatan XML; kaskade kapasitas; tanpa retrosesi → nol baris, bukan galat; `PREMIUM_SPREADED_NET` menurut hasil bacaan precondition |
| `PohonKlaim.Simpan` | baris adjustment + spreading + retro dalam **satu transaksi**; `STS_REJECT` diisi `0` **oleh aksi Save ke Outstanding**; pewarisan 8 kolom dari baris pertama peserta *(`SetIndexAdjustmentList`)* |
| dokumen | sensus `.DocumentList` *(**ad**)* dari `SaveOutStandingLife_Act` + `Section\DocumentLife.xml` → langkah migrasi **`010`** kolom isi `T_CLAIMLF_DOCUMENT` dengan path XML per kolom; gerbang simpan menyebut **peserta mana** *(pesan XML: "The document hasn't been uploaded person number …")* |

**Frontend:** daftar baris per peserta *(status kata)*; formulir baris + tiga field bank *(dari
`AdjustmentDetail_Section`: himpunan kolom dan gerbang tombol diambil dari `<pyCondition>` dan
`<pyActivity>` yang dibaca ulang)*; tabel spreading per treaty-year dan retro per reinsurer
*(`RetroDetailClaimLife`: `REINSURERNAME`, `PERCENTSHARE`, `Amount`)*; uang dan rate tetap teks.

**Yang dikunci test:** rumus dengan angka contoh hitung-tangan *(termasuk kasus `÷1000`, koma,
kaskade dua treaty-year: klaim lebih besar dari kapasitas pertama)*; ambiguitas rate dilaporkan;
pewarisan 8 kolom; gerbang dokumen menyebut peserta; satu transaksi; test db *(SKIP)*: dua peserta ×
dua putaran = empat baris tertelusur; spreading dibekukan; `SELECT` dokumen biasa.

### Tiket 06 — Validasi DOL per `Type` + `ContentNote` *(0/5)*

**XML:** `Activity\ValidasiDOL_Act.xml` *(59.747 byte — baca kedua cabang beserta aksi
precondition-nya)*, `Activity\SaveOutStandingLife_Act.xml` *(penurunan `ContentNote` dari
`BusinessCode`)*, `Activity\LoadDataPeserta_Act.xml` *(sumber `Type`)*. Brief batch §4-06 berlaku:
`ValidasiDOL(tipe, dol, peserta)` memakai jendela valuasi **yang sudah disalin ke peserta klaim**
*(`T_CLAIMLF_PREMIUMLIST_DETAIL`, teruji pulang-pergi di `4c10059`)*. Yang dibaca ulang dari XML:
argumen `@addCalendar` tiap cabang, bentuk `local.Begin`/`local.Expired`, pesan `"Invalid DOL"`,
dan tabel `BusinessCode → ContentNote` **persis** *(bila berbeda dari `CONTEXT.md` L1–L21, XML
menang, §1.2)*. Satuan pergeseran `+1` tetap `[dugaan: hari]` sampai dibaca dari dokumentasi Pega.

### Tiket 04 — Mesin status per baris *(0/8)*

**XML:** `Activity\SaveOutStandingLife_Act.xml` *(penulis `0`)*,
`Komite Claim Life\Activity\KomitePostAdjustment.xml` *(penulis `1` dan `2`, gerbang
`KomiteCount == KomiteLoop`)*, `Activity\RejectOSClaimLife_Act.xml` *(penulis `2` sisi Admin)*,
`Activity\SetSTS_Reject.xml` *(penurunan status klaim → baris)*. Brief batch §4-04 berlaku:
`StatusBaris` tertutup, `Klaim.StatusTurunan()` dihitung bukan disimpan, `Transisi` hanya dari
Outstanding, kode `4` tidak pernah ditulis *(arti `[terbuka — work owner]`)*, pencerminan
`STS_REJECT`/`ACCEPTED_NO` dua tingkat **dalam satu transaksi**. Dari XML dibaca ulang **arah
pencerminan** `SetSTS_Reject` *(`.STS_REJECT = Primary.STS_REJECT` — klaim → baris atau
sebaliknya)* dan dicatat; bila tiket menyatakan arah yang berbeda, ralat.

### Tiket 05 — Reject Outstanding oleh Admin *(0/7)*

**XML:** `Activity\RejectOSClaimLife_Act.xml` *(menulis `2` ke `.STS_REJECT` **dan**
`…PremiumListDetail(idx).STS_REJECT`)*, `Section\RejectOSClaimLife_Sec.xml`,
`Section\AdjustmentDetail_Section.xml` *(gerbang `pyPosition=='ReasLifeAdmin' && CLAIM_NO!='' &&
.STS_REJECT==0`)*. Brief batch §4-05 berlaku *(pintu `…/tolak`, `WajibPeran`, **ab** `AUTH_STUB`)*.
Dari XML dibaca ulang apakah activity itu juga menyentuh tanggal/kolom lain *(mis. alasan, pelaku)*
— bila ya, ikut ditulis dan ditest.

### Tiket 15 — Hapus klaim: popup konfirmasi, kaskade, baris work *(0/12)*

**XML:** tiket ini penyimpangan sadar berbentuk relasional; rule Pega yang menghapus kasus *(cari di
`Activity\` dan `Section\` kata `Delete`/`Hapus`/`pxDelete`; tulis hasilnya walau nol)*. Pohon yang
berlaku: spec §2b RALAT D — kaskade **lima tingkat** klaim → peserta → adjustment → spreading →
spreading retro, plus dokumen, plus baris `T_WORK_CLAIM`, **satu transaksi**. `services.Dampak(ctx,
id)` menghitung jumlah baris per jenis **sebelum** menghapus; `Hapus(ctx, pelaku, id)` menolak bila
klaim sudah dikirim ke Komite *(`KOMITE_ID` terisi)* kecuali XML/tiket berkata lain; keputusan **y**
*(ORA-02292 bila ada penunjuk dari luar)* tetap. Frontend: popup Ya/Batal dengan rincian angka.
Test db *(SKIP)*: `TestHapusMengkaskadeSampaiCicit` yang ada diperluas sampai retro dan dokumen.

### Tiket 07 — Penegakan peran di `services` + wewenang kirim-Komite per `Type` *(0/6)*

**XML:** `Section\AdjustmentDetail_Section.xml` *(tiga gerbang `<pyCondition>`: Reject
Outstanding; jalur Komite `pyPosition=='ReasLifeSPV' || Type='TP' || Type='TR'`; Save to
Outstanding `pyPosition=='ReasLifeSPV'`)*, `Flow\Register_Flow.xml` *(15 kemunculan `pyPosition`
— peran per tahap)*, `Activity\LoadDataPeserta_Act.xml` *(`pyWorkPage.Type =
pyWorkPage.PolicyDataLife.Type` — sumber otoritatif `PolicyDataLife.Type`)*.
⚠️ Gerbang "Save to Outstanding" di XML menuntut **`ReasLifeSPV`**, sedangkan tiket 03 menulis
"Sebagai ReasLifeAdmin … menyimpannya ke Outstanding" — **periksa ulang di XML** *(Section dan Flow)*
siapa yang boleh menyimpan ke Outstanding; bila tiket 03/05/07 saling bertentangan, **XML menang**,
ralat ketiganya. **Backend:** `WajibPeran` dipanggil di **setiap** fungsi transisi *(test: tiap
fungsi layanan yang mengubah status memanggilnya — penjaga statik)*; sumber peran tetap **ab**
*(`AUTH_STUB`)* sampai IAM/ADR-U-0030; frontend menyembunyikan kontrol sebagai kenyamanan.

### Tiket 08 — Tahap Medical Check & Claim Analis + jalur balik *(0/7)*

**XML:** `Flow\Register_Flow.xml` *(empat tahap: `Assignment2` Input Register, `Assignment1`
Outstanding Claim, `Assignment3` Medical Check, `Assignment4` Claim Analis — baca transisi dan
`pyPosition` tiap tahap)*, `When\IsSendtoAdmin.xml` *(`SendtoAdmin = 1`; tiga titik pengembalian)*,
`When\IsSendtoMedical.xml` *(kondisi dibaca ulang — tiket menandainya `[keputusan work owner]`
karena "tidak terbaca dari tag"; §1.2-d: baca berkasnya utuh, mungkin terbaca)*,
`Section\MedicalCheckClaimLife.xml`, FlowAction `MEDICALCHECK`, `AKSEPTASICLAIMLIFE`.
**Backend:** `models.Tahap` tertutup *(dari Flow)*; `services.Tahap.Ajukan/Kembalikan(ctx, pelaku,
id, ke)` dengan peran per tahap *(07)*; kolom `PY_POSITION`, `SENDTO_ADMIN`, `SENDTO_MEDICAL` di
`T_WORK_CLAIM` *(sudah ada di `001`)*; status baris **tidak** berubah di tahap 2–3 *(ADR-0011)*.
**Frontend:** dua layar mengikuti Section, kontrol kembalikan.

### Tiket 09 — Jejak audit setiap transisi dan jalur balik *(0/6)*

**XML:** `RDBList\UpdateOsAkseptasiClaimLife_sql.xml` *(hanya `CREATEOPNAME` + empat tanggal —
bukti sistem lama **tidak** merekam pelaku/waktu transisi)*, `When\IsSendtoAdmin.xml`,
`When\IsSendtoMedical.xml`, `RDBList\InsertLogServiceClaim.xml` *(`INSERT INTO
pooldata.monitoring_klaim_log` — log layanan, bukan jejak keputusan; **jangan** dipakai sebagai
tabel audit)*, `Activity\RejectOSClaimLife_Act.xml` dan `Activity\SendEmailKlaimLF.xml`
*(`OperatorID.pyUserIdentifier`/`pyUserName` sebagai data pelaku)*. Ini **penyimpangan sadar**
*(ADR-0007)*: XML dibaca untuk **daftar transisi** yang harus direkam, bukan untuk bentuk jejaknya.
**Backend:** tabel jejak baru → langkah migrasi bernomor baru *(sebutkan di tiket 14)*; satu tempat
perekaman di `services` *(penjaga statik: setiap transisi memanggil `Jejak.Rekam`)*; per baris
adjustment; pelaku = `Pelaku.AkunID`. **Frontend:** riwayat pada klaim.

### Tiket 10 — Kontrak Komite: penyerahan kasus *(0/18)*

**XML:** `Activity\CreateKMTLife_Act.xml` *(121.652 byte — sepuluh langkah; baca **seluruh**
Property-Set muatan `childPageKomite.*` dan sisi induk `.IsKomite`, `.KomiteNo`, `.TotalKomite`)*,
`ReportDefinition\FilterEmailKomiteWithLimit.xml` *(roster: parameter `LIMIT_BOTTOM`,
`STS_KLAIM`; filter dibaca utuh)*, `Section\ClaimComite.xml`, `Harness\Committe_Life.xml`.
**Backend:** muatan penyerahan sebagai `models.PenyerahanKomite` *(uang teks + mata uang)*;
`services.Komite.Serahkan(ctx, pelaku, klaimID, adjID)` — wewenang per `Type` *(07)*, baris harus
Outstanding, tiga field bank wajib *(tiket 03 AC 57)*; tulis baris `T_WORK_CLAIM` `KMT-xxxxxx`
*(`RakitPengenalWork(AwalanKomite, …)` sudah ada)* dengan `COVER_KEY` = ID klaim, dan
`T_CLAIMLF_ADJUSTMENT.KOMITE_ID`. ⚠️ Tabel `T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST`, dan kolom
`COVER_KEY` **milik tiket 00 Komite Claim Life** *(tiket 14 relasi 8·9·10)* — periksa
`.scratch/komite-claim-life/issues/00-*.md`; bila belum ada migrasinya di `APP_RNM`, executor
**tidak** membuatnya di modul ini tanpa **af** `[DIPUTUSKAN]` *(§5)*: penulis rekam Komite dipisah
di balik antarmuka `services.PencatatKomite`, dan AC yang bergantung padanya ditandai
`[terbuka — Komite tiket 00]`. Roster: pemilihan meniru filter ReportDefinition atas tabel roster
*(`EMAILKOMITE` — katalognya `[data DBA]` dibaca dulu, agregat saja)*. Nasib `CLMNO`
`[terbuka — work owner]` tetap. **Frontend:** kontrol "Send ke Komite" + pesan.

### Tiket 12 — Efek keluar asinkron + antre-ulang + flag lingkungan *(0/12)*

**XML:** `Activity\InsertGoogleStorage_Act.xml` *(160.027 byte)*, `GetUrlGoogleStorage_Act.xml`,
`DeleteGoogleStorage_Act.xml`, `RDBList\GetTokenStorage_SQL.xml` *(token dari
`pooldata.GET_TOKEN_STORAGE(...)` — **procedure**: keputusan **o** melarang memanggilnya; tulis
`[terbuka — DBA/work owner]` sumber token pengganti, jangan tiru pemanggilannya)*,
`Activity\SendEmailKlaimLF.xml`, `Activity\serviceInsertArasapasClaimLife_act.xml`,
`RDBList\InsertLogServiceClaim.xml`. **Backend:** antarmuka `EfekKeluar` *(unggah, unduh-URL,
hapus, email, Arasapas)* + antrean asinkron dengan antre-ulang dan pencatatan kegagalan ke jejak
audit *(09)*; alamat layanan **hanya** dari `M_LINK_SERVICE` *(ADR-0013; katalognya `[data DBA]`)*
atau env; flag lingkungan di `config` *(`IS_PEGA_PROD` sudah ada)*; implementasi nyata **tidak**
dipanggil di test — stub. ⛔ Menghubungkan endpoint nyata, kredensial storage, atau pengiriman
email sungguhan = §9.

### Tiket 11 — Terima & tampilkan hasil keputusan Komite *(0/8)*

**XML:** `Komite Claim Life\Activity\KomitePostAdjustment.xml` *(penulis `1`/`2`, gerbang
`KomiteCount == KomiteLoop`; tiket ini **membaca**, tidak menjalankannya)*,
`RDBList\UpdateOsAkseptasiClaimLife_sql.xml`, `Activity\SetIndexAdjustmentList.xml`. **Backend:**
pembaca keputusan lewat rantai `KOMITE_ID → T_WORK_CLAIM → T_GENERAL_KOMITE → T_KOMITE_KOMITELIST`
*(bila tabelnya ada — lihat 10)*; `AcceptStatus` dipetakan ke `STS_REJECT` **di batas**; baris
lanjutan mewarisi 8 kolom; status klaim turunan *(04)*. Penulisnya **Komite 05** di modul lain —
yang tidak dapat diuji ujung-ke-ujung ditandai. **Frontend:** hasil dan riwayat putaran.

### Tiket 13 — Migrasi data penuh Claim Life *(0/6)*

**XML:** `RDBList\InsertJsonKlaimLife_sql.xml` *(INSERT flat 51 kolom — sumber pemetaan)*,
`RDBList\UpdateOsAkseptasiClaimLife_sql.xml` *(55 kolom)*, `Activity\SavePesertaClaim.xml`.
**Backend:** pembaca `OS_AKSEPTASI_KLAIM_LIFE` → `BongkarBarisLama` *(sudah ada)* → pohon baru,
per klaim dalam satu transaksi, **seluruh baris** ikut; skrip rekonsiliasi *(cacah baris per
tingkat, jumlah uang per mata uang, tanggal min/max)*; berjalan di skema uji atas tiruan
`OS_AKSEPTASI_KLAIM_LIFE` *(62 kolom, fixture `UJI-*`)*. ⛔ Menjalankan ke data nyata, DDL produksi,
cutover = §9; tiket ditutup hanya oleh work owner + DBA.

### Tiket 02 — Register klaim + penomoran *(7/26, `claimed`)*

**XML** *(bila disentuh lagi)*: `Flow\Register_Flow.xml`, `RDBList\GetSequenceNumber_SQL.xml`.
Tidak disentuh sesi ini kecuali: **o** menjadi `[DIPUTUSKAN]` *(lalu `PenomorCounter` dari
`SUMBER-PENOMORAN-DBA.md` termasuk perakitan format)*, atau tiket 03/06 memerlukan kolom pendaftaran
tambahan *(mis. `SEX`, `AGE`, `PERIOD_MM` peserta untuk spreading — bila `kolomSalin` bertambah,
perbarui `TestUrutanKolomSalinDikunci` dan tiruan skema uji bersama-sama)*. AC 15–17 tetap terbuka.

### Tiket 01 dan 14 — `claimed`, menunggu Oracle

Persis brief ronde 6 §4 Langkah A bila G1 terbuka. Setiap langkah migrasi baru sesi ini *(`010`
dokumen; tabel jejak audit 09; tabel Komite bila **af**)* dicatat satu baris di tiket 14 bab
Implementasi, dan penjaga cacah `CREATE` diperbarui — bukan dilonggarkan.

---

## 5. KEPUTUSAN WORK OWNER — keadaan

| | Keputusan | Keadaan |
| ---: | --- | --- |
| **XML menang** | §1.2 — logika bisnis mengikuti XML; tiket diralat dengan bukti | **`[DIPUTUSKAN]` 26 September 2026 malam** |
| **per modul** | §1.1 — seluruh tiket modul dalam satu sesi | **`[DIPUTUSKAN]` 26 September 2026 malam** |
| aa, z1, ab, ae1 | sequence `SEQ_WORK_CLAIM` + `LPAD(6)`; kolom `CURRENCY` di `002`; `AUTH_STUB`; `CASEID` = pengenal work | diterapkan di `753cef2`/`f439107` |
| o1–o3 | penomoran klaim tanpa procedure; format `<prefix>K<kode>.MM.YYYY.<5 digit>` dirakit aplikasi | `[USULAN]` — tiket 02 AC 2, 3, 7–11 menunggu |
| ac | baris negatif jurnal balik: kolom bertanda | `[USULAN]` |
| ad | kolom isi `T_CLAIMLF_DOCUMENT` dari sensus `.DocumentList` → migrasi `010` | `[USULAN]` — dikerjakan executor di tiket 03 |
| **af** *(baru)* | tabel `T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST`, kolom `COVER_KEY`: dibuat di modul Claim Life *(tiket 10)* bila tiket 00 Komite Claim Life belum berjalan — atau ditunggu | `[USULAN]` — tanpa ini tiket 10/11 sebagian terbuka |
| **ag** *(baru)* | sumber `OUTWARDRATEID` per plan *(kunci rate)* dari tabel produk relasional modul Master Product Name Life, bukan JSON | `[USULAN]` — tanpa ini rate di balik antarmuka `SumberRate` |
| **ah** *(baru)* | aturan pemilihan bila `RATE_LIFE` memuat lebih dari satu baris cocok *(kombinasi ganda 2.693)* — Pega tanpa `ORDER BY` | `[terbuka — Product+UW]`; Go melaporkan, tidak memilih |
| s′, v2, y, j | tetap sebagaimana brief ronde 5–7 | `[USULAN]` / `[terbuka]` |
| k, l | tetap | berlaku |

Baris `[USULAN]` disahkan dengan mengganti kata itu menjadi `[DIPUTUSKAN]` **di berkas ini**. Yang
masih `[USULAN]` saat sesi berjalan **tidak ditebak**.

---

## 6. FAKTA KATALOG `[data DBA — dibaca sendiri, agregat saja]` YANG MENGUBAH CARA MENULIS KODE

Brief batch §3 tetap berlaku *(peserta 66,8 juta baris; `EDMSTATUS`; tanggal valuasi; `STS_REJECT`
0/1/2/4)*, dengan ralat dan tambahan:

| Fakta | Akibat |
| --- | --- |
| `RATE_LIFE` **VIEW** 8 kolom teks atas `M_RATE_LIFE.JSONDATA`; `RATE` berkoma di 92 % baris; nol index | pembaca: `ReplaceAll(",", ".")` → `ParseDecimal`; kunci `IDUSEDBY` *(teks)*; urutkan; laporkan ganda; tiruan skema uji berkolom sama |
| `RETROCESSIONLIFE` VIEW = `M_RETROCESSIONLIFE` ⋈ `M_REINSURANCETYPE` lewat `JSONDATA.TREATYTYPEID`; `M_RETROCESSIONLIFE` **juga** punya kolom bertipe *(`PERCENTSHARE`, `RATE`, `COMMISION`, `OVR_COMM` `NUMBER`; `IDTREATYYEAR_LIFE` `VARCHAR2(10)`; tanggal `VARCHAR2(10)`)* | pembaca tetap lewat **view** *(itu yang Pega baca)*; di DEV **13/13 baris** kolom bertipe **sama persis** dengan nilai JSON-nya *(`PERCENTSHARE`, `RATE`, `COMMISION`, `OVR_COMM`, `IDTREATYYEAR_LIFE`, `TREATYSTARTDATE`)* — membaca kolom bertipe langsung dari `M_RETROCESSIONLIFE` adalah pilihan yang sah tetapi **keputusan work owner** *(`[USULAN]` ai: view, karena itulah yang Pega baca)* |
| DEV: 13 baris retrosesi, 5 treaty-year id *(1000032–1000036)*, **nol** cocok dengan `TREATYYEAR_LIFE.ID` *(1000078 = 2018, 1000079 = 2025)* | fixture uji harus konsisten; di DEV jalur "nol baris spreading" yang akan terjadi |
| `TREATYYEAR_LIFE` 7 kolom, **tanpa** `IDR`/`USD` | sumber kapasitas `.IDR`/`.USD` `[data DBA]`; sampai jelas, `Treaty.TahunTreaty` mengembalikan kapasitas sebagai `*Money` yang **boleh nil**, dan kaskade dengan kapasitas nil **menggagalkan** spreading secara terang |
| `PRODUCT_LIFE` VIEW 36 kolom atas `M_PRODUCT_LIFE.JSONDATA` *(`OUTWARDRATEID` tingkat produk, bukan per plan)* | tidak dibaca langsung *(§1.2-c)*; lihat **ag** |
| Ralat: pembaca `M_PRODUCT_LIFE.JSONDATA` di Pega adalah `GetProductLife` *(baris pertama SQL)*, bukan `GetJsonProductLife` *(`select * from treatyyear_life`)* | tabel §3 brief batch dan katalog diralat; larangan AC 38 berlaku pada **pembacaan JSON produk**, apa pun nama rule-nya |

---

## 7. VERIFIKASI INDEPENDEN `4c10059` — 26 September 2026 malam

Tereproduksi: 5 berkas +387/−16; **112 PASS · 0 FAIL · 24 SKIP** *(pesan seragam)*; 88 modul; tiruan
`M_LIFE_PREMIUM_DETAIL` 24 kolom `kolomSalin` + `EDMSTATUS` dengan tipe katalog, tanpa `KTP`,
fixture `UJI-*` dua peserta *(NULL dan `Batal`)* dan `TestPesertaBatalTidakDapatDidaftarkan`;
`STNC` kini `TO_CHAR` *(cacat 02-L ditutup; selisih tipe `STNC_TREATY VARCHAR2(64)` di `003` tetap
`[terbuka]`)*; `TestUrutanKolomSalinDikunci` mengunci 24 posisi; penjaga AC 29 kini membedakan
**menyaring** dari **menyebut** *(cacah penyebut 2)*; tiruan `RETROCESSIONLIFE` bertipe asli dan
`TREATYYEAR_LIFE`; `Bongkar` membuang keempat tiruan; nol kebocoran. Blocker `RATE_LIFE` yang
dinyatakan executor **ditutup** oleh pembacaan katalog dan XML di §4-03 dan §6.

---

## 8. GAYA KODE — tetap mengikat

Brief batch §6. Tambahan: setiap fungsi yang meniru langkah XML menyebut **rule dan langkahnya** di
komentar kepala *(`// SpreadingClaimLife_Act langkah 8.2.1.9.2`)*; pembulatan ditulis eksplisit
dengan konstanta bernama; frontend `.tsx`, uang dan rate teks.

---

## 9. YANG MEMERLUKAN PERSETUJUAN MANUSIA — berhenti dan tanya

Brief batch §7, ditambah: memanggil procedure apa pun *(`GET_TOKEN_STORAGE`, penomoran)* ·
membaca `m_product_life.JSONDATA` atau JSON produk mana pun · membuat tabel Komite tanpa **af** ·
menghubungkan layanan luar nyata *(storage, email, Arasapas)* atau menyimpan kredensialnya ·
menjalankan migrasi data *(13)* atau `-migrate` ke skema mana pun selain user kosong dari DBA ·
mengubah penyimpangan sadar §1.2-c · `git push`.

---

## 10. TELEMETRI EKSEKUSI — bab wajib di laporan akhir

| Besaran | Cara ukur |
| --- | --- |
| Per tiket: AC ditutup / total; AC diralat menurut XML *(cacah, dengan bukti)*; AC ditutup oleh XML *(§1.2-d)*; AC terbuka dengan pemilik | tiket |
| Per tiket: rule XML dibaca ulang *(cacah berkas, byte)*, perbedaan tiket↔XML ditemukan | bab pembacaan ulang |
| Per tiket: SHA titik tetap dan commit; berkas dibuat/diubah; baris berisi | `git diff --stat` |
| Per tiket: test murni · test db *(SKIP/PASS)* · test JS sebelum/sesudah | `go test -v` |
| Keputusan §5 yang dipakai vs ditunggu, dan AC yang terkena | baris pertama tiap bab |
| Query ke `M_LIFE_PREMIUM_DETAIL` ber-index dan berbatas; nol pembacaan JSON produk; nol procedure | test statik |
| Sub-agen review per tiket: token, panggilan | laporan harness |
| Token sesi utama · biaya · jam dinding | **⛔ tidak diukur** — nyatakan |

---

**Langkah 0 sesi ini:** `git add PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE.md
.scratch/claim-life/KATALOG-TABEL-PESERTA-DAN-TREATY.md` **saja** → commit `docs: brief modul Claim
Life; katalog RATE_LIFE, treaty, dan produk`. ⚠️ Saat brief ini ditulis, working tree **sudah memuat
pekerjaan tiket 03 yang sedang berjalan** *(`010_kolom_t_claimlf_document.sql` + `_down`,
`migrasi.go`, `strukturkolom_test.go`, `STRUKTUR-TABEL-CLAIM-LIFE.md`, `services/adjustment.go` +
test)* — berkas-berkas itu **tidak** ikut commit Langkah 0 dan syarat "`git status --porcelain`
kosong" **tidak** berlaku untuk Langkah 0 ini; ia berlaku lagi sesudah commit tiket 03. Titik tetap
tiket 03 tetap `4c10059`. Sebelum menulis satu baris spreading, baca §4-03 dan terapkan ralatnya.

*Disusun 26 September 2026 malam dari verifikasi independen `4c10059` (uji dijalankan ulang, diff 5
berkas dibaca utuh), katalog `RATE_LIFE`/`PRODUCT_LIFE`/`RETROCESSIONLIFE`/`TREATYYEAR_LIFE` di
instance pengembangan (definisi view dan agregat, nol baris data), dan pembacaan
`SpreadingClaimLife_Act`, `GetRateRetro`, `GetRetroLife_SQL` serta baris pertama SQL `GetProductLife`
dan `GetJsonProductLife` dari korpus. Aturan §1 adalah kata work owner malam ini.*
