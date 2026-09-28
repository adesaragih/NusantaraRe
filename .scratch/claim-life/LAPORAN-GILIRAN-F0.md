# Laporan giliran F0.3 → A3 Register

**Aturan giliran** *(brief lanjutan 8 §1)*: nol teks ke manusia sesudah tiap commit paket; catatan
kemajuan ditulis di berkas ini; satu pesan ke manusia hanya sesudah A3 kelompok Register ter-commit,
atau sesudah berhenti sah *(gerbang yang hanya manusia dapat buka, disebut namanya)*.

Titik tetap: **`9186078`** *(docs: brief lanjutan 8)*.

---

## Verifikasi F0.2 yang diminta work owner

Dijalankan sebelum giliran ini dimulai; hasilnya sejalan dengan tabel §0 brief 8.

| Klaim | Diperiksa | Hasil |
| --- | --- | --- |
| 12 berkas, 5 baru | `git show --numstat 5b49954` | ✅ 12 diubah, 5 baru, +1.082/−353 |
| axios dilepas | `package.json` + impor di `src/` | ✅ nol; sisa hanya kata di komentar dan nama uji |
| build 41 modul, 158,83 kB | `npm run build` | ✅ |
| 76 uji JS; Go 261 · 0 · 34 | vitest; `go test -tags=db -v` | ✅ |

---

## F0.3 — Shell dan menu

**Isi:** Shell referensi (sidebar terlipat + laci ponsel, topbar, menu profil, palet Ctrl+K,
`PagarGalat`), `lib/lipatMenu.ts` + `components/KelompokMenu.tsx` + `hooks/useHalaman.ts` disalin,
`MENU` + `MODUL_LAIN_TERLARANG` di `labels.ts`, `App.tsx` F0.2 dibuang.

**Menu — dua butir, keduanya berbukti:**

| Butir | Bukti |
| --- | --- |
| kelompok `Claim Life` | `Flow/Register_Flow.xml:270` `<pyWorkTypeName>ClaimLife</pyWorkTypeName>` |
| `Inbox Claim Life` | `[tidak ada di korpus]` — kosakata `InboxPremiumList` / `Struktur_InboxClaimLife`; nol harness portal di ekspor |
| `Register` | VERBATIM `Flow/Register_Flow.xml:155` `<pyLabel>Register</pyLabel>` |

**Penjaga yang dibuktikan menggigit:** butir menu ketiga (`Detail & Tutup`) -> dua uji jatuh; menu
modul lain (`NB Treaty In`) -> uji jatuh.

**Ralat penjaga saya sendiri:** larangan nama modul memeriksa SELURUH teks berkas dan menuduh
`ShellProps` karena memuat "Prop". Disempitkan ke LABEL MENU - penjaga yang menuduh nama tipe akan
membuat orang mengganti nama tipenya, bukan memperbaiki menunya.

**Telemetri:** uji JS 76 -> 97, modul build 41 -> 44, Go tidak tersentuh. XML dibaca ulang: 1
(`Register_Flow.xml`, baris 155 · 270 · 268 · 313 · 343 · 358). Taksiran token +-40k (taksiran).

---

## at/au — kolom TAHAP dan TGL_CREATE  ·  SHA `cdd7b98`

Migrasi 016; `serahTerimaSah` dari peta PERAN menjadi peta TAHAP; `TahapDanPeran` membaca keduanya
dalam satu kueri; penjaga optimis `UPDATE` memakai `NVL(TAHAP, …)`; pendaftaran menulis `TAHAP`,
`TGL_CREATE`, `CREATE_OP`, dan memakai SATU jam. Ralat tiket 08 dan 14; STRUKTUR diperbarui.

XML dibaca ulang: `Register_Flow.xml` (358 · 343 · 268 · 313), `InputOSClaimLife.xml` (21349 ·
21404 · 21433 · 21839 · 21863), `InboxPremiumList.xml` (736). Go 261 → 263 PASS.

---

## F0.4 backend — rute kotak masuk  ·  SHA `6b1618c`

`GET /api/klaim-life?tahap=&halaman=&ukuran=`; `repository/inbox.go` + `services/inbox.go` +
`handlers/inbox.go`. Worklist (pribadi) vs workbasket (bersama) dari `Register_Flow.xml`
1511 · 1508 · 1300 · 1351 · 993 · 1072 · 1119 · 1198 — kedelapan baris diperiksa satu per satu.

⛔ **RALAT ANGKA DI PESAN COMMIT `6b1618c`**: pesannya menulis "Go 263 → 274 PASS"; angka
sebenarnya **272**. Dicatat di sini alih-alih menulis ulang sejarah. Cacah sesudah paket ini tetap
272 PASS · 0 FAIL · 34 SKIP.

XML dibaca ulang: `InboxPremiumList.xml` (592 · 721 · 733 · 736 · 751 · 765 · 942),
`Register_Flow.xml` (delapan baris perutean).

---

## F0.4 frontend — halaman Inbox

`pages/InboxClaimLife.tsx` + `ambilKotakMasuk` di `api.ts` + gaya. Empat tab = empat assignment,
judul VERBATIM `pyTaskName`, terjemahan Indonesia DI SAMPING. Tab yang perannya tidak dipegang
TIDAK dirender. Label kolom VERBATIM `InboxPremiumList.xml`, dan ujinya MEMBUKA korpus pada baris
yang tiap label sebut.

Uji JS 97 → 114; build 45 modul.

---

## F0.5 — panduan  ·  SHA `16d567f`

`PANDUAN-MENJALANKAN.txt` bab 4 diperbarui: layar masuk stub, Shell + menu dua butir beserta
buktinya, Inbox empat tab, pembedaan antrian PRIBADI vs BERSAMA. Cacah uji JS 5 -> 114.

Satu kalimat yang sengaja ditulis terang: tab Admin akan KOSONG sampai pemakai mendaftarkan klaim
sendiri (worklist = `CREATE_OP` miliknya). Tanpa itu orang melaporkannya sebagai bug.

---

## A3 kelompok 1 — REGISTER

**Backend:** `repository/rujukan.go` (tiga pembaca berbatas), `services/rujukan.go` (gerbang,
himpunan tertutup), `handlers/rujukan.go` + rute `GET /api/rujukan/{jenis}?cari=`.

**Frontend:** `components/PanelDataPolis.tsx` (himpunan medan PERSIS section), label `REGISTER`
dengan 15 nomor baris, `cariRujukan` di `api.ts`, `RegisterKlaim.tsx` memakai label VERBATIM.

**Dua temuan yang mengubah kode:**

1. Penjaga arsitektur `TestHandlersTidakMengimporRepository` menangkap `handlers/rujukan.go`
   mengimpor repository. Diperbaiki dengan ALIAS di services, bukan struct kembar.
2. Urutan pemeriksaan salah: jenis rujukan yang salah ketik dijawab "database belum dikonfigurasi".
   Permintaan diperiksa lebih dulu, baru infrastrukturnya.

**XML dibaca ulang (6 berkas):** `InputRegisterClaimLife.xml` (15 baris label + `pyValue`/`pyFormat`
tiap medan), `BrowseCedingCoLife_RD.xml`, `BrowseBusinessLife_RD.xml`, `BrowseMarketingOfficer_RD.xml`,
`InputOSClaimLife.xml`, `Register_Flow.xml`.

**Telemetri:** Go 272 -> 274 PASS · 0 FAIL · 34 SKIP; JS 114 -> 136; build 46 modul.
Taksiran token giliran ini **+-1,1 juta** (taksiran; angka sejati tidak terlihat dari dalam sesi).

---

## F0.6 — form login ditiadakan  ·  SHA `daa9b2c`

Halaman `Masuk` + ujinya dibuang; identitas dari env Vite (`VITE_AUTH_STUB`, `VITE_STUB_PELAKU`,
`VITE_STUB_PERAN`). Bawaan peran KETIGA-tiganya. Tombol `Keluar` dibuang. JS 136 → 129.

## Ralat av  ·  SHA `1255809`

⛔ **RALAT PEMBACAAN SAYA SENDIRI.** Empat medan (`Type`, `Marketing Officer`, `Ceding`,
`Class of Business`) `pyReadOnly` true / `Read-only` di Pega, terikat `.PolicyDataLife.*`. Ronde
sebelumnya saya membangunnya sebagai dropdown + rute + tiga pembaca tabel — **kode mati di tiga
lapis, nol `.tsx` memakainya**. Sebabnya: blok kontrol berdiri **sebelum** labelnya di DOM, dan
jendela pembacaan saya menghadap **ke depan**. Itu grep dengan langkah tambahan, bukan pohon.

Dibuang: `repository/rujukan.go`, `services/rujukan.go` (+uji), `handlers/rujukan.go`, rutenya,
`api.ts cariRujukan`. Tiket 02 + PARITAS §6 + `OQ-untuk-tim.md` (4 pertanyaan).

⛔ **RALAT ANGKA DI PESAN COMMIT `1255809`**: pesannya menulis "Go tetap 274"; angka sebenarnya
**272** — membuang `rujukan_test.go` menghapus 2 uji. Dicatat di sini, bukan dengan menulis ulang
sejarah.

## A3 Outstanding — butir aw

Rute `POST /api/klaim-life/{id}/tahap/{tujuan}` (tujuan KATA), `handlers/tahap.go`,
`pages/OutstandingClaimLife.tsx`, `OUTSTANDING` + `TOMBOL_OS` di `labels.ts` dengan 16 nomor baris.

**Pertanyaan pohon terjawab:** `SendtoAdmin_Act` b338 berprasyarat
`pyPosition=="ReasLifeMedicalAdvisor"` (WhenFalse=3 = lewati) padahal layar Outstanding dipegang
Admin → kedua tombol **tidak menulis apa pun** di Pega. Cacat rule warisan; ditiru **maksudnya**,
cacatnya dilaporkan OQ-C. Ralat tiket 08.

**Telemetri:** Go 272 → 274 PASS · 0 FAIL · 34 SKIP; JS 129 → 150; build 46 modul. XML dibaca ulang
(4 berkas): `InputRegisterClaimLife.xml` (blok kontrol read-only), `InputOSClaimLife.xml` (16 baris),
`Register_Flow.xml` (8 penyambung), `setDetailClaim_act.xml`. Taksiran token **±1,4 juta** giliran
ini (taksiran).

---

## Giliran lanjutan 10 — paket 1: Register(2) backend (`1c97724`)

### Cara membaca yang dipakai, dan kenapa

Alat baca pohon ditulis lebih dulu *(scratchpad `pohon.py`, `langkah.py`, `kontrol.py`,
`jendela.py`)*, dan **diuji jangkarnya** sebelum dipercaya: versi pertama melaporkan
`<pagedata>` di baris 809 padahal ia di baris 3 — nomor barisnya datang dari batas buffer,
bukan dari XML. Diganti `xml.parsers.expat` + `CurrentLineNumber`, lalu dicocokkan dengan
`sed -n '3p;5p'`. Alat yang nomor barisnya meleset menghasilkan bukti yang meleset.

Dua kali alat itu sendiri salah dan ketahuan karena hasilnya ganjil:
`pyStepsActivityName` *(bukan `pyStepsMethod`)* yang memuat metode, dan
`pyStepsRepeatDefHasRepeat` bernilai **`EMBEDDED`**, bukan `true` — sehingga kedua loop
`DeletePesertaClaimLife` sempat tidak terlihat.

### Empat temuan yang MEMBANTAH brief lanjutan 10 §1

| # | Brief menyebut | XML sebenarnya | Bukti |
| --- | --- | --- | --- |
| 1 | `DeletePesertaClaimLife` = hapus peserta, rute `DELETE .../peserta/{id}` bergerbang tahap+pemegang | Activity ini **tidak menghapus apa pun**. Ia mengindeks ulang `.AdjustmentList(*).IndexPremiumList = Local.IndexPremium` di dalam dua loop `EMBEDDED`, lalu `Obj-Save`. Penghapusan barisnya dikerjakan **klien** | `Activity/DeletePesertaClaimLife.xml` b225 *(`Local.IndexPremium = .pxListSubscript`)*, b337, b417/b583 `ULANG(EMBEDDED)`, b441 `Obj-Save` |
| 2 | kelompok **Register** | Tombolnya berdiri di layar **Outstanding** | `Section/InputOSClaimLife.xml` b17909, b18039. `SelectAllClaimLife_act` juga: b16633, b16710 |
| 3 | *(tidak disebut)* | Tombol `DELETE` hanya **terlihat** saat `CLAIM_NO` kosong, dan kliknya menjalankan **dua** aksi berurutan: `deleteRow` di klien lalu `refresh` yang memanggil activity-nya. Konfirmasi **dimatikan** | `pyLabel = DELETE` b17865 · `pyAction = deleteRow` b18017 · `pyNextGenGridDeleteConfirm = false` b18021 · `pyAction = refresh` b18032 · `pyActivity` b18039 · `pyUserData/pyCondition` b18082 |
| 4 | `GET /api/peserta-life?pl=&nama=` *(dua kriteria)* | **TIGA** kriteria; sertifikat **tidak** di-uppercase | `RDBList/GetPesertaClaim_sql1.xml:85`; `@toUpperCase` hanya pada `CARI3` b405 |

