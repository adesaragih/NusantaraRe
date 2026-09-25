# -*- coding: utf-8 -*-
r"""
bangkitkan-dari-erd-v2.py -- struktur data, tabel datar, dan ERD HTML,
KETIGANYA dibangkitkan dari SATU sumber: erd-v2.json.

    masukan : alat/erd-v2.json          <- hasil urai-erd-v2.py, salinan apa adanya
    keluaran: 4-erd-dan-tabel-datar/ERD-TREATY-IN-DAN-EDM.html
              4-erd-dan-tabel-datar/STRUKTUR-DATA-ERD-V2.md
              4-erd-dan-tabel-datar/TABEL-DATAR-ERD-V2.md

KENAPA KETIGANYA DARI SATU SUMBER
    Pemilik proses menetapkan Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx sebagai
    sumber, "PERSIS, TIDAK ADA PERUBAHAN", dan meminta struktur data, tabel datar,
    dan ERD HTML ketiganya mengikutinya. Tiga berkas yang ditulis tangan dari satu
    gambar akan berbeda dalam sebulan. Tiga berkas yang DIBANGKITKAN dari satu
    berkas tidak bisa berbeda -- dan bila berbeda, ALAT INI yang salah.

SATU HAL YANG DITAMBAHKAN, dan ia diminta
    GARIS PENGHUBUNG. Legenda berkas sumbernya berbunyi:

        "Garis -- garis tegak turun dari kotak induk, lalu siku mendatar
         menyentuh kotak anak. Anak menjorok empat kolom dari induknya."

    Di berkas xlsx-nya garis itu TIDAK ADA -- diperiksa satu per satu, tepi selnya
    hanya membingkai kotak, tidak menghubungkannya. Yang ada hanya indentasi.
    Di HTML ini garis itu DIGAMBAR, persis seperti bunyi legendanya.

APA YANG TIDAK DIUBAH
    Nama tabel, kardinalitas, PK, FK, ON DELETE, jalur Pega, padanan model baru,
    tingkat bukti, penanda BERSAMA, urutan kotak, kedalaman, pembagian empat
    lembar, Daftar Relasi, dan Catatan & Batas -- seluruhnya apa adanya.
    Termasuk BANNER-nya: sumbernya POTRET SISTEM LAMA, dan keluarannya
    mengatakan hal yang sama.

PENOLAKAN dilaporkan satu baris per sebab.

Pakai:  PYTHONIOENCODING=utf-8 python alat/bangkitkan-dari-erd-v2.py
"""
import html, io, os, json, sys
from collections import Counter, OrderedDict, defaultdict

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
ERD = os.path.join(AKAR, '4-erd-dan-tabel-datar')
SUMBER = os.path.join(AKAR, 'alat', 'erd-v2.json')
HTML = os.path.join(ERD, 'ERD-TREATY-IN-DAN-EDM.html')
MD_STRUKTUR = os.path.join(ERD, 'STRUKTUR-DATA-ERD-V2.md')
MD_DATAR = os.path.join(ERD, 'TABEL-DATAR-ERD-V2.md')

tolak = Counter()
d = json.load(io.open(SUMBER, encoding='utf-8'))
e = html.escape
lembar = d['lembar']
relasi = d['relasi']
catatan = d['catatan']

# ══ Pohon per lembar, disusun dari kedalaman yang sudah terbaca ════════════
def pohon(kotak):
    akar, tumpuk = [], []
    for k in kotak:
        n = dict(k, anak=[])
        while tumpuk and tumpuk[-1]['dalam'] >= n['dalam']:
            tumpuk.pop()
        if not tumpuk:
            if n['dalam'] != 0:
                tolak['kotak berkedalaman > 0 tanpa induk di atasnya'] += 1
            akar.append(n)
        else:
            tumpuk[-1]['anak'].append(n)
        tumpuk.append(n)
    return akar


# ══ Kumpulan tabel unik, untuk kedua berkas markdown ═══════════════════════
tabel = OrderedDict()
di_lembar = defaultdict(list)
for nama, L in lembar.items():
    for k in L['kotak']:
        di_lembar[k['nama']].append(nama)
        if k['nama'] not in tabel:
            tabel[k['nama']] = k

