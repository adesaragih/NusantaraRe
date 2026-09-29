/**
 * Pagar galat React di AKAR aplikasi.
 *
 * # Kenapa ada (ronde 87)
 *
 * Satu sel grid yang melempar meruntuhkan seluruh pohon React menjadi
 * HALAMAN PUTIH — tanpa pesan, tanpa jalan keluar selain menebak. Yang
 * dilihat pemilik proyek hanyalah layar kosong setelah login yang berhasil.
 *
 * Pagar ini tidak memperbaiki galatnya — ia membuat galat BERBUNYI: pesan
 * yang menyebut apa yang pecah, dan satu tombol muat ulang. Class component
 * dipakai karena React hanya menyediakan componentDidCatch di sana.
 */
import { Component, type ReactNode } from "react";

interface Keadaan {
  galat: Error | null;
}

export class PagarGalat extends Component<{ children: ReactNode }, Keadaan> {
  state: Keadaan = { galat: null };

  static getDerivedStateFromError(galat: Error): Keadaan {
    return { galat };
  }

  render() {
    if (!this.state.galat) {
      return this.props.children;
    }
    return (
      <div style={{ padding: "40px", maxWidth: "640px", margin: "0 auto" }}>
        <h1 style={{ fontSize: "18px", marginBottom: "12px" }}>
          Terjadi kesalahan pada tampilan
        </h1>
        <p style={{ marginBottom: "8px" }}>
          Halaman berhenti karena galat berikut — data Anda tidak hilang,
          dan backend tetap berjalan:
        </p>
        <pre
          style={{
            // Token, bukan hex: kotak terang tetap terang di mode gelap
            // sementara teksnya ikut terang — tak terbaca.
            background: "var(--surface-2)",
            padding: "12px",
            borderRadius: "6px",
            whiteSpace: "pre-wrap",
            marginBottom: "16px",
          }}
        >
          {String(this.state.galat.message || this.state.galat)}
        </pre>
        <button type="button" onClick={() => window.location.reload()}>
          Muat ulang halaman
        </button>
      </div>
    );
  }
}
