---
status: tertahan
---

# 04: Keputusan tercatat pada jenjang yang benar, dan orang lain ditolak

*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Pemegang jenjang aktif memutus dengan komentarnya; siapa pun yang bukan
pemegang jenjang aktif **ditolak di server** dengan pencocokan identitas penuh. Penolakan
menutup sirkulasi dan menyisakan jenjang di bawahnya berkeadaan `TIDAK_SAMPAI` — bukan
"menolak" atas nama orang yang tidak pernah bertindak.

Operasi `CatatKeputusan`; jalur pencatatan keputusan; keadaan jenjang dan keadaan sirkulasi.

**Persyaratan:** `S-017`, `S-018`, `S-019`, `S-020`, `S-021`, `S-022`, `S-023`, `S-026`,
`S-027`, `S-028`, `S-029`, `S-030`, `S-068`.

**Tidak termasuk:** akibat pada klaim dan tulis-balik — tiket `05`. Versi usulan dan penanda
usulan — tiket `06`. Nomor akseptasi — tiket `08`.

**Jalur gagal:** pemanggil bukan pemegang jenjang aktif → `403` menyatakan giliran bukan
miliknya, **berbeda bentuk** dari `503` · keputusan kedua atas jenjang yang sudah memutus →
`409`, catatan yang ada tidak berubah · keputusan atas sirkulasi yang sudah selesai → `409` ·
nilai penanda di luar dua keadaan → `422` · kunci idempotensi sama dengan muatan berbeda →
`409` · gangguan di tengah → `503`, nol keputusan tersimpan.

**Uji:** P5-01, P5-02 (urutan giliran, **paritas** — dirancang lulus) · P5-02b (dua penentu
berselisih, **dirancang gagal**) · P5-14, P5-15 (peran dan jabatan dari roster) · P5-05
(penanda dua-keadaan) · BARU untuk `TIDAK_SAMPAI`, ketetapan keputusan, dan idempotensi.

**Menggantikan:** `KomiteRouter`·1–4 — penugasan lewat hitungan jenjang, **dibuang** (`K5-1`) ·
`KomiteRouter`·6, 6.1 — baris pertama yang belum memutuskan, **dipertahankan sebagai
paritas** · `KomiteRouter` — empat cabang mati, **dibuang** · `IsKomiteLoop` — gerbang loop,
**dibuang** · `DateApproval` dan `DateApprove` — dua kolom untuk satu fakta, dilebur jadi satu.

**Blocked by:**

- `02` — Sirkulasi usulan pembayaran lahir lengkap dengan jenjangnya


**Dasar:** DECIDED(`K5-1`, `K5-6`, `E-2`, `E-1`, keputusan beku no. 1, 2, 3, 4, ADR-0030).
EVIDENCED: `F-24` — nol `pyPrivilegeName` berisi pada 59 rule; GRILL-06 P6-3 — pemeriksaan
wewenang memasang pesan dan tidak menghentikan apa pun.

- [ ] Identitas dicocokkan dengan **kesetaraan penuh**; pengguna yang identitasnya mengandung
      pemegang sebagai potongan teks **ditolak**.
- [ ] Jenjang aktif adalah jenjang berderajat terendah yang belum memutuskan, dan **hitungan
      jenjang tidak pernah** menjadi syarat penugasan maupun penutupan.
- [ ] Sirkulasi berjalan punya **tepat satu** jenjang aktif — dihitung per sirkulasi. Ini yang
      mewujudkan invarian `I-7`.
- [ ] Penolakan menutup sirkulasi, dan seluruh jenjang berderajat lebih tinggi berkeadaan
      `TIDAK_SAMPAI` **tanpa komentar dan tanpa waktu milik orang lain**. Invarian `I-5`.
- [ ] Persetujuan jenjang terakhir menutup sirkulasi; penutupan lahir dari **keadaan
      keputusan**, bukan dari perbandingan dua bilangan. Invarian `I-6`.
- [ ] Keputusan tidak dapat disunting sesudah tercatat. Invarian `I-9`.
- [ ] Tabel jenjang memuat **tepat satu** kolom waktu keputusan — pencarian katalog skema
      menemukan tidak ada yang kedua. Invarian `I-10`.
- [ ] Peran dan jabatan disalin pada saat keputusan diambil; mutasi jabatan sesudahnya tidak
      mengubah riwayat.
- [ ] Galat gangguan dan penolakan wewenang **berbeda bentuk** di permukaan.
- [ ] Komentar, penanda bersyarat, dan penanda usulan yang sudah diisi **tidak hilang** ketika
      penyimpanan gagal.

**Ketidakpastian:** Tidak ada.
