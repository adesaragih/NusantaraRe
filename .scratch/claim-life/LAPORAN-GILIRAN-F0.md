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
