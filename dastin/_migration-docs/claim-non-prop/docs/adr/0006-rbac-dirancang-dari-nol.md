---
status: accepted
label: DECIDED
---

# RBAC sistem baru dirancang dari nol, bukan diwarisi

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Modul ini tidak memiliki otorisasi di tingkat rule: `pyPrivilegeName` kosong di seluruh 279 berkas, tidak ada satu pun `Rule-Access-When`, dan seluruh elemen UI bervisibilitas `ALWAYS` — termasuk kedua jalur penutupan klaim. Kontrol akses yang benar-benar berjalan bertumpu pada access group dan routing flow yang tidak ikut ter-export. Kami menetapkan model otorisasi sistem baru dirancang dari nol bersama pemilik proses, dan tidak ada satu pun aturan akses yang boleh diklaim sebagai warisan sistem lama.

## Consequences

Ini **menambah cakupan pekerjaan**, bukan sekadar catatan. Perancangan RBAC menjadi pekerjaan tersendiri dengan pemilik keputusan di sisi bisnis, bukan turunan otomatis dari analisis XML.

Empat label peran memang terbukti ada di `DataTransform\InsertChronology_DT.xml` — `Claim Admin`, `Claim Dept. Head`, `Operational Director`, `Technical Director` — tetapi itu label yang dicatatkan ke jejak audit, bukan aturan yang menegakkan siapa boleh melakukan apa. Pemetaannya ke orang pun bersifat hardcode per nama operator, sehingga tidak dapat dijadikan dasar model peran.

Sampai RBAC baru disepakati, setiap pernyataan tentang "siapa boleh melakukan apa" dalam dokumen mana pun berstatus hipotesis dan wajib berlabel EXTERNAL.
