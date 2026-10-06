"""Ringkas Struktur_MenuNBTreatyIn.xlsx (pohon pemanggilan dari menu portal) -> menu.json.

HANYA MEMBACA. Kolom sheet `Struktur`: 14 kolom nomor bertingkat, lalu Jenis Rule, Nama Rule,
Shape ID / Dipanggil dari, XML, Class, Status XML.
"""
import json
import os
import sys

import openpyxl

KORPUS = os.environ.get("NBTREATYIN_KORPUS", r"D:/NUSARE DEV/NusantaraRe/NB Treaty In (Done)")
DIR = os.path.dirname(os.path.abspath(__file__))


def main():
    p = os.path.join(KORPUS, "Struktur_MenuNBTreatyIn.xlsx")
    ws = openpyxl.load_workbook(p, read_only=True, data_only=True)["Struktur"]
    rows = list(ws.iter_rows(min_row=2, values_only=True))
    atas, dalam = [], 0
    for r in rows:
        lv = next((i for i in range(14) if r[i] is not None), None)
        if lv is None:
            continue
        dalam = max(dalam, lv + 1)
        if lv <= 2:
            atas.append({"no": str(r[lv]), "jenis": r[14] or "", "nama": r[15] or "", "dari": r[16] or ""})
    out = {
        "baris": len(rows), "kedalaman": dalam, "atas": atas,
        "catatan": ("Akar tunggal: Harness `SFAPortalOpportunities` (menu portal). Satu-satunya jalan ke alur "
                    "realisasi adalah tombol *Create* portal (`createWork` → `Flow/InputRealizationTreatyIn`). "
                    "Tidak ada butir menu lain di pohon ini — sistem baru tidak boleh menambah menu di luar itu."),
    }
    tujuan = sys.argv[1] if len(sys.argv) > 1 else os.path.join(DIR, "menu.json")
    with open(tujuan, "w", encoding="utf-8", newline="\n") as fh:
        json.dump(out, fh, ensure_ascii=False, indent=1)
    print(len(atas), "baris tingkat 1-3 ->", tujuan)


if __name__ == "__main__":
    main()
