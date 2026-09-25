> Modul  : Treaty In Adjustment · ronde C
> Dibuat : 2026-09-24
> Sifat  : sidang atas klaim bahan sebelumnya — **DIKUATKAN · DIPERLUAS · SALAH KAPRAH · DITUTUP**
> Masukan: `PENGETAHUAN-PENGGRILL-ADJUSTMENT.md` (peta, **bukan bukti**), `PENGETAHUAN.md`,
>          `PETA-SUMBER-INDUK.md`, `KEPUTUSAN-GRILLING-ADJUSTMENT.md`

# Sidang ronde C

Enam klaim disidangkan. Seluruhnya disapu ulang atas ekspor dengan perkakas terkalibrasi; tidak satu
pun diterima karena tertulis.

---

## SC-1 — *"pro rata: 10 dari 53 langkah mati"* → **DIPERLUAS, dan perluasannya mengubah jawabannya**

**Klaim asal:** berkas penggrill §5.6 — pro rata termasuk "kemampuan yang mati", lima blok `//`
ditambah lima langkah bersarang.

**Hasil sidang:** angkanya **benar** — `tools/pre.py TreatyEDMProRateCalculation` mencetak blok mati
pada 2.1, 2.8, 2.9, 2.10, dan 8, beserta anaknya. Tetapi penggolongannya sebagai *"kemampuan mati"*
**kurang**, dan kekurangannya mengubah putusan:

> Pro rata **tidak mati**. Ia dipanggil (`TreatyEDMCalculateDifference` langkah 11, bersyarat
> `IsProRate == true`), sebagian besar langkahnya hidup, dan labelnya jujur. Yang membuatnya tidak
> berarti adalah **masukannya dipaku pada konstanta** (NC-01).

Bedanya bukan istilah. *"Kemampuan mati yang dinyalakan"* tunduk pada `METODE` §3.8 — rencana
penyalaan. *"Kemampuan hidup yang masukannya dipaku"* tidak: tidak ada yang perlu dinyalakan, karena
ia sudah berjalan; yang tidak ada adalah **kemampuan menyatakan tanggal berlaku**. Itu sebabnya
`GRL-15` membawa **atribut**, bukan **saklar**.

---

## SC-2 — *"`EDMEffective` hanya dibaca `TreatyCalculateProratePct`"* → **SALAH KAPRAH**

**Klaim asal:** berkas penggrill §8, catatan E3a.

**Hasil sidang:** kalimat itu menyebut **pembaca** dan diam tentang **penulis** — dan penulisnyalah
yang menentukan. `tools/tulis.py sapu .EDMEffective` mengembalikan lima penulisan; sesudah
`Param.EDMEffective` dipisahkan dari `TreatyIn.EDMEffective`, penulis yang sebenarnya **satu**:
`TreatyInSetEditPre` 1.5 → `= TreatyIn.Commencement`.

**Pelajaran yang berdiri sendiri:** penyapu mencocokkan **akhiran** jalur, sehingga `Param.X` ikut
terjaring bersama `TreatyIn.X`. Hasil sapuan yang tidak memisahkan keduanya akan **melaporkan lebih
banyak penulis daripada yang ada**. Dicatat di `01-TEMUAN.md` NC-01 dan diusulkan sebagai catatan
pemakaian `tulis.py`.

---

## SC-3 — *"share fakultatif: kemampuan yang belum pernah menyala"* → **DIPERLUAS**

**Klaim asal:** `PEMILAHAN-SISA-GRILLING.md` cabang E, E3b — `TreatyEDMDifferenceDeduction` langkah
4 dan 5 ber-blok `//`, berlabel *"Facultative share not enable yet in adjustment"*.

**Hasil sidang:** benar, dan **tidak lengkap**. Ada **ketiadaan kedua yang berdiri sendiri**:
`TreatyInDifferenceFacShare` — satu-satunya penulis akar `FacultativeShare` dan
`FacultativeShareBrokerage` — **tidak pernah dipanggil** dari rantai `TreatyEDMCalculateDifference`
(NC-04).

