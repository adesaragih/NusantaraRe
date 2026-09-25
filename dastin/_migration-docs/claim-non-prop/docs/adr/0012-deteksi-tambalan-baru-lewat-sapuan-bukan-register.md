---
status: accepted
label: DECIDED
---

# Tambalan baru dideteksi lewat sapuan ulang, bukan lewat register

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Tambalan per-case dan per-identitas dibekukan sejak tanggal kesepakatan. **Selama pekerjaan migrasi ini berjalan, pembekuan berlaku mutlak — tidak ada pengecualian.** Yang menjadi **alat deteksi** bukan register, melainkan **sapuan ulang otomatis** atas setiap ekspor XML baru, dibandingkan terhadap `BLUEPRINT.md` bagian 7.5. Register pengecualian tetap dibuat, tetapi statusnya pelengkap.

## Consequences

Pembekuan total akan dilanggar diam-diam; yang tercatat lebih berguna daripada yang dilarang. Tetapi register yang diisi manusia punya cacat yang sama dengan `CNPStatusCase` — bisa tidak pernah diisi, dan ketiadaannya tidak menghasilkan sinyal apa pun. **Kebasian analisis tidak boleh bergantung pada kepatuhan orang.**

Sapuan ulang mendeteksi tambalan baru apakah registernya diisi atau tidak. Selisih antara hasil sapuan dan bagian 7.5 adalah daftar tambalan yang masuk sesudah ekspor terakhir.

Konsekuensi kedua: **setiap artefak membawa stempel asal** — tanggal ekspor XML yang menjadi dasarnya dan jumlah berkas yang disapu. Tanpa itu, tidak ada yang dapat menilai dokumen ini berlaku untuk keadaan kapan. Stempel sudah dipasang di seluruh artefak per 18 September 2026 (279 berkas, ekspor 2026-09-08/09).

## Pengecualian dihapuskan — 19 September 2026 (menutup A13)

Versi pertama memberi pengecualian untuk *"tambalan yang menahan pembayaran atau menghentikan operasi"*, lalu menyisakan pertanyaan terbuka: **siapa berwenang menyetujuinya** (`_selesai/OPEN-QUESTIONS.md` A13). Pertanyaan itu ditutup dengan **menghapus pengecualiannya**, bukan dengan menjawabnya.

Selama pekerjaan ini berjalan, pembekuan berlaku mutlak. Tidak ada pengecualian, jadi tidak ada yang perlu berwenang menyetujuinya, dan tidak ada yang perlu ditanyakan ke luar.

**Register pengecualian tetap ada, dan tetap kosong.** Ia tidak dibongkar: bila kelak ada yang mengisinya, saat itulah wewenang penyetujunya ditetapkan — oleh orang yang mengisinya, pada saat ia mengisinya. Register yang kosong adalah pernyataan yang dapat diperiksa; ketiadaan register bukan.

Yang **tidak** berubah: alat deteksinya tetap sapuan ulang, bukan register. Pembekuan mutlak justru mempertajam alasan itu — semakin tegas larangannya, semakin kecil kemungkinan pelanggarnya mencatat sendiri.

## Tambahan 18 September 2026 — sapuan mencakup tiga lapisan, bukan satu

Versi pertama menyebut "sapuan ulang atas setiap ekspor XML baru". **Itu terlalu sempit.** Tiga jalur yang menggerakkan sistem tidak berjejak di XML sama sekali (lihat `BLUEPRINT.md` bagian 14), termasuk tanggal tertanam di dalam `PROC_GENERATE_SEQUENCE_NUMBER` yang merupakan kelas tambalan yang sama dengan 29 langkah di bagian 7.

Sapuan berkala mencakup:

| Lapisan | Yang dicari |
|---|---|
| **Ekspor XML** | literal `CLMNP-…`, `IDMaster`, identitas operator/orang, pola `@contains` atas teks bebas |
| **DDL** | tanggal dan periode tertanam, nilai mati di dalam procedure, kolom baru, constraint yang hilang |
| **Daftar database link** | ketergantungan lintas basis data yang baru muncul |

Lapisan ketiga ditambahkan setelah `V_MST_USER_TEKNIS` terbaca. Sebelum itu tidak ada yang tahu perlu menyapunya — yang justru menjadi alasan terkuat mengapa daftar lapisan ini sendiri harus ditinjau setiap kali sumber baru masuk.
