# 03: Tangga tiga jenjang dan empat cabang putusan — nilai kosong berarti DISETUJUI

**Status:** ready-for-agent
**Blocked by:** 02
**Menutup:** AC 1 · 2 · 3 · 4 · 5 · 6 · 7 · 8 · 9 · 10 · 84 *(11 AC)* — US 7–14

## Hasil & nilai pengguna

Hari ini persetujuan realisasi treaty melewati **lima posisi**, dan dua di antaranya —
Ketua Kelompok dan Direktur — ⚠️ sudah lama tidak dipakai. ⛔ Lebih rawan lagi: ada **dua aturan
berbeda bernama sama** untuk memutuskan "sudah disetujui atau belum", dan keduanya **berlawanan**
pada nilai kosong.

Sesudah tiket ini, persetujuan berjalan lewat **tiga jenjang** yang jelas, dan ⭐ setiap penolakan
punya akibat yang pasti: **admin menolak ⇒ berkas selesai sebagai ditolak; atasan menolak ⇒ berkas
kembali ke admin.**

## Area codebase

- Lapisan service: mesin tahap persetujuan
- Lapisan service: penilaian penanda persetujuan

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| ⭐ Aturan **yang hidup** | `DecisionTable\isApproved.xml` — `= 0` ⇒ `No`, bawaan ⇒ `YES`; kolomnya bertipe **teks** |
| ⛔ Aturan yang **tidak dipakai** | `When\isApproved.xml` — menguji `= 1`; **nol** kotak Decision menyambung ke sana |
| Kotak putusan | `Flow\InputRealizationTreatyIn.xml` — **enam** Decision, seluruhnya menyambung ke `Rule-Declare-DecisionTable` |
| Perpindahan tahap | `DataTransform\InboxPolicyTreatyIn_postDT.xml` · `DeptHeadTreatyIn_UW_postDT.xml` |

## ADR terkait

- **ADR-0007** — jejak audit setiap transisi dan setiap jalur balik

## Acceptance criteria

- [ ] **AC 1** — `0` diperlakukan **ditolak**
- [ ] **AC 2** — nilai apa pun selain `0`, ⭐ **termasuk kosong**, diperlakukan **disetujui**
- [ ] **AC 3** — perbandingan dilakukan sebagai **teks**
- [ ] **AC 4** — aturan berasal dari **tabel keputusan**, bukan dari aturan bernama sama yang menguji `= 1`
- [ ] **AC 5** — admin menolak ⇒ berkas **diselesaikan sebagai ditolak**
- [ ] **AC 6** — atasan menolak ⇒ berkas **kembali ke admin**
- [ ] **AC 7** — admin menyetujui ⇒ naik ke jenjang kedua
- [ ] **AC 8** — jenjang kedua menyetujui ⇒ naik ke jenjang ketiga
- [ ] **AC 9** — jenjang ketiga menyetujui ⇒ realisasi **selesai**; ⛔ tidak ada jenjang keempat
- [ ] **AC 10** — dua posisi yang dibuang **tidak ada**; berkas tidak pernah dirutekan ke sana
- [ ] **AC 84** — nilai kosong pada penanda **tidak** menghentikan alur

## Perintah verifikasi

1. Ajukan berkas yang penandanya **belum pernah diisi** — ⭐ ia diperlakukan **disetujui**, bukan
   ditolak, dan **tidak** melempar galat.
2. Tolak sebagai admin — ⭐ berkas **selesai**.
3. Tolak sebagai atasan — ⭐ berkas **kembali ke admin**, bukan selesai.
4. Setujui tiga kali berturut-turut — ⭐ realisasi **selesai** di jenjang ketiga.

## Catatan

⚠️ **Perbedaan dua aturan itu nyata, bukan gaya penulisan.** Aturan yang tidak dipakai
memperlakukan nilai kosong sebagai **tidak disetujui**; yang hidup memperlakukannya **disetujui**.
⛔ Membangun dari yang salah **membalik perilaku** pada setiap berkas yang penandanya belum pernah
diisi.