Temuan 3 dibaca dengan **menaiki** pohon dari `pyActivity` ke blok pembungkusnya — bukan jendela
maju. `pyCondition` berdiri di `<pyUserData>` milik **sel**, sedangkan `pyModes` yang memuat
tombolnya berakhir di b17938, jauh **sebelum** kondisi itu. Jendela maju dari label `DELETE` akan
berhenti sebelum b18082 dan melewatkan syarat tampilnya — pengulangan persis kekeliruan butir av.

### `UploadCSVClaimLife_Act` — tiga langkah, nol logika bisnis

`Page-Remove TempWorkPage` b250 · `Page-New TempWorkPage` b340 · `Call pxUploadCSVResults` b488.

Ia **memanggil mesin bawaan platform**, bukan aturan Nusantara Re: tidak ada pemetaan kolom, tidak
ada validasi, tidak ada penulisan peserta di rule ini. Karena itu **tidak ada rute unggah** yang
dibangun di paket ini. Bila kelak diperlukan, yang harus digrilling lebih dulu adalah pemetaan
kolom CSV-nya — dan itu **tidak ada di korpus** *(OQ-F)*.

### Yang dibangun

| Lapis | Isi |
| --- | --- |
| `repository/pilihpeserta.go` **baru** | `UmurPeserta`, `ShareNusantaraReTeks` — dua aturan; ujinya mengunci **asimetri** keduanya |
| `repository/pesertapolis.go` | `kolomSalin` 24 → **28** ekspresi *(+`SHARE_NUSANTARA_RE_GROSS`, `AGE`, `ENTRY_AGE`, `CURRENT_AGE` — **bahan**, bukan isi)*; `sqlCariPeserta` dipisah jadi fungsi murni; `Cari` menerima sertifikat + nama |
| `repository/kolompeserta.go` | `AGE` masuk daftar tulis/baca — kolomnya **sudah ada** di migrasi 003, hanya tidak pernah terisi |
| `models/klaimlife.go` | `Peserta.Umur` **teks**, sebab kosong bukan nol *(ADR-U-0027)* |
| `services/peserta.go`, `handlers/register.go` | penyaring diteruskan; `sertifikat=` dan `nama=` |

**Nol migrasi baru.** Nomor `017`–`029` masih utuh.

### Dua penjaga diperbaiki — keduanya dibuktikan menyala dulu

1. **`TestUrutanKolomSalinDikunci`** memakai `strings.Contains`, sehingga
   `SHARE_NUSANTARA_RE_GROSS` lolos sebagai `SHARE_NUSANTARA_RE` dan `ENTRY_AGE` lolos sebagai
   `AGE` — dua pasangan yang justru baru ditambahkan. Diganti pencocokan nama **persis**.
   *Dibuktikan:* menukar kedua kolom share → penjaga menyebut posisi 20 **dan** 24.

2. **`TestQueryTabelPesertaSelaluBerindexDanBerbatas`** membaca sampai backtick pertama, sehingga
   SQL yang kini dirakit terbaca terpotong dan ia **menuduh query yang berpagar lengkap**.
   Diperlebar ke seluruh fungsi, dan komentar dibuang sebelum dipindai.

   ⚠️ Lalu ditemukan **lubang yang lebih dalam, dan dinyatakan terbuka secara tertulis**: untuk SQL
   rakitan, pemindaian teks tidak dapat membedakan penyaring **wajib** dari penyaring
   **bersyarat**. `CERTIFICATE_NO LIKE` yang hanya terpasang bila kotaknya terisi sudah cukup
   memuaskannya. *Dibuktikan:* mencabut `PL_NUMBER = :1` → penjaga **tetap hijau**.

   Karena itu penjaga kini **menuntut** setiap SQL rakitan punya uji yang menyebut nama fungsinya.
   *Dibuktikan:* berkas uji perakitnya disingkirkan → penjaga menuduh, menyebut `sqlCariPeserta`.

   Yang benar-benar menjaga bentuknya: **`TestSQLCariPesertaSelaluBerpagar`**, dibuktikan **merah**
   untuk kedua cacat *(penyaring dicabut; batas hasil dicabut)*.

### Ralat cara menghitung uji

Perintah yang benar adalah `go test ./... -tags db`. Tanpa tag `db`, seluruh `*_db_test.go`
**tidak dikompilasi**, dan larinya melaporkan `0 SKIP` — bukan 34. Itu hijau yang tidak menguji
apa yang dikira diuji. Dicatat di `29a7ebf`.

### Telemetri paket 1

| Ukuran | Nilai |
| --- | --- |
| Commit | `1c97724` *(sesudah `29a7ebf` Langkah 0; `e3d537a`+`53f9ab0` butir ax)* |
| Go | **274 → 278 PASS · 0 FAIL · 34 SKIP** |
| `go vet` | bersih |
| Berkas Go baru | 3 *(`pilihpeserta.go`, `pilihpeserta_test.go`, `caripeserta_test.go`)* |
| Migrasi baru | **0** |
| Mutasi pembuktian penjaga | **5** *(2 tukar kolom, 2 cabut pagar SQL, 1 singkirkan berkas uji)* — seluruhnya dipulihkan |
| XML dibaca sebagai pohon | 5 activity + 2 section + 1 Connect-SQL |

---

## Giliran lanjutan 10 — paket 2: layar Register (`a203601`)

| Hal | Isi |
| --- | --- |
| Label baru | `Certificate No` b16277 **`<pyLabelFieldValue>`** · `Search` b16553 **`<pyLabel>`** |
| Tag TIDAK disamakan | diperiksa: **nol** `<pyLabelPreview>Certificate No` di seluruh berkas. Keduanya berdiri di daftar `medanCari` terpisah dengan tagnya sendiri — memaksanya masuk daftar `pyLabelPreview` berarti menguji tag yang tidak ada, lalu melonggarkan ujinya sampai lulus |
| Ralat tombol | `Find Insured` b7057 **bukan** penjalan pencarian; ia pembuka panelnya. Yang memanggil activity adalah `Search` b16553. Sebelumnya layar memakai label b7057 untuk tombol yang menjalankan pencarian |
| Huruf besar | dikerjakan **server** *(b405)*, bukan layar — mengubah ketikan orang saat ia mengetik membuat kotaknya terasa rusak, dan aturannya tetap satu rumah |
| Kotak kosong | **tidak** dikirim sebagai parameter kosong; ujinya mengunci bentuk URL dan dibuktikan merah ketika kosong ikut terkirim |
| `lib/pilihSemua.ts` | penjungkit **tiga** keadaan b247 — `''` → `'true'` → `'false'` → `'true'`. Keadaannya **teks**: `''` berarti belum pernah ditekan dan harus mencentang, `'false'` berarti baru dilepas. Boolean tidak dapat membedakannya |

⚠️ **`pilihSemua.ts` BELUM PUNYA PEMANGGIL, dan itu dinyatakan, bukan didiamkan.** Tombolnya
berdiri di layar **Outstanding** *(`InputOSClaimLife.xml` b16633, b16710, b24489)*, dan grid
pesertanya baru dibangun di kelompok Detail & Tutup. Diperiksa: **nol** `Select All` dan **nol**
`IsAccept` di `InputRegisterClaimLife.xml` — memasangnya di Register akan mengarang tombol yang
Pega tidak punya. Aturannya ditulis sekarang bersama ujinya supaya tidak lahir dari ingatan ketika
gridnya tiba. **Work owner boleh memveto** dan menundanya sampai gridnya ada.

**Telemetri:** JS 150 → **169** · tsc bersih · build **46** modul · 1 mutasi pembuktian
*(label dikarang `'Cari Peserta'` → penjaga menyala)*.

---

## Giliran lanjutan 10 — paket 3: Detail & Tutup (1) (`27c5e56`)

Dua validasi tanggal berambang produk. Keduanya **berbentuk sama persis**, jadi keduanya memakai
satu fungsi — `SelisihHari` + `PenandaBatasHari`.

| Validasi | Dari → ke | Ambang | Penanda |
| --- | --- | --- | --- |
| `ValidasiClaimReceived_Act` | `DATE_OF_LOSS` → `CLAIM_RECEIVED_DATE` b561 | `MAXEXPIREDCLAIM` b540 | `.MAXCLAIM_RECEIVED` b582 |
| `ValidasiSTNC_Act` | `EFFECTIVE_DATE` → `RECEIVED_DATE` b598 | `MAXDATARECEIVE` b577 | `.STNC` b627 → kolom `STNC_CLAIM` *(migrasi 002)* |

⚠️ Propertinya `.MAXDATARECEIVED` **ber-D**, kolom sumbernya `MAXDATARECEIVE` **tanpa D**. Mudah
tertukar, dan tertukarnya tidak berbunyi.

Hasilnya **penanda, bukan bool**: kosong = sah; terisi = tidak sah **dan** isinya tanggal yang
melanggar, `dd/MM/yyyy`. Bool akan membuang keterangan *mana* tanggalnya.

`@addCalendar` dengan **seluruh** argumen nol *(b498, b519, b526, b556)* tidak menggeser apa pun —
ia menormalkan nilai menjadi tanggal. Selisihnya karena itu dihitung dari tanggalnya saja;
menghitungnya berjam membuat dua tanggal berselisih 25 jam terbaca **1** hari dan yang 23 jam
terbaca **0**.

### Ambangnya tidak ditebak

`RDBList/GetProductName.xml:84` membacanya dari tabel produk warisan lewat
`PolicyDataLife.ProductNameID`. `PolicyDataLife` **menunggu modul PremiumList Life** *(av)*, dan
tabelnya `[data DBA]` OQ-001. Fungsinya karena itu **menerima** ambang, dan ambang **kosong
menjawab galat** — menganggapnya nol akan menandai hampir setiap klaim tidak sah *(ADR-U-0027)*.

### Penjaga yang menyala dan TERNYATA BENAR

`TestMasterViewTidakDisentuh` menuduh berkas baru itu menyebut nama tabel warisan. Pemeriksaan:
namanya ada di **string pesan galat**, bukan di komentar — dan pesan galat mendarat di log,
sedangkan nama objek warisan tidak pernah boleh masuk log atau artefak. **Yang diperbaiki pesannya,
bukan penjaganya.** Keterangan lengkapnya pindah ke komentar kepala berkas, tempat yang memang
untuk itu.

**Telemetri:** Go 278 → **284 PASS · 0 FAIL · 34 SKIP** · `go vet` bersih · **0** migrasi baru.

### Yang BELUM dikerjakan dari §2, dan sebabnya

`ClaimLifeDetailGCNM` *(Save Adjustment 22590)*, `ShowEditClaimLife` + `EditDateClaimLife_Section`,
`SetSTS_Reject`, `SetIndexAdjustmentList`, `CloseClaim`, `DocumentLife` — belum dibaca sebagai
pohon di giliran ini. Ia **tidak** dikerjakan setengah dari ingatan; kelompok berikutnya
membacanya lebih dulu, sebagaimana lima activity paket 1 dibaca.

---

## Giliran lanjutan 11 — paket 0: envelope galat (`b18de5b`)

**DUA** cacat sejenis, bukan satu, dan keduanya menelan pesan yang backend kirim dengan benar.

### Cacat 1 — kunci envelope tidak sama

`handlers.galat` menulis `{"galat": …}` sejak tiket 01 *(`klaimlife.go:47`)*; `services/api.ts`
membaca `o.error`. Akibatnya **setiap** pesan backend — 401, 403, 409, 503 — jatuh ke teks bawaan
*"Permintaan ditolak backend"*.

⛔ **Sebab ia bertahan, dan ini yang paling layak dicatat: ADA ujinya, dan ujinya ikut keliru.**
`klien.test.ts:147` menyuapkan `{"error": …}` lalu menuntut kalimatnya lolos — dan ia **hijau**,
sebab klien memang membaca `error`. Yang diuji bukan kontrak dengan backend melainkan kontrak klien
dengan dirinya sendiri. Komentar `api.ts` pun menyebut envelope `{"error"}`, sehingga **kode,
komentar, dan uji ketiganya salah bersama-sama** — dan kesalahan yang konsisten tidak berbunyi.

### Cacat 2 — pembaca berbentuk transport yang sudah dibuang

`pages/RegisterKlaim.tsx` punya `pesanGalat` **sendiri** berbentuk axios
*(`e.response.data.galat`)*, padahal axios dibuang di F0.2. Bentuk itu tidak pernah cocok dengan
apa pun, sehingga setiap penolakan backend di layar Register tampil sebagai *"Gagal menghubungi
server."*

⛔ Itu bukan sekadar pesan yang hilang melainkan pesan yang **menyesatkan**: ia menuduh jaringan
padahal backend menjawab, dan menjawab dengan sebab yang tepat. Orang yang membacanya akan
memeriksa koneksi, bukan datanya.

### Yang dibangun

