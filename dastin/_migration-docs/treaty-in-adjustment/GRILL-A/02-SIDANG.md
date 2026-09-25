> Modul  : Treaty In Adjustment · Ronde A · 2026-09-24
> Peran  : interogator
> Masukan: `KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md` (G1-G4) · `TEMUAN-ADJUSTMENT-DITUNDA.md` · `METODE-GRILLING.md` Bagian VIII · bukti mentah pada 708 XML
> Status : TERBUKA
> Sifat  : TAMBAH-SAJA

# 02 · SIDANG ATAS KLAIM BAHAN SEBELUMNYA

Klaim yang dibawa masuk ronde ini, diadili satu per satu. Status mengikuti pola contoh:
**DIKUATKAN · DIPERLUAS · DIPERKECIL · SALAH KAPRAH · DITUTUP**.

---

## G1 · Addendum adalah `VERSI_KONTRAK`, bukan entitas tersendiri — **DIKUATKAN**

Buktinya bertambah, bukan berkurang: 19 properti `Int-treaty_in_edm` seluruhnya lapisan kepala
kontrak ditambah tiga penanda (`OLDID`, `EDMState`, `EDMMaterialType`). Diuji ulang lewat §8.6 Q1
di `06-PUTUSAN` GRL-04 dan bertahan — **tetapi bertahan sebagai gerbang, bukan atas dasar uji
ekspor**; lihat adjudikasinya di sana.

## G2 · "Nilai selisih historis tidak dapat direproduksi" — **SALAH KAPRAH**

Mesin selisih **tidak pernah** membaca lewat `OLDID`. Ia membaca `TreatyIn.OLDDATA`, potret penuh
kontrak yang dibekukan ke dalam `JSONDATA` addendum saat addendum dibuat
(`Activity/TreatyInSetAddendumToHistory.xml`). Bagian pertama G2 — bahwa `OLDID` menunjuk sasaran
bergerak — tetap benar; kesimpulannya tidak berlaku.

Sudah ditandai di kepala `KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md` dan `SEAM-ADJUSTMENT.md`, dan butir
eskalasi induk yang bersandar padanya sudah diralat.

## G3 · Tabel selisih sempit, kosakata besaran tertutup — **DIKUATKAN, dan daftarnya dituangkan**

Sesi ini menghitung ulang: **33 akar / 165 jalur daun**, bukan 31/161, dan **daftarnya** ditulis
(`PENGETAHUAN.md` §5.6). Selisih akarnya `LimitShareSummaryList` dan `LimitFacShareSummaryList`.
Sifat "hanya uang atau persentase" tidak dibantah apa pun.

## G4 · Ruang pengenal tidak bertabrakan — **SALAH KAPRAH pada bagian "tabrakan mustahil"**

Bentuk `‹kontrak›/Rnn` memang membuat pengenal addendum tidak pernah sama dengan pengenal kontrak.
Tetapi **tabrakan antar addendum sangat mungkin** — NA-02. Yang mustahil hanya tabrakan lintas
tabel, dan itu bukan yang dijaminkan §8.3.

Pertanyaan G4 yang sengaja ditunda — *"untuk revisi kedua, `OLDID` menunjuk apa?"* — **DITUTUP**:
`TreatyInRevisi_post` langkah 5 menetapkan `OLDID <- TreatyIn.ID` **sebelum** `ID` diubah, sehingga
`OLDID` menunjuk **baris yang dipilih pengguna**. Rantainya karena itu tidak seragam.

---

## Butir `TEMUAN-ADJUSTMENT-DITUNDA.md`

