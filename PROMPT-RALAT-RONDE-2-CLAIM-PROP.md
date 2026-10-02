# PERINTAH PERBAIKAN — Ralat Grilling Claim Prop Ronde 2

Tempel seluruh blok di bawah ini ke Claude eksekutor. Tidak perlu skill grilling — ini pekerjaan ralat, bukan ronde baru.

---

Korpus `D:\XML\RNM_BRD\` READ-ONLY. Tulis HANYA ke `D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\`. Folder `D:\XML\nusantara-re\` DI-BLACKLIST. Setiap klaim WAJIB punya bukti `path berkas + nama rule + nomor step Pega`. **JANGAN pernah mengutip nomor baris XML.** Yang tidak ada di korpus = `[terbuka]`, JANGAN ditebak.

Berkas yang diperbaiki: `.scratch/claim-prop/grilling-ronde-2.md`. Jangan tulis ulang bab yang tidak disebut di bawah.

Lingkup korpus Claim Prop = **270 berkas**: Activity 112 · RDBList 54 · Section 36 · ReportDefinition 19 · FlowAction 12 · DataTransform 11 · Harness 11 · ConnectREST 6 · When 6 · Flow 1 · DecisionTable 1 · SystemSettings 1. Sebutkan cakupan di setiap sensus — "36 dari 36 Section", bukan "seluruh Section".

---

## TUGAS 1 — P18 tuntas. Dua laporanmu sendiri salah. Susun ulang R7.

Kontradiksi yang kamu angkat sudah saya selesaikan dengan menghitung anak dari setiap elemen `pyActionConditions` di 270 berkas. Hasilnya:

```
Section/AdjustmentDetail.xml           10 baris
Section/OutstandingClaim.xml            6 baris
Section/InputAcceptation_Est.xml        1 baris
Section/AdjustmentDetail_Section.xml    0
FlowAction/AdjustmentDetail.xml         0
--- tidak ada berkas lain yang berisi ---
```

Dua kesalahan yang harus kamu ralat:

1. Penyisiran Section melaporkan `pyActionConditions` **kosong di seluruh 36 berkas**. Itu tidak benar — tiga Section berisi, total 17 baris.
2. Penyisiran Harness melaporkan `AdjustmentDetail` punya **dua** baris. Salah dua kali: jumlahnya 10, dan **tidak ada Harness bernama `AdjustmentDetail`**. Berkas itu Section.

Penyebabnya **tabrakan nama**: tiga berkas berbeda bernama `AdjustmentDetail` tersebar di `Section/`, `FlowAction/`, dan `Section/AdjustmentDetail_Section.xml`.

Yang harus dikerjakan:

- **R7 dibatalkan.** Ia disimpulkan dari premis "pyActionConditions kosong di mana-mana" yang ternyata palsu. Susun ulang aturan baca gerbang tombol dari tiga Section yang benar-benar berisi, lalu tulis aturan barunya di bab aturan baca sebelum menyimpulkan apa pun.
- Periksa **setiap** kesimpulan ronde 2 yang bersandar pada R7 dan ralat satu per satu.
- **Aturan kerja baru, berlaku seterusnya:** rujuk berkas dengan **path lengkap dari akar modul** (`Section/AdjustmentDetail.xml`), tidak pernah dengan nama saja. Sebelum menyimpulkan "rule X begini", pastikan dulu tidak ada berkas bernama sama di folder lain.

## TUGAS 2 — P19 jumlahnya salah

`When/IsFire.xml` memang tidak ada. Tapi ia dirujuk di **6 berkas**, bukan dua:

```
Activity/PrintDLATreatyIn.xml
Activity/PrintFileAcceptance.xml
Section/MstAdjusterConsultant.xml
Section/InputAcceptation.xml
Section/OutstandingClaim.xml
Harness/MstAdjusterConsultant.xml
```

Ralat angkanya. Lalu catat untuk tiap rujukan: apa yang digerbangi, dan apa akibatnya bila rule itu tidak pernah ada — termasuk apakah Pega memperlakukan when yang hilang sebagai benar atau salah. Kalau tidak terbaca dari korpus, tandai `[terbuka]`, jangan tebak.

## TUGAS 3 — P20 turunkan dari pemblokir

Temuannya benar: `ResponseCode` **nol kemunculan** di 270 berkas, sedangkan `ResponseMsg` dieja benar. Satu ekspresi memuat dua ejaan sekaligus.

Tapi akibat yang kamu tulis terlalu besar. Dua tempat pemakaiannya:

```
Activity/HitServiceToKasir_Act.xml  step 9.8   tanpa when
Activity/HitServiceToKasir_Act.xml  step 13.5  tanpa when

  StatusKasir = @if(StatusServiceKasir.ReponseCode = 1,
                    "Akseptasi Sudah Masuk ke Kasir",
                    StatusServiceKasir.ResponseMsg)
