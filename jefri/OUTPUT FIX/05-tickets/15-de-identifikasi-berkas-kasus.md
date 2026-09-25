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

**Status:** ready-for-agent

- [ ] Kelima berkas menghasilkan fixture turunan; berkas mentah **tidak** ikut masuk repositori
- [ ] Penghapusan memakai **daftar kunci eksplisit**; tidak ada pencocokan pola, substring, atau regex kata kunci di jalur penyaringan
- [ ] Kedua belas kunci pada daftar BUANG hilang seluruhnya dari fixture
- [ ] **Nol nama orang** dan **nol alamat** di fixture, diperiksa dengan pencocokan **batas kata**
- [ ] **Setiap angka identik dengan sumbernya** — dibuktikan pembandingan otomatis antara fixture dan berkas asal, bukan pemeriksaan mata
- [ ] Struktur dokumen dipertahankan; hanya nilai kunci identitas yang hilang, bukan bentuknya
- [ ] Daftar kunci yang dibuang didokumentasikan bersama fixture, sehingga penambahan berikutnya dapat ditinjau
