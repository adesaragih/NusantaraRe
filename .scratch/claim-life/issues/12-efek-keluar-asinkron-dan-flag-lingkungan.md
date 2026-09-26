# 12: Efek keluar asinkron + antre-ulang + flag lingkungan

**Status:** claimed

**Blocked by:** 09 (jejak audit) — kegagalan efek keluar harus masuk jalur audit, bukan hanya log

## Hasil & nilai pengguna

Sebagai **ReasLifeAdmin**, saya dapat mengunggah dan mengunduh dokumen pendukung klaim; sebagai
**ReasLifeSPV**, anggota Komite menerima notifikasi saat kasus diserahkan. Dan sebagai **pengguna
mana pun**, alur klaim saya **tidak pernah tertahan** karena layanan luar sedang gagal — kegagalan
itu tercatat dan diantre ulang, bukan hilang diam-diam.
*(User story 32–36 di spec)*

## Area codebase

`internal/services` (interface efek keluar + orkestrasi asinkron + antre-ulang),
`internal/repository` (implementasi pemanggil), `internal/config` (alamat layanan + flag
lingkungan), `internal/handlers` (endpoint unggah/unduh), `frontend/` (kontrol dokumen).

## Rule Pega sumber

| Efek | Rule | Identitas |
| --- | --- | --- |
| Unggah berkas | `Claim Life/Activity/InsertGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `INSERTGOOGLESTORAGE_ACT` / `RULE-OBJ-ACTIVITY`, 160.027 byte |
| Ambil URL / hapus | `GetUrlGoogleStorage_Act.xml`, `DeleteGoogleStorage_Act.xml` | kelas sama |
| Token penyimpanan | `Claim Life/RDBList/GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!GETTOKENSTORAGE_SQL` / `RULE-CONNECT-SQL` |
| Email | `Claim Life/Activity/SendEmailKlaimLF.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `SENDEMAILKLAIMLF` / `RULE-OBJ-ACTIVITY` |
| Arasapas | `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` | `[terverifikasi]` **satu-satunya salinan di korpus** — OQ-035 |
| ⚠️ ~~Konversi ke produksi lewat payload JSON~~ — **DIBUANG 2026-09-16** | ~~`Claim Life/ConnectREST/convertJsonNusareToProductionClaimLife.xml`~~ | `[keputusan work owner]` hilir **`SELECT` langsung dari tabel klaim** (spec §2b, §12) |
| Pencatatan layanan | `Claim Life/RDBList/InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` / `RNM!INSERTLOGSERVICECLAIM` / `RULE-CONNECT-SQL` → `INSERT INTO pooldata.monitoring_klaim_log` |

`[terverifikasi]` Token penyimpanan **tidak** diterbitkan aplikasi — ia datang dari Oracle:
`pooldata.GET_TOKEN_STORAGE({UploadDoc.App}, {OperatorID.pyUserIdentifier}, {UploadDoc.Kodestring OUT}, {DocAPI.ResponseMsg OUT})`.

## ADR terkait

**ADR-0008** (asinkron, tidak memblokir, antre-ulang), **ADR-0010** (penyimpanan tetap Google
Storage), **ADR-0013** (**resolusi endpoint lewat `M_LINK_SERVICE`** — meralat **ADR-0004** untuk
alamat layanan keluar), **ADR-0005** (flag lingkungan), **ADR-0007** (kegagalan masuk jalur audit).

## Acceptance criteria

- [x] ⚠️ **Diselaraskan 2026-09-16:** efek keluar tinggal **tiga** — berkas, email, Arasapas.
      Konversi ke produksi **tidak lagi mengirim payload JSON**; hilir membaca **langsung dari
      tabel klaim**. Test yang menemukan payload JSON dikirim keluar **gagal**. *(AC 55 spec;
      penyimpangan sadar 1)*
- [x] Seluruh efek keluar berada **di balik interface**, sehingga dapat diganti dalam test.
- [x] Kegagalan efek keluar mana pun **tidak menahan** transisi status klaim. *(AC 19 spec)*
- [ ] Kegagalan tercatat di **jalur audit**, bukan hanya di log layanan, dan **dapat diantre ulang**.
      *(AC 20 spec)*
- [x] Di lingkungan non-production, klaim **tetap tersimpan**; ketiga efek keluar tidak berjalan.
      *(AC 21 spec)*
