# 03: Tangga tiga jenjang dan empat cabang putusan — nilai kosong berarti DISETUJUI

**Status:** selesai *(putaran 2, 03-10-2026: AC 84 ditutup lewat K6, AC 8 dikukuhkan K2; implementasi
2026-10-03 semula: sebagian; semula: ready-for-agent)*
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

- [x] **AC 1** — `0` diperlakukan **ditolak**
- [x] **AC 2** — nilai apa pun selain `0`, ⭐ **termasuk kosong**, diperlakukan **disetujui**
- [x] **AC 3** — perbandingan dilakukan sebagai **teks**
- [x] **AC 4** — aturan berasal dari **tabel keputusan**, bukan dari aturan bernama sama yang menguji `= 1`
- [x] **AC 5** — admin menolak ⇒ berkas **diselesaikan sebagai ditolak**
- [x] **AC 6** — atasan menolak ⇒ berkas **kembali ke admin**
- [x] **AC 7** — admin menyetujui ⇒ naik ke jenjang kedua
- [x] **AC 8** — jenjang kedua menyetujui ⇒ naik ke jenjang ketiga
- [x] **AC 9** — jenjang ketiga menyetujui ⇒ realisasi **selesai**; ⛔ tidak ada jenjang keempat
- [x] **AC 10** — dua posisi yang dibuang **tidak ada**; berkas tidak pernah dirutekan ke sana
- [x] **AC 84** — nilai kosong pada penanda **tidak** menghentikan alur *di tabel keputusan*; di layar
      submit tanpa Approval ditolak 422 (K6 — RALAT putaran 2 di bawah)

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

## ⛔ Pertentangan WO lawan XML — dicatat 2026-10-03

**AC 8** *"jenjang kedua menyetujui ⇒ naik ke jenjang ketiga"* `[keputusan work owner]` bertentangan
dengan XML: `Decision13` (`ToTREATYDEPTHEAD`, `pyWorkPage.LetterNo` diisi `CekLimitTreatyAcc_Act` bila
|TotalPremium| × kurs > 200.000.000) dan `Decision8` (`NopolisEmpty`) membiarkan Sec Head
**menyelesaikan** berkas bernomor di bawah batas. Aturan prompt implementasi: ikuti work owner, catat.
⇒ Sec Head menyetujui **selalu** naik ke Dept Head; `CekLimitTreatyAcc_Act` dan tombol Generate Sec Head
tidak dibangun. ~~Mohon konfirmasi work owner bila batas 200 juta tetap dikehendaki.~~ ⇒ **Dijawab K2
(PROMPT-NB-TREATY-IN-PUTARAN-2 bab 2, 03-10-2026): WO lama berlaku** — Sec Head menyetujui → selalu
Dept Head; batas 200 juta XML **tidak** dibangun (alasan (b) bab 4: keputusan work owner K2).
Pertentangannya tetap tercatat di sini: XML `Flow\InputRealizationTreatyIn.xml` `Decision13`
(`When\ToTREATYDEPTHEAD`, `pyWorkPage.LetterNo` diisi `Activity\CekLimitTreatyAcc_Act.xml` bila
|TotalPremium| × kurs > 200.000.000) dan `Decision8` (`NopolisEmpty`) membiarkan Sec Head
menyelesaikan berkas bernomor di bawah batas. Uji: `tangga_test.go`
TestTanggaSecHeadMenyetujuiSelaluNaikKeDeptHead, `handlers/alur_test.go` TestTanggaPenuhDanNomorPolisSekali.

**AC 84** — pada tingkat tabel keputusan, IsApproved kosong = disetujui (diuji). Di layar, `ListSuggest`
mewajibkan Approval dan tombol Submit hanya tampil untuk IsApproved 1/0 (XML), sehingga submit dengan
Approval kosong dijawab validasi 422 — ~~dicatat sebagai sebagian.~~ ⇒ **Dijawab K6 (03-10-2026): Ya,
submit tanpa Approval ditolak 422, mengikuti layar XML.**

## ⛔ RALAT putaran 2 (P4, 03-10-2026) — AC 84

Bunyi lama (spec AC 84): *"Nilai kosong pada penanda persetujuan **tidak** menghentikan alur. Test yang
menemukan galat **gagal**."* Bunyi baru (K6): **"Pada tabel keputusan `DecisionTable\isApproved`,
`IsApproved` kosong diperlakukan *disetujui* tanpa galat (AC 2); di layar, submit tanpa Approval
DITOLAK 422 di ketiga jenjang."** Bukti XML layar: `Section\ListSuggest` mewajibkan `.IsApproved`;
tombol Submit `DetailPolicyTreatyIn`/`DetailDeptHeadTreatyIn_UW` tampil hanya untuk IsApproved 1/0
(`models.TombolUntuk`). Uji: `tangga_test.go` (kosong = disetujui), `handlers/alur_test.go`
TestMedanWajibMenahanKirimDanSimpan (admin), `handlers/logika_test.go`
TestSubmitTanpaApprovalDitolakDiSetiapJenjang (admin, Sec Head, Dept Head: 422, tanpa riwayat, tanpa
perpindahan). ⇒ AC 84 ✅.