```

Tidak ada satu pun step yang bercabang pada nilai itu. Yang terdampak hanya **teks status yang tampil**, bukan alur. Tidak ada transaksi yang gagal karenanya.

Ralat jadi cacat kosmetik. P20 keluar dari daftar pemblokir `/to-spec`, kecuali kamu bisa menunjukkan dari `ConnectREST/SendAcceptationToKasir.xml` bahwa pemetaan responsnya memang menyasar nama yang salah eja — **kalau tidak terbaca, tulis `[terbuka]`, jangan simpulkan**.

## TUGAS 4 — angka gerbang UI mati tidak cocok

Kamu melaporkan **57**. Sensus saya atas 36 Section:

```
180  elemen punya pyCondition
119  pyVisible = OTHER     → syarat berlaku
 44  pyVisible = ALWAYS    → syarat TIDAK berlaku   (pasti mati)
 17  pyVisible kosong      → arti default tidak terbaca dari korpus
```

44 pasti, 61 batas atas. **57 tidak jatuh di keduanya** — jelaskan dari mana angkamu, lalu rekonsiliasi.

Kemudian: arti `pyVisible` **kosong** adalah pertanyaan yang sejenis dengan kode arah sisi-kosong yang kamu temukan sendiri. Jangan diputuskan. Jadikan satu OQ:

> `[terbuka]` Bila `pyVisible` kosong sementara `pyCondition` terisi, apakah syaratnya berlaku? 17 elemen di Section, ditambah 3 kasus kode arah sisi-kosong di Activity. Cara menutup: buka satu elemen di Pega Designer dan lihat apakah kotak kondisinya aktif.

Daftar triase Q1 **jangan** dipakai sebelum ini selesai — batas bawah dan batas atasnya beda 17.

## TUGAS 5 — yang sudah saya verifikasi BENAR, jangan diulang

| Temuan | Status |
| --- | --- |
| Aturan `pyStepsBlockName`: hanya `//` berarti remark | **benar**. Nilai se-Claim Prop: `//` 68 · `EXIT` 2 · `END` · `FAIL` · `TO` · `UP` · `.` masing-masing 1. Sensus 68 remark tetap sah. |
| P21 — empat RequestType hidup tanpa rule | **benar**, keempatnya tidak ada dan keempatnya dirujuk |
| P22 — `RDBList/Update_T_Storage_SQL.xml` dua format tanggal | **benar**: `EXPDATE` pakai `DD/MM/YYYY HH24:MI:SS`, `TANGGAL_UPLOAD` pakai `MM/DD/YYYY HH24:MI:SS`, satu statement |
| Kategori keempat kode arah (salah satu sisi kosong) | **benar**, gabungkan OQ-nya dengan TUGAS 4 |
| 84 gerbang Activity tidak berlaku | **benar**, cocok dengan sensus saya |

## ATURAN SELAMA RALAT

1. Jangan buka ronde baru. Jangan tambah temuan baru kecuali muncul sendiri saat meralat — kalau muncul, taruh di bab terpisah "Temuan sampingan saat ralat".
2. Setiap ralat ditulis sebagai blok `⚠️ RALAT <tanggal>` di dalam bab aslinya, **jangan hapus teks lama** — supaya jejak kesalahannya terbaca.
3. Setiap angka sensus wajib menyebut penyebutnya: "44 dari 180", bukan "44".
4. Sebelum menulis "tidak ada di korpus", sebutkan berapa berkas yang kamu cari dan di folder mana. Klaim ketiadaan tanpa cakupan = ditolak.
5. JANGAN tulis kode Go/React. JANGAN sentuh korpus.

## SELESAI BILA

Kelima tugas tertulis di `grilling-ronde-2.md`, R7 tersusun ulang, dan daftar pemblokir diperbarui: P18 tertutup, P20 keluar dari daftar, P19 angkanya benar, dan OQ `pyVisible` kosong masuk sebagai pemblokir baru untuk daftar triase Q1.

Laporkan ringkas: apa yang diralat, apa yang tetap, dan pemblokir apa yang tersisa sebelum `/to-spec`.
