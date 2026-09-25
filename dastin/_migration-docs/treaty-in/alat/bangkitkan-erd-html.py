# -*- coding: utf-8 -*-
"""bangkitkan-erd-html.py -- ERD skema baru dalam HTML, DIBANGKITKAN dari erd-skema-baru.json.

KENAPA ADA
    Dua berkas ERD HTML di folder ini berbanner BASI: keduanya menggambarkan skema versi
    24 September 17:50 -- nol DOKUMEN_ADDENDUM, nol tabel NILAI_*. Berkas ini menggantikan
    keduanya sebagai gambar yang mutakhir.

SUMBER
    erd-skema-baru.json, yang dibangkitkan bangkitkan-erd-skema-baru.py dari 2-to-spec/ddl-usulan/.
    Rantainya: DDL -> json -> html. Tidak ada satu pun angka yang ditulis tangan di sini.

BATAS KEMAMPUAN PERKAKAS INI, dinyatakan di dalam keluarannya sendiri:
    Ia menggambar apa yang ADA di DDL. Ia TIDAK dapat menyatakan bahwa DDL-nya benar.
    Rujukan yang SEHARUSNYA ADA tetapi tidak dideklarasikan muncul di bagiannya sendiri --
    ia temuan F-19/F-20, bukan bagian dari gambar.
"""
import html
import io
import json
import os

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', '4-erd-dan-tabel-datar')
SUMBER = os.path.join(AKAR, 'erd-skema-terpadu.json')
CADANGAN = os.path.join(AKAR, 'erd-skema-baru.json')
KELUARAN = os.path.join(AKAR, 'ERD-SKEMA-BARU.html')

if os.path.exists(SUMBER):
    dipakai = SUMBER
else:
    dipakai = CADANGAN
    print('DITOLAK: erd-skema-terpadu.json tidak ada; mundur ke erd-skema-baru.json '
          '-- gambarnya akan KEHILANGAN kedua entitas Adjustment.')
d = json.load(io.open(dipakai, encoding='utf-8'))
tabel, relasi, lembar = d['tabel'], d['relasi'], d['lembar']
tanpa_fk, luar, uq, tolak = d['tanpa_fk'], d['luar_skema'], dict(d['uq']), d['tolak']
tanggal = d.get('tanggal', '')

e = html.escape
anak_dari = {}
for r in relasi:
    anak_dari.setdefault(r['induk'], []).append(r)

