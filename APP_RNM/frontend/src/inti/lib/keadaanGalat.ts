/**
 * Klasifikasi galat layar — membedakan tiga keadaan yang MIRIP di mata
 * pemakai tetapi menuntut tindakan yang sama sekali berbeda.
 *
 * | Keadaan | Yang dilihat pemakai sebelumnya | Yang harus ia lakukan |
 * |---|---|---|
 * | backend mati | "Belum ada data" | jalankan backend-nya |
 * | tabel belum ada | "Belum ada data" | tunggu DDL/migrasi |
 * | benar-benar kosong | "Belum ada data" | boleh menambah record |
 *
 * Ketiganya dahulu tampil sebagai kalimat yang SAMA, dan yang pertama
 * bahkan tanpa satu pun tanda merah: `App.tsx` menelan galat yang bukan
 * `ApiFailure`, dan galat jaringan memang bukan ApiFailure. Akibatnya
 * "Limit Treaty In kosong" tidak dapat dibedakan dari "backend tidak
 * berjalan" — dan salah diagnosa itu sudah terjadi.
 *
 * Modul ini MURNI (tanpa React, tanpa fetch, tanpa storage) supaya dapat
 * diuji langsung; komponen `Gagal` hanya merendernya.
 */

import type { Bahasa } from "./teksUI";

/** Bentuk minimal ApiFailure yang dibutuhkan — sengaja bukan impor kelas,
 *  supaya modul ini tetap murni dan dapat diuji di node. */
interface GalatApiSeperti {
  status: number;
  detail: { code?: string; message?: string };
}

export type JenisGalat =
  | "backend-mati"
  | "belum-tersedia"
  | "galat-api"
  | "tidak-dikenal";

export interface KeadaanGalat {
  jenis: JenisGalat;
  pesan: string;
  /** Langkah konkret yang dapat dilakukan pembaca; kosong bila tidak ada. */
  petunjuk?: string;
}

function miripApiFailure(g: unknown): g is GalatApiSeperti {
  if (typeof g !== "object" || g === null) return false;
  const o = g as Record<string, unknown>;
  return typeof o.status === "number" && typeof o.detail === "object" &&
    o.detail !== null;
}

/**
 * Galat jaringan: `fetch` menolak SEBELUM ada respons.
 *
 * Peramban berbeda memberi kalimat berbeda untuk keadaan yang sama —
 * Chrome "Failed to fetch", Firefox "NetworkError when attempting to fetch
 * resource.", Safari "Load failed" — jadi yang diperiksa ketiganya, bukan
 * satu ejaan yang kebetulan dipakai mesin pengembang.
 */
function galatJaringan(g: unknown): boolean {
  if (!(g instanceof TypeError)) return false;
  const p = g.message.toLowerCase();
  return (
    p.includes("failed to fetch") ||
    p.includes("networkerror") ||
    p.includes("load failed") ||
    p.includes("fetch")
  );
}

/** Kode yang `request()` pasang ketika jawaban bukan JSON sama sekali. */
export const KODE_BACKEND_MATI = "BACKEND_TIDAK_TERJANGKAU";

function backendMati(bahasa: Bahasa): KeadaanGalat {
  if (bahasa === "en") {
    return {
      jenis: "backend-mati",
      pesan:
        "Backend not connected (127.0.0.1:8080). The screen below is empty " +
        "NOT because there is no data, but because the request did not " +
        "reach the backend.",
      petunjuk:
        "Start the backend in PowerShell: Set-Location APP_RNM ; " +
        ". .\\muat-env.ps1 ; go run .\\cmd\\api — then reload this " +
        "page. Go does NOT read .env by itself; muat-env.ps1 loads it " +
        "into that window.",
    };
  }
  return {
    jenis: "backend-mati",
    pesan:
      "Backend tidak terhubung (127.0.0.1:8080). Layar di bawah kosong " +
      "BUKAN karena datanya tidak ada, melainkan karena permintaannya " +
      "tidak sampai ke backend.",
    petunjuk:
      "Jalankan backend di PowerShell: Set-Location APP_RNM ; " +
      ". .\\muat-env.ps1 ; go run .\\cmd\\api — lalu muat ulang " +
      "halaman ini. Go TIDAK membaca .env sendiri; muat-env.ps1 yang " +
      "memuatnya ke jendela itu.",
  };
}