- [x] Tidak ada host, endpoint, atau kredensial sebagai literal di kode.
- [x] Kegagalan konfigurasi dapat dibedakan dari kegagalan jaringan, agar antre-ulang tidak berputar
      sia-sia.
- [ ] Dokumen yang sudah diunggah dapat diunduh kembali.
- [ ] Alamat endpoint keluar di-resolve lewat **runtime lookup** ke `M_LINK_SERVICE` dengan kunci
      `(KATEGORI_1, KATEGORI_2)` — untuk Arasapas: `("Klaim", "insertClaimLife")`.
- [x] **Tidak ada URL** sebagai literal, konstanta, **maupun env var** di kode. Yang boleh menjadi
      konstanta hanyalah **kunci kategori**. *(**ADR-0013**)*
- [x] Pemisahan dev–prod terjadi lewat **isi tabel per-database**, bukan lewat percabangan di kode.
- [ ] Bila kunci kategori tidak ditemukan di `M_LINK_SERVICE`, kegagalan **terang-terangan** dan
      masuk jalur audit — bukan diam-diam melewati efek keluar.

## Catatan penutupan (2026-09-14)

**OQ-047 TERTUTUP** `[terverifikasi + data DBA]` — kontrak resolusi endpoint diketahui penuh:
`Claim Life/Activity/GetLinkService.xml` (`ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` /
`RULE-OBJ-ACTIVITY`) melakukan `Obj-Browse` atas `M_LINK_SERVICE` dengan
`.KATEGORI_1 = Param.Kategori_1 AND .KATEGORI_2 = Param.Kategori_2`, mengambil `.URL`, lalu
`Connect-REST`.

`[terverifikasi]` Kunci Claim — Life terbaca di
`Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml`:
`Kategori_1 = "Klaim"`, `Kategori_2 = "insertClaimLife"`.
`[data DBA]` Isi tabel 19 baris; endpoint Life →
`http://10.100.10.75:7315/Nusare-Integration-WS/resources1/restws/NusareClaim/insertClaimLife`.

**OQ-018 TERTUTUP untuk Claim — Life** `[terverifikasi]` — **nol** URL endpoint bisnis ter-hardcode
di modul ini; seluruh `http(s)://` yang ada hanyalah `pyHelpURI` ke `community.pega.com` (132+) dan
3 tautan penampil dokumen Office. Pembedaan lingkungan memakai
`Claim Life/When/IsPEGAPROD.xml` (`@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN`):
`pzProductionLevel = "5"`.

**Flag lingkungan — keputusan spesifikasi tetap berlaku:** flag hanya menggerbangi **efek keluar**,
tidak pernah **penyimpanan**. Ini penyimpangan sadar dari Pega (di sana `IsPEGAPROD` juga
menggerbangi simpan utama); tanpa itu lingkungan non-production tidak dapat dipakai menguji.

`[terbuka]` **OQ-035** — kepemilikan `serviceInsertArasapasClaimLife_act`. **Tidak memblokir.**

## Catatan

`[terverifikasi]` "Tidak memblokir" adalah **paritas**, bukan perubahan: `Claim Life` hanya punya
**pencatatan** (`InsertLogServiceClaim`), bukan gerbang keberhasilan seperti `IsSuccessHitService`
di konteks facultative. Yang **ditambahkan** adalah antre-ulang.

⚠️ Konsekuensi yang diterima (**ADR-0008**): sebuah klaim dapat mencapai keadaan akhir sementara
efek keluarnya masih tertunda.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Pembacaan ulang XML — 27 September 2026

Jendela `D:\XML\RNM_BRD\Claim Life\**\*.xml`. Empat rule dibaca; satu baris per rule.

| Rule (path) | Baris pecahan | Yang dipastikan |
| --- | --- | --- |
| `Activity/GetLinkService.xml` | 371 · 491 · 517 · 393 | `Obj-Browse` disaring `.KATEGORI_1` dan `.KATEGORI_2`, mengambil `.URL` |
| idem | **701 · 705-706** | `ResponLink.URL = linkService.pxResults(1).URL` pada langkah ber-`pyStepsPreCondition` **KOSONG** |
| `Activity/serviceInsertArasapasClaimLife_act.xml` | 497 · **540-541** · 610 | `Call GetLinkService`; `Kategori_1 = "Klaim"`, `Kategori_2 = "insertClaimLife"`; `Connect-REST` |
| `When/IsPEGAPROD.xml` | **164 · 318** | `pxProcess.pzProductionLevel = "5"` |
| `Activity/CreateKMTLife_Act.xml` | **2081 · 2213** | `Obj-Save` lalu `Call SendEmailKlaimLF` — efek keluar MENYUSUL simpan |