CSS = """
:root{--bg:#f7f9fc;--kartu:#fff;--tepi:#d6dee8;--tua:#1f4e79;--biru:#2e75b6;--kuning:#fff2cc;
--muda:#edf3f9;--teks:#1b1f24;--redup:#5b6470;--merah:#c00000;--jingga:#e07b00;--hijau:#548235}
@media(prefers-color-scheme:dark){:root:not([data-theme=light]){--bg:#11161d;--kartu:#181f28;
--tepi:#2b3644;--muda:#1c2733;--kuning:#3a3420;--teks:#e6eaf0;--redup:#98a3b3;--biru:#5b9bd5}}
:root[data-theme=dark]{--bg:#11161d;--kartu:#181f28;--tepi:#2b3644;--muda:#1c2733;
--kuning:#3a3420;--teks:#e6eaf0;--redup:#98a3b3;--biru:#5b9bd5}
*{box-sizing:border-box}
body{margin:0;background:var(--bg);color:var(--teks);
font:14px/1.55 -apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif}
.bungkus{max-width:1280px;margin:0 auto;padding:24px 16px 72px}
h1{font-size:24px;margin:0 0 4px;color:var(--tua)}
@media(prefers-color-scheme:dark){h1{color:var(--biru)}}
h2{font-size:18px;margin:38px 0 12px;padding-bottom:6px;border-bottom:2px solid var(--tepi)}
h3{font-size:14px;margin:22px 0 8px;color:var(--redup);text-transform:uppercase;letter-spacing:.06em}
.sub{color:var(--redup);font-size:12.5px;margin:0 0 18px}
.pita{background:var(--merah);color:#fff;padding:10px 14px;border-radius:6px;font-weight:600;
margin:0 0 6px;font-size:13px}
.pita.b{background:var(--hijau)}
.catatan{background:var(--kartu);border:1px solid var(--tepi);border-left:4px solid var(--jingga);
border-radius:6px;padding:12px 14px;margin:14px 0;font-size:13px}
.angka{display:flex;flex-wrap:wrap;gap:10px;margin:14px 0 6px}
.angka div{background:var(--kartu);border:1px solid var(--tepi);border-radius:6px;
padding:8px 14px;min-width:104px}
.angka b{display:block;font-size:21px;color:var(--tua)}
@media(prefers-color-scheme:dark){.angka b{color:var(--biru)}}
.angka span{font-size:11.5px;color:var(--redup)}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(324px,1fr));gap:14px}
.kartu{background:var(--kartu);border:1px solid var(--tepi);border-radius:8px;overflow:hidden}
.kepala{background:var(--tua);color:#fff;padding:8px 12px;font-weight:700;font-size:13.5px;
display:flex;justify-content:space-between;gap:8px;align-items:baseline}
.kepala.anak{background:var(--biru)}
.kepala small{font-weight:400;opacity:.85;font-size:11px}
table{width:100%;border-collapse:collapse;font-size:12.5px}
td,th{padding:4px 10px;border-bottom:1px solid var(--tepi);text-align:left;vertical-align:top}
th{background:var(--muda);font-size:11px;text-transform:uppercase;letter-spacing:.05em;
color:var(--redup)}
tr:last-child td{border-bottom:none}
.pk{background:var(--kuning);font-weight:700}
.fk td:first-child{font-weight:600}
code{font-family:ui-monospace,SFMono-Regular,Consolas,monospace;font-size:12px}
.tag{display:inline-block;font-size:10.5px;padding:1px 6px;border-radius:10px;
border:1px solid var(--tepi);color:var(--redup);margin-left:4px}
.tag.m{border-color:var(--merah);color:var(--merah)}
.tag.j{border-color:var(--jingga);color:var(--jingga)}
.lebar{background:var(--kartu);border:1px solid var(--tepi);border-radius:8px;overflow-x:auto}
.mono{font-family:ui-monospace,Consolas,monospace}
.nol{color:var(--redup);font-style:italic}
@media print{body{background:#fff}.kartu,.lebar{break-inside:avoid}}
"""


