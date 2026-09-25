---
status: aktif
golongan: perubahan
---

# 21: Perhitungan yang gagal menghasilkan keterangan apa yang gagal, bukan angka nol dan bukan kurs satu

*Asal: `DAFTAR-PEKERJAAN.md` `P-28` · `ADR-0035` · `SPEC-INVARIAN.md` `INV-38`, `INV-42`.*

**What to build:** **PK** yang perhitungannya tidak dapat diselesaikan — kurs tidak ada, pembagi nol,
nilai sumber kosong — melihat **keterangan apa yang gagal**. Ia **tidak** melihat angka nol, dan
**tidak** melihat hasil yang diam-diam memakai kurs satu.

**Kenapa begini:** Ini instans `ADR-0035` yang paling sering terjadi, dan sistem lama melakukannya di
dua tempat sekaligus. `SaveTreatyInDetail_Act` menulis **`DEDUCTION1 = 0`, `DEDUCTION2 = 0`,
`SPREAD_RNM_SHARE_VALUE = 0`** ke tabel datar pada setiap penyimpanan, dan `DetailCalculationROL`
menulis **`9989998`** sebagai persentase ketika limitnya nol. Keduanya **kegagalan yang menyamar
sebagai nilai**, dan bedanya cuma satu: yang kedua begitu ganjil sampai orang curiga, yang pertama
**sama sekali tidak dapat dibedakan dari nol yang benar**.

> **Nol yang sebenarnya kegagalan tidak dapat dibedakan lagi setelah tersimpan.** Itu sebab `G2`
> menempatkan `ADR-0035` di **butir (b)** — perubahan yang menentukan bentuk dan tidak dapat dipasang
> belakangan.

**Persyaratan:** `ADR-0035` · `INV-42` (nilai uang tidak pernah memuat angka penanda kegagalan;
ketidakmampuan menghitung menghasilkan **kosong**, bukan angka) · `INV-38` (kurs tidak pernah nol,
tidak pernah dipaksa satu) · `INV-37` (`NILAI_IDR` terisi ⟹ kurs, tanggal kurs, sumber kurs terisi)

**Tidak termasuk:** **Penulisan nol oleh jalur migrasi** — baris warisan yang **sudah** memuat nol
atau `9989998` dibawa **apa adanya** (`ADR-0042`), dan pemisahannya dari nol yang benar adalah
pekerjaan irisan migrasi batch 2. Irisan ini melarang sistem baru **membuat** yang baru.
**Berapa baris warisan yang terkena** — **`Uji AQ`**, belum dijalankan.
**Nasib kolom `PERSEN_ROL` itu sendiri** — `F-15`, menunggu satu kalimat pemilik proses.

**Jalur gagal:** Konversi tanpa kurs -> hasilnya **kosong beserta keterangannya**, bukan nol ·
Pembagian dengan nol -> keterangan, bukan angka penanda · `NILAI_IDR` terisi sementara kursnya kosong
-> **ditolak** `INV-37` · Sebuah kolom uang memuat `9989998` -> **ditolak** `INV-42`.

**Uji:** **Negatif:** simpan nilai IDR tanpa kurs; paksakan kurs satu; sisipkan `9989998` ke kolom
persentase; bagi dengan nol lewat jalur perhitungan.
**Positif — dan ia yang menangkap larangan yang terlalu lebar:** sebuah besaran yang **memang bernilai
nol** — potongan nol persen yang disepakati — **tersimpan sebagai nol dan diterima**. Sebuah
constraint yang menolak setiap nol lulus keempat uji negatif di atas dan **menolak data yang sah**.

**Menggantikan:** `SaveTreatyInDetail_Act` yang menulis nol tetap ke tiga kolom datar, dan
`DetailCalculationROL` yang menulis `9989998` — keduanya tercatat di `TETAPAN-DI-KODE.md` §2b dan §3.
Golongannya **PERUBAHAN**: perilakunya sengaja dibalik, sehingga **cocok-tidaknya terhadap data lama
berarti kebalikan** dari yang biasa.

**Blocked by:** 20

**Dasar:**
```
EVIDENCED(SaveTreatyInDetail_Act@ekspor-2026-09 - DEDUCTION1=0, DEDUCTION2=0, SPREAD_RNM_SHARE_VALUE=0)
        EVIDENCED(DetailCalculationROL@ekspor-2026-09 - @if(Local.TotalLimit==0, 9989998, ...))
        DECIDED(ADR-0035, INV-42)
```

- [ ] `INV-42` terpasang dan menolak angka penanda kegagalan pada kolom uang dan persentase
- [ ] `INV-38` menolak kurs nol **dan** kurs yang dipaksakan satu
- [ ] keterangan kegagalan **tersimpan dan terbaca**, bukan hanya muncul sesaat di layar
- [ ] uji positif lulus: nol yang memang disepakati tetap tersimpan sebagai nol
- [ ] daftar angka penanda yang dilarang **menyebut `9989998` dengan namanya**, bukan hanya "angka besar"
- [ ] pembatasan lingkup tertulis: baris warisan dibawa apa adanya, dan `Uji AQ` yang menghitungnya
