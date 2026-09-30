// Kontrol dasar yang dipakai lebih dari satu layar.
//
// `Field` sebelumnya tinggal di src/App.tsx. Ia dipindah ke sini APA ADANYA
// ketika layar penawaran mulai memerlukannya: menyalinnya akan membuat dua
// definisi yang harus diubah berpasangan, dan pasangan seperti itu selalu
// terlepas. Halaman login tidak berubah — JSX-nya tetap memanggil `Field`
// dengan prop yang sama.

import {
  useCallback,
  useEffect,
  useId,
  useRef,
  useState,
  type ReactNode,
} from "react";

import { ApiFailure } from '../../klien'
import { klasifikasiGalat } from "../../lib/keadaanGalat";
import { useBahasaUI, useTeksUI } from "./bahasaUI";

import { keInputTanggal, dariInputTanggal } from "../../lib/tanggalInput";

export function Field({
  label,
  value,
  onChange,
  error,
  required,
  type = "text",
  autoFocus,
  placeholder,
  readOnly,
  ikon,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  error?: string;
  required?: boolean;
  type?: string;
  autoFocus?: boolean;
  /** Bentuk yang diharapkan, mis. DD-MM-YYYY untuk tanggal (NFR-14). */
  placeholder?: string;
  /** Jejak simpan / identitas baris: DITAMPILKAN, tidak dapat diisi
   *  (Q-50/Q-43b, r165). `readOnly`, bukan `disabled`: nilai yang
   *  di-disable tidak terbaca pembaca layar dan tidak ikut terkirim. */
  readOnly?: boolean;
  /**
   * Glif pemandu DI DALAM medan (§2.2 ronde 272, B.7).
   *
   * ⚠️ OPSIONAL, dan itu bukan kenyamanan — ia pagar biaya. Komponen ini
   * dipakai puluhan layar; memasang pembungkusnya tanpa syarat akan
   * menggeser SETIAP medan di aplikasi, dan tidak satu pun uji akan
   * berbunyi karena yang berubah hanya jarak. Tanpa `ikon`, keluaran
   * fungsi ini sama persis dengan sebelum ronde 272.
   */
  ikon?: ReactNode;
}) {
  const kotak = (
    <input
      className={
        "field__input" +
        (error ? " field__input--error" : "") +
        (readOnly ? " field__input--readonly" : "")
      }
      type={type}
      value={value}
      placeholder={placeholder}
      autoFocus={autoFocus}
      readOnly={readOnly}
      onChange={(e) => onChange(e.target.value)}
    />
  );
  return (
    <div className="field">
      <label className="field__label">
        {label}
        {required && <span className="field__req">*</span>}
      </label>
      {ikon ? (
        <div className="field__wrap-ikon">
          <span className="field__ikon">{ikon}</span>
          {kotak}
        </div>
      ) : (
        kotak
      )}
      {error && <div className="field__error">{error}</div>}
    </div>
  );
}

/**
 * Kotak SANDI dengan tombol mata (ronde 246/247, permintaan langsung
 * pemilik proyek 19-09-2026) — bawaan tersembunyi, tombol mengungkap.
 *
 * Komponen TERPISAH dari `Field`, bukan menambah prop show/hide ke
 * `Field`: `Field` dipakai 100+ titik panggil di seluruh aplikasi untuk
 * field teks/angka biasa, sementara `type="password"` HANYA dipakai SATU
 * tempat (form Login). Menambah logika visibilitas ke `Field` akan
 * membebani seratusan pemanggil lain yang tidak pernah memakainya.
 *
 * Struktur meniru `Field` persis (div.field > label > kontrol > error)
 * supaya CSS `.field`/`.field__label`/`.field__input`/`.field__error`
 * yang sudah ada tetap berlaku tanpa duplikasi aturan.
 */
export function FieldSandi({
  label,
  value,
  onChange,
  error,
  required,
  autoFocus,
  ikon,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  error?: string;
  required?: boolean;
  autoFocus?: boolean;
  /** Glif pemandu di dalam medan — lihat catatan yang sama pada `Field`. */
  ikon?: ReactNode;
}) {
  const [terlihat, setTerlihat] = useState(false);
  return (
    <div className="field">
      <label className="field__label">
        {label}
        {required && <span className="field__req">*</span>}
      </label>
      {/* Dua kelas DIGABUNG, bukan kelas ketiga yang dikarang:
          `field__sandi-wrap` sudah memberi `position: relative` + ruang
          KANAN untuk tombol mata, dan `field__wrap-ikon` menambah ruang
          KIRI untuk glifnya. Medan sandi butuh keduanya. */}
      <div className={"field__sandi-wrap" + (ikon ? " field__wrap-ikon" : "")}>
        {ikon && <span className="field__ikon">{ikon}</span>}
        <input
          className={"field__input" + (error ? " field__input--error" : "")}
          type={terlihat ? "text" : "password"}
          value={value}
          autoFocus={autoFocus}
          onChange={(e) => onChange(e.target.value)}
        />
        <button
          type="button"
          className="field__toggle-sandi"
          aria-label={terlihat ? "Sembunyikan sandi" : "Tampilkan sandi"}
          title={terlihat ? "Sembunyikan sandi" : "Tampilkan sandi"}
          onClick={() => setTerlihat((v) => !v)}
        >
          {terlihat ? <IkonMataCoret /> : <IkonMata />}
        </button>
      </div>
      {error && <div className="field__error">{error}</div>}
    </div>
  );
}

