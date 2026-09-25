# -*- coding: utf-8 -*-
r"""
buat-erd-struktur-treatyin.py — ERD tabel datar T_… , DIBANGKITKAN DARI STRUKTUR.

Bedanya dengan `buat-skema-treaty-masuk.py`: alat itu menggambar ENTITAS RANCANGAN
(nama §16, 31 kotak) dari definisi yang ditulis tangan. Alat INI menggambar TABEL DATAR
(kepala TREATY_IN + anak T_…) dan **seluruh isi kotaknya dibaca dari hasil sapuan XML** —
nama field, kedalaman, dan cacah rujukan tidak diketik ulang.

Masukan (keduanya turunan, harus dibangkitkan lebih dulu):
  4-erd-dan-tabel-datar/peta-nama-tabel-treatyin.tsv   <- buat-peta-nama-tabel-treatyin.py
  4-erd-dan-tabel-datar/datar-treatyin-lama.csv        <- buat-pohon-treatyin.py

Keluaran:
  4-erd-dan-tabel-datar/ERD-STRUKTUR-TREATYIN.html

Jalankan:  PYTHONIOENCODING=utf-8 python alat/buat-erd-struktur-treatyin.py
"""
import os, io, csv, collections

KELUARAN = os.path.join(r'D:\XML_NURE', '_migration-docs', 'treaty-in',
                        '4-erd-dan-tabel-datar')
P_TSV = os.path.join(KELUARAN, 'peta-nama-tabel-treatyin.tsv')
P_CSV = os.path.join(KELUARAN, 'datar-treatyin-lama.csv')

WARNA = {
    'akar':    ('1F4E79', 'AKAR \u2014 kepalanya memakai tabel TREATY_IN yang SUDAH ADA'),
    'limit':   ('2E75B6', 'LAYER \u2014 Limits[] dan rincian proporsionalnya'),
    'share':   ('31859C', 'BAGIAN \u2014 Share[] dan penyebarannya'),
    'fac':     ('C55A11', 'FAKULTATIF \u2014 FacultativeShareList[] dan reasuradurnya'),
    'kepala':  ('7030A0', 'DAFTAR TINGKAT KONTRAK \u2014 menggantung langsung pada kontrak'),
    'bersama': ('548235', 'DIPAKAI BERSAMA modul lain'),
    'edm':     ('C00000', 'ADDENDUM DAN SELISIH \u2014 melintasi sekat ke modul Adjustment'),
    'migrasi': ('808080', 'JEMBATAN MIGRASI \u2014 tidak punya asal di pohon Pega'),
    'lubang':  ('BF8F00', 'BELUM ADA ISINYA \u2014 digambar sebagai lubang, bukan sebagai kotak'),
}
URUT = ['akar', 'limit', 'share', 'fac', 'kepala', 'bersama', 'edm', 'migrasi', 'lubang']

KREM, PUCAT, PITA_N = 'FFF2CC', 'EDF3F9', 'FCE4D6'
MAKS_FIELD = 9          # field yang digambar per kotak; sisanya diringkas

# Relasi yang MELINTASI SEKAT Treaty In <-> Adjustment. Didaftar SATU PER SATU,
# tidak disimpulkan dari kelompok: `ERD.md` menetapkan hanya dua yang melintas, dan
# kelompok `edm` juga memuat hal-hal yang tetap di dalam Treaty In
# (T_TREATY_REVISION, T_TREATY_VALUE_BEFORE_PRORATE).
LINTAS_SEKAT = {'T_TREATY_VALUE_DIFFERENCE'}


# ------------------------------------------------------------------ baca masukan
def baca():
    with io.open(P_TSV, encoding='utf-8') as h:
        tabel = [r for r in csv.DictReader(h, delimiter='\t')
                 if not r['NAMA_T'].startswith('\u2014')]
    simpul, anak = {}, collections.defaultdict(list)
    with io.open(P_CSV, encoding='utf-8-sig') as h:
        for r in csv.DictReader(h):
            simpul[r['JALUR']] = r
            if '.' in r['JALUR']:
                anak[r['JALUR'].rsplit('.', 1)[0]].append(r)
    return tabel, simpul, anak


TABEL, SIMPUL, ANAK = baca()
PETA = dict((t['NAMA_T'], t) for t in TABEL)


