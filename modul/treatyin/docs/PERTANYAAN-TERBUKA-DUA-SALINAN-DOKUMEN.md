# Pertanyaan terbuka — dua dokumen bersalinan dua, dan keduanya BUKAN duplikat

**Diajukan kepada pemilik spec 5 Oktober 2026.**

⛔ **Nol berkas dihapus, nol isi diganti.** Ronde ini diminta *"nyatakan mana yang induk, lalu
ganti isi yang satunya dengan satu baris penunjuk"* — dan pengukuran membatalkan premisnya untuk
berkas pertama. Keduanya diserahkan sebagai pertanyaan.

---

## 1 · `KOREKSI-ERD-VERSUS-POOLDATA.md` — dua dokumen BERBEDA, satu nama

| Salinan | Baris | Isi |
| --- | ---: | --- |
| `D:\NUSANTARA RE APP\treaty-in\4-erd-dan-tabel-datar\` | **810** | §1–§28, ditulis bertahap sejak 3 Oktober 2026 |
| `D:\XML_NURE\_migration-docs\treaty-in\4-erd-dan-tabel-datar\` | **220** | §0–§7 |

⛔ **Yang 220 baris BUKAN versi lama dari yang 810.** Kedelapan babnya diadu satu per satu dengan
salinan repo, dan **nol** di antaranya ada di sana:

```
## 0 · Satu kalimat yang perlu dikoreksi lebih dulu        TIDAK ADA di repo
## 1 · Tabel yang SUNGGUH ADA di POOLDATA                  TIDAK ADA di repo
## 2 · Isi tab — ada di JSON, bukan di tabel               TIDAK ADA di repo
## 3 · Pola struktur — SATU, bukan dua                     TIDAK ADA di repo
## 4 · `TreatyLeader` — tiga keadaan, bukan dua            TIDAK ADA di repo
## 5 · Yang masih belum terukur                            TIDAK ADA di repo
## 6 · Yang berubah bagi rekonsiliasi sebelumnya           TIDAK ADA di repo
## 7 · Cara memeriksa ulang berkas ini                     TIDAK ADA di repo
```

Keduanya **menjawab pertanyaan yang sama dengan penomoran yang berbeda** — beberapa temuannya
bahkan sejalan (`TreatyLeader` tiga keadaan; tabel yang sungguh ada di POOLDATA). Tetapi tidak
satu pun kalimatnya berbagi, jadi mengganti isi salah satunya dengan penunjuk **membuang isi yang
tidak ada di mana pun lagi**.

**Yang ditanyakan:**

1. Apakah keduanya **digabung** menjadi satu dokumen — dan bila ya, siapa yang memutuskan
   penomoran babnya, sebab §0–§7 dan §1–§28 bertabrakan?
2. Atau keduanya **dibedakan namanya**, misalnya salah satunya menjadi
   `KOREKSI-ERD-VERSUS-POOLDATA-SAPUAN-AWAL.md`?
3. Mana yang menjadi tempat koreksi berikutnya ditulis? *(Ronde-ronde terakhir diarahkan ke
   salinan repo yang 810 baris, dan §1–§28 di sana adalah hasilnya.)*

---

## 2 · `SPEC-MODEL-DATA.md` — yang LEBIH BARU ada di LUAR repo

| Salinan | Baris | Diubah | Akhir baris |
| --- | ---: | --- | --- |
| `D:\NUSANTARA RE APP\treaty-in\` | 2.604 | **24 September 2026** | CRLF |
| `D:\XML_NURE\_migration-docs\treaty-in\` | **2.791** | **2 Oktober 2026** | LF |

⚠️ **Selisihnya jauh lebih kecil daripada yang terlihat.** `diff` mentah melaporkan 5.395 baris
berbeda — itu akibat CRLF lawan LF, bukan isi. Sesudah akhir barisnya disamakan, yang benar-benar
berbeda **193 baris** dari ~2.700; keduanya ~93% sama.

⛔ **Yang perlu diperhatikan bukan besarnya selisih, melainkan arahnya:** salinan yang **lebih
baru dan lebih panjang ada DI LUAR repo**. Siapa pun yang membaca salinan repo membaca dokumen
yang delapan hari lebih tua.

**Yang ditanyakan:**

1. Mana yang **induk**? *(Ronde ini tidak memilih: menyunting `SPEC-MODEL-DATA.md` dilarang.)*
2. Bila yang di `XML_NURE` induknya, apakah salinan repo **disegarkan** darinya — dan siapa yang
   menjalankannya?
3. Apakah akhir barisnya diseragamkan sekalian, supaya `diff` berikutnya menunjukkan isi dan
   bukan CRLF?

---

## 3 · ⚠️ Tagihan yang menunggu jawaban nomor 2

[`KEPUTUSAN-PENYELARASAN-REPO.md` §19.3](KEPUTUSAN-PENYELARASAN-REPO.md) mencatat bahwa
`SPEC-MODEL-DATA.md` **masih menyatakan model datanya SATU**, sementara sejak §19 dua tabel
(`NILAI_SELISIH`, `NILAI_SEBELUM_PRO_RATE`) berdiri di luar `treatyin`.

Pernyataan itu perlu diperbarui — **dan sekarang ia perlu diperbarui di salinan yang mana, dan
itu belum dapat dijawab sebelum nomor 2 di atas dijawab.** Satu tagihan menunggu satu keputusan,
dan keduanya milik pemilik spec.
