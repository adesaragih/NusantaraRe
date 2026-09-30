---
status: accepted
tanggal: 2026-09-15
sumber: grilling Komite Claim Life Ronde 1 + Ronde 2 (`.scratch/komite-claim-life/`), keputusan work owner
---

# Keputusan komite hanya boleh diambil pemilik `KomiteID` pada tingkat berjalan

Pada setiap tingkat tangga persetujuan, **hanya pemilik
`KomiteList(KomiteCount).KomiteID`** yang boleh menyimpan keputusan. Pengguna lain **ditolak di
lapisan layanan**, meskipun ia dapat membuka kasusnya.

Ini **penyimpangan sadar**: Pega hanya *menempatkan* tugas, ia tidak *menegakkan* siapa yang
memutuskan.

## Keadaan sekarang `[terverifikasi]`

Penempatan tugas dilakukan router:

`Komite Claim Life/Activity/KomiteRouter.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY`, 26.387 byte):

```
Property-Set   param.AssignTo = .KomiteID        (baris ~294)
  precondition .KomiteAproval == 0               (baris ~382)
```

Flow memakai `<pyImplementation>WorkList` dengan `<pyRouteTo>Custom`. Artinya kasus **muncul di
antrean** pemilik `KomiteID` — penempatan.

**Yang tidak ada:** sapuan 47 berkas modul tidak menemukan satu pun pemeriksaan bahwa pengguna yang
menyimpan keputusan adalah pemilik `KomiteID` tingkat berjalan. Keputusan masuk lewat dropdown wajib
di `Komite Claim Life/Section/ShowTransfer.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `SHOWTRANSFER` / `RULE-HTML-SECTION`) baris 32607 — dan
`AcceptStatus` **tidak ditulis rule mana pun** (nol `<PropertiesName>…AcceptStatus</PropertiesName>`).

Jadi di sistem lama, **siapa pun yang dapat membuka layar dapat menyimpan keputusan** untuk tingkat
itu. Tidak ditemukan bukti bahwa itu disengaja.

## Considered Options

- **Tegakkan di lapisan layanan: hanya pemilik `KomiteID` tingkat berjalan** — dipilih
- Paritas (mengandalkan penempatan worklist saja) — ditolak: memindahkan lubang wewenang ke sistem
  baru, dan di Go tidak ada "worklist" yang secara kebetulan membatasi — tanpa penegakan eksplisit
  endpoint terbuka bagi siapa pun yang terautentikasi
- Tegakkan di UI saja — ditolak: sama dengan tidak menegakkan; API tetap terbuka

## Consequences

- `internal/services` memeriksa identitas pemanggil terhadap `KomiteList(KomiteCount).KomiteID`
  **sebelum** menyimpan keputusan; gagal → tolak, bukan abaikan diam-diam.
- **Roster menjadi sumber wewenang**, bukan sekadar daftar penerima email. Menambah/menonaktifkan
  baris `EMAILKOMITE` mengubah **siapa yang boleh memutuskan** — perubahannya layak masuk jejak
  audit (**ADR-0007**).
- Penegakan bergantung pada pemetaan **`KomiteID` → identitas akun**. Bentuk pemetaan itu belum
  ditetapkan; ia bagian dari Identity & Access yang `[terverifikasi]` **ABSENT** dari korpus
  (`discovery/context-map.md`).
- Sejalan semangat **ADR-0002** (peran ditegakkan di lapisan layanan, bukan di visibilitas layar)
  dan **ADR-0012** (wewenang eksplisit, bukan efek samping UI).
- Konsekuensi operasional: bila pemilik `KomiteID` berhalangan, kasus akan macet — **kecuali** lewat
  jalur eskalasi manual yang ditetapkan di §"Pengecualian sah" di bawah. Di sistem lama, orang lain
  bisa menyelesaikannya tanpa jejak; kini perpindahannya eksplisit dan terekam.

## Pengecualian sah: eskalasi manual naik satu tingkat (2026-09-15)

`[keputusan work owner]` Penegakan per `KomiteID` **tetap berlaku**, dengan **satu** pengecualian
yang ditetapkan eksplisit:

> **Bila anggota komite pada tingkat berjalan berhalangan, kasus dipindahkan NAIK SATU TINGKAT**
> untuk diaksep oleh tingkat di atasnya. Perpindahan itu adalah **aksi manual admin**, bukan
> otomatis, dan bukan "siapa saja boleh".

Ini menjawab konsekuensi operasional yang dicatat di atas — kasus **tidak** macet ketika pemilik
`KomiteID` absen.

**Yang tetap dilarang:** pemutus di tingkat **sama** yang bukan pemilik `KomiteID`, dan eskalasi
**turun** tingkat.

### Consequences tambahan

- **Inbox komite hanya menampilkan kasus sesuai posisi** — roster tingkat itu. Bukan seluruh antrean
  komite.
- Eskalasi adalah **tindakan yang direkam**: siapa yang memindahkan, kapan, dari tingkat mana ke
  tingkat mana. Ia mengubah siapa yang berwenang, jadi ia masuk jejak audit (**ADR-0007**).
- `KomiteCount` **ikut naik** saat eskalasi — konsekuensinya tingkat yang dilewati **tidak pernah
  memberi keputusan**, sehingga entri `KomiteAproval` untuk tingkat itu kosong. Bentuk pencatatannya
  adalah keputusan implementasi.
- ⚠️ Eskalasi **memperpendek** tangga persetujuan yang sebenarnya. Bila kebijakan menuntut jumlah
  persetujuan minimum, eskalasi melanggarnya — `[pertanyaan terbuka]`, tetapi **tidak memblokir**:
  perilaku lama pun tidak memaksakan minimum.

## OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-021** | RBAC lintas konteks belum ditetapkan; pemetaan `KomiteID` → akun bergantung padanya |
| **OQ-007** | Otorisasi Komite secara umum — korpus tidak memuat rule otorisasi |