/**
 * Kotak TANGGAL — pemilih tanggal, bukan ketik bebas (ronde 117).
 *
 * Perintah pemilik proyek: seluruh input tanggal di aplikasi memakai
 * `<input type="date">`. Ini SATU-SATUNYA tempat konversi formatnya hidup;
 * konversi tanggal yang disalin ke tiap form akan berbeda suatu hari, dan
 * yang berbeda adalah TANGGAL — bidang yang selisih satu harinya tidak
 * terlihat sampai seseorang merekonsiliasi.
 *
 * # Batas yang dijaga
 *
 * `value` dan `onChange` tetap berbicara dalam bentuk KABEL (DD-MM-YYYY,
 * NFR-14). Pemanggil TIDAK perlu tahu bahwa kontrolnya kini `type="date"`,
 * dan tidak satu pun kontrak API, kolom database, maupun format tampilan
 * tabel berubah karenanya. Yang berubah hanya kontrol input-nya.
 *
 * # Kenapa `required` di sini hanya PETUNJUK
 *
 * Sama seperti `Field` dan `Area`: bintangnya penanda, bukan gerbang.
 * Penolakan wajib-isi dan rentang datang dari backend (`offer.ParseDate`,
 * BR-29d/BR-29e). Menggandakan aturannya di React menghasilkan dua sumber
 * kebenaran yang akan berselisih.
 */
export function FieldTanggal({
  label,
  value,
  onChange,
  error,
  required,
  autoFocus,
}: {
  label: string;
  /** Bentuk KABEL (DD-MM-YYYY), atau bentuk existing yang masih bercampur. */
  value: string;
  /** Menerima bentuk KABEL (DD-MM-YYYY); kosong berarti dikosongkan. */
  onChange: (v: string) => void;
  error?: string;
  required?: boolean;
  autoFocus?: boolean;
}) {
  return (
    <div className="field">
      <label className="field__label">
        {label}
        {required && <span className="field__req">*</span>}
      </label>
      <input
        className={"field__input" + (error ? " field__input--error" : "")}
        type="date"
        value={keInputTanggal(value)}
        autoFocus={autoFocus}
        onChange={(e) => onChange(dariInputTanggal(e.target.value))}
      />
      {error && <div className="field__error">{error}</div>}
    </div>
  );
}

/**
 * Kotak teks panjang, untuk `Exclusions` dan `Special Conditions` yang di
 * existing memang satu `pxTextArea` (UI_EXISTING_REFERENSI §7.2, §7.3).
 */
export function Area({
  label,
  value,
  onChange,
  error,
  required,
  baris = 4,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  error?: string;
  /**
   * Menandai label dengan bintang. Nilainya datang dari daftar field wajib
   * BACKEND (`/offers/meta`), bukan dari daftar yang ditulis di React —
   * BR-29d/BR-29e menempatkan otoritas wajib-isi di backend.
   *
   * Bintang di sini PETUNJUK, bukan gerbang: tombol Simpan tetap hidup dan
   * penolakan tetap datang dari backend beserta pesan per fieldnya.
   */
  required?: boolean;
  baris?: number;
}) {
  return (
    <div className="field field--lebar">
      <label className="field__label">
        {label}
        {required && <span className="field__req">*</span>}
      </label>
      <textarea
        className={"field__input" + (error ? " field__input--error" : "")}
        rows={baris}
        value={value}
        onChange={(e) => onChange(e.target.value)}
      />
      {error && <div className="field__error">{error}</div>}
    </div>
  );
}

/** Satu pilihan pada `Pilih`. */
export interface Opsi {
  value: string;
  label: string;
}

/**
 * Daftar pilihan, untuk field yang di existing memang pxDropdown —
 * mis. `Treaty Group` (§4.8) dan `Reinsurance Type` (§4.6).
 *
 * Pilihan datang dari endpoint daftar referensi; TIDAK ada daftar yang
 * ditulis di dalam React. Sumber Pega tidak memuat isi `pyPromptValues`
 * (§11), sehingga daftar yang ditulis di sini pasti karangan.
 *
 * Nilai yang sudah tersimpan tetapi tidak ada di daftar TIDAK dibuang: ia
 * ditambahkan sebagai pilihan terpisah. Membuangnya akan mengubah data
 * pengguna tanpa memberi tahu, hanya karena daftar referensinya berubah.
 */