**Sensus gerbang keberhasilan — dua cara, keduanya nol.** Cara 1: berkas bernama
`IsSuccess*`/`SuccessHit*` → **0**. Cara 2: berkas yang menyebutnya di isinya → **0**. Klaim tiket
*"`Claim Life` hanya punya pencatatan, bukan gerbang keberhasilan"* **sahih** `[terverifikasi]`.

⭐ **Temuan yang menajamkan tiket:** `InsertLogServiceClaim` hanya punya **satu** pemanggil —
`serviceInsertArasapasClaimLife_act`. Jadi berkas dan email **tidak tercatat sama sekali** di
sistem lama. "Kegagalan masuk jalur audit" karena itu **perbaikan atas dua dari tiga efek**, bukan
paritas seperti bunyi catatan tiket.

**Pertanyaan yang XML tidak jawab:** apa yang terjadi bila `M_LINK_SERVICE` memuat **lebih dari
satu** baris untuk satu kunci — Pega mengambil `pxResults(1)` tanpa memeriksa cacahnya.
`[terbuka — work owner]`

### Ralat menurut XML — 27 September 2026

| Butir | Teks lama | Teks baru | Bukti |
| --- | --- | --- | --- |
| AC lingkungan non-production | *"**keempat** efek keluar tidak berjalan"* | *"**ketiga** efek keluar tidak berjalan"* | AC pertama tiket ini sendiri sudah menetapkan **tiga** sejak penyelarasan 2026-09-16; kata "keempat" tertinggal dari sebelum konversi JSON dibuang |

## Implementasi — 27 September 2026 (tiket 12)

**Status: `claimed`** — **8 dari 12 AC tertutup.** Titik tetap `237055b`.

| Angka | Nilai | Cara 1 | Cara 2 | Sepakat? | Label |
| --- | --- | --- | --- | --- | --- |
| AC tertutup | **8 dari 12** | `grep -c '^- \[x\]'` → 8 | `grep -cE '^- \[[x ]\]'` → 12, dikurangi `grep -c '^- \[ \]'` → 4 | ✅ | `[terverifikasi]` |
| test Go | **214 PASS · 0 FAIL · 32 SKIP** | cacah awalan `--- PASS`/`--- FAIL`/`--- SKIP` dari `go test -tags=db ./internal/... -v` → **246** | `grep -rhoE '^func Test[A-Za-z0-9_]+' internal/ --include=*_test.go \| wc -l` → **246** | ✅ | `[terverifikasi]` |
| test JS | **11 PASS** | `npm test` → `Tests 11 passed` | `grep -c '  it(' src/services/api.test.ts` → 11 | ✅ | `[terverifikasi]` |
| modul frontend | **88** | `npm run build` → `88 modules transformed` | — | **belum punya data** | `[dugaan]` |

⛔ **Dua ralat cara hitung, keduanya milik saya.**
**(1)** Angka `--- SKIP` hanya muncul dengan `-tags=db`; perintah auditnya semula tidak menyebut
tag itu, sehingga pembaca yang menjalankannya apa adanya akan melihat SKIP = 0.
**(2)** Cara 2 semula memakai `sort -u` dan melaporkan selisih 1 — padahal `sort -u` membuang
`TestPagarSkemaUjiMenuntutDuaSyarat`, nama yang memang ada di **dua** paket dan memang dijalankan
dua kali. Tanpa dedup keduanya **sepakat**. CLAUDE.md §4 bab 4a melarang memilih salah satu bila
dua cara berselisih; di tiket 10 saya memilih 245 dan melabelinya `[terverifikasi]`. Itu
pelanggaran, dan ralatnya ditulis juga di tiket 10.
**(3)** Modul frontend: cara keduanya *"8 berkas sumber + dependensi"* bukan cara hitung — ia
tidak menghasilkan 88 dan tidak punya perintah. Ditulis **belum punya data**, bukan dipaksakan.