`error` **sengaja tidak ikut diterima**. Menerima kedua kunci akan menambal gejalanya dan
menyembunyikan sebabnya: sejak itu kedua sisi tidak pernah dipaksa bertemu lagi.

Kontraknya dikunci **dua sisi**, sebab cacatnya tidak berbunyi di satu sisi mana pun — backend
benar, klien benar menurut komentarnya sendiri, dan hanya **pertemuannya** yang salah:

| Berkas | Yang dikunci |
| --- | --- |
| `internal/handlers/envelopegalat_test.go` | kunci `galat`, **tepat satu** kunci, nol `error` |
| `frontend/src/services/envelopegalat.test.ts` | kalimat sampai; `error` **ditolak**; kosong tidak menggantikan teks bawaan |
| `frontend/src/services/envelopegalat.guard.test.ts` | nol pembaca berbentuk axios di seluruh `src`; kunci hanya dibaca di satu berkas |

Seluruh penjaga **dibuktikan menyala** lebih dulu pada cacat aslinya: `api.ts` dikembalikan membaca
`error` → 4 uji merah; pembaca axios dikembalikan → penjaga axios merah.

**Telemetri:** Go 285 → **287** · JS 169 → **179** · tsc bersih · build 46 modul.

---

## Giliran lanjutan 11 — paket 1: gerbang Close Claim (`f3443c6`)

`Section/CloseClaim_Section.xml` b1081 `Close Claim` → `pyActivity` b1101 `ProtectCloseClaim_act`.
Tombol itu **tidak punya aksi lain**: gerbangnya **adalah** aksinya.

| Baris | Langkah | Isi |
| ---: | --- | --- |
| 236 | `Property-Set` | `ProtectLife.CARI1 = ""` |
| 370 · 812 | `Property-Set` + `ULANG(EMBEDDED)` | atas seluruh peserta; b398 `idx`, b445 `name` |
| 484 | prasyarat **b608 `.STS_REJECT!=1`** | `WhenTrue=2` lanjut / `WhenFalse=3` lewati → b511 tandai, b557 susun pesan |
| 646 | `Page-Set-Messages` | prasyarat b756 `ProtectLife.CARI1==1` |
| 836 | `Call FinishAssignment` | prasyarat **b992 `ProtectLife.CARI1==""`** |

⚠️ **`STS_REJECT == 1` berarti DIAKSEP**, bukan ditolak — dipastikan dari `models.KodeAksep = "1"`
yang sudah ada di kode, bukan dari nama kolomnya. Namanya menyesatkan dan itu warisan.

⛔ Prasyaratnya `!= 1`, **bukan `== 0`**. Peserta berstatus **ditolak** *("2")* juga bukan 1, jadi
ia menahan pula. Memperlakukan "ditolak" sebagai "selesai" akan menutup klaim yang barisnya belum
diputus ulang. Ada ujinya.

**Seluruh** penghalang dilaporkan: `Page-Set-Messages` b646 adalah **sublangkah** b370, jadi ia
berjalan di dalam loop. Melaporkan satu saja memaksa pemakai menutup berulang kali dan menemukan
satu penghalang baru tiap kali.

Klaim **tanpa peserta boleh ditutup** — ditiru apa adanya dan dinyatakan; kita tidak menambahkan
larangan yang XML tidak punya.

### Penyimpangan sadar: pesannya memakai nomor sertifikat, bukan nama

Pega menyusun kalimatnya dari `.NAME_OF_INSURED`. Kita **sengaja tidak menyimpan** nama tertanggung
*(`kolomSalin` tiket 02)*; mengambilnya kembali dari tabel warisan hanya demi sebuah pesan berarti
membatalkan keputusan itu. Nomor sertifikat menunjuk baris yang sama persis dan tidak memuat nama
siapa pun; **sisa kalimatnya VERBATIM**. Ada uji yang akan gagal bila seseorang menambahkan medan
nama.

### Yang DINYATAKAN belum ada

Sisi `Call FinishAssignment` *(b836)*. Tahap tujuan sesudah tutup belum dibaca dari `Flow/`, dan
menebaknya akan memindahkan kasus ke tempat yang salah **tanpa satu pun galat**. Karena itu rutenya
`GET`, bukan `POST` — dan metodenya dikunci uji, supaya tombol tidak menjanjikan lebih daripada
yang ia lakukan.

### Dua bacaan XML untuk kelompok berikutnya

| Aksi | Temuan | Pemiliknya |
| --- | --- | --- |
| `SetIndexAdjustmentList` b542–b743 | putaran **baru mewarisi delapan angka** dari `.AdjustmentList(1)`, bukan menghitung ulang | Akseptasi |
| `SetSTS_Reject` b233 | halaman langkahnya **`.DiagnoseList`** *(bukan AdjustmentList)*, loop `EMBEDDED`, `.STS_REJECT = Primary.STS_REJECT` | Medis — daftar diagnosisnya belum ada |

**Telemetri:** Go 287 → **295 PASS · 0 FAIL · 34 SKIP** · `go vet` bersih · **0** migrasi baru.

⚠️ `go vet` menemukan apa yang test **tidak** temukan: `t.Context()` menuntut go1.24 sedangkan modul
ini go1.22. Testnya hijau, vet-nya merah. Diganti `context.Background()`.

---

## Giliran lanjutan 11 — paket 3: lima total menunggu rule yang hilang (`639e97a`)

### Temuan: rujukan menggantung pada **angka uang**

`Section/ClaimLifeDetailGCNM.xml` memanggil **`CheckTotalAdjustmentClaim` sepuluh kali**, tetapi
activity itu **tidak punya satu pun berkas rule di seluruh korpus**. Diperiksa dua cara:
`find -iname "*CheckTotalAdjustment*"` → nol hasil; `grep -rl` → satu-satunya berkas yang
menyebutnya adalah section itu sendiri.

Kesepuluh pemanggilan menempel pada **lima medan total uang**, dua per medan:

| Medan | `pyLabelPreview` | Pemanggilan |
| --- | ---: | --- |
| `Total Share Nusantara Re` | b20914 | b20970, b21091 |
| `Total Sum Insured` | b21201 | b21262, b21377 |
| `Total Sum Reasured` | b21488 | b21546, b21664 |
| `Total Share Retro` | b21775 | b21836, b21951 |
| `Total Claim Amount` | b22063 | b22120, b22238 |

⛔ **Kelimanya tidak dijumlahkan sendiri.** "Total" terdengar seperti penjumlahan kolomnya, tetapi
pertanyaan yang menentukan tidak terjawab: baris yang **mana** yang ikut — seluruhnya, atau hanya
yang `IsCheck`? Termasuk yang `STS_REJECT = 2`? Jalur tolak **mencabut `IsCheck`*
*(`RejectOSClaimLife_Act`)*, jadi jawabannya berpengaruh nyata — dan ini angka uang *(ADR-U-0003)*.

Kelima medannya **tetap tampil** dengan penanda dan menyebut nama rule-nya. Paritas yang tampak
lengkap padahal tidak adalah paritas yang tidak akan dicari lagi *(pelajaran butir av)*.

**Dua uji menjaganya dari dua arah**, keduanya dibuktikan: menambahkan angka pada salah satu total
→ merah; dan uji kedua memeriksa bahwa rujukannya **ada** (10×) **dan** berkasnya **tidak ada** —
bila ekspornya kelak dilengkapi, uji itulah yang gagal, dan gagalnya **kabar baik**.

Dilaporkan **OQ-H**; PARITAS baris 18 diralat *(Koreksi 2 benar bahwa ia tidak ada, tetapi
melewatkan bagian terpenting: ia **dirujuk**)*.

**Telemetri:** JS 179 → **193** · build 46 → **47** modul *(panelnya terpakai, bukan kode mati)*.

---

## Giliran lanjutan 11 — paket 4: `Edit Date` terhubung ke layar (`d10427e`)

### Temuan: rute tanpa pemanggil

`PUT /api/klaim-life/{id}/peserta/{pesertaId}/tanggal-kejadian` ada di backend **sejak tiket 06**,
lengkap dengan `ValidasiDOL` dan ujinya. Tetapi **nol pemanggil di React**. Jalurnya karena itu
tidak pernah dapat dijalankan orang — dan itu **tidak berbunyi di uji mana pun**: backend hijau,
layar hijau, fiturnya tidak ada.

Sumbernya: b14115 `Edit Date` → b14144 `pyLocalAction ShowEditClaimLife` → `EditDateClaimLife_Section`;
validasinya `ValidasiDOL_Act`, dipanggil section yang sama di b11177 dan b11298.

⛔ Kontrolnya berdiri **per peserta**, bukan per klaim: `ValidasiDOL_Act` berkelas
`Int-LIFE_PREMIUM_DETAIL` dan menempelkan galatnya pada `.DATE_OF_LOSS` **peserta**. Galatnya pun
disimpan per peserta — satu kotak galat bersama akan menuding peserta yang salah.

`TanggalKejadian` kini menyeberang di JSON peserta. Sebelumnya medannya sengaja tidak dikirim,
sehingga kotak tanggal selalu terbuka **kosong** — dan kosong terbaca *"belum diisi"* padahal
mungkin sudah. Penjaga ikut dipasang di uji marshaller itu: **nama orang tetap tidak menyeberang**,
sebab medan pada marshaller mudah bertambah satu per satu.

### Balapan yang ditutup sebelum sempat terjadi

Dua pengubahan beruntun dapat berlomba: yang kedua tiba lebih dulu, lalu `ambilKlaimLife` milik
yang **pertama** mendarat dan menimpa layar dengan potret kedaluwarsa — tanggal yang baru saja
disimpan tampak hilang, dan pemakai mengetiknya lagi. Kotaknya kini terkunci selagi permintaannya
terbang, **per peserta** *(penanda bersama akan mengunci kotak peserta lain tanpa sebab yang
terlihat)*.

⚠️ Uji sempat merah karena **ujinya**, bukan kodenya: `new Response('', {status: 204})` ditolak
konstruktor — spesifikasi Fetch menuntut 204 berbadan `null`.

**Telemetri:** Go 295 → **296 PASS · 0 FAIL · 34 SKIP** · JS 193 → **199** · build 47 modul.

---

## Ralat giliran lanjutan 11 — sesudah telaah dua sumbu

| # | Klaim ronde pertama | Yang benar | Bukti |
| ---: | --- | --- | --- |
| 1 | Panel lima total berdiri di tingkat **klaim**, berjudul "Total klaim" | **Milik PESERTA.** `ClaimLifeDetailGCNM.xml` berkelas `Int-LIFE_PREMIUM_DETAIL` dan medannya terikat properti ber-TITIK pada halaman itu | b84; `.TotalShareRNM` b20921, `.TotalSumInsured` b21208, `.TotalSumReasured` b21495, `.TotalShareRetro` b21782 |
| 2 | "dua pemanggilan per medan *(`postValue` + `refresh`)*" | **Satu** pemanggilan per medan; `postValue` b20940 tidak membawa `pyActivity`. Kedua kemunculan adalah `refresh` yang sama, terserialisasi di `pyModes` dan `pyActionSets` | b20940, b20952, b21080 |
| 3 | "tombol `Close Claim` tidak punya aksi lain" | **Dua** aksi sekali klik: `refresh` → activity **dan** `closeContainer`; ada pula konfirmasi yang terlewat | b1101, b1129, b499 |
| 4 | `TrimSpace` pada kode status gerbang tutup | **Melonggarkan gerbang uang**: `" 1 "` menutup klaim di Go, tidak di Pega. Dibuang; dibandingkan persis, dan ada ujinya | b608 `.STS_REJECT!=1` |
| 5 | Penjaga "nama orang tidak menyeberang" menyenaraikan **tiga kata terlarang** | Daftar larangan selalu kalah dari nama yang belum terpikirkan *(`namaPeserta` lolos)*. Diganti: **himpunan kunci dikunci utuh** | — |

⚠️ **Sebab nomor 1 layak dicatat**, sebab ia pengulangan butir av dalam bentuk lain: labelnya
dibaca, **ikatannya** tidak. `Total Claim Amount` terdengar seperti total sebuah klaim; yang
menentukan justru titik di depan `.TotalClaimAmount` dan kelas section-nya.

**Yang ditambahkan sesudah telaah:** `Find Disease` kini **dirender** sebagai `BelumTersedia`
bernama *(sebelumnya labelnya ada tetapi tidak pernah tampil — uji cacah label membuatnya tampak
selesai)*; dan gerbang tutup kini punya **pemanggil React** — sebelumnya ia persis "rute tanpa
pemanggil" yang paket 4 sendiri kecam.

---

## Giliran lanjutan 11 — paket 5: penanda DIPILIH (`bc012f8`)

### Cacat ketiga dengan bentuk yang sama persis

Jalur **pendaftaran** menulis `IsCheck: "1"`; gerbang **akseptasi** menuntut `"true"`. Setiap
peserta yang baru didaftarkan karena itu **ditolak** saat hendak diaksep — daftarkan klaim, tekan
`Save Adjustment`, tertolak. **Nol uji merah.**