def jalur_list(t):
    j = t.get('JALUR_PEGA', '').strip()
    return [x.strip() for x in j.split('|') if x.strip()] if j else []


def skalar(jalur):
    return [r['NAMA'] for r in sorted(ANAK.get(jalur, []), key=lambda r: r['NAMA'])
            if r['JENIS'] == 'Skalar']


def kedalaman_pohon(t):
    """kedalaman simpul terdangkal yang ditambatkan di pohon Pega"""
    js = jalur_list(t)
    if not js:
        return None
    d = [int(SIMPUL[j]['KEDALAMAN']) for j in js if j in SIMPUL]
    return min(d) if d else None


def kedalaman(t):
    """Kedalaman gambar diambil dari RANTAI INDUK, bukan langsung dari jalur Pega.

    Sejak kepala memakai tabel yang sudah ada (`TREATY_IN`, menambat ke halaman
    `TreatyIn` di tingkat 0), rantai induk SAMA PERSIS dengan kedalaman pohon.
    `beda_kedalaman()` memeriksa itu, dan hanya empat tabel yang menyimpang —
    seluruhnya bersebab tertulis.
    """
    d, n = 0, t['NAMA_T']
    while PETA.get(n, {}).get('INDUK'):
        n = PETA[n]['INDUK']
        d += 1
        if d > 8:
            break
    return d


# Tabel yang rantai induknya SENGAJA lebih dalam dari pohonnya, di luar selisih +1.
# Didaftar supaya penyimpangan baru ketahuan, bukan tenggelam di antara yang wajar.
SELISIH_DISENGAJA = {
    'T_TREATY_REVISION':
        'skalar akar dipecah keluar — menambat ke TreatyIn (tingkat 0), '
        'digambar sebagai anak TREATY_IN',
    'T_TREATY_HAZARD_LIMIT':
        '10 skalar akar diputar jadi daftar — menambat ke TreatyIn (tingkat 0)',
    'T_TREATY_VALUE_DIFFERENCE':
        'induknya T_TREATY_REVISION yang sendirinya sudah dipecah keluar',
    # arah sebaliknya: tambatannya LEBIH DALAM dari tempat sebenarnya, karena
    # bentuk InstallmentList hanya terbaca di dalam salinan ValueDifference.
    # Di akarnya ia tingkat 2; lewat salinan ia terbaca di tingkat 3.
    'T_TREATY_INSTALLMENT_ITEM':
        'bentuknya HANYA terbaca lewat salinan ValueDifference — '
        'struktur-treatyin-lama.md §5.2',
}
# Kepala = TREATY_IN yang menambat ke `TreatyIn` (tingkat 0), jadi rantai induk
# dan kedalaman pohon berjalan seiring. Dulu nilainya 1 karena ada kotak
# T_WORK_TREATY_IN (pyWorkPage) di atasnya; kotak itu sudah tidak ada.
SELISIH_WAJAR = 0


def beda_kedalaman():
    """kembalikan (tak_terduga, disengaja) — keduanya dicetak, hanya yang pertama
    yang jadi alasan memeriksa ulang"""
    tak_terduga, disengaja = [], []
    for t in TABEL:
        dp = kedalaman_pohon(t)
        if dp is None:
            continue
        selisih = kedalaman(t) - dp
        if selisih == SELISIH_WAJAR:
            continue
        baris = (t['NAMA_T'], kedalaman(t), dp, selisih)
        (disengaja if t['NAMA_T'] in SELISIH_DISENGAJA else tak_terduga).append(baris)
    return tak_terduga, disengaja


def ref(t):
    return sum(int(SIMPUL[j]['REF']) for j in jalur_list(t) if j in SIMPUL)


# ------------------------------------------------------------------ kolom nyata
# `TREATY_IN` adalah tabel yang SUDAH ADA, dan kolomnya bukan tebakan: argumen
# prosedur POOLDATA.PEGA_TREATY_IN menyebutnya satu per satu. Dibaca, bukan diketik.
P_KELUAR = os.path.join(KELUARAN, 'datar-treatyin-muatan-keluar.csv')


def kolom_nyata(prosedur='PEGA_TREATY_IN'):
    if not os.path.exists(P_KELUAR):
        return []
    out = []
    with io.open(P_KELUAR, encoding='utf-8-sig') as h:
        for r in csv.DictReader(h):
            if r['PROSEDUR'] != prosedur:
                continue
            a = r['DIISI_DARI'].strip()
            if a.startswith('TreatyIn.') and '.' not in a[9:]:
                out.append(a[9:])
    return out


