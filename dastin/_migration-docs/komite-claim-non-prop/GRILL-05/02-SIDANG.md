> Modul  : Komite Claim Non Prop · Ronde 05 · 2026-09-20
> Peran  : interogator
> Masukan: `FAKTA-A1B.md` (F-13 … F-24, ronde 4, 2026-09-20) · `PUTUSAN-01.md` §8 (2026-09-20) · bukti baru pada `01-TEMUAN.md`
> Status : DITUTUP 2026-09-20
> Sifat  : TAMBAH-SAJA

# 02 · SIDANG ATAS TEMUAN RONDE SEBELUMNYA

Keluaran ronde 4 berbentuk **fakta** `F-13 … F-24`, bukan dakwaan. Sidang di sini menguji
tiap fakta terhadap langkah yang baru dibaca ronde ini, dengan lima putusan yang sama.
Satu blok tambahan menyidangkan perlakuan K-02 pada `PUTUSAN-01.md` §8, karena bukti ronde
ini menjatuhkannya.

---

## F-13 · `KomiteLoop` = jumlah baris hasil laporan — **DIPERLUAS**

**Bukti yang diperiksa ulang.** `CreateChildKomiteCNP_Act.xml` sub-langkah 26.13 dan 26.14,
keduanya menulis `ChildWorkPage.KomiteLoop = @Utilities.SizeOfPropertyList(pyReportContentPage.pxResults)`;
sub-langkah 26.15 menimpanya dengan `1` untuk jalur `VINCENTVERNANDO_1`. Nilai awalnya `1`,
ditetapkan langkah 9.

**Apa yang dilewatkan ronde 4.** Bahwa nilai itu ditulis **dua kali berturut-turut** oleh
dua sub-langkah berbeda, dan bahwa pembandingnya — `When/IsKomiteLoop.xml` — berbunyi
`.AcceptStatus = "1"` **dan** `.KomiteCount <= .KomiteLoop`, dengan `KomiteCount` lahir
bernilai `1`. Lingkar karena itu berjalan tepat `KomiteLoop` kali, bukan `KomiteLoop+1`.

**Perlakuan.** Simpan jumlah jenjang sebagai turunan dari baris roster yang benar-benar
terbentuk, bukan dari jumlah baris sumbernya.

---

## F-14 · Penyaring roster `A AND C AND B`, urut `.DEGREE ASC` — **DIKUATKAN dan DIPERLUAS**

**Bukti yang diperiksa ulang.** `KomiteRouter.xml` langkah 6.1 memilih baris
ber-`KomiteAproval==0` dari daftar yang sudah terurut itu, lalu keluar dari iterasi
(transisi `6/6`).

**Apa yang dilewatkan ronde 4.** Ronde 4 mencatat urutannya tanpa memeriksa siapa yang
membacanya. Pembacanya adalah router, dan router berhenti pada kecocokan **pertama**.
`.DEGREE ASC` karena itu benar-benar urutan giliran, bukan sekadar urutan tampilan — ia
menentukan siapa memutuskan lebih dulu. Yang bertambah: penentu kedua, `.KomiteCount`,
hidup berdampingan tanpa pernah diperiksa terhadapnya (N-01).

**Perlakuan.** Bawa `.DEGREE` sebagai urutan giliran yang dinyatakan, dan jadikan hitungan
jenjang turunan dari roster, bukan penentu tandingan.

---

## F-15 · Perbandingan ambang dimatikan, tetapan `0` dan `25000001` — **DIKUATKAN dan DIPERLUAS**

**Bukti yang diperiksa ulang.** Sub-langkah 26.3 memasang parameter laporan **tanpa**
`LIMIT_BOTTOM` sama sekali, berkomentar "Untuk sementara dihilangkan
(Param.LIMIT_BOTTOM==Local.TotalValueAdjust)". Sub-langkah 26.4 memasang `0`, 26.5 memasang
`25000001`, 26.6 memasang `0` untuk jalur bersyarat.

**Apa yang dilewatkan ronde 4.** Sub-langkah 26.2 berbunyi
`set Local.TotalValueAdjust = Local.TotalValueAdjust` — penetapan nilai kepada dirinya
sendiri, sisa dari perbandingan yang dimatikan itu. Dan 26.6 berjalan **sesudah** 26.5,
sehingga menimpanya (N-10).

**Perlakuan.** Bawa ambang sebagai data (H-1) dan jangan bawa satu pun dari ketiga sisa ini.

---

## F-16 · Tetapan `Flagkomite`, cabang mati langkah 14 — **DIPERLUAS**

**Bukti yang diperiksa ulang.** Langkah 10 (empat tetapan dan `Flagkomite=0`), langkah 12,
13, dan 14. Cabang mati terbukti: `>Local.LimitPersenMax && <=Local.LimitPersenMaxDivHead`
dengan kedua tetapan bernilai `30.00`.