Kedua sisi benar menurut dirinya sendiri: pendaftaran benar menurut komentarnya, akseptasi benar
menurut XML. Hanya **pertemuannya** yang salah — dan tidak ada satu pun uji yang memaksa keduanya
bertemu. Itu bentuk yang **sama persis** dengan dua cacat sebelumnya di giliran ini:

| # | Cacat | Bentuknya |
| ---: | --- | --- |
| 1 | envelope galat *(`galat` vs `error`)* | dua sisi, masing-masing benar sendiri |
| 2 | rute tanpa pemanggil *(tanggal kejadian)* | satu sisi ada, sisi lain tidak pernah dibuat |
| 3 | penanda dipilih *(`"1"` vs `"true"`)* | dua sisi, masing-masing benar sendiri |

XML menjawab tegas, dan yang keliru **komentar kita sendiri**: `SetIndexAdjustmentList.xml` b328
menuliskan `true`, dan `SavePesertaClaim.xml` b1812/b1997/b2413 mengujinya `.IsCheck=="true"`.
Nilai `"1"` tidak punya dasar XML mana pun.

Kini satu konstanta *(`models.PenandaDipilih`)* dan satu pembaca *(`services.PesertaDipilih`)*,
ditambah **penjaga statik**: nol berkas sumber menulis penanda itu dengan nilai harfiah. Penjaga
statiknya diperlukan sebab jalur penulis kedua tidak akan membuat uji perilaku mana pun merah
sampai ada yang mencoba mengaksep peserta yang lahir dari jalur itu.

### Ralat catatan giliran sebelumnya

Laporan paket 1 menulis `SetIndexAdjustmentList` *"belum dibangun, milik Akseptasi"*. **Keliru** —
ia **sudah ditiru lengkap**:

| Langkah | Baris XML | Di kode |
| --- | --- | --- |
| 1 `.IsCheck = true` | b281, b328 | `services/adjustment.go:112`, dengan kutipan barisnya |
| 2 `.IndexPremiumList` | b436 | **tidak perlu** — penunjuk balik posisi digantikan FK relasional |
| 3 warisan delapan kolom dari `.AdjustmentList(1)` | b570–b743 | `services.WarisiKolom`, delapan kolom persis |

Yang saya lakukan sebelumnya adalah membaca XML-nya lalu **tidak memeriksa apakah kode sudah
punya jawabannya**. Itu kesalahan arah sebaliknya dari biasanya, dan sama mahalnya: ia hampir
membuat pekerjaan yang sudah ada dikerjakan dua kali.

**Telemetri:** Go 297 → **300 PASS · 0 FAIL · 34 SKIP** · JS 201 · vet bersih · 0 migrasi baru.

---

## Giliran lanjutan 12 — paket 0: enam total peserta, dan OQ-H yang saya simpulkan salah

### Temuan: kesimpulan yang benar bukti-buktinya, salah kesimpulannya

Giliran lalu saya melaporkan **OQ-H**: `CheckTotalAdjustmentClaim` dirujuk sepuluh kali oleh
`ClaimLifeDetailGCNM.xml` dan **nol** berkas rule-nya ada di korpus, sehingga *"kelima total tidak
dapat ditiru — baris mana yang ikut belum terjawab, dan itu angka uang"*. Layar karena itu
menampilkan **enam medan uang** sebagai "belum tersedia".

Faktanya benar. **Kesimpulannya keliru**, dan dua kali keliru:

**1. Yang hilang hanya pemanggil *refresh*.** Yang **menghitung** keenam angka ada di korpus, dua
kali, dengan rumus yang sama persis:

| Activity | Dipanggil | Langkah |
| --- | --- | --- |
| `SavePesertaClaim.xml` | `Submit` b27369 | 8 b4002 *(ULANG peserta, `EMBEDDED` b4843)* · 8.1 b4221 *(ULANG `.AdjustmentList`, `EMBEDDED` b4570)* · 8.2 b4592 *(tulis ke peserta b4616–b4743)* |
| `SaveOutStandingLife_Act.xml` | `Save to RNM` b21102 | 23 b10638 · 23.1 b10841 *(`EMBEDDED` b11046)* · 23.2 b11067 *(b11091–b11217)* |

Pertanyaan yang saya sebut "tidak terjawab" **terjawab di XML**: kedua langkah penjumlahnya
**berprasyarat KOSONG** *(`<pyStepsPreCondParamsWhen/>` b4562)*. Jadi **seluruh** baris
`.AdjustmentList` ikut — termasuk yang `STS_REJECT = 2`, termasuk yang `IsCheck`-nya dicabut.
Dugaan wajar *("tentu yang ditolak tidak ikut")* justru yang salah, dan itulah kenapa menebak
berbahaya di kedua arah — bukan hanya arah "mengarang angka".

**2. Totalnya ENAM, bukan lima.** `Total Ceding Retention` b20629 / `.TotalCedingRetention` b20636
tidak pernah masuk daftar saya, sebab saya mencacah lewat **rujukan** `CheckTotalAdjustmentClaim` —
dan ia satu-satunya total yang **tidak punya aksi refresh**. Mencacah lewat pemanggil, bukan lewat
label. Enam hari medan uang itu tidak ada di layar dan **tidak satu pun uji berbunyi**; uji saya
justru mengunci cacah **lima**, sehingga kekurangannya tampak disengaja.

⚠️ **Sebab kedua kekeliruan itu satu**: saya berhenti pada rule yang **namanya tertulis di section**,
lalu menyimpulkan tentang **medannya** — tanpa menanyakan *siapa lagi yang menulis properti itu*.
Bentuk yang sama persis dengan *"Close Claim tidak punya aksi lain"*: berhenti di aksi pertama yang
membawa activity, lalu menyimpulkan tentang seluruh tombol. Dua kali dalam dua giliran.

### Yang dibangun

Keenam total dihitung `models.HitungTotalPeserta`, dirakit `services.KlaimLife.Ambil` **saat Detail
dibaca**, dan **tidak disimpan**. Nol kolom `TOTAL_*`; dua penjaga statik menolak migrasi yang
menambahkannya dan menolak `repository` yang menulisnya.

⚠️ **Penyimpangan sadar, dinyatakan**: Pega **menyimpan** hasilnya saat `Submit`/`Save to RNM`,
sehingga layar Pega dapat **basi** sesudah putaran atau akseptasi sampai tombol simpan ditekan lagi.
Di sini tidak pernah basi. Ditanyakan balik ke pemilik ekspor kalau-kalau ada laporan yang justru
mengandalkan angka tersimpan.

### Uji yang dibuat gagal lebih dulu, lalu dipulihkan

| Penjaga | Dibuat gagal dengan | Berbunyi |
| --- | --- | --- |
| `baris DITOLAK ikut dijumlah` | menambahkan saringan `KodeStatus == KodeDitolak` | `JumlahKlaim = "6 IDR", mau "12 IDR"` |
| `TiapKolomTotalMembacaKolomnyaSendiri` | menukar `SUM_INSURED` ↔ `SUM_REASURED` | tiga uji sekaligus |
| `EnamKolomTotalDikunci` | menghapus `CEDING_RETENTION` | `kolomTotal = 5, mau 6` |
| `MigrasiTidakMenyimpanTotalPeserta` | migrasi palsu ber-`TOTAL_CLAIM_AMOUNT` | menyebut berkas dan kolomnya |
| `RepositoryTidakMenulisTotalPeserta` | `p.Total = models.TotalPeserta{}` di repository | menyebut berkas dan barisnya |
| label `Total Ceding Retention` b20629 | digeser ke b20630 | membaca `<pyColumnSorting>` |
| rumus penjumlah di XML | mencari `STS_REJECT + local.Total` | gagal, seperti seharusnya |

Dan **satu penjaga lama berbunyi tanpa diminta**: kunci himpunan kunci JSON peserta menolak medan
`total` yang baru. Itu tugasnya — medan baru pada marshaller harus **dilihat** orang, bukan
menyelinap. Saya tuliskan namanya di sana, bukan melonggarkan penjaganya.

### Kontrak dua sisi, sebab ini cacat keempat yang berbentuk sama

Nama keenam medan JSON dikunci **di kedua sisi**: `TestNamaJSONTotalPesertaDikunci` di Go dan blok
`kontrak JSON TotalPeserta` di `PanelTotalPeserta.test.ts`, masing-masing memuat daftar yang sama
dan menyebut pasangannya. Tanpa itu, nama medan uang yang berganti di satu sisi membuat React
membaca `undefined` dan menampilkan sel **kosong** — tanpa satu pun galat, persis seperti envelope
`galat`, rute tanpa pemanggil, dan `IsCheck` `"1"`.

**Telemetri:** Go 300 → **307 PASS · 0 FAIL · 34 SKIP** · JS 201 → **207** · `tsc` bersih ·
`go vet` bersih · build **47** modul · **nol** migrasi baru.

---

## Giliran lanjutan 12 — kelompok 1: `Close Claim` menutup, dan daftar dokumen tampil

### Bagian A — penutupan kasus (butir bb)

Alur dibaca **utuh**, dan yang ditemukan membalik catatan lama: `Flow/Register_Flow.xml` punya 12
konektor dan **tidak satu pun bernama `CloseClaim`**. Ia **local action**, dan cacah berkasnya yang
menentukan segalanya — `pyLocalAction>CloseClaim` ada di **tepat dua** section:

| Section | Baris | Tahap |
| --- | --- | --- |
| `InputOSClaimLife.xml` | b22837, b22988 | Outstanding Claim |
| `InputAkseptasiClaimLife.xml` | b21457, b21602 | Claim Analis |

Hanya shape **End1** yang menetapkan status kerja: `<rowdata REPEATINGINDEX="End1">` b883,
`pyMOId End1` b885, `Data-MO-Event-End` b901, `pyWorkStatus` **`Resolved-Completed`** b899.
Sembilan shape lain kosong.

⚠️ **Satu jebakan pembacaan yang hampir menjerat saya.** `grep` atas `pyShapeType` dan
`pyWorkStatus` berdampingan menampilkan `Resolved-Completed` b899 tepat di bawah
`Data-MO-Event-Start` b867 — seolah status itu milik shape **Start**. Ia bukan: b883 membuka
`rowdata` baru. Pohonnya harus dibaca lewat batas `rowdata`, bukan lewat kedekatan baris. Aturan
yang sama yang pernah menyembunyikan dua loop di `DeletePesertaClaimLife`.

**Yang dibangun:** migrasi `017` *(`STATUS_WORK VARCHAR2(32)`, NULL = belum ditutup, satu-satunya
nilai VERBATIM b899)* · `POST /api/klaim-life/{id}/tutup` → 409 berisi **seluruh** penghalang ·
satu transaksi mengisi status, **mengosongkan `TAHAP`**, dan merekam jejak · tombol hanya pada kedua
tahap itu, konfirmasi **VERBATIM b499**, berhasil → jendela ditutup *(`closeContainer` b1129)*.

### Satu pintu untuk tujuh rute, sebab tujuh tempat adalah tujuh tempat untuk lupa

Sesudah tutup, **setiap** rute pengubah harus menolak. Ditulis ulang di tujuh berkas, aturan itu
suatu hari hanya akan ada di enam — dan yang ketujuh **tidak akan berbunyi**, sebab tiap berkas
hijau sendirian. Bentuk cacat yang sudah terjadi tiga kali di modul ini.

Jadi: satu fungsi `PastikanKasusTerbuka`, dan **dua** penjaga statik yang saling menutup arah:

| Penjaga | Arah yang dijaga |
| --- | --- |
| `TestSetiapLayananPengubahMemeriksaKasusTerbuka` | tiap layanan di daftar **memanggilnya** |
| `TestDaftarLayananPengubahMencakupSeluruhRutePengubah` | tiap layanan **bertransaksi** ada di daftar, atau punya alasan tertulis bernama |

Yang kedua **langsung menemukan dua berkas** yang saya lewatkan — `services.go` *(mendefinisikan
`DalamTransaksi`)* dan `statusbaris.go`. Yang kedua nyata: `Status.Ubah` dan `Status.Tolak`
sama-sama menyalurkan ke badan `ubah`, jadi penjaganya saya **pindahkan ke sana** dan cabut dari
`tolak.go`. Penjaga di dua pintu menuju satu ruang adalah dua tempat untuk lupa.

⚠️ **Dan penjaga lama menangkap kesalahan saya.** `TestUbahStatusMenjagaPagarnya` merah: saya
menaruh pemeriksaan itu **sebelum** pemeriksaan bentuk permintaan, sehingga permintaan tanpa
pengenal peserta dijawab *"ORACLE_DSN belum dikonfigurasi"* alih-alih *"pengenal wajib diisi"* —
penjaga yang benar, diletakkan di tempat yang membuat galat lain berbohong. Dipindahkan ke sesudah
pemeriksaan bentuk.

### Bagian B — daftar dokumen (`DocumentLife`)