| Butir | Status | Alasan |
|---|---|---|
| **A-1** nomor revisi dari potongan teks | **DIPERLUAS** | bukan sekadar "dari teks": offsetnya **meleset satu** (NA-02) |
| **A-2** `ROWNUM = 1` mendahului `ORDER BY` | **DIPERKECIL** | `HASIL1` hanya diuji kosong/tidak; baris mana pun yang kembali menghasilkan uji yang sama |
| **B-1** picker meng-UNION tanpa pembeda | **DIKUATKAN, DIPERLUAS** | ditambah temuan saringan `ProportionType` yang timpang: varian XOL menyaring, varian non-XOL tidak |
| **D-1** `PositionUsername` di kelas addendum | **DIKUATKAN** | sama seperti induk; keputusan induk berlaku wajar, dikonfirmasi saat rekonsiliasi ADR |
| **G4** offset tetap | **DIPERLUAS** | menjadi NA-02 |

---

## `METODE-GRILLING.md` Bagian VIII

Bagian VIII disusun sebelum `PENGETAHUAN.md` ditulis. Menurut §3.7 metode itu sendiri, temuan
membatalkan klaim **fakta**, bukan keputusan **rancangan** — dan seluruh butir di bawah adalah
klaim fakta.

| Butir | Status | Yang benar |
|---|---|---|
| **§8.1** *"323 berkas BYTE-IDENTIK — `pzInsKey` sama, ruleset sama, versi sama"* | **SALAH KAPRAH, dua kali** | **nol** pasang byte-identik (urutan elemen dan dua cap ekspor selalu berbeda); dan irisan ber-`pzInsKey` sama adalah **317**, bukan 323 (NA-03) |
| **§8.3** gerbang Pengenal: *"tabrakan mustahil"* | **SALAH KAPRAH sebagai fakta** | tabrakan antar addendum mungkin (NA-02). Bentuk `‹kontrak›/Rnn` boleh bertahan **sebagai rancangan**; penjaga tabrakannya digali di cabang B |
| **§8.4** *"`OLDID` menunjuk KONTRAK; selisih historis tidak dapat direproduksi"* | **SALAH KAPRAH** | `OLDID` menunjuk **baris yang dipilih** — bisa revisi; dan selisih dibaca dari `OLDDATA`, bukan dari `OLDID` |
| **§8.4** *"penomoran hanya `R01..R99`"* | **SALAH KAPRAH** | patah pada revisi **kesepuluh**, bukan keseratus (NA-02) |
| **§8.4** `ROWNUM` | **DIPERKECIL** | nilainya tidak pernah dipakai |
| **§8.6 Q2** *"untuk revisi kedua, `OLDID` menunjuk apa?"* | **DITUTUP** | terjawab dari aturan — lihat G4 di atas |
| **§8.2** aturan sisi 323/56 | **DIKUATKAN** | dipakai sebagai kolom tetap di Lacak TDA (NA-04) |
| **§8.5** empat aturan bisnis pemilik proses | **DIKUATKAN** | dipakai sebagai dasar GRL-03; tidak dibantah apa pun dari ekspor |

**Tidak ada butir Bagian VIII yang dibantah sebagai keputusan rancangan.** Yang gugur seluruhnya
klaim fakta, dan ketiga putusan §8.3 tetap berdiri sebagai rancangan.

## TDA-07 · "Membuka kontrak dari picker Adjustment menulis kembali ke baris kontrak" — **SALAH KAPRAH**

Temuan milik sesi ini sendiri, disidangkan dengan pagar yang sama seperti klaim orang lain.
Pencabutannya lolos penyisiran tertutup (NA-06): syarat aktif, pemanggil tersapu habis, dan kedua
kontrol yang menyalakannya tidak dapat dicapai dari layar addendum.

Penggantinya lebih sempit dan **bukan** temuan Adjustment: dua tombol di layar penawaran
mengosongkan status akseptasi kontrak **dan menyimpannya** dalam satu tindakan, tanpa pembatasan
apa pun.

## ADR-0052 dan ADR-0055 terhadap revisi-di-tempat — **DIKUATKAN, dan sudah menanganinya**

Kedua tombol itu berpasangan dengan cabang 2 `Akseptasi_DT` (`RevisionState == 1`: Admin ->
SecHead -> Resolve Complete). Jadi ia **alur revisi-di-tempat milik Treaty In yang memang
dirancang**, bukan semata kehilangan status.