/** Teks cadangan penolakan tanpa kalimat backend. */
const DITOLAK_TANPA_PESAN: Readonly<Record<Bahasa, string>> = {
  id: "Permintaan ditolak backend",
  en: "Request rejected by the backend",
};

/** Kode yang `kegagalanDari` (`inti/klien.ts`) pasang pada jawaban JSON backend. */
const KODE_DITOLAK_BACKEND = "DITOLAK_BACKEND";

/**
 * Apakah jawaban ini membawa kalimat backend sendiri di amplop `{"galat": …}`.
 *
 * ⛔ Aturannya SAMA dengan `pesanGalat` di `inti/klien.ts` — hanya kode
 * DITOLAK_BACKEND, pesan tak kosong — disalin, bukan diimpor, karena modul
 * ini murni; `keadaanGalat.test.ts` mengunci keduanya tetap sepakat. Pesan
 * BACKEND_TIDAK_TERJANGKAU ditulis KLIEN (`bukanJSON`), bukan backend.
 */
function membawaGalatBackend(g: GalatApiSeperti): boolean {
  return (
    g.detail.code === KODE_DITOLAK_BACKEND &&
    typeof g.detail.message === "string" &&
    g.detail.message !== ""
  );
}

/**
 * `bahasa` memilih teks BUATAN KLIEN (panel backend mati, teks cadangan);
 * kalimat backend selalu tampil apa adanya. Bawaan `id` (`lib/teksUI.ts`).
 */
export function klasifikasiGalat(galat: unknown, bahasa: Bahasa = "id"): KeadaanGalat | null {
  if (galat == null) return null;

  if (galatJaringan(galat)) return backendMati(bahasa);

  if (miripApiFailure(galat)) {
    const pesan = galat.detail.message ?? DITOLAK_TANPA_PESAN[bahasa];

    /* Backend mati LEWAT PROXY tidak menolak fetch: proxy pengembangan
       (Vite) maupun reverse-proxy produksi menjawab 502/503/504 — atau
       500 dengan badan yang bukan JSON, yang `request()` ubah menjadi
       kode BACKEND_TIDAK_TERJANGKAU. Bagi pemakai keadaannya SAMA dengan
       koneksi ditolak, jadi pesannya pun sama.

       ⛔ 502/503/504 yang MEMBAWA `{"galat": …}` datang dari backend
       sendiri — ia menyala dan menolak dengan alasan (mis. 503 master
       kurs rusak, lanjutan 6 Treaty Contract Out). Menyebutnya "backend
       tidak terhubung" menelan alasannya dan menyuruh pemakai menyalakan
       backend yang sudah menyala.

       ⚠️ Proxy dev Vite 5.4 menjawab upstream mati dengan 500 BERBADAN
       KOSONG — bukan 502/503/504, dan bukan badan bukan-JSON. `kegagalanDari`
       menjadikannya DITOLAK_BACKEND tanpa pesan, sehingga ia jatuh ke
       galat-api "Permintaan ditolak backend" — sebelum dan sesudah
       lanjutan 6. Terbuka, dilaporkan. */
    if (
      galat.detail.code === KODE_BACKEND_MATI ||
      (!membawaGalatBackend(galat) &&
        (galat.status === 502 ||
          galat.status === 503 ||
          galat.status === 504))
    ) {
      return backendMati(bahasa);
    }
    if (galat.detail.code === "NOT_PROVISIONED") {
      return { jenis: "belum-tersedia", pesan };
    }
    return { jenis: "galat-api", pesan };
  }

  return {
    jenis: "tidak-dikenal",
    pesan: galat instanceof Error ? galat.message : String(galat),
  };
}