# ═══════════════════════════════ HTML ══════════════════════════════════════
CSS = """
:root{--bg:#f7f9fc;--kartu:#fff;--tepi:#d6dee8;--tua:#1F4E79;--biru:#2E75B6;
--kuning:#FFF2CC;--muda:#EDF3F9;--teks:#1b1f24;--redup:#595959;--merah:#C00000;
--jingga:#e07b00;--garis:#8fa8c4}
@media(prefers-color-scheme:dark){:root:not([data-theme=light]){--bg:#11161d;
--kartu:#181f28;--tepi:#2b3644;--muda:#1c2733;--kuning:#3a3420;--teks:#e6eaf0;
--redup:#98a3b3;--garis:#41576f}}
:root[data-theme=dark]{--bg:#11161d;--kartu:#181f28;--tepi:#2b3644;--muda:#1c2733;
--kuning:#3a3420;--teks:#e6eaf0;--redup:#98a3b3;--garis:#41576f}
*{box-sizing:border-box}
body{margin:0;background:var(--bg);color:var(--teks);
font:13.5px/1.5 -apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif}
.bungkus{max-width:1220px;margin:0 auto;padding:22px 16px 72px}
h1{font-size:23px;margin:0 0 4px;color:var(--tua)}
@media(prefers-color-scheme:dark){h1{color:var(--biru)}}
h2{font-size:17px;margin:34px 0 4px;padding-bottom:6px;border-bottom:2px solid var(--tepi)}
.sub{color:var(--redup);font-size:12.5px;margin:0 0 16px}
.pita{background:var(--merah);color:#fff;padding:10px 14px;border-radius:6px;
font-weight:700;margin:0 0 6px;font-size:13px}
.pita.k{background:var(--kuning);color:#000;font-weight:600;font-size:12px}
.catatan{background:var(--kartu);border:1px solid var(--tepi);
border-left:4px solid var(--jingga);border-radius:6px;padding:11px 13px;
margin:12px 0;font-size:12.5px}
.angka{display:flex;flex-wrap:wrap;gap:9px;margin:12px 0}
.angka div{background:var(--kartu);border:1px solid var(--tepi);border-radius:6px;
padding:7px 13px;min-width:96px}
.angka b{display:block;font-size:17px;color:var(--tua)}
@media(prefers-color-scheme:dark){.angka b{color:var(--biru)}}
.angka span{font-size:10.5px;color:var(--redup);text-transform:uppercase;letter-spacing:.05em}

/* ── POHON: garis tegak dari induk, siku mendatar menyentuh anak ── */
.pohon,.pohon ul{list-style:none;margin:0;padding:0}
.pohon ul{margin-left:30px;position:relative}
.pohon ul::before{content:"";position:absolute;left:0;top:0;bottom:0;
border-left:2px solid var(--garis)}
.pohon ul>li{position:relative;padding-left:26px}
.pohon ul>li::before{content:"";position:absolute;left:0;top:19px;width:24px;
border-top:2px solid var(--garis)}
.pohon ul>li:last-child::after{content:"";position:absolute;left:-2px;top:21px;
bottom:0;border-left:2px solid var(--bg)}
.pohon>li{margin-bottom:6px}
.pohon li{margin:0 0 6px}

.kotak{background:var(--kartu);border:2px solid var(--tua);border-radius:5px;
overflow:hidden;max-width:820px}
.kotak .k1{background:var(--tua);color:#fff;font-weight:700;font-size:12.5px;
padding:5px 10px;display:flex;flex-wrap:wrap;gap:8px;align-items:baseline}
.kotak.anak{border-color:var(--biru)}
.kotak.anak .k1{background:var(--biru)}
.kotak .k1 em{font-style:normal;font-weight:400;font-size:11px;opacity:.92}
.kotak .k2{background:var(--kuning);color:#000;font-weight:700;font-size:11.5px;
padding:4px 10px;font-family:ui-monospace,Consolas,monospace}
.kotak .k3{background:var(--muda);color:var(--redup);font-size:11px;padding:4px 10px}
.tag{display:inline-block;border:1px solid currentColor;border-radius:3px;
padding:0 5px;font-size:9.5px;font-weight:700;letter-spacing:.04em}
.tag.b{color:#ffd966}.tag.m{color:#ffb3b3}
.bukti{float:right;font-size:9.5px;letter-spacing:.04em;opacity:.85}
table{border-collapse:collapse;width:100%;font-size:12px}
th{background:var(--tua);color:#fff;text-align:left;padding:6px 9px;
font-size:10.5px;text-transform:uppercase;letter-spacing:.05em}
td{padding:5px 9px;border-top:1px solid var(--tepi);vertical-align:top}
tr:nth-child(even) td{background:var(--muda)}
td.kun{background:var(--kuning)!important}
.lebar{background:var(--kartu);border:1px solid var(--tepi);border-radius:8px;overflow-x:auto}
code{font-family:ui-monospace,Consolas,monospace;font-size:11.5px}
@media print{body{background:#fff}.kotak,.lebar{break-inside:avoid}
.pohon ul>li:last-child::after{border-left-color:#fff}}
@media(max-width:640px){.pohon ul{margin-left:14px}.pohon ul>li{padding-left:16px}
.pohon ul>li::before{width:14px}}
"""