KOLOM_NYATA = kolom_nyata()
KEPALA_TABEL = 'TREATY_IN'

# Dua tabel MENGAMBIL SEBAGIAN skalar akar, bukan seluruhnya. Tanpa daftar ini
# keduanya menambat ke `TreatyIn` dan menumpahkan ke-99 skalarnya ke dalam kotak.
# Nama-nama di bawah DIPERIKSA keberadaannya di pohon oleh periksa_pilihan().
FIELD_PILIHAN = {
    'T_TREATY_REVISION': ['OLDID', 'EDMState', 'EDMEffective', 'EDMMaterialType',
                          'RevisionState', 'RevisionDate', 'AddendumPremi'],
    'T_TREATY_HAZARD_LIMIT': ['Earthquake', 'CurrencyEarthquake', 'FloodJab',
                              'CurrencyFloodJab', 'FloodNation', 'CurrencyFloodNat',
                              'RSMDLimit', 'CurrencyRSMD', 'MaxCoGroup',
                              'MaxCoNonGroup'],
}


def periksa_pilihan():
    akar = set(skalar('TreatyIn'))
    salah = [(t, f) for t, fs in FIELD_PILIHAN.items() for f in fs if f not in akar]
    if salah:
        raise SystemExit('FIELD_PILIHAN menyebut skalar yang TIDAK ADA di akar:\n  ' +
                         '\n  '.join('%s: %s' % x for x in salah))
    return sum(len(v) for v in FIELD_PILIHAN.values())


N_PILIHAN = periksa_pilihan()


# ------------------------------------------------------------------ susun kotak
def pk_dari(nama_t):
    return 'ID_' + nama_t[2:] if nama_t.startswith('T_') else 'ID_' + nama_t


def baris_kotak(t):
    """('pk'|'fk'|'f'|'x'|'n', teks)"""
    b = [('pk', 'PK   %s' % pk_dari(t['NAMA_T']))]
    if t['INDUK']:
        ku = t['KUNCI'] if t['KUNCI'] and not t['KUNCI'].startswith('\u2014') \
            else pk_dari(t['INDUK'])
        b.append(('fk', 'FK   %-22s \u2192  %s' % (ku, t['INDUK'])))

    js = jalur_list(t)

    # Kepala digambar khusus: kolom yang BENAR-BENAR ADA di tabel Oracle dipisahkan
    # dari skalar yang hanya hidup di M_TREATY_IN.JSONDATA. Menggabung keduanya akan
    # membuat 79 field tampak seolah sudah punya kolom.
    if t['NAMA_T'] == KEPALA_TABEL and KOLOM_NYATA:
        semua = set(skalar('TreatyIn'))
        nyata = [k for k in KOLOM_NYATA if k != 'ID']
        b.append(('x', '\u25b8    %d KOLOM NYATA di tabel TREATY_IN yang sudah ada'
                  % len(KOLOM_NYATA)))
        for nm in nyata:
            b.append(('f', '     ' + nm))
        sisa = len(semua - set(KOLOM_NYATA))
        b.append(('n', '     \u2500\u2500 di bawah ini TIDAK punya kolom \u2500\u2500'))
        b.append(('n', '     %d skalar akar lain hanya ada di M_TREATY_IN.JSONDATA' % sisa))
        b.append(('n', '     ViewState \u00b7 EDMMaterialType \u00b7 RNMShare \u00b7 '
                       'Exclusions \u00b7 Earthquake \u00b7 \u2026'))
        return b

    if t['NAMA_T'] in FIELD_PILIHAN:
        b.append(('x', '\u25b8    dipecah dari skalar akar TreatyIn'))
        for nm in FIELD_PILIHAN[t['NAMA_T']]:
            b.append(('f', '     ' + nm))
        return b

    if len(js) > 1:
        b.append(('x', '\u25b8    PERAN            %d daftar sebentuk digabung' % len(js)))
    for j in js[:1] if len(js) > 1 else js:
        f = skalar(j)
        if not f:
            continue
        for nm in f[:MAKS_FIELD]:
            b.append(('f', '     ' + nm))
        if len(f) > MAKS_FIELD:
            b.append(('n', '     \u2026 +%d field lagi' % (len(f) - MAKS_FIELD)))
    if len(js) > 1:
        b.append(('n', '     peran: ' + ', '.join(x.rsplit('.', 1)[-1] for x in js)[:96]))
    if not js:
        b.append(('n', '     tidak punya asal di pohon Pega'))

    for j in js:
        r = SIMPUL.get(j)
        if r and r['KELAS_PEGA'].startswith('(tidak'):
            b.append(('n', '     \u26a0 kelas TIDAK dideklarasikan: ' + j.rsplit('.', 1)[-1]))
    # Konvensi kedalaman mengikuti struktur-treatyin-lama.md: FIELD-nya yang disebut
    # tingkat 4, jadi WADAH terdalam ber-KEDALAMAN 3 di CSV. Jangan uji == 4 di sini:
    # tidak ada satu pun wadah ber-KEDALAMAN 4, dan ujinya tidak akan pernah kena.
    if js and js[0] in SIMPUL and int(SIMPUL[js[0]]['KEDALAMAN']) == 3:
        b.append(('n', '     \u25c4 fieldnya ada di KEDALAMAN 4 \u2014 terdalam di sistem lama'))
    return b