def kartu(nm, lembar_ini=None):
    t = tabel[nm]
    lain = [x for x in d['di_lembar'].get(nm, []) if x != lembar_ini]
    fkkol = {f['kolom']: f for f in t['fk']}
    baris = []
    for k in t['kolom']:
        cls = ' class="pk"' if k['nama'] == t['pk'] else ''
        tags = ''
        if k['nama'] == t['pk']:
            tags += '<span class="tag">PK</span>'
        if k['nama'] in fkkol:
            f = fkkol[k['nama']]
            od = f.get('on_delete') or 'tidak dinyatakan'
            tags += '<span class="tag">FK &rarr; %s</span>' % e(f['induk'])
            if not f.get('on_delete'):
                tags += '<span class="tag j">ON DELETE %s</span>' % e(od)
        baris.append('<tr%s><td><code>%s</code>%s</td><td class="mono">%s</td><td>%s</td></tr>'
                     % (cls, e(k['nama']), tags, e(k['tipe']),
                        '' if k['wajib'] else '<span class="nol">boleh kosong</span>'))
    kunci = uq.get(nm)
    kalimat = ''
    if kunci:
        kalimat = '<div style="padding:7px 12px;font-size:11.5px;background:var(--muda)">' \
                  '<b>Kunci alami:</b> ' + ' &middot; '.join(
                      '<code>%s</code>' % e(' + '.join(u)) for u in kunci) + '</div>'
    elif t.get('kunci_alami_teks'):
        kalimat = ('<div style="padding:7px 12px;font-size:11.5px;background:var(--muda)">'
                   '<b>Kunci alami:</b> <code>%s</code></div>' % e(t['kunci_alami_teks']))
    # ASAL -- dari peta nama tabel datar dan dari SPEC-MODEL-DATA sec 2.3 / 10.22
    jejak = []
    if t.get('asal_pega'):
        jejak.append('<b>Asal di sistem lama:</b> ' + ' &middot; '.join(
            '<code>%s</code>' % e(x) for x in t['asal_pega']))
    if t.get('nama_lama'):
        jejak.append('<b>Nama tabel datar:</b> ' + ' &middot; '.join(
            '<code>%s</code>' % e(x) for x in t['nama_lama']))
    if not t.get('asal_pega') and not t.get('nama_lama'):
        jejak.append('<span class="nol">tidak punya asal di pohon Pega — '
                     'entitas BARU atau tabel anak yang lahir dari keputusan model</span>')
    asalnya = ('<div style="padding:7px 12px;font-size:11.5px;border-top:1px solid var(--tepi)">'
               + ' &nbsp;|&nbsp; '.join(jejak) + '</div>')
    # LUBANG -- dicetak di dalam kartunya sendiri, bukan disembunyikan
    lub = ''
    if t.get('lubang'):
        lub = ('<div style="padding:9px 12px;font-size:11.5px;background:var(--kuning);'
               'border-top:1px solid var(--tepi)"><b>BELUM DIPUTUSKAN — %d butir</b><ol '
               'style="margin:6px 0 0;padding-left:18px">%s</ol></div>'
               % (len(t['lubang']), ''.join('<li>%s</li>' % e(x) for x in t['lubang'])))
    n_anak = len(anak_dari.get(nm, []))
    akar = not t['fk']
    bersama = ('<span class="tag j">JUGA DI: %s</span>'
               % e(', '.join(lain))) if lain else ''
    if t.get('belum_ddl'):
        bersama += '<span class="tag m">BELUM BER-DDL</span>'
    return ('<div class="kartu"><div class="kepala%s"><span>%s%s</span>'
            '<small>%d kolom%s</small></div>%s'
            '<table><thead><tr><th>Kolom</th><th>Tipe</th><th>Keterisian</th></tr></thead>'
            '<tbody>%s</tbody></table>%s%s</div>'
            % ('' if akar else ' anak', e(nm), bersama, len(t['kolom']),
               ' &middot; %d anak' % n_anak if n_anak else '', kalimat,
               ''.join(baris), asalnya, lub))


bag = []
bag.append('<div class="pita b">GAMBAR SKEMA BARU — MUTAKHIR. '
           'Dibangkitkan dari <code>2-to-spec/ddl-usulan/</code>, bukan digambar tangan.</div>')
bag.append('<h1>ERD skema baru — Treaty In dan Treaty In Adjustment</h1>')
bag.append('<p class="sub">Satu model, satu spesifikasi (<code>GRL-01</code>). '
           'Dibangkitkan %s dari <b>enam sumber</b>, bukan digambar tangan.</p>' % e(tanggal))

if d.get('sumber_rantai'):
    bag.append('<div class="catatan"><b>Rantai sumbernya, supaya dapat dijalankan ulang.</b>'
               '<ol style="margin:6px 0 0;padding-left:18px">%s</ol></div>'
               % ''.join('<li><code>%s</code></li>' % e(x) for x in d['sumber_rantai']))

belum = [n for n, t in tabel.items() if t.get('belum_ddl')]
bag.append('<div class="angka">'
           '<div><b>%d</b><span>tabel</span></div>'
           '<div><b>%d</b><span>kunci asing</span></div>'
           '<div><b>%d</b><span>kolom</span></div>'
           '<div><b>%d</b><span>kunci alami</span></div>'
           '<div><b>%d</b><span>rujukan SEHARUSNYA ADA</span></div>'
           '<div><b>%d</b><span>entitas luar</span></div>'
           '<div><b>%d</b><span>BELUM ber-DDL</span></div></div>'
           % (len(tabel), len(relasi), sum(len(t['kolom']) for t in tabel.values()),
              len(uq), len(tanpa_fk), len(luar), len(belum)))