| Berkas | Isi |
| --- | --- |
| `services/efekkeluar.go` *(baru)* | lingkungan, resolver, tiga efek bernama, antrean, `Penyalur` |
| `services/efekkeluar_test.go` *(baru)* | 11 kasus |
| `services/efekkeluar_statik_test.go` *(baru)* | 2 penjaga |
| `services/komite.go` | `Penyalur` dipasang **sesudah** transaksi |
| `services/komite_db_test.go` | kasus: efek gagal tidak menggagalkan penyerahan |

⭐ **Pemanggil produksi `Penyalur` ada sejak lahir.** Tiket 07 melahirkan `WajibWewenangKomite`
tanpa pemanggil dan AC-nya menggantung satu tiket; kekeliruan itu tidak diulang di sini.

**Penjaga, tiap-tiap dibuktikan gagal lalu dipulihkan — kelima bentuk yang ADR-U-0013 larang:**

| Mutasi | Hasil |
| --- | --- |
| URL sebagai **literal** | merah, alamatnya disebut |
| URL sebagai **konstanta** | merah |
| URL lewat **env var** | merah |
| payload JSON keluar | merah |
| percabangan lingkungan di berkas penyimpanan | merah |
| gerbang lingkungan **dihapus seluruhnya** | merah — *"gerbangnya hilang, bukan pindah"* |

⛔ **Mutasi env var semula memberi hijau palsu.** Sisipan memakai `os.Getenv` tanpa impor `os`,
sehingga **kompilasinya gagal dan testnya tidak pernah jalan** — dan `grep` atas keluaran kosong
terbaca seperti lulus. Diulang dengan impor yang benar; barulah ia merah. Pelajaran yang sudah
tercatat, terulang, dan tertangkap.

⛔ **Penjaga lingkungan semula salah tuduh.** Ia melarang kata `BukanProduksi` muncul sama sekali,
dan langsung menangkap `NewPenyalur(BukanProduksi, …)` — pemasangan yang justru benar: itu
**argumen** penyusun penyalur, bukan gerbang atas simpan. Menyebut sebuah nilai bukan bercabang
atasnya. Kekeliruan yang sama persis dengan `KodeStatus:` di tiket 10; polanya dipersempit ke
percabangan, dan ditambahi tuntutan bahwa gerbangnya memang **ada** di `Salurkan`.

### ⛔ AC yang belum tertutup — 4

| AC | Sebab |
| --- | --- |
| kegagalan tercatat di **jalur audit** dan dapat diantre ulang | separuh: mekanismenya ada dan berjalan — `CatatanEfekGagal`, `LayakUlang`, `Antrean` — tetapi **tabelnya** menunggu butir **am** |
| kunci kategori tak ditemukan → gagal terang **dan masuk jalur audit** | "gagal terang" **tertutup** dan terbukti; "masuk jalur audit" menunggu **am** yang sama |
| dokumen yang sudah diunggah dapat diunduh kembali | menuntut storage nyata **dan** prosedur `GET_TOKEN_STORAGE` — keduanya butuh **persetujuan manusia**. Tidak dikarang |
| alamat di-resolve lewat **runtime lookup** ke `M_LINK_SERVICE` | ⛔ **dicabut centangnya sesudah review.** Kontraknya ada dan diuji — `ResolverEndpoint`, `KunciArasapasLife`, `AlamatLayanan` beserta penjaga hasil kosongnya — tetapi satu-satunya implementasi adalah `ResolverBelumDiputuskan`. Pembacaan `M_LINK_SERVICE` menuntut **persetujuan manusia**: ia tabel milik DBA. Readernya **sengaja tidak ditulis** — reader tanpa pemanggil adalah jebakan `WajibWewenangKomite` tiket 07 sekali lagi |

⚠️ **Yang sengaja TIDAK dibangun:** endpoint unggah/unduh. Ia hanya dapat menjawab 501 hari ini,
dan pintu menuju ruang yang belum ada adalah permukaan tanpa isi. Dinyatakan, bukan dilewatkan
diam-diam.

### Hasil /code-review — titik tetap `237055b`

Kedua review menemukan hal yang tidak akan saya temukan sendiri, dan keduanya **membangun
elakannya lalu menjalankannya** — bukan menduga.

**Spec — 1 cacat logika + 4 centang tidak sah + 2 kode mati.**

