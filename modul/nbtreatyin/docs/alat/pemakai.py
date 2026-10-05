"""Siapa memakai sebuah nama (properti, parameter, kolom) di korpus NB Treaty In - HANYA MEMBACA.

Bukti "nol pemakai" / "pemakai tunggal" untuk butir bab 4 prompt putaran 2 (AC 26 `RNM_SHARE`,
`BROKERAGE`; AC 66 `isFOR`). Memakai ulang `korpus.muat()` dan `graf.Graf(...).terjangkau()`:
setiap rule yang teksnya (salinan section tertanam dibuang, sama dengan graf.py) cocok dengan
pola dicetak beserta tanda TERJANGKAU / tak-terjangkau dari titik masuk nyata.

    python pemakai.py "\\bisFOR\\b" "RNM_SHARE" "\\bBROKERAGE\\b"

Pola adalah regex Python, peka huruf besar-kecil (nama properti Pega peka huruf).
"""
from __future__ import annotations

import html
import os
import re
import sys
import xml.etree.ElementTree as ET

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

import korpus  # noqa: E402
from graf import Graf  # noqa: E402


def pemakai(data, jangkau, pola):
    rx = re.compile(pola)
    for d in data:
        akar = ET.parse(os.path.join(korpus.KORPUS, d["berkas"])).getroot()
        teks = html.unescape(ET.tostring(korpus._tanpa_tanaman(akar), encoding="unicode"))
        cocok = sorted({m.group(0) for m in rx.finditer(teks)})
        if cocok:
            yield (d["jenis"], d["nama"]) in jangkau, d["jenis"], d["nama"], cocok


def main(pola_pola):
    sys.stdout.reconfigure(encoding="utf-8")
    data = korpus.muat()
    jangkau = Graf(data).terjangkau()
    print(f"{len(data)} rule, {len(jangkau)} terjangkau")
    for pola in pola_pola:
        hasil = list(pemakai(data, jangkau, pola))
        print(f"== {pola}: {len(hasil)} rule ({sum(1 for h in hasil if h[0])} terjangkau)")
        for terjangkau, jenis, nama, cocok in hasil:
            print(f"  {'TERJANGKAU    ' if terjangkau else 'tak-terjangkau'}  {jenis}/{nama}: {', '.join(cocok[:12])}")


if __name__ == "__main__":
    main(sys.argv[1:])