export function Pilih({
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
  opsi: Opsi[];
  error?: string;
  required?: boolean;
  kosong?: string;
}) {
  const teksUI = useTeksUI();
  const teksKosong = kosong ?? teksUI.pilihKosong;
  const asing = value !== "" && !opsi.some((o) => o.value === value);
  return (
    <div className="field">
      <label className="field__label">
        {label}
        {required && <span className="field__req">*</span>}
      </label>
      <select
        className={"field__input" + (error ? " field__input--error" : "")}
        value={value}
        onChange={(e) => onChange(e.target.value)}
      >
        <option value="">{teksKosong}</option>
        {asing && (
          <option value={value}>{value} {teksUI.tidakDiDaftar}</option>
        )}
        {opsi.map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
      {error && <div className="field__error">{error}</div>}
    </div>
  );
}

/**
 * Menampilkan kegagalan APA ADANYA.
 *
 * Pesan berasal dari envelope backend: `message` beserta SELURUH field yang
 * gagal (BR-29e) — bukan yang pertama saja. Tidak ada pesan buatan React yang
 * menggantikannya; pesan buatan sendiri menyembunyikan alasan sebenarnya dan
 * membuat pengguna menebak.
 *
 * Kegagalan yang BUKAN envelope backend (jaringan mati, body tak terbaca)
 * juga ditampilkan apa adanya, bukan diubah menjadi "terjadi kesalahan".
 */
export function Gagal({ galat }: { galat: unknown }) {
  const bahasa = useBahasaUI();
  const teks = useTeksUI();
  const k = klasifikasiGalat(galat, bahasa);
  if (!k) return null;

  /* Backend mati mendapat panel TERSENDIRI, bukan alert merah biasa:
     penyebabnya bukan data dan bukan permintaan pemakai, dan yang ia
     butuhkan bukan "coba lagi" melainkan satu perintah konkret. Tanpa
     pembedaan ini, layar kosong terbaca sebagai "belum ada data" — salah
     diagnosa yang sudah terjadi (ronde 88). */
  if (k.jenis === "backend-mati") {
    return (
      <div className="alert alert--warn" role="alert">
        <strong>{k.pesan}</strong>
        {k.petunjuk && (
          <p style={{ marginTop: 8 }}>
            <code>{k.petunjuk}</code>
          </p>
        )}
        <p style={{ marginTop: 8 }}>
          <button
            className="btn btn--ghost btn--sm"
            onClick={() => window.location.reload()}
          >
            {teks.muatUlang}
          </button>
        </p>
      </div>
    );
  }

  if (galat instanceof ApiFailure && galat.fields.length > 0) {
    return (
      <div className="alert alert--error">
        {k.pesan}
        <ul>
          {galat.fields.map((f) => (
            <li key={f.field}>{f.message}</li>
          ))}
        </ul>
      </div>
    );
  }
  return <div className="alert alert--error">{k.pesan}</div>;
}

export function BelumTersedia({ apa }: { apa: string }) {
  // KETERANGAN TEKNISNYA TIDAK DITAMPILKAN (keputusan pemilik proyek ronde
  // 67). Nama tabel, nomor migrasi, dan kolom yang hilang adalah bahasa
  // pengembang; pemakai layar ini tidak dapat berbuat apa pun dengannya.
  // Sebabnya tetap SAMPAI ke pengembang — backend sudah mencatatnya di log
  // server pada saat penolakan terjadi.
  //
  // Parameter `sebab` DIBUANG ronde 70: ia diterima lalu diabaikan, dan
  // parameter yang diterima tetapi tidak dipakai membuat pemanggil percaya
  // ia berpengaruh.
  // Kata "belum tersedia" DIBUANG (butir 13, ronde 72): ia bercerita
  // tentang keadaan pembangunan sistem, dan pengguna tidak dapat berbuat
  // apa pun atasnya.
  //
  // Yang TIDAK boleh dilakukan sebagai gantinya: menampilkan tabel kosong
  // atau kalimat "tidak ada data". Penolakan sistem dan data yang memang
  // kosong terlihat sama di layar tetapi berarti berlawanan — pengguna yang
  // membaca "tidak ada data" akan menyimpulkan tidak ada realisasi pada
  // periode itu, dan menyimpulkan itu tentang uang. Panel ini karena itu
  // tetap ada; yang berubah hanya bahasanya, menjadi keadaan + tindakan.
  return (
    <div className="alert alert--info">
      <strong>{apa} tidak dapat dimuat saat ini.</strong>
      {/* <div style={{ marginTop: 6 }}>
        Coba beberapa saat lagi. Bila tetap berulang, hubungi tim aplikasi —
        data Anda tidak terpengaruh.
      </div> */}
    </div>
  );
}

/**
 * Kerangka SATU popup, dipakai SELURUH modal aplikasi.
 *
 * Ia ada karena tiga jalan keluar itu tidak boleh ditulis ulang per modal.
 * Sebelumnya setiap modal menyusun backdrop dan tombolnya sendiri, dan yang
 * terjadi pada susunan seperti itu selalu sama: satu modal punya tombol X,
 * yang lain tidak; satu menutup dengan Escape, yang lain memaksa tetikus.
 * Di sini ketiganya ada SEKALI, sehingga tidak ada modal yang dapat
 * kehilangannya:
 *
 *   1. tombol X di pojok kanan atas
 *   2. klik backdrop
 *   3. tombol Batal di kaki modal
 *
 * Escape ditambahkan sebagai jalan KEEMPAT — bukan pengganti salah satu di
 * atas. Modal yang hanya dapat ditutup dengan tetikus menghambat pemakai yang
 * bekerja dari papan tuts, dan pada layar pengisian data itu justru pemakai
 * yang paling sering membuka modal.
 */
export function Modal({
  judul,
  onTutup,
  onKirim,
  aksi,
  labelBatal = "Cancel",
  lebar,
  penuh,
  children,
}: {
  judul: string;
  onTutup: () => void;
  /**
   * Bila diisi, isi modal dirender sebagai `<form>` sehingga Enter di kotak
   * isian mengirimnya. Modal yang hanya menampilkan keterangan tidak perlu
   * form, dan membuatnya form berarti menjanjikan pengiriman yang tidak ada.
   */
  onKirim?: () => void;
  /** Tombol utama di kaki modal. Tombol Batal sudah disediakan di sini. */
  aksi?: ReactNode;
  labelBatal?: string;
  /**
   * Modal LEBAR, untuk isi berfield banyak yang memakai `.form-grid`.
   *
   * Lebar bawaan 460px cukup untuk satu kolom. Dua kolom di dalam 460px
   * menghasilkan kotak isian yang lebih pendek daripada isi yang diketik ke
   * dalamnya — jadi yang melebar adalah modalnya, bukan gridnya yang
   * dipaksa menyempit.
   */
  lebar?: boolean;
  /**
   * Modal SELEBAR LAYAR, untuk isi yang membelah dua panel berdampingan.
   *
   * Berbeda dari `lebar`, yang berhenti di 1100px. Angka itu cukup untuk
   * form berkolom banyak, tetapi tidak untuk DUA panel yang masing-masing
   * memuat tabel: panel kanan tersisa sekitar 700px sementara sebagian
   * gridnya berisi tujuh kolom.
   *
   * Ronde 268 (§1) — satu-satunya pemakainya List Description, dan
   * permintaannya memang "LAYAR PENUH SPLIT JADI 2".
   */
  penuh?: boolean;
  children: ReactNode;
}) {
  const idJudul = useId();
  const teksModal = useTeksUI();

  // Ronde 240 — transisi KELUAR. `onTutup` sesungguhnya (yang meng-unmount
  // modal ini di pemanggil) ditunda 140ms supaya `.modal--keluar`/
  // `.modal__backdrop--keluar` (styles.css) sempat terlihat, bukan lenyap
  // seketika. Guard `menutup` mencegah timer ganda bila pengguna menekan
  // X lalu Escape sebelum unmount sungguhan terjadi.
  const [menutup, setMenutup] = useState(false);
  // Guard timer lewat ref, BUKAN di dalam updater `setMenutup` (temuan
  // code-review ronde 241): StrictMode (aktif, main.tsx) memanggil
  // updater state DUA KALI di development justru untuk menangkap efek
  // samping yang bersembunyi di sana — timer yang dijadwalkan di dalam
  // updater akan terjadwal dua kali dan memanggil `onTutup` dua kali.
  // Ref yang sama juga membersihkan timer saat Modal lenyap lewat jalur
  // lain sebelum 140ms berlalu.
  const timerTutup = useRef<number | undefined>(undefined);
  const mulaiTutup = useCallback(() => {
    setMenutup(true);
    if (timerTutup.current !== undefined) return;
    timerTutup.current = window.setTimeout(onTutup, 140);
  }, [onTutup]);
  useEffect(
    () => () => {
      if (timerTutup.current !== undefined)
        window.clearTimeout(timerTutup.current);
    },
    [],
  );

  // Escape menutup modal. Pendengarnya dipasang di document, bukan di kotak
  // modalnya: fokus dapat berada di mana saja di dalam modal — termasuk di
  // kotak isian yang baru dibuat — dan pendengar yang menempel pada satu
  // elemen akan melewatkan tekanan tombol yang tidak mengenai elemen itu.
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") mulaiTutup();
    };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [mulaiTutup]);

  const isi = (
    <>
      <div className="modal__head">
        <h3 className="modal__title" id={idJudul}>
          {judul}
        </h3>
        {/* Tombol tutup. `type="button"` WAJIB: di dalam <form> tombol tanpa
            type adalah tombol kirim, sehingga X justru akan menyimpan. */}
        <button
          type="button"
          className="modal__close"
          aria-label={teksModal.tutup}
          title={teksModal.tutupEsc}
          onClick={mulaiTutup}
        >
          <IkonTutup />
        </button>
      </div>

      <div className="modal__body">{children}</div>

      <div className="modal__actions">
        <button type="button" className="btn btn--ghost" onClick={mulaiTutup}>
          {labelBatal}
        </button>
        {aksi}
      </div>
    </>
  );

  // Klik backdrop menutup. Yang diperiksa adalah SASARAN klik, bukan
  // stopPropagation di kotak modal: menahan propagasi di kotak modal
  // mematikan pula klik lain yang kelak dipasang di atas backdrop, dan
  // matinya tidak berbunyi.
  const backdrop = (anak: ReactNode) => (
    <div
      className={
        "modal__backdrop" + (menutup ? " modal__backdrop--keluar" : "")
      }
      role="presentation"
      onClick={(e) => {
        if (e.target === e.currentTarget) mulaiTutup();
      }}
    >
      {anak}
    </div>
  );

  if (onKirim) {
    return backdrop(
      <form
        className={
          "modal" +
          (lebar ? " modal--lebar" : "") +
          (penuh ? " modal--penuh" : "") +
          (menutup ? " modal--keluar" : "")
        }
        role="dialog"
        aria-modal="true"
        aria-labelledby={idJudul}
        onSubmit={(e) => {
          e.preventDefault();
          onKirim();
        }}
      >
        {isi}
      </form>,
    );
  }

  return backdrop(
    <div
      className={
        "modal" +
        (lebar ? " modal--lebar" : "") +
        (penuh ? " modal--penuh" : "") +
        (menutup ? " modal--keluar" : "")
      }
      role="dialog"
      aria-modal="true"
      aria-labelledby={idJudul}
    >
      {isi}
    </div>,
  );
}

