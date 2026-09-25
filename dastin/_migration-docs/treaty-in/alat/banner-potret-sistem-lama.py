# -*- coding: utf-8 -*-
r"""
banner-potret-sistem-lama.py -- menyisipkan satu baris penanda di kepala TIAP LEMBAR
pada berkas yang merupakan POTRET SISTEM LAMA, bukan rancangan.

KENAPA ADA
    Kedua berkas berjudul di dalamnya "rancangan tabel" dan berkolom "padanan di
    model baru". Keduanya membuatnya MUDAH DIKIRA RANCANGAN -- dan 64 tiket hendak
    dikerjakan, sebagiannya oleh orang yang belum pernah membuka berkas lain.

APA YANG DAPAT DIBUKTIKAN PERKAKAS INI
    Bahwa tiap lembar yang berhasil disunting MEMUAT banner, dan bahwa banner itu
    tidak disisipkan dua kali.

APA YANG TIDAK
    Ia TIDAK dapat menyunting berkas yang sedang TERBUKA di Excel. Bila gagal, ia
    MENYEBUT NAMA BERKASNYA dan BERHENTI -- ia TIDAK menyimpan dengan nama lain,
    sebab salinan bernama lain melahirkan berkas kelima yang usianya berbeda lagi.

    Ia juga TIDAK memeriksa ISI lembarnya. Sebuah lembar yang sebenarnya sudah
    rancangan baru tetap akan diberi banner potret bila berkasnya ada di daftar.
    Daftar berkasnya DITULIS TANGAN di bawah, dan itu disengaja.

PENOLAKAN dilaporkan satu baris per sebab.

Pakai:  PYTHONIOENCODING=utf-8 python alat/banner-potret-sistem-lama.py
"""
import io, os, sys
from collections import Counter

from openpyxl import load_workbook
from openpyxl.styles import Font, PatternFill, Alignment

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
FOLDER = os.path.join(AKAR, '4-erd-dan-tabel-datar')

BARIS_1 = 'POTRET SISTEM LAMA, 24 September 2026. BUKAN RANCANGAN.'
BARIS_2 = ('Skema yang dibangun ada di 2-to-spec/KAMUS-KOLOM.md dan 2-to-spec/ddl-usulan/. '
           'Gambar skema barunya: 4-erd-dan-tabel-datar/ERD-SKEMA-BARU.xlsx')

# Daftar DITULIS TANGAN. Perkakas ini tidak menebak berkas mana yang potret.
POTRET = [
    ('Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx',
     '44 tabel sistem lama, dari pohon clipboard Pega'),
    ('_arsip/ERD-TREATY-IN-DAN-EDM.xlsx',
     '44 tabel sistem lama, bentuk ERD -- DIARSIPKAN 25 Sep 2026'),
    ('_arsip/Diagram-Skema-Tabel-TreatyIn-dan-EDM.xlsx',
     'pendahulu v2 -- potret sistem lama -- DIARSIPKAN 25 Sep 2026'),
]
# Berkas BASI: skema baru, tetapi versi lama. Banner berbeda bunyinya.
BASI = [
    ('Diagram-Skema-Tabel-TreatyMasuk.xlsx',
     'ERD skema baru versi 24 Sep 17:50 -- NOL DOKUMEN_ADDENDUM, NOL tabel NILAI_*'),
]
BASI_1 = 'BASI — versi 24 September 2026 17:50. JANGAN DIPAKAI.'
BASI_2 = ('Yang mutakhir: 4-erd-dan-tabel-datar/ERD-SKEMA-BARU.xlsx, dibangkitkan dari '
          '2-to-spec/ddl-usulan/. Berkas ini NOL DOKUMEN_ADDENDUM dan NOL tabel NILAI_*.')

MERAH = PatternFill('solid', fgColor='C00000')
KUNING = PatternFill('solid', fgColor='FFC000')
PUTIH13 = Font(name='Arial', size=13, bold=True, color='FFFFFF')
HITAM9 = Font(name='Arial', size=9, bold=True, color='000000')
KIRI = Alignment(horizontal='left', vertical='center')

tolak = Counter()
hasil = []


def pasang(nama, ket, b1, b2, warna1, warna2):
    p = os.path.join(FOLDER, nama)
    if not os.path.exists(p):
        tolak['berkas tidak ada'] += 1
        hasil.append((nama, 'TIDAK ADA', 0, ket))
        return
    kunci = os.path.join(FOLDER, '~$' + nama)
    terbuka = os.path.exists(kunci)
    try:
        wb = load_workbook(p)
    except Exception as e:
        tolak['berkas tidak dapat dibaca'] += 1
        hasil.append((nama, 'TIDAK TERBACA: %s' % e, 0, ket))
        return
    n = 0
    lewat = 0
    for ws in wb.worksheets:
        if str(ws.cell(1, 1).value or '').startswith(b1[:20]):
            lewat += 1
            tolak['lembar yang SUDAH berbanner (tidak disisipkan dua kali)'] += 1
            continue
        ws.insert_rows(1, 2)
        for baris, teks, fill, font in ((1, b1, warna1, PUTIH13), (2, b2, warna2, HITAM9)):
            s = ws.cell(baris, 1, teks)
            s.fill = fill
            s.font = font
            s.alignment = KIRI
            for c in range(2, 40):
                ws.cell(baris, c).fill = fill
        ws.row_dimensions[1].height = 22
        ws.row_dimensions[2].height = 16
        n += 1
    try:
        wb.save(p)
        hasil.append((nama, 'BERBANNER di %d lembar%s' % (n, ' (%d sudah ada)' % lewat if lewat else ''), n, ket))
    except Exception as e:
        tolak['berkas TIDAK DAPAT DITIMPA (kemungkinan terbuka di Excel)'] += 1
        hasil.append((nama, 'GAGAL DISIMPAN%s — %s' % (
            ' [berkas kunci ~$ ADA]' if terbuka else '', e), 0, ket))


print('=== banner-potret-sistem-lama.py ===')
print()
for nama, ket in POTRET:
    pasang(nama, ket, BARIS_1, BARIS_2, MERAH, KUNING)
for nama, ket in BASI:
    pasang(nama, ket, BASI_1, BASI_2, MERAH, KUNING)

print('HASIL per berkas:')
for nama, st, n, ket in hasil:
    print('   %-46s %s' % (nama, st))
    print('   %-46s   isinya: %s' % ('', ket))
print()
print('DITOLAK / TIDAK DIKERJAKAN -- satu baris per sebab:')
for k, v in tolak.most_common():
    print('   %-60s %d' % (k, v))
if not tolak:
    print('   (tidak ada)')
print()
gagal = [h for h in hasil if h[1].startswith('GAGAL')]
if gagal:
    print('BERKAS YANG TIDAK DAPAT DITIMPA -- disebut namanya, TIDAK disimpan dengan nama lain:')
    for nama, st, n, ket in gagal:
        print('   %s' % nama)
    print()
    print('   Tutup berkasnya di Excel, lalu jalankan ulang perkakas ini.')
    print('   Menyimpannya dengan nama lain akan melahirkan berkas kelima yang usianya')
    print('   berbeda lagi -- dan lima berkas ERD berusia berbeda adalah bahaya tersendiri.')
print()
print('-- BATAS PERKAKAS INI --')
print('   Ia tidak memeriksa ISI lembar. Daftar berkas potret DITULIS TANGAN di dalamnya,')
print('   dan itu disengaja: menebaknya dari nama berkas akan salah pada berkas berikutnya.')
sys.exit(1 if gagal else 0)
