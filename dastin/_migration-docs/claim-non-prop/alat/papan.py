# -*- coding: utf-8 -*-
"""Bangun README papan tiket dari field `status:` di tiap berkas tiket.

DIJALANKAN DARI MANA SAJA. Ia MENULIS satu berkas dan hanya satu:
    .scratch/claim-non-prop-lapisan-data/issues/README.md
Berkas tiketnya sendiri TIDAK PERNAH disentuh - ia hanya dibaca.

Berbeda dari alat/periksa-penamaan.py, yang hanya membaca dan tidak menulis
apa pun. Bedanya disebut di sini supaya tidak perlu ditebak dari kodenya.

angun README papan dari field `status:` di tiap berkas tiket.

Status adalah field, bukan lokasi folder. Folder hanya kerapian.
Pembangkit ini MEMBACA berkas tiket; ia tidak pernah menghapusnya.
"""
import io, os, re

ISS = r'D:\XML_NURE\.scratch\claim-non-prop-lapisan-data\issues'
baca = lambda p: io.open(p, encoding='utf-8').read()

T = {}
for folder in (ISS, os.path.join(ISS, '_selesai'), os.path.join(ISS, '_tertahan'), os.path.join(ISS, '_mati')):
    if not os.path.isdir(folder):
        continue
    for f in sorted(os.listdir(folder)):
        if not f.endswith('.md') or f == 'README.md':
            continue
        t = baca(os.path.join(folder, f))
        m = re.search(r'^# (\d\d): (.*)$', t, re.M)
        nn, judul = m.group(1), m.group(2).strip()
        status = re.search(r'^status: (\S+)', t, re.M)
        luar = re.search(r'^ditahan-luar: \[(.*)\]', t, re.M)
        asal = re.search(r'Asal: `(T-\d+)`', t)
        T[nn] = dict(
            nn=nn, judul=judul, berkas=f,
            folder=os.path.basename(folder) if folder != ISS else '',
            status=status.group(1) if status else 'aktif',
            luar=[x.strip() for x in luar.group(1).split(',')] if luar else [],
            tunggu=[x.strip() for x in re.search(r'^menunggu-luar: \[(.*)\]', t, re.M).group(1).split(',')]
                   if re.search(r'^menunggu-luar: \[(.*)\]', t, re.M) else [],
            asal=asal.group(1) if asal else '—',
            blok=re.findall(r'^- `(\d\d)` —', t, re.M),           # hidup saja; yang dicoret tak cocok
            blok_mati=re.findall(r'^- ~~`(\d\d)`~~ \*\((\w+)\)\*', t, re.M))

AKTIF   = [k for k in sorted(T) if T[k]['status'] == 'aktif']
TERTAHN = [k for k in sorted(T) if T[k]['status'] == 'tertahan']
INSTAN  = [k for k in sorted(T) if T[k]['status'] == 'menunggu-instance']
SELESAI = [k for k in sorted(T) if T[k]['status'].startswith('selesai')]
MATI    = [k for k in sorted(T) if T[k]['status'] == 'mati']
PAPAN   = sorted(AKTIF + TERTAHN + INSTAN)

# beku transitif: apa saja yang berakar pada tiket yang ditahan ADR-0028
def akar_luar(nn, lihat=None):
    lihat = lihat or set()
    if nn in lihat:
        return set()
    lihat.add(nn)
    keluar = set(T[nn]['luar'])
    for b in T[nn]['blok']:
        if b in T and T[b]['status'] in ('aktif', 'menunggu-instance'):
            keluar |= akar_luar(b, lihat)
    return keluar

beku = {k: akar_luar(k) for k in PAPAN}
mulai = [k for k in PAPAN if not beku[k] and not T[k]['blok'] and T[k]['status'] == 'aktif']
per_penahan = {}
for k in PAPAN:
    for x in beku[k]:
        per_penahan.setdefault(x, []).append(k)

R = ['# Papan tiket — lapisan data Claim Non Prop', '',
 '**Tiket: aktif %d · menunggu instance %d · tertahan %d · selesai %d · mati %d**'
 % (len(AKTIF), len(INSTAN), len(TERTAHN), len(SELESAI), len(MATI)), '',
 '> Pencacah ini menghitung **tiket**. Register REQ punya pencacahnya sendiri dengan format yang mirip; '
 'keduanya diberi label bendanya supaya tidak pernah tertukar.', '',
 '## Yang sebenarnya dapat dimulai', '',
 '**%d dari %d tiket di papan.** Sisanya tertahan sesuatu di **luar papan** — keputusan yang belum '
 'dikonfirmasi atau REQ yang belum kembali — bukan tertahan tiket lain. Angka ini, bukan "aktif %d", '
 'yang menggambarkan keadaan proyek.' % (len(mulai), len(PAPAN), len(AKTIF)), '',
 '| Penahan dari luar papan | Tiket yang dibekukannya | Ditahan langsung |', '|---|---|---|']
for x in sorted(per_penahan):
    ls = sorted(per_penahan[x])
    lang = [k for k in ls if x in T[k]['luar']]
    R.append('| **%s** | %d | langsung: %s |' % (x, len(ls), ' '.join('`%s`' % k for k in lang)))
