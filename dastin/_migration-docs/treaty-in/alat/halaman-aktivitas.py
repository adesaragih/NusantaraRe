# -*- coding: utf-8 -*-
r"""
halaman-aktivitas.py -- mencetak daftar HALAMAN sebuah aturan Pega beserta KELASNYA.

KENAPA ADA
    Penulis properti yang sasarannya memakai nama halaman relatif -- `Primary.X`,
    `Local.X`, `SaveData.X` -- tidak dapat dibaca oleh sapuan penulis maupun oleh
    pembangun pohon, sebab keduanya bekerja atas JALUR PENUH. Tanpa kelas
    halamannya, `Primary.PremiumEarnedList` dapat berarti dua entitas yang berbeda.

NAMA TAG YANG DIPAKAI, dicocokkan lebih dulu ke `datar-nama-elemen.csv`:
    pxPageName   1.162 kemunculan / 431 berkas   (Harness, Section, Activity, Model,
                                                  Report-Definition, When)
    pyClassName  1.141 kemunculan / 708 berkas
    pyPageName       4 kemunculan /   3 berkas   -- bentuk kedua, ikut dibaca

BATAS
    Ia membaca DEKLARASI, bukan isi saat jalan. Halaman yang kelasnya ditetapkan
    saat berjalan (misal lewat Page-New berkelas parameter) TIDAK terbaca di sini,
    dan itu dilaporkan sebagai penolakan, bukan didiamkan.

Pakai:  python alat/halaman-aktivitas.py "Treaty In/Activity/DetailCalculation.xml"
"""
import sys, os, xml.etree.ElementTree as ET

AKAR = r'D:\XML_NURE'

def main(rel):
    p = os.path.join(AKAR, rel.replace('/', os.sep))
    if not os.path.exists(p):
        print('berkas tidak ada: %s' % p); return
    root = ET.parse(p).getroot()
    pasang, tolak_tanpa_kelas, tolak_tanpa_nama = [], 0, 0
    for el in root.iter():
        nm = el.find('pxPageName')
        if nm is None:
            nm = el.find('pyPageName')
        kl = el.find('pyClassName')
        if nm is None and kl is None:
            continue
        n = (nm.text or '').strip() if nm is not None else ''
        k = (kl.text or '').strip() if kl is not None else ''
        if not n:
            tolak_tanpa_nama += 1; continue
        if not k:
            tolak_tanpa_kelas += 1
        pasang.append((n, k or '(kelas kosong)'))
    print('=== halaman-aktivitas.py : %s ===' % os.path.basename(p))
    print('HALAMAN TERDEKLARASI: %d' % len(pasang))
    for n, k in pasang:
        print('   %-20s %s' % (n, k))
    print()
    print('DITOLAK -- satu baris per sebab:')
    print('   %-44s %d' % ('elemen berkelas tetapi tanpa nama halaman', tolak_tanpa_nama))
    print('   %-44s %d' % ('halaman bernama tetapi tanpa kelas', tolak_tanpa_kelas))

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print(__doc__)
    else:
        main(sys.argv[1])
