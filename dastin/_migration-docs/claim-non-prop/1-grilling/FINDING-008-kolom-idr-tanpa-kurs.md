# FINDING-008 (usulan) — Kolom IDR dapat terisi tanpa kurs pernah dipakai

**Status: USULAN.** Ia naik menjadi temuan tetap hanya bila ramalan di bagian 4 terukur benar. Sampai itu, ia anomali E14 yang diusulkan naik.

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, empat lapisan. Sapuan S10, dijalankan 18 September 2026 — `SAPUAN-S10-S16.md` §3.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut.

Dokumen ini **tidak** menilai siapa pun dan tidak menyatakan ada kerugian. Ia menyatakan satu hal yang terbaca di kode: sebuah kolom dapat memuat angka yang tidak sesuai labelnya, dan tidak ada apa pun di baris itu yang memperlihatkannya.

## 1. Apa yang terbaca

Satu kolom, `ClaimAmountIDR` pada `SpreadingRisk`, diisi oleh **dua rule dengan perlakuan yang berbeda**.

**Penulis pertama — menguji mata uang, lalu mengalikan kurs.** `Activity\CountLossAllocation_act.xml` baris **9585**:

```
pyWorkPage.ClaimData.SpreadingRisk(<LAST>).ClaimAmountIDR
  <- @if(Local.Currency=="IDR",
         pyWorkPage.ClaimData.SpreadingRisk(<LAST>).ClaimEstimation,
         pyWorkPage.ClaimData.SpreadingRisk(<LAST>).ClaimEstimation * Local.Kurs)
```

**Penulis kedua — menyalin apa adanya.** `Activity\SetActualPremium_ACT.xml` baris **681**:

```
pyWorkPage.ClaimData.SpreadingRisk(<LAST>).ClaimAmountIDR
  <- pyWorkPage.ClaimData.SpreadingRisk(<LAST>).TotalClaim
```

Penulis kedua **tidak menguji mata uang** dan **tidak mengalikan kurs**. Baris yang sama juga menerima `ClaimAmountAdjust`, `ClaimEstimation`, dan `TotalSpread` dari rule itu (660, 702, 723), seluruhnya berupa penyalinan langsung.

## 2. Mengapa ini bukan sekadar perbedaan gaya

`TotalClaim` tidak dijamin rupiah. Ia diisi mesin alokasi di baris 9279 sebagai turunan `Local.ClaimValue`, yang mengikuti mata uang klaimnya. Bila klaim berjalan dalam mata uang selain rupiah, penulis kedua menempatkan nilai mata uang asing di kolom yang bernama rupiah.

Dan tidak ada penanda yang memperlihatkannya. Baris itu memang punya `KursIDR`, tetapi penulis kedua **tidak menyentuhnya** — kursnya tetap berisi apa pun yang ditinggalkan penulis sebelumnya, atau kosong. Pembaca hilir tidak punya cara membedakan baris yang dikonversi dari yang tidak.

## 3. Apa yang disentuhnya

| Dokumen | Kaitannya |
|---|---|
| **FINDING-006** | Kurs yang tidak ditemukan bernilai satu. Keduanya cacat pada jalur yang sama: di sana kurs **gagal ditemukan** lalu diperlakukan `1`; di sini kurs **tidak pernah dicari**. Akibat pada angkanya sama — nilai asing masuk kolom rupiah tanpa dikonversi |
| **ADR-0014** | Menetapkan perlakuan kurs. Perlakuan yang ditetapkan mengandaikan ada **satu** jalur konversi; di sini ada dua, dan salah satunya melewatinya |
| **ADR-0029** (accepted 19 Sep 2026) | Menguatkannya dari arah yang berbeda dari S4. S4 menyatakan sebagian nilai uang tidak punya padanan rupiah; temuan ini menyatakan padanan rupiah yang **ada** pun belum tentu hasil konversi |
| **E14** | Butir anomali asalnya |
| **AK-2a / AK-2b** | Bila terukur benar, nilai lamanya mengikuti keduanya: dimigrasi apa adanya, dan selisihnya didaftar |

## 4. Ramalan yang dapat diuji

Bila temuan ini benar, data produksi memuat baris dengan sifat berikut:

> Ada baris `SpreadingRisk` yang **`ClaimAmountIDR`-nya sama persis dengan nilai mata uang aslinya**, sementara `Currency` baris itu **bukan** `IDR`.

Kesamaan persis itu yang menentukan. Bila nilainya dikonversi, `ClaimAmountIDR` akan berbeda dari nilai aslinya sebesar kursnya; hanya baris yang **tidak** dikonversi yang menyimpan angka yang identik.

**Yang menggugurkan temuan ini**: nol baris seperti itu. Bila setiap baris bermata-uang-asing punya `ClaimAmountIDR` yang berbeda dari nilai aslinya, berarti `SetActualPremium_ACT` tidak pernah dijalankan pada baris non-rupiah, dan cacatnya teoretis belaka.

**Yang memperkuatnya**: baris seperti itu ada **dan** `KursIDR`-nya tidak konsisten dengan selisihnya.

Diukur oleh **REQ-034**.

## 5. Yang tidak dinyatakan dokumen ini

Berapa banyak baris yang terdampak — belum terukur. Apakah angkanya pernah dipakai dalam jurnal — pertanyaan akuntansi, bukan pertanyaan kode. Siapa yang menulis rule mana dan kapan — di luar lingkup; dokumen ini tentang perilaku kode, bukan tentang orang.

**Pembatasan rancangan dicabut 19 September 2026.** Kalimat lama berbunyi *"sampai REQ-034 kembali, tidak boleh ada rancangan yang mengandaikan temuan ini benar maupun gugur"*. Itu tidak lagi berlaku, dan sebabnya bukan karena jawabannya datang.

**Rancangannya sudah benar di kedua keadaan.** ADR-0029 mewajibkan setiap nilai uang menyimpan nilai asli, kode mata uang, nilai rupiah, kurs yang dipakai, dan asal-usul kurs itu — di baris yang sama. Dengan bentuk itu **nilai rupiah tidak dapat lahir tanpa kurs yang menyertainya**, sehingga keadaan yang dijelaskan dokumen ini tidak dapat terjadi di sistem baru, entah ia terbukti pernah terjadi di sistem lama atau tidak.

Yang tersisa untuk REQ-034 karena itu **bukan pertanyaan rancangan melainkan pertanyaan data lama**: berapa baris yang perlu diperiksa saat migrasi, dan bagaimana baris seperti itu diperlakukan — AK-2a dan AK-2b. Status dokumen ini tetap **USULAN**: ia naik menjadi temuan tetap bila ramalan bagian 4 terukur benar, dan gugur bila nol.

Status E14 di `_selesai/OPEN-QUESTIONS.md`: **ditutup 19 September 2026**, naik ke dokumen ini.