`LoadDocumentLife_ACT.xml` dibaca sebagai pohon. **Nomor barisnya berbeda jauh dari brief** — brief
menyebut b755/b861/b1164/b1377/b1516; bacaan ini b388/b495/b797/b1011/b1149. Yang dipakai bacaan
ini, selisihnya dicatat.

Dua temuan yang mengubah apa yang dibangun, **keduanya prasyarat ber-`WhenTrue=3` (LEWATI)**:

1. **b1404 `.DOCUMENT==""`** menggerbangi seluruh langkah 1.4. Saringan browse-nya
   `Field .KATEGORI_1` `=` `Value .DOCUMENT` — dan kolom `DOCUMENT` **tidak ada** di
   `T_CLAIMLF_PREMIUMLIST_DETAIL`. Saringan itu **tidak dapat ditiru**; yang dipakai FK
   `PREMIUM_LIST_DETAIL_ID`. Dilaporkan **OQ-J**.
2. **b1310 `DataImage.URLImage==""`** menggerbangi langkah 1.4.2 — baris yang URL penyimpanannya
   kosong **tidak ikut ditambahkan** ke daftar. Karena Google Storage belum tersambung, URL-nya
   selalu kosong; meniru gerbang itu membuat daftar **selalu kosong**, dan layar akan berkata
   *"tidak ada dokumen"* untuk peserta yang dokumennya lengkap. Barisnya **tetap tampil** dengan
   penanda. Penyimpangan sadar, dinyatakan di kode, di PARITAS, dan di OQ-J.

Keempat tombol section itu *(`Refresh` b611, `Add attachment` b1245, `View Office Online` b3502,
`Delete` b4288 — yang terakhir menjalankan **dua** aksi seperti `Close Claim`)* hadir sebagai
`BelumTersedia` bernama; unggah/unduh/hapus milik kelompok Dokumen.

### Penjaga yang dibuat gagal lebih dulu, lalu dipulihkan

| Penjaga | Dibuat gagal dengan | Berbunyi |
| --- | --- | --- |
| `SetiapLayananPengubahMemeriksaKasusTerbuka` | cabut panggilan dari `dol.go` | `dol.go (Set) tidak memanggil…` |
| `bolehTutupDiLayar` | tambahkan `Medical Check` ke daftar tahap | dua uji sekaligus |
| status kerja dibandingkan persis | sisipkan `.trim()` | `expected true to be false` |
| `barisDokumen` memakai KATEGORI_2 | tukar ke `kategori1` | `expected 'SALAH' to be 'BENAR'` |
| himpunan kunci baris dokumen | selundupkan medan `namaTertanggung` | himpunan kunci berubah |
| `NamaJSONDokumenDikunci` | cabut tag `json:"namaFile"` | `NamaFile` muncul berhuruf besar |

Dan **dua penjaga lama berbunyi tanpa diminta, keduanya benar**: `strukturkolom_test` menolak DDL
`STATUS_WORK` yang belum ada di STRUKTUR-TABEL *(didokumentasikan, bukan dilonggarkan)*, dan
`TestNamaTabelDokumenLamaHanyaUntukWarisan` menolak nama warisan `DOCUMENT_CLAIM` di komentar saya
yang tidak menyebutnya warisan *(komentarnya yang diperbaiki)*.

**Telemetri:** Go 307 → **316 PASS · 0 FAIL · 34 SKIP** · JS 207 → **224** · `tsc` bersih ·
`go vet` bersih · build 47 → **48** modul · migrasi **017** *(dari keputusan bb yang tercatat)*.

---

## Giliran lanjutan 12 — kelompok 2: aturan dokumen yang dapat diputuskan, dan yang tidak

### Apa yang dibangun, dan kenapa hanya itu

Sembilan activity kelompok ini dibaca sebagai pohon. Empat di antaranya — `InsertGoogleStorage_Act`,
`GetUrlGoogleStorage_Act`, `DeleteGoogleStorage_Act`, `SendEmailWithAttachments` — **menyambung ke
layanan luar nyata**, dan menyambungkannya menuntut persetujuan manusia. Jadi yang dibangun giliran
ini adalah seluruh aturan yang **dapat diputuskan tanpa layanan itu**, dan sisanya **dinyatakan**:

| Aturan | Sumber | Keadaan |
| --- | --- | --- |
| Tabel MIME **48 baris** + `otherwise` | `DecisionTable/GetMimeType.xml` b290–b337 / b417–b464, bawaan b89 | ✅ disalin utuh, cacahnya dikunci |
| MIME pemanggil **menang** atas tabel | prasyarat b586 `Param.MIME==""` `WhenTrue=2` LANJUT | ✅ |
| Kunci kelompok `DL-`+angka, **tak pernah ditimpa** | `SaveAttachLife.xml` b595–596, b615–616 | ✅ |
| Pengenal baris `yyyyMMddhhmmssSSS` Asia/Jakarta | `InsertDocument_Act.xml` b647–648 | ✅ *(termasuk cacatnya — di bawah)* |
| **Unggah dulu, baris kemudian** | prasyarat b1283 `T_STORAGE_ID==""` `WhenTrue=3` LEWATI | ✅ aturannya |
| Hapus penyimpanan dilewati bila kosong; baris **tetap** dihapus | b472 `WhenTrue=3`; `Obj-Delete` b513 tanpa prasyarat | ✅ aturannya |
| Unggah / unduh / hapus sebagai **rute** | `InsertGoogleStorage_Act` b1023 dst. | ⛔ menunggu persetujuan penyambungan |

### Cacat warisan yang ditiru, bukan diperbaiki

`.ID` baris dokumen adalah `@CurrentDate("yyyyMMddhhmmssSSS","Asia/Jakarta")` b648. Pada pola
Java/Pega, **`hh` adalah jam 12-jam**. Jadi pukul 14:05:09.123 dan pukul 02:05:09.123 menghasilkan
pengenal yang **sama persis**.

Itu cacat, dan saya **tidak memperbaikinya**. Memperbaikinya membuat pengenal baris baru berbeda
bentuk dari pengenal baris lama, dan keduanya hidup di kolom yang sama — pengenal yang tidak dapat
diurutkan bersama pendahulunya adalah harga yang lebih mahal daripada tabrakan yang hanya mungkin
dalam milidetik yang sama. Ia **dikunci uji** justru supaya tidak diam-diam "dirapikan" orang, dan
uji itu berbunyi ketika saya menggantinya ke jam 24 sebagai percobaan.

### ⛔ Ralat kedua hari ini atas OQ saya sendiri

OQ-J yang saya tulis beberapa jam lalu bertanya *"kolom apa `.DOCUMENT` itu?"*. **Terjawab, oleh
XML yang belum saya baca saat menulisnya**: `SaveAttachLife` b596 **membuat**nya —
`"DL-" + angka(@CurrentDateTime())` bila kosong — dan meneruskannya sebagai `KATEGORI_1` b1467.
Ia kunci kelompok dokumen milik peserta, bukan misteri.

⚠️ Sebabnya **sama persis dengan ralat OQ-H pagi ini**: saya membaca activity **pembaca** lalu
menyimpulkan tentang sebuah medan, tanpa menanyakan siapa yang **menulis** medan itu. Dua kali dalam
satu hari, dengan bentuk yang identik. Saya catat bentuknya, bukan hanya kejadiannya: *sebelum
menyatakan sebuah medan tidak terbaca, cari penulisnya, bukan hanya pembacanya.*

OQ-J tetap terbuka tetapi jauh lebih sempit, dan pertanyaannya kini milik **migrasi data**: apakah
`KATEGORI_1` pada dokumen warisan selalu cocok dengan satu peserta? Yang tidak cocok tidak akan
terbawa FK, dan itu harus diketahui **sebelum** migrasi.

### Penjaga yang dibuat gagal lebih dulu, lalu dipulihkan

| Penjaga | Dibuat gagal dengan | Berbunyi |
| --- | --- | --- |
| cacah tabel MIME | hapus baris `"et"` | `petaMime = 47 baris, mau 48` |
| MIME pemanggil menang | balik urutannya | `= "application/pdf", mau image/png` |
| kunci kelompok tak ditimpa | cabut cabang `sudahAda` | kunci baru menimpa yang lama |
| cacat jam 12-jam ditiru | ganti pola ke jam 24 | dua pengenal tidak lagi bertabrakan |

**Telemetri:** Go 316 → **322 PASS · 0 FAIL · 34 SKIP** · JS **224** *(tak berubah — kelompok ini
belum menyentuh layar)* · `tsc` bersih · `go vet` bersih · build 48 modul · **nol** migrasi baru.

---

## Giliran lanjutan 12 — kelompok 3: pencarian diagnosa berbatas atas 97.586 baris

### Batasnya bukan karangan kami

`DISEASE_LIFE` berisi **97.586 baris**, dan setiap angka pembatasnya ada di rule:

| Angka | Sumber |
| ---: | --- |
| `pyMaxRecords` **500** | `ReportDefinition/BrowseDiseaseLife_RD.xml` b659 |
| `pyPageSize` **50** | b514 / b747 |
| `pyQueryTimeoutValue` 30 | b661 |
| urut `.Number` ASC | b598 `pySortOrder` 1 |

Penyaringnya **`A AND B`** b535/b754 — A `.ICD_Code` `Contains` `Param.ICD_Code`,
B `.Disease` `Contains` `Param.Disease` — dan keduanya dinaikkan ke huruf besar lebih dulu
*(`SearchDiagnose_act.xml` b255, b302)*.

⛔ **`CARI1` adalah KODE dan `CARI2` adalah NAMA**, dan itu dibaca dari pemetaannya
*(`Diagnose_Section.xml` b1645, b1651)* — bukan ditebak dari namanya. "CARI1/CARI2" tidak
menyebutkan apa pun, dan menukarnya membuat setiap pencarian gagal dengan cara yang terlihat
seperti *"datanya memang tidak ada"*.

⚠️ **Pencarian dengan kedua kata kunci kosong SAH**, dan itu ditiru: di Pega `Contains ""` cocok
dengan semua baris, dan yang menahannya hanya `pyMaxRecords`. Justru kombinasi itulah yang paling
harus berbatas, dan ujinya menguji **seluruh** kombinasi termasuk yang kosong.

⚠️ Dan ketika hasil menyentuh batas, layar **mengatakannya**. Daftar terpotong yang diam terbaca
sebagai daftar lengkap; pemakai akan menyimpulkan diagnosanya tidak ada, lalu berhenti mencari.

### Dua hal yang SENGAJA tidak dibangun, dan sebabnya

**1. Tombol `Choose` b2509 → `SetDisease` b2528.** `.DiagnoseList` *(`ClaimLifeDetailGCNM.xml`
b3923)* adalah **`RepeatGrid`** — banyak diagnosa per peserta. Tetapi
`T_CLAIMLF_PREMIUMLIST_DETAIL` hanya punya `DISEASE` dan `ICD_CODE` **tunggal** *(migrasi 003)*.
**Satu lawan banyak.** Memilih diagnosa ke kolom tunggal berarti membuang diagnosa kedua dan
seterusnya — diam-diam. **OQ-K.2**, dan ia keputusan skema, bukan keputusan executor.

**2. Nama kolom `DISEASE_LIFE`.** Ekspor memuat nama **properti Pega** tetapi **tidak** memuat
`Rule-Obj-Class`-nya, jadi pemetaan kelas-ke-tabel tidak ada. `DISEASE` dan `ICD_CODE` dugaan kuat
*(sama persis dengan tabel saudaranya)*; `.Number` → `NUMBER_` **terbuka** — `NUMBER` kata cadangan
Oracle, jadi kolomnya pasti bernama lain. Ketiganya dikumpulkan di **satu blok konstanta** supaya
koreksi DBA adalah satu suntingan, dan ada uji yang menagih OQ-K.1 supaya pertanyaannya tidak hilang
bersama giliran ini.

### Penjaga yang dibuat gagal lebih dulu, lalu dipulihkan

| Penjaga | Dibuat gagal dengan | Berbunyi |
| --- | --- | --- |
| `QueryPenyakitSelaluBerbatas` | cabut `FETCH FIRST` | kriteria kosong tanpa batas |
| `KriteriaPenyakitDisambungAND` | ganti `AND` → `OR` | `query memakai OR` |
| nilai lewat bind, bukan tempel | tempelkan `k.Nama` ke teks | `nilai kriteria tertempel` |
| jepitan batas di klien | cabut `Math.min` | `expected '1000000' to be '500'` |

**Telemetri:** Go 322 → **328 PASS · 0 FAIL · 34 SKIP** · JS 224 → **234** · `tsc` bersih ·
`go vet` bersih · build 48 → **49** modul · **nol** migrasi baru.

---

## Giliran lanjutan 12 — kelompok 4: tiap layar menawarkan perpindahannya sendiri

### Temuan: backendnya sudah lengkap, layarnya yang tidak

Keempat perpindahan yang kelompok Medis dan Akseptasi perlukan **sudah sah** di
`models.serahTerimaSah` sejak F0.4, dan rutenya `POST …/tahap/{tujuan}` sudah ada. Yang tidak ada
adalah **tombolnya**: hanya layar Outstanding yang punya, dan Medical Check serta Claim Analis tidak
punya satu pun.