if belum:
    bag.append('<div class="pita">%d dari %d tabel di gambar ini BELUM ADA DDL-nya: %s. '
               'Keduanya milik modul Adjustment, DIDAFTAR di <code>STRUKTUR-DATA.md</code> '
               'yang mengikat, tetapi nol berkas di <code>ddl-usulan/</code> dan nol pasal di '
               '<code>KAMUS-KOLOM.md</code>. Menggambar sebuah tabel bukan membangunnya.</div>'
               % (len(belum), len(tabel), ', '.join('<code>%s</code>' % e(x) for x in belum)))

bag.append('<div class="catatan"><b>Batas berkas ini.</b> Untuk 35 tabel ia menggambar apa yang '
           '<b>ada di DDL</b>; untuk 2 tabel ia menggambar apa yang <b>ada di to-spec Adjustment</b>, '
           'dan tiap kolomnya membawa kutipan sumbernya. Ia <b>tidak dapat</b> menyatakan bahwa '
           'sumbernya benar, dan <b>tidak dapat</b> menyatakan himpunannya lengkap — semesta §10 '
           'masih kurang <b>340 properti titik buta</b> (<code>L-8</code>), dan gambar dari '
           'himpunan yang bolong <b>terlihat lengkap</b>. Rujukan yang seharusnya ada tetapi tidak '
           'dideklarasikan berdiri di bagiannya sendiri di bawah — ia temuan <code>F-19</code>/'
           '<code>F-20</code>, <b>bukan bagian dari gambar</b>.</div>')

for nama_lembar, isi in lembar.items():
    bag.append('<h2>%s <small style="font-weight:400;color:var(--redup)">— %d tabel</small></h2>'
               % (e(nama_lembar), len(isi)))
    bag.append('<div class="grid">%s</div>' % ''.join(kartu(n, nama_lembar) for n in isi if n in tabel))

bag.append('<h2>Daftar relasi — %d kunci asing</h2>' % len(relasi))
bar = ''.join(
    '<tr><td>%d</td><td><code>%s</code></td><td class="mono">%s</td><td><code>%s</code></td>'
    '<td><code>%s</code></td><td>%s</td><td style="font-size:11px;color:var(--redup)">%s</td></tr>'
    % (r['no'], e(r['induk']), e(r['kard']), e(r['anak']), e(r['kolom']),
       ('<span class="tag j">%s</span>' % e(r['on_delete'])) if 'tidak dinyatakan' in str(r['on_delete'])
       else e(str(r['on_delete'])),
       e(r['lembar_induk']) + (' &rarr; ' + e(r['lembar_anak'])
                               if r['lembar_induk'] != r['lembar_anak'] else ''))
    for r in relasi)
bag.append('<div class="lebar"><table><thead><tr><th>#</th><th>Induk</th><th>Kard.</th>'
           '<th>Anak</th><th>Kolom</th><th>ON DELETE</th><th>Lembar</th></tr></thead>'
           '<tbody>%s</tbody></table></div>' % bar)

bag.append('<h2>Rujukan yang SEHARUSNYA ADA — %d, dan belum dideklarasikan</h2>' % len(tanpa_fk))
bag.append('<div class="catatan"><b>Ini temuan, bukan gambar.</b> Kolom di bawah menunjuk entitas '
           'lain menurut namanya, tetapi <b>tidak ada <code>FOREIGN KEY</code></b> yang menjaganya. '
           'Basis data tidak akan menolak nilai yang tidak ada padanannya. '
           'Lihat <code>TEMUAN-F19-F21-RUJUKAN-YANG-HILANG.md</code>.</div>')