/**
 * Modal penyunting SATU BARIS grid.
 *
 * Grid tetap read-only dan tombol `Ubah` membuka form ini (keputusan pemilik
 * proyek; pola existing §3.2 juga demikian). Tombolnya `Terapkan`, BUKAN
 * `Simpan`: himpunan diganti utuh lewat satu PUT, sehingga baris yang
 * diterapkan belum tersimpan sampai `Simpan` di layar induk ditekan.
 * Menamainya `Simpan` akan berbohong tentang apa yang sudah terjadi.
 *
 * Bentuk popupnya sendiri milik `Modal` — X, backdrop, Batal, dan Escape
 * datang dari sana, sehingga keenam pemanggil modal baris ini tidak dapat
 * berbeda satu dengan yang lain.
 */
export function ModalBaris({
  judul,
  onTutup,
  onTerapkan,
  children,
}: {
  judul: string;
  onTutup: () => void;
  onTerapkan: () => void;
  children: ReactNode;
}) {
  return (
    <Modal
      judul={judul}
      onTutup={onTutup}
      onKirim={onTerapkan}
      aksi={<button className="btn btn--primary">Apply</button>}
    >
      {children}
    </Modal>
  );
}

/**
 * Strip tab.
 *
 * Tab dipakai karena existing memang memakai tab pada layar ini, dan hanya
 * pada layar ini beserta Exclusion (§7.1): delapan tab Treaty In Offer.
 * Nama tab di sini VERBATIM dari §7.2/§7.3.
 */