Bentuk cacat yang sama dengan `Edit Date` giliran lalu — jalur backend lengkap, teruji, dan **tidak
dapat dijalankan siapa pun**. Ia tidak berbunyi di uji mana pun: backend hijau, layar hijau, fiturnya
tidak ada. Kali ini tiga jalur sekaligus.

### Cacah berkas, bukan tata letak

| Layar | Tombol | Baris | Tujuan |
| --- | --- | ---: | --- |
| `InputOSClaimLife` | `Send Back to Register` | b21404 | Input Register |
| | `Send to Medical Check` | b21839 | Medical Check |
| `MedicalCheckClaimLife` | `Send Back to Admin` | b20256 | Outstanding |
| | `Send to Claim Analyst` | b21151 | Claim Analis |
| `InputAkseptasiClaimLife` | `Send Back to Admin` | b20221 | Outstanding |
| | `Send Back to Medical` | b20467 | Medical Check |

⛔ **Menyalin tombol antarlayar membuka jalur yang di sistem lama tidak ada** — dan `SerahTerimaSah`
akan menolaknya 409, sehingga yang lahir hanya tombol yang selalu gagal. Ujinya mengunci cacahnya
dari kedua arah, plus satu yang menjaga hal yang lebih halus: **tidak satu pun tombol menunjuk
tahapnya sendiri**.

⚠️ Dan dua kalimat yang menuju tahap yang **sama** sengaja tidak disatukan: `Send to Medical Check`
*(Outstanding)* dan `Send Back to Medical` *(Claim Analis)* keduanya menuju Medical Check.
Menukarnya tidak akan membuat satu pun uji perilaku merah — jadi ada uji yang mengunci teksnya.

⚠️ Sesudah perpindahan, klaimnya **dibaca ulang**. Tanpa itu layar tetap menawarkan tombol tahap
LAMA, dan tiap kliknya ditolak 409 — pemakai akan menyimpulkan perpindahannya gagal padahal
berhasil.

**Telemetri:** Go **328 PASS · 0 FAIL · 34 SKIP** *(tak berubah — kelompok ini murni layar)* ·
JS 234 → **240** · `tsc` bersih · build 49 → **50** modul · **nol** migrasi baru.

---

## Giliran lanjutan 12 — kelompok 5: label yang dikarang, dan yang dikunci ke XML

### ⛔ Temuan: satu tombol memakai teks yang tidak ada di korpus mana pun

Tombol penyerahan ke Komite di layar Detail berbunyi **`Send ke Komite`** — campuran
Indonesia-Inggris yang **tidak ada di satu pun berkas korpus**. Labelnya ada di XML dan tidak
pernah diambil: `Section/ClaimComite.xml` **b7033** `<pyLabel>Send Claim to Committee</pyLabel>`.

Itu bukan cacat kecil. Pemakai sistem lama mencari kalimat yang sama, dan penguji penerimaan
membandingkan layar lama dengan layar baru **kata demi kata** — tombol yang namanya berbeda terbaca
sebagai fitur yang berbeda. Dan ia **tidak berbunyi di uji mana pun**: tombolnya bekerja, rutenya
benar, hanya namanya yang karangan.

Kini `TOMBOL_KOMITE.serahkan` mengambil dari XML, dan ada uji yang **membaca baris b7033 langsung
dari korpus** — bukan menyalinnya ke dalam uji. Dibuktikan: mengembalikan teks lama membuat uji itu
merah seketika.

⚠️ Dicatat pula: section-nya bernama `ClaimComite` sedangkan `pyRuleName`-nya `ClaimComiteeLife`
(b139). Dua ejaan, keduanya salah eja, keduanya warisan — disebut apa adanya supaya dapat dicari.

### Yang TIDAK dikerjakan kelompok ini, dan sebabnya

- Penyerahan ke Komite **sudah ada** sejak tiket 10 *(`services.Penyerahan.Serahkan`,
  `POST …/adjustment/{adjId}/komite`)*; yang kurang hanya namanya.
- `SendEmailKlaimLF` → outbox: **menyambung ke layanan email nyata**, dan itu menuntut persetujuan
  manusia. `EMAILKOMITE` pun berisi nama dan alamat email orang — dibaca saat jalan, tidak pernah
  disalin ke fixture, tiket, atau log.
- Berkas `komite_*` **tidak disentuh**, sesuai brief.

**Telemetri:** Go **328 PASS · 0 FAIL · 34 SKIP** · JS 240 → **241** · `tsc` bersih ·
build **50** modul · **nol** migrasi baru.

---

## Giliran lanjutan 12 — paket telaah: `/code-review` dua sumbu atas enam commit

Titik tetap `8907fba`. Dua penelaah berjalan terpisah — satu atas **standar**, satu atas **spec** —
supaya tidak saling mencemari. Temuannya saya periksa satu per satu; **dua saya bantah dengan
bukti**, sisanya saya perbaiki.

### ⛔ Cacat lintas-lapis KELIMA, dan saya yang membuatnya

`models.Dokumen.Tanggal` bertipe `*time.Time`, sehingga kolom `DATE` yang NULL menyeberang sebagai
**`null`**. `api.ts` menyatakannya **`tanggal: string`**. `PanelDokumenPeserta` memeriksa
`b.tanggal === ''` — yang **tidak pernah menyala untuk `null`** — sehingga tanggal yang memang
kosong tampil sebagai **sel kosong**, bukan penanda `—`. Persis kebalikan dari yang ADR-U-0027
minta, dan `tsc` tidak dapat menangkapnya: tipe yang berbohong tentang data dari jaringan tetap
dikompilasi.

⚠️ **Ia lahir di giliran yang sama ketika saya mendaftarkan keempat pendahulunya** — envelope
`galat`, rute tanpa pemanggil, penanda `IsCheck`, nama medan total. Menamai pola tidak membuat saya
kebal terhadapnya. Yang berubah sekarang: `string | null` di klien, normalisasi di **satu** tempat,
dan uji yang gagal bila normalisasi itu dicabut.

### Lubang penjaga yang penjaganya sendiri tidak lihat

`TestDaftarLayananPengubahMencakupSeluruhRutePengubah` menemukan calonnya lewat
`strings.Contains(…, "DalamTransaksi(ctx")`. **`hapus.go` mengubah** (`DELETE /api/klaim-life/{id}`)
tetapi **tidak memanggilnya sama sekali** — ia lolos penemuan, dan hanya aman karena kebetulan saya
tulis tangan di daftar. Layanan pengubah baru yang berbentuk seperti `hapus.go` akan lolos dari
**kedua** penjaga.

Penjaga ketiga kini memakai definisi yang tidak dapat diakali dengan menulis kode berbeda bentuk:
**tabel rute**. Apa pun yang terdaftar dengan metode selain `GET` adalah pengubah, titik — dan ia
harus terdaftar beserta berkas layanannya atau punya alasan pengecualian tertulis. Dibuktikan:
menambahkan satu rute `POST` palsu membuatnya merah seketika.

### Komentar saya sendiri yang berbohong

`totalpeserta.go` menulis `nilai.Currency = kurs` dengan komentar *"yang berbeda sudah ditolak di
atas"*. **Komentar itu keliru**: pemeriksaan di atas hanya membaca `CurrencyID` **baris**, tidak
pernah membaca `Currency` tiap nilai uang. Satu kolom bermata uang lain karena itu dijumlahkan
diam-diam di bawah mata uang yang salah — dan `Money.Add` yang seharusnya menolaknya justru
**dilucuti lebih dulu**. Kini tiap kolom diperiksa sendiri; yang tanpa label mengikut penampung,
yang berbeda ditolak.

### Dua temuan yang saya BANTAH, dengan bukti

| Temuan | Bantahan |
| --- | --- |
| *"label `totalCedingRetention` ditambahkan tanpa kutipan"* | Kutipannya **ada**, `labels.ts` b333–338: `b20629 pyLabelPreview -> .TotalCedingRetention b20636` |
| *"`GetListKomiteLife` dijatuhkan diam-diam"* | Ia **sudah ditiru** sejak tiket 10 — `services/komite.go` mengutipnya `[terverifikasi]` di empat tempat *(pecahan baris 449-450, 655, 790)* |

### Temuan yang saya terima sebagian, dengan alasan dipertajam

**Penghapusan dokumen.** Penelaah benar bahwa `Obj-Delete` b513 **tanpa prasyarat**, jadi baris dapat
dihapus tanpa memanggil layanan luar. Tetapi alasan penundaan tetap berdiri, dan kini saya nyatakan
lebih tepat: **tidak ada baris yang dapat dihapus tanpa `T_STORAGE_ID`**, sebab penyimpanannya
digerbangi b1283 — setiap baris yang ada pasti punya berkas. Menghapus barisnya saja akan
meninggalkan berkas **yatim di penyimpanan**, selamanya, tanpa penunjuk.

**Aturan dokumen tanpa pemanggil.** Benar, dan ini bentuk yang sama dengan "rute tanpa pemanggil".
Bedanya satu: di sini ketiadaan pemanggil **dinyatakan** — kini di kepala kedua berkasnya, bukan
hanya di laporan ini.

**Dialog diagnosa per peserta.** Benar bahwa ia berdiri di dalam perulangan peserta — dan itu
memang letaknya menurut XML *(b5061 di section berkelas `Int-LIFE_PREMIUM_DETAIL`)*. Tetapi klaim
grup berpeserta 500 akan merender 500 formulir sekaligus. Kini ia **tombol yang membuka**, persis
seperti `showHarness` b5081: lebih setia, bukan kurang.

**Telemetri:** Go 328 → **329 PASS · 0 FAIL · 34 SKIP** · JS 241 → **242** · `tsc`, `vet`, `gofmt`
bersih · build 50 modul · nol migrasi baru.

---

## Giliran tiga modul 1 — §5 penyatuan, §1 OQ-K.1, §2 A bagian 1

### §5 — penyatuan tiga cabang

| Langkah | SHA | Hasil |
| --- | --- | --- |
| `merge --no-ff modul/komite-claim-life` | `df1353d` | tanpa konflik |
| `merge --no-ff modul/premiumlist-life` | `c1d3b8c` | **dua konflik**, keduanya diselesaikan dengan mengambil **kedua sisi** |
| `--ff-only` kedua worktree | `c1d3b8c` | ketiganya sejajar |

Konfliknya tepat di tempat yang brief duga: `migrasi_test.go` dan `strukturkolom_test.go`. Dan
keduanya **benar**, dengan cara yang saling melengkapi:

- sisi PremiumList **menyempitkan LINGKUP** penjaga kaskade *(001–049 milik Claim Life; kebijakan
  kedua modul berbeda dan keduanya disengaja)*;
- sisi Komite **memperbaiki ISI**-nya *(relasi 9 masuk daftar; penjaga itu sempat menahan
  perbaikan cacat `013`)*.

⚠️ Kedua pelajaran itu **satu**, dan saya satukan di komentarnya: *daftar yang ditulis sebelum modul
kedua lahir berhenti menjadi penjaga dan mulai menjadi pagar.* Yang satu menyempitkan lingkupnya,
yang lain memperbaiki isinya — keduanya perlu.

`letakStruktur` kini **tiga** dokumen. Dua di antaranya menggambarkan tabel yang sama
*(`T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST` — batas antara dua konteks)*, dijaga
`TestDokumenSTRUKTURSepakatAtasTabelBersama`.

**Telemetri:** Go 329 → **334 PASS · 0 FAIL · 34 SKIP** · JS **242** · 50 modul.

### §1 butir 1 — ⛔ tebakan saya, dan namanya ada di korpus sepanjang waktu

`kolomNomorPenyakit = "NUMBER_"` saya tulis beralasan *"`NUMBER` kata cadangan Oracle, jadi kolomnya
pasti bernama lain — dan nama itu tidak ada di korpus"*.

Separuh pertamanya benar. **Separuh keduanya salah**:

```
ReportDefinition/BrowseDiseaseLife_RD.xml
  b598  <pyFieldName>.Number</pyFieldName>
  b599  <pyFieldLabel>ID</pyFieldLabel>      <- namanya, di baris berikutnya
```

Baris b599 **ada di keluaran grep saya sendiri** saat saya membaca report definition itu untuk
mengambil `pyMaxRecords` dan penyaringnya. Saya memperlakukannya sebagai **label layar**.

Katalog DEV membenarkan: `POOLDATA.DISEASE_LIFE` berkolom `ID`, `ICD_CODE`, `DISEASE`.

⚠️ Pelajarannya sempit dan tajam, dan berbeda dari yang sudah saya catat: *sebelum menyatakan
sesuatu "tidak ada di korpus", periksa apa yang **sudah terbaca** — bukan hanya apa yang sudah
dicari.* Uji penggantinya menyebut tebakan lama **dengan namanya**, supaya ia tidak kembali lewat
"perapian" berikutnya.

### §2 A bagian 1 — ⛔ ralat kedua: grid yang saya baca separuh

