# 15: De-identifikasi berkas kasus menjadi fixture rekonsiliasi

**What to build:** Lima berkas kasus nyata berubah menjadi fixture yang **aman disimpan di
repositori** — identitas hilang, **setiap angka utuh**.

Kelimanya berformat **JSON** dan memuat data pelanggan: nama tertanggung, alamat, serta **nama
operator** yang ikut terbawa metadata di kelimanya. Aturan proyek melarang nama orang dan data
pelanggan masuk artefak mana pun, termasuk fixture dan test — jadi berkas mentahnya **tidak boleh**
dipakai langsung oleh tiket rekonsiliasi.

Nilainya tetap besar: angka-angka di dalamnya adalah **masukan rekonsiliasi yang sesungguhnya**,
mencakup beberapa lini bisnis dan ketiga siklus.

## ⛔ Penyaringan berbasis pola DILARANG

Ini kriteria terpenting tiket ini, dan alasannya konkret. Menyaring dengan pencocokan kata kunci akan
**membuang angka yang justru diuji**: `TotalSumInsured` tertangkap filter *"insured"*, sementara
`CedingRetention`, `ShareCeding`, dan `ShareOfCeding` tertangkap filter *"ceding"* — keempatnya
**angka yang wajib dipertahankan**.

Ini jebakan yang sama bentuknya dengan kekeliruan pencocokan awalan yang sudah dua kali terjadi di
proyek ini. Bedanya, kali ini taruhannya angka rekonsiliasi.

**Gunakan daftar kunci eksplisit, ditinjau satu per satu.**

**BUANG (12 kunci, disetujui):**
`InsuredName` · `InsuredID` · `ASMAddress` · `SelectedLocationAddress` · `RoadName` · `ASMZipCode` ·
`RiskZipCode` · `Email` · `pxCreateOpName` · `MarketingName` · `CedingCoName` · `CoinsName`

**PERTAHANKAN:** seluruh field bernilai angka — termasuk keempat yang tertangkap filter palsu di atas.

**Bila ditemukan field identitas lain saat implementasi:** tambahkan ke daftar lewat **tinjauan
manual**. **Jangan beralih ke filter pola.**

**Blocked by:** None (can start immediately)

**Status:** ready-for-human — kriteria terpenuhi 01-10-2026, menunggu tinjauan work owner atas A16–A19

- [ ] Kelima berkas menghasilkan fixture turunan; berkas mentah **tidak** ikut masuk repositori
- [ ] Penghapusan memakai **daftar kunci eksplisit**; tidak ada pencocokan pola, substring, atau regex kata kunci di jalur penyaringan
- [ ] Kedua belas kunci pada daftar BUANG hilang seluruhnya dari fixture
- [ ] **Nol nama orang** dan **nol alamat** di fixture, diperiksa dengan pencocokan **batas kata**
- [ ] **Setiap angka identik dengan sumbernya** — dibuktikan pembandingan otomatis antara fixture dan berkas asal, bukan pemeriksaan mata
- [ ] Struktur dokumen dipertahankan; hanya nilai kunci identitas yang hilang, bukan bentuknya
- [ ] Daftar kunci yang dibuang didokumentasikan bersama fixture, sehingga penambahan berikutnya dapat ditinjau

## Comments

### 2026-10-01 — alat jadi, fixture BELUM masuk repositori (agent)

Kode: `APP_RNM/modul/nbfacin/backend/alat/deidentifikasi/` (paket + program `jalankan`).

| Kriteria | Keadaan |
| --- | --- |
| Daftar kunci eksplisit, tanpa pola | ✅ `DaftarBuang` 12 kunci; nilai dikosongkan, kunci tetap |
| Struktur dipertahankan | ✅ disalin token demi token: urutan kunci dan teks angka asli |
| Setiap angka identik, dibuktikan otomatis | ✅ `Periksa` membandingkan setiap daun; diuji dengan hasil yang sengaja dirusak |
| Nol nama/alamat, batas kata | ⚠️ `Periksa` mencari setiap nilai yang dibuang sebagai kata utuh di medan lain — **menemukan kebocoran** |
| Kelima berkas → fixture | ⛔ dijalankan ke scratchpad (di luar repo); **tidak** disimpan di repo |