export function StripTab<T extends string>({
  tab,
  aktif,
  onPilih,
  label,
}: {
  tab: readonly T[];
  aktif: T;
  onPilih: (t: T) => void;
  /**
   * Bentuk TAMPIL sebuah tab. Tanpa ini nilainya sendiri yang tampil —
   * benar untuk tab yang nilainya sudah berupa nama ("Portfolio"), keliru
   * untuk tab yang nilainya kode alur ("STAGE1_ADMIN").
   */
  label?: (t: T) => string;
}) {
  return (
    <div className="tabs" role="tablist">
      {tab.map((t) => (
        <button
          key={t}
          type="button"
          role="tab"
          aria-selected={t === aktif}
          className={"tabs__item" + (t === aktif ? " tabs__item--aktif" : "")}
          onClick={() => onPilih(t)}
        >
          {label ? label(t) : t}
        </button>
      ))}
    </div>
  );
}

/**
 * Blok bertitel yang bertumpuk — bukan tab.
 *
 * Existing memakai pola ini untuk `Status`, `Summary of Limit`,
 * `Total All Layers`, dan belasan panel lain (§7.5), jadi pengelompokan form
 * penawaran mengikutinya alih-alih menambah tab yang tidak ada di existing.
 */
export function Panel({
  judul,
  catatan,
  children,
}: {
  judul: string;
  catatan?: string;
  children: ReactNode;
}) {
  return (
    <section className="panel">
      <h4 className="panel__title">{judul}</h4>
      {catatan && <div className="panel__note">{catatan}</div>}
      {children}
    </section>
  );
}

