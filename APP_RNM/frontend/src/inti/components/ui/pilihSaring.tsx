/**
 * PilihSaring — dropdown yang DAPAT DIFILTER dengan mengetik (combobox).
 *
 * Satu kotak: panah (atau klik / ArrowDown) membuka SELURUH daftar, mengetik
 * menyaringnya, klik / Enter memilih, Escape / fokus pergi menutup. Pilihan
 * HANYA berubah saat satu butir dipilih — teks yang diketik tanpa memilih
 * dikembalikan ke pilihan terakhir, sehingga tidak ada keadaan "setengah
 * terpilih" yang sampai ke Save.
 *
 * Penyaringan dikerjakan PEMANGGIL lewat `onCari` (mis. ke server, karena
 * master bisa lebih besar daripada satu jawaban), dengan jeda ketik 250ms;
 * `onCari('')` = daftar lengkap / awal. Pemakai pertama: Treaty Contract Out
 * (keputusan work owner 30-09-2026: "drop down tapi tetep bisa di filter").
 *
 * ⛔ Bukan `<datalist>`: Firefox tidak memberinya panah, dan di form Edit
 * datalist hanya menampilkan butir yang cocok dengan teks terpilih.
 */
import { useEffect, useId, useRef, useState } from "react";

import { useTeksUI } from "./bahasaUI";
import { IkonChevron } from "./dasar";

export interface OpsiSaring {
  value: string;
  label: string;
  /** Teks kecil di sisi kanan butir (mis. ID) — membedakan nama kembar. */
  keterangan?: string;
}

/** Butir sorotan berikutnya untuk ArrowUp/ArrowDown, dijepit di dalam daftar. */
export function sorotBerikut(sorot: number, jumlah: number, arah: 1 | -1): number {
  if (jumlah <= 0) return 0;
  return Math.min(jumlah - 1, Math.max(0, sorot + arah));
}

/** Kata cari saat daftar DIBUKA: teks pilihan terakhir berarti "tampilkan semua". */
export function kataSaatBuka(ketik: string, teksTerpilih: string): string {
  return ketik === teksTerpilih ? "" : ketik;
}

const JEDA_KETIK_MS = 250;

export function PilihSaring({
  label,
  value,
  teksTerpilih,
  opsi,
  onCari,
  onPilih,
  required,
  memuat = false,
}: {
  label: string;
  /** Nilai terpilih (mis. ID); kosong = belum ada. */
  value: string;
  /** Teks yang tampil untuk nilai terpilih. */
  teksTerpilih: string;
  opsi: readonly OpsiSaring[];
  /** Muat ulang `opsi` untuk kata ini; `""` = daftar awal. */
  onCari: (kata: string) => void;
  onPilih: (o: OpsiSaring) => void;
  required?: boolean;
  /** true selama `onCari` belum menjawab. */
  memuat?: boolean;
}) {
  const teks = useTeksUI();
  const [ketik, setKetik] = useState(teksTerpilih);
  const [terbuka, setTerbuka] = useState(false);
  const [sorot, setSorot] = useState(0);
  const jeda = useRef<number | undefined>(undefined);
  const idInput = useId();
  const idDaftar = useId();

  // Pilihan berubah dari luar (dipilih, atau form lain dibuka): teks ikut.
  useEffect(() => {
    setKetik(teksTerpilih);
  }, [teksTerpilih, value]);
  useEffect(() => () => window.clearTimeout(jeda.current), []);

  function buka(): void {
    setTerbuka(true);
    setSorot(0);
    window.clearTimeout(jeda.current);
    onCari(kataSaatBuka(ketik, teksTerpilih));
  }
  function tutup(): void {
    window.clearTimeout(jeda.current);
    setTerbuka(false);
    setKetik(teksTerpilih);
  }
  function pilih(o: OpsiSaring): void {
    onPilih(o);
    setKetik(o.label);
    setTerbuka(false);
  }
  function ubahKetik(t: string): void {
    setKetik(t);
    setTerbuka(true);
    setSorot(0);
    window.clearTimeout(jeda.current);
    jeda.current = window.setTimeout(() => onCari(t.trim()), JEDA_KETIK_MS);
  }

  return (
    <div
      className="field pilih-saring"
      onBlur={(e) => {
        if (!e.currentTarget.contains(e.relatedTarget as Node | null)) tutup();
      }}
    >
      <label className="field__label" htmlFor={idInput}>
        {label}
        {required && <span className="field__req">*</span>}
      </label>
      <div className="pilih-saring__kotak">
        <input
          id={idInput}
          className="field__input"
          role="combobox"
          aria-expanded={terbuka}
          aria-controls={idDaftar}
          aria-autocomplete="list"
          aria-required={required}
          autoComplete="off"
          placeholder={teks.ketikUntukMenyaring}
          value={ketik}
          onClick={() => {
            if (!terbuka) buka();
          }}
          onChange={(e) => {
            ubahKetik(e.target.value);
          }}
          onKeyDown={(e) => {
            if (e.key === "ArrowDown") {
              e.preventDefault();
              if (terbuka) setSorot((s) => sorotBerikut(s, opsi.length, 1));
              else buka();
            } else if (e.key === "ArrowUp") {
              e.preventDefault();
              setSorot((s) => sorotBerikut(s, opsi.length, -1));
            } else if (e.key === "Enter" && terbuka) {
              e.preventDefault();
              const o = opsi[sorot];
              if (o !== undefined) pilih(o);
            } else if (e.key === "Escape" && terbuka) {
              // Menutup DAFTAR saja - bukan popup tempatnya berada.
              e.preventDefault();
              e.stopPropagation();
              tutup();
            }
          }}
        />
        <button
          type="button"
          className="pilih-saring__panah"
          aria-label={teks.bukaDaftar}
          tabIndex={-1}
          onMouseDown={(e) => {
            e.preventDefault();
          }}
          onClick={() => {
            if (terbuka) tutup();
            else buka();
          }}
        >
          <IkonChevron />
        </button>
      </div>
      {terbuka && (
        <ul id={idDaftar} role="listbox" className="pilih-saring__daftar">
          {memuat ? (
            <li className="pilih-saring__kosong">{teks.memuat}</li>
          ) : opsi.length === 0 ? (
            <li className="pilih-saring__kosong">{teks.tidakCocok}</li>
          ) : (
            opsi.map((o, i) => (
              <li
                key={o.value}
                role="option"
                aria-selected={o.value === value}
                className={"pilih-saring__opsi" + (i === sorot ? " pilih-saring__opsi--sorot" : "")}
                onMouseDown={(e) => {
                  // Pilih SEBELUM blur menutup daftar.
                  e.preventDefault();
                  pilih(o);
                }}
                onMouseEnter={() => {
                  setSorot(i);
                }}
              >
                <span>{o.label}</span>
                {o.keterangan !== undefined && <span className="pilih-saring__ket">{o.keterangan}</span>}
              </li>
            ))
          )}
        </ul>
      )}
    </div>
  );
}
