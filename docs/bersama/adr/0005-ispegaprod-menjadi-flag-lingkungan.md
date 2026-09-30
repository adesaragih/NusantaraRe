---
status: accepted
tanggal: 2026-09-14
sumber: grilling Ronde 1 Q6, keputusan work owner
---

# `IsPEGAPROD` tidak ditiru sebagai rule; diganti flag lingkungan eksplisit

Di Pega, tiga efek keluar Claim — Life digerbangi rule `When` bernama `IsPEGAPROD` yang
**kondisinya tidak terbaca dari korpus**. Di sistem baru, gerbang itu menjadi **flag lingkungan
eksplisit** (`ENV=production`) yang menggerbangi ketiga efek yang sama — sehingga perilakunya
terbaca, bukan tersembunyi di dalam rule.

Keputusan ini tentang **apakah efek keluar dijalankan**. Dari mana alamatnya berasal adalah
keputusan terpisah — lihat **ADR-0013** (alamat di-lookup runtime dari `M_LINK_SERVICE`; ini
menggantikan ADR-0004). Keduanya **sengaja tidak digabung**.

## Keadaan sekarang `[terverifikasi]`

`Claim Life/When/IsPEGAPROD.xml` → `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN`.
`<pyLabel>`-nya hanya berisi template kosong `[first value][relation][second value]` — **apa yang
diujinya tidak diketahui** (**OQ-029**). Identitas ini juga **terdaftar berkonflik** di
`discovery/inventory/_oq011-konflik-isi.md` entri **#305** — isinya berbeda antar modul.

Tiga efek keluar yang digerbanginya di modul ini:

| Efek | Rule |
| --- | --- |
| Unggah berkas | `Claim Life/Activity/InsertGoogleStorage_Act.xml` |
| Simpan utama | `Claim Life/Activity/SaveOutStandingLife_Act.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT`, 642.787 byte) |
| Email | `Claim Life/Activity/SendEmailKlaimLF.xml` |

Perintah audit:
```
grep -rl "IsPEGAPROD" "Claim Life" --include="*.xml"
```

## Considered Options

- **Flag lingkungan eksplisit** (`ENV=production`) — dipilih
- Migrasikan `IsPEGAPROD` apa adanya — **tidak mungkin**: kondisinya tidak terbaca, dan rule-nya
  berkonflik antar modul (OQ-011 #305). Menyalin salah satu varian berarti menebak
- Jalankan efek keluar tanpa gerbang — ditolak: menghapus perilaku yang jelas disengaja

## Consequences

- Perilaku non-production menjadi **terbaca dari konfigurasi**, bukan dari isi rule.
- Ketiga efek keluar digerbangi **satu** flag yang sama, sebagaimana keadaan sekarang. Bila
  kemudian ternyata ketiganya perlu gerbang berbeda, itu perubahan berikutnya — bukan yang
  diputuskan di sini.
- `SaveOutStandingLife_Act` adalah **simpan utama**, bukan sekadar efek samping. Menggerbanginya
  dengan flag lingkungan berarti **di non-production data tidak tersimpan lewat jalur itu** —
  konsekuensi ini perlu disadari saat menyiapkan lingkungan uji.

## OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-029** (terjawab untuk Claim — Life) | Kondisi asli `IsPEGAPROD` tetap tidak diketahui — keputusan ini **menghindari**, bukan menjawab. Varian di 13 modul lain tetap terbuka |
| **OQ-011** entri **#305** | Versi `IsPEGAPROD` mana yang berlaku di production |
| **OQ-018** (terjawab untuk Claim — Life) | `jboss1073` = production, `jboss117` = dev, sistem mirroring |