**Apa yang dilewatkan ronde 4.** Dua hal. Pertama, `Local.CekLimitPersen` diisi dari
`pyWorkPage.TreatyInMaster.RNMShare` — persentase bagian, bukan nilai uang; maka dua dimensi
bercampur dalam satu penanda. Kedua, `Flagkomite` tetap `0` bukan hanya ketika nilai melebihi
50 juta, melainkan **juga** ketika `RNMShare` melebihi 30 — dan pada keadaan itu
`Param.LIMIT_BOTTOM` tidak pernah diisi sebelum laporan dipanggil (N-09).

**Perlakuan.** Nyatakan aturan seleksi sebagai tabel dua dimensi — nilai dan bagian — dengan
setiap kombinasi punya hasil, tidak ada kombinasi tanpa aturan.

---

## F-17 · `Local.TotalValueAdjust` memuat baris terakhir — **DIKUATKAN**

**Bukti yang diperiksa ulang.** Langkah 11 beriterasi atas
`pyWorkPage.ClaimData.AdjustmentList` dan tiap putaran menulis
`Local.TotalValueAdjust = .ValueAdjustment`, menimpa yang sebelumnya.

**Apa yang dilewatkan ronde 4.** Tidak ada yang penting. Yang bertambah hanya bahwa
nilai itu kemudian dipakai langkah 12 dan 13 sebagai bahan penentu kelas komite — sehingga
kekeliruan "terakhir, bukan total" langsung menjadi kekeliruan wewenang, bukan sekadar
tampilan.

**Perlakuan.** Hitung total sebagai total, dan jadikan ini uji paritas yang dirancang gagal.

---

## F-18 · `.JABATAN` ditimpa oleh `KomiteList = .ComiteeClaim` — **DIPERLUAS**

**Bukti yang diperiksa ulang.** Sub-langkah 26.8.1 menulis
`Primary.ComiteeClaim(<LAST>).IDKomite = .ID`, sementara 26.8.2 menulis
`ChildWorkPage.KomiteList(<LAST>).IDKomite = .JABATAN`. Sub-langkah 26.13 kemudian menyalin
`ChildWorkPage.KomiteList = .ComiteeClaim`.

**Apa yang dilewatkan ronde 4.** Tiga hal. Satu, kolom pemenangnya adalah `.ID`, bukan
`.JABATAN`. Dua, 26.8.2 **dilewati seluruhnya** ketika `Local.Subjectivity==true`, sehingga
pada jalur bersyarat `KomiteList` hanya pernah terisi lewat penyalinan 26.13. Tiga,
sub-langkah 26.12 kemudian menimpa `IDKomite`, `KomitePost`, dan `Initial` untuk empat nama
orang tertentu — sehingga nilai dari basis data hanya bertahan bagi orang di luar keempatnya.

**Perlakuan.** Satu sumber untuk jabatan, dinyatakan di roster, tanpa penimpaan per orang.

---

## F-19 · `.OPERATOR_ID=="DARTO"` → `"CHRISTINEANGELINA"` — **DIPERLUAS**

**Bukti yang diperiksa ulang.** Sub-langkah 26.8.3, persis seperti dicatat ronde 4.

**Apa yang dilewatkan ronde 4.** Sub-langkah 26.12 memuat **empat** pemetaan orang→jabatan
yang ditulis di dalam rule: `CHRISTINEANGELINA` → "Claim Dept. Head"/"CA",
`CHRISTOPMARHASAK` → "Technic Div. Head"/"CM", `Himawan` → "Operational Director"/"HY",
`NANDINA` → "Technical Director"/"NC". Satu di antaranya ditulis dengan huruf campuran
sementara tiga lainnya huruf besar — pembandingnya adalah kesamaan teks.

**Perlakuan.** Pindahkan seluruh pemetaan orang→jabatan ke data (D-5), termasuk keempat baris
ini, dan jangan bawa satu pun nama orang ke dalam kode.

---

## F-20 · Jalur tutup/tolak memaku satu penyetuju — **DIKUATKAN dan DIPERLUAS**

**Bukti yang diperiksa ulang.** `CreateChildKomiteCloseNP_Act.xml` langkah 8 menulis enam
medan `KomiteList(1)` dan lima medan `ComiteeClaim(1)` dengan nama, surel, dan jabatan yang
dipaku, lalu `KomiteLoop = 1`. Langkah 6 menetapkan `IsReject` dari
`TypeComentAnalysis=="5"`.

**Apa yang dilewatkan ronde 4.** Bahwa `IsCloseFile` pada langkah yang sama ditulis dalam
dua ejaan berbeda untuk anak dan induk (N-05), dan bahwa langkah 12 menetapkan
`pyWorkPage.CNPStatusCase = "COMITEE ACCEPTANCE (DEPT. HEAD)"` — sebuah keadaan klaim yang
ditulis sebagai teks di dalam rule, sejalan dengan jabatan yang dipakunya.

**Perlakuan.** Roster berdasarkan peran (H-2), dan keadaan klaim sebagai kode, bukan kalimat.

---