// =========================================================================
// IKON
// =========================================================================
//
// Project ini TIDAK punya pustaka ikon, dan tidak boleh menambah dependensi.
// Karena itu bentuknya ditulis sebagai SVG sebaris di sini: satu pustaka ikon
// membawa ratusan bentuk untuk empat yang benar-benar dipakai, dan keempatnya
// hanya berupa beberapa garis.
//
// Seluruh ikon memakai `currentColor`, sehingga warnanya ditentukan CSS di
// tempat pemakaian — bukan ditanam di dalam berkas ini. Tidak ada satu pun
// warna merek yang ditulis di sini.

/** Sifat bersama seluruh ikon garis 24x24. */
function sifatIkon(ukuran: number) {
  return {
    width: ukuran,
    height: ukuran,
    viewBox: "0 0 24 24",
    fill: "none",
    stroke: "currentColor",
    strokeWidth: 2,
    strokeLinecap: "round" as const,
    strokeLinejoin: "round" as const,
    "aria-hidden": true,
  };
}

/** Silang penutup modal. */
export function IkonTutup({ ukuran = 16 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)}>
      <path d="M18 6 6 18M6 6l12 12" />
    </svg>
  );
}

/**
 * Tombol lipat panel samping.
 *
 * Bentuknya panel bergaris kiri, bukan tiga garis mendatar: tiga garis
 * mendatar di layar lebar berarti "buka menu yang tersembunyi", sementara di
 * sini panelnya TIDAK tersembunyi — ia hanya menyempit. Pada layar sempit,
 * tempat panel memang tersembunyi, bentuk ini tetap terbaca sebagai panel.
 */
export function IkonPanel({ ukuran = 17 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.8}>
      <rect x="3" y="4" width="18" height="16" rx="2" />
      <path d="M9 4v16" />
    </svg>
  );
}

/** Mata terbuka — tombol "Tampilkan sandi" (ronde 246/247). */
export function IkonMata({ ukuran = 18 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)}>
      <path d="M2 12s3.6-7 10-7 10 7 10 7-3.6 7-10 7-10-7-10-7Z" />
      <circle cx="12" cy="12" r="3" />
    </svg>
  );
}

/** Mata dicoret — tombol "Sembunyikan sandi", sandi SEDANG terlihat. */
export function IkonMataCoret({ ukuran = 18 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)}>
      <path d="M2 12s3.6-7 10-7 10 7 10 7-3.6 7-10 7-10-7-10-7Z" />
      <circle cx="12" cy="12" r="3" />
      <path d="M3 3l18 18" />
    </svg>
  );
}

/**
 * Tiga ikon §2.2 ronde 272 (B.7) — pemandu medan Login + glif tombolnya.
 *
 * Ronde 270 `#1` mencatat ketiganya sebagai yang MASIH kurang terhadap
 * acuan `REFERENCE_UI/FORM_LOGIN_E-FACULTATIVE.png`, dan B.7 menuntut
 * diukur ulang alih-alih dipercaya: diperiksa 23-09-2026, `<svg>` di
 * kartu Login berjumlah NOL. Jadi memang masih kurang.
 *
 * Ketiganya memakai `sifatIkon` yang sudah ada — `currentColor`,
 * `aria-hidden`, viewBox 24. Nol paket npm, nol berkas gambar, dan
 * warnanya mengikuti teks di sekitarnya sehingga nol hex baru.
 *
 * `aria-hidden` SENGAJA, dan `aria-label` sengaja TIDAK ada: ketiganya
 * berdiri di samping label yang sudah terbaca ("Login", "Password",
 * "Sign in"). Memberi mereka nama membuat pembaca layar mengucapkan
 * labelnya dua kali.
 */

/** Sosok orang — pemandu medan nama pengguna. */
export function IkonPengguna({ ukuran = 17 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.8}>
      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
      <circle cx="12" cy="7" r="4" />
    </svg>
  );
}

/** Gembok — pemandu medan kata sandi. */
export function IkonKunci({ ukuran = 17 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.8}>
      <rect x="4" y="11" width="16" height="10" rx="2" />
      <path d="M8 11V7a4 4 0 0 1 8 0v4" />
    </svg>
  );
}

/** Panah masuk pintu — glif tombol utama Login. */
export function IkonMasuk({ ukuran = 16 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.9}>
      <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4" />
      <path d="M10 17l5-5-5-5M15 12H3" />
    </svg>
  );
}

/**
 * Kaca pembesar — pemandu kendali cari di topbar (TV2-442, ronde 277).
 *
 * Tinggal DI SINI bersama delapan saudaranya, bukan sebagai `<svg>`
 * sekali-pakai di `App.tsx`: `sifatIkon` memegang viewBox, ketebalan
 * goresan, dan `aria-hidden` yang sama untuk semuanya. Glif yang
 * digambar di tempat lain menyimpang ukurannya tanpa ada yang melihat.
 */
