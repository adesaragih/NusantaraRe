# Prompt — STEP D3 + D4 (penutup FASE A: catatan modul, glossary final, context-map, understanding-report)

**Status skill Matt Pocock (sudah diverifikasi manusia):** skill `/wayfinder` (dan
`/grill-with-docs`, `/to-spec`, `/to-tickets`, `/setup-matt-pocock-skills`) TERPASANG
(mattpocock-skills@1.2.3) tetapi frontmatter-nya `disable-model-invocation: true` -> skill HANYA
bisa dipicu oleh manusia (Anda), tidak bisa dipanggil model. Karena itu FASE A dijalankan sebagai
metodologi discovery proyek ini sendiri (bukan alur tiket wayfinder), dan itu SAH. Skill penulisan
FASE B tetap dipicu manusia sendiri nanti.

**Cara pakai:** salin SELURUH isi blok kode di bawah ke Claude Code, satu kali. Tidak ada yang
perlu diganti. Mengerjakan **5 task sekaligus**: catatan modul (20) + glossary final + context-map
+ understanding-report + paket handoff ke FASE B.

---

```
Jalankan STEP D3 dan D4 FASE A Discovery migrasi Nusantara Re. D0-D2 sudah selesai & direview:
20 modul terinventaris (D1) dan 20 konteks tertelusur end-to-end (D2, lihat
discovery/D2-CLOSING-REPORT.md). Ini SINTESIS dari artefak yang sudah ada, bukan telusur baru.

==== STATUS SKILL (sudah ditetapkan manusia - JANGAN diubah, JANGAN mencoba memanggil skill) ====
Diverifikasi: skill Matt Pocock /wayfinder, /grill-with-docs, /to-spec, /to-tickets,
/setup-matt-pocock-skills TERPASANG (mattpocock-skills@1.2.3) tetapi ber-frontmatter
disable-model-invocation: true -> HANYA bisa dipicu manusia, TIDAK boleh kamu panggil via Skill
tool. Jangan mencoba memanggilnya, jangan memasang ulang, dan JANGAN mereplikasi alur internal
skill wayfinder (peta decision-ticket di issue tracker).
Keputusan manusia: D3+D4 dikerjakan sebagai SINTESIS metodologi FASE A proyek ini (format
modules/*.md, glossary.md, context-map.md, understanding-report.md - sama seperti D0-D2). Ini SAH
karena bukan alur tiket wayfinder, melainkan runbook proyek yang sudah berjalan.
Aturan tambahan yang Anda minta: bila di kemudian hari sebuah prompt MEMINTA kamu menjalankan skill
Matt Pocock dan kamu TIDAK bisa (karena disable-model-invocation atau sebab lain), BERHENTI dan
beri tahu saya cara menjalankannya manual (ketik sendiri command /mattpocock-skills:<nama> di
Claude Code) - jangan diam-diam mereplikasi alur skill.

Kerjakan KELIMA task berikut, berurutan. Jangan berhenti di task 1-3; selesaikan sampai task 5.
Bila konteks membengkak, /handoff ke discovery/handoff/handoff-d3d4.md lalu lanjut sesi baru dari
task tertunda - bukan mengurangi cakupan.

SUMBER & BATAS:
- Baca: discovery/README.md, discovery/D1-CLOSING-REPORT.md, discovery/D2-CLOSING-REPORT.md,
  seluruh discovery/inventory/*.md, seluruh discovery/flows/*.md, discovery/glossary-draft.md,
  discovery/open-questions.md, discovery/inventory/_oq011-konflik-isi.md.
- Tulis HANYA ke OUTPUT_HASIL_RNM. Korpus D:\XML\RNM_BRD\ & D:\XML\nusantara-re\ READ-ONLY.
- SINTESIS dari artefak D1/D2 yang ada. Boleh buka korpus HANYA untuk mengisi celah tersisa, dengan
  bukti path+rule. TIDAK membuat ADR/BRD/spec/tiket (FASE B). TIDAK merancang Go/React/skema.
  TIDAK menilai benar/salah.

ATURAN WAJIB:
- Setiap klaim wajib bukti path+rule; rule ditulis class/nama/tipe (dari <pxInsName>).
- Nomor OQ diambil dari register open-questions.md, bukan dari ingatan.
- Yang tak pasti tetap "belum terverifikasi" + di open-questions.md (pemilik peran + apa yang
  diblokir). Label [terverifikasi]/[dugaan]/[pertanyaan terbuka]. Angka disertai perintah audit.
- Jangan menyatukan identitas berkonflik; jangan menebak arti kode/status/rumus/SP/JSON.

==== TASK 1 - CATATAN PEMAHAMAN PER MODUL (D3), 20 modul ====
Untuk tiap modul tulis discovery/modules/<nama modul>.md, disintesis dari inventory + flows-nya:
peran modul; proses/fitur utama (rujuk flows/<konteks>.md); entitas & tabel data yang disentuh;
integrasi eksternal; ketergantungan ke modul lain (dengan bukti, mis. Claim baca master Treaty;
Master Life tulis M_TREATY_IN); batasan & batas pengetahuan; OQ yang menyentuh modul. Label
[terverifikasi]/[dugaan]/[pertanyaan terbuka] tiap poin.

==== TASK 2 - FINALISASI GLOSSARY (D3) ====
Finalkan discovery/glossary-draft.md menjadi discovery/glossary.md: istilah domain (Indonesia &
Inggris apa adanya), kode/enumerasi ("arti belum terverifikasi" bila korpus tak menjelaskan), nama
workbasket, singkatan belum terjabarkan. Tiap entri wajib kolom Bukti (path+rule). Jangan mengarang
terjemahan istilah Indonesia, kepanjangan singkatan, atau arti field dari caption. Tandai
glossary.md sebagai kandidat seed CONTEXT.md untuk FASE B.

==== TASK 3 - CONTEXT-MAP (D4) ====
Tulis discovery/context-map.md: bounded context dari 20 modul + 20 flows, modul mana milik konteks
mana, cakupan bukti tiap konteks full/partial/absent. Sertakan temuan struktural D2: facultative =
satu ruleset berdiskriminator siklus (IsNB/IsRenewal/IsEDM); Claim = tiga basis kode terpisah;
Treaty In & Adjustment berbagi 320 identitas; Treaty Contract Out = master term INWARD (bukan
outward, OQ-022 terjawab); Komite = tangga approval berbasis data. Petakan kebocoran batas
antar-konteks (dengan bukti). Konteks yang tak bisa diturunkan dari korpus (mis. Identity & Access)
-> tandai absent + pertanyaan terbuka. JANGAN usulkan paket Go (itu keputusan FASE B/ADR).

==== TASK 4 - UNDERSTANDING-REPORT (D4) ====
Tulis discovery/understanding-report.md: narasi "bagaimana aplikasi ini bekerja" per konteks dan
antar-konteks (offer->akseptasi->produksi; klaim register->akseptasi->komite; treaty master &
realisasi; life). Sertakan: peta integrasi eksternal total; daftar batas pengetahuan total (70 SP,
8 objek db-link, 10 objek JSON, 49 DecisionTable tanpa baris, 5 keluarga When tak terbaca, 533
identitas berkonflik); daftar risiko migrasi (korpus mungkin dari dev; rumus di luar korpus; RBAC
harus dirancang; uang non-float; identitas/host ter-hardcode; AI pihak ketiga di alur underwriting).
Semua berbukti.

==== TASK 5 - PAKET HANDOFF KE FASE B ====
Tulis discovery/D3-D4-CLOSING-REPORT.md (BUKAN WAYFINDER-CLOSING.md, karena skill wayfinder tidak
dijalankan): ringkas hasil FASE A (D0-D4); tegaskan GATE - apa yang manusia harus review & putuskan
sebelum FASE B. Susun daftar seluruh OQ terbuka per pemilik peran (DBA/Product+UW/Actuarial/Finance/
IAM) dengan penanda pemblokir FASE B. Usulkan URUTAN FASE B berbasis bukti: konteks mana paling siap
untuk grilling lebih dulu (paling sedikit OQ pemblokir, bukti paling utuh) dan mana yang harus
menunggu jawaban bisnis.
Karena skill FASE B (/grill-with-docs, /to-spec, /to-tickets) hanya bisa dipicu manusia, sertakan
juga bagian "CARA MENJALANKAN FASE B SECARA MANUAL": command persis yang harus SAYA ketik sendiri
(mis. /mattpocock-skills:grill-with-docs) dan prasyaratnya - termasuk catatan bahwa skill ini
mungkin meminta issue tracker sedangkan D:\XML\RNM_BRD bukan repo git; sarankan opsi tracker
(mis. Markdown lokal di .scratch/<konteks>/ sesuai runbook) untuk saya putuskan. JANGAN memulai
grilling/spec/tiket sendiri.
Perbarui README section 9: D3 & D4 selesai, FASE A ditutup, menunggu GATE manusia sebelum FASE B.

==== VERIFIKASI SEBELUM LAPOR SELESAI ====
- 20 file modules/*.md + glossary.md + context-map.md + understanding-report.md +
  D3-D4-CLOSING-REPORT.md ada.
- Tiap klaim ada bukti path+rule; nomor OQ dari register; angka ada perintah audit.
- Tidak ada ADR/BRD/spec/tiket; tidak ada desain Go/React/skema; tidak ada pemanggilan skill.
- Korpus & nusantara-re tidak tersentuh (buktikan lewat git status / cek timestamp).

LAPORAN AKHIR (ringkas): status kelima task; jumlah bounded context + cakupan buktinya; total OQ per
pemilik peran + berapa pemblokir FASE B; usulan urutan FASE B + command manual yang harus saya ketik.
Lalu tanya saya: apakah saya (manusia) akan menjalankan /mattpocock-skills:grill-with-docs sendiri
pada konteks yang direkomendasikan, atau review paket discovery lengkap dulu.

Jangan mulai FASE B dan jangan memanggil skill apa pun di sesi ini.
```
