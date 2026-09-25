# -*- coding: utf-8 -*-
r"""
tulis-blok-erd-md.py -- menulis BLOK YANG DIBANGKITKAN ke dalam ERD.md.

    masukan : 4-erd-dan-tabel-datar/erd-skema-baru.json  (dari bangkitkan-erd-skema-baru.py)
              alat/cocok-relasi.json                     (dari cocok-relasi-erd-versus-ddl.py)
    keluaran: blok di antara penanda di dalam ERD.md

APA YANG DIBANGKITKAN, DAN APA YANG TIDAK -- ini pemisahan pokoknya
    DIBANGKITKAN : daftar FOREIGN KEY yang benar-benar ada di ddl-usulan/, dan
                   selisih dua arahnya terhadap sec 2.
    TIDAK, DAN TIDAK BOLEH : sec 2 itu sendiri. Ia memuat KEPUTUSAN -- perilaku
                   hapus per relasi -- yang TIDAK ADA di DDL. Membangkitkan ulang
                   sec 2 dari DDL akan MENGHAPUS ketiga puluh enam keputusan itu
                   tanpa satu galat pun. Itu sebab blok ini DITAMBAHKAN di
                   sebelahnya, bukan menggantikannya.

PENOLAKAN dilaporkan satu baris per sebab.
"""
import io, os, json
from collections import Counter

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
ERD = os.path.join(AKAR, '4-erd-dan-tabel-datar', 'ERD.md')
J1 = os.path.join(AKAR, '4-erd-dan-tabel-datar', 'erd-skema-baru.json')
J2 = os.path.join(AKAR, 'alat', 'cocok-relasi.json')
AWAL = '<!-- DIBANGKITKAN:cocok-silang-ddl -->'
AKHIR = '<!-- /DIBANGKITKAN:cocok-silang-ddl -->'
TANGGAL = '25 September 2026'

tolak = Counter()
D1 = json.load(io.open(J1, encoding='utf-8'))
D2 = json.load(io.open(J2, encoding='utf-8'))

L = []
w = L.append
w(AWAL)
w('')
w('## 2z. Cocok-silang terhadap `2-to-spec/ddl-usulan/` — **DIBANGKITKAN**')
w('')
w('> **Blok ini dibangkitkan `alat/tulis-blok-erd-md.py`, %s. Jangan disunting dengan tangan.**' % TANGGAL)
w('> Sumbernya `2-to-spec/ddl-usulan/` — berkas yang akan dibangun. **§2 di atas TIDAK**')
w('> **dibangkitkan**, dan itu disengaja: ia memuat **keputusan perilaku hapus** yang tidak ada di')
w('> DDL, dan membangkitkannya ulang akan menghapus keputusan itu tanpa satu galat pun.')
w('')
w('### 2z.1 Kunci asing yang benar-benar ada di DDL — %d' % len(D1['relasi']))
w('')
w('| # | Induk | Anak | Kolom | Kard. | ON DELETE di DDL |')
w('|---:|---|---|---|---|---|')
for r in D1['relasi']:
    w('| %d | `%s` | `%s` | `%s` | %s | %s |'
      % (r['no'], r['induk'], r['anak'], r['kolom'], r['kard'],
         r['on_delete'] if r['on_delete'] != '(tidak dinyatakan)' else '**(tidak dinyatakan)**'))
w('')

hd = D2['hilang_di_ddl']
he = D2['hilang_di_erd']
w('### 2z.2 Selisih dua arah terhadap §2')
w('')
w('| Arah | Berapa |')
w('|---|---:|')
w('| relasi di §2 (dalam skema) yang **tidak punya `FOREIGN KEY`** | **%d** |' % len(hd))
w('| `FOREIGN KEY` di DDL yang **tidak ada di §2** | **%d** |' % len(he))
w('| kolom yang §2 sebut dan **tidak ada di `KAMUS-KOLOM.md`** | **%d** |' % len(D2['kolom_hilang']))
w('| kunci asing yang **membawa `ON DELETE`** di DDL | **%d dari %d** |'
  % (D2['on_delete_di_ddl'], len(D1['relasi'])))
w('')
w('**Relasi di §2 tanpa kunci asing — %d:**' % len(hd))
w('')
w('| Induk | Anak | `[hapus: …]` yang §2 putuskan | Catatan §2 |')
w('|---|---|---|---|')
for r in hd:
    w('| `%s` | `%s` | `%s` | %s |' % (r['induk'], r['anak'], r['hapus'], r['catatan'] or '—'))
w('')
w('**Kunci asing di DDL yang tidak ada di §2 — %d:**' % len(he))
w('')
w('| Induk | Anak | Kolom |')
w('|---|---|---|')
for r in he:
    w('| `%s` | `%s` | `%s` |' % (r['induk'], r['anak'], r['kolom']))
w('')

tf = D1.get('tanpa_fk', [])
harus = [x for x in tf if x[2] == 'SEHARUSNYA ADA']
w('### 2z.3 Kolom rujukan tanpa kunci asing — %d, di antaranya **%d seharusnya ada**'
  % (len(tf), len(harus)))