`[terverifikasi]` Dijalankan atas kelima `DDL\P-5 *.txt` (JSON): verifikasi angka lolos di kelimanya.
Nilai yang dibuang **muncul lagi di 9 medan di luar daftar**: `pxCreateOperator`,
`AccumulationDescription`, `AccumulationCode`, `PIC`, `PICSuggest`, `CommentSuggest`,
`PropertiItemNote`, `SobName`, `TopRiskLocation`. Ditambah 48–128 nama kunci teks per berkas sebagai
kandidat tinjauan. Penambahan ke daftar BUANG = tinjauan manual (`../PERTANYAAN-LANJUTAN.md` butir 4). Uji mutasi: 5/5.

⚠️ `DDL\CONTOH` berisi 115 berkas XML lain (bukan JSON) — di luar lingkup tiket ini.

### 2026-10-01 — tindak lanjut review (agent)

- Program `jalankan` kini **menolak menulis** berkas yang kebocorannya tidak nol dan keluar dengan kode 1
  (dijalankan ulang atas dua berkas P-5: 7 dan 18 kebocoran, nol berkas ditulis).
- Flag `-kandidat` mencetak NAMA kunci bernilai teks (tanpa nilai) untuk tinjauan manual.
- ⚠️ `Periksa` hanya memburu nilai milik kunci BUANG. Nama atau alamat yang tidak pernah berada di kunci
  BUANG tidak terdeteksi otomatis — tinjauan manual atas daftar kandidat tetap wajib.

### 2026-10-01 — sembilan kunci ditambahkan; masih ada dua medan bocor (butir 32)

`DaftarBuang` kini 21 kunci (`TestDaftarBuangDisetujui`). Dijalankan ulang atas kelima `DDL\P-5 *.txt`
(keluaran ke scratchpad, di luar repositori): kebocoran **2 / 2 / 1 / 0 / 6**. Sisa bocor di **dua medan
baru**: `Comment` dan `OperatorID` (di bawah `QuotationData` / `OldData`). Alat menolak menulis keempat
berkas yang bocor, sesuai rancangan. **Fixture belum disimpan** — butir 32 mensyaratkan nol kebocoran di
kelima berkas. Keduanya menunggu keputusan work owner (tiket ini melarang filter pola). Daftar kandidat:
139 nama kunci unik, diserahkan untuk tinjauan; `pxUpdateOperator`/`pxUpdateOpName` 0 kemunculan (A15).

### 2026-10-01 — nol kebocoran; fixture menunggu satu keputusan (butir 36)

`Comment` dan `OperatorID` ditambahkan: `DaftarBuang` 23 kunci, kebocoran **0** di kelima berkas P-5.
Tujuh kunci mencurigakan ditinjau lewat profil bentuk tanpa nilai — tidak satu pun identitas
(rinciannya di `../KEPUTUSAN-30-09-2026.md` bawah butir 37). Kriteria "berkas mentah tidak ikut masuk
repositori" terpenuhi; **fixture belum disimpan** karena hasilnya masih memuat nomor kasus/polis dan teks
bebas yang tidak terperiksa alat ("Yang belum diputuskan" di register).

### 2026-10-01 — fixture disimpan (butir 38)

Kelima fixture di `backend/services/premium/testdata/kasus/` (nama netral, `README.md` di sana).
Kebocoran 0, dicek dua cara; penjaga `TestFixtureKasusSudahDibersihkan` terbukti gagal atas berkas mentah.
Alat kini membedakan `DaftarBuang` (identitas, sumber deteksi), `DaftarKosongkan` (teks bebas, alamat), dan
`JalurBuang` (jalur persis) — tetap daftar eksplisit, tanpa pola. Rincian dan A16–A19: register butir 38.

| Kriteria | Keadaan |
| --- | --- |
| Fixture turunan; berkas mentah tidak masuk repositori | ✅ |
| Daftar kunci eksplisit, tanpa pola | ✅ (`JalurBuang` jalur persis, bukan pola) |
| Setiap angka utuh | ✅ — seluruh angka berupa teks; daun di luar daftar identik (`Periksa` + cara kedua) |