## F-21 · Gerbang hanya menandai; `@hasMessages` menghentikan pembuatan — **SALAH KAPRAH**

**Bukti yang diperiksa ulang.** `CreateChildKomiteCNP_Act.xml` langkah 29 dan
`CreateChildKomiteCloseNP_Act.xml` langkah 10, beserta langkah-langkah **sesudahnya**.

**Apa yang dilewatkan ronde 4.** Aksi pra-syaratnya adalah *lewati-step*, bukan *keluar
activity*. Yang berhenti hanyalah pemanggilan `pxAddChildWork`. Langkah berikutnya tetap
membuka kunci klaim, membaca `pxCoveredInsKeys(<LAST>)` seolah sirkulasi baru saja lahir,
menulis nomornya ke klaim, menyimpan, dan mengirim surat (N-07).

**Perlakuan.** Perluas H-5: yang menolak bukan hanya langkah pembuatan, melainkan seluruh
urutan yang bersandar pada keberhasilannya — satu transaksi, gagal bersama.

---

## F-22 · Sirkulasi dapat lahir tanpa jenjang — **DIKUATKAN, sebabnya ditemukan**

**Bukti yang diperiksa ulang.** Rantai 26.3 → 26.4/26.5/26.6 → 26.7, dan langkah 10–14 yang
menentukan `Flagkomite`.

**Apa yang dilewatkan ronde 4.** Sebabnya. Ronde 4 mencatat bahwa hal itu **mungkin**; ronde
ini menunjukkan **kapan**: ketika `Flagkomite` tetap `0`, tak satu pun sub-langkah pemasang
ambang berjalan, dan laporan dipanggil dengan `Param.LIMIT_BOTTOM` yang tidak pernah diisi.

**Perlakuan.** H-6 tetap berlaku tanpa perubahan — ditolak saat dibuat. Yang bertambah:
perkara ini tidak lagi menunggu bukti apa pun.

---

## F-23 · `VINCENTVERNANDO_1` menulis `KomiteAproval=""` — **DITUTUP**

**Bukti yang diperiksa ulang.** Sub-langkah 26.8 dilewati seluruhnya untuk operator itu;
26.9 dan 26.10 memasang satu anggota dengan surel `klaim5@nusantarare.com`, `KomitePost="VC"`,
dan `KomiteAproval=""`; 26.15 memaksa `KomiteLoop=1`.

**Apa yang dilewatkan ronde 4.** Bahwa router menyeleksi dengan `==0` sedangkan yang ditulis
`""` (N-11) — dan bahwa seluruh jalur itu bergantung pada nama satu operator pembuat case.

**Perlakuan.** Tidak ada. H-7 sudah menetapkan jalur ini tidak dimigrasikan; perkara ditutup
oleh ketetapan, bukan oleh bukti.

---

## F-24 · Nol `pyPrivilege` pada empat rule A-1b — **DIPERLUAS**

**Bukti yang diperiksa ulang.** Sapuan diperluas dari empat rule ke seluruh 59 berkas folder
Komite: elemen `pyPrivilegeName` muncul **26 kali** dan **kosong pada seluruh 26 kemunculan**;
tidak ada satu pun `pyPrivilegeName` berisi. `KomiteRouter.xml` menyatakan
`pyPrivilegeClass` = `ASM-FW-GCNMFW-Work-Komite` dengan `pyPrivilegeName` kosong.

**Apa yang dilewatkan ronde 4.** Lingkupnya. Ketiadaan wewenang bukan sifat empat rule
pembentuk, melainkan sifat **seluruh modul**: tidak ada satu pun rule yang menuntut hak.

**Perlakuan.** ADR-0006 (RBAC dari nol) berlaku penuh untuk modul ini; tidak ada hak yang
dapat diwarisi dari sistem lama karena tidak ada yang tertulis.

---

## Tambahan · Perlakuan K-02 pada `PUTUSAN-01.md` §8 — **DIKUATKAN**

**Bukti yang diperiksa ulang.** `PUTUSAN-01.md` §8 menguatkan K-02 "dengan syarat" dan
menuliskan perlakuan: tetapkan aturan jenjang aktif sebagai "baris pertama yang belum
memutuskan". `KomiteRouter.xml` langkah 6.1 — pra-syarat `.KomiteAproval==0` dengan
`WhenFalse=3`, transisi `WhenTrue=6` dan `WhenFalse=6` — memberlakukan persis itu.

**Apa yang dilewatkan.** Perlakuan itu ditulis tanpa membuka `KomiteRouter.xml`; ia benar
tetapi bersandar pada dugaan. Kini ia bersandar pada langkah. Yang benar-benar terlewat
adalah penentu keduanya: empat langkah keranjang yang memilih dari `.KomiteCount` (N-01,
N-02), yang tidak disebut sama sekali oleh perlakuan itu.

**Perlakuan.** Naikkan dari tafsir menjadi ketetapan, dan tambahkan bagian yang hilang:
satu penentu jenjang aktif, hitungan sebagai turunannya.