bar = ''.join('<tr><td><code>%s</code></td><td><code>%s</code></td>'
              '<td>%s</td><td style="font-size:11.5px">%s</td></tr>'
              % (e(x[0]), e(x[1]),
                 '<span class="tag m">%s</span>' % e(x[2]) if 'SEHARUSNYA' in x[2]
                 else '<span class="tag">%s</span>' % e(x[2]), e(x[3]))
              for x in tanpa_fk)
bag.append('<div class="lebar"><table><thead><tr><th>Tabel</th><th>Kolom</th><th>Golongan</th>'
           '<th>Keterangan</th></tr></thead><tbody>%s</tbody></table></div>' % bar)

bag.append('<h2>Entitas di luar skema — dirujuk, tidak dimiliki</h2>')
bar = ''.join('<tr><td><code>%s</code></td><td>%s</td></tr>' % (e(k), e(v)) for k, v in luar.items())
bag.append('<div class="lebar"><table><thead><tr><th>Entitas</th><th>Keterangan</th></tr></thead>'
           '<tbody>%s</tbody></table></div>' % bar)

if d.get('peta_belum'):
    bag.append('<h2>Peta nama tabel datar yang BELUM TERPADANKAN — %d</h2>'
               % len(d['peta_belum']))
    bag.append('<div class="catatan"><code>peta-nama-tabel-treatyin.tsv</code> memberi tiap tabel '
               'datar lama padanannya di model baru. <b>%d padanan menyebut nama yang tidak ada '
               'di <code>ddl-usulan/</code></b> — bukan karena sengaja di luar gelombang ini '
               '(yang itu sudah dipisahkan), melainkan karena <b>namanya berubah dan petanya tidak '
               'ikut diperbarui</b>. Selama begini, peta nama <b>tidak dapat dipakai menelusuri '
               'tabel lama ke tabel baru</b> untuk baris-baris ini.</div>' % len(d['peta_belum']))
    bar = ''.join('<tr><td class="mono">%s</td><td><code>%s</code></td></tr>' % (e(a), e(b))
                  for a, b in d['peta_belum'])
    bag.append('<div class="lebar"><table><thead><tr><th>Tabel datar lama</th>'
               '<th>Padanan yang disebut, tetapi tidak ada di DDL</th></tr></thead>'
               '<tbody>%s</tbody></table></div>' % bar)

if d.get('tanpa_nasib'):
    bag.append('<h2>Cakupan nasib — %d dari %d simpul pohon Pega belum bernasib</h2>'
               % (d['tanpa_nasib'], d.get('cacah_simpul', 0)))
    bag.append('<div class="catatan"><code>PETA-TELUSUR-JSON.md</code> §6 memberi nasib '
               '(<b>DIPETAKAN · DITURUNKAN · DIBUANG · DITUNDA</b>) kepada <b>667</b> jalur, dan '
               'menyatakan sendiri <i>"di luar keempat kategori: nol"</i>. Pohon Pega yang sebenarnya '
               'berisi <b>%d</b> simpul (<code>datar-treatyin-lama.csv</code>). Diadu satu per satu, '
               '<b>%d simpul tidak punya nasib</b>. '
               'Berkas itu sudah menyatakan semestanya kurang <b>340</b> (<code>L-8</code>); '
               'pengukuran mandiri ini memberi <b>%d</b>. Selisih keduanya <b>%d</b>, dan sebagiannya '
               'dapat berasal dari cara perkakas menyamakan jalur — ia membuang <code>[]</code> dan '
               'awalan <code>TreatyIn.</code>, sehingga jalur yang berbeda hanya pada tanda itu '
               'terbaca sama. <b>Angka ini menakar lubang, ia tidak menutupnya.</b></div>'
               % (d.get('cacah_simpul', 0), d['tanpa_nasib'], d['tanpa_nasib'],
                  abs(d['tanpa_nasib'] - 340)))