def kotak_html(k):
    anak = k['dalam'] > 0
    tags = []
    if k['kard']:
        tags.append('<span class="tag b">%s</span>' % e(k['kard']))
    if k['bersama']:
        tags.append('<span class="tag b">BERSAMA → %s</span>' % e(k['bersama']))
    k2 = ('PK&nbsp;&nbsp; <b>%s</b>' % e(k['pk'])) if k['pk'] else ''
    if k['fk']:
        f = k['fk']
        k2 = ('FK&nbsp;&nbsp; <b>%s</b> &rarr; %s.%s &nbsp;&nbsp;<span class="tag m">%s</span>'
              % (e(f['kolom']), e(f['induk']), e(f['kolom_induk']), e(f['on_delete'])))
    if not k2:
        tolak['kotak tanpa PK maupun FK yang terbaca'] += 1
        k2 = '<i>(kunci tidak terbaca di sumbernya)</i>'
    bukti = ('<span class="bukti">%s</span>' % e(k['bukti'])) if k['bukti'] else ''
    pad = (' &nbsp;<b>[%s]</b>' % e(k['padanan'])) if k['padanan'] else ''
    return ('<div class="kotak%s"><div class="k1"><span>%s</span>%s</div>'
            '<div class="k2">%s</div>'
            '<div class="k3">%s&larr; %s%s</div></div>'
            % (' anak' if anak else '', e(k['nama']), ''.join(tags), k2,
               bukti, e(k['jalur']) or '<i>(tidak ada jalur Pega)</i>', pad))


def daftar_html(simpul):
    out = []
    for n in simpul:
        anak = ('<ul>%s</ul>' % daftar_html(n['anak'])) if n['anak'] else ''
        out.append('<li>%s%s</li>' % (kotak_html(n), anak))
    return ''.join(out)


bag = []
bag.append('<div class="pita">%s</div>' % e(d['banner'][0]))
bag.append('<div class="pita k">%s</div>' % e(d['banner'][1]))
bag.append('<h1>ERD Treaty In dan Treaty In EDM — empat lembar</h1>')
bag.append('<p class="sub">Dibangkitkan dari <code>%s</code> lewat '
           '<code>alat/urai-erd-v2.py</code> &rarr; <code>alat/erd-v2.json</code> &rarr; '
           '<code>alat/bangkitkan-dari-erd-v2.py</code>. '
           '<b>Isinya tidak diubah satu baris pun.</b></p>' % e(d['sumber']))

n_kotak = sum(len(L['kotak']) for L in lembar.values())
bag.append('<div class="angka">'
           '<div><b>%d</b><span>tabel unik</span></div>'
           '<div><b>%d</b><span>kotak di 4 lembar</span></div>'
           '<div><b>%d</b><span>relasi</span></div>'
           '<div><b>%d</b><span>dipakai bersama</span></div></div>'
           % (len(tabel), n_kotak, len(relasi),
              sum(1 for n in di_lembar if len(di_lembar[n]) > 1)))

bag.append('<div class="catatan"><b>Satu hal ditambahkan terhadap berkas sumbernya: '
           'GARIS PENGHUBUNG.</b> Legenda sumbernya berbunyi <i>"garis tegak turun dari kotak '
           'induk, lalu siku mendatar menyentuh kotak anak"</i> — tetapi di berkas xlsx-nya '
           'garis itu <b>tidak ada</b>; diperiksa satu per satu, tepi selnya hanya membingkai '
           'kotak, dan hierarkinya hanya terbaca dari indentasi. Di halaman ini garis itu '
           '<b>digambar</b>, persis seperti bunyi legendanya. Selain itu, tidak ada yang '
           'diubah.</div>')