export function IkonCari({ ukuran = 15 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.9}>
      <circle cx="11" cy="11" r="7" />
      <path d="M20 20l-4.3-4.3" />
    </svg>
  );
}

/*
 * Ikon bingkai aplikasi — desain `workpage-template.html` (29-09-2026).
 *
 * Template memuat Lucide dari CDN; di sini bentuknya DISALIN dari
 * `lucide-static@1.48.0` (lisensi ISC) — versi yang sama dengan template —
 * supaya aplikasi tetap nol paket ikon dan tetap jalan tanpa internet.
 * Ketebalan 1.8 mengikuti `svg.lucide { stroke-width: 1.8 }` di template.
 */

/** `house` — butir Beranda. */
export function IkonRumah({ ukuran = 20 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.8}>
      <path d="M15 21v-8a1 1 0 0 0-1-1h-4a1 1 0 0 0-1 1v8" />
      <path d="M3 10a2 2 0 0 1 .709-1.528l7-6a2 2 0 0 1 2.582 0l7 6A2 2 0 0 1 21 10v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
    </svg>
  );
}

/** `menu` — tombol ciutkan/buka menu di topbar. */
export function IkonMenu({ ukuran = 20 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.8}>
      <path d="M4 5h16M4 12h16M4 19h16" />
    </svg>
  );
}

/** `chevron-down` — kelompok menu terbuka/terlipat, pemicu menu profil. */
export function IkonChevron({ ukuran = 16 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.8}>
      <path d="m6 9 6 6 6-6" />
    </svg>
  );
}

/** `inbox` — kartu tahap Input Register di Beranda. */
export function IkonKotakMasuk({ ukuran = 16 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.8}>
      <polyline points="22 12 16 12 14 15 10 15 8 12 2 12" />
      <path d="M5.45 5.11 2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z" />
    </svg>
  );
}

/** `hourglass` — kartu tahap Outstanding Claim di Beranda. */
export function IkonJamPasir({ ukuran = 16 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.8}>
      <path d="M5 22h14M5 2h14" />
      <path d="M17 22v-4.172a2 2 0 0 0-.586-1.414L12 12l-4.414 4.414A2 2 0 0 0 7 17.828V22" />
      <path d="M7 2v4.172a2 2 0 0 0 .586 1.414L12 12l4.414-4.414A2 2 0 0 0 17 6.172V2" />
    </svg>
  );
}

/** `stethoscope` — kartu tahap Medical Check di Beranda. */
export function IkonStetoskop({ ukuran = 16 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.8}>
      <path d="M11 2v2M5 2v2" />
      <path d="M5 3H4a2 2 0 0 0-2 2v4a6 6 0 0 0 12 0V5a2 2 0 0 0-2-2h-1" />
      <path d="M8 15a6 6 0 0 0 12 0v-3" />
      <circle cx="20" cy="10" r="2" />
    </svg>
  );
}

/** `file-search` — kartu tahap Claim Analis di Beranda. */
export function IkonBerkasCari({ ukuran = 16 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.8}>
      <path d="M6 22a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h8a2.4 2.4 0 0 1 1.704.706l3.588 3.588A2.4 2.4 0 0 1 20 8v12a2 2 0 0 1-2 2z" />
      <path d="M14 2v5a1 1 0 0 0 1 1h5" />
      <circle cx="11.5" cy="14.5" r="2.5" />
      <path d="M13.3 16.3 15 18" />
    </svg>
  );
}

/** `sun` — tombol tema di topbar saat tema gelap (tekan = terang). */
export function IkonMatahari({ ukuran = 20 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.8}>
      <circle cx="12" cy="12" r="4" />
      <path d="M12 2v2M12 20v2m-7.07-17.07 1.41 1.41m11.32 11.32 1.41 1.41M2 12h2M20 12h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41" />
    </svg>
  );
}

/** `moon` — tombol tema di topbar saat tema terang (tekan = gelap). */
export function IkonBulan({ ukuran = 20 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.8}>
      <path d="M20.985 12.486a9 9 0 1 1-9.473-9.472c.405-.022.617.46.402.803a6 6 0 0 0 8.268 8.268c.344-.215.825-.004.803.401" />
    </svg>
  );
}

/** `shield-check` — kartu Peran Anda di Beranda. */
export function IkonPerisai({ ukuran = 16 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.8}>
      <path d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z" />
      <path d="m9 12 2 2 4-4" />
    </svg>
  );
}

/** Kotak kosong, untuk keadaan "belum ada data". */
export function IkonKosong({ ukuran = 30 }: { ukuran?: number }) {
  return (
    <svg {...sifatIkon(ukuran)} strokeWidth={1.5}>
      <path d="M4 14h4l1.5 3h5L16 14h4" />
      <path d="M4 14 6.5 6h11L20 14v4a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2v-4Z" />
    </svg>
  );
}

// =========================================================================
// KEADAAN MEMUAT DAN KEADAAN KOSONG
// =========================================================================

