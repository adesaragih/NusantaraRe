/**
 * KelompokMenu — satu kelompok sidebar yang dapat dilipat (ronde 103, A1–A5).
 *
 * # Kenapa <button>, bukan <div onClick>
 *
 * Aksesibilitas di sini bukan hiasan (A4): tombol asli mendapat fokus
 * keyboard, Enter/Space, dan peran yang benar TANPA satu baris kode —
 * sedangkan div yang bisa diklik menutup jalan bagi pemakai keyboard, dan
 * setiap perbaikannya (tabIndex, onKeyDown, role) adalah menulis ulang apa
 * yang <button> sudah punya. `aria-expanded` memberi tahu pembaca layar
 * keadaan buka/tutupnya.
 *
 * # Keadaan
 *
 * Logika buka/tutup hidup di src/lib/lipatMenu.ts (teruji murni):
 * bawaan TERBUKA, kelompok ber-halaman-aktif SELALU terbuka, dan lipatan
 * tersimpan di localStorage dengan semua jalan gagal jatuh ke bawaan.
 *
 * # Bertingkat dua (A5)
 *
 * children adalah ReactNode, jadi KelompokMenu boleh memuat KelompokMenu
 * lain — struktur yang CITRIX akan pakai. Menu Treaty yang ada TIDAK
 * ditambah tingkatnya: perubahan ini tampilan, bukan susunan.
 */
import { useState, type ReactNode } from "react";

import {
  bacaLipatan,
  simpanLipatan,
  terbukaKah,
  type GudangMini,
} from "../lib/lipatMenu";

// localStorage disentuh lewat fungsi supaya kegagalannya (mode privat,
// kebijakan peramban) terkurung — bacaLipatan/simpanLipatan sudah menelan
// galatnya dan jatuh ke bawaan.
function gudang(): GudangMini | null {
  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

export function KelompokMenu({
  nama,
  memuatAktif,
  anak = false,
  children,
}: {
  /** Judul kelompok — juga kunci penyimpanan lipatannya. */
  nama: string;
  /** true bila halaman AKTIF ada di dalam kelompok ini (A2). */
  memuatAktif: boolean;
  /** true untuk kelompok di dalam kelompok (tingkat dua — A5). */
  anak?: boolean;
  children: ReactNode;
}) {
  const [terbuka, setTerbuka] = useState(() =>
    terbukaKah(nama, bacaLipatan(gudang()), memuatAktif),
  );

  const balik = () => {
    setTerbuka((v) => {
      const baru = !v;
      const l = bacaLipatan(gudang());
      l[nama] = baru;
      simpanLipatan(gudang(), l);
      return baru;
    });
  };

  return (
    <div className={"kelompok" + (anak ? " kelompok--anak" : "")}>
      <button
        type="button"
        className="kelompok__judul"
        aria-expanded={terbuka}
        onClick={balik}
      >
        <span
          className={
            "kelompok__panah" + (terbuka ? " kelompok__panah--buka" : "")
          }
          aria-hidden="true"
        >
          ▸
        </span>
        <span className="kelompok__teks">{nama}</span>
      </button>
      {terbuka && <div className="kelompok__isi">{children}</div>}
    </div>
  );
}