for nama, L in lembar.items():
    ak = pohon(L['kotak'])
    bag.append('<h2>%s <small style="font-weight:400;color:var(--redup)">— %d kotak</small></h2>'
               % (e(nama), len(L['kotak'])))
    if L.get('subjudul'):
        bag.append('<p class="sub">%s</p>' % e(L['subjudul']))
    bag.append('<ul class="pohon">%s</ul>' % daftar_html(ak))

# ── Daftar Relasi ──────────────────────────────────────────────────────────
bag.append('<h2>Daftar Relasi <small style="font-weight:400;color:var(--redup)">— %d</small></h2>'
           % len(relasi))
if relasi:
    kol = list(relasi[0].keys())
    th = ''.join('<th>%s</th>' % e(c) for c in kol)
    tr = ''.join('<tr>%s</tr>' % ''.join(
        '<td%s>%s</td>' % (' class="kun"' if c == 'ON DELETE' else '', e(str(r.get(c, ''))))
        for c in kol) for r in relasi)
    bag.append('<div class="lebar"><table><thead><tr>%s</tr></thead><tbody>%s</tbody>'
               '</table></div>' % (th, tr))
else:
    tolak['Daftar Relasi kosong'] += 1

# ── Legenda, disalin dari lembar gambar ────────────────────────────────────
leg = []
for nama, L in lembar.items():
    pass
bag.append('<h2>Keterangan — disalin dari lembar gambar</h2>')
bag.append('<div class="catatan">'
           '<p><b>Kotak</b> — baris 1: nama tabel, kardinalitas, dan apakah dipakai bersama '
           'lembar lain · baris 2 (kuning): PK atau FK ke induk · baris 3 (biru muda): jalur Pega '
           'asal, padanan di model baru, dan tingkat buktinya.</p>'
           '<p><b>Garis</b> — garis tegak turun dari kotak induk, lalu siku mendatar menyentuh '
           'kotak anak. Anak menjorok dari induknya.</p>'
           '<p><b>Bukti</b> — <code>JALUR-PENUH</code>: jalur Pega lengkap ditemukan di badan '
           'seksi (kuat). <code>DAUN-RELATIF</code>: hanya nama ruas terakhir yang ditemukan, '
           'misalnya <code>.COBList</code> (lemah, mungkin memungut nama milik entitas lain).</p>'
           '<p><b>BATAS</b> — sapuan ini hanya membaca Section. Activity, Data Transform, Report '
           'Definition, dan RDB List <b>TIDAK</b> disapu. Nol rujukan berarti <b>BELUM '
           'TERPERIKSA</b>, bukan tidak dipakai.</p>'
           '<p><b>ON DELETE CASCADE</b> mengikuti konvensi <code>contooh.xlsx</code>. Ia '
           '<b>USULAN rancangan</b>, bukan perilaku sistem lama.</p></div>')

# ── Catatan & Batas ────────────────────────────────────────────────────────
bag.append('<h2>Catatan &amp; Batas — disalin apa adanya</h2>')
tr = ''.join('<tr><td style="width:24%%"><b>%s</b></td><td>%s</td></tr>' % (e(a), e(b))
             for a, b in catatan if a)
bag.append('<div class="lebar"><table><tbody>%s</tbody></table></div>' % tr)

doc = ('<!DOCTYPE html><html lang="id"><head><meta charset="utf-8">'
       '<meta name="viewport" content="width=device-width,initial-scale=1">'
       '<title>ERD Treaty In dan EDM</title><style>%s</style></head>'
       '<body><div class="bungkus">%s</div></body></html>' % (CSS, ''.join(bag)))
io.open(HTML, 'w', encoding='utf-8').write(doc)

# ════════════════════════ STRUKTUR-DATA-ERD-V2.md ══════════════════════════
m = []
m.append('# Struktur data — **dibangkitkan dari `%s`**' % d['sumber'])
m.append('')
m.append('> ### %s' % d['banner'][0])
m.append('>')
m.append('> %s' % d['banner'][1])
m.append('')
m.append('**Dibangkitkan:** [`alat/bangkitkan-dari-erd-v2.py`](../alat/bangkitkan-dari-erd-v2.py) '
         '· **Sumber tunggal:** `alat/erd-v2.json`')