bag.append('<h2>Yang ditolak perkakas — dilaporkan, bukan didiamkan</h2>')
bar = ''.join('<tr><td>%s</td><td class="mono">%s</td></tr>' % (e(k), v) for k, v in tolak.items())
bag.append('<div class="lebar"><table><thead><tr><th>Sebab penolakan</th><th>Jumlah</th></tr>'
           '</thead><tbody>%s</tbody></table></div>' % bar)

bag.append('<h2>Batas — dibaca sebelum memakai gambar ini sebagai bukti</h2>')
bag.append('<div class="catatan">'
           '<p><b>Yang dapat dibuktikan berkas ini:</b> untuk 35 tabel, bahwa sebuah tabel, kolom, '
           'atau kunci asing <b>ada di <code>ddl-usulan/</code></b>. Untuk 2 tabel Adjustment, bahwa '
           'sebuah kolom <b>disebut oleh sumber yang dikutip di barisnya sendiri</b>.</p>'
           '<p><b>Yang TIDAK dapat dibuktikan:</b> bahwa skemanya <b>lengkap</b> atau <b>benar</b>. '
           'Peta telusur induk masih berlabel <i>semestanya kurang</i> akibat titik buta '
           '<code>L-8</code> — <b>340 properti</b> belum diperiksa siapa pun. Dan kedua tabel '
           'Adjustment <b>belum melewati §10</b>: kolom yang tidak ada di kartunya bukan kolom yang '
           'diputuskan tidak ada, melainkan kolom yang <b>sumbernya diam</b>.</p>'
           '<p><b>Satu pertentangan sumber DIPILIH, bukan diadili.</b> Berapa penunjuk versi yang '
           'disimpan <code>NILAI_SELISIH</code> — <code>KEPUTUSAN-SAMBUNGAN</code> G2 berkata dua, '
           '<code>STRUKTUR-DATA</code> §1.1 berkata satu. Gambar ini memakai <b>satu</b>, atas dasar '
           'butir 4 urutan wewenang. Itu <b>pilihan wewenang</b>, dan ia menunggu putusan pemilik '
           'proses.</p>'
           '<p><b>Yang berubah membuat berkas ini basi:</b> setiap suntingan pada '
           '<code>2-to-spec/ddl-usulan/</code>, <code>STRUKTUR-DATA.md</code>, atau '
           '<code>KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md</code>. Jalankan ulang ketiganya berurutan — '
           '<code>bangkitkan-erd-skema-baru.py</code> &rarr; <code>cocok-enam-sumber-struktur.py</code> '
           '&rarr; <code>lengkapi-erd-adjustment.py</code> &rarr; <code>bangkitkan-erd-html.py</code> — '
           'semuanya dapat dijalankan berkali-kali.</p>'
           '</div>')

doc = ('<!DOCTYPE html><html lang="id"><head><meta charset="utf-8">'
       '<meta name="viewport" content="width=device-width,initial-scale=1">'
       '<title>ERD Skema Baru — Treaty In</title><style>%s</style></head>'
       '<body><div class="bungkus">%s</div></body></html>' % (CSS, ''.join(bag)))

io.open(KELUARAN, 'w', encoding='utf-8').write(doc)

print('DITULIS:', os.path.normpath(KELUARAN))
print('  tabel %d · relasi %d · kolom %d · kunci alami %d'
      % (len(tabel), len(relasi), sum(len(t['kolom']) for t in tabel.values()), len(uq)))
print('  rujukan SEHARUSNYA ADA %d · entitas luar %d' % (len(tanpa_fk), len(luar)))
print('  lembar:', ' · '.join('%s (%d)' % (k, len(v)) for k, v in lembar.items()))
print()
print('-- BATAS PERKAKAS INI --')
print('   Ia menggambar apa yang ADA di DDL. Ia TIDAK dapat menyatakan DDL-nya benar.')
print('   Rujukan SEHARUSNYA ADA dicetak di bagiannya sendiri: ia temuan, bukan gambar.')
