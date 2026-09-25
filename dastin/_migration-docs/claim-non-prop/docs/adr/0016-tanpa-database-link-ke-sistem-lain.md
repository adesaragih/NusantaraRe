---
status: accepted
label: DECIDED
---

# Tanpa database link ke sistem lain — dan tanpa data pegawai sama sekali

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.
> Berlaku hanya untuk keadaan sistem pada kedua ekspor itu.

Sistem baru **tidak mewarisi ketergantungan lintas basis data**.

**Modul Claim Non Prop tidak mengambil data HRD sama sekali** — tidak lewat database link, tidak lewat API, tidak lewat salinan tersinkron. Pelaku sebuah tindakan disimpan sebagai **potret**: pengenal operator dan nama **sebagaimana tercatat pada saat tindakan terjadi**. Tidak ada pencarian ke HRD, tidak ada penyegaran, tidak ada salinan yang disinkronkan.

Sebabnya bukan sekadar lingkup. **Pelaku suatu tindakan adalah fakta yang terjadi pada satu saat, dan fakta itu tidak berubah ketika orangnya pindah bagian atau berhenti.** Menyimpannya sebagai rujukan hidup ke HRD membuat catatan lama ikut berubah ketika sumbernya berubah — itu merusak audit, bukan memperkayanya.

Sistem lama memakai `hrdasm.v_hrd_mst@asmd.sinarmas.co.id` — database link langsung ke basis data sistem HRD, di dalam view `V_MST_USER_TEKNIS`.

## Consequences

Tiga alasan, tidak satu pun soal selera:

1. Database link **melewati kontrol akses aplikasi HRD sepenuhnya**.
2. Ia **mengikat kita pada struktur tabel internal** sistem lain. Bila mereka mengubah kolom, sistem kita rusak tanpa ada yang memberi tahu.
3. Bila tautan putus, view mengembalikan **kosong** — daftar pengguna menghilang tanpa pesan galat. Itu pola yang sama persis dengan `RETURN 1` pada kurs (`FINDING-006` bagian 3): kegagalan yang menyamar menjadi hasil yang sah.

Bila HRD tidak dapat menyediakan API maupun mekanisme sinkronisasi, itu percakapan yang dibawa ke luar. Rancangannya **tidak menunggu** jawaban itu, dan tidak dibuat seolah database link adalah salah satu pilihan.

Apakah tautan yang ada sekarang resmi disepakati **tidak lagi ditanyakan**: ia tautan milik sistem lama, dan sistem baru tidak memerlukan data di ujungnya.

## Penyelesaian K1 — 19 September 2026, sebagai **K1a**

`_selesai/OPEN-QUESTIONS.md` K1 mencatat bahwa ADR ini melarang sesuatu yang sudah berjalan produksi: `V_MST_USER_TEKNIS` memakai `hrdasm.v_hrd_mst@asmd.sinarmas.co.id`, satu-satunya pemakaian `@<host>` di keempat puluh sembilan berkas DDL (D20). Dua cabang ditulis sejajar: **K1a** (ADR tetap, yang lama dibiarkan) dan **K1b** (ADR direvisi dengan pengecualian bernama).

**Dipilih K1a, dan alasannya menghapus biaya yang semula melekat padanya.**

K1a semula mahal karena ia menuntut *"tetapkan cara lain memenuhi kebutuhan data HRD"* — pekerjaan baru yang belum dilingkupi ADR mana pun. Biaya itu **tidak ada**, karena kebutuhan itu tidak ada: dengan aturan potret di atas, modul ini tidak pernah memerlukan data HRD. Larangan ADR ini karena itu **tidak pernah bertabrakan dengan kebutuhan kita**.

| | Keadaan |
|---|---|
| `V_MST_USER_TEKNIS` | **tidak dimiliki, tidak dimigrasi.** View HRD, bukan entitas klaim — ADR-0026 (batas kepemilikan mengikuti nama class) |
| View lama beserta link-nya | dibiarkan hidup di sistem lama; bukan milik kita, dan tidak ada yang perlu kita gantikan |
| ADR-0023 dan ADR-0028 | **tidak ikut ditinjau.** Peninjauan itu konsekuensi K1b, dan K1b tidak dipilih |
| D20 | tetap tercatat sebagai fakta; kedudukannya kini **catatan atas sistem lama**, bukan pekerjaan |

**Rujukan ADR-0024 dicabut.** Versi pertama bagian atas ADR ini menyebut *"ADR-0024 tentang pemisahan pengguna, jabatan, dan keanggotaan komite"*. ADR-0024 adalah **kunci alami akseptasi**; rujukan itu salah sejak ditulis, dan kalimat yang memuatnya sudah hilang bersama janji HRD.
