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
  error,
  sembunyikanLabel = false,
  jedaMs = JEDA_KETIK_MS,
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
  /** Pesan galat di bawah medan — bentuk yang sama dengan `Pilih`/`Field`. */
  error?: string;
  /**
   * Sembunyikan labelnya; namanya pindah ke `aria-label`.
   *
   * ⛔ Untuk SEL TABEL: `<th>` kolomnya sudah menamai medan itu, jadi label
   * kedua di dalam sel menggandakan judul kolom dan menaikkan tinggi baris.
   * Mengosongkan `label` saja akan membuat medannya anonim bagi pembaca
   * layar — karena itu disembunyikan, bukan dihapus.
   */
  sembunyikanLabel?: boolean;
  /**
   * Jeda ketik sebelum `onCari` dipanggil, milidetik.
   *
   * ⭐ Bawaannya 250ms karena pemakai pertama menyaring DI SERVER. Daftar
   * yang sudah ada di klien tidak perlu menunggu siapa pun: `PilihCari`
   * memberi `0`, dan saringannya terasa seketika saat diketik.
   */
  jedaMs?: number;
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
    jeda.current = window.setTimeout(() => onCari(t.trim()), jedaMs);
  }

  return (
    <div
      className="field pilih-saring"
      onBlur={(e) => {
        if (!e.currentTarget.contains(e.relatedTarget as Node | null)) tutup();
      }}
    >
      {!sembunyikanLabel && (
        <label className="field__label" htmlFor={idInput}>
          {label}
          {required && <span className="field__req">*</span>}
        </label>
      )}
      <div className="pilih-saring__kotak">
        <input
          id={idInput}
          className={"field__input" + (error ? " field__input--error" : "")}
          role="combobox"
          aria-expanded={terbuka}
          aria-controls={idDaftar}
          aria-label={sembunyikanLabel ? label : undefined}
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
      {error && <div className="field__error">{error}</div>}
    </div>
  );
}

/**
 * Cocokkah satu butir dengan kata yang diketik?
 *
 * ⛔ Pencocokan BAGIAN TEKS, bukan awalan — pemakai mengetik potongan yang
 * diingatnya (`"QS"`, `"181"`), bukan selalu huruf pertama. Tanpa huruf
 * besar-kecil, dan `keterangan` (biasanya kode/ID) ikut dicari sebab di
 * sanalah nomor kontrak berada.
 */
export function cocokSaring(o: OpsiSaring, kata: string): boolean {
  const k = kata.trim().toLowerCase();
  if (k === "") return true;
  return o.label.toLowerCase().includes(k) || (o.keterangan ?? "").toLowerCase().includes(k);
}

/**
 * PilihCari — `Pilih` yang DAPAT DIKETIK untuk mencari, saringan di KLIEN.
 *
 * ⛔ Permintaan pemilik proses 8 Oktober 2026 untuk seluruh dropdown Treaty
 * In dan Adjustment: *"kondisi saat melakukan pengetikannya seharusnya
 * terlihat layaknya melakukan mencari, kemudian data yang keluar adalah yang
 * 100% mirip dengan yang diketik"*.
 *
 * ⭐ API-nya SENGAJA sama persis dengan `Pilih` (`label` · `value` ·
 * `onChange` · `opsi` · `error` · `required` · `kosong`), supaya penggantian
 * di layar cukup satu baris impor dan nol perubahan di badan komponennya.
 * Dua puluhan pemakaian ditukar tanpa menyentuh logika satu pun di antaranya.
 *
 * ⚠️ BUKAN isian bebas: teks yang diketik tanpa memilih dikembalikan ke
 * pilihan terakhir (`PilihSaring` yang menjaganya). Itu beda pokoknya dari
 * `<datalist>` yang dibuang 8 Oktober 2026 — di sana ketikan apa pun lolos
 * menjadi nilai.
 */
export function PilihCari({
  label,
  value,
  onChange,
  opsi,
  error,
  required,
  kosong,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  opsi: readonly { value: string; label: string }[];
  error?: string;
  required?: boolean;
  kosong?: string;
}) {
  const teksUI = useTeksUI();
  const [kata, setKata] = useState("");
  const teksKosong = kosong ?? teksUI.pilihKosong;

  // ⛔ Nilai TERSIMPAN yang tidak ada di daftar tetap ditawarkan, bukan
  // dijatuhkan — dropdown yang diam-diam mengosongkan nilai lama terbaca
  // sebagai data yang hilang.
  const asing = value !== "" && !opsi.some((o) => o.value === value);
  const semua: OpsiSaring[] = [
    { value: "", label: teksKosong },
    ...(asing ? [{ value, label: value }] : []),
    ...opsi.map((o) => ({ value: o.value, label: o.label })),
  ];
  const terpilih = semua.find((o) => o.value === value);

  return (
    <PilihSaring
      label={label}
      value={value}
      teksTerpilih={terpilih?.label ?? teksKosong}
      opsi={semua.filter((o) => cocokSaring(o, kata))}
      onCari={setKata}
      onPilih={(o) => {
        onChange(o.value);
      }}
      required={required}
      error={error}
      jedaMs={0}
    />
  );
}