/**
 * Penanda sedang memuat.
 *
 * Ia menggantikan tulisan "Memuat..." yang tidak bergerak. Teks yang tidak
 * bergerak tidak dapat dibedakan dari layar yang macet, dan pemakai yang
 * menyangka layarnya macet akan menekan tombolnya dua kali.
 */
export function Memuat({ pesan }: { pesan?: string }) {
  const teks = useTeksUI();
  return (
    <div className="memuat" role="status">
      <span className="spinner" aria-hidden="true" />
      {pesan ?? teks.memuat}
    </div>
  );
}

/**
 * Keadaan kosong: ikon sederhana beserta keterangan.
 *
 * `pesan` menyebut APA yang belum ada, `petunjuk` menyebut apa yang dapat
 * dilakukan. Keduanya dipisah karena keduanya dibaca pada saat berbeda —
 * yang pertama sekilas, yang kedua ketika pemakai memutuskan bertindak.
 */
export function Kosong({
  pesan,
  petunjuk,
}: {
  pesan: string;
  petunjuk?: string;
}) {
  return (
    <div className="kosong">
      <span className="kosong__ikon" aria-hidden="true">
        <IkonKosong />
      </span>
      <div className="kosong__pesan">{pesan}</div>
      {petunjuk && <div className="kosong__petunjuk">{petunjuk}</div>}
    </div>
  );
}

/**
 * Baris tabel palsu selagi data dimuat.
 *
 * Bentuknya mengikuti bentuk tabel yang akan datang, sehingga tinggi dan
 * lebar kolom tidak melompat ketika data tiba. Satu baris "Memuat..." di
 * tengah tabel selalu melompat, karena tinggi satu baris tidak sama dengan
 * tinggi daftarnya.
 *
 * Jumlah barisnya semata pengisi ruang; ia TIDAK mengaku sebagai jumlah data
 * yang akan datang.
 */
export function BarisSkeleton({
  kolom,
  baris = 5,
}: {
  kolom: number;
  baris?: number;
}) {
  return (
    <>
      {Array.from({ length: baris }, (_, i) => (
        <tr key={i} aria-hidden="true">
          {Array.from({ length: kolom }, (_, j) => (
            <td key={j}>
              <span className="skeleton" />
            </td>
          ))}
        </tr>
      ))}
    </>
  );
}

/**
 * Besar halaman SERAGAM untuk seluruh daftar.
 *
 * Sebelum ronde 69 ada empat besaran berbeda di empat layar, sehingga
 * "halaman 3" berarti baris yang berbeda tergantung menu mana yang dibuka.
 */
export const UKURAN_HALAMAN = 20;

/**
 * Kontrol halaman — satu bentuk untuk seluruh daftar.
 *
 * # Kenapa ini ada
 *
 * Enam belas menu Treaty Master merender lewat satu komponen tabel yang
 * meminta halaman pertama saja. Akibatnya terukur: 16.686 baris nyata, hanya
 * 304 yang dapat dicapai dari layar. Angka "5807 record" tetap tampil di
 * toolbar — memberi tahu pemakai bahwa datanya ada, tanpa memberi jalan
 * untuk mencapainya.
 *
 * `total` boleh tidak diketahui (server yang hanya mengirim "masih ada
 * halaman berikutnya"). Dalam keadaan itu jangkauan dan jumlah halaman TIDAK
 * ditampilkan — menghitungnya dari panjang halaman berarti mengarang angka.
 */
export function Halaman({
  halaman,
  ukuran,
  total,
  adaLagi,
  onPindah,
}: {
  halaman: number;
  ukuran: number;
  total?: number;
  adaLagi?: boolean;
  onPindah: (h: number) => void;
}) {
  const tahu = typeof total === "number";
  const totalHalaman = tahu ? Math.max(1, Math.ceil(total! / ukuran)) : 0;
  const kini = tahu ? Math.min(halaman, totalHalaman) : halaman;
  const awal = total === 0 ? 0 : (kini - 1) * ukuran + 1;
  const akhir = tahu ? Math.min(kini * ukuran, total!) : 0;
  const bisaMaju = tahu ? kini < totalHalaman : !!adaLagi;
  const teks = useTeksUI();

  return (
    <div className="pager">
      <span className="muted">
        {tahu
          ? total === 0
            ? teks.tidakAdaBaris
            : teks.menampilkan(awal, akhir, total!)
          : teks.halaman(kini)}
      </span>
      <div className="pager__aksi">
        <button
          className="btn btn--ghost btn--sm"
          disabled={kini <= 1}
          onClick={() => onPindah(kini - 1)}
        >
          {teks.sebelumnya}
        </button>
        {tahu && <span className="muted">{teks.halamanDari(kini, totalHalaman)}</span>}
        <button
          className="btn btn--ghost btn--sm"
          disabled={!bisaMaju}
          onClick={() => onPindah(kini + 1)}
        >
          {teks.berikutnya}
        </button>
      </div>
    </div>
  );
}