# Tulang punggung digambar di lajur AKAR apa pun KELOMPOK-nya di TSV.
# Sebabnya tata letak, bukan penggolongan: T_GENERAL_TREATY_IN adalah INDUK hampir
# semua kotak, dan menggambarnya di lajur `bersama` menempatkannya di BAWAH anak-anaknya.
# TREATY_IN memakai NAMA TABEL YANG SUDAH ADA di Oracle, jadi ia tanpa awalan T_.
TULANG_PUNGGUNG = ['TREATY_IN', 'T_TREATY_REVISION']


def kotak():
    out, sudah = [], set()
    punggung = [PETA[n] for n in TULANG_PUNGGUNG if n in PETA]
    sudah.update(TULANG_PUNGGUNG)
    out.append(('akar', punggung))
    for kel in URUT:
        if kel == 'akar':
            continue
        ts = [t for t in TABEL if t['KELOMPOK'] == kel and t['NAMA_T'] not in sudah]
        if ts:
            out.append((kel, ts))
    return out


# ------------------------------------------------------------------ relasi
def relasi():
    R = []
    for t in TABEL:
        if not t['INDUK']:
            continue
        ku = t['KUNCI'] if t['KUNCI'] and not t['KUNCI'].startswith('\u2014') \
            else pk_dari(t['INDUK'])
        kard = '1:1' if 'SHARED PK' in t['KUNCI'] else '1:N'
        hapus = 'ikut hapus' if t['KELOMPOK'] not in ('edm', 'migrasi') else 'tolak'
        R.append((t['KELOMPOK'], t['INDUK'], t['NAMA_T'], ku, kard, hapus,
                  t['ASAL_PEGA'][:150],
                  t['NAMA_T'] in LINTAS_SEKAT))
    return R


RELASI = relasi()


# ------------------------------------------------------------------ gambar SVG
def esc(s):
    return s.replace('&', '&amp;').replace('<', '&lt;').replace('>', '&gt;')


