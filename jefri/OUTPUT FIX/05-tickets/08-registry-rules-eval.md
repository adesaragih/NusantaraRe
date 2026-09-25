# 08: Registry predikat — satu seam untuk seluruh rule `When`

**What to build:** Seluruh gerbang keputusan sistem lama hidup di **satu registry** yang dapat
ditanyai dengan nama, dan setiap predikat menyebut rule Pega asalnya. Tidak ada salinan logika
gerbang yang tersebar dan bisa menyimpang.

⚠️ **Kondisi rule `When` ada di DUA tag, dan keduanya wajib dibaca.** Membaca hanya yang pertama
membuat rule tampak "tanpa kondisi" padahal kondisinya ada — **170 berkas (28,3 %)** menyembunyikan
kondisinya di tag kedua: 13 tagnya kosong, 157 berisi teks placeholder. Seluruh **601** berkas punya
tag kedua terisi, dan **tidak ada satu pun** rule yang kondisinya benar-benar tidak terbaca.

Klaim lama "lima rule `When` kondisinya kosong" berasal dari membaca satu tag saja, dan **keliru**.

Predikat individual **bukan** seam. Satu pintu masuk, diuji table-driven.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] Satu pintu masuk evaluasi bernama, melayani seluruh registry — bukan satu fungsi publik per predikat
- [ ] Kondisi dibaca dari **kedua** tag; rule yang kondisinya hanya ada di tag kedua tetap terbaca benar
- [ ] Teks placeholder pada tag pertama **tidak** diperlakukan sebagai kondisi
- [ ] **Nama predikat yang tidak dikenal → gagal keras**, bukan mengembalikan salah — salah ketik tidak boleh berubah menjadi gerbang yang selalu tertutup
- [ ] Setiap predikat menyebut rule Pega asalnya dalam komentar
- [ ] Uji table-driven mencakup contoh dari kedua kelompok: yang kondisinya terbaca di tag pertama, dan yang tersembunyi di tag kedua