| Temuan | Tindakan |
| --- | --- |
| ⛔ **CACAT: antrean berputar selamanya.** `LayakDicobaUlang` mengecualikan 3 galat — dan melewatkan justru **ketiga galat yang benar-benar diproduksi di produksi** (`…BelumDisetujui`), yang keadaannya permanen sampai seorang manusia bertindak. Persis yang AC larang | **diperbaiki**, keenam galat permanen didaftar; dibuktikan merah pada mutasi |
| flag lingkungan **tidak pernah terpasang**: `LingkunganDariFlag` nol pemanggil produksi, `komite.go` menyetel `BukanProduksi` harfiah → efek keluar tidak berjalan **di mana pun**, sehingga AC-nya benar secara hampa | **diperbaiki** — `Service.DenganLingkungan` dipasang di `cmd/api` dari `config.IsPegaProd`, plus test perilaku yang memeriksa nilainya MENGALIR |
| AC runtime lookup dicentang padahal hanya antarmuka | **dicabut centangnya**; alasannya ditulis di tabel AC terbuka |
| `HasilSalur` dibuang utuh di satu-satunya pemanggil produksi — bertentangan dengan komentarnya sendiri *"dilaporkan, bukan ditelan"* | **diperbaiki** — `HasilEfekTerakhir()` |
| `ErrJaringan` nol produsen; `DariTingkatProduksi` mengurai kolom yang aplikasi ini tidak pernah baca | **keduanya dibuang**; pengetahuan `pzProductionLevel = "5"` hidup di bab XML dengan nomor baris |

**Standards — 3 elakan penjaga TERBUKTI + 4 pelanggaran sensus.**

| Temuan | Tindakan |
| --- | --- |
| ⛔ **penjaga lingkungan dielakkan** dengan menaikkan gerbang ke variabel: `prod := …Produksi(); if !prod { return nil }` — menyimpan dilewati di non-produksi, **hijau** | penjaga tekstual **dicabut**, diganti penjaga **perilaku** |
| ⛔ **pemeriksaan "gerbangnya masih ada" dikalahkan sebuah literal teks**: `var _ = "if !p.lingkungan.Produksi()"` sementara gerbang aslinya dihapus — efek keluar menyala di setiap lingkungan, **hijau** | diganti penjaga perilaku berpasangan |
| ⛔ **penjaga URL dielakkan** tiga cara sekaligus: penyambungan, `Sprintf("%s://%s",…)`, dan `os.Getenv("ARASAPAS_TUJUAN")` yang namanya tidak memuat satu pun kata yang dicari; `.env` tidak pernah dibaca | pola diganti `://` (bukan `http://`); env jadi **aturan arsitektur** — nol `os.Getenv` di luar `internal/config`; berkas `.env*` ikut dibaca |
| penjaga payload hanya memeriksa **nama** — `json.Marshal` + POST lolos | diganti pemeriksaan **perilaku**: nol klien HTTP keluar |
| §4a: memilih 245 saat dua cara berselisih | **diralat** di tiket 12 **dan** tiket 10 — selisihnya cacat cara hitung saya (`sort -u`), bukan selisih sungguhan |
| §4 rule 9: perintah audit SKIP tanpa `-tags=db`; "88 modul" cara keduanya bukan cara hitung | **diperbaiki**; modul ditulis **belum punya data** |
| §4 rule 1: kutipan storage/email tanpa nomor baris | **diturunkan ke `[dugaan]`** — tiket 12 memang belum membacanya baris demi baris, dan `[terverifikasi]` hanya untuk yang dibaca sesi ini |
| bau: Data Clumps, Mysterious Name (`Lingkungan.Produksi()` vs konstanta `Produksi`) | `CatatanEfekGagal` menyematkan `MuatanEfek`; metode jadi `AdalahProduksi()` |

⭐ **Pelajaran yang saya tulis di kodenya:** pola atas teks sumber selalu dapat dielakkan oleh
penulisan ulang yang setara. Tiga elakan terbukti mengalahkan penjaga tekstual saya; yang
menggantikannya penjaga **perilaku** — memanggil kodenya dan memeriksa apa yang terjadi — dan
ketiga elakan itu kini tertangkap semua, apa pun bentuk penulisannya.

**lanjut dari sini:** tiket 12 selesai. Empat AC menunggu: dua butir **am**, satu persetujuan
storage, satu persetujuan membaca `M_LINK_SERVICE`. Berikutnya tiket 11.