m.append('')
m.append('> **Berkas ini tidak ditulis tangan, dan tidak boleh disunting.** Ia salinan struktur '
         'yang ada di berkas sumbernya — nama tabel, kardinalitas, kunci, jalur Pega, padanan '
         'model baru, dan tingkat bukti seluruhnya apa adanya. Sunting sumbernya, lalu jalankan '
         'ulang alatnya.')
m.append('')
m.append('> ### BERKAS INI GABUNGAN DUA MODUL — pecahannya ada, dan dibangkitkan dari sumber '
         'yang SAMA')
m.append('>')
m.append('> | Modul | Berkas | Lembar |')
m.append('> |---|---|---|')
m.append('> | **Treaty In** | [`../STRUKTUR-DATA-TREATY-IN.md`](../STRUKTUR-DATA-TREATY-IN.md) | '
         '*Treaty In Prop* · *Treaty In Non Prop* |')
m.append('> | **Treaty In Adjustment** | '
         '[`../../treaty-in-adjustment/STRUKTUR-DATA-TREATY-IN-ADJUSTMENT.md`]'
         '(../../treaty-in-adjustment/STRUKTUR-DATA-TREATY-IN-ADJUSTMENT.md) | '
         '*Treaty In EDM Prop* · *Treaty In EDM Non Prop* |')
m.append('>')
m.append('> Keduanya dibangkitkan [`alat/pisah-struktur-erd-v2-per-modul.py`]'
         '(../alat/pisah-struktur-erd-v2-per-modul.py) dari `alat/erd-v2.json` yang sama. '
         '**34 dari 37 tabel dipakai kedua modul** dan karena itu muncul di kedua pecahan — '
         'itu kenyataan sistem lama, bukan penggandaan.')
m.append('')
m.append('---')
m.append('')
m.append('## 1. Hitungan')
m.append('')
m.append('| | Jumlah |')
m.append('|---|---:|')
m.append('| Tabel unik di keempat lembar | **%d** |' % len(tabel))
m.append('| Kotak tergambar (satu tabel dapat muncul di beberapa lembar) | **%d** |' % n_kotak)
m.append('| Relasi di lembar *Daftar Relasi* | **%d** |' % len(relasi))
m.append('| Tabel yang dipakai **bersama** lebih dari satu lembar | **%d** |'
         % sum(1 for n in di_lembar if len(di_lembar[n]) > 1))
for nama, L in lembar.items():
    m.append('| Lembar *%s* | %d kotak |' % (nama, len(L['kotak'])))
m.append('')
m.append('## 2. Entitas — satu baris per tabel unik')
m.append('')
m.append('| Tabel | Kard. | Kunci utama | Kunci asing | ON DELETE | Induk | Muncul di lembar |')
m.append('|---|---|---|---|---|---|---|')
for nm, k in tabel.items():
    f = k['fk'] or {}
    m.append('| `%s` | %s | %s | %s | %s | %s | %s |'
             % (nm, k['kard'] or '—',
                '`%s`' % k['pk'] if k['pk'] else '*tidak dinyatakan*',
                '`%s`' % f['kolom'] if f else '—',
                f.get('on_delete', '—') or '—',
                '`%s`' % f['induk'] if f else '**akar**',
                ', '.join(x.replace('Treaty In ', '') for x in di_lembar[nm])))
m.append('')
m.append('> **Kolom *Kunci utama* berbunyi *tidak dinyatakan* untuk %d dari %d tabel, dan itu '
         'salinan setia.** Berkas sumbernya mencetak **PK hanya pada kotak akar**; pada kotak anak '
         'baris keduanya dipakai untuk **FK**. Maka kekosongan ini berarti **sumbernya diam**, '
         'bukan **tabelnya tanpa kunci utama** — lembar *Daftar Relasi* sumbernya sendiri merujuk '
         '`‹TABEL›.ID` sebagai sasaran tiap kunci asing.'
         % (sum(1 for k in tabel.values() if not k['pk']), len(tabel)))
