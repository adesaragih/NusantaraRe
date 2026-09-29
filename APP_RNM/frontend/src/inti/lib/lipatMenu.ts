/**
 * Keadaan lipatan kelompok sidebar - ronde 103, BAGIAN A.
 *
 * Tiga keputusan, ketiganya teruji di lipatMenu.test.ts:
 *
 *   1. BAWAAN TERBUKA - itulah tampilan sidebar sebelum ronde ini; melipat
 *      semuanya diam-diam berarti menyembunyikan menu dari pengguna lama.
 *   2. Kelompok ber-halaman-AKTIF selalu terbuka, MENGALAHKAN lipatan
 *      tersimpan (A2): pengguna yang memuat ulang tidak boleh mendapati
 *      halamannya "hilang" di balik kelompok terlipat.
 *   3. Penyimpanan tidak dapat dipercaya (A3): localStorage bisa diblokir,
 *      penuh, atau berisi sampah. Semua jalan gagal JATUH ke bawaan tanpa
 *      melempar - sidebar yang meledak karena JSON rusak lebih buruk
 *      daripada lipatan yang terlupa.
 */
const KUNCI = "sidebar-lipatan";

/** Potongan Storage yang dipakai - memudahkan uji tanpa peramban. */
export interface GudangMini {
  getItem(k: string): string | null;
  setItem(k: string, v: string): void;
}

/** nama kelompok -> false berarti DILIPAT. Tidak tercatat = terbuka. */
export type Lipatan = Record<string, boolean>;

export function terbukaKah(
  nama: string,
  lipatan: Lipatan,
  memuatAktif: boolean,
): boolean {
  if (memuatAktif) return true;
  return lipatan[nama] !== false;
}

export function bacaLipatan(gudang: GudangMini | null): Lipatan {
  if (!gudang) return {};
  try {
    const mentah = gudang.getItem(KUNCI);
    if (!mentah) return {};
    const isi: unknown = JSON.parse(mentah);
    if (typeof isi !== "object" || isi === null || Array.isArray(isi)) {
      return {};
    }
    const bersih: Lipatan = {};
    for (const [k, v] of Object.entries(isi)) {
      if (typeof v === "boolean") bersih[k] = v;
    }
    return bersih;
  } catch {
    return {};
  }
}

export function simpanLipatan(gudang: GudangMini | null, l: Lipatan): void {
  if (!gudang) return;
  try {
    gudang.setItem(KUNCI, JSON.stringify(l));
  } catch {
    // Kuota penuh / mode privat: lipatan tidak diingat, dan itu bukan galat.
  }
}
