# -*- coding: utf-8 -*-
r"""
banner-html-basi.py -- menyisipkan banner ke berkas ERD berbentuk HTML.

KENAPA BANNER, BUKAN DIBANGKITKAN ULANG -- dan ini keputusan, bukan kemalasan
    Membangkitkan ulang bentuk HTML berarti memelihara PENAMPIL KEDUA atas sumber
    yang sama. Dua penampil dapat menyimpang, dan yang menyimpang diam-diam adalah
    yang jarang dibuka. Satu penampil yang dibangkitkan -- ERD-SKEMA-BARU.xlsx --
    ditambah berkas lama yang BERTANDA BASI lebih aman daripada dua penampil yang
    sama-sama mengaku mutakhir.

    Bila kelak HTML dituntut, ia dibangkitkan dari erd-skema-baru.json yang SAMA,
    bukan ditulis tangan.

APA YANG DAPAT DIBUKTIKAN PERKAKAS INI
    Bahwa tiap berkas dalam daftarnya memuat banner, dan tidak dua kali.
APA YANG TIDAK
    Ia tidak membaca isi HTML-nya. Daftarnya ditulis tangan.

PENOLAKAN dilaporkan satu baris per sebab.
"""
import io, os
from collections import Counter

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
FOLDER = os.path.join(AKAR, '4-erd-dan-tabel-datar')
PENANDA = 'data-banner="nure-basi"'

BERKAS = [
    ('ERD-TREATY-MASUK.html', 'BASI',
     'ERD skema baru versi 24 September 2026 17:50',
     'NOL <code>DOKUMEN_ADDENDUM</code> · NOL tabel <code>NILAI_*</code> · '
     'kini skema baru punya <b>35 tabel</b> dan <b>38 kunci asing</b>.'),
    ('ERD-STRUKTUR-TREATYIN.html', 'POTRET',
     'POTRET SISTEM LAMA, 24 September 2026',
     'Ia menggambarkan <b>pohon clipboard Pega</b>, bukan skema yang dibangun.'),
]

BANNER = '''<div %s style="position:sticky;top:0;z-index:9999;font-family:Arial,sans-serif">
  <div style="background:#C00000;color:#fff;font-size:17px;font-weight:bold;padding:10px 14px">
    %s &mdash; %s. JANGAN DIPAKAI SEBAGAI RANCANGAN.
  </div>
  <div style="background:#FFC000;color:#000;font-size:12px;padding:8px 14px">
    %s<br>
    Yang <b>mutakhir dan mengikat</b>:
    <code>2-to-spec/KAMUS-KOLOM.md</code> (nama dan tipe kolom) &middot;
    <code>2-to-spec/ddl-usulan/</code> (tabel dan kunci asing) &middot;
    <code>4-erd-dan-tabel-datar/ERD-SKEMA-BARU.xlsx</code> (gambar, <b>dibangkitkan</b> dari DDL).
  </div>
</div>
'''

tolak = Counter()
hasil = []
for nama, jenis, judul, ket in BERKAS:
    p = os.path.join(FOLDER, nama)
    if not os.path.exists(p):
        tolak['berkas tidak ada'] += 1
        hasil.append((nama, 'TIDAK ADA'))
        continue
    s = io.open(p, encoding='utf-8', errors='replace').read()
    if PENANDA in s:
        tolak['berkas yang SUDAH berbanner (tidak disisipkan dua kali)'] += 1
        hasil.append((nama, 'sudah berbanner — dilewati'))
        continue
    b = BANNER % (PENANDA, jenis, judul, ket)
    low = s.lower()
    i = low.find('<body')
    if i >= 0:
        j = s.find('>', i)
        baru = s[:j + 1] + '\n' + b + s[j + 1:]
        cara = 'disisipkan sesudah <body>'
    else:
        baru = b + s
        cara = 'DILEKATKAN DI KEPALA — tidak ada tag <body>'
        tolak['HTML tanpa tag <body>'] += 1
    try:
        io.open(p, 'w', encoding='utf-8').write(baru)
        hasil.append((nama, 'BERBANNER — %s' % cara))
    except Exception as e:
        tolak['berkas TIDAK DAPAT DITIMPA'] += 1
        hasil.append((nama, 'GAGAL DISIMPAN — %s' % e))

print('=== banner-html-basi.py ===')
for nama, st in hasil:
    print('   %-34s %s' % (nama, st))
print()
print('DITOLAK -- satu baris per sebab:')
for k, v in tolak.most_common():
    print('   %-58s %d' % (k, v))
if not tolak:
    print('   (tidak ada)')
print()
print('-- BATAS PERKAKAS INI --')
print('   Ia tidak membaca isi HTML-nya; daftarnya ditulis tangan. Dan banner')
print('   TIDAK membuat berkasnya benar -- ia hanya membuat kebasiannya terlihat.')