def svg():
    X0, W, BARIS, PITA = 300.0, 800.0, 19.0, 20.0
    y = 84.0
    pitas, boxes = [], []
    posisi = {}
    for kel, ts in kotak():
        warna, judul = WARNA[kel]
        pitas.append((X0, y, W, '#' + warna, judul))
        y += PITA + 18
        for t in ts:
            d = min(kedalaman(t), 4)
            x = X0 + d * 18
            w = W - d * 18
            br = baris_kotak(t)
            h = PITA + BARIS * len(br)
            kepala = t['NAMA_T']
            if t['PADANAN_DDL'] and not t['PADANAN_DDL'].startswith('\u2014'):
                kepala += '      \u2192  ' + t['PADANAN_DDL']
            n = ref(t)
            if n:
                kepala += '   n=%d' % n
            boxes.append((x, y, w, h, '#' + warna, kepala, br))
            posisi[t['NAMA_T']] = (x, y + h / 2.0)
            y += h + 16
        y += 14

    tinggi = y + 24
    o = []
    o.append('<svg width="1180" height="%d" viewBox="0 0 1180 %d" '
             'xmlns="http://www.w3.org/2000/svg" '
             'font-family="Arial, Helvetica, sans-serif">' % (int(tinggi), int(tinggi)))
    o.append('<style>.kp{font-size:11px;font-weight:700;fill:#fff;text-anchor:middle}'
             '.br{font-family:ui-monospace,Consolas,Menlo,monospace;font-size:10.2px}'
             '.pita{font-size:11.5px;font-weight:700;fill:#fff}</style>')

    for kel, induk, anak_, ku, kard, hapus, cat, lintas in RELASI:
        if induk not in posisi or anak_ not in posisi:
            continue
        xi, yi = posisi[induk]
        xa, ya = posisi[anak_]
        if xa <= xi:
            continue
        bx = xa - 9
        wr = '#C00000' if lintas else '#1F4E79'
        lb = '2.6' if lintas else '1.9'
        ds = ' stroke-dasharray="6,3"' if lintas else ''
        o.append('<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" '
                 'stroke-width="%s"%s/>' % (bx, yi, bx, ya, wr, lb, ds))
        o.append('<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" '
                 'stroke-width="%s"%s/>' % (bx, ya, xa, ya, wr, lb, ds))

    for x, yy, w, wr, tk in pitas:
        o.append('<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>'
                 % (x, yy, w, PITA, wr))
        o.append('<text class="pita" x="%.1f" y="%.1f">%s</text>'
                 % (x + 8, yy + 14.4, esc(tk)))

    for x, yy, w, h, wr, judul, br in boxes:
        o.append('<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#%s" '
                 'stroke="%s" stroke-width="1.8"/>' % (x, yy, w, h, PUCAT, wr))
        o.append('<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>'
                 % (x, yy, w, PITA, wr))
        o.append('<text class="kp" x="%.1f" y="%.1f">%s</text>'
                 % (x + w / 2.0, yy + 14.4, esc(judul)))
        yb = yy + PITA
        for jenis, teks in br:
            if jenis == 'pk':
                o.append('<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#%s"/>'
                         % (x + 1.4, yb, w - 2.8, BARIS, KREM))
            elif jenis == 'x':
                o.append('<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#%s"/>'
                         % (x + 1.4, yb, w - 2.8, BARIS, PITA_N))
            berat = ' font-weight="700"' if jenis in ('pk', 'fk', 'x') else ''
            wt = ' fill="%s"' % wr if jenis == 'x' else (
                ' fill="#5b6b80"' if jenis == 'n' else '')
            o.append('<text class="br" x="%.1f" y="%.1f"%s%s>%s</text>'
                     % (x + 10, yb + 13.6, berat, wt, esc(teks)))
            yb += BARIS
    o.append('</svg>')
    return '\n'.join(o), int(tinggi)


SVG, TINGGI = svg()

baris_rel = '\n'.join(
    '<tr%s><td>%d</td><td>%s</td><td><code>%s</code></td><td><code>%s</code></td>'
    '<td><code>%s</code></td><td>%s</td><td>%s</td><td>%s</td></tr>'
    % (' class="sekat"' if r[7] else '', i, r[0] + (' · LINTAS SEKAT' if r[7] else ''),
       r[1], r[2], r[3], r[4], r[5], esc(r[6]))
    for i, r in enumerate(RELASI, 1))

baris_tab = '\n'.join(
    '<tr><td><code>%s</code></td><td>%s</td><td>%s</td><td><code>%s</code></td>'
    '<td>%s</td></tr>'
    % (t['NAMA_T'], t['PADANAN_DDL'], t['KELOMPOK'],
       t.get('JALUR_PEGA', '') or '\u2014', esc(t['ASAL_PEGA'][:120]))
    for t in TABEL)

n_lintas = len([r for r in RELASI if r[7]])
DALAM_MAKS = max(kedalaman(t) for t in TABEL)
n_dalam4 = len([t for t in TABEL if kedalaman(t) == DALAM_MAKS])