Akibatnya pada putusan: **menyalakan blok mati tidak menutup ketiadaan kedua.** Bila E3b dijawab
"nyalakan", dua dari empat akar tetap kosong dan tidak ada yang akan tahu sampai angkanya kembali
dan salah. Itu sebab `GRL-16` menjawab pada tingkat **bentuk mesin**, bukan pada tingkat saklar.

---

## SC-4 — *"`ActualValue.*` (108) memakai entitas yang sama, bukan entitas baru"* → **DIKUATKAN, dan batasnya dinyatakan**

**Klaim asal:** `PETA-TELUSUR-JSON.md` §1 induk, dikutip `PETA-SUMBER-INDUK.md` dengan tanda ⚠.

**Hasil sidang:** **benar dan tetap berlaku** — `ActualValue` memang tidak melahirkan entitas baru.
Tetapi ia menjawab pertanyaan yang berbeda dari C3: ia menjawab *"berapa entitas"*, bukan *"disimpan,
diturunkan, atau dibuang"*. Penilaian `PETA-SUMBER-INDUK.md` — **sebagian** — tepat.

Yang ronde ini tambahkan dan tidak ada di sana: `ActualValue` **berarti dua hal yang berlawanan arah
menurut `EDMState`** (NC-05). Berkas induk tidak dapat melihatnya karena ia menghitung jalur, bukan
membaca prasyarat langkah yang menulisnya.

---

## SC-5 — *"ADR-0043 belum pernah saya baca"* → **DITUTUP**

**Klaim asal:** `KEPUTUSAN-GRILLING-ADJUSTMENT.md`, penutup daftar sumber induk, beserta
perintahnya: *"Ia dibaca **sebelum** I2 diajukan, bukan sesudah."*

**Hasil sidang:** ketiganya dibaca sebelum ronde ini disusun. Hasilnya **mengubah bentuk I2**, bukan
mengisinya: payungnya turun menjadi **KONFIRMASI** dan anggaran GRILL turun (NC-06).

Dan ADR-0043 memuat satu hal yang **tidak terduga dari namanya** dan langsung menyentuh `GRL-17`:
ukuran keberhasilan migrasi *"setiap angka identik, bukan mirip, bukan dalam toleransi"*. Dibaca
polos, memberikan nomor urut baru melanggarnya — sehingga `GRL-17` wajib menyatakan mengapa tidak.

> Perintah membaca sumber induk lebih dulu terbayar di sini dengan cara yang tidak diperkirakan:
> bukan dengan menjawab butirnya, melainkan dengan **memunculkan tabrakan yang tidak akan terlihat
> dari sisi Adjustment mana pun**.

---

## SC-6 — *"mesin selisih adalah permukaan khas Adjustment"* → **SALAH KAPRAH**

**Klaim asal:** tersirat di seluruh bahan — `TreatyEDM*` diperlakukan sebagai aturan khas jalur
addendum, termasuk di `PENGETAHUAN.md` §5.6 yang menyebutnya *"khas jalur addendum"*.

**Hasil sidang:** `tools/cat56.py` mengembalikan **56 aturan** khas, dan **tidak satu pun** dari
`TreatyEDMCalculateDifference`, `TreatyEDMDifferencePremium`, `TreatyEDMDifferenceShare`,
`TreatyEDMDifferenceDeduction`, `TreatyEDMProRateCalculation`, maupun `TreatyCalculateProratePct` ada
di dalamnya. Hanya `TreatyInSetEditPre` yang 56-KHAS di antara aturan yang ronde ini sentuh.

**Awalan nama bukan penanda sisi.** `TreatyEDM*` berarti *"aturan yang menangani addendum"*, bukan
*"aturan yang hanya ada di ekspor Adjustment"*. Ini instans `METODE` §2.0 di tempat yang baru:
struktur yang terlihat — di sini, **konvensi penamaan** — bukan struktur yang berlaku.

**Akibatnya pada ke mana temuan dikirim:** lima dari tujuh temuan ronde ini adalah **temuan Treaty
In**, dan ronde ini memutuskan nasibnya untuk jalur addendum saja.