OQ-K.2 saya tulis: *"`.DiagnoseList` RepeatGrid (banyak) lawan kolom tunggal — satu lawan banyak,
keputusan skema"*. Jawabannya ada **dua baris di bawah** tempat saya berhenti:

```
b4690  <pyLabel>Add</pyLabel>     -> addRow    b4700/b4841
b6160  <pyLabel>Delete</pyLabel>  -> deleteRow b6170
```

Grid memang **dapat** berarti tampilan satu baris. Grid ber-`Add` **dan** ber-`Delete` **tidak
dapat**. Saya berhenti pada bentuk tampilannya tanpa membaca tombolnya — bentuk kesalahan yang sama
dengan `Close Claim` *("berhenti di aksi pertama")* dan OQ-H *("berhenti pada rule yang namanya
tertulis")*. **Ketiga kalinya.**

Migrasi `018` lahir: `T_CLAIMLF_DIAGNOSE`, FK peserta `ON DELETE CASCADE`, `URUTAN` =
`.pxListSubscript`, sequence, index.

⛔ **Lebar `DISEASE` 1000 terukur, bukan selera.** Sumbernya `DISEASE_LIFE.DISEASE` `VARCHAR2(1000)`
dengan isi terpanjang **290** `[data DBA]`. `VARCHAR2(255)` di kolom peserta **terbukti kurang** —
ia akan menolak nama penyakit yang sah. `018` melebarkannya pula.

⛔ **Kaskade dengan bukti**: `SetDisease.xml` b389 menutup dengan `Obj-Save pyWorkPage`, bukan
menyimpan halaman diagnosa sendiri — daftarnya hidup **di dalam** halaman peserta. Daftar penjaga
**diperbarui**, bukan dilawan *(pelajaran tiket 00 Komite, dua jam sebelumnya)*.

**OQ-L dibuka, lebih sempit:** daftar pilihan `GROUPDIAGNOSE`. Dropdown b5863 ber-`pyListSource`
**`associated`** — daftarnya hidup pada **rule properti** di kelas `Data-DiagnoseLife`, dan rule itu
tidak ada di ekspor *(`GROUPDIAGNOSE` muncul di tepat satu berkas)*. Kolomnya dibuat; nilainya tidak
dikarang.

**Gerbang sunting dicatat:** `pyDisabledWhen` `.STS_REJECT=='1' || .STS_REJECT=='2'` muncul **empat
kali** di dalam grid — baris yang sudah diaksep atau ditolak tidak dapat disunting lagi.

**Telemetri:** Go **334 PASS · 0 FAIL · 34 SKIP** · JS **242** · migrasi `001`–`018`, `030`,
`050`–`056`.

## Giliran tiga modul 2 — paket 0 dan paket 1 (diagnosa bd)

### Paket 0 — `06bf711` *(dokumen; nol kode produksi)*

Tiket 08 mendapat bab bertanggal *"Diagnosa banyak per peserta"*; baris 94-nya — yang masih
berbunyi *"bila work owner menyetujui **al**"* — diralat. Tiket 14 mendapat bab *"Migrasi 018"*.
PARITAS 24/25/13c menyebut rute yang akan ada, dan 13c **meralat dirinya sendiri**.

`TestKaskadeHanyaPadaEmpatRelasi` → `TestKaskadeHanyaPadaRelasiTerdaftar`. Daftarnya enam berkas
sejak 018; nama yang menyebut **angka** berbohong setiap kali relasi berikutnya lahir, nama yang
menyebut **aturannya** tidak. Nama lamanya ditulis di komentar supaya pencarian atasnya tetap
sampai. Dibuktikan merah: cabut `ON DELETE CASCADE` dari 018 → `ada=false, mau=true`.

### Paket 1 — §2 A bagian 2: rute dan layar diagnosa

**Yang dibaca lebih dulu**: `ClaimLifeDetailGCNM.xml` b3218–b6260 *(pohon grid)*, `SetDisease.xml`,
`SetSTS_Reject.xml`, `SearchDiagnose_act.xml`, ketiga pemuat `ClaimLifeDetailGCNM`, dan ketiga
section bergrid peserta. Rinciannya di tiket 08.

**Tiga temuan yang mengubah kode**, dan ketiganya datang dari turun satu tingkat lagi:

**1. Layar diagnosa terjangkau dari TIGA tahap, bukan satu.** Daftar rule yang **memuat**
`ClaimLifeDetailGCNM` tidak memuat `MedicalCheckClaimLife` — dan berhenti di situ berarti
menyimpulkan Medical Advisor tidak dapat menyunting diagnosa. Satu tingkat ke bawah membantahnya:
`ViewClaimDetailLifeGCNM` adalah `pyEditAction` grid peserta di `InputOSClaimLife` b18252,
`MedicalCheckClaimLife` b17416, dan `InputAkseptasiClaimLife` b17387 — grid yang **sama**
*(b15924, b15764, b15735)*. Ketiga pemuat ber-`pyWhenName` kosong dan ber-`pyPrivilegeName` kosong.
Bentuk kekeliruan yang sama untuk **keempat** kalinya *(Close Claim, OQ-H, OQ-K.2)*.

**2. Gerbangnya tujuh kali, bukan empat.** Bacaan pertama menyisir rentang grid saja. Tiga sisanya
menjaga `.ADMIN_NOTES` b2628, `.RECOMMENDATION` b7335, dan `.NOTES` b15234 — artinya peserta yang
sudah diputus membekukan **seluruh** isian layar Detail. Yang menemukannya **uji**, bukan pembacaan
ulang: `TestGerbangDiagnosaVERBATIMDariKorpus` membaca berkas korpus langsung dan menagih angkanya.

**3. `SetSTS_Reject` tidak dipanggil rute tolak/akseptasi** — ia aktivitas **pra-muat** b3224 wilayah
`S5`, dan grid bersarang di dalamnya. Yang ditiru invariannya, bukan mekanismenya. Ralat atas brief,
ditulis di tiket 08.

### ⛔ Cacat lintas-lapis KEENAM — dicegat sebelum berlayar

`models.Diagnosa` ronde pertama ditulis **tanpa satu pun tag JSON**, sementara `api.ts` sudah
mendeklarasikan `kodeIcd`, `pesertaId`, `groupDiagnose`. Go akan mengirim `KodeICD`, `PesertaID`,
`GroupDiagnose`; React membaca `undefined`; grid tampil dengan tiga kolom kosong dan **nol galat di
kedua sisi**. Ditutup `TestNamaJSONDiagnosaDikunci` *(dibuktikan merah dengan tag `kodeICD`)*.

Dan yang **ketujuh** hampir menyusulnya di berkas yang sama: `bolehUbahDiagnosa` membaca
`peserta.kodeStatus`, yang tidak pernah diseberangkan Go — gerbang layar akan selalu terbuka, dan
pemakai belajar mengabaikan 409. `KodeStatus` kini menyeberang, dan penjaga himpunan kunci
menyala lebih dulu.

### Dua penjaga yang menuduh hal yang benar — dipersempit dengan bukti

`TestNolNamaTabelTelanjangDiQuery` menuduh sebuah **kalimat prosa** *("menirunya dengan N UPDATE
berarti daftar berlubang")* sebagai `UPDATE <nama tabel>`. Komentar kini dibuang sebelum pencocokan
— sebagaimana `polaKomentar` di `models/kodestatus_test.go` sudah lebih dulu melakukannya, dengan
sebab yang sama persis. Dibuktikan **masih menggigit**: `FROM T_CLAIMLF_DIAGNOSE` telanjang yang
ditanam di `diagnosa.go` tetap tertangkap.

⚠️ Penjaga yang menuduh hal yang benar akan **dilonggarkan** orang, bukan dipatuhi — jadi ia
dipersempit sekarang, bukan dibiarkan sampai seseorang menuliskan pengecualian untuk berkasnya.

### Penjaga yang menyala — dan tidak satu pun dilonggarkan

| Penjaga | Kenapa menyala | Yang dikerjakan |
| --- | --- | --- |
| `TestPesertaJSONMembawaTanggalKejadian` | `diagnosa` lalu `kodeStatus` masuk marshaller | keduanya **ditulis** di daftar, bukan daftarnya dilonggarkan |
| `TestSetiapRuteNonGETPunyaPenjagaKasusTertutup` | tiga rute baru | ketiganya didaftarkan ke `rutePengubah` |
| `TestSetiapLayananPengubahMemeriksaKasusTerbuka` | layanan baru | `diagnosa.go` → `pagari` didaftarkan |
| `TestGerbangDiagnosaVERBATIMDariKorpus` | angka 4 salah | diralat menjadi **7**, dengan ketujuh barisnya disebut |
| `TestNamaJSONDiagnosaDikunci` | tag JSON hilang | tujuh tag eksplisit ditambahkan |
| `TestNolNamaTabelTelanjangDiQuery` | menuduh prosa | **penjaganya** yang dipersempit, lalu dibuktikan masih menggigit |

### Instrumen yang dibuktikan merah lebih dulu

`TestKaskadeHanyaPadaRelasiTerdaftar` *(cabut kaskade 018)* · `TestNolNamaTabelTelanjangDiQuery`
*(tanam `FROM T_CLAIMLF_DIAGNOSE`)* · `TestPencerminanDiagnosaDiTransaksiKeputusan` *(ganti
`baru.KodeStatus` → `sasaranKodeLama`)* · `TestPagariDiagnosaMemeriksaKetigaGerbang` *(cabut
`DiagnosaTerkunci`)* · `TestNamaJSONDiagnosaDikunci` *(tag `kodeICD`)*.

### TELEMETRI EKSEKUSI

| Hal | Paket 0 | Paket 1 |
| --- | --- | --- |
| Commit | `06bf711` | *(berikutnya)* |
| Go | 334 PASS · 0 FAIL · 34 SKIP | **356 PASS · 0 FAIL · 37 SKIP** |
| JS | 242 | **252** |
| `gofmt` / `go vet` / `go vet -tags db` | bersih | bersih |
| `tsc` / `npm run build` | bersih | bersih · 184,79 kB |
| Migrasi | `001`–`018`, `030`, `050`–`056` | tidak bertambah |
| Kebocoran | nol | nol |

⚠️ Angka SKIP naik 34 → 37 karena **tiga** uji `db` baru; ketiganya melewati selama `ORACLE_DSN`
belum dikonfigurasi. **Melewati bukan lulus** — dan `-migrate` belum dijalankan work owner *(brief
§4)*, jadi ketiganya memang belum pernah menyentuh Oracle.

### Paket 2 — §2 B: dokumen lewat outbox (butir be)

`UNGGAHAN_DIR` lahir *(config + `.env.example` + disambung `cmd/api`)*; migrasi **019**
`T_CLAIMLF_STORAGE`; `SisipDokumen`/`SatuDokumen`/`HapusDokumen`/kartu berkas di repository;
`services.Unggahan` dengan `Unggah`/`Unduh`/`Hapus`; `PelaksanaBerkasLokal` sebagai pelaksana
**stub**; tiga rute; layar `PanelDokumenPeserta` bertombol hidup dan penanda
*"URL menunggu penyambungan"* dicabut untuk baris yang sudah tertaut.

**Utang yang lunas**: kepala `models/dokumenbaru.go` sejak kelompok Dokumen berbunyi *"BELUM PUNYA
SATU PUN PEMANGGIL DI LUAR UJI ... harus dibaca sebagai UTANG"*. Pemanggilnya kini ada.

⛔ **Satu antarmuka berubah**: `PelaksanaEfek.Laksanakan` kini menerima `*repository.Tx`. Sebabnya
keras — efek `storage-unggah` menulis kartu berkas **dan** mengisi `T_STORAGE_ID`, dan kedua
tulisan itu tidak boleh berdiri sendirian bila penuntasan barisnya gagal. Pola yang sama dengan
`Jejak.Rekam`. Aman dilakukan: `PekerjaEfek` ternyata punya **nol pemanggil** — dibangun, tidak
pernah disambung.

### ⛔ CACAT LINTAS-LAPIS KETUJUH — dan yang paling buruk sejauh ini

Enam pendahulunya membuat layar **diam**. Yang ini membuat layar **berbohong**.

Pengenal dokumen adalah cap waktu `@CurrentDate("yyyyMMddhhmmssSSS")` *(`InsertDocument_Act.xml`
b648)* — **17 angka**, ≈2,0e16. `Number.MAX_SAFE_INTEGER` di JavaScript **9.007.199.254.740.991**
≈9,0e15. Jadi **setiap** pengenal dokumen berada di luar jangkauan aman, dan `JSON.parse`
membulatkannya diam-diam:

```
20260927103000123  ->  20260927103000124
```

Akibatnya tautan unduh menunjuk dokumen yang **tidak ada**, penghapusan mengenai baris yang salah
atau tidak ada, dan **nol galat di kedua sisi**. Kolomnya `NUMBER(19)` dan tipenya `int64` — Go
benar, TypeScript benar, hanya pertemuannya yang salah. Persis bentuk yang sudah enam kali terjadi.

Kini menyeberang sebagai **teks** *(`json:"id,string"`)*, dikunci dua sisi:
`TestPengenalDokumenMenyeberangSebagaiTeks` di Go dan uji tautan di
`PanelDokumenPeserta.test.ts` di sini.

⚠️ **Yang menemukannya uji, bukan pembacaan ulang.** Uji itu membandingkan tautan yang dirakit
dengan tautan yang diharap, dan angkanya berbeda **satu**. Tidak ada pembacaan kode yang akan
menangkapnya — pembulatan itu sah menurut kedua bahasa.

⚠️ `Diagnosa.ID` **tidak** ikut berubah, dan itu disengaja: ia dari sequence yang mulai dari 1,
jauh di dalam jangkauan aman. Yang menentukan bukan tipe Go-nya melainkan **besar nilainya** —
dan itu ditulis di kedua tempat supaya tidak ada yang "menyeragamkan" keduanya nanti.

### TELEMETRI EKSEKUSI — paket 2

| Hal | Isi |
| --- | --- |
| Go | **373 PASS · 0 FAIL · 37 SKIP** *(dari 356/37)* |
| JS | **254** *(dari 252)* |
| `gofmt` · `go vet` · `go vet -tags db` · `tsc` · build | bersih · 186,49 kB |
| Migrasi | `001`–`019`, `030`, `050`–`056` |
| Kebocoran | nol — fixture memakai `UJI-berkas.pdf`, nol nama orang |

## Giliran tiga modul 3 — paket 0: penyimpanan diadu dengan RDB-nya

Brief §1 menyuruh membaca **dua** rule yang belum dibaca. Keduanya membantah kode yang baru saja
ditulis, dan keduanya dengan kekeliruan yang **sama bentuknya**: membaca satu dari dua rule penulis
lalu menyimpulkan tentang tabelnya — kelima kalinya di modul ini.

**1. `TANGGAL_UPLOAD`** ditulis `Update_T_Storage_SQL.xml` b89, bukan oleh INSERT yang jadi sumber
019. Migrasi **020** menambahkannya. ⚠️ Dan rule itu memuat **dua bentuk tanggal berbeda** dalam
satu pernyataan — `EXPDATE` `DD/MM`, `TANGGAL_UPLOAD` `MM/DD` — sehingga nilai warisan keduanya
tidak dapat dibedakan untuk tanggal 1–12. Dicatat untuk A4.

**2. `IMAGEID` salah rumus.** Saya memakai pengenal dokumen *(cap waktu)*. Rumusnya ada di korpus
sejak awal: `GenerateImageID_SQL.xml` b85 `STANDARD_HASH('ASMPP'||FF9||SYS_GUID(),'MD5')`, dipanggil
`InsertGoogleStorage_Act.xml` b2226. Bedanya **bukan kosmetik**: cap waktu dapat **ditebak**, jadi
siapa pun yang tahu kapan sebuah berkas diunggah dapat menyusun kunci penyimpanannya.

⚠️ `'ASMPP'` **lima** huruf — bukan `ASMAPP`, awalan **token** penyimpanan (butir an). Keduanya
berdampingan di modul ini, berbeda satu huruf.

**Akibat yang ikut diperbaiki**: nama berkas lokal semula dirakit dari `T_STORAGE_ID`; sejak
`IMAGEID` benar keduanya berbeda, dan setiap unduhan akan gagal seketika rumusnya dipasang. Kini
berkunci pengenal dokumen lewat **satu** fungsi yang dipakai penulis dan pembacanya.

### ⛔ Satu uji yang gagal SEKALI lalu lulus tiga kali

`TestBatasUkuranTepatMasihDiterima` menulis 25 MiB ke disk sungguhan, dua kali. Ia gagal sekali di
tengah jalannya seluruh suite lalu lulus tiga kali berturut-turut.

Uji yang gagal secara acak akan **diabaikan** orang, bukan dipatuhi — persis prinsip yang sama
dengan penjaga yang menuduh hal yang benar. Batasnya kini dapat dikecilkan uji *(`DenganBatas`,
pola yang sama dengan `DenganFolder` yang sudah ada)*, dan **angka 25 MiB dijaga uji tersendiri**
yang juga membuktikan batas itu tidak dapat dimatikan lewat nol atau negatif.

### Dikunci dua sisi

`ImageIDDari` diadu dengan Oracle sendiri: `TestImageIDGoSamaDenganStandardHashOracle` menghitung
`STANDARD_HASH(:1,'MD5')` di Oracle atas **empat** masukan tetap dan membandingkannya. Uji murni
hanya dapat membuktikan Go konsisten dengan dirinya sendiri; ini yang mengadunya dengan pihak yang
sebenarnya. Literal tetapnya dihitung ulang bebas dengan Python — cocok.

### TELEMETRI EKSEKUSI — giliran 3 paket 0

| Hal | Isi |
| --- | --- |
| Go | **378 PASS · 0 FAIL · 38 SKIP** *(dari 373/37)* |
| JS | 254 *(tidak berubah)* |
| `gofmt` · `go vet` · `go vet -tags db` · `tsc` · build | bersih |
| Migrasi | `001`–`020`, `030`, `050`–`056` |
| Kebocoran | nol |

## Giliran 3 paket 2 — sensus paritas akhir, dan celah yang ditutup

135 rule di dua belas folder korpus dicacah **dengan skrip**, bukan dengan ingatan: nama tiap rule
dicari di seluruh korpus *(di luar berkasnya sendiri)* untuk menemukan pemanggilnya, lalu dicari
di `APP_RNM/` untuk menemukan padanannya. Hasilnya bab *"Sensus akhir 28-09-2026"* di PARITAS —
**satu baris per rule**, nol baris "nanti".

| Keadaan | Cacah |
| --- | ---: |
| ✅ ada padanan kode | 85 |
| ⚠️ tercatat di dokumen, belum berkode | 41 |
| dihakimi tangan sesudah rule-nya dibaca | 9 |
| ⛔ residu *(nol pemanggil, nol padanan)* | **0** |

### Delapan celah, dibaca satu per satu

**Dua dibangun.** `SetCurrencyID_Act` + `GetCurrencyID` — dan yang ditemukannya bukan sekadar
rule yang belum ditiru, melainkan **cacat yang sudah hidup**: `CURRENCYID` selama ini hanya
**dibawa** *(`WarisiKolom` menyalinnya dari baris sebelumnya)* dan tidak pernah **diterbitkan**.
Baris **pertama** seorang peserta karena itu lahir tanpa pengenal mata uang, dan
`HitungTotalPeserta` menolaknya dengan *"mata uang beragam"* — kalimat yang benar tentang hal yang
salah, dan yang membacanya akan mencari kesalahan di tempat yang bukan sebabnya.

⚠️ Ditiru di jalur **baca**, meniru **letaknya** bukan hanya hasilnya: rule itu
`pyPreDataTransform` `AdjustmentDetail_Section.xml` b1206 — ia berjalan setiap kali layar
adjustment dimuat.

**Satu ditolak dengan sebab.** `SetDisableAddButton` melakukan **satu** hal: `.IsCheck = false`
*(b139–141)*, sebagai pra-muat section `ClaimLifeDetail` *(b19169)*. Artinya **membuka layar
mencabut penanda "dipilih"**. Invariannya sudah kami punya lewat jalur yang benar —
`RejectOSClaimLife_Act` mencabut `IsCheck` saat baris **ditolak** *(tiket 05)*. Meniru
mekanismenya berarti setiap orang yang **melihat** layar meng-unselect pesertanya.

**Tiga milik modul PremiumList** *(`getMaxPagination_Act`/`_sql`, `InboxPremiumList_Claim`)*, dan
satu lagi — `DetailPolisLife` — adalah butir **av** yang menunggu PremiumList tiket 04.

**Dua bagian cabang `ContentNote != DEATH`.** `CountPesertaAkseptasiLife_SQL` dan kembarannya:
namanya "Count" tetapi ia **SELECT** — pencari baris akseptasi yang sudah ada untuk orang yang
sama. Prasyarat pemanggilnya `SavePesertaClaim.xml` **b1997**
`.IsCheck=="true" && ContentNote="DEATH"`, jadi ia bagian cabang yang PARITAS baris 1 sudah tandai
belum dibangun — **bukan celah baru**. ⛔ `NAME_OF_INSURED` dan `DOB` di query itu data orang:
dibaca saat jalan, tidak pernah disalin ke fixture.

### TELEMETRI EKSEKUSI — giliran 3 paket 1 dan 2

| Hal | A4 | Sensus |
| --- | --- | --- |
| Commit | `7cd7889` | *(berikutnya)* |
| Go | 387 PASS · 0 FAIL · 38 SKIP | **389 PASS · 0 FAIL · 38 SKIP** |
| JS | 254 | 254 |
| `gofmt` · `vet` · `vet -tags db` · `tsc` · build | bersih | bersih |
| Migrasi | tidak bertambah | tidak bertambah |
| Kebocoran | nol | nol |

## Giliran 6 paket 0 — bg: Beranda, 17 kelompok, PaletMenu, pola grid

**Sumbernya REFERENSI_UI, bukan korpus XML** — dan itu ditandai di setiap berkas. Beranda adalah
pengganti layar awal portal: PremiumList punya `PremiumLife_harness` *(kelas `Data-Portal`)*,
sedangkan **Claim Life tidak punya harness portal yang terekspor**. Bentuk Beranda karena itu
keputusan kami, `[kerangka aplikasi, bukan menu Pega]`.

| Unsur | Isi |
| --- | --- |
| Sidebar | **17 kelompok** = nama folder korpus apa adanya, termasuk ejaan janggalnya *(`Endorsment Fac In` tanpa `e`, `Komite Claim FacIn` tanpa spasi)* |
| Butir | **empat**, seluruhnya berbukti: Claim Life *(dua)*, PremiumList Life, Komite Claim Life |
| Kelompok tanpa butir | **empat belas**, berdiri terlipat, berketerangan `belum dimigrasi` |
| `PaletMenu` Ctrl+K | komponen tersendiri, membaca `lib/daftarMenu.ts` |
| Beranda | kartu per modul + cacah antrean Claim Life dari endpoint yang **sudah ada** |
| Pola grid | `exportXlsx` di-port dan **dipakai** Inbox Claim Life |

### ⛔ Logo: diperiksa, lalu TIDAK disalin

Brief mengizinkan menyalin `assets/logo-topbar.png` *"hanya bila itu logo perusahaan (bukan tulisan
Treaty)"*. Berkasnya dibuka dan dilihat: ia **wordmark "e-treaty"** — nama produk aplikasi Treaty,
bukan logo perusahaan. Tidak disalin; topbar tetap tanpa gambar.

### ⛔ Lima dari enam berkas pola grid TIDAK di-port

Brief melarang kode mati, dan lima di antaranya akan menjadi itu:

| Berkas | Sebab |
| --- | --- |
| `RecordForm`, `BilahSaringRegistry`, `saringRegistry` | bergantung tipe `MasterEntity`/`MasterField` — registry master yang aplikasi ini **tidak punya** |
| `CariSebaris`, `PilihCari` | bergantung `@tanstack/react-query` **dan** endpoint lookup `?q=` yang belum ada |
| **`exportXlsx`** | **nol impor** — mandiri, dan langsung berguna. **Di-port**, dengan dua penyesuaian karena `noUncheckedIndexedAccess` menyala di tsconfig kami dan tidak di sana; nol perubahan perilaku |

### ⚠️ Satu penyimpangan ekspor yang disengaja, arahnya berlawanan dengan kolomnya

Kolomnya **tepat yang tampil**, berurutan sama. Tetapi **tanggalnya keluar apa adanya (RFC 3339)**,
bukan dalam bentuk layar `28-09-2026 14:03` — sebab berkas lembar-sebar **diurutkan**, dan bentuk
layar itu diurutkan sebagai teks: Desember mendahului Februari.

### Dua penjaga lama yang menagih dengan benar, dan berpindah bersama kodenya

- *"butir sidebar TEPAT dua"* membaca literal `const BUTIR` yang kini tidak ada. Diganti tiga uji
  yang menjaga kebenaran barunya: **17** kelompok, **4** butir, **14** tanpa butir.
- *"fokus kembali sesudah palet ditutup"* mencari `fokusSebelum` di `Shell.tsx`; palet pindah ke
  berkasnya sendiri. Penjaganya ikut pindah — **dan bertambah tiga**: panah ber-`preventDefault`,
  sorotan melingkar, daftar kosong yang menyebut kuerinya. Ditambah satu yang baru: **palet
  membaca `daftarMenu`, bukan DOM** — sejak bg empat belas kelompok terlipat, dan `KelompokMenu`
  melepas anak kelompok terlipat dari DOM.

### TELEMETRI EKSEKUSI — giliran 6 paket 0

| Hal | Isi |
| --- | --- |
| Go | 415 PASS · 0 FAIL · 38 SKIP *(tidak berubah — nol kode Go)* |
| JS | **292** *(dari 254)* |
| `tsc` · build | bersih · 190,60 kB |
| Migrasi | tidak bertambah |
| Kebocoran | nol |