HTML = u"""<!DOCTYPE html>
<html lang="id"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>ERD Struktur Treaty In</title>
<style>
:root{--tinta:#14243a;--redup:#5b6b80;--garis:#d6dfea;--latar:#f5f7fa;--kertas:#fff;
      --sekat:#C00000;--pohon:#1F4E79;}
@media (prefers-color-scheme: dark){:root:not([data-theme="light"]){
  --tinta:#e8edf4;--redup:#9fb0c4;--garis:#2b3644;--latar:#10151c;--kertas:#161c25;}}
:root[data-theme="dark"]{--tinta:#e8edf4;--redup:#9fb0c4;--garis:#2b3644;
  --latar:#10151c;--kertas:#161c25;}
*{box-sizing:border-box}
body{margin:0;background:var(--latar);color:var(--tinta);
  font:15px/1.6 -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif}
.bungkus{max-width:1240px;margin:0 auto;padding:28px 16px 64px}
h1{font-size:25px;margin:0 0 6px;letter-spacing:-.3px}
h2{font-size:18px;margin:34px 0 10px;padding-bottom:6px;border-bottom:2px solid var(--garis)}
.sub{color:var(--redup);font-size:13.5px;margin:0 0 18px}
.kartu{display:flex;flex-wrap:wrap;gap:10px;margin:0 0 20px}
.kartu div{background:var(--kertas);border:1px solid var(--garis);border-radius:9px;
  padding:9px 14px;min-width:132px}
.kartu b{display:block;font-size:21px;line-height:1.2}
.kartu span{color:var(--redup);font-size:11.5px;text-transform:uppercase;
  letter-spacing:.5px}
.papan{background:var(--kertas);border:1px solid var(--garis);border-radius:11px;
  padding:10px;overflow:auto}
svg{display:block}
.lgd{display:flex;flex-wrap:wrap;gap:8px 16px;margin:14px 0 0;font-size:12.5px;
  color:var(--redup)}
.lgd i{display:inline-block;width:12px;height:12px;border-radius:3px;
  margin-right:6px;vertical-align:-1px}
table{width:100%%;border-collapse:collapse;background:var(--kertas);
  border:1px solid var(--garis);border-radius:9px;overflow:hidden;font-size:13px}
th,td{padding:7px 9px;text-align:left;border-bottom:1px solid var(--garis);
  vertical-align:top}
th{background:var(--latar);font-size:11.5px;text-transform:uppercase;
  letter-spacing:.4px;color:var(--redup)}
tr:last-child td{border-bottom:0}
tr.sekat td{background:rgba(192,0,0,.07)}
code{font-family:ui-monospace,Consolas,Menlo,monospace;font-size:12px}
.nota{background:var(--kertas);border-left:4px solid var(--pohon);
  border-radius:0 9px 9px 0;padding:12px 16px;margin:16px 0;font-size:13.5px}
.nota.merah{border-left-color:var(--sekat)}
</style></head><body><div class="bungkus">

<h1>ERD Struktur &mdash; Treaty In</h1>
<p class="sub">Tabel datar <code>T_&hellip;</code> mengikuti <b>bentuk pohon
<code>TreatyIn</code> sistem lama</b>. Seluruh isi kotak dibaca dari sapuan
<b>329 berkas XML</b> &mdash; nama field dan cacah rujukan tidak diketik ulang.</p>

<div class="kartu">
  <div><b>%(n_tabel)d</b><span>tabel datar</span></div>
  <div><b>%(n_rel)d</b><span>relasi induk&ndash;anak</span></div>
  <div><b>%(n_lintas)d</b><span>melintasi sekat</span></div>
  <div><b>%(n_dalam4)d</b><span>di tingkat terdalam</span></div>
  <div><b>985</b><span>simpul pohon</span></div>
</div>

<div class="nota"><b>Indentasi kotak = kedalaman simpulnya di pohon Pega.</b>
Kotak yang menjorok satu tingkat adalah anak dari kotak di atasnya yang lebih kiri.
Baris krem adalah kunci utama, baris tebal kunci tamu, baris oranye penanda bahwa
tabel itu <b>menggabungkan beberapa daftar sebentuk</b> dan membedakannya lewat kolom
<code>PERAN</code>.</div>

<h2>Peta</h2>
<div class="papan">%(svg)s</div>
<div class="lgd">
  <span><i style="background:#1F4E79"></i>akar</span>
  <span><i style="background:#2E75B6"></i>layer</span>
  <span><i style="background:#31859C"></i>bagian</span>
  <span><i style="background:#C55A11"></i>fakultatif</span>
  <span><i style="background:#7030A0"></i>daftar tingkat kontrak</span>
  <span><i style="background:#548235"></i>dipakai bersama</span>
  <span><i style="background:#C00000"></i>addendum &amp; selisih &mdash; garis putus merah
    melintasi sekat</span>
  <span><i style="background:#808080"></i>jembatan migrasi</span>
  <span><i style="background:#BF8F00"></i>belum ada isinya</span>
</div>

<div class="nota"><b>Satu relasi lintas-sekat di sini, dua di <code>ERD.md</code> &mdash;
dan itu bukan pertentangan.</b> Peta ini hanya memuat tabel yang punya asal di pohon
<code>TreatyIn</code>. Relasi lintas-sekat kedua di <code>ERD.md</code> berpangkal pada
<code>BESARAN_DAPAT_DISESUAIKAN</code>, sebuah <b>tabel acuan rancangan</b> yang tidak
pernah ada di sistem lama, sehingga ia tidak punya kotak di peta ini. Untuk rancangan,
<code>ERD.md</code> yang mengikat.</div>

<div class="nota merah"><b>Dua halaman sengaja TIDAK menjadi tabel.</b>
<code>TreatyIn.OLDDATA</code> (142 simpul) karena aturan bisnisnya: data lama
<b>diselect, tidak disimpan ulang</b>. <code>TreatyIn.ActualValue</code> (152 simpul)
karena namanya tidak cocok dengan isinya &mdash; delapan isian menerima selisih.
Keduanya nyata dan hidup di sistem lama; yang tidak dibawa adalah bentuk
penyimpanannya, bukan datanya.</div>

<h2>Daftar relasi</h2>
<table><thead><tr><th>#</th><th>Kelompok</th><th>Induk</th><th>Anak</th>
<th>Kunci tamu</th><th>Kardinalitas</th><th>Perilaku hapus</th><th>Asal di pohon Pega</th>
</tr></thead><tbody>
%(baris_rel)s
</tbody></table>

<h2>Tambatan tiap tabel ke pohon</h2>
<table><thead><tr><th><code>NAMA_T</code></th><th>Padanan DDL</th><th>Kelompok</th>
<th>Jalur Pega</th><th>Keterangan</th></tr></thead><tbody>
%(baris_tab)s
</tbody></table>

<p class="sub" style="margin-top:26px">Dibangkitkan
<code>alat/buat-erd-struktur-treatyin.py</code> dari
<code>peta-nama-tabel-treatyin.tsv</code> dan <code>datar-treatyin-lama.csv</code>.
Bila gambar ini berbeda dari sumbernya, <b>alatnya yang salah</b>.</p>

</div></body></html>
""" % dict(n_tabel=len(TABEL), n_rel=len(RELASI), n_lintas=n_lintas,
           n_dalam4=n_dalam4, svg=SVG, baris_rel=baris_rel, baris_tab=baris_tab)