m.append('')
m.append('## 3. Asal di sistem lama dan padanannya di model baru')
m.append('')
m.append('| Tabel | Jalur Pega | Padanan di model baru | Tingkat bukti |')
m.append('|---|---|---|---|')
for nm, k in tabel.items():
    m.append('| `%s` | %s | %s | %s |'
             % (nm, ('`%s`' % k['jalur']) if k['jalur'] else '*(tidak ada)*',
                ('**%s**' % k['padanan']) if k['padanan'] else '*(tidak ada)*',
                k['bukti'] or '—'))
m.append('')
m.append('## 4. Pohon per lembar')
m.append('')


def pohon_md(simpul, out, dalam=0):
    for n in simpul:
        f = n['fk'] or {}
        kunci = ('PK `%s`' % n['pk']) if n['pk'] else \
                ('FK `%s` → `%s`.`%s` · %s' % (f.get('kolom'), f.get('induk'),
                                               f.get('kolom_induk'), f.get('on_delete')))
        out.append('%s- **`%s`**%s — %s' % ('  ' * dalam, n['nama'],
                                            ' *%s*' % n['kard'] if n['kard'] else '', kunci))
        pohon_md(n['anak'], out, dalam + 1)


for nama, L in lembar.items():
    m.append('### %s — %d kotak' % (nama, len(L['kotak'])))
    m.append('')
    pohon_md(pohon(L['kotak']), m)
    m.append('')
m.append('## 5. Catatan & Batas — disalin apa adanya dari sumbernya')
m.append('')
m.append('| | |')
m.append('|---|---|')
for a, b in catatan:
    if a:
        m.append('| **%s** | %s |' % (a, b.replace('|', '\\|')))
m.append('')
io.open(MD_STRUKTUR, 'w', encoding='utf-8').write('\n'.join(m) + '\n')

# ═══════════════════════════ TABEL-DATAR-ERD-V2.md ═════════════════════════
t = []
t.append('# Tabel datar — **dibangkitkan dari `%s`**' % d['sumber'])
t.append('')
t.append('> ### %s' % d['banner'][0])
t.append('>')
t.append('> %s' % d['banner'][1])
t.append('')
t.append('**Dibangkitkan:** [`alat/bangkitkan-dari-erd-v2.py`](../alat/bangkitkan-dari-erd-v2.py) '
         '· **Sumber tunggal:** `alat/erd-v2.json` — sama dengan '
         '[`STRUKTUR-DATA-ERD-V2.md`](STRUKTUR-DATA-ERD-V2.md) dan '
         '[`ERD-TREATY-IN-DAN-EDM.html`](ERD-TREATY-IN-DAN-EDM.html)')
t.append('')
t.append('> **Ketiganya tidak dapat berbeda isinya**, sebab ketiganya dibangkitkan dari satu '
         'berkas yang sama. Bila berbeda, **alatnya yang salah**.')
t.append('')
t.append('---')
t.append('')
t.append('## 1. Satu baris per tabel datar')
t.append('')
t.append('| # | Tabel datar | Satu baris mewakili | Kunci | Induk | Kunci tamu | ON DELETE |')
t.append('|---:|---|---|---|---|---|---|')
for i, (nm, k) in enumerate(tabel.items(), 1):
    f = k['fk'] or {}
    wakil = k['jalur'] or '*(tidak ada jalur Pega — lihat §3)*'
    t.append('| %d | `%s` | `%s` | %s | %s | %s | %s |'
             % (i, nm, wakil,
                '`%s`' % k['pk'] if k['pk'] else '*tidak dinyatakan*',
                '`%s`' % f['induk'] if f else '**akar**',
                '`%s`' % f['kolom'] if f else '—',
                f.get('on_delete', '—') or '—'))
t.append('')
t.append('## 2. Cara mengisinya — asal tiap baris')
t.append('')
t.append('| Tabel datar | Diisi dari jalur Pega | Padanan model baru | Tingkat bukti |')
t.append('|---|---|---|---|')
for nm, k in tabel.items():
    t.append('| `%s` | %s | %s | %s |'
             % (nm, ('`%s`' % k['jalur']) if k['jalur'] else '*(tidak ada)*',
                ('**%s**' % k['padanan']) if k['padanan'] else '*(tidak ada)*',
                k['bukti'] or '—'))
