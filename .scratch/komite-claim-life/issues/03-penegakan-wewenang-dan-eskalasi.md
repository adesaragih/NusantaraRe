# 03: Penegakan wewenang per `KomiteID` + eskalasi naik satu tingkat

**Status:** ready-for-agent

**Blocked by:** 02 (mesin tangga) — gerbang perlu tindakan nyata untuk dijaga

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin keputusan pada sebuah tingkat **hanya** dapat disimpan oleh
anggota komite yang memang ditunjuk tingkat itu — dan penolakannya terjadi di lapisan layanan,
sehingga tidak dapat dilewati dengan memanggil API langsung.

Sebagai **admin komite**, saya ingin memindahkan kasus **naik satu tingkat** bila anggota tingkat
berjalan berhalangan, supaya kasus tidak macet menunggu satu orang.
*(User story 12, 16, 30 di spec)*

⚠️ **Penyimpangan sadar.** Sistem lama tidak menegakkan apa pun.

## Area codebase

`internal/services` (pemeriksaan wewenang sebelum menyimpan keputusan; aksi eskalasi),
`internal/handlers` (identitas pemanggil), `internal/repository` (pencatatan eskalasi),
`frontend/` (kontrol eskalasi untuk admin; kontrol keputusan hanya bagi yang berwenang).

## Rule Pega sumber

| Rule | Identitas | Keadaan sekarang |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomiteRouter.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` hanya **menempatkan** tugas: `param.AssignTo = .KomiteID`, `pyImplementation = WorkList` |
| `Komite Claim Life/Section/ShowTransfer.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `SHOWTRANSFER` / `RULE-HTML-SECTION` | `[terverifikasi]` keputusan masuk lewat dropdown wajib — **tanpa** pemeriksaan pemilik |

`[terverifikasi]` **Tidak ada penegakan di sistem lama.** Sapuan 47 berkas modul tidak menemukan satu
pun pemeriksaan bahwa penyimpan keputusan adalah pemilik `KomiteID` tingkat berjalan; dan
`AcceptStatus` **tidak ditulis rule mana pun** (nol `<PropertiesName>…AcceptStatus</PropertiesName>`).
Di Pega, routing `WorkList` hanya menaruh kasus di antrean — **penempatan, bukan penegakan**.

`[terverifikasi]` Komite Claim Life **tidak memuat identitas orang ter-hardcode** (0 berkas) —
berbeda dari Komite Claim FacIn (4) dan Komite Claim Prop (3). Wewenangnya memang berbasis data.

## ADR terkait

**ADR-0014** (penegakan per `KomiteID` + pengecualian eskalasi naik satu tingkat), **ADR-0002**
(peran ditegakkan di lapisan layanan), **ADR-0007** (eskalasi dan perubahan roster masuk jejak
audit).

## Acceptance criteria

- [ ] Pengguna yang **bukan** pemilik `KomiteList(KomiteCount).KomiteID` **ditolak** saat menyimpan
      keputusan, meskipun ia dapat membuka kasusnya. *(AC 8 spec)*
- [ ] Penolakan terjadi **di lapisan layanan**, dan tetap terjadi meskipun kontrol UI ditampilkan.
      *(AC 9 spec)*
- [ ] Admin dapat memindahkan kasus **naik satu tingkat**; eskalasi **turun** ditolak.
      *(AC 11 spec)*
- [ ] Eskalasi tercatat: **siapa** memindahkan, **kapan**, dari tingkat mana ke tingkat mana.
      *(AC 12 spec)*
- [ ] Pemutus di tingkat **yang sama** yang bukan pemilik `KomiteID` tetap ditolak — eskalasi bukan
      pintu belakang.
- [ ] Perubahan roster tercatat, karena ia mengubah **siapa yang berwenang**.
- [ ] Tidak ada nama orang ter-hardcode di lapisan mana pun.

## Catatan

⚠️ `[keputusan work owner]` Eskalasi **memperpendek** tangga: tingkat yang dilewati tidak pernah
memberi keputusan, sehingga entri `KomiteAproval` untuk tingkat itu kosong. `KomiteCount` ikut naik.
Bentuk pencatatannya adalah keputusan implementasi.

`[terbuka]` **OQ-007 / OQ-021** — model RBAC lintas konteks belum ditetapkan. **Asumsi tiket ini:
satu `KomiteID` memetakan ke satu identitas akun yang dapat diautentikasi.** Bila pemetaan ternyata
banyak-ke-banyak, aturan penegakan perlu ditinjau ulang.

`[terverifikasi]` Identity & Access ditandai **ABSENT** dari korpus (`discovery/context-map.md`) —
nol rule otorisasi. Ia dibangun dari nol.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
