# 12: Akseptasi lini financial — filter ambang limit, tanpa tangga berjenjang

**What to build:** Kasus lini financial menemukan jabatan berwenangnya lewat **penyaringan ambang
limit**, bukan lewat tangga bertingkat. Lini ini **tidak bereskalasi**.

⛔ **Bentuk tabelnya berbeda, dan menyamakannya membuat tangga financial macet total.** Tabel limit
financial **tidak punya** kolom jabatan atasan, tidak punya kolom antrean, dan kolom limitnya bukan
ambang tunggal melainkan **empat ambang per jenis pertanggungan**: bond, credit CL, credit NCL, trade.

⛔ **Ejaan jabatannya PAKAI SPASI**, berbeda dari bentuk standar yang tanpa spasi. Bila kode
mencocokkan dengan token tanpa spasi, **tidak ada approver financial yang pernah ditemukan**.
**Normalisasi ejaan dilarang** — menghapus spasi agar seragam adalah perubahan perilaku, bukan
migrasi.

**Yang tidak ada** adalah eskalasi berjenjang dan antrean per jabatan. **Yang tetap ada dan wajib
diport** adalah **filter ambangnya**: tiga rule SQL menyaring dengan `WHERE <kolom_limit> <= nilai`,
masing-masing untuk satu jenis pertanggungan, tanpa penggabungan tabel dan tanpa pengurutan.
Membuang `WHERE` menghapus satu-satunya kontrol wewenang yang dimiliki lini ini.

Keluarannya karena itu **daftar jabatan yang limitnya menampung nilai** — bukan satu jabatan tujuan
berikutnya.

⚠️ Kedua mekanisme **tidak disatukan di balik satu abstraksi**. Memaksa bentuk ini ke dalam bentuk
tangga berarti mengarang langkah naik yang tidak ada.

**Blocked by:** 11

**Status:** ready-for-agent

- [ ] Satu query per jenis pertanggungan (bond, credit CL, credit NCL), memakai kolom ambangnya masing-masing
- [ ] **`WHERE` dipertahankan apa adanya**; tidak ada penyaringan tambahan di luar ambang limit
- [ ] Pencocokan jabatan memakai ejaan **berspasi**; tidak ada normalisasi
- [ ] Keluaran berupa **daftar jabatan**, bukan jabatan tujuan berikutnya
- [ ] Tidak ada langkah naik dan tidak ada token antrean untuk lini ini
- [ ] Kedua bentuk tetap dua jalur terpisah di kode; kasus uji membuktikan bentuk standar **tidak** dipakai untuk lini financial dan sebaliknya
- [ ] Fixture bentuk ini juga **tanpa kolom nama**