t.append('')
t.append('## 3. Batas — dibaca sebelum memakai daftar ini')
t.append('')
# Sebuah tabel dapat membawa KEDUA tanda -- baris buktinya berbunyi
# "DAUN-RELATIF/JALUR-PENUH". Karena itu keempat golongan dibuat SALING LEPAS,
# dan jumlahnya diperiksa harus sama dengan cacah tabel.
gol = Counter()
for k in tabel.values():
    b = k['bukti'] or ''
    p, l = 'JALUR-PENUH' in b, 'DAUN-RELATIF' in b
    gol['keduanya' if (p and l) else 'kuat' if p else 'lemah' if l else 'tanpa'] += 1
if sum(gol.values()) != len(tabel):
    tolak['golongan bukti yang tidak menjumlah ke cacah tabel'] += 1
t.append('| | |')
t.append('|---|---|')
t.append('| Bukti **kuat** saja — `JALUR-PENUH` | **%d** tabel |' % gol['kuat'])
t.append('| Bukti **lemah** saja — `DAUN-RELATIF` | **%d** tabel |' % gol['lemah'])
t.append('| Membawa **kedua** tanda | **%d** tabel |' % gol['keduanya'])
t.append('| Tanpa tingkat bukti tertulis | **%d** tabel |' % gol['tanpa'])
t.append('| **Jumlah** | **%d** |' % sum(gol.values()))
t.append('')
t.append('> **Keempat golongan saling lepas, dan jumlahnya diperiksa sama dengan cacah tabel.** '
         'Sebuah tabel dapat membawa **kedua** tanda — baris buktinya berbunyi '
         '`DAUN-RELATIF/JALUR-PENUH` — sehingga menjumlahkan "kuat" dan "lemah" begitu saja '
         'menghasilkan angka yang melampaui cacah tabelnya.')
t.append('')
t.append('> **`DAUN-RELATIF` lemah, dan sebabnya disebut di sumbernya:** hanya nama ruas terakhir '
         'yang ditemukan — misalnya `.COBList` — sehingga ia **mungkin memungut nama milik '
         'entitas lain**.')
t.append('')
t.append('> **Nol rujukan berarti BELUM TERPERIKSA, bukan tidak dipakai.** Sapuan yang '
         'menghasilkan daftar ini hanya membaca **Section**. Activity, Data Transform, Report '
         'Definition, dan RDB List **tidak disapu**.')
t.append('')
t.append('> **`ON DELETE CASCADE` adalah USULAN rancangan, bukan perilaku sistem lama.** Ia '
         'mengikuti konvensi `contooh.xlsx`.')
t.append('')
io.open(MD_DATAR, 'w', encoding='utf-8').write('\n'.join(t) + '\n')

# ══════════════════════════ CETAK ══════════════════════════════════════════
P = print
P('=== bangkitkan-dari-erd-v2.py ===')
P()
P('SUMBER TUNGGAL : alat/erd-v2.json  (dari %s)' % d['sumber'])
P()
P('DITULIS:')
for p in (HTML, MD_STRUKTUR, MD_DATAR):
    P('   %-34s %7d bita' % (os.path.basename(p), os.path.getsize(p)))
P()
P('ISI, dan ketiga berkas memuat angka yang sama:')
P('   tabel unik                 : %d' % len(tabel))
P('   kotak di empat lembar      : %d' % n_kotak)
for nama, L in lembar.items():
    P('      %-24s %d' % (nama, len(L['kotak'])))
P('   relasi                     : %d' % len(relasi))
P('   dipakai bersama > 1 lembar : %d' % sum(1 for n in di_lembar if len(di_lembar[n]) > 1))
P('   golongan bukti (saling lepas)    : kuat %d · lemah %d · keduanya %d · tanpa %d'
  % (gol['kuat'], gol['lemah'], gol['keduanya'], gol['tanpa']))
P()
P('DITOLAK / TIDAK DIGAMBAR -- satu baris per sebab:')
for k, v in tolak.most_common():
    P('   %-64s %d' % (k, v))
if not tolak:
    P('   (tidak ada)')
P()
P('-- BATAS PERKAKAS INI --')
P('   Ia MENYALIN sumbernya, ia tidak menilai isinya. Sumbernya berbanner')
P('   POTRET SISTEM LAMA, dan ketiga keluarannya mengatakan hal yang sama di kepalanya.')
P('   Satu hal ditambahkan dan ia diminta: GARIS PENGHUBUNG di HTML, yang di')
P('   berkas xlsx-nya dijanjikan legenda tetapi TIDAK ADA.')
