# R01: Tracer — alur masuk renewal sampai kasus diterima tangga akseptasi

**What to build:** Seorang underwriter memasukkan **No. Polis + Renewal Date + Note**, menekan **OK**,
dan sebuah kasus renewal tercipta — dengan data polis lama sudah terisi sebagai nilai awal — lalu
**diterima tangga akseptasi tanpa penyesuaian apa pun**.

Ini irisan **tertipis yang lengkap** untuk renewal. Begitu ia hijau, seluruh premisnya terbukti: bahwa
renewal hanyalah pintu masuk lain menuju mesin yang sudah ada. Tiket RNW berikutnya tinggal menambah
layar.

⚠️ **Gerbang masuknya harus ditetapkan, bukan diport.** `[terverifikasi]` predikat gerbang masuk
siklus New Business **dirujuk nol kali** di korpus renewal, dan flow renewal **tidak dirujuk berkas
mana pun** — penentunya ada di konfigurasi work type/portal yang **tidak ikut terekspor**. Karena
tidak ada perilaku terekam untuk direproduksi, sistem baru menetapkannya eksplisit sesuai alur di atas.

`[terverifikasi]` Asal, `FlowAction\Renewal_FlowAct`: section `InputRenewal`, pre-activity
`InputOfferFacInEngineer_preACT`, post-activity `SetValidateDate_PostAct`, transform
`AddToListSuggestOfferFacIn_DT`. Flow `InputRenewalFacultativeIn` memuat 75 shape dan merujuk
**19 rule `When`** — seluruhnya predikat routing yang diwarisi dari NB, tidak ada predikat gerbang-masuk
di antaranya.

Berkas gambar `halaman depan renewal.JPG` **ada** sebagai referensi layar; isinya tidak dibaca sebagai
fakta.

**Blocked by:** **NB-08** (registry predikat `rules.Eval`) · **NB-11** (`acceptance.Next`)

**Status:** ready-for-agent

- [ ] Layar masuk menerima **No. Polis**, **Renewal Date**, dan **Note**; tombol **OK** memicu pembuatan kasus
- [ ] Kasus baru tercipta dengan pembeda siklus **`StatusBusiness = 2`**
- [ ] **Data polis lama tersalin sebagai nilai awal** pada kasus baru — underwriter cukup menyunting yang berubah
- [ ] Kasus hasil **diterima `acceptance.Next` tanpa penyesuaian** — tidak ada cabang khusus renewal di tangga akseptasi
- [ ] Ke-19 predikat routing dievaluasi lewat registry NB (`rules.Eval`), **bukan** salinan predikat baru
- [ ] Gerbang masuk **ditetapkan eksplisit** dan didokumentasikan sebagai keputusan rancangan — bukan diklaim sebagai perilaku terekam
- [ ] Setiap transisi menyebut rule Pega asalnya dalam komentar
- [ ] Tidak ada modul perhitungan baru yang dibuat oleh tiket ini