p = os.path.join(KELUARAN, 'ERD-STRUKTUR-TREATYIN.html')
with io.open(p, 'w', encoding='utf-8') as h:
    h.write(HTML)

print('tabel digambar      : %d' % len(TABEL))
print('kolom nyata TREATY_IN : %d  (dibaca dari argumen PEGA_TREATY_IN)' % len(KOLOM_NYATA))
print('skalar akar dipecah   : %d  (seluruhnya ADA di pohon)' % N_PILIHAN)
print('relasi induk-anak   : %d' % len(RELASI))
print('  melintasi sekat   : %d' % n_lintas)
print('tingkat terdalam    : %d  (berisi %d tabel)' % (DALAM_MAKS, n_dalam4))
tak_terduga, disengaja = beda_kedalaman()
print('kedalaman rantai = pohon + %d' % SELISIH_WAJAR)
print('  menyimpang, DISENGAJA : %d' % len(disengaja))
for nm, dr, dp, sl in disengaja:
    print('     %-30s rantai=%d pohon=%d  — %s'
          % (nm, dr, dp, SELISIH_DISENGAJA[nm]))
print('  menyimpang, TAK TERDUGA : %d' % len(tak_terduga))
for nm, dr, dp, sl in tak_terduga:
    print('     ⚠ %-28s rantai=%d pohon=%d  selisih=%+d' % (nm, dr, dp, sl))
print('tinggi SVG          : %d px' % TINGGI)
print('---')
print(p)