ADR-0052 butir 4 menghapus pemendekannya; ADR-0055 menetapkan bahwa perubahan atas versi yang sudah
disetujui **melahirkan versi baru yang mulai dari `DRAFT`**. Keduanya bersama-sama sudah memutuskan
nasib alur ini di sistem baru, dan butir 1 daftar eskalasi induk sudah menuliskannya pada paragraf
"Di sistem baru". **Tidak ada butir eskalasi baru yang perlu dibuat** — yang perlu hanyalah koreksi
atas satu kalimat peredam di butir yang sudah ada (`07-AUDIT` §3).

## TDA-11 · "Daftar pilihan menyatukan kontrak dan addendum tanpa saringan keadaan" — **DIKUATKAN, dengan sumber yang diperbaiki**

Butirnya berdiri, tetapi bukti yang semula saya pakai tidak lengkap. Seksi picker memuat **dua**
tanda yang bertentangan pada grid yang sama:

| Tanda | Isinya | Hidup? |
|---|---|---|
| `pyPageListProperty = TempMasterList.pxResults`, diisi `pyDeferLoadRetrievalActivity = TreatyLoadMasterJoinEdm` | `select … from treaty_in UNION select … from treaty_in_edm`, **tanpa `WHERE`** | **ya** |
| `pyGridProps/pyRDName = BrowseTREATY_IN` dengan `pyRDParams` memuat `StatusAkseptasi = "Resolve Complete"` | hanya `TREATY_IN`, tersaring | **tidak dapat hidup** |

Yang memutuskan bukan mana yang tertulis lebih rapi, melainkan **kolom yang tampil**: grid
menampilkan `.CARI1` sampai `.CARI7`, yaitu alias RDB List. `BrowseTREATY_IN` tidak mengembalikan
kolom bernama demikian; bila ia sumbernya, seluruh kolom akan kosong. Karena grid itu menampilkan
data, sumbernya daftar halaman (`METODE` §5.1).

**Akibat yang tidak terduga, dan ia menyambung ke §7.4.** Karena tidak ada saringan keadaan sama
sekali, kontrak yang sedang direvisi-di-tempat — `StatusAkseptasi = ""`, `RevisionState = 1` —
**tetap muncul** di daftar. Addendum yang dibuat darinya mewarisi rantai dua tingkat **dan** lahir
terkunci. Dugaan awal bahwa daftar itu menyaring `Resolve Complete` **dibantah oleh berkasnya
sendiri**.

Nasibnya tidak diputuskan di ronde A: daftar pilihan milik **cabang B**. Yang dikunci di sini hanya
faktanya.

## ADR-0052 · Premis "dua penyimpangan yang berjalan" — **SALAH KAPRAH pada satu dari dua**

Konteks ADR-0052 menyebut dua penyimpangan dari rantai empat tingkat. Diuji satu per satu:

* **Jalur revisi — berdiri.** `RevisionState = 1` memang memotong rantai menjadi dua tingkat, dan
  cabang 2 `Akseptasi_DT` memang dirancang untuknya. Yang dipertajam: pada **addendum** ia tidak
  dipilih siapa pun, melainkan **diwarisi** (NA-12).
* **Jalur penerima tugas kelompok — SALAH KAPRAH.** Ia bukan penyimpangan yang berjalan. Tidak ada
  aturan di kedua ekspor yang dapat menulis `Position = "ReasTreatyInGroupLeader"`; nilai yang
  mungkin hanya lima, dan itu bukan salah satunya (NA-10). Enam tempat yang menanganinya semuanya
  **pembaca**, dan cabang yang memindahkan keadaan itu ke tingkat berikutnya **mati**
  (`pyDisabled = true`).

Keempat keputusan ADR-0052 tetap berlaku; yang keliru adalah **pernyataan fakta di Konteksnya**,
dan `METODE` §3.7 melarang ADR memutuskan fakta tentang sistem lama. Koreksinya diajukan sebagai
**`REV-2`**, usulan revisi — bukan suntingan diam-diam. Putusannya **GRL-07**.