w('')
w('Pembangkit DDL menurunkan kunci asing dari **pola nama `ID_<ENTITAS>`**. Setiap kolom rujukan')
w('yang dinamai lain karena itu **tidak memperoleh kunci asing**, dan ketiadaannya **tidak berbunyi**')
w('di mana pun — tabelnya tetap berdiri dan DDL-nya tetap sah.')
w('')
w('| Tabel | Kolom | Kedudukan | Sasaran yang dimaksud |')
w('|---|---|---|---|')
for t, k, ked, sas in tf:
    w('| `%s` | `%s` | %s | %s |' % (t, k, ('**%s**' % ked) if ked == 'SEHARUSNYA ADA' else ked, sas))
w('')
w('### 2z.4 Batas blok ini')
w('')
w('Ia mencocokkan **relasi**, bukan kolom bukan-kunci, dan **tidak** memeriksa kardinalitas yang §2')
w('tulis (`1` / `o` / `<`) terhadap `UNIQUE` di DDL — itu pemeriksaan **ketiga**, dan ia **belum ada**.')
w('')
w('Dan seperti setiap pemeriksaan di modul ini: **"cocok" bukan "lengkap"**. Semesta §10 masih')
w('kurang **340 properti titik buta** (`L-8`, ditagih `M-4`).')
w('')
w(AKHIR)

blok = '\n'.join(L)
s = io.open(ERD, encoding='utf-8').read()
if AWAL in s and AKHIR in s:
    a = s.index(AWAL); b = s.index(AKHIR) + len(AKHIR)
    s = s[:a] + blok + s[b:]
    cara = 'blok lama DIGANTI'
else:
    anc = '## 3. Atribut bukan-kunci'
    if anc not in s:
        tolak['penambat "## 3." tidak ditemukan'] += 1
        cara = 'GAGAL — penambat tidak ada'
    else:
        s = s.replace(anc, blok + '\n\n---\n\n' + anc, 1)
        cara = 'blok BARU disisipkan sebelum §3'
io.open(ERD, 'w', encoding='utf-8').write(s)

print('=== tulis-blok-erd-md.py ===')
print('   %s' % cara)
print('   kunci asing didaftar      : %d' % len(D1['relasi']))
print('   relasi §2 tanpa FK        : %d' % len(hd))
print('   FK tanpa relasi di §2     : %d' % len(he))
print('   kolom rujukan tanpa FK    : %d  (%d seharusnya ada)' % (len(tf), len(harus)))

# ══ blok kedua: gambar mermaid sec 6, DIBANGKITKAN ══════════════════
AWAL2 = '<!-- DIBANGKITKAN:gambar-mermaid -->'
AKHIR2 = '<!-- /DIBANGKITKAN:gambar-mermaid -->'
KARD = {'1:1': '||--||', '1:N': '||--o{'}
G = [AWAL2, '']
G.append('> **Dibangkitkan dari `2-to-spec/ddl-usulan/`, %s.** Ia memuat **tepat** kunci asing' % TANGGAL)
G.append('> yang ada di DDL — %d — dan karena itu **tidak dapat basi**. Relasi yang §2' % len(D1['relasi']))
G.append('> putuskan tetapi **belum punya kunci asing** (§2z.2) **tidak muncul di sini**, dan')
G.append('> ketidakmunculannya itulah yang membuat gambar ini berguna sebagai pemeriksa.')
G.append('')
G.append('```mermaid')
G.append('erDiagram')
for r in D1['relasi']:
    G.append('    %s %s %s : \"%s\"' % (r['induk'], KARD.get(r['kard'], '||--o{'),
                                       r['anak'], r['kolom']))
G.append('```')
G.append('')
G.append(AKHIR2)
gambar = chr(10).join(G)

s2 = io.open(ERD, encoding='utf-8').read()
if AWAL2 in s2 and AKHIR2 in s2:
    a2 = s2.index(AWAL2); b2 = s2.index(AKHIR2) + len(AKHIR2)
    s2 = s2[:a2] + gambar + s2[b2:]
    cara2 = 'gambar mermaid DIGANTI'
else:
    import re as _re
    m2 = _re.search(r'(## 6\. Gambar[^\n]*\n)(\n?```mermaid.*?```\n)', s2, _re.S)
    if m2:
        s2 = s2[:m2.start(2)] + chr(10) + gambar + chr(10) + s2[m2.end(2):]
        cara2 = 'gambar mermaid lama DIGANTI blok yang dibangkitkan'
    else:
        cara2 = 'GAGAL — blok mermaid §6 tidak ditemukan'
        tolak['blok mermaid §6 tidak ditemukan'] += 1
io.open(ERD, 'w', encoding='utf-8').write(s2)
print('   %s' % cara2)

print()
print('DITOLAK -- satu baris per sebab:')
for k, v in tolak.most_common():
    print('   %-56s %d' % (k, v))
if not tolak:
    print('   (tidak ada)')
print()
print('-- BATAS PERKAKAS INI --')
print('   Ia TIDAK membangkitkan sec 2. Sec 2 memuat KEPUTUSAN perilaku hapus yang')
print('   tidak ada di DDL; membangkitkannya ulang akan menghapusnya tanpa satu galat pun.')