R += ['', ('Dapat dimulai sekarang: ' + ' '.join('`%s`' % k for k in mulai) + '.') if mulai
       else '**Tidak ada yang dapat dimulai sekarang.** Empat sapuan sumber yang tidak menyentuh '
            'gerbang sudah dikerjakan sampai habis; sisanya menunggu konfirmasi ADR-0028 dan REQ-032.', '',
 'Diterbitkan 18 September 2026 dari `_migration-docs/claim-non-prop/TICKETS.md`, '
 'dibangun ulang dari field `status:` tiap berkas.',
 'Tracker belum terkonfigurasi (`/setup-matt-pocock-skills` belum dijalankan), jadi papannya berkas lokal.', '',
 '## Aturan papan', '',
 '- **Status adalah field di dalam berkas tiket** (`status: aktif | tertahan | selesai | selesai-sebagian`), '
 'bukan lokasi foldernya. Folder `_selesai/` dan `_tertahan/` hanya kerapian; pembangkit membaca field, '
 'dan **tidak pernah menghapus berkas tiket**.',
 '- **Penahan dari luar papan ditulis dengan nama aslinya** — `ADR-0028`, `REQ-032` — tidak diterjemahkan '
 'jadi nomor tiket. Yang menahan dari luar harus terlihat berasal dari luar.',
 '- **Penahan yang sudah selesai dicoret, tidak dihapus**, supaya rantainya tetap terbaca: '
 '`~~05~~ *(selesai)*`.',
 '- **Nomor yang lompat bukan kekeliruan.** Penomoran tidak disusun ulang, karena menyusunnya ulang '
 'memutus setiap rujukan yang sudah ada.', '',
 '## Papan', '',
 '| # | Asal | Tiket | Ditahan oleh |', '|---|---|---|---|']
for k in PAPAN:
    v = T[k]
    tahan = ['**%s**' % x for x in v['luar']] + ['*(%s)*' % x for x in v['tunggu']]
    tahan += ['`%s`' % b for b in v['blok']]
    tahan += ['~~`%s`~~' % b for b, _ in v['blok_mati']]
    R.append('| `%s` | %s | %s | %s |' % (k, v['asal'], v['judul'], ' '.join(tahan) or '— *dapat dimulai*'))
R += ['', '## Menunggui, tidak menahan', '',
 'Tiket ini **boleh jalan**. Jawaban yang ditunggu mengubah satu hal yang tiketnya sudah menyebut '
 '— sebuah constraint, sebuah aturan penolakan, sebuah jalur migrasi — bukan rancangannya. '
 'Dipisahkan dari penahan sungguhan supaya papan tidak berteriak serigala.', '',
 '| # | Menunggu | Yang berubah bila jawabannya lain |', '|---|---|---|']
for k in PAPAN:
    for x in T[k]['tunggu']:
        R.append('| `%s` | **%s** | lihat *Blocked by* di berkasnya |' % (k, x))

R += ['', '## Dapat dimulai hari pertama', '']
if mulai:
    for k in mulai:
        R.append('- `%s` %s — %s' % (k, T[k]['asal'], T[k]['judul']))
else:
    R.append('**Tidak ada.** Setiap tiket yang tersisa tertahan sesuatu di luar papan. '
             'Itu keadaan yang sah dan bukan kebuntuan kerja: yang dapat dikerjakan tanpa '
             'menyentuh gerbang sudah dikerjakan sampai habis.')
R += ['', '---', '', '## Menunggu instance — bukan aktif', '',
 'Rancangannya lengkap; yang kurang mesinnya. Uji ini baru dapat hijau setelah DDL dijalankan di '
 'instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan '
 'pekerjaan yang tidak dapat dimulai siapa pun.', '', '| # | Tiket |', '|---|---|']
for k in INSTAN:
    R.append('| `%s` | %s |' % (k, T[k]['judul']))
R += ['', '## Tertahan di luar papan', '', '| # | Tiket | Ditahan |', '|---|---|---|']
for k in sorted(T):
    if T[k]['status'] == 'tertahan':
        R.append('| `%s` | %s | %s |' % (k, T[k]['judul'], ' '.join('**%s**' % x for x in T[k]['luar'])))
R += ['', '## Mati — dibatalkan, bukan ditunda', '',
 'Tiket yang premisnya gugur. Berkasnya tidak dihapus: tiket yang pernah ada adalah bukti '
 'bahwa premisnya diperiksa, bukan dilewatkan. Nomornya tidak dipakai ulang.', '',
 '| # | Tiket | Sebab |', '|---|---|---|'] if MATI else []
for k in MATI:
    R.append('| `%s` | %s | premis gugur — lihat kepala berkasnya di `_mati/` |' % (k, T[k]['judul']))
R += ['', '## Selesai', '', '| # | Tiket | Status |', '|---|---|---|']
for k in SELESAI:
    R.append('| `%s` | %s | %s |' % (k, T[k]['judul'], T[k]['status']))
R.append('')
io.open(os.path.join(ISS, 'README.md'), 'w', encoding='utf-8').write('\n'.join(R))

print('papan  : aktif %d · tertahan %d · selesai %d · mati %d' % (len(AKTIF), len(TERTAHN), len(SELESAI), len(MATI)))
print('mulai  : %s' % ' '.join(mulai))
for x in sorted(per_penahan):
    print('beku %-9s: %d' % (x, len(per_penahan[x])))
